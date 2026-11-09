package linktokenadmin

import (
	"math/big"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
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
		{name: "unknown kind", mutate: func(in *ExecuteInput) { in.Kind = "redeem" }, wantErr: `unknown authorization kind "redeem"`},
		{name: "authorization without owner", mutate: func(in *ExecuteInput) { in.Authorization = "auth-1" }, wantErr: "authorization:"},
		{name: "no registry source", mutate: func(in *ExecuteInput) { in.Registry = RegistryInput{} }, wantErr: "URL or FactoryCid is required"},
		{name: "both registry sources", mutate: func(in *ExecuteInput) { in.Registry.URL = "http://registry" }, wantErr: "not both"},
		{
			name: "registry URL for a transfer",
			mutate: func(in *ExecuteInput) {
				in.Kind = KindTransfer
				in.Registry = RegistryInput{URL: "http://registry"}
			},
			wantErr: "mint and burn only",
		},
		{name: "mint with inputs", mutate: func(in *ExecuteInput) { in.InputHoldingCids = []string{"h1"} }, wantErr: "a mint takes no input holdings"},
		{name: "allocate without settlement", mutate: func(in *ExecuteInput) { in.Kind = KindAllocate }, wantErr: "allocation settlement parameters are required"},
		{name: "settlement on a mint", mutate: func(in *ExecuteInput) { in.Allocation = &AllocationInput{} }, wantErr: "apply only to an allocate"},
		{name: "timeout on a mint", mutate: func(in *ExecuteInput) { h := int64(1); in.TransferTimeoutHours = &h }, wantErr: "applies only to a transfer"},
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

func TestActAs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		kind         Kind
		counterparty string
		submitter    string
		want         []string
	}{
		{name: "mint to ccipOwner", kind: KindMint, counterparty: ccipOwner, submitter: submitter, want: []string{submitter}},
		{name: "mint to third party co-signs", kind: KindMint, counterparty: recipient, submitter: submitter, want: []string{submitter, recipient}},
		{name: "burn from third party co-signs", kind: KindBurn, counterparty: recipient, submitter: submitter, want: []string{submitter, recipient}},
		{name: "recipient submitting alone", kind: KindMint, counterparty: recipient, submitter: recipient, want: []string{recipient}},
		{name: "transfer never co-signs", kind: KindTransfer, counterparty: recipient, submitter: submitter, want: []string{submitter}},
		{name: "allocate never co-signs", kind: KindAllocate, counterparty: recipient, submitter: submitter, want: []string{submitter}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			auth := authorization{ccipOwner: ccipOwner, counterparty: tc.counterparty}
			require.Equal(t, tc.want, actAs(kinds[tc.kind], auth, tc.submitter))
		})
	}
}

func TestHoldingOwner(t *testing.T) {
	t.Parallel()

	auth := authorization{ccipOwner: ccipOwner, counterparty: recipient}
	require.Equal(t, recipient, holdingOwner(kinds[KindBurn], auth), "a burn consumes the holder's holdings")
	require.Equal(t, ccipOwner, holdingOwner(kinds[KindTransfer], auth), "a transfer spends ccipOwner's holdings")
	require.Equal(t, ccipOwner, holdingOwner(kinds[KindAllocate], auth), "an allocation spends ccipOwner's holdings")
}

func holdings(amounts map[string]int64) []holding {
	out := make([]holding, 0, len(amounts))
	for cid, amount := range amounts {
		out = append(out, holding{contractID: cid, amount: big.NewRat(amount, 1)})
	}

	return out
}

func TestSelectHoldings(t *testing.T) {
	t.Parallel()

	candidates := holdings(map[string]int64{"h7": 7, "h5": 5, "h3": 3, "h1": 1})

	tests := []struct {
		name    string
		kind    Kind
		amount  types.NUMERIC
		want    []string
		wantErr string
	}{
		{name: "burn hits maxAmount exactly", kind: KindBurn, amount: "8.0", want: []string{"h7", "h1"}},
		{name: "burn takes everything that fits", kind: KindBurn, amount: "16.0", want: []string{"h7", "h5", "h3", "h1"}},
		{name: "burn refuses to fall short of maxAmount", kind: KindBurn, amount: "9.0", wantErr: "not the authorized burn amount"},
		{name: "burn with nothing small enough", kind: KindBurn, amount: "0.5", wantErr: "no unlocked holdings fit"},
		{name: "transfer covers the amount", kind: KindTransfer, amount: "10.0", want: []string{"h7", "h5"}},
		{name: "allocate exact cover", kind: KindAllocate, amount: "7.0", want: []string{"h7"}},
		{name: "transfer short of funds", kind: KindTransfer, amount: "17.0", wantErr: "below the authorization amount"},
		{name: "invalid amount", kind: KindTransfer, amount: "abc", wantErr: "invalid authorization amount"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := selectHoldings(tc.kind, candidates, tc.amount)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)

				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestSelectHoldingsBurnNeverUnderBurnsSilently(t *testing.T) {
	t.Parallel()

	_, err := selectHoldings(KindBurn, holdings(map[string]int64{"a": 60, "b": 50, "c": 50}), "100.0")
	require.ErrorContains(t, err, "not the authorized burn amount", "greedy pick of 60 must not burn 60 of an approved 100")
}

// ledgerArgs encodes a command's arguments as the ledger receives them and decodes them into T.
func ledgerArgs[T any](t *testing.T, cmd *model.ExerciseCommand) T {
	t.Helper()

	var out T
	require.NoError(t, ledger.RecordToStruct(ledger.MapToValue(cmd.Arguments).GetRecord(), &out))

	return out
}

func TestExecuteCommand(t *testing.T) {
	t.Parallel()

	args := executeArgs{submitter: submitter, adminCid: "admin-cid", factoryCid: "factory-cid", inputs: []string{"h1", "h2"}}
	inputs := []types.CONTRACT_ID{"h1", "h2"}

	t.Run("mint", func(t *testing.T) {
		t.Parallel()
		cmd, err := executeCommand(ExecuteInput{Kind: KindMint}, "auth-cid", executeArgs{submitter: submitter, adminCid: "admin-cid", factoryCid: "factory-cid"})
		require.NoError(t, err)
		require.Equal(t, "ExecuteMint", cmd.Choice)
		require.Equal(t, "auth-cid", cmd.ContractID)
		got := ledgerArgs[tokenadmin.ExecuteMint](t, cmd)
		require.Equal(t, types.PARTY(submitter), got.Submitter)
		require.Equal(t, types.CONTRACT_ID("admin-cid"), got.TokenAdminCid)
		require.Equal(t, types.CONTRACT_ID("factory-cid"), got.BurnMintFactoryCid)
	})

	t.Run("burn", func(t *testing.T) {
		t.Parallel()
		cmd, err := executeCommand(ExecuteInput{Kind: KindBurn}, "auth-cid", args)
		require.NoError(t, err)
		require.Equal(t, "ExecuteBurn", cmd.Choice)
		got := ledgerArgs[tokenadmin.ExecuteBurn](t, cmd)
		require.Equal(t, types.PARTY(submitter), got.Submitter)
		require.Equal(t, types.CONTRACT_ID("factory-cid"), got.BurnMintFactoryCid)
		require.Equal(t, inputs, got.InputHoldingCids)
	})

	t.Run("transfer", func(t *testing.T) {
		t.Parallel()
		timeout := int64(24)
		cmd, err := executeCommand(ExecuteInput{Kind: KindTransfer, TransferTimeoutHours: &timeout}, "auth-cid", args)
		require.NoError(t, err)
		require.Equal(t, "ExecuteTransfer", cmd.Choice)
		got := ledgerArgs[tokenadmin.ExecuteTransfer](t, cmd)
		require.Equal(t, types.CONTRACT_ID("factory-cid"), got.TransferFactoryCid)
		require.Equal(t, inputs, got.InputHoldingCids)
		require.NotNil(t, got.TransferTimeoutHours)
		require.Equal(t, types.INT64(24), *got.TransferTimeoutHours)
	})

	t.Run("transfer default timeout", func(t *testing.T) {
		t.Parallel()
		cmd, err := executeCommand(ExecuteInput{Kind: KindTransfer}, "auth-cid", args)
		require.NoError(t, err)
		require.Nil(t, ledgerArgs[tokenadmin.ExecuteTransfer](t, cmd).TransferTimeoutHours)
	})

	t.Run("allocate", func(t *testing.T) {
		t.Parallel()
		in := ExecuteInput{Kind: KindAllocate, Allocation: &AllocationInput{
			SettlementExecutor: recipient, SettlementRefID: "ref-1", AllocateBeforeHours: 1, SettleBeforeHours: 2,
		}}
		cmd, err := executeCommand(in, "auth-cid", args)
		require.NoError(t, err)
		require.Equal(t, "ExecuteAllocate", cmd.Choice)
		got := ledgerArgs[tokenadmin.ExecuteAllocate](t, cmd)
		require.Equal(t, types.CONTRACT_ID("factory-cid"), got.AllocationFactoryCid)
		require.Equal(t, types.PARTY(recipient), got.SettlementExecutor)
		require.Equal(t, types.TEXT("ref-1"), got.SettlementRefId)
		require.Equal(t, types.INT64(1), got.AllocateBeforeHours)
		require.Equal(t, types.INT64(2), got.SettleBeforeHours)
		require.Equal(t, inputs, got.InputHoldingCids)
	})
}

func TestCleanupCommands(t *testing.T) {
	t.Parallel()

	want := map[Kind]string{
		KindMint:     "CleanupExpiredMintAuthorization",
		KindBurn:     "CleanupExpiredBurnAuthorization",
		KindTransfer: "CleanupExpiredTransferAuthorization",
		KindAllocate: "CleanupExpiredAllocationAuthorization",
	}
	require.Len(t, kinds, len(want))
	for kind, choice := range want {
		cmd := kinds[kind].cleanup("auth-cid", types.PARTY(submitter))
		require.Equal(t, choice, cmd.Choice, kind)
		require.Equal(t, "auth-cid", cmd.ContractID, kind)
		require.Equal(t, types.PARTY(submitter), ledgerArgs[tokenadmin.CleanupExpiredMintAuthorization](t, cmd).Submitter, kind)
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

func TestResultContractIDs(t *testing.T) {
	t.Parallel()

	cid := func(c string) *apiv2.Value { return &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: c}} }
	list := func(vs ...*apiv2.Value) *apiv2.Value {
		return &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: vs}}}
	}
	// TransferInstructionResult{output = Pending instructionCid, senderChangeCids, meta}
	transferResult := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "output", Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
			Constructor: "TransferInstructionResult_Pending",
			Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "transferInstructionCid", Value: cid("instruction")},
			}}}},
		}}}},
		{Label: "senderChangeCids", Value: list(cid("change-1"), cid("change-2"))},
		{Label: "meta", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{}}}},
	}}}}

	require.Equal(t, []string{"instruction", "change-1", "change-2"}, resultContractIDs(transferResult))
	require.Equal(t, []string{"h1", "h2"}, resultContractIDs(list(cid("h1"), cid("h2"))))
	require.Empty(t, resultContractIDs(&apiv2.Value{Sum: &apiv2.Value_Unit{}}))
	require.Empty(t, resultContractIDs(nil))
}

func TestDedupDisclosures(t *testing.T) {
	t.Parallel()

	a := &apiv2.DisclosedContract{ContractId: "a", SynchronizerId: "first"}
	got := dedupDisclosures(a, nil, &apiv2.DisclosedContract{ContractId: "b"}, &apiv2.DisclosedContract{ContractId: "a", SynchronizerId: "second"})
	require.Len(t, got, 2)
	require.Equal(t, "first", got[0].GetSynchronizerId())
	require.Equal(t, "b", got[1].GetContractId())
}
