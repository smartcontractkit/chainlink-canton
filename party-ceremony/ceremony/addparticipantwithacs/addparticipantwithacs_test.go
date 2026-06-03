package addparticipantwithacs

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/replication"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

var errNotImplemented = errors.New("not implemented")

const (
	testPartyID   = "test-party::1220aabbccdd"
	testSyncID    = "global"
	testSourceUID = "p2"
	testTargetUID = "p3"
)

type mockReplicationClient struct {
	participantID string
}

func (m *mockReplicationClient) GetParticipantUID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockReplicationClient) GetParticipantID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockReplicationClient) ExportAcs(_ context.Context, _ []string, _ string, _ int64) ([]byte, error) {
	return []byte("acs-snapshot"), nil
}

func (m *mockReplicationClient) ImportAcs(context.Context, []byte, string) error { return nil }

func (m *mockReplicationClient) DisconnectSynchronizer(context.Context, string) error { return nil }

func (m *mockReplicationClient) ReconnectSynchronizer(context.Context, string) error { return nil }

func (m *mockReplicationClient) ListConnectedSynchronizers(context.Context) ([]client.SynchronizerInfo, error) {
	return nil, nil
}

func (m *mockReplicationClient) ClearPartyOnboardingFlag(context.Context, string, string, int64) (bool, error) {
	return true, nil
}

func (m *mockReplicationClient) LookupOffsetByTime(context.Context, int64) (int64, error) {
	return 42, nil
}

func (m *mockReplicationClient) GenerateSigningKey(context.Context, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) RegisterKmsSigningKey(context.Context, string, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) GetNamespaceFingerprint(context.Context, string, string, []string) (string, error) {
	return "", nil
}

func (m *mockReplicationClient) GetNamespaceKeyName(context.Context, string, []string) (string, error) {
	return "", nil
}

func (m *mockReplicationClient) GetProtocolKeyFingerprint(context.Context, []string) (string, string, error) {
	return "", "", nil
}

func (m *mockReplicationClient) Authorize(context.Context, uint32, *protov30.TopologyMapping, string, bool, ...string) (*protov30.SignedTopologyTransaction, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) SignTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) ([]*protov30.SignedTopologyTransaction, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) AddTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) error {
	return nil
}

func (m *mockReplicationClient) DNSExists(context.Context, string, string) (bool, error) {
	return false, nil
}

func (m *mockReplicationClient) NSDExists(context.Context, string, string) (bool, error) {
	return false, nil
}

func (m *mockReplicationClient) P2PExists(context.Context, string, string) (bool, error) {
	return false, nil
}

func (m *mockReplicationClient) GetDNS(context.Context, string, string) (*client.DNSState, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) GetP2P(context.Context, string, string) (*client.P2PState, error) {
	return nil, errNotImplemented
}

func (m *mockReplicationClient) ListDecentralizedNamespaces(context.Context, string) ([]*client.DNSState, error) {
	return nil, nil
}

func (m *mockReplicationClient) UploadDar(context.Context, []byte) (string, error) {
	return "", nil
}

func newReplicationDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockReplicationClient{participantID: participantID},
		Logger: logger.Nop(),
	}
}

func TestReplicationStepErr(t *testing.T) {
	t.Parallel()

	opErr := errors.New("participant ID mismatch: expected p2, got p1")

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, replicationStepErr("p1", "p2", nil))
	})

	t.Run("wrong actor returns threshold", func(t *testing.T) {
		t.Parallel()
		err := replicationStepErr("p1", "p2", opErr)
		require.ErrorIs(t, err, ErrThresholdNotMet)
		assert.NotContains(t, err.Error(), "participant ID mismatch")
	})

	t.Run("designated actor preserves error", func(t *testing.T) {
		t.Parallel()
		err := replicationStepErr("p2", "p2", opErr)
		require.Error(t, err)
		assert.False(t, errors.Is(err, ErrThresholdNotMet))
		assert.Contains(t, err.Error(), "participant ID mismatch")
	})
}

func TestExportAcsOp_WrongRunnerParticipantMismatch(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	exportIn := replication.ExportAcsInput{
		ParticipantID:    testSourceUID,
		PartyIDs:         []string{testPartyID},
		SynchronizerID:   testSyncID,
		LedgerOffset:     10,
		TimestampSeconds: 1,
	}

	_, err := operations.ExecuteOperation(b, replication.ExportAcsOp, newReplicationDeps("p1"), exportIn)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "participant ID mismatch")
}

func TestExportAcsOp_NonSourceRunnerUsesCachedReport(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	exportIn := replication.ExportAcsInput{
		ParticipantID:    testSourceUID,
		PartyIDs:         []string{testPartyID},
		SynchronizerID:   testSyncID,
		LedgerOffset:     10,
		TimestampSeconds: 1,
	}

	sourceReport, err := operations.ExecuteOperation(b, replication.ExportAcsOp, newReplicationDeps(testSourceUID), exportIn)
	require.NoError(t, err)
	require.NotEmpty(t, sourceReport.Output.AcsSnapshotB64)

	targetReport, err := operations.ExecuteOperation(b, replication.ExportAcsOp, newReplicationDeps(testTargetUID), exportIn)
	require.NoError(t, err, "non-source runner should reuse cached export report")
	assert.Equal(t, sourceReport.Output.AcsSnapshotB64, targetReport.Output.AcsSnapshotB64)
	assert.Equal(t, sourceReport.Output.SHA256, targetReport.Output.SHA256)
}
