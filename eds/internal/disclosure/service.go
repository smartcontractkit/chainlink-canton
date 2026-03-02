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
	CCVs                  []contracts.InstanceAddress
}

// The DisclosureService returns explicit disclosures for CCIP contracts.
// It uses a store.ContractStore to retrieve active contracts from the ledger and provides explicit disclosures for them.
// It is configured with a fixed list of InstanceAddresses for all CCIP contracts.
type DisclosureService struct {
	contractStore store.ContractStore

	perPartyRouterFactory contracts.InstanceAddress
	onRamp                contracts.InstanceAddress
	offRamp               contracts.InstanceAddress
	globalConfig          contracts.InstanceAddress
	tokenAdminRegistry    contracts.InstanceAddress
	rmnRemote             contracts.InstanceAddress
	ccvs                  []contracts.InstanceAddress

	// Contains all configured instance addresses, to allow looking up if a requested disclosure should be returned.
	allContracts map[contracts.InstanceAddress]struct{}
}

func NewDisclosureService(ctx context.Context, config DisclosureServiceConfig) *DisclosureService {
	// Create a map of all instance addresses
	allContracts := make(map[contracts.InstanceAddress]struct{}, 6+len(config.CCVs))
	allContracts[config.PerPartyRouterFactory] = struct{}{}
	allContracts[config.OnRamp] = struct{}{}
	allContracts[config.OffRamp] = struct{}{}
	allContracts[config.GlobalConfig] = struct{}{}
	allContracts[config.TokenAdminRegistry] = struct{}{}
	allContracts[config.RMNRemote] = struct{}{}
	for _, ccv := range config.CCVs {
		allContracts[ccv] = struct{}{}
	}

	return &DisclosureService{
		contractStore: config.ContractStore,

		perPartyRouterFactory: config.PerPartyRouterFactory,
		onRamp:                config.OnRamp,
		offRamp:               config.OffRamp,
		globalConfig:          config.GlobalConfig,
		tokenAdminRegistry:    config.TokenAdminRegistry,
		rmnRemote:             config.RMNRemote,
		ccvs:                  config.CCVs,

		allContracts: allContracts,
	}
}

type CCIPSendRequest struct {
	CCVs []contracts.InstanceAddress
}

type CCIPSendDisclosures struct {
	OnRamp             *apiv2.DisclosedContract
	GlobalConfig       *apiv2.DisclosedContract
	TokenAdminRegistry *apiv2.DisclosedContract
	RMNRemote          *apiv2.DisclosedContract
	CCVs               map[contracts.InstanceAddress]*apiv2.DisclosedContract
}

// GetDisclosure returns the explicit disclosure for a given contract address, or an error if the contract is not found or not configured for disclosure.
func (s *DisclosureService) GetDisclosure(ctx context.Context, address contracts.InstanceAddress) (*apiv2.DisclosedContract, error) {
	// Only allow looking up contracts that have explicitly been configured.
	// Since there could be multiple contracts of the same template deployed on any given network, the ContractStore will
	// return all of them. This could expose contracts that should be kept private, therefore we're only returning
	// explicitly configured contracts.
	if _, ok := s.allContracts[address]; !ok {
		return nil, fmt.Errorf("contract %s not found", address)
	}

	activeContract := s.contractStore.GetContract(ctx, address)
	if activeContract == nil {
		return nil, fmt.Errorf("contract %s not found", address)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       activeContract.GetCreatedEvent().GetTemplateId(),
		ContractId:       activeContract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: activeContract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   activeContract.GetSynchronizerId(),
	}, nil
}

func (s *DisclosureService) GetCCIPSendDisclosures(ctx context.Context, request CCIPSendRequest) (CCIPSendDisclosures, error) {
	onRamp, err := s.GetDisclosure(ctx, s.onRamp)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("onRamp: %w", err)
	}
	globalConfig, err := s.GetDisclosure(ctx, s.globalConfig)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("globalConfig: %w", err)
	}
	tokenAdminRegistry, err := s.GetDisclosure(ctx, s.tokenAdminRegistry)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("tokenAdminRegistry: %w", err)
	}
	rmnRemote, err := s.GetDisclosure(ctx, s.rmnRemote)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("rmnRemote: %w", err)
	}

	// CCVs
	ccvs := make(map[contracts.InstanceAddress]*apiv2.DisclosedContract, len(request.CCVs))
	for _, requestedCCV := range request.CCVs {
		ccv, err := s.GetDisclosure(ctx, requestedCCV)
		if err != nil {
			return CCIPSendDisclosures{}, fmt.Errorf("ccv: %w", err)
		}

		ccvs[requestedCCV] = ccv
	}

	return CCIPSendDisclosures{
		OnRamp:             onRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		CCVs:               ccvs,
	}, nil
}

type CCIPExecuteRequest struct {
	CCVs []contracts.InstanceAddress
}

type CCIPExecuteDisclosures struct {
	OffRamp            *apiv2.DisclosedContract
	GlobalConfig       *apiv2.DisclosedContract
	TokenAdminRegistry *apiv2.DisclosedContract
	RMNRemote          *apiv2.DisclosedContract
	CCVs               map[contracts.InstanceAddress]*apiv2.DisclosedContract
}

func (s *DisclosureService) GetCCIPExecuteDisclosures(ctx context.Context, request CCIPExecuteRequest) (CCIPExecuteDisclosures, error) {
	offRamp, err := s.GetDisclosure(ctx, s.offRamp)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("offRamp: %w", err)
	}
	globalConfig, err := s.GetDisclosure(ctx, s.globalConfig)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("globalConfig: %w", err)
	}
	tokenAdminRegistry, err := s.GetDisclosure(ctx, s.tokenAdminRegistry)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("tokenAdminRegistry: %w", err)
	}
	rmnRemote, err := s.GetDisclosure(ctx, s.rmnRemote)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("rmnRemote: %w", err)
	}

	// CCVs
	ccvs := make(map[contracts.InstanceAddress]*apiv2.DisclosedContract, len(request.CCVs))
	for _, requestedCCV := range request.CCVs {
		ccv, err := s.GetDisclosure(ctx, requestedCCV)
		if err != nil {
			return CCIPExecuteDisclosures{}, fmt.Errorf("ccv: %w", err)
		}

		ccvs[requestedCCV] = ccv
	}

	return CCIPExecuteDisclosures{
		OffRamp:            offRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		CCVs:               ccvs,
	}, nil
}

type PerPartyRouterFactoryRequest struct{}

type PerPartyRouterFactoryDisclosures struct {
	PerPartyRouterFactory *apiv2.DisclosedContract
}

func (s *DisclosureService) GetPerPartyRouterFactory(ctx context.Context, _ PerPartyRouterFactoryRequest) (PerPartyRouterFactoryDisclosures, error) {
	perPartyRouterFactory, err := s.GetDisclosure(ctx, s.perPartyRouterFactory)
	if err != nil {
		return PerPartyRouterFactoryDisclosures{}, fmt.Errorf("perPartyRouterFactory: %w", err)
	}

	return PerPartyRouterFactoryDisclosures{
		PerPartyRouterFactory: perPartyRouterFactory,
	}, nil
}
