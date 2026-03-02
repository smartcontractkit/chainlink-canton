package accessors

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
)

type factory struct {
	lggr          logger.Logger
	helper        map[string]*ccip.BlockchainInfo
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
	cantonConfig, ok := f.helper[strSelector]
	if !ok {
		return nil, fmt.Errorf("canton config not found for chain %d", chainSelector)
	}

	readerConfig, ok := f.readerConfigs[strSelector]
	if !ok {
		return nil, fmt.Errorf("canton reader config not found for chain %d", chainSelector)
	}

	sourceReader, err := sourcereader.NewSourceReader(
		logger.Named(f.lggr, fmt.Sprintf("CantonSourceReader.%d", chainSelector)),
		cantonConfig.GRPCLedgerAPIURL,
		cantonConfig.JWT,
		sourcereader.ReaderConfig{
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

func NewFactory(lggr logger.Logger, helper map[string]*ccip.BlockchainInfo, readerConfigs map[string]sourcereader.ReaderConfig) chainaccess.AccessorFactory {
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
