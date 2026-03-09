package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

func TestTokenAdminRegistry_SetPoolViaMCMS(t *testing.T) {
	t.Parallel()

	// Use shared TAR environment
	sharedEnv := GetSharedTAREnvironment(t)
	participant := sharedEnv.Participant
	mcmsPkgID := sharedEnv.McmsPkgID
	tarPkgID := sharedEnv.TarPkgID
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	// Deploy MCMS with minDelay=0 for testing
	mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	// Deploy TokenAdminRegistry with empty tokenConfigs
	tarInstanceID := "tar-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("test-token"),
	}

	tarCid := createTokenAdminRegistryEmpty(t, participant, tarPkgID, ccipOwner, tarInstanceID)

	// Set up admin for the instrument so SetPool can succeed
	tarCid = exerciseProposeAndAcceptAdmin(t, participant, tarPkgID, ccipOwner, tarCid, testInstrumentId)

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := tokenadminregistry.NewContract(tarPkgID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Build set_pool operation
	poolReg := &tokenadminregistry.PoolRegistration{
		PoolOwner:      types.PARTY(ccipOwner),
		PoolInstanceId: types.TEXT("test-pool-001"),
	}
	encodedSetPool, err := tarContract.Encoder().TokenAdminRegistrySetPoolMCMSParams(
		tokenadminregistry.TokenAdminRegistrySetPoolMCMSParams{
			InstrumentId: testInstrumentId,
			TokenPool:    poolReg,
		},
	)
	require.NoError(t, err)
	setPoolData := encodedSetPool.OperationData

	// Build timelock calls
	calls := []mcms.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("SetPool"),
		OperationData:         types.TEXT(setPoolData),
	}}
	salt := uuid.New().String()[:8]

	// Encode schedule params
	scheduleParams := mcms.ScheduleBatchParams{
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		DelaySecs:   types.INT64(0),
	}
	scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

	// Build MCMS proposal and sign
	proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
	proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
		AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	// SetRoot
	mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	// ScheduleBatch
	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

	// ExecuteScheduledBatch - this exercises the set_pool choice on TokenAdminRegistry via MCMS
	newTarCid := executeScheduledBatchReturningTargetCid(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{tarInstanceAddr: tarCid}, tarPkgID, "TokenAdminRegistry")

	// Verify the TokenAdminRegistry was updated (new contract created)
	require.NotEqual(t, tarCid, newTarCid, "TokenAdminRegistry contract should have been updated")

	// Verify pool was set by querying the updated contract
	config := queryTokenConfigFromContract(t, participant, tarPkgID, newTarCid, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should exist after set_pool")
	require.NotNil(t, config.TokenPool, "pool should be set after MCMS execution")
	require.Equal(t, string(poolReg.PoolOwner), string(config.TokenPool.PoolOwner))
	require.Equal(t, string(poolReg.PoolInstanceId), string(config.TokenPool.PoolInstanceId))
}

func TestTokenAdminRegistry_ClearPoolViaMCMS(t *testing.T) {
	t.Parallel()

	// Use shared TAR environment
	sharedEnv := GetSharedTAREnvironment(t)
	participant := sharedEnv.Participant
	mcmsPkgID := sharedEnv.McmsPkgID
	tarPkgID := sharedEnv.TarPkgID
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-clear-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	tarInstanceID := "tar-clear-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("clear-token"),
	}

	// Create TokenAdminRegistry and set up admin
	tarCid := createTokenAdminRegistryEmpty(t, participant, tarPkgID, ccipOwner, tarInstanceID)
	tarCid = exerciseProposeAndAcceptAdmin(t, participant, tarPkgID, ccipOwner, tarCid, testInstrumentId)

	// Set initial pool directly (not via MCMS) so we have something to clear
	initialPool := &tokenadminregistry.PoolRegistration{
		PoolOwner:      types.PARTY(ccipOwner),
		PoolInstanceId: types.TEXT("initial-pool"),
	}
	tarCid = exerciseSetPoolDirectly(t, participant, tarPkgID, ccipOwner, tarCid, testInstrumentId, initialPool)

	// Verify pool is initially set
	config := queryTokenConfigFromContract(t, participant, tarPkgID, tarCid, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should exist")
	require.NotNil(t, config.TokenPool, "pool should be set initially")

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := tokenadminregistry.NewContract(tarPkgID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Clear pool by setting to None via MCMS
	encodedClearPool, err := tarContract.Encoder().TokenAdminRegistrySetPoolMCMSParams(
		tokenadminregistry.TokenAdminRegistrySetPoolMCMSParams{
			InstrumentId: testInstrumentId,
			TokenPool:    nil, // nil = clear pool
		},
	)
	require.NoError(t, err)
	clearPoolData := encodedClearPool.OperationData

	calls := []mcms.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("SetPool"),
		OperationData:         types.TEXT(clearPoolData),
	}}
	salt := uuid.New().String()[:8]

	scheduleParams := mcms.ScheduleBatchParams{
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		DelaySecs:   types.INT64(0),
	}
	scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

	proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
	proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
		AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)
	newTarCid := executeScheduledBatchReturningTargetCid(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{tarInstanceAddr: tarCid}, tarPkgID, "TokenAdminRegistry")
	_ = mcmsCid

	// Verify pool was cleared
	config = queryTokenConfigFromContract(t, participant, tarPkgID, newTarCid, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should still exist")
	require.Nil(t, config.TokenPool, "pool should be cleared after MCMS execution")
}

func TestTokenAdminRegistry_ProposeAdminViaMCMS(t *testing.T) {
	t.Parallel()

	// Use shared TAR environment
	sharedEnv := GetSharedTAREnvironment(t)
	participant := sharedEnv.Participant
	mcmsPkgID := sharedEnv.McmsPkgID
	tarPkgID := sharedEnv.TarPkgID
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-propose-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	mcmsCid := createMCMSMultiRole(t, participant, mcmsPkgID, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	tarInstanceID := "tar-propose-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("propose-token"),
	}

	// Create TokenAdminRegistry (no config needed -- ProposeAdministrator creates config entries)
	tarCid := createTokenAdminRegistryEmpty(t, participant, tarPkgID, ccipOwner, tarInstanceID)

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := tokenadminregistry.NewContract(tarPkgID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Propose a new admin via MCMS
	newAdmin := types.PARTY(ccipOwner) // For simplicity, propose self as admin
	encodedPropose, err := tarContract.Encoder().TokenAdminRegistryProposeAdministratorMCMSParams(
		tokenadminregistry.TokenAdminRegistryProposeAdministratorMCMSParams{
			InstrumentId: testInstrumentId,
			NewAdmin:     newAdmin,
		},
	)
	require.NoError(t, err)
	proposeData := encodedPropose.OperationData

	calls := []mcms.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("ProposeAdministrator"),
		OperationData:         types.TEXT(proposeData),
	}}
	salt := uuid.New().String()[:8]

	scheduleParams := mcms.ScheduleBatchParams{
		Calls:       calls,
		Predecessor: types.TEXT(ZeroHash),
		Salt:        types.TEXT(salt),
		DelaySecs:   types.INT64(0),
	}
	scheduleChoice := MustEncodeScheduleBatch(t, mcmsEncoder, scheduleParams)

	proposerMultisigID := MakeMcmsId(mcmsInstanceAddr, MCMSRoleProposer)
	proposal := NewMCMSProposal(int(chainID), proposerMultisigID, 0, false).
		AddOperation(mcmsInstanceAddr, scheduleChoice.Choice, scheduleChoice.OperationData).
		Build()

	validUntil := time.Now().Add(1 * time.Hour)
	signatures, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	mcmsCid = setRootWithRole(t, participant, mcmsPkgID, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, mcmsPkgID, ccipOwner, mcmsCid, proposal.Operations[0], opProof)
	newTarCid := executeScheduledBatchReturningTargetCid(t, participant, mcmsPkgID, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{tarInstanceAddr: tarCid}, tarPkgID, "TokenAdminRegistry")
	_ = mcmsCid

	// Verify pending admin was set
	config := queryTokenConfigFromContract(t, participant, tarPkgID, newTarCid, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should exist")
	require.NotNil(t, config.PendingAdmin, "pendingAdmin should be set after propose_administrator")
	require.Equal(t, string(newAdmin), string(*config.PendingAdmin))
}

// Helper functions

func createTokenAdminRegistryEmpty(
	t *testing.T,
	participant canton.Participant,
	tarPkgID string,
	owner string,
	instanceID string,
) string {
	t.Helper()

	tarContract := tokenadminregistry.TokenAdminRegistry{
		InstanceId:   types.TEXT(instanceID),
		Owner:        types.PARTY(owner),
		TokenConfigs: types.GENMAP{}, // Empty initially
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  tarPkgID,
							ModuleName: "CCIP.TokenAdminRegistry",
							EntityName: "TokenAdminRegistry",
						},
						CreateArguments: ledger.ConvertToRecord(tarContract),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	return res.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
}

// exerciseSetPoolDirectly exercises the TokenAdminRegistry_SetPool choice directly (not via MCMS).
// Used to set up initial state for tests.
func exerciseSetPoolDirectly(
	t *testing.T,
	participant canton.Participant,
	tarPkgID string,
	owner string,
	tarCid string,
	instrumentId splice.InstrumentId,
	pool *tokenadminregistry.PoolRegistration,
) string {
	t.Helper()

	setPoolArgs := tokenadminregistry.TokenAdminRegistrySetPool{
		InstrumentId: instrumentId,
		TokenPool:    pool,
		Caller:       types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  tarPkgID,
							ModuleName: "CCIP.TokenAdminRegistry",
							EntityName: "TokenAdminRegistry",
						},
						ContractId:     tarCid,
						Choice:         "TokenAdminRegistry_SetPool",
						ChoiceArgument: ledger.MapToValue(setPoolArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	// Return the new contract ID
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TokenAdminRegistry" {
			return created.GetContractId()
		}
	}
	t.Fatal("no TokenAdminRegistry contract created after SetPool")

	return ""
}

// exerciseProposeAndAcceptAdmin sets up a token config entry with an admin for the given instrumentId.
// It exercises ProposeAdministrator (creates the config entry with pendingAdmin) then AcceptAdminRole (promotes to admin).
func exerciseProposeAndAcceptAdmin(
	t *testing.T,
	participant canton.Participant,
	tarPkgID string,
	owner string,
	tarCid string,
	instrumentId splice.InstrumentId,
) string {
	t.Helper()

	proposeArgs := tokenadminregistry.TokenAdminRegistryProposeAdministrator{
		InstrumentId: instrumentId,
		NewAdmin:     types.PARTY(owner),
		Caller:       types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  tarPkgID,
							ModuleName: "CCIP.TokenAdminRegistry",
							EntityName: "TokenAdminRegistry",
						},
						ContractId:     tarCid,
						Choice:         "TokenAdminRegistry_ProposeAdministrator",
						ChoiceArgument: ledger.MapToValue(proposeArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	var newTarCid string
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TokenAdminRegistry" {
			newTarCid = created.GetContractId()
		}
	}
	require.NotEmpty(t, newTarCid, "no TokenAdminRegistry contract created after ProposeAdministrator")

	acceptArgs := tokenadminregistry.TokenAdminRegistryAcceptAdminRole{
		InstrumentId: instrumentId,
		Caller:       types.PARTY(owner),
	}

	res, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  tarPkgID,
							ModuleName: "CCIP.TokenAdminRegistry",
							EntityName: "TokenAdminRegistry",
						},
						ContractId:     newTarCid,
						Choice:         "TokenAdminRegistry_AcceptAdminRole",
						ChoiceArgument: ledger.MapToValue(acceptArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TokenAdminRegistry" {
			return created.GetContractId()
		}
	}
	t.Fatal("no TokenAdminRegistry contract created after AcceptAdminRole")

	return ""
}

// executeScheduledBatchReturningTargetCid executes a scheduled batch and returns the new CID of a target contract.
func executeScheduledBatchReturningTargetCid(
	t *testing.T,
	participant canton.Participant,
	mcmsPkgID string,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcms.TimelockCall,
	predecessor string,
	salt string,
	targetCids map[string]string,
	targetPkgID string,
	targetEntityName string,
) string {
	t.Helper()

	executeArgs := mcms.ExecuteScheduledBatch{
		Submitter:   types.PARTY(owner),
		OpId:        types.TEXT(opID),
		Calls:       calls,
		Predecessor: types.TEXT(predecessor),
		Salt:        types.TEXT(salt),
		TargetCids:  toContractIDMap(targetCids),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  mcmsPkgID,
							ModuleName: "MCMS.Main",
							EntityName: "MCMS",
						},
						ContractId:     mcmsCid,
						Choice:         "ExecuteScheduledBatch",
						ChoiceArgument: ledger.MapToValue(executeArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	// Find the new target contract CID
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil {
			if created.GetTemplateId().GetEntityName() == targetEntityName && created.GetTemplateId().GetPackageId() == targetPkgID {
				return created.GetContractId()
			}
		}
	}
	t.Fatalf("no %s contract created after ExecuteScheduledBatch", targetEntityName)

	return ""
}

// queryTokenConfigFromContract queries token config from a specific contract instance.
func queryTokenConfigFromContract(
	t *testing.T,
	participant canton.Participant,
	tarPkgID string,
	tarCid string,
	instrumentId splice.InstrumentId,
) *tokenadminregistry.TokenConfig {
	t.Helper()

	// Get the contract by ID
	contracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{
		PackageId:  tarPkgID,
		ModuleName: "CCIP.TokenAdminRegistry",
		EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	for _, contract := range contracts {
		if contract.GetCreatedEvent().GetContractId() != tarCid {
			continue
		}

		// Found our contract, parse tokenConfigs
		for _, field := range contract.GetCreatedEvent().GetCreateArguments().GetFields() {
			if field.GetLabel() == "tokenConfigs" {
				return parseTokenConfigFromGenMap(t, field.GetValue(), instrumentId)
			}
		}
	}

	return nil
}

func parseTokenConfigFromGenMap(t *testing.T, tokenConfigs *apiv2.Value, instrumentId splice.InstrumentId) *tokenadminregistry.TokenConfig {
	t.Helper()

	genMap := tokenConfigs.GetGenMap()
	if genMap == nil {
		// Empty map or no entries
		return nil
	}

	for _, entry := range genMap.GetEntries() {
		key := entry.GetKey()
		value := entry.GetValue()

		// Check if key matches our instrumentId
		keyRecord := key.GetRecord()
		if keyRecord == nil {
			continue
		}

		var keyAdmin, keyId string
		for _, field := range keyRecord.GetFields() {
			switch field.GetLabel() {
			case "admin":
				keyAdmin = field.GetValue().GetParty()
			case "id":
				keyId = field.GetValue().GetText()
			}
		}

		if keyAdmin == string(instrumentId.Admin) && keyId == string(instrumentId.Id) {
			return parseTokenConfig(t, value)
		}
	}

	return nil
}

func parseTokenConfig(t *testing.T, value *apiv2.Value) *tokenadminregistry.TokenConfig {
	t.Helper()

	record := value.GetRecord()
	if record == nil {
		t.Fatal("TokenConfig value is not a Record")
	}

	config := &tokenadminregistry.TokenConfig{}

	for _, field := range record.GetFields() {
		switch field.GetLabel() {
		case "admin":
			if opt := field.GetValue().GetOptional(); opt != nil && opt.GetValue() != nil {
				admin := types.PARTY(opt.GetValue().GetParty())
				config.Admin = &admin
			}
		case "pendingAdmin":
			if opt := field.GetValue().GetOptional(); opt != nil && opt.GetValue() != nil {
				pending := types.PARTY(opt.GetValue().GetParty())
				config.PendingAdmin = &pending
			}
		case "tokenPool":
			if opt := field.GetValue().GetOptional(); opt != nil && opt.GetValue() != nil {
				poolRecord := opt.GetValue().GetRecord()
				if poolRecord != nil {
					pool := &tokenadminregistry.PoolRegistration{}
					for _, pf := range poolRecord.GetFields() {
						switch pf.GetLabel() {
						case "poolOwner":
							pool.PoolOwner = types.PARTY(pf.GetValue().GetParty())
						case "poolInstanceId":
							pool.PoolInstanceId = types.TEXT(pf.GetValue().GetText())
						}
					}
					config.TokenPool = pool
				}
			}
		}
	}

	return config
}
