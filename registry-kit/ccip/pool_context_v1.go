package ccip

import (
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
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
			string(ratelimiter.RateLimiterContextKey): {
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
