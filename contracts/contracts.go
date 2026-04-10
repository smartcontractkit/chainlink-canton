package contracts

import (
	"embed"
	"fmt"
	"slices"
)

//go:embed dars
var Dars embed.FS

//go:embed dependencies
var Dependencies embed.FS

type Package string

const (
	Coin = Package("coin")

	MCMS         = Package("mcms")
	MCMSTest     = Package("mcms-test")
	GlobalConfig = Package("globalconfig")

	CCIPClient               = Package("ccip-client")
	CCIPCommon               = Package("ccip-common")
	CCIPSender               = Package("ccip-sender")
	CCIPReceiver             = Package("ccip-receiver")
	CCIPCommitteeVerifier    = Package("ccip-committeeverifier")
	CCIPFeeQuoter            = Package("ccip-feequoter")
	CCIPTokenAdminRegistry   = Package("ccip-tokenadminregistry")
	CCIPOnRamp               = Package("ccip-onramp")
	CCIPOffRamp              = Package("ccip-offramp")
	CCIPPoolInterfaces       = Package("ccip-tokenpool-interfaces")
	CCIPLockReleaseTokenPool = Package("ccip-lockreleasetokenpool")
	CCIPPerPartyRouter       = Package("ccip-perpartyrouter")
	CCIPExecutor             = Package("ccip-executor")
	CCIPFactory              = Package("ccip-factory")
	CCIPRMN                  = Package("ccip-rmn")
	CCIPTest                 = Package("ccip-test")
	CCIPEDSRegistry          = Package("ccip-eds-registry")

	SpliceApiTokenBurnMintV1            = Package("splice-api-token-burn-mint-v1")
	SpliceApiTokenHoldingV1             = Package("splice-api-token-holding-v1")
	SpliceApiTokenMetadataV1            = Package("splice-api-token-metadata-v1")
	SpliceApiTokenTransferInstructionV1 = Package("splice-api-token-transfer-instruction-v1")
)

const CurrentVersion = "current"

var Versions map[Package][]string = map[Package][]string{
	Coin: []string{"0.0.1", CurrentVersion},

	MCMS:         []string{"0.0.1", CurrentVersion},
	MCMSTest:     []string{"0.0.1", CurrentVersion},
	GlobalConfig: []string{"1.0.0", "2.0.0", CurrentVersion},

	CCIPClient:               []string{"0.0.1", CurrentVersion},
	CCIPCommon:               []string{"0.0.1", CurrentVersion},
	CCIPSender:               []string{"0.0.1", CurrentVersion},
	CCIPReceiver:             []string{"0.0.1", CurrentVersion},
	CCIPCommitteeVerifier:    []string{"0.0.1", CurrentVersion},
	CCIPFeeQuoter:            []string{"0.0.1", CurrentVersion},
	CCIPTokenAdminRegistry:   []string{"0.0.1", CurrentVersion},
	CCIPOnRamp:               []string{"0.0.1", CurrentVersion},
	CCIPOffRamp:              []string{"0.0.1", CurrentVersion},
	CCIPPoolInterfaces:       []string{"0.0.1", CurrentVersion},
	CCIPLockReleaseTokenPool: []string{"0.0.1", CurrentVersion},
	CCIPPerPartyRouter:       []string{"0.0.1", CurrentVersion},
	CCIPExecutor:             []string{"0.0.1", CurrentVersion},
	CCIPFactory:              []string{"0.0.1", CurrentVersion},
	CCIPRMN:                  []string{"0.0.1", CurrentVersion},
	CCIPTest:                 []string{"0.0.1", CurrentVersion},
	CCIPEDSRegistry:          []string{"0.0.1", CurrentVersion},

	SpliceApiTokenBurnMintV1:            []string{"1.0.0"},
	SpliceApiTokenHoldingV1:             []string{"1.0.0"},
	SpliceApiTokenMetadataV1:            []string{"1.0.0"},
	SpliceApiTokenTransferInstructionV1: []string{"1.0.0"},
}

func GetDar(packageName Package, version string) ([]byte, error) {
	availableVersions, ok := Versions[packageName]
	if !ok {
		return nil, fmt.Errorf("no available versions for package %s", packageName)
	}

	if !slices.Contains(availableVersions, version) {
		return nil, fmt.Errorf("version %s not found for package %s", version, packageName)
	}

	path := fmt.Sprintf("dars/%s-%s.dar", packageName, version)
	data, err := Dars.ReadFile(path)
	if err != nil {
		// Try to read from dependencies if not found in dars
		data, err = Dependencies.ReadFile(fmt.Sprintf("dependencies/%s-%s.dar", packageName, version))
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded DAR file %s: %w", path, err)
		}
	}

	return data, nil
}

var OutputDirs = map[Package][]string{
	Coin: []string{"coin"},

	MCMS:     []string{"mcms"},
	MCMSTest: []string{"mcms", "mcmstest"},

	CCIPClient:               []string{"ccip", "client"},
	CCIPReceiver:             []string{"ccip", "ccipreceiver"},
	CCIPSender:               []string{"ccip", "ccipsender"},
	CCIPCommitteeVerifier:    []string{"ccip", "ccvs"},
	CCIPCommon:               []string{"ccip", "common"},
	CCIPFeeQuoter:            []string{"ccip", "feequoter"},
	CCIPPoolInterfaces:       []string{"ccip", "interfaces"},
	CCIPLockReleaseTokenPool: []string{"ccip", "lockreleasetokenpool"},
	CCIPOffRamp:              []string{"ccip", "offramp"},
	CCIPOnRamp:               []string{"ccip", "onramp"},
	CCIPPerPartyRouter:       []string{"ccip", "perpartyrouter"},
	CCIPExecutor:             []string{"ccip", "executor"},
	CCIPFactory:              []string{"ccip", "factory"},
	CCIPRMN:                  []string{"ccip", "rmn"},
	CCIPTokenAdminRegistry:   []string{"ccip", "tokenadminregistry"},
	CCIPEDSRegistry:          []string{"ccip", "edsregistry"},

	SpliceApiTokenBurnMintV1:            []string{"splice", "splice_api_token_burn_mint_v1"},
	SpliceApiTokenHoldingV1:             []string{"splice", "splice_api_token_holding_v1"},
	SpliceApiTokenMetadataV1:            []string{"splice", "splice_api_token_metadata_v1"},
	SpliceApiTokenTransferInstructionV1: []string{"splice", "splice_api_token_transfer_instruction_v1"},
}
