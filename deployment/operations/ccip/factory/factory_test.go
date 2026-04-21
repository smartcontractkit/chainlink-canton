package factory

import (
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	ccvsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	mcmsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestDeployCommitteeVerifierChoiceArg_EncodesOptionalParty(t *testing.T) {
	t.Parallel()

	admin := types.PARTY("owner-party")
	args := factorybindings.DeployCommitteeVerifier{
		Contract: ccvsbindings.CommitteeVerifier{
			InstanceId:                   types.TEXT("committeeverifier-1"),
			Owner:                        types.PARTY("owner-party"),
			CcipOwner:                    types.PARTY("owner-party"),
			VersionTag:                   types.TEXT("e9a05a20"),
			AllowListAdmin:               &admin,
			MessageSentObservers:         nil,
			StorageLocations:             []types.TEXT{"ipfs://test"},
			StorageLocationsAdmin:        types.PARTY("owner-party"),
			PendingStorageLocationsAdmin: types.PARTY("owner-party"),
			RemoteChainConfigs:           types.GENMAP{},
			SignerConfigs:                types.GENMAP{},
			Deps: ccvsbindings.CommitteeVerifierDeps{
				RmnRemote: mcmsbindings.RawInstanceAddress{Unpack: types.TEXT("rmn@test")},
			},
		},
	}

	choiceArg := ledger.MapToValue(args.ToMap())
	require.NotNil(t, choiceArg)

	contractMap, ok := args.ToMap()["contract"].(map[string]any)
	require.True(t, ok)
	allowListMap, ok := contractMap["allowListAdmin"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "optional", allowListMap["_type"])

	contractField := findRecordField(t, choiceArg.GetRecord(), "contract")
	allowListField := findRecordField(t, contractField.GetRecord(), "allowListAdmin")
	t.Logf("allowListAdmin field: %#v", allowListField)
	require.NotNil(t, allowListField.GetOptional(), "allowListAdmin should encode as Optional")
	require.NotNil(t, allowListField.GetOptional().Value, "allowListAdmin should encode as Some")
	require.Equal(t, "owner-party", allowListField.GetOptional().Value.GetParty())
}

func TestDeployCommitteeVerifierChoiceArg_EncodesOptionalPartyNone(t *testing.T) {
	t.Parallel()

	args := factorybindings.DeployCommitteeVerifier{
		Contract: ccvsbindings.CommitteeVerifier{
			InstanceId:                   types.TEXT("committeeverifier-1"),
			Owner:                        types.PARTY("owner-party"),
			CcipOwner:                    types.PARTY("owner-party"),
			VersionTag:                   types.TEXT("e9a05a20"),
			AllowListAdmin:               nil,
			MessageSentObservers:         nil,
			StorageLocations:             []types.TEXT{"ipfs://test"},
			StorageLocationsAdmin:        types.PARTY("owner-party"),
			PendingStorageLocationsAdmin: types.PARTY("owner-party"),
			RemoteChainConfigs:           types.GENMAP{},
			SignerConfigs:                types.GENMAP{},
			Deps: ccvsbindings.CommitteeVerifierDeps{
				RmnRemote: mcmsbindings.RawInstanceAddress{Unpack: types.TEXT("rmn@test")},
			},
		},
	}

	choiceArg := ledger.MapToValue(args.ToMap())
	require.NotNil(t, choiceArg)

	contractField := findRecordField(t, choiceArg.GetRecord(), "contract")
	allowListField := findRecordField(t, contractField.GetRecord(), "allowListAdmin")
	require.NotNil(t, allowListField.GetOptional(), "allowListAdmin should encode as Optional")
	require.Nil(t, allowListField.GetOptional().Value, "allowListAdmin should encode as None")
}

func findRecordField(t *testing.T, record *apiv2.Record, label string) *apiv2.Value {
	t.Helper()

	require.NotNil(t, record)
	for _, field := range record.Fields {
		if field.Label == label {
			return field.Value
		}
	}

	t.Fatalf("field %q not found", label)
	return nil
}
