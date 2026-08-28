package changesets

import (
	"context"
	"encoding/hex"
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
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
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestDeployTokenPool deploys TAR, then runs DeployLockReleaseTokenPool (deploy pool + register with TAR in one changeset),
// then verifies TAR config by querying GetTokenConfig and asserting tokenPoolOwner is set.
func TestDeployTokenPool(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	party := participant.PartyID

	commonDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
	require.NoError(t, err)
	tarDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: commonDar},
			{Bytes: tarDar},
			{Bytes: poolDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload dar files")

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)

	// Deploy TAR so we have an instance address for register-with-TAR
	tarAddrRef, err := cld_ops.ExecuteOperation(bundle, token_admin_registry.Deploy, *cantonChain, contract.DeployInput[core.TokenAdminRegistry]{
		Template: core.TokenAdminRegistry{
			CcipOwner:  types.PARTY(party),
			InstanceId: "",
			EntryCount: 0,
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy TAR")
	require.NotEmpty(t, tarAddrRef.Output.Address, "TAR address")

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.AddressRefStore.Add(tarAddrRef.Output))

	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        ds.Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("AMT"),
	}

	// Deploy pool and register with TAR in one changeset
	_, err = (DeployLockReleaseTokenPool{}).Apply(*env, CantonCSDeps[DeployLockReleaseTokenPoolConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployLockReleaseTokenPoolConfig{
			CcipOwner:                         party,
			PoolOwner:                         party,
			InstrumentId:                      instrumentId,
			Decimals:                          6,
			Qualifier:                         "AMT",
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tarAddrRef.Output.Address),
		},
	})
	require.NoError(t, err, "deploy token pool and register with TAR")

	tokenConfig, err := findTokenConfigByInstanceAddress(t.Context(), participant.LedgerServices.State, party, instrumentId)
	require.NoError(t, err, "find TokenConfig contract in ACS")
	require.NotNil(t, tokenConfig.TokenPool, "TokenConfig should have a registered token pool")
	require.Equal(t, types.PARTY(party), tokenConfig.TokenPool.PoolOwner, "TokenConfig tokenPool.poolOwner should match pool owner")
}

func findTokenConfigByInstanceAddress(
	ctx context.Context,
	stateClient apiv2.StateServiceClient,
	ccipParty string,
	instrumentId splice_api_token_holding_v1.InstrumentId,
) (*core.TokenConfig, error) {
	activeContract, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		stateClient,
		[]string{ccipParty},
		core.TokenConfig{}.GetTemplateID(),
		contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).RawInstanceAddress(types.PARTY(ccipParty)).InstanceAddress(),
	)
	if err != nil {
		return nil, err
	}

	return bindings.UnmarshalCreatedEvent[core.TokenConfig](activeContract.GetCreatedEvent())
}
