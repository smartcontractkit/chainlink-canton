// Tool to generate ECDSA secp256k1 test signatures for Daml tests.
// Usage: go run scripts/gen_test_signatures/main.go <hex-encoded-message>
//
// The input should be the encoded message (without version tag).
// Matches EVM: signers sign keccak256(versionTag || messageId)
// where messageId = keccak256(encodedMessage).
//
// Outputs Daml code with:
// - 20 public keys (uncompressed, 65 bytes each, starting with 04)
// - 20 r values (32 bytes each)
// - 20 s values (32 bytes each)

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <hex-encoded-message>")
	}

	encodedMsgHex := os.Args[1]
	encodedMsgBytes, err := hex.DecodeString(encodedMsgHex)
	if err != nil {
		log.Fatalf("Failed to decode encoded message: %v", err)
	}

	versionTag, _ := hex.DecodeString("e9a05a20")

	// Compute messageId = keccak256(encodedMessage)
	messageId := crypto.Keccak256(encodedMsgBytes)
	fmt.Printf("-- messageId (keccak256 of encoded message): %s\n", hex.EncodeToString(messageId))

	// Compute signedHash = keccak256(versionTag || messageId)
	preimage := append(versionTag, messageId...)
	hash := crypto.Keccak256(preimage)
	fmt.Printf("-- signedHash (keccak256 of versionTag || messageId): %s\n", hex.EncodeToString(hash))
	fmt.Printf("-- Encoded message length: %d bytes\n\n", len(encodedMsgBytes))

	numKeys := 20

	// Generate key pairs
	privateKeys := make([]*ecdsa.PrivateKey, 0, numKeys)
	for range numKeys {
		pk, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		privateKeys = append(privateKeys, pk)
	}

	// Output public keys
	fmt.Println("pubkeys : [BytesHex]")
	fmt.Println("pubkeys = [")
	for i, pk := range privateKeys {
		pub := pk.Public().(*ecdsa.PublicKey)
		pubBytes := crypto.FromECDSAPub(pub)
		comma := ","
		if i == numKeys-1 {
			comma = ""
		}
		fmt.Printf("        \"%s\"%s\n", hex.EncodeToString(pubBytes), comma)
	}
	fmt.Println("    ]")

	// Sign and collect r/s values
	rValues := make([][]byte, 0, numKeys)
	sValues := make([][]byte, 0, numKeys)
	for _, pk := range privateKeys {
		sig, err := crypto.Sign(hash, pk)
		if err != nil {
			log.Fatalf("Failed to sign: %v", err)
		}
		// sig is 65 bytes: R (32) || S (32) || V (1)
		rValues = append(rValues, sig[:32])
		sValues = append(sValues, sig[32:64])
	}

	// Output r values
	fmt.Println("\nrValues : [BytesHex]")
	fmt.Println("rValues = [")
	for i, r := range rValues {
		comma := ","
		if i == numKeys-1 {
			comma = ""
		}
		fmt.Printf("        \"%s\"%s\n", hex.EncodeToString(r), comma)
	}
	fmt.Println("    ]")

	// Output s values
	fmt.Println("\nsValues : [BytesHex]")
	fmt.Println("sValues = [")
	for i, s := range sValues {
		comma := ","
		if i == numKeys-1 {
			comma = ""
		}
		fmt.Printf("        \"%s\"%s\n", hex.EncodeToString(s), comma)
	}
	fmt.Println("    ]")
}
