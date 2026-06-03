package tests

import (
	"path/filepath"
	"testing"

	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

func (s *CeremonyTestSuite) runContractDeployFlow(
	t *testing.T,
	partyID string,
	label string,
	kmsCfgs []client.KMSConfig,
) (operations.SequenceReport[contractdeploy.ContractDeployInput, contractdeploy.ContractDeployOutput], []*recordingKMSAPI) {
	t.Helper()
	require.Len(t, kmsCfgs, len(s.chain.Participants), "one KMS config per participant")

	ledgerClients := make([]client.LedgerClient, len(s.chain.Participants))
	for i, p := range s.chain.Participants {
		lc, conn := s.NewLedgerClient(p)
		t.Cleanup(func() { _ = conn.Close() })
		ledgerClients[i] = lc
	}

	darDir := filepath.Join("..", "..", "..", "contracts", "dars")
	darLoader := ledger.FileDARLoader(darDir)

	recorders := make([]*recordingKMSAPI, len(s.chain.Participants))
	deps := make([]ledger.ContractDeployDeps, len(s.chain.Participants))
	for i, p := range s.chain.Participants {
		adminConn := s.NewAdminConn(p)
		t.Cleanup(func() { _ = adminConn.Close() })

		var kmsAPI client.AWSKMSAPI
		if kmsCfgs[i].ProtocolKeyID != "" {
			require.NotNil(t, s.KMS, "KMS registry is required for contract deploy KMS signing")
			recorders[i] = newRecordingKMSAPI(s.KMS.AWSClient())
			kmsAPI = recorders[i]
		}

		deps[i] = ledger.ContractDeployDeps{
			AdminClient:  s.Actors[i].deps.Client,
			LedgerClient: ledgerClients[i],
			DARLoader:    darLoader,
			SignerFactory: client.NewTransactionSignerFactory(
				s.Actors[i].deps.Client,
				cryptoadminv30.NewVaultServiceClient(adminConn),
				kmsCfgs[i],
				kmsAPI,
			),
			Logger: logger.Test(t),
			UserID: s.chain.Participants[i].UserID,
		}
		t.Logf("Participant %d contract deploy signer ready (%s)", i+1, label)
	}

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: partyID,
		SynchronizerID:       s.SynchronizerID,
		Packages:             []contractdeploy.PackageRef{{Name: "test-test", Version: "current"}},
		TemplateModule:       "Main",
		TemplateEntity:       "DisclosedTarget",
		ContractArgs:         buildContractArgs(t, partyID),
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	t.Logf("%s run 1: p1 uploads DARs (1/3)", label)
	_, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(), "run 1: expected threshold-not-met (1/3)")

	t.Logf("%s run 2: p2 uploads DARs (2/3)", label)
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(), "run 2: expected threshold-not-met (2/3)")

	t.Logf("%s run 3: p3 uploads DARs + prepares + signs (1/3)", label)
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[2], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(), "run 3: expected threshold-not-met for signing (1/3)")

	t.Logf("%s run 4: p1 signs (2/3)", label)
	_, err = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error(), "run 4: expected threshold-not-met for signing (2/3)")

	t.Logf("%s run 5: p2 signs (3/3) + executes + verifies", label)
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[1], input)
	require.NoError(t, err, "run 5: expected full success")

	assert.NotEmpty(t, sr.Output.PackageIDs, "should have at least one package ID")
	assert.NotEmpty(t, sr.Output.PreparedTransactionHash, "should have a prepared transaction hash")
	assert.NotEmpty(t, sr.Output.ContractID, "should have a deployed contract ID")

	t.Logf("%s run 6: p1 idempotency check", label)
	srCached, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, deps[0], input)
	require.NoError(t, err, "run 6: idempotent re-run should succeed")
	assert.Equal(t, sr.Output.PackageIDs, srCached.Output.PackageIDs, "cached package IDs should match")
	assert.Equal(t, sr.Output.ContractID, srCached.Output.ContractID, "cached contract ID should match")
	s.assertReportsDoNotContainKMS(sharedReporter)

	return sr, recorders
}

func assertRecordersMatchKMSConfig(t *testing.T, kmsCfgs []client.KMSConfig, recorders []*recordingKMSAPI) {
	t.Helper()
	require.Len(t, recorders, len(kmsCfgs), "one recorder slot per participant")
	for i, kmsCfg := range kmsCfgs {
		if kmsCfg.ProtocolKeyID == "" {
			require.Nil(t, recorders[i], "participant %d should use vault signer", i+1)
			continue
		}
		require.NotNil(t, recorders[i], "participant %d should use KMS signer", i+1)
		recorders[i].assertUsed(t)
	}
}
