package disclosure

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
)

type DisclosureServiceConfig struct {
	ContractStore store.ContractStore

	// Contracts
	PerPartyRouterFactory contracts.InstanceAddress
	OnRamp                contracts.InstanceAddress
	OffRamp               contracts.InstanceAddress
	GlobalConfig          contracts.InstanceAddress
	TokenAdminRegistry    contracts.InstanceAddress
	RMNRemote             contracts.InstanceAddress
	DefaultCCV            contracts.InstanceAddress
}

type DisclosureService struct {
	contractStore store.ContractStore

	perPartyRouterFactory contracts.InstanceAddress
	onRamp                contracts.InstanceAddress
	offRamp               contracts.InstanceAddress
	globalConfig          contracts.InstanceAddress
	tokenAdminRegistry    contracts.InstanceAddress
	rmnRemote             contracts.InstanceAddress
	defaultCCV            contracts.InstanceAddress
}

func NewDisclosureService(ctx context.Context, config DisclosureServiceConfig) *DisclosureService {
	return &DisclosureService{
		contractStore: config.ContractStore,

		perPartyRouterFactory: config.PerPartyRouterFactory,
		onRamp:                config.OnRamp,
		offRamp:               config.OffRamp,
		globalConfig:          config.GlobalConfig,
		tokenAdminRegistry:    config.TokenAdminRegistry,
		rmnRemote:             config.RMNRemote,
		defaultCCV:            config.DefaultCCV,
	}
}

type CCIPSendRequest struct{}

type CCIPSendDisclosures struct {
	OnRamp             *apiv2.DisclosedContract
	GlobalConfig       *apiv2.DisclosedContract
	TokenAdminRegistry *apiv2.DisclosedContract
	RMNRemote          *apiv2.DisclosedContract
	DefaultCCV         *apiv2.DisclosedContract
}

func (s *DisclosureService) GetCCIPSendDisclosures(ctx context.Context, _ CCIPSendRequest) (CCIPSendDisclosures, error) {
	onRamp, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.onRamp))
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("onRamp: %w", err)
	}
	globalConfig, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.globalConfig))
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("globalConfig: %w", err)
	}
	tokenAdminRegistry, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.tokenAdminRegistry))
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("tokenAdminRegistry: %w", err)
	}
	rmnRemote, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.rmnRemote))
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("rmnRemote: %w", err)
	}
	defaultCCV, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.defaultCCV))
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("defaultCCV: %w", err)
	}

	return CCIPSendDisclosures{
		OnRamp:             onRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		DefaultCCV:         defaultCCV,
	}, nil
}

type CCIPExecuteRequest struct{}

type CCIPExecuteDisclosures struct {
	OffRamp            *apiv2.DisclosedContract
	GlobalConfig       *apiv2.DisclosedContract
	TokenAdminRegistry *apiv2.DisclosedContract
	RMNRemote          *apiv2.DisclosedContract
	DefaultCCV         *apiv2.DisclosedContract
}

func (s *DisclosureService) GetCCIPExecuteDisclosures(ctx context.Context, _ CCIPExecuteRequest) (CCIPExecuteDisclosures, error) {
	offRamp, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.offRamp))
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("offRamp (%s): %w", s.offRamp, err)
	}
	globalConfig, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.globalConfig))
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("globalConfig (%s): %w", s.globalConfig, err)
	}
	tokenAdminRegistry, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.tokenAdminRegistry))
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("tokenAdminRegistry (%s): %w", s.tokenAdminRegistry, err)
	}
	rmnRemote, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.rmnRemote))
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("rmnRemote (%s): %w", s.rmnRemote, err)
	}
	defaultCCV, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.defaultCCV))
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("defaultCCV (%s): %w", s.defaultCCV, err)
	}

	return CCIPExecuteDisclosures{
		OffRamp:            offRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		DefaultCCV:         defaultCCV,
	}, nil
}

type PerPartyRouterFactoryRequest struct{}

type PerPartyRouterFactoryDisclosures struct {
	PerPartyRouterFactory *apiv2.DisclosedContract
}

func (s *DisclosureService) GetPerPartyRouterFactory(ctx context.Context, _ PerPartyRouterFactoryRequest) (PerPartyRouterFactoryDisclosures, error) {
	perPartyRouterFactory, err := getDisclosedContract(s.contractStore.GetContract(ctx, s.perPartyRouterFactory))
	if err != nil {
		return PerPartyRouterFactoryDisclosures{}, fmt.Errorf("perPartyRouterFactory (%s): %w", s.perPartyRouterFactory, err)
	}

	return PerPartyRouterFactoryDisclosures{
		PerPartyRouterFactory: perPartyRouterFactory,
	}, nil
}

func getDisclosedContract(contract *apiv2.ActiveContract) (*apiv2.DisclosedContract, error) {
	if contract == nil {
		return nil, fmt.Errorf("contract not found")
	}

	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}
