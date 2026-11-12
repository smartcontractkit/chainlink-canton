package tests

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	splice "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	spliceMeta "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	registryholding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/utility/registry_holding_v0"
	rkledger "github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

const (
	linkTokenAdminInstrumentID = "LINK"
	linkTokenAdminTestChainID  = int64(1)
)

// TestLinkTokenAdminRegistry_ApproveAndExecute drives the governed lifecycle
// against the real DA Registry: MCMS bypasser dispatch of Approve*, SetPaused,
// SetObservers and WithdrawAuthorization, then submitter-driven ExecuteMint and
// ExecuteBurn on the bootstrap AllocationFactory, for the registrar itself and
// for a third-party recipient. ccipOwner never acts at the root of Execute*.
//
//nolint:paralleltest // Exclusive CTF environment
func TestLinkTokenAdminRegistry_ApproveAndExecute(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
	participant := env.Chain.Participants[0]
	party := participant.PartyID

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	// The registry utility stack plus link-token-admin (which embeds the splice
	// token interfaces, chainlink-api and mcms-api) and mcms-core for the MCMS
	// contract itself.
	utilityDars, err := registry.LoadUtilityDARs()
	require.NoError(t, err, "load registry utility DARs")
	tokenAdminDar, err := contracts.GetDar(contracts.LinkTokenAdmin, contracts.DevVersion)
	require.NoError(t, err)
	mcmsDar, err := contracts.GetDar(contracts.MCMSCore, contracts.DevVersion)
	require.NoError(t, err)
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), append(utilityDars, tokenAdminDar, mcmsDar), participant)
	require.NoError(t, err)

	ctx := t.Context()
	client := rkledger.NewCTFClient(participant)

	bootstrap, err := registry.BootstrapServices(ctx, client, party, linkTokenAdminInstrumentID)
	require.NoError(t, err)

	instrumentId := splice.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT(linkTokenAdminInstrumentID),
	}

	adminInstanceID := "link-token-admin-itest-" + uuid.NewString()[:8]
	adminInstanceAddr := fmt.Sprintf("%s@%s", adminInstanceID, party)

	liveAdminCid := createLinkTokenAdmin(t, ctx, client, party, adminInstanceID, instrumentId)

	// The recipient must co-sign through extraActors when it is not the registrar.
	recipient := testhelpers.AllocateParty(t, participant, "lta-recipient")
	testhelpers.GrantCanActAs(t, participant, recipient)

	// A submitter distinct from the registrar proves Execute* needs no root-level
	// ccipOwner authority.
	submitter := testhelpers.AllocateParty(t, participant, "lta-submitter")
	testhelpers.GrantCanActAs(t, participant, submitter)

	signers := createSigners(t)
	sortedSigners := SortSignersByAddress(signers)
	bypasserNonce := 0
	cfg := New2of3Config(signers)

	baseMcmsID := "mcms-lta-" + uuid.NewString()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, party)
	mcmsCid := createMCMSMultiRole(t, participant, party, linkTokenAdminTestChainID, baseMcmsID, cfg, 0, nil)
	mcmsEncoder := NewMCMSEncoder()

	adminContract := tokenadmin.NewContract(fmt.Sprintf("#%s", tokenadmin.PackageName), "Link.TokenAdmin", "LinkTokenAdmin")

	instDisclosed, err := registry.DiscloseByID(ctx, client, party, bootstrap.InstrumentConfiguration)
	require.NoError(t, err)
	mintCtx := registry.MintChoiceContext(bootstrap.InstrumentConfiguration, false)

	// Governed mint to the registrar itself.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveMint", encodeAdminParams(t, adminContract, tokenadmin.ApproveMintParams{
			AuthorizationId: types.TEXT("lta-auth-mint-1"),
			Recipient:       types.PARTY(party),
			Amount:          types.NUMERIC("10.0"),
			Reference:       types.TEXT("itest-mint-1"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	mintAuthCid := findMintAuthorization(t, ctx, participant, party, "lta-auth-mint-1")

	// The named submitter must act: ccipOwner's authority alone is not enough.
	_, err = client.SubmitExerciseMulti(ctx, []string{party}, tokenadmin.LinkMintAuthorization{}, mintAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			Submitter:          types.PARTY(submitter),
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed})
	require.Error(t, err, "ExecuteMint must require the named submitter")

	mintedHoldingCid := executeMintOnRegistry(t, ctx, client, []string{submitter}, party, mintAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, party, "10.0", mintCtx, instDisclosed)

	// Governed burn round-trip of the registrar-held holding.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveBurn", encodeAdminParams(t, adminContract, tokenadmin.ApproveBurnParams{
			AuthorizationId: types.TEXT("lta-auth-burn-1"),
			Holder:          types.PARTY(party),
			MaxAmount:       types.NUMERIC("10.0"),
			Reference:       types.TEXT("itest-burn-1"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	burnAuthCid := findBurnAuthorization(t, ctx, participant, party, "lta-auth-burn-1")

	executeBurnOnRegistry(t, ctx, client, []string{submitter}, party, burnAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, []string{mintedHoldingCid}, mintCtx, instDisclosed)

	// Governed mint to the third-party recipient: the extraActors path against
	// the real registry factory.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveMint", encodeAdminParams(t, adminContract, tokenadmin.ApproveMintParams{
			AuthorizationId: types.TEXT("lta-auth-mint-3"),
			Recipient:       types.PARTY(recipient),
			Amount:          types.NUMERIC("5.0"),
			Reference:       types.TEXT("itest-mint-3"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	recipientAuthCid := findMintAuthorization(t, ctx, participant, party, "lta-auth-mint-3")

	recipientHoldingCid := executeMintOnRegistry(t, ctx, client, []string{submitter, recipient}, party, recipientAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, recipient, "5.0", mintCtx, instDisclosed)

	// Burn the recipient-held holding: the burn-side extraActors path.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveBurn", encodeAdminParams(t, adminContract, tokenadmin.ApproveBurnParams{
			AuthorizationId: types.TEXT("lta-auth-burn-2"),
			Holder:          types.PARTY(recipient),
			MaxAmount:       types.NUMERIC("5.0"),
			Reference:       types.TEXT("itest-burn-2"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	recipientBurnAuthCid := findBurnAuthorization(t, ctx, participant, party, "lta-auth-burn-2")

	executeBurnOnRegistry(t, ctx, client, []string{submitter, recipient}, party, recipientBurnAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, []string{recipientHoldingCid}, mintCtx, instDisclosed)

	// Pause enforcement: a fresh approval that must not survive rotation or pause.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveMint", encodeAdminParams(t, adminContract, tokenadmin.ApproveMintParams{
			AuthorizationId: types.TEXT("lta-auth-mint-2"),
			Recipient:       types.PARTY(party),
			Amount:          types.NUMERIC("30.0"),
			Reference:       types.TEXT("itest-mint-2"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	pausedAuthCid := findMintAuthorization(t, ctx, participant, party, "lta-auth-mint-2")

	// SetObservers rotates the admin cid without pausing, isolating the fetch
	// failure from the paused assert.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"SetObservers", encodeAdminParams(t, adminContract, tokenadmin.SetObserversParams{Observers: nil}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	staleAdminCid := liveAdminCid
	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)
	require.NotEqual(t, staleAdminCid, liveAdminCid, "SetObservers should rotate the admin cid")

	factoryDisclosed, err := registry.DiscloseByID(ctx, client, party, bootstrap.AllocationFactory)
	require.NoError(t, err)
	pausedAuthDisclosed, err := registry.DiscloseByID(ctx, client, party, pausedAuthCid)
	require.NoError(t, err)
	_, err = client.SubmitExerciseMulti(ctx, []string{submitter}, tokenadmin.LinkMintAuthorization{}, pausedAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			Submitter:          types.PARTY(submitter),
			TokenAdminCid:      types.CONTRACT_ID(staleAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed, pausedAuthDisclosed, factoryDisclosed})
	require.ErrorContains(t, err, "CONTRACT_NOT_FOUND", "ExecuteMint with a stale admin cid must fail at fetch")
	require.NotContains(t, err.Error(), "token admin", "stale admin cid must fail before the admin asserts")

	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"SetPaused", encodeAdminParams(t, adminContract, tokenadmin.SetPausedParams{Paused: types.BOOL(true)}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)

	_, err = client.SubmitExerciseMulti(ctx, []string{submitter}, tokenadmin.LinkMintAuthorization{}, pausedAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			Submitter:          types.PARTY(submitter),
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		executeDisclosures(t, ctx, client, party, pausedAuthCid, liveAdminCid, bootstrap.AllocationFactory, instDisclosed))
	require.ErrorContains(t, err, "token admin paused", "ExecuteMint with a paused admin must fail on the pause assert")

	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"SetPaused", encodeAdminParams(t, adminContract, tokenadmin.SetPausedParams{Paused: types.BOOL(false)}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)

	unpausedHoldingCid := executeMintOnRegistry(t, ctx, client, []string{submitter}, party, pausedAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, party, "30.0", mintCtx, instDisclosed)
	require.NotEmpty(t, unpausedHoldingCid)

	// Governance withdraws an authorization through the authorization's own
	// MCMS entrypoint.
	mcmsCid = dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"ApproveMint", encodeAdminParams(t, adminContract, tokenadmin.ApproveMintParams{
			AuthorizationId: types.TEXT("lta-auth-mint-cancel"),
			Recipient:       types.PARTY(party),
			Amount:          types.NUMERIC("1.0"),
			Reference:       types.TEXT("itest-mint-cancel"),
			ValidForSeconds: types.INT64(3600),
		}),
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	cancelAuthCid := findMintAuthorization(t, ctx, participant, party, "lta-auth-mint-cancel")
	cancelAuthInstanceAddr := fmt.Sprintf("%s@%s", "lta-auth-mint-cancel", party)

	dispatchViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		"WithdrawAuthorization", "", cancelAuthInstanceAddr, cancelAuthCid)

	_, active := lookupMintAuthorization(t, ctx, participant, party, "lta-auth-mint-cancel")
	require.False(t, active, "withdrawn authorization must be archived")

	adminDisclosed, err := registry.DiscloseByID(ctx, client, party, liveAdminCid)
	require.NoError(t, err)
	_, err = client.SubmitExerciseMulti(ctx, []string{submitter}, tokenadmin.LinkMintAuthorization{}, cancelAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			Submitter:          types.PARTY(submitter),
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed, adminDisclosed, factoryDisclosed})
	require.ErrorContains(t, err, "CONTRACT_NOT_FOUND", "ExecuteMint on a withdrawn authorization must fail")
}

func createLinkTokenAdmin(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	party, instanceID string,
	instrumentId splice.InstrumentId,
) string {
	t.Helper()

	res, err := client.SubmitCreate(ctx, party, tokenadmin.LinkTokenAdmin{
		CcipOwner:    types.PARTY(party),
		InstanceId:   types.TEXT(instanceID),
		InstrumentId: instrumentId,
		Paused:       false,
		Observers:    nil,
	})
	require.NoError(t, err)

	cid, ok := rkledger.CreatedContractID(res.GetTransaction(), "LinkTokenAdmin")
	require.True(t, ok, "LinkTokenAdmin not created")

	return cid
}

// dispatchViaMCMS bypasser-executes one call on an MCMS target (the admin or
// an authorization) and returns the refreshed MCMS cid. `nonce` is the bypasser
// root's operation count before this dispatch. `targetCid` must be the current
// cid: consuming choices rotate it, so the caller re-queries afterwards.
func dispatchViaMCMS(
	t *testing.T,
	participant canton.Participant,
	party string,
	mcmsEncoder mcmsApi.MCMSEncoder,
	mcmsCid, mcmsInstanceAddr string,
	nonce int,
	sortedSigners []*MCMSSigner,
	functionName, operationData string,
	targetInstanceAddr, targetCid string,
) string {
	t.Helper()

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(targetInstanceAddr),
		FunctionName:          types.TEXT(functionName),
		OperationData:         types.TEXT(operationData),
	}}

	bypasserChoice := MustEncodeBypasserExecuteBatch(t, mcmsEncoder, mcmsApi.BypasserExecuteBatchParams{Calls: calls})

	bypasserMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleBypasser)
	proposal := NewMCMSProposal(int(linkTokenAdminTestChainID), bypasserMultisigID, nonce, false).
		AddOperation(mcmsInstanceAddr, bypasserChoice.Choice, bypasserChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	mcmsCid = setRootWithRole(t, participant, party, mcmsCid, "Bypasser", proposal, validUntil, signatures)

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)

	return bypasserExecuteBatch(t, participant, party, mcmsCid, map[string]string{
		targetInstanceAddr: targetCid,
	}, proposal.Operations[0], opProof)
}

func encodeAdminParams(t *testing.T, adminContract *tokenadmin.Contract, params any) string {
	t.Helper()

	encoded, err := encodeAdminChoice(adminContract, params)
	require.NoError(t, err)

	return encoded.OperationData
}

func encodeAdminChoice(adminContract *tokenadmin.Contract, params any) (*bind.EncodedChoice, error) {
	encoder := adminContract.Encoder()

	switch p := params.(type) {
	case tokenadmin.ApproveMintParams:
		return encoder.ApproveMintParams(p)
	case tokenadmin.ApproveBurnParams:
		return encoder.ApproveBurnParams(p)
	case tokenadmin.SetPausedParams:
		return encoder.SetPausedParams(p)
	case tokenadmin.SetObserversParams:
		return encoder.SetObserversParams(p)
	}

	return nil, fmt.Errorf("unsupported LinkTokenAdmin choice params: %T", params)
}

// executeDisclosures returns what a submitter that is a stakeholder of neither
// the authorization, the admin nor the registry factory must disclose alongside
// the instrument configuration.
func executeDisclosures(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	ccipOwner, authCid, liveAdminCid, burnMintFactoryCid string,
	instDisclosed *apiv2.DisclosedContract,
) []*apiv2.DisclosedContract {
	t.Helper()

	disclosed := make([]*apiv2.DisclosedContract, 0, 4)
	disclosed = append(disclosed, instDisclosed)
	for _, cid := range []string{authCid, liveAdminCid, burnMintFactoryCid} {
		d, err := registry.DiscloseByID(ctx, client, ccipOwner, cid)
		require.NoError(t, err, "disclose %s", cid)
		disclosed = append(disclosed, d)
	}

	return disclosed
}

// activeRegistryHoldings lists the registry Holdings visible to the client's
// participant party (ccipOwner, a stakeholder of every LINK holding). The
// submitter's own transaction view omits holdings it is not a stakeholder of,
// so the tests assert through the ACS instead.
func activeRegistryHoldings(t *testing.T, ctx context.Context, client rkledger.Client) map[string]*registryholding.Holding {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.Participant(), contracts.IdentifierFromBinding(registryholding.Holding{}))
	require.NoError(t, err)

	holdings := make(map[string]*registryholding.Holding, len(active))
	for _, ac := range active {
		holding, err := bindings.UnmarshalCreatedEvent[registryholding.Holding](ac.GetCreatedEvent())
		require.NoError(t, err)
		holdings[ac.GetCreatedEvent().GetContractId()] = holding
	}

	return holdings
}

// executeMintOnRegistry submits ExecuteMint as actAs[0] and returns the created
// registry Holding cid. `actAs` also holds the recipient when it differs from
// ccipOwner.
func executeMintOnRegistry(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	actAs []string,
	ccipOwner, mintAuthCid, liveAdminCid, burnMintFactoryCid, expectedOwner, expectedAmount string,
	mintCtx spliceMeta.ChoiceContext,
	instDisclosed *apiv2.DisclosedContract,
) string {
	t.Helper()

	before := activeRegistryHoldings(t, ctx, client)

	_, err := client.SubmitExerciseMulti(ctx, actAs, tokenadmin.LinkMintAuthorization{}, mintAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			Submitter:          types.PARTY(actAs[0]),
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(burnMintFactoryCid),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		executeDisclosures(t, ctx, client, ccipOwner, mintAuthCid, liveAdminCid, burnMintFactoryCid, instDisclosed))
	require.NoError(t, err, "ExecuteMint one-step against the registry factory")

	var minted []string
	for cid, holding := range activeRegistryHoldings(t, ctx, client) {
		if _, seen := before[cid]; seen {
			continue
		}
		require.Equal(t, expectedOwner, string(holding.Owner), "minted holding owner")
		assertDecimalEqual(t, holding.Amount, expectedAmount)
		minted = append(minted, cid)
	}
	require.Len(t, minted, 1, "ExecuteMint must create exactly one registry Holding")

	return minted[0]
}

// executeBurnOnRegistry submits ExecuteBurn as actAs[0] and asserts every input
// holding is archived. `actAs` also holds the holder when it differs from
// ccipOwner.
func executeBurnOnRegistry(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	actAs []string,
	ccipOwner, burnAuthCid, liveAdminCid, burnMintFactoryCid string,
	inputHoldingCids []string,
	mintCtx spliceMeta.ChoiceContext,
	instDisclosed *apiv2.DisclosedContract,
) {
	t.Helper()

	inputs := make([]types.CONTRACT_ID, 0, len(inputHoldingCids))
	for _, cid := range inputHoldingCids {
		inputs = append(inputs, types.CONTRACT_ID(cid))
	}

	before := activeRegistryHoldings(t, ctx, client)
	for _, cid := range inputHoldingCids {
		require.Contains(t, before, cid, "input holding %s must be active before ExecuteBurn", cid)
	}

	// The submitter is not a stakeholder of the input holdings either.
	disclosed := executeDisclosures(t, ctx, client, ccipOwner, burnAuthCid, liveAdminCid, burnMintFactoryCid, instDisclosed)
	for _, cid := range inputHoldingCids {
		holdingDisclosed, err := registry.DiscloseByID(ctx, client, ccipOwner, cid)
		require.NoError(t, err, "disclose input holding")
		disclosed = append(disclosed, holdingDisclosed)
	}

	_, err := client.SubmitExerciseMulti(ctx, actAs, tokenadmin.LinkBurnAuthorization{}, burnAuthCid, "ExecuteBurn",
		tokenadmin.ExecuteBurn{
			Submitter:          types.PARTY(actAs[0]),
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(burnMintFactoryCid),
			InputHoldingCids:   inputs,
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		disclosed)
	require.NoError(t, err, "ExecuteBurn against the registry factory")

	after := activeRegistryHoldings(t, ctx, client)
	for _, cid := range inputHoldingCids {
		require.NotContains(t, after, cid, "holding %s should be archived by ExecuteBurn", cid)
	}
}

// assertDecimalEqual compares two decimal strings by value, independent of scale.
func assertDecimalEqual(t *testing.T, got types.NUMERIC, want string) {
	t.Helper()

	gotRat, ok := new(big.Rat).SetString(string(got))
	require.True(t, ok, "invalid decimal %q", string(got))
	wantRat, ok := new(big.Rat).SetString(want)
	require.True(t, ok, "invalid decimal %q", want)
	require.Equal(t, 0, gotRat.Cmp(wantRat), "decimal %s, want %s", string(got), want)
}

func findLinkTokenAdminCid(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party, instanceID string,
) string {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, contracts.IdentifierFromBinding(tokenadmin.LinkTokenAdmin{}))
	require.NoError(t, err)

	for _, ac := range active {
		admin, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkTokenAdmin](ac.GetCreatedEvent())
		require.NoError(t, err)
		if string(admin.InstanceId) == instanceID && string(admin.CcipOwner) == party {
			return ac.GetCreatedEvent().GetContractId()
		}
	}

	t.Fatalf("no active LinkTokenAdmin with instanceId %s", instanceID)

	return ""
}

func findMintAuthorization(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party, authorizationID string,
) string {
	t.Helper()

	cid, ok := lookupMintAuthorization(t, ctx, participant, party, authorizationID)
	require.True(t, ok, "no active LinkMintAuthorization with instanceId %s", authorizationID)

	return cid
}

func lookupMintAuthorization(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party, authorizationID string,
) (string, bool) {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, contracts.IdentifierFromBinding(tokenadmin.LinkMintAuthorization{}))
	require.NoError(t, err)

	for _, ac := range active {
		auth, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkMintAuthorization](ac.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.CcipOwner) == party && string(auth.InstanceId) == authorizationID {
			return ac.GetCreatedEvent().GetContractId(), true
		}
	}

	return "", false
}

func findBurnAuthorization(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party, authorizationID string,
) string {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, contracts.IdentifierFromBinding(tokenadmin.LinkBurnAuthorization{}))
	require.NoError(t, err)

	for _, ac := range active {
		auth, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkBurnAuthorization](ac.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.CcipOwner) == party && string(auth.InstanceId) == authorizationID {
			return ac.GetCreatedEvent().GetContractId()
		}
	}

	t.Fatalf("no active LinkBurnAuthorization with instanceId %s", authorizationID)

	return ""
}
