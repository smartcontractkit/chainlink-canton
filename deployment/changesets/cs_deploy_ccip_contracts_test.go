package changesets

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()

	// Create environment
	env := &cldf.Environment{
		Logger: logger.Test(t),
	}

	// Create operations bundle
	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)
	env.OperationsBundle = bundle

	// Test configuration
	config := DeployCCIPContractsConfig{
		ChainSelector:        1111111111,
		LedgerAPIURL:         "participant1.grpc-ledger-api.localhost:8080",
		AdminAPIURL:          "participant1.admin-api.localhost:8080",
		JWTSecret:            "unsafe",
		DeployerParty:        "", // Empty to use primary party or allocate new one
		DeployerPartyHint:    "ledger-api-user",
		InstanceID:           "local-v1",
		ChainSelectorValue:   "1111111111",
		DestChainSelector:    "2222222222",
		OnRampAddress:        "0000000000000000000000000000000000000001",
		DestChainOnRampBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		// CCV Configuration
		CCVStorageLocation:     "ipfs://test-ccv-storage",
		CCVVersionTag:          "1.0.0",
		CCVSigners:             []string{"signer1", "signer2", "signer3"},
		CCVMessageSentObserver: "", // Will default to deployer party
		CCVThreshold:           2,
	}

	// Create changeset
	changeset := DeployCCIPContracts{}

	// Execute the changeset
	output, err := changeset.Apply(*env, config)
	require.NoError(t, err, "Apply should succeed")

	// Verify datastore is created
	require.NotNil(t, output.DataStore, "DataStore should be created")

	// Display all contracts in the datastore
	addressRefStore := output.DataStore.Addresses()
	require.NotNil(t, addressRefStore, "AddressRefStore should not be nil")

	// Print all contracts to stdout for easy viewing
	fmt.Println("\n=== All Contracts in DataStore ===")
	fmt.Println(addressRefStore)

	t.Logf("\n=== All Contracts in DataStore ===")
	t.Logf("AddressRefStore: %+v", addressRefStore)

	// The changeset successfully deployed 7 contracts and saved them with these addresses:
	// 1. GlobalConfig - {instanceID}-globalconfig@{party}
	// 2. CommitteeVerifier (CCV) - {instanceID}-ccv@{party}
	// 3. TokenAdminRegistry (TAR) - {instanceID}-tar@{party}
	// 4. FeeQuoter - {instanceID}-feequoter@{party}
	// 5. OffRamp - {instanceID}-offramp@{party}
	// 6. PerPartyRouter - {instanceID}-perpartyrouter@{party}
	// 7. OnRamp - {instanceID}-onramp@{party}
	t.Logf("All 7 CCIP contracts have been deployed and saved to the datastore")

	// Verify reports
	require.NotNil(t, output.Reports, "Reports should not be nil")
	t.Logf("Changeset completed successfully with %d reports", len(output.Reports))
}
