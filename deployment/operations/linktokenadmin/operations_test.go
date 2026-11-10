package linktokenadmin

import (
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	spliceHolding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
)

const (
	ccipOwner = "ccip-owner::1220aa"
	submitter = "submitter::1220bb"
	recipient = "recipient::1220cc"
)

func staticRegistry() RegistryInput {
	return RegistryInput{FactoryCid: "factory-cid"}
}

func TestExecuteInputValidate(t *testing.T) {
	t.Parallel()

	valid := ExecuteInput{Kind: KindMint, Authorization: "auth-1@" + ccipOwner, Registry: staticRegistry()}

	tests := []struct {
		name    string
		mutate  func(*ExecuteInput)
		wantErr string
	}{
		{name: "valid mint", mutate: func(*ExecuteInput) {}},
		{name: "valid mint via registry URL", mutate: func(in *ExecuteInput) { in.Registry = RegistryInput{URL: "http://registry"} }},
		{
			name: "valid burn",
			mutate: func(in *ExecuteInput) {
				in.Kind = KindBurn
				in.InputHoldingCids = []string{"h1"}
			},
		},
		{name: "unknown kind", mutate: func(in *ExecuteInput) { in.Kind = "redeem" }, wantErr: `unknown authorization kind "redeem"`},
		{name: "transfer not executable", mutate: func(in *ExecuteInput) { in.Kind = KindTransfer }, wantErr: "execute supports mint and burn, not transfer"},
		{name: "allocate not executable", mutate: func(in *ExecuteInput) { in.Kind = KindAllocate }, wantErr: "execute supports mint and burn, not allocate"},
		{name: "authorization without owner", mutate: func(in *ExecuteInput) { in.Authorization = "auth-1" }, wantErr: "authorization:"},
		{name: "no registry source", mutate: func(in *ExecuteInput) { in.Registry = RegistryInput{} }, wantErr: "URL or FactoryCid is required"},
		{name: "both registry sources", mutate: func(in *ExecuteInput) { in.Registry.URL = "http://registry" }, wantErr: "not both"},
		{name: "mint with inputs", mutate: func(in *ExecuteInput) { in.InputHoldingCids = []string{"h1"} }, wantErr: "a mint takes no input holdings"},
		{name: "burn without inputs", mutate: func(in *ExecuteInput) { in.Kind = KindBurn }, wantErr: "a burn requires InputHoldingCids"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			in := valid
			tc.mutate(&in)
			err := in.Validate()
			if tc.wantErr == "" {
				require.NoError(t, err)

				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestCleanupInputValidate(t *testing.T) {
	t.Parallel()

	require.NoError(t, CleanupInput{Kind: KindBurn, Authorization: "auth-1@" + ccipOwner}.Validate())
	require.ErrorContains(t, CleanupInput{Kind: "x", Authorization: "auth-1@" + ccipOwner}.Validate(), "unknown authorization kind")
	require.ErrorContains(t, CleanupInput{Kind: KindBurn, Authorization: "@"}.Validate(), "authorization:")
}

func TestCheckCounterparty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		counterparty string
		wantErr      string
	}{
		{name: "ccipOwner", counterparty: ccipOwner},
		{name: "submitter", counterparty: submitter},
		{name: "third party", counterparty: recipient, wantErr: "neither ccipOwner " + ccipOwner + " nor the submitter " + submitter},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := checkCounterparty(KindMint, authorization{ccipOwner: ccipOwner, counterparty: tc.counterparty}, submitter)
			if tc.wantErr == "" {
				require.NoError(t, err)

				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestDecodeAuthorization(t *testing.T) {
	t.Parallel()

	party := func(p string) *apiv2.Value { return &apiv2.Value{Sum: &apiv2.Value_Party{Party: p}} }
	text := func(s string) *apiv2.Value { return &apiv2.Value{Sum: &apiv2.Value_Text{Text: s}} }
	ev := &apiv2.CreatedEvent{
		ContractId: "auth-cid",
		CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "ccipOwner", Value: party(ccipOwner)},
			{Label: "instanceId", Value: text("auth-1")},
			{Label: "adminInstanceId", Value: text("admin-1")},
			{Label: "instrumentId", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "admin", Value: party(ccipOwner)},
				{Label: "id", Value: text("LINK")},
			}}}}},
			{Label: "holder", Value: party(recipient)},
			{Label: "maxAmount", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "5.0"}}},
			{Label: "reference", Value: text("ref")},
			{Label: "expiresAt", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 1}}},
		}},
	}

	auth, err := kinds[KindBurn].decode(ev)
	require.NoError(t, err)
	require.Equal(t, authorization{
		contractID:      "auth-cid",
		ccipOwner:       ccipOwner,
		adminInstanceID: "admin-1",
		instrumentID:    spliceHolding.InstrumentId{Admin: ccipOwner, Id: "LINK"},
		counterparty:    recipient,
		amount:          "5.0",
	}, auth)
}

func TestDedupDisclosures(t *testing.T) {
	t.Parallel()

	a := &apiv2.DisclosedContract{ContractId: "a", SynchronizerId: "first"}
	got := dedupDisclosures(a, nil, &apiv2.DisclosedContract{ContractId: "b"}, &apiv2.DisclosedContract{ContractId: "a", SynchronizerId: "second"})
	require.Len(t, got, 2)
	require.Equal(t, "first", got[0].GetSynchronizerId())
	require.Equal(t, "b", got[1].GetContractId())
}

func TestCleanupCommands(t *testing.T) {
	t.Parallel()

	want := map[Kind]struct{ template, choice string }{
		KindMint:     {"LinkMintAuthorization", "CleanupExpiredMintAuthorization"},
		KindBurn:     {"LinkBurnAuthorization", "CleanupExpiredBurnAuthorization"},
		KindTransfer: {"LinkTransferAuthorization", "CleanupExpiredTransferAuthorization"},
		KindAllocate: {"LinkAllocationAuthorization", "CleanupExpiredAllocationAuthorization"},
	}
	require.Len(t, kinds, len(want), "every authorization kind has a cleanup")
	for kind, w := range want {
		cmd := kinds[kind].cleanup("auth-cid", types.PARTY(submitter))
		require.Equal(t, w.choice, cmd.Choice, "%s cleanup choice", kind)
		require.Equal(t, "auth-cid", cmd.ContractID, "%s cleanup targets the authorization", kind)
		require.Equal(t, "#link-token-admin:Link.TokenAdmin:"+w.template, cmd.TemplateID, "%s cleanup template", kind)
		require.Contains(t, cleanupOps, kind, "%s cleanup operation registered", kind)
	}
}
