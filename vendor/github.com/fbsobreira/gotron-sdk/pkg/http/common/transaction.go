package common

import (
	"encoding/hex"
)

type Transaction struct {
	Visible    bool     `json:"visible"`
	TxID       string   `json:"txID"`
	RawData    RawData  `json:"raw_data"`
	RawDataHex string   `json:"raw_data_hex"`
	Signature  []string `json:"signature"`
}

func (t *Transaction) AddSignatureBytes(signatureBytes []byte) {
	signatureHex := hex.EncodeToString(signatureBytes)
	t.AddSignature(signatureHex)
}

func (t *Transaction) AddSignature(signatureHex string) {
	t.Signature = append(t.Signature, signatureHex)
}

// https://github.com/tronprotocol/java-tron/blob/c136a26f8140b17f2c05df06fb5efb1bb47d3baa/protocol/src/main/protos/core/Tron.proto#L389
type Return struct {
	// populated in transaction info results
	ContractRet string `json:"contractRet"`
	// either SUCESS (intentionally mispelled) or FAILED, populated in trigger results
	Ret string `json:"ret"`
}

type ExecutedTransaction struct {
	Transaction
	Ret []Return `json:"ret,omitempty"`
}

type NewContract struct {
	OriginAddress              string   `json:"origin_address,omitempty"`                // Contract creator address
	ContractAddress            string   `json:"contract_address,omitempty"`              // Contract address
	ABI                        *JSONABI `json:"abi,omitempty"`                           // ABI
	Bytecode                   string   `json:"bytecode,omitempty"`                      // Bytecode
	CallValue                  int64    `json:"call_value,omitempty"`                    // The amount of TRX passed into the contract when deploying the contract
	ConsumeUserResourcePercent int64    `json:"consume_user_resource_percent,omitempty"` // Proportion of user energy consumption
	Name                       string   `json:"name,omitempty"`                          // contract name
	OriginEnergyLimit          int64    `json:"origin_energy_limit,omitempty"`           // Each transaction is allowed to consume the maximum energy of the contract creator
	CodeHash                   string   `json:"code_hash,omitempty"`                     // code hash
}

type ParameterValue struct {
	OwnerAddress    string       `json:"owner_address,omitempty"`
	ToAddress       string       `json:"to_address,omitempty"`
	Data            string       `json:"data,omitempty"`
	ContractAddress string       `json:"contract_address,omitempty"`
	Amount          int64        `json:"amount,omitempty"`
	NewContract     *NewContract `json:"new_contract,omitempty"`
}

type Parameter struct {
	Value   ParameterValue `json:"value,omitempty"`
	TypeUrl string         `json:"type_url,omitempty"`
}

type Contract struct {
	Parameter Parameter `json:"parameter,omitempty"`
	Type      string    `json:"type,omitempty"`
}

type RawData struct {
	Contract      []Contract `json:"contract,omitempty"`
	RefBlockBytes string     `json:"ref_block_bytes,omitempty"`
	RefBlockHash  string     `json:"ref_block_hash,omitempty"`
	Expiration    int64      `json:"expiration,omitempty"`
	FeeLimit      int64      `json:"fee_limit,omitempty"`
	Timestamp     int64      `json:"timestamp,omitempty"`
}
