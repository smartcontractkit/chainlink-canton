package ccip

import (
	burnminttokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	ccipcore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// RegistryPoolSendExtraContextV1 builds pool extraContext for LockOrBurn / ReleaseFromTicket
// using latest CCIP bindings (integration-tests lane compatibility).
func RegistryPoolSendExtraContextV1(
	rateLimiterCID, allocationFactoryCID, instrumentConfigCID string,
	enableResultContracts bool,
) splice_api_token_metadata_v1.ChoiceContext {
	emptyList := []splice_api_token_metadata_v1.AnyValue{}
	nestedValues := map[string]splice_api_token_metadata_v1.AnyValue{
		registry.CtxKeyInstrumentConfiguration: {AVContractId: new(types.CONTRACT_ID(instrumentConfigCID))},
		registry.CtxKeyIssuerCredentials:       {AVList: &emptyList},
	}
	if enableResultContracts {
		nestedValues[registry.CtxKeyEnableResultContracts] = splice_api_token_metadata_v1.AnyValue{
			AVBool: new(types.BOOL(true)),
		}
	}

	return splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{
			string(ccipcore.RateLimiterKey): {
				AVContractId: new(types.CONTRACT_ID(rateLimiterCID)),
			},
			string(burnminttokenpool.BurnMintFactoryContextKey): {
				AVContractId: new(types.CONTRACT_ID(allocationFactoryCID)),
			},
			string(burnminttokenpool.BurnMintFactoryExtraArgsContextValuesContextKey): {
				AVMap: &nestedValues,
			},
		},
	}
}
