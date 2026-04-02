package disclosure

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
)

type DisclosureServiceConfig struct {
	ContractStore          store.ContractStore
	InstrumentHoldingStore store.InstrumentHoldingStore

	// Contracts
	PerPartyRouterFactory contracts.InstanceAddress
	OnRamp                contracts.InstanceAddress
	OffRamp               contracts.InstanceAddress
	GlobalConfig          contracts.InstanceAddress
	TokenAdminRegistry    contracts.InstanceAddress
	RMNRemote             contracts.InstanceAddress
	FeeQuoter             contracts.InstanceAddress
	DefaultExecutor       contracts.InstanceAddress
	CCVs                  []contracts.InstanceAddress

	TokenPools                                   []contracts.InstanceAddress
	TokenPoolInboundRateLimiters                 []contracts.InstanceAddress
	TokenPoolRateLimiterCustomBlockConfirmations []contracts.InstanceAddress
	TokenPoolOutboundRateLimiters                []contracts.InstanceAddress
}

// The DisclosureService returns explicit disclosures for CCIP contracts.
// It uses a store.ContractStore to retrieve active contracts from the ledger and provides explicit disclosures for them.
// It is configured with a fixed list of InstanceAddresses for all CCIP contracts.
type DisclosureService struct {
	contractStore          store.ContractStore
	instrumentHoldingStore store.InstrumentHoldingStore

	perPartyRouterFactory contracts.InstanceAddress
	onRamp                contracts.InstanceAddress
	offRamp               contracts.InstanceAddress
	globalConfig          contracts.InstanceAddress
	tokenAdminRegistry    contracts.InstanceAddress
	rmnRemote             contracts.InstanceAddress
	feeQuoter             contracts.InstanceAddress
	defaultExecutor       contracts.InstanceAddress
	ccvs                  []contracts.InstanceAddress

	// Contains all configured instance addresses, to allow looking up if a requested disclosure should be returned.
	allContracts map[contracts.InstanceAddress]struct{}
}

func NewDisclosureService(ctx context.Context, config DisclosureServiceConfig) *DisclosureService {
	// Create a map of all instance addresses
	allContracts := make(
		map[contracts.InstanceAddress]struct{},
		7+ // perPartyRouterFactory, onRamp, offRamp, globalConfig, tokenAdminRegistry, rmnRemote, feeQuoter
			len(config.CCVs)+
			len(config.TokenPools)+
			len(config.TokenPoolInboundRateLimiters)+
			len(config.TokenPoolRateLimiterCustomBlockConfirmations)+
			len(config.TokenPoolOutboundRateLimiters),
	)
	allContracts[config.PerPartyRouterFactory] = struct{}{}
	allContracts[config.OnRamp] = struct{}{}
	allContracts[config.OffRamp] = struct{}{}
	allContracts[config.GlobalConfig] = struct{}{}
	allContracts[config.TokenAdminRegistry] = struct{}{}
	allContracts[config.RMNRemote] = struct{}{}
	allContracts[config.FeeQuoter] = struct{}{}
	allContracts[config.DefaultExecutor] = struct{}{}
	for _, ccv := range config.CCVs {
		allContracts[ccv] = struct{}{}
	}
	for _, tokenPool := range config.TokenPools {
		allContracts[tokenPool] = struct{}{}
	}
	for _, tokenPoolInboundRateLimiter := range config.TokenPoolInboundRateLimiters {
		allContracts[tokenPoolInboundRateLimiter] = struct{}{}
	}
	for _, tokenPoolRateLimiterCustomBlockConfirmations := range config.TokenPoolRateLimiterCustomBlockConfirmations {
		allContracts[tokenPoolRateLimiterCustomBlockConfirmations] = struct{}{}
	}
	for _, tokenPoolOutboundRateLimiter := range config.TokenPoolOutboundRateLimiters {
		allContracts[tokenPoolOutboundRateLimiter] = struct{}{}
	}

	return &DisclosureService{
		contractStore:          config.ContractStore,
		instrumentHoldingStore: config.InstrumentHoldingStore,

		perPartyRouterFactory: config.PerPartyRouterFactory,
		onRamp:                config.OnRamp,
		offRamp:               config.OffRamp,
		globalConfig:          config.GlobalConfig,
		tokenAdminRegistry:    config.TokenAdminRegistry,
		rmnRemote:             config.RMNRemote,
		feeQuoter:             config.FeeQuoter,
		defaultExecutor:       config.DefaultExecutor,
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
	FeeQuoter          *apiv2.DisclosedContract
	DefaultExecutor    *apiv2.DisclosedContract
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
	defaultExecutor, err := s.GetDisclosure(ctx, s.defaultExecutor)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("defaultExecutor: %w", err)
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

	feeQuoter, err := s.GetDisclosure(ctx, s.feeQuoter)
	if err != nil {
		return CCIPSendDisclosures{}, fmt.Errorf("feeQuoter: %w", err)
	}

	return CCIPSendDisclosures{
		OnRamp:             onRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		FeeQuoter:          feeQuoter,
		DefaultExecutor:    defaultExecutor,
		CCVs:               ccvs,
	}, nil
}

type CCIPExecuteRequest struct {
	Message *protocol.Message
	CCVs    []contracts.InstanceAddress
}

type CCIPExecuteDisclosures struct {
	// These disclosures are always needed for execution.
	OffRamp            *apiv2.DisclosedContract
	GlobalConfig       *apiv2.DisclosedContract
	TokenAdminRegistry *apiv2.DisclosedContract
	RMNRemote          *apiv2.DisclosedContract
	CCVs               map[contracts.InstanceAddress]*apiv2.DisclosedContract

	// These disclosures are optional, only will be returned in the event
	// there is a token transfer.
	TokenPool                                  *apiv2.DisclosedContract
	TokenPoolHolding                           *apiv2.DisclosedContract
	InboundRateLimiter                         *apiv2.DisclosedContract
	InboundCustomBlockConfirmationsRateLimiter *apiv2.DisclosedContract
	OutboundRateLimiter                        *apiv2.DisclosedContract
}

// tokenPoolRelatedDisclosures groups optional execute disclosures resolved from a token transfer message.
type tokenPoolRelatedDisclosures struct {
	TokenPool                                  *apiv2.DisclosedContract
	TokenPoolHolding                           *apiv2.DisclosedContract
	InboundRateLimiter                         *apiv2.DisclosedContract
	InboundCustomBlockConfirmationsRateLimiter *apiv2.DisclosedContract
	OutboundRateLimiter                        *apiv2.DisclosedContract
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

	extras, err := s.getTokenPoolRelatedDisclosures(ctx, request.Message)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("can't get token pool and holding disclosure: %w", err)
	}

	return CCIPExecuteDisclosures{
		OffRamp:            offRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		CCVs:               ccvs,
		TokenPool:          extras.TokenPool,
		TokenPoolHolding:   extras.TokenPoolHolding,
		InboundRateLimiter: extras.InboundRateLimiter,
		InboundCustomBlockConfirmationsRateLimiter: extras.InboundCustomBlockConfirmationsRateLimiter,
		OutboundRateLimiter:                        extras.OutboundRateLimiter,
	}, nil
}

func (s *DisclosureService) getTokenPoolRelatedDisclosures(ctx context.Context, message *protocol.Message) (tokenPoolRelatedDisclosures, error) {
	if message == nil {
		// Nothing to do, can't provide disclosures for token pool and holdings if we don't have message data
		return tokenPoolRelatedDisclosures{}, nil
	}

	// To get the token pool, we need to query the token admin registry for the token pool instance address.
	if message.TokenTransfer == nil || len(message.TokenTransfer.DestTokenAddress) == 0 {
		// If a message doesn't have a token transfer there will be no disclosures for either the token pool or the instrument holding.
		return tokenPoolRelatedDisclosures{}, nil
	}

	// get the created event of the token admin registry via the update store.
	tarActiveContract := s.contractStore.GetContract(ctx, s.tokenAdminRegistry)
	if tarActiveContract == nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("token admin registry contract not found in update store (instance: %v)", s.tokenAdminRegistry)
	}
	tarCreatedEvent, err := bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenAdminRegistry](tarActiveContract.GetCreatedEvent())
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't unmarshal tokenAdminRegistry created event: %w", err)
	}

	// get the token pool instance address from the created event.
	// Note that the destTokenAddress is the hashed instrumentId, which is what the mapping
	// in the TokenAdminRegistry uses as a key.
	tokenPoolInstanceAddress, instrumentID, err := getTokenPoolAddressAndInstrumentID(tarCreatedEvent, message.TokenTransfer.DestTokenAddress)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get token pool instance address from token admin registry createdEvent: %w", err)
	}

	tokenPoolDisclosure, err := s.GetDisclosure(ctx, tokenPoolInstanceAddress)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get token pool disclosure: %w", err)
	}

	tpActiveContract := s.contractStore.GetContract(ctx, tokenPoolInstanceAddress)
	if tpActiveContract == nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("token pool contract not found in update store (instance: %v)", tokenPoolInstanceAddress)
	}

	tpCreatedEvent, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](tpActiveContract.GetCreatedEvent())
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't unmarshal token pool created event: %w", err)
	}

	inboundRateLimiterInstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress, outboundRateLimiterInstanceAddress, err := getRateLimiterInstanceAddresses(uint64(message.SourceChainSelector), tpCreatedEvent)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get rate limiter instance addresses from token pool: %w", err)
	}

	inboundRateLimiterDisclosure, err := s.GetDisclosure(ctx, inboundRateLimiterInstanceAddress)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get inbound rate limiter disclosure: %w", err)
	}

	inboundCustomBlockConfirmationsRateLimiterDisclosure, err := s.GetDisclosure(ctx, inboundCustomBlockConfirmationsRateLimiterInstanceAddress)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get inbound custom block confirmations rate limiter disclosure: %w", err)
	}

	outboundRateLimiterDisclosure, err := s.GetDisclosure(ctx, outboundRateLimiterInstanceAddress)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get outbound rate limiter disclosure: %w", err)
	}

	// look up the instrument holding disclosure for the instrument id.
	instrumentHoldingDisclosure, err := s.instrumentHoldingStore.GetInstrumentHolding(ctx, instrumentID)
	if err != nil {
		return tokenPoolRelatedDisclosures{}, fmt.Errorf("can't get instrument holding disclosure: %w", err)
	}

	return tokenPoolRelatedDisclosures{
		TokenPool:          tokenPoolDisclosure,
		TokenPoolHolding:   instrumentHoldingDisclosure,
		InboundRateLimiter: inboundRateLimiterDisclosure,
		InboundCustomBlockConfirmationsRateLimiter: inboundCustomBlockConfirmationsRateLimiterDisclosure,
		OutboundRateLimiter:                        outboundRateLimiterDisclosure,
	}, nil
}

func getRateLimiterInstanceAddresses(remoteChainSelector uint64, tpCreatedEvent *lockreleasetokenpool.LockReleaseTokenPool) (inboundRateLimiterInstanceAddress contracts.InstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress contracts.InstanceAddress, outboundRateLimiterInstanceAddress contracts.InstanceAddress, err error) {
	// look up the remote chain config for the remoteChainSelector provided
	remoteChainConfigAny, ok := tpCreatedEvent.RemoteChainConfigs[fmt.Sprintf("%d.", remoteChainSelector)]
	if !ok {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf(
			"remote chain config not found for remote chain selector: %d, keys: %+v", remoteChainSelector, slices.Collect(maps.Keys(tpCreatedEvent.RemoteChainConfigs)))
	}

	remoteChainConfigMap, ok := remoteChainConfigAny.(map[string]any)
	if !ok {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf("remote chain config is not a map[string]any")
	}

	// unmarshal the remote chain config using ledger.MapToStruct
	var remoteChainConfig lockreleasetokenpool.RemoteChainConfig
	err = ledger.MapToStruct(remoteChainConfigMap, &remoteChainConfig)
	if err != nil {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf("failed to unmarshal remote chain config: %w", err)
	}

	// Create the instance addresses by combining with the poolOwner
	// Unpack is already of the form <instance-id>@<party-id>, so we need to parse it as such prior to calculating the instance address.
	rawInboundRateLimiterInstanceAddress, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.InboundRateLimiter.Unpack))
	if err != nil {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf("failed to parse inbound rate limiter instance address: %w", err)
	}
	inboundRateLimiterInstanceAddress = rawInboundRateLimiterInstanceAddress.InstanceAddress()
	rawInboundCustomBlockConfirmationsRateLimiterInstanceAddress, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter.Unpack))
	if err != nil {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf("failed to parse inbound custom block confirmations rate limiter instance address: %w", err)
	}
	inboundCustomBlockConfirmationsRateLimiterInstanceAddress = rawInboundCustomBlockConfirmationsRateLimiterInstanceAddress.InstanceAddress()
	rawOutboundRateLimiterInstanceAddress, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.OutboundRateLimiter.Unpack))
	if err != nil {
		return contracts.InstanceAddress{}, contracts.InstanceAddress{}, contracts.InstanceAddress{}, fmt.Errorf("failed to parse outbound rate limiter instance address: %w", err)
	}
	outboundRateLimiterInstanceAddress = rawOutboundRateLimiterInstanceAddress.InstanceAddress()

	return inboundRateLimiterInstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress, outboundRateLimiterInstanceAddress, nil
}

func getTokenPoolAddressAndInstrumentID(tarCreatedEvent *tokenadminregistry.TokenAdminRegistry, destTokenAddress []byte) (contracts.InstanceAddress, splice_api_token_holding_v1.InstrumentId, error) {
	expectedHashedInstrumentID := contracts.BytesToInstanceAddress(destTokenAddress)

	if len(tarCreatedEvent.TokenConfigs) == 0 {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf(
			"token admin registry has no token configs, expected at least one",
		)
	}

	// tokenConfigs : Map.Map BytesHex TokenConfig — keys are hex encodings of the same bytes as destTokenAddress.
	hexNo0x := strings.TrimPrefix(expectedHashedInstrumentID.Hex(), "0x")
	v, ok := tarCreatedEvent.TokenConfigs[hexNo0x]
	if !ok {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf(
			"token pool instance address not found in token admin registry data (dest token address: %s): %+v",
			expectedHashedInstrumentID.String(),
			tarCreatedEvent.TokenConfigs,
		)
	}

	vMap, ok := v.(map[string]any)
	if !ok {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf(
			"token config for hashed instrument id %s has unexpected type %T, should be map[string]any", hexNo0x, v)
	}

	var tokenConfig tokenadminregistry.TokenConfig
	if err := ledger.MapToStruct(vMap, &tokenConfig); err != nil {
		return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf(
			"token config for hashed instrument id %s: %w", hexNo0x, err)
	}

	poolInstanceAddress := contracts.InstanceID(tokenConfig.TokenPool.PoolInstanceId).
		RawInstanceAddress(tokenConfig.TokenPool.PoolOwner).
		InstanceAddress()

	return poolInstanceAddress, tokenConfig.InstrumentId, nil
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
