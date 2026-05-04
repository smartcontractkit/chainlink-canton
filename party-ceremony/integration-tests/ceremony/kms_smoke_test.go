package tests

import (
	"encoding/base64"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"
)

type KMSHarnessSmokeSuite struct {
	CeremonyTestSuite
}

func (s *KMSHarnessSmokeSuite) TestKMSHarnessSmoke() {
	t := s.T()

	kmsCfg := s.kmsConfigFor(0, "smoke")
	namespaceKey, err := s.Actors[0].deps.Client.RegisterKmsSigningKey(
		t.Context(),
		kmsCfg.NamespaceKeyID,
		s.uniqueName("smoke-ns"),
		[]cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE},
	)
	require.NoError(t, err, "register smoke namespace KMS signing key")

	_, err = s.Actors[0].deps.Client.RegisterKmsSigningKey(
		t.Context(),
		kmsCfg.ProtocolKeyID,
		s.uniqueName("smoke-protocol"),
		[]cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL},
	)
	require.NoError(t, err, "register smoke protocol KMS signing key")
	s.assertKMSKeysRegistered(0, kmsCfg)

	namespaceFP, err := helpers.GetPublicKeyFingerprint(namespaceKey.GetPublicKey())
	require.NoError(t, err, "compute smoke namespace fingerprint")
	keyBytes, err := proto.Marshal(namespaceKey)
	require.NoError(t, err, "marshal smoke namespace key")

	reporter := operations.NewMemoryReporter()
	bundle := operations.NewBundle(t.Context, logger.Test(t), reporter)
	_, err = operations.ExecuteOperation(bundle, topology.ProposeNamespaceDelegationOp, s.depsFor(0, kmsCfg), topology.ProposeNSDInput{
		ParticipantID:  s.Actors[0].uid,
		SigningKeyB64:  base64.StdEncoding.EncodeToString(keyBytes),
		Namespace:      namespaceFP,
		SynchronizerID: s.SynchronizerID,
	})
	require.NoError(t, err, "propose smoke namespace delegation through KMS-backed Canton signing")
	s.assertReportsDoNotContainKMS(reporter)
}
