package integration

import (
	"encoding/binary"
	"strconv"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/util"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
)

func getCantonConfig(t *testing.T, in *ccv.Cfg, cantonSelector string) *ccip.Config {
	for _, v := range in.Verifier {
		if v.ChainFamily == blockchain.FamilyCanton {
			cantonConfig, err := util.OpaqueToConcreteStrict[ccip.Config](v.CantonConfigs)
			require.NoError(t, err)

			return cantonConfig
		}
	}
	require.FailNowf(t, "no canton config found for selector %s", cantonSelector)

	return nil
}

func getCantonGRPCURL(t *testing.T, in *ccv.Cfg, cantonSelector string) string {
	for _, b := range in.Blockchains {
		if b.Type == blockchain.TypeCanton {
			details, err := chainsel.GetChainDetailsByChainIDAndFamily(b.ChainID, b.Type)
			require.NoError(t, err)
			if strconv.FormatUint(details.ChainSelector, 10) != cantonSelector {
				continue
			}

			return b.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants[0].GRPCLedgerAPIURL
		}
	}
	require.FailNowf(t, "no canton chain found for selector %s", cantonSelector)

	return ""
}

//nolint:paralleltest // can't parallelize this test because it modifies the state of the chain
func TestIntegration_SourceReader_CurseDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestIntegration_SourceReader_CurseDetection test in short mode")
	}

	ccv.RegisterImplFactory(chainsel.FamilyCanton, cantondevenv.NewImplFactory())

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	ctx := ccv.Plog.WithContext(t.Context())
	harness, err := tcapi.NewTestHarness(
		ctx,
		configPath,
		in,
		chainsel.FamilyEVM,
		chainsel.FamilyCanton,
	)
	require.NoError(t, err)

	cantonChain := devenvtests.GetChain(t, blockchain.TypeCanton, in, harness)
	chainMap, err := harness.Lib.ChainsMap(ctx)
	require.NoError(t, err)
	cantonImpl := chainMap[cantonChain.ChainSelector()]
	require.NotNil(t, cantonImpl)

	lggr := logger.Test(t)

	rmnRemoteRef, err := in.CLDF.DataStore.Addresses().Get(datastore.NewAddressRefKey(cantonImpl.ChainSelector(), datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	require.NoError(t, err)

	cantonSelectorStr := strconv.FormatUint(cantonChain.ChainSelector(), 10)
	cantonConfig := getCantonConfig(t, in, cantonSelectorStr)
	readerConfig := cantonConfig.ReaderConfigs[cantonSelectorStr]
	blockchainInfo := cantonConfig.BlockchainInfos[cantonSelectorStr]
	authProvider, err := blockchainInfo.Auth.NewProvider(ctx)
	require.NoError(t, err)
	t.Logf("readerConfig: %+v", readerConfig)
	grpcURL := getCantonGRPCURL(t, in, cantonSelectorStr)
	t.Logf("grpcURL: %s", grpcURL)
	sourceReader, err := sourcereader.NewSourceReader(
		lggr,
		grpcURL,
		readerConfig,
		contracts.HexToInstanceAddress(rmnRemoteRef.Address),
		grpc.WithTransportCredentials(authProvider.TransportCredentials()),
		grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
	)
	require.NoError(t, err)

	// should return no cursed subjects initially
	subjects, err := sourceReader.GetRMNCursedSubjects(ctx)
	require.NoError(t, err)
	require.Empty(t, subjects)

	// TODO: this should be a helper function in protocol.
	var subject protocol.Bytes16
	binary.BigEndian.PutUint64(subject[8:], cantonChain.ChainSelector())

	require.NoError(t, cantonChain.Curse(ctx, [][16]byte{subject}))

	subjects, err = sourceReader.GetRMNCursedSubjects(ctx)
	require.NoError(t, err)
	require.Equal(t, []protocol.Bytes16{subject}, subjects)

	// Uncurse the subject and verify that it is no longer cursed
	require.NoError(t, cantonChain.Uncurse(ctx, [][16]byte{subject}))

	// should return no cursed subjects
	subjects, err = sourceReader.GetRMNCursedSubjects(ctx)
	require.NoError(t, err)
	require.Empty(t, subjects)
}
