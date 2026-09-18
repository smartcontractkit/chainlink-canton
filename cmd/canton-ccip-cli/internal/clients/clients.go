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
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/authentication"
	"github.com/smartcontractkit/chainlink-canton/authentication/providers/authorizationcode"
	"github.com/smartcontractkit/chainlink-canton/authentication/providers/clientcredentials"
	"github.com/smartcontractkit/chainlink-canton/authentication/providers/static"
	cfgpkg "github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/config"
	"github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/evmledger"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/eds/api/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/eds/api/ccv"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/eds/api/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/eds/api/tokenpool"
	oapiTokenMetadata "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// Bundle holds every constructed client/handle used by the commands.
type Bundle struct {
	Profile *cfgpkg.NetworkProfile
	Config  *cfgpkg.UserConfig

	// Canton
	Participant canton.Participant
	// Mapping of owner party -> URL for CCIP EDS APIs
	edsURLs map[types.PARTY]string
	// Mapping of owner party -> URL for CIP-56 Token Standard APIs
	tokenStandardURLs map[types.PARTY]string

	// EVM
	ETHClient  *ethclient.Client
	ETHAddress common.Address
	EthAuth    *bind.TransactOpts
	EthChainID *big.Int

	IndexerClient *indexerclient.IndexerClient

	// Explorers
	CCIPExplorerURL   string
	EVMExplorerURL    string
	CantonExplorerURL string
}

type EDSClients struct {
	CCIPEDS      oapiCCIP.ClientWithResponsesInterface
	CCVEDS       oapiCCV.ClientWithResponsesInterface
	ExecutorEDS  oapiExecutor.ClientWithResponsesInterface
	TokenPoolEDS oapiTokenPool.ClientWithResponsesInterface
}

type TokenStandardClients struct {
	TransferInstructionClient oapiTransferInstruction.ClientWithResponsesInterface
	MetadataClient            oapiTokenMetadata.ClientWithResponsesInterface
}

// New builds a Bundle for the given profile + user config.
func New(ctx context.Context, profile *cfgpkg.NetworkProfile, cfg *cfgpkg.UserConfig) (*Bundle, error) {
	bundle := &Bundle{
		Profile: profile,
		Config:  cfg,
	}
	var err error

	if !cfg.Canton.Disabled {
		// --- Canton ---
		var authProvider authentication.Provider
		switch cfg.Canton.AuthType {
		case "static":
			authProvider = static.NewInsecureStaticProvider(cfg.Canton.AuthJWT)
		case "clientCredentials":
			authProvider, err = clientcredentials.NewDiscoveryProvider(ctx, cfg.Canton.AuthServerURL, cfg.Canton.AuthClientID, cfg.Canton.AuthClientSecret)
			if err != nil {
				return nil, fmt.Errorf("create clientCredentials provider: %w", err)
			}
		case "authorizationCode":
			fallthrough
		default:
			authProvider, err = authorizationcode.NewDiscoveryProvider(ctx, cfg.Canton.AuthServerURL, cfg.Canton.AuthClientID)
			if err != nil {
				return nil, fmt.Errorf("create authorizationCode provider: %w", err)
			}
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
		bundle.Participant = cantonChain.Participants[0]

		// --- EDS URLs ---
		bundle.edsURLs = make(map[types.PARTY]string)
		for party, url := range profile.EDSURLs {
			bundle.edsURLs[types.PARTY(party)] = url
		}
		for party, url := range cfg.Canton.EDSURLs {
			bundle.edsURLs[types.PARTY(party)] = url
		}

		// --- Token Standard URLs ---
		bundle.tokenStandardURLs = make(map[types.PARTY]string)
		for party, url := range profile.TokenStandardURLs {
			bundle.tokenStandardURLs[types.PARTY(party)] = url
		}
		for party, url := range cfg.Canton.TokenStandardURLs {
			bundle.tokenStandardURLs[types.PARTY(party)] = url
		}
	}

	bundle.IndexerClient, err = indexerclient.NewIndexerClient(profile.IndexerURL, &http.Client{Timeout: 15 * time.Second})
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
	bundle.ETHClient, err = ethclient.DialContext(ctx, cfg.EVM.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial EVM rpc: %w", err)
	}
	bundle.EthChainID = new(big.Int).SetUint64(chainID)
	// The private key is optional: when it is not configured, EVM transactions
	// must be signed with a Ledger device via the --ledger flag on the commands.
	if cfg.EVM.PrivateKeyHex != "" {
		pk, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.EVM.PrivateKeyHex, "0x"))
		if err != nil {
			return nil, fmt.Errorf("parse EVM private key: %w", err)
		}
		publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("invalid EVM private key")
		}
		bundle.ETHAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
		bundle.EthAuth, err = bind.NewKeyedTransactorWithChainID(pk, bundle.EthChainID)
		if err != nil {
			return nil, fmt.Errorf("create EVM transactor: %w", err)
		}
	}

	// --- Explorers ---
	bundle.CCIPExplorerURL = profile.CCIPExlorerURL
	if cfg.CCIPExplorerURL != "" {
		bundle.CCIPExplorerURL = cfg.CCIPExplorerURL
	}
	bundle.EVMExplorerURL = profile.EVMExplorerURL
	if cfg.EVMExplorerURL != "" {
		bundle.EVMExplorerURL = cfg.EVMExplorerURL
	}
	bundle.CantonExplorerURL = profile.CantonExplorerURL
	if cfg.CantonExplorerURL != "" {
		bundle.CantonExplorerURL = cfg.CantonExplorerURL
	}

	return bundle, nil
}

// UseLedgerEVM connects to a Ledger device, derives the account at the given
// derivation path and replaces the EVM signer with one that signs transactions
// using the Ledger. The returned function closes the device connection.
func (b *Bundle) UseLedgerEVM(ctx context.Context, pathOrIndex string) (func(), error) {
	auth, address, closeLedger, err := evmledger.NewTransactor(ctx, pathOrIndex, b.EthChainID)
	if err != nil {
		return nil, err
	}
	b.EthAuth = auth
	b.ETHAddress = address

	return closeLedger, nil
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

func (b *Bundle) GetEDSClients(party types.PARTY) (EDSClients, error) {
	url, ok := b.edsURLs[party]
	if !ok {
		return EDSClients{}, fmt.Errorf("no EDS URL found for party %q", party)
	}

	var (
		clients EDSClients
		err     error
	)
	clients.CCIPEDS, err = oapiCCIP.NewClientWithResponses(url)
	if err != nil {
		return EDSClients{}, fmt.Errorf("create CCIP EDS client: %w", err)
	}
	clients.CCVEDS, err = oapiCCV.NewClientWithResponses(url)
	if err != nil {
		return EDSClients{}, fmt.Errorf("create CCV EDS client: %w", err)
	}
	clients.ExecutorEDS, err = oapiExecutor.NewClientWithResponses(url)
	if err != nil {
		return EDSClients{}, fmt.Errorf("create executor EDS client: %w", err)
	}
	clients.TokenPoolEDS, err = oapiTokenPool.NewClientWithResponses(url)
	if err != nil {
		return EDSClients{}, fmt.Errorf("create token pool EDS client: %w", err)
	}

	return clients, nil
}

func (b *Bundle) GetTokenStandardClients(party types.PARTY) (TokenStandardClients, error) {
	var (
		clients TokenStandardClients
		err     error
	)

	url, ok := b.tokenStandardURLs[party]
	if ok {
		clients.TransferInstructionClient, err = oapiTransferInstruction.NewClientWithResponses(url)
		if err != nil {
			return TokenStandardClients{}, fmt.Errorf("create transferInstruction EDS client: %w", err)
		}
		clients.MetadataClient, err = oapiTokenMetadata.NewClientWithResponses(url)
		if err != nil {
			return TokenStandardClients{}, fmt.Errorf("create metadata EDS client: %w", err)
		}

		return clients, nil
	}

	// If no override is found and party is DSO, use Validator API
	if party == b.Profile.DSOPartyID {
		_, clients.MetadataClient, clients.TransferInstructionClient, err = testhelpers.NewValidatorAPIClients(b.Participant)
		if err != nil {
			return TokenStandardClients{}, fmt.Errorf("create validator API clients: %w", err)
		}

		return clients, nil
	}

	return TokenStandardClients{}, fmt.Errorf("no Token Standard URL found for party %q", party)
}
