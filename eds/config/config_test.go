package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		config            string
		want              *Config
		wantErr           bool
		wantValidationErr bool
	}{
		{
			name: "valid config",
			//language=toml
			config: `
chain_selector = "8706591216959472610"

[server]
	host = "0.0.0.0"
	port = 8088

[node]
	url = "localhost:8545"
	max_retries = 10
	[node.auth]
		type = "insecureStatic"
		user_id = "local-user"
		jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

[monitoring]
	enabled = true
	[monitoring.beholder]
		insecure_connection = true
		otel_exporter_grpc_endpoint = "beholder:4317"
		otel_exporter_http_endpoint = "http://beholder:4318/v1/traces"
		log_streaming_enabled = true
		metric_reader_interval = 15
		trace_sample_ratio = 0.1
		trace_batch_timeout = 5


[ccip_api]
	enabled = true
	[ccip_api.per_party_router_factory]
		party_id = "ccipOwner"
		instance_address = "0x02e3bd794d451e6e9962f0c3ca7de62049bce82d040e431ccd3c4d77acaa31cf"
	[ccip_api.on_ramp]
		party_id = "ccipOwner"
		instance_address = "0xcdc30a75c93191f39bff4236dec958f5b000525268cc383c3cde09529d8fb8f1"
	[ccip_api.off_ramp]
		party_id = "ccipOwner"
		instance_address = "0xbce16338511449abfc9246d5320aa76e36650d6388afa8dd4f431f0b30fce3eb"
	[ccip_api.global_config]
		party_id = "ccipOwner"
		instance_address = "0xd44fd347e9fa0f5d8b655d92fbd2865f97c0cbce7517b35bf4b25410947019f1"
	[ccip_api.token_admin_registry]
		party_id = "ccipOwner"
		instance_address = "0x7d27bb6077ef84ed2c4fce13a1a7bb1f7cf48c2b76a3e6773eb6240382276826"
	[ccip_api.rmn_remote]
		party_id = "ccipOwner"
		instance_address = "0x7f2ebf216e26051335a9e132d6c013a771d8406378011ca057a6222d3fea1ee5"
	[ccip_api.fee_quoter]
		party_id = "ccipOwner"
		instance_address = "0x9d13995ee04c3e9c441dc1abf307bdc0160119390d77dca2dd95d5604d902fc4"

[ccv_api]
	enabled = true
	[[ccv_api.ccvs]]
		party_id = "ccvOwner"
		instance_address = "0xad5d98a90ea7dbba634111605ccc4e5c1dca73a460c403070b49284e950aebf2"

[executor_api]
	enabled = false

[token_pool_api]
	enabled = true
	[[token_pool_api.token_pools]]
		type = "lockRelease"
		party_id = "tokenPoolOwner"
		instance_address = "0xcd5fe3362a873da7d7ac7b0ae7aa23761d2c8ea7c3872dcfbc715fc8e92f0dec"
		pool_owner = "tokenPoolOwner"
		token_standard_url = "localhost:8545"
		[token_pool_api.token_pools.token_standard_auth]
			type = "insecureStatic"
			jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"
	[[token_pool_api.token_pools]]
		type = "burnMint"
		party_id = "tokenPoolOwner"
		instance_address = "0x44f3b1f70058285992aaffa899d0015ea4d9c0b5cba4ed3a90f2c99b5ca30011"
		pool_owner = "tokenPoolOwner"
		[token_pool_api.token_pools.transfer_preapproval]
			context_key = "transfer-preapproval"
			template_id = "#splice-amulet:Splice.AmuletRules:TransferPreapproval"
	`,
			want: &Config{
				ChainSelector: "8706591216959472610",
				Server: ServerConfig{
					Host: "0.0.0.0",
					Port: 8088,
				},
				GlobalAPIConfig: GlobalAPIConfig{
					MaxBatchSize: 1024, // Default
				},
				CCIPAPIConfig: CCIPAPIConfig{
					Enabled: true,
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
					FeeQuoter: ContractIdentifier{
						PartyID:         "ccipOwner",
						InstanceAddress: contracts.HexToInstanceAddress("0x9d13995ee04c3e9c441dc1abf307bdc0160119390d77dca2dd95d5604d902fc4"),
					},
				},
				CCVAPIConfig: CCVAPIConfig{
					Enabled: true,
					CCVs: []CCV{
						{
							ContractIdentifier: ContractIdentifier{
								PartyID:         "ccvOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0xad5d98a90ea7dbba634111605ccc4e5c1dca73a460c403070b49284e950aebf2"),
							},
						},
					},
				},
				ExecutorAPIConfig: ExecutorAPIConfig{
					Enabled: false,
				},
				TokenPoolAPIConfig: TokenPoolAPIConfig{
					Enabled: true,
					TokenPools: []TokenPool{
						{
							Type: TokenPoolTypeLockRelease,
							ContractIdentifier: ContractIdentifier{
								PartyID:         "tokenPoolOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0xcd5fe3362a873da7d7ac7b0ae7aa23761d2c8ea7c3872dcfbc715fc8e92f0dec"),
							},
							PoolOwner:        "tokenPoolOwner",
							TokenStandardURL: new("localhost:8545"),
							TokenStandardAuthConfig: &commonconfig.AuthConfig{
								Type: commonconfig.AuthTypeInsecureStatic,
								JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
							},
						}, {
							Type: TokenPoolTypeBurnMint,
							ContractIdentifier: ContractIdentifier{
								PartyID:         "tokenPoolOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0x44f3b1f70058285992aaffa899d0015ea4d9c0b5cba4ed3a90f2c99b5ca30011"),
							},
							PoolOwner: "tokenPoolOwner",
							TransferPreapproval: &TransferPreapproval{
								ContextKey: "transfer-preapproval",
								TemplateId: "#splice-amulet:Splice.AmuletRules:TransferPreapproval",
							},
						},
					},
				},
				Node: NodeConfig{
					URL: "localhost:8545",
					AuthConfig: commonconfig.AuthConfig{
						Type:   commonconfig.AuthTypeInsecureStatic,
						UserID: "local-user",
						JWT:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
					},
					MaxRetries: 10,
				},
				Monitoring: MonitoringConfig{
					Enabled: true,
					Beholder: BeholderConfig{
						InsecureConnection:       true,
						OtelExporterGRPCEndpoint: "beholder:4317",
						OtelExporterHTTPEndpoint: "http://beholder:4318/v1/traces",
						LogStreamingEnabled:      true,
						MetricReaderInterval:     15,
						TraceSampleRatio:         0.1,
						TraceBatchTimeout:        5,
					},
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
			require.Equal(t, tt.want, got)
			// Validate
			if err := got.Validate(); (err != nil) != tt.wantValidationErr {
				t.Errorf("Validate() error = %v, wantValidationErr %v", err, tt.wantValidationErr)
			}
		})
	}
}

func TestConfig_Merge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		base    *Config
		config  *Config
		want    *Config
		wantErr bool
	}{
		{
			name: "valid merge",
			base: &Config{
				ChainSelector: "123", // Should be overridden
				Server: ServerConfig{ // Should not be overridden
					Host: "localhost",
					Port: 8080,
				},
				CCVAPIConfig: CCVAPIConfig{
					CCVs: []CCV{
						{ // Should not be overridden
							ContractIdentifier: ContractIdentifier{
								PartyID:         "ccvOwner1",
								InstanceAddress: contracts.HexToInstanceAddress("0x1"),
							},
						},
					},
				},
				TokenPoolAPIConfig: TokenPoolAPIConfig{
					Enabled: false,
					TokenPools: []TokenPool{
						{ // Should be overridden
							Type: TokenPoolTypeLockRelease,
							ContractIdentifier: ContractIdentifier{
								PartyID:         "tokenPoolOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0x2222"),
							},
						},
					},
				},
			},
			config: &Config{
				ChainSelector: "456",
				Server:        ServerConfig{},
				TokenPoolAPIConfig: TokenPoolAPIConfig{
					Enabled: true,
					TokenPools: []TokenPool{
						{
							Type: TokenPoolTypeBurnMint,
							ContractIdentifier: ContractIdentifier{
								PartyID:         "tokenPoolOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0x2222"),
							},
						},
					},
				},
			},
			want: &Config{
				ChainSelector: "456",
				Server: ServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				CCVAPIConfig: CCVAPIConfig{
					CCVs: []CCV{
						{ // Should not be overridden
							ContractIdentifier: ContractIdentifier{
								PartyID:         "ccvOwner1",
								InstanceAddress: contracts.HexToInstanceAddress("0x1"),
							},
						},
					},
				},
				TokenPoolAPIConfig: TokenPoolAPIConfig{
					Enabled: true,
					TokenPools: []TokenPool{
						{
							Type: TokenPoolTypeBurnMint,
							ContractIdentifier: ContractIdentifier{
								PartyID:         "tokenPoolOwner",
								InstanceAddress: contracts.HexToInstanceAddress("0x2222"),
							},
						},
					},
				},
			},
			wantErr: false,
		}, {
			name: "nil override",
			base: &Config{
				ChainSelector: "1111",
				Server:        ServerConfig{},
				Node: NodeConfig{
					URL: "localhost:1234",
				},
			},
			config: nil,
			want: &Config{
				ChainSelector: "1111",
				Server:        ServerConfig{},
				Node: NodeConfig{
					URL: "localhost:1234",
				},
			},
			wantErr: false,
		}, {
			name: "nil base",
			base: nil,
			config: &Config{
				ChainSelector: "1111",
				Server:        ServerConfig{},
				Node: NodeConfig{
					URL: "localhost:1234",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.base.Merge(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)

			b, _ := toml.Marshal(got)
			fmt.Println(string(b))
		})
	}
}
