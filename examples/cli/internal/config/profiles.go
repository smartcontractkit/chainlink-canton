package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
)

// NetworkProfile contains all static, in-code configuration that differs
// between mainnet and testnet (chain selectors, well-known addresses,
// well-known Daml party IDs, instrument IDs).
type NetworkProfile struct {
	Name string

	// Service URLs
	EDSURL     string
	IndexerURL string

	// Canton
	CantonSelector   uint64
	DSOPartyID       types.PARTY
	CCIPOwnerPartyID types.PARTY

	// Instruments
	AmuletInstrumentID *splice_api_token_holding_v1.InstrumentId
	LinkInstrumentID   *splice_api_token_holding_v1.InstrumentId

	// EVM
	EthSelector          uint64
	EthRouterAddress     common.Address
	EthTokenAddress      common.Address // Address of the token used for bridging (e.g. LINK token)
	EthLinkTokenAddress  common.Address // LINK token address used for fee payment
	EthEmptyAddress      common.Address
	NoExecutionTag       common.Address
	CCIPReceiverContract common.Address

	// Explorers
	CCIPExlorerURL    string
	EVMExplorerURL    string
	CantonExplorerURL string
}

// Networks holds all supported network profiles.
var Networks = map[string]*NetworkProfile{
	"mainnet": Mainnet,
	"testnet": Testnet,
}

// Get returns the profile for the given network name or an error.
func Get(name string) (*NetworkProfile, error) {
	p, ok := Networks[name]
	if !ok {
		return nil, fmt.Errorf("unknown network %q (supported: mainnet, testnet)", name)
	}

	return p, nil
}

// Mainnet profile.
var Mainnet = &NetworkProfile{
	Name: "mainnet",

	EDSURL:     "https://eds.ccip.chain.link",
	IndexerURL: "https://indexer-1.ccip.chain.link",

	CantonSelector:   chainsel.CANTON_MAINNET.Selector,
	DSOPartyID:       types.PARTY("DSO::1220b1431ef217342db44d516bb9befde802be7d8899637d290895fa58880f19accc"),
	CCIPOwnerPartyID: types.PARTY("ccipOwner::122012714685760dc1927c4cfe119ce2126c48756154e95c06f5c181da05a5519093"),

	AmuletInstrumentID: &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY("DSO::1220b1431ef217342db44d516bb9befde802be7d8899637d290895fa58880f19accc"),
		Id:    "Amulet",
	},
	LinkInstrumentID: &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY("ccipOwner::122012714685760dc1927c4cfe119ce2126c48756154e95c06f5c181da05a5519093"),
		Id:    "link-token",
	},

	EthSelector:          chainsel.ETHEREUM_MAINNET.Selector,
	EthRouterAddress:     common.HexToAddress("0x80226fc0Ee2b096224EeAc085Bb9a8cba1146f7D"),
	EthTokenAddress:      common.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA"),
	EthLinkTokenAddress:  common.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA"),
	EthEmptyAddress:      common.HexToAddress(""),
	NoExecutionTag:       common.HexToAddress("0xEBa517d200000000000000000000000000000000"),
	CCIPReceiverContract: common.HexToAddress("0x6da8b12B70bfc9246EdBAb53b1C4dA93cF88CaDb"),

	CCIPExlorerURL:    "https://ccip.chain.link",
	EVMExplorerURL:    "https://etherscan.io",
	CantonExplorerURL: "https://lighthouse.xyz",
}

// Testnet profile.
var Testnet = &NetworkProfile{
	Name: "testnet",

	EDSURL:     "https://eds.testnet.ccip.chain.link",
	IndexerURL: "https://indexer-1.testnet.ccip.chain.link",

	CantonSelector: chainsel.CANTON_TESTNET.Selector,
	DSOPartyID: types.PARTY(
		"DSO::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337",
	),
	CCIPOwnerPartyID: types.PARTY(
		"ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551",
	),

	AmuletInstrumentID: &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(
			"DSO::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337",
		),
		Id: "Amulet",
	},
	LinkInstrumentID: &splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(
			"ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551",
		),
		Id: "link-token",
	},

	EthSelector:          chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	EthRouterAddress:     common.HexToAddress("0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59"),
	EthTokenAddress:      common.HexToAddress("0xeEe6675b20fE5950eb51361b93021D076289F612"),
	EthLinkTokenAddress:  common.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789"),
	EthEmptyAddress:      common.HexToAddress(""),
	NoExecutionTag:       common.HexToAddress("0xEBa517d200000000000000000000000000000000"),
	CCIPReceiverContract: common.HexToAddress("0x1E340B34Cc732a71f5e59804da7E645a52e10E1B"),

	CCIPExlorerURL:    "https://ccip.chain.link",
	EVMExplorerURL:    "https://sepolia.etherscan.io",
	CantonExplorerURL: "https://lighthouse.testnet.cantonloop.com",
}
