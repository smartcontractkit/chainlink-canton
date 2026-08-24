package config

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

func validContractIdentifier() ContractIdentifier {
	return ContractIdentifier{
		PartyID:         "party",
		InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
	}
}

func TestConfigValidation(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		s       any
		wantErr bool
	}
	type structSuite struct {
		structName string
		tests      []test
	}

	tests := []structSuite{
		{
			structName: "Factory",
			tests: []test{
				{
					name: "Type invalid type",
					s: Factory{
						Type: "invalidtype",
					},
					wantErr: true,
				}, {
					name: "Type URL valid",
					s: Factory{
						Type:                    FactoryTypeURL,
						TokenStandardURL:        new("http://eds.chain.link"),
						TokenStandardAuthConfig: nil,
					},
					wantErr: false,
				}, {
					name: "Type URL valid",
					s: Factory{
						Type:             FactoryTypeURL,
						TokenStandardURL: new("https://eds.chain.link"),
						TokenStandardAuthConfig: &commonconfig.AuthConfig{
							Type: commonconfig.AuthTypeInsecureStatic,
							JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
						},
					},
					wantErr: false,
				}, {
					name: "Type URL invalid",
					s: Factory{
						Type:             FactoryTypeURL,
						TokenStandardURL: nil, // Required if Type = url
					},
					wantErr: true,
				}, {
					name: "Type URL invalid",
					s: Factory{
						Type:             FactoryTypeURL,
						TokenStandardURL: nil,
						TokenStandardAuthConfig: &commonconfig.AuthConfig{ // Must not be specified, unless TokenStandardURL is specified as well
							Type: commonconfig.AuthTypeInsecureStatic,
							JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
						},
					},
					wantErr: true,
				}, {
					name: "Type urlRequests valid",
					s: Factory{
						Type:             FactoryTypeURLRequests,
						TokenStandardURL: new("https://eds.chain.link"),
						TokenStandardAuthConfig: &commonconfig.AuthConfig{
							Type: commonconfig.AuthTypeInsecureStatic,
							JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
						},
					},
					wantErr: false,
				}, {
					name: "Type urlRequests valid",
					s: Factory{
						Type:                    FactoryTypeURLRequests,
						TokenStandardURL:        new("https://eds.chain.link"),
						TokenStandardAuthConfig: nil,
					},
					wantErr: false,
				}, {
					name: "Type urlRequests invalid",
					s: Factory{
						Type:             FactoryTypeURLRequests,
						TokenStandardURL: nil,
					},
					wantErr: true,
				}, {
					name: "Type Disabled valid",
					s: Factory{
						Type:             FactoryTypeDisabled,
						TokenStandardURL: nil,
					},
					wantErr: false,
				}, {
					name: "Type Disabled invalid",
					s: Factory{
						Type:             FactoryTypeDisabled,
						TokenStandardURL: new("https://eds.chain.link"),
					},
					wantErr: true,
				}, {
					name: "Type Address valid",
					s: Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: false,
				}, {
					name: "Type Address invalid",
					s: Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      nil, // missing
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           nil, // missing
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: nil, // missing
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: Factory{
						Type:             FactoryTypeAddress,
						TemplateId:       new("#package.module.entity"),
						Party:            new("partyid"),
						InstanceAddress:  new(contracts.HexToInstanceAddress("0x1234")),
						TokenStandardURL: new("https://eds.chain.link"), // must not be specified
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
						TokenStandardAuthConfig: &commonconfig.AuthConfig{
							Type: commonconfig.AuthTypeInsecureStatic,
							JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
						}, // must not be specified
					},
					wantErr: true,
				},
			},
		}, {
			structName: "TokenPool",
			tests: []test{
				{
					name: "valid",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:                TokenPoolTypeBurnMint,
						PoolOwner:           "owner",
						Factory:             nil,
						TransferPreapproval: nil,
					},
					wantErr: false,
				}, {
					name: "valid Factory",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:      TokenPoolTypeBurnMint,
						PoolOwner: "owner",
						Factory: &Factory{
							Type: FactoryTypeDisabled,
						},
						TransferPreapproval: nil,
					},
					wantErr: false,
				}, {
					name: "valid Factory",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:      TokenPoolTypeLockRelease,
						PoolOwner: "owner",
						Factory: &Factory{
							Type: FactoryTypeDisabled,
						},
						TransferPreapproval: nil,
					},
					wantErr: false,
				}, {
					name: "valid",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:      TokenPoolTypeBurnMint,
						PoolOwner: "owner",
						Factory: &Factory{
							Type: FactoryTypeDisabled,
						},
						TransferPreapproval: &TransferPreapproval{
							ContextKey: "preapproval",
							TemplateId: "package:module:entity",
						},
					},
					wantErr: false,
				}, {
					name: "invalid TransferPreapproval",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:      TokenPoolTypeBurnMint,
						PoolOwner: "owner",
						Factory: &Factory{
							Type: FactoryTypeDisabled,
						},
						TransferPreapproval: &TransferPreapproval{
							ContextKey: "",
							TemplateId: "",
						},
					},
					wantErr: true,
				}, {
					name: "invalid Factory",
					s: TokenPool{
						ContractIdentifier: ContractIdentifier{
							PartyID:         "owner",
							InstanceAddress: contracts.HexToInstanceAddress("0x1234"),
						},
						Type:      TokenPoolTypeBurnMint,
						PoolOwner: "owner",
						Factory: &Factory{
							Type: FactoryTypeAddress,
							// Missing InstanceAddress
						},
						TransferPreapproval: nil,
					},
					wantErr: true,
				},
			},
		}, {
			structName: "CCIPAPIConfig",
			tests: []test{
				{
					name: "disabled with empty identifiers valid",
					s: CCIPAPIConfig{
						Enabled: false,
					},
					wantErr: false,
				}, {
					name: "enabled with empty OnRamp invalid",
					s: CCIPAPIConfig{
						Enabled:               true,
						PerPartyRouterFactory: validContractIdentifier(),
						OffRamp:               validContractIdentifier(),
						GlobalConfig:          validContractIdentifier(),
						TokenAdminRegistry:    validContractIdentifier(),
						RMNRemote:             validContractIdentifier(),
						FeeQuoter:             validContractIdentifier(),
					},
					wantErr: true,
				}, {
					name: "enabled with all fields populated valid",
					s: CCIPAPIConfig{
						Enabled:               true,
						PerPartyRouterFactory: validContractIdentifier(),
						OnRamp:                validContractIdentifier(),
						OffRamp:               validContractIdentifier(),
						GlobalConfig:          validContractIdentifier(),
						TokenAdminRegistry:    validContractIdentifier(),
						RMNRemote:             validContractIdentifier(),
						FeeQuoter:             validContractIdentifier(),
					},
					wantErr: false,
				},
			},
		},
	}

	for _, suite := range tests {
		t.Run(suite.structName, func(t *testing.T) {
			t.Parallel()
			for _, tt := range suite.tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					validate := validator.New(validator.WithRequiredStructEnabled())
					err := validate.Struct(tt.s)
					if (err != nil) != tt.wantErr {
						t.Errorf("Validate: error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				})
			}
		})
	}
}

func TestCCIPAPIConfigConfigValidate(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		ChainSelector: "8706591216959472610",
		Server: ServerConfig{
			Host:                "0.0.0.0",
			Port:                8088,
			MaxRequestSizeBytes: 1024,
		},
		Node: NodeConfig{
			URL: "http://localhost:8545",
			AuthConfig: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeInsecureStatic,
				JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			},
		},
		GlobalAPIConfig: GlobalAPIConfig{
			MaxBatchSize: 1024,
		},
		CCIPAPIConfig: CCIPAPIConfig{
			Enabled: false,
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: error = %v, want nil", err)
	}
}

func TestPoolOnlyEDSConfigValidate(t *testing.T) {
	t.Parallel()

	poolAddr := contracts.HexToInstanceAddress("0x1bb561b507d97da922c4d13e356e9d6006f02d34ef2ca8c16459e51bf1db51c9")
	linkRegistryAddr := contracts.HexToInstanceAddress("0x690ac078553724389b1de2d9b6f88345cf57f4fc8c3188e51cc0326e11639bf9")
	partySender := "participant2-localparty-1::1220aadfaafa35194855870540bd335c44b32595fe7545e01769b51610dc71b4a6ff"

	cfg := &Config{
		ChainSelector: "8706591216959472610",
		Server: ServerConfig{
			Host:                "0.0.0.0",
			Port:                8089,
			MaxRequestSizeBytes: 1024,
		},
		Node: NodeConfig{
			URL: "http://localhost:8545",
			AuthConfig: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeInsecureStatic,
				JWT:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			},
		},
		GlobalAPIConfig: GlobalAPIConfig{
			MaxBatchSize: 1024,
		},
		CCIPAPIConfig: CCIPAPIConfig{
			Enabled: false,
		},
		CCVAPIConfig: CCVAPIConfig{
			Enabled: false,
		},
		ExecutorAPIConfig: ExecutorAPIConfig{
			Enabled: false,
		},
		TokenPoolAPIConfig: TokenPoolAPIConfig{
			Enabled: true,
			TokenPools: map[string]TokenPool{
				poolAddr.Hex(): {
					Type: TokenPoolTypeBurnMint,
					ContractIdentifier: ContractIdentifier{
						PartyID:         partySender,
						InstanceAddress: poolAddr,
					},
					PoolOwner: partySender,
					Factory: &Factory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#link.module.entity"),
						Party:           new(partySender),
						InstanceAddress: new(linkRegistryAddr),
					},
				},
			},
		},
		TokenStandardAPIConfig: TokenStandardAPIConfig{
			Enabled: false,
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: error = %v, want nil", err)
	}
}
