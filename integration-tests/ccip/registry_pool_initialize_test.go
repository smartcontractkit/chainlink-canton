package tests

import (
	"context"
	"encoding/hex"
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/model"
	damlledger "github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	registryrl "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestInitializeBurnMintTokenPool exercises the Initialize choice on the registry-pools
// BurnMintTokenPool end to end against a real Canton participant: deploy the pool (Create),
// deploy a TokenAdminRegistry, then exercise Initialize with one lane. Initialize should
// propose+accept the admin role, deploy the lane's three rate limiters, wire the lane via
// ApplyChainUpdates, and register the pool with TAR - all in one choice.
//
// This is the self-issued path: instrumentId.admin == admin, so ProposeAdministrator's
// caller == instrumentId.admin check passes and a single actAs party can authorize the
// whole choice (controller poolOwner, admin - here poolOwner == admin).
func TestInitializeBurnMintTokenPool(t *testing.T) {
	t.Parallel()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)
	participant := cantonChain.Participants[0]
	party := participant.PartyID
	ctx := t.Context()

	observerParty := testhelpers.AllocateParty(t, participant, "observer")

	// Upload DARs. The registry-pools packages declare CCIP core and the registry
	// rate-limiter as data-dependencies, but upload them explicitly too so vetting
	// doesn't depend on how the build happened to bundle them.
	coreDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
	require.NoError(t, err)
	rateLimiterDar, err := contracts.GetDar(contracts.CCIPRegistryRateLimiterV2, contracts.DevVersion)
	require.NoError(t, err)
	poolDar, err := contracts.GetDar(contracts.CCIPRegistryBurnMintTokenPoolV2, contracts.DevVersion)
	require.NoError(t, err)
	_, err = participant.AdminServices.Package.UploadDar(ctx, &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: coreDar},
			{Bytes: rateLimiterDar},
			{Bytes: poolDar},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	require.NoError(t, err, "failed to upload dar files")

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), reporter)

	// Deploy TAR - Initialize registers the pool against this instance.
	tarAddrRef, err := cld_ops.ExecuteOperation(bundle, token_admin_registry.Deploy, *cantonChain, contract.DeployInput[core.TokenAdminRegistry]{
		Template: core.TokenAdminRegistry{
			CcipOwner:  types.PARTY(party),
			InstanceId: "",
			EntryCount: 0,
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "deploy TAR")
	tarInstanceAddress := contracts.HexToInstanceAddress(tarAddrRef.Output.Address)
	tarCID, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, []string{party}, core.TokenAdminRegistry{}.GetTemplateID(), tarInstanceAddress)
	require.NoError(t, err, "find TAR contract ID")

	// Self-issued instrument: instrumentId.admin must equal admin for
	// ProposeAdministrator's caller == instrumentId.admin check to pass.
	instrumentId := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    "E2ETOKEN",
	}

	poolInstanceID := "e2e-bnm-pool"
	dummyAddr := func(label string) chainlinkapi.RawInstanceAddress {
		return contracts.NewRawInstanceAddress(contracts.InstanceID(label), types.PARTY(party)).Binding()
	}

	pool := registrybnm.BurnMintTokenPool{
		InstanceId:         types.TEXT(poolInstanceID),
		PoolOwner:          types.PARTY(party),
		CcipOwner:          types.PARTY(party),
		InstrumentId:       instrumentId,
		Decimals:           10,
		Observers:          []types.PARTY{types.PARTY(observerParty)},
		PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{},
		TransferTimeout: registrybnm.TransferTimeout{
			Indefinite: &types.UNIT{},
		},
		Deps: registrybnm.BurnMintTokenPoolDeps{
			TokenAdminRegistry: dummyAddr("dummy-tar"),
			RmnRemote:          dummyAddr("dummy-rmn"),
			FeeQuoter:          dummyAddr("dummy-fq"),
		},
	}

	poolInstanceAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(poolInstanceID), types.PARTY(party)).InstanceAddress()

	lane := registrybnm.LaneDeploySpec{
		RemoteChainSelector: "1.0",
		RemotePools:         []types.TEXT{"1fb477ac89df394bead4f46ad754f8aec70cc0e4"},
		RemoteTokenAddress:  "cbae4ea0c4c503a582af009bb2b30b75badc1e32",
		FinalityConfig:      ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
		Inbound: registrybnm.RateLimiterDeploySpec{
			InstanceId: "e2e-lane-1-in",
			IsEnabled:  true,
			Capacity:   "1000.0",
			Rate:       "10.0",
		},
		Outbound: registrybnm.RateLimiterDeploySpec{
			InstanceId: "e2e-lane-1-out",
			IsEnabled:  true,
			Capacity:   "1000.0",
			Rate:       "10.0",
		},
		InboundCustomFinality: registrybnm.RateLimiterDeploySpec{
			InstanceId: "e2e-lane-1-in-custom",
			IsEnabled:  false,
			Capacity:   "0.0",
			Rate:       "0.0",
		},
	}

	initArgs := registrybnm.Initialize{
		TokenAdminRegistryCid:  types.CONTRACT_ID(tarCID),
		ExistingTokenConfigCid: nil,
		Admin:                  types.PARTY(party),
		Lanes:                  []registrybnm.LaneDeploySpec{lane},
	}
	_, err = submitCreateAndExercise(ctx, participant, party, pool, "Initialize", initArgs)
	require.NoError(t, err, "create pool and exercise Initialize atomically")

	// TokenConfig was registered: admin == poolOwner, tokenPool points at this pool.
	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).
		RawInstanceAddress(types.PARTY(party)).InstanceAddress()
	tokenConfig := findActiveContract[core.TokenConfig](t, ctx, participant, party, core.TokenConfig{}.GetTemplateID(), tokenConfigAddress)
	require.NotNil(t, tokenConfig.Admin)
	require.Equal(t, types.PARTY(party), *tokenConfig.Admin)
	require.Nil(t, tokenConfig.PendingAdmin)
	require.NotNil(t, tokenConfig.TokenPool)
	require.Equal(t, types.PARTY(party), tokenConfig.TokenPool.PoolOwner)
	require.Equal(t, types.TEXT(poolInstanceID), tokenConfig.TokenPool.PoolInstanceId)

	// The pool was re-created (ApplyChainUpdates archives+recreates) with the lane wired up.
	pool2 := findActiveContract[registrybnm.BurnMintTokenPool](t, ctx, participant, party, pool.GetTemplateID(), poolInstanceAddress)
	require.Equal(t, []types.PARTY{types.PARTY(observerParty)}, pool2.Observers)
	require.Len(t, pool2.RemoteChainConfigs, 1)

	// All three rate limiters for the lane were deployed with the observers set.
	for _, spec := range []registrybnm.RateLimiterDeploySpec{lane.Inbound, lane.Outbound, lane.InboundCustomFinality} {
		rlAddress := contracts.NewRawInstanceAddress(contracts.InstanceID(spec.InstanceId), types.PARTY(party)).InstanceAddress()
		rl := findActiveContract[registryrl.RateLimiter](t, ctx, participant, party, registryrl.RateLimiter{}.GetTemplateID(), rlAddress)
		require.Equal(t, types.TEXT(poolInstanceID), rl.PoolInstanceId)
		require.Equal(t, types.PARTY(party), rl.PoolOwner)
		require.Equal(t, []types.PARTY{types.PARTY(observerParty)}, rl.Observers)
		require.Equal(t, spec.IsEnabled, rl.IsEnabled)
		require.Equal(t, spec.Capacity, rl.Capacity)
		require.Equal(t, spec.Rate, rl.Rate)
	}
}

// creatable is satisfied by every generated binding template struct.
type creatable interface {
	model.CreateCommander
	GetTemplateID() string
}

// submitCreateAndExercise submits a single CreateAndExercise command: payload is
// created and choice is exercised on the resulting contract in one atomic
// transaction, authorized by actAs. This is the real usage pattern Initialize is
// designed for - a caller never needs to see the pool's contractId, since it's
// created and exercised on in the same command.
func submitCreateAndExercise(
	ctx context.Context,
	participant canton.Participant,
	actAs string,
	payload creatable,
	choice string,
	choiceArg any,
) (*apiv2.SubmitAndWaitForTransactionResponse, error) {
	createCmd := payload.CreateCommand()
	return participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			ActAs:     []string{actAs},
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_CreateAndExercise{
					CreateAndExercise: &apiv2.CreateAndExerciseCommand{
						TemplateId:      contracts.IdentifierFromBinding(payload),
						CreateArguments: damlledger.MapToRecord(createCmd.Arguments),
						Choice:          choice,
						ChoiceArgument:  damlledger.MapToValue(choiceArg),
					},
				},
			}},
		},
	})
}

func findActiveContract[T model.CreateCommander](
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party string,
	templateID string,
	addr contracts.InstanceAddress,
) *T {
	t.Helper()

	ac, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, []string{party}, templateID, addr)
	require.NoError(t, err)

	v, err := bindings.UnmarshalCreatedEvent[T](ac.GetCreatedEvent())
	require.NoError(t, err)

	return v
}
