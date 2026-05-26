package changesets

import (
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	mcms_bindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	mcms_ops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TestConfigureGlobalConfig_DirectExecution deploys all chain contracts, then exercises
// ConfigureGlobalConfig with TimelockConfig==nil (direct execution) and verifies the
// dest chain config update is applied without producing any MCMS proposals.
func TestConfigureGlobalConfig_DirectExecution(t *testing.T) {
	t.Parallel()

	cantonChain, bundle, env := setupCantonEnv(t)
	participant := cantonChain.Participants[0]
	party := participant.PartyID

	uploadChainContractDARs(t, participant)

	chainSelector := types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10))
	deployOut, err := cld_ops.ExecuteSequence(bundle, sequences.DeployChainContracts, *cantonChain, sequences.DeployChainContractsParams{
		CCIPOwnerParty: party,
		RMNOwnerParty:  party,
		CommitteeVerifiers: []sequences.CommitteeVerifierParams{{
			Template: ccvs.CommitteeVerifier{
				Owner:                        types.PARTY(party),
				CcipOwner:                    types.PARTY(party),
				VersionTag:                   "test-v1",
				StorageLocations:             []types.TEXT{"ipfs://test"},
				StorageLocationsAdmin:        types.PARTY(party),
				PendingStorageLocationsAdmin: types.PARTY(party),
			},
		}},
		GlobalConfig: sequences.GlobalConfigParams{
			Template: common.GlobalConfig{
				ChainSelector: chainSelector,
			},
		},
		RMNRemote: sequences.RMNRemoteParams{
			Template: rmn.RMNRemote{},
		},
		NativeInstrumentId: splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(party),
			Id:    "LINK",
		},
	})
	require.NoError(t, err, "deploy chain contracts")

	gcAddr, ccvAddr := extractAddresses(t, deployOut.Output.Addresses)

	out, err := ConfigureGlobalConfig{}.Apply(*env, CantonCSDeps[ConfigureGlobalConfigConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: ConfigureGlobalConfigConfig{
			InstanceAddress: gcAddr.InstanceAddress(),
			DestChainUpdates: []common.DestChainConfigArgs{
				makeDestChainConfig("999", ccvAddr),
			},
		},
	})
	require.NoError(t, err)
	assert.Empty(t, out.MCMSTimelockProposals, "direct execution should not produce proposals")
}

// TestConfigureGlobalConfig_MCMSProposal deploys a GlobalConfig and MCMS contract, then exercises
// ConfigureGlobalConfig with TimelockConfig set. Instead of executing on-chain, the changeset
// encodes the operation and returns an MCMSTimelockProposal with the correct structure.
func TestConfigureGlobalConfig_MCMSProposal(t *testing.T) {
	t.Parallel()

	cantonChain, bundle, env := setupCantonEnv(t)
	participant := cantonChain.Participants[0]
	party := participant.PartyID

	uploadDARs(t, participant, contracts.CCIPCommon, contracts.MCMS)

	gcAddrRef := deployGlobalConfig(t, bundle, *cantonChain, party)
	gcRawAddr, err := contracts.RawInstanceAddressFromString(gcAddrRef.Labels.List()[0])
	require.NoError(t, err)

	mcmsAddrRef := deployMCMSContract(t, bundle, *cantonChain, party)
	mcmsRawAddr, err := contracts.RawInstanceAddressFromString(mcmsAddrRef.Labels.List()[0])
	require.NoError(t, err)

	out, err := ConfigureGlobalConfig{}.Apply(*env, CantonCSDeps[ConfigureGlobalConfigConfig]{
		ChainSelector: chainsel.CANTON_LOCALNET.Selector,
		Participant:   0,
		Config: ConfigureGlobalConfigConfig{
			InstanceAddress:    gcRawAddr.InstanceAddress(),
			RawInstanceAddress: gcRawAddr.String(),
			DestChainUpdates: []common.DestChainConfigArgs{{
				DestChainSelector:         "999",
				IsEnabled:                 true,
				AddressBytesLength:        20,
				TokenReceiverAllowed:      true,
				BaseExecutionGasCost:      21000,
				OffRampAddress:            "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
				MessageNetworkFeeUSDCents: "100",
				TokenNetworkFeeUSDCents:   "50",
			}},
			TimelockConfig: &MCMSTimelockConfig{
				MinDelay:         10 * time.Minute,
				Description:      "test configure global config via MCMS",
				OverridePrevRoot: true,
				Action:           mcms_types.TimelockActionSchedule,
				MCMSContract: cantonmcms.MCMSContractInfo{
					RawInstanceAddress: mcmsRawAddr,
					InstanceAddress:    mcmsRawAddr.InstanceAddress(),
				},
				Role: cantonsdk.TimelockRoleProposer,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, out.MCMSTimelockProposals, 1, "MCMS path should produce exactly one proposal")

	proposal := out.MCMSTimelockProposals[0]
	assert.Equal(t, "v1", proposal.Version)
	assert.Equal(t, "test configure global config via MCMS", proposal.Description)
	assert.Equal(t, mcms_types.TimelockActionSchedule, proposal.Action)
	require.Len(t, proposal.Operations, 1)

	batchOp := proposal.Operations[0]
	assert.Equal(t, mcms_types.ChainSelector(chainsel.CANTON_LOCALNET.Selector), batchOp.ChainSelector)
	require.Len(t, batchOp.Transactions, 1, "one dest chain update should produce one transaction")

	tx := batchOp.Transactions[0]
	assert.Equal(t, gcRawAddr.InstanceAddress().Hex(), tx.To)
	assert.NotEmpty(t, tx.Data, "encoded operation data should not be empty")

	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	assert.Equal(t, gcRawAddr.String(), af.TargetInstanceAddress)
	assert.Equal(t, "ApplyDestChainConfigUpdates", af.FunctionName)
	assert.NotEmpty(t, af.OperationData)
}

// setupCantonEnv creates a CTF Canton network and returns the chain, operations bundle, and environment.
func setupCantonEnv(t *testing.T) (*canton.Chain, cld_ops.Bundle, *cldf.Environment) {
	t.Helper()

	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: 1,
		Once:               &sync.Once{},
	}).Initialize(t.Context())
	require.NoError(t, err)

	cantonChain := bc.(*canton.Chain)

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), reporter)

	env := &cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{bc}),
		OperationsBundle: bundle,
	}

	return cantonChain, bundle, env
}

// uploadDARs uploads the specified contract packages to the participant node.
func uploadDARs(t *testing.T, participant canton.Participant, packages ...contracts.Package) {
	t.Helper()

	// Upload one DAR per request: after DAR consolidation, the combined request can exceed
	// Canton admin API's default 10 MiB gRPC message limit even though each DAR is valid.
	for _, pkg := range packages {
		dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
		require.NoError(t, err, "failed to get DAR for %s", pkg)

		_, err = participant.AdminServices.Package.UploadDar(t.Context(), &participantv30.UploadDarRequest{
			Dars: []*participantv30.UploadDarRequest_UploadDarData{
				{Bytes: dar},
			},
			VetAllPackages:     true,
			SynchronizeVetting: true,
		})
		require.NoError(t, err, "failed to upload DAR for %s", pkg)
	}
}

// uploadChainContractDARs uploads all CCIP chain contract DARs needed by DeployChainContracts.
func uploadChainContractDARs(t *testing.T, participant canton.Participant) {
	t.Helper()
	uploadDARs(t, participant,
		contracts.CCIPCommon,
		contracts.CCIPOffRamp,
		contracts.CCIPOnRamp,
		contracts.CCIPTokenAdminRegistry,
		contracts.CCIPCommitteeVerifier,
		contracts.CCIPLockReleaseTokenPool,
		contracts.CCIPPerPartyRouter,
		contracts.CCIPRMN,
	)
}

func deployGlobalConfig(t *testing.T, bundle cld_ops.Bundle, chain canton.Chain, party string) datastore.AddressRef {
	t.Helper()

	chainSelector := types.NUMERIC(strconv.FormatUint(chainsel.CANTON_LOCALNET.Selector, 10))
	result, err := cld_ops.ExecuteOperation(bundle, global_config.Deploy, chain, opcontract.DeployInput[common.GlobalConfig]{
		Template: common.GlobalConfig{
			CcipOwner:     types.PARTY(party),
			ChainSelector: chainSelector,
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "failed to deploy GlobalConfig")
	require.NotEmpty(t, result.Output.Address)

	return result.Output
}

func deployMCMSContract(t *testing.T, bundle cld_ops.Bundle, chain canton.Chain, party string) datastore.AddressRef {
	t.Helper()

	emptyRoleState := mcms_bindings.RoleState{
		Config: mcms_bindings.MultisigConfig{
			Signers:      nil,
			GroupQuorums: make([]types.INT64, 32),
			GroupParents: make([]types.INT64, 32),
		},
		SeenHashes: map[types.TEXT]types.TIMESTAMP{},
		ExpiringRoot: mcms_bindings.ExpiringRoot{
			Root:    "",
			OpCount: 0,
		},
		RootMetadata: mcms_bindings.RootMetadata{},
	}

	result, err := cld_ops.ExecuteOperation(bundle, mcms_ops.Deploy, chain, opcontract.DeployInput[mcms_bindings.MCMS]{
		Template: mcms_bindings.MCMS{
			Owner:              types.PARTY(party),
			ChainId:            1,
			Proposer:           emptyRoleState,
			Canceller:          emptyRoleState,
			Bypasser:           emptyRoleState,
			BlockedFunctions:   nil,
			TimelockTimestamps: map[types.TEXT]types.TIMESTAMP{},
		},
		OwnerParty: types.PARTY(party),
	})
	require.NoError(t, err, "failed to deploy MCMS contract")
	require.NotEmpty(t, result.Output.Address)

	return result.Output
}

// extractAddresses finds the GlobalConfig and CommitteeVerifier addresses from the
// deploy chain contracts output.
func extractAddresses(t *testing.T, addresses []datastore.AddressRef) (contracts.RawInstanceAddress, mcms_bindings.RawInstanceAddress) {
	t.Helper()

	var gcRaw contracts.RawInstanceAddress
	var ccvRaw mcms_bindings.RawInstanceAddress

	for _, addr := range addresses {
		raw, err := contracts.RawInstanceAddressFromString(addr.Labels.List()[0])
		require.NoError(t, err)

		switch string(addr.Type) {
		case "CantonGlobalConfig":
			gcRaw = raw
		case "CommitteeVerifier":
			ccvRaw = mcms_bindings.RawInstanceAddress{Unpack: types.TEXT(raw.String())}
		}
	}

	require.NotEmpty(t, gcRaw.String(), "GlobalConfig address not found in deploy output")
	require.NotEmpty(t, ccvRaw.Unpack, "CommitteeVerifier address not found in deploy output")

	return gcRaw, ccvRaw
}

func makeDestChainConfig(destChainSelector string, ccvAddr mcms_bindings.RawInstanceAddress) common.DestChainConfigArgs {
	return common.DestChainConfigArgs{
		DestChainSelector:         types.NUMERIC(destChainSelector),
		IsEnabled:                 true,
		AddressBytesLength:        20,
		TokenReceiverAllowed:      true,
		BaseExecutionGasCost:      21000,
		OffRampAddress:            "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		DefaultCCVs:               []mcms_bindings.RawInstanceAddress{ccvAddr},
		MessageNetworkFeeUSDCents: "100",
		TokenNetworkFeeUSDCents:   "50",
	}
}
