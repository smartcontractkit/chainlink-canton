package changesets

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"testing"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
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

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"

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
	dars := []struct {
		contract contracts.Package
		name     string
	}{
		{contracts.CCIPCommon, "common"},
		{contracts.CCIPOffRamp, "offramp"},
		{contracts.CCIPOnRamp, "onramp"},
		{contracts.CCIPTokenAdminRegistry, "token admin registry"},
		{contracts.CCIPCommitteeVerifier, "committee verifier"},
		{contracts.CCIPLockReleaseTokenPool, "token pool"},
		{contracts.CCIPPerPartyRouter, "per-party router"},
	}

	darData := make([]*participantv30.UploadDarRequest_UploadDarData, 0, len(dars))
	for _, dar := range dars {
		darBytes, err := contracts.GetDar(dar.contract, contracts.CurrentVersion)
		require.NoError(t, err, "failed to get %s dar file", dar.name)
		darData = append(darData, &participantv30.UploadDarRequest_UploadDarData{Bytes: darBytes})
	}

	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars:               darData,
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload DAR files")

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
	versionTag := "49ff34ed"
	_ = ccvSignerPubKeys // The signers are set during lane deployment
	// ccvID := versionTag + "@" + user.PrimaryParty

	chainSelector := types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10))
	config := CantonCSDeps[DeployChainContractsConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty: ccipOwnerParty,
				CommitteeVerifiers: []sequences.CommitteeVerifierParams{
					{
						Template: ccvs.CommitteeVerifier{
							Owner:                        types.PARTY(ccipOwnerParty),
							CcipOwner:                    types.PARTY(ccipOwnerParty),
							VersionTag:                   types.TEXT(versionTag),
							MessageSentObserver:          types.PARTY(ccipOwnerParty),
							StorageLocations:             []types.TEXT{"ipfs://test-receive"},
							StorageLocationsAdmin:        types.PARTY(ccipOwnerParty),
							PendingStorageLocationsAdmin: types.PARTY(ccipOwnerParty),
						},
					},
				},
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: chainSelector,
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						CcipOwner:      "", // Populated by the sequence
						RmnOwner:       types.PARTY(ccipOwnerParty),
						CursedSubjects: nil,
					},
				},
			},
		},
	}

	deployChainContracts := DeployChainContracts{}

	out, err := deployChainContracts.Apply(*env, config)
	require.NoError(t, err)

	addresses := out.DataStore.Addresses().Filter()
	for i, address := range addresses {
		fmt.Printf("Deployed Address %d: ChainSelector=%d, Type=%s, Version=%s, Address=%s, Qualifier=%s\n", i, address.ChainSelector, address.Type, address.Version, address.Address, address.Qualifier)
	}
}
