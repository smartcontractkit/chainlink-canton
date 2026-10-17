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

// TestLinkTokenAdminRegistry_ApproveAndExecute drives the full governed
// lifecycle against the real DA Registry: MCMS-governed ApproveMint/ApproveBurn
// dispatched through the bypasser role, then operator ExecuteMint/ExecuteBurn
// one-step BurnMintFactory_BurnMint on the bootstrap AllocationFactory — both
// self-hosted (registrar as recipient) and third-party recipient (extraActors
// path) — plus the Execute* pause enforcement (stale admin cid, paused admin,
// unpause).
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

	// Third-party recipient: the realistic mainnet path, where the minted
	// holding's owner differs from the registrar and the factory demands the
	// recipient's authority through extraActors.
	recipient := testhelpers.AllocateParty(t, participant, "lta-recipient")
	testhelpers.GrantCanActAs(t, participant, recipient)

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

	// Governed mint to the registrar itself: ApproveMint via MCMS, then
	// ExecuteMint one-step against the real registry AllocationFactory.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "ApproveMint",
		tokenadmin.ApproveMintParams{
			Recipient:       types.PARTY(party),
			Amount:          types.NUMERIC("10.0"),
			Reference:       types.TEXT("itest-mint-1"),
			ValidForSeconds: types.INT64(3600),
		},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	mintAuthCid := findMintAuthorization(t, ctx, participant, party, "itest-mint-1")

	mintedHoldingCid := executeMintOnRegistry(t, ctx, client, []string{party}, mintAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, party, "10.0", mintCtx, instDisclosed)

	// Governed burn round-trip of the registrar-held holding.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "ApproveBurn",
		tokenadmin.ApproveBurnParams{
			Holder:          types.PARTY(party),
			MaxAmount:       types.NUMERIC("10.0"),
			Reference:       types.TEXT("itest-burn-1"),
			ValidForSeconds: types.INT64(3600),
		},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	burnAuthCid := findBurnAuthorization(t, ctx, participant, party, "itest-burn-1")

	executeBurnOnRegistry(t, ctx, client, []string{party}, burnAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, []string{mintedHoldingCid}, mintCtx, instDisclosed)

	// Governed mint to the third-party recipient: the extraActors path against
	// the real registry factory.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "ApproveMint",
		tokenadmin.ApproveMintParams{
			Recipient:       types.PARTY(recipient),
			Amount:          types.NUMERIC("5.0"),
			Reference:       types.TEXT("itest-mint-3"),
			ValidForSeconds: types.INT64(3600),
		},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	recipientAuthCid := findMintAuthorization(t, ctx, participant, party, "itest-mint-3")

	recipientHoldingCid := executeMintOnRegistry(t, ctx, client, []string{party, recipient}, recipientAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, recipient, "5.0", mintCtx, instDisclosed)

	// Burn the recipient-held holding: the burn-side extraActors path.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "ApproveBurn",
		tokenadmin.ApproveBurnParams{
			Holder:          types.PARTY(recipient),
			MaxAmount:       types.NUMERIC("5.0"),
			Reference:       types.TEXT("itest-burn-2"),
			ValidForSeconds: types.INT64(3600),
		},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	recipientBurnAuthCid := findBurnAuthorization(t, ctx, participant, party, "itest-burn-2")

	executeBurnOnRegistry(t, ctx, client, []string{party, recipient}, recipientBurnAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, []string{recipientHoldingCid}, mintCtx, instDisclosed)

	// Pause enforcement: a fresh approval that must not survive rotation or pause.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "ApproveMint",
		tokenadmin.ApproveMintParams{
			Recipient:       types.PARTY(party),
			Amount:          types.NUMERIC("30.0"),
			Reference:       types.TEXT("itest-mint-2"),
			ValidForSeconds: types.INT64(3600),
		},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	pausedAuthCid := findMintAuthorization(t, ctx, participant, party, "itest-mint-2")

	// SetObservers rotates the admin cid without pausing, isolating the fetch
	// failure from the paused assert.
	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "SetObservers",
		tokenadmin.SetObserversParams{Observers: nil},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	staleAdminCid := liveAdminCid
	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)
	require.NotEqual(t, staleAdminCid, liveAdminCid, "SetObservers should rotate the admin cid")

	_, err = client.SubmitExerciseMulti(ctx, []string{party}, tokenadmin.LinkMintAuthorization{}, pausedAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			TokenAdminCid:      types.CONTRACT_ID(staleAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed})
	require.ErrorContains(t, err, "CONTRACT_NOT_FOUND", "ExecuteMint with a stale admin cid must fail at fetch")
	require.NotContains(t, err.Error(), "token admin", "stale admin cid must fail before the admin asserts")

	mcmsCid = dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "SetPaused",
		tokenadmin.SetPausedParams{Paused: types.BOOL(true)},
		adminInstanceAddr, liveAdminCid)
	bypasserNonce++

	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)

	_, err = client.SubmitExerciseMulti(ctx, []string{party}, tokenadmin.LinkMintAuthorization{}, pausedAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(bootstrap.AllocationFactory),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed})
	require.ErrorContains(t, err, "token admin paused", "ExecuteMint with a paused admin must fail on the pause assert")

	dispatchOnAdminViaMCMS(t, participant, party, mcmsEncoder, mcmsCid, mcmsInstanceAddr, bypasserNonce, sortedSigners,
		adminContract, "SetPaused",
		tokenadmin.SetPausedParams{Paused: types.BOOL(false)},
		adminInstanceAddr, liveAdminCid)

	liveAdminCid = findLinkTokenAdminCid(t, ctx, participant, party, adminInstanceID)

	unpausedHoldingCid := executeMintOnRegistry(t, ctx, client, []string{party}, pausedAuthCid, liveAdminCid,
		bootstrap.AllocationFactory, party, "30.0", mintCtx, instDisclosed)
	require.NotEmpty(t, unpausedHoldingCid)
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

// dispatchOnAdminViaMCMS bypasser-executes a single LinkTokenAdmin choice and
// returns the refreshed MCMS cid. `nonce` is the bypasser root's operation
// count before this dispatch — SetRoot requires preOpCount to match it, so the
// caller must pass an incrementing counter. `liveAdminCid` must be the current
// admin cid: consuming choices (SetPaused, SetObservers) rotate it, so the
// caller must re-query afterwards.
func dispatchOnAdminViaMCMS(
	t *testing.T,
	participant canton.Participant,
	party string,
	mcmsEncoder mcmsApi.MCMSEncoder,
	mcmsCid, mcmsInstanceAddr string,
	nonce int,
	sortedSigners []*MCMSSigner,
	adminContract *tokenadmin.Contract,
	functionName string,
	params any,
	adminInstanceAddr, liveAdminCid string,
) string {
	t.Helper()

	encoded, err := encodeAdminChoice(adminContract, params)
	require.NoError(t, err)

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(adminInstanceAddr),
		FunctionName:          types.TEXT(functionName),
		OperationData:         types.TEXT(encoded.OperationData),
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
		adminInstanceAddr: liveAdminCid,
	}, proposal.Operations[0], opProof)
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

// executeMintOnRegistry submits ExecuteMint and returns the created registry
// Holding cid, asserting its owner and amount. This is the one-step
// BurnMintFactory_BurnMint path against the real registry AllocationFactory;
// `actAs` must cover ccipOwner plus the recipient when they differ.
func executeMintOnRegistry(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	actAs []string,
	mintAuthCid, liveAdminCid, allocationFactoryCid, expectedOwner, expectedAmount string,
	mintCtx spliceMeta.ChoiceContext,
	instDisclosed *apiv2.DisclosedContract,
) string {
	t.Helper()

	res, err := client.SubmitExerciseMulti(ctx, actAs, tokenadmin.LinkMintAuthorization{}, mintAuthCid, "ExecuteMint",
		tokenadmin.ExecuteMint{
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(allocationFactoryCid),
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed})
	require.NoError(t, err, "ExecuteMint one-step against the registry AllocationFactory")

	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "Holding" {
			holding, err := bindings.UnmarshalCreatedEvent[registryholding.Holding](created)
			require.NoError(t, err)
			require.Equal(t, expectedOwner, string(holding.Owner), "minted holding owner")
			assertDecimalEqual(t, holding.Amount, expectedAmount)

			return created.GetContractId()
		}
	}

	t.Fatal("no registry Holding created by ExecuteMint")

	return ""
}

// executeBurnOnRegistry submits ExecuteBurn and asserts every input holding is
// archived. `actAs` must cover ccipOwner plus the holder when they differ.
func executeBurnOnRegistry(
	t *testing.T,
	ctx context.Context,
	client rkledger.Client,
	actAs []string,
	burnAuthCid, liveAdminCid, allocationFactoryCid string,
	inputHoldingCids []string,
	mintCtx spliceMeta.ChoiceContext,
	instDisclosed *apiv2.DisclosedContract,
) {
	t.Helper()

	inputs := make([]types.CONTRACT_ID, 0, len(inputHoldingCids))
	for _, cid := range inputHoldingCids {
		inputs = append(inputs, types.CONTRACT_ID(cid))
	}

	res, err := client.SubmitExerciseMulti(ctx, actAs, tokenadmin.LinkBurnAuthorization{}, burnAuthCid, "ExecuteBurn",
		tokenadmin.ExecuteBurn{
			TokenAdminCid:      types.CONTRACT_ID(liveAdminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(allocationFactoryCid),
			InputHoldingCids:   inputs,
			RegistryContext:    mintCtx,
			Meta:               spliceMeta.Metadata{},
		},
		[]*apiv2.DisclosedContract{instDisclosed})
	require.NoError(t, err, "ExecuteBurn against the registry AllocationFactory")

	burned := make(map[string]bool)
	for _, event := range res.GetTransaction().GetEvents() {
		if archived := event.GetArchived(); archived != nil && archived.GetTemplateId().GetEntityName() == "Holding" {
			burned[archived.GetContractId()] = true
		}
	}
	for _, cid := range inputHoldingCids {
		require.True(t, burned[cid], "holding %s should be archived by ExecuteBurn", cid)
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
	party, reference string,
) string {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, contracts.IdentifierFromBinding(tokenadmin.LinkMintAuthorization{}))
	require.NoError(t, err)

	for _, ac := range active {
		auth, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkMintAuthorization](ac.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.CcipOwner) == party && string(auth.Reference) == reference {
			return ac.GetCreatedEvent().GetContractId()
		}
	}

	t.Fatalf("no active LinkMintAuthorization with reference %s", reference)

	return ""
}

func findBurnAuthorization(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party, reference string,
) string {
	t.Helper()

	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, participant, contracts.IdentifierFromBinding(tokenadmin.LinkBurnAuthorization{}))
	require.NoError(t, err)

	for _, ac := range active {
		auth, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkBurnAuthorization](ac.GetCreatedEvent())
		require.NoError(t, err)
		if string(auth.CcipOwner) == party && string(auth.Reference) == reference {
			return ac.GetCreatedEvent().GetContractId()
		}
	}

	t.Fatalf("no active LinkBurnAuthorization with reference %s", reference)

	return ""
}
