package devenv

import (
	"bytes"
	"context"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"maps"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	ccipChangesets "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	ccipOffchain "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	ccvservices "github.com/smartcontractkit/chainlink-ccv/build/devenv/services"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	mcmsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonadapters "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	cantonchangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	executor2 "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

var (
	_                       cciptestinterfaces.CCIP17              = &Chain{}
	_                       cciptestinterfaces.CCIP17Configuration = &Chain{}
	_                       ccv.ImplFactory                        = &ImplFactory{}
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

const (
	// #nosec G101 -- fixed test datastore qualifier, not a credential
	cantonDestTokenQualifier = "TEST (LockReleaseTokenPool 2.0.0 [default] to BurnMintTokenPool 2.0.0 [default])"
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

	lastSentMessage               protocol.Message
	lastSentVerifierDestAddress   protocol.UnknownAddress
	lastSentVerifierBlob          []byte
	lastSentHasVerificationInputs bool
	cfg                           *ccv.Cfg
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

	out, err := cantonchangesets.DeployCCIPFactory{}.Apply(*env, cantonchangesets.CantonCSDeps[cantonchangesets.DeployCCIPFactoryConfig]{
		ChainSelector: selector,
		Participant:   0,
		Config: cantonchangesets.DeployCCIPFactoryConfig{
			Params: cantonchangesets.DeployCCIPFactoryParams{
				OwnerParty: participant.PartyID,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("deploy CCIPFactory for chain %d: %w", selector, err)
	}

	return out.DataStore.Seal(), nil
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

// DeployContractsForSelector implements cciptestinterfaces.CCIP17Configuration.
func (c *Chain) DeployContractsForSelector(ctx context.Context, env *deployment.Environment, selector uint64, topology *ccipOffchain.EnvironmentTopology) (datastore.DataStore, error) {
	return ccv.DeployContractsForSelector(ctx, env, c, selector, topology)
}

func (c *Chain) GetConnectionProfile(env *deployment.Environment, selector uint64) (lanes.ChainDefinition, lanes.CommitteeVerifierRemoteChainInput, error) {
	// TODO this is currently not populated by populateAddressesV2
	globalConfig, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	if err != nil {
		return lanes.ChainDefinition{}, lanes.CommitteeVerifierRemoteChainInput{}, fmt.Errorf("failed to get GlobalConfig address for chain %d: %w", selector, err)
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
	if env != nil && env.DataStore != nil {
		cantonadapters.SetRuntimeDataStore(env.DataStore)
	} else if c.e != nil && c.e.DataStore != nil {
		cantonadapters.SetRuntimeDataStore(c.e.DataStore)
	}

	defaultFeeQuoterCfg := cantonadapters.DefaultCantonFeeQuoterDestChainConfig()

	return cciptestinterfaces.ChainLaneProfile{
		AddressBytesLength:   32,
		BaseExecutionGasCost: 1,
		FeeQuoterDestChainConfig: ccipadapters.FeeQuoterDestChainConfig{
			OverrideExistingConfig:      defaultFeeQuoterCfg.OverrideExistingConfig,
			IsEnabled:                   defaultFeeQuoterCfg.IsEnabled,
			MaxDataBytes:                defaultFeeQuoterCfg.MaxDataBytes,
			MaxPerMsgGasLimit:           defaultFeeQuoterCfg.MaxPerMsgGasLimit,
			DestGasOverhead:             defaultFeeQuoterCfg.DestGasOverhead,
			DestGasPerPayloadByteBase:   defaultFeeQuoterCfg.DestGasPerPayloadByteBase,
			ChainFamilySelector:         cantonadapters.CantonFamilySelector,
			DefaultTokenFeeUSDCents:     defaultFeeQuoterCfg.DefaultTokenFeeUSDCents,
			DefaultTokenDestGasOverhead: defaultFeeQuoterCfg.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:           defaultFeeQuoterCfg.DefaultTxGasLimit,
			NetworkFeeUSDCents:          defaultFeeQuoterCfg.NetworkFeeUSDCents,
			LinkFeeMultiplierPercent:    defaultFeeQuoterCfg.V2Params.LinkFeeMultiplierPercent,
			USDPerUnitGas:               defaultFeeQuoterCfg.V2Params.USDPerUnitGas,
		},
		ExecutorDestChainConfig: ccipadapters.ExecutorDestChainConfig{
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
		GasForVerification: 50_000,
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
) (types.CONTRACT_ID, types.CONTRACT_ID, []*ledgerv2.DisclosedContract, error) {
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

	return types.CONTRACT_ID(feeHoldingCID), types.CONTRACT_ID(tokenHoldingCID), []*ledgerv2.DisclosedContract{
		disclosedFeeHolding,
		disclosedTokenHolding,
	}, nil
}

func mintAmuletHolding(
	ctx context.Context,
	participant canton.Participant,
	ownerParty string,
	amount string,
) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
	scanClient, metadataClient, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return "", nil, err
	}
	holdingCID, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount)
	if err != nil {
		return "", nil, fmt.Errorf("mint fee holding: %w", err)
	}
	disclosedHolding, err := testhelpers.GetDisclosedContractById(ctx, participant, holdingCID)
	if err != nil {
		return "", nil, fmt.Errorf("get disclosed holding by id: %w", err)
	}

	return types.CONTRACT_ID(holdingCID), disclosedHolding, nil
}

func (c *Chain) GetTokenTransferConfigs(
	env *deployment.Environment,
	selector uint64,
	remoteSelectors []uint64,
	topology *ccipOffchain.EnvironmentTopology,
) ([]tokenscore.TokenTransferConfig, error) {
	applicableCombos := devenvcommon.FilterTokenCombinations(
		devenvcommon.AllTokenCombinations(), topology, nil, nil,
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
					RemotePool:                               &pair.remote,
					DefaultFinalityInboundRateLimiterConfig:  tokenscore.RateLimiterConfigFloatInput{},
					DefaultFinalityOutboundRateLimiterConfig: tokenscore.RateLimiterConfigFloatInput{},
					CustomFinalityInboundRateLimiterConfig:   tokenscore.RateLimiterConfigFloatInput{},
					CustomFinalityOutboundRateLimiterConfig:  tokenscore.RateLimiterConfigFloatInput{},
					OutboundCCVs:                             ccvRefs,
					InboundCCVs:                              ccvRefs,
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

// Curse implements cciptestinterfaces.CCIP17.
func (c *Chain) Curse(ctx context.Context, subjects [][16]byte) error {
	rmnRemoteRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(c.chainDetails.ChainSelector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	if err != nil {
		return fmt.Errorf("get rmn remote address: %w", err)
	}

	instanceAddr := contracts.HexToInstanceAddress(rmnRemoteRef.Address)

	c.logger.Info().
		Uint64("chainSelector", c.chainDetails.ChainSelector).
		Int("numSubjects", len(subjects)).
		Msg("Cursing subjects on chain")
	for _, subject := range subjects {
		_, err := operations.ExecuteOperation(c.e.OperationsBundle, rmn_remote.Curse, c.chain, contract.ChoiceInput[rmn.Curse]{
			InstanceAddress: instanceAddr,
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

	instanceAddr := contracts.HexToInstanceAddress(rmnRemoteRef.Address)

	c.logger.Info().
		Uint64("chainSelector", c.chainDetails.ChainSelector).
		Int("numSubjects", len(subjects)).
		Msg("Uncursing subjects on chain")
	for _, subject := range subjects {
		_, err := operations.ExecuteOperation(c.e.OperationsBundle, rmn_remote.Uncurse, c.chain, contract.ChoiceInput[rmn.Uncurse]{
			InstanceAddress: instanceAddr,
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

	holdingContracts, err := testhelpers.ListActiveContractsByInterfaceId(ctx, participant, &ledgerv2.Identifier{
		PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
	if err != nil {
		return nil, fmt.Errorf("list active token holdings: %w", err)
	}

	total := big.NewInt(0)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	scaleRat := new(big.Rat).SetInt(scale)
	for _, holding := range holdingContracts {
		views := holding.GetCreatedEvent().GetInterfaceViews()
		if len(views) == 0 || views[0].GetViewValue() == nil {
			continue
		}
		fields := views[0].GetViewValue().GetFields()
		if len(fields) < 4 {
			continue
		}
		ownerField := fields[0].GetValue().GetParty()
		amountRaw := strings.TrimSpace(fields[2].GetValue().GetNumeric())
		locked := fields[3].GetValue().GetOptional().GetValue() != nil

		if ownerField != ownerParty || amountRaw == "" || locked {
			continue
		}

		amountRat, ok := new(big.Rat).SetString(amountRaw)
		if !ok {
			return nil, fmt.Errorf("invalid holding amount %q", fields[2].GetValue().GetNumeric())
		}
		amountRat.Mul(amountRat, scaleRat)
		if !amountRat.IsInt() {
			return nil, fmt.Errorf("holding amount scale exceeds 10: %q", fields[2].GetValue().GetNumeric())
		}
		total.Add(total, amountRat.Num())
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

// SendMessage implements cciptestinterfaces.CCIP17.
func (c *Chain) SendMessage(ctx context.Context, dest uint64, fields cciptestinterfaces.MessageFields, opts cciptestinterfaces.MessageOptions) (cciptestinterfaces.MessageSentEvent, error) {
	participant := c.chain.Participants[0]
	party := participant.PartyID

	// Router for sender party.
	_, disclosedRouter, err := c.DeployPerPartyRouter(ctx, participant, party)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("deploy per-party router: %w", err)
	}

	seqNo := c.nextSeq + 1

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
	activeContract, err := contract.FindActiveContractByInstanceAddress(
		ctx, participant.LedgerServices.State, party, ccipsender.CCIPSender{}.GetTemplateID(), senderInstanceAddress)
	disclosedCCIPSender := convertToDisclosedContract(activeContract)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve ccip sender disclosed contract: %w", err)
	}
	if ccipSenderCID == "" {
		ccipSenderCID = types.CONTRACT_ID(disclosedCCIPSender.GetContractId())
	}
	if ccipSenderCID == "" {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolved empty ccip sender contract ID")
	}

	// collect ccv instance addresses so we can request their disclosures from EDS.
	var ccvInstanceAddresses []contracts.InstanceAddress
	for _, ccvItem := range opts.CCVs {
		ccvInstanceAddress := contracts.BytesToInstanceAddress(ccvItem.CCVAddress.Bytes())
		ccvInstanceAddresses = append(ccvInstanceAddresses, ccvInstanceAddress)
	}

	// get send disclosures and build the ccipsender.Send struct.
	sendDisclosures, err := c.GetDisclosuresForSend(ctx, ccvInstanceAddresses)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get disclosures for send: %w", err)
	}
	if len(sendDisclosures.CCVContractIDs) != len(opts.CCVs) {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("expected %d ccv disclosures returned from EDS, got %d", len(opts.CCVs), len(sendDisclosures.CCVContractIDs))
	}

	ccvExtraArgs := make([]ccipclient.CCVExtraArg, 0, len(opts.CCVs))
	ccvSendInputs := make([]ccipsender.CCVSendInput, 0, len(opts.CCVs))
	var fallbackVerifierDestAddress protocol.UnknownAddress
	var fallbackVerifierBlob []byte
	for i := range opts.CCVs {
		// Get the raw instance address of the CCV from the datastore
		addrRef, err := getByAddress(c.e.DataStore, ccvInstanceAddresses[i].String())
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("failed to get address ref for ccv %d from env datastore: %w", i, err)
		}
		if len(addrRef.Labels.List()) == 0 {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("no labels found for ccv %d", i)
		}
		rawInstanceAddress := addrRef.Labels.List()[0]
		rawAddrBinding := mcmsbindings.RawInstanceAddress{Unpack: types.TEXT(rawInstanceAddress)}
		ccvExtraArgs = append(ccvExtraArgs, ccipclient.CCVExtraArg{
			CcvAddress: rawAddrBinding,
			CcvArgs:    types.TEXT(""),
		})
		ccvSendInputs = append(ccvSendInputs, ccipsender.CCVSendInput{
			CcvAddress:      rawAddrBinding,
			CcvCid:          types.CONTRACT_ID(sendDisclosures.CCVContractIDs[i].ContractId),
			CcvExtraContext: common.CCIPContext{},
		})
	}

	hasTokenTransfer := fields.TokenAmount.Amount != nil && fields.TokenAmount.Amount.Sign() > 0
	var feeSenderInputCIDs []types.CONTRACT_ID
	var preferredTokenHoldingCID types.CONTRACT_ID
	var mintedHoldingDisclosures []*ledgerv2.DisclosedContract
	if hasTokenTransfer {
		feeHoldingCID, tokenHoldingCID, holdingDisclosures, mintErr := mintTwoAmuletHoldings(ctx, participant, party, "1000000.00")
		if mintErr != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("mint two amulet holdings for send: %w", mintErr)
		}
		feeSenderInputCIDs = []types.CONTRACT_ID{feeHoldingCID}
		preferredTokenHoldingCID = tokenHoldingCID
		mintedHoldingDisclosures = holdingDisclosures
	} else {
		feeHoldingCID, disclosedHolding, mintErr := mintAmuletHolding(ctx, participant, party, "1000000.00")
		if mintErr != nil {
			return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("mint fee holding for send: %w", mintErr)
		}
		feeSenderInputCIDs = []types.CONTRACT_ID{feeHoldingCID}
		mintedHoldingDisclosures = []*ledgerv2.DisclosedContract{disclosedHolding}
	}
	registryAdmin, err := testhelpers.ResolveRegistryAdmin(ctx, participant)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve registry admin: %w", err)
	}
	feeTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}

	_, _, transferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("create transfer instruction client: %w", err)
	}
	feeTransferFactorycid, feeTransferFactoryDisclosures, feeTransferFactoryChoiceContextValue, err := testhelpers.GetTransferFactory(
		ctx,
		transferClient,
		registryAdmin,
		party,
		party,
	)
	if err != nil {
		return cciptestinterfaces.MessageSentEvent{}, fmt.Errorf("resolve fee transfer factory from scan-proxy: %w", err)
	}

	feeTokenInput := interfaces.TokenInput{
		TransferFactory: types.CONTRACT_ID(feeTransferFactorycid),
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: splice_api_token_metadata_v1.ChoiceContext{Values: testhelpers.ExtractChoiceContextValues(feeTransferFactoryChoiceContextValue)},
			Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
		},
		TokenPoolHoldings: []types.CONTRACT_ID{},
	}

	var messageTokenTransfer *ccipclient.TokenTransfer
	var tokenTransferInput *ccipsender.TokenTransferInput
	var tokenTransferDisclosures []*ledgerv2.DisclosedContract
	ccipReceiveGasLimit := opts.ExecutionGasLimit
	if hasTokenTransfer {
		messageTokenTransfer, tokenTransferInput, tokenTransferDisclosures, err = c.buildTokenTransferSendInput(
			ctx,
			participant,
			transferClient,
			registryAdmin,
			party,
			dest,
			preferredTokenHoldingCID,
		)
		if err != nil {
			return cciptestinterfaces.MessageSentEvent{}, err
		}
		// Keep send fee path deterministic for token-transfer devenv tests.
		ccipReceiveGasLimit = 0
	}

	sendArgs := ccipsender.Send{
		Context:                  sendDisclosures.SendContext,
		RouterCid:                types.CONTRACT_ID(disclosedRouter.ContractId),
		DestinationChainSelector: types.NUMERIC(fmt.Sprintf("%d", dest)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver:      types.TEXT(hex.EncodeToString(fields.Receiver)),
			Payload:       types.TEXT(hex.EncodeToString(fields.Data)),
			TokenTransfer: messageTokenTransfer,
			FeeToken:      feeTokenInstrument,
			ExtraArgs: ccipclient.ExtraArgs{
				V3: &ccipclient.GenericExtraArgsV3{
					GasLimit: types.INT64(ccipReceiveGasLimit),
					Ccvs:     ccvExtraArgs,
					Executor: ccipclient.ExecutorExtraArg{
						ExecutorUseDefault: &ccipclient.ExecutorUseDefault{
							ExecutorArgs: types.TEXT(""),
						},
					},
					TokenReceiver: types.TEXT(""),
					TokenArgs:     types.TEXT(""),
				},
			},
		},
		FeeTokenInput: ccipsender.FeeTokenInput{
			SenderInputCids: feeSenderInputCIDs,
			TokenInput:      feeTokenInput,
		},
		CcvSendInputs:      ccvSendInputs,
		TokenTransferInput: tokenTransferInput,
		ExecutorInput: &ccipsender.ExecutorInput{
			ExecutorCid:          types.CONTRACT_ID(sendDisclosures.ExecutorContractID.ContractId),
			ExecutorExtraContext: common.CCIPContext{},
		},
	}
	sendArgsMap := sendArgs.ToMap()

	// sender, router + all other disclosed contracts
	disclosedContracts := make([]*ledgerv2.DisclosedContract, 0, 2+len(sendDisclosures.DisclosedContracts))
	disclosedContracts = append(disclosedContracts, disclosedCCIPSender, disclosedRouter)
	disclosedContracts = append(disclosedContracts, sendDisclosures.DisclosedContracts...)
	disclosedContracts = append(disclosedContracts, mintedHoldingDisclosures...)
	disclosedContracts = append(disclosedContracts, feeTransferFactoryDisclosures...)
	disclosedContracts = append(disclosedContracts, tokenTransferDisclosures...)
	disclosedContracts = testhelpers.DeduplicateDisclosedContracts(disclosedContracts...)
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
		ReceiptIssuers: nil, // TODO: add them later, not currently needed
	}
	if foundEncodedMessage {
		event.Message = new(decodedMessage)
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

func (c *Chain) buildTokenTransferSendInput(
	ctx context.Context,
	participant canton.Participant,
	transferClient transferInstructionV1.ClientWithResponsesInterface,
	registryAdmin string,
	party string,
	dest uint64,
	senderInputCID types.CONTRACT_ID,
) (*ccipclient.TokenTransfer, *ccipsender.TokenTransferInput, []*ledgerv2.DisclosedContract, error) {
	const tokenTransferAmountDecimal = "0.0000001000"

	poolRefs := c.e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(c.chainDetails.ChainSelector),
		datastore.AddressRefByType(datastore.ContractType("LockReleaseTokenPool")),
		datastore.AddressRefByQualifier(cantonDestTokenQualifier),
	)
	if len(poolRefs) != 1 {
		return nil, nil, nil, fmt.Errorf("expected exactly one LockReleaseTokenPool, got %d", len(poolRefs))
	}
	poolRef := poolRefs[0]
	// TODO: replace with EDS ccipSend token disclosure whenever available
	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		contracts.HexToInstanceAddress(poolRef.Address),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve lock/release pool %s: %w", poolRef.Address, err)
	}
	parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse lock/release pool %s: %w", poolRef.Address, err)
	}
	remoteCfgAny, ok := parsedPool.RemoteChainConfigs[fmt.Sprintf("%d.", dest)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("missing remote chain config for %d", dest)
	}
	var outboundRateLimiterInstanceAddress contracts.InstanceAddress
	switch remoteCfg := remoteCfgAny.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		rawOutboundRateLimiter, err := contracts.RawInstanceAddressFromString(string(remoteCfg.OutboundRateLimiter.Unpack))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("parse outbound rate limiter raw address for %d: %w", dest, err)
		}
		outboundRateLimiterInstanceAddress = rawOutboundRateLimiter.InstanceAddress()
	case map[string]any:
		rawOutboundRateLimiter, ok := remoteCfg["outboundRateLimiter"].(map[string]any)
		if !ok {
			return nil, nil, nil, fmt.Errorf("missing outbound rate limiter in remote chain config for %d", dest)
		}
		rawOutboundRateLimiterText, ok := rawOutboundRateLimiter["unpack"].(string)
		if !ok || rawOutboundRateLimiterText == "" {
			return nil, nil, nil, fmt.Errorf("missing outbound rate limiter address in remote chain config for %d", dest)
		}
		parsedRawOutboundRateLimiter, err := contracts.RawInstanceAddressFromString(rawOutboundRateLimiterText)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("parse outbound rate limiter raw address for %d: %w", dest, err)
		}
		outboundRateLimiterInstanceAddress = parsedRawOutboundRateLimiter.InstanceAddress()
	default:
		return nil, nil, nil, fmt.Errorf("unexpected remote chain config type %T for %d", remoteCfgAny, dest)
	}
	// TODO: maybe replace outboundRateLimiterInstanceAddress with value from EDS?? EDS TP rate limit is still 0x0
	activeOutboundRateLimiter, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		common.RateLimiter{}.GetTemplateID(),
		outboundRateLimiterInstanceAddress,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve outbound rate limiter contract: %w", err)
	}
	transferFactoryCIDRaw, transferFactoryDisclosures, transferFactoryChoiceContextValue, err := testhelpers.GetTransferFactory(
		ctx,
		transferClient,
		registryAdmin,
		party,
		party,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("resolve token transfer factory from scan-proxy: %w", err)
	}
	transferFactoryCID := types.CONTRACT_ID(transferFactoryCIDRaw)
	transferFactoryContext := splice_api_token_metadata_v1.ChoiceContext{
		Values: testhelpers.ExtractChoiceContextValues(transferFactoryChoiceContextValue),
	}
	tokenInput := interfaces.TokenInput{
		TransferFactory: transferFactoryCID,
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: transferFactoryContext,
			Meta:    splice_api_token_metadata_v1.Metadata{Values: types.TEXTMAP{}},
		},
		TokenPoolHoldings: []types.CONTRACT_ID{},
	}

	return &ccipclient.TokenTransfer{
			Token:  parsedPool.InstrumentId,
			Amount: types.NUMERIC(tokenTransferAmountDecimal),
		}, &ccipsender.TokenTransferInput{
			SenderInputCids: []types.CONTRACT_ID{senderInputCID},
			TokenPoolCid:    types.CONTRACT_ID(activePool.GetCreatedEvent().GetContractId()),
			PoolExtraContext: common.CCIPContext{
				Values: types.TEXTMAP{
					"rate-limiter": common.AnyValue{AVContractId: new(types.CONTRACT_ID(activeOutboundRateLimiter.GetCreatedEvent().GetContractId()))},
				},
			},
			TokenInput: tokenInput,
		}, append(
			[]*ledgerv2.DisclosedContract{
				convertToDisclosedContract(activePool),
				convertToDisclosedContract(activeOutboundRateLimiter),
			},
			transferFactoryDisclosures...,
		), nil
}

func getByAddress(ds datastore.DataStore, address string) (datastore.AddressRef, error) {
	all, err := ds.Addresses().Fetch()
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	for _, addr := range all {
		if addr.Address == address {
			return addr, nil
		}
	}

	return datastore.AddressRef{}, fmt.Errorf("address %s not found", address)
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
