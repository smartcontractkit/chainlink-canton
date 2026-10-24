package adapters

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

const (
	testMCMSOwnerParty = "mcms-owner::1220aa"
	testOtherParty     = "deployer::1220bb"
	testMCMSRawAddress = "mcms-unit@" + testMCMSOwnerParty
	testSignerAddress  = "0x1375DC8A4c1476E6628b03216545e5cdcbff3f84"
)

func testMCMSRoleRef(t *testing.T, contractType cldf.ContractType) datastore.AddressRef {
	t.Helper()

	raw, err := contracts.RawInstanceAddressFromString(testMCMSRawAddress)
	require.NoError(t, err)

	return datastore.AddressRef{
		Address:       raw.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(raw.String()),
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Type:          datastore.ContractType(contractType),
		Version:       ccipdeploy.MCMSVersion,
		Qualifier:     ccipdeploymentutils.CLLQualifier,
	}
}

// The deployer is not the MCMS owner, so UpdateMCMSConfig builds a proposal without the ledger.
func testCantonBlockChains() chain.BlockChains {
	return chain.NewBlockChainsFromSlice([]chain.BlockChain{canton.Chain{
		Selector:     chainsel.CANTON_LOCALNET.Selector,
		Participants: []canton.Participant{{PartyID: testOtherParty}},
	}})
}

func testSingleSignerConfig() mcms_types.Config {
	return mcms_types.Config{Quorum: 1, Signers: []common.Address{common.HexToAddress(testSignerAddress)}}
}

func TestCantonMCMSDeployer_UpdateMCMSConfig_ProposalDriven(t *testing.T) {
	t.Parallel()

	bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
	report, err := cldf_ops.ExecuteSequence(bundle, (&CantonMCMSDeployer{}).UpdateMCMSConfig(), testCantonBlockChains(),
		ccipdeploy.UpdateMCMSConfigInputPerChainWithSelector{
			MCMConfig: testSingleSignerConfig(),
			MCMContracts: []datastore.AddressRef{
				testMCMSRoleRef(t, ccipdeploymentutils.ProposerManyChainMultisig),
				testMCMSRoleRef(t, ccipdeploymentutils.BypasserManyChainMultisig),
			},
			ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)

	batch := report.Output.BatchOps[0]
	require.Equal(t, mcms_types.ChainSelector(chainsel.CANTON_LOCALNET.Selector), batch.ChainSelector)
	require.Len(t, batch.Transactions, 2, "one self-dispatched SetConfig per role")

	raw, err := contracts.RawInstanceAddressFromString(testMCMSRawAddress)
	require.NoError(t, err)

	// numSigners | addrLen addr signerIndex signerGroup | numQuorums quorums | numParents parents | clearRoot
	setConfigParamsHex := "01" + "14" + "1375dc8a4c1476e6628b03216545e5cdcbff3f84" + "00000000" + "00000000" +
		"20" + "00000001" + strings.Repeat("00000000", 31) +
		"20" + strings.Repeat("00000000", 32) +
		"00"

	for i, wantRoleByte := range []string{"00", "02"} {
		tx := batch.Transactions[i]
		require.Equal(t, raw.InstanceAddress().Hex(), tx.To)
		require.Equal(t, "CantonMCMS", tx.ContractType)
		require.Equal(t, wantRoleByte+setConfigParamsHex, hex.EncodeToString(tx.Data))

		var af cantonsdk.AdditionalFields
		require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
		require.Equal(t, testMCMSRawAddress, af.TargetInstanceAddress, "self-dispatch targets the MCMS itself")
		require.Equal(t, "SetConfig", af.FunctionName)
		require.Empty(t, af.TargetTemplateID, "self-dispatch calls must not carry a target contract")
	}
}

func TestCantonMCMSDeployer_UpdateMCMSConfig_Errors(t *testing.T) {
	t.Parallel()

	missingLabel := testMCMSRoleRef(t, ccipdeploymentutils.ProposerManyChainMultisig)
	missingLabel.Labels = datastore.NewLabelSet()

	tests := []struct {
		name   string
		chains chain.BlockChains
		ref    datastore.AddressRef
		msg    string
	}{
		{
			name:   "chain not in environment",
			chains: chain.NewBlockChainsFromSlice(nil),
			ref:    testMCMSRoleRef(t, ccipdeploymentutils.ProposerManyChainMultisig),
			msg:    "not found or has no participants",
		},
		{
			name:   "unsupported contract type",
			chains: testCantonBlockChains(),
			ref:    testMCMSRoleRef(t, ccipdeploymentutils.RBACTimelock),
			msg:    `unsupported MCMS contract type "RBACTimelock"`,
		},
		{
			name:   "missing raw instance address label",
			chains: testCantonBlockChains(),
			ref:    missingLabel,
			msg:    "missing the raw instance address label",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
			_, err := cldf_ops.ExecuteSequence(bundle, (&CantonMCMSDeployer{}).UpdateMCMSConfig(), tt.chains,
				ccipdeploy.UpdateMCMSConfigInputPerChainWithSelector{
					MCMConfig:     testSingleSignerConfig(),
					MCMContracts:  []datastore.AddressRef{tt.ref},
					ChainSelector: chainsel.CANTON_LOCALNET.Selector,
				},
			)
			require.ErrorContains(t, err, tt.msg)
		})
	}
}

func TestCantonMCMSDeployer_Registered(t *testing.T) {
	t.Parallel()

	deployer, ok := ccipdeploy.GetRegistry().GetDeployer(chainsel.FamilyCanton, ccipdeploy.MCMSVersion)
	require.True(t, ok)
	require.IsType(t, &CantonMCMSDeployer{}, deployer)
}

func TestCantonMCMSDeployer_UpdateMCMSConfigChangesetVerify(t *testing.T) {
	t.Parallel()

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.AddressRefStore.Add(testMCMSRoleRef(t, ccipdeploymentutils.ProposerManyChainMultisig)))
	env := cldf.Environment{BlockChains: testCantonBlockChains(), DataStore: ds.Seal()}

	validInput := func() ccipdeploy.UpdateMCMSConfigInput {
		return ccipdeploy.UpdateMCMSConfigInput{
			AdapterVersion: ccipdeploy.MCMSVersion,
			MCMS: ccipmcms.Input{
				TimelockAction: mcms_types.TimelockActionSchedule,
				ValidUntil:     3759765795,
				Qualifier:      ccipdeploymentutils.CLLQualifier,
			},
			Chains: map[uint64]ccipdeploy.UpdateMCMSConfigInputPerChain{
				chainsel.CANTON_LOCALNET.Selector: {
					MCMConfig: testSingleSignerConfig(),
					MCMContracts: []datastore.AddressRef{{
						Type:      datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig),
						Version:   ccipdeploy.MCMSVersion,
						Qualifier: ccipdeploymentutils.CLLQualifier,
					}},
				},
			},
		}
	}

	changeset := ccipdeploy.UpdateMCMSConfig(ccipdeploy.GetRegistry(), nil)
	require.NoError(t, changeset.VerifyPreconditions(env, validInput()))

	timelockRef := validInput()
	timelockRef.Chains[chainsel.CANTON_LOCALNET.Selector] = ccipdeploy.UpdateMCMSConfigInputPerChain{
		MCMConfig: testSingleSignerConfig(),
		MCMContracts: []datastore.AddressRef{{
			Type:      datastore.ContractType(ccipdeploymentutils.RBACTimelock),
			Version:   ccipdeploy.MCMSVersion,
			Qualifier: ccipdeploymentutils.CLLQualifier,
		}},
	}
	require.ErrorContains(t, changeset.VerifyPreconditions(env, timelockRef), "type of contract needs to be mcms")
}

func TestCantonMCMSDeployer_TimelockHooks(t *testing.T) {
	t.Parallel()

	deployer := &CantonMCMSDeployer{}
	chains := testCantonBlockChains()
	selector := chainsel.CANTON_LOCALNET.Selector

	t.Run("FinalizeDeployMCMS is a no-op", func(t *testing.T) {
		t.Parallel()

		bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
		report, err := cldf_ops.ExecuteSequence(bundle, deployer.FinalizeDeployMCMS(), chains,
			ccipdeploy.MCMSDeploymentConfigPerChainWithAddress{ChainSelector: selector})
		require.NoError(t, err)
		require.Empty(t, report.Output.BatchOps)
	})

	t.Run("GrantAdminRoleToTimelock is not supported", func(t *testing.T) {
		t.Parallel()

		bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
		_, err := cldf_ops.ExecuteSequence(bundle, deployer.GrantAdminRoleToTimelock(), chains,
			ccipdeploy.GrantAdminRoleToTimelockConfigPerChainWithSelector{ChainSelector: selector})
		require.ErrorContains(t, err, "GrantAdminRoleToTimelock is not supported on Canton")
	})

	t.Run("DeployMCMS is not supported", func(t *testing.T) {
		t.Parallel()

		bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
		_, err := cldf_ops.ExecuteSequence(bundle, deployer.DeployMCMS(), chains,
			ccipdeploy.MCMSDeploymentConfigPerChainWithAddress{ChainSelector: selector})
		require.ErrorContains(t, err, "no integer Daml chain ID")
	})

	t.Run("1.6-only methods are not supported", func(t *testing.T) {
		t.Parallel()

		bundle := cldf_ops.NewBundle(t.Context, logger.Test(t), cldf_ops.NewMemoryReporter())
		_, err := cldf_ops.ExecuteSequence(bundle, deployer.SetOCR3Config(), chains, ccipdeploy.SetOCR3ConfigInput{})
		require.ErrorContains(t, err, "SetOCR3Config is 1.6-only")

		_, err = cldf_ops.ExecuteSequence(bundle, deployer.DeployChainContracts(), chains,
			ccipdeploy.ContractDeploymentConfigPerChainWithAddress{})
		require.ErrorContains(t, err, "DeployChainContracts is 1.6-only")
	})
}
