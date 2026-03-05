// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hybrid_with_external_minter_token_pool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type HybridTokenPoolAbstractGroupUpdate struct {
	RemoteChainSelector uint64
	Group               uint8
	RemoteChainSupply   *big.Int
}

type PoolLockOrBurnInV1 struct {
	Receiver            []byte
	RemoteChainSelector uint64
	OriginalSender      common.Address
	Amount              *big.Int
	LocalToken          common.Address
}

type PoolLockOrBurnOutV1 struct {
	DestTokenAddress []byte
	DestPoolData     []byte
}

type PoolReleaseOrMintInV1 struct {
	OriginalSender          []byte
	RemoteChainSelector     uint64
	Receiver                common.Address
	SourceDenominatedAmount *big.Int
	LocalToken              common.Address
	SourcePoolAddress       []byte
	SourcePoolData          []byte
	OffchainTokenData       []byte
}

type PoolReleaseOrMintOutV1 struct {
	DestinationAmount *big.Int
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

type TokenPoolChainUpdate struct {
	RemoteChainSelector       uint64
	RemotePoolAddresses       [][]byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var HybridWithExternalMinterTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGroup\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockedTokens\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRebalancer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"provideLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRebalancer\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"updateGroups\",\"inputs\":[{\"name\":\"groupUpdates\",\"type\":\"tuple[]\",\"internalType\":\"structHybridTokenPoolAbstract.GroupUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawLiquidity\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GroupUpdated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityAdded\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityMigrated\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumHybridTokenPoolAbstract.Group\"},{\"name\":\"remoteChainSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidityRemoved\",\"inputs\":[{\"name\":\"rebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RebalancerSet\",\"inputs\":[{\"name\":\"oldRebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRebalancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidGroupUpdate\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"group\",\"type\":\"uint8\",\"internalType\":\"enumHybridTokenPoolAbstract.Group\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LiquidityAmountCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x61012080604052346102c65761589c803803809161001d82856104db565b8339810160c0828203126102c657610034826104fe565b60208301516001600160a01b038116939192918482036102c65761005a60408201610512565b60608201519091906001600160401b0381116102c65781019380601f860112156102c6578451946001600160401b0386116104c5578560051b9060208201966100a660405198896104db565b87526020808801928201019283116102c657602001905b8282106104ad575050506100df60a06100d8608084016104fe565b92016104fe565b92331561049c57600180546001600160a01b031916331790558615801561048b575b801561047a575b6104695760805260c05260405163313ce56760e01b8152602081600481895afa6000918161042d575b50610402575b5060a052600480546001600160a01b0319166001600160a01b03929092169190911790558051151560e08190526102df575b506101008190526040516321df0da760e01b815290602090829060049082906001600160a01b03165afa9081156102d357600091610294575b506001600160a01b03169081810361027d576040516151db90816106c182396080518181816103b501528181610fcd01528181611c1101528181611def01528181612a8a01528181612c4701528181612fbe0152818161313e015281816131b601526132c5015260a051818181611f2a015281816130c901528181613d0a0152613d8d015260c0518181816111da01528181611cac0152612b26015260e05181818161116a01528181611cf001526127ba01526101005181818161020c01528181610ebd01528181610fa701528181611fe50152612e240152f35b63f902523f60e01b60005260045260245260446000fd5b90506020813d6020116102cb575b816102af602093836104db565b810103126102c6576102c0906104fe565b386101a2565b600080fd5b3d91506102a2565b6040513d6000823e3d90fd5b90602090604051906102f183836104db565b60008252600036813760e051156103f15760005b825181101561036c576001906001600160a01b036103238286610520565b51168561032f82610562565b61033c575b505001610305565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a13885610334565b5092905060005b81518110156103e7576001906001600160a01b036103918285610520565b511680156103e157846103a382610660565b6103b1575b50505b01610373565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a138846103a8565b506103ab565b5050506020610169565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036104165750610137565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d602011610461575b81610449602093836104db565b810103126102c65761045a90610512565b9038610131565b3d915061043c565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03821615610108565b506001600160a01b03841615610101565b639b15e16f60e01b60005260046000fd5b602080916104ba846104fe565b8152019101906100bd565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b038211908210176104c557604052565b51906001600160a01b03821682036102c657565b519060ff821682036102c657565b80518210156105345760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b80548210156105345760005260206000200190600090565b600081815260036020526040902054801561065957600019810181811161064357600254600019810191908211610643578181036105f2575b50505060025480156105dc57600019016105b681600261054a565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b61062b61060361061493600261054a565b90549060031b1c928392600261054a565b819391549060031b91821b91600019901b19161790565b9055600052600360205260406000205538808061059b565b634e487b7160e01b600052601160045260246000fd5b5050600090565b806000526003602052604060002054156000146106ba57600254680100000000000000008110156104c5576106a1610614826001859401600255600261054a565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714613317575080630a861f2a1461327e578063181f5a77146131da57806321df0da71461316b578063240028e8146130ed57806324f65ee714613091578063319ac1011461302c5780633317bbcc14612f4657806339077537146129b7578063432a6ba3146129655780634c5ef0ed1461290257806354c8a4f31461278857806362ddd3c4146127055780636cfd15531461265a5780636d3d1a581461260857806379ba50971461251f5780637d54534e146124745780638926f54f146124255780638da5cb5b146123d3578063962d40201461225f5780639a4575b914611b6a578063a42a7b8b146119de578063a7cd63b71461190c578063acfecf91146117e7578063af58d59f1461177f578063b0f479a11461172d578063b7946580146116d7578063c0d78655146115d6578063c4bffe2b14611488578063c75eea9c146113c1578063cf7401f3146111fe578063dc0bd9711461118f578063e0351e1314611134578063e7e62f8514610d16578063e8a1da1714610431578063eb521a4c14610325578063f2fde38b146102355763f3667517146101c157600080fd5b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305773ffffffffffffffffffffffffffffffffffffffff610281613463565b610289613e97565b163381146102fb57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057600435801561040757610365613be9565b6040517f23b872dd000000000000000000000000000000000000000000000000000000006020820152336024820152306044820152606480820183905281526103d9906103b360848261353c565b7f00000000000000000000000000000000000000000000000000000000000000006145d4565b6040519081527fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208860203392a2005b7fa90c0d190000000000000000000000000000000000000000000000000000000060005260046000fd5b346102305761043f3661363d565b91909261044a613e97565b6000905b828210610b6d5750505060009063ffffffff4216907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee184360301925b81811015610b6b576000918160051b86013585811215610b675786019061012082360312610b6757604051956104bf87613504565b823567ffffffffffffffff81168103610b63578752602083013567ffffffffffffffff8111610b635783019536601f88011215610b63578635966105028861388e565b97610510604051998a61353c565b8089526020808a019160051b83010190368211610b5f5760208301905b828210610b2c575050505060208801968752604084013567ffffffffffffffff8111610b285761056090369086016135ee565b926040890193845261058a610578366060880161377e565b9560608b0196875260c036910161377e565b9660808a0197885261059c865161448d565b6105a6885161448d565b84515115610b00576105c267ffffffffffffffff8b5116614e0c565b15610ac95767ffffffffffffffff8a5116815260076020526040812061070287516fffffffffffffffffffffffffffffffff604082015116906106bd6fffffffffffffffffffffffffffffffff6020830151169151151583608060405161062881613504565b858152602081018c905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff84161773ffffffff0000000000000000000000000000000060808b901b1617178555565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176001830155565b61082889516fffffffffffffffffffffffffffffffff604082015116906107e36fffffffffffffffffffffffffffffffff6020830151169151151583608060405161074c81613504565b858152602081018c9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166fffffffffffffffffffffffffffffffff85161773ffffffff0000000000000000000000000000000060808c901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176003830155565b6004865191019080519067ffffffffffffffff8211610a9c5761084b8354613971565b601f8111610a61575b50602090601f83116001146109c2576108a292918591836109b7575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b885180518210156108da57906108d46001926108cd8367ffffffffffffffff8f51169261395d565b5190613ee2565b016108a5565b5050977f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c29391999750956109a867ffffffffffffffff600197969498511692519351915161097461093f60405196879687526101006020880152610100870190613404565b9360408601906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60a08401906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b0390a10193919392909261048a565b015190508f80610870565b83855281852091907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08416865b818110610a495750908460019594939210610a12575b505050811b0190556108a5565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080610a05565b929360206001819287860151815501950193016109ef565b610a8c9084865260208620601f850160051c81019160208610610a92575b601f0160051c0190613b90565b8e610854565b9091508190610a7f565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60249067ffffffffffffffff8b51167f1d5ad3c5000000000000000000000000000000000000000000000000000000008252600452fd5b807f8579befe0000000000000000000000000000000000000000000000000000000060049252fd5b8680fd5b813567ffffffffffffffff8111610b5b57602091610b5083928336918901016135ee565b81520191019061052d565b8a80fd5b8880fd5b8580fd5b8380fd5b005b909267ffffffffffffffff610b8e610b8986868699979961390e565b61381b565b1692610b9984614b40565b15610ce857836000526007602052610bb76005604060002001614947565b9260005b8451811015610bf357600190866000526007602052610bec6005604060002001610be5838961395d565b5190614c6b565b5001610bbb565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610c328154613971565b9081610ca5575b5050018054906000815581610c84575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a101909193929361044e565b6000526020600020908101905b81811015610c495760008155600101610c91565b81601f60009311600114610cbd5750555b8880610c39565b81835260208320610cd891601f01861c810190600101613b90565b8082528160208120915555610cb6565b837f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760043567ffffffffffffffff811161023057610d65903690600401613730565b610d6d613e97565b60005b818110610d7957005b610d8481838561394d565b9067ffffffffffffffff610d978361381b565b16600052600b60205260ff6040600020541660208301359060028210801561023057600091600281101561110757831480156110d7575b61105c57506102305767ffffffffffffffff610de98461381b565b16600052600b6020526040600020926000937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541660ff8416179055604081013580610e78575b50610e3b9061381b565b926102305760019267ffffffffffffffff167f1d1eeb97006356bf772500dc592e232d913119a3143e8452f60e5c98b6a29ca1600080a301610d70565b6000945082610f8b576040517f40c10f1900000000000000000000000000000000000000000000000000000000815230600482015260248101829052602081604481897f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff165af18015610f8057918491610e3b9493610f52575b505b610f128361381b565b96507fbbaa9aea43e3358cd56e894ad9620b8a065abcffab21357fb0702f222480fccc602067ffffffffffffffff6000996040519485521692a390610e31565b610f729060203d8111610f79575b610f6a818361353c565b810190613b78565b5089610f07565b503d610f60565b6040513d88823e3d90fd5b84602073ffffffffffffffffffffffffffffffffffffffff60247f0000000000000000000000000000000000000000000000000000000000000000610ff186827f00000000000000000000000000000000000000000000000000000000000000006142e4565b60405194859384927f42966c68000000000000000000000000000000000000000000000000000000008452886004850152165af18015610f8057918491610e3b949361103e575b50610f09565b6110559060203d8111610f7957610f6a818361353c565b5089611038565b9067ffffffffffffffff906110708661381b565b90507fe2017d610000000000000000000000000000000000000000000000000000000060005216600452156110a85760245260446000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b5061110167ffffffffffffffff6110ed8761381b565b166000526006602052604060002054151590565b15610dce565b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526021600452fd5b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102305760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057611235613486565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc3601126102305760405161126b81613520565b60243580151581036102305781526044356fffffffffffffffffffffffffffffffff811681036102305760208201526064356fffffffffffffffffffffffffffffffff8116810361023057604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c36011261023057604051906112f282613520565b608435801515810361023057825260a4356fffffffffffffffffffffffffffffffff8116810361023057602083015260c4356fffffffffffffffffffffffffffffffff8116810361023057604083015273ffffffffffffffffffffffffffffffffffffffff600954163314158061139f575b61137157610b6b92614122565b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b5073ffffffffffffffffffffffffffffffffffffffff60015416331415611364565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305767ffffffffffffffff611401613486565b611409613ac5565b5016600052600760205261148461142b6114266040600020613af0565b61425f565b6040519182918291909160806fffffffffffffffffffffffffffffffff8160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610230576040516005548082528160208101600560005260206000209260005b8181106115bd5750506114e89250038261353c565b80519061150d6114f78361388e565b92611505604051948561353c565b80845261388e565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401920136833760005b815181101561156d578067ffffffffffffffff61155a6001938561395d565b5116611566828761395d565b520161153b565b5050906040519182916020830190602084525180915260408301919060005b81811061159a575050500390f35b825167ffffffffffffffff1684528594506020938401939092019160010161158c565b84548352600194850194869450602090930192016114d3565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305761160d613463565b611615613e97565b73ffffffffffffffffffffffffffffffffffffffff81169081156116ad57600480547fffffffffffffffffffffffff000000000000000000000000000000000000000081169390931790556040805173ffffffffffffffffffffffffffffffffffffffff93841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849190a1005b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057611484611719611714613486565b613b56565b604051918291602083526020830190613404565b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602073ffffffffffffffffffffffffffffffffffffffff60045416604051908152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305767ffffffffffffffff6117bf613486565b6117c7613ac5565b5016600052600760205261148461142b6114266002604060002001613af0565b346102305767ffffffffffffffff6117fe366136ad565b929091611809613e97565b1690611822826000526006602052604060002054151590565b156118de5781600052600760205261185360056040600020016118463686856135b7565b6020815191012090614c6b565b15611897577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d769192611892604051928392602084526020840191613a86565b0390a2005b6118da906040519384937f74f23c7c0000000000000000000000000000000000000000000000000000000085526004850152604060248501526044840191613a86565b0390fd5b507f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760405160025490818152602081018092600260005260206000209060005b8181106119c8575050508161196f91038261353c565b6040519182916020830190602084525180915260408301919060005b818110611999575050500390f35b825173ffffffffffffffffffffffffffffffffffffffff1684528594506020938401939092019160010161198b565b8254845260209093019260019283019201611959565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305767ffffffffffffffff611a1e613486565b166000526007602052611a376005604060002001614947565b8051907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0611a7d611a678461388e565b93611a75604051958661353c565b80855261388e565b0160005b818110611b5957505060005b8151811015611ad55780611aa36001928461395d565b516000526008602052611ab960406000206139c4565b611ac3828661395d565b52611ace818561395d565b5001611a8d565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b828210611b0e57505050500390f35b91936020611b49827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc060019597998495030186528851613404565b9601920192018594939192611aff565b806060602080938701015201611a81565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760043567ffffffffffffffff81116102305760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126102305760606020604051611be7816134e8565b828152015260848101611bf981613830565b73ffffffffffffffffffffffffffffffffffffffff807f00000000000000000000000000000000000000000000000000000000000000001691160361221357506024810177ffffffffffffffff00000000000000000000000000000000611c5f8261381b565b60801b16604051907f2cbc26bb000000000000000000000000000000000000000000000000000000008252600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156120a3576000916121f4575b506121ca57611cee60448301613830565b7f0000000000000000000000000000000000000000000000000000000000000000612174575b5067ffffffffffffffff611d278261381b565b16611d3f816000526006602052604060002054151590565b1561214757602073ffffffffffffffffffffffffffffffffffffffff60045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa9081156120a3576000916120dd575b5073ffffffffffffffffffffffffffffffffffffffff1633036120af5767ffffffffffffffff916064611dd48361381b565b91013592839116806000526007602052604060002090611e2e7f00000000000000000000000000000000000000000000000000000000000000009273ffffffffffffffffffffffffffffffffffffffff8416968791614ebb565b6040805173ffffffffffffffffffffffffffffffffffffffff87168152602081018590527fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da81789449190a267ffffffffffffffff611e888461381b565b16600052600b60205260ff6040600020541660028110156110a857600114611fc1575b611f90611f2061171485877ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae108767ffffffffffffffff611eea8561381b565b6040805173ffffffffffffffffffffffffffffffffffffffff9690961686523360208701528501929092521691606090a261381b565b61148460405160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260208152611f5e60408261353c565b60405192611f6b846134e8565b8352602083019081526040519384936020855251604060208601526060850190613404565b90517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152613404565b91602073ffffffffffffffffffffffffffffffffffffffff6024849561200b6000967f000000000000000000000000000000000000000000000000000000000000000080936142e4565b60405195869384927f42966c68000000000000000000000000000000000000000000000000000000008452896004850152165af19384156120a3577ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae10611f209461171494611f9097612084575b50935050935091611eab565b61209c9060203d602011610f7957610f6a818361353c565b5087612078565b6040513d6000823e3d90fd5b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d60201161213f575b816120f66020938361353c565b8101031261213b57519073ffffffffffffffffffffffffffffffffffffffff82168203612138575073ffffffffffffffffffffffffffffffffffffffff611da2565b80fd5b5080fd5b3d91506120e9565b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff16806000526003602052604060002054611d14577fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b61220d915060203d602011610f7957610f6a818361353c565b83611cdd565b61223173ffffffffffffffffffffffffffffffffffffffff91613830565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346102305760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760043567ffffffffffffffff8111610230576122ae90369060040161360c565b9060243567ffffffffffffffff8111610230576122cf903690600401613730565b9060443567ffffffffffffffff8111610230576122f0903690600401613730565b73ffffffffffffffffffffffffffffffffffffffff60095416331415806123b1575b611371578386148015906123a7575b61237d5760005b86811061233157005b80612377612345610b896001948b8b61390e565b61235083898961394d565b61237161236961236186898b61394d565b92369061377e565b91369061377e565b91614122565b01612328565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b5080861415612321565b5073ffffffffffffffffffffffffffffffffffffffff60015416331415612312565b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602061246a67ffffffffffffffff6110ed613486565b6040519015158152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610230577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174602073ffffffffffffffffffffffffffffffffffffffff6124e3613463565b6124eb613e97565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000006009541617600955604051908152a1005b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760005473ffffffffffffffffffffffffffffffffffffffff811633036125de577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602073ffffffffffffffffffffffffffffffffffffffff60095416604051908152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057612691613463565b612699613e97565b73ffffffffffffffffffffffffffffffffffffffff80600a54921691827fffffffffffffffffffffffff0000000000000000000000000000000000000000821617600a55167f64187bd7b97e66658c91904f3021d7c28de967281d18b1a20742348afdd6a6b3600080a3005b3461023057612713366136ad565b61271e929192613e97565b67ffffffffffffffff8216612740816000526006602052604060002054151590565b1561275b5750610b6b926127559136916135b7565b90613ee2565b7f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b34610230576127b06127b861279c3661363d565b94916127a9939193613e97565b36916138a6565b9236916138a6565b7f0000000000000000000000000000000000000000000000000000000000000000156128d85760005b8251811015612854578073ffffffffffffffffffffffffffffffffffffffff61280c6001938661395d565b5116612817816149aa565b612823575b50016127e1565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a18461281c565b5060005b8151811015610b6b578073ffffffffffffffffffffffffffffffffffffffff6128836001938561395d565b511680156128d25761289481614dac565b6128a1575b505b01612858565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a183612899565b5061289b565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b346102305760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057612939613486565b60243567ffffffffffffffff81116102305760209161295f61246a9236906004016135ee565b90613851565b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602073ffffffffffffffffffffffffffffffffffffffff600a5416604051908152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305760043567ffffffffffffffff811161023057806004016101007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8336030112610230576000604051612a378161349d565b52612a64612a5a612a55612a4e60c48601856137ca565b36916135b7565b613c96565b6064840135613d8a565b9060848301612a7281613830565b73ffffffffffffffffffffffffffffffffffffffff807f0000000000000000000000000000000000000000000000000000000000000000169116036122135750602483019077ffffffffffffffff00000000000000000000000000000000612ad98361381b565b60801b16604051907f2cbc26bb000000000000000000000000000000000000000000000000000000008252600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa9081156120a357600091612f27575b506121ca5767ffffffffffffffff612b6e8361381b565b16612b86816000526006602052604060002054151590565b1561214757602073ffffffffffffffffffffffffffffffffffffffff60045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa9081156120a357600091612f08575b50156120af57612bfe8261381b565b90612c1460a486019261295f612a4e85856137ca565b15612ec15750508067ffffffffffffffff612c2f849361381b565b16806000526007602052600260406000200190612c867f00000000000000000000000000000000000000000000000000000000000000009273ffffffffffffffffffffffffffffffffffffffff8416958691614ebb565b6040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018790527f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c9190a267ffffffffffffffff612ce08361381b565b16600052600b602052604460ff60406000205416950194612d0086613830565b9060028110156110a857612db55760209573ffffffffffffffffffffffffffffffffffffffff612d74612d6e7ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0966080968a67ffffffffffffffff973087821603612da4575b50505061381b565b92613830565b60405196875233898801521660408601528560608601521692a280604051612d9b8161349d565b52604051908152f35b612dad92613c0a565b8b8a81612d66565b6040517f40c10f1900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101859052905060208180604481010381600073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af180156120a35760209573ffffffffffffffffffffffffffffffffffffffff612d74612d6e7ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc09660809667ffffffffffffffff96612ea4575b5061381b565b612eba908c3d8e11610f7957610f6a818361353c565b508b612e9e565b612ecb92506137ca565b6118da6040519283927f24eb47e5000000000000000000000000000000000000000000000000000000008452602060048501526024840191613a86565b612f21915060203d602011610f7957610f6a818361353c565b85612bef565b612f40915060203d602011610f7957610f6a818361353c565b85612b57565b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610230576040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa80156120a357600090612ff9575b602090604051908152f35b506020813d602011613024575b816130136020938361353c565b810103126102305760209051612fee565b3d9150613006565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305767ffffffffffffffff61306c613486565b16600052600b60205260ff6040600020541660405160028210156110a8576020918152f35b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610230576020613126613463565b73ffffffffffffffffffffffffffffffffffffffff807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346102305760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102305761148460405161321a60608261353c565b602781527f4879627269645769746845787465726e616c4d696e746572546f6b656e506f6f60208201527f6c20312e362e30000000000000000000000000000000000000000000000000006040820152604051918291602083526020830190613404565b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610230576004358015610407576132be613be9565b6132e981337f0000000000000000000000000000000000000000000000000000000000000000613c0a565b6040519081527fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171960203392a2005b346102305760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261023057600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361023057817faff2afbf00000000000000000000000000000000000000000000000000000000602093149081156133da575b81156133b0575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014836133a9565b7f0e64dd2900000000000000000000000000000000000000000000000000000000811491506133a2565b919082519283825260005b84811061344e5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b8060208092840101518282860101520161340f565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361023057565b6004359067ffffffffffffffff8216820361023057565b6020810190811067ffffffffffffffff8211176134b957604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040810190811067ffffffffffffffff8211176134b957604052565b60a0810190811067ffffffffffffffff8211176134b957604052565b6060810190811067ffffffffffffffff8211176134b957604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176134b957604052565b67ffffffffffffffff81116134b957601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9291926135c38261357d565b916135d1604051938461353c565b829481845281830111610230578281602093846000960137010152565b9080601f8301121561023057816020613609933591016135b7565b90565b9181601f840112156102305782359167ffffffffffffffff8311610230576020808501948460051b01011161023057565b60407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126102305760043567ffffffffffffffff811161023057816136869160040161360c565b929092916024359067ffffffffffffffff8211610230576136a99160040161360c565b9091565b60407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126102305760043567ffffffffffffffff81168103610230579160243567ffffffffffffffff811161023057826023820112156102305780600401359267ffffffffffffffff84116102305760248483010111610230576024019190565b9181601f840112156102305782359167ffffffffffffffff8311610230576020808501946060850201011161023057565b35906fffffffffffffffffffffffffffffffff8216820361023057565b91908260609103126102305760405161379681613520565b809280359081151582036102305760406137c591819385526137ba60208201613761565b602086015201613761565b910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe181360301821215610230570180359067ffffffffffffffff82116102305760200191813603831361023057565b3567ffffffffffffffff811681036102305790565b3573ffffffffffffffffffffffffffffffffffffffff811681036102305790565b9067ffffffffffffffff61360992166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b67ffffffffffffffff81116134b95760051b60200190565b92916138b18261388e565b936138bf604051958661353c565b602085848152019260051b810191821161023057915b8183106138e157505050565b823573ffffffffffffffffffffffffffffffffffffffff81168103610230578152602092830192016138d5565b919081101561391e5760051b0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b919081101561391e576060020190565b805182101561391e5760209160051b010190565b90600182811c921680156139ba575b602083101461398b57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f1691613980565b90604051918260008254926139d884613971565b8084529360018116908115613a4657506001146139ff575b506139fd9250038361353c565b565b90506000929192526020600020906000915b818310613a2a5750509060206139fd92820101386139f0565b6020919350806001915483858901015201910190918492613a11565b602093506139fd9592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b820101386139f0565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b60405190613ad282613504565b60006080838281528260208201528260408201528260608201520152565b90604051613afd81613504565b60806001829460ff81546fffffffffffffffffffffffffffffffff8116865263ffffffff81861c16602087015260a01c161515604085015201546fffffffffffffffffffffffffffffffff81166060840152811c910152565b67ffffffffffffffff16600052600760205261360960046040600020016139c4565b90816020910312610230575180151581036102305790565b818110613b9b575050565b60008155600101613b90565b81810292918115918404141715613bba57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff600a5416330361137157565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff909216602483015260448201929092526139fd91613c9182606481015b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810184528361353c565b6145d4565b80518015613d0657602003613cc857805160208281019183018390031261023057519060ff8211613cc8575060ff1690565b6118da906040519182917f953576f7000000000000000000000000000000000000000000000000000000008352602060048401526024830190613404565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff8211613bba57565b60ff16604d8111613bba57600a0a90565b8115613d5b570490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff811692828414613e9057828411613e665790613dcf91613d2c565b91604d60ff8416118015613e2d575b613df757505090613df161360992613d40565b90613ba7565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b50613e3783613d40565b8015613d5b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048411613dde565b613e6f91613d2c565b91604d60ff841611613df757505090613e8a61360992613d40565b90613d51565b5050505090565b73ffffffffffffffffffffffffffffffffffffffff600154163303613eb857565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156116ad5767ffffffffffffffff81516020830120921691826000526007602052613f17816005604060002001614e66565b156140de5760005260086020526040600020815167ffffffffffffffff81116134b957613f448254613971565b601f81116140ac575b506020601f8211600114613fe65791613fc0827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea9593613fd695600091613fdb575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055604051918291602083526020830190613404565b0390a2565b905084015138613f8f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169083600052806000209160005b818110614094575092613fd69492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea98961061405d575b5050811b019055611719565b8501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690553880614051565b9192602060018192868a015181550194019201614016565b6140d890836000526020600020601f840160051c81019160208510610a9257601f0160051c0190613b90565b38613f4d565b50906118da6040519283927f393b8ad20000000000000000000000000000000000000000000000000000000084526004840152604060248401526044830190613404565b67ffffffffffffffff166000818152600660205260409020549092919015614224579161422160e0926141ed856141797f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b9761448d565b846000526007602052614190816040600020614714565b6141998361448d565b8460005260076020526141b3836002604060002001614714565b60405194855260208501906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60808301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565ba1565b827f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b91908203918211613bba57565b614267613ac5565b506fffffffffffffffffffffffffffffffff6060820151166fffffffffffffffffffffffffffffffff80835116916142c460208501936142be6142b163ffffffff87511642614252565b8560808901511690613ba7565b90614d9f565b808210156142dd57505b16825263ffffffff4216905290565b90506142ce565b919091811580156143d9575b15614355576040517f095ea7b300000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff909316602484015260448301919091526139fd9190613c918260648101613c65565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff841660248201526020818060448101038173ffffffffffffffffffffffffffffffffffffffff86165afa9081156120a35760009161445b575b50156142f0565b90506020813d602011614485575b816144766020938361353c565b81010312610230575138614454565b3d9150614469565b80511561452d576fffffffffffffffffffffffffffffffff6040820151166fffffffffffffffffffffffffffffffff602083015116106144ca5750565b60649061452b604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565bfd5b6fffffffffffffffffffffffffffffffff604082015116158015906145b5575b6145545750565b60649061452b604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b506fffffffffffffffffffffffffffffffff602082015116151561454d565b73ffffffffffffffffffffffffffffffffffffffff614663911691604092600080855193614602878661353c565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564602086015260208151910182855af13d1561470c573d916146478361357d565b926146548751948561353c565b83523d6000602085013e615102565b8051908161467057505050565b602080614681938301019101613b78565b156146895750565b608490517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606091615102565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c199161484d606092805461475163ffffffff8260801c1642614252565b908161488c575b50506fffffffffffffffffffffffffffffffff600181602086015116928281541680851060001461488457508280855b16167fffffffffffffffffffffffffffffffff000000000000000000000000000000008254161781556148018651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601517fffffffffffffffffffffffffffffffff0000000000000000000000000000000060809190911b16939092166fffffffffffffffffffffffffffffffff1692909217910155565b61422160405180926fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b838091614788565b6fffffffffffffffffffffffffffffffff916148c18392836148ba6001880154948286169560801c90613ba7565b9116614d9f565b8082101561494057505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781553880614758565b90506148cb565b906040519182815491828252602082019060005260206000209260005b8181106149795750506139fd9250038361353c565b8454835260019485019487945060209093019201614964565b805482101561391e5760005260206000200190600090565b6000818152600360205260409020548015614b39577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111613bba57600254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211613bba57818103614aca575b5050506002548015614a9b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01614a58816002614992565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600255600052600360205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b614b21614adb614aec936002614992565b90549060031b1c9283926002614992565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b90556000526003602052604060002055388080614a1f565b5050600090565b6000818152600660205260409020548015614b39577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111613bba57600554907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211613bba57818103614c31575b5050506005548015614a9b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01614bee816005614992565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600555600052600660205260006040812055600190565b614c53614c42614aec936005614992565b90549060031b1c9283926005614992565b90556000526006602052604060002055388080614bb5565b9060018201918160005282602052604060002054801515600014614d96577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101818111613bba578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201918211613bba57818103614d5f575b50505080548015614a9b577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190614d208282614992565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b614d7f614d6f614aec9386614992565b90549060031b1c92839286614992565b905560005283602052604060002055388080614ce8565b50505050600090565b91908201809211613bba57565b80600052600360205260406000205415600014614e0657600254680100000000000000008110156134b957614ded614aec8260018594016002556002614992565b9055600254906000526003602052604060002055600190565b50600090565b80600052600660205260406000205415600014614e0657600554680100000000000000008110156134b957614e4d614aec8260018594016005556005614992565b9055600554906000526006602052604060002055600190565b6000828152600182016020526040902054614b3957805490680100000000000000008210156134b95782614ea4614aec846001809601855584614992565b905580549260005201602052604060002055600190565b9182549060ff8260a01c161580156150fa575b6150f4576fffffffffffffffffffffffffffffffff82169160018501908154614f1363ffffffff6fffffffffffffffffffffffffffffffff83169360801c1642614252565b9081615056575b505084811061500a5750838310614f74575050614f496fffffffffffffffffffffffffffffffff928392614252565b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416179055565b5460801c91614f838185614252565b927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810190808211613bba57614fd1614fd69273ffffffffffffffffffffffffffffffffffffffff96614d9f565b613d51565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b828573ffffffffffffffffffffffffffffffffffffffff927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b8286929396116150ca57615071926142be9160801c90613ba7565b808410156150c55750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff0000000000000000000000000000000016178655923880614f1a565b61507c565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215614ece565b9192901561517d5750815115615116575090565b3b1561511f5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b8251909150156151905750805190602001fd5b6118da906040519182917f08c379a000000000000000000000000000000000000000000000000000000000835260206004840152602483019061340456fea164736f6c634300081a000a",
}

var HybridWithExternalMinterTokenPoolABI = HybridWithExternalMinterTokenPoolMetaData.ABI

var HybridWithExternalMinterTokenPoolBin = HybridWithExternalMinterTokenPoolMetaData.Bin

func DeployHybridWithExternalMinterTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *HybridWithExternalMinterTokenPool, error) {
	parsed, err := HybridWithExternalMinterTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HybridWithExternalMinterTokenPoolBin), backend, minter, token, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HybridWithExternalMinterTokenPool{address: address, abi: *parsed, HybridWithExternalMinterTokenPoolCaller: HybridWithExternalMinterTokenPoolCaller{contract: contract}, HybridWithExternalMinterTokenPoolTransactor: HybridWithExternalMinterTokenPoolTransactor{contract: contract}, HybridWithExternalMinterTokenPoolFilterer: HybridWithExternalMinterTokenPoolFilterer{contract: contract}}, nil
}

type HybridWithExternalMinterTokenPool struct {
	address common.Address
	abi     abi.ABI
	HybridWithExternalMinterTokenPoolCaller
	HybridWithExternalMinterTokenPoolTransactor
	HybridWithExternalMinterTokenPoolFilterer
}

type HybridWithExternalMinterTokenPoolCaller struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type HybridWithExternalMinterTokenPoolSession struct {
	Contract     *HybridWithExternalMinterTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type HybridWithExternalMinterTokenPoolCallerSession struct {
	Contract *HybridWithExternalMinterTokenPoolCaller
	CallOpts bind.CallOpts
}

type HybridWithExternalMinterTokenPoolTransactorSession struct {
	Contract     *HybridWithExternalMinterTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type HybridWithExternalMinterTokenPoolRaw struct {
	Contract *HybridWithExternalMinterTokenPool
}

type HybridWithExternalMinterTokenPoolCallerRaw struct {
	Contract *HybridWithExternalMinterTokenPoolCaller
}

type HybridWithExternalMinterTokenPoolTransactorRaw struct {
	Contract *HybridWithExternalMinterTokenPoolTransactor
}

func NewHybridWithExternalMinterTokenPool(address common.Address, backend bind.ContractBackend) (*HybridWithExternalMinterTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(HybridWithExternalMinterTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindHybridWithExternalMinterTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPool{address: address, abi: abi, HybridWithExternalMinterTokenPoolCaller: HybridWithExternalMinterTokenPoolCaller{contract: contract}, HybridWithExternalMinterTokenPoolTransactor: HybridWithExternalMinterTokenPoolTransactor{contract: contract}, HybridWithExternalMinterTokenPoolFilterer: HybridWithExternalMinterTokenPoolFilterer{contract: contract}}, nil
}

func NewHybridWithExternalMinterTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*HybridWithExternalMinterTokenPoolCaller, error) {
	contract, err := bindHybridWithExternalMinterTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolCaller{contract: contract}, nil
}

func NewHybridWithExternalMinterTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*HybridWithExternalMinterTokenPoolTransactor, error) {
	contract, err := bindHybridWithExternalMinterTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolTransactor{contract: contract}, nil
}

func NewHybridWithExternalMinterTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*HybridWithExternalMinterTokenPoolFilterer, error) {
	contract, err := bindHybridWithExternalMinterTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolFilterer{contract: contract}, nil
}

func bindHybridWithExternalMinterTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HybridWithExternalMinterTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HybridWithExternalMinterTokenPool.Contract.HybridWithExternalMinterTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.HybridWithExternalMinterTokenPoolTransactor.contract.Transfer(opts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.HybridWithExternalMinterTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HybridWithExternalMinterTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.contract.Transfer(opts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetAllowList(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetAllowList(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetAllowListEnabled(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetAllowListEnabled(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetCurrentInboundRateLimiterState(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetCurrentInboundRateLimiterState(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetGroup(opts *bind.CallOpts, remoteChainSelector uint64) (uint8, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getGroup", remoteChainSelector)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetGroup(remoteChainSelector uint64) (uint8, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetGroup(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetGroup(remoteChainSelector uint64) (uint8, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetGroup(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetLockedTokens(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getLockedTokens")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetLockedTokens() (*big.Int, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetLockedTokens(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetLockedTokens() (*big.Int, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetLockedTokens(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetMinter() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetMinter(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetMinter() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetMinter(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRateLimitAdmin(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRateLimitAdmin(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRebalancer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRebalancer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRebalancer() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRebalancer(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRebalancer() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRebalancer(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRemotePools(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRemotePools(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRemoteToken(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRemoteToken(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRmnProxy(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRmnProxy(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetRouter() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRouter(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetRouter(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetSupportedChains(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetSupportedChains(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetToken() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetToken(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetToken(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetTokenDecimals(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _HybridWithExternalMinterTokenPool.Contract.GetTokenDecimals(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsRemotePool(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsRemotePool(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsSupportedChain(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsSupportedChain(&_HybridWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsSupportedToken(&_HybridWithExternalMinterTokenPool.CallOpts, token)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.IsSupportedToken(&_HybridWithExternalMinterTokenPool.CallOpts, token)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) Owner() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.Owner(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) Owner() (common.Address, error) {
	return _HybridWithExternalMinterTokenPool.Contract.Owner(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SupportsInterface(&_HybridWithExternalMinterTokenPool.CallOpts, interfaceId)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SupportsInterface(&_HybridWithExternalMinterTokenPool.CallOpts, interfaceId)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HybridWithExternalMinterTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) TypeAndVersion() (string, error) {
	return _HybridWithExternalMinterTokenPool.Contract.TypeAndVersion(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _HybridWithExternalMinterTokenPool.Contract.TypeAndVersion(&_HybridWithExternalMinterTokenPool.CallOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.AcceptOwnership(&_HybridWithExternalMinterTokenPool.TransactOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.AcceptOwnership(&_HybridWithExternalMinterTokenPool.TransactOpts)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.AddRemotePool(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.AddRemotePool(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ApplyAllowListUpdates(&_HybridWithExternalMinterTokenPool.TransactOpts, removes, adds)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ApplyAllowListUpdates(&_HybridWithExternalMinterTokenPool.TransactOpts, removes, adds)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ApplyChainUpdates(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ApplyChainUpdates(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.LockOrBurn(&_HybridWithExternalMinterTokenPool.TransactOpts, lockOrBurnIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.LockOrBurn(&_HybridWithExternalMinterTokenPool.TransactOpts, lockOrBurnIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "provideLiquidity", amount)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ProvideLiquidity(&_HybridWithExternalMinterTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ProvideLiquidity(&_HybridWithExternalMinterTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ReleaseOrMint(&_HybridWithExternalMinterTokenPool.TransactOpts, releaseOrMintIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.ReleaseOrMint(&_HybridWithExternalMinterTokenPool.TransactOpts, releaseOrMintIn)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.RemoveRemotePool(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.RemoveRemotePool(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfig(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfig(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "setChainRateLimiterConfigs", remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfigs(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfigs(&_HybridWithExternalMinterTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRateLimitAdmin(&_HybridWithExternalMinterTokenPool.TransactOpts, rateLimitAdmin)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRateLimitAdmin(&_HybridWithExternalMinterTokenPool.TransactOpts, rateLimitAdmin)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "setRebalancer", rebalancer)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRebalancer(&_HybridWithExternalMinterTokenPool.TransactOpts, rebalancer)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRebalancer(&_HybridWithExternalMinterTokenPool.TransactOpts, rebalancer)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRouter(&_HybridWithExternalMinterTokenPool.TransactOpts, newRouter)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.SetRouter(&_HybridWithExternalMinterTokenPool.TransactOpts, newRouter)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.TransferOwnership(&_HybridWithExternalMinterTokenPool.TransactOpts, to)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.TransferOwnership(&_HybridWithExternalMinterTokenPool.TransactOpts, to)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "updateGroups", groupUpdates)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) UpdateGroups(groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) UpdateGroups(groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.UpdateGroups(&_HybridWithExternalMinterTokenPool.TransactOpts, groupUpdates)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactor) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.contract.Transact(opts, "withdrawLiquidity", amount)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.WithdrawLiquidity(&_HybridWithExternalMinterTokenPool.TransactOpts, amount)
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolTransactorSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _HybridWithExternalMinterTokenPool.Contract.WithdrawLiquidity(&_HybridWithExternalMinterTokenPool.TransactOpts, amount)
}

type HybridWithExternalMinterTokenPoolAllowListAddIterator struct {
	Event *HybridWithExternalMinterTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolAllowListAdd)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolAllowListAdd)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolAllowListAddIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolAllowListAdd)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*HybridWithExternalMinterTokenPoolAllowListAdd, error) {
	event := new(HybridWithExternalMinterTokenPoolAllowListAdd)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolAllowListRemoveIterator struct {
	Event *HybridWithExternalMinterTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolAllowListRemove)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolAllowListRemove)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolAllowListRemoveIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolAllowListRemove)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*HybridWithExternalMinterTokenPoolAllowListRemove, error) {
	event := new(HybridWithExternalMinterTokenPoolAllowListRemove)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolChainAddedIterator struct {
	Event *HybridWithExternalMinterTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolChainAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolChainAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainAddedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolChainAddedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolChainAdded)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseChainAdded(log types.Log) (*HybridWithExternalMinterTokenPoolChainAdded, error) {
	event := new(HybridWithExternalMinterTokenPoolChainAdded)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolChainConfiguredIterator struct {
	Event *HybridWithExternalMinterTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolChainConfigured)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolChainConfigured)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolChainConfiguredIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolChainConfigured)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseChainConfigured(log types.Log) (*HybridWithExternalMinterTokenPoolChainConfigured, error) {
	event := new(HybridWithExternalMinterTokenPoolChainConfigured)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolChainRemovedIterator struct {
	Event *HybridWithExternalMinterTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolChainRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolChainRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolChainRemovedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolChainRemoved)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseChainRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolChainRemoved, error) {
	event := new(HybridWithExternalMinterTokenPoolChainRemoved)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolConfigChangedIterator struct {
	Event *HybridWithExternalMinterTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolConfigChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolConfigChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolConfigChangedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolConfigChanged)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseConfigChanged(log types.Log) (*HybridWithExternalMinterTokenPoolConfigChanged, error) {
	event := new(HybridWithExternalMinterTokenPoolConfigChanged)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolGroupUpdatedIterator struct {
	Event *HybridWithExternalMinterTokenPoolGroupUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolGroupUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolGroupUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolGroupUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolGroupUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolGroupUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolGroupUpdated struct {
	RemoteChainSelector uint64
	Group               uint8
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterTokenPoolGroupUpdatedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "GroupUpdated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolGroupUpdatedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "GroupUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolGroupUpdated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "GroupUpdated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolGroupUpdated)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "GroupUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseGroupUpdated(log types.Log) (*HybridWithExternalMinterTokenPoolGroupUpdated, error) {
	event := new(HybridWithExternalMinterTokenPoolGroupUpdated)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "GroupUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator struct {
	Event *HybridWithExternalMinterTokenPoolInboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolInboundRateLimitConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolInboundRateLimitConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolInboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "InboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolInboundRateLimitConsumed)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseInboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterTokenPoolInboundRateLimitConsumed, error) {
	event := new(HybridWithExternalMinterTokenPoolInboundRateLimitConsumed)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolLiquidityAddedIterator struct {
	Event *HybridWithExternalMinterTokenPoolLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolLiquidityAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolLiquidityAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolLiquidityAdded struct {
	Rebalancer common.Address
	Amount     *big.Int
	Raw        types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterTokenPoolLiquidityAddedIterator, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "LiquidityAdded", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolLiquidityAddedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityAdded, rebalancer []common.Address) (event.Subscription, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "LiquidityAdded", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolLiquidityAdded)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseLiquidityAdded(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityAdded, error) {
	event := new(HybridWithExternalMinterTokenPoolLiquidityAdded)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolLiquidityMigratedIterator struct {
	Event *HybridWithExternalMinterTokenPoolLiquidityMigrated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolLiquidityMigratedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolLiquidityMigrated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolLiquidityMigrated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolLiquidityMigratedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolLiquidityMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolLiquidityMigrated struct {
	RemoteChainSelector uint64
	Group               uint8
	RemoteChainSupply   *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterLiquidityMigrated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterTokenPoolLiquidityMigratedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "LiquidityMigrated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolLiquidityMigratedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "LiquidityMigrated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchLiquidityMigrated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityMigrated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}
	var groupRule []interface{}
	for _, groupItem := range group {
		groupRule = append(groupRule, groupItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "LiquidityMigrated", remoteChainSelectorRule, groupRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolLiquidityMigrated)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityMigrated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseLiquidityMigrated(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityMigrated, error) {
	event := new(HybridWithExternalMinterTokenPoolLiquidityMigrated)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolLiquidityRemovedIterator struct {
	Event *HybridWithExternalMinterTokenPoolLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolLiquidityRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolLiquidityRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolLiquidityRemoved struct {
	Rebalancer common.Address
	Amount     *big.Int
	Raw        types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterTokenPoolLiquidityRemovedIterator, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "LiquidityRemoved", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolLiquidityRemovedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityRemoved, rebalancer []common.Address) (event.Subscription, error) {

	var rebalancerRule []interface{}
	for _, rebalancerItem := range rebalancer {
		rebalancerRule = append(rebalancerRule, rebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "LiquidityRemoved", rebalancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolLiquidityRemoved)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseLiquidityRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityRemoved, error) {
	event := new(HybridWithExternalMinterTokenPoolLiquidityRemoved)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolLockedOrBurnedIterator struct {
	Event *HybridWithExternalMinterTokenPoolLockedOrBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolLockedOrBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolLockedOrBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolLockedOrBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolLockedOrBurnedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolLockedOrBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolLockedOrBurned struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolLockedOrBurnedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolLockedOrBurnedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "LockedOrBurned", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolLockedOrBurned)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseLockedOrBurned(log types.Log) (*HybridWithExternalMinterTokenPoolLockedOrBurned, error) {
	event := new(HybridWithExternalMinterTokenPoolLockedOrBurned)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator struct {
	Event *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "OutboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseOutboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed, error) {
	event := new(HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator struct {
	Event *HybridWithExternalMinterTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolOwnershipTransferRequested)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*HybridWithExternalMinterTokenPoolOwnershipTransferRequested, error) {
	event := new(HybridWithExternalMinterTokenPoolOwnershipTransferRequested)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolOwnershipTransferredIterator struct {
	Event *HybridWithExternalMinterTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolOwnershipTransferredIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolOwnershipTransferred)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*HybridWithExternalMinterTokenPoolOwnershipTransferred, error) {
	event := new(HybridWithExternalMinterTokenPoolOwnershipTransferred)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator struct {
	Event *HybridWithExternalMinterTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolRateLimitAdminSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolRateLimitAdminSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolRateLimitAdminSet)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*HybridWithExternalMinterTokenPoolRateLimitAdminSet, error) {
	event := new(HybridWithExternalMinterTokenPoolRateLimitAdminSet)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolRebalancerSetIterator struct {
	Event *HybridWithExternalMinterTokenPoolRebalancerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolRebalancerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolRebalancerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolRebalancerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolRebalancerSetIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolRebalancerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolRebalancerSet struct {
	OldRebalancer common.Address
	NewRebalancer common.Address
	Raw           types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterRebalancerSet(opts *bind.FilterOpts, oldRebalancer []common.Address, newRebalancer []common.Address) (*HybridWithExternalMinterTokenPoolRebalancerSetIterator, error) {

	var oldRebalancerRule []interface{}
	for _, oldRebalancerItem := range oldRebalancer {
		oldRebalancerRule = append(oldRebalancerRule, oldRebalancerItem)
	}
	var newRebalancerRule []interface{}
	for _, newRebalancerItem := range newRebalancer {
		newRebalancerRule = append(newRebalancerRule, newRebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "RebalancerSet", oldRebalancerRule, newRebalancerRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolRebalancerSetIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "RebalancerSet", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRebalancerSet, oldRebalancer []common.Address, newRebalancer []common.Address) (event.Subscription, error) {

	var oldRebalancerRule []interface{}
	for _, oldRebalancerItem := range oldRebalancer {
		oldRebalancerRule = append(oldRebalancerRule, oldRebalancerItem)
	}
	var newRebalancerRule []interface{}
	for _, newRebalancerItem := range newRebalancer {
		newRebalancerRule = append(newRebalancerRule, newRebalancerItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "RebalancerSet", oldRebalancerRule, newRebalancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolRebalancerSet)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RebalancerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseRebalancerSet(log types.Log) (*HybridWithExternalMinterTokenPoolRebalancerSet, error) {
	event := new(HybridWithExternalMinterTokenPoolRebalancerSet)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RebalancerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolReleasedOrMintedIterator struct {
	Event *HybridWithExternalMinterTokenPoolReleasedOrMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolReleasedOrMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolReleasedOrMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolReleasedOrMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolReleasedOrMintedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolReleasedOrMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolReleasedOrMinted struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Recipient           common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolReleasedOrMintedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolReleasedOrMintedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "ReleasedOrMinted", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolReleasedOrMinted)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseReleasedOrMinted(log types.Log) (*HybridWithExternalMinterTokenPoolReleasedOrMinted, error) {
	event := new(HybridWithExternalMinterTokenPoolReleasedOrMinted)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolRemotePoolAddedIterator struct {
	Event *HybridWithExternalMinterTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolRemotePoolAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolRemotePoolAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolRemotePoolAddedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolRemotePoolAdded)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*HybridWithExternalMinterTokenPoolRemotePoolAdded, error) {
	event := new(HybridWithExternalMinterTokenPoolRemotePoolAdded)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator struct {
	Event *HybridWithExternalMinterTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolRemotePoolRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolRemotePoolRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolRemotePoolRemoved)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolRemotePoolRemoved, error) {
	event := new(HybridWithExternalMinterTokenPoolRemotePoolRemoved)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type HybridWithExternalMinterTokenPoolRouterUpdatedIterator struct {
	Event *HybridWithExternalMinterTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *HybridWithExternalMinterTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HybridWithExternalMinterTokenPoolRouterUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(HybridWithExternalMinterTokenPoolRouterUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *HybridWithExternalMinterTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *HybridWithExternalMinterTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type HybridWithExternalMinterTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &HybridWithExternalMinterTokenPoolRouterUpdatedIterator{contract: _HybridWithExternalMinterTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _HybridWithExternalMinterTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(HybridWithExternalMinterTokenPoolRouterUpdated)
				if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*HybridWithExternalMinterTokenPoolRouterUpdated, error) {
	event := new(HybridWithExternalMinterTokenPoolRouterUpdated)
	if err := _HybridWithExternalMinterTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _HybridWithExternalMinterTokenPool.abi.Events["AllowListAdd"].ID:
		return _HybridWithExternalMinterTokenPool.ParseAllowListAdd(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["AllowListRemove"].ID:
		return _HybridWithExternalMinterTokenPool.ParseAllowListRemove(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["ChainAdded"].ID:
		return _HybridWithExternalMinterTokenPool.ParseChainAdded(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["ChainConfigured"].ID:
		return _HybridWithExternalMinterTokenPool.ParseChainConfigured(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["ChainRemoved"].ID:
		return _HybridWithExternalMinterTokenPool.ParseChainRemoved(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["ConfigChanged"].ID:
		return _HybridWithExternalMinterTokenPool.ParseConfigChanged(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["GroupUpdated"].ID:
		return _HybridWithExternalMinterTokenPool.ParseGroupUpdated(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["InboundRateLimitConsumed"].ID:
		return _HybridWithExternalMinterTokenPool.ParseInboundRateLimitConsumed(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["LiquidityAdded"].ID:
		return _HybridWithExternalMinterTokenPool.ParseLiquidityAdded(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["LiquidityMigrated"].ID:
		return _HybridWithExternalMinterTokenPool.ParseLiquidityMigrated(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["LiquidityRemoved"].ID:
		return _HybridWithExternalMinterTokenPool.ParseLiquidityRemoved(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["LockedOrBurned"].ID:
		return _HybridWithExternalMinterTokenPool.ParseLockedOrBurned(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["OutboundRateLimitConsumed"].ID:
		return _HybridWithExternalMinterTokenPool.ParseOutboundRateLimitConsumed(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _HybridWithExternalMinterTokenPool.ParseOwnershipTransferRequested(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _HybridWithExternalMinterTokenPool.ParseOwnershipTransferred(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _HybridWithExternalMinterTokenPool.ParseRateLimitAdminSet(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["RebalancerSet"].ID:
		return _HybridWithExternalMinterTokenPool.ParseRebalancerSet(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["ReleasedOrMinted"].ID:
		return _HybridWithExternalMinterTokenPool.ParseReleasedOrMinted(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _HybridWithExternalMinterTokenPool.ParseRemotePoolAdded(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _HybridWithExternalMinterTokenPool.ParseRemotePoolRemoved(log)
	case _HybridWithExternalMinterTokenPool.abi.Events["RouterUpdated"].ID:
		return _HybridWithExternalMinterTokenPool.ParseRouterUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (HybridWithExternalMinterTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (HybridWithExternalMinterTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (HybridWithExternalMinterTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (HybridWithExternalMinterTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (HybridWithExternalMinterTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (HybridWithExternalMinterTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (HybridWithExternalMinterTokenPoolGroupUpdated) Topic() common.Hash {
	return common.HexToHash("0x1d1eeb97006356bf772500dc592e232d913119a3143e8452f60e5c98b6a29ca1")
}

func (HybridWithExternalMinterTokenPoolInboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0x50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c")
}

func (HybridWithExternalMinterTokenPoolLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (HybridWithExternalMinterTokenPoolLiquidityMigrated) Topic() common.Hash {
	return common.HexToHash("0xbbaa9aea43e3358cd56e894ad9620b8a065abcffab21357fb0702f222480fccc")
}

func (HybridWithExternalMinterTokenPoolLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (HybridWithExternalMinterTokenPoolLockedOrBurned) Topic() common.Hash {
	return common.HexToHash("0xf33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae10")
}

func (HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0xff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944")
}

func (HybridWithExternalMinterTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (HybridWithExternalMinterTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (HybridWithExternalMinterTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (HybridWithExternalMinterTokenPoolRebalancerSet) Topic() common.Hash {
	return common.HexToHash("0x64187bd7b97e66658c91904f3021d7c28de967281d18b1a20742348afdd6a6b3")
}

func (HybridWithExternalMinterTokenPoolReleasedOrMinted) Topic() common.Hash {
	return common.HexToHash("0xfc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0")
}

func (HybridWithExternalMinterTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (HybridWithExternalMinterTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (HybridWithExternalMinterTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (_HybridWithExternalMinterTokenPool *HybridWithExternalMinterTokenPool) Address() common.Address {
	return _HybridWithExternalMinterTokenPool.address
}

type HybridWithExternalMinterTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetGroup(opts *bind.CallOpts, remoteChainSelector uint64) (uint8, error)

	GetLockedTokens(opts *bind.CallOpts) (*big.Int, error)

	GetMinter(opts *bind.CallOpts) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

	IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	UpdateGroups(opts *bind.TransactOpts, groupUpdates []HybridTokenPoolAbstractGroupUpdate) (*types.Transaction, error)

	WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*HybridWithExternalMinterTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*HybridWithExternalMinterTokenPoolAllowListRemove, error)

	FilterChainAdded(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*HybridWithExternalMinterTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*HybridWithExternalMinterTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*HybridWithExternalMinterTokenPoolConfigChanged, error)

	FilterGroupUpdated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterTokenPoolGroupUpdatedIterator, error)

	WatchGroupUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolGroupUpdated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error)

	ParseGroupUpdated(log types.Log) (*HybridWithExternalMinterTokenPoolGroupUpdated, error)

	FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolInboundRateLimitConsumedIterator, error)

	WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseInboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterTokenPoolInboundRateLimitConsumed, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterTokenPoolLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityAdded, rebalancer []common.Address) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityAdded, error)

	FilterLiquidityMigrated(opts *bind.FilterOpts, remoteChainSelector []uint64, group []uint8) (*HybridWithExternalMinterTokenPoolLiquidityMigratedIterator, error)

	WatchLiquidityMigrated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityMigrated, remoteChainSelector []uint64, group []uint8) (event.Subscription, error)

	ParseLiquidityMigrated(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityMigrated, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, rebalancer []common.Address) (*HybridWithExternalMinterTokenPoolLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLiquidityRemoved, rebalancer []common.Address) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolLiquidityRemoved, error)

	FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolLockedOrBurnedIterator, error)

	WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error)

	ParseLockedOrBurned(log types.Log) (*HybridWithExternalMinterTokenPoolLockedOrBurned, error)

	FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator, error)

	WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseOutboundRateLimitConsumed(log types.Log) (*HybridWithExternalMinterTokenPoolOutboundRateLimitConsumed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*HybridWithExternalMinterTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*HybridWithExternalMinterTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*HybridWithExternalMinterTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*HybridWithExternalMinterTokenPoolRateLimitAdminSet, error)

	FilterRebalancerSet(opts *bind.FilterOpts, oldRebalancer []common.Address, newRebalancer []common.Address) (*HybridWithExternalMinterTokenPoolRebalancerSetIterator, error)

	WatchRebalancerSet(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRebalancerSet, oldRebalancer []common.Address, newRebalancer []common.Address) (event.Subscription, error)

	ParseRebalancerSet(log types.Log) (*HybridWithExternalMinterTokenPoolRebalancerSet, error)

	FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolReleasedOrMintedIterator, error)

	WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error)

	ParseReleasedOrMinted(log types.Log) (*HybridWithExternalMinterTokenPoolReleasedOrMinted, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*HybridWithExternalMinterTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*HybridWithExternalMinterTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*HybridWithExternalMinterTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*HybridWithExternalMinterTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *HybridWithExternalMinterTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*HybridWithExternalMinterTokenPoolRouterUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
