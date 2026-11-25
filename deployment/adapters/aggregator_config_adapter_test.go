package adapters

import (
	"context"
	"encoding/hex"
	"sort"
	"testing"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
)

func TestCantonCommitteeVerifierOnchainUnsupportedWrites(t *testing.T) {
	t.Parallel()

	adapter := &CantonCommitteeVerifierOnchain{}
	env := deployment.Environment{}

	t.Run("ApplySignatureConfigs", func(t *testing.T) {
		t.Parallel()

		ops, err := adapter.ApplySignatureConfigs(context.Background(), env, 0, "", ccvadapters.SignatureConfigChange{})
		require.Nil(t, ops)
		require.EqualError(t, err, "apply signature configs is not supported on Canton")
	})

	t.Run("SetAllowedFinalityConfig", func(t *testing.T) {
		t.Parallel()

		ops, err := adapter.SetAllowedFinalityConfig(context.Background(), env, 0, "", false, false, 0)
		require.Nil(t, ops)
		require.EqualError(t, err, "set allowed finality config is not supported on Canton")
	})

	t.Run("ApplyAllowlistUpdates", func(t *testing.T) {
		t.Parallel()

		ops, err := adapter.ApplyAllowlistUpdates(context.Background(), env, 0, "", 0, false, nil, nil)
		require.Nil(t, ops)
		require.EqualError(t, err, "apply allowlist updates is not supported on Canton")
	})
}

func TestSignatureConfigsFromCommitteeVerifierKeepsSignerPubKeys(t *testing.T) {
	t.Parallel()
	key, err := gethcrypto.GenerateKey()
	require.NoError(t, err)
	pubkeyHex := hex.EncodeToString(gethcrypto.FromECDSAPub(&key.PublicKey))
	addr := gethcrypto.PubkeyToAddress(key.PublicKey).Hex()

	cv := &committeeverifier.CommitteeVerifier{
		SignerConfigs: map[types.NUMERIC]committeeverifier.SignatureConfig{
			"5009297550715157269": {
				SourceChainSelector: "5009297550715157269",
				Threshold:           1,
				SignerKeys:          []types.TEXT{types.TEXT("0x" + pubkeyHex)},
			},
		},
	}

	configs, err := signatureConfigsFromCommitteeVerifier(cv)
	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, uint64(5009297550715157269), configs[0].SourceChainSelector)
	assert.Equal(t, []string{addr}, configs[0].Signers)
	assert.Equal(t, map[string]string{addr: pubkeyHex}, configs[0].SignerPubKeys,
		"the 65-byte ledger key is kept alongside its derived address")
}

func TestDeployedCommitteeVerifierChains(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()
	for _, ref := range []datastore.AddressRef{
		{ChainSelector: 111, Type: datastore.ContractType(committee_verifier.ContractType), Version: committee_verifier.Version, Address: "cv-1", Qualifier: "default"},
		{ChainSelector: 111, Type: datastore.ContractType(committee_verifier.ContractType), Version: committee_verifier.Version, Address: "cv-2", Qualifier: "secondary"},
		{ChainSelector: 222, Type: datastore.ContractType(committee_verifier.ContractType), Version: committee_verifier.Version, Address: "cv-3"},
		{ChainSelector: 333, Type: "OtherContract", Version: committee_verifier.Version, Address: "other"},
	} {
		require.NoError(t, ds.Addresses().Add(ref))
	}

	chains := (&CantonCommitteeVerifierOnchain{}).DeployedCommitteeVerifierChains(ds.Seal())
	sort.Slice(chains, func(i, j int) bool { return chains[i] < chains[j] })
	assert.Equal(t, []uint64{111, 222}, chains, "every qualifier counts, other contract types do not")
	assert.Empty(t, (&CantonCommitteeVerifierOnchain{}).DeployedCommitteeVerifierChains(nil))
}
