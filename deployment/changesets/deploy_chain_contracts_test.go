package changesets

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

func TestDeployChainContracts(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	ccipOwnerParty := participant.PartyID

	// Upload Dars
	uploadDARs(t, participant,
		contracts.CCIPRuntime,
		contracts.CCIPCore,
		contracts.CCIPCommitteeVerifier,
		contracts.CCIPLockReleaseTokenPool,
	)

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
	)
	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	ccvSignerPubKeys := make([]types.TEXT, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, types.TEXT(pubKeyHex))
	}
	versionTag := "e9a05a20"
	_ = ccvSignerPubKeys // The signers are set during lane deployment
	// ccvID := versionTag + "@" + user.PrimaryParty

	chainSelector := types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10))
	config := CantonCSDeps[DeployChainContractsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: ccipOwnerParty,
				RMNOwnerParty:  ccipOwnerParty,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: committeeverifier.CommitteeVerifier{
							Owner:                        types.PARTY(ccipOwnerParty),
							CcipOwner:                    types.PARTY(ccipOwnerParty),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObservers:         nil,
							StorageLocations:             []types.TEXT{"ipfs://test-receive"},
							StorageLocationsAdmin:        types.PARTY(ccipOwnerParty),
							PendingStorageLocationsAdmin: types.PARTY(ccipOwnerParty),
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: core.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: chainSelector,
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: core.RMNRemote{
						CursedSubjects: nil,
					},
				},
				NativeInstrumentId: splice_api_token_holding_v1.InstrumentId{
					Admin: types.PARTY(ccipOwnerParty),
					Id:    "LINK",
				},
			},
		},
	}

	out, err := DeployChainContracts{}.Apply(*env, config)
	require.NoError(t, err)

	addresses := out.DataStore.Addresses().Filter()
	for i, address := range addresses {
		fmt.Printf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier)
	}
}
