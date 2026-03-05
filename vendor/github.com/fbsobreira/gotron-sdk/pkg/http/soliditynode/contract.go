package soliditynode

import (
	"encoding/hex"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/http/common"
)

type TriggerConstantContractRequest struct {
	OwnerAddress     string `json:"owner_address"`     // Owner address that triggers the contract. If visible=true, use base58check format, otherwise use hex format
	ContractAddress  string `json:"contract_address"`  // Smart contract address. If visible=true, use base58check format, otherwise use hex format
	FunctionSelector string `json:"function_selector"` // Function call, must not be left blank
	Parameter        string `json:"parameter"`         // ABI encoded hex string of parameters
	Data             string `json:"data"`              // The bytecode of the contract or the data for interacting with smart contracts, including the contract function and parameters
	CallValue        int64  `json:"call_value"`        // Amount of TRX transferred with this transaction, measured in SUN
	CallTokenValue   int64  `json:"call_token_value"`  // Amount of TRC10 token transferred with this transaction
	TokenId          int64  `json:"token_id"`          // TRC10 token id
	Visible          bool   `json:"visible"`           // Whether the address is in base58check format
}

type TriggerConstantContractResponse struct {
	Result         ReturnEnergyEstimate        `json:"result"`
	EnergyUsed     int64                       `json:"energy_used"`     // Estimated energy consumption, including the basic energy consumption and penalty energy consumption
	EnergyPenalty  int64                       `json:"energy_penalty"`  // The penalty energy consumption
	ConstantResult []string                    `json:"constant_result"` // []	Result list
	Transaction    *common.ExecutedTransaction `json:"transaction"`     // Transaction information, refer to GetTransactionByID
}

func (tc *Client) TriggerConstantContract(from, contractAddress address.Address, method string, params []any) (*TriggerConstantContractResponse, error) {
	paramBytes, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode params: %w", err)
	}
	tcRequest := TriggerConstantContractRequest{
		OwnerAddress:     from.String(),
		ContractAddress:  contractAddress.String(),
		FunctionSelector: method,
		Parameter:        hex.EncodeToString(paramBytes),
		Visible:          true,
	}
	response := TriggerConstantContractResponse{}
	err = tc.Post("/triggerconstantcontract", tcRequest, &response)
	if err != nil {
		return nil, err
	}
	if !response.Result.Result {
		return &response, fmt.Errorf("failed to trigger constant contract, code: %s, message: %s", response.Result.Code, response.Result.Message)
	}

	return &response, nil
}

type EnergyEstimateRequest struct {
	OwnerAddress     string `json:"owner_address"`     // Owner address that triggers the contract. If visible=true, use base58check format, otherwise use hex format
	ContractAddress  string `json:"contract_address"`  // Smart contract address. If visible=true, use base58check format, otherwise use hex format
	FunctionSelector string `json:"function_selector"` // Function call, must not be left blank
	Parameter        string `json:"parameter"`
	Data             string `json:"data"`             // The bytecode of the contract or the data for interacting with smart contracts, including the contract function and parameters
	CallValue        int64  `json:"call_value"`       // Amount of TRX transferred with this transaction, measured in SUN
	CallTokenValue   int64  `json:"call_token_value"` // Amount of TRC10 token transferred with this transaction
	TokenId          int64  `json:"token_id"`         // TRC10 token id
	Visible          bool   `json:"visible"`          // Whether the address is in base58check format
}

type ReturnEnergyEstimate struct {
	Result  bool   `json:"result"`  // Is the estimate successful
	Code    string `json:"code"`    // (enum)   response code, an enum type
	Message string `json:"message"` // Result message
}

type EnergyEstimateResult struct {
	Result         ReturnEnergyEstimate `json:"result"`          // Run result
	EnergyRequired int64                `json:"energy_required"` // Estimated energy to run the contract
}

func (tc *Client) EstimateEnergy(from, contractAddress address.Address, method string, params []any, tAmount int64) (*EnergyEstimateResult, error) {
	paramBytes, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode params: %w", err)
	}
	reqBody := EnergyEstimateRequest{
		OwnerAddress:     from.String(),
		ContractAddress:  contractAddress.String(),
		FunctionSelector: method,
		Parameter:        hex.EncodeToString(paramBytes),
		CallValue:        tAmount,
		Visible:          true,
	}

	response := EnergyEstimateResult{}
	err = tc.Post("/estimateenergy", reqBody, &response)
	if err != nil {
		return nil, err
	}
	if !response.Result.Result {
		return &response, fmt.Errorf("failed to estimate energy, code: %s, message: %s", response.Result.Code, response.Result.Message)
	}

	return &response, nil
}
