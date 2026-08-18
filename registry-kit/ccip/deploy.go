package ccip

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ratelimiter"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

// PoolDeployDeps holds CCIP contract bindings required to create registrar-owned token pools.
type PoolDeployDeps struct {
	CcipOwner          string
	TokenAdminRegistry contracts.RawInstanceAddress
	RMNRemote          contracts.RawInstanceAddress
	FeeQuoter          contracts.RawInstanceAddress
}

// DeployInboundRateLimiterForOwner creates an inbound rate limiter with ownerParty as actAs.
func DeployInboundRateLimiterForOwner(
	ctx context.Context,
	client ledger.Client,
	ownerParty string,
	template ratelimiter.RateLimiter,
) (contracts.RawInstanceAddress, error) {
	res, err := client.SubmitCreate(ctx, ownerParty, template)
	if err != nil {
		return "", fmt.Errorf("deploy inbound rate limiter: %w", err)
	}
	if _, ok := ledger.CreatedContractID(res.GetTransaction(), "RateLimiter"); !ok {
		return "", fmt.Errorf("RateLimiter not created")
	}

	raw, err := contracts.RawInstanceAddressFromString(fmt.Sprintf("%s@%s", template.InstanceId, ownerParty))
	if err != nil {
		return "", err
	}

	return raw, nil
}

// DeployOutboundRateLimiterForOwner creates an outbound rate limiter with ownerParty as actAs.
func DeployOutboundRateLimiterForOwner(
	ctx context.Context,
	client ledger.Client,
	ownerParty string,
	template ratelimiter.RateLimiter,
) (contracts.RawInstanceAddress, error) {
	res, err := client.SubmitCreate(ctx, ownerParty, template)
	if err != nil {
		return "", fmt.Errorf("deploy outbound rate limiter: %w", err)
	}
	if _, ok := ledger.CreatedContractID(res.GetTransaction(), "RateLimiter"); !ok {
		return "", fmt.Errorf("RateLimiter not created")
	}

	raw, err := contracts.RawInstanceAddressFromString(fmt.Sprintf("%s@%s", template.InstanceId, ownerParty))
	if err != nil {
		return "", err
	}

	return raw, nil
}

// DeployBurnMintPoolForOwner creates a BurnMintTokenPool with poolOwner as actAs.
func DeployBurnMintPoolForOwner(
	ctx context.Context,
	client ledger.Client,
	deps PoolDeployDeps,
	poolOwner string,
	instrumentID splice_api_token_holding_v1.InstrumentId,
	poolInstanceID string,
	remoteChainConfigs map[types.NUMERIC]burnminttokenpool.RemoteChainConfig,
	feeConfigs ...map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig,
) (contracts.RawInstanceAddress, error) {
	if remoteChainConfigs == nil {
		remoteChainConfigs = map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{}
	}
	tokenTransferFeeConfigs := map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{}
	if len(feeConfigs) > 0 && feeConfigs[0] != nil {
		tokenTransferFeeConfigs = feeConfigs[0]
	}

	pool := burnminttokenpool.BurnMintTokenPool{
		InstanceId:              types.TEXT(poolInstanceID),
		PoolOwner:               types.PARTY(poolOwner),
		CcipOwner:               types.PARTY(deps.CcipOwner),
		InstrumentId:            instrumentID,
		Decimals:                types.INT64(10),
		RemoteChainConfigs:      remoteChainConfigs,
		TokenTransferFeeConfigs: tokenTransferFeeConfigs,
		TransferTimeout:         burnminttokenpool.TransferTimeout{RelativeHours: new(types.INT64(24))},
		Deps: burnminttokenpool.BurnMintTokenPoolDeps{
			TokenAdminRegistry: deps.TokenAdminRegistry.Binding(),
			RmnRemote:          deps.RMNRemote.Binding(),
			FeeQuoter:          deps.FeeQuoter.Binding(),
		},
	}

	res, err := client.SubmitCreate(ctx, poolOwner, pool)
	if err != nil {
		return "", fmt.Errorf("deploy burn mint token pool: %w", err)
	}
	if _, ok := ledger.CreatedContractID(res.GetTransaction(), "BurnMintTokenPool"); !ok {
		return "", fmt.Errorf("BurnMintTokenPool not created")
	}

	return contracts.NewRawInstanceAddress(contracts.InstanceID(poolInstanceID), types.PARTY(poolOwner)), nil
}
