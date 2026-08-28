package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// discoveryRemoteSelector is an arbitrary remote chain selector used to key the pool's single
// RemoteChainConfigs entry; it just needs to match between the pool's config and the test message.
const discoveryRemoteSelector uint64 = 16015286601757825753

// TestPoolDiscoveryRegistry_BurnMint deploys a BurnMintTokenPool and verifies PoolDiscoveryService
// picks it up via ACS polling, then serves a real send disclosure for it, without any statically
// configured TokenPools entry.
//
// One party acts as both the pool owner and EDS's registry-observer party. That is what makes the
// pool and its rate limiter visible to EDS here - in a real deployment those are separate parties,
// and visibility instead depends on the pool and rate limiter naming the observer party.
func TestPoolDiscoveryRegistry_BurnMint(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t)
	participant := env.Chain.Participants[0]
	party := participant.PartyID

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	runtimeDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
	require.NoError(t, err)
	coreDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	// This single participant plays both roles a real deployment splits: the pool deployer and the
	// host of EDS's registry-observer party. The latter must have both pool DARs vetted to be an
	// informee on either pool type, which is what PoolDiscoveryService subscribes to.
	otherPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{runtimeDar, coreDar, poolDar, otherPoolDar}, participant)
	require.NoError(t, err)

	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("BNM"),
	}

	poolInstanceID := "test-pool-bnm-discovery"
	rateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployOutbound, env.Chain, contract.DeployInput[ratelimiter.RateLimiter]{
		OwnerParty: types.PARTY(party),
		Template: ratelimiter.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceID),
			PoolOwner:           types.PARTY(party),
			RemoteChainSelector: types.NUMERIC(strconv.FormatUint(discoveryRemoteSelector, 10)),
			Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Outbound,
			Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
			IsEnabled:           true,
			Capacity:            types.NUMERIC("10000000000"),
			Rate:                types.NUMERIC("10000000000"),
			Tokens:              types.NUMERIC("10000000000"),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "deploy outbound rate limiter")
	rateLimiterRawAddr, err := contracts.RawInstanceAddressFromString(rateLimiterOut.Output.Labels.List()[0])
	require.NoError(t, err, "parse rate limiter raw address")

	deployOut, err := cld_ops.ExecuteOperation(bundle, burn_mint_token_pool.Deploy, env.Chain, contract.DeployInput[burnminttokenpool.BurnMintTokenPool]{
		Template: burnminttokenpool.BurnMintTokenPool{
			InstanceId:   types.TEXT(poolInstanceID),
			CcipOwner:    types.PARTY(party),
			PoolOwner:    types.PARTY(party),
			InstrumentId: instrumentId,
			Decimals:     10,
			RemoteChainConfigs: map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{
				types.NUMERIC(strconv.FormatUint(discoveryRemoteSelector, 10)): {
					RemotePools:        []types.TEXT{"7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3"},
					RemoteTokenAddress: "acdafefb07bff5b120b7afa6ea777cf7eabacc0d",
					InboundCCVs:        []chainlinkapi.RawInstanceAddress{},
					OutboundCCVs:       []chainlinkapi.RawInstanceAddress{},
					FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
					InboundRateLimiter: rateLimiterRawAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: rateLimiterRawAddr.Binding(),
					OutboundRateLimiter:                        rateLimiterRawAddr.Binding(),
				},
			},
			TokenTransferFeeConfigs: map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{},
			PoolReceiveContext:      splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
			TransferTimeout:         burnminttokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy burn mint token pool")
	require.NotEmpty(t, deployOut.Output.Labels.List(), "missing raw pool label in deploy output")

	rawPoolAddr, err := contracts.RawInstanceAddressFromString(deployOut.Output.Labels.List()[0])
	require.NoError(t, err, "parse raw pool label")
	poolInstanceAddress := rawPoolAddr.InstanceAddress()

	edsPort := runDiscoveryEDS(t, env, party)
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "failed to create Token Pool API client")

	msg := oapiCommon.Message{
		DestinationChainSelector: strconv.FormatUint(discoveryRemoteSelector, 10),
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: "1.0",
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(instrumentId.Admin),
				Id:    string(instrumentId.Id),
			},
		},
	}
	assertPoolDiscovered(t, tokenPoolAPIClient, poolInstanceAddress, msg)
}

// TestPoolDiscoveryRegistry_LockRelease mirrors TestPoolDiscoveryRegistry_BurnMint for LockReleaseTokenPool.
func TestPoolDiscoveryRegistry_LockRelease(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t)
	participant := env.Chain.Participants[0]
	party := participant.PartyID

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	runtimeDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
	require.NoError(t, err)
	coreDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	// This single participant plays both roles a real deployment splits: the pool deployer and the
	// host of EDS's registry-observer party. The latter must have both pool DARs vetted to be an
	// informee on either pool type, which is what PoolDiscoveryService subscribes to.
	otherPoolDar, err := contracts.GetDar(contracts.CCIPBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{runtimeDar, coreDar, poolDar, otherPoolDar}, participant)
	require.NoError(t, err)

	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("LNR"),
	}

	poolInstanceID := "test-pool-lnr-discovery"
	rateLimiterOut, err := cld_ops.ExecuteOperation(bundle, rate_limiter.DeployOutbound, env.Chain, contract.DeployInput[ratelimiter.RateLimiter]{
		OwnerParty: types.PARTY(party),
		Template: ratelimiter.RateLimiter{
			PoolInstanceId:      types.TEXT(poolInstanceID),
			PoolOwner:           types.PARTY(party),
			RemoteChainSelector: types.NUMERIC(strconv.FormatUint(discoveryRemoteSelector, 10)),
			Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Outbound,
			Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
			IsEnabled:           true,
			Capacity:            types.NUMERIC("10000000000"),
			Rate:                types.NUMERIC("10000000000"),
			Tokens:              types.NUMERIC("10000000000"),
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
	})
	require.NoError(t, err, "deploy outbound rate limiter")
	rateLimiterRawAddr, err := contracts.RawInstanceAddressFromString(rateLimiterOut.Output.Labels.List()[0])
	require.NoError(t, err, "parse rate limiter raw address")

	deployOut, err := cld_ops.ExecuteOperation(bundle, lock_release_token_pool.Deploy, env.Chain, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
		Template: lockreleasetokenpool.LockReleaseTokenPool{
			InstanceId:   types.TEXT(poolInstanceID),
			CcipOwner:    types.PARTY(party),
			PoolOwner:    types.PARTY(party),
			InstrumentId: instrumentId,
			Decimals:     10,
			RemoteChainConfigs: map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{
				types.NUMERIC(strconv.FormatUint(discoveryRemoteSelector, 10)): {
					RemotePools:        []types.TEXT{"7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3"},
					RemoteTokenAddress: "acdafefb07bff5b120b7afa6ea777cf7eabacc0d",
					InboundCCVs:        []chainlinkapi.RawInstanceAddress{},
					OutboundCCVs:       []chainlinkapi.RawInstanceAddress{},
					FinalityConfig:     ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
					InboundRateLimiter: rateLimiterRawAddr.Binding(),
					InboundCustomBlockConfirmationsRateLimiter: rateLimiterRawAddr.Binding(),
					OutboundRateLimiter:                        rateLimiterRawAddr.Binding(),
				},
			},
			TokenTransferFeeConfigs: map[types.NUMERIC]lockreleasetokenpool.TokenTransferFeeConfig{},
			PoolReceiveContext:      splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
			TransferTimeout:         lockreleasetokenpool.TransferTimeout{Indefinite: new(types.UNIT)},
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy lock release token pool")
	require.NotEmpty(t, deployOut.Output.Labels.List(), "missing raw pool label in deploy output")

	rawPoolAddr, err := contracts.RawInstanceAddressFromString(deployOut.Output.Labels.List()[0])
	require.NoError(t, err, "parse raw pool label")
	poolInstanceAddress := rawPoolAddr.InstanceAddress()

	edsPort := runDiscoveryEDS(t, env, party)
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(fmt.Sprintf("http://localhost:%d", edsPort))
	require.NoError(t, err, "failed to create Token Pool API client")

	msg := oapiCommon.Message{
		DestinationChainSelector: strconv.FormatUint(discoveryRemoteSelector, 10),
		TokenTransfer: &oapiCommon.TokenTransfer{
			Amount: "1.0",
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(instrumentId.Admin),
				Id:    string(instrumentId.Id),
			},
		},
	}
	assertPoolDiscovered(t, tokenPoolAPIClient, poolInstanceAddress, msg)
}

// runDiscoveryEDS runs EDS with TokenPool discovery enabled and no statically configured pools,
// naming registryObserverParty as its registry observer, and returns the port it's listening on.
func runDiscoveryEDS(t *testing.T, env testhelpers.TestEnvironment, registryObserverParty string) int {
	t.Helper()

	participant := env.Chain.Participants[0]
	edsToken, _ := participant.TokenSource.Token()
	edsPort := freeport.GetOne(t)

	go func() {
		log.Info().Msg("Running EDS...")
		err := service.RunEDS(t.Context(), log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel), &config.Config{
			ChainSelector: strconv.FormatUint(env.Chain.ChainSelector(), 10),
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: uint16(edsPort),
			},
			Node: config.NodeConfig{
				URL: participant.Endpoints.GRPCLedgerAPIURL,
				AuthConfig: commonconfig.AuthConfig{
					Type:   commonconfig.AuthTypeInsecureStatic,
					UserID: participant.UserID,
					JWT:    edsToken.AccessToken,
				},
				MaxRetries: 0,
			},
			TokenPoolAPIConfig: config.TokenPoolAPIConfig{
				Enabled:    true,
				TokenPools: map[string]config.TokenPool{},
			},
			RegistryAPIConfig: config.RegistryAPIConfig{
				Enabled: true,
				PartyID: registryObserverParty,
			},
		})
		log.Info().Err(err).Msg("EDS terminated")
		if !errors.Is(err, context.Canceled) {
			log.Error().Err(err).Msg("EDS server exited with error")
			t.Fail()
		}
	}()

	// wait for EDS to start up
	time.Sleep(1 * time.Second)

	return edsPort
}

// assertPoolDiscovered confirms the pool at poolInstanceAddress is unknown to the TokenPool API
// before discovery, then waits past PoolDiscoveryService's 5s poll ticker for a real send
// disclosure fetch against msg to succeed - proving the pool is not just registered, but actually
// servable (its own contract, and its outbound rate limiter, are both visible to EDS).
func assertPoolDiscovered(t *testing.T, tokenPoolAPIClient oapiTokenPool.ClientWithResponsesInterface, poolInstanceAddress contracts.InstanceAddress, msg oapiCommon.Message) {
	t.Helper()

	ctx := t.Context()

	resp, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, poolInstanceAddress.String(), oapiTokenPool.TokenPoolSendRequest{Message: msg})
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode(), "pool should not be known before discovery: %s", string(resp.Body))

	require.Eventually(t, func() bool {
		r, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, poolInstanceAddress.String(), oapiTokenPool.TokenPoolSendRequest{Message: msg})
		return err == nil && r.StatusCode() == 200
	}, 20*time.Second, 500*time.Millisecond, "pool should have been discovered and able to serve a real send disclosure")

	sendResp, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, poolInstanceAddress.String(), oapiTokenPool.TokenPoolSendRequest{Message: msg})
	require.NoError(t, err)
	require.Equal(t, 200, sendResp.StatusCode(), "expected disclosure fetch to succeed once discovered: %s", string(sendResp.Body))
	require.NotNil(t, sendResp.JSON200)
	require.Equal(t, poolInstanceAddress.Hex(), sendResp.JSON200.InstanceAddress)
	require.NotEmpty(t, sendResp.JSON200.ContractId)
	require.NotEmpty(t, sendResp.JSON200.DisclosedContracts, "expected the pool and its rate limiter to be disclosed")
}
