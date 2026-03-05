package fullnode

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	eABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/http/common"
)

type DeployContractRequest struct {
	OwnerAddress               string `json:"owner_address,omitempty"`
	ABI                        string `json:"abi,omitempty"`
	Bytecode                   string `json:"bytecode,omitempty"`
	Parameter                  string `json:"parameter,omitempty"`
	Name                       string `json:"name,omitempty"`
	CallValue                  int    `json:"call_value,omitempty"`
	FeeLimit                   int    `json:"fee_limit,omitempty"`
	ConsumeUserResourcePercent int    `json:"consume_user_resource_percent,omitempty"`
	OriginEnergyLimit          int    `json:"origin_energy_limit,omitempty"`
	Visible                    bool   `json:"visible,omitempty"`
}

// DeployContract (and any other RPC calls that end up with CreateSmartContract data?) has an additional "contract_address" convenience field,
// which is populated from data already in common.Transaction, see
// https://github.com/tronprotocol/java-tron/blob/c136a26f8140b17f2c05df06fb5efb1bb47d3baa/framework/src/main/java/org/tron/core/services/http/Util.java#L232
type DeployContractResponse struct {
	common.Transaction
	ContractAddress string `json:"contract_address"`
}

func (tc *Client) DeployContract(ownerAddress address.Address, contractName, abiJson, bytecode string, oeLimit, curPercent, feeLimit int, params []interface{}) (*DeployContractResponse, error) {
	parsedABI, err := eABI.JSON(bytes.NewReader([]byte(abiJson)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	if params == nil {
		params = []interface{}{}
	}

	encodedParams, err := parsedABI.Pack("", params...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode params: %w", err)
	}

	reqBody := DeployContractRequest{
		OwnerAddress:               ownerAddress.String(),
		ABI:                        abiJson,
		Bytecode:                   bytecode,
		Parameter:                  hex.EncodeToString(encodedParams),
		Name:                       contractName,
		FeeLimit:                   feeLimit,
		ConsumeUserResourcePercent: curPercent,
		OriginEnergyLimit:          oeLimit,
		Visible:                    true,
	}

	response := DeployContractResponse{}
	if err := tc.Post("/deploycontract", reqBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type GetContractRequest struct {
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
}

type GetContractResponse struct {
	OriginAddress              string          `json:"origin_address,omitempty"`                // Contract creator address
	ContractAddress            string          `json:"contract_address,omitempty"`              // Contract address
	ABI                        *common.JSONABI `json:"abi,omitempty"`                           // ABI
	Bytecode                   string          `json:"bytecode,omitempty"`                      // Bytecode
	CallValue                  int64           `json:"call_value,omitempty"`                    // The amount of TRX passed into the contract when deploying the contract
	ConsumeUserResourcePercent int64           `json:"consume_user_resource_percent,omitempty"` // Proportion of user energy consumption
	Name                       string          `json:"name,omitempty"`                          // contract name
	OriginEnergyLimit          int64           `json:"origin_energy_limit,omitempty"`           // Each transaction is allowed to consume the maximum energy of the contract creator
	CodeHash                   string          `json:"code_hash,omitempty"`                     // code hash
}

func (tc *Client) GetContract(contractAddress address.Address) (*GetContractResponse, error) {
	contractInfo := GetContractResponse{}
	err := tc.Post("/getcontract",
		&GetContractRequest{
			Value:   contractAddress.String(),
			Visible: true,
		}, &contractInfo)

	if err != nil {
		return nil, err
	}

	if contractInfo.ABI == nil {
		return nil, errors.New("could not get contract ABI")
	}

	return &contractInfo, nil
}

type TriggerSmartContractRequest struct {
	OwnerAddress     string `json:"owner_address"`     // Address that triggers the contract, converted to a hex string
	ContractAddress  string `json:"contract_address"`  // Contract address, converted to a hex string
	FunctionSelector string `json:"function_selector"` // Function call, must not be left blank
	Parameter        string `json:"parameter"`         // ABI encoded hex string of parameters
	Data             string `json:"data"`              // The data for interacting with smart contracts, including the contract function and parameters
	FeeLimit         int32  `json:"fee_limit"`         // Maximum TRX consumption, measured in SUN
	CallValue        int64  `json:"call_value"`        // Amount of TRX transferred with this transaction, measured in SUN
	CallTokenValue   int64  `json:"call_token_value"`  // Amount of TRC10 token transferred with this transaction
	TokenId          int64  `json:"token_id"`          // TRC 10 token id
	// typo in spec? json:"Permission_id" https://developers.tron.network/reference/triggersmartcontract
	PermissionId int32 `json:"permission_id"` // for multi-signature
	Visible      bool  `json:"visible"`       // Whether the address is in base58check format
}

type TriggerResult struct {
	Result bool `json:"result"`
}

type TriggerSmartContractResponse struct {
	Result      TriggerResult       `json:"result"`
	Transaction *common.Transaction `json:"transaction"`
}

func (tc *Client) TriggerSmartContract(from, contractAddress address.Address, method string, params []any, feeLimit int32, tAmount int64) (*TriggerSmartContractResponse, error) {
	paramBytes, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode params: %w", err)
	}
	tcRequest := TriggerSmartContractRequest{
		OwnerAddress:     from.String(),
		ContractAddress:  contractAddress.String(),
		FunctionSelector: method,
		Parameter:        hex.EncodeToString(paramBytes),
		FeeLimit:         feeLimit,
		CallValue:        tAmount,
		Visible:          true,
	}
	contractResponse := TriggerSmartContractResponse{}
	err = tc.Post("/triggersmartcontract", tcRequest, &contractResponse)
	if err != nil {
		return nil, err
	}

	return &contractResponse, nil
}

type BroadcastResponse struct {
	Result  bool   `json:"result"`
	Code    string `json:"code"`
	TxID    string `json:"txid"`
	Message string `json:"message"`
}

func (tc *Client) BroadcastTransaction(reqBody *common.Transaction) (*BroadcastResponse, error) {
	if reqBody == nil {
		return nil, errors.New("empty body")
	}

	if len(reqBody.TxID) < 1 {
		return nil, fmt.Errorf("empty transaction ID in request")
	}

	if len(reqBody.Signature) < 1 {
		return nil, fmt.Errorf("no signatures")
	}

	response := BroadcastResponse{}
	err := tc.Post("/broadcasttransaction", reqBody, &response)

	if err != nil {
		return nil, err
	}

	if !response.Result {
		return &response, fmt.Errorf("broadcasting failed. Code: %s, Message: %s", response.Code, response.Message)
	}

	return &response, nil
}
