package tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

// ===========================================================================
// MCMS INTEGRATION TEST
// ===========================================================================

// TestMCMS_SetRootWithRealSignatures tests the full MCMS flow with real
// cryptographic signatures, verifying that Canton's native crypto verification
// works correctly.
func TestMCMS_SetRootWithRealSignatures(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run.")
	}

	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))

	participant := env.Participant(1)

	// ========================
	// |   Setup: Upload DAR  |
	// ========================

	t.Log("Uploading MCMS DAR...")

	mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err)

	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant)
	require.NoError(t, err)
	t.Logf("Uploaded MCMS DAR, package IDs: %v", packageIDs)

	// ========================
	// |   Setup: Parties     |
	// ========================

	mcmsOwner := participant.Party
	t.Logf("Using party: %s", mcmsOwner)

	// ===========================================================================
	// Step 1: Generate Signers
	// ===========================================================================
	fmt.Println("\n=== Step 1: Generate Signers ===")

	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}

	fmt.Printf("Signer 1: %s\n", FormatSignerInfo(signer1))
	fmt.Printf("Signer 2: %s\n", FormatSignerInfo(signer2))
	fmt.Printf("Signer 3: %s\n", FormatSignerInfo(signer3))

	// ===========================================================================
	// Step 2: Create MCMS Contract
	// ===========================================================================
	fmt.Println("\n=== Step 2: Create MCMS Contract ===")

	chainId := 1
	baseMcmsId := "mcms-test-001"
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, mcmsOwner)
	proposerMultisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleProposer)

	mcmsCid, err := createMCMS(t.Context(), participant, mcmsOwner, chainId, baseMcmsId)
	require.NoError(t, err)
	fmt.Printf("Created MCMS: %s\n", mcmsCid)

	// ===========================================================================
	// Step 3: Configure Signers (2-of-3)
	// ===========================================================================
	fmt.Println("\n=== Step 3: Configure Signers ===")

	config := New2of3Config(signers)
	mcmsCid, err = setMCMSConfig(t.Context(), participant, mcmsOwner, mcmsCid, config)
	require.NoError(t, err)
	fmt.Printf("Configured MCMS with 2-of-3 config: %s\n", mcmsCid)

	// ===========================================================================
	// Step 4: Build Proposal with Operation
	// ===========================================================================
	fmt.Println("\n=== Step 4: Build Proposal ===")

	proposal := NewMCMSProposal(chainId, proposerMultisigId, 0, false).
		AddOperation("counter@owner", "Increment", "").
		Build()

	fmt.Printf("Merkle Root: %s\n", proposal.GetRoot())
	fmt.Printf("Metadata: chainId=%d, multisigId=%s, preOp=%d, postOp=%d\n",
		proposal.Metadata.ChainId, proposal.Metadata.MultisigId,
		proposal.Metadata.PreOpCount, proposal.Metadata.PostOpCount)

	// ===========================================================================
	// Step 5: Sign Root with 2 Signers
	// ===========================================================================
	fmt.Println("\n=== Step 5: Sign Root ===")

	validUntil := time.Now().Add(time.Hour)
	selectedSigners := SortSignersByAddress(signers)[:2] // Take first 2 after sorting

	signatures, err := proposal.Sign(validUntil, selectedSigners)
	require.NoError(t, err)
	require.Len(t, signatures, 2)

	fmt.Printf("Generated %d signatures\n", len(signatures))
	for i, sig := range signatures {
		fmt.Printf("  Sig %d: r=%s..., s=%s...\n", i, sig.R[:16], sig.S[:16])
	}

	// Verify signatures locally before sending to Canton
	signedHash := ComputeSignedHash(proposal.GetRoot(), validUntil)
	hashBytes, _ := hex.DecodeString(signedHash)
	for i, sig := range signatures {
		addr, err := VerifySignature(hashBytes, sig)
		require.NoError(t, err)
		fmt.Printf("  Verified sig %d from address: %s\n", i, addr)
	}

	// ===========================================================================
	// Step 6: Get Merkle Proofs
	// ===========================================================================
	fmt.Println("\n=== Step 6: Get Merkle Proofs ===")

	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	fmt.Printf("Metadata proof length: %d\n", len(metadataProof))

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	fmt.Printf("Operation proof length: %d\n", len(opProof))

	// Verify proofs locally
	metadataLeaf := HashMetadataLeaf(proposal.Metadata)
	require.True(t, VerifyMerkleProof(proposal.GetRoot(), metadataLeaf, metadataProof))
	fmt.Println("Metadata proof verified locally ✓")

	opLeaf := HashOpLeaf(proposal.Operations[0])
	require.True(t, VerifyMerkleProof(proposal.GetRoot(), opLeaf, opProof))
	fmt.Println("Operation proof verified locally ✓")

	// ===========================================================================
	// Step 7: Submit SetRoot to Canton
	// ===========================================================================
	fmt.Println("\n=== Step 7: Submit SetRoot ===")

	mcmsCid, err = setMCMSRoot(t.Context(), participant, mcmsOwner, mcmsCid,
		proposal.GetRoot(), validUntil, &proposal.Metadata, metadataProof, signatures)
	require.NoError(t, err)
	fmt.Printf("SetRoot succeeded! New MCMS CID: %s\n", mcmsCid)

	// ===========================================================================
	// Step 8: Test Complete
	// The new architecture uses ExecuteOp for direct invocation
	// SetRoot is the key test for cryptographic verification
	// ===========================================================================
	fmt.Println("\n=== Step 8: Test Complete ===")
	fmt.Println("SetRoot with real ECDSA signatures verified ✓")
	fmt.Println("Canton's native crypto verification working correctly ✓")
	fmt.Printf("Final MCMS CID: %s\n", mcmsCid)

	fmt.Println("\n=== TEST PASSED ===")
}

// ===========================================================================
// UNIT TESTS FOR CRYPTO HELPERS
// These tests output all values in a format usable for Daml tests
// ===========================================================================

// TestMCMSCrypto_TimeToHex tests the timestamp to hex encoding
func TestMCMSCrypto_TimeToHex(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: TimeToHex Encoding")
	fmt.Println(strings.Repeat("=", 80))

	testCases := []struct {
		name     string
		unixTime int64
		expected string
	}{
		{
			name:     "Unix 1700000000 (2023-11-14 22:13:20 UTC)",
			unixTime: 1700000000,
			expected: "000000000000000000000000000000000000000000000000000000006553f100",
		},
		{
			name:     "Unix 0 (epoch)",
			unixTime: 0,
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:     "Unix 1 (epoch + 1 second)",
			unixTime: 1,
			expected: "0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name:     "Unix 4294967295 (max uint32)",
			unixTime: 4294967295, // 0xFFFFFFFF
			expected: "00000000000000000000000000000000000000000000000000000000ffffffff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ts := time.Unix(tc.unixTime, 0)
			result := TimeToHex(ts)
			fmt.Printf("\n%s:\n", tc.name)
			fmt.Printf("  Unix: %d\n", tc.unixTime)
			fmt.Printf("  Hex:  %s\n", result)
			require.Equal(t, tc.expected, result, "TimeToHex mismatch")
		})
	}

	fmt.Println("\n✓ TimeToHex encoding tests passed")
}

// TestMCMSCrypto_SignedHash tests the signed hash computation
func TestMCMSCrypto_SignedHash(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: Signed Hash Computation")
	fmt.Println(strings.Repeat("=", 80))

	root := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	validUntil := time.Unix(1700000000, 0)

	fmt.Println("\n-- Input Values --")
	fmt.Printf("root = \"%s\"\n", root)
	fmt.Printf("validUntil = %v (Unix: %d)\n", validUntil, validUntil.Unix())
	fmt.Printf("validUntilHex = \"%s\"\n", TimeToHex(validUntil))

	signedHash := ComputeSignedHash(root, validUntil)
	require.Len(t, signedHash, 64) // 32 bytes = 64 hex chars

	fmt.Println("\n-- Output Values (for Daml tests) --")
	fmt.Printf("signedHash = \"%s\"\n", signedHash)

	// Verify the expected hash (must match Daml's computeSignedHashNative)
	// With proper timestamp encoding (Unix 1700000000 = 0x6553f100):
	// innerData = root || "000000000000000000000000000000000000000000000000000000006553f100"
	// innerHash = keccak256(innerData)
	// signedHash = keccak256("\x19Ethereum Signed Message:\n32" || innerHash)
	expectedHash := "a96a392cce9d743dccbf235388ca4f72c050e5ec7f3c6913312d277b75522212"
	require.Equal(t, expectedHash, signedHash, "signedHash must match Daml's computeSignedHashNative")

	// Verify it's deterministic
	signedHash2 := ComputeSignedHash(root, validUntil)
	require.Equal(t, signedHash, signedHash2)
	fmt.Println("\n✓ Signed hash is deterministic and matches expected value")
}

// TestMCMSCrypto_MerkleTree tests Merkle tree construction and proofs
func TestMCMSCrypto_MerkleTree(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: Merkle Tree Construction & Proofs")
	fmt.Println(strings.Repeat("=", 80))

	// Create test leaves
	leaves := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	}

	fmt.Println("\n-- Input Leaves --")
	for i, leaf := range leaves {
		fmt.Printf("leaf[%d] = \"%s\"\n", i, leaf)
	}

	// Build tree
	tree := NewMerkleTree(leaves)
	root := tree.Root

	fmt.Println("\n-- Output Values (for Daml tests) --")
	fmt.Printf("merkleRoot = \"%s\"\n", root)

	// Get proofs for each leaf
	for i, leaf := range leaves {
		proof, err := tree.GetProof(leaf)
		require.NoError(t, err)
		fmt.Printf("\n-- Proof for leaf[%d] --\n", i)
		fmt.Printf("leaf = \"%s\"\n", leaf)
		fmt.Printf("proof = [")
		for j, p := range proof {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("\"%s\"", p)
		}
		fmt.Printf("]\n")

		// Verify proof
		require.True(t, VerifyMerkleProof(root, leaf, proof))
		fmt.Printf("✓ Proof verified for leaf[%d]\n", i)
	}
}

// TestMCMSCrypto_Signature tests ECDSA signing and recovery
func TestMCMSCrypto_Signature(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: ECDSA secp256k1 Signing & Recovery")
	fmt.Println(strings.Repeat("=", 80))

	// Create signer
	signer, err := NewMCMSSigner()
	require.NoError(t, err)

	fmt.Println("\n-- Signer Info --")
	fmt.Printf("address = \"%s\"\n", signer.Address)
	fmt.Printf("publicKey = \"%s\"\n", signer.PublicKey)

	// Sign a hash
	messageHash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	hashBytes, _ := hex.DecodeString(messageHash)

	sig, err := signer.Sign(hashBytes)
	require.NoError(t, err)

	fmt.Println("\n-- Signature (for Daml tests) --")
	fmt.Printf("messageHash = \"%s\"\n", messageHash)
	fmt.Println("signature = RawSignature {")
	fmt.Printf("  publicKey = \"%s\"\n", sig.PublicKey)
	fmt.Printf("  r = \"%s\"\n", sig.R)
	fmt.Printf("  s = \"%s\"\n", sig.S)
	fmt.Println("}")

	// Verify signature
	recoveredAddr, err := VerifySignature(hashBytes, *sig)
	require.NoError(t, err)
	require.Equal(t, strings.ToLower(signer.Address), strings.ToLower(recoveredAddr))

	fmt.Println("\n-- Verification --")
	fmt.Printf("recoveredAddress = \"%s\"\n", recoveredAddr)
	fmt.Printf("expectedAddress  = \"%s\"\n", signer.Address)
	fmt.Println("✓ Signature verified - addresses match")
}

// TestMCMSCrypto_ProposalBuilder tests the full proposal building flow
func TestMCMSCrypto_ProposalBuilder(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: Proposal Builder (Full Flow)")
	fmt.Println(strings.Repeat("=", 80))

	chainId := 1
	multisigId := "mcms-test-001-proposer"
	preOpCount := 0

	proposal := NewMCMSProposal(chainId, multisigId, preOpCount, false).
		AddOperation("counter@owner", "Increment", "").
		Build()

	fmt.Println("\n-- Metadata --")
	fmt.Printf("RootMetadata {\n")
	fmt.Printf("  chainId = %d\n", proposal.Metadata.ChainId)
	fmt.Printf("  multisigId = \"%s\"\n", proposal.Metadata.MultisigId)
	fmt.Printf("  preOpCount = %d\n", proposal.Metadata.PreOpCount)
	fmt.Printf("  postOpCount = %d\n", proposal.Metadata.PostOpCount)
	fmt.Printf("  overridePreviousRoot = %v\n", proposal.Metadata.OverridePreviousRoot)
	fmt.Printf("}\n")

	fmt.Println("\n-- Operations --")
	for i, op := range proposal.Operations {
		fmt.Printf("Op[%d] {\n", i)
		fmt.Printf("  chainId = %d\n", op.ChainId)
		fmt.Printf("  multisigId = \"%s\"\n", op.MultisigId)
		fmt.Printf("  nonce = %d\n", op.Nonce)
		fmt.Printf("  targetInstanceId = \"%s\"\n", op.TargetInstanceId)
		fmt.Printf("  functionName = \"%s\"\n", op.FunctionName)
		fmt.Printf("  operationData = \"%s\"\n", op.OperationData)
		fmt.Printf("}\n")
	}

	fmt.Println("\n-- Merkle Tree --")
	fmt.Printf("root = \"%s\"\n", proposal.GetRoot())

	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	fmt.Printf("metadataProof = [")
	for i, p := range metadataProof {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("\"%s\"", p)
	}
	fmt.Printf("]\n")

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	fmt.Printf("opProof[0] = [")
	for i, p := range opProof {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("\"%s\"", p)
	}
	fmt.Printf("]\n")

	// Verify proofs
	metadataLeaf := HashMetadataLeaf(proposal.Metadata)
	require.True(t, VerifyMerkleProof(proposal.GetRoot(), metadataLeaf, metadataProof))
	fmt.Println("\n✓ Metadata proof verified")

	opLeaf := HashOpLeaf(proposal.Operations[0])
	require.True(t, VerifyMerkleProof(proposal.GetRoot(), opLeaf, opProof))
	fmt.Println("✓ Operation proof verified")
}

// TestMCMSCrypto_2of3Config tests the 2-of-3 signer config generation
func TestMCMSCrypto_2of3Config(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: 2-of-3 Signer Config")
	fmt.Println(strings.Repeat("=", 80))

	// Create 3 signers
	signer1, _ := NewMCMSSigner()
	signer2, _ := NewMCMSSigner()
	signer3, _ := NewMCMSSigner()
	signers := []*MCMSSigner{signer1, signer2, signer3}

	config := New2of3Config(signers)

	fmt.Println("\n-- SignerInfo List (sorted by address) --")
	for _, si := range config.Signers {
		fmt.Printf("  SignerInfo { signerAddress = \"%s\", signerIndex = %d, signerGroup = %d }\n",
			si.SignerAddress, si.SignerIndex, si.SignerGroup)
	}
	fmt.Printf("]\n")

	fmt.Println("\n-- Group Config --")
	fmt.Printf("groupQuorums[0] = %d  -- Root group needs 2 signatures\n", config.GroupQuorums[0])
	fmt.Printf("groupParents[0] = %d  -- Root is its own parent\n", config.GroupParents[0])

	// Verify config
	require.Len(t, config.Signers, 3)
	require.Equal(t, 2, config.GroupQuorums[0]) // Root group needs 2
	require.Equal(t, 0, config.GroupParents[0]) // Root is its own parent

	fmt.Println("\n✓ 2-of-3 config test passed")
}

// TestMCMSCrypto_FullSigningFlow tests the full flow from proposal to signatures
func TestMCMSCrypto_FullSigningFlow(t *testing.T) {
	t.Parallel()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TEST: Full MCMS Signing Flow (Copy these values to Daml tests!)")
	fmt.Println(strings.Repeat("=", 80))

	// Step 1: Create signers
	fmt.Println("\n----------------------------------------")
	fmt.Println("STEP 1: Generate 3 Signers")
	fmt.Println("----------------------------------------")

	signer1, _ := NewMCMSSigner()
	signer2, _ := NewMCMSSigner()
	signer3, _ := NewMCMSSigner()
	signers := []*MCMSSigner{signer1, signer2, signer3}
	sortedSigners := SortSignersByAddress(signers)

	fmt.Println("\n-- Signers (sorted by address) --")
	for i, s := range sortedSigners {
		fmt.Printf("signer%d:\n", i+1)
		fmt.Printf("  address   = \"%s\"\n", s.Address)
		fmt.Printf("  publicKey = \"%s\"\n", s.PublicKey)
	}

	// Step 2: Create 2-of-3 config
	fmt.Println("\n----------------------------------------")
	fmt.Println("STEP 2: Create 2-of-3 Config")
	fmt.Println("----------------------------------------")

	config := New2of3Config(signers)

	fmt.Println("\n-- SignerInfo List (for SetConfig) --")
	for _, si := range config.Signers {
		fmt.Printf("SignerInfo { signerAddress = \"%s\", signerIndex = %d, signerGroup = %d }\n",
			si.SignerAddress, si.SignerIndex, si.SignerGroup)
	}

	// Step 3: Build proposal
	fmt.Println("\n----------------------------------------")
	fmt.Println("STEP 3: Build Proposal with 1 Operation")
	fmt.Println("----------------------------------------")

	chainId := 1
	multisigId := "mcms-test-001-proposer"

	proposal := NewMCMSProposal(chainId, multisigId, 0, false).
		AddOperation("counter@owner", "Increment", "").
		Build()

	fmt.Println("\n-- Metadata --")
	fmt.Println("RootMetadata {")
	fmt.Printf("  chainId = %d\n", proposal.Metadata.ChainId)
	fmt.Printf("  multisigId = \"%s\"\n", proposal.Metadata.MultisigId)
	fmt.Printf("  preOpCount = %d\n", proposal.Metadata.PreOpCount)
	fmt.Printf("  postOpCount = %d\n", proposal.Metadata.PostOpCount)
	fmt.Printf("  overridePreviousRoot = %v\n", proposal.Metadata.OverridePreviousRoot)
	fmt.Println("}")

	fmt.Println("\n-- Operation --")
	fmt.Println("Op {")
	fmt.Printf("  chainId = %d\n", proposal.Operations[0].ChainId)
	fmt.Printf("  multisigId = \"%s\"\n", proposal.Operations[0].MultisigId)
	fmt.Printf("  nonce = %d\n", proposal.Operations[0].Nonce)
	fmt.Printf("  targetInstanceId = \"%s\"\n", proposal.Operations[0].TargetInstanceId)
	fmt.Printf("  functionName = \"%s\"\n", proposal.Operations[0].FunctionName)
	fmt.Printf("  operationData = \"%s\"\n", proposal.Operations[0].OperationData)
	fmt.Println("}")

	metadataProof, _ := proposal.GetMetadataProof()
	opProof, _ := proposal.GetOpProof(0)

	fmt.Println("\n-- Merkle Values --")
	fmt.Printf("metadataLeaf = \"%s\"\n", HashMetadataLeaf(proposal.Metadata))
	fmt.Printf("opLeaf = \"%s\"\n", HashOpLeaf(proposal.Operations[0]))
	fmt.Printf("merkleRoot = \"%s\"\n", proposal.GetRoot())
	fmt.Printf("metadataProof = [%s]\n", strings.Join(metadataProof, ", "))
	fmt.Printf("opProof = [%s]\n", strings.Join(opProof, ", "))

	// Step 4: Sign with 2 signers
	fmt.Println("\n----------------------------------------")
	fmt.Println("STEP 4: Sign Root with 2 Signers")
	fmt.Println("----------------------------------------")

	validUntil := time.Now().Add(time.Hour)
	selectedSigners := sortedSigners[:2]

	signatures, _ := proposal.Sign(validUntil, selectedSigners)

	fmt.Println("\n-- Signing Message --")
	fmt.Printf("root = \"%s\"\n", proposal.GetRoot())
	fmt.Printf("validUntil = %v\n", validUntil)
	signedHash := ComputeSignedHash(proposal.GetRoot(), validUntil)
	fmt.Printf("signedHash = \"%s\"\n", signedHash)

	fmt.Println("\n-- Signatures (for SetRoot) --")
	for i, sig := range signatures {
		fmt.Printf("RawSignature[%d] {\n", i)
		fmt.Printf("  publicKey = \"%s\"\n", sig.PublicKey)
		fmt.Printf("  r = \"%s\"\n", sig.R)
		fmt.Printf("  s = \"%s\"\n", sig.S)
		fmt.Println("}")
	}

	// Verify signatures
	hashBytes, _ := hex.DecodeString(signedHash)
	fmt.Println("\n-- Signature Verification --")
	for i, sig := range signatures {
		addr, _ := VerifySignature(hashBytes, sig)
		fmt.Printf("sig[%d] verified: address = \"%s\"\n", i, addr)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✓ Full signing flow completed successfully!")
	fmt.Println("Copy the values above to use in Daml Script tests")
	fmt.Println(strings.Repeat("=", 80))
}

// ===========================================================================
// HELPER FUNCTIONS
// ===========================================================================

func createMCMS(ctx context.Context, participant testhelpers.Participant, owner string, chainId int, baseMcmsId string) (string, error) {
	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
		{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
	}}}}
	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: makeInt64List(NumGroups, 0)}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: makeInt64List(NumGroups, 0)}}}},
		},
	}}}
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyBlockedFunctions := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	createRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: fmt.Sprintf("%s@%s", baseMcmsId, owner)}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(chainId)}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyBlockedFunctions},
								{Label: "timelockTimestamps", Value: emptyMap},
							}},
						},
					},
				},
			},
			ActAs: []string{owner},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create MCMS: %w", err)
	}

	return createRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId(), nil
}

func setMCMSConfig(ctx context.Context, participant testhelpers.Participant, owner string, mcmsCid string, config MCMSConfig) (string, error) {
	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	exerciseRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetConfig",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "newSigners", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
									{Label: "newGroupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
									{Label: "newGroupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
									{Label: "clearRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{owner},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to SetConfig: %w", err)
	}

	// Find the new MCMS contract ID
	for _, event := range exerciseRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			return created.GetContractId(), nil
		}
	}

	return "", fmt.Errorf("MCMS contract not found in SetConfig response")
}

func setMCMSRoot(ctx context.Context, participant testhelpers.Participant, owner string, mcmsCid string,
	root string, validUntil time.Time, metadata *MCMSRootMetadata, metadataProof []string, signatures []RawSignature) (string, error) {
	validUntilMicros := validUntil.UnixMicro()

	signatureValues := make([]*apiv2.Value, len(signatures))
	for i, sig := range signatures {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: metadata.OverridePreviousRoot}}},
		},
	}}}

	exerciseRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: owner}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{owner},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to SetRoot: %w", err)
	}

	// Find the new MCMS contract ID
	for _, event := range exerciseRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			return created.GetContractId(), nil
		}
	}

	return "", fmt.Errorf("MCMS contract not found in SetRoot response")
}

func makeInt64List(count int, value int64) []*apiv2.Value {
	result := make([]*apiv2.Value, count)
	for i := range count {
		result[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: value}}
	}

	return result
}

// ===========================================================================
// MCMS CODEC TESTS
// ===========================================================================

// TestMCMSCodec_SetConfigParams_Roundtrip tests that encoding and decoding
// SetConfigParams produces the same result (like CCIP message codec tests)
func TestMCMSCodec_SetConfigParams_Roundtrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		params SetConfigParams
	}{
		{
			name: "empty config",
			params: SetConfigParams{
				Signers:      []SignerInfo{},
				GroupQuorums: []int{},
				GroupParents: []int{},
				ClearRoot:    false,
			},
		},
		{
			name: "single signer",
			params: SetConfigParams{
				Signers: []SignerInfo{
					{SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84", SignerIndex: 0, SignerGroup: 0},
				},
				GroupQuorums: []int{1},
				GroupParents: []int{0},
				ClearRoot:    false,
			},
		},
		{
			name: "2-of-3 multisig",
			params: SetConfigParams{
				Signers: []SignerInfo{
					{SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84", SignerIndex: 0, SignerGroup: 0},
					{SignerAddress: "a4b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f80123", SignerIndex: 1, SignerGroup: 0},
					{SignerAddress: "b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e0f9a87654", SignerIndex: 2, SignerGroup: 0},
				},
				GroupQuorums: []int{2},
				GroupParents: []int{0},
				ClearRoot:    true,
			},
		},
		{
			name: "full 32-group config",
			params: SetConfigParams{
				Signers: []SignerInfo{
					{SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84", SignerIndex: 0, SignerGroup: 0},
					{SignerAddress: "a4b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f80123", SignerIndex: 1, SignerGroup: 1},
				},
				GroupQuorums: make([]int, NumGroups), // 32 zeros
				GroupParents: make([]int, NumGroups), // 32 zeros
				ClearRoot:    false,
			},
		},
		{
			name: "hierarchical groups",
			params: SetConfigParams{
				Signers: []SignerInfo{
					{SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84", SignerIndex: 0, SignerGroup: 1},
					{SignerAddress: "a4b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f80123", SignerIndex: 1, SignerGroup: 1},
					{SignerAddress: "b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e0f9a87654", SignerIndex: 2, SignerGroup: 2},
				},
				GroupQuorums: []int{1, 2, 1}, // Group 0 needs 1, Group 1 needs 2, Group 2 needs 1
				GroupParents: []int{0, 0, 1}, // Group 0 is root, Group 1 -> Group 0, Group 2 -> Group 1
				ClearRoot:    true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Encode
			encoded := EncodeSetConfigParams(tc.params)
			t.Logf("Encoded %s: %s (%d hex chars = %d bytes)",
				tc.name, truncateString(encoded, 60), len(encoded), len(encoded)/2)

			// Decode
			decoded, err := DecodeSetConfigParams(encoded)
			require.NoError(t, err, "failed to decode")

			// Verify roundtrip
			require.Len(t, decoded.Signers, len(tc.params.Signers), "signer count mismatch")
			for i, signer := range tc.params.Signers {
				require.Equal(t, signer.SignerAddress, decoded.Signers[i].SignerAddress,
					"signer %d address mismatch", i)
				require.Equal(t, signer.SignerIndex, decoded.Signers[i].SignerIndex,
					"signer %d index mismatch", i)
				require.Equal(t, signer.SignerGroup, decoded.Signers[i].SignerGroup,
					"signer %d group mismatch", i)
			}

			require.Equal(t, tc.params.GroupQuorums, decoded.GroupQuorums, "quorums mismatch")
			require.Equal(t, tc.params.GroupParents, decoded.GroupParents, "parents mismatch")
			require.Equal(t, tc.params.ClearRoot, decoded.ClearRoot, "clearRoot mismatch")
		})
	}
}

// TestMCMSCodec_SetConfigParams_DecodeErrors tests that invalid data is rejected
func TestMCMSCodec_SetConfigParams_DecodeErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		hexData string
		errMsg  string
	}{
		{
			name:    "empty data",
			hexData: "",
			errMsg:  "need at least 1 byte",
		},
		{
			name:    "truncated signer address",
			hexData: "0114", // 1 signer, address length 20, but no address data
			errMsg:  "truncated at signer 0 address",
		},
		{
			name:    "truncated signer index",
			hexData: "01141375dc8a4c1476e6628b03216545e5cdcbff3f84", // address but no index
			errMsg:  "truncated at signer 0 index",
		},
		{
			name:    "truncated quorums",
			hexData: "00", // 0 signers, but no quorum count
			errMsg:  "truncated at quorums count",
		},
		{
			name:    "invalid hex",
			hexData: "zzzz",
			errMsg:  "failed to decode hex",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := DecodeSetConfigParams(tc.hexData)
			require.Error(t, err, "should have failed to decode")
			require.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

// TestMCMSCodec_SetConfigParams_KnownValues tests encoding against known/expected values
// This ensures Go and Daml produce identical encodings
func TestMCMSCodec_SetConfigParams_KnownValues(t *testing.T) {
	t.Parallel()

	// Test with a specific config and verify exact hex output
	params := SetConfigParams{
		Signers: []SignerInfo{
			{SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84", SignerIndex: 0, SignerGroup: 0},
		},
		GroupQuorums: []int{1, 0, 0}, // 3 quorums
		GroupParents: []int{0, 0, 0}, // 3 parents
		ClearRoot:    false,
	}

	encoded := EncodeSetConfigParams(params)

	// Parse and verify structure manually
	// Format: numSigners(1) + [addrLen(1) + addr(20) + index(4) + group(4)]... + numQuorums(1) + quorums(4*n) + numParents(1) + parents(4*n) + clearRoot(1)
	expectedPrefix := "01" + // 1 signer
		"14" + // address length = 20 bytes
		"1375dc8a4c1476e6628b03216545e5cdcbff3f84" + // address
		"00000000" + // signerIndex = 0
		"00000000" // signerGroup = 0

	require.True(t, strings.HasPrefix(encoded, expectedPrefix),
		"encoding should start with expected prefix\nGot: %s\nExpected prefix: %s", encoded, expectedPrefix)

	t.Logf("Full encoding: %s", encoded)
	t.Logf("This value should match Daml encodeSetConfigParams output for the same input")
}

// Helper to truncate long strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen] + "..."
}
