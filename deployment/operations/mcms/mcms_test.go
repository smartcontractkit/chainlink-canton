package mcms

import (
	"sync"
	"testing"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
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

	// TODO Expose via CLDF/CTF
	token, err := cantonChain.Participants[0].JWTProvider.Token(t.Context())
	require.NoError(t, err)
	bindingClient, err := client.NewDamlClient(token, cantonChain.Participants[0].Endpoints.GRPCLedgerAPIURL).
		WithAdminAddress(cantonChain.Participants[0].Endpoints.AdminAPIURL).
		Build(t.Context())
	require.NoError(t, err, "failed to create Daml binding client")
	t.Cleanup(bindingClient.Close)

	// Upload Dar
	mcmdDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err, "failed to get MCMS dar file")
	err = bindingClient.PackageMng.UploadDarFile(t.Context(), mcmdDar, "")
	require.NoError(t, err, "failed to upload MCMS dar file")

	// Get primary party
	user, err := bindingClient.UserMng.GetUser(t.Context(), "user-participant1")
	require.NoError(t, err, "failed to get user")

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	deps := dependencies.CantonDeps{
		Chain:         *cantonChain,
		BindingClient: bindingClient,
		Party:         user.PrimaryParty,
	}

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

	var mcmsInstanceID contracts.InstanceID
	t.Run("Deploy", func(t *testing.T) {
		result, err := cld_ops.ExecuteOperation(bundle, Deploy, deps, contract.DeployInput[mcms.MCMS]{
			ChainSelector: cantonChain.Selector,
			ActAs:         []string{user.PrimaryParty},
			Template: mcms.MCMS{
				Owner:        types.PARTY(user.PrimaryParty),
				Role:         mcms.RoleProposer,
				ChainId:      types.INT64(chainID),
				McmsId:       types.TEXT(mcmsID),
				Config:       config,
				SeenHashes:   nil,
				ExpiringRoot: mcms.ExpiringRoot{},
				RootMetadata: mcms.RootMetadata{},
			},
			OwnerParty: types.PARTY(user.PrimaryParty),
		})
		require.NoError(t, err, "failed to deploy MCMS")
		mcmsInstanceID = contracts.InstanceID(result.Output.Address)
		require.Truef(t, mcmsInstanceID.Valid(), "instance ID is not valid: %s", mcmsInstanceID.String())
		t.Logf("Deployed MCMS, InstanceID: %s", mcmsInstanceID.String())
	})

	if !mcmsInstanceID.Valid() {
		t.Fatalf("instance ID is not valid: %s", mcmsInstanceID.String())
	}

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

		result, err := cld_ops.ExecuteOperation(bundle, SetConfig, deps, contract.ChoiceInput[mcms.SetConfig]{
			ChainSelector:   cantonChain.Selector,
			InstanceAddress: mcmsInstanceID.InstanceAddress(),
			ActAs:           []string{user.PrimaryParty},
			Args: mcms.SetConfig{
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
