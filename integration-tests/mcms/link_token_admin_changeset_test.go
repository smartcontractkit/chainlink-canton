package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestLinkTokenAdmin_ChangesetsThroughMCMS deploys LinkTokenAdmin through the factory, approves a
// mint, then withdraws it. Each step is a CLD changeset whose proposal runs on the ccipOwner MCMS timelock.
func TestLinkTokenAdmin_ChangesetsThroughMCMS(t *testing.T) {
	t.Parallel()

	sharedEnv := GetSharedCCIPMCMSEnvironment(t)
	participant := sharedEnv.Participant
	submitter := sharedEnv.CcipOwner
	selector := chainsel.CANTON_LOCALNET.Selector
	uid := uuid.NewString()[:8]

	tokenAdminDar, err := contracts.GetDar(contracts.LinkTokenAdmin, contracts.DevVersion)
	require.NoError(t, err)
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{tokenAdminDar}, participant)
	require.NoError(t, err)

	// ccipOwner is not the participant's ActAs party, so every changeset must emit a proposal.
	ccipOwner := testhelpers.AllocateParty(t, participant, "lta-owner-"+uid)
	testhelpers.GrantCanActAs(t, participant, ccipOwner)
	testhelpers.GrantCanReadAs(t, participant, ccipOwner)
	mcmsParties := []string{ccipOwner}
	deployer := participant
	deployer.ReadAsPartyIDs = mcmsParties

	_, mcmsDS := deployMCMSWithRoleRefs(t, participant, ccipOwner, sharedEnv.Config, sharedEnv.Config)

	factoryRaw := contracts.NewRawInstanceAddress(contracts.InstanceID("factory-lta-"+uid), types.PARTY(ccipOwner))
	createCCIPFactory(t, participant, ccipOwner, factoryRaw.InstanceID())

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Merge(mcmsDS))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		Address:       factoryRaw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(factoryRaw.String()),
		ChainSelector: selector,
		Type:          datastore.ContractType(factoryops.ContractType),
		Version:       factoryops.Version,
		Qualifier:     dsutils.QualifierCCIP,
	}))

	runProposal := func(t *testing.T, apply func(cldf.Environment) (cldf.ChangesetOutput, error)) cldf.ChangesetOutput {
		t.Helper()

		out, err := apply(updateMCMSConfigEnv(t, deployer, ds.Seal()))
		require.NoError(t, err)
		require.Len(t, out.MCMSTimelockProposals, 1, "ccipOwner is MCMS-governed, so the changeset must emit one proposal")
		executeTimelockProposal(t, participant, submitter, mcmsParties, sharedEnv.SortedSigners[:2], out.MCMSTimelockProposals[0])

		return out
	}

	adminRaw := contracts.NewRawInstanceAddress(contracts.InstanceID("link-token-admin-"+uid), types.PARTY(ccipOwner))
	deployOut := runProposal(t, func(e cldf.Environment) (cldf.ChangesetOutput, error) {
		cfg := changesets.CantonCSDeps[changesets.DeployLinkTokenAdminFromFactoryConfig]{
			ChainSelector: selector,
			Config: changesets.DeployLinkTokenAdminFromFactoryConfig{
				CCIPOwnerParty: ccipOwner,
				Description:    "deploy LinkTokenAdmin",
				Params: sequences.DeployLinkTokenAdminFromFactoryParams{
					InstanceID:   types.TEXT(adminRaw.InstanceID()),
					InstrumentID: splice_api_token_holding_v1.InstrumentId{Id: "LINK"},
				},
			},
		}
		if err := (changesets.DeployLinkTokenAdminFromFactory{}).VerifyPreconditions(e, cfg); err != nil {
			return cldf.ChangesetOutput{}, err
		}

		return changesets.DeployLinkTokenAdminFromFactory{}.Apply(e, cfg)
	})
	require.NoError(t, ds.Merge(deployOut.DataStore.Seal()))

	admin := activeByInstance[tokenadmin.LinkTokenAdmin](t, participant, mcmsParties, tokenadmin.LinkTokenAdmin{}.GetTemplateID(), adminRaw)
	require.Equal(t, splice_api_token_holding_v1.InstrumentId{Admin: types.PARTY(ccipOwner), Id: "LINK"}, admin.InstrumentId)

	govern := func(t *testing.T, description string, actions ...sequences.LinkTokenAdminAction) {
		t.Helper()

		runProposal(t, func(e cldf.Environment) (cldf.ChangesetOutput, error) {
			cfg := changesets.CantonCSDeps[changesets.GovernLinkTokenAdminConfig]{
				ChainSelector: selector,
				Config: changesets.GovernLinkTokenAdminConfig{
					Actions:     actions,
					Description: description,
				},
			}
			if err := (changesets.GovernLinkTokenAdmin{}).VerifyPreconditions(e, cfg); err != nil {
				return cldf.ChangesetOutput{}, err
			}

			return changesets.GovernLinkTokenAdmin{}.Apply(e, cfg)
		})
	}

	authorizationRaw := contracts.NewRawInstanceAddress(contracts.InstanceID("lta-cs-mint-"+uid), types.PARTY(ccipOwner))
	mintTemplate := tokenadmin.LinkMintAuthorization{}.GetTemplateID()

	// SetObservers and SetPaused each replace the admin contract; ApproveMint succeeds only if the batch
	// carries the latest admin through in order and leaves it unpaused.
	govern(t, "set observers, pause, unpause, approve LINK mint",
		sequences.LinkTokenAdminAction{SetObservers: &tokenadmin.SetObservers{Observers: []types.PARTY{types.PARTY(submitter)}}},
		sequences.LinkTokenAdminAction{SetPaused: &tokenadmin.SetPaused{Paused: true}},
		sequences.LinkTokenAdminAction{SetPaused: &tokenadmin.SetPaused{Paused: false}},
		sequences.LinkTokenAdminAction{ApproveMint: &tokenadmin.ApproveMint{
			AuthorizationId: types.TEXT(authorizationRaw.InstanceID()),
			Recipient:       types.PARTY(ccipOwner),
			Amount:          "10.5",
			Reference:       "changeset-itest",
			ValidForSeconds: 3600,
		}},
	)
	admin = activeByInstance[tokenadmin.LinkTokenAdmin](t, participant, mcmsParties, tokenadmin.LinkTokenAdmin{}.GetTemplateID(), adminRaw)
	require.Equal(t, []types.PARTY{types.PARTY(submitter)}, admin.Observers)
	require.False(t, bool(admin.Paused))
	auth := activeByInstance[tokenadmin.LinkMintAuthorization](t, participant, mcmsParties, mintTemplate, authorizationRaw)
	assertDecimalEqual(t, auth.Amount, "10.5")
	require.Equal(t, types.PARTY(ccipOwner), auth.Recipient)
	require.Equal(t, types.TEXT(adminRaw.InstanceID()), auth.AdminInstanceId)

	// A second approval with the live id must be rejected before signing.
	err = changesets.GovernLinkTokenAdmin{}.VerifyPreconditions(updateMCMSConfigEnv(t, deployer, ds.Seal()), changesets.CantonCSDeps[changesets.GovernLinkTokenAdminConfig]{
		ChainSelector: selector,
		Config: changesets.GovernLinkTokenAdminConfig{
			Actions: []sequences.LinkTokenAdminAction{{ApproveMint: &tokenadmin.ApproveMint{
				AuthorizationId: types.TEXT(authorizationRaw.InstanceID()),
				Recipient:       types.PARTY(ccipOwner),
				Amount:          "1",
				ValidForSeconds: 3600,
			}}},
			Description: "approve LINK mint again",
		},
	})
	require.ErrorContains(t, err, "already live")

	govern(t, "withdraw LINK mint authorization", sequences.LinkTokenAdminAction{Withdraw: &sequences.LinkAuthorizationRef{
		Kind:            linkadmingov.KindMint,
		AuthorizationID: authorizationRaw.InstanceID(),
	}})
	_, err = opcontract.FindActiveContractByInstanceAddress(t.Context(), participant.LedgerServices.State, mcmsParties, mintTemplate, authorizationRaw.InstanceAddress())
	require.ErrorContains(t, err, "no active contract found", "withdrawn authorization must be archived")
}

func activeByInstance[T model.CreateCommander](t *testing.T, participant canton.Participant, parties []string, templateID string, raw contracts.RawInstanceAddress) *T {
	t.Helper()

	active, err := opcontract.FindActiveContractByInstanceAddress(t.Context(), participant.LedgerServices.State, parties, templateID, raw.InstanceAddress())
	require.NoError(t, err, "find %s", raw)
	decoded, err := bindings.UnmarshalCreatedEvent[T](active.GetCreatedEvent())
	require.NoError(t, err)

	return decoded
}
