package ceremony_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// buildDNSTransactionB64 creates a base64-encoded SignedTopologyTransaction
// containing a DecentralizedNamespaceDefinition mapping.
func buildDNSTransactionB64(t *testing.T, ns string, threshold int32, owners []string) string {
	t.Helper()

	mapping := &protov30.TopologyMapping{
		Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
			DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
				DecentralizedNamespace: ns,
				Threshold:              threshold,
				Owners:                 owners,
			},
		},
	}

	innerTx := &protov30.TopologyTransaction{
		Operation: protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE,
		Serial:    1,
		Mapping:   mapping,
	}
	innerBytes, err := proto.Marshal(innerTx)
	require.NoError(t, err)

	stx := &protov30.SignedTopologyTransaction{
		Transaction: innerBytes,
		Proposal:    true,
	}
	stxBytes, err := proto.Marshal(stx)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(stxBytes)
}

// buildP2PTransactionB64 creates a base64-encoded SignedTopologyTransaction
// containing a PartyToParticipant mapping.
func buildP2PTransactionB64(t *testing.T, party string, threshold uint32, participants []string) string {
	t.Helper()

	hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(participants))
	for i, p := range participants {
		hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
			ParticipantUid: p,
			Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
		}
	}

	mapping := &protov30.TopologyMapping{
		Mapping: &protov30.TopologyMapping_PartyToParticipant{
			PartyToParticipant: &protov30.PartyToParticipant{
				Party:        party,
				Threshold:    threshold,
				Participants: hostingParticipants,
			},
		},
	}

	innerTx := &protov30.TopologyTransaction{
		Operation: protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE,
		Serial:    2,
		Mapping:   mapping,
	}
	innerBytes, err := proto.Marshal(innerTx)
	require.NoError(t, err)

	stx := &protov30.SignedTopologyTransaction{
		Transaction: innerBytes,
		Proposal:    true,
	}
	stxBytes, err := proto.Marshal(stx)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(stxBytes)
}

func TestSummarizeTopologyTx_DNS(t *testing.T) {
	t.Parallel()

	owners := []string{"fp-a", "fp-b", "fp-c"}
	txB64 := buildDNSTransactionB64(t, "test-ns::12345", 2, owners)

	detail, err := ceremony.SummarizeTopologyTx(txB64, "abc123hash", "PAR::p1::xyz")
	require.NoError(t, err)

	assert.Equal(t, "DecentralizedNamespaceDefinition", detail.MappingType)
	assert.Equal(t, "TOPOLOGY_CHANGE_OP_ADD_REPLACE", detail.Operation)
	assert.Equal(t, uint32(1), detail.Serial)
	assert.Equal(t, "test-ns::12345", detail.DNSNamespace)
	assert.Equal(t, int32(2), detail.DNSThreshold)
	assert.Equal(t, owners, detail.DNSOwners)
	assert.Equal(t, "abc123hash", detail.ProposalHash)
	assert.Equal(t, "PAR::p1::xyz", detail.SignerIdentity)
	assert.Equal(t, 0, detail.ExistingSignatures)
}

func TestSummarizeTopologyTx_P2P(t *testing.T) {
	t.Parallel()

	participants := []string{"p1", "p2", "p3"}
	txB64 := buildP2PTransactionB64(t, "party::ns123", 2, participants)

	detail, err := ceremony.SummarizeTopologyTx(txB64, "hashP2P", "p2")
	require.NoError(t, err)

	assert.Equal(t, "PartyToParticipant", detail.MappingType)
	assert.Equal(t, uint32(2), detail.Serial)
	assert.Equal(t, "party::ns123", detail.P2PParty)
	assert.Equal(t, uint32(2), detail.P2PThreshold)
	assert.Len(t, detail.P2PParticipants, 3)
	assert.Contains(t, detail.P2PParticipants[0], "p1")
	assert.Contains(t, detail.P2PParticipants[0], "CONFIRMATION")
}

func TestSummarizeTopologyTx_InvalidBase64(t *testing.T) {
	t.Parallel()

	_, err := ceremony.SummarizeTopologyTx("not-valid-base64!!!", "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decoding base64")
}

func TestSummarizeTopologyTx_CorruptProto(t *testing.T) {
	t.Parallel()

	corruptB64 := base64.StdEncoding.EncodeToString([]byte("corrupt data"))
	_, err := ceremony.SummarizeTopologyTx(corruptB64, "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}
