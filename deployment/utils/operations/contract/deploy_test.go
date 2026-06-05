package contract

import (
	"reflect"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/coin"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/core"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func Test_setInstanceID(t *testing.T) {
	t.Parallel()

	type args struct {
		template   core.Template
		instanceID contracts.InstanceID
	}
	tests := []struct {
		name    string
		args    args
		want    core.Template
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
			name: "sets InstanceId field (MCMS template)",
			args: args{
				template: mcmsCore.MCMS{
					InstanceId: types.TEXT("old"),
				},
				instanceID: contracts.InstanceID("mcms-001"),
			},
			want: mcmsCore.MCMS{
				InstanceId: types.TEXT("mcms-001"),
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
