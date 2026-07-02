package tests

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/core"
	splice "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func TestSetPoolViaMCMS(t *testing.T) {
	t.Parallel()

	// Use shared TAR environment
	sharedEnv := GetSharedTAREnvironment(t)
	participant := sharedEnv.Participant
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	// Deploy MCMS with minDelay=0 for testing
	mcmsCid := createMCMSMultiRole(t, participant, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	// Deploy TokenAdminRegistry with no TokenConfig entries.
	tarInstanceID := "tar-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("test-token"),
	}

	tarCid := createTokenAdminRegistryEmpty(t, participant, ccipOwner, tarInstanceID)

	// Set up admin for the instrument so SetPool can succeed
	tarCid, tokenConfigCid := exerciseProposeAndAcceptAdmin(t, participant, ccipOwner, tarCid, testInstrumentId)
	tokenConfigInstanceAddr := makeTokenConfigInstanceAddr(ccipOwner, testInstrumentId)

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := core.NewContract(fmt.Sprintf("#%s", core.PackageName), "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Build set_pool operation
	poolReg := &core.PoolRegistration{
		PoolOwner:      types.PARTY(ccipOwner),
		PoolInstanceId: types.TEXT("test-pool-001"),
	}
	encodedSetPool, err := tarContract.Encoder().SetPoolParams(
		core.SetPoolParams{
			InstrumentId: testInstrumentId,
			TokenPool:    poolReg,
		},
	)
	require.NoError(t, err)
	setPoolData := encodedSetPool.OperationData

	// Build timelock calls
	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("SetPool"),
		OperationData:         types.TEXT(setPoolData),
	}}
	salt := uuid.New().String()[:8]

	// Encode schedule params
	scheduleParams := mcmsApi.ScheduleBatchParams{
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
	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	// ScheduleBatch
	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, proposal.Operations[0], opProof)

	// ExecuteScheduledBatch - this exercises the set_pool choice on TokenAdminRegistry via MCMS
	newTokenConfigCid := executeScheduledBatchReturningTargetCid(t, participant, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{
		tarInstanceAddr:         tarCid,
		tokenConfigInstanceAddr: tokenConfigCid,
	}, "TokenConfig")

	require.NotEqual(t, tokenConfigCid, newTokenConfigCid, "TokenConfig contract should have been updated")

	// Verify pool was set by querying the standalone TokenConfig contract
	config := queryTokenConfig(t, participant, ccipOwner, testInstrumentId)
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
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-clear-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	mcmsCid := createMCMSMultiRole(t, participant, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	tarInstanceID := "tar-clear-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("clear-token"),
	}

	// Create TokenAdminRegistry and set up admin
	tarCid := createTokenAdminRegistryEmpty(t, participant, ccipOwner, tarInstanceID)
	tarCid, tokenConfigCid := exerciseProposeAndAcceptAdmin(t, participant, ccipOwner, tarCid, testInstrumentId)
	tokenConfigInstanceAddr := makeTokenConfigInstanceAddr(ccipOwner, testInstrumentId)

	// Set initial pool directly (not via MCMS) so we have something to clear
	initialPool := &core.PoolRegistration{
		PoolOwner:      types.PARTY(ccipOwner),
		PoolInstanceId: types.TEXT("initial-pool"),
	}
	tokenConfigCid = exerciseSetPoolDirectly(t, participant, ccipOwner, tarCid, tokenConfigCid, testInstrumentId, initialPool)

	// Verify pool is initially set
	config := queryTokenConfig(t, participant, ccipOwner, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should exist")
	require.NotNil(t, config.TokenPool, "pool should be set initially")

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := core.NewContract(fmt.Sprintf("#%s", core.PackageName), "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Clear pool by setting to None via MCMS
	encodedClearPool, err := tarContract.Encoder().SetPoolParams(
		core.SetPoolParams{
			InstrumentId: testInstrumentId,
			TokenPool:    nil, // nil = clear pool
		},
	)
	require.NoError(t, err)
	clearPoolData := encodedClearPool.OperationData

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("SetPool"),
		OperationData:         types.TEXT(clearPoolData),
	}}
	salt := uuid.New().String()[:8]

	scheduleParams := mcmsApi.ScheduleBatchParams{
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

	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, proposal.Operations[0], opProof)
	newTokenConfigCid := executeScheduledBatchReturningTargetCid(t, participant, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{
		tarInstanceAddr:         tarCid,
		tokenConfigInstanceAddr: tokenConfigCid,
	}, "TokenConfig")
	_ = mcmsCid
	require.NotEqual(t, tokenConfigCid, newTokenConfigCid, "TokenConfig contract should have been updated")

	// Verify pool was cleared
	config = queryTokenConfig(t, participant, ccipOwner, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should still exist")
	require.Nil(t, config.TokenPool, "pool should be cleared after MCMS execution")
}

func TestTokenAdminRegistry_ProposeAdminViaMCMS(t *testing.T) {
	t.Parallel()

	// Use shared TAR environment
	sharedEnv := GetSharedTAREnvironment(t)
	participant := sharedEnv.Participant
	mcmsEncoder := sharedEnv.McmsEncoder
	ccipOwner := sharedEnv.CcipOwner
	cfg := sharedEnv.Config
	sortedSigners := sharedEnv.SortedSigners

	chainID := int64(1)
	baseMcmsID := "mcms-tar-propose-" + uuid.New().String()[:8]
	mcmsInstanceAddr := fmt.Sprintf("%s@%s", baseMcmsID, ccipOwner)

	mcmsCid := createMCMSMultiRole(t, participant, ccipOwner, chainID, baseMcmsID, cfg, 0, nil)

	tarInstanceID := "tar-propose-" + uuid.New().String()[:8]
	tarInstanceAddr := fmt.Sprintf("%s@%s", tarInstanceID, ccipOwner)
	testInstrumentId := splice.InstrumentId{
		Admin: types.PARTY(ccipOwner),
		Id:    types.TEXT("propose-token"),
	}

	// Create TokenAdminRegistry (no config needed -- ProposeAdministrator creates config entries)
	tarCid := createTokenAdminRegistryEmpty(t, participant, ccipOwner, tarInstanceID)

	// Create MCMS encoder for TokenAdminRegistry
	tarContract := core.NewContract(fmt.Sprintf("#%s", core.PackageName), "CCIP.TokenAdminRegistry", "TokenAdminRegistry")

	// Propose a new admin via MCMS
	newAdmin := types.PARTY(ccipOwner) // For simplicity, propose self as admin
	encodedPropose, err := tarContract.Encoder().ProposeAdministrator(
		core.ProposeAdminParams{
			InstrumentId: testInstrumentId,
			NewAdmin:     newAdmin,
		},
	)
	require.NoError(t, err)
	proposeData := encodedPropose.OperationData

	calls := []mcmsApi.TimelockCall{{
		TargetInstanceAddress: types.TEXT(tarInstanceAddr),
		FunctionName:          types.TEXT("ProposeAdministrator"),
		OperationData:         types.TEXT(proposeData),
	}}
	salt := uuid.New().String()[:8]

	scheduleParams := mcmsApi.ScheduleBatchParams{
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

	mcmsCid = setRootWithRole(t, participant, ccipOwner, mcmsCid, "Proposer", proposal, validUntil, signatures)

	opID := HashTimelockOpId(UnwrapTimelockCalls(calls), ZeroHash, salt)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	mcmsCid = scheduleBatch(t, participant, ccipOwner, mcmsCid, proposal.Operations[0], opProof)
	newTokenConfigCid := executeScheduledBatchReturningTargetCid(t, participant, ccipOwner, mcmsCid, opID, calls, ZeroHash, salt, map[string]string{tarInstanceAddr: tarCid}, "TokenConfig")
	_ = mcmsCid
	require.NotEmpty(t, newTokenConfigCid, "TokenConfig should be created")

	// Verify pending admin was set
	config := queryTokenConfig(t, participant, ccipOwner, testInstrumentId)
	require.NotNil(t, config, "TokenConfig should exist")
	require.NotNil(t, config.PendingAdmin, "pendingAdmin should be set after propose_administrator")
	require.Equal(t, string(newAdmin), string(*config.PendingAdmin))
}

// Helper functions

func createTokenAdminRegistryEmpty(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instanceID string,
) string {
	t.Helper()

	tarContract := core.TokenAdminRegistry{
		InstanceId: types.TEXT(instanceID),
		CcipOwner:  types.PARTY(owner),
		EntryCount: 0,
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{
					Create: &apiv2.CreateCommand{
						TemplateId:      contracts.IdentifierFromBinding(core.TokenAdminRegistry{}),
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

// exerciseSetPoolDirectly exercises the SetPool choice directly (not via MCMS).
// Used to set up initial state for tests.
func exerciseSetPoolDirectly(
	t *testing.T,
	participant canton.Participant,
	owner string,
	tarCid string,
	tokenConfigCid string,
	instrumentId splice.InstrumentId,
	pool *core.PoolRegistration,
) string {
	t.Helper()

	setPoolArgs := core.SetPool{
		TokenConfigCid: types.CONTRACT_ID(tokenConfigCid),
		InstrumentId:   instrumentId,
		TokenPool:      pool,
		Caller:         types.PARTY(owner),
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(core.TokenAdminRegistry{}),
						ContractId:     tarCid,
						Choice:         "SetPool",
						ChoiceArgument: ledger.MapToValue(setPoolArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	// Return the new TokenConfig contract ID
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TokenConfig" {
			return created.GetContractId()
		}
	}
	t.Fatal("no TokenConfig contract created after SetPool")

	return ""
}

// exerciseProposeAndAcceptAdmin sets up a token config entry with an admin for the given instrumentId.
// It exercises ProposeAdministrator (creates the config entry with pendingAdmin) then AcceptAdminRole (promotes to admin).
func exerciseProposeAndAcceptAdmin(
	t *testing.T,
	participant canton.Participant,
	owner string,
	tarCid string,
	instrumentId splice.InstrumentId,
) (string, string) {
	t.Helper()

	proposeArgs := core.ProposeAdministrator{
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
						TemplateId:     contracts.IdentifierFromBinding(core.TokenAdminRegistry{}),
						ContractId:     tarCid,
						Choice:         "ProposeAdministrator",
						ChoiceArgument: ledger.MapToValue(proposeArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	var newTarCid string
	var tokenConfigCid string
	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created == nil {
			continue
		} else if created.GetTemplateId().GetEntityName() == "TokenAdminRegistry" {
			newTarCid = created.GetContractId()
		} else if created.GetTemplateId().GetEntityName() == "TokenConfig" {
			tokenConfigCid = created.GetContractId()
		}
	}
	require.NotEmpty(t, newTarCid, "no TokenAdminRegistry contract created after ProposeAdministrator")
	require.NotEmpty(t, tokenConfigCid, "no TokenConfig contract created after ProposeAdministrator")

	acceptArgs := core.AcceptAdminRole{
		TokenConfigCid: types.CONTRACT_ID(tokenConfigCid),
		InstrumentId:   instrumentId,
		Caller:         types.PARTY(owner),
	}

	res, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(core.TokenAdminRegistry{}),
						ContractId:     newTarCid,
						Choice:         "AcceptAdminRole",
						ChoiceArgument: ledger.MapToValue(acceptArgs),
					},
				},
			}},
			ActAs: []string{owner},
		},
	})
	require.NoError(t, err)

	for _, event := range res.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TokenConfig" {
			return newTarCid, created.GetContractId()
		}
	}
	t.Fatal("no TokenConfig contract created after AcceptAdminRole")

	return "", ""
}

// executeScheduledBatchReturningTargetCid executes a scheduled batch and returns the new CID of a target contract.
func executeScheduledBatchReturningTargetCid(
	t *testing.T,
	participant canton.Participant,
	owner string,
	mcmsCid string,
	opID string,
	calls []mcmsApi.TimelockCall,
	predecessor string,
	salt string,
	targetCids map[string]string,
	targetEntityName string,
) string {
	t.Helper()

	executeArgs := mcmsCore.ExecuteScheduledBatch{
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
						TemplateId:     contracts.IdentifierFromBinding(mcmsCore.MCMS{}),
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
			if created.GetTemplateId().GetEntityName() == targetEntityName && created.GetPackageName() == core.PackageName {
				return created.GetContractId()
			}
		}
	}
	t.Fatalf("no %s contract created after ExecuteScheduledBatch", targetEntityName)

	return ""
}

// queryTokenConfig queries the standalone TokenConfig for an instrument and registry owner.
func queryTokenConfig(
	t *testing.T,
	participant canton.Participant,
	owner string,
	instrumentId splice.InstrumentId,
) *core.TokenConfig {
	t.Helper()

	contracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, contracts.IdentifierFromBinding(core.TokenAdminRegistry{}))
	require.NoError(t, err)

	expectedInstanceID := encodeInstrumentId(instrumentId)
	for _, contract := range contracts {
		createdEvent := contract.GetCreatedEvent()
		config, err := bindings.UnmarshalCreatedEvent[core.TokenConfig](createdEvent)
		require.NoError(t, err)
		if string(config.InstanceId) != expectedInstanceID {
			continue
		}
		if string(config.RegistryOwner) != owner {
			continue
		}
		if config.InstrumentId != instrumentId {
			continue
		}

		return config
	}

	return nil
}

func makeTokenConfigInstanceAddr(owner string, instrumentId splice.InstrumentId) string {
	return fmt.Sprintf("%s@%s", encodeInstrumentId(instrumentId), owner)
}

// encodeInstrumentId matches Canton's CCIP.MessageCodecV1.encodeInstrumentId
// which computes: keccak256(toHex(id <> "@" <> partyToText admin))
func encodeInstrumentId(instrumentId splice.InstrumentId) string {
	combined := string(instrumentId.Id) + "@" + string(instrumentId.Admin)
	hash := crypto.Keccak256([]byte(combined))

	return hex.EncodeToString(hash)
}
