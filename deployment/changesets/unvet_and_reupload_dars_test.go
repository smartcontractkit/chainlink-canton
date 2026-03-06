package changesets

import (
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/google/uuid"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestUnvetAndReuploadDARs(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	ownerParty := participant.PartyID

	// ---- Setup: upload v1 and create a contract instance ----

	gcV1Dar, err := contracts.GetDar(contracts.GlobalConfig, "1.0.0")
	require.NoError(t, err)

	v1Validated, err := participant.AdminServices.Package.ValidateDar(t.Context(), &participantv30.ValidateDarRequest{
		Data: gcV1Dar,
	})
	require.NoError(t, err)
	v1PkgID := v1Validated.GetMainPackageId()
	require.NotEmpty(t, v1PkgID)

	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: gcV1Dar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err)

	emptyIntMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{}}}
	createV1Res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  v1PkgID,
							ModuleName: "MCMS.GlobalConfig.V1",
							EntityName: "GlobalConfigV1",
						},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-gc"}}},
							{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ownerParty}}},
							{Label: "chainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 1}}},
							{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "aa"}}},
							{Label: "destChainConfigs", Value: emptyIntMap},
							{Label: "sourceChainConfigs", Value: emptyIntMap},
						}},
					},
				},
			}},
			ActAs: []string{ownerParty},
		},
	})
	require.NoError(t, err)
	v1ContractID := createV1Res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	require.NotEmpty(t, v1ContractID)

	gcV2Dar, err := contracts.GetDar(contracts.GlobalConfig, "2.0.0")
	require.NoError(t, err)

	v2Validated, err := participant.AdminServices.Package.ValidateDar(t.Context(), &participantv30.ValidateDarRequest{
		Data: gcV2Dar,
	})
	require.NoError(t, err)
	v2PkgID := v2Validated.GetMainPackageId()
	require.NotEmpty(t, v2PkgID)
	require.NotEqual(t, v1PkgID, v2PkgID)

	// ---- Step 1: Upload + vet new v2 DAR (both v1 and v2 now vetted) ----
	_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: gcV2Dar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err)

	listAfterUpload, err := participant.AdminServices.Package.ListDars(t.Context(), &participantv30.ListDarsRequest{})
	require.NoError(t, err)
	afterUploadIDs := darMainPackageIDs(listAfterUpload)
	assert.Contains(t, afterUploadIDs, v2PkgID, "v2 should be listed after upload")
	assert.Contains(t, afterUploadIDs, v1PkgID, "v1 should still be listed")

	// ---- Step 2: Archive old v1 contract (v1 still vetted, so this works) ----
	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  v1PkgID,
							ModuleName: "MCMS.GlobalConfig.V1",
							EntityName: "GlobalConfigV1",
						},
						ContractId:     v1ContractID,
						Choice:         "Archive",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{}}},
					},
				},
			}},
			ActAs: []string{ownerParty},
		},
	})
	require.NoError(t, err, "should archive v1 contract while v1 package is still vetted")

	// ---- Step 3: Create new v2 contract ----
	createV2Res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  v2PkgID,
							ModuleName: "MCMS.GlobalConfig.V2",
							EntityName: "GlobalConfigV2",
						},
						CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-gc"}}},
							{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ownerParty}}},
							{Label: "chainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 1}}},
							{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "aa"}}},
							{Label: "destChainConfigs", Value: emptyIntMap},
							{Label: "sourceChainConfigs", Value: emptyIntMap},
							{Label: "newFeatureEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{}}}},
						}},
					},
				},
			}},
			ActAs: []string{ownerParty},
		},
	})
	require.NoError(t, err, "should create v2 contract from new package")
	v2ContractID := createV2Res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	require.NotEmpty(t, v2ContractID)

	// ---- Step 4: Unvet old v1 DAR ----
	_, err = participant.AdminServices.Package.UnvetDar(t.Context(), &participantv30.UnvetDarRequest{
		MainPackageId: v1PkgID,
	})
	require.NoError(t, err)

	// ---- Verify: v2 contract is operational after full upgrade ----
	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  v2PkgID,
							ModuleName: "MCMS.GlobalConfig.V2",
							EntityName: "GlobalConfigV2",
						},
						ContractId: v2ContractID,
						Choice:     "GetInstanceIdV2",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "viewer", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ownerParty}}},
						}}}},
					},
				},
			}},
			ActAs: []string{ownerParty},
		},
	})
	require.NoError(t, err, "v2 contract should be operational after full upgrade")
}

func darMainPackageIDs(resp *participantv30.ListDarsResponse) []string {
	var ids []string
	for _, d := range resp.GetDars() {
		if m := d.GetMain(); m != "" {
			ids = append(ids, m)
		}
	}
	return ids
}
