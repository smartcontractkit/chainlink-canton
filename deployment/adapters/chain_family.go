package adapters

import (
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonsequences "github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

var _ ccipadapters.ChainFamily = (*CantonChainFamilyAdapter)(nil)

// CantonChainFamilyAdapter bridges CCIP tooling's ChainFamily API
// into chainlink-canton's native deployment sequence implementation.
type CantonChainFamilyAdapter struct{}

func NewCantonChainFamilyAdapter() *CantonChainFamilyAdapter {
	return &CantonChainFamilyAdapter{}
}

func (c *CantonChainFamilyAdapter) ConfigureChainForLanes() *operations.Sequence[ccipadapters.ConfigureChainForLanesInput, seqcore.OnChainOutput, chain.BlockChains] {
	return cantonsequences.ConfigureChainForLanesAdapter
}

func (c *CantonChainFamilyAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

