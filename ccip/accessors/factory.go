package accessors

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"github.com/BurntSushi/toml"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

const CantonConfigPathEnv = "CANTON_CONFIG_PATH"

func init() {
	chainaccess.Register(chainsel.FamilyCanton, CreateCantonAccessorFactory)
}

func loadConfig(path string) (*ccip.Config, error) {
	var cfg ccip.Config
	if md, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %w", path, err)
	} else if len(md.Undecoded()) > 0 {
		return nil, fmt.Errorf("unknown fields in config: %v", md.Undecoded())
	}

	return &cfg, nil
}

func CreateCantonAccessorFactory(lggr logger.Logger, genericConfig chainaccess.GenericConfig) (chainaccess.AccessorFactory, error) {
	configPath, ok := os.LookupEnv(CantonConfigPathEnv)
	if !ok {
		configPath = ccip.DefaultCantonConfigPath
	}

	cantonConfig, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Don't log full config to avoid leaking sensitive fields.
	lggr.Infow("loaded canton config",
		"numChains", len(cantonConfig.BlockchainInfos),
		"numReaderConfigs", len(cantonConfig.ReaderConfigs),
	)

	return NewFactory(lggr, cantonConfig.BlockchainInfos, cantonConfig.ReaderConfigs, genericConfig.RMNRemoteAddresses), nil
}

type factory struct {
	lggr               logger.Logger
	helper             map[string]ccip.BlockchainInfo
	readerConfigs      map[string]sourcereader.ReaderConfig
	rmnRemoteAddresses map[string]string
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

	rmnRemoteAddressStr, ok := f.rmnRemoteAddresses[strSelector]
	if !ok {
		return nil, fmt.Errorf("RMN instance address not found for chain %d", chainSelector)
	}

	authProvider, err := blockchainInfo.Auth.NewProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth provider for chain %d: %w", chainSelector, err)
	}

	sourceReader, err := sourcereader.NewSourceReader(
		logger.Named(f.lggr, fmt.Sprintf("CantonSourceReader.%d", chainSelector)),
		blockchainInfo.GRPCLedgerAPIURL,
		readerConfig,
		contracts.HexToInstanceAddress(rmnRemoteAddressStr),
		grpc.WithTransportCredentials(authProvider.TransportCredentials()),
		grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create source reader: %w", err)
	}

	return newAccessor(sourceReader), nil
}

func NewFactory(
	lggr logger.Logger,
	helper map[string]ccip.BlockchainInfo,
	readerConfigs map[string]sourcereader.ReaderConfig,
	rmnRemoteAddresses map[string]string,
) chainaccess.AccessorFactory {
	return &factory{
		lggr:               lggr,
		helper:             helper,
		readerConfigs:      readerConfigs,
		rmnRemoteAddresses: rmnRemoteAddresses,
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

func (a *accessor) SourceReader() (chainaccess.SourceReader, error) {
	return a.sourceReader, nil
}

func (a *accessor) DestinationReader() (chainaccess.DestinationReader, error) {
	return nil, fmt.Errorf("not supported")
}

func (a *accessor) ContractTransmitter() (chainaccess.ContractTransmitter, error) {
	return nil, fmt.Errorf("not supported")
}
