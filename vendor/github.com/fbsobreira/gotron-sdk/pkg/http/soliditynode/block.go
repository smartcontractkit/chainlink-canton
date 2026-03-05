package soliditynode

import (
	"errors"

	"github.com/fbsobreira/gotron-sdk/pkg/http/common"
)

type BlockHeaderRaw struct {
	Timestamp        int64  `json:"timestamp,omitempty"`
	TxTrieRoot       string `json:"txTrieRoot,omitempty"`
	ParentHash       string `json:"parentHash,omitempty"`
	Number           int64  `json:"number,omitempty"`
	WitnessId        int64  `json:"witness_id,omitempty"`
	WitnessAddress   string `json:"witness_address,omitempty"`
	Version          int32  `json:"version,omitempty"`
	AccountStateRoot string `json:"accountStateRoot,omitempty"`
}

type BlockHeader struct {
	RawData          *BlockHeaderRaw `json:"raw_data,omitempty"`
	WitnessSignature string          `json:"witness_signature,omitempty"`
}

type Block struct {
	BlockID      string                       `json:"blockID"`
	Transactions []common.ExecutedTransaction `json:"transactions,omitempty"`
	BlockHeader  *BlockHeader                 `json:"block_header,omitempty"`
}

func (tc *Client) GetNowBlock() (*Block, error) {
	block := Block{}
	err := tc.Get("/getnowblock", &block)
	if err != nil {
		return nil, err
	}
	if block.BlockHeader == nil {
		return nil, errors.New("failed to retrieve block header")
	}

	return &block, nil
}

type GetBlockByNumRequest struct {
	Num int32 `json:"num"` // defined as int32 in https://developers.tron.network/reference/wallet-getblockbynum
}

func (tc *Client) GetBlockByNum(num int32) (*Block, error) {
	block := Block{}
	err := tc.Post("/getblockbynum",
		&GetBlockByNumRequest{
			Num: num,
		}, &block)
	if err != nil {
		return nil, err
	}
	if block.BlockHeader == nil {
		return nil, errors.New("failed to retrieve block header")
	}

	return &block, nil
}
