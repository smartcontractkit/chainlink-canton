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
	// Shared

	ChainlinkAPI = Package("chainlink-api")

	// LINK Token

	Link = Package("link")

	// MCMS

	MCMSAPI  = Package("mcms-api")
	MCMSCore = Package("mcms-core")

	// CCIP - V2

	CCIPAPIV2                  = Package("ccip-api-v2")
	CCIPCodecV2                = Package("ccip-codec-v2")
	CCIPTicketsV2              = Package("ccip-tickets-v2")
	CCIPEventsV2               = Package("ccip-events-v2")
	CCIPRateLimiterV2          = Package("ccip-rate-limiter-v2")
	CCIPClientAPIV2            = Package("ccip-client-api-v2")
	CCIPUtilsV2                = Package("ccip-utils-v2")
	CCIPCoreV2                 = Package("ccip-core-v2")
	CCIPExtensionAPIV2         = Package("ccip-extension-api-v2")
	CCIPRuntimeV2              = Package("ccip-runtime-v2")
	CCIPSenderV2               = Package("ccip-sender-v2")
	CCIPReceiverV2             = Package("ccip-receiver-v2")
	CCIPCommitteeVerifierV2    = Package("ccip-committee-verifier-v2")
	CCIPExecutorV2             = Package("ccip-executor-v2")
	CCIPLockReleaseTokenPoolV2 = Package("ccip-lock-release-token-pool-v2")
	CCIPBurnMintTokenPoolV2    = Package("ccip-burn-mint-token-pool-v2")
	CCIPFactoryV2              = Package("ccip-factory-v2")

	// CCIP - Legacy

	CCIPCore                 = Package("ccip-core")
	CCIPExtensionAPI         = Package("ccip-extension-api")
	CCIPRuntime              = Package("ccip-runtime")
	CCIPSender               = Package("ccip-sender")
	CCIPReceiver             = Package("ccip-receiver")
	CCIPCommitteeVerifier    = Package("ccip-committee-verifier")
	CCIPExecutor             = Package("ccip-executor")
	CCIPLockReleaseTokenPool = Package("ccip-lock-release-token-pool")
	CCIPBurnMintTokenPool    = Package("ccip-burn-mint-token-pool")
	CCIPFactory              = Package("ccip-factory")

	// Token Standard

	SpliceApiTokenBurnMintV1            = Package("splice-api-token-burn-mint-v1")
	SpliceApiTokenHoldingV1             = Package("splice-api-token-holding-v1")
	SpliceApiTokenMetadataV1            = Package("splice-api-token-metadata-v1")
	SpliceApiTokenTransferInstructionV1 = Package("splice-api-token-transfer-instruction-v1")

	// Tests

	Coin         = Package("coin")
	GlobalConfig = Package("globalconfig")
	MCMSTest     = Package("mcms-test")
	CCIPTest     = Package("ccip-test")
)

const DevVersion = "dev"

// ReleasedVersions contains all Dars that should be placed in the `dars/released` directory.
// To release a new version of a package:
//  1. Bump the version in the package's daml.yaml file.
//  2. Update the package's daml.yaml `upgrades:` field to point at the previous version in `dars/released`
//  3. Add the new version to the list of released versions for that package in this map.
//  4. Run `make contracts` to regenerate
//
// CI will prevent any of the already committed release artifacts from being changed/altered. If a released package
// is therefore updated without its version being bumped, CI is expected to fail.
var ReleasedVersions map[Package][]string = map[Package][]string{
	CCIPAPIV2:                  []string{"2.0.0"},
	CCIPBurnMintTokenPoolV2:    []string{"2.0.0", "2.1.0"},
	CCIPClientAPIV2:            []string{"2.0.0"},
	CCIPCodecV2:                []string{"2.0.0"},
	CCIPCommitteeVerifierV2:    []string{"2.0.0"},
	CCIPCoreV2:                 []string{"2.0.0", "2.1.0"},
	CCIPEventsV2:               []string{"2.0.0"},
	CCIPExecutorV2:             []string{"2.0.0"},
	CCIPExtensionAPIV2:         []string{"2.0.0"},
	CCIPFactoryV2:              []string{"2.0.0", "2.1.0"},
	CCIPLockReleaseTokenPoolV2: []string{"2.0.0", "2.1.0"},
	CCIPRateLimiterV2:          []string{"2.0.0"},
	CCIPReceiverV2:             []string{"2.0.0", "2.1.0"},
	CCIPRuntimeV2:              []string{"2.0.0", "2.1.0"},
	CCIPSenderV2:               []string{"2.0.0", "2.1.0"},
	CCIPTicketsV2:              []string{"2.0.0"},
	CCIPUtilsV2:                []string{"2.0.0", "2.1.0"},
	ChainlinkAPI:               []string{"2.0.0"},
	Link:                       []string{"2.0.0", "2.1.0"},
	MCMSAPI:                    []string{"1.0.0"},
	MCMSCore:                   []string{"2.0.0"},
}

const (
	ReleaseDir = "released"
	DevDir     = "dev"
)

var Versions map[Package][]string = map[Package][]string{
	ChainlinkAPI: append(ReleasedVersions[ChainlinkAPI], DevVersion),

	Link: append(ReleasedVersions[Link], DevVersion),

	MCMSAPI:  append(ReleasedVersions[MCMSAPI], DevVersion),
	MCMSCore: append(ReleasedVersions[MCMSCore], DevVersion),

	CCIPAPIV2:                  append(ReleasedVersions[CCIPAPIV2], DevVersion),
	CCIPCodecV2:                append(ReleasedVersions[CCIPCodecV2], DevVersion),
	CCIPTicketsV2:              append(ReleasedVersions[CCIPTicketsV2], DevVersion),
	CCIPEventsV2:               append(ReleasedVersions[CCIPEventsV2], DevVersion),
	CCIPRateLimiterV2:          append(ReleasedVersions[CCIPRateLimiterV2], DevVersion),
	CCIPClientAPIV2:            append(ReleasedVersions[CCIPClientAPIV2], DevVersion),
	CCIPUtilsV2:                append(ReleasedVersions[CCIPUtilsV2], DevVersion),
	CCIPCoreV2:                 append(ReleasedVersions[CCIPCoreV2], DevVersion),
	CCIPExtensionAPIV2:         append(ReleasedVersions[CCIPExtensionAPIV2], DevVersion),
	CCIPRuntimeV2:              append(ReleasedVersions[CCIPRuntimeV2], DevVersion),
	CCIPSenderV2:               append(ReleasedVersions[CCIPSenderV2], DevVersion),
	CCIPReceiverV2:             append(ReleasedVersions[CCIPReceiverV2], DevVersion),
	CCIPCommitteeVerifierV2:    append(ReleasedVersions[CCIPCommitteeVerifierV2], DevVersion),
	CCIPExecutorV2:             append(ReleasedVersions[CCIPExecutorV2], DevVersion),
	CCIPLockReleaseTokenPoolV2: append(ReleasedVersions[CCIPLockReleaseTokenPoolV2], DevVersion),
	CCIPBurnMintTokenPoolV2:    append(ReleasedVersions[CCIPBurnMintTokenPoolV2], DevVersion),
	CCIPFactoryV2:              append(ReleasedVersions[CCIPFactoryV2], DevVersion),

	SpliceApiTokenBurnMintV1:            []string{"1.0.0"},
	SpliceApiTokenHoldingV1:             []string{"1.0.0"},
	SpliceApiTokenMetadataV1:            []string{"1.0.0"},
	SpliceApiTokenTransferInstructionV1: []string{"1.0.0"},

	Coin:         []string{DevVersion},
	GlobalConfig: []string{DevVersion},
	MCMSTest:     []string{DevVersion},
}

// versionDir maps a DAR version string to its artifact subdirectory.
func versionDir(version string) string {
	if version == DevVersion {
		return DevDir
	}

	return ReleaseDir
}

func darPath(packageName Package, version string) string {
	return fmt.Sprintf("dars/%s/%s-%s.dar", versionDir(version), packageName, version)
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

var LegacyVersions map[Package][]string = map[Package][]string{
	CCIPBurnMintTokenPool:    []string{"1.0.0", "2.0.0"},
	CCIPCommitteeVerifier:    []string{"1.0.0", "2.0.0"},
	CCIPCore:                 []string{"1.0.0", "2.0.0"},
	CCIPExecutor:             []string{"1.0.0", "2.0.0"},
	CCIPExtensionAPI:         []string{"1.0.0", "2.0.0"},
	CCIPFactory:              []string{"1.0.0", "2.0.0"},
	CCIPLockReleaseTokenPool: []string{"1.0.0", "2.0.0"},
	CCIPReceiver:             []string{"1.0.0", "2.0.0"},
	CCIPRuntime:              []string{"1.0.0", "2.0.0"},
	CCIPSender:               []string{"1.0.0", "2.0.0"},
	ChainlinkAPI:             []string{"1.0.0", "2.0.0"},
	Coin:                     []string{"0.0.1"},
	GlobalConfig:             []string{"1.0.0", "2.0.0"},
	Link:                     []string{"0.0.1", "2.0.0"},
	MCMSAPI:                  []string{"1.0.0"},
	MCMSCore:                 []string{"1.0.0", "2.0.0"},
	MCMSTest:                 []string{"1.0.0"},
}

func GetLegacyDar(packageName Package, version string) ([]byte, error) {
	availableVersions, ok := LegacyVersions[packageName]
	if !ok {
		return nil, fmt.Errorf("no available versions for package %s", packageName)
	}

	if !slices.Contains(availableVersions, version) {
		return nil, fmt.Errorf("version %s not found for package %s", version, packageName)
	}

	path := fmt.Sprintf("dars/legacy/%s-%s.dar", packageName, version)
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

var BindingsOutputDirs = map[Package][]string{
	Coin: []string{"coin"},
	Link: []string{"link"},

	ChainlinkAPI: []string{"chainlink", "chainlinkapi"},
	MCMSAPI:      []string{"mcms", "api"},
	MCMSCore:     []string{"mcms", "core"},
	MCMSTest:     []string{"mcms", "mcmstest"},

	CCIPAPIV2:                  []string{"ccip", "ccipapi"},
	CCIPCodecV2:                []string{"ccip", "ccipcodec"},
	CCIPTicketsV2:              []string{"ccip", "tickets"},
	CCIPEventsV2:               []string{"ccip", "events"},
	CCIPRateLimiterV2:          []string{"ccip", "ratelimiter"},
	CCIPClientAPIV2:            []string{"ccip", "clientapi"},
	CCIPCoreV2:                 []string{"ccip", "core"},
	CCIPExtensionAPIV2:         []string{"ccip", "extensionapi"},
	CCIPRuntimeV2:              []string{"ccip", "ccipruntime"},
	CCIPReceiverV2:             []string{"ccip", "receiver"},
	CCIPSenderV2:               []string{"ccip", "sender"},
	CCIPCommitteeVerifierV2:    []string{"ccip", "committeeverifier"},
	CCIPExecutorV2:             []string{"ccip", "executor"},
	CCIPLockReleaseTokenPoolV2: []string{"ccip", "lockreleasetokenpool"},
	CCIPBurnMintTokenPoolV2:    []string{"ccip", "burnminttokenpool"},
	CCIPFactoryV2:              []string{"ccip", "factory"},

	SpliceApiTokenBurnMintV1:            []string{"splice", "splice_api_token_burn_mint_v1"},
	SpliceApiTokenHoldingV1:             []string{"splice", "splice_api_token_holding_v1"},
	SpliceApiTokenMetadataV1:            []string{"splice", "splice_api_token_metadata_v1"},
	SpliceApiTokenTransferInstructionV1: []string{"splice", "splice_api_token_transfer_instruction_v1"},
}
