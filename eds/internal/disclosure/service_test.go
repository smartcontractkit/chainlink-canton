package disclosure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/ledger"
)

// MockContractQuerier implements ContractQuerier for testing
type MockContractQuerier struct {
	Contracts []*ledger.ActiveContract
	Error     error
}

func (m *MockContractQuerier) GetAllContractsForParty(ctx context.Context, party string) ([]*ledger.ActiveContract, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Contracts, nil
}

func createMockContract(moduleName, entityName, environmentID, contractID string, blob []byte) *ledger.ActiveContract {
	return &ledger.ActiveContract{
		CreatedEvent: &apiv2.CreatedEvent{
			ContractId: contractID,
			TemplateId: &apiv2.Identifier{
				PackageId:  "test-package-id",
				ModuleName: moduleName,
				EntityName: entityName,
			},
			CreateArguments: &apiv2.Record{
				Fields: []*apiv2.RecordField{
					{
						Label: "instanceId",
						Value: &apiv2.Value{
							Sum: &apiv2.Value_Text{Text: environmentID},
						},
					},
				},
			},
			CreatedEventBlob: blob,
		},
		SynchronizerID: "test-sync-id",
	}
}

func TestService_GetCCIPSendDisclosures(t *testing.T) {
	envConfig := &config.EnvironmentsConfig{
		Environments: map[string]config.EnvironmentConfig{
			"mainnet-v1": {
				Party:       "ccip-owner::123",
				Description: "Production",
				Contracts: config.ContractIdentifiers{
					Router:    "mainnet-v1-router",
					OnRamp:    "mainnet-v1-onramp",
					FeeQuoter: "mainnet-v1-feequoter",
				},
			},
		},
	}

	t.Run("returns disclosures when all contracts found", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				createMockContract("CCIP.Router", "Router", "mainnet-v1-router", "router-id-123", []byte("router-blob")),
				createMockContract("CCIP.OnRamp", "OnRamp", "mainnet-v1-onramp", "onramp-id-456", []byte("onramp-blob")),
				createMockContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1-feequoter", "feequoter-id-789", []byte("feequoter-blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		disclosures, err := svc.GetCCIPSendDisclosures(context.Background(), "mainnet-v1")
		require.NoError(t, err)

		assert.Equal(t, "mainnet-v1", disclosures.EnvironmentID)
		assert.NotNil(t, disclosures.Contracts.Router)
		assert.NotNil(t, disclosures.Contracts.OnRamp)
		assert.NotNil(t, disclosures.Contracts.FeeQuoter)

		assert.Equal(t, "router-id-123", disclosures.Contracts.Router.ContractID)
		assert.Equal(t, "CCIP.Router", disclosures.Contracts.Router.TemplateID.ModuleName)
		assert.Equal(t, "Router", disclosures.Contracts.Router.TemplateID.EntityName)
	})

	t.Run("returns error for unknown environment", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{}
		svc := NewService(mockQuerier, envConfig)

		_, err := svc.GetCCIPSendDisclosures(context.Background(), "unknown-env")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown environment")
	})

	t.Run("returns error when router not found", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				createMockContract("CCIP.OnRamp", "OnRamp", "mainnet-v1-onramp", "onramp-id", []byte("onramp-blob")),
				createMockContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1-feequoter", "feequoter-id", []byte("feequoter-blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		_, err := svc.GetCCIPSendDisclosures(context.Background(), "mainnet-v1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "router not found")
	})

	t.Run("filters by instanceId correctly", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				createMockContract("CCIP.Router", "Router", "testnet-router", "router-testnet", []byte("blob")),         // Wrong env
				createMockContract("CCIP.Router", "Router", "mainnet-v1-router", "router-mainnet", []byte("blob")),      // Correct env
				createMockContract("CCIP.OnRamp", "OnRamp", "mainnet-v1-onramp", "onramp-mainnet", []byte("blob")),
				createMockContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1-feequoter", "feequoter-mainnet", []byte("blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		disclosures, err := svc.GetCCIPSendDisclosures(context.Background(), "mainnet-v1")
		require.NoError(t, err)

		assert.Equal(t, "router-mainnet", disclosures.Contracts.Router.ContractID)
	})

	t.Run("ignores package ID for upgrade resilience", func(t *testing.T) {
		contract := createMockContract("CCIP.Router", "Router", "mainnet-v1-router", "router-id", []byte("blob"))
		contract.CreatedEvent.TemplateId.PackageId = "different-package-id-after-upgrade"

		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				contract,
				createMockContract("CCIP.OnRamp", "OnRamp", "mainnet-v1-onramp", "onramp-id", []byte("blob")),
				createMockContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1-feequoter", "feequoter-id", []byte("blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		disclosures, err := svc.GetCCIPSendDisclosures(context.Background(), "mainnet-v1")
		require.NoError(t, err)

		assert.Equal(t, "router-id", disclosures.Contracts.Router.ContractID)
		assert.Equal(t, "different-package-id-after-upgrade", disclosures.Contracts.Router.TemplateID.PackageID)
	})
}

func TestService_GetCCIPExecuteDisclosures(t *testing.T) {
	envConfig := &config.EnvironmentsConfig{
		Environments: map[string]config.EnvironmentConfig{
			"mainnet-v1": {
				Party:       "ccip-owner::123",
				Description: "Production",
				Contracts: config.ContractIdentifiers{
					OffRamp:            "mainnet-v1-offramp",
					CCV:                "mainnet-v1-ccv",
					TokenAdminRegistry: "mainnet-v1-tar",
				},
			},
		},
	}

	t.Run("returns disclosures when all contracts found", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				createMockContract("CCIP.OffRamp", "OffRamp", "mainnet-v1-offramp", "offramp-id", []byte("offramp-blob")),
				createMockContract("CCIP.CommitteeVerifier", "CommitteeVerifier", "mainnet-v1-ccv", "ccv-id", []byte("ccv-blob")),
				createMockContract("CCIP.TokenAdminRegistry", "TokenAdminRegistry", "mainnet-v1-tar", "tar-id", []byte("tar-blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		disclosures, err := svc.GetCCIPExecuteDisclosures(context.Background(), "mainnet-v1")
		require.NoError(t, err)

		assert.Equal(t, "mainnet-v1", disclosures.EnvironmentID)
		assert.NotNil(t, disclosures.Contracts.OffRamp)
		assert.NotNil(t, disclosures.Contracts.CCV)
		assert.NotNil(t, disclosures.Contracts.TokenAdminRegistry)

		assert.Equal(t, "offramp-id", disclosures.Contracts.OffRamp.ContractID)
		assert.Equal(t, "ccv-id", disclosures.Contracts.CCV.ContractID)
		assert.Equal(t, "tar-id", disclosures.Contracts.TokenAdminRegistry.ContractID)
	})

	t.Run("returns error when offRamp not found", func(t *testing.T) {
		mockQuerier := &MockContractQuerier{
			Contracts: []*ledger.ActiveContract{
				createMockContract("CCIP.CommitteeVerifier", "CommitteeVerifier", "mainnet-v1-ccv", "ccv-id", []byte("ccv-blob")),
				createMockContract("CCIP.TokenAdminRegistry", "TokenAdminRegistry", "mainnet-v1-tar", "tar-id", []byte("tar-blob")),
			},
		}

		svc := NewService(mockQuerier, envConfig)
		_, err := svc.GetCCIPExecuteDisclosures(context.Background(), "mainnet-v1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "offRamp not found")
	})
}

func TestExtractInstanceID(t *testing.T) {
	t.Run("extracts instanceId from record", func(t *testing.T) {
		record := &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{
					Label: "owner",
					Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: "alice"}},
				},
				{
					Label: "instanceId",
					Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "mainnet-v1"}},
				},
			},
		}

		envID := ExtractInstanceID(record)
		assert.Equal(t, "mainnet-v1", envID)
	})

	t.Run("returns empty string when instanceId not found", func(t *testing.T) {
		record := &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{
					Label: "owner",
					Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: "alice"}},
				},
			},
		}

		envID := ExtractInstanceID(record)
		assert.Equal(t, "", envID)
	})

	t.Run("returns empty string for nil record", func(t *testing.T) {
		envID := ExtractInstanceID(nil)
		assert.Equal(t, "", envID)
	})

	t.Run("returns empty string when instanceId is not text", func(t *testing.T) {
		record := &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{
					Label: "instanceId",
					Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 123}},
				},
			},
		}

		envID := ExtractInstanceID(record)
		assert.Equal(t, "", envID)
	})
}

func TestContractTypeToModule(t *testing.T) {
	testCases := []struct {
		contractType ContractType
		moduleName   string
		entityName   string
	}{
		{ContractTypeRouter, "CCIP.Router", "Router"},
		{ContractTypeOnRamp, "CCIP.OnRamp", "OnRamp"},
		{ContractTypeFeeQuoter, "CCIP.FeeQuoter", "FeeQuoter"},
		{ContractTypeOffRamp, "CCIP.OffRamp", "OffRamp"},
		{ContractTypeCCV, "CCIP.CommitteeVerifier", "CommitteeVerifier"},
		{ContractTypeTokenAdminRegistry, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"},
		{ContractTypeTokenPool, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.contractType), func(t *testing.T) {
			info, ok := ContractTypeToModule[tc.contractType]
			assert.True(t, ok)
			assert.Equal(t, tc.moduleName, info.ModuleName)
			assert.Equal(t, tc.entityName, info.EntityName)
		})
	}
}
