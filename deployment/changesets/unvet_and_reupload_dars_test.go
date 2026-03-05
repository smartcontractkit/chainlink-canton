package changesets

import (
	"sync"
	"testing"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestUnvetAndReuploadDARs(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]

	// Upload the initial contract DAR (v1).
	gcV1Dar, err := contracts.GetDar(contracts.GlobalConfig, "1.0.0")
	require.NoError(t, err)
	v1Validated, err := participant.AdminServices.Package.ValidateDar(t.Context(), &participantv30.ValidateDarRequest{
		Data: gcV1Dar,
	})
	require.NoError(t, err)
	v1MainPkgID := v1Validated.GetMainPackageId()
	require.NotEmpty(t, v1MainPkgID)

	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: gcV1Dar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err)

	// Simulate a contract change by using the upgraded DAR (v2).
	gcV2Dar, err := contracts.GetDar(contracts.GlobalConfig, "2.0.0")
	require.NoError(t, err)
	v2Validated, err := participant.AdminServices.Package.ValidateDar(t.Context(), &participantv30.ValidateDarRequest{
		Data: gcV2Dar,
	})
	require.NoError(t, err)
	v2MainPkgID := v2Validated.GetMainPackageId()
	require.NotEmpty(t, v2MainPkgID)
	require.NotEqual(t, v1MainPkgID, v2MainPkgID)

	reporter := cldops.NewMemoryReporter()
	bundle := cldops.NewBundle(t.Context, logger.Test(t), reporter)
	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	cs := UnvetAndReuploadDARs{}
	_, err = cs.Apply(*env, CantonCSDeps[UnvetAndReuploadDARsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: UnvetAndReuploadDARsConfig{
			DARs:                  [][]byte{gcV2Dar},
			MainPackageIDsToUnvet: []string{v1MainPkgID},
			SynchronizeVetting:    true,
		},
	})
	require.NoError(t, err)

	listResp, err := participant.AdminServices.Package.ListDars(t.Context(), &participantv30.ListDarsRequest{})
	require.NoError(t, err)

	foundV2 := false
	for _, dar := range listResp.GetDars() {
		if dar.GetMain() == v2MainPkgID {
			foundV2 = true
			break
		}
	}
	require.True(t, foundV2, "expected upgraded v2 DAR to be present after reupload")
}
