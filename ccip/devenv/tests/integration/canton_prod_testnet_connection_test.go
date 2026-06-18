package integration

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // register canton impl factory
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestIntegration_CantonProdTestnet_Connection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping prod testnet connection smoke test in short mode")
	}
	if os.Getenv("CANTON_GRPC_URL") == "" {
		t.Skip("CANTON_GRPC_URL unset: not configured for real Canton testnet")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	configPath := filepath.Join("../..", env.ConfigPath())
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	var cantonBlockchain *blockchain.Input
	for _, bc := range in.Blockchains {
		if bc.Type == blockchain.TypeCanton && bc.ChainID == "TestNet" {
			cantonBlockchain = bc
			break
		}
	}
	require.NotNil(t, cantonBlockchain, "need Canton TestNet blockchain in %s", configPath)

	ctx := t.Context()
	chainBC, _, err := cantondevenv.NewCLDF(ctx, cantonBlockchain)
	require.NoError(t, err)

	chain, ok := chainBC.(*canton.Chain)
	require.True(t, ok, "expected *canton.Chain, got %T", chainBC)
	require.NotEmpty(t, chain.Participants)

	participant := chain.Participants[0]
	require.NotEmpty(t, participant.PartyID, "participant should be connected with a party ID")
	t.Logf("connected participant user_id=%s party_id=%s grpc=%s", participant.UserID, participant.PartyID, participant.Endpoints.GRPCLedgerAPIURL)

	holdings, err := testhelpers.ListHoldingsForInstrument(
		ctx,
		participant,
		nil,
		testhelpers.WithHoldingOwner(participant.PartyID),
	)
	require.NoError(t, err)
	t.Logf("listed %d holdings for party %s", len(holdings), participant.PartyID)
	for _, holding := range holdings {
		t.Logf("holding contract_id=%s instrument=%s amount=%s", holding.ContractID, holding.View.InstrumentId.Id, holding.Amount)
	}
}
