package changesets

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestUpdateCantonTokenPoolFinalityConfig wires one Sepolia lane with WaitForFinality on each Canton
// pool type, then updates the lane's finality config to BlockDepth 1 via the changeset and checks
// that only the finality config changed on the pool.
func TestUpdateCantonTokenPoolFinalityConfig(t *testing.T) {
	t.Parallel()

	selector := chainsel.CANTON_LOCALNET.Selector
	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	bc, err := cantonProvider.NewCTFChainProvider(t, selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	party := types.PARTY(participant.PartyID)

	factoryDar, err := contracts.GetDar(contracts.CCIPFactoryV2, contracts.DevVersion)
	require.NoError(t, err)
	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars:               []*participantv30.UploadDarRequest_UploadDarData{{Bytes: factoryDar}},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "upload factory DAR")

	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())
	ds := datastore.NewMemoryDataStore()

	factoryQualifier := dsutils.QualifierCCIP
	factoryRef, err := cld_ops.ExecuteOperation(bundle, factoryops.Deploy, *cantonChain, contract.DeployInput[factorybindings.CCIPFactory]{
		Qualifier: &factoryQualifier,
		Template: factorybindings.CCIPFactory{
			Owner:             party,
			McmsParty:         party,
			UsedInstanceIds:   map[types.TEXT]types.BOOL{},
			DeployedContracts: map[types.TEXT]types.CONTRACT_ID{},
		},
		OwnerParty: party,
	})
	require.NoError(t, err, "deploy factory")
	require.NoError(t, ds.AddressRefStore.Add(factoryRef.Output))

	tarRef, err := cld_ops.ExecuteOperation(bundle, token_admin_registry.Deploy, *cantonChain, contract.DeployInput[core.TokenAdminRegistry]{
		Template:   core.TokenAdminRegistry{CcipOwner: party},
		OwnerParty: party,
	})
	require.NoError(t, err, "deploy TAR")
	require.NoError(t, ds.AddressRefStore.Add(tarRef.Output))

	// The pool only records RMNRemote and FeeQuoter addresses; neither deploy nor lane setup reads them.
	for _, dep := range []struct {
		instanceID   contracts.InstanceID
		contractType datastore.ContractType
		version      *semver.Version
	}{
		{"rmn-finality-test", datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version},
		{"fq-finality-test", datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version},
	} {
		raw := dep.instanceID.RawInstanceAddress(party)
		require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
			ChainSelector: selector,
			Type:          dep.contractType,
			Version:       dep.version,
			Address:       raw.InstanceAddress().Hex(),
			Labels:        datastore.NewLabelSet(raw.String()),
		}))
	}

	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        ds.Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}
	rateLimit := &tokenadapters.RateLimiterConfigFloatInput{IsEnabled: true, Capacity: 100, Rate: 1}

	//nolint:paralleltest // subtests share one Canton network; parallel networks exhaust Docker
	for _, poolType := range []string{"BurnMintTokenPool", "LockReleaseTokenPool", "RegistryBurnMintTokenPool"} {
		t.Run(poolType, func(t *testing.T) {
			// A fresh operations bundle per changeset run, so the reporter does not return a
			// previous run's result.
			applyFinalityUpdate := func(cfg UpdateCantonTokenPoolFinalityConfigConfig) (cldf.ChangesetOutput, error) {
				env.OperationsBundle = cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

				return (UpdateCantonTokenPoolFinalityConfig{}).Apply(*env, CantonCSDeps[UpdateCantonTokenPoolFinalityConfigConfig]{
					ChainSelector: selector,
					Participant:   0,
					Config:        cfg,
				})
			}

			deployOut, err := cld_ops.ExecuteSequence(bundle, sequences.DeployTokenPoolForToken, env.BlockChains, tokenadapters.DeployTokenPoolInput{
				ChainSelector:      selector,
				PoolType:           poolType,
				TokenPoolVersion:   semver.MustParse("2.0.0"),
				TokenPoolQualifier: strings.ToLower(strings.ReplaceAll(poolType, "TokenPool", "")) + "-finality",
				TokenRef: &datastore.AddressRef{
					Address: "link-token-finality-test",
					Labels: datastore.NewLabelSet(
						"instrument-admin:"+string(party),
						"instrument-id:LINK-"+poolType,
						"ccip-owner:"+string(party),
					),
				},
				ExistingDataStore: ds.Seal(),
			})
			require.NoError(t, err, "deploy pool")
			for _, ref := range deployOut.Output.Addresses {
				require.NoError(t, ds.AddressRefStore.Add(ref))
			}
			poolRef := deployOut.Output.Addresses[0]
			require.Equal(t, datastore.ContractType(poolType), poolRef.Type)

			_, err = cld_ops.ExecuteSequence(bundle, sequences.ConfigureTokenForTransfers, env.BlockChains, tokenadapters.ConfigureTokenForTransfersInput{
				ChainSelector:    selector,
				TokenPoolAddress: poolRef.Address,
				PoolType:         poolType,
				RemoteChains: map[uint64]tokenadapters.RemoteChainConfig[[]byte, string]{
					remoteSelector: {
						RemotePool:                common.HexToAddress("0x7e3febbdaf80e7e96c1ae107508ec3fafc36d7f3").Bytes(),
						RemoteToken:               common.HexToAddress("0xacdafefb07bff5b120b7afa6ea777cf7eabacc0d").Bytes(),
						RemoteDecimals:            18,
						InboundRateLimiterConfig:  rateLimit,
						OutboundRateLimiterConfig: rateLimit,
					},
				},
				AllowedFinalityConfig: finality.Config{WaitForFinality: true},
				ExistingDataStore:     ds.Seal(),
			})
			require.NoError(t, err, "configure pool lane")

			laneBefore := loadCantonPoolLane(t, participant, poolRef.Type, poolRef.Address, remoteSelector)
			require.Equal(t, "WaitForFinality", laneBefore.Finality.GetVariantTag())

			out, err := applyFinalityUpdate(UpdateCantonTokenPoolFinalityConfigConfig{
				TokenPoolRef:          poolRef,
				RemoteChainSelectors:  []uint64{remoteSelector},
				AllowedFinalityConfig: finality.Config{BlockDepth: 1},
			})
			require.NoError(t, err, "update lane finality config")
			require.Empty(t, out.MCMSTimelockProposals, "direct execution must not produce MCMS proposals")

			laneAfter := loadCantonPoolLane(t, participant, poolRef.Type, poolRef.Address, remoteSelector)
			require.Equal(t, "BlockDepth", laneAfter.Finality.GetVariantTag())
			require.NotNil(t, laneAfter.Finality.BlockDepth)
			require.Equal(t, types.INT64(1), *laneAfter.Finality.BlockDepth)
			// Everything except the finality config must be preserved.
			require.Equal(t, laneBefore.RemotePools, laneAfter.RemotePools)
			require.Equal(t, laneBefore.RemoteTokenAddress, laneAfter.RemoteTokenAddress)
			require.Equal(t, laneBefore.InboundCCVs, laneAfter.InboundCCVs)
			require.Equal(t, laneBefore.OutboundCCVs, laneAfter.OutboundCCVs)
			require.Equal(t, laneBefore.InboundRateLimiter, laneAfter.InboundRateLimiter)
			require.Equal(t, laneBefore.InboundCustomRateLimiter, laneAfter.InboundCustomRateLimiter)
			require.Equal(t, laneBefore.OutboundRateLimiter, laneAfter.OutboundRateLimiter)

			// A rerun for a lane already at the target finality is a no-op.
			rerunOut, err := applyFinalityUpdate(UpdateCantonTokenPoolFinalityConfigConfig{
				TokenPoolRef:          poolRef,
				RemoteChainSelectors:  []uint64{remoteSelector},
				AllowedFinalityConfig: finality.Config{BlockDepth: 1},
			})
			require.NoError(t, err, "rerun finality config update")
			require.Empty(t, rerunOut.MCMSTimelockProposals)

			// A lane that is not configured on the pool is rejected.
			_, err = applyFinalityUpdate(UpdateCantonTokenPoolFinalityConfigConfig{
				TokenPoolRef:          poolRef,
				RemoteChainSelectors:  []uint64{12345},
				AllowedFinalityConfig: finality.Config{BlockDepth: 1},
			})
			require.ErrorContains(t, err, "remote chain 12345 is not configured on token pool")

			// A config with several finality modes is rejected in preconditions.
			_, err = applyFinalityUpdate(UpdateCantonTokenPoolFinalityConfigConfig{
				TokenPoolRef:          poolRef,
				RemoteChainSelectors:  []uint64{remoteSelector},
				AllowedFinalityConfig: finality.Config{WaitForFinality: true, WaitForSafe: true},
			})
			require.ErrorContains(t, err, "exactly one finality mode must be set")

			if poolType != "BurnMintTokenPool" {
				return
			}

			// A lane whose custom finality inbound rate limiter is unset cannot move to a
			// faster-than-finality config. The DAML only validates that limiter address when the
			// finality config requires it, so it can be empty while the lane waits for finality.
			lane := loadCantonPoolLane(t, participant, poolRef.Type, poolRef.Address, remoteSelector)
			_, err = cld_ops.ExecuteOperation(bundle, burn_mint_token_pool.ApplyChainUpdates, *cantonChain, contract.ChoiceInput[burnminttokenpool.ApplyChainUpdates]{
				InstanceAddress: contracts.HexToInstanceAddress(poolRef.Address),
				Args: burnminttokenpool.ApplyChainUpdates{
					RemoteChainSelectorsToRemove: []types.NUMERIC{lane.RemoteChainSelectorKey},
					ChainsToAdd: []burnminttokenpool.ChainUpdate{{
						RemoteChainSelector: lane.RemoteChainSelectorKey,
						RemotePools:         lane.RemotePools,
						RemoteTokenAddress:  lane.RemoteTokenAddress,
						InboundCCVs:         lane.InboundCCVs,
						OutboundCCVs:        lane.OutboundCCVs,
						FinalityConfig:      ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
						InboundRateLimiter:  lane.InboundRateLimiter,
						InboundCustomBlockConfirmationsRateLimiter: chainlinkapi.RawInstanceAddress{},
						OutboundRateLimiter:                        lane.OutboundRateLimiter,
					}},
				},
			})
			require.NoError(t, err, "rebuild lane without custom finality rate limiter")

			_, err = applyFinalityUpdate(UpdateCantonTokenPoolFinalityConfigConfig{
				TokenPoolRef:          poolRef,
				RemoteChainSelectors:  []uint64{remoteSelector},
				AllowedFinalityConfig: finality.Config{BlockDepth: 1},
			})
			require.ErrorContains(t, err, "custom finality inbound rate limiter is not configured")
		})
	}
}

type cantonLaneSnapshot struct {
	RemoteChainSelectorKey   types.NUMERIC
	RemotePools              []types.TEXT
	RemoteTokenAddress       types.TEXT
	InboundCCVs              []chainlinkapi.RawInstanceAddress
	OutboundCCVs             []chainlinkapi.RawInstanceAddress
	Finality                 ccipcodec.FinalityConfig
	InboundRateLimiter       chainlinkapi.RawInstanceAddress
	InboundCustomRateLimiter chainlinkapi.RawInstanceAddress
	OutboundRateLimiter      chainlinkapi.RawInstanceAddress
}

// loadCantonPoolLane queries the pool from the ledger and returns the lane config for remoteSelector.
func loadCantonPoolLane(
	t *testing.T,
	participant canton.Participant,
	poolType datastore.ContractType,
	poolAddressHex string,
	remoteSelector uint64,
) cantonLaneSnapshot {
	t.Helper()

	ctx := t.Context()
	queryParties := contract.LedgerQueryParties(participant)
	poolAddress := contracts.HexToInstanceAddress(poolAddressHex)
	selectorStr := strconv.FormatUint(remoteSelector, 10)

	snapshot := func(selectorKey types.NUMERIC, cfg any) cantonLaneSnapshot {
		switch lane := cfg.(type) {
		case lockreleasetokenpool.RemoteChainConfig:
			return newCantonLaneSnapshot(selectorKey, lane.RemotePools, lane.RemoteTokenAddress, lane.InboundCCVs, lane.OutboundCCVs, lane.FinalityConfig, lane.InboundRateLimiter, lane.InboundCustomBlockConfirmationsRateLimiter, lane.OutboundRateLimiter)
		case burnminttokenpool.RemoteChainConfig:
			return newCantonLaneSnapshot(selectorKey, lane.RemotePools, lane.RemoteTokenAddress, lane.InboundCCVs, lane.OutboundCCVs, lane.FinalityConfig, lane.InboundRateLimiter, lane.InboundCustomBlockConfirmationsRateLimiter, lane.OutboundRateLimiter)
		case registrybnm.RemoteChainConfig:
			return newCantonLaneSnapshot(selectorKey, lane.RemotePools, lane.RemoteTokenAddress, lane.InboundCCVs, lane.OutboundCCVs, lane.FinalityConfig, lane.InboundRateLimiter, lane.InboundCustomBlockConfirmationsRateLimiter, lane.OutboundRateLimiter)
		default:
			t.Fatalf("unsupported lane config type %T", cfg)
		}

		return cantonLaneSnapshot{}
	}

	switch poolType {
	case datastore.ContractType("LockReleaseTokenPool"):
		activePool, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, queryParties,
			lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), poolAddress)
		require.NoError(t, err)
		pool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
		require.NoError(t, err)
		for key, cfg := range pool.RemoteChainConfigs {
			if whole, _, _ := strings.Cut(string(key), "."); whole == selectorStr {
				return snapshot(key, cfg)
			}
		}
	case datastore.ContractType("BurnMintTokenPool"):
		activePool, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, queryParties,
			burnminttokenpool.BurnMintTokenPool{}.GetTemplateID(), poolAddress)
		require.NoError(t, err)
		pool, err := bindings.UnmarshalCreatedEvent[burnminttokenpool.BurnMintTokenPool](activePool.GetCreatedEvent())
		require.NoError(t, err)
		for key, cfg := range pool.RemoteChainConfigs {
			if whole, _, _ := strings.Cut(string(key), "."); whole == selectorStr {
				return snapshot(key, cfg)
			}
		}
	case datastore.ContractType("RegistryBurnMintTokenPool"):
		activePool, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, queryParties,
			registrybnm.BurnMintTokenPool{}.GetTemplateID(), poolAddress)
		require.NoError(t, err)
		pool, err := bindings.UnmarshalCreatedEvent[registrybnm.BurnMintTokenPool](activePool.GetCreatedEvent())
		require.NoError(t, err)
		for key, cfg := range pool.RemoteChainConfigs {
			if whole, _, _ := strings.Cut(string(key), "."); whole == selectorStr {
				return snapshot(key, cfg)
			}
		}
	default:
		t.Fatalf("unsupported pool type %q", poolType)
	}

	t.Fatalf("lane %d not found on pool %s", remoteSelector, poolAddressHex)

	return cantonLaneSnapshot{}
}

func newCantonLaneSnapshot(
	selectorKey types.NUMERIC,
	remotePools []types.TEXT,
	remoteTokenAddress types.TEXT,
	inboundCCVs, outboundCCVs []chainlinkapi.RawInstanceAddress,
	finality ccipcodec.FinalityConfig,
	inboundRateLimiter, inboundCustomRateLimiter, outboundRateLimiter chainlinkapi.RawInstanceAddress,
) cantonLaneSnapshot {
	return cantonLaneSnapshot{
		RemoteChainSelectorKey:   selectorKey,
		RemotePools:              remotePools,
		RemoteTokenAddress:       remoteTokenAddress,
		InboundCCVs:              inboundCCVs,
		OutboundCCVs:             outboundCCVs,
		Finality:                 finality,
		InboundRateLimiter:       inboundRateLimiter,
		InboundCustomRateLimiter: inboundCustomRateLimiter,
		OutboundRateLimiter:      outboundRateLimiter,
	}
}
