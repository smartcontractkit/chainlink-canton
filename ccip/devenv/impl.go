package devenv

import (
	"bytes"
	"context"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipChangesets "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/chainreg"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	ccvservices "github.com/smartcontractkit/chainlink-ccv/build/devenv/services"
	ccipOffchain "github.com/smartcontractkit/chainlink-ccv/deployment"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	ccipcore "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonadapters "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	cantonchangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	executor2 "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TODO: move this to share between devenv and integration tests
// and define a const for LINK instrument
const AMTInstrument = types.TEXT("Amulet")

var (
	_                       cciptestinterfaces.CCIP17              = &Chain{}
	_                       cciptestinterfaces.CCIP17Configuration = &Chain{}
	_                       chainreg.ImplFactory                   = &ImplFactory{}
	cantonTokenPoolVersion                                         = semver.MustParse("2.0.0")
	cantonDeployDarPackages                                        = []contracts.Package{
		contracts.CCIPFactory,
		contracts.CCIPCommon,
		contracts.CCIPReceiver,
		contracts.CCIPOffRamp,
		contracts.CCIPOnRamp,
		contracts.CCIPTokenAdminRegistry,
		contracts.CCIPCommitteeVerifier,
		contracts.CCIPPerPartyRouter,
		contracts.CCIPFeeQuoter,
		contracts.CCIPRMN,
		contracts.CCIPSender,
		contracts.CCIPExecutor,
		contracts.CCIPTest,
		contracts.CCIPLockReleaseTokenPool,
	}
)

type poolConfigKey struct {
	poolType    datastore.ContractType
	poolVersion string
	qualifier   string
}

type ImplFactory struct{}

func NewImplFactory() *ImplFactory {
	return &ImplFactory{}
}

// New implements [chainimpl.ImplFactory].
func (i *ImplFactory) New(ctx context.Context, lggr zerolog.Logger, env *deployment.Environment, chainSelector uint64) (cciptestinterfaces.CCIP17, error) {
	return New(ctx, lggr, env, chainSelector)
}

// NewEmpty implements [chainimpl.ImplFactory].
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

func (i *ImplFactory) DefaultSignerKey(keys ccvservices.BootstrapKeys) string {
	return keys.ECDSAPublicKey
}

func (i *ImplFactory) DefaultFeeAggregator(env *deployment.Environment, chainSelector uint64) string {
	chain, ok := env.BlockChains.CantonChains()[chainSelector]
	if !ok || len(chain.Participants) == 0 {
		return ""
	}

	return chain.Participants[0].PartyID
}

func (i *ImplFactory) SupportsFunding() bool {
	return false
}

func (i *ImplFactory) SupportsBootstrapExecutor() bool {
	return true
}

func (i *ImplFactory) GenerateTransmitterKey() (string, error) {
	_, privateKey, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return "", fmt.Errorf("generate ed25519 transmitter key: %w", err)
	}

	return hex.EncodeToString(privateKey), nil
}

func (i *ImplFactory) TransmitterAddress(privateKeyHex string) (protocol.UnknownAddress, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return protocol.UnknownAddress{}, fmt.Errorf("invalid Canton private key hex: %w", err)
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return protocol.UnknownAddress{}, fmt.Errorf("invalid Canton private key length: %d", len(privateKeyBytes))
	}

	publicKey := ed25519.PrivateKey(privateKeyBytes).Public().(ed25519.PublicKey)

	return protocol.UnknownAddress(publicKey), nil
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

	// Send setup prerequisites
	routerAddress       contracts.InstanceAddress
	senderAddress       contracts.InstanceAddress
	registryAdmin       string
	validatorAPIClients validatorAPIClients

	feeTokenInstrument      splice_api_token_holding_v1.InstrumentId
	nextFeeCID              string // holding CID to be used as fee on next message send
	transferTokenInstrument *splice_api_token_holding_v1.InstrumentId
	nextTransferCID         string // holding CID to be used as transfer on next message send

	// verifierObs is injected post-construction by test runners (see SetVerifierObservation).
	// Required by ConfirmExecOnDest to fetch verifier results from aggregator/indexer.
	verifierObs VerifierObservation

	// partyMutexes serializes manual execution per receiver party so that
	// concurrent ConfirmExecOnDest calls cannot race on PerPartyRouter CID consumption.
	partyMutexes   map[string]*sync.Mutex
	partyMutexesMu sync.Mutex
}

type validatorAPIClients struct {
	scanClient     scanProxy.ClientWithResponsesInterface
	metadataClient tokenMetadataV1.ClientWithResponsesInterface
	transferClient transferInstructionV1.ClientWithResponsesInterface
}

func New(ctx context.Context, logger zerolog.Logger, e *deployment.Environment, chainSelector uint64) (*Chain, error) {
	chainFamily, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("get chain family for chain %d: %w", chainSelector, err)
	}
	if chainFamily != chainsel.FamilyCanton {
		return nil, fmt.Errorf("chain %d is not a canton chain", chainSelector)
	}
	chainDetails, err := chainsel.GetChainDetails(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("get chain details for chain %d: %w", chainSelector, err)
	}
	chain := e.BlockChains.CantonChains()[chainDetails.ChainSelector]

	return &Chain{
		e:            e,
		chain:        chain,
		chainDetails: chainDetails,
		logger:       logger,
		partyMutexes: make(map[string]*sync.Mutex),
	}, nil
}

func NewEmptyCCIP17Canton(logger zerolog.Logger) *Chain {
	return &Chain{
		logger:       logger,
		partyMutexes: make(map[string]*sync.Mutex),
	}
}

func (c *Chain) ChainSelector() uint64 {
	return c.chainDetails.ChainSelector
}

// ChainFamily implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) ChainFamily() string {
	return chainsel.FamilyCanton
}

// ConfigureNodes implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) ConfigureNodes(ctx context.Context, blockchain *blockchain.Input) (string, error) {
	return "", nil // TODO: implement
}

func uploadAndVetDar(ctx context.Context, participant canton.Participant, pkg contracts.Package) error {
	dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
	if err != nil {
		return fmt.Errorf("failed to get %s dar file: %w", pkg, err)
	}
	_, err = participant.LedgerServices.Admin.PackageManagement.UploadDarFile(ctx, &adminv2.UploadDarFileRequest{
		DarFile:       dar,
		VettingChange: adminv2.UploadDarFileRequest_VETTING_CHANGE_VET_ALL_PACKAGES,
	})
	if err != nil {
		return fmt.Errorf("failed to upload %s dar file: %w", pkg, err)
	}

	return nil
}

func (c *Chain) PreDeployContractsForSelector(ctx context.Context, env *deployment.Environment, selector uint64, _ *ccipOffchain.EnvironmentTopology) (datastore.DataStore, error) {
	chain, ok := env.BlockChains.CantonChains()[selector]
	if !ok || len(chain.Participants) == 0 {
		return nil, fmt.Errorf("canton chain %d not found or has no participants", selector)
	}

	participant := chain.Participants[0]
	for _, pkg := range cantonDeployDarPackages {
		if err := uploadAndVetDar(ctx, participant, pkg); err != nil {
			return nil, err
		}
	}

	runningDS := datastore.NewMemoryDataStore()
	owner := participant.PartyID

	for _, qual := range []string{dsutils.QualifierCCIP, dsutils.QualifierCCV, dsutils.QualifierRMN} {
		out, err := cantonchangesets.DeployCCIPFactory{}.Apply(*env, cantonchangesets.CantonCSDeps[cantonchangesets.DeployCCIPFactoryConfig]{
			ChainSelector: selector,
			Participant:   0,
			Config: cantonchangesets.DeployCCIPFactoryConfig{
				Params: cantonchangesets.DeployCCIPFactoryParams{
					OwnerParty: owner,
					MCMSParty:  owner,
					Qualifier:  qual,
					InstanceID: "ccip-factory-" + qual,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("deploy CCIPFactory %q for chain %d: %w", qual, selector, err)
		}
		if err := runningDS.Merge(out.DataStore.Seal()); err != nil {
			return nil, fmt.Errorf("merge factory %q datastore: %w", qual, err)
		}
	}

	return runningDS.Seal(), nil
}

func (c *Chain) GetDeployChainContractsCfg(env *deployment.Environment, selector uint64, _ *ccipOffchain.EnvironmentTopology) (ccipChangesets.DeployChainContractsPerChainCfg, error) {
	chain, ok := env.BlockChains.CantonChains()[selector]
	if !ok || len(chain.Participants) == 0 {
		return ccipChangesets.DeployChainContractsPerChainCfg{}, fmt.Errorf("canton chain %d not found or has no participants", selector)
	}

	return ccipChangesets.DeployChainContractsPerChainCfg{
		DeployerContract: fmt.Sprintf("canton:%s", chain.Participants[0].PartyID),
		DeployerKeyOwned: true,
	}, nil
}

func (c *Chain) PostDeployContractsForSelector(ctx context.Context, env *deployment.Environment, selector uint64, _ *ccipOffchain.EnvironmentTopology) (datastore.DataStore, error) {
	chain, ok := env.BlockChains.CantonChains()[selector]
	if !ok || len(chain.Participants) == 0 {
		return nil, fmt.Errorf("canton chain %d not found or has no participants", selector)
	}

	feeQuoterRef, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	))
	if err != nil {
		return nil, fmt.Errorf("resolve FeeQuoter for chain %d: %w", selector, err)
	}

	participant := chain.Participants[0]
	registryAdmin, err := testhelpers.ResolveRegistryAdmin(ctx, participant)
	if err != nil {
		return nil, fmt.Errorf("resolve registry admin: %w", err)
	}

	feeQuoterAddress := contracts.HexToInstanceAddress(feeQuoterRef.Address)
	_, err = operations.ExecuteOperation(env.OperationsBundle, feequoterop.ApplyPriceUpdatersUpdate, chain, contract.ChoiceInput[feequoter.ApplyPriceUpdatersUpdate]{
		InstanceAddress: feeQuoterAddress,
		Args: feequoter.ApplyPriceUpdatersUpdate{
			AddedPriceUpdaters: []types.PARTY{types.PARTY(participant.PartyID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("apply fee quoter price updaters update: %w", err)
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, feequoterop.UpdatePrices, chain, contract.ChoiceInput[feequoter.UpdatePrices]{
		InstanceAddress: feeQuoterAddress,
		Args: feequoter.UpdatePrices{
			PriceUpdates: feequoter.PriceUpdates{
				TokenPriceUpdates: []feequoter.TokenPriceUpdate{{
					InstrumentId: splice_api_token_holding_v1.InstrumentId{
						Admin: types.PARTY(registryAdmin),
						Id:    types.TEXT("Amulet"),
					},
					UsdPerToken: types.NUMERIC("100000000"),
				}},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("update fee quoter prices: %w", err)
	}

	return datastore.NewMemoryDataStore().Seal(), nil
}

func (c *Chain) GetConnectionProfile(env *deployment.Environment, selector uint64) (lanes.ChainDefinition, lanes.CommitteeVerifierRemoteChainInput, error) {
	// TODO this is currently not populated by populateAddressesV2
	globalConfig, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	if err != nil {
		return lanes.ChainDefinition{}, lanes.CommitteeVerifierRemoteChainInput{}, fmt.Errorf("failed to get GlobalConfig address for chain %d: %w", selector, err)
	}
	feeQuoter, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType(feequoterop.ContractType), feequoterop.Version, ""))
	if err != nil {
		return lanes.ChainDefinition{}, lanes.CommitteeVerifierRemoteChainInput{}, fmt.Errorf("failed to get FeeQuoter address for chain %d: %w", selector, err)
	}
	c.logger.Debug().Str("GlobalConfig", globalConfig.Address).Msg("Resolved GlobalConfig")

	registryAdmin, err := testhelpers.ResolveRegistryAdmin(context.Background(), c.chain.Participants[0])
	if err != nil {
		return lanes.ChainDefinition{}, lanes.CommitteeVerifierRemoteChainInput{}, fmt.Errorf("resolve registry admin for token prices: %w", err)
	}

	chainDefinition := lanes.ChainDefinition{
		Selector:           selector,
		AddressBytesLength: 32,
		TokenPrices: map[string]*big.Int{
			fmt.Sprintf("%s:%s", registryAdmin, "Amulet"): big.NewInt(100_000_000),
		},
		DefaultExecutor: datastore.AddressRef{
			ChainSelector: selector,
			Type:          datastore.ContractType(executor2.ContractType),
			Version:       executor2.Version,
			Qualifier:     devenvcommon.DefaultExecutorQualifier,
		},
		ExecutorDestChainConfig: lanes.ExecutorDestChainConfig{
			Enabled: false,
		},
		DefaultInboundCCVs: []datastore.AddressRef{
			{
				ChainSelector: selector,
				Type:          datastore.ContractType(committee_verifier.ContractType),
				Version:       committee_verifier.Version,
				Qualifier:     devenvcommon.DefaultCommitteeVerifierQualifier,
			},
		},
		DefaultOutboundCCVs: []datastore.AddressRef{
			{
				ChainSelector: selector,
				Type:          datastore.ContractType(committee_verifier.ContractType),
				Version:       committee_verifier.Version,
				Qualifier:     devenvcommon.DefaultCommitteeVerifierQualifier,
			},
		},
		BaseExecutionGasCost: 1,
		FeeQuoter:            contracts.HexToInstanceAddress(feeQuoter.Address).Bytes(),
		CantonLaneConfig: &lanes.CantonLaneConfig{
			GlobalConfig: globalConfig,
		},
	}
	cvConfig := lanes.CommitteeVerifierRemoteChainInput{
		GasForVerification: 50_000,
	}

	return chainDefinition, cvConfig, nil
}

func (c *Chain) GetChainLaneProfile(env *deployment.Environment, selector uint64) (cciptestinterfaces.ChainLaneProfile, error) {
	defaultFeeQuoterCfg := cantonadapters.DefaultCantonFeeQuoterDestChainConfig()

	baseExecutionGasCost := uint32(1)
	tokenReceiverAllowed := false
	gasForVerification := uint32(50_000)

	return cciptestinterfaces.ChainLaneProfile{
		BaseExecutionGasCost: &baseExecutionGasCost,
		TokenReceiverAllowed: &tokenReceiverAllowed,
		AllowedFinalityConfig: &finality.Config{
			WaitForFinality: true,
		},
		FeeQuoterDestChainConfig: ccipChangesets.FeeQuoterDestChainConfigOverrides{
			OverrideExistingConfig:      defaultFeeQuoterCfg.OverrideExistingConfig,
			IsEnabled:                   &defaultFeeQuoterCfg.IsEnabled,
			MaxDataBytes:                &defaultFeeQuoterCfg.MaxDataBytes,
			MaxPerMsgGasLimit:           &defaultFeeQuoterCfg.MaxPerMsgGasLimit,
			DestGasPerPayloadByteBase:   &defaultFeeQuoterCfg.DestGasPerPayloadByteBase,
			DefaultTokenFeeUSDCents:     &defaultFeeQuoterCfg.DefaultTokenFeeUSDCents,
			DefaultTokenDestGasOverhead: &defaultFeeQuoterCfg.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:           &defaultFeeQuoterCfg.DefaultTxGasLimit,
			NetworkFeeUSDCents:          &defaultFeeQuoterCfg.NetworkFeeUSDCents,
			LinkFeeMultiplierPercent:    &defaultFeeQuoterCfg.V2Params.LinkFeeMultiplierPercent,
			USDPerUnitGas:               defaultFeeQuoterCfg.V2Params.USDPerUnitGas,
		},
		ExecutorDestChainConfig: &ccipadapters.ExecutorDestChainConfig{
			Enabled: true,
		},
		DefaultExecutorQualifier: devenvcommon.DefaultExecutorQualifier,
		DefaultInboundCCVs: []datastore.AddressRef{
			{
				ChainSelector: selector,
				Type:          datastore.ContractType(committee_verifier.ContractType),
				Version:       committee_verifier.Version,
				Qualifier:     devenvcommon.DefaultCommitteeVerifierQualifier,
			},
		},
		DefaultOutboundCCVs: []datastore.AddressRef{
			{
				ChainSelector: selector,
				Type:          datastore.ContractType(committee_verifier.ContractType),
				Version:       committee_verifier.Version,
				Qualifier:     devenvcommon.DefaultCommitteeVerifierQualifier,
			},
		},
		GasForVerification: &gasForVerification,
	}, nil
}

func (c *Chain) PostConnect(env *deployment.Environment, selector uint64, remoteSelectors []uint64) error {
	return nil
}

func (c *Chain) GetSupportedPools() []devenvcommon.PoolCapability {
	return []devenvcommon.PoolCapability{
		{
			PoolType:    devenvcommon.LockReleaseTokenPoolType,
			PoolVersion: cantonTokenPoolVersion,
		},
	}
}

func (c *Chain) GetTokenExpansionConfigs(
	env *deployment.Environment,
	selector uint64,
	combos []devenvcommon.TokenCombination,
) ([]tokenscore.TokenExpansionInputPerChain, error) {
	chain, ok := env.BlockChains.CantonChains()[selector]
	if !ok || len(chain.Participants) == 0 {
		return nil, fmt.Errorf("canton chain %d not found or has no participants", selector)
	}
	registryAdmin, err := testhelpers.ResolveRegistryAdmin(context.Background(), chain.Participants[0])
	if err != nil {
		return nil, fmt.Errorf("resolve registry admin for token expansion: %w", err)
	}

	seen := make(map[poolConfigKey]struct{})
	var configs []tokenscore.TokenExpansionInputPerChain

	for _, combo := range combos {
		for _, poolRef := range []datastore.AddressRef{combo.LocalPoolAddressRef(), combo.RemotePoolAddressRef()} {
			if poolRef.Type != datastore.ContractType(devenvcommon.LockReleaseTokenPoolType) {
				continue
			}
			key := poolConfigKey{
				poolType:    poolRef.Type,
				poolVersion: poolRef.Version.String(),
				qualifier:   poolRef.Qualifier,
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			tokenAddress := hex.EncodeToString(gethcrypto.Keccak256([]byte("Amulet@" + registryAdmin)))
			configs = append(configs, tokenscore.TokenExpansionInputPerChain{
				TokenPoolVersion: poolRef.Version,
				DeployTokenPoolInput: &tokenscore.DeployTokenPoolInput{
					TokenRef: &datastore.AddressRef{
						Address:   tokenAddress,
						Type:      datastore.ContractType("Token"),
						Qualifier: poolRef.Qualifier,
						Labels: datastore.NewLabelSet(
							"instrument-admin:"+registryAdmin,
							"instrument-id:Amulet",
						),
					},
					PoolType:           string(poolRef.Type),
					TokenPoolQualifier: poolRef.Qualifier,
				},
			})
		}
	}

	return configs, nil
}

func (c *Chain) PostTokenDeploy(
	env *deployment.Environment,
	selector uint64,
	deployedRefs []datastore.AddressRef,
) error {
	hasLockReleasePool := false
	for _, ref := range deployedRefs {
		if ref.Type == datastore.ContractType(devenvcommon.LockReleaseTokenPoolType) {
			hasLockReleasePool = true
			break
		}
	}
	if !hasLockReleasePool {
		return nil
	}

	chain := env.BlockChains.CantonChains()[selector]
	if len(chain.Participants) == 0 {
		return fmt.Errorf("canton chain %d has no participants", selector)
	}
	ctx := context.Background()
	if env.GetContext != nil {
		ctx = env.GetContext()
	}
	if _, _, _, err := mintTwoAmuletHoldings(ctx, chain.Participants[0], chain.Participants[0].PartyID, "1000000.00"); err != nil {
		return fmt.Errorf("seed AMT liquidity: %w", err)
	}

	return nil
}

func mintTwoAmuletHoldings(
	ctx context.Context,
	participant canton.Participant,
	ownerParty string,
	amount string,
) (types.CONTRACT_ID, types.CONTRACT_ID, []*apiv2.DisclosedContract, error) {
	scanClient, metadataClient, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return "", "", nil, err
	}
	feeHoldingCID, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount)
	if err != nil {
		return "", "", nil, fmt.Errorf("mint fee holding: %w", err)
	}
	tokenHoldingCID, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount)
	if err != nil {
		return "", "", nil, fmt.Errorf("mint token-transfer holding: %w", err)
	}
	if feeHoldingCID == tokenHoldingCID {
		return "", "", nil, fmt.Errorf("mint returned same holding cid twice: %s", feeHoldingCID)
	}
	disclosedFeeHolding, err := testhelpers.GetDisclosedContractById(ctx, participant, feeHoldingCID)
	if err != nil {
		return "", "", nil, fmt.Errorf("get disclosed fee holding by id: %w", err)
	}
	disclosedTokenHolding, err := testhelpers.GetDisclosedContractById(ctx, participant, tokenHoldingCID)
	if err != nil {
		return "", "", nil, fmt.Errorf("get disclosed token holding by id: %w", err)
	}

	return types.CONTRACT_ID(feeHoldingCID), types.CONTRACT_ID(tokenHoldingCID), []*apiv2.DisclosedContract{
		disclosedFeeHolding,
		disclosedTokenHolding,
	}, nil
}

func (c *Chain) GetTokenTransferConfigs(
	env *deployment.Environment,
	selector uint64,
	remoteSelectors []uint64,
	topology *ccipOffchain.EnvironmentTopology,
) ([]tokenscore.TokenTransferConfig, error) {
	applicableCombos := devenvcommon.FilterTokenCombinations(
		devenvcommon.AllTokenCombinations(), topology, env.DataStore, append([]uint64{selector}, remoteSelectors...),
	)
	hasAddressRef := func(chainSelector uint64, ref datastore.AddressRef) bool {
		_, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(
			chainSelector,
			ref.Type,
			ref.Version,
			ref.Qualifier,
		))

		return err == nil
	}
	merged := make(map[poolConfigKey]tokenscore.TokenTransferConfig)

	for _, combo := range applicableCombos {
		for _, pair := range []struct {
			local, remote datastore.AddressRef
			ccvQuals      []string
		}{
			{combo.LocalPoolAddressRef(), combo.RemotePoolAddressRef(), combo.LocalPoolCCVQualifiers()},
			{combo.RemotePoolAddressRef(), combo.LocalPoolAddressRef(), combo.RemotePoolCCVQualifiers()},
		} {
			if pair.local.Type != datastore.ContractType(devenvcommon.LockReleaseTokenPoolType) {
				continue
			}
			if !hasAddressRef(selector, pair.local) {
				continue
			}

			remoteChains := make(map[uint64]tokenscore.RemoteChainConfig[*datastore.AddressRef, datastore.AddressRef])
			for _, rs := range remoteSelectors {
				family, err := chainsel.GetSelectorFamily(rs)
				if err != nil || family != chainsel.FamilyEVM {
					continue
				}
				if !hasAddressRef(rs, pair.remote) {
					continue
				}

				ccvRefs := make([]datastore.AddressRef, 0, len(pair.ccvQuals))
				for _, qualifier := range pair.ccvQuals {
					ccvRefs = append(ccvRefs, datastore.AddressRef{
						Type:      datastore.ContractType(committee_verifier.ContractType),
						Version:   committee_verifier.Version,
						Qualifier: qualifier,
					})
				}

				remoteChains[rs] = tokenscore.RemoteChainConfig[*datastore.AddressRef, datastore.AddressRef]{
					RemotePool:                &pair.remote,
					InboundRateLimiterConfig:  nil,
					OutboundRateLimiterConfig: nil,
					OutboundCCVs:              ccvRefs,
					InboundCCVs:               ccvRefs,
					// TODO: what to set for these?
					// RemoteToken: nil,
					// RemoteDecimals: 0,
					// OutboundCCVsToAddAboveThreshold: nil,
					// InboundCCVsToAddAboveThreshold: nil,
					// TokenTransferFeeConfig: tokenscore.TokenTransferFeeConfig{},
				}
			}
			if len(remoteChains) == 0 {
				continue
			}

			cfg := tokenscore.TokenTransferConfig{
				ChainSelector: selector,
				TokenPoolRef:  pair.local,
				TokenRef: datastore.AddressRef{
					Type:      datastore.ContractType("Token"),
					Qualifier: pair.local.Qualifier,
				},
				RegistryRef: datastore.AddressRef{
					Type:    datastore.ContractType(token_admin_registry.ContractType),
					Version: token_admin_registry.Version,
				},
				RemoteChains:          remoteChains,
				AllowedFinalityConfig: finality.Config{WaitForFinality: true},
			}

			key := poolConfigKey{
				poolType:    cfg.TokenPoolRef.Type,
				poolVersion: cfg.TokenPoolRef.Version.String(),
				qualifier:   cfg.TokenPoolRef.Qualifier,
			}

			if existing, ok := merged[key]; ok {
				maps.Copy(existing.RemoteChains, cfg.RemoteChains)
				merged[key] = existing
			} else {
				merged[key] = cfg
			}
		}
	}

	configs := make([]tokenscore.TokenTransferConfig, 0, len(merged))
	for _, cfg := range merged {
		configs = append(configs, cfg)
	}

	return configs, nil
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

// ExposeMetrics implements cciptestinterfaces.CCIP17.
func (c *Chain) ExposeMetrics(ctx context.Context, source, dest uint64) ([]string, *prometheus.Registry, error) {
	return nil, nil, nil // TODO: implement
}

// GetEOAReceiverAddress implements cciptestinterfaces.CCIP17.
func (c *Chain) GetEOAReceiverAddress() (protocol.UnknownAddress, error) {
	if len(c.chain.Participants) == 0 {
		return nil, fmt.Errorf("no canton participants configured")
	}

	receiver := contracts.HashedPartyFromString(c.chain.Participants[0].PartyID)

	return protocol.UnknownAddress(receiver.Bytes()), nil
}

// GetExpectedNextSequenceNumber implements cciptestinterfaces.CCIP17.
func (c *Chain) GetExpectedNextSequenceNumber(ctx context.Context, to uint64) (uint64, error) {
	return c.nextSeq + 1, nil
}

// GetMaxDataBytes implements cciptestinterfaces.CCIP17.
func (c *Chain) GetMaxDataBytes(ctx context.Context, remoteChainSelector uint64) (uint32, error) {
	return 0, nil // TODO: implement
}

// GetSenderAddress implements cciptestinterfaces.CCIP17.
func (c *Chain) GetSenderAddress() (protocol.UnknownAddress, error) {
	return protocol.UnknownAddress{}, nil // TODO: implement
}

// GetTokenBalance implements cciptestinterfaces.CCIP17.
func (c *Chain) GetTokenBalance(ctx context.Context, address, tokenAddress protocol.UnknownAddress) (*big.Int, error) {
	if len(c.chain.Participants) == 0 {
		return nil, fmt.Errorf("no canton participants configured")
	}

	participant := c.chain.Participants[0]
	ownerParty := participant.PartyID
	if len(address) > 0 {
		for _, p := range c.chain.Participants {
			if bytes.Equal(contracts.HashedPartyFromString(p.PartyID).Bytes(), []byte(address)) {
				participant = p
				ownerParty = p.PartyID

				break
			}
		}
	}

	// TODO: this currently gets all holdings. Differentiate by tokenAddress
	// need to map tokenAddress -> splice_api_token_holding_v1.InstrumentId
	totalRat, err := testhelpers.GetHoldingsBalance(ctx, participant, nil,
		testhelpers.WithHoldingOwner(ownerParty),
		testhelpers.WithUnlockedHoldingsOnly(),
	)
	if err != nil {
		return nil, fmt.Errorf("get holdings balance: %w", err)
	}

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	scaleRat := new(big.Rat).SetInt(scale)
	scaled := new(big.Rat).Mul(totalRat, scaleRat)
	if !scaled.IsInt() {
		return nil, fmt.Errorf("holding balance scale exceeds 10 decimals for total %s", totalRat.FloatString(12))
	}

	return scaled.Num(), nil
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

// ConfirmSendOnSource implements cciptestinterfaces.CCIP17.
func (c *Chain) ConfirmSendOnSource(ctx context.Context, to uint64, key cciptestinterfaces.MessageEventKey, timeout time.Duration) (cciptestinterfaces.MessageSentEvent, error) {
	if key.MessageID == (protocol.Bytes32{}) && key.SeqNum == 0 {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("MessageEventKey must have MessageID or SeqNum set")
	}
	if key.SeqNum != 0 {
		return c.waitOneSentEventBySeqNo(ctx, to, key.SeqNum, timeout)
	}
	deadline := time.Now().Add(timeout)
	for {
		if c.lastSentEvent.MessageID == key.MessageID {
			return c.lastSentEvent, nil
		}
		if timeout > 0 && time.Now().After(deadline) {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("timed out waiting for sent event by message ID")
		}
		select {
		case <-ctx.Done():
			return cciptestinterfaces.MessageSentEvent{}, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// ConfirmExecOnDest implements cciptestinterfaces.CCIP17.
//
// Canton has no executor service today, so finalization on Canton dest happens
// via ManuallyExecuteMessage. To keep the gun and e2e tests direction-agnostic,
// this method drives the full path:
//  1. Lock per receiver party (PerPartyRouter CID is consumed on every Execute).
//  2. Idempotency: if an ExecutionStateChanged for (from, seqNo, messageID) already
//     exists on the ledger, parse and return it without re-executing.
//  3. Otherwise fetch the verifier result via verifier observation (aggregator +
//     indexer), translate the verifier dest address to its hashed instance address,
//     and call ManuallyExecuteMessage.
//
// Both SeqNum AND MessageID must be set on key: they key the idempotency lookup and
// the verifier-result fetch respectively. EVM-side ConfirmExecOnDest is permissive
// and accepts either; Canton requires both. Test runners must wire verifier
// observation (see WireVerifierObservationFromLib) before calling this method.
func (c *Chain) ConfirmExecOnDest(ctx context.Context, from uint64, key cciptestinterfaces.MessageEventKey, timeout time.Duration) (cciptestinterfaces.ExecutionStateChangedEvent, error) {
	if key.MessageID == (protocol.Bytes32{}) || key.SeqNum == 0 {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("MessageEventKey must have both MessageID and SeqNum")
	}
	if !c.verifierObs.wired() {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("verifier observation not wired (test runner must call WireVerifierObservationFromLib)")
	}
	if len(c.chain.Participants) == 0 {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("no participants on chain")
	}

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	receiverParty := c.chain.Participants[0].PartyID
	unlock := c.lockForParty(receiverParty)
	defer unlock()

	if ev, found, err := c.findExistingExecutionState(ctx, from, key.SeqNum, key.MessageID); err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("idempotency lookup: %w", err)
	} else if found {
		c.logger.Info().
			Uint64("from", from).
			Uint64("seqNo", key.SeqNum).
			Str("messageID", hex.EncodeToString(key.MessageID[:])).
			Str("state", ev.State.String()).
			Msg("ConfirmExecOnDest idempotent return: ExecutionStateChanged already present")

		return ev, nil
	}

	vr, err := c.fetchVerifierResult(ctx, key.MessageID)
	if err != nil {
		return cciptestinterfaces.ExecutionStateChangedEvent{}, fmt.Errorf("fetch verifier result: %w", err)
	}

	return c.ManuallyExecuteMessage(
		ctx, vr.Message, 0,
		[]protocol.UnknownAddress{vr.HashedVerifierDestAddr},
		[][]byte{vr.CCVData},
	)
}

// SetupSend sets up Canton sender specific prerequisites for sending a message.
// eg: deploy per-party router, deploy ccipsender contract...
func (c *Chain) SetupSend(
	ctx context.Context,
	feeAmountPerMessage uint64,
	transferAmountPerMessage uint64,
) error {
	participant := c.chain.Participants[0]
	party := participant.PartyID

	// Deploy contracts //
	// Router for sender party.
	routerAddress, err := c.DeployPerPartyRouter(ctx, participant, party)
	if err != nil {
		return fmt.Errorf("failed to deploy per-party router: %w", err)
	}

	// Deploy a sender-owned CCIPSender contract.
	senderInstanceID := contracts.MustNewInstanceID("devenv-ccipsender")
	out, err := operations.ExecuteOperation(c.e.OperationsBundle, sender.Deploy, c.chain, contract.DeployInput[ccipsender.CCIPSender]{
		Qualifier: nil,
		Template: ccipsender.CCIPSender{
			InstanceId: types.TEXT(senderInstanceID),
			Owner:      types.PARTY(party),
		},
		OwnerParty: types.PARTY(party),
	})
	if err != nil {
		return fmt.Errorf("failed to deploy ccip sender contract: %w", err)
	}
	senderAddress := contracts.HexToInstanceAddress(out.Output.Address)
	registryAdmin, err := testhelpers.ResolveRegistryAdmin(ctx, participant)
	if err != nil {
		return fmt.Errorf("resolve registry admin: %w", err)
	}
	c.registryAdmin = registryAdmin
	c.routerAddress = routerAddress
	c.senderAddress = senderAddress

	// Clients setup //
	_, err = c.getValidatorAPIClients()
	if err != nil {
		return fmt.Errorf("get validator API clients: %w", err)
	}

	// Tokens setup //
	// Setup fee holding
	feeTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    AMTInstrument,
	}

	baseFilters := make([]testhelpers.Filter, 0, 3)
	baseFilters = append(baseFilters,
		testhelpers.WithHoldingOwner(party),
		testhelpers.WithUnlockedHoldingsOnly(),
	)
	feeRows, err := testhelpers.ListHoldingsForInstrument(ctx, participant, &feeTokenInstrument, baseFilters...)
	if err != nil {
		return fmt.Errorf("list fee holdings for setup: %w", err)
	}
	feeMin := new(big.Rat).SetUint64(feeAmountPerMessage)
	selectedFee, err := testhelpers.SelectHoldingsForInstrument(feeRows, []*big.Rat{feeMin})
	if err != nil || len(selectedFee) == 0 {
		return fmt.Errorf("select fee holding for setup: %w", err)
	}

	c.feeTokenInstrument = feeTokenInstrument
	c.nextFeeCID = selectedFee[0].ContractID

	// End setup, no transfer holdings needed
	if transferAmountPerMessage == 0 {
		return nil
	}

	// Setup transfer holdings
	// TODO: add support for LINK instrument
	transferTokenInstrument := &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    AMTInstrument,
	}

	// Select transfer holding using another call here because fee/token might use different instruments
	filters := append(baseFilters, testhelpers.ExcludeCIDs([]string{c.nextFeeCID})) // this filter is here in case fee and transfer use the same instrument
	transferRows, err := testhelpers.ListHoldingsForInstrument(ctx, participant, transferTokenInstrument, filters...)
	if err != nil {
		return fmt.Errorf("list transfer holdings for setup: %w", err)
	}
	transferMin := new(big.Rat).SetUint64(transferAmountPerMessage)
	selectedTransfer, err := testhelpers.SelectHoldingsForInstrument(transferRows, []*big.Rat{transferMin})
	if err != nil || len(selectedTransfer) == 0 {
		return fmt.Errorf("select transfer holding for setup: %w", err)
	}

	c.transferTokenInstrument = transferTokenInstrument
	c.nextTransferCID = selectedTransfer[0].ContractID

	return nil
}

// MintTokens mint tokens for transfer and fees. To be used on devenv tests only.
// this method won't work in staging/prod tests
// TODO: add support for LINK instrument
func (c *Chain) MintTokens(ctx context.Context, amount uint64) error {
	participant := c.chain.Participants[0]
	party := participant.PartyID

	validatorAPIClients, err := c.getValidatorAPIClients()
	if err != nil {
		return fmt.Errorf("get validator API clients: %w", err)
	}

	_, err = testhelpers.MintAMT(
		ctx,
		participant,
		validatorAPIClients.metadataClient,
		validatorAPIClients.transferClient,
		validatorAPIClients.scanClient,
		party,
		strconv.FormatUint(amount, 10),
	)
	if err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	return nil
}

// SendMessage implements cciptestinterfaces.CCIP17.
func (c *Chain) SendMessage(ctx context.Context, dest uint64, fields cciptestinterfaces.MessageFields, dataProvider cciptestinterfaces.ExtraArgsDataProvider, messageVersion uint8) (cciptestinterfaces.MessageSentEvent, error) {
	opts, ok := dataProvider.(cciptestinterfaces.MessageOptions)
	if !ok {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("canton SendMessage only supports cciptestinterfaces.MessageOptions, got %T", dataProvider)
	}
	if messageVersion != 3 {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("canton SendMessage only supports message version 3, got %d", messageVersion)
	}
	var unset contracts.InstanceAddress
	if c.routerAddress == unset || c.senderAddress == unset {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
			"canton SendMessage: router or sender address is unset; call SetupSend once before trying to send messages",
		)
	}

	hasTokenTransfer := c.transferTokenInstrument != nil && fields.TokenAmount.Amount != nil && fields.TokenAmount.Amount.Sign() > 0

	participant := c.chain.Participants[0]
	party := participant.PartyID

	if strings.TrimSpace(c.nextFeeCID) == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
			"canton SendMessage: next fee holding CID is unset; call SetupSend after minting",
		)
	}
	if hasTokenTransfer {
		if strings.TrimSpace(c.nextTransferCID) == "" {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
				"canton SendMessage: next transfer holding CID is unset; call SetupSend after minting",
			)
		}
		if c.nextFeeCID == c.nextTransferCID {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf(
				"canton SendMessage: fee and transfer holding CIDs must be different contracts",
			)
		}
	}

	feeFactoryChoice, err := testhelpers.GetTransferFactoryV2(
		ctx,
		c.validatorAPIClients.transferClient,
		c.registryAdmin,
		splice_api_token_transfer_instruction_v1.Transfer{
			Sender:           types.PARTY(party),
			Receiver:         types.PARTY(party),
			Amount:           types.NUMERIC("100.00"), // TransferFactory API doesn't actually do anything with this value, anything should work here
			InstrumentId:     c.feeTokenInstrument,
			InputHoldingCids: []types.CONTRACT_ID{types.CONTRACT_ID(c.nextFeeCID)},
			Meta:             splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}}, // The API doesn't actually use that value
		},
	)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve fee transfer factory from scan-proxy: %w", err)
	}
	feeTransferFactorycid := feeFactoryChoice.FactoryID
	feeTransferFactoryDisclosures := feeFactoryChoice.DisclosedContracts
	feeTransferFactoryChoiceContextRaw := feeFactoryChoice.ChoiceContextData
	feeTransferFactoryChoiceContextValue, err := testhelpers.ChoiceContextFromData(feeTransferFactoryChoiceContextRaw)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("parse fee transfer factory choice context: %w", err)
	}

	// Collect Disclosures
	outgoingMessage := oapiCommon.Message{
		DestinationChainSelector: strconv.FormatUint(dest, 10),
		FeeToken: oapiCommon.InstrumentId{
			Admin: oapiCommon.PartyId(c.feeTokenInstrument.Admin),
			Id:    string(c.feeTokenInstrument.Id),
		},
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{Type: oapiCommon.Empty}, // Using default Executor
		Payload:  hex.EncodeToString(fields.Data),
		Receiver: hex.EncodeToString(fields.Receiver),
	}
	if hasTokenTransfer {
		// TODO: move this to the other line also branching by tokenTransfer
		outgoingMessage.TokenTransfer = &oapiCommon.TokenTransfer{
			Amount: fields.TokenAmount.Amount.String(),
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(c.transferTokenInstrument.Admin),
				Id:    string(c.transferTokenInstrument.Id),
			},
		}
	}
	// TODO come up with a better way of doing this
	routerCid, err := contract.FindActiveContractIDByInstanceAddress(ctx, participant.LedgerServices.State, []string{party}, perpartyrouter.PerPartyRouter{}.GetTemplateID(), c.routerAddress)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("find active contract ID for router at address %s: %w", c.routerAddress, err)
	}
	var disclosedContracts []*apiv2.DisclosedContract
	sendArgs := ccipsender.Send{
		RouterCid:                types.CONTRACT_ID(routerCid),
		DestinationChainSelector: types.NUMERIC(strconv.FormatUint(dest, 10)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver: types.TEXT(hex.EncodeToString(fields.Receiver)),
			Payload:  types.TEXT(hex.EncodeToString(fields.Data)),
			FeeToken: c.feeTokenInstrument,
			ExtraArgs: ccipclient.ExtraArgs{
				V3: &ccipclient.GenericExtraArgsV3{
					GasLimit: types.INT64(opts.ExecutionGasLimit),
					Executor: ccipclient.ExecutorExtraArg{
						ExecutorUseDefault: &ccipclient.ExecutorUseDefault{},
					},
				},
			},
		},
		FeeTokenInput: ccipsender.FeeTokenInput{
			SenderInputCids:         []types.CONTRACT_ID{types.CONTRACT_ID(c.nextFeeCID)},
			FeeTokenTransferFactory: types.CONTRACT_ID(feeTransferFactorycid),
			FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
				Context: splice_api_token_metadata_v1.ChoiceContext{Values: testhelpers.ExtractChoiceContextValues(feeTransferFactoryChoiceContextValue)},
				Meta:    splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
			},
		},
	}

	// Sender-required CCVs
	var senderRequiredCCVs []string
	for _, ccvItem := range opts.CCVs {
		senderRequiredCCVs = append(senderRequiredCCVs, contracts.BytesToInstanceAddress(ccvItem.CCVAddress.Bytes()).String())
	}

	// Token Pool
	var tokenPoolRequiredCCVs []string
	if hasTokenTransfer {
		// Get Token Pool
		token := contracts.EncodeInstrumentID(c.feeTokenInstrument)
		tokenPoolAddress, err := c.GetTokenPoolForToken(ctx, token)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("get token pool for token %s: %w", token.String(), err)
		}

		// Query Token Pool API
		tokenPoolSendDisclosure, err := c.GetTokenPoolSendDisclosure(ctx, outgoingMessage, tokenPoolAddress.InstanceAddress())
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get Token Pool Send disclosure for token pool at address %s: %w", tokenPoolAddress.String(), err)
		}
		sendArgs.Message.TokenTransfer = &ccipclient.TokenTransfer{
			Token: splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(outgoingMessage.TokenTransfer.Token.Admin),
				Id:    types.TEXT(outgoingMessage.TokenTransfer.Token.Id),
			},
			Amount: types.NUMERIC(outgoingMessage.TokenTransfer.Amount),
		}
		sendArgs.TokenTransferInput = &ccipsender.TokenTransferInput{
			SenderInputCids:  []types.CONTRACT_ID{types.CONTRACT_ID(c.nextTransferCID)},
			TokenPoolCid:     types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
			PoolExtraContext: tokenPoolSendDisclosure.ChoiceContext,
		}
		tokenPoolRequiredCCVs = tokenPoolSendDisclosure.RequiredCCVs
		disclosedContracts = append(disclosedContracts, tokenPoolSendDisclosure.DisclosedContracts...)
	}

	// CCIP
	ccipSendDisclosure, err := c.GetCCIPSendDisclosure(ctx, outgoingMessage, senderRequiredCCVs, tokenPoolRequiredCCVs)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get CCIP Send disclosure: %w", err)
	}
	sendArgs.Context = ccipSendDisclosure.ChoiceContext
	sendArgs.FeeTokenInput.FeeTokenConfigCid = types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid)
	disclosedContracts = append(disclosedContracts, ccipSendDisclosure.DisclosedContracts...)

	// CCVs
	var allCCVs []string
	for _, v := range ccipSendDisclosure.CCVs {
		// The value returned by the API can be a RawInstanceAddress or InstanceAddress, depending on what was provided as input
		var ccvAddress contracts.InstanceAddress
		if rawInstanceAddress, err := contracts.RawInstanceAddressFromString(v); err == nil {
			ccvAddress = rawInstanceAddress.InstanceAddress()
		} else {
			ccvAddress = contracts.HexToInstanceAddress(v)
		}

		ccvSendDisclosure, err := c.GetCCVSendDisclosure(ctx, outgoingMessage, ccvAddress)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get CCV Send disclosure for CCV %q: %w", v, err)
		}
		sendArgs.CcvSendInputs = append(sendArgs.CcvSendInputs, ccipsender.CCVSendInput{
			CcvAddress:      ccvSendDisclosure.Address.Binding(),
			CcvCid:          types.CONTRACT_ID(ccvSendDisclosure.ContractId),
			CcvExtraContext: ccvSendDisclosure.ChoiceContext,
		})
		sendArgs.Message.ExtraArgs.V3.Ccvs = append(sendArgs.Message.ExtraArgs.V3.Ccvs, ccipclient.CCVExtraArg{
			CcvAddress: ccvSendDisclosure.Address.Binding(),
			CcvArgs:    "",
		})
		allCCVs = append(allCCVs, ccvSendDisclosure.Address.InstanceAddress().Hex())
		disclosedContracts = append(disclosedContracts, ccvSendDisclosure.DisclosedContracts...)
	}

	// Executor
	if ccipSendDisclosure.Executor != nil {
		var executorAddress contracts.InstanceAddress
		if rawInstanceAddress, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor); err == nil {
			executorAddress = rawInstanceAddress.InstanceAddress()
		} else {
			executorAddress = contracts.HexToInstanceAddress(*ccipSendDisclosure.Executor)
		}

		executorSendDisclosure, err := c.GetExecutorSendDisclosure(ctx, outgoingMessage, executorAddress, allCCVs)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get Executor Send disclosure for Executor %q: %w", *ccipSendDisclosure.Executor, err)
		}
		sendArgs.ExecutorInput = &ccipsender.ExecutorInput{
			ExecutorCid:          types.CONTRACT_ID(executorSendDisclosure.ContractId),
			ExecutorExtraContext: executorSendDisclosure.ChoiceContext,
		}
		disclosedContracts = append(disclosedContracts, executorSendDisclosure.DisclosedContracts...)
	}

	// Fee Token
	disclosedContracts = append(disclosedContracts, feeTransferFactoryDisclosures...)

	// TODO deduplicating disclosed contracts shouldn't be required once we move away from native from token transfers
	disclosedContracts = testhelpers.DeduplicateDisclosedContracts(disclosedContracts...)

	// Call CCIPSend
	ccipSendReport, err := operations.ExecuteOperation(c.e.OperationsBundle, sender.Send, c.chain, contract.ChoiceInput[ccipsender.Send]{
		InstanceAddress:    c.senderAddress,
		Args:               sendArgs,
		MCMSEnabled:        false,
		DisclosedContracts: contract.DisclosedContractsFromProto(disclosedContracts),
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("execute CCIP Send: %w", err)
	}
	c.logger.Info().Str("UpdateID", ccipSendReport.Output.ExecInfo.UpdateID).Msg("CCIP Send executed")
	update, err := participant.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: ccipSendReport.Output.ExecInfo.UpdateID,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						participant.PartyID: {
							Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
										TemplateFilter: &apiv2.TemplateFilter{
											TemplateId:              contracts.TemplateIDFromBinding(ccipcore.CCIPMessageSent{}).ToLedgerIdentifier(),
											IncludeCreatedEventBlob: true,
										},
									},
								},
								{
									IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{
										InterfaceFilter: &apiv2.InterfaceFilter{
											InterfaceId:             testhelpers.HoldingV1InterfaceID(),
											IncludeInterfaceView:    true,
											IncludeCreatedEventBlob: false,
										},
									},
								},
							},
						},
					},
					Verbose: true,
				},
			},
		},
	})
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get update %q: %w", ccipSendReport.Output.ExecInfo.UpdateID, err)
	}

	parsedSend, err := parseFirstCCIPMessageSentFromLedgerEvents(update.GetTransaction().GetEvents(), c.nextSeq)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, err
	}

	// Set next holdings
	err = c.setNextHoldings(update.GetTransaction().GetEvents(), hasTokenTransfer, fields.TokenAmount.Amount)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("set next holdings: %w", err)
	}

	event := cciptestinterfaces.MessageSentEvent{
		MessageID:      parsedSend.messageID,
		ReceiptIssuers: nil, // TODO: add them later, not currently needed
	}
	if parsedSend.foundEncodedMessage {
		event.Message = new(parsedSend.decodedMessage)
	}
	c.nextSeq = parsedSend.seqNo
	c.lastSentDest = dest
	c.lastSentSeq = parsedSend.seqNo
	c.lastSentEvent = event

	return event, nil
}

// ccipMessageSentFromSendUpdate is the outcome of parsing a send transaction's first CCIPMessageSent
// created event: wire identifiers, updated sequence, and optional decoded message body.
type ccipMessageSentFromSendUpdate struct {
	messageID           [32]byte
	seqNo               uint64
	foundEncodedMessage bool
	decodedMessage      protocol.Message
}

// parseFirstCCIPMessageSentFromLedgerEvents scans events for the first created CCIPMessageSent
// (module and entity name from generated bindings). previousSeq is the chain's last known sequence
// before this send; returned seqNo is either previousSeq+1 or the sequence from the encoded payload
// when that path is present.
func parseFirstCCIPMessageSentFromLedgerEvents(events []*apiv2.Event, previousSeq uint64) (ccipMessageSentFromSendUpdate, error) {
	messageSentTemplateID := contracts.TemplateIDFromBinding(ccipcore.CCIPMessageSent{})

	// Find CCIPMessageSent event in the events
	var created *apiv2.CreatedEvent
	for _, event := range events {
		evCreated := event.GetCreated()
		if evCreated == nil || evCreated.GetTemplateId() == nil {
			continue
		}
		tid := evCreated.GetTemplateId()
		if tid.GetModuleName() != messageSentTemplateID.ModuleName || tid.GetEntityName() != messageSentTemplateID.EntityName {
			continue
		}
		created = evCreated

		break
	}
	if created == nil {
		return ccipMessageSentFromSendUpdate{}, fmt.Errorf("no CCIPMessageSent event found in sender transaction")
	}

	parsed, err := bindings.UnmarshalCreatedEvent[ccipcore.CCIPMessageSent](created)
	if err != nil {
		return ccipMessageSentFromSendUpdate{}, fmt.Errorf("unmarshal CCIPMessageSent created event: %w", err)
	}
	evt := parsed.Event

	var (
		foundEventMessageID, foundEncodedMessage bool
		eventMessageID, computedMessageID        [32]byte
		decodedMessage                           protocol.Message
	)

	if mid := string(evt.MessageId); mid != "" {
		decoded, err := hex.DecodeString(mid)
		if err != nil || len(decoded) != 32 {
			return ccipMessageSentFromSendUpdate{}, fmt.Errorf("decode messageId from CCIPMessageSent event: %w", err)
		}
		copy(eventMessageID[:], decoded)
		foundEventMessageID = true
	}

	if enc := string(evt.EncodedMessage); enc != "" {
		encodedMessage, err := hex.DecodeString(enc)
		if err != nil {
			return ccipMessageSentFromSendUpdate{}, fmt.Errorf("decode encodedMessage from CCIPMessageSent event: %w", err)
		}
		decodedMessagePtr, err := protocol.DecodeMessage(encodedMessage)
		if err != nil {
			return ccipMessageSentFromSendUpdate{}, fmt.Errorf("decode protocol message from encodedMessage: %w", err)
		}
		decodedMessage = *decodedMessagePtr
		computedHash := gethcrypto.Keccak256(encodedMessage)
		if len(computedHash) != len(computedMessageID) {
			return ccipMessageSentFromSendUpdate{}, fmt.Errorf("computed encodedMessage hash has invalid length: %d", len(computedHash))
		}
		copy(computedMessageID[:], computedHash)
		foundEncodedMessage = true
	}

	seqNo := previousSeq + 1
	var messageID [32]byte
	if foundEncodedMessage {
		// Prefer the sequence from the encoded message since local process state can be stale
		// across repeated test runs against a long-lived environment.
		seqNo = uint64(decodedMessage.SequenceNumber)
		messageID = computedMessageID
	} else if foundEventMessageID {
		messageID = eventMessageID
	} else {
		return ccipMessageSentFromSendUpdate{}, fmt.Errorf("CCIPMessageSent event missing both messageId and encodedMessage")
	}

	return ccipMessageSentFromSendUpdate{
		messageID:           messageID,
		seqNo:               seqNo,
		foundEncodedMessage: foundEncodedMessage,
		decodedMessage:      decodedMessage,
	}, nil
}

// setNextHoldings sets fee and transfer holdings on Chain struct for next MessageSend
// based on the events from the send update.
func (c *Chain) setNextHoldings(events []*apiv2.Event, hasTokenTransfer bool, tokenAmount *big.Int) error {
	party := c.chain.Participants[0].PartyID

	spentCIDs := []string{c.nextFeeCID}
	if hasTokenTransfer && c.nextTransferCID != "" {
		spentCIDs = append(spentCIDs, c.nextTransferCID)
	}
	refreshFilters := []testhelpers.Filter{
		testhelpers.WithHoldingOwner(party),
		testhelpers.WithUnlockedHoldingsOnly(),
		testhelpers.ExcludeCIDs(spentCIDs),
	}

	// Select next fee holding
	freshFeeHoldings, err := testhelpers.ListedHoldingsFromTransactionEventsForInstrument(
		events,
		c.feeTokenInstrument,
		refreshFilters...,
	)
	if err != nil {
		return fmt.Errorf("parse fee holdings from send update: %w", err)
	}
	pickedFee, err := testhelpers.SelectHoldingsForInstrument(freshFeeHoldings, []*big.Rat{big.NewRat(0, 1)}) // TODO: we'll have to consider real fee here once we go load tests
	if err != nil {
		return fmt.Errorf("refresh next fee holding from update: %w", err)
	}
	c.nextFeeCID = pickedFee[0].ContractID
	c.logger.Info().Str("NextFeeCID", c.nextFeeCID).Msg("Selected next fee holding")
	if payload, err := json.MarshalIndent(freshFeeHoldings, "", "  "); err != nil {
		c.logger.Warn().Err(err).Msg("marshal freshFeeHoldings for log")
	} else {
		c.logger.Debug().
			RawJSON("freshFeeHoldings", payload).
			Msg("Fresh fee holdings parsed from send update")
	}

	if !hasTokenTransfer {
		return nil
	}

	// Select next transfer holding
	refreshFilters = append(refreshFilters[:len(refreshFilters)-1], testhelpers.ExcludeCIDs(append(spentCIDs, c.nextFeeCID))) // includes one more item in the filter to exclude the fee holding
	freshTransferHoldings, err := testhelpers.ListedHoldingsFromTransactionEventsForInstrument(
		events,
		*c.transferTokenInstrument,
		refreshFilters...,
	)
	if err != nil {
		return fmt.Errorf("parse transfer holdings from send update: %w", err)
	}
	transferValue := new(big.Rat).SetInt(tokenAmount)
	pickedTransfer, err := testhelpers.SelectHoldingsForInstrument(freshTransferHoldings, []*big.Rat{transferValue})
	if err != nil {
		return fmt.Errorf("refresh next transfer holding from update: %w", err)
	}
	c.nextTransferCID = pickedTransfer[0].ContractID
	c.logger.Info().Str("NextTransferCID", c.nextTransferCID).Msg("Selected next transfer holding")

	return nil
}

func (c *Chain) waitOneSentEventBySeqNo(ctx context.Context, to, seq uint64, timeout time.Duration) (cciptestinterfaces.MessageSentEvent, error) {
	deadline := time.Now().Add(timeout)

	for {
		if c.lastSentDest == to && c.lastSentSeq == seq {
			return c.lastSentEvent, nil
		}

		if timeout > 0 && time.Now().After(deadline) {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("timed out waiting for sent event: dest=%d seq=%d", to, seq)
		}

		select {
		case <-ctx.Done():
			return cciptestinterfaces.MessageSentEvent{}, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// getValidatorAPIClients gets validator API clients from Chain struct cache or
// creates them if they are not already created.
func (c *Chain) getValidatorAPIClients() (validatorAPIClients, error) {
	participant := c.chain.Participants[0]
	if c.validatorAPIClients.scanClient != nil {
		return c.validatorAPIClients, nil
	}

	// Clients setup
	scanClient, metadataClient, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return validatorAPIClients{}, fmt.Errorf("creating validator API clients: %w", err)
	}

	c.validatorAPIClients = validatorAPIClients{
		scanClient:     scanClient,
		metadataClient: metadataClient,
		transferClient: transferClient,
	}

	return c.validatorAPIClients, nil
}
