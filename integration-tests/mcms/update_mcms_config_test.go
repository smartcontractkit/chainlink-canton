package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcmscore "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestMCMS_UpdateMCMSConfigAdapter runs the shared ccip UpdateMCMSConfig changeset on Canton,
// directly and through a proposal.
func TestMCMS_UpdateMCMSConfigAdapter(t *testing.T) {
	t.Parallel()

	sharedEnv := GetSharedCCIPMCMSEnvironment(t)
	participant := sharedEnv.Participant
	party := sharedEnv.CcipOwner

	t.Run("direct", func(t *testing.T) {
		t.Parallel()

		cancellerSigners := SortSignersByAddress(createSigners(t))
		raw, ds := deployMCMSWithRoleRefs(t, participant, party, sharedEnv.Config, New2of3Config(cancellerSigners))
		requireRoleSignerAddresses(t, participant, []string{party}, raw, "canceller", signerAddressesOf(cancellerSigners))
		newSigners := SortSignersByAddress(createSigners(t))

		out, err := ccipdeploy.UpdateMCMSConfig(ccipdeploy.GetRegistry(), nil).Apply(
			updateMCMSConfigEnv(t, participant, ds),
			updateProposerConfigInput(newSigners),
		)
		require.NoError(t, err, "UpdateMCMSConfig changeset")
		require.Empty(t, out.MCMSTimelockProposals, "the owner can act directly, so no proposal is expected")

		requireRoleSignerAddresses(t, participant, []string{party}, raw, "proposer", signerAddressesOf(newSigners))
		requireRoleSignerAddresses(t, participant, []string{party}, raw, "bypasser", signerAddressesOf(sharedEnv.SortedSigners))
	})

	t.Run("proposal", func(t *testing.T) {
		t.Parallel()

		ownerParty := testhelpers.AllocateParty(t, participant, "mcms-owner-"+uuid.NewString()[:8])
		testhelpers.GrantCanActAs(t, participant, ownerParty)
		testhelpers.GrantCanReadAs(t, participant, ownerParty)
		mcmsParties := []string{ownerParty}

		raw, ds := deployMCMSWithRoleRefs(t, participant, ownerParty, sharedEnv.Config, sharedEnv.Config)
		newSigners := SortSignersByAddress(createSigners(t))

		// Like prod: the deployer only reads as the owner, so the change needs a proposal.
		deployer := participant
		deployer.ReadAsPartyIDs = mcmsParties

		out, err := ccipdeploy.UpdateMCMSConfig(ccipdeploy.GetRegistry(), nil).Apply(
			updateMCMSConfigEnv(t, deployer, ds),
			updateProposerConfigInput(newSigners),
		)
		require.NoError(t, err, "UpdateMCMSConfig changeset")
		require.Len(t, out.MCMSTimelockProposals, 1)
		requireRoleSignerAddresses(t, participant, mcmsParties, raw, "proposer", signerAddressesOf(sharedEnv.SortedSigners))

		timelockProposal := out.MCMSTimelockProposals[0]
		require.Len(t, timelockProposal.Operations, 1)
		require.Len(t, timelockProposal.Operations[0].Transactions, 1)

		var af cantonsdk.AdditionalFields
		require.NoError(t, json.Unmarshal(timelockProposal.Operations[0].Transactions[0].AdditionalFields, &af))
		require.Equal(t, raw.String(), af.TargetInstanceAddress, "SetConfig is self-dispatched to the MCMS")
		require.Equal(t, mcmsops.SelfDispatchSetConfigFunction, af.FunctionName)

		executeTimelockProposal(t, participant, party, mcmsParties, sharedEnv.SortedSigners[:2], timelockProposal)

		requireRoleSignerAddresses(t, participant, mcmsParties, raw, "proposer", signerAddressesOf(newSigners))
		requireRoleSignerAddresses(t, participant, mcmsParties, raw, "bypasser", signerAddressesOf(sharedEnv.SortedSigners))
	})
}

// deployMCMSWithRoleRefs deploys an MCMS with the DeployAndConfigureMCMS changeset.
// The canceller gets cancellerConfig through a SetConfig after the deploy.
func deployMCMSWithRoleRefs(
	t *testing.T,
	participant canton.Participant,
	owner string,
	config MCMSConfig,
	cancellerConfig MCMSConfig,
) (contracts.RawInstanceAddress, datastore.DataStore) {
	t.Helper()

	deployer := participant
	deployer.PartyID = owner
	out, err := changesets.DeployAndConfigureMCMS{}.Apply(
		updateMCMSConfigEnv(t, deployer, datastore.NewMemoryDataStore().Seal()),
		changesets.CantonCSDeps[changesets.DeployAndConfigureMCMSConfig]{
			ChainSelector: chainsel.CANTON_LOCALNET.Selector,
			Config: changesets.DeployAndConfigureMCMSConfig{Params: sequences.DeployAndConfigureMCMSParams{
				OwnerParty:    owner,
				InstanceID:    "mcms-upd-" + uuid.NewString()[:8],
				ChainID:       1,
				InitialConfig: mcmsConfigParams(config),
				RoleConfigs: []sequences.MCMSRoleConfigParams{
					{Role: mcmsApi.RoleCanceller, Config: mcmsConfigParams(cancellerConfig)},
				},
			}},
		},
	)
	require.NoError(t, err, "DeployAndConfigureMCMS changeset")

	ds := out.DataStore.Seal()
	raw, err := dsutils.MCMSRawInstanceAddress(ds, chainsel.CANTON_LOCALNET.Selector, ccipdeploymentutils.CLLQualifier)
	require.NoError(t, err)

	return raw, ds
}

func mcmsConfigParams(config MCMSConfig) sequences.MCMSConfigParams {
	params := sequences.MCMSConfigParams{
		Signers:      make([]mcmsApi.SignerInfo, len(config.Signers)),
		GroupQuorums: make([]types.INT64, len(config.GroupQuorums)),
		GroupParents: make([]types.INT64, len(config.GroupParents)),
	}
	for i, s := range config.Signers {
		params.Signers[i] = mcmsApi.SignerInfo{
			SignerAddress: types.TEXT(s.SignerAddress),
			SignerIndex:   types.INT64(s.SignerIndex),
			SignerGroup:   types.INT64(s.SignerGroup),
		}
	}
	for i, q := range config.GroupQuorums {
		params.GroupQuorums[i] = types.INT64(q)
	}
	for i, p := range config.GroupParents {
		params.GroupParents[i] = types.INT64(p)
	}

	return params
}

func updateMCMSConfigEnv(t *testing.T, participant canton.Participant, ds datastore.DataStore) cldf.Environment {
	t.Helper()

	return cldf.Environment{
		Logger:     logger.Test(t),
		GetContext: t.Context,
		DataStore:  ds,
		BlockChains: chain.NewBlockChainsFromSlice([]chain.BlockChain{&canton.Chain{
			Selector:     chainsel.CANTON_LOCALNET.Selector,
			Participants: []canton.Participant{participant},
		}}),
		OperationsBundle: cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter()),
	}
}

func updateProposerConfigInput(signers []*MCMSSigner) ccipdeploy.UpdateMCMSConfigInput {
	addresses := make([]common.Address, len(signers))
	for i, s := range signers {
		addresses[i] = common.HexToAddress(s.Address)
	}

	return ccipdeploy.UpdateMCMSConfigInput{
		AdapterVersion: ccipdeploy.MCMSVersion,
		MCMS: ccipmcms.Input{
			TimelockAction: mcms_types.TimelockActionSchedule,
			ValidUntil:     uint32(time.Now().Add(time.Hour).Unix()),
			Qualifier:      ccipdeploymentutils.CLLQualifier,
			Description:    "Update Canton MCMS proposer config",
		},
		Chains: map[uint64]ccipdeploy.UpdateMCMSConfigInputPerChain{
			chainsel.CANTON_LOCALNET.Selector: {
				MCMConfig: mcms_types.Config{Quorum: 2, Signers: addresses},
				MCMContracts: []datastore.AddressRef{{
					Type:      datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig),
					Version:   mcmsops.Version,
					Qualifier: ccipdeploymentutils.CLLQualifier,
				}},
			},
		},
	}
}

// executeTimelockProposal signs, schedules and executes a Schedule proposal.
func executeTimelockProposal(
	t *testing.T,
	participant canton.Participant,
	submitter string,
	mcmsParties []string,
	signers []*MCMSSigner,
	timelockProposal mcmscore.TimelockProposal,
) {
	t.Helper()

	selector := mcms_types.ChainSelector(chainsel.CANTON_LOCALNET.Selector)

	proposal, _, err := timelockProposal.Convert(t.Context(), map[mcms_types.ChainSelector]sdk.TimelockConverter{
		selector: cantonsdk.NewTimelockConverter(),
	})
	require.NoError(t, err, "convert timelock proposal")

	inspector := cantonsdk.NewInspector(participant.LedgerServices.State, mcmsParties, cantonsdk.TimelockRoleProposer)
	signable, err := mcmscore.NewSignable(&proposal, map[mcms_types.ChainSelector]sdk.Inspector{selector: inspector})
	require.NoError(t, err, "create signable")
	for _, signer := range signers {
		_, err = signable.SignAndAppend(mcmscore.NewPrivateKeySigner(signer.PrivateKey))
		require.NoError(t, err, "sign proposal")
	}
	quorumMet, err := signable.ValidateSignatures(t.Context())
	require.NoError(t, err, "validate signatures")
	require.True(t, quorumMet, "quorum not met")

	encoders, err := proposal.GetEncoders()
	require.NoError(t, err, "get encoders")
	encoder, ok := encoders[selector].(*cantonsdk.Encoder)
	require.True(t, ok, "canton encoder")
	executor, err := cantonsdk.NewExecutor(
		encoder, inspector, participant.LedgerServices.Command, submitter, mcmsParties, cantonsdk.TimelockRoleProposer,
	)
	require.NoError(t, err, "create executor")
	executable, err := mcmscore.NewExecutable(&proposal, map[mcms_types.ChainSelector]sdk.Executor{selector: executor})
	require.NoError(t, err, "create executable")

	_, err = executable.SetRoot(t.Context(), selector)
	require.NoError(t, err, "SetRoot")
	for i := range proposal.Operations {
		_, err = executable.Execute(t.Context(), i)
		require.NoError(t, err, "ScheduleBatch operation %d", i)
	}

	timelockExecutable, err := mcmscore.NewTimelockExecutable(t.Context(), &timelockProposal, map[mcms_types.ChainSelector]sdk.TimelockExecutor{
		selector: cantonsdk.NewTimelockExecutor(participant.LedgerServices.Command, participant.LedgerServices.State, submitter, mcmsParties),
	})
	require.NoError(t, err, "create timelock executable")
	for i := range timelockProposal.Operations {
		_, err = timelockExecutable.Execute(t.Context(), i)
		require.NoError(t, err, "ExecuteScheduledBatch operation %d", i)
	}
}

func signerAddressesOf(signers []*MCMSSigner) []string {
	addresses := make([]string, len(signers))
	for i, s := range signers {
		addresses[i] = s.Address
	}

	return addresses
}

// requireRoleSignerAddresses checks one role's signers on the active MCMS contract.
func requireRoleSignerAddresses(
	t *testing.T,
	participant canton.Participant,
	parties []string,
	raw contracts.RawInstanceAddress,
	roleField string,
	want []string,
) {
	t.Helper()

	active, err := opcontract.FindActiveContractByInstanceAddress(
		t.Context(), participant.LedgerServices.State, parties, mcmsCore.MCMS{}.GetTemplateID(), raw.InstanceAddress(),
	)
	require.NoError(t, err, "find MCMS contract")

	var got []string
	for _, field := range active.GetCreatedEvent().GetCreateArguments().GetFields() {
		if field.GetLabel() != roleField {
			continue
		}
		for _, roleStateField := range field.GetValue().GetRecord().GetFields() {
			if roleStateField.GetLabel() != "config" {
				continue
			}
			for _, configField := range roleStateField.GetValue().GetRecord().GetFields() {
				if configField.GetLabel() != "signers" {
					continue
				}
				for _, signer := range configField.GetValue().GetList().GetElements() {
					for _, signerField := range signer.GetRecord().GetFields() {
						if signerField.GetLabel() == "signerAddress" {
							got = append(got, signerField.GetValue().GetText())
						}
					}
				}
			}
		}
	}
	require.Equal(t, want, got, "%s signers", roleField)
}
