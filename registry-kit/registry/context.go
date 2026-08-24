package registry

import (
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"

	"github.com/smartcontractkit/go-daml/pkg/types"
)

const (
	CtxKeyInstrumentConfiguration = "utility.digitalasset.com/instrument-configuration"
	CtxKeyIssuerCredentials       = "utility.digitalasset.com/issuer-credentials"
	CtxKeyEnableResultContracts   = "utility.digitalasset.com/enable-result-contracts"
)

func emptyMetadata() splice_api_token_metadata_v1.Metadata {
	return splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}}
}

// MintChoiceContext builds Registry mint choice context for AllocationFactory exercises.
func MintChoiceContext(instrumentConfigCID string, enableResultContracts bool) splice_api_token_metadata_v1.ChoiceContext {
	emptyList := []splice_api_token_metadata_v1.AnyValue{}
	values := map[string]splice_api_token_metadata_v1.AnyValue{
		CtxKeyInstrumentConfiguration: {AVContractId: new(types.CONTRACT_ID(instrumentConfigCID))},
		CtxKeyIssuerCredentials:       {AVList: &emptyList},
	}
	if enableResultContracts {
		values[CtxKeyEnableResultContracts] = splice_api_token_metadata_v1.AnyValue{AVBool: new(types.BOOL(true))}
	}

	return splice_api_token_metadata_v1.ChoiceContext{Values: values}
}
