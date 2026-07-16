package adapters

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	mcmsreaderapi "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipshared "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain/shared"
	ccvadapters "github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	ccvshared "github.com/smartcontractkit/chainlink-ccv/deployment/shared"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
)

var tokenPoolVersions = []string{
	"1.6.1",
	"2.0.0",
}

func init() {
	// Map the JD proto Canton chain type to its chain-selectors family so ccv's
	// signing-key fetch (fetch_signing_keys) keeps Canton verifier chain configs
	// instead of skipping them as an unsupported type. Without this, canton NOPs never
	// reach state inference and on-chain canton signers fail to resolve to a NOP.
	ccvshared.RegisterChainTypeFamily(nodev1.ChainType_CHAIN_TYPE_CANTON, chainsel.FamilyCanton)

	// Canton's signer identity is the raw secp256k1 public key (OnchainSigningPubKey),
	// not the EVM-derived address (OnchainSigningAddress). Register a reader in both
	// CCV's and CCIP's shared packages so fetch_signing_keys indexes Canton NOPs under
	// the canton family with the pubkey.
	ccvshared.RegisterSigningIdentityReader(chainsel.FamilyCanton, cantonSigningIdentityReader{})
	ccipshared.RegisterSigningIdentityReader(chainsel.FamilyCanton, cantonSigningIdentityReader{})

	// Canonicalise canton committee signer addresses (raw secp256k1 pubkey ->
	// derived 20-byte ECDSA address) so ccv state inference can match on-chain
	// committee signers back to canton NOPs.
	ccvshared.RegisterAddressNormalizer(chainsel.FamilyCanton, normalizeCantonSignerAddress)

	// Register the onchain adapters
	ccipadapters.GetDeployChainContractsRegistry().Register(chainsel.FamilyCanton, &CantonDeployChainContractsAdapter{})
	ccipadapters.GetChainFamilyRegistry().RegisterChainFamily(chainsel.FamilyCanton, &cantonChainFamilyWithDataStoreCache{})
	ccipdeploy.GetTransferOwnershipRegistry().RegisterAdapter(chainsel.FamilyCanton, ccipdeploy.MCMSVersion, &CantonTransferOwnershipAdapter{})

	// Register the offchain adapters into the ccv per-type singleton registries.
	// The legacy GetRegistry().Register(ChainAdapters{}) shim was removed in the
	// ccv adapter-registration cleanup, so each adapter type registers on its own.
	// TokenVerifier is not implemented for canton yet.
	ccvadapters.GetAggregatorRegistry().Register(chainsel.FamilyCanton, &CantonAggregatorConfigAdapter{})
	ccvadapters.GetExecutorRegistry().Register(chainsel.FamilyCanton, &CantonExecutorConfigAdapter{})
	ccvadapters.GetVerifierRegistry().Register(chainsel.FamilyCanton, &CantonVerifierJobConfigAdapter{})
	ccvadapters.GetIndexerRegistry().Register(chainsel.FamilyCanton, &CantonIndexerConfigAdapter{})
	ccvadapters.GetCommitteeVerifierOnchainRegistry().Register(chainsel.FamilyCanton, &CantonCommitteeVerifierOnchain{})

	lanes.GetLaneAdapterRegistry().RegisterLaneAdapter(chainsel.FamilyCanton, semver.MustParse("2.0.0"), CantonLaneAdapter{})
	mcmsreaderapi.GetRegistry().RegisterMCMSReader(chainsel.FamilyCanton, &CantonMCMSReader{})
	ccipadapters.GetCommitteeVerifierContractRegistry().Register(chainsel.FamilyCanton, &CantonCommitteeVerifierContractAdapter{})

	// Register the canton token adapter for the canton family.
	for _, version := range tokenPoolVersions {
		tokenscore.GetTokenAdapterRegistry().RegisterTokenAdapter(chainsel.FamilyCanton, semver.MustParse(version), CantonTokenAdapter{})
	}

	// Register the curse adapter for the canton family.
	fastcurse.GetCurseRegistry().RegisterNewCurse(fastcurse.CurseRegistryInput{
		CursingFamily:       chainsel.FamilyCanton,
		CursingVersion:      rmn_remote.Version,
		CurseAdapter:        NewCantonCurseAdapter(),
		CurseSubjectAdapter: NewCantonCurseSubjectAdapter(),
	})
}
