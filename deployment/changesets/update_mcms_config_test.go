package changesets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

const (
	mcmsUnitOwnerParty   = "participant::owner"
	mcmsUnitOtherParty   = "participant::other"
	mcmsUnitSignerPrefix = "00000000000000000000000000000000000000"
)

func mcmsUnitTestSigners(prefix string, count int64) []mcmsApi.SignerInfo {
	signers := make([]mcmsApi.SignerInfo, count)
	for i := range signers {
		signers[i] = mcmsApi.SignerInfo{
			SignerAddress: types.TEXT(fmt.Sprintf("%s%c%c", mcmsUnitSignerPrefix, prefix[0], 'a'+byte(i))),
			SignerIndex:   types.INT64(i),
			SignerGroup:   0,
		}
	}

	return signers
}

func mcmsUnitTestEnv(t *testing.T, participantParty string, seedMCMSRef bool) cldf.Environment {
	t.Helper()

	ds := datastore.NewMemoryDataStore()
	if seedMCMSRef {
		raw, err := contracts.RawInstanceAddressFromString("mcms-upd-unit@" + mcmsUnitOwnerParty)
		require.NoError(t, err)
		ref := newMCMSRoleAddressRef(
			chainsel.CANTON_LOCALNET.Selector, raw,
			datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig), "CLLCCIP",
		)
		require.NoError(t, ds.AddressRefStore.Add(ref))
	}

	return cldf.Environment{
		BlockChains: chain.NewBlockChainsFromSlice([]chain.BlockChain{&canton.Chain{
			ChainMetadata: canton.ChainMetadata{Selector: chainsel.CANTON_LOCALNET.Selector},
			Participants:  []canton.Participant{{PartyID: participantParty}},
		}}),
		DataStore: ds.Seal(),
	}
}

func TestUpdateMCMSConfig_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	validConfig := MCMSConfigParams{
		Signers:      mcmsUnitTestSigners("aa", 3),
		GroupQuorums: []types.INT64{2},
		GroupParents: []types.INT64{0},
	}

	tests := []struct {
		name          string
		env           cldf.Environment
		roleConfigs   []MCMSRoleConfigParams
		errorContains string
	}{
		{
			name: "happy path",
			env:  mcmsUnitTestEnv(t, mcmsUnitOwnerParty, true),
			roleConfigs: []MCMSRoleConfigParams{{
				Role:   mcmsApi.RoleProposer,
				Config: validConfig,
			}},
		},
		{
			name:          "empty role configs",
			env:           mcmsUnitTestEnv(t, mcmsUnitOwnerParty, true),
			roleConfigs:   nil,
			errorContains: "roleConfigs is required",
		},
		{
			name: "oversized group config",
			env:  mcmsUnitTestEnv(t, mcmsUnitOwnerParty, true),
			roleConfigs: []MCMSRoleConfigParams{{
				Role: mcmsApi.RoleProposer,
				Config: MCMSConfigParams{
					Signers:      mcmsUnitTestSigners("aa", 3),
					GroupQuorums: make([]types.INT64, mcmsGroupCount+1),
					GroupParents: []types.INT64{0},
				},
			}},
			errorContains: "build MCMS config for role Proposer",
		},
		{
			name: "missing MCMS ref",
			env:  mcmsUnitTestEnv(t, mcmsUnitOwnerParty, false),
			roleConfigs: []MCMSRoleConfigParams{{
				Role:   mcmsApi.RoleProposer,
				Config: validConfig,
			}},
			errorContains: "must be deployed first",
		},
		{
			name: "proposal-driven participant rejected",
			env:  mcmsUnitTestEnv(t, mcmsUnitOtherParty, true),
			roleConfigs: []MCMSRoleConfigParams{{
				Role:   mcmsApi.RoleProposer,
				Config: validConfig,
			}},
			errorContains: "cannot ActAs MCMS owner party",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := UpdateMCMSConfig{}.VerifyPreconditions(tt.env, CantonCSDeps[UpdateMCMSConfigConfig]{
				ChainSelector: chainsel.CANTON_LOCALNET.Selector,
				Participant:   0,
				Config: UpdateMCMSConfigConfig{
					Params: UpdateMCMSConfigParams{RoleConfigs: tt.roleConfigs},
				},
			})
			if tt.errorContains == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.errorContains)
		})
	}
}

// TestUpdateMCMSConfig_DirectExecution deploys an MCMS via the DeployAndConfigureMCMS
// changeset, then exercises UpdateMCMSConfig with a new signer set and verifies the
// proposer role state on-chain.
func TestUpdateMCMSConfig_DirectExecution(t *testing.T) {
	t.Parallel()

	cantonChain, _, env := setupCantonEnv(t)
	participant := cantonChain.Participants[0]
	party := participant.PartyID

	uploadDARs(t, participant, contracts.MCMSCore)

	initialSigners := mcmsUnitTestSigners("aa", 3)
	initialConfig := MCMSConfigParams{
		Signers:      initialSigners,
		GroupQuorums: []types.INT64{2},
		GroupParents: []types.INT64{0},
	}

	deployOut, err := DeployAndConfigureMCMS{}.Apply(*env, CantonCSDeps[DeployAndConfigureMCMSConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployAndConfigureMCMSConfig{
			Params: DeployAndConfigureMCMSParams{
				OwnerParty:    party,
				ChainID:       1,
				InitialConfig: initialConfig,
				RoleConfigs: []MCMSRoleConfigParams{
					{Role: mcmsApi.RoleProposer, Config: initialConfig},
					{Role: mcmsApi.RoleCanceller, Config: initialConfig},
					{Role: mcmsApi.RoleBypasser, Config: initialConfig},
				},
			},
		},
	})
	require.NoError(t, err, "deploy and configure MCMS")

	updatedEnv := *env
	updatedEnv.DataStore = deployOut.DataStore.Seal()

	rawInstanceAddress, err := dsutils.MCMSRawInstanceAddress(updatedEnv.DataStore, chainsel.CANTON_LOCALNET.Selector, "CLLCCIP")
	require.NoError(t, err, "resolve deployed MCMS raw instance address")

	newSigners := mcmsUnitTestSigners("bb", 3)
	deps := CantonCSDeps[UpdateMCMSConfigConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: UpdateMCMSConfigConfig{
			Params: UpdateMCMSConfigParams{
				RoleConfigs: []MCMSRoleConfigParams{{
					Role: mcmsApi.RoleProposer,
					Config: MCMSConfigParams{
						Signers:      newSigners,
						GroupQuorums: []types.INT64{3},
						GroupParents: []types.INT64{0},
						ClearRoot:    true,
					},
				}},
			},
		},
	}

	require.NoError(t, UpdateMCMSConfig{}.VerifyPreconditions(updatedEnv, deps))

	out, err := UpdateMCMSConfig{}.Apply(updatedEnv, deps)
	require.NoError(t, err, "update MCMS config")
	assert.Empty(t, out.MCMSTimelockProposals, "direct execution should not produce proposals")

	requireRoleSigners(t, participant, party, rawInstanceAddress, "proposer", signerAddresses(newSigners))
	requireRoleSigners(t, participant, party, rawInstanceAddress, "bypasser", signerAddresses(initialSigners))
}

func signerAddresses(signers []mcmsApi.SignerInfo) []string {
	addrs := make([]string, len(signers))
	for i, s := range signers {
		addrs[i] = string(s.SignerAddress)
	}

	return addrs
}

func requireRoleSigners(
	t *testing.T,
	participant canton.Participant,
	party string,
	rawInstanceAddress contracts.RawInstanceAddress,
	roleField string,
	want []string,
) {
	t.Helper()

	active, err := opcontract.FindActiveContractByInstanceAddress(
		t.Context(), participant.LedgerServices.State, []string{party},
		mcmsCore.MCMS{}.GetTemplateID(), rawInstanceAddress.InstanceAddress(),
	)
	require.NoError(t, err, "find updated MCMS contract")

	var got []string
	for _, field := range active.GetCreatedEvent().GetCreateArguments().GetFields() {
		if field.GetLabel() != roleField {
			continue
		}
		for _, roleField := range field.GetValue().GetRecord().GetFields() {
			if roleField.GetLabel() != "config" {
				continue
			}
			for _, configField := range roleField.GetValue().GetRecord().GetFields() {
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
	require.Equal(t, want, got, "%s signers after SetConfig", roleField)
}
