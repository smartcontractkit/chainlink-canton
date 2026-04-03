package mcms

import (
	"sync"
	"testing"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

//nolint:paralleltest // Cannot run in parallel due to shared state
func TestMCMSOps(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)
	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	primaryParty := participant.PartyID

	// Upload Dar
	mcmdDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get MCMS dar file")
	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: mcmdDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload MCMS dar file")

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)

	chainID := int64(1)
	mcmsID := "test-mcms-1"

	// Create a simple 2-of-3 multisig config
	// This is a minimal config for testing
	signers := []mcms.SignerInfo{
		{
			SignerAddress: types.TEXT("0x1111111111111111111111111111111111111111"),
			SignerIndex:   types.INT64(0),
			SignerGroup:   types.INT64(0),
		},
		{
			SignerAddress: types.TEXT("0x2222222222222222222222222222222222222222"),
			SignerIndex:   types.INT64(1),
			SignerGroup:   types.INT64(0),
		},
		{
			SignerAddress: types.TEXT("0x3333333333333333333333333333333333333333"),
			SignerIndex:   types.INT64(2),
			SignerGroup:   types.INT64(0),
		},
	}

	// Create group quorums (32 groups, 2-of-3 for group 0)
	groupQuorums := make([]types.INT64, 32)
	groupQuorums[0] = types.INT64(2) // 2-of-3 for group 0

	// Create group parents (32 groups, all 0 for flat structure)
	groupParents := make([]types.INT64, 32)

	config := mcms.MultisigConfig{
		Signers:      signers,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
	}

	// Create RoleState for each role (multi-role MCMS structure)
	roleState := mcms.RoleState{
		Config:       config,
		SeenHashes:   nil,
		ExpiringRoot: mcms.ExpiringRoot{},
		RootMetadata: mcms.RootMetadata{},
	}

	var mcmsInstanceAddress contracts.InstanceAddress
	t.Run("Deploy", func(t *testing.T) {
		result, err := cld_ops.ExecuteOperation(bundle, Deploy, *cantonChain, contract.DeployInput[mcms.MCMS]{
			Template: mcms.MCMS{
				Owner:              types.PARTY(primaryParty),
				InstanceId:         types.TEXT(mcmsID),
				ChainId:            types.INT64(chainID),
				Proposer:           roleState,
				Canceller:          roleState,
				Bypasser:           roleState,
				MinDelay:           0,
				BlockedFunctions:   nil,
				TimelockTimestamps: nil,
			},
			OwnerParty: types.PARTY(primaryParty),
		})
		require.NoError(t, err, "failed to deploy MCMS")
		mcmsInstanceAddress = contracts.HexToInstanceAddress(result.Output.Address)
		t.Logf("Deployed MCMS, InstanceAddress: %s", mcmsInstanceAddress.String())
	})

	t.Run("SetConfig", func(t *testing.T) {
		// Create updated config with a new signer
		newSigners := []mcms.SignerInfo{
			{
				SignerAddress: types.TEXT("0x1111111111111111111111111111111111111111"),
				SignerIndex:   types.INT64(0),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x2222222222222222222222222222222222222222"),
				SignerIndex:   types.INT64(1),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x3333333333333333333333333333333333333333"),
				SignerIndex:   types.INT64(2),
				SignerGroup:   types.INT64(0),
			},
			{
				SignerAddress: types.TEXT("0x4444444444444444444444444444444444444444"),
				SignerIndex:   types.INT64(3),
				SignerGroup:   types.INT64(0),
			},
		}

		// Update quorum to 3-of-4
		newGroupQuorums := make([]types.INT64, 32)
		newGroupQuorums[0] = types.INT64(3) // 3-of-4 for group 0

		result, err := cld_ops.ExecuteOperation(bundle, SetConfig, *cantonChain, contract.ChoiceInput[mcms.SetConfig]{
			InstanceAddress: mcmsInstanceAddress,
			Args: mcms.SetConfig{
				TargetRole:      mcms.RoleProposer, // Target the proposer role
				NewSigners:      newSigners,
				NewGroupQuorums: newGroupQuorums,
				NewGroupParents: groupParents,      // Keep same parent structure
				ClearRoot:       types.BOOL(false), // Don't clear root
			},
		})
		require.NoError(t, err, "failed to set MCMS config")
		t.Logf("Set MCMS config, UpdateID: %s", result.Output.ExecInfo.UpdateID)
	})
}
