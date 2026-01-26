package deployment

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	CantonCCIPOnRampType             deployment.ContractType = "CantonCCIPOnRamp"
	CantonCCIPOffRampType            deployment.ContractType = "CantonCCIPOffRamp"
	CantonCCIPFeeQuoterType          deployment.ContractType = "CantonCCIPFeeQuoter"
	CantonCCIPTokenAdminRegistryType deployment.ContractType = "CantonCCIPTokenAdminRegistry"
	CantonCCIPCommitteeVerifierType  deployment.ContractType = "CantonCCIPCommitteeVerifier"
	CantonCCIPPerPartyRouterType     deployment.ContractType = "CantonCCIPPerPartyRouter"
	CantonCCIPCCVRegistryType        deployment.ContractType = "CantonCCIPCCVRegistry"
	CantonCCIPCommonType             deployment.ContractType = "CantonCCIPCommon"
	CantonCCIPGlobalConfigType       deployment.ContractType = "CantonCCIPGlobalConfig"

	CantonLinkTokenRegistryType deployment.ContractType = "CantonLinkTokenRegistry"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)
