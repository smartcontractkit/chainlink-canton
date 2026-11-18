package sequences

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"testing"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	registryratelimiter "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// TestRegistryBurnMintTokenPool_proposalMode runs the registry pool deploy and lane setup with a participant
// that can only read as ccipOwner, so every exercise becomes an MCMS transaction. The ledger is empty: the
// pool is not active yet, so lane setup falls back to the datastore as it does while a deploy proposal is pending.
func TestRegistryBurnMintTokenPool_proposalMode(t *testing.T) {
	t.Parallel()

	selector := chainsel.CANTON_LOCALNET.Selector
	ccipOwner := types.PARTY("ccip-owner::1220aa")
	observers := []types.PARTY{"eds-a::1220bb", "eds-b::1220cc"}
	chains := chain.NewBlockChainsFromSlice([]chain.BlockChain{cldfcanton.Chain{
		Selector: selector,
		Participants: []cldfcanton.Participant{{
			PartyID:        "operator::1220dd",
			ReadAsPartyIDs: []string{string(ccipOwner)},
			LedgerServices: cldfcanton.LedgerServiceClients{State: emptyLedgerState(t)},
		}},
	}})
	bundle := operations.NewBundle(t.Context, logger.Test(t), operations.NewMemoryReporter())

	ds := datastore.NewMemoryDataStore()
	addRef := func(instanceID contracts.InstanceID, contractType datastore.ContractType, version *semver.Version, qualifier string) contracts.RawInstanceAddress {
		raw := instanceID.RawInstanceAddress(ccipOwner)
		require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
			ChainSelector: selector,
			Type:          contractType,
			Version:       version,
			Qualifier:     qualifier,
			Address:       raw.InstanceAddress().Hex(),
			Labels:        datastore.NewLabelSet(raw.String()),
		}))

		return raw
	}
	factoryRaw := addRef("ccipfactory", datastore.ContractType(factoryops.ContractType), factoryops.Version, dsutils.QualifierCCIP)
	tarRaw := addRef("tar", datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version, "")
	rmnRaw := addRef("rmn", datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, "")
	fqRaw := addRef("fq", datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, "")

	tokenLabels := datastore.NewLabelSet(
		"instrument-admin:"+string(ccipOwner), "instrument-id:LINK", "ccip-owner:"+string(ccipOwner),
		"observer:"+string(observers[0]), "observer:"+string(observers[1]),
	)
	deployOut, err := operations.ExecuteSequence(bundle, DeployTokenPoolForToken, chains, tokenadapters.DeployTokenPoolInput{
		ChainSelector:      selector,
		PoolType:           registryBurnMintPoolType.String(),
		TokenPoolVersion:   semver.MustParse("2.0.0"),
		TokenPoolQualifier: "link",
		TokenRef:           &datastore.AddressRef{Address: "link-token", Labels: tokenLabels},
		ExistingDataStore:  ds.Seal(),
	})
	require.NoError(t, err)

	poolRaw := contracts.InstanceID("registryburnminttokenpool-link").RawInstanceAddress(ccipOwner)
	deployTx := findProposalTx(t, deployOut.Output.BatchOps, "DeployRegistryBurnMintTokenPool")
	require.Equal(t, factoryRaw.InstanceAddress().Hex(), deployTx.To)
	var poolParams factorybindings.DeployRegistryBurnMintTokenPoolParams
	require.NoError(t, poolParams.UnmarshalHex(hex.EncodeToString(deployTx.Data)))
	require.Equal(t, factorybindings.DeployRegistryBurnMintTokenPoolParams{
		InstanceId:         types.TEXT(poolRaw.InstanceID()),
		InstrumentId:       splice_api_token_holding_v1.InstrumentId{Admin: ccipOwner, Id: "LINK"},
		Decimals:           10,
		Observers:          observers,
		TokenAdminRegistry: tarRaw.Binding(),
		FeeQuoter:          fqRaw.Binding(),
		RmnRemote:          rmnRaw.Binding(),
		PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
		TransferTimeout:    registrybnm.TransferTimeout{RelativeHours: new(types.INT64(24))},
	}, poolParams)

	for _, ref := range deployOut.Output.Addresses {
		require.NoError(t, ds.AddressRefStore.Add(ref))
	}
	tokenRef := deployOut.Output.Addresses[2]
	require.Equal(t, datastore.ContractType("Token"), tokenRef.Type)
	// The configure input has no observer labels: the limiters must still get the pool's observers.
	configureTokenRef := tokenRef
	configureTokenRef.Labels = datastore.NewLabelSet(
		"instrument-admin:"+string(ccipOwner), "instrument-id:LINK", "ccip-owner:"+string(ccipOwner),
	)

	remoteSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	remoteSelectorKey := types.NUMERIC("16015286601757825753")
	rateLimit := &tokenadapters.RateLimiterConfigFloatInput{IsEnabled: true, Capacity: 100, Rate: 1}
	configureOut, err := operations.ExecuteSequence(bundle, ConfigureTokenForTransfers, chains, tokenadapters.ConfigureTokenForTransfersInput{
		ChainSelector:    selector,
		TokenPoolAddress: poolRaw.InstanceAddress().Hex(),
		PoolType:         registryBurnMintPoolType.String(),
		TokenRef:         configureTokenRef,
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
	require.NoError(t, err)
	require.Len(t, configureOut.Output.BatchOps, 1)
	txs := configureOut.Output.BatchOps[0].Transactions
	require.Len(t, txs, 4, "three rate limiters, then ApplyChainUpdates")
	require.Len(t, configureOut.Output.Addresses, 3)

	limiterRaws := make([]contracts.RawInstanceAddress, 3)
	for i, want := range []struct {
		direction      registryratelimiter.RateLimitDirection
		mode           registryratelimiter.RateLimitMode
		capacity, rate types.NUMERIC
	}{
		// 100 tokens and 1 token per second at 10 decimals; inbound adds 10% headroom.
		{registryratelimiter.RateLimitDirectionRateLimitDirection_Outbound, registryratelimiter.RateLimitModeRateLimitMode_DefaultFinality, "1000000000000", "10000000000"},
		{registryratelimiter.RateLimitDirectionRateLimitDirection_Inbound, registryratelimiter.RateLimitModeRateLimitMode_DefaultFinality, "1100000000000", "11000000000"},
		{registryratelimiter.RateLimitDirectionRateLimitDirection_Inbound, registryratelimiter.RateLimitModeRateLimitMode_CustomFinality, "1100000000000", "11000000000"},
	} {
		limiterRaws[i], err = dsutils.GetRawInstanceAddressFromAddressRef(configureOut.Output.Addresses[i])
		require.NoError(t, err)
		requireFunctionName(t, txs[i], "DeployRegistryRateLimiter")
		require.Equal(t, factoryRaw.InstanceAddress().Hex(), txs[i].To)
		var limiterParams factorybindings.DeployRegistryRateLimiterParams
		require.NoError(t, limiterParams.UnmarshalHex(hex.EncodeToString(txs[i].Data)))
		require.Equal(t, factorybindings.DeployRegistryRateLimiterParams{
			InstanceId:          types.TEXT(limiterRaws[i].InstanceID()),
			PoolInstanceId:      types.TEXT(poolRaw.InstanceID()),
			RemoteChainSelector: remoteSelectorKey,
			Direction:           want.direction,
			Mode:                want.mode,
			IsEnabled:           true,
			Capacity:            want.capacity,
			Rate:                want.rate,
			Observers:           observers,
		}, limiterParams)
	}

	requireFunctionName(t, txs[3], "ApplyChainUpdates")
	require.Equal(t, poolRaw.InstanceAddress().Hex(), txs[3].To)
	var applyParams registrybnm.ApplyChainUpdatesParams
	require.NoError(t, applyParams.UnmarshalHex(hex.EncodeToString(txs[3].Data)))
	require.Equal(t, registrybnm.ApplyChainUpdatesParams{
		RemoteChainSelectorsToRemove: []types.NUMERIC{},
		ChainsToAdd: []registrybnm.ChainUpdate{{
			RemoteChainSelector: remoteSelectorKey,
			RemotePools:         []types.TEXT{"0000000000000000000000007e3febbdaf80e7e96c1ae107508ec3fafc36d7f3"},
			RemoteTokenAddress:  "000000000000000000000000acdafefb07bff5b120b7afa6ea777cf7eabacc0d",
			InboundCCVs:         []chainlinkapi.RawInstanceAddress{},
			OutboundCCVs:        []chainlinkapi.RawInstanceAddress{},
			FinalityConfig:      ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
			OutboundRateLimiter: limiterRaws[0].Binding(),
			InboundRateLimiter:  limiterRaws[1].Binding(),
			InboundCustomBlockConfirmationsRateLimiter: limiterRaws[2].Binding(),
		}},
	}, applyParams)
}

func findProposalTx(t *testing.T, batches []mcms_types.BatchOperation, functionName string) mcms_types.Transaction {
	t.Helper()

	var found []mcms_types.Transaction
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			var af cantonsdk.AdditionalFields
			require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
			if af.FunctionName == functionName {
				found = append(found, tx)
			}
		}
	}
	require.Len(t, found, 1, "proposal transactions for %s", functionName)

	return found[0]
}

func requireFunctionName(t *testing.T, tx mcms_types.Transaction, want string) {
	t.Helper()

	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	require.Equal(t, want, af.FunctionName)
}

// emptyLedgerState answers every active-contract query with no contracts.
func emptyLedgerState(t *testing.T) *mocks.MockStateServiceClient {
	t.Helper()

	state := mocks.NewMockStateServiceClient(t)
	state.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil).Maybe()
	state.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(context.Context, *apiv2.GetActiveContractsRequest, ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse], error) {
			stream := mocks.NewMockServerStreamingClient[apiv2.GetActiveContractsResponse](t)
			stream.EXPECT().Recv().Return(nil, io.EOF).Once()
			stream.EXPECT().CloseSend().Return(nil).Once()

			return stream, nil
		}).Maybe()

	return state
}
