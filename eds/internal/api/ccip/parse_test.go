package ccip

import (
	"fmt"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/chainlink/chainlinkapi"
)

func TestParseGlobalConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		createdEvent *apiv2.CreatedEvent
		want         *GlobalConfig
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			createdEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.GlobalConfig{
					InstanceId:    "globalconfig",
					CcipOwner:     "ccipOwner",
					ChainSelector: "1",
					DestChainConfigs: map[types.NUMERIC]core.DestChainConfig2{
						types.NUMERIC("111"): {
							IsEnabled:            true,
							AddressBytesLength:   32,
							TokenReceiverAllowed: false,
							BaseExecutionGasCost: 100,
							OffRampAddress:       "deadbeef",
							DefaultExecutor:      nil,
							LaneMandatedCCVs: []chainlinkapi.RawInstanceAddress{
								contracts.NewRawInstanceAddress("ccv1", "owner").Binding(),
							},
							DefaultCCVs: []chainlinkapi.RawInstanceAddress{
								contracts.NewRawInstanceAddress("ccv2", "owner").Binding(),
							},
							MessageNetworkFeeUSDCents: "1.0",
							TokenNetworkFeeUSDCents:   "400",
						},
					},
					SourceChainConfigs: map[types.NUMERIC]core.SourceChainConfig2{
						types.NUMERIC("222"): {
							IsEnabled: true,
							OnRampAddresses: []types.TEXT{
								types.TEXT("123456789"),
							},
							DefaultCCVs: nil,
							LaneMandatedCCVs: []chainlinkapi.RawInstanceAddress{
								contracts.NewRawInstanceAddress("ccv3", "owner").Binding(),
							},
						},
					},
				}),
			},
			want: &GlobalConfig{
				Address: contracts.NewRawInstanceAddress("globalconfig", "ccipOwner"),
				DestChainConfigs: map[uint64]DestChainConfig{
					111: {
						IsEnabled:       true,
						DefaultExecutor: nil,
						LaneMandatedCCVs: []contracts.RawInstanceAddress{
							contracts.NewRawInstanceAddress("ccv1", "owner"),
						},
						DefaultCCVs: []contracts.RawInstanceAddress{
							contracts.NewRawInstanceAddress("ccv2", "owner"),
						},
					},
				},
				SourceChainConfigs: map[uint64]SourceChainConfig{
					222: {
						IsEnabled:   true,
						DefaultCCVs: []contracts.RawInstanceAddress{},
						LaneMandatedCCVs: []contracts.RawInstanceAddress{
							contracts.NewRawInstanceAddress("ccv3", "owner"),
						},
					},
				},
			},
			wantErr: assert.NoError,
		}, {
			// Validate that create arguments of a wrong type will lead to an error being returned.
			name: "invalid GlobalConfig type",
			createdEvent: &apiv2.CreatedEvent{
				CreateArguments: bindings.MarshalTemplateToRecord(core.FeeQuoter{
					InstanceId: "feequoter",
					CcipOwner:  "ccipOwner",
				}),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.ErrorContains(t, err, "priceUpdaters") // Error message should contain the unknown source fields, i.e. the CreateArguments
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseGlobalConfig(tt.createdEvent)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseGlobalConfig(%v)", tt.createdEvent)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseGlobalConfig(%v)", tt.createdEvent)
		})
	}
}
