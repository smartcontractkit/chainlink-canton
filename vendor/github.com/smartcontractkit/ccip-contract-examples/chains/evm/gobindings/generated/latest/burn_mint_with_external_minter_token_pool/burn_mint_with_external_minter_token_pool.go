// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burn_mint_with_external_minter_token_pool

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

var BurnMintWithExternalMinterTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"localTokenDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"allowlist\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyAllowListUpdates\",\"inputs\":[{\"name\":\"removes\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"adds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyChainUpdates\",\"inputs\":[{\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"chainsToAdd\",\"type\":\"tuple[]\",\"internalType\":\"structTokenPool.ChainUpdate[]\",\"components\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"remoteTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllowList\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllowListEnabled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentInboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentOutboundRateLimiterState\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.TokenBucket\",\"components\":[{\"name\":\"tokens\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdated\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRateLimitAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemotePools\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemoteToken\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRmnProxy\",\"inputs\":[],\"outputs\":[{\"name\":\"rmnProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getToken\",\"inputs\":[],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockOrBurn\",\"inputs\":[{\"name\":\"lockOrBurnIn\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnInV1\",\"components\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"originalSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.LockOrBurnOutV1\",\"components\":[{\"name\":\"destTokenAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destPoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseOrMint\",\"inputs\":[{\"name\":\"releaseOrMintIn\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintInV1\",\"components\":[{\"name\":\"originalSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceDenominatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainTokenData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"components\":[{\"name\":\"destinationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeRemotePool\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfig\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setChainRateLimiterConfigs\",\"inputs\":[{\"name\":\"remoteChainSelectors\",\"type\":\"uint64[]\",\"internalType\":\"uint64[]\"},{\"name\":\"outboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundConfigs\",\"type\":\"tuple[]\",\"internalType\":\"structRateLimiter.Config[]\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateLimitAdmin\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"AllowListAdd\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowListRemove\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"remoteToken\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainConfigured\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChainRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigChanged\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LockedOrBurned\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutboundRateLimitConsumed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateLimitAdminSet\",\"inputs\":[{\"name\":\"rateLimitAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReleasedOrMinted\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemotePoolRemoved\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllowListNotEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BucketOverfilled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotARampOnRouter\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotTransferToSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainAlreadyExists\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ChainNotAllowed\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"CursedByRMN\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DisabledNonZeroRateLimit\",\"inputs\":[{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidDecimalArgs\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidRateLimitRate\",\"inputs\":[{\"name\":\"rateLimiterConfig\",\"type\":\"tuple\",\"internalType\":\"structRateLimiter.Config\",\"components\":[{\"name\":\"isEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"capacity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}]},{\"type\":\"error\",\"name\":\"InvalidRemoteChainDecimals\",\"inputs\":[{\"name\":\"sourcePoolData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidRemotePoolForChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidSourcePoolAddress\",\"inputs\":[{\"name\":\"sourcePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeProposedOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonExistentChain\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnlyCallableByOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OverflowDetected\",\"inputs\":[{\"name\":\"remoteDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"localDecimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"remoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnerCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyAdded\",\"inputs\":[{\"name\":\"remoteChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remotePoolAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"SenderNotAllowed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMaxCapacityExceeded\",\"inputs\":[{\"name\":\"capacity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}]},{\"type\":\"error\",\"name\":\"TokenRateLimitReached\",\"inputs\":[{\"name\":\"minWaitInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]}]",
	Bin: "0x610120806040523461029c57614e2d803803809161001d82856104b1565b8339810160c08282031261029c57610034826104d4565b60208301516001600160a01b0381169391929184820361029c5761005a604082016104e8565b60608201519091906001600160401b03811161029c5781019380601f8601121561029c578451946001600160401b03861161049b578560051b9060208201966100a660405198896104b1565b875260208088019282010192831161029c57602001905b828210610483575050506100df60a06100d8608084016104d4565b92016104d4565b92331561047257600180546001600160a01b0319163317905586158015610461575b8015610450575b61043f5760805260c05260405163313ce56760e01b8152602081600481895afa60009181610403575b506103d8575b5060a052600480546001600160a01b0319166001600160a01b03929092169190911790558051151560e08190526102b5575b506001600160a01b03166101008190526040516321df0da760e01b815290602090829060049082905afa9081156102a95760009161026a575b506001600160a01b031690818103610253576040516147969081610697823960805181818161169a0152818161188c0152818161265d0152818161283001528181612b280152612ba0015260a051818181611aca01528181612ab3015281816135960152613619015260c051818181610c630152818161173501526126f9015260e051818181610bf30152818161177901526123de0152610100518181816101bf0152818161191301526129230152f35b63f902523f60e01b60005260045260245260446000fd5b90506020813d6020116102a1575b81610285602093836104b1565b8101031261029c57610296906104d4565b386101a2565b600080fd5b3d9150610278565b6040513d6000823e3d90fd5b90602090604051906102c783836104b1565b60008252600036813760e051156103c75760005b8251811015610342576001906001600160a01b036102f982866104f6565b51168561030582610538565b610312575b5050016102db565b7f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a1388561030a565b5092905060005b81518110156103bd576001906001600160a01b0361036782856104f6565b511680156103b7578461037982610636565b610387575b50505b01610349565b7f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a1388461037e565b50610381565b5050506020610169565b6335f4a7b360e01b60005260046000fd5b60ff1660ff82168181036103ec5750610137565b6332ad3e0760e11b60005260045260245260446000fd5b9091506020813d602011610437575b8161041f602093836104b1565b8101031261029c57610430906104e8565b9038610131565b3d9150610412565b6342bcdf7f60e11b60005260046000fd5b506001600160a01b03821615610108565b506001600160a01b03841615610101565b639b15e16f60e01b60005260046000fd5b60208091610490846104d4565b8152019101906100bd565b634e487b7160e01b600052604160045260246000fd5b601f909101601f19168101906001600160401b0382119082101761049b57604052565b51906001600160a01b038216820361029c57565b519060ff8216820361029c57565b805182101561050a5760209160051b010190565b634e487b7160e01b600052603260045260246000fd5b805482101561050a5760005260206000200190600090565b600081815260036020526040902054801561062f57600019810181811161061957600254600019810191908211610619578181036105c8575b50505060025480156105b2576000190161058c816002610520565b8154906000199060031b1b19169055600255600052600360205260006040812055600190565b634e487b7160e01b600052603160045260246000fd5b6106016105d96105ea936002610520565b90549060031b1c9283926002610520565b819391549060031b91821b91600019901b19161790565b90556000526003602052604060002055388080610571565b634e487b7160e01b600052601160045260246000fd5b5050600090565b80600052600360205260406000205415600014610690576002546801000000000000000081101561049b576106776105ea8260018594016002556002610520565b9055600254906000526003602052604060002055600190565b5060009056fe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a714612c6857508063181f5a7714612bc457806321df0da714612b55578063240028e814612ad757806324f65ee714612a7b57806339077537146125895780634c5ef0ed1461252657806354c8a4f3146123ac57806362ddd3c4146123295780636d3d1a58146122d757806379ba5097146121ee5780637d54534e146121435780638926f54f146120e05780638da5cb5b1461208e578063962d402014611f1a5780639a4575b9146115f3578063a42a7b8b14611467578063a7cd63b714611395578063acfecf9114611270578063af58d59f14611208578063b0f479a1146111b6578063b794658014611160578063c0d786551461105f578063c4bffe2b14610f11578063c75eea9c14610e4a578063cf7401f314610c87578063dc0bd97114610c18578063e0351e1314610bbd578063e8a1da17146102d8578063f2fde38b146101e85763f36675171461017457600080fd5b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35773ffffffffffffffffffffffffffffffffffffffff610234612db4565b61023c61373b565b163381146102ae57807fffffffffffffffffffffffff0000000000000000000000000000000000000000600054161760005573ffffffffffffffffffffffffffffffffffffffff600154167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278600080a3005b7fdad89dca0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101e3576102e636612f8e565b9190926102f161373b565b6000905b828210610a145750505060009063ffffffff4216907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee184360301925b81811015610a12576000918160051b86013585811215610a0e5786019061012082360312610a0e576040519561036687612e55565b823567ffffffffffffffff81168103610a0a578752602083013567ffffffffffffffff8111610a0a5783019536601f88011215610a0a578635966103a9886131df565b976103b7604051998a612e8d565b8089526020808a019160051b83010190368211610a065760208301905b8282106109d3575050505060208801968752604084013567ffffffffffffffff81116109cf576104079036908601612f3f565b926040890193845261043161041f36606088016130cf565b9560608b0196875260c03691016130cf565b9660808a019788526104438651613b88565b61044d8851613b88565b845151156109a75761046967ffffffffffffffff8b51166143c7565b156109705767ffffffffffffffff8a511681526007602052604081206105a987516fffffffffffffffffffffffffffffffff604082015116906105646fffffffffffffffffffffffffffffffff602083015116915115158360806040516104cf81612e55565b858152602081018c905260408101849052606081018690520152855474ff000000000000000000000000000000000000000091151560a01b919091167fffffffffffffffffffffff0000000000000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff84161773ffffffff0000000000000000000000000000000060808b901b1617178555565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176001830155565b6106cf89516fffffffffffffffffffffffffffffffff6040820151169061068a6fffffffffffffffffffffffffffffffff602083015116915115158360806040516105f381612e55565b858152602081018c9052604081018490526060810186905201526002860180547fffffffffffffffffffffff000000000000000000000000000000000000000000166fffffffffffffffffffffffffffffffff85161773ffffffff0000000000000000000000000000000060808c901b161791151560a01b74ff000000000000000000000000000000000000000016919091179055565b60809190911b7fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff91909116176003830155565b6004865191019080519067ffffffffffffffff8211610943576106f283546132c2565b601f8111610908575b50602090601f831160011461086957610749929185918361085e575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b88518051821015610781579061077b6001926107748367ffffffffffffffff8f5116926132ae565b5190613786565b0161074c565b5050977f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c293919997509561084f67ffffffffffffffff600197969498511692519351915161081b6107e660405196879687526101006020880152610100870190612d55565b9360408601906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60a08401906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b0390a101939193929092610331565b015190508f80610717565b83855281852091907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08416865b8181106108f057509084600195949392106108b9575b505050811b01905561074c565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e80806108ac565b92936020600181928786015181550195019301610896565b6109339084865260208620601f850160051c81019160208610610939575b601f0160051c01906134c9565b8e6106fb565b9091508190610926565b6024847f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b60249067ffffffffffffffff8b51167f1d5ad3c5000000000000000000000000000000000000000000000000000000008252600452fd5b807f8579befe0000000000000000000000000000000000000000000000000000000060049252fd5b8680fd5b813567ffffffffffffffff8111610a02576020916109f78392833691890101612f3f565b8152019101906103d4565b8a80fd5b8880fd5b8580fd5b8380fd5b005b909267ffffffffffffffff610a35610a3086868699979961325f565b61318d565b1692610a40846140fb565b15610b8f57836000526007602052610a5e6005604060002001613f02565b9260005b8451811015610a9a57600190866000526007602052610a936005604060002001610a8c83896132ae565b5190614226565b5001610a62565b5093909491959250806000526007602052600560406000206000815560006001820155600060028201556000600382015560048101610ad981546132c2565b9081610b4c575b5050018054906000815581610b2b575b5050907f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599166020600193604051908152a10190919392936102f5565b6000526020600020908101905b81811015610af05760008155600101610b38565b81601f60009311600114610b645750555b8880610ae0565b81835260208320610b7f91601f01861c8101906001016134c9565b8082528160208120915555610b5d565b837f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760206040517f000000000000000000000000000000000000000000000000000000000000000015158152f35b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e35760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357610cbe612dd7565b60607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc3601126101e357604051610cf481612e71565b60243580151581036101e35781526044356fffffffffffffffffffffffffffffffff811681036101e35760208201526064356fffffffffffffffffffffffffffffffff811681036101e357604082015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7c3601126101e35760405190610d7b82612e71565b60843580151581036101e357825260a4356fffffffffffffffffffffffffffffffff811681036101e357602083015260c4356fffffffffffffffffffffffffffffffff811681036101e357604083015273ffffffffffffffffffffffffffffffffffffffff6009541633141580610e28575b610dfa57610a12926139c6565b7f8e4a23d6000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b5073ffffffffffffffffffffffffffffffffffffffff60015416331415610ded565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35767ffffffffffffffff610e8a612dd7565b610e92613416565b50166000526007602052610f0d610eb4610eaf6040600020613441565b613b03565b6040519182918291909160806fffffffffffffffffffffffffffffffff8160a084019582815116855263ffffffff6020820151166020860152604081015115156040860152826060820151166060860152015116910152565b0390f35b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e3576040516005548082528160208101600560005260206000209260005b818110611046575050610f7192500382612e8d565b805190610f96610f80836131df565b92610f8e6040519485612e8d565b8084526131df565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe060208401920136833760005b8151811015610ff6578067ffffffffffffffff610fe3600193856132ae565b5116610fef82876132ae565b5201610fc4565b5050906040519182916020830190602084525180915260408301919060005b818110611023575050500390f35b825167ffffffffffffffff16845285945060209384019390920191600101611015565b8454835260019485019486945060209093019201610f5c565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357611096612db4565b61109e61373b565b73ffffffffffffffffffffffffffffffffffffffff811690811561113657600480547fffffffffffffffffffffffff000000000000000000000000000000000000000081169390931790556040805173ffffffffffffffffffffffffffffffffffffffff93841681529190921660208201527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f16849190a1005b7f8579befe0000000000000000000000000000000000000000000000000000000060005260046000fd5b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357610f0d6111a261119d612dd7565b6134a7565b604051918291602083526020830190612d55565b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602073ffffffffffffffffffffffffffffffffffffffff60045416604051908152f35b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35767ffffffffffffffff611248612dd7565b611250613416565b50166000526007602052610f0d610eb4610eaf6002604060002001613441565b346101e35767ffffffffffffffff61128736612ffe565b92909161129261373b565b16906112ab826000526006602052604060002054151590565b15611367578160005260076020526112dc60056040600020016112cf368685612f08565b6020815191012090614226565b15611320577f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76919261131b6040519283926020845260208401916133d7565b0390a2005b611363906040519384937f74f23c7c00000000000000000000000000000000000000000000000000000000855260048501526040602485015260448401916133d7565b0390fd5b507f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760405160025490818152602081018092600260005260206000209060005b81811061145157505050816113f8910382612e8d565b6040519182916020830190602084525180915260408301919060005b818110611422575050500390f35b825173ffffffffffffffffffffffffffffffffffffffff16845285945060209384019390920191600101611414565b82548452602090930192600192830192016113e2565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35767ffffffffffffffff6114a7612dd7565b1660005260076020526114c06005604060002001613f02565b8051907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06115066114f0846131df565b936114fe6040519586612e8d565b8085526131df565b0160005b8181106115e257505060005b815181101561155e578061152c600192846132ae565b5160005260086020526115426040600020613315565b61154c82866132ae565b5261155781856132ae565b5001611516565b826040518091602082016020835281518091526040830190602060408260051b8601019301916000905b82821061159757505050500390f35b919360206115d2827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc060019597998495030186528851612d55565b9601920192018594939192611588565b80606060208093870101520161150a565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760043567ffffffffffffffff81116101e35760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126101e3576060602060405161167081612e39565b8281520152608481016116828161316c565b73ffffffffffffffffffffffffffffffffffffffff807f000000000000000000000000000000000000000000000000000000000000000016911603611ece57506024810177ffffffffffffffff000000000000000000000000000000006116e88261318d565b60801b16604051907f2cbc26bb000000000000000000000000000000000000000000000000000000008252600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa908115611d5e57600091611eaf575b50611e85576117776044830161316c565b7f0000000000000000000000000000000000000000000000000000000000000000611e2f575b5067ffffffffffffffff6117b08261318d565b166117c8816000526006602052604060002054151590565b15611e0257602073ffffffffffffffffffffffffffffffffffffffff60045416916024604051809481937fa8d87a3b00000000000000000000000000000000000000000000000000000000835260048301525afa908115611d5e57600091611d98575b5073ffffffffffffffffffffffffffffffffffffffff163303611d6a5767ffffffffffffffff91606461185d8361318d565b910135928391168060005260076020526118b4604060002073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016958691614476565b6040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018490527fff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da81789449190a273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169281158015611cc4575b15611c40576040517f095ea7b3000000000000000000000000000000000000000000000000000000006020820190815273ffffffffffffffffffffffffffffffffffffffff86166024830152604480830185905282529490611a16906119a4606482612e8d565b6000806040988951936119b78b86612e8d565b602085527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646020860152519082885af13d15611c38573d906119f882612ece565b91611a058a519384612e8d565b82523d6000602084013e5b856146bd565b805180611b97575b5050602060009160248751809481937f42966c680000000000000000000000000000000000000000000000000000000083528860048401525af18015611b8c5793610f0d937ff33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae106060611ac29561119d95611b2c9a99611b5d575b5067ffffffffffffffff611aab8661318d565b1693895191825233602083015289820152a261318d565b9180519060ff7f000000000000000000000000000000000000000000000000000000000000000016602083015260208252611afd8183612e8d565b805193611b0985612e39565b845260208401918252805194859460208652518260208701526060860190612d55565b9151907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08584030190850152612d55565b611b7e9060203d602011611b85575b611b768183612e8d565b810190613723565b508a611a98565b503d611b6c565b85513d6000823e3d90fd5b90602080611ba9938301019101613723565b15611bb5578580611a1e565b608485517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152fd5b606090611a10565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f5361666545524332303a20617070726f76652066726f6d206e6f6e2d7a65726f60448201527f20746f206e6f6e2d7a65726f20616c6c6f77616e6365000000000000000000006064820152fd5b506040517fdd62ed3e00000000000000000000000000000000000000000000000000000000815230600482015273ffffffffffffffffffffffffffffffffffffffff85166024820152602081604481855afa908115611d5e57600091611d2c575b501561193d565b90506020813d602011611d56575b81611d4760209383612e8d565b810103126101e3575185611d25565b3d9150611d3a565b6040513d6000823e3d90fd5b7f728fe07b000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b6020813d602011611dfa575b81611db160209383612e8d565b81010312611df657519073ffffffffffffffffffffffffffffffffffffffff82168203611df3575073ffffffffffffffffffffffffffffffffffffffff61182b565b80fd5b5080fd5b3d9150611da4565b7fa9902c7e0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff1680600052600360205260406000205461179d577fd0d259760000000000000000000000000000000000000000000000000000000060005260045260246000fd5b7f53ad11d80000000000000000000000000000000000000000000000000000000060005260046000fd5b611ec8915060203d602011611b8557611b768183612e8d565b83611766565b611eec73ffffffffffffffffffffffffffffffffffffffff9161316c565b7f961c9a4f000000000000000000000000000000000000000000000000000000006000521660045260246000fd5b346101e35760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760043567ffffffffffffffff81116101e357611f69903690600401612f5d565b9060243567ffffffffffffffff81116101e357611f8a903690600401613081565b9060443567ffffffffffffffff81116101e357611fab903690600401613081565b73ffffffffffffffffffffffffffffffffffffffff600954163314158061206c575b610dfa57838614801590612062575b6120385760005b868110611fec57005b80612032612000610a306001948b8b61325f565b61200b83898961329e565b61202c61202461201c86898b61329e565b9236906130cf565b9136906130cf565b916139c6565b01611fe3565b7f568efce20000000000000000000000000000000000000000000000000000000060005260046000fd5b5080861415611fdc565b5073ffffffffffffffffffffffffffffffffffffffff60015416331415611fcd565b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602073ffffffffffffffffffffffffffffffffffffffff60015416604051908152f35b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602061213967ffffffffffffffff612125612dd7565b166000526006602052604060002054151590565b6040519015158152f35b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e3577f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174602073ffffffffffffffffffffffffffffffffffffffff6121b2612db4565b6121ba61373b565b16807fffffffffffffffffffffffff00000000000000000000000000000000000000006009541617600955604051908152a1005b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760005473ffffffffffffffffffffffffffffffffffffffff811633036122ad577fffffffffffffffffffffffff00000000000000000000000000000000000000006001549133828416176001551660005573ffffffffffffffffffffffffffffffffffffffff3391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3005b7f02b543c60000000000000000000000000000000000000000000000000000000060005260046000fd5b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602073ffffffffffffffffffffffffffffffffffffffff60095416604051908152f35b346101e35761233736612ffe565b61234292919261373b565b67ffffffffffffffff8216612364816000526006602052604060002054151590565b1561237f5750610a1292612379913691612f08565b90613786565b7f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b346101e3576123d46123dc6123c036612f8e565b94916123cd93919361373b565b36916131f7565b9236916131f7565b7f0000000000000000000000000000000000000000000000000000000000000000156124fc5760005b8251811015612478578073ffffffffffffffffffffffffffffffffffffffff612430600193866132ae565b511661243b81613f65565b612447575b5001612405565b60207f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756691604051908152a184612440565b5060005b8151811015610a12578073ffffffffffffffffffffffffffffffffffffffff6124a7600193856132ae565b511680156124f6576124b881614367565b6124c5575b505b0161247c565b60207f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d891604051908152a1836124bd565b506124bf565b7f35f4a7b30000000000000000000000000000000000000000000000000000000060005260046000fd5b346101e35760407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35761255d612dd7565b60243567ffffffffffffffff81116101e357602091612583612139923690600401612f3f565b906131a2565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e35760043567ffffffffffffffff81116101e35780600401906101007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126101e357600060405161260a81612dee565b5261263761262d61262861262160c485018661311b565b3691612f08565b613522565b6064830135613616565b90608481016126458161316c565b73ffffffffffffffffffffffffffffffffffffffff807f000000000000000000000000000000000000000000000000000000000000000016911603611ece5750602481019277ffffffffffffffff000000000000000000000000000000006126ac8561318d565b60801b16604051907f2cbc26bb000000000000000000000000000000000000000000000000000000008252600482015260208160248173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165afa908115611d5e57600091612a5c575b50611e855767ffffffffffffffff6127418561318d565b16612759816000526006602052604060002054151590565b15611e0257602073ffffffffffffffffffffffffffffffffffffffff60045416916044604051809481937f83826b2b00000000000000000000000000000000000000000000000000000000835260048301523360248301525afa908115611d5e57600091612a3d575b5015611d6a576127d18461318d565b906127e760a4840192612583612621858561311b565b156129f65750506044829167ffffffffffffffff6128048661318d565b16806000526007602052612858600260406000200173ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016958691614476565b6040805173ffffffffffffffffffffffffffffffffffffffff86168152602081018790527f50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c9190a201926129086020846128b18761316c565b60405193849283927f40c10f19000000000000000000000000000000000000000000000000000000008452600484016020909392919373ffffffffffffffffffffffffffffffffffffffff60408201951681520152565b0381600073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000165af18015611d5e5760209573ffffffffffffffffffffffffffffffffffffffff6129a96129a37ffc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc09660809667ffffffffffffffff966129d9575b5061318d565b9261316c565b60405196875233898801521660408601528560608601521692a2806040516129d081612dee565b52604051908152f35b6129ef908c3d8e11611b8557611b768183612e8d565b508b61299d565b612a00925061311b565b6113636040519283927f24eb47e50000000000000000000000000000000000000000000000000000000084526020600485015260248401916133d7565b612a56915060203d602011611b8557611b768183612e8d565b856127c2565b612a75915060203d602011611b8557611b768183612e8d565b8561272a565b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602060405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e3576020612b10612db4565b73ffffffffffffffffffffffffffffffffffffffff807f0000000000000000000000000000000000000000000000000000000000000000169116146040519015158152f35b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b346101e35760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357610f0d604051612c04606082612e8d565b602981527f4275726e4d696e745769746845787465726e616c4d696e746572546f6b656e5060208201527f6f6f6c20312e362e3000000000000000000000000000000000000000000000006040820152604051918291602083526020830190612d55565b346101e35760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101e357600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036101e357817faff2afbf0000000000000000000000000000000000000000000000000000000060209314908115612d2b575b8115612d01575b5015158152f35b7f01ffc9a70000000000000000000000000000000000000000000000000000000091501483612cfa565b7f0e64dd290000000000000000000000000000000000000000000000000000000081149150612cf3565b919082519283825260005b848110612d9f5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8460006020809697860101520116010190565b80602080928401015182828601015201612d60565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036101e357565b6004359067ffffffffffffffff821682036101e357565b6020810190811067ffffffffffffffff821117612e0a57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040810190811067ffffffffffffffff821117612e0a57604052565b60a0810190811067ffffffffffffffff821117612e0a57604052565b6060810190811067ffffffffffffffff821117612e0a57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117612e0a57604052565b67ffffffffffffffff8111612e0a57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b929192612f1482612ece565b91612f226040519384612e8d565b8294818452818301116101e3578281602093846000960137010152565b9080601f830112156101e357816020612f5a93359101612f08565b90565b9181601f840112156101e35782359167ffffffffffffffff83116101e3576020808501948460051b0101116101e357565b60407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126101e35760043567ffffffffffffffff81116101e35781612fd791600401612f5d565b929092916024359067ffffffffffffffff82116101e357612ffa91600401612f5d565b9091565b60407ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8201126101e35760043567ffffffffffffffff811681036101e3579160243567ffffffffffffffff81116101e357826023820112156101e35780600401359267ffffffffffffffff84116101e357602484830101116101e3576024019190565b9181601f840112156101e35782359167ffffffffffffffff83116101e357602080850194606085020101116101e357565b35906fffffffffffffffffffffffffffffffff821682036101e357565b91908260609103126101e3576040516130e781612e71565b809280359081151582036101e3576040613116918193855261310b602082016130b2565b6020860152016130b2565b910152565b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101e3570180359067ffffffffffffffff82116101e3576020019181360383136101e357565b3573ffffffffffffffffffffffffffffffffffffffff811681036101e35790565b3567ffffffffffffffff811681036101e35790565b9067ffffffffffffffff612f5a92166000526007602052600560406000200190602081519101209060019160005201602052604060002054151590565b67ffffffffffffffff8111612e0a5760051b60200190565b9291613202826131df565b936132106040519586612e8d565b602085848152019260051b81019182116101e357915b81831061323257505050565b823573ffffffffffffffffffffffffffffffffffffffff811681036101e357815260209283019201613226565b919081101561326f5760051b0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b919081101561326f576060020190565b805182101561326f5760209160051b010190565b90600182811c9216801561330b575b60208310146132dc57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916132d1565b9060405191826000825492613329846132c2565b80845293600181169081156133975750600114613350575b5061334e92500383612e8d565b565b90506000929192526020600020906000915b81831061337b57505090602061334e9282010138613341565b6020919350806001915483858901015201910190918492613362565b6020935061334e9592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b82010138613341565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b6040519061342382612e55565b60006080838281528260208201528260408201528260608201520152565b9060405161344e81612e55565b60806001829460ff81546fffffffffffffffffffffffffffffffff8116865263ffffffff81861c16602087015260a01c161515604085015201546fffffffffffffffffffffffffffffffff81166060840152811c910152565b67ffffffffffffffff166000526007602052612f5a6004604060002001613315565b8181106134d4575050565b600081556001016134c9565b818102929181159184041417156134f357565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80518015613592576020036135545780516020828101918301839003126101e357519060ff8211613554575060ff1690565b611363906040519182917f953576f7000000000000000000000000000000000000000000000000000000008352602060048401526024830190612d55565b50507f000000000000000000000000000000000000000000000000000000000000000090565b9060ff8091169116039060ff82116134f357565b60ff16604d81116134f357600a0a90565b81156135e7570490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b907f00000000000000000000000000000000000000000000000000000000000000009060ff82169060ff81169282841461371c578284116136f2579061365b916135b8565b91604d60ff84161180156136b9575b6136835750509061367d612f5a926135cc565b906134e0565b9091507fa9cb113d0000000000000000000000000000000000000000000000000000000060005260045260245260445260646000fd5b506136c3836135cc565b80156135e7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04841161366a565b6136fb916135b8565b91604d60ff84161161368357505090613716612f5a926135cc565b906135dd565b5050505090565b908160209103126101e3575180151581036101e35790565b73ffffffffffffffffffffffffffffffffffffffff60015416330361375c57565b7f2b5c74de0000000000000000000000000000000000000000000000000000000060005260046000fd5b908051156111365767ffffffffffffffff815160208301209216918260005260076020526137bb816005604060002001614421565b156139825760005260086020526040600020815167ffffffffffffffff8111612e0a576137e882546132c2565b601f8111613950575b506020601f821160011461388a5791613864827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea959361387a9560009161387f575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055604051918291602083526020830190612d55565b0390a2565b905084015138613833565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169083600052806000209160005b81811061393857509261387a9492600192827f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea989610613901575b5050811b0190556111a2565b8501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c1916905538806138f5565b9192602060018192868a0151815501940192016138ba565b61397c90836000526020600020601f840160051c8101916020851061093957601f0160051c01906134c9565b386137f1565b50906113636040519283927f393b8ad20000000000000000000000000000000000000000000000000000000084526004840152604060248401526044830190612d55565b67ffffffffffffffff166000818152600660205260409020549092919015613ac85791613ac560e092613a9185613a1d7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b97613b88565b846000526007602052613a34816040600020613ccf565b613a3d83613b88565b846000526007602052613a57836002604060002001613ccf565b60405194855260208501906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b60808301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565ba1565b827f1e670e4b0000000000000000000000000000000000000000000000000000000060005260045260246000fd5b919082039182116134f357565b613b0b613416565b506fffffffffffffffffffffffffffffffff6060820151166fffffffffffffffffffffffffffffffff8083511691613b686020850193613b62613b5563ffffffff87511642613af6565b85608089015116906134e0565b9061435a565b80821015613b8157505b16825263ffffffff4216905290565b9050613b72565b805115613c28576fffffffffffffffffffffffffffffffff6040820151166fffffffffffffffffffffffffffffffff60208301511610613bc55750565b606490613c26604051917f8020d12400000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565bfd5b6fffffffffffffffffffffffffffffffff60408201511615801590613cb0575b613c4f5750565b606490613c26604051917fd68af9cc00000000000000000000000000000000000000000000000000000000835260048301906fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b506fffffffffffffffffffffffffffffffff6020820151161515613c48565b7f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1991613e086060928054613d0c63ffffffff8260801c1642613af6565b9081613e47575b50506fffffffffffffffffffffffffffffffff6001816020860151169282815416808510600014613e3f57508280855b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416178155613dbc8651151582907fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff0000000000000000000000000000000000000000835492151560a01b169116179055565b60408601517fffffffffffffffffffffffffffffffff0000000000000000000000000000000060809190911b16939092166fffffffffffffffffffffffffffffffff1692909217910155565b613ac560405180926fffffffffffffffffffffffffffffffff60408092805115158552826020820151166020860152015116910152565b838091613d43565b6fffffffffffffffffffffffffffffffff91613e7c839283613e756001880154948286169560801c906134e0565b911661435a565b80821015613efb57505b83547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9290911692909216167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116174260801b73ffffffff00000000000000000000000000000000161781553880613d13565b9050613e86565b906040519182815491828252602082019060005260206000209260005b818110613f3457505061334e92500383612e8d565b8454835260019485019487945060209093019201613f1f565b805482101561326f5760005260206000200190600090565b60008181526003602052604090205480156140f4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116134f357600254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116134f357818103614085575b5050506002548015614056577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01614013816002613f4d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600255600052600360205260006040812055600190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b6140dc6140966140a7936002613f4d565b90549060031b1c9283926002613f4d565b81939154907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9060031b92831b921b19161790565b90556000526003602052604060002055388080613fda565b5050600090565b60008181526006602052604090205480156140f4577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116134f357600554907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116134f3578181036141ec575b5050506005548015614056577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff016141a9816005613f4d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b19169055600555600052600660205260006040812055600190565b61420e6141fd6140a7936005613f4d565b90549060031b1c9283926005613f4d565b90556000526006602052604060002055388080614170565b9060018201918160005282602052604060002054801515600014614351577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181116134f3578254907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82019182116134f35781810361431a575b50505080548015614056577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01906142db8282613f4d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82549160031b1b191690555560005260205260006040812055600190565b61433a61432a6140a79386613f4d565b90549060031b1c92839286613f4d565b9055600052836020526040600020553880806142a3565b50505050600090565b919082018092116134f357565b806000526003602052604060002054156000146143c15760025468010000000000000000811015612e0a576143a86140a78260018594016002556002613f4d565b9055600254906000526003602052604060002055600190565b50600090565b806000526006602052604060002054156000146143c15760055468010000000000000000811015612e0a576144086140a78260018594016005556005613f4d565b9055600554906000526006602052604060002055600190565b60008281526001820160205260409020546140f45780549068010000000000000000821015612e0a578261445f6140a7846001809601855584613f4d565b905580549260005201602052604060002055600190565b9182549060ff8260a01c161580156146b5575b6146af576fffffffffffffffffffffffffffffffff821691600185019081546144ce63ffffffff6fffffffffffffffffffffffffffffffff83169360801c1642613af6565b9081614611575b50508481106145c5575083831061452f5750506145046fffffffffffffffffffffffffffffffff928392613af6565b16167fffffffffffffffffffffffffffffffff00000000000000000000000000000000825416179055565b5460801c9161453e8185613af6565b927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101908082116134f35761458c6145919273ffffffffffffffffffffffffffffffffffffffff9661435a565b6135dd565b7fd0c8d23a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b828573ffffffffffffffffffffffffffffffffffffffff927f1a76572a000000000000000000000000000000000000000000000000000000006000526004526024521660445260646000fd5b8286929396116146855761462c92613b629160801c906134e0565b808410156146805750825b85547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff164260801b73ffffffff00000000000000000000000000000000161786559238806144d5565b614637565b7f9725942a0000000000000000000000000000000000000000000000000000000060005260046000fd5b50505050565b508215614489565b9192901561473857508151156146d1575090565b3b156146da5790565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152fd5b82519091501561474b5750805190602001fd5b611363906040519182917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401526024830190612d5556fea164736f6c634300081a000a",
}

var BurnMintWithExternalMinterTokenPoolABI = BurnMintWithExternalMinterTokenPoolMetaData.ABI

var BurnMintWithExternalMinterTokenPoolBin = BurnMintWithExternalMinterTokenPoolMetaData.Bin

func DeployBurnMintWithExternalMinterTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, minter common.Address, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *types.Transaction, *BurnMintWithExternalMinterTokenPool, error) {
	parsed, err := BurnMintWithExternalMinterTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintWithExternalMinterTokenPoolBin), backend, minter, token, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnMintWithExternalMinterTokenPool{address: address, abi: *parsed, BurnMintWithExternalMinterTokenPoolCaller: BurnMintWithExternalMinterTokenPoolCaller{contract: contract}, BurnMintWithExternalMinterTokenPoolTransactor: BurnMintWithExternalMinterTokenPoolTransactor{contract: contract}, BurnMintWithExternalMinterTokenPoolFilterer: BurnMintWithExternalMinterTokenPoolFilterer{contract: contract}}, nil
}

type BurnMintWithExternalMinterTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnMintWithExternalMinterTokenPoolCaller
	BurnMintWithExternalMinterTokenPoolTransactor
	BurnMintWithExternalMinterTokenPoolFilterer
}

type BurnMintWithExternalMinterTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnMintWithExternalMinterTokenPoolSession struct {
	Contract     *BurnMintWithExternalMinterTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintWithExternalMinterTokenPoolCallerSession struct {
	Contract *BurnMintWithExternalMinterTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnMintWithExternalMinterTokenPoolTransactorSession struct {
	Contract     *BurnMintWithExternalMinterTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnMintWithExternalMinterTokenPoolRaw struct {
	Contract *BurnMintWithExternalMinterTokenPool
}

type BurnMintWithExternalMinterTokenPoolCallerRaw struct {
	Contract *BurnMintWithExternalMinterTokenPoolCaller
}

type BurnMintWithExternalMinterTokenPoolTransactorRaw struct {
	Contract *BurnMintWithExternalMinterTokenPoolTransactor
}

func NewBurnMintWithExternalMinterTokenPool(address common.Address, backend bind.ContractBackend) (*BurnMintWithExternalMinterTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintWithExternalMinterTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintWithExternalMinterTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPool{address: address, abi: abi, BurnMintWithExternalMinterTokenPoolCaller: BurnMintWithExternalMinterTokenPoolCaller{contract: contract}, BurnMintWithExternalMinterTokenPoolTransactor: BurnMintWithExternalMinterTokenPoolTransactor{contract: contract}, BurnMintWithExternalMinterTokenPoolFilterer: BurnMintWithExternalMinterTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnMintWithExternalMinterTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnMintWithExternalMinterTokenPoolCaller, error) {
	contract, err := bindBurnMintWithExternalMinterTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolCaller{contract: contract}, nil
}

func NewBurnMintWithExternalMinterTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintWithExternalMinterTokenPoolTransactor, error) {
	contract, err := bindBurnMintWithExternalMinterTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolTransactor{contract: contract}, nil
}

func NewBurnMintWithExternalMinterTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintWithExternalMinterTokenPoolFilterer, error) {
	contract, err := bindBurnMintWithExternalMinterTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolFilterer{contract: contract}, nil
}

func bindBurnMintWithExternalMinterTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintWithExternalMinterTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintWithExternalMinterTokenPool.Contract.BurnMintWithExternalMinterTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.BurnMintWithExternalMinterTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.BurnMintWithExternalMinterTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintWithExternalMinterTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetAllowList(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetAllowList(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetAllowListEnabled(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetAllowListEnabled(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetMinter() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetMinter(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetMinter() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetMinter(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRateLimitAdmin(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRateLimitAdmin(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRemotePools(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRemotePools(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRemoteToken(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRemoteToken(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRmnProxy(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRmnProxy(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRouter(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetRouter(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetSupportedChains(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetSupportedChains(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetToken(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetToken(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetTokenDecimals(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.GetTokenDecimals(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsRemotePool(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsRemotePool(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsSupportedChain(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsSupportedChain(&_BurnMintWithExternalMinterTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsSupportedToken(&_BurnMintWithExternalMinterTokenPool.CallOpts, token)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.IsSupportedToken(&_BurnMintWithExternalMinterTokenPool.CallOpts, token)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) Owner() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.Owner(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.Owner(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SupportsInterface(&_BurnMintWithExternalMinterTokenPool.CallOpts, interfaceId)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SupportsInterface(&_BurnMintWithExternalMinterTokenPool.CallOpts, interfaceId)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintWithExternalMinterTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.TypeAndVersion(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.TypeAndVersion(&_BurnMintWithExternalMinterTokenPool.CallOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.AcceptOwnership(&_BurnMintWithExternalMinterTokenPool.TransactOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.AcceptOwnership(&_BurnMintWithExternalMinterTokenPool.TransactOpts)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.AddRemotePool(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.AddRemotePool(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintWithExternalMinterTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintWithExternalMinterTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ApplyChainUpdates(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ApplyChainUpdates(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.LockOrBurn(&_BurnMintWithExternalMinterTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.LockOrBurn(&_BurnMintWithExternalMinterTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ReleaseOrMint(&_BurnMintWithExternalMinterTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.ReleaseOrMint(&_BurnMintWithExternalMinterTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.RemoveRemotePool(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.RemoveRemotePool(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "setChainRateLimiterConfigs", remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfigs(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) SetChainRateLimiterConfigs(remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetChainRateLimiterConfigs(&_BurnMintWithExternalMinterTokenPool.TransactOpts, remoteChainSelectors, outboundConfigs, inboundConfigs)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetRateLimitAdmin(&_BurnMintWithExternalMinterTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetRateLimitAdmin(&_BurnMintWithExternalMinterTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetRouter(&_BurnMintWithExternalMinterTokenPool.TransactOpts, newRouter)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.SetRouter(&_BurnMintWithExternalMinterTokenPool.TransactOpts, newRouter)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.TransferOwnership(&_BurnMintWithExternalMinterTokenPool.TransactOpts, to)
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintWithExternalMinterTokenPool.Contract.TransferOwnership(&_BurnMintWithExternalMinterTokenPool.TransactOpts, to)
}

type BurnMintWithExternalMinterTokenPoolAllowListAddIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolAllowListAdd)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolAllowListAdd)
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

func (it *BurnMintWithExternalMinterTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolAllowListAddIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolAllowListAdd)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnMintWithExternalMinterTokenPoolAllowListAdd, error) {
	event := new(BurnMintWithExternalMinterTokenPoolAllowListAdd)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolAllowListRemove)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolAllowListRemove)
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

func (it *BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolAllowListRemove)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnMintWithExternalMinterTokenPoolAllowListRemove, error) {
	event := new(BurnMintWithExternalMinterTokenPoolAllowListRemove)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolChainAddedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolChainAdded)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolChainAdded)
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

func (it *BurnMintWithExternalMinterTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolChainAddedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolChainAdded)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainAdded, error) {
	event := new(BurnMintWithExternalMinterTokenPoolChainAdded)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolChainConfiguredIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolChainConfigured)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolChainConfigured)
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

func (it *BurnMintWithExternalMinterTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolChainConfiguredIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolChainConfigured)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainConfigured, error) {
	event := new(BurnMintWithExternalMinterTokenPoolChainConfigured)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolChainRemovedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolChainRemoved)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolChainRemoved)
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

func (it *BurnMintWithExternalMinterTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolChainRemovedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolChainRemoved)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainRemoved, error) {
	event := new(BurnMintWithExternalMinterTokenPoolChainRemoved)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolConfigChangedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolConfigChanged)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolConfigChanged)
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

func (it *BurnMintWithExternalMinterTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolConfigChangedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolConfigChanged)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnMintWithExternalMinterTokenPoolConfigChanged, error) {
	event := new(BurnMintWithExternalMinterTokenPoolConfigChanged)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed)
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

func (it *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "InboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "InboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseInboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed, error) {
	event := new(BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "InboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolLockedOrBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolLockedOrBurned)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolLockedOrBurned)
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

func (it *BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolLockedOrBurned struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "LockedOrBurned", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "LockedOrBurned", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolLockedOrBurned)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseLockedOrBurned(log types.Log) (*BurnMintWithExternalMinterTokenPoolLockedOrBurned, error) {
	event := new(BurnMintWithExternalMinterTokenPoolLockedOrBurned)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "LockedOrBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed)
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

func (it *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed struct {
	RemoteChainSelector uint64
	Token               common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "OutboundRateLimitConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "OutboundRateLimitConsumed", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseOutboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed, error) {
	event := new(BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OutboundRateLimitConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested)
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

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolOwnershipTransferred)
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

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolOwnershipTransferred)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferred, error) {
	event := new(BurnMintWithExternalMinterTokenPoolOwnershipTransferred)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolRateLimitAdminSet)
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

func (it *BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolRateLimitAdminSet)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnMintWithExternalMinterTokenPoolRateLimitAdminSet, error) {
	event := new(BurnMintWithExternalMinterTokenPoolRateLimitAdminSet)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolReleasedOrMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolReleasedOrMinted)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolReleasedOrMinted)
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

func (it *BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolReleasedOrMinted struct {
	RemoteChainSelector uint64
	Token               common.Address
	Sender              common.Address
	Recipient           common.Address
	Amount              *big.Int
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "ReleasedOrMinted", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "ReleasedOrMinted", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolReleasedOrMinted)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseReleasedOrMinted(log types.Log) (*BurnMintWithExternalMinterTokenPoolReleasedOrMinted, error) {
	event := new(BurnMintWithExternalMinterTokenPoolReleasedOrMinted)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "ReleasedOrMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolRemotePoolAdded)
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

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolRemotePoolAdded)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnMintWithExternalMinterTokenPoolRemotePoolAdded, error) {
	event := new(BurnMintWithExternalMinterTokenPoolRemotePoolAdded)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolRemotePoolRemoved)
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

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolRemotePoolRemoved)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnMintWithExternalMinterTokenPoolRemotePoolRemoved, error) {
	event := new(BurnMintWithExternalMinterTokenPoolRemotePoolRemoved)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator struct {
	Event *BurnMintWithExternalMinterTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintWithExternalMinterTokenPoolRouterUpdated)
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
		it.Event = new(BurnMintWithExternalMinterTokenPoolRouterUpdated)
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

func (it *BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintWithExternalMinterTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator{contract: _BurnMintWithExternalMinterTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintWithExternalMinterTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintWithExternalMinterTokenPoolRouterUpdated)
				if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnMintWithExternalMinterTokenPoolRouterUpdated, error) {
	event := new(BurnMintWithExternalMinterTokenPoolRouterUpdated)
	if err := _BurnMintWithExternalMinterTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintWithExternalMinterTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseAllowListAdd(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseAllowListRemove(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseChainAdded(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseChainConfigured(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseChainRemoved(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseConfigChanged(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["InboundRateLimitConsumed"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseInboundRateLimitConsumed(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["LockedOrBurned"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseLockedOrBurned(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["OutboundRateLimitConsumed"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseOutboundRateLimitConsumed(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseOwnershipTransferred(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseRateLimitAdminSet(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["ReleasedOrMinted"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseReleasedOrMinted(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseRemotePoolAdded(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseRemotePoolRemoved(log)
	case _BurnMintWithExternalMinterTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnMintWithExternalMinterTokenPool.ParseRouterUpdated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintWithExternalMinterTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnMintWithExternalMinterTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnMintWithExternalMinterTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnMintWithExternalMinterTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnMintWithExternalMinterTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnMintWithExternalMinterTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0x50f6fbee3ceedce6b7fd7eaef18244487867e6718aec7208187efb6b7908c14c")
}

func (BurnMintWithExternalMinterTokenPoolLockedOrBurned) Topic() common.Hash {
	return common.HexToHash("0xf33bc26b4413b0e7f19f1ea739fdf99098c0061f1f87d954b11f5293fad9ae10")
}

func (BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed) Topic() common.Hash {
	return common.HexToHash("0xff0133389f9bb82d5b9385826160eaf2328039f6fa950eeb8cf0836da8178944")
}

func (BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintWithExternalMinterTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintWithExternalMinterTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnMintWithExternalMinterTokenPoolReleasedOrMinted) Topic() common.Hash {
	return common.HexToHash("0xfc5e3a5bddc11d92c2dc20fae6f7d5eb989f056be35239f7de7e86150609abc0")
}

func (BurnMintWithExternalMinterTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnMintWithExternalMinterTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnMintWithExternalMinterTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (_BurnMintWithExternalMinterTokenPool *BurnMintWithExternalMinterTokenPool) Address() common.Address {
	return _BurnMintWithExternalMinterTokenPool.address
}

type BurnMintWithExternalMinterTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetMinter(opts *bind.CallOpts) (common.Address, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

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

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetChainRateLimiterConfigs(opts *bind.TransactOpts, remoteChainSelectors []uint64, outboundConfigs []RateLimiterConfig, inboundConfigs []RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnMintWithExternalMinterTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnMintWithExternalMinterTokenPoolAllowListRemove, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnMintWithExternalMinterTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnMintWithExternalMinterTokenPoolConfigChanged, error)

	FilterInboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumedIterator, error)

	WatchInboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseInboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterTokenPoolInboundRateLimitConsumed, error)

	FilterLockedOrBurned(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolLockedOrBurnedIterator, error)

	WatchLockedOrBurned(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolLockedOrBurned, remoteChainSelector []uint64) (event.Subscription, error)

	ParseLockedOrBurned(log types.Log) (*BurnMintWithExternalMinterTokenPoolLockedOrBurned, error)

	FilterOutboundRateLimitConsumed(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumedIterator, error)

	WatchOutboundRateLimitConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed, remoteChainSelector []uint64) (event.Subscription, error)

	ParseOutboundRateLimitConsumed(log types.Log) (*BurnMintWithExternalMinterTokenPoolOutboundRateLimitConsumed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintWithExternalMinterTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnMintWithExternalMinterTokenPoolRateLimitAdminSet, error)

	FilterReleasedOrMinted(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolReleasedOrMintedIterator, error)

	WatchReleasedOrMinted(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolReleasedOrMinted, remoteChainSelector []uint64) (event.Subscription, error)

	ParseReleasedOrMinted(log types.Log) (*BurnMintWithExternalMinterTokenPoolReleasedOrMinted, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnMintWithExternalMinterTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintWithExternalMinterTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnMintWithExternalMinterTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintWithExternalMinterTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintWithExternalMinterTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnMintWithExternalMinterTokenPoolRouterUpdated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
