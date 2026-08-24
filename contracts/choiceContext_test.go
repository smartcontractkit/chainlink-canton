package contracts

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
)

func TestCCIPContextFromData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		ccipContext splice_api_token_metadata_v1.ChoiceContext
		wantErr     bool
	}{
		{
			name: "all values",
			ccipContext: splice_api_token_metadata_v1.ChoiceContext{
				Values: map[string]splice_api_token_metadata_v1.AnyValue{
					"text": {
						AVText: new(types.TEXT("Hello, world!")),
					},
					"int": {
						AVInt: new(types.INT64(42)),
					},
					"decimal": {
						AVDecimal: new(types.NUMERIC("123.456")),
					},
					"bool": {
						AVBool: new(types.BOOL(true)),
					},
					"date": {
						AVDate: new(types.DATE(time.Now())),
					},
					"time": {
						AVTime: new(types.TIMESTAMP(time.Now())),
					},
					"reltime": {
						AVRelTime: new(types.RELTIME(time.Second * 777)),
					},
					"party": {
						AVParty: new(types.PARTY("party123")),
					},
					"contractId": {
						AVContractId: new(types.CONTRACT_ID("00929f6a675c128bb54bef77a8cad8d330badd4330cf71f6a19b596b01f4b66b2dca121220cd655ce27fae64bdf6a203fcd987727de81c3762e46b1882ce2c2fbf0a72c64f")),
					},
					"list": {
						AVList: new([]splice_api_token_metadata_v1.AnyValue{
							{AVText: new(types.TEXT("element1"))},
							{AVText: new(types.TEXT("element2"))},
						}),
					},
					"map": {
						AVMap: new(map[string]splice_api_token_metadata_v1.AnyValue{
							"key1": {AVText: new(types.TEXT("value1"))},
							"key2": {AVText: new(types.TEXT("value2"))},
						}),
					},
				},
			},
			wantErr: false,
		}, {
			name: "empty context",
			ccipContext: splice_api_token_metadata_v1.ChoiceContext{
				Values: map[string]splice_api_token_metadata_v1.AnyValue{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			jsonBytes, err := tt.ccipContext.MarshalJSON()
			require.NoError(t, err)
			var data map[string]any
			err = json.Unmarshal(jsonBytes, &data)
			require.NoError(t, err)

			got, err := ChoiceContextFromData(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CCIPContextFromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Must ignore the unexported fields of time.Time
			if !cmp.Equal(tt.ccipContext, got, cmpopts.IgnoreUnexported(types.DATE{}, types.TIMESTAMP{})) {
				require.Equal(t, tt.ccipContext, got)
			}
		})
	}
}
