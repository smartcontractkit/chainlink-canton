package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: `
chain_selector = 8706591216959472610

[server]
host = "0.0.0.0"
port = 8088

[contracts]
[contracts.per_party_router_factory]
party_id = "ccipOwner"
instance_address = "0x02e3bd794d451e6e9962f0c3ca7de62049bce82d040e431ccd3c4d77acaa31cf"
[contracts.on_ramp]
party_id = "ccipOwner"
instance_address = "0xcdc30a75c93191f39bff4236dec958f5b000525268cc383c3cde09529d8fb8f1"
[contracts.off_ramp]
party_id = "ccipOwner"
instance_address = "0xbce16338511449abfc9246d5320aa76e36650d6388afa8dd4f431f0b30fce3eb"
[contracts.global_config]
party_id = "ccipOwner"
instance_address = "0xd44fd347e9fa0f5d8b655d92fbd2865f97c0cbce7517b35bf4b25410947019f1"
[contracts.token_admin_registry]
party_id = "ccipOwner"
instance_address = "0x7d27bb6077ef84ed2c4fce13a1a7bb1f7cf48c2b76a3e6773eb6240382276826"
[contracts.rmn_remote]
party_id = "ccipOwner"
instance_address = "0x7f2ebf216e26051335a9e132d6c013a771d8406378011ca057a6222d3fea1ee5"
[[contracts.ccvs]]
party_id = "ccvOwner"
instance_address = "0x44f3b1f70058285992aaffa899d0015ea4d9c0b5cba4ed3a90f2c99b5ca30011"
[[contracts.ccvs]]
party_id = "ccvOwner"
instance_address = "0xad5d98a90ea7dbba634111605ccc4e5c1dca73a460c403070b49284e950aebf2"

[node]
url = "localhost:8545"
max_retries = 10

[node.auth]
type = "static"
user_id = "local-user"
jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	`,
			want: &Config{
				ChainSelector: 8706591216959472610,
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 8088,
				},
				Contracts: Contracts{
					PerPartyRouterFactory: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0x02e3bd794d451e6e9962f0c3ca7de62049bce82d040e431ccd3c4d77acaa31cf"),
					},
					OnRamp: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0xcdc30a75c93191f39bff4236dec958f5b000525268cc383c3cde09529d8fb8f1"),
					},
					OffRamp: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0xbce16338511449abfc9246d5320aa76e36650d6388afa8dd4f431f0b30fce3eb"),
					},
					GlobalConfig: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0xd44fd347e9fa0f5d8b655d92fbd2865f97c0cbce7517b35bf4b25410947019f1"),
					},
					TokenAdminRegistry: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0x7d27bb6077ef84ed2c4fce13a1a7bb1f7cf48c2b76a3e6773eb6240382276826"),
					},
					RMNRemote: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0x7f2ebf216e26051335a9e132d6c013a771d8406378011ca057a6222d3fea1ee5"),
					},
					CCVs: []ContractIdentifier{
						{
							PartyID:         "ccvOwner",
							InstanceAddress: contracts.HexToInstanceAddress("0x44f3b1f70058285992aaffa899d0015ea4d9c0b5cba4ed3a90f2c99b5ca30011"),
						}, {
							PartyID:         "ccvOwner",
							InstanceAddress: contracts.HexToInstanceAddress("0xad5d98a90ea7dbba634111605ccc4e5c1dca73a460c403070b49284e950aebf2"),
						},
					},
				},
				Node: NodeConfig{
					URL: "localhost:8545",
					AuthConfig: AuthConfig{
						Type:   "static",
						UserID: "local-user",
						JWT:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
					},
					MaxRetries: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Read(strings.NewReader(tt.config))
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
