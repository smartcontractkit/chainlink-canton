package contracts

import (
	"embed"
	"fmt"
	"slices"
)

//go:embed dars
var Dars embed.FS

//go:embed dependencies/splice
var SpliceDependencies embed.FS

type Package string

const (
	Coin = Package("coin")
	Link = Package("link")

	ChainlinkAPI         = Package("chainlink-api")
	ChainlinkInstanceAPI = ChainlinkAPI

	MCMSAPI      = Package("mcms-api")
	MCMSCore     = Package("mcms-core")
	MCMS         = MCMSCore
	MCMSTest     = Package("mcms-test")
	GlobalConfig = Package("globalconfig")

	CCIPCore                 = Package("ccip-core")
	CCIPExtensionAPI         = Package("ccip-extension-api")
	CCIPRuntime              = Package("ccip-runtime")
	CCIPMain                 = CCIPRuntime
	CCIPSender               = Package("ccip-sender")
	CCIPReceiver             = Package("ccip-receiver")
	CCIPCommitteeVerifier    = Package("ccip-committee-verifier")
	CCIPExecutor             = Package("ccip-executor")
	CCIPLockReleaseTokenPool = Package("ccip-lock-release-token-pool")
	CCIPBurnMintTokenPool    = Package("ccip-burn-mint-token-pool")
	CCIPFactory              = Package("ccip-factory")
	CCIPTest                 = Package("ccip-test")

	CCIPClient             = CCIPCore
	CCIPCommon             = CCIPCore
	CCIPFeeQuoter          = CCIPCore
	CCIPTokenAdminRegistry = CCIPCore
	CCIPRMN                = CCIPCore
	CCIPOnRamp             = CCIPRuntime
	CCIPOffRamp            = CCIPRuntime
	CCIPPerPartyRouter     = CCIPRuntime
	CCIPPoolInterfaces     = CCIPExtensionAPI

	SpliceApiTokenBurnMintV1            = Package("splice-api-token-burn-mint-v1")
	SpliceApiTokenHoldingV1             = Package("splice-api-token-holding-v1")
	SpliceApiTokenMetadataV1            = Package("splice-api-token-metadata-v1")
	SpliceApiTokenTransferInstructionV1 = Package("splice-api-token-transfer-instruction-v1")
)

const CurrentVersion = "current"

// ReleaseDir is the frozen production DAR snapshot (e.g. v2_0_0).
// Individual packages keep their own semver in the filename; multiple package
// versions (e.g. globalconfig-1.0.0 and globalconfig-2.0.0) live in the same
// release directory.
const ReleaseDir = "v2_0_0"

var Versions map[Package][]string = map[Package][]string{
	Coin:         []string{CurrentVersion},
	Link:         []string{"2.0.0", CurrentVersion},
	ChainlinkAPI: []string{"2.0.0", CurrentVersion},
	MCMSAPI:      []string{"2.0.0", CurrentVersion},
	MCMSCore:     []string{"2.0.0", CurrentVersion},
	MCMSTest:     []string{CurrentVersion},
	GlobalConfig: []string{"1.0.0", "2.0.0", CurrentVersion},

	CCIPCore:                 []string{"2.0.0", CurrentVersion},
	CCIPExtensionAPI:         []string{"2.0.0", CurrentVersion},
	CCIPRuntime:              []string{"2.0.0", CurrentVersion},
	CCIPSender:               []string{"2.0.0", CurrentVersion},
	CCIPReceiver:             []string{"2.0.0", CurrentVersion},
	CCIPCommitteeVerifier:    []string{"2.0.0", CurrentVersion},
	CCIPExecutor:             []string{"2.0.0", CurrentVersion},
	CCIPLockReleaseTokenPool: []string{"2.0.0", CurrentVersion},
	CCIPBurnMintTokenPool:    []string{"2.0.0", CurrentVersion},
	CCIPFactory:              []string{"2.0.0", CurrentVersion},
	CCIPTest:                 []string{CurrentVersion},

	SpliceApiTokenBurnMintV1:            []string{"1.0.0"},
	SpliceApiTokenHoldingV1:             []string{"1.0.0"},
	SpliceApiTokenMetadataV1:            []string{"1.0.0"},
	SpliceApiTokenTransferInstructionV1: []string{"1.0.0"},
}

// VersionDir maps a DAR version string to its artifact subdirectory.
func VersionDir(version string) string {
	if version == CurrentVersion {
		return CurrentVersion
	}

	return ReleaseDir
}

func darPath(packageName Package, version string) string {
	return fmt.Sprintf("dars/%s/%s-%s.dar", VersionDir(version), packageName, version)
}

func GetDar(packageName Package, version string) ([]byte, error) {
	availableVersions, ok := Versions[packageName]
	if !ok {
		return nil, fmt.Errorf("no available versions for package %s", packageName)
	}

	if !slices.Contains(availableVersions, version) {
		return nil, fmt.Errorf("version %s not found for package %s", version, packageName)
	}

	path := darPath(packageName, version)
	data, err := Dars.ReadFile(path)
	if err != nil {
		// Try to read from Splice dependencies if not found in dars
		data, err = SpliceDependencies.ReadFile(fmt.Sprintf("dependencies/splice/%s-%s.dar", packageName, version))
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded DAR file %s: %w", path, err)
		}
	}

	return data, nil
}

var OutputDirs = map[Package][]string{
	Coin: []string{"coin"},
	Link: []string{"link"},

	ChainlinkAPI: []string{"chainlink", "chainlinkapi"},
	MCMSAPI:      []string{"mcms", "api"},
	MCMSCore:     []string{"mcms", "core"},
	MCMSTest:     []string{"mcms", "mcmstest"},

	CCIPCore:                 []string{"ccip", "core"},
	CCIPExtensionAPI:         []string{"ccip", "extensionapi"},
	CCIPRuntime:              []string{"ccip", "ccipruntime"},
	CCIPReceiver:             []string{"ccip", "receiver"},
	CCIPSender:               []string{"ccip", "sender"},
	CCIPCommitteeVerifier:    []string{"ccip", "committeeverifier"},
	CCIPExecutor:             []string{"ccip", "executor"},
	CCIPLockReleaseTokenPool: []string{"ccip", "lockreleasetokenpool"},
	CCIPBurnMintTokenPool:    []string{"ccip", "burnminttokenpool"},
	CCIPFactory:              []string{"ccip", "factory"},

	SpliceApiTokenBurnMintV1:            []string{"splice", "splice_api_token_burn_mint_v1"},
	SpliceApiTokenHoldingV1:             []string{"splice", "splice_api_token_holding_v1"},
	SpliceApiTokenMetadataV1:            []string{"splice", "splice_api_token_metadata_v1"},
	SpliceApiTokenTransferInstructionV1: []string{"splice", "splice_api_token_transfer_instruction_v1"},
}
