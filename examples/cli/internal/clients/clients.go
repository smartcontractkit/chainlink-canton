// Package clients constructs the various RPC / HTTP clients required by the
// CLI commands.
package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	chainsel "github.com/smartcontractkit/chain-selectors"
	indexerclient "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/client"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"

	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"

	cfgpkg "github.com/smartcontractkit/chainlink-canton/examples/cli/internal/config"
)

// Bundle holds every constructed client/handle used by the commands.
type Bundle struct {
	Profile *cfgpkg.NetworkProfile
	Config  *cfgpkg.UserConfig

	Participant canton.Participant

	// EDS clients
	CCIPEDS              oapiCCIP.ClientWithResponsesInterface
	CCVEDS               oapiCCV.ClientWithResponsesInterface
	ExecutorEDS          oapiExecutor.ClientWithResponsesInterface
	TokenPoolEDS         oapiTokenPool.ClientWithResponsesInterface
	AmuletTransferClient oapiTransferInstruction.ClientWithResponsesInterface
	LinkTransferClient   oapiTransferInstruction.ClientWithResponsesInterface
	IndexerClient        *indexerclient.IndexerClient

	// EVM
	ETHClient  *ethclient.Client
	ETHAddress common.Address
	EthAuth    *bind.TransactOpts
	EthChainID *big.Int

	// Explorers
	CCIPExplorerURL   string
	EVMExplorerURL    string
	CantonExplorerURL string
}

// New builds a Bundle for the given profile + user config.
func New(ctx context.Context, profile *cfgpkg.NetworkProfile, cfg *cfgpkg.UserConfig) (*Bundle, error) {
	// --- Canton ---
	authProvider, err := authorizationcode.NewDiscoveryProvider(ctx, cfg.Canton.AuthServerURL, cfg.Canton.AuthClientID)
	if err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}
	if _, err := authProvider.TokenSource().Token(); err != nil {
		return nil, fmt.Errorf("retrieve initial token: %w", err)
	}

	rpcCfg := provider.RPCChainProviderConfig{
		Participants: []provider.ParticipantConfig{{
			Endpoints: provider.Endpoints{
				JSONLedgerAPIURL: "json",
				GRPCLedgerAPIURL: cfg.Canton.ParticipantGRPCLedgerAPIURL,
				ValidatorAPIURL:  cfg.Canton.ValidatorAPIURL,
			},
			UserID:       cfg.Canton.UserID,
			PartyID:      cfg.Canton.PartyID,
			AuthProvider: authProvider,
		}},
	}
	ch, err := provider.NewRPCChainProvider(profile.CantonSelector, rpcCfg).Initialize(ctx)
	if err != nil {
		return nil, fmt.Errorf("init canton chain: %w", err)
	}
	cantonChain, ok := ch.(*canton.Chain)
	if !ok {
		return nil, fmt.Errorf("unexpected chain type %T", ch)
	}
	if len(cantonChain.Participants) == 0 {
		return nil, fmt.Errorf("no participants configured")
	}
	participant := cantonChain.Participants[0]

	// --- Validator API clients ---
	_, _, amuletTransferClient, err := testhelpers.NewValidatorAPIClients(participant)
	if err != nil {
		return nil, fmt.Errorf("create validator API clients: %w", err)
	}

	// --- EDS clients ---
	ccipEds, err := oapiCCIP.NewClientWithResponses(profile.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("create CCIP EDS client: %w", err)
	}
	ccvEds, err := oapiCCV.NewClientWithResponses(profile.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("create CCV EDS client: %w", err)
	}
	execEds, err := oapiExecutor.NewClientWithResponses(profile.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("create executor EDS client: %w", err)
	}
	tokenPoolEds, err := oapiTokenPool.NewClientWithResponses(profile.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("create token pool EDS client: %w", err)
	}
	transferInstructionEds, err := oapiTransferInstruction.NewClientWithResponses(profile.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("create transferInstruction EDS client: %w", err)
	}

	idx, err := indexerclient.NewIndexerClient(profile.IndexerURL, &http.Client{Timeout: 15 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("create indexer client: %w", err)
	}

	// --- EVM ---
	chainIDStr, err := chainsel.GetChainIDFromSelector(profile.EthSelector)
	if err != nil {
		return nil, fmt.Errorf("resolve EVM chain id: %w", err)
	}
	chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse EVM chain id: %w", err)
	}
	ethClient, err := ethclient.DialContext(ctx, cfg.EVM.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial EVM rpc: %w", err)
	}
	pk, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.EVM.PrivateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse EVM private key: %w", err)
	}
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid EVM private key")
	}
	ethAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	chainIDBig := new(big.Int).SetUint64(chainID)
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainIDBig)
	if err != nil {
		return nil, fmt.Errorf("create EVM transactor: %w", err)
	}

	// --- Explorers ---
	ccipExplorerURL := profile.CCIPExlorerURL
	if cfg.CCIPExplorerURL != "" {
		ccipExplorerURL = cfg.CCIPExplorerURL
	}
	evmExplorerURL := profile.EVMExplorerURL
	if cfg.EVMExplorerURL != "" {
		evmExplorerURL = cfg.EVMExplorerURL
	}
	cantonExplorerURL := profile.CantonExplorerURL
	if cfg.CantonExplorerURL != "" {
		cantonExplorerURL = cfg.CantonExplorerURL
	}

	return &Bundle{
		Profile:              profile,
		Config:               cfg,
		Participant:          participant,
		CCIPEDS:              ccipEds,
		CCVEDS:               ccvEds,
		ExecutorEDS:          execEds,
		TokenPoolEDS:         tokenPoolEds,
		AmuletTransferClient: amuletTransferClient,
		LinkTransferClient:   transferInstructionEds,
		IndexerClient:        idx,
		ETHClient:            ethClient,
		ETHAddress:           ethAddress,
		EthAuth:              auth,
		EthChainID:           chainIDBig,
		CCIPExplorerURL:      ccipExplorerURL,
		EVMExplorerURL:       evmExplorerURL,
		CantonExplorerURL:    cantonExplorerURL,
	}, nil
}

func (b *Bundle) CCIPExplorerLink(msgId string) string {
	return fmt.Sprintf("%s/msg/0x%s", strings.TrimSuffix(b.CCIPExplorerURL, "/"), strings.TrimPrefix(msgId, "0x"))
}

func (b *Bundle) EVMExplorerLink(tx string) string {
	return fmt.Sprintf("%s/tx/%s", strings.TrimSuffix(b.EVMExplorerURL, "/"), strings.TrimPrefix(tx, "0x"))
}

func (b *Bundle) CantonExplorerLink(update string) string {
	return fmt.Sprintf("%s/transactions/%s", strings.TrimSuffix(b.CantonExplorerURL, "/"), strings.TrimPrefix(update, "0x"))
}
