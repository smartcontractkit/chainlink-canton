package contracts

import (
	"embed"
	"path/filepath"
)

//go:embed coin ccip dependencies mcms test multi-package.yaml
var Embed embed.FS

type Package string

const (
	Coin = Package("coin")

	MCMS     = Package("mcms")
	MCMSTest = Package("mcms_test")

	CCIPCommon               = Package("ccip_common")
	CCIPReceiver             = Package("ccip_receiver")
	CCIPCommitteeVerifier    = Package("ccip_committeeverifier")
	CCIPFeeQuoter            = Package("ccip_feequoter")
	CCIPTokenAdminRegistry   = Package("ccip_tokenadminregistry")
	CCIPOnRamp               = Package("ccip_onramp")
	CCIPOffRamp              = Package("ccip_offramp")
	CCIPPoolInterfaces       = Package("ccip_pool_interfaces")
	CCIPLockReleaseTokenPool = Package("ccip_lockreleasetokenpool")
	CCIPPerPartyRouter       = Package("ccip_perpartyrouter")
	CCIPTest                 = Package("ccip_test")
)

var Contracts map[Package]string = map[Package]string{
	Coin:     filepath.Join("coin"),
	MCMS:     filepath.Join("mcms"),
	MCMSTest: filepath.Join("mcms", "test"),

	CCIPCommon:               filepath.Join("ccip", "common"),
	CCIPReceiver:             filepath.Join("ccip", "ccipreceiver"),
	CCIPCommitteeVerifier:    filepath.Join("ccip", "ccvs"),
	CCIPFeeQuoter:            filepath.Join("ccip", "feequoter"),
	CCIPTokenAdminRegistry:   filepath.Join("ccip", "tokenAdminRegistry"),
	CCIPOnRamp:               filepath.Join("ccip", "onramp"),
	CCIPOffRamp:              filepath.Join("ccip", "offramp"),
	CCIPPoolInterfaces:       filepath.Join("ccip", "pools", "interfaces"),
	CCIPLockReleaseTokenPool: filepath.Join("ccip", "pools", "lockReleaseTokenPool"),
	CCIPPerPartyRouter:       filepath.Join("ccip", "perpartyrouter"),
	CCIPTest:                 filepath.Join("ccip", "test"),
}
