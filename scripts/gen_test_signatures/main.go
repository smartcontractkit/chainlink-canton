// Tool to generate ECDSA secp256k1 test signatures for Daml tests.
// Usage: go run scripts/gen_test_signatures/main.go <hex-encoded-preimage>
//
// The preimage should be: versionTag || encodedMessage
// For CCIP: "49ff34ed" || encodedMessageV1
//
// Outputs Daml code with:
// - 20 public keys (uncompressed, 65 bytes each, starting with 04)
// - 20 r values (32 bytes each)
// - 20 s values (32 bytes each)
// - The keccak256 hash of the input preimage

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
		log.Fatal("Usage: go run main.go <hex-encoded-preimage>")
	}

	preimageHex := os.Args[1]
	preimageBytes, err := hex.DecodeString(preimageHex)
	if err != nil {
		log.Fatalf("Failed to decode preimage: %v", err)
	}

	// Compute keccak256 hash of the preimage
	hash := crypto.Keccak256(preimageBytes)
	fmt.Printf("-- Hash of preimage: %s\n", hex.EncodeToString(hash))
	fmt.Printf("-- Preimage length: %d bytes\n\n", len(preimageBytes))

	numKeys := 20

	// Generate key pairs
	var privateKeys []*ecdsa.PrivateKey
	for i := 0; i < numKeys; i++ {
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
	var rValues, sValues [][]byte
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
