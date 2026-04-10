package adapters

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	datastore2 "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutil "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

var _ lanes.LaneAdapter = (*CantonLaneAdapter)(nil)

type CantonLaneAdapter struct{}

func (c CantonLaneAdapter) ConfigureLaneLegAsSource() *cldfops.Sequence[lanes.UpdateLanesInput, seqcore.OnChainOutput, cldfchain.BlockChains] {
	return sequences.ConfigureLaneLegAsSource
}

func (c CantonLaneAdapter) ConfigureLaneLegAsDest() *cldfops.Sequence[lanes.UpdateLanesInput, seqcore.OnChainOutput, cldfchain.BlockChains] {
	return sequences.ConfigureLaneLegAsDest
}

func (c CantonLaneAdapter) DisableRemoteChain() *cldfops.Sequence[lanes.DisableRemoteChainInput, seqcore.OnChainOutput, cldfchain.BlockChains] {
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
	return DefaultCantonFeeQuoterDestChainConfig()
}

func (c CantonLaneAdapter) GetDefaultGasPrice() *big.Int {
	return big.NewInt(38)
}
