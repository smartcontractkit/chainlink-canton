package adapters

import (
	"encoding/binary"
	"math/big"

	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	datastore2 "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	seq_core "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutil "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

var _ lanes.LaneAdapter = &CantonLaneAdapter{}

type CantonLaneAdapter struct{}

func (c CantonLaneAdapter) ConfigureLaneLegAsSource() *operations.Sequence[lanes.UpdateLanesInput, seq_core.OnChainOutput, chain.BlockChains] {
	return sequences.ConfigureLaneLegAsSource
}

func (c CantonLaneAdapter) ConfigureLaneLegAsDest() *operations.Sequence[lanes.UpdateLanesInput, seq_core.OnChainOutput, chain.BlockChains] {
	return sequences.ConfigureLaneLegAsDest
}

func (c CantonLaneAdapter) DisableRemoteChain() *operations.Sequence[lanes.DisableRemoteChainInput, seq_core.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func (c CantonLaneAdapter) GetOnRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore2.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(onramp.ContractType),
		Version:       onramp.Version,
	}, chainSelector, dsutil.ToInstanceAddressBytes)
}

func (c CantonLaneAdapter) GetOffRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore2.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(offramp.ContractType),
		Version:       offramp.Version,
	}, chainSelector, dsutil.ToInstanceAddressBytes)
}

func (c CantonLaneAdapter) GetRouterAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return nil, nil
}

func (c CantonLaneAdapter) GetFQAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore2.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(fee_quoter.ContractType),
		Version:       fee_quoter.Version,
	}, chainSelector, dsutil.ToInstanceAddressBytes)
}

func (c CantonLaneAdapter) GetFeeQuoterDestChainConfig() lanes.FeeQuoterDestChainConfig {
	// TODO update Canton values
	return lanes.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      false,
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         binary.BigEndian.Uint32([]byte{0xdf, 0xaf, 0xaf, 0x4b}), // bytes4(keccak256("CCIP ChainFamilySelector Canton"))
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V1Params:                    nil,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 90,
			USDPerUnitGas:            big.NewInt(38),
		},
	}
}

func (c CantonLaneAdapter) GetDefaultGasPrice() *big.Int {
	// TODO calculate this
	return big.NewInt(38)
}
