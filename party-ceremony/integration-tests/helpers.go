package integrationtests

import (
	"sync"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
)

var defaultNetworkOnce = &sync.Once{}

func LoadChainWithCTF(t *testing.T, numberOfValidators int) (*canton.Chain, error) {
	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: numberOfValidators,
		Once:               defaultNetworkOnce,
	}).Initialize(t.Context())
	require.NoError(t, err, "Failed to initialize CTF chain provider")

	return bc.(*canton.Chain), nil
}
