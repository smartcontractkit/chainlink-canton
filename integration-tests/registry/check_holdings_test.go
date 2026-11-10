package registry_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link"
	spliceholding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	splicemetadata "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	rkledger "github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

//nolint:paralleltest // Exclusive CTF environment
func TestCheckHoldings(t *testing.T) {
	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
	participant := env.Chain.Participants[0]
	party := participant.PartyID

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		testhelpers.ContractCleanup(t, ctx, env.Chain.Participants)
	})

	dars, err := registry.LoadUtilityDARs()
	require.NoError(t, err)
	linkDAR, err := contracts.GetDar(contracts.Link, contracts.DevVersion)
	require.NoError(t, err)
	dars = append(dars, linkDAR)
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, participant)
	require.NoError(t, err)

	client := rkledger.NewCTFClient(participant)
	bootstrap, err := registry.BootstrapServices(t.Context(), client, party, "LINK")
	require.NoError(t, err)

	recipient := testhelpers.AllocateParty(t, participant, "check-holdings-recipient")
	testhelpers.GrantCanActAs(t, participant, recipient)
	_, err = registry.MintViaAllocationFactory(t.Context(), client, bootstrap, party, "10.0")
	require.NoError(t, err)
	_, err = registry.MintViaAllocationFactory(t.Context(), client, bootstrap, recipient, "5.0")
	require.NoError(t, err)

	token, err := participant.TokenSource.Token()
	require.NoError(t, err)
	binary, configPath := buildCheckHoldingsCLI(t, participant.Endpoints.GRPCLedgerAPIURL,
		participant.Endpoints.ValidatorAPIURL, participant.UserID, party, token.AccessToken)

	t.Run("registry holdings and expected supply", func(t *testing.T) {
		output, err := runCheckHoldings(t, binary, configPath, party, "15.0")
		require.NoError(t, err, output)
		require.Contains(t, output, "15.0000000000")
		require.Contains(t, output, party)
		require.NotContains(t, output, "Findings")
	})

	t.Run("expected supply mismatch", func(t *testing.T) {
		output, err := runCheckHoldings(t, binary, configPath, party, "16.0")
		require.Error(t, err)
		require.Contains(t, output, "supply mismatch")
		require.Contains(t, output, "1 finding(s)")
	})

	t.Run("party must match registrar", func(t *testing.T) {
		output, err := runCheckHoldings(t, binary, configPath, recipient, "15.0")
		require.Error(t, err)
		require.Contains(t, output, "run as the registrar")
	})

	_, err = client.SubmitCreate(t.Context(), party, link.LinkHolding{
		HoldingOwner:        types.PARTY(party),
		HoldingAdmin:        types.PARTY(party),
		HoldingInstrumentId: spliceholding.InstrumentId{Id: "LINK", Admin: types.PARTY(party)},
		HoldingAmount:       "2.0",
		Meta:                splicemetadata.Metadata{},
	})
	require.NoError(t, err)

	t.Run("legacy link holding", func(t *testing.T) {
		output, err := runCheckHoldings(t, binary, configPath, party, "15.0")
		require.Error(t, err)
		require.Contains(t, output, "legacy LINK not drained")
		require.Contains(t, output, "Link.Token:LinkHolding")
		require.Contains(t, output, "1 finding(s)")
	})
}

func buildCheckHoldingsCLI(t *testing.T, grpcURL, validatorURL, userID, party, accessToken string) (string, string) {
	t.Helper()

	root, err := filepath.Abs("../..")
	require.NoError(t, err)
	binary := filepath.Join(t.TempDir(), "canton-ccip-cli")
	build := exec.CommandContext(t.Context(), "go", "build", "-o", binary, "./cmd/canton-ccip-cli")
	build.Dir = root
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))

	config, err := json.Marshal(map[string]any{
		"canton": map[string]any{
			"authType":                    "insecureStatic",
			"authJWT":                     accessToken,
			"participantGRPCLedgerAPIURL": grpcURL,
			"validatorAPIURL":             validatorURL,
			"userID":                      userID,
			"partyID":                     party,
		},
		"evm": map[string]string{"rpcURL": "http://127.0.0.1:8545"},
	})
	require.NoError(t, err)
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configPath, config, 0o600))

	return binary, configPath
}

func runCheckHoldings(t *testing.T, binary, configPath, admin, expectedSupply string) (string, error) {
	t.Helper()

	command := exec.CommandContext(t.Context(), binary,
		"--config", configPath,
		"--network", "testnet",
		"canton", "check-holdings",
		"--instrument", "LINK",
		"--admin", admin,
		"--expected-supply", expectedSupply,
	)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	err := command.Run()
	t.Logf("check-holdings output:\n%s", output.String())

	return output.String(), err
}
