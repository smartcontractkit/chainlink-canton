package config

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

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
			structName: "TransferFactory",
			tests: []test{
				{
					name: "Type invalid type",
					s: TransferFactory{
						Type: "invalidtype",
					},
					wantErr: true,
				}, {
					name: "Type URL valid",
					s: TransferFactory{
						Type:                    FactoryTypeURL,
						TokenStandardURL:        new("http://eds.chain.link"),
						TokenStandardAuthConfig: nil,
					},
					wantErr: false,
				}, {
					name: "Type URL valid",
					s: TransferFactory{
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
					s: TransferFactory{
						Type:             FactoryTypeURL,
						TokenStandardURL: nil,
					},
					wantErr: true,
				}, {
					name: "Type URL invalid",
					s: TransferFactory{
						Type:             FactoryTypeURL,
						TokenStandardURL: nil,
					},
					wantErr: true,
				}, {
					name: "Type Disabled valid",
					s: TransferFactory{
						Type:             FactoryTypeDisabled,
						TokenStandardURL: nil,
					},
					wantErr: false,
				}, {
					name: "Type Disabled invalid",
					s: TransferFactory{
						Type:             FactoryTypeDisabled,
						TokenStandardURL: new("https://eds.chain.link"),
					},
					wantErr: true,
				}, {
					name: "Type Address valid",
					s: TransferFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: false,
				}, {
					name: "Type Address invalid",
					s: TransferFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      nil, // missing
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: TransferFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           nil, // missing
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: TransferFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: nil, // missing
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: TransferFactory{
						Type:             FactoryTypeAddress,
						TemplateId:       new("#package.module.entity"),
						Party:            new("partyid"),
						InstanceAddress:  new(contracts.HexToInstanceAddress("0x1234")),
						TokenStandardURL: new("https://eds.chain.link"), // must not be specified
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: TransferFactory{
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
			structName: "BurnMintFactory",
			tests: []test{
				{
					name: "Type invalid type",
					s: BurnMintFactory{
						Type: "invalidtype",
					},
					wantErr: true,
				}, {
					name: "Type disabled valid",
					s: BurnMintFactory{
						Type: FactoryTypeDisabled,
					},
					wantErr: false,
				}, {
					name: "Type address valid",
					s: BurnMintFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: false,
				}, {
					name: "Type Address invalid",
					s: BurnMintFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      nil, // missing
						Party:           new("partyid"),
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: BurnMintFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           nil, // missing
						InstanceAddress: new(contracts.HexToInstanceAddress("0x1234")),
					},
					wantErr: true,
				}, {
					name: "Type Address invalid",
					s: BurnMintFactory{
						Type:            FactoryTypeAddress,
						TemplateId:      new("#package.module.entity"),
						Party:           new("partyid"),
						InstanceAddress: nil, // missing
					},
					wantErr: true,
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
					fmt.Println(err)
					if (err != nil) != tt.wantErr {
						t.Errorf("Validate: error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				})
			}
		})
	}
}
