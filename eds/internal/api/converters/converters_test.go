package converters

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func TestResolveRawOrHashedInstanceAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address common.RawOrHashedAddress
		want    contracts.InstanceAddress
		wantErr bool
	}{
		{
			name:    "instanceAddress",
			address: InstanceAddressAsRawOrHashedAddress(contracts.HexToInstanceAddress("0x1234")),
			want:    contracts.HexToInstanceAddress("0x1234"),
			wantErr: false,
		}, {
			name:    "rawInstanceAddress",
			address: RawInstanceAddressAsRawOrHashedAddress(must(contracts.RawInstanceAddressFromString("onramp@ccipOwner"))),
			want:    must(contracts.RawInstanceAddressFromString("onramp@ccipOwner")).InstanceAddress(),
			wantErr: false,
		}, {
			name: "invalidInstanceAddress",
			address: func() common.RawOrHashedAddress {
				var address common.RawOrHashedAddress
				_ = address.FromInstanceAddress("0xasdf")

				return address
			}(),
			want:    contracts.InstanceAddress{},
			wantErr: true,
		}, {
			name: "invalidRawInstanceAddress",
			address: func() common.RawOrHashedAddress {
				var address common.RawOrHashedAddress
				_ = address.FromRawInstanceAddress("invalidAddress")

				return address
			}(),
			want:    contracts.InstanceAddress{},
			wantErr: true,
		}, {
			name:    "emptyAddress",
			address: common.RawOrHashedAddress{},
			want:    contracts.InstanceAddress{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ResolveRawOrHashedAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveRawOrHashedInstanceAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveRawOrHashedInstanceAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
