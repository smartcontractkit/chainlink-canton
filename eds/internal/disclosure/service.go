package disclosure

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"slices"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

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
	ccvs                  []contracts.InstanceAddress
	tokenPools            []contracts.InstanceAddress

	// Contains all configured instance addresses, to allow looking up if a requested disclosure should be returned.
	allContracts map[contracts.InstanceAddress]struct{}
}

func NewDisclosureService(ctx context.Context, config DisclosureServiceConfig) *DisclosureService {
	// Create a map of all instance addresses
	allContracts := make(map[contracts.InstanceAddress]struct{}, 7+len(config.CCVs)+len(config.TokenPools))
	allContracts[config.PerPartyRouterFactory] = struct{}{}
	allContracts[config.OnRamp] = struct{}{}
	allContracts[config.OffRamp] = struct{}{}
	allContracts[config.GlobalConfig] = struct{}{}
	allContracts[config.TokenAdminRegistry] = struct{}{}
	allContracts[config.RMNRemote] = struct{}{}
	allContracts[config.FeeQuoter] = struct{}{}
	for _, ccv := range config.CCVs {
		allContracts[ccv] = struct{}{}
	}
	for _, tokenPool := range config.TokenPools {
		allContracts[tokenPool] = struct{}{}
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
		ccvs:                  config.CCVs,
		tokenPools:            config.TokenPools,

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

	// These disclosures are optional, only will be returned in the event
	// there is a token transfer.
	CCVs                                       map[contracts.InstanceAddress]*apiv2.DisclosedContract
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

	tokenPoolDisclosure, instrumentHoldingDisclosure, inboundRateLimiterDisclosure, inboundCustomBlockConfirmationsRateLimiterDisclosure, outboundRateLimiterDisclosure, err := s.getTokenPoolRelatedDisclosures(ctx, request.Message)
	if err != nil {
		return CCIPExecuteDisclosures{}, fmt.Errorf("can't get token pool and holding disclosure: %w", err)
	}

	return CCIPExecuteDisclosures{
		OffRamp:            offRamp,
		GlobalConfig:       globalConfig,
		TokenAdminRegistry: tokenAdminRegistry,
		RMNRemote:          rmnRemote,
		CCVs:               ccvs,
		TokenPool:          tokenPoolDisclosure,
		TokenPoolHolding:   instrumentHoldingDisclosure,
		InboundRateLimiter: inboundRateLimiterDisclosure,
		InboundCustomBlockConfirmationsRateLimiter: inboundCustomBlockConfirmationsRateLimiterDisclosure,
		OutboundRateLimiter:                        outboundRateLimiterDisclosure,
	}, nil
}

func (s *DisclosureService) getTokenPoolRelatedDisclosures(
	ctx context.Context,
	message *protocol.Message,
) (
	tokenPoolDisclosure *apiv2.DisclosedContract,
	instrumentHoldingDisclosure *apiv2.DisclosedContract,
	inboundRateLimiterDisclosure *apiv2.DisclosedContract,
	inboundCustomBlockConfirmationsRateLimiterDisclosure *apiv2.DisclosedContract,
	outboundRateLimiterDisclosure *apiv2.DisclosedContract,
	err error,
) {
	if message == nil {
		// Nothing to do, can't provide disclosures for token pool and holdings if we don't have message data
		return nil, nil, nil, nil, nil, nil
	}

	// To get the token pool, we need to query the token admin registry for the token pool instance address.
	if message.TokenTransfer != nil &&
		len(message.TokenTransfer.DestTokenAddress) > 0 {
		// get the created event of the token admin registry via the update store.
		tarActiveContract := s.contractStore.GetContract(ctx, s.tokenAdminRegistry)
		if tarActiveContract == nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get tokenAdminRegistry disclosure from update store: %w", err)
		}
		tarCreatedEvent, err := bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenAdminRegistry](tarActiveContract.GetCreatedEvent())
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't unmarshal tokenAdminRegistry created event: %w", err)
		}

		// get the token pool instance address from the created event.
		// Note that the destTokenAddress is the hashed instrumentId, which is what the mapping
		// in the TokenAdminRegistry uses as a key.
		tokenPoolInstanceAddress, instrumentID, err := getTokenPoolAddressAndInstrumentID(tarCreatedEvent, message.TokenTransfer.DestTokenAddress)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get token pool instance address from token admin registry createdEvent: %w", err)
		}

		tokenPoolDisclosure, err = s.GetDisclosure(ctx, tokenPoolInstanceAddress)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get token pool disclosure: %w", err)
		}

		// TODO: get the rate limiters from the tp created event.
		tpActiveContract := s.contractStore.GetContract(ctx, tokenPoolInstanceAddress)
		if tpActiveContract == nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get token pool disclosure from update store: %w", err)
		}

		tpCreatedEvent, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](tpActiveContract.GetCreatedEvent())
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't unmarshal token pool created event: %w", err)
		}

		inboundRateLimiterInstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress, outboundRateLimiterInstanceAddress, err := getRateLimiterInstanceAddresses(
			uint64(message.SourceChainSelector),
			tpCreatedEvent,
			tpCreatedEvent.PoolOwner,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get rate limiter instance addresses from token pool: %w", err)
		}

		inboundRateLimiterDisclosure, err = s.GetDisclosure(ctx, inboundRateLimiterInstanceAddress)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get inbound rate limiter disclosure: %w", err)
		}

		inboundCustomBlockConfirmationsRateLimiterDisclosure, err = s.GetDisclosure(ctx, inboundCustomBlockConfirmationsRateLimiterInstanceAddress)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get inbound custom block confirmations rate limiter disclosure: %w", err)
		}

		outboundRateLimiterDisclosure, err = s.GetDisclosure(ctx, outboundRateLimiterInstanceAddress)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get outbound rate limiter disclosure: %w", err)
		}

		// look up the instrument holding disclosure for the instrument id.
		instrumentHoldingDisclosure, err = s.instrumentHoldingStore.GetInstrumentHolding(ctx, instrumentID)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("can't get instrument holding disclosure: %w", err)
		}

		return tokenPoolDisclosure, instrumentHoldingDisclosure, inboundRateLimiterDisclosure, inboundCustomBlockConfirmationsRateLimiterDisclosure, outboundRateLimiterDisclosure, nil
	}

	// If a message doesn't have a token transfer there will be no disclosures for either the token pool or the instrument holding.
	return nil, nil, nil, nil, nil, nil
}

func getRateLimiterInstanceAddresses(
	remoteChainSelector uint64,
	tpCreatedEvent *lockreleasetokenpool.LockReleaseTokenPool,
	poolOwner types.PARTY,
) (inboundRateLimiterInstanceAddress contracts.InstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress contracts.InstanceAddress, outboundRateLimiterInstanceAddress contracts.InstanceAddress, err error) {
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
	inboundRateLimiterInstanceAddress = contracts.InstanceID(remoteChainConfig.InboundRateLimiter.Unpack).
		RawInstanceAddress(poolOwner).
		InstanceAddress()
	inboundCustomBlockConfirmationsRateLimiterInstanceAddress = contracts.InstanceID(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter.Unpack).
		RawInstanceAddress(poolOwner).
		InstanceAddress()
	outboundRateLimiterInstanceAddress = contracts.InstanceID(remoteChainConfig.OutboundRateLimiter.Unpack).
		RawInstanceAddress(poolOwner).
		InstanceAddress()

	return inboundRateLimiterInstanceAddress, inboundCustomBlockConfirmationsRateLimiterInstanceAddress, outboundRateLimiterInstanceAddress, nil
}

// TODO: this is really janky, how to properly test this?
func getTokenPoolAddressAndInstrumentID(tarCreatedEvent *tokenadminregistry.TokenAdminRegistry, destTokenAddress []byte) (contracts.InstanceAddress, splice_api_token_holding_v1.InstrumentId, error) {
	expectedHashedInstrumentID := contracts.BytesToInstanceAddress(destTokenAddress)

	// -- | Maps keccak256(InstrumentId) to token configuration
	// tokenConfigs : Map.Map BytesHex TokenConfig
	for hashedInstrumentID, v := range tarCreatedEvent.TokenConfigs {
		fmt.Printf("processing map: %+v\n", v)

		// check if this instrument ID corresponds to the one we expect.
		instrumentIDBytes, err := hex.DecodeString(hashedInstrumentID)
		if err != nil {
			continue
		}
		if !bytes.Equal(instrumentIDBytes, expectedHashedInstrumentID.Bytes()) {
			fmt.Printf("expected hashed instrument ID: %s, but got: %s, skipping\n", expectedHashedInstrumentID.String(), hashedInstrumentID)
			continue
		}

		vMap, ok := v.(map[string]any)
		if !ok {
			fmt.Printf("v is not a map[string]any, skipping\n")
			continue
		}

		var tokenConfig tokenadminregistry.TokenConfig
		err = ledger.MapToStruct(vMap, &tokenConfig)
		if err != nil {
			fmt.Printf("failed to decode token config: %s, skipping\n", err.Error())
			continue
		}

		// Construct the pool instance address from the token config.
		poolInstanceAddress := contracts.InstanceID(tokenConfig.TokenPool.PoolInstanceId).
			RawInstanceAddress(tokenConfig.TokenPool.PoolOwner).
			InstanceAddress()

		return poolInstanceAddress, tokenConfig.InstrumentId, nil
	}

	return contracts.InstanceAddress{}, splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf(
		"token pool instance address not found in token admin registry data (dest token address: %s): %+v",
		expectedHashedInstrumentID.String(),
		tarCreatedEvent.TokenConfigs,
	)
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
