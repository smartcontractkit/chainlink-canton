package contracts

import (
	"embed"
	"fmt"
	"slices"
)

//go:embed dars
var Dars embed.FS

type Package string

const (
	Coin = Package("coin")

	MCMS     = Package("mcms")
	MCMSTest = Package("mcms-test")

	CCIPCommon               = Package("ccip-common")
	CCIPReceiver             = Package("ccip-receiver")
	CCIPCommitteeVerifier    = Package("ccip-committeeverifier")
	CCIPFeeQuoter            = Package("ccip-feequoter")
	CCIPTokenAdminRegistry   = Package("ccip-tokenadminregistry")
	CCIPOnRamp               = Package("ccip-onramp")
	CCIPOffRamp              = Package("ccip-offramp")
	CCIPPoolInterfaces       = Package("ccip-tokenpool-interfaces")
	CCIPLockReleaseTokenPool = Package("ccip-lockreleasetokenpool")
	CCIPPerPartyRouter       = Package("ccip-perpartyrouter")
	CCIPTest                 = Package("ccip-test")
)

const CurrentVersion = "current"

var Versions map[Package][]string = map[Package][]string{
	Coin: []string{"0.0.1", CurrentVersion},

	MCMS:     []string{"1.0.0", CurrentVersion},
	MCMSTest: []string{"1.0.0", CurrentVersion},

	CCIPCommon:               []string{"1.0.0", CurrentVersion},
	CCIPReceiver:             []string{"1.0.0", CurrentVersion},
	CCIPCommitteeVerifier:    []string{"1.0.0", CurrentVersion},
	CCIPFeeQuoter:            []string{"1.0.0", CurrentVersion},
	CCIPTokenAdminRegistry:   []string{"1.0.0", CurrentVersion},
	CCIPOnRamp:               []string{"1.0.0", CurrentVersion},
	CCIPOffRamp:              []string{"1.0.0", CurrentVersion},
	CCIPPoolInterfaces:       []string{"1.0.0", CurrentVersion},
	CCIPLockReleaseTokenPool: []string{"1.0.0", CurrentVersion},
	CCIPPerPartyRouter:       []string{"1.0.0", CurrentVersion},
	CCIPTest:                 []string{"1.0.0", CurrentVersion},
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
		return nil, fmt.Errorf("failed to read embedded DAR file %s: %w", path, err)
	}

	return data, nil
}
