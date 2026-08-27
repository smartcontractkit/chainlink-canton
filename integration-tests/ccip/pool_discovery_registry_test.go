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
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestPoolDiscoveryRegistry_BurnMint deploys a BurnMintTokenPool naming a single party as both
// its rateLimitAdmin (the pool's only observer field) and EDS's registry-observer party, then
// verifies PoolDiscoveryService picks it up via ACS polling without any statically configured
// TokenPools entry.
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
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{runtimeDar, coreDar, poolDar}, participant)
	require.NoError(t, err)

	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("BNM"),
	}

	deployOut, err := cld_ops.ExecuteOperation(bundle, burn_mint_token_pool.Deploy, env.Chain, contract.DeployInput[burnminttokenpool.BurnMintTokenPool]{
		Template: burnminttokenpool.BurnMintTokenPool{
			CcipOwner:    types.PARTY(party),
			PoolOwner:    types.PARTY(party),
			InstrumentId: instrumentId,
			Decimals:     6,
			// Naming `party` here (the pool's only observer field) as the same party EDS uses as
			// its registry observer is what makes the pool visible to EDS's party-scoped ACS query.
			RateLimitAdmin:          new(types.PARTY(party)),
			RemoteChainConfigs:      map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{},
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

	assertPoolDiscovered(t, tokenPoolAPIClient, poolInstanceAddress)
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
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{runtimeDar, coreDar, poolDar}, participant)
	require.NoError(t, err)

	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("LNR"),
	}

	deployOut, err := cld_ops.ExecuteOperation(bundle, lock_release_token_pool.Deploy, env.Chain, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
		Template: lockreleasetokenpool.LockReleaseTokenPool{
			CcipOwner:    types.PARTY(party),
			PoolOwner:    types.PARTY(party),
			InstrumentId: instrumentId,
			Decimals:     6,
			// Naming `party` here (the pool's only observer field) as the same party EDS uses as
			// its registry observer is what makes the pool visible to EDS's party-scoped ACS query.
			RateLimitAdmin:          new(types.PARTY(party)),
			RemoteChainConfigs:      map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig{},
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

	assertPoolDiscovered(t, tokenPoolAPIClient, poolInstanceAddress)
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
// before discovery, then waits past PoolDiscoveryService's 5s poll ticker for it to become known.
func assertPoolDiscovered(t *testing.T, tokenPoolAPIClient oapiTokenPool.ClientWithResponsesInterface, poolInstanceAddress contracts.InstanceAddress) {
	t.Helper()

	ctx := t.Context()

	resp, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, poolInstanceAddress.String(), oapiTokenPool.TokenPoolSendRequest{})
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode(), "pool should not be known before discovery: %s", string(resp.Body))

	require.Eventually(t, func() bool {
		resp, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, poolInstanceAddress.String(), oapiTokenPool.TokenPoolSendRequest{})
		if err != nil {
			return false
		}

		return resp.StatusCode() != 404
	}, 20*time.Second, 500*time.Millisecond, "pool should have been discovered and registered by PoolDiscoveryService")
}
