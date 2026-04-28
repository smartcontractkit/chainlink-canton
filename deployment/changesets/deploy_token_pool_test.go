package changesets

import (
	"context"
	"encoding/hex"
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestDeployTokenPool deploys TAR, then runs DeployTokenPool (deploy pool + register with TAR in one changeset),
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

	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	tarDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
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
	tarAddrRef, err := cld_ops.ExecuteOperation(bundle, token_admin_registry.Deploy, *cantonChain, contract.DeployInput[tokenadminregistry.TokenAdminRegistry]{
		Template: tokenadminregistry.TokenAdminRegistry{
			Owner:      types.PARTY(party),
			InstanceId: "",
			EntryCount: 0,
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy TAR")
	require.NotEmpty(t, tarAddrRef.Output.Address, "TAR address")

	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("AMT"),
	}

	// Deploy pool and register with TAR in one changeset
	_, err = (DeployTokenPool{}).Apply(*env, CantonCSDeps[DeployTokenPoolConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployTokenPoolConfig{
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
) (*tokenadminregistry.TokenConfig, error) {
	activeContract, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		stateClient,
		ccipParty,
		tokenadminregistry.TokenConfig{}.GetTemplateID(),
		contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).RawInstanceAddress(types.PARTY(ccipParty)).InstanceAddress(),
	)
	if err != nil {
		return nil, err
	}

	return bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenConfig](activeContract.GetCreatedEvent())
}
