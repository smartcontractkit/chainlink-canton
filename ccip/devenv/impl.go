package devenv

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	chainsel "github.com/smartcontractkit/chain-selectors"
	evmadapters "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/executor"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/lock_release_token_pool"
	dsutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/go-daml/pkg/auth"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/deployments"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	onrampbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	cantonadapters "github.com/smartcontractkit/chainlink-canton/ccip/devenv/adapters"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	canton_lock_release_token_pool "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var (
	_ cciptestinterfaces.CCIP17              = &Chain{}
	_ cciptestinterfaces.CCIP17Configuration = &Chain{}
	_ ccv.ImplFactory                        = &ImplFactory{}
)

type ImplFactory struct{}

func NewImplFactory() *ImplFactory {
	return &ImplFactory{}
}

// New implements [registry.ImplFactory].
func (i *ImplFactory) New(ctx context.Context, cfg *ccv.Cfg, lggr zerolog.Logger, env *deployment.Environment, bc *blockchain.Input) (cciptestinterfaces.CCIP17, error) {
	return New(ctx, cfg, lggr, env, bc.ChainID)
}

// NewEmpty implements [registry.ImplFactory].
func (i *ImplFactory) NewEmpty() cciptestinterfaces.CCIP17Configuration {
	return NewEmptyCCIP17Canton(
		log.
			Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.DebugLevel).
			With().
			Fields(map[string]any{"component": "Canton"}).
			Logger(),
	)
}

type Chain struct {
	e            *deployment.Environment
	chain        canton.Chain
	logger       zerolog.Logger
	chainDetails chainsel.ChainDetails

	nextSeq       uint64
	lastSentDest  uint64
	lastSentSeq   uint64
	lastSentEvent cciptestinterfaces.MessageSentEvent

	lastSentMessage               protocol.Message
	lastSentVerifierDestAddress   protocol.UnknownAddress
	lastSentVerifierBlob          []byte
	lastSentHasVerificationInputs bool
	cfg                           *ccv.Cfg
}

func (c *Chain) ChainSelector() uint64 {
	return c.chainDetails.ChainSelector
}

func New(ctx context.Context, cfg *ccv.Cfg, logger zerolog.Logger, e *deployment.Environment, chainID string) (*Chain, error) {
	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(chainID, chainsel.FamilyCanton)
	if err != nil {
		return nil, fmt.Errorf("get chain details for chain %s: %w", chainID, err)
	}
	chain := e.BlockChains.CantonChains()[chainDetails.ChainSelector]

	return &Chain{
		e:            e,
		chain:        chain,
		chainDetails: chainDetails,
		logger:       logger,
		cfg:          cfg,
	}, nil
}

func NewEmptyCCIP17Canton(logger zerolog.Logger) *Chain {
	return &Chain{
		logger: logger,
	}
}

// ChainFamily implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) ChainFamily() string {
	return chainsel.FamilyCanton
}

// ConfigureNodes implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) ConfigureNodes(ctx context.Context, blockchain *blockchain.Input) (string, error) {
	return "", nil // TODO: implement
}

// DeployContractsForSelector implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) DeployContractsForSelector(ctx context.Context, env *deployment.Environment, selector uint64, topology *deployments.EnvironmentTopology) (datastore.DataStore, error) {
	// Only using a single participant for now
	chain := env.BlockChains.CantonChains()[selector]
	participant := chain.Participants[0]

	l := c.logger
	l.Info().Msg("Configuring contracts for selector")
	l.Info().Any("Selector", selector).Msg("Deploying for chain selectors")
	runningDS := datastore.NewMemoryDataStore()

	l.Info().Uint64("Selector", selector).Msg("Configuring per-chain contracts bundle")
	bundle := operations.NewBundle(
		func() context.Context { return context.Background() },
		env.Logger,
		operations.NewMemoryReporter(),
	)
	env.OperationsBundle = bundle

	l.Info().Msg("Uploading and vetting CCIP DARs...")
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get common dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: commonDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload common dar file")
	}
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get offramp dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: offRampDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload offramp dar file")
	}
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get onramp dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: onRampDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload onramp dar file")
	}
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get token admin registry dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: tokenAdminRegistryDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload token admin registry dar file")
	}
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get committee verifier dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: committeeVerifierDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload committee verifier dar file")
	}
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get per-party router dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: perPartyRouterDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload per-party router dar file")
	}
	feeQuoterDar, err := contracts.GetDar(contracts.CCIPFeeQuoter, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee quoter dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: feeQuoterDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload fee quoter dar file")
	}
	rmnDar, err := contracts.GetDar(contracts.CCIPRMN, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get rmn dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: rmnDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload rmn dar file")
	}
	ccipSenderDar, err := contracts.GetDar(contracts.CCIPSender, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get ccip sender dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: ccipSenderDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload ccip sender dar file")
	}
	ccipExecutorDar, err := contracts.GetDar(contracts.CCIPExecutor, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get ccip executor dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: ccipExecutorDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload ccip executor dar file")
	}
	ccipTestDar, err := contracts.GetDar(contracts.CCIPTest, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get ccip test dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: ccipTestDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload ccip test dar file")
	}

	tokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get token pool dar file")
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile: tokenPoolDar,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload token pool dar file")
	}

	l.Info().Any("selector", selector).Any("party", participant.PartyID).Msg("Deploying chain contracts")
	config := cantonChangesets.CantonCSDeps[cantonChangesets.DeployChainContractsConfig]{
		ChainSelector: selector,
		Participant:   0,
		Config: cantonChangesets.DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty:     participant.PartyID,
				CommitteeVerifiers: nil,
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "", // Populated by the sequence
						ChainSelector: types.NUMERIC(strconv.FormatUint(selector, 10)),
					},
				},
				FeeQuoterConfig: sequences.FeeQuoterParams{
					Template: feequoter.FeeQuoter{
						PriceUpdaters: []types.PARTY{types.PARTY(participant.PartyID)},
					},
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						RmnOwner:       types.PARTY(participant.PartyID),
						CursedSubjects: nil,
					},
				},
			},
		},
	}

	remoteChainFeeConfigs := types.GENMAP{}
	for _, bc := range env.BlockChains.All() {
		sel := bc.ChainSelector()
		if sel == selector {
			continue
		}
		remoteChainFeeConfigs[strconv.FormatUint(sel, 10)] = ccvs.CCVFeeConfig{
			FeeUSDCents:        types.NUMERIC("0"),
			GasForVerification: types.INT64(0),
			PayloadSizeBytes:   types.INT64(0),
		}.ToMap()
	}

	// Get committees
	for qualifier, committeeConfig := range topology.NOPTopology.Committees {
		storageLocations := make([]types.TEXT, len(committeeConfig.StorageLocations))
		for i, location := range committeeConfig.StorageLocations {
			storageLocations[i] = types.TEXT(location)
		}
		cv := sequences.CommitteeVerifierParams{
			Qualifier: qualifier,
			Template: ccvs.CommitteeVerifier{
				Owner:                        types.PARTY(participant.PartyID), // TODO: use different ccv owner?
				CcipOwner:                    types.PARTY(participant.PartyID),
				VersionTag:                   types.TEXT("49ff34ed"),
				MessageSentObserver:          types.PARTY(participant.PartyID),
				StorageLocations:             storageLocations,
				StorageLocationsAdmin:        types.PARTY(participant.PartyID),
				PendingStorageLocationsAdmin: types.PARTY(participant.PartyID),
				Deps:                         ccvs.CommitteeVerifierDeps{}, // Set by sequence
				// MUST be a real GENMAP, not a Go map.
				RemoteChainFeeConfigs: remoteChainFeeConfigs,
			},
		}
		config.Config.Params.CommitteeVerifiers = append(config.Config.Params.CommitteeVerifiers, cv)
	}

	out, err := cantonChangesets.DeployChainContracts{}.Apply(*env, config)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy chain contracts for selector %d: %w", selector, err)
	}
	err = runningDS.Merge(out.DataStore.Seal())
	if err != nil {
		return nil, err
	}

	const defaultLockReleaseQualifier = "TEST (LockReleaseTokenPool 1.7.0 [default] to BurnMintTokenPool 1.7.0 [default])"

	// Use registry admin for Amulet instrument so sender holdings and pool instrument align.
	lockPoolInstrumentAdmin := participant.PartyID
	if registryAdmin, adminErr := resolveRegistryAdmin(ctx, participant); adminErr == nil && registryAdmin != "" {
		lockPoolInstrumentAdmin = registryAdmin
	}

	// Deploy and TAR-register the default lock/release pool for Canton source token transfers.
	relativeHours := types.INT64(24)
	lockPoolInstrumentID := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(lockPoolInstrumentAdmin),
		Id:    types.TEXT("Amulet"),
	}
	lockPoolReceiveContext := common.CCIPContext{
		Values: types.TEXTMAP{},
	}
	lockPoolTransferTimeout := lockreleasetokenpool.TransferTimeout{
		RelativeHours: &relativeHours,
	}
	lockPoolQualifier := defaultLockReleaseQualifier
	tokenAdminRegistryRef, tarErr := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	))
	if tarErr != nil {
		return nil, fmt.Errorf("failed to get token admin registry for token pool deployment: %w", tarErr)
	}
	rmnRemoteRef, rmnErr := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(rmn_remote.ContractType),
		rmn_remote.Version,
		"",
	))
	if rmnErr != nil {
		return nil, fmt.Errorf("failed to get rmn remote for token pool deployment: %w", rmnErr)
	}
	feeQuoterRef, fqErr := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(fee_quoter.ContractType),
		fee_quoter.Version,
		"",
	))
	if fqErr != nil {
		return nil, fmt.Errorf("failed to get fee quoter for token pool deployment: %w", fqErr)
	}
	var tokenAdminRegistryRawAddr common.RawInstanceAddress
	if labels := tokenAdminRegistryRef.Labels.List(); len(labels) > 0 {
		if rawAddr, parseErr := contracts.RawInstanceAddressFromString(labels[0]); parseErr == nil {
			tokenAdminRegistryRawAddr = rawAddr.Binding()
		}
	}
	var rmnRemoteRawAddr common.RawInstanceAddress
	if labels := rmnRemoteRef.Labels.List(); len(labels) > 0 {
		if rawAddr, parseErr := contracts.RawInstanceAddressFromString(labels[0]); parseErr == nil {
			rmnRemoteRawAddr = rawAddr.Binding()
		}
	}
	var feeQuoterRawAddr common.RawInstanceAddress
	if labels := feeQuoterRef.Labels.List(); len(labels) > 0 {
		if rawAddr, parseErr := contracts.RawInstanceAddressFromString(labels[0]); parseErr == nil {
			feeQuoterRawAddr = rawAddr.Binding()
		}
	}
	lockPoolChangeset, deployErr := (cantonChangesets.DeployTokenPool{}).Apply(*env, cantonChangesets.CantonCSDeps[cantonChangesets.DeployTokenPoolConfig]{
		ChainSelector: selector,
		Participant:   0,
		Config: cantonChangesets.DeployTokenPoolConfig{
			CcipOwner:                         participant.PartyID,
			PoolOwner:                         participant.PartyID,
			InstrumentId:                      lockPoolInstrumentID,
			Decimals:                          18,
			Qualifier:                         lockPoolQualifier,
			PoolReceiveContext:                lockPoolReceiveContext,
			TransferTimeout:                   lockPoolTransferTimeout,
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
			TokenAdminRegistryRawAddress:      tokenAdminRegistryRawAddr,
			RmnRemoteRawAddress:               rmnRemoteRawAddr,
			FeeQuoterRawAddress:               feeQuoterRawAddr,
		},
	})
	if deployErr != nil {
		return nil, fmt.Errorf("failed to deploy canton lock/release pool changeset: %w", deployErr)
	}
	if mergeErr := runningDS.Merge(lockPoolChangeset.DataStore.Seal()); mergeErr != nil {
		return nil, fmt.Errorf("failed to merge deployed canton lock/release pool datastore: %w", mergeErr)
	}
	// ccv up token-transfer setup resolves the Canton remote pool by the canonical
	// LockReleaseTokenPool type + qualifier. Alias the deployed CantonLockReleaseTokenPool
	// ref to that key so remote pool lookup succeeds during EVM lane configuration.
	lockPoolRef, getErr := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(canton_lock_release_token_pool.ContractType),
		canton_lock_release_token_pool.Version,
		lockPoolQualifier,
	))
	if getErr == nil && lockPoolRef.Address != "" {
		_ = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
			Address:       lockPoolRef.Address,
			Labels:        lockPoolRef.Labels,
			Type:          datastore.ContractType("LockReleaseTokenPool"),
			Version:       lock_release_token_pool.Version,
			Qualifier:     lockPoolQualifier,
			ChainSelector: selector,
		})
	} else {
		l.Warn().Err(getErr).Msg("Lock/release pool alias was not added because no deployed pool address ref was found")
	}

	l.Info().Msg("Deployed and registered lock/release pool via DeployTokenPool changeset")
	// Pre-seed real AMT holdings for the lock/release pool owner so e2e token flows use
	// existing liquidity (same model as ccip_execute_token_test.go) and never mint at execute time.
	if seedErr := seedAMTLiquidity(ctx, participant, participant.PartyID, "1000000.00"); seedErr != nil {
		return nil, fmt.Errorf("failed to pre-seed AMT liquidity for lock/release pool owner: %w", seedErr)
	}

	// ccv up expects datastore refs for all destination token-pool combos during lane config.
	// Canton does not deploy all BurnMint pool variants, so we add synthetic refs for the missing ones.
	// Do not add a synthetic LockReleaseTokenPool ref: the real lock/release pool is already deployed.
	for i, combo := range devenvcommon.AllTokenCombinations() {
		addressRef := combo.DestPoolAddressRef()
		if addressRef.Type == datastore.ContractType(lock_release_token_pool.ContractType) {
			continue
		}
		_ = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
			Address:       contracts.MustNewInstanceID("dst-token-pool-" + strconv.Itoa(i)).RawInstanceAddress(types.PARTY(participant.PartyID)).InstanceAddress().Hex(),
			Type:          addressRef.Type,
			Version:       addressRef.Version,
			Qualifier:     addressRef.Qualifier,
			ChainSelector: selector,
		})
	}

	// Add executor refs, storing raw instance addresses as labels so that
	// SendMessage can recover the instanceId needed to match the executor service.
	executorRawAddr := contracts.MustNewInstanceID("executor-1").RawInstanceAddress(types.PARTY(participant.PartyID))
	err = runningDS.AddressRefStore.Add(datastore.AddressRef{
		Address:       executorRawAddr.InstanceAddress().Hex(),
		Labels:        datastore.NewLabelSet(executorRawAddr.String()),
		Type:          datastore.ContractType(executor.ContractType),
		Version:       executor.Version,
		Qualifier:     devenvcommon.DefaultExecutorQualifier,
		ChainSelector: selector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add executor address ref: %w", err)
	}
	executorProxyRawAddr := contracts.MustNewInstanceID("executor-proxy-1").RawInstanceAddress(types.PARTY(participant.PartyID))
	err = runningDS.AddressRefStore.Add(datastore.AddressRef{
		Address:       executorProxyRawAddr.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(executorProxyRawAddr.String()),
		Type:          datastore.ContractType(executor.ProxyType),
		Version:       executor.Version,
		Qualifier:     devenvcommon.DefaultExecutorQualifier,
		ChainSelector: selector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add executor proxy address ref: %w", err)
	}

	env.DataStore = runningDS.Seal()

	return runningDS.Seal(), nil
}

// ConnectContractsWithSelectors implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) ConnectContractsWithSelectors(ctx context.Context, env *deployment.Environment, selector uint64, remoteSelectors []uint64, committees *deployments.EnvironmentTopology) error {
	l := c.logger
	l.Info().Uint64("FromSelector", selector).Any("ToSelectors", remoteSelectors).Msg("Connecting contracts with selectors")
	bundle := operations.NewBundle(
		func() context.Context { return context.Background() },
		env.Logger,
		operations.NewMemoryReporter(),
	)
	env.OperationsBundle = bundle

	formatFunc := func(ref datastore.AddressRef) (contracts.InstanceAddress, error) {
		return contracts.HexToInstanceAddress(ref.Address), nil
	}

	// Get InstanceAddresses of Canton contracts
	globalConfig, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
		Type: datastore.ContractType(global_config.ContractType),
	}, selector, formatFunc)
	if err != nil {
		return fmt.Errorf("failed to get global config address for chain %d: %w", selector, err)
	}
	feeQuoter, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
		Type: datastore.ContractType(fee_quoter.ContractType),
	}, selector, formatFunc)
	if err != nil {
		return fmt.Errorf("failed to get fee quoter address for chain %d: %w", selector, err)
	}
	onRamp, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
		Type: datastore.ContractType(onramp.ContractType),
	}, selector, formatFunc)
	if err != nil {
		return fmt.Errorf("failed to get on ramp address for chain %d: %w", selector, err)
	}
	offRamp, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
		Type: datastore.ContractType(offramp.ContractType),
	}, selector, formatFunc)
	if err != nil {
		return fmt.Errorf("failed to get off ramp address for chain %d: %w", selector, err)
	}

	config := cantonChangesets.CantonCSDeps[cantonChangesets.ConfigureChainForLanesConfig]{
		ChainSelector: selector,
		Participant:   0,
		Config: cantonChangesets.ConfigureChainForLanesConfig{
			Input: sequences.ConfigureChainForLanesInput{
				ChainSelector:      selector,
				GlobalConfig:       globalConfig,
				FeeQuoter:          feeQuoter,
				OnRamp:             onRamp,
				OffRamp:            offRamp,
				CommitteeVerifiers: nil,
				RemoteChains:       make(map[uint64]adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress], len(remoteSelectors)),
			},
		},
	}

	// Configure outbound defaults: use the default committee verifier as both the outbound CCV and executor,
	// matching the reference integration test.
	var committeeVerifierRawAddr contracts.RawInstanceAddress
	ccvRef, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(committee_verifier.ContractType),
		committee_verifier.Version,
		devenvcommon.DefaultCommitteeVerifierQualifier,
	))
	if err == nil && len(ccvRef.Labels.List()) > 0 {
		committeeVerifierRawAddr = contracts.RawInstanceAddress(ccvRef.Labels.List()[0])
	}

	for _, remoteSelector := range remoteSelectors {
		// TODO: should be moved to the ChainFamily interface.
		var addressBytesLength uint8
		family, err := chainsel.GetSelectorFamily(remoteSelector)
		if err != nil {
			return fmt.Errorf("failed to get selector family for chain %d: %w", remoteSelector, err)
		}
		var chainFamily adapters.ChainFamily
		switch family {
		case chainsel.FamilyEVM:
			addressBytesLength = 20
			chainFamily = &evmadapters.ChainFamilyAdapter{}
		case chainsel.FamilyCanton:
			addressBytesLength = 32
			chainFamily = cantonadapters.NewChainFamilyAdapter(&evmadapters.ChainFamilyAdapter{})
		default:
			return fmt.Errorf("unsupported family %s for chain %d", family, remoteSelector)
		}

		remoteOnRamp, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
			Type:    datastore.ContractType(onramp.ContractType),
			Version: onramp.Version,
		}, remoteSelector, chainFamily.AddressRefToBytes)
		if err != nil {
			return fmt.Errorf("failed to get on ramp address for remote chain %d: %w", remoteSelector, err)
		}
		remoteOffRamp, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
			Type:    datastore.ContractType(offramp.ContractType),
			Version: offramp.Version,
		}, remoteSelector, chainFamily.AddressRefToBytes)
		if err != nil {
			return fmt.Errorf("failed to get off ramp address for remote chain %d: %w", remoteSelector, err)
		}
		localExecutorRef, err := env.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				selector,
				datastore.ContractType(executor.ProxyType),
				executor.Version,
				devenvcommon.DefaultExecutorQualifier,
			),
		)
		if err != nil {
			return fmt.Errorf("failed to get default executor for source chain %d: %w", selector, err)
		}
		// Normalize executor to 32 bytes (left-padded) for Canton config encoding.
		normalizedSourceExecutor := contracts.HexToInstanceAddress(localExecutorRef.Address).Hex()
		remoteChainConfig := adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
			AllowTrafficFrom:         true,
			OnRamps:                  [][]byte{remoteOnRamp},
			OffRamp:                  remoteOffRamp,
			DefaultInboundCCVs:       []contracts.RawInstanceAddress{committeeVerifierRawAddr},
			LaneMandatedInboundCCVs:  nil,
			DefaultOutboundCCVs:      []contracts.RawInstanceAddress{committeeVerifierRawAddr},
			LaneMandatedOutboundCCVs: nil,
			DefaultExecutor:          contracts.RawInstanceAddress(normalizedSourceExecutor),
			FeeQuoterDestChainConfig: adapters.FeeQuoterDestChainConfig{NetworkFeeUSDCents: 0, DefaultTokenFeeUSDCents: 0},
			ExecutorDestChainConfig:  adapters.ExecutorDestChainConfig{},
			AddressBytesLength:       addressBytesLength,
			BaseExecutionGasCost:     0,
		}
		config.Config.Input.RemoteChains[remoteSelector] = remoteChainConfig
	}

	for qualifier, committee := range committees.NOPTopology.Committees {
		// Get CommitteeVerifier address for this qualifier
		committeeVerifier, err := dsutils.FindAndFormatRef(env.DataStore, datastore.AddressRef{
			Type:      datastore.ContractType(committee_verifier.ContractType),
			Qualifier: qualifier,
		}, selector, formatFunc)
		if err != nil {
			return fmt.Errorf("failed to get committee verifier address with qualifier %s for chain %d: %w", qualifier, selector, err)
		}

		committeeVerifierConfig := adapters.CommitteeVerifierConfig[contracts.InstanceAddress]{
			CommitteeVerifier: []contracts.InstanceAddress{committeeVerifier},
			RemoteChains:      make(map[uint64]adapters.CommitteeVerifierRemoteChainConfig),
		}

		// Configure all remote chains with the respective signers
		for _, remoteSelector := range remoteSelectors {
			chainCfg, ok := committee.ChainConfigs[strconv.FormatUint(remoteSelector, 10)]
			if !ok {
				return fmt.Errorf("chain selector %d not found in committee %q", remoteSelector, qualifier)
			}
			// For each of the NOPs in this committee, get their Canton-family signer.
			// Since the Canton CommitteeVerifier requires the (uncompressed) signer pubkey to be set on-chain,
			// nop.SignerAddressByFamily[chainsel.FamilyCanton] must contain the signer's pubkey, NOT address
			signers := make([]string, 0, len(chainCfg.NOPAliases))
			for _, alias := range chainCfg.NOPAliases {
				nop, ok := committees.NOPTopology.GetNOP(alias)
				if !ok {
					return fmt.Errorf("NOP alias %q not found for committee %q chain %d", alias, qualifier, remoteSelector)
				}
				signer, ok := nop.SignerAddressByFamily[chainsel.FamilyCanton]
				if !ok {
					return fmt.Errorf("no Canton pubkey signer found for NOP alias %q", alias)
				}
				signers = append(signers, signer)
			}
			committeeVerifierConfig.RemoteChains[remoteSelector] = adapters.CommitteeVerifierRemoteChainConfig{
				AllowlistEnabled:          false,
				AddedAllowlistedSenders:   nil,
				RemovedAllowlistedSenders: nil,
				FeeUSDCents:               0,
				GasForVerification:        0,
				PayloadSizeBytes:          0,
				SignatureConfig: adapters.CommitteeVerifierSignatureQuorumConfig{
					Signers:   signers,
					Threshold: chainCfg.Threshold,
				},
			}
		}
		config.Config.Input.CommitteeVerifiers = append(config.Config.Input.CommitteeVerifiers, committeeVerifierConfig)
	}

	_, err = cantonChangesets.ConfigureChainForLanes{}.Apply(*env, config)
	if err != nil {
		return fmt.Errorf("failed to configure chain for lanes: %w", err)
	}

	// Configure lock/release pool remote chain configs only at lane-setup time.
	const defaultLockReleaseQualifier = "TEST (LockReleaseTokenPool 1.7.0 [default] to BurnMintTokenPool 1.7.0 [default])"
	const reverseLockReleaseQualifier = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"
	chain := env.BlockChains.CantonChains()[selector]
	participant := chain.Participants[0]
	addresses := env.DataStore.Addresses()

	lockPoolRef, lockPoolErr := addresses.Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(lock_release_token_pool.ContractType),
		lock_release_token_pool.Version,
		defaultLockReleaseQualifier,
	))
	if lockPoolErr != nil || lockPoolRef.Address == "" {
		return fmt.Errorf("get lock/release pool address ref for lane remote config update: %w", lockPoolErr)
	}

	activePool, activePoolErr := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		contracts.HexToInstanceAddress(lockPoolRef.Address),
	)
	if activePoolErr != nil {
		return fmt.Errorf("find active lock/release pool for lane remote config update: %w", activePoolErr)
	}
	parsedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if parseErr != nil {
		return fmt.Errorf("parse active lock/release pool for lane remote config update: %w", parseErr)
	}

	ccvRefs := []common.RawInstanceAddress{}
	if committeeVerifierRawAddr != "" {
		ccvRefs = append(ccvRefs, committeeVerifierRawAddr.Binding())
	}

	updatedDataStore := datastore.NewMemoryDataStore()
	if mergeErr := updatedDataStore.Merge(env.DataStore); mergeErr != nil {
		return fmt.Errorf("clone datastore before rate limiter address persistence: %w", mergeErr)
	}

	updates := make([]lockreleasetokenpool.ChainUpdate, 0, len(remoteSelectors))
	for _, remoteSelector := range remoteSelectors {
		remotePoolHex := ""
		for _, candidate := range addresses.Filter(
			datastore.AddressRefByChainSelector(remoteSelector),
			datastore.AddressRefByType(datastore.ContractType("BurnMintTokenPool")),
		) {
			if candidate.Qualifier == reverseLockReleaseQualifier {
				remotePoolHex = candidate.Address
				break
			}
		}
		if remotePoolHex == "" {
			return fmt.Errorf("missing remote BurnMintTokenPool with qualifier %q for selector %d", reverseLockReleaseQualifier, remoteSelector)
		}

		remoteTokenHex := ""
		for _, candidate := range addresses.Filter(
			datastore.AddressRefByChainSelector(remoteSelector),
			datastore.AddressRefByType(datastore.ContractType("BurnMintERC20WithDrip")),
		) {
			if candidate.Qualifier == reverseLockReleaseQualifier {
				remoteTokenHex = candidate.Address
				break
			}
		}
		if remoteTokenHex == "" {
			return fmt.Errorf("missing remote BurnMintERC20WithDrip token with qualifier %q for selector %d", reverseLockReleaseQualifier, remoteSelector)
		}

		remoteSelectorKey := strconv.FormatUint(remoteSelector, 10)
		selectorForInstanceID := strings.ReplaceAll(remoteSelectorKey, ".", "-")

		outboundInstanceID := fmt.Sprintf("devenv-outbound-rl-%s", selectorForInstanceID)
		outboundQualifier := fmt.Sprintf("%s-outbound-%s", defaultLockReleaseQualifier, remoteSelectorKey)
		outboundOut, outboundErr := (cantonChangesets.DeployRateLimiter{}).Apply(*env, cantonChangesets.CantonCSDeps[cantonChangesets.DeployRateLimiterConfig]{
			ChainSelector: selector,
			Participant:   0,
			Config: cantonChangesets.DeployRateLimiterConfig{
				PoolOwner:           string(parsedPool.PoolOwner),
				PoolInstanceID:      string(parsedPool.InstanceId),
				RemoteChainSelector: remoteSelectorKey,
				Direction:           common.RateLimitDirectionRateLimitDirection_Outbound,
				Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
				InstanceID:          outboundInstanceID,
				Qualifier:           outboundQualifier,
				IsEnabled:           true,
				Capacity:            "999999999999999999",
				Rate:                "999999999999999999",
				Tokens:              "999999999999999999",
			},
		})
		if outboundErr != nil {
			return fmt.Errorf("create outbound rate limiter for selector %d: %w", remoteSelector, outboundErr)
		}
		if mergeErr := updatedDataStore.Merge(outboundOut.DataStore.Seal()); mergeErr != nil {
			return fmt.Errorf("merge outbound rate limiter datastore output for selector %d: %w", remoteSelector, mergeErr)
		}
		outboundRef, outboundRefErr := outboundOut.DataStore.Seal().Addresses().Get(datastore.NewAddressRefKey(
			selector,
			datastore.ContractType("RateLimiter"),
			cantonChangesets.RateLimiterVersion,
			outboundQualifier,
		))
		if outboundRefErr != nil {
			return fmt.Errorf("resolve outbound rate limiter address ref for selector %d: %w", remoteSelector, outboundRefErr)
		}
		if len(outboundRef.Labels.List()) == 0 {
			return fmt.Errorf("missing outbound rate limiter raw address label for selector %d", remoteSelector)
		}
		outboundRawAddr, parseOutboundRawErr := contracts.RawInstanceAddressFromString(outboundRef.Labels.List()[0])
		if parseOutboundRawErr != nil {
			return fmt.Errorf("parse outbound rate limiter raw address for selector %d: %w", remoteSelector, parseOutboundRawErr)
		}
		outboundRL := outboundRawAddr.Binding()

		inboundInstanceID := fmt.Sprintf("devenv-inbound-rl-%s", selectorForInstanceID)
		inboundQualifier := fmt.Sprintf("%s-inbound-%s", defaultLockReleaseQualifier, remoteSelectorKey)
		inboundOut, inboundErr := (cantonChangesets.DeployRateLimiter{}).Apply(*env, cantonChangesets.CantonCSDeps[cantonChangesets.DeployRateLimiterConfig]{
			ChainSelector: selector,
			Participant:   0,
			Config: cantonChangesets.DeployRateLimiterConfig{
				PoolOwner:           string(parsedPool.PoolOwner),
				PoolInstanceID:      string(parsedPool.InstanceId),
				RemoteChainSelector: remoteSelectorKey,
				Direction:           common.RateLimitDirectionRateLimitDirection_Inbound,
				Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
				InstanceID:          inboundInstanceID,
				Qualifier:           inboundQualifier,
				IsEnabled:           true,
				Capacity:            "999999999999999999",
				Rate:                "999999999999999999",
				Tokens:              "999999999999999999",
			},
		})
		if inboundErr != nil {
			return fmt.Errorf("create inbound rate limiter for selector %d: %w", remoteSelector, inboundErr)
		}
		if mergeErr := updatedDataStore.Merge(inboundOut.DataStore.Seal()); mergeErr != nil {
			return fmt.Errorf("merge inbound rate limiter datastore output for selector %d: %w", remoteSelector, mergeErr)
		}
		inboundRef, inboundRefErr := inboundOut.DataStore.Seal().Addresses().Get(datastore.NewAddressRefKey(
			selector,
			datastore.ContractType("RateLimiter"),
			cantonChangesets.RateLimiterVersion,
			inboundQualifier,
		))
		if inboundRefErr != nil {
			return fmt.Errorf("resolve inbound rate limiter address ref for selector %d: %w", remoteSelector, inboundRefErr)
		}
		if len(inboundRef.Labels.List()) == 0 {
			return fmt.Errorf("missing inbound rate limiter raw address label for selector %d", remoteSelector)
		}
		inboundRawAddr, parseInboundRawErr := contracts.RawInstanceAddressFromString(inboundRef.Labels.List()[0])
		if parseInboundRawErr != nil {
			return fmt.Errorf("parse inbound rate limiter raw address for selector %d: %w", remoteSelector, parseInboundRawErr)
		}
		inboundRL := inboundRawAddr.Binding()

		updates = append(updates, lockreleasetokenpool.ChainUpdate{
			RemoteChainSelector: types.NUMERIC(remoteSelectorKey),
			RemotePools:         []types.TEXT{types.TEXT(canonicalCantonRemotePoolHex(remotePoolHex))},
			RemoteTokenAddress:  types.TEXT(strings.ToLower(strings.TrimPrefix(remoteTokenHex, "0x"))),
			InboundCCVs:         ccvRefs,
			OutboundCCVs:        ccvRefs,
			InboundRateLimiter:  inboundRL,
			OutboundRateLimiter: outboundRL,
		})
	}

	if len(updates) > 0 {
		updateArgs := lockreleasetokenpool.ApplyChainUpdates{
			RemoteChainSelectorsToRemove: []types.NUMERIC{},
			ChainsToAdd:                  updates,
		}
		poolInstanceAddr := contracts.InstanceID(string(parsedPool.InstanceId)).RawInstanceAddress(parsedPool.PoolOwner).InstanceAddress()
		_, exerciseErr := operations.ExecuteOperation(
			env.OperationsBundle,
			canton_lock_release_token_pool.ApplyChainUpdates,
			dependencies.CantonDeps{Chain: chain},
			contract.ChoiceInput[lockreleasetokenpool.ApplyChainUpdates]{
				ChainSelector:   selector,
				InstanceAddress: poolInstanceAddr,
				ActAs:           []string{string(parsedPool.PoolOwner)},
				Args:            updateArgs,
			},
		)
		if exerciseErr != nil {
			return fmt.Errorf("apply lock/release pool chain updates operation for connected lanes: %w", exerciseErr)
		}
	}

	env.DataStore = updatedDataStore.Seal()

	return nil
}

// DeployLocalNetwork implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) DeployLocalNetwork(ctx context.Context, bcs *blockchain.Input) (*blockchain.Output, error) {
	c.logger.
		Info().
		Int("NumberOfCantonValidators", bcs.NumberOfCantonValidators).
		Msg("Deploying Canton network")
	out, err := blockchain.NewBlockchainNetwork(bcs)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain network: %w", err)
	}

	return out, nil
}

// FundAddresses implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) FundAddresses(ctx context.Context, bc *blockchain.Input, addresses []protocol.UnknownAddress, nativeAmount *big.Int) error {
	return nil // TODO: implement
}

// FundNodes implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) FundNodes(ctx context.Context, cls []*simple_node_set.Input, bc *blockchain.Input, linkAmount, nativeAmount *big.Int) error {
	return nil // TODO: implement
}

// Curse implements cciptestinterfaces.CCIP17.
func (c *Chain) Curse(ctx context.Context, subjects [][16]byte) error {
	rmnRemoteRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	if err != nil {
		return fmt.Errorf("get rmn remote address: %w", err)
	}

	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: 0,
	}
	instanceAddr := contracts.HexToInstanceAddress(rmnRemoteRef.Address)
	party := c.chain.Participants[0].PartyID

	c.logger.Info().
		Uint64("chainSelector", c.chainDetails.ChainSelector).
		Int("numSubjects", len(subjects)).
		Msg("Cursing subjects on chain")
	for _, subject := range subjects {
		_, err := operations.ExecuteOperation(c.e.OperationsBundle, rmn_remote.Curse, deps, contract.ChoiceInput[rmn.Curse]{
			ChainSelector:   c.chainDetails.ChainSelector,
			InstanceAddress: instanceAddr,
			ActAs:           []string{party},
			Args: rmn.Curse{
				Subject: types.TEXT(hex.EncodeToString(subject[:])),
			},
		})
		if err != nil {
			return fmt.Errorf("curse subject: %w", err)
		}
		c.logger.Info().
			Uint64("chainSelector", c.chainDetails.ChainSelector).
			Msg("Cursed chain")
	}

	return nil
}

// Uncurse implements cciptestinterfaces.CCIP17.
func (c *Chain) Uncurse(ctx context.Context, subjects [][16]byte) error {
	rmnRemoteRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	if err != nil {
		return fmt.Errorf("get rmn remote address: %w", err)
	}

	deps := dependencies.CantonDeps{
		Chain:       c.chain,
		Participant: 0,
	}
	instanceAddr := contracts.HexToInstanceAddress(rmnRemoteRef.Address)
	party := c.chain.Participants[0].PartyID

	c.logger.Info().
		Uint64("chainSelector", c.chainDetails.ChainSelector).
		Int("numSubjects", len(subjects)).
		Msg("Uncursing subjects on chain")
	for _, subject := range subjects {
		_, err := operations.ExecuteOperation(c.e.OperationsBundle, rmn_remote.Uncurse, deps, contract.ChoiceInput[rmn.Uncurse]{
			ChainSelector:   c.chainDetails.ChainSelector,
			InstanceAddress: instanceAddr,
			ActAs:           []string{party},
			Args: rmn.Uncurse{
				Subject: types.TEXT(hex.EncodeToString(subject[:])),
			},
		})
		if err != nil {
			return fmt.Errorf("uncurse subject: %w", err)
		}
		c.logger.Info().
			Uint64("chainSelector", c.chainDetails.ChainSelector).
			Msg("Uncursed chain")
	}

	return nil
}

// ExposeMetrics implements cciptestinterfaces.CCIP17.
func (c *Chain) ExposeMetrics(ctx context.Context, source, dest uint64) ([]string, *prometheus.Registry, error) {
	return nil, nil, nil // TODO: implement
}

// GetEOAReceiverAddress implements cciptestinterfaces.CCIP17.
func (c *Chain) GetEOAReceiverAddress() (protocol.UnknownAddress, error) {
	return protocol.UnknownAddress{}, nil // TODO: implement
}

// GetExpectedNextSequenceNumber implements cciptestinterfaces.CCIP17.
func (c *Chain) GetExpectedNextSequenceNumber(ctx context.Context, to uint64) (uint64, error) {
	return c.nextSeq + 1, nil
}

// GetMaxDataBytes implements cciptestinterfaces.CCIP17.
func (c *Chain) GetMaxDataBytes(ctx context.Context, remoteChainSelector uint64) (uint32, error) {
	return 0, nil // TODO: implement
}

// GetRoundRobinUser implements cciptestinterfaces.CCIP17.
func (c *Chain) GetRoundRobinUser() func() *bind.TransactOpts {
	return nil // TODO: implement
}

// GetSenderAddress implements cciptestinterfaces.CCIP17.
func (c *Chain) GetSenderAddress() (protocol.UnknownAddress, error) {
	return protocol.UnknownAddress{}, nil // TODO: implement
}

// GetTokenBalance implements cciptestinterfaces.CCIP17.
func (c *Chain) GetTokenBalance(ctx context.Context, address, tokenAddress protocol.UnknownAddress) (*big.Int, error) {
	ownerParty, err := c.resolvePartyFromHashedAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("resolve owner party from receiver address: %w", err)
	}
	participant := c.chain.Participants[c.participantIndexForParty(ownerParty)]

	var targetInstrument *splice_api_token_holding_v1.InstrumentId
	if len(tokenAddress.Bytes()) > 0 {
		instrument, resolveErr := c.resolveInstrumentIDForRemoteToken(ctx, tokenAddress)
		if resolveErr != nil {
			return nil, fmt.Errorf("resolve instrument for token address %s: %w", tokenAddress.String(), resolveErr)
		}
		targetInstrument = instrument
	}

	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for holdings query: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{
						{
							IdentifierFilter: &ledgerv2.CumulativeFilter_InterfaceFilter{
								InterfaceFilter: &ledgerv2.InterfaceFilter{
									InterfaceId: &ledgerv2.Identifier{
										PackageId:  "#splice-api-token-holding-v1",
										ModuleName: "Splice.Api.Token.HoldingV1",
										EntityName: "Holding",
									},
									IncludeInterfaceView: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active holding contracts: %w", err)
	}
	defer stream.CloseSend()

	total := big.NewInt(0)
	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, fmt.Errorf("receive active holding contracts: %w", recvErr)
		}

		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}
		view := entry.ActiveContract.GetCreatedEvent().GetInterfaceViews()
		if len(view) == 0 {
			continue
		}
		viewRecord := view[0].GetViewValue()
		if viewRecord == nil {
			continue
		}

		ownerField := viewRecord.GetFields()
		if len(ownerField) < 3 {
			continue
		}
		owner := ownerField[0].GetValue().GetParty()
		if owner != ownerParty {
			continue
		}

		if targetInstrument != nil {
			instrumentRecord := ownerField[1].GetValue().GetRecord()
			if instrumentRecord == nil || len(instrumentRecord.GetFields()) < 2 {
				continue
			}
			holdingAdmin := instrumentRecord.GetFields()[0].GetValue().GetParty()
			holdingID := instrumentRecord.GetFields()[1].GetValue().GetText()
			if holdingAdmin != string(targetInstrument.Admin) || holdingID != string(targetInstrument.Id) {
				continue
			}
		}

		amountNumeric := ownerField[2].GetValue().GetNumeric()
		amountInt, parseErr := damlNumericToBigInt(amountNumeric)
		if parseErr != nil {
			return nil, fmt.Errorf("parse holding amount %q: %w", amountNumeric, parseErr)
		}
		total.Add(total, amountInt)
	}

	return total, nil
}

// NativeBalance implements cciptestinterfaces.CCIP17.
// Canton does not have an EVM-like native token balance surface.
func (c *Chain) NativeBalance(ctx context.Context, address protocol.UnknownAddress) (*big.Int, error) {
	return big.NewInt(0), nil
}

// TransferNative implements cciptestinterfaces.CCIP17.
// Canton does not support native transfers through this API.
func (c *Chain) TransferNative(ctx context.Context, from, to protocol.UnknownAddress, amount *big.Int) error {
	return nil
}

// GetUserNonce implements cciptestinterfaces.CCIP17.
func (c *Chain) GetUserNonce(ctx context.Context, userAddress protocol.UnknownAddress) (uint64, error) {
	return 0, nil // TODO: implement
}

func normalizeNumericForCompare(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	if strings.Contains(v, ".") {
		parts := strings.SplitN(v, ".", 2)
		frac := strings.TrimRight(parts[1], "0")
		if frac == "" {
			return parts[0]
		}
	}
	if strings.ContainsAny(v, "eE") {
		if f, _, err := big.ParseFloat(v, 10, 256, big.ToZero); err == nil {
			if i, _ := f.Int(nil); i != nil {
				return i.String()
			}
		}
	}

	return strings.TrimSuffix(v, ".")
}

func canonicalCantonRemotePoolHex(remotePoolHex string) string {
	clean := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(remotePoolHex), "0x"))
	if clean == "" {
		return clean
	}
	if len(clean) > 64 {
		return clean[len(clean)-64:]
	}
	if len(clean) < 64 {
		return strings.Repeat("0", 64-len(clean)) + clean
	}

	return clean
}

func findRemoteChainConfigBySelector(remoteChainConfigs map[string]any, selectorKey string) (any, bool) {
	selectorNorm := normalizeNumericForCompare(selectorKey)
	for rawKey, cfg := range remoteChainConfigs {
		if normalizeNumericForCompare(rawKey) == selectorNorm {
			return cfg, true
		}
	}

	return nil, false
}

func remoteChainConfigFromAny(v any) (lockreleasetokenpool.RemoteChainConfig, bool) {
	switch cfg := v.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		return cfg, true
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		out := lockreleasetokenpool.RemoteChainConfig{
			InboundCCVs:  []common.RawInstanceAddress{},
			OutboundCCVs: []common.RawInstanceAddress{},
			RemotePools:  []types.TEXT{},
		}
		if raw, ok := m["inboundRateLimiter"]; ok {
			if unpack, err := extractRawRateLimiterAddress(raw); err == nil && unpack != "" {
				out.InboundRateLimiter = common.RawInstanceAddress{Unpack: types.TEXT(unpack)}
			}
		}
		if raw, ok := m["outboundRateLimiter"]; ok {
			if unpack, err := extractRawRateLimiterAddress(raw); err == nil && unpack != "" {
				out.OutboundRateLimiter = common.RawInstanceAddress{Unpack: types.TEXT(unpack)}
			}
		}
		if remoteTokenAddress, ok := m["remoteTokenAddress"]; ok {
			switch rv := remoteTokenAddress.(type) {
			case string:
				out.RemoteTokenAddress = types.TEXT(rv)
			case map[string]any:
				if val, ok := rv["value"].(string); ok {
					out.RemoteTokenAddress = types.TEXT(val)
				} else if data, ok := rv["data"].(map[string]any); ok {
					if val, ok := data["value"].(string); ok {
						out.RemoteTokenAddress = types.TEXT(val)
					} else if text, ok := data["text"].(string); ok {
						out.RemoteTokenAddress = types.TEXT(text)
					}
				}
			default:
				out.RemoteTokenAddress = types.TEXT(fmt.Sprint(rv))
			}
		}
		if remoteRaw, ok := m["remotePools"].([]any); ok {
			for _, rp := range remoteRaw {
				out.RemotePools = append(out.RemotePools, types.TEXT(fmt.Sprint(rp)))
			}
		}

		return out, true
	default:
		return lockreleasetokenpool.RemoteChainConfig{}, false
	}
}

func resolveRateLimiterForSend(
	ctx context.Context,
	participant canton.Participant,
	parsedTokenPool *lockreleasetokenpool.LockReleaseTokenPool,
	destSelectorKey string,
	destSelectorNumericKey string,
	resolveDisclosedByAddress func(templateID string, address contracts.InstanceAddress) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error),
) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
	var expectedRawAddress string
	selectorNorm := normalizeNumericForCompare(destSelectorKey)
	if cfgAny, ok := findRemoteChainConfigBySelector(parsedTokenPool.RemoteChainConfigs, destSelectorKey); ok {
		if cfg, cfgOK := remoteChainConfigFromAny(cfgAny); cfgOK {
			outboundRateLimiterRaw := cfg.OutboundRateLimiter
			if rawAddress, rawErr := extractRawRateLimiterAddress(outboundRateLimiterRaw); rawErr == nil {
				expectedRawAddress = rawAddress
			}
			cid, disclosure, err := resolveRateLimiterFromRawAddressForSend(outboundRateLimiterRaw, resolveDisclosedByAddress)
			if err == nil {
				return cid, disclosure, nil
			}
		}
	}
	if cfgAny, ok := parsedTokenPool.RemoteChainConfigs[destSelectorNumericKey]; ok && expectedRawAddress == "" {
		if cfg, cfgOK := remoteChainConfigFromAny(cfgAny); cfgOK {
			if rawAddress, rawErr := extractRawRateLimiterAddress(cfg.OutboundRateLimiter); rawErr == nil {
				expectedRawAddress = rawAddress
			}
		}
	}

	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return "", nil, fmt.Errorf("get ledger end for rate limiter fallback lookup: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{
						{
							IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &ledgerv2.TemplateFilter{
									TemplateId: &ledgerv2.Identifier{
										PackageId:  "#ccip-common",
										ModuleName: "CCIP.RateLimiter",
										EntityName: "RateLimiter",
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("query active rate limiters for fallback lookup: %w", err)
	}
	defer stream.CloseSend()

	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return "", nil, fmt.Errorf("receive active rate limiters for fallback lookup: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		parsedRateLimiter, parseErr := bindings.UnmarshalCreatedEvent[common.RateLimiter](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			continue
		}
		if parsedRateLimiter.Direction != common.RateLimitDirectionRateLimitDirection_Outbound {
			continue
		}
		if string(parsedRateLimiter.PoolOwner) != string(parsedTokenPool.PoolOwner) || string(parsedRateLimiter.PoolInstanceId) != string(parsedTokenPool.InstanceId) {
			continue
		}
		if normalizeNumericForCompare(string(parsedRateLimiter.RemoteChainSelector)) != selectorNorm {
			continue
		}
		if expectedRawAddress != "" {
			candidateRawAddress := fmt.Sprintf("%s@%s", parsedRateLimiter.InstanceId, parsedRateLimiter.PoolOwner)
			if candidateRawAddress != expectedRawAddress {
				continue
			}
		}

		return types.CONTRACT_ID(entry.ActiveContract.GetCreatedEvent().GetContractId()), convertToDisclosedContract(entry.ActiveContract), nil
	}

	return "", nil, fmt.Errorf("missing outbound rate limiter for destination selector %s", destSelectorKey)
}

func resolveRateLimiterFromRawAddressForSend(
	outboundRateLimiterRaw any,
	resolveDisclosedByAddress func(templateID string, address contracts.InstanceAddress) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error),
) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
	rawAddressString, err := extractRawRateLimiterAddress(outboundRateLimiterRaw)
	if err != nil {
		return "", nil, err
	}
	rawAddr, err := contracts.RawInstanceAddressFromString(rawAddressString)
	if err != nil {
		return "", nil, fmt.Errorf("parse outbound rate limiter raw instance address: %w", err)
	}

	return resolveDisclosedByAddress(common.RateLimiter{}.GetTemplateID(), rawAddr.InstanceAddress())
}

func extractRawRateLimiterAddress(outboundRateLimiterRaw any) (string, error) {
	switch value := outboundRateLimiterRaw.(type) {
	case common.RawInstanceAddress:
		return string(value.Unpack), nil
	case map[string]any:
		m := value
		if data, ok := value["data"].(map[string]any); ok {
			m = data
		}
		unpack, ok := m["unpack"].(string)
		if !ok || unpack == "" {
			return "", fmt.Errorf("missing unpack in outbound rate limiter map")
		}

		return unpack, nil
	default:
		return "", fmt.Errorf("unexpected outbound rate limiter type %T", outboundRateLimiterRaw)
	}
}

func (c *Chain) resolvePartyFromHashedAddress(ctx context.Context, address protocol.UnknownAddress) (string, error) {
	target := contracts.BytesToHashedParty(address.Bytes())
	for _, participant := range c.chain.Participants {
		if contracts.HashedPartyFromString(participant.PartyID) == target {
			return participant.PartyID, nil
		}
	}

	for _, participant := range c.chain.Participants {
		token, tokenErr := participant.TokenSource.Token()
		if tokenErr != nil {
			continue
		}
		conn, connErr := grpc.NewClient(
			participant.Endpoints.GRPCLedgerAPIURL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithPerRPCCredentials(auth.NewBearerToken(token.AccessToken)),
		)
		if connErr != nil {
			continue
		}
		pmClient := adminv2.NewPartyManagementServiceClient(conn)
		resp, listErr := pmClient.ListKnownParties(ctx, &adminv2.ListKnownPartiesRequest{})
		_ = conn.Close()
		if listErr != nil {
			continue
		}
		for _, details := range resp.GetPartyDetails() {
			if contracts.HashedPartyFromString(details.GetParty()) == target {
				return details.GetParty(), nil
			}
		}
	}

	return "", fmt.Errorf("no party found for hashed address %s", target.Hex())
}

func (c *Chain) participantIndexForParty(party string) int {
	for i, participant := range c.chain.Participants {
		if participant.PartyID == party {
			return i
		}
	}

	return 0
}

func (c *Chain) resolveInstrumentIDForRemoteToken(
	ctx context.Context,
	tokenAddress protocol.UnknownAddress,
) (*splice_api_token_holding_v1.InstrumentId, error) {
	participant := c.chain.Participants[0]
	normalizeTokenHex := func(raw string) string {
		clean := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(raw), "0x"))
		if len(clean) > 40 {
			clean = clean[len(clean)-40:]
		}
		return clean
	}

	ledgerEnd, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for token pool query: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{
						{
							IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &ledgerv2.TemplateFilter{
									TemplateId: &ledgerv2.Identifier{
										PackageId:  "#ccip-lockreleasetokenpool",
										ModuleName: "CCIP.LockReleaseTokenPool",
										EntityName: "LockReleaseTokenPool",
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active lock/release token pools: %w", err)
	}
	defer stream.CloseSend()

	want := normalizeTokenHex(tokenAddress.String())
	var found *splice_api_token_holding_v1.InstrumentId

	for {
		resp, recvErr := stream.Recv()
		if errors.Is(recvErr, io.EOF) {
			break
		}
		if recvErr != nil {
			return nil, fmt.Errorf("receive lock/release token pools: %w", recvErr)
		}

		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}
		parsed, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](entry.ActiveContract.GetCreatedEvent())
		if parseErr != nil {
			return nil, fmt.Errorf("parse lock/release token pool created event: %w", parseErr)
		}

		for _, cfgAny := range parsed.RemoteChainConfigs {
			cfg, ok := remoteChainConfigFromAny(cfgAny)
			if !ok {
				continue
			}
			remoteToken := normalizeTokenHex(string(cfg.RemoteTokenAddress))
			if remoteToken != want {
				continue
			}
			if found != nil && (found.Admin != parsed.InstrumentId.Admin || found.Id != parsed.InstrumentId.Id) {
				return nil, fmt.Errorf("multiple lock/release pools mapped token %s to different instruments", tokenAddress.String())
			}
			inst := parsed.InstrumentId
			found = &inst

			break
		}
	}

	if found == nil {
		// Devenv fallback: when mapping wasn't populated yet, use the single active lock/release pool instrument.
		var singleton *splice_api_token_holding_v1.InstrumentId
		ledgerEnd2, err2 := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
		if err2 != nil {
			return nil, fmt.Errorf("no lock/release pool mapping found for token %s", tokenAddress.String())
		}
		stream2, err2 := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
			ActiveAtOffset: ledgerEnd2.GetOffset(),
			EventFormat: &ledgerv2.EventFormat{
				FiltersByParty: map[string]*ledgerv2.Filters{
					participant.PartyID: {
						Cumulative: []*ledgerv2.CumulativeFilter{
							{
								IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
									TemplateFilter: &ledgerv2.TemplateFilter{
										TemplateId: &ledgerv2.Identifier{
											PackageId:  "#ccip-lockreleasetokenpool",
											ModuleName: "CCIP.LockReleaseTokenPool",
											EntityName: "LockReleaseTokenPool",
										},
										IncludeCreatedEventBlob: true,
									},
								},
							},
						},
					},
				},
				Verbose: true,
			},
		})
		if err2 == nil {
			defer stream2.CloseSend()
			for {
				resp2, recvErr2 := stream2.Recv()
				if errors.Is(recvErr2, io.EOF) {
					break
				}
				if recvErr2 != nil {
					break
				}
				entry2, ok2 := resp2.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
				if !ok2 {
					continue
				}
				parsed2, parseErr2 := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](entry2.ActiveContract.GetCreatedEvent())
				if parseErr2 != nil {
					continue
				}
				inst := parsed2.InstrumentId
				if singleton != nil && (singleton.Admin != inst.Admin || singleton.Id != inst.Id) {
					singleton = nil
					break
				}
				singleton = &inst
			}
		}
		if singleton != nil {
			return singleton, nil
		}

		return nil, fmt.Errorf("no lock/release pool mapping found for token %s", tokenAddress.String())
	}

	return found, nil
}

func damlNumericToBigInt(v string) (*big.Int, error) {
	if v == "" {
		return nil, fmt.Errorf("empty numeric value")
	}
	if !strings.Contains(v, ".") {
		n, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer numeric %q", v)
		}

		return n, nil
	}

	parts := strings.SplitN(v, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid decimal numeric %q", v)
	}
	intPart := parts[0]
	if intPart == "" {
		intPart = "0"
	}
	if intPart == "-" {
		intPart = "-0"
	}
	n, ok := new(big.Int).SetString(intPart, 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer component %q in numeric %q", intPart, v)
	}

	return n, nil
}

func (c *Chain) findLatestActiveContractByInstanceAddress(
	ctx context.Context,
	participant canton.Participant,
	templateID string,
	address contracts.InstanceAddress,
) (*ledgerv2.ActiveContract, error) {
	active, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, templateID, address)
	if err == nil {
		return active, nil
	}
	if !strings.Contains(err.Error(), "multiple active contracts found") {
		return nil, err
	}

	parts := strings.Split(templateID, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid template ID for fallback lookup %q: %w", templateID, err)
	}
	packageID, moduleName, entityName := parts[0], parts[1], parts[2]

	lookupCtx, cancelLookup := context.WithTimeout(ctx, 20*time.Second)
	defer cancelLookup()

	ledgerEnd, endErr := participant.LedgerServices.State.GetLedgerEnd(lookupCtx, &ledgerv2.GetLedgerEndRequest{})
	if endErr != nil {
		return nil, fmt.Errorf("get ledger end for fallback lookup: %w", endErr)
	}
	stream, streamErr := participant.LedgerServices.State.GetActiveContracts(lookupCtx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEnd.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{
						{
							IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &ledgerv2.TemplateFilter{
									TemplateId: &ledgerv2.Identifier{
										PackageId:  packageID,
										ModuleName: moduleName,
										EntityName: entityName,
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if streamErr != nil {
		return nil, fmt.Errorf("get active contracts for fallback lookup: %w", streamErr)
	}
	defer stream.CloseSend()

	var latestMatch *ledgerv2.ActiveContract
	for {
		resp, recvErr := stream.Recv()
		if recvErr != nil {
			if lookupCtx.Err() != nil {
				return nil, fmt.Errorf("fallback lookup timed out while reading active contracts for %s: %w", address.String(), lookupCtx.Err())
			}
			if errors.Is(recvErr, io.EOF) {
				break
			}

			return nil, fmt.Errorf("receive active contracts for fallback lookup: %w", recvErr)
		}
		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		createArgs := created.GetCreateArguments()
		if createArgs == nil {
			continue
		}
		var instanceIDText string
		for _, f := range createArgs.GetFields() {
			if f.GetLabel() == "instanceId" {
				instanceIDText = f.GetValue().GetText()
				break
			}
		}
		if instanceIDText == "" || len(created.GetSignatories()) != 1 {
			continue
		}
		gotAddr := contracts.InstanceID(instanceIDText).RawInstanceAddress(types.PARTY(created.GetSignatories()[0])).InstanceAddress()
		if gotAddr != address {
			continue
		}
		if latestMatch == nil || created.GetOffset() > latestMatch.GetCreatedEvent().GetOffset() {
			latestMatch = entry.ActiveContract
		}
	}
	if latestMatch == nil {
		return nil, err
	}

	return latestMatch, nil
}

func (c *Chain) resolveDisclosedContractByAddress(
	ctx context.Context,
	participant canton.Participant,
	templateID string,
	address contracts.InstanceAddress,
) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
	active, err := c.findLatestActiveContractByInstanceAddress(ctx, participant, templateID, address)
	if err != nil {
		return "", nil, err
	}

	return types.CONTRACT_ID(active.GetCreatedEvent().GetContractId()), convertToDisclosedContract(active), nil
}

// SendMessage implements cciptestinterfaces.CCIP17.
func (c *Chain) SendMessage(ctx context.Context, dest uint64, fields cciptestinterfaces.MessageFields, opts cciptestinterfaces.MessageOptions) (cciptestinterfaces.MessageSentEvent, error) {
	participant := c.chain.Participants[0]
	party := participant.PartyID

	seqNo := c.nextSeq + 1
	findLatestActiveByAddress := func(templateID string, address contracts.InstanceAddress) (*ledgerv2.ActiveContract, error) {
		return c.findLatestActiveContractByInstanceAddress(ctx, participant, templateID, address)
	}
	resolveDisclosedByAddress := func(templateID string, address contracts.InstanceAddress) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
		return c.resolveDisclosedContractByAddress(ctx, participant, templateID, address)
	}

	// Resolve commonly required contracts and disclosures.
	onRampRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(onramp.ContractType), onramp.Version, ""))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get source onramp address: %w", err)
	}
	globalConfigRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get global config address: %w", err)
	}
	tokenAdminRegistryRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version, ""))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get token admin registry address: %w", err)
	}
	rmnRemoteRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get rmn remote address: %w", err)
	}
	feeQuoterRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get fee quoter address: %w", err)
	}

	onRampCID, disclosedOnRamp, err := resolveDisclosedByAddress(onrampbindings.OnRamp{}.GetTemplateID(), contracts.HexToInstanceAddress(onRampRef.Address))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve onramp disclosed contract: %w", err)
	}
	globalConfigCID, disclosedGlobalConfig, err := resolveDisclosedByAddress(common.GlobalConfig{}.GetTemplateID(), contracts.HexToInstanceAddress(globalConfigRef.Address))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve global config disclosed contract: %w", err)
	}
	tokenAdminRegistryCID, disclosedTokenAdminRegistry, err := resolveDisclosedByAddress(tokenadminregistry.TokenAdminRegistry{}.GetTemplateID(), contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve token admin registry disclosed contract: %w", err)
	}
	rmnRemoteCID, disclosedRMNRemote, err := resolveDisclosedByAddress(rmn.RMNRemote{}.GetTemplateID(), contracts.HexToInstanceAddress(rmnRemoteRef.Address))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve rmn remote disclosed contract: %w", err)
	}
	feeQuoterCID, disclosedFeeQuoter, err := resolveDisclosedByAddress(feequoter.FeeQuoter{}.GetTemplateID(), contracts.HexToInstanceAddress(feeQuoterRef.Address))
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve fee quoter disclosed contract: %w", err)
	}

	// Router for sender party.
	routerAddress, err := c.DeployPerPartyRouter(ctx, party)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("deploy per-party router: %w", err)
	}
	routerCID, disclosedRouter, err := resolveDisclosedByAddress(perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve per-party router disclosed contract: %w", err)
	}

	// Deploy a sender-owned CCIPSender contract.
	senderInstanceID := contracts.MustNewInstanceID("devenv-ccipsender")
	createSenderRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					CreateArguments: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
						{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: senderInstanceID.String()}}},
						{Label: "owner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: party}}},
					}},
				}},
			}},
			ActAs: []string{party},
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("deploy ccip sender contract: %w", err)
	}
	senderInstanceAddress := senderInstanceID.RawInstanceAddress(types.PARTY(party)).InstanceAddress()
	var ccipSenderCID types.CONTRACT_ID
	for _, ev := range createSenderRes.GetTransaction().GetEvents() {
		created := ev.GetCreated()
		if created == nil || created.GetTemplateId() == nil {
			continue
		}
		if created.GetTemplateId().GetEntityName() == "CCIPSender" && created.GetContractId() != "" {
			ccipSenderCID = types.CONTRACT_ID(created.GetContractId())
			break
		}
	}
	_, disclosedCCIPSender, err := resolveDisclosedByAddress(ccipsender.CCIPSender{}.GetTemplateID(), senderInstanceAddress)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve ccip sender disclosed contract: %w", err)
	}
	if ccipSenderCID == "" {
		ccipSenderCID = types.CONTRACT_ID(disclosedCCIPSender.GetContractId())
	}
	if ccipSenderCID == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolved empty ccip sender contract ID")
	}

	// Deploy a minimal Canton-side executor that implements CCIP.Interfaces.Executor.IExecutor.
	// Use the same instance ID as the registered executor proxy so that the derived address
	// matches what the executor service expects (stored in CLDF labels during DeployContractsForSelector).
	executorProxyRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(executor.ProxyType),
			executor.Version,
			devenvcommon.DefaultExecutorQualifier,
		),
	)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get executor proxy address ref: %w", err)
	}
	labels := executorProxyRef.Labels.List()
	if len(labels) == 0 {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("executor proxy address ref has no labels; raw instance address is required")
	}
	executorProxyRawParts := strings.SplitN(labels[0], "@", 2)
	if len(executorProxyRawParts) != 2 {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("executor proxy label %q is not a valid raw instance address", labels[0])
	}
	executorInstanceID := contracts.InstanceID(executorProxyRawParts[0])
	destSelectorText := strconv.FormatUint(dest, 10)
	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-executor", ModuleName: "CCIP.Executor", EntityName: "Executor"},
					CreateArguments: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
						{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: executorInstanceID.String()}}},
						{Label: "owner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: party}}},
						{Label: "maxCCVsPerMsg", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 10}}},
						{Label: "dynamicConfig", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
							{Label: "feeAggregator", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Optional{Optional: &ledgerv2.Optional{}}}},
							{Label: "minBlockConfirmations", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
							{Label: "ccvAllowlistEnabled", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: false}}},
						}}}}},
						{Label: "allowedCCVs", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
						{
							Label: "remoteChainConfigs",
							Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{
								Entries: []*ledgerv2.GenMap_Entry{
									{
										Key: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: destSelectorText}},
										Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
											{Label: "feeUSDCents", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "0"}}},
											{Label: "enabled", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: true}}},
										}}}},
									},
								},
							}}},
						},
					}},
				}},
			}},
			ActAs: []string{party},
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("deploy executor contract: %w", err)
	}
	executorInstanceAddress := executorInstanceID.RawInstanceAddress(types.PARTY(party)).InstanceAddress()
	executorCID, disclosedExecutor, err := resolveDisclosedByAddress("#ccip-executor:CCIP.Executor:Executor", executorInstanceAddress)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve executor disclosed contract: %w", err)
	}
	if executorCID == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolved empty executor contract ID")
	}

	// Keep fee token setup local and deterministic for devenv sends.
	feeTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("devenv-fee-token"),
	}
	linkTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("link-token"),
	}
	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter"},
					ContractId: string(feeQuoterCID),
					Choice:     "ApplyFeeTokenUpdates",
					ChoiceArgument: ledger.MapToValue(feequoter.ApplyFeeTokenUpdates{
						FeeTokensToRemove: []splice_api_token_holding_v1.InstrumentId{},
						FeeTokensToAdd: []feequoter.FeeTokenArgs{
							{
								InstrumentId:      feeTokenInstrument,
								PremiumMultiplier: types.NUMERIC("100000000"),
							},
						},
					}),
				}},
			}},
			ActAs:              []string{party},
			DisclosedContracts: []*ledgerv2.DisclosedContract{disclosedFeeQuoter},
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("apply fee token updates: %w", err)
	}

	// ApplyFeeTokenUpdates archives/recreates FeeQuoter; refresh CID+disclosure before next exercise.
	feeQuoterCID, disclosedFeeQuoter, err = resolveDisclosedByAddress(
		feequoter.FeeQuoter{}.GetTemplateID(),
		contracts.HexToInstanceAddress(feeQuoterRef.Address),
	)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("refresh fee quoter after fee token update: %w", err)
	}

	_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-feequoter", ModuleName: "CCIP.FeeQuoter", EntityName: "FeeQuoter"},
					ContractId: string(feeQuoterCID),
					Choice:     "UpdatePrices",
					ChoiceArgument: ledger.MapToValue(feequoter.UpdatePrices{
						PriceUpdates: feequoter.PriceUpdates{
							TokenPriceUpdates: []feequoter.TokenPriceUpdate{
								{InstrumentId: feeTokenInstrument, UsdPerToken: types.NUMERIC("100000000")},
								{InstrumentId: linkTokenInstrument, UsdPerToken: types.NUMERIC("100000000")},
							},
							GasPriceUpdates: []feequoter.GasPriceUpdate{},
						},
						Caller: types.PARTY(party),
					}),
				}},
			}},
			ActAs:              []string{party},
			DisclosedContracts: []*ledgerv2.DisclosedContract{disclosedFeeQuoter},
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("update fee token prices: %w", err)
	}
	// UpdatePrices also archives/recreates FeeQuoter; refresh again for Send context/disclosures.
	feeQuoterCID, disclosedFeeQuoter, err = resolveDisclosedByAddress(
		feequoter.FeeQuoter{}.GetTemplateID(),
		contracts.HexToInstanceAddress(feeQuoterRef.Address),
	)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("refresh fee quoter after price update: %w", err)
	}

	defaultCCVRef, defaultCCVRefErr := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(committee_verifier.ContractType),
			committee_verifier.Version,
			devenvcommon.DefaultCommitteeVerifierQualifier,
		),
	)

	senderRequiredCCVs := make([]common.RawInstanceAddress, 0, len(opts.CCVs))
	ccvSendInputs := make([]ccipsender.CCVSendInput, 0, len(opts.CCVs))
	disclosedVerifierContracts := make([]*ledgerv2.DisclosedContract, 0, len(opts.CCVs))
	receiptIssuers := make([]protocol.UnknownAddress, 0, len(opts.CCVs)+2)
	var fallbackVerifierDestAddress protocol.UnknownAddress
	var fallbackVerifierBlob []byte
	for _, ccvItem := range opts.CCVs {
		verifierAddress := contracts.BytesToInstanceAddress(ccvItem.CCVAddress.Bytes())
		// In this devenv e2e lane we use a single default committee verifier.
		// Prefer datastore canonical address to avoid stale/rotated test-input drift.
		if len(opts.CCVs) == 1 && defaultCCVRefErr == nil && defaultCCVRef.Address != "" {
			verifierAddress = contracts.HexToInstanceAddress(defaultCCVRef.Address)
		}
		activeVerifier, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, ccvs.CommitteeVerifier{}.GetTemplateID(), verifierAddress)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve verifier by address %s: %w", verifierAddress.String(), err)
		}
		parsedVerifier, err := bindings.UnmarshalCreatedEvent[ccvs.CommitteeVerifier](activeVerifier.GetCreatedEvent())
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse committee verifier contract: %w", err)
		}
		var rawAddr contracts.RawInstanceAddress
		// Prefer canonical raw address from datastore labels (same source used in integration tests).
		ccvRef, err := c.e.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				c.chainDetails.ChainSelector,
				datastore.ContractType(committee_verifier.ContractType),
				committee_verifier.Version,
				devenvcommon.DefaultCommitteeVerifierQualifier,
			),
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve committee verifier address ref: %w", err)
		}
		if strings.EqualFold(ccvRef.Address, verifierAddress.String()) && len(ccvRef.Labels.List()) > 0 {
			rawAddr, err = contracts.RawInstanceAddressFromString(ccvRef.Labels.List()[0])
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse committee verifier raw address label: %w", err)
			}
		} else {
			rawAddr, err = contracts.RawInstanceAddressFromString(fmt.Sprintf("%s@%s", parsedVerifier.InstanceId, parsedVerifier.Owner))
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("construct verifier raw address: %w", err)
			}
		}
		senderRequiredCCVs = append(senderRequiredCCVs, rawAddr.Binding())
		ccvSendInputs = append(ccvSendInputs, ccipsender.CCVSendInput{
			CcvCid:          types.CONTRACT_ID(activeVerifier.GetCreatedEvent().GetContractId()),
			VerifierArgs:    types.TEXT(hex.EncodeToString(ccvItem.Args)),
			CcvExtraContext: common.CCIPContext{},
		})
		disclosedVerifierContracts = append(disclosedVerifierContracts, convertToDisclosedContract(activeVerifier))
		receiptIssuers = append(receiptIssuers, protocol.UnknownAddress(verifierAddress.Bytes()))
		if len(fallbackVerifierDestAddress) == 0 {
			fallbackVerifierDestAddress = protocol.UnknownAddress(verifierAddress.Bytes())
		}
		if len(fallbackVerifierBlob) == 0 {
			versionTagBytes, decodeErr := hex.DecodeString(string(parsedVerifier.VersionTag))
			if decodeErr == nil {
				fallbackVerifierBlob = versionTagBytes
			}
		}
	}
	if len(opts.Executor) > 0 {
		receiptIssuers = append(receiptIssuers, opts.Executor)
	}
	receiptIssuers = append(receiptIssuers, protocol.UnknownAddress(contracts.HexToInstanceAddress(onRampRef.Address).Bytes()))

	var tokenTransferInput *ccipsender.TokenTransferInput
	var disclosedTokenPool *ledgerv2.DisclosedContract
	var disclosedRateLimiter *ledgerv2.DisclosedContract
	var disclosedTransferFactoryContracts []*ledgerv2.DisclosedContract
	if fields.TokenAmount.Amount != nil {
		if fields.TokenAmount.Amount.Sign() < 0 {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token amount must be non-negative")
		}

		// Use the default Canton->EVM lock/release pool pair wired by devenv topology.
		//nolint:gosec // Topology qualifier string, not credentials.
		const tokenPoolQualifier = "TEST (LockReleaseTokenPool 1.7.0 [default] to BurnMintTokenPool 1.7.0 [default])"
		const remoteDestBurnMintQualifier = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"
		tokenPoolRef, err := c.e.DataStore.Addresses().Get(
			datastore.NewAddressRefKey(
				c.chainDetails.ChainSelector,
				datastore.ContractType(lock_release_token_pool.ContractType),
				lock_release_token_pool.Version,
				tokenPoolQualifier,
			),
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve source token pool address ref: %w", err)
		}
		tokenPoolAddress := contracts.HexToInstanceAddress(tokenPoolRef.Address)
		tokenPoolCID, tokenPoolDisclosure, err := resolveDisclosedByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), tokenPoolAddress)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve token pool disclosed contract: %w", err)
		}
		activeTokenPool, err := findLatestActiveByAddress(
			lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
			tokenPoolAddress,
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve active token pool by configured address: %w", err)
		}
		parsedTokenPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activeTokenPool.GetCreatedEvent())
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse token pool contract: %w", err)
		}
		// Token transfer must use the TAR bound to the selected pool deps.
		poolTARRaw := strings.TrimSpace(string(parsedTokenPool.Deps.TokenAdminRegistry.Unpack))
		if poolTARRaw != "" {
			poolTARRawAddr, parseErr := contracts.RawInstanceAddressFromString(poolTARRaw)
			if parseErr != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse selected pool token admin registry address %q: %w", poolTARRaw, parseErr)
			}
			poolTARCID, poolTARDisclosure, resolveErr := resolveDisclosedByAddress(
				tokenadminregistry.TokenAdminRegistry{}.GetTemplateID(),
				poolTARRawAddr.InstanceAddress(),
			)
			if resolveErr != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve selected pool token admin registry disclosed contract: %w", resolveErr)
			}
			tokenAdminRegistryCID = poolTARCID
			disclosedTokenAdminRegistry = poolTARDisclosure
		}
		destSelectorKey := strconv.FormatUint(dest, 10)
		remoteConfigKeys := make([]string, 0, len(parsedTokenPool.RemoteChainConfigs))
		for k := range parsedTokenPool.RemoteChainConfigs {
			remoteConfigKeys = append(remoteConfigKeys, k)
		}
		c.logger.Debug().
			Str("DestSelector", destSelectorKey).
			Strs("RemoteChainConfigKeys", remoteConfigKeys).
			Msg("Resolved lock/release pool remote chain config keys before token transfer send")
		destCfgAny, hasDestCfg := findRemoteChainConfigBySelector(parsedTokenPool.RemoteChainConfigs, destSelectorKey)
		if !hasDestCfg {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("missing lock/release pool remote chain config for selector %s (available keys: %v)", destSelectorKey, remoteConfigKeys)
		}
		destCfg, cfgOK := remoteChainConfigFromAny(destCfgAny)
		if !cfgOK {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("invalid lock/release pool remote chain config for selector %s", destSelectorKey)
		}
		if len(destCfg.RemotePools) == 0 || strings.TrimSpace(string(destCfg.RemoteTokenAddress)) == "" {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("incomplete lock/release pool remote chain config for selector %s", destSelectorKey)
		}
		if strings.TrimSpace(string(destCfg.OutboundRateLimiter.Unpack)) == "" || strings.TrimSpace(string(destCfg.InboundRateLimiter.Unpack)) == "" {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("missing rate limiter config in lock/release pool for selector %s", destSelectorKey)
		}
		if fields.TokenAmount.Amount.Sign() == 0 {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token amount must be greater than zero for token transfer")
		}
		tokenDestSelectorKey := strconv.FormatUint(dest, 10)
		tokenDestSelectorNumericKey := tokenDestSelectorKey + ".0"
		rateLimiterCID, rateLimiterDisclosure, err := resolveRateLimiterForSend(
			ctx,
			participant,
			parsedTokenPool,
			tokenDestSelectorKey,
			tokenDestSelectorNumericKey,
			resolveDisclosedByAddress,
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve outbound rate limiter for token transfer: %w", err)
		}
		disclosedRateLimiter = rateLimiterDisclosure

		adminParty := string(parsedTokenPool.InstrumentId.Admin)
		poolOwnerParty := string(parsedTokenPool.PoolOwner)
		holdings, err := listHoldingContracts(ctx, participant)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("list sender holdings: %w", err)
		}
		senderHoldingCIDs, _ := selectUnlockedHoldingCIDs(holdings, party, adminParty, string(parsedTokenPool.InstrumentId.Id))
		if len(senderHoldingCIDs) == 0 {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("no unlocked sender holdings found for %s/%s", parsedTokenPool.InstrumentId.Admin, parsedTokenPool.InstrumentId.Id)
		}
		senderInputHoldingIDs := make([]string, 0, len(senderHoldingCIDs))
		for _, cid := range senderHoldingCIDs {
			senderInputHoldingIDs = append(senderInputHoldingIDs, string(cid))
		}
		transferFactoryCID, transferFactoryDisclosures, transferFactoryCtx, err := getTransferFactoryFromScanProxy(
			ctx,
			participant,
			adminParty,
			party,
			poolOwnerParty,
			fields.TokenAmount.Amount.String(),
			adminParty,
			string(parsedTokenPool.InstrumentId.Id),
			senderInputHoldingIDs,
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get transfer factory for token transfer send: %w", err)
		}
		disclosedTransferFactoryContracts = transferFactoryDisclosures

		poolContext := common.CCIPContext{
			Values: types.TEXTMAP{
				"rate-limiter": common.AnyValue{AVContractId: &rateLimiterCID},
			},
		}
		tokenTransferInput = &ccipsender.TokenTransferInput{
			TokenPoolCid: tokenPoolCID,
			TokenInput: interfaces.TokenInput{
				TransferFactory: transferFactoryCID,
				ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
					Context: transferFactoryCtx,
					Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
				},
				TokenPoolHoldings: nil,
			},
			SenderInputCids:   senderHoldingCIDs,
			Amount:            types.NUMERIC(fields.TokenAmount.Amount.String()),
			TokenInstrumentId: parsedTokenPool.InstrumentId,
			PoolExtraContext:  poolContext,
		}
		disclosedTokenPool = tokenPoolDisclosure
		receiptIssuers = append(receiptIssuers, protocol.UnknownAddress(tokenPoolAddress.Bytes()))
	}

	sendContext := common.CCIPContext{
		Values: types.TEXTMAP{
			"on-ramp":              common.AnyValue{AVContractId: &onRampCID},
			"global-config":        common.AnyValue{AVContractId: &globalConfigCID},
			"token-admin-registry": common.AnyValue{AVContractId: &tokenAdminRegistryCID},
			"fee-quoter":           common.AnyValue{AVContractId: &feeQuoterCID},
			"rmn-remote":           common.AnyValue{AVContractId: &rmnRemoteCID},
		},
	}

	sendArgs := ccipsender.Send{
		Context:           sendContext,
		RouterCid:         routerCID,
		DestChainSelector: types.NUMERIC(fmt.Sprintf("%d.0", dest)),
		Receiver:          types.TEXT(hex.EncodeToString(fields.Receiver)),
		Payload:           types.TEXT(hex.EncodeToString(fields.Data)),
		ExtraArgs: ccipsender.CantonExtraArgsV1{
			GasLimit:           types.INT64(opts.ExecutionGasLimit),
			BlockConfirmations: nil,
			SenderRequiredCCVs: senderRequiredCCVs,
			ExecutorCid:        executorCID,
			ExecutorArgs:       nil,
			TokenReceiver:      nil,
			TokenArgs:          types.TEXT(""),
		},
		FeeToken: feeTokenInstrument,
		FeeTokenInput: interfaces.TokenInput{
			TransferFactory: routerCID,
			ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
				Context: splice_api_token_metadata_v1.ChoiceContext{Values: types.TEXTMAP{}},
				Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
			},
			TokenPoolHoldings: nil,
		},
		FeeTokenHoldingCids: nil,
		TokenTransfer:       tokenTransferInput,
		CcvSendInputs:       ccvSendInputs,
	}
	sendArgsMap := sendArgs.ToMap()
	if onRampCID == "" || globalConfigCID == "" || tokenAdminRegistryCID == "" || feeQuoterCID == "" || rmnRemoteCID == "" || routerCID == "" || ccipSenderCID == "" || executorCID == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
			"empty contract ID before send: ccipSender=%q router=%q executor=%q onRamp=%q globalConfig=%q tokenAdminRegistry=%q feeQuoter=%q rmnRemote=%q",
			ccipSenderCID, routerCID, executorCID, onRampCID, globalConfigCID, tokenAdminRegistryCID, feeQuoterCID, rmnRemoteCID,
		)
	}

	disclosedContracts := make([]*ledgerv2.DisclosedContract, 0, 9+len(disclosedVerifierContracts))
	disclosedContracts = append(disclosedContracts,
		disclosedCCIPSender,
		disclosedExecutor,
		disclosedRouter,
		disclosedOnRamp,
		disclosedGlobalConfig,
		disclosedTokenAdminRegistry,
		disclosedRMNRemote,
		disclosedFeeQuoter,
	)
	if disclosedTokenPool != nil {
		disclosedContracts = append(disclosedContracts, disclosedTokenPool)
	}
	if disclosedRateLimiter != nil {
		disclosedContracts = append(disclosedContracts, disclosedRateLimiter)
	}
	disclosedContracts = append(disclosedContracts, disclosedTransferFactoryContracts...)
	disclosedContracts = append(disclosedContracts, disclosedVerifierContracts...)
	for _, dc := range disclosedContracts {
		if dc != nil && dc.GetContractId() == "" {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("empty disclosed contract ID before ccip sender send")
		}
	}

	sendRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
					TemplateId:     &ledgerv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					ContractId:     string(ccipSenderCID),
					Choice:         "Send",
					ChoiceArgument: ledger.MapToValue(sendArgsMap),
				}},
			}},
			ActAs:              []string{party},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("submit ccip send via ccip sender: %w", err)
	}

	var messageID [32]byte
	var eventMessageID [32]byte
	var computedMessageID [32]byte
	var decodedMessage protocol.Message
	foundCCIPMessageSent := false
	foundEventMessageID := false
	foundEncodedMessage := false
	for _, event := range sendRes.GetTransaction().GetEvents() {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil || created.GetTemplateId().GetEntityName() != "CCIPMessageSent" {
			continue
		}
		foundCCIPMessageSent = true
		fields := created.GetCreateArguments().GetFields()
		for _, f := range fields {
			if f.GetLabel() != "event" || f.GetValue().GetRecord() == nil {
				continue
			}
			for _, eventField := range f.GetValue().GetRecord().GetFields() {
				switch eventField.GetLabel() {
				case "messageId":
					decoded, err := hex.DecodeString(eventField.GetValue().GetText())
					if err != nil || len(decoded) != 32 {
						return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("decode messageId from CCIPMessageSent event: %w", err)
					}
					copy(eventMessageID[:], decoded)
					foundEventMessageID = true
				case "encodedMessage":
					encodedMessage, err := hex.DecodeString(eventField.GetValue().GetText())
					if err != nil {
						return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("decode encodedMessage from CCIPMessageSent event: %w", err)
					}
					decodedMessagePtr, err := protocol.DecodeMessage(encodedMessage)
					if err != nil {
						return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("decode protocol message from encodedMessage: %w", err)
					}
					decodedMessage = *decodedMessagePtr
					computedHash := gethcrypto.Keccak256(encodedMessage)
					if len(computedHash) != len(computedMessageID) {
						return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("computed encodedMessage hash has invalid length: %d", len(computedHash))
					}
					copy(computedMessageID[:], computedHash)
					foundEncodedMessage = true
				}
			}
		}

		break
	}
	if !foundCCIPMessageSent {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("no CCIPMessageSent event found in sender transaction")
	}
	if foundEncodedMessage {
		// Prefer the sequence from the encoded message since local process state can be stale
		// across repeated test runs against a long-lived environment.
		seqNo = uint64(decodedMessage.SequenceNumber)
	}
	if foundEncodedMessage {
		messageID = computedMessageID
	} else if foundEventMessageID {
		messageID = eventMessageID
	} else {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("CCIPMessageSent event missing both messageId and encodedMessage")
	}

	event := cciptestinterfaces.MessageSentEvent{
		MessageID:      messageID,
		ReceiptIssuers: receiptIssuers,
	}
	if foundEncodedMessage {
		msg := decodedMessage
		event.Message = &msg
	}
	c.nextSeq = seqNo
	c.lastSentDest = dest
	c.lastSentSeq = seqNo
	c.lastSentEvent = event
	c.lastSentMessage = decodedMessage
	c.lastSentVerifierDestAddress = fallbackVerifierDestAddress
	c.lastSentVerifierBlob = append([]byte(nil), fallbackVerifierBlob...)
	c.lastSentHasVerificationInputs = foundEncodedMessage && len(fallbackVerifierDestAddress) > 0 && len(fallbackVerifierBlob) > 0

	return event, nil
}

// SendMessageWithNonce implements cciptestinterfaces.CCIP17.
func (c *Chain) SendMessageWithNonce(ctx context.Context, dest uint64, fields cciptestinterfaces.MessageFields, opts cciptestinterfaces.MessageOptions, sender *bind.TransactOpts, nonce *uint64, disableTokenAmountCheck bool) (cciptestinterfaces.MessageSentEvent, error) {
	return c.SendMessage(ctx, dest, fields, opts)
}

// WaitOneExecEventBySeqNo implements cciptestinterfaces.CCIP17.
func (c *Chain) WaitOneExecEventBySeqNo(ctx context.Context, from, seq uint64, timeout time.Duration) (cciptestinterfaces.ExecutionStateChangedEvent, error) {
	return cciptestinterfaces.ExecutionStateChangedEvent{}, nil // TODO: implement
}

// WaitOneSentEventBySeqNo implements cciptestinterfaces.CCIP17.
func (c *Chain) WaitOneSentEventBySeqNo(ctx context.Context, to, seq uint64, timeout time.Duration) (cciptestinterfaces.MessageSentEvent, error) {
	deadline := time.Now().Add(timeout)

	for {
		if c.lastSentDest == to && c.lastSentSeq == seq {
			return c.lastSentEvent, nil
		}

		if time.Now().After(deadline) {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("timed out waiting for sent event: dest=%d seq=%d", to, seq)
		}

		select {
		case <-ctx.Done():
			return cciptestinterfaces.MessageSentEvent{}, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// LastSentVerificationInput returns locally cached message + verifier proof input for fallback manual execution.
func (c *Chain) LastSentVerificationInput() (protocol.Message, protocol.UnknownAddress, []byte, bool) {
	return c.lastSentMessage, c.lastSentVerifierDestAddress, append([]byte(nil), c.lastSentVerifierBlob...), c.lastSentHasVerificationInputs
}

func convertToDisclosedContract(active *ledgerv2.ActiveContract) *ledgerv2.DisclosedContract {
	if active == nil || active.GetCreatedEvent() == nil {
		return nil
	}

	created := active.GetCreatedEvent()

	return &ledgerv2.DisclosedContract{
		TemplateId:       created.GetTemplateId(),
		ContractId:       created.GetContractId(),
		CreatedEventBlob: created.GetCreatedEventBlob(),
		SynchronizerId:   active.GetSynchronizerId(),
	}
}
