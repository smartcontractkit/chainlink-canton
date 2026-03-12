package devenv

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
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
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/deployments"
	"github.com/smartcontractkit/chainlink-ccv/protocol"

	cantonadapters "github.com/smartcontractkit/chainlink-canton/ccip/devenv/adapters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
						OnRampAddress: "", // TODO ?
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
	const reverseLockReleaseQualifier = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"

	// Add synthetic destination pool refs needed by transfer-token lane configuration for Canton selector.
	// Avoid creating LockReleaseTokenPool refs here to prevent clashes with the real deployed lock/release pool.
	for i, combo := range devenvcommon.AllTokenCombinations() {
		addressRef := combo.DestPoolAddressRef()
		if addressRef.Type == datastore.ContractType(lock_release_token_pool.ContractType) {
			continue
		}
		err = runningDS.AddressRefStore.Add(datastore.AddressRef{
			Address:       contracts.MustNewInstanceID("dst-token-pool-" + strconv.Itoa(i)).RawInstanceAddress(types.PARTY(participant.PartyID)).InstanceAddress().Hex(),
			Type:          addressRef.Type,
			Version:       addressRef.Version,
			Qualifier:     addressRef.Qualifier,
			ChainSelector: selector,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add dst token pool address ref: %w", err)
		}
	}

	// Deploy and TAR-register the default lock/release pool for Canton source token transfers.
	chainPoolConfigs := types.GENMAP{}
	chainFeeConfigs := types.GENMAP{}
	remoteTokens := types.GENMAP{}
	defaultCCVRef, err := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(committee_verifier.ContractType),
		committee_verifier.Version,
		devenvcommon.DefaultCommitteeVerifierQualifier,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to get default committee verifier for lock/release pool config: %w", err)
	}
	outboundCCVs := make([]common.RawInstanceAddress, 0, 1)
	if labels := defaultCCVRef.Labels.List(); len(labels) > 0 {
		rawAddr, parseErr := contracts.RawInstanceAddressFromString(labels[0])
		if parseErr == nil {
			outboundCCVs = append(outboundCCVs, rawAddr.Binding())
		}
	}
	// For Canton -> EVM sends from LockRelease source pool, pick the EVM token/pool pair
	// whose destination counterpart is LockRelease.
	const defaultBurnMintQualifier = reverseLockReleaseQualifier
	sourceDS := runningDS.Addresses()
	for _, bc := range env.BlockChains.All() {
		remoteSelector := bc.ChainSelector()
		if remoteSelector == selector {
			continue
		}
		selectorKey := strconv.FormatUint(remoteSelector, 10)
		chainFeeConfigs[selectorKey] = lockreleasetokenpool.PoolFeeConfig{
			FeeUSDCents:       types.NUMERIC("0"),
			DestGasOverhead:   types.INT64(0),
			DestBytesOverhead: types.INT64(0),
		}

		var tokenRef *datastore.AddressRef
		var firstBurnMint *datastore.AddressRef
		var poolRef *datastore.AddressRef
		var firstBurnMintPool *datastore.AddressRef
		collectTokenRef := func(store datastore.AddressRefStore) {
			for _, candidate := range store.Filter(
				datastore.AddressRefByChainSelector(remoteSelector),
				datastore.AddressRefByType(datastore.ContractType("BurnMintERC20WithDrip")),
			) {
				if firstBurnMint == nil {
					c := candidate
					firstBurnMint = &c
				}
				if candidate.Qualifier == defaultBurnMintQualifier {
					c := candidate
					tokenRef = &c
					return
				}
			}
		}
		collectPoolRef := func(store datastore.AddressRefStore) {
			for _, candidate := range store.Filter(
				datastore.AddressRefByChainSelector(remoteSelector),
				datastore.AddressRefByType(datastore.ContractType("BurnMintTokenPool")),
			) {
				if firstBurnMintPool == nil {
					c := candidate
					firstBurnMintPool = &c
				}
				if candidate.Qualifier == defaultBurnMintQualifier {
					c := candidate
					poolRef = &c
					return
				}
			}
		}
		collectTokenRef(sourceDS)
		collectPoolRef(sourceDS)
		if tokenRef == nil && env.DataStore != nil {
			collectTokenRef(env.DataStore.Addresses())
		}
		if poolRef == nil && env.DataStore != nil {
			collectPoolRef(env.DataStore.Addresses())
		}
		if tokenRef == nil {
			tokenRef = firstBurnMint
		}
		if poolRef == nil {
			poolRef = firstBurnMintPool
		}
		remotePools := []types.TEXT{}
		if poolRef != nil {
			remotePools = []types.TEXT{types.TEXT(canonicalCantonRemotePoolHex(poolRef.Address))}
		}
		chainPoolConfigs[selectorKey] = lockreleasetokenpool.ChainPoolConfig{
			InboundCCVs:  outboundCCVs,
			OutboundCCVs: outboundCCVs,
			RemotePools:  remotePools,
		}
		if tokenRef == nil {
			continue
		}
		remoteTokens[selectorKey] = strings.TrimPrefix(tokenRef.Address, "0x")
	}
	relativeHours := types.INT64(24)
	lockPoolInstrumentAdmin := participant.PartyID
	if metadataClient, metadataErr := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			token, tokenErr := participant.TokenSource.Token()
			if tokenErr != nil {
				return fmt.Errorf("retrieve participant token: %w", tokenErr)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			return nil
		}),
	); metadataErr == nil {
		if registryResp, registryErr := metadataClient.GetRegistryInfoWithResponse(ctx); registryErr == nil && registryResp.StatusCode() == http.StatusOK && registryResp.JSON200 != nil && registryResp.JSON200.AdminId != "" {
			lockPoolInstrumentAdmin = registryResp.JSON200.AdminId
		}
	}
	lockPoolTemplate := lockreleasetokenpool.LockReleaseTokenPool{
		CcipOwner: types.PARTY(participant.PartyID),
		PoolOwner: types.PARTY(participant.PartyID),
		InstrumentId: splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(lockPoolInstrumentAdmin),
			Id:    types.TEXT("Amulet"),
		},
		Decimals:         types.INT64(18),
		ChainPoolConfigs: chainPoolConfigs,
		ChainFeeConfigs:  chainFeeConfigs,
		RemoteTokens:     remoteTokens,
		PoolReceiveContext: common.CCIPContext{
			Values: types.TEXTMAP{},
		},
		TransferTimeout: lockreleasetokenpool.TransferTimeout{
			RelativeHours: &relativeHours,
		},
	}
	lockPoolQualifier := defaultLockReleaseQualifier
	lockPoolOutAddress := ""
	lockPoolOutLabels := datastore.NewLabelSet()
	existingLockPoolRefFromEnv := false
	if env.DataStore != nil {
		existingRef, getErr := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
			selector,
			datastore.ContractType("LockReleaseTokenPool"),
			lock_release_token_pool.Version,
			defaultLockReleaseQualifier,
		))
		if getErr == nil && existingRef.Address != "" {
			existingLockPoolRefFromEnv = true
			lockPoolOutAddress = existingRef.Address
			lockPoolOutLabels = existingRef.Labels
		}
	}
	if !existingLockPoolRefFromEnv {
		lockPoolOut, deployErr := operations.ExecuteOperation(
			env.OperationsBundle,
			canton_lock_release_token_pool.Deploy,
			dependencies.CantonDeps{Chain: chain},
			contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
				ChainSelector: selector,
				Qualifier:     &lockPoolQualifier,
				ActAs:         []string{participant.PartyID},
				Template:      lockPoolTemplate,
				OwnerParty:    types.PARTY(participant.PartyID),
			},
		)
		if deployErr != nil {
			return nil, fmt.Errorf("failed to deploy canton lock/release pool: %w", deployErr)
		}
		lockPoolOutAddress = lockPoolOut.Output.Address
		lockPoolOutLabels = lockPoolOut.Output.Labels
		if err = runningDS.AddressRefStore.Add(lockPoolOut.Output); err != nil {
			return nil, fmt.Errorf("failed to store deployed canton lock/release pool address ref: %w", err)
		}
	}
	// Add alias expected by CCV token combinations/tests.
	if err = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
		Address:       lockPoolOutAddress,
		Labels:        lockPoolOutLabels,
		Type:          datastore.ContractType(lock_release_token_pool.ContractType),
		Version:       lock_release_token_pool.Version,
		Qualifier:     defaultLockReleaseQualifier,
		ChainSelector: selector,
	}); err != nil {
		return nil, fmt.Errorf("failed to upsert lock/release pool alias address ref: %w", err)
	}
	if err = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
		Address:       lockPoolOutAddress,
		Labels:        lockPoolOutLabels,
		Type:          datastore.ContractType("LockReleaseTokenPool"),
		Version:       lock_release_token_pool.Version,
		Qualifier:     defaultLockReleaseQualifier,
		ChainSelector: selector,
	}); err != nil {
		return nil, fmt.Errorf("failed to upsert lock/release pool canonical alias address ref: %w", err)
	}
	if err = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
		Address:       lockPoolOutAddress,
		Labels:        lockPoolOutLabels,
		Type:          datastore.ContractType(lock_release_token_pool.ContractType),
		Version:       lock_release_token_pool.Version,
		Qualifier:     reverseLockReleaseQualifier,
		ChainSelector: selector,
	}); err != nil {
		return nil, fmt.Errorf("failed to upsert reverse lock/release pool alias address ref: %w", err)
	}
	if err = runningDS.AddressRefStore.Upsert(datastore.AddressRef{
		Address:       lockPoolOutAddress,
		Labels:        lockPoolOutLabels,
		Type:          datastore.ContractType("LockReleaseTokenPool"),
		Version:       lock_release_token_pool.Version,
		Qualifier:     reverseLockReleaseQualifier,
		ChainSelector: selector,
	}); err != nil {
		return nil, fmt.Errorf("failed to upsert reverse lock/release pool canonical alias address ref: %w", err)
	}
	activePool, activeErr := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		contracts.HexToInstanceAddress(lockPoolOutAddress),
	)
	if activeErr != nil {
		return nil, fmt.Errorf("failed to find active deployed lock/release pool: %w", activeErr)
	}
	parsedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse deployed lock/release pool: %w", parseErr)
	}
	tokenAdminRegistryRef, tarErr := runningDS.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	))
	if tarErr != nil {
		return nil, fmt.Errorf("failed to get token admin registry for pool registration: %w", tarErr)
	}
	if !existingLockPoolRefFromEnv {
		_, regErr := operations.ExecuteSequence(
			env.OperationsBundle,
			sequences.RegisterTokenPool,
			dependencies.CantonDeps{Chain: chain},
			sequences.RegisterTokenPoolInput{
				TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
				InstrumentId:                      parsedPool.InstrumentId,
				PoolInstanceID:                    string(parsedPool.InstanceId),
				CcipParty:                         participant.PartyID,
				PoolOwnerParty:                    participant.PartyID,
			},
		)
		if regErr != nil {
			return nil, fmt.Errorf("failed to register lock/release pool in token admin registry: %w", regErr)
		}
	}
	// Pre-seed real AMT holdings for the lock/release pool owner so e2e token flows use
	// existing liquidity (same model as ccip_execute_token_test.go) and never mint at execute time.
	if seedErr := seedAMTLiquidity(ctx, participant, participant.PartyID, "1000000.00"); seedErr != nil {
		return nil, fmt.Errorf("failed to pre-seed AMT liquidity for lock/release pool owner: %w", seedErr)
	}

	// If the pool has no outbound rate limiters but has remote chain configs, create rate limiters
	// for each remote chain and update the pool so Canton→EVM token sends can exercise LockOrBurn.
	if len(parsedPool.OutboundRateLimiters) == 0 && len(parsedPool.ChainPoolConfigs) > 0 {
		poolOwnerParty := string(parsedPool.PoolOwner)
		poolInstanceId := string(parsedPool.InstanceId)
		newOutbound := make(types.GENMAP)
		newInbound := make(types.GENMAP)
		var createCommands []*ledgerv2.Command
		for selectorKey := range parsedPool.ChainPoolConfigs {
			// selectorKey may be DAML numeric text (e.g. "1234.") and instance IDs reject "."
			selectorKeyForInstanceID := strings.ReplaceAll(selectorKey, ".", "-")
			outboundInstanceId := "devenv-outbound-rl-" + selectorKeyForInstanceID
			inboundInstanceId := "devenv-inbound-rl-" + selectorKeyForInstanceID
			nowMicro := time.Now().UnixMicro()
			createCommands = append(createCommands,
				&ledgerv2.Command{Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
						{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: outboundInstanceId}}},
						{Label: "poolInstanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: poolInstanceId}}},
						{Label: "poolOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: poolOwnerParty}}},
						{Label: "remoteChainSelector", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: selectorKey}}},
						{Label: "direction", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Enum{Enum: &ledgerv2.Enum{Constructor: "RateLimitDirection_Outbound"}}}},
						{Label: "mode", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Enum{Enum: &ledgerv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
						{Label: "isEnabled", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: true}}},
						{Label: "capacity", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "999999999999999999"}}},
						{Label: "rate", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "999999999999999999"}}},
						{Label: "tokens", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "0"}}},
						{Label: "lastUpdated", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Timestamp{Timestamp: nowMicro}}},
					}},
				}}},
				&ledgerv2.Command{Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
						{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: inboundInstanceId}}},
						{Label: "poolInstanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: poolInstanceId}}},
						{Label: "poolOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: poolOwnerParty}}},
						{Label: "remoteChainSelector", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: selectorKey}}},
						{Label: "direction", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Enum{Enum: &ledgerv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Enum{Enum: &ledgerv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
						{Label: "isEnabled", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: true}}},
						{Label: "capacity", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "999999999999999999"}}},
						{Label: "rate", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "999999999999999999"}}},
						{Label: "tokens", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "0"}}},
						{Label: "lastUpdated", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Timestamp{Timestamp: nowMicro}}},
					}},
				}}},
			)
			outboundRaw := contracts.MustNewInstanceID(outboundInstanceId).RawInstanceAddress(types.PARTY(parsedPool.PoolOwner))
			inboundRaw := contracts.MustNewInstanceID(inboundInstanceId).RawInstanceAddress(types.PARTY(parsedPool.PoolOwner))
			newOutbound[selectorKey] = outboundRaw.Binding()
			newInbound[selectorKey] = inboundRaw.Binding()
		}
		_, createErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
			Commands: &ledgerv2.Commands{
				CommandId: uuid.New().String(),
				Commands:  createCommands,
				ActAs:     []string{poolOwnerParty},
			},
		})
		if createErr != nil {
			return nil, fmt.Errorf("create rate limiters for lock/release pool: %w", createErr)
		}
		updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateRateLimiters{
			NewOutboundRateLimiters: newOutbound,
			NewInboundRateLimiters:  newInbound,
		}
		poolContractId := activePool.GetCreatedEvent().GetContractId()
		_, exerciseErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
			Commands: &ledgerv2.Commands{
				CommandId: uuid.New().String(),
				Commands: []*ledgerv2.Command{{
					Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
						TemplateId:     activePool.GetCreatedEvent().GetTemplateId(),
						ContractId:     poolContractId,
						Choice:         "LockReleaseTokenPool_UpdateRateLimiters",
						ChoiceArgument: ledger.MapToValue(updateArgs.ToMap()),
					}},
				}},
				ActAs: []string{poolOwnerParty},
			},
		})
		if exerciseErr != nil {
			return nil, fmt.Errorf("update lock/release pool rate limiters: %w", exerciseErr)
		}
		l.Info().Msg("Configured outbound and inbound rate limiters for Canton lock/release pool")
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
			DefaultInboundCCVs:       nil,
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

	// Keep lock/release remotePools aligned with currently deployed remote BurnMint pools.
	// This must happen at lane-config time because remote refs may not exist yet during initial deploy.
	const defaultLockReleaseQualifier = "TEST (LockReleaseTokenPool 1.7.0 [default] to BurnMintTokenPool 1.7.0 [default])"
	const reverseLockReleaseQualifier = "TEST (BurnMintTokenPool 1.7.0 [default] to LockReleaseTokenPool 1.7.0 [default])"
	chain := env.BlockChains.CantonChains()[selector]
	participant := chain.Participants[0]
	activePools := make([]*ledgerv2.ActiveContract, 0, 2)
	if lockPoolRef, lockPoolErr := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType("LockReleaseTokenPool"),
		lock_release_token_pool.Version,
		defaultLockReleaseQualifier,
	)); lockPoolErr == nil && lockPoolRef.Address != "" {
		lockPoolAddress := contracts.HexToInstanceAddress(lockPoolRef.Address)
		activePool, activePoolErr := contract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
			lockPoolAddress,
		)
		if activePoolErr != nil {
			return fmt.Errorf("find active lock/release pool for lane remotePools update: %w", activePoolErr)
		}
		activePools = append(activePools, activePool)
	}
	if len(activePools) == 0 {
		ledgerEnd, endErr := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
		if endErr != nil {
			return fmt.Errorf("get ledger end for lock/release pool lane update lookup: %w", endErr)
		}
		stream, streamErr := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
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
		if streamErr != nil {
			return fmt.Errorf("query active lock/release pools for lane remotePools update: %w", streamErr)
		}
		defer stream.CloseSend()
		for {
			resp, recvErr := stream.Recv()
			if errors.Is(recvErr, io.EOF) {
				break
			}
			if recvErr != nil {
				return fmt.Errorf("receive active lock/release pools for lane remotePools update: %w", recvErr)
			}
			entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
				continue
			}
			activePools = append(activePools, entry.ActiveContract)
		}
	}

	outboundCCVs := []common.RawInstanceAddress{}
	if committeeVerifierRawAddr != "" {
		outboundCCVs = append(outboundCCVs, committeeVerifierRawAddr.Binding())
	}
	for _, activePool := range activePools {
		parsedPool, parseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
		if parseErr != nil {
			return fmt.Errorf("parse active lock/release pool for lane remotePools update: %w", parseErr)
		}
		updatedChainPoolConfigs := types.GENMAP{}
		for k, v := range parsedPool.ChainPoolConfigs {
			updatedChainPoolConfigs[k] = v
		}
		needsPoolUpdate := false
		for _, remoteSelector := range remoteSelectors {
			remotePoolHex := ""
			var fallbackPool *datastore.AddressRef
			for _, candidate := range env.DataStore.Addresses().Filter(
				datastore.AddressRefByChainSelector(remoteSelector),
				datastore.AddressRefByType(datastore.ContractType("BurnMintTokenPool")),
			) {
				cand := candidate
				if fallbackPool == nil {
					fallbackPool = &cand
				}
				if candidate.Qualifier == reverseLockReleaseQualifier {
					remotePoolHex = candidate.Address
					break
				}
			}
			if remotePoolHex == "" && fallbackPool != nil {
				remotePoolHex = fallbackPool.Address
			}
			if remotePoolHex == "" {
				continue
			}
			selectorKey := strconv.FormatUint(remoteSelector, 10)
			updatedChainPoolConfigs[selectorKey] = lockreleasetokenpool.ChainPoolConfig{
				InboundCCVs:  outboundCCVs,
				OutboundCCVs: outboundCCVs,
				RemotePools:  []types.TEXT{types.TEXT(canonicalCantonRemotePoolHex(remotePoolHex))},
			}
			needsPoolUpdate = true
		}
		if !needsPoolUpdate {
			continue
		}
		updateArgs := lockreleasetokenpool.LockReleaseTokenPoolUpdateChainPoolConfigs{
			NewChainPoolConfigs: updatedChainPoolConfigs,
		}
		_, exerciseErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
			Commands: &ledgerv2.Commands{
				CommandId: uuid.New().String(),
				Commands: []*ledgerv2.Command{{
					Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
						TemplateId:     activePool.GetCreatedEvent().GetTemplateId(),
						ContractId:     activePool.GetCreatedEvent().GetContractId(),
						Choice:         "LockReleaseTokenPool_UpdateChainPoolConfigs",
						ChoiceArgument: ledger.MapToValue(updateArgs.ToMap()),
					}},
				}},
				ActAs: []string{string(parsedPool.PoolOwner)},
			},
		})
		if exerciseErr != nil {
			return fmt.Errorf("update lock/release pool remotePools for connected lanes: %w", exerciseErr)
		}
	}

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

	for key, outboundRateLimiterRaw := range parsedTokenPool.OutboundRateLimiters {
		if normalizeNumericForCompare(key) != selectorNorm {
			continue
		}
		if rawAddress, rawErr := extractRawRateLimiterAddress(outboundRateLimiterRaw); rawErr == nil && expectedRawAddress == "" {
			expectedRawAddress = rawAddress
		}
		cid, disclosure, err := resolveRateLimiterFromRawAddressForSend(outboundRateLimiterRaw, resolveDisclosedByAddress)
		if err == nil {
			return cid, disclosure, nil
		}
	}

	if outboundRateLimiterRaw, ok := parsedTokenPool.OutboundRateLimiters[destSelectorKey]; ok {
		if rawAddress, rawErr := extractRawRateLimiterAddress(outboundRateLimiterRaw); rawErr == nil {
			expectedRawAddress = rawAddress
		}
		cid, disclosure, err := resolveRateLimiterFromRawAddressForSend(outboundRateLimiterRaw, resolveDisclosedByAddress)
		if err == nil {
			return cid, disclosure, nil
		}
	}
	if outboundRateLimiterRaw, ok := parsedTokenPool.OutboundRateLimiters[destSelectorNumericKey]; ok {
		if rawAddress, rawErr := extractRawRateLimiterAddress(outboundRateLimiterRaw); rawErr == nil && expectedRawAddress == "" {
			expectedRawAddress = rawAddress
		}
		cid, disclosure, err := resolveRateLimiterFromRawAddressForSend(outboundRateLimiterRaw, resolveDisclosedByAddress)
		if err == nil {
			return cid, disclosure, nil
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

	want := strings.ToLower(strings.TrimPrefix(tokenAddress.String(), "0x"))
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

		for _, remoteTokenValue := range parsed.RemoteTokens {
			remoteToken := strings.ToLower(strings.TrimPrefix(fmt.Sprint(remoteTokenValue), "0x"))
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
	if strings.TrimRight(parts[1], "0") != "" {
		return nil, fmt.Errorf("non-integer DAML numeric %q cannot be represented as big.Int", v)
	}
	intPart := parts[0]
	if intPart == "" {
		intPart = "0"
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
	routerAddress, err := c.DeployPerPartyRouter(ctx, c.participantIndexForParty(party), party)
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
					TemplateId: &ledgerv2.Identifier{PackageId: "#ccip-test", ModuleName: "TestExecutor", EntityName: "TestExecutor"},
					CreateArguments: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
						{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: executorInstanceID.String()}}},
						{Label: "owner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: party}}},
						{Label: "minBlockConfirmations", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
						{Label: "maxCCVsPerMsg", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 10}}},
						{Label: "ccvAllowlistEnabled", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: false}}},
						{Label: "allowedCCVs", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
						{
							Label: "remoteChainFeeUSDCents",
							Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{
								Entries: []*ledgerv2.GenMap_Entry{
									{
										Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: destSelectorText}},
										Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "0.0"}},
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
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("deploy test executor contract: %w", err)
	}
	executorInstanceAddress := executorInstanceID.RawInstanceAddress(types.PARTY(party)).InstanceAddress()
	executorCID, disclosedExecutor, err := resolveDisclosedByAddress("#ccip-test:TestExecutor:TestExecutor", executorInstanceAddress)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve test executor disclosed contract: %w", err)
	}
	if executorCID == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolved empty test executor contract ID")
	}

	// Keep fee token setup local and deterministic for devenv sends.
	feeTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(party),
		Id:    types.TEXT("devenv-fee-token"),
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
								PremiumMultiplier: types.NUMERIC("1.0"),
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
								{InstrumentId: feeTokenInstrument, UsdPerToken: types.NUMERIC("1.0")},
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
		destSelectorKey := strconv.FormatUint(dest, 10)
		destSelectorNumericKey := destSelectorKey + "."
		_, hasDestRemoteToken := parsedTokenPool.RemoteTokens[destSelectorKey]
		if !hasDestRemoteToken {
			_, hasDestRemoteToken = parsedTokenPool.RemoteTokens[destSelectorNumericKey]
		}
		_, hasDestPoolCfg := parsedTokenPool.ChainPoolConfigs[destSelectorKey]
		if !hasDestPoolCfg {
			_, hasDestPoolCfg = parsedTokenPool.ChainPoolConfigs[destSelectorNumericKey]
		}
		_, hasDestOutboundRateLimiter := parsedTokenPool.OutboundRateLimiters[destSelectorKey]
		if !hasDestOutboundRateLimiter {
			_, hasDestOutboundRateLimiter = parsedTokenPool.OutboundRateLimiters[destSelectorNumericKey]
		}
		_, hasDestInboundRateLimiter := parsedTokenPool.InboundRateLimiters[destSelectorKey]
		if !hasDestInboundRateLimiter {
			_, hasDestInboundRateLimiter = parsedTokenPool.InboundRateLimiters[destSelectorNumericKey]
		}
		if !hasDestRemoteToken || !hasDestPoolCfg {
			keys := func(m types.GENMAP) []string {
				out := make([]string, 0, len(m))
				for k := range m {
					out = append(out, k)
				}
				sort.Strings(out)
				return out
			}
			if os.Getenv("CCIP_DEBUG_POOL_CONFIG") == "1" {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
					"missing lock/release pool destination config for selector %s: hasRemoteToken=%t hasPoolCfg=%t remoteTokenKeys=%v poolCfgKeys=%v",
					destSelectorKey,
					hasDestRemoteToken,
					hasDestPoolCfg,
					keys(parsedTokenPool.RemoteTokens),
					keys(parsedTokenPool.ChainPoolConfigs),
				)
			}
			if !hasDestRemoteToken && hasDestPoolCfg {
				const preferredQualifierRTPatch = remoteDestBurnMintQualifier
				resolveRemoteTokenForDest := func(destSelector uint64) (string, error) {
					var fallback *datastore.AddressRef
					for _, candidate := range c.e.DataStore.Addresses().Filter(
						datastore.AddressRefByChainSelector(destSelector),
						datastore.AddressRefByType(datastore.ContractType("BurnMintERC20WithDrip")),
					) {
						cand := candidate
						if fallback == nil {
							fallback = &cand
						}
						if candidate.Qualifier == preferredQualifierRTPatch {
							return strings.ToLower(strings.TrimPrefix(candidate.Address, "0x")), nil
						}
					}
					if fallback == nil {
						return "", fmt.Errorf("no BurnMintERC20WithDrip token refs found for destination selector %d", destSelector)
					}

					return strings.ToLower(strings.TrimPrefix(fallback.Address, "0x")), nil
				}
				remoteTokenHex, rtErr := resolveRemoteTokenForDest(dest)
				if rtErr != nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve remote token address for destination selector %d: %w", dest, rtErr)
				}
				created := activeTokenPool.GetCreatedEvent()
				if created == nil || created.GetCreateArguments() == nil || created.GetTemplateId() == nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token pool created event is missing required fields for remote token patch")
				}
				replaceRemoteTokens := false
				newFields := make([]*ledgerv2.RecordField, 0, len(created.GetCreateArguments().GetFields()))
				for _, f := range created.GetCreateArguments().GetFields() {
					if f.GetLabel() != "remoteTokens" {
						newFields = append(newFields, f)
						continue
					}
					replaceRemoteTokens = true
					var entries []*ledgerv2.GenMap_Entry
					if gm := f.GetValue().GetGenMap(); gm != nil {
						entries = append(entries, gm.GetEntries()...)
					}
					updated := false
					for _, e := range entries {
						keyNumeric := e.GetKey().GetNumeric()
						if keyNumeric == destSelectorKey || keyNumeric == destSelectorNumericKey {
							e.Value = &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: remoteTokenHex}}
							updated = true
							break
						}
					}
					if !updated {
						entries = append(entries, &ledgerv2.GenMap_Entry{
							Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: destSelectorNumericKey}},
							Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: remoteTokenHex}},
						})
					}
					newFields = append(newFields, &ledgerv2.RecordField{
						Label: "remoteTokens",
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{Entries: entries}}},
					})
				}
				if !replaceRemoteTokens {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token pool create arguments missing remoteTokens field")
				}
				_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
					Commands: &ledgerv2.Commands{
						CommandId: uuid.New().String(),
						Commands: []*ledgerv2.Command{{
							Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
								TemplateId: created.GetTemplateId(),
								CreateArguments: &ledgerv2.Record{
									Fields: newFields,
								},
							}},
						}},
						ActAs: []string{string(parsedTokenPool.PoolOwner)},
					},
				})
				if err != nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("patch token pool remoteTokens for destination selector %d: %w", dest, err)
				}
				activeTokenPool, err = findLatestActiveByAddress(
					lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
					tokenPoolAddress,
				)
				if err != nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve patched active token pool by configured address after remote token patch: %w", err)
				}
				if activeTokenPool == nil || activeTokenPool.GetCreatedEvent() == nil || activeTokenPool.GetCreatedEvent().GetContractId() == "" {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("patched token pool missing created event or contract id after remote token patch")
				}
				tokenPoolCID = types.CONTRACT_ID(activeTokenPool.GetCreatedEvent().GetContractId())
				tokenPoolDisclosure = convertToDisclosedContract(activeTokenPool)
				if tokenPoolDisclosure == nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to build disclosed token pool after remote token patch")
				}
				parsedTokenPool, err = bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activeTokenPool.GetCreatedEvent())
				if err != nil {
					return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse patched token pool after remote token patch: %w", err)
				}
				hasDestRemoteToken = true
			}
			if hasDestRemoteToken && hasDestPoolCfg {
				goto tokenPoolConfigReady
			}

			const preferredQualifier = remoteDestBurnMintQualifier
			resolveRemoteTokenForDest := func(destSelector uint64) (string, error) {
				var fallback *datastore.AddressRef
				for _, candidate := range c.e.DataStore.Addresses().Filter(
					datastore.AddressRefByChainSelector(destSelector),
					datastore.AddressRefByType(datastore.ContractType("BurnMintERC20WithDrip")),
				) {
					cand := candidate
					if fallback == nil {
						fallback = &cand
					}
					if candidate.Qualifier == preferredQualifier {
						return strings.ToLower(strings.TrimPrefix(candidate.Address, "0x")), nil
					}
				}
				if fallback == nil {
					return "", fmt.Errorf("no BurnMintERC20WithDrip token refs found for destination selector %d", destSelector)
				}

				return strings.ToLower(strings.TrimPrefix(fallback.Address, "0x")), nil
			}
			resolveRemotePoolForDest := func(destSelector uint64) (string, error) {
				var fallback *datastore.AddressRef
				for _, candidate := range c.e.DataStore.Addresses().Filter(
					datastore.AddressRefByChainSelector(destSelector),
					datastore.AddressRefByType(datastore.ContractType("BurnMintTokenPool")),
				) {
					cand := candidate
					if fallback == nil {
						fallback = &cand
					}
					if candidate.Qualifier == remoteDestBurnMintQualifier {
						return strings.ToLower(strings.TrimPrefix(candidate.Address, "0x")), nil
					}
				}
				if fallback == nil {
					return "", fmt.Errorf("no BurnMintTokenPool refs found for destination selector %d", destSelector)
				}

				return strings.ToLower(strings.TrimPrefix(fallback.Address, "0x")), nil
			}
			remoteTokenHex, err := resolveRemoteTokenForDest(dest)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve remote token address for destination selector %d: %w", dest, err)
			}
			remotePoolHex, err := resolveRemotePoolForDest(dest)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve remote pool address for destination selector %d: %w", dest, err)
			}
			updatedChainPoolConfigs := types.GENMAP{}
			for k, v := range parsedTokenPool.ChainPoolConfigs {
				updatedChainPoolConfigs[k] = v
			}
			updatedChainPoolConfigs[destSelectorKey] = lockreleasetokenpool.ChainPoolConfig{
				InboundCCVs:  senderRequiredCCVs,
				OutboundCCVs: senderRequiredCCVs,
				RemotePools:  []types.TEXT{types.TEXT(canonicalCantonRemotePoolHex(remotePoolHex))},
			}
			updatedChainFeeConfigs := types.GENMAP{}
			for k, v := range parsedTokenPool.ChainFeeConfigs {
				updatedChainFeeConfigs[k] = v
			}
			updatedChainFeeConfigs[destSelectorKey] = lockreleasetokenpool.PoolFeeConfig{
				FeeUSDCents:       types.NUMERIC("0"),
				DestGasOverhead:   types.INT64(0),
				DestBytesOverhead: types.INT64(0),
			}
			updatedRemoteTokens := types.GENMAP{}
			for k, v := range parsedTokenPool.RemoteTokens {
				updatedRemoteTokens[k] = v
			}
			updatedRemoteTokens[destSelectorKey] = types.TEXT(remoteTokenHex)
			updatedOutboundRateLimiters := types.GENMAP{}
			for k, v := range parsedTokenPool.OutboundRateLimiters {
				updatedOutboundRateLimiters[k] = v
			}
			updatedInboundRateLimiters := types.GENMAP{}
			for k, v := range parsedTokenPool.InboundRateLimiters {
				updatedInboundRateLimiters[k] = v
			}
			hasDestOutboundRateLimiter := false
			if _, ok := updatedOutboundRateLimiters[destSelectorKey]; ok {
				hasDestOutboundRateLimiter = true
			} else if _, ok := updatedOutboundRateLimiters[destSelectorNumericKey]; ok {
				hasDestOutboundRateLimiter = true
			}
			hasDestInboundRateLimiter := false
			if _, ok := updatedInboundRateLimiters[destSelectorKey]; ok {
				hasDestInboundRateLimiter = true
			} else if _, ok := updatedInboundRateLimiters[destSelectorNumericKey]; ok {
				hasDestInboundRateLimiter = true
			}
			if !hasDestOutboundRateLimiter {
				var fallback any
				for _, v := range updatedOutboundRateLimiters {
					fallback = v
					break
				}
				if fallback != nil {
					updatedOutboundRateLimiters[destSelectorKey] = fallback
				}
			}
			if !hasDestInboundRateLimiter {
				var fallback any
				for _, v := range updatedInboundRateLimiters {
					fallback = v
					break
				}
				if fallback != nil {
					updatedInboundRateLimiters[destSelectorKey] = fallback
				}
			}
			updatedPool := lockreleasetokenpool.LockReleaseTokenPool{
				InstanceId:           parsedTokenPool.InstanceId,
				CcipOwner:            parsedTokenPool.CcipOwner,
				PoolOwner:            parsedTokenPool.PoolOwner,
				InstrumentId:         parsedTokenPool.InstrumentId,
				Decimals:             parsedTokenPool.Decimals,
				ChainPoolConfigs:     updatedChainPoolConfigs,
				ChainFeeConfigs:      updatedChainFeeConfigs,
				RemoteTokens:         updatedRemoteTokens,
				OutboundRateLimiters: updatedOutboundRateLimiters,
				InboundRateLimiters:  updatedInboundRateLimiters,
				PoolReceiveContext:   parsedTokenPool.PoolReceiveContext,
				TransferTimeout:      parsedTokenPool.TransferTimeout,
			}
			bundle := operations.NewBundle(
				func() context.Context { return context.Background() },
				c.e.Logger,
				operations.NewMemoryReporter(),
			)
			_, deployErr := operations.ExecuteOperation(
				bundle,
				canton_lock_release_token_pool.Deploy,
				dependencies.CantonDeps{Chain: c.chain},
				contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
					ChainSelector: c.chain.ChainSelector(),
					Qualifier: func() *string {
						q := tokenPoolQualifier
						return &q
					}(),
					ActAs:      []string{string(updatedPool.PoolOwner)},
					Template:   updatedPool,
					OwnerParty: updatedPool.PoolOwner,
				},
			)
			if deployErr != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("configure lock/release pool for destination selector %d: %w", dest, deployErr)
			}
			tokenPoolCID, tokenPoolDisclosure, err = resolveDisclosedByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), tokenPoolAddress)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve updated token pool disclosed contract: %w", err)
			}
			activeTokenPool, err = contract.FindActiveContractByInstanceAddress(
				ctx,
				participant.LedgerServices.State,
				participant.PartyID,
				lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
				tokenPoolAddress,
			)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve updated active token pool by configured address: %w", err)
			}
			parsedTokenPool, err = bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activeTokenPool.GetCreatedEvent())
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse updated token pool contract: %w", err)
			}
		}
		{
			created := activeTokenPool.GetCreatedEvent()
			if created == nil || created.GetCreateArguments() == nil || created.GetTemplateId() == nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token pool created event is missing required fields for rate limiter patch")
			}

			var outboundFallback any
			for k, v := range parsedTokenPool.OutboundRateLimiters {
				if k == destSelectorKey || k == destSelectorNumericKey {
					continue
				}
				outboundFallback = v
				break
			}
			if outboundFallback == nil {
				for _, v := range parsedTokenPool.OutboundRateLimiters {
					outboundFallback = v
					break
				}
			}
			if outboundFallback == nil {
				outboundKeys := make([]string, 0, len(parsedTokenPool.OutboundRateLimiters))
				for k := range parsedTokenPool.OutboundRateLimiters {
					outboundKeys = append(outboundKeys, k)
				}
				sort.Strings(outboundKeys)
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("missing outbound rate limiter config for destination selector %d (outboundKeys=%v)", dest, outboundKeys)
			}
			var inboundFallback any
			for k, v := range parsedTokenPool.InboundRateLimiters {
				if k == destSelectorKey || k == destSelectorNumericKey {
					continue
				}
				inboundFallback = v
				break
			}
			if inboundFallback == nil {
				for _, v := range parsedTokenPool.InboundRateLimiters {
					inboundFallback = v
					break
				}
			}
			if inboundFallback == nil {
				inboundKeys := make([]string, 0, len(parsedTokenPool.InboundRateLimiters))
				for k := range parsedTokenPool.InboundRateLimiters {
					inboundKeys = append(inboundKeys, k)
				}
				sort.Strings(inboundKeys)
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("missing inbound rate limiter config for destination selector %d (inboundKeys=%v)", dest, inboundKeys)
			}

			replaceOutbound := false
			replaceInbound := false
			newFields := make([]*ledgerv2.RecordField, 0, len(created.GetCreateArguments().GetFields()))
			for _, f := range created.GetCreateArguments().GetFields() {
				switch f.GetLabel() {
				case "outboundRateLimiters":
					replaceOutbound = true
					var entries []*ledgerv2.GenMap_Entry
					if gm := f.GetValue().GetGenMap(); gm != nil {
						entries = append(entries, gm.GetEntries()...)
					}
					updated := false
					for _, e := range entries {
						keyNumeric := e.GetKey().GetNumeric()
						if keyNumeric == destSelectorKey || keyNumeric == destSelectorNumericKey {
							e.Value = ledger.MapToValue(outboundFallback)
							updated = true
							break
						}
					}
					if !updated {
						entries = append(entries, &ledgerv2.GenMap_Entry{
							Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: destSelectorNumericKey}},
							Value: ledger.MapToValue(outboundFallback),
						})
					}
					newFields = append(newFields, &ledgerv2.RecordField{
						Label: "outboundRateLimiters",
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{Entries: entries}}},
					})
				case "inboundRateLimiters":
					replaceInbound = true
					var entries []*ledgerv2.GenMap_Entry
					if gm := f.GetValue().GetGenMap(); gm != nil {
						entries = append(entries, gm.GetEntries()...)
					}
					updated := false
					for _, e := range entries {
						keyNumeric := e.GetKey().GetNumeric()
						if keyNumeric == destSelectorKey || keyNumeric == destSelectorNumericKey {
							e.Value = ledger.MapToValue(inboundFallback)
							updated = true
							break
						}
					}
					if !updated {
						entries = append(entries, &ledgerv2.GenMap_Entry{
							Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: destSelectorNumericKey}},
							Value: ledger.MapToValue(inboundFallback),
						})
					}
					newFields = append(newFields, &ledgerv2.RecordField{
						Label: "inboundRateLimiters",
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{Entries: entries}}},
					})
				default:
					newFields = append(newFields, f)
				}
			}
			if !replaceOutbound || !replaceInbound {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token pool create arguments missing rate limiter fields")
			}

			_, err = participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
				Commands: &ledgerv2.Commands{
					CommandId: uuid.New().String(),
					Commands: []*ledgerv2.Command{{
						Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
							TemplateId: created.GetTemplateId(),
							CreateArguments: &ledgerv2.Record{
								Fields: newFields,
							},
						}},
					}},
					ActAs: []string{string(parsedTokenPool.PoolOwner)},
				},
			})
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("patch token pool rate limiters for destination selector %d: %w", dest, err)
			}
			tokenPoolCID, tokenPoolDisclosure, err = resolveDisclosedByAddress(
				lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
				tokenPoolAddress,
			)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve patched token pool disclosed contract after rate limiter patch: %w", err)
			}
			activeTokenPool, err = findLatestActiveByAddress(
				lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
				tokenPoolAddress,
			)
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve patched active token pool by configured address after rate limiter patch: %w", err)
			}
			parsedTokenPool, err = bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activeTokenPool.GetCreatedEvent())
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse patched token pool after rate limiter patch: %w", err)
			}
		}
	tokenPoolConfigReady:
		if fields.TokenAmount.Amount.Sign() == 0 {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("token amount must be greater than zero for token transfer")
		}
		tokenDestSelectorKey := strconv.FormatUint(dest, 10)
		tokenDestSelectorNumericKey := tokenDestSelectorKey + "."
		rateLimiterCID, rateLimiterDisclosure, err := resolveRateLimiterForSend(
			ctx,
			participant,
			parsedTokenPool,
			tokenDestSelectorKey,
			tokenDestSelectorNumericKey,
			resolveDisclosedByAddress,
		)
		if err != nil {
			findActiveRateLimiterRaw := func(direction common.RateLimitDirection) (any, error) {
				ledgerEnd, leErr := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
				if leErr != nil {
					return nil, fmt.Errorf("get ledger end for rate limiter fallback: %w", leErr)
				}
				stream, streamErr := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
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
				if streamErr != nil {
					return nil, fmt.Errorf("query active rate limiters for fallback: %w", streamErr)
				}
				defer stream.CloseSend()
				selectorNorm := normalizeNumericForCompare(tokenDestSelectorKey)
				var fallback any
				for {
					resp, recvErr := stream.Recv()
					if errors.Is(recvErr, io.EOF) {
						break
					}
					if recvErr != nil {
						return nil, fmt.Errorf("receive active rate limiters for fallback: %w", recvErr)
					}
					entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
					if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
						continue
					}
					parsedRateLimiter, parseErr := bindings.UnmarshalCreatedEvent[common.RateLimiter](entry.ActiveContract.GetCreatedEvent())
					if parseErr != nil {
						continue
					}
					if parsedRateLimiter.Direction != direction {
						continue
					}
					if string(parsedRateLimiter.PoolOwner) != string(parsedTokenPool.PoolOwner) || string(parsedRateLimiter.PoolInstanceId) != string(parsedTokenPool.InstanceId) {
						continue
					}
					raw := common.RawInstanceAddress{
						Unpack: types.TEXT(fmt.Sprintf("%s@%s", parsedRateLimiter.InstanceId, parsedRateLimiter.PoolOwner)),
					}
					if normalizeNumericForCompare(string(parsedRateLimiter.RemoteChainSelector)) == selectorNorm {
						return raw, nil
					}
					if fallback == nil {
						fallback = raw
					}
				}
				if fallback != nil {
					return fallback, nil
				}
				return nil, fmt.Errorf("no active rate limiter found for direction %s", direction)
			}

			outboundFallback, outErr := findActiveRateLimiterRaw(common.RateLimitDirectionRateLimitDirection_Outbound)
			if outErr != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve outbound rate limiter for token transfer: %w", err)
			}
			inboundFallback, inErr := findActiveRateLimiterRaw(common.RateLimitDirectionRateLimitDirection_Inbound)
			if inErr != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve outbound rate limiter for token transfer: %w", err)
			}
			created := activeTokenPool.GetCreatedEvent()
			if created == nil || created.GetCreateArguments() == nil || created.GetTemplateId() == nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve outbound rate limiter for token transfer: %w", err)
			}
			replaceOutbound := false
			replaceInbound := false
			newFields := make([]*ledgerv2.RecordField, 0, len(created.GetCreateArguments().GetFields()))
			for _, f := range created.GetCreateArguments().GetFields() {
				switch f.GetLabel() {
				case "outboundRateLimiters":
					replaceOutbound = true
					var entries []*ledgerv2.GenMap_Entry
					if gm := f.GetValue().GetGenMap(); gm != nil {
						entries = append(entries, gm.GetEntries()...)
					}
					updated := false
					for _, e := range entries {
						keyNumeric := e.GetKey().GetNumeric()
						if normalizeNumericForCompare(keyNumeric) == normalizeNumericForCompare(tokenDestSelectorKey) {
							e.Value = ledger.MapToValue(outboundFallback)
							updated = true
							break
						}
					}
					if !updated {
						entries = append(entries, &ledgerv2.GenMap_Entry{
							Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: tokenDestSelectorNumericKey}},
							Value: ledger.MapToValue(outboundFallback),
						})
					}
					newFields = append(newFields, &ledgerv2.RecordField{
						Label: "outboundRateLimiters",
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{Entries: entries}}},
					})
				case "inboundRateLimiters":
					replaceInbound = true
					var entries []*ledgerv2.GenMap_Entry
					if gm := f.GetValue().GetGenMap(); gm != nil {
						entries = append(entries, gm.GetEntries()...)
					}
					updated := false
					for _, e := range entries {
						keyNumeric := e.GetKey().GetNumeric()
						if normalizeNumericForCompare(keyNumeric) == normalizeNumericForCompare(tokenDestSelectorKey) {
							e.Value = ledger.MapToValue(inboundFallback)
							updated = true
							break
						}
					}
					if !updated {
						entries = append(entries, &ledgerv2.GenMap_Entry{
							Key:   &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: tokenDestSelectorNumericKey}},
							Value: ledger.MapToValue(inboundFallback),
						})
					}
					newFields = append(newFields, &ledgerv2.RecordField{
						Label: "inboundRateLimiters",
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_GenMap{GenMap: &ledgerv2.GenMap{Entries: entries}}},
					})
				default:
					newFields = append(newFields, f)
				}
			}
			if replaceOutbound && replaceInbound {
				_, createErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
					Commands: &ledgerv2.Commands{
						CommandId: uuid.New().String(),
						Commands: []*ledgerv2.Command{{
							Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
								TemplateId: created.GetTemplateId(),
								CreateArguments: &ledgerv2.Record{
									Fields: newFields,
								},
							}},
						}},
						ActAs: []string{string(parsedTokenPool.PoolOwner)},
					},
				})
				if createErr == nil {
					tokenPoolCID, tokenPoolDisclosure, _ = resolveDisclosedByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), tokenPoolAddress)
					activeTokenPool, _ = findLatestActiveByAddress(lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), tokenPoolAddress)
					if activeTokenPool != nil && activeTokenPool.GetCreatedEvent() != nil {
						if reparsed, reparseErr := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activeTokenPool.GetCreatedEvent()); reparseErr == nil {
							parsedTokenPool = reparsed
							rateLimiterCID, rateLimiterDisclosure, err = resolveRateLimiterForSend(
								ctx,
								participant,
								parsedTokenPool,
								tokenDestSelectorKey,
								tokenDestSelectorNumericKey,
								resolveDisclosedByAddress,
							)
						}
					}
				}
			}
			if err != nil {
				return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve outbound rate limiter for token transfer: %w", err)
			}
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
		DestChainSelector: types.NUMERIC(fmt.Sprintf("%d", dest)),
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
