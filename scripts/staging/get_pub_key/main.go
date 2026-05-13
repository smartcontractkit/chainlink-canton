package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Your type-2 EIP-1559 transaction fields
	chainID := big.NewInt(11155111)

	to := common.HexToAddress("0xd41f256d605324aA437762e4F2Ae30EbEf0BF228")
	nonce := uint64(46)
	gasLimit := uint64(201411)
	value := big.NewInt(0)
	maxPriorityFeePerGas := big.NewInt(1000000)
	maxFeePerGas := big.NewInt(1000026)

	data := mustDecodeHex("0x96f4e9f90000000000000000000000000000000000000000000000008c4ae44a1d2f2023000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000207bbbcf8b7915d25de57df173fde169109f27d636f369c41e6652ebec8970863b000000000000000000000000000000000000000000000000000000000000001868656c6c6f2066726f6d2065766d20746f2063616e746f6e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003ca69dd4aa000186a000010114fbfddbd496995ae27687cdb75e492e6911608a7f000014eba517d200000000000000000000000000000000000000000000000000")

	// Signature fields from cast tx
	r := new(big.Int)
	s := new(big.Int)
	r.SetString("d5f4a81e182e96ba653833985a4f5ea5d2c7f8dbb184e4b46d8bac301fc10de2", 16)
	s.SetString("48f564a28d58405a129a25e3483ea19373160057e77c9fe57ae395049d75c444", 16)
	yParity := byte(1) // same as v for type-2 recovery, must be 0 or 1

	// Build the unsigned tx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  maxPriorityFeePerGas,
		GasFeeCap:  maxFeePerGas,
		Gas:        gasLimit,
		To:         &to,
		Value:      value,
		Data:       data,
		AccessList: types.AccessList{},
	})

	// Compute the actual signing hash (NOT the tx hash)
	signer := types.LatestSignerForChainID(chainID)
	digest := signer.Hash(tx)

	// Build 65-byte Ethereum signature: R || S || V
	sig := make([]byte, 65)
	copy(sig[0:32], pad32(r.Bytes()))
	copy(sig[32:64], pad32(s.Bytes()))
	sig[64] = yParity // 0 or 1

	// Recover secp256k1 public key
	pubBytes, err := crypto.Ecrecover(digest.Bytes(), sig)
	if err != nil {
		log.Fatalf("recover pubkey: %v", err)
	}

	// Unmarshal for validation / address derivation
	pubKey, err := crypto.UnmarshalPubkey(pubBytes)
	if err != nil {
		log.Fatalf("unmarshal pubkey: %v", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// 65-byte uncompressed pubkey, starts with 0x04
	uncompressed := crypto.FromECDSAPub(pubKey)

	// 33-byte compressed pubkey
	compressed := crypto.CompressPubkey(pubKey)

	fmt.Println("Signing digest:     0x" + hex.EncodeToString(digest.Bytes()))
	fmt.Println("Recovered address:  " + recoveredAddr.Hex())
	fmt.Println("Expected from:      0x574F139c67b1B3b7BDb9e462A5e80E5664A17105")
	fmt.Println("Uncompressed pubkey:", "0x"+hex.EncodeToString(uncompressed))
	fmt.Println("Compressed pubkey:  ", "0x"+hex.EncodeToString(compressed))
}

func mustDecodeHex(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatalf("decode hex: %v", err)
	}
	return b
}

func pad32(b []byte) []byte {
	if len(b) > 32 {
		log.Fatalf("value longer than 32 bytes")
	}
	if len(b) == 32 {
		return b
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
