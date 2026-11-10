package sequences

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
)

const (
	testChainSelector = uint64(9_999_001)
	testCCIPOwner     = "ccipOwner::1220abcd"
	testAdminAddress  = contracts.RawInstanceAddress("link-token-admin-1@" + testCCIPOwner)
)

type hexDecoder interface{ UnmarshalHex(string) error }

func TestGovernLinkTokenAdmin_proposalRoundTrip(t *testing.T) {
	t.Parallel()

	mint := tokenadmin.ApproveMint{AuthorizationId: "auth-mint", Recipient: "alice::1220", Amount: "12.5", Reference: "ref-mint", ValidForSeconds: 3600}
	burn := tokenadmin.ApproveBurn{AuthorizationId: "auth-burn", Holder: "bob::1220", MaxAmount: "0.0000000001", Reference: "ref-burn", ValidForSeconds: 60}
	transfer := tokenadmin.ApproveTransfer{AuthorizationId: "auth-transfer", Receiver: "carol::1220", Amount: "1000000", Reference: "ref-transfer", ValidForSeconds: 7200}
	allocate := tokenadmin.ApproveAllocate{AuthorizationId: "auth-allocate", Receiver: "dave::1220", Amount: "3.25", TransferLegId: "leg-1", Reference: "ref-allocate", ValidForSeconds: 900}
	observers := tokenadmin.SetObservers{Observers: []types.PARTY{"obs1::1220", "obs2::1220"}}

	actions := []LinkTokenAdminAction{
		{SetPaused: &tokenadmin.SetPaused{Paused: false}},
		{ApproveMint: &mint},
		{ApproveBurn: &burn},
		{ApproveTransfer: &transfer},
		{ApproveAllocate: &allocate},
		{SetObservers: &observers},
		{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindMint, AuthorizationID: "old-mint"}},
		{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindBurn, AuthorizationID: "old-burn"}},
		{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindTransfer, AuthorizationID: "old-transfer"}},
		{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindAllocate, AuthorizationID: "old-allocate"}},
	}

	out, err := operations.ExecuteSequence(newTestBundle(t), GovernLinkTokenAdmin, canton.Chain{Selector: testChainSelector}, GovernLinkTokenAdminInput{
		AdminAddress:   testAdminAddress,
		ProposalDriven: true,
		Actions:        actions,
	})
	require.NoError(t, err)
	require.Len(t, out.Output.BatchOps, 1, "all actions go into one atomic batch")
	batch := out.Output.BatchOps[0]
	require.Equal(t, mcms_types.ChainSelector(testChainSelector), batch.ChainSelector)
	require.Len(t, batch.Transactions, len(actions), "one transaction per action, in order")

	adminTemplate := tokenadmin.LinkTokenAdmin{}.GetTemplateID()
	authorization := func(id string) contracts.RawInstanceAddress {
		return contracts.RawInstanceAddress(id + "@" + testCCIPOwner)
	}

	want := []struct {
		function string
		target   contracts.RawInstanceAddress
		template string
		decoded  hexDecoder
		params   any
	}{
		{"SetPaused", testAdminAddress, adminTemplate, &tokenadmin.SetPausedParams{}, &tokenadmin.SetPausedParams{Paused: false}},
		{"ApproveMint", testAdminAddress, adminTemplate, &tokenadmin.ApproveMintParams{}, &tokenadmin.ApproveMintParams{
			AuthorizationId: mint.AuthorizationId, Recipient: mint.Recipient, Amount: mint.Amount, Reference: mint.Reference, ValidForSeconds: mint.ValidForSeconds,
		}},
		{"ApproveBurn", testAdminAddress, adminTemplate, &tokenadmin.ApproveBurnParams{}, &tokenadmin.ApproveBurnParams{
			AuthorizationId: burn.AuthorizationId, Holder: burn.Holder, MaxAmount: burn.MaxAmount, Reference: burn.Reference, ValidForSeconds: burn.ValidForSeconds,
		}},
		{"ApproveTransfer", testAdminAddress, adminTemplate, &tokenadmin.ApproveTransferParams{}, &tokenadmin.ApproveTransferParams{
			AuthorizationId: transfer.AuthorizationId, Receiver: transfer.Receiver, Amount: transfer.Amount, Reference: transfer.Reference, ValidForSeconds: transfer.ValidForSeconds,
		}},
		{"ApproveAllocate", testAdminAddress, adminTemplate, &tokenadmin.ApproveAllocateParams{}, &tokenadmin.ApproveAllocateParams{
			AuthorizationId: allocate.AuthorizationId, Receiver: allocate.Receiver, Amount: allocate.Amount, TransferLegId: allocate.TransferLegId, Reference: allocate.Reference, ValidForSeconds: allocate.ValidForSeconds,
		}},
		{"SetObservers", testAdminAddress, adminTemplate, &tokenadmin.SetObserversParams{}, &tokenadmin.SetObserversParams{Observers: observers.Observers}},
		{"WithdrawAuthorization", authorization("old-mint"), tokenadmin.LinkMintAuthorization{}.GetTemplateID(), nil, nil},
		{"WithdrawAuthorization", authorization("old-burn"), tokenadmin.LinkBurnAuthorization{}.GetTemplateID(), nil, nil},
		{"WithdrawAuthorization", authorization("old-transfer"), tokenadmin.LinkTransferAuthorization{}.GetTemplateID(), nil, nil},
		{"WithdrawAuthorization", authorization("old-allocate"), tokenadmin.LinkAllocationAuthorization{}.GetTemplateID(), nil, nil},
	}

	for i, w := range want {
		tx := batch.Transactions[i]
		var af cantonsdk.AdditionalFields
		require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
		require.Equal(t, w.function, af.FunctionName, "tx %d function", i)
		require.Equal(t, w.target.String(), af.TargetInstanceAddress, "tx %d target", i)
		require.Equal(t, w.target.InstanceAddress().Hex(), tx.To, "tx %d to", i)
		require.Equal(t, w.template, af.TargetTemplateID, "tx %d template", i)

		if w.decoded == nil {
			require.Empty(t, tx.Data, "tx %d: WithdrawAuthorization takes no operation data", i)
			continue
		}
		require.NoError(t, w.decoded.UnmarshalHex(hex.EncodeToString(tx.Data)), "tx %d decode", i)
		require.Equal(t, w.params, w.decoded, "tx %d params round trip", i)
	}
}

func TestGovernLinkTokenAdmin_validation(t *testing.T) {
	t.Parallel()

	validMint := tokenadmin.ApproveMint{AuthorizationId: "a", Recipient: "r::1", Amount: "1.0", ValidForSeconds: 1}
	cases := []struct {
		name    string
		input   GovernLinkTokenAdminInput
		wantErr string
	}{
		{"no actions", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress}, "at least one action"},
		{"bad admin address", GovernLinkTokenAdminInput{AdminAddress: "no-at-sign", Actions: []LinkTokenAdminAction{{ApproveMint: &validMint}}}, "adminAddress"},
		{"empty action", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{}}}, "exactly one call must be set, got 0"},
		{"two calls in one action", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveMint: &validMint, SetPaused: &tokenadmin.SetPaused{Paused: true},
		}}}, "exactly one call must be set, got 2"},
		{"zero amount", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveBurn: &tokenadmin.ApproveBurn{AuthorizationId: "a", Holder: "h::1", MaxAmount: "0.0", ValidForSeconds: 1},
		}}}, "must be a positive decimal"},
		{"unparsable amount", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveTransfer: &tokenadmin.ApproveTransfer{AuthorizationId: "a", Receiver: "r::1", Amount: "ten", ValidForSeconds: 1},
		}}}, "must be a positive decimal"},
		{"more than 10 fractional digits", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveMint: &tokenadmin.ApproveMint{AuthorizationId: "a", Recipient: "r::1", Amount: "1.00000000009", ValidForSeconds: 1},
		}}}, "at most 10 fractional digits"},
		{"fraction syntax", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveMint: &tokenadmin.ApproveMint{AuthorizationId: "a", Recipient: "r::1", Amount: "1/3", ValidForSeconds: 1},
		}}}, "must be a positive decimal"},
		{"authorization id with @", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindMint, AuthorizationID: "a@b"},
		}}}, "must not contain '@'"},
		{"withdraw of an authorization approved in the same proposal", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{
			{ApproveMint: &validMint},
			{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindMint, AuthorizationID: string(validMint.AuthorizationId)}},
		}}, "approved earlier in this proposal"},
		{"non-positive validity", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveAllocate: &tokenadmin.ApproveAllocate{AuthorizationId: "a", Receiver: "r::1", Amount: "1", ValidForSeconds: 0},
		}}}, "validForSeconds must be positive"},
		{"missing authorization id", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveMint: &tokenadmin.ApproveMint{Recipient: "r::1", Amount: "1", ValidForSeconds: 1},
		}}}, "authorizationId is required"},
		{"missing recipient", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			ApproveMint: &tokenadmin.ApproveMint{AuthorizationId: "a", Amount: "1", ValidForSeconds: 1},
		}}}, "must be a full party id"},
		{"same authorization approved twice", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{
			{ApproveMint: &validMint}, {ApproveMint: &validMint},
		}}, "approved twice in this proposal"},
		{"same authorization withdrawn twice", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{
			{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindBurn, AuthorizationID: "b"}},
			{Withdraw: &LinkAuthorizationRef{Kind: linkadmingov.KindBurn, AuthorizationID: "b"}},
		}}, "withdrawn twice in this proposal"},
		{"unknown withdraw kind", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			Withdraw: &LinkAuthorizationRef{Kind: "recovery", AuthorizationID: "a"},
		}}}, `unknown authorization kind "recovery"`},
		{"observer without fingerprint", GovernLinkTokenAdminInput{AdminAddress: testAdminAddress, Actions: []LinkTokenAdminAction{{
			SetObservers: &tokenadmin.SetObservers{Observers: []types.PARTY{"alice"}},
		}}}, "setObservers"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := operations.ExecuteSequence(newTestBundle(t), GovernLinkTokenAdmin, canton.Chain{Selector: testChainSelector}, tc.input)
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestDeployLinkTokenAdminFromFactory_proposalRoundTrip(t *testing.T) {
	t.Parallel()

	factoryRaw := contracts.RawInstanceAddress("factory-1@" + testCCIPOwner)
	observers := []types.PARTY{"obs::1220"}
	out, err := operations.ExecuteSequence(newTestBundle(t), DeployLinkTokenAdminFromFactory, canton.Chain{Selector: testChainSelector}, DeployLinkTokenAdminFromFactoryParams{
		CCIPOwnerParty:    testCCIPOwner,
		FactoryAddressRef: datastore.AddressRef{Labels: datastore.NewLabelSet(factoryRaw.String())},
		ProposalDriven:    true,
		InstanceID:        "link-token-admin-1",
		InstrumentID:      splice_api_token_holding_v1.InstrumentId{Id: "LINK"},
		Observers:         observers,
	})
	require.NoError(t, err)

	require.Len(t, out.Output.Addresses, 1)
	ref := out.Output.Addresses[0]
	require.Equal(t, datastore.ContractType(linkadmingov.ContractType), ref.Type)
	require.Equal(t, linkadmingov.Version, ref.Version)
	require.Equal(t, testAdminAddress.InstanceAddress().String(), ref.Address)
	require.Equal(t, []string{testAdminAddress.String()}, ref.Labels.List())

	require.Len(t, out.Output.BatchOps, 1)
	require.Len(t, out.Output.BatchOps[0].Transactions, 1)
	tx := out.Output.BatchOps[0].Transactions[0]
	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	require.Equal(t, "DeployLinkTokenAdmin", af.FunctionName)
	require.Equal(t, factoryRaw.String(), af.TargetInstanceAddress)
	require.Equal(t, string(factoryops.ContractType), tx.ContractType)

	var decoded factorybindings.DeployLinkTokenAdminParams
	require.NoError(t, decoded.UnmarshalHex(hex.EncodeToString(tx.Data)))
	require.Equal(t, factorybindings.DeployLinkTokenAdminParams{
		InstanceId:   "link-token-admin-1",
		InstrumentId: splice_api_token_holding_v1.InstrumentId{Admin: testCCIPOwner, Id: "LINK"},
		Observers:    observers,
	}, decoded)
}

func TestDeployLinkTokenAdminFromFactory_rejectsBadInstrument(t *testing.T) {
	t.Parallel()

	base := DeployLinkTokenAdminFromFactoryParams{
		CCIPOwnerParty:    testCCIPOwner,
		FactoryAddressRef: datastore.AddressRef{Labels: datastore.NewLabelSet("factory-1@" + testCCIPOwner)},
		ProposalDriven:    true,
	}

	_, err := operations.ExecuteSequence(newTestBundle(t), DeployLinkTokenAdminFromFactory, canton.Chain{Selector: testChainSelector}, base)
	require.ErrorContains(t, err, "InstrumentID.Id is required")

	otherAdmin := base
	otherAdmin.InstrumentID = splice_api_token_holding_v1.InstrumentId{Admin: "someone::else", Id: "LINK"}
	_, err = operations.ExecuteSequence(newTestBundle(t), DeployLinkTokenAdminFromFactory, canton.Chain{Selector: testChainSelector}, otherAdmin)
	require.ErrorContains(t, err, "instrument admin must match ccipOwnerParty")

	badObserver := base
	badObserver.InstrumentID = splice_api_token_holding_v1.InstrumentId{Id: "LINK"}
	badObserver.Observers = []types.PARTY{"alice"}
	_, err = operations.ExecuteSequence(newTestBundle(t), DeployLinkTokenAdminFromFactory, canton.Chain{Selector: testChainSelector}, badObserver)
	require.ErrorContains(t, err, "observers")
}

func newTestBundle(t *testing.T) operations.Bundle {
	t.Helper()

	return operations.NewBundle(t.Context, logger.Test(t), operations.NewMemoryReporter())
}
