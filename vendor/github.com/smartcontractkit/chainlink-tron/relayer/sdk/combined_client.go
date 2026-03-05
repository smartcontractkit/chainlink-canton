package sdk

import (
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/http/common"
	"github.com/fbsobreira/gotron-sdk/pkg/http/fullnode"
	"github.com/fbsobreira/gotron-sdk/pkg/http/soliditynode"
)

//go:generate mockery --name CombinedClient --output ../mocks/
type CombinedClient interface {
	FullNodeClient
	TriggerConstantContractFullNode(from, contractAddress address.Address, method string, params []any) (*soliditynode.TriggerConstantContractResponse, error)
	GetNowBlockFullNode() (*soliditynode.Block, error)
	GetBlockByNumFullNode(num int32) (*soliditynode.Block, error)
	GetAccountFullNode(accountAddress address.Address) (*soliditynode.GetAccountResponse, error)
	GetTransactionInfoByIdFullNode(txhash string) (*soliditynode.TransactionInfo, error)

	FullNodeClient() *fullnode.Client
	SolidityClient() *soliditynode.Client
}

type combinedClient struct {
	*fullnode.Client
	solidityClient *soliditynode.Client
}

func NewCombinedClient(fullnodeClient *fullnode.Client, solidityClient *soliditynode.Client) CombinedClient {
	return &combinedClient{
		Client:         fullnodeClient,
		solidityClient: solidityClient,
	}
}

func CreateCombinedClient(fullnodeUrl, soliditynodeUrl *url.URL) (CombinedClient, error) {
	return CreateCombinedClientWithTimeout(fullnodeUrl, soliditynodeUrl, 15*time.Second)
}

func CreateCombinedClientWithTimeout(fullnodeUrl, soliditynodeUrl *url.URL, timeout time.Duration) (CombinedClient, error) {
	httpClient := CreateHttpClientWithTimeout(timeout)
	fullnodeClient := fullnode.NewClient(fullnodeUrl.String(), httpClient)
	soliditynodeClient := soliditynode.NewClient(soliditynodeUrl.String(), httpClient)

	return NewCombinedClient(fullnodeClient, soliditynodeClient), nil
}

func (g *combinedClient) FullNodeClient() *fullnode.Client {
	return g.Client
}

func (g *combinedClient) SolidityClient() *soliditynode.Client {
	return g.solidityClient
}

// We manually override methods that we want to use the solidity client for (all read methods).
// We also provide the fullnode versions of these methods for flexibility

// GetAccount from BASE58 address using solidity client
func (g *combinedClient) GetAccount(accountAddress address.Address) (*soliditynode.GetAccountResponse, error) {
	return g.solidityClient.GetAccount(accountAddress)
}

// GetAccount from BASE58 address using fullnode client
func (g *combinedClient) GetAccountFullNode(accountAddress address.Address) (*soliditynode.GetAccountResponse, error) {
	return g.Client.GetAccount(accountAddress)
}

// GetTransactionInfoByID returns transaction receipt by ID using solidity client
func (g *combinedClient) GetTransactionInfoById(txhash string) (*soliditynode.TransactionInfo, error) {
	return g.solidityClient.GetTransactionInfoById(txhash)
}

// GetTransactionInfoByID returns transaction receipt by ID using fullnode client
func (g *combinedClient) GetTransactionInfoByIdFullNode(txhash string) (*soliditynode.TransactionInfo, error) {
	return g.Client.GetTransactionInfoById(txhash)
}

// TriggerConstantContract and return tx result using solidity client
func (g *combinedClient) TriggerConstantContract(from, contractAddress address.Address, method string, params []any) (*soliditynode.TriggerConstantContractResponse, error) {
	return g.solidityClient.TriggerConstantContract(from, contractAddress, method, params)
}

// TriggerConstantContract and return tx result using solidity client
func (g *combinedClient) TriggerConstantContractFullNode(from, contractAddress address.Address, method string, params []any) (*soliditynode.TriggerConstantContractResponse, error) {
	return g.Client.TriggerConstantContract(from, contractAddress, method, params)
}

// GetNowBlock return TIP block using solidity client
func (g *combinedClient) GetNowBlock() (*soliditynode.Block, error) {
	return g.solidityClient.GetNowBlock()
}

// GetNowBlock return TIP block using fullnode client
func (g *combinedClient) GetNowBlockFullNode() (*soliditynode.Block, error) {
	return g.Client.GetNowBlock()
}

// GetBlockByNum block from number using solidity client
func (g *combinedClient) GetBlockByNum(num int32) (*soliditynode.Block, error) {
	return g.solidityClient.GetBlockByNum(num)
}

// GetBlockByNum block from number using fullnode client
func (g *combinedClient) GetBlockByNumFullNode(num int32) (*soliditynode.Block, error) {
	return g.Client.GetBlockByNum(num)
}

type validatedCombinedClient struct {
	orig    CombinedClient
	chainID *big.Int

	mu   sync.RWMutex
	done bool
	err  error
}

// NewValidatedCombinedClient returns a CombinedClient which wraps c and gates usage with lazy chain ID validation.
func NewValidatedCombinedClient(c CombinedClient, chainID *big.Int) CombinedClient {
	return &validatedCombinedClient{orig: c, chainID: chainID}
}

func (c *validatedCombinedClient) validate() error {
	c.mu.RLock()
	if c.done {
		defer c.mu.RUnlock()
		return c.err
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done { // double check since we re-locked
		return c.err
	}

	// check client chain id matches config chain id
	blockInfo, err := c.orig.GetBlockByNum(0)
	if err != nil {
		return fmt.Errorf("error getting genesis block info: %w", err)
	}

	c.done = true
	// last 4 bytes of genesis block is the chain id for Tron testnets and mainnet
	chainIdHex := blockInfo.BlockID[len(blockInfo.BlockID)-8:]
	chainId, ok := new(big.Int).SetString(chainIdHex, 16)
	if !ok {
		c.err = fmt.Errorf("couldn't parse chain id %s from genesis block", chainIdHex)
		return c.err
	}
	if chainId.Cmp(c.chainID) != 0 {
		c.err = fmt.Errorf("client chain id %s does not match config chain id %s", chainId, c.chainID)
	}
	return c.err
}

func (c *validatedCombinedClient) TriggerConstantContract(from, contractAddress address.Address, method string, params []any) (*soliditynode.TriggerConstantContractResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.TriggerConstantContract(from, contractAddress, method, params)
}

func (c *validatedCombinedClient) EstimateEnergy(from, contractAddress address.Address, method string, params []any, tAmount int64) (*soliditynode.EnergyEstimateResult, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.EstimateEnergy(from, contractAddress, method, params, tAmount)
}

func (c *validatedCombinedClient) GetNowBlock() (*soliditynode.Block, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetNowBlock()
}

func (c *validatedCombinedClient) GetBlockByNum(num int32) (*soliditynode.Block, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetBlockByNum(num)
}

func (c *validatedCombinedClient) GetAccount(accountAddress address.Address) (*soliditynode.GetAccountResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetAccount(accountAddress)
}

func (c *validatedCombinedClient) GetTransactionInfoById(txhash string) (*soliditynode.TransactionInfo, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetTransactionInfoById(txhash)
}

func (c *validatedCombinedClient) DeployContract(ownerAddress address.Address, contractName, abiJson, bytecode string, oeLimit, curPercent, feeLimit int, params []interface{}) (*fullnode.DeployContractResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.DeployContract(ownerAddress, contractName, abiJson, bytecode, oeLimit, curPercent, feeLimit, params)
}

func (c *validatedCombinedClient) GetContract(address address.Address) (*fullnode.GetContractResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetContract(address)
}

func (c *validatedCombinedClient) TriggerSmartContract(from, contractAddress address.Address, method string, params []any, feeLimit int32, tAmount int64) (*fullnode.TriggerSmartContractResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.TriggerSmartContract(from, contractAddress, method, params, feeLimit, tAmount)
}

func (c *validatedCombinedClient) Transfer(fromAddress, toAddress address.Address, amount int64) (*common.Transaction, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.Transfer(fromAddress, toAddress, amount)
}

func (c *validatedCombinedClient) BroadcastTransaction(reqBody *common.Transaction) (*fullnode.BroadcastResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.BroadcastTransaction(reqBody)
}

func (c *validatedCombinedClient) GetEnergyPrices() (*fullnode.EnergyPrices, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetEnergyPrices()
}

func (c *validatedCombinedClient) TriggerConstantContractFullNode(from, contractAddress address.Address, method string, params []any) (*soliditynode.TriggerConstantContractResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.TriggerConstantContract(from, contractAddress, method, params)
}

func (c *validatedCombinedClient) GetNowBlockFullNode() (*soliditynode.Block, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetNowBlockFullNode()
}

func (c *validatedCombinedClient) GetBlockByNumFullNode(num int32) (*soliditynode.Block, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetBlockByNumFullNode(num)
}

func (c *validatedCombinedClient) GetAccountFullNode(accountAddress address.Address) (*soliditynode.GetAccountResponse, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetAccountFullNode(accountAddress)
}

func (c *validatedCombinedClient) GetTransactionInfoByIdFullNode(txhash string) (*soliditynode.TransactionInfo, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c.orig.GetTransactionInfoByIdFullNode(txhash)
}

func (c *validatedCombinedClient) FullNodeClient() *fullnode.Client { return c.orig.FullNodeClient() }

func (c *validatedCombinedClient) SolidityClient() *soliditynode.Client {
	return c.orig.SolidityClient()
}
