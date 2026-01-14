package disclosure

import (
	"context"
	"encoding/base64"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/ledger"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/types"
)

type ContractType string

const (
	ContractTypeRouter             ContractType = "Router"
	ContractTypeOnRamp             ContractType = "OnRamp"
	ContractTypeFeeQuoter          ContractType = "FeeQuoter"
	ContractTypeOffRamp            ContractType = "OffRamp"
	ContractTypeCCV                ContractType = "CommitteeVerifier"
	ContractTypeTokenAdminRegistry ContractType = "TokenAdminRegistry"
	ContractTypeTokenPool          ContractType = "LockReleaseTokenPool"
)

type ModuleInfo struct {
	ModuleName string
	EntityName string
}

// maps contract types to module:entity names
var ContractTypeToModule = map[ContractType]ModuleInfo{
	ContractTypeRouter:             {"CCIP.Router", "Router"},
	ContractTypeOnRamp:             {"CCIP.OnRamp", "OnRamp"},
	ContractTypeFeeQuoter:          {"CCIP.FeeQuoter", "FeeQuoter"},
	ContractTypeOffRamp:            {"CCIP.OffRamp", "OffRamp"},
	ContractTypeCCV:                {"CCIP.CommitteeVerifier", "CommitteeVerifier"},
	ContractTypeTokenAdminRegistry: {"CCIP.TokenAdminRegistry", "TokenAdminRegistry"},
	ContractTypeTokenPool:          {"CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"},
}

type ContractQuerier interface {
	GetAllContractsForParty(ctx context.Context, party string) ([]*ledger.ActiveContract, error)
}

type Service struct {
	querier   ContractQuerier
	envConfig *config.EnvironmentsConfig
}

func NewService(querier ContractQuerier, envConfig *config.EnvironmentsConfig) *Service {
	return &Service{
		querier:   querier,
		envConfig: envConfig,
	}
}

func (s *Service) GetCCIPSendDisclosures(ctx context.Context, environmentID string) (*types.CCIPSendDisclosures, error) {
	env, ok := s.envConfig.GetEnvironment(environmentID)
	if !ok {
		return nil, fmt.Errorf("unknown environment: %s", environmentID)
	}

	contracts, err := s.querier.GetAllContractsForParty(ctx, env.Party)
	if err != nil {
		return nil, fmt.Errorf("failed to query contracts: %w", err)
	}

	router, err := s.findContract(contracts, ContractTypeRouter, env.Contracts.Router)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	onRamp, err := s.findContract(contracts, ContractTypeOnRamp, env.Contracts.OnRamp)
	if err != nil {
		return nil, fmt.Errorf("onRamp not found: %w", err)
	}

	feeQuoter, err := s.findContract(contracts, ContractTypeFeeQuoter, env.Contracts.FeeQuoter)
	if err != nil {
		return nil, fmt.Errorf("feeQuoter not found: %w", err)
	}

	return &types.CCIPSendDisclosures{
		EnvironmentID: environmentID,
		Contracts: types.CCIPSendContracts{
			Router:    router,
			OnRamp:    onRamp,
			FeeQuoter: feeQuoter,
		},
	}, nil
}

func (s *Service) GetCCIPExecuteDisclosures(ctx context.Context, environmentID string) (*types.CCIPExecuteDisclosures, error) {
	env, ok := s.envConfig.GetEnvironment(environmentID)
	if !ok {
		return nil, fmt.Errorf("unknown environment: %s", environmentID)
	}

	contracts, err := s.querier.GetAllContractsForParty(ctx, env.Party)
	if err != nil {
		return nil, fmt.Errorf("failed to query contracts: %w", err)
	}

	offRamp, err := s.findContract(contracts, ContractTypeOffRamp, env.Contracts.OffRamp)
	if err != nil {
		return nil, fmt.Errorf("offRamp not found: %w", err)
	}

	ccv, err := s.findContract(contracts, ContractTypeCCV, env.Contracts.CCV)
	if err != nil {
		return nil, fmt.Errorf("ccv not found: %w", err)
	}

	tokenAdminRegistry, err := s.findContract(contracts, ContractTypeTokenAdminRegistry, env.Contracts.TokenAdminRegistry)
	if err != nil {
		return nil, fmt.Errorf("tokenAdminRegistry not found: %w", err)
	}

	return &types.CCIPExecuteDisclosures{
		EnvironmentID: environmentID,
		Contracts: types.CCIPExecuteContracts{
			OffRamp:            offRamp,
			CCV:                ccv,
			TokenAdminRegistry: tokenAdminRegistry,
		},
	}, nil
}

// find contract by type and environmentId, matching module:entity (ignores packageId for upgrade resilience)
func (s *Service) findContract(
	contracts []*ledger.ActiveContract,
	contractType ContractType,
	expectedEnvID string,
) (*types.DisclosedContract, error) {
	moduleInfo, ok := ContractTypeToModule[contractType]
	if !ok {
		return nil, fmt.Errorf("unknown contract type: %s", contractType)
	}

	for _, c := range contracts {
		event := c.CreatedEvent
		templateID := event.GetTemplateId()

		if templateID.GetModuleName() != moduleInfo.ModuleName ||
			templateID.GetEntityName() != moduleInfo.EntityName {
			continue
		}

		envID := ExtractEnvironmentID(event.GetCreateArguments())
		if envID != expectedEnvID {
			continue
		}

		return &types.DisclosedContract{
			ContractID: event.GetContractId(),
			TemplateID: types.TemplateID{
				PackageID:  templateID.GetPackageId(),
				ModuleName: templateID.GetModuleName(),
				EntityName: templateID.GetEntityName(),
			},
			CreatedEventBlob: base64.StdEncoding.EncodeToString(event.GetCreatedEventBlob()),
			SynchronizerID:   c.SynchronizerID,
		}, nil
	}

	return nil, fmt.Errorf("no contract found with type %s and environmentId %s", contractType, expectedEnvID)
}

func ExtractEnvironmentID(args *apiv2.Record) string {
	if args == nil {
		return ""
	}
	for _, field := range args.GetFields() {
		if field.GetLabel() == "environmentId" {
			if textVal, ok := field.GetValue().GetSum().(*apiv2.Value_Text); ok {
				return textVal.Text
			}
		}
	}
	return ""
}
