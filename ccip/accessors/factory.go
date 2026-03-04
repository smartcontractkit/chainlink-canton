package accessors

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"
)

type factory struct {
	lggr          logger.Logger
	helper        map[string]ccip.BlockchainInfo
	readerConfigs map[string]sourcereader.ReaderConfig
}

// GetAccessor implements chainaccess.AccessorFactory.
func (f *factory) GetAccessor(ctx context.Context, chainSelector protocol.ChainSelector) (chainaccess.Accessor, error) {
	if f.helper == nil {
		return nil, fmt.Errorf("canton ccip config is not set - can't get accessor for chain %d", chainSelector)
	}

	if f.readerConfigs == nil {
		return nil, fmt.Errorf("canton reader configs are not set - can't get accessor for chain %d", chainSelector)
	}

	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get selector family for %d - update chain-selectors library?: %w", chainSelector, err)
	}
	if family != chainsel.FamilyCanton {
		return nil, fmt.Errorf("skipping chain, only canton is supported for chain %d, family %s", chainSelector, family)
	}

	strSelector := strconv.FormatUint(uint64(chainSelector), 10)
	blockchainInfo, ok := f.helper[strSelector]
	if !ok {
		return nil, fmt.Errorf("canton config not found for chain %d", chainSelector)
	}

	readerConfig, ok := f.readerConfigs[strSelector]
	if !ok {
		return nil, fmt.Errorf("canton reader config not found for chain %d", chainSelector)
	}

	tokenSource, err := newTokenSource(ctx, blockchainInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token source for chain %d: %w", chainSelector, err)
	}

	sourceReader, err := sourcereader.NewSourceReader(
		logger.Named(f.lggr, fmt.Sprintf("CantonSourceReader.%d", chainSelector)),
		blockchainInfo.GRPCLedgerAPIURL,
		tokenSource,
		sourcereader.ReaderConfig{
			NodeOperatorParty:         readerConfig.NodeOperatorParty,
			CCIPOwnerParty:            readerConfig.CCIPOwnerParty,
			CCIPMessageSentTemplateID: readerConfig.CCIPMessageSentTemplateID,
			Authority:                 readerConfig.Authority,
		},
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: make this configurable
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create source reader: %w", err)
	}

	return newAccessor(sourceReader), nil
}

// newTokenSource builds an oauth2.TokenSource from the BlockchainInfo auth configuration.
// When Auth.Type is empty or "static", it falls back to the top-level JWT field for backward compatibility.
func newTokenSource(ctx context.Context, info ccip.BlockchainInfo) (oauth2.TokenSource, error) {
	authType := info.Auth.Type
	if authType == "" {
		authType = ccip.AuthTypeStatic
	}

	switch authType {
	case ccip.AuthTypeStatic:
		jwt := info.Auth.JWT
		if jwt == "" {
			jwt = info.JWT
		}
		if jwt == "" {
			return nil, fmt.Errorf("static auth requires a JWT token (set auth.jwt or top-level jwt)")
		}
		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: jwt}), nil

	case ccip.AuthTypeClientCredentials:
		if info.Auth.AuthURL == "" || info.Auth.ClientID == "" || info.Auth.ClientSecret == "" {
			return nil, fmt.Errorf("clientCredentials auth requires auth_url, client_id, and client_secret")
		}
		provider, err := clientcredentials.NewDiscoveryProvider(ctx, info.Auth.AuthURL, info.Auth.ClientID, info.Auth.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("clientCredentials provider: %w", err)
		}
		return provider.TokenSource(), nil

	case ccip.AuthTypeAuthorizationCode:
		if info.Auth.AuthURL == "" || info.Auth.ClientID == "" {
			return nil, fmt.Errorf("authorizationCode auth requires auth_url and client_id")
		}
		provider, err := authorizationcode.NewDiscoveryProvider(ctx, info.Auth.AuthURL, info.Auth.ClientID)
		if err != nil {
			return nil, fmt.Errorf("authorizationCode provider: %w", err)
		}
		return provider.TokenSource(), nil

	default:
		return nil, fmt.Errorf("unsupported auth type: %q (expected static, clientCredentials, or authorizationCode)", authType)
	}
}

func NewFactory(lggr logger.Logger, helper map[string]ccip.BlockchainInfo, readerConfigs map[string]sourcereader.ReaderConfig) chainaccess.AccessorFactory {
	return &factory{
		lggr:          lggr,
		helper:        helper,
		readerConfigs: readerConfigs,
	}
}

type accessor struct {
	sourceReader chainaccess.SourceReader
}

func newAccessor(sourceReader chainaccess.SourceReader) chainaccess.Accessor {
	return &accessor{
		sourceReader: sourceReader,
	}
}

func (a *accessor) SourceReader() chainaccess.SourceReader {
	return a.sourceReader
}
