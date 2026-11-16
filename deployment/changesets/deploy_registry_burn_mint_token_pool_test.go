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
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	registryratelimiter "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestDeployRegistryBurnMintTokenPool deploys a registry burn/mint pool through the CCIPFactory with the
// token-pool sequences, wires one lane, and checks the pool's lane, its rate limiters and its TAR registration.
func TestDeployRegistryBurnMintTokenPool(t *testing.T) {
	t.Parallel()

	selector := chainsel.CANTON_LOCALNET.Selector
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

	// owner == mcmsParty makes the factory MCMS-controlled while the operator can still act for it directly.
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
		{"rmn-registry-pool-test", datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version},
		{"fq-registry-pool-test", datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version},
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

	instrumentID := splice_api_token_holding_v1.InstrumentId{Admin: party, Id: "LINK"}
	qualifier := "registry-link"
	deployOut, err := cld_ops.ExecuteSequence(bundle, sequences.DeployTokenPoolForToken, chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}), tokenadapters.DeployTokenPoolInput{
		ChainSelector:      selector,
		PoolType:           "RegistryBurnMintTokenPool",
		TokenPoolVersion:   semver.MustParse("2.0.0"),
		TokenPoolQualifier: qualifier,
		TokenRef: &datastore.AddressRef{
			Address: "link-token",
			Labels:  datastore.NewLabelSet("instrument-admin:"+string(party), "instrument-id:LINK"),
		},
		ExistingDataStore: ds.Seal(),
	})
	require.NoError(t, err, "deploy registry pool")
	require.Empty(t, deployOut.Output.BatchOps, "direct execution must not produce MCMS batches")
	for _, ref := range deployOut.Output.Addresses {
		require.NoError(t, ds.AddressRefStore.Add(ref))
	}
	poolRef := deployOut.Output.Addresses[0]
	require.Equal(t, datastore.ContractType("RegistryBurnMintTokenPool"), poolRef.Type)

	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	rateLimit := &tokenadapters.RateLimiterConfigFloatInput{IsEnabled: true, Capacity: 100, Rate: 1}
	configureOut, err := cld_ops.ExecuteSequence(bundle, sequences.ConfigureTokenForTransfers, chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}), tokenadapters.ConfigureTokenForTransfersInput{
		ChainSelector:    selector,
		TokenPoolAddress: poolRef.Address,
		PoolType:         "RegistryBurnMintTokenPool",
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
	require.NoError(t, err, "configure registry pool lane")
	require.Empty(t, configureOut.Output.BatchOps, "direct execution must not produce MCMS batches")
	require.Len(t, configureOut.Output.Addresses, 3, "one lane deploys three rate limiters")

	queryParties := contract.LedgerQueryParties(participant)
	activePool, err := contract.FindActiveContractByInstanceAddress(t.Context(), participant.LedgerServices.State, queryParties,
		registrybnm.BurnMintTokenPool{}.GetTemplateID(), contracts.HexToInstanceAddress(poolRef.Address))
	require.NoError(t, err, "find registry pool")
	pool, err := bindings.UnmarshalCreatedEvent[registrybnm.BurnMintTokenPool](activePool.GetCreatedEvent())
	require.NoError(t, err)
	require.Equal(t, party, pool.PoolOwner)
	require.Equal(t, party, pool.CcipOwner)
	require.Equal(t, instrumentID, pool.InstrumentId)

	require.Len(t, pool.RemoteChainConfigs, 1, "pool must have exactly the configured lane")
	var lane registrybnm.RemoteChainConfig
	for selectorKey, cfg := range pool.RemoteChainConfigs {
		t.Logf("pool lane key %q", selectorKey)
		selectorText, _, _ := strings.Cut(string(selectorKey), ".")
		require.Equal(t, strconv.FormatUint(remoteSelector, 10), selectorText)
		lane = cfg
	}
	limiters := []types.TEXT{lane.InboundRateLimiter.Unpack, lane.OutboundRateLimiter.Unpack, lane.InboundCustomBlockConfirmationsRateLimiter.Unpack}
	for _, address := range limiters {
		raw, err := contracts.RawInstanceAddressFromString(string(address))
		require.NoError(t, err)
		activeLimiter, err := contract.FindActiveContractByInstanceAddress(t.Context(), participant.LedgerServices.State, queryParties,
			registryratelimiter.RateLimiter{}.GetTemplateID(), raw.InstanceAddress())
		require.NoError(t, err, "find registry rate limiter %s", address)
		limiter, err := bindings.UnmarshalCreatedEvent[registryratelimiter.RateLimiter](activeLimiter.GetCreatedEvent())
		require.NoError(t, err)
		require.Equal(t, party, limiter.PoolOwner)
		require.Equal(t, pool.InstanceId, limiter.PoolInstanceId)
	}

	tokenConfig, err := findTokenConfigByInstanceAddress(t.Context(), participant.LedgerServices.State, string(party), instrumentID)
	require.NoError(t, err, "find TokenConfig")
	require.NotNil(t, tokenConfig.TokenPool, "pool must be registered with TAR")
	require.Equal(t, party, tokenConfig.TokenPool.PoolOwner)
	require.Equal(t, pool.InstanceId, tokenConfig.TokenPool.PoolInstanceId)
}
