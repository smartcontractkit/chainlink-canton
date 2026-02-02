package contract

import (
	"reflect"
	"testing"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/coin"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func Test_setInstanceID(t *testing.T) {
	t.Parallel()

	type args struct {
		template   common.Template
		instanceID contracts.InstanceID
	}
	tests := []struct {
		name    string
		args    args
		want    common.Template
		wantErr bool
	}{
		{
			name: "sets instance ID successfully",
			args: args{
				template: coin.CoinRegistry{
					InstanceId: types.TEXT("abc"),
				},
				instanceID: contracts.InstanceID("testID"),
			},
			want: coin.CoinRegistry{
				InstanceId: types.TEXT("testID"),
			},
			wantErr: false,
		}, {
			name: "sets instance ID successfully on a pointer value",
			args: args{
				template: &coin.CoinRegistry{
					InstanceId: types.TEXT("abc"),
				},
				instanceID: contracts.InstanceID("testID"),
			},
			want: &coin.CoinRegistry{
				InstanceId: types.TEXT("testID"),
			},
			wantErr: false,
		}, {
			name: "returns error when no InstanceId field",
			args: args{
				template: coin.CoinHolding{
					Issuer: "testIssuer",
				},
				instanceID: contracts.InstanceID("testID"),
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "returns error when InstanceId field is of wrong type",
			args: args{
				template:   testTemplateInvalid{},
				instanceID: contracts.InstanceID("testID"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := setInstanceID(tt.args.template, tt.args.instanceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("setInstanceID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setInstanceID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type testTemplateInvalid struct{}

func (t testTemplateInvalid) CreateCommand() *model.CreateCommand {
	panic("not implemented")
}

func (t testTemplateInvalid) GetTemplateID() string {
	panic("not implemented")
}
