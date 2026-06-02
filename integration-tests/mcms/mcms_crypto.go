package tests

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	mrand "math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms"
)

// ===========================================================================
// TYPES (matching Canton MCMS.Types)
// ===========================================================================

// MCMSRole matches Canton Role type
type MCMSRole int

const (
	MCMSRoleProposer  MCMSRole = iota // 0
	MCMSRoleCanceller                 // 1
	MCMSRoleBypasser                  // 2
)

func (r MCMSRole) String() string {
	switch r {
	case MCMSRoleProposer:
		return "proposer"
	case MCMSRoleCanceller:
		return "canceller"
	case MCMSRoleBypasser:
		return "bypasser"
	default:
		return "unknown"
	}
}

// SignerInfo matches Canton SignerInfo
type SignerInfo struct {
	SignerAddress string // EVM address (40 hex chars, lowercase, no 0x prefix)
	SignerIndex   int    // Position in signers list (0-indexed)
	SignerGroup   int    // Group assignment (0-31)
}

// MCMSConfig matches Canton MultisigConfig
type MCMSConfig struct {
	Signers      []SignerInfo
	GroupQuorums []int // Length 32
	GroupParents []int // Length 32
}

// MCMSOp matches Canton Op
type MCMSOp struct {
	ChainId               int
	MultisigId            string
	Nonce                 int
	TargetInstanceAddress string
	FunctionName          string
	OperationData         string // hex encoded
}

// TimelockCall matches Canton TimelockCall
type TimelockCall struct {
	TargetInstanceAddress string
	FunctionName          string
	OperationData         string // hex encoded
}

// ZeroHash represents "no predecessor" in timelock operations
const ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"

var (
	// opLeafDomainSeparator = keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_CANTON")
	opLeafDomainSeparator = func() string {
		data := []byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_CANTON")
		return hex.EncodeToString(crypto.Keccak256(data))
	}()
	// metadataLeafDomainSeparator = keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA_CANTON")
	metadataLeafDomainSeparator = func() string {
		data := []byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA_CANTON")
		return hex.EncodeToString(crypto.Keccak256(data))
	}()
)

// ===========================================================================
// TYPE CONVERSION HELPERS (local <-> generated mcms bindings)
// ===========================================================================

// UnwrapTimelockCall extracts plain Go types from mcms.TimelockCall for hashing.
func UnwrapTimelockCall(call mcms.TimelockCall) TimelockCall {
	return TimelockCall{
		TargetInstanceAddress: string(call.TargetInstanceAddress),
		FunctionName:          string(call.FunctionName),
		OperationData:         string(call.OperationData),
	}
}

// UnwrapTimelockCalls extracts plain Go types from mcms.TimelockCall slice for hashing.
// This is needed because HashTimelockOpId requires plain Go types, not DAML-wrapped types.
func UnwrapTimelockCalls(calls []mcms.TimelockCall) []TimelockCall {
	result := make([]TimelockCall, len(calls))
	for i, call := range calls {
		result[i] = UnwrapTimelockCall(call)
	}

	return result
}

// NewMCMSEncoder creates an MCMS encoder for the given package ID.
// Usage: encoder := NewMCMSEncoder()
//
//	choice, _ := encoder.ScheduleBatch(params)
//	proposal.AddOperation(instanceAddress, choice.Choice, choice.OperationData)
func NewMCMSEncoder() mcms.MCMSEncoder {
	return mcms.NewContract(fmt.Sprintf("#%s", mcms.PackageName), "MCMS.Main", "MCMS").Encoder()
}

// MustEncodeScheduleBatch encodes ScheduleBatchParams and returns the full EncodedChoice.
// Use choice.Choice for the function name, choice.OperationData for hex data.
func MustEncodeScheduleBatch(t testing.TB, encoder mcms.MCMSEncoder, params mcms.ScheduleBatchParams) *bind.EncodedChoice {
	t.Helper()
	choice, err := encoder.ScheduleBatch(params)
	require.NoError(t, err, "failed to encode ScheduleBatchParams")

	return choice
}

// MustEncodeCancelBatch encodes CancelBatchParams and returns the full EncodedChoice.
// Use choice.Choice for the function name, choice.OperationData for hex data.
func MustEncodeCancelBatch(t testing.TB, encoder mcms.MCMSEncoder, params mcms.CancelBatchParams) *bind.EncodedChoice {
	t.Helper()
	choice, err := encoder.CancelBatch(params)
	require.NoError(t, err, "failed to encode CancelBatchParams")

	return choice
}

// MustEncodeBypasserExecuteBatch encodes BypasserExecuteBatchParams and returns the full EncodedChoice.
// Use choice.Choice for the function name, choice.OperationData for hex data.
func MustEncodeBypasserExecuteBatch(t testing.TB, encoder mcms.MCMSEncoder, params mcms.BypasserExecuteBatchParams) *bind.EncodedChoice {
	t.Helper()
	choice, err := encoder.BypasserExecuteBatch(params)
	require.NoError(t, err, "failed to encode BypasserExecuteBatchParams")

	return choice
}

// MustEncodeSetConfigParams encodes SetConfigParams and returns the full EncodedChoice.
// Use choice.Choice for the function name, choice.OperationData for hex data.
func MustEncodeSetConfigParams(t testing.TB, encoder mcms.MCMSEncoder, params mcms.SetConfigParams) *bind.EncodedChoice {
	t.Helper()
	choice, err := encoder.SetConfigParams(params)
	require.NoError(t, err, "failed to encode SetConfigParams")

	return choice
}

// ToBindingSignerInfo converts local SignerInfo to mcms.SignerInfo binding type.
func ToBindingSignerInfo(si SignerInfo) mcms.SignerInfo {
	return mcms.SignerInfo{
		SignerAddress: types.TEXT(si.SignerAddress),
		SignerIndex:   types.INT64(si.SignerIndex),
		SignerGroup:   types.INT64(si.SignerGroup),
	}
}

// ToBindingSignerInfos converts a slice of local SignerInfo to mcms.SignerInfo binding types.
func ToBindingSignerInfos(signers []SignerInfo) []mcms.SignerInfo {
	result := make([]mcms.SignerInfo, len(signers))
	for i, si := range signers {
		result[i] = ToBindingSignerInfo(si)
	}

	return result
}

// ToINT64Slice converts []int to []types.INT64.
func ToINT64Slice(ints []int) []types.INT64 {
	result := make([]types.INT64, len(ints))
	for i, v := range ints {
		result[i] = types.INT64(v)
	}

	return result
}

// ToBindingSetConfigParams converts local config to mcms.SetConfigParams binding type.
func ToBindingSetConfigParams(signers []SignerInfo, groupQuorums, groupParents []int, clearRoot bool) mcms.SetConfigParams {
	return mcms.SetConfigParams{
		Signers:      ToBindingSignerInfos(signers),
		GroupQuorums: ToINT64Slice(groupQuorums),
		GroupParents: ToINT64Slice(groupParents),
		ClearRoot:    types.BOOL(clearRoot),
	}
}

// MCMSRootMetadata matches Canton RootMetadata
type MCMSRootMetadata struct {
	ChainId              int
	MultisigId           string
	PreOpCount           int
	PostOpCount          int
	OverridePreviousRoot bool
}

// RawSignature matches Canton RawSignature format
type RawSignature struct {
	PublicKey string // Uncompressed public key (130 hex chars: 04 || X || Y)
	R         string // Signature r component (64 hex chars)
	S         string // Signature s component (64 hex chars)
}

// MakeMcmsId creates multisigId from MCMS instanceAddress and role (e.g., "mcms-001@partyId-proposer")
func MakeMcmsId(instanceAddress string, role MCMSRole) string {
	return instanceAddress + "-" + role.String()
}

// MCMSSigner wraps a private key with MCMS-specific functionality
type MCMSSigner struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string // lowercase, no 0x prefix
	PublicKey  string // uncompressed, 130 hex chars with 04 prefix
}

// NewMCMSSigner creates a new signer with generated key
func NewMCMSSigner() (*MCMSSigner, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return NewMCMSSignerFromKey(privateKey), nil
}

// NewMCMSSignerFromSeed creates a signer with deterministically generated key from a seed.
// This is useful for generating consistent test vectors across multiple test runs.
func NewMCMSSignerFromSeed(seed int64) (*MCMSSigner, error) {
	src := mrand.NewSource(seed)
	rng := mrand.New(src) //nolint:gosec // Deterministic key for testing

	keyBytes := make([]byte, 32)
	for i := range keyBytes {
		keyBytes[i] = byte(rng.Intn(256)) //nolint:gosec // intentional overflow
	}

	privateKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create key from seed: %w", err)
	}

	return NewMCMSSignerFromKey(privateKey), nil
}

// NewMCMSSignerFromKey creates a signer from existing private key
func NewMCMSSignerFromKey(privateKey *ecdsa.PrivateKey) *MCMSSigner {
	pubKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &MCMSSigner{
		PrivateKey: privateKey,
		Address:    strings.ToLower(address.Hex()[2:]), // Remove 0x, lowercase
		PublicKey:  hex.EncodeToString(pubKeyBytes),    // Includes 04 prefix
	}
}

// ToSignerInfo creates SignerInfo for config
func (s *MCMSSigner) ToSignerInfo(index, group int) SignerInfo {
	return SignerInfo{
		SignerAddress: s.Address,
		SignerIndex:   index,
		SignerGroup:   group,
	}
}

// Sign signs a message hash and returns RawSignature
func (s *MCMSSigner) Sign(messageHash []byte) (*RawSignature, error) {
	sig, err := crypto.Sign(messageHash, s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	// sig is 65 bytes: R (32) || S (32) || V (1)
	r := hex.EncodeToString(sig[:32])
	sigS := hex.EncodeToString(sig[32:64])

	return &RawSignature{
		PublicKey: s.PublicKey,
		R:         r,
		S:         sigS,
	}, nil
}

// SortSignersByAddress sorts signers by address (required by MCMS)
func SortSignersByAddress(signers []*MCMSSigner) []*MCMSSigner {
	sorted := make([]*MCMSSigner, len(signers))
	copy(sorted, signers)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Address < sorted[j].Address
	})

	return sorted
}

// ===========================================================================
// HASH COMPUTATION (matching Canton Crypto.daml)
// ===========================================================================

// ComputeSignedHash computes the signed hash for SetRoot
// Matches Canton's computeSignedHashNative:
//
//	innerHash = keccak256(root || timeToHex(validUntil))
//	signedHash = keccak256("\x19Ethereum Signed Message:\n32" || innerHash)
func ComputeSignedHash(root string, validUntil time.Time) string {
	// Inner data: root || validUntil as hex
	validUntilHex := TimeToHex(validUntil)
	innerData, _ := hex.DecodeString(root + validUntilHex)
	innerHash := crypto.Keccak256(innerData)

	// EIP-191 prefix: "\x19Ethereum Signed Message:\n32"
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	prefixedData := append(prefix, innerHash...)

	return hex.EncodeToString(crypto.Keccak256(prefixedData))
}

// TimeToHex converts time to 32-byte hex (matching ABI encoding of uint32)
// Encodes Unix timestamp in seconds, left-padded to 32 bytes
func TimeToHex(t time.Time) string {
	unixSeconds := t.Unix()
	return PadLeft32(fmt.Sprintf("%x", unixSeconds))
}

// HashOpLeaf hashes an operation to get its Merkle leaf
// Matches Canton's hashOpLeafNative with domain separator and length-prefixed encoding for variable-length fields
func HashOpLeaf(op MCMSOp) string {
	multisigIdHex := TextToHex(op.MultisigId)
	targetAddressHex := TextToHex(op.TargetInstanceAddress)
	functionNameHex := TextToHex(op.FunctionName)

	// Domain separator followed by length prefixes (byte count for all fields)
	encoded := opLeafDomainSeparator +
		PadLeft32(IntToHex(op.ChainId)) +
		PadLeft32(IntToHex(len(op.MultisigId))) + // Length prefix for multisigId (UTF-8 byte count)
		multisigIdHex +
		PadLeft32(IntToHex(op.Nonce)) +
		PadLeft32(IntToHex(len(op.TargetInstanceAddress))) + // Length prefix for targetInstanceAddress (UTF-8 byte count)
		targetAddressHex +
		PadLeft32(IntToHex(len(op.FunctionName))) + // Length prefix for functionName (UTF-8 byte count)
		functionNameHex +
		PadLeft32(IntToHex(len(op.OperationData)/2)) + // Length prefix for operationData (byte count)
		op.OperationData

	data, _ := hex.DecodeString(encoded)

	return hex.EncodeToString(crypto.Keccak256(data))
}

// HashMetadataLeaf hashes metadata to get its Merkle leaf
// Matches Canton's hashMetadataLeafNative with domain separator and length-prefixed encoding
func HashMetadataLeaf(meta MCMSRootMetadata) string {
	overrideFlag := "00"
	if meta.OverridePreviousRoot {
		overrideFlag = "01"
	}

	multisigIdHex := TextToHex(meta.MultisigId)

	// Domain separator followed by length prefixes (byte count for all fields)
	encoded := metadataLeafDomainSeparator +
		PadLeft32(IntToHex(meta.ChainId)) +
		PadLeft32(IntToHex(len(meta.MultisigId))) + // Length prefix for multisigId (UTF-8 byte count)
		multisigIdHex +
		PadLeft32(IntToHex(meta.PreOpCount)) +
		PadLeft32(IntToHex(meta.PostOpCount)) +
		overrideFlag

	data, _ := hex.DecodeString(encoded)

	return hex.EncodeToString(crypto.Keccak256(data))
}

// HashTimelockOpId computes the operation ID for timelock operations
// Matches Canton's hashTimelockOpId with length-prefixed encoding to prevent collision attacks
// Note: predecessor and salt are BytesHex (already hex-encoded), used directly without conversion.
func HashTimelockOpId(calls []TimelockCall, predecessor, salt string) string {
	var sb strings.Builder

	// Length prefix for calls array (32-byte padded)
	sb.WriteString(PadLeft32(IntToHex(len(calls))))

	// Encode each call with length prefixes
	for _, call := range calls {
		// Length prefix for targetInstanceAddress (UTF-8 byte count)
		sb.WriteString(PadLeft32(IntToHex(len(call.TargetInstanceAddress))))
		sb.WriteString(TextToHex(call.TargetInstanceAddress))
		// Length prefix for functionName (UTF-8 byte count)
		sb.WriteString(PadLeft32(IntToHex(len(call.FunctionName))))
		sb.WriteString(TextToHex(call.FunctionName))
		// Length prefix for operationData (byte count = hex length / 2)
		opData := EncodeOperationDataForHash(call.OperationData)
		sb.WriteString(PadLeft32(IntToHex(len(opData) / 2)))
		sb.WriteString(opData)
	}

	// Length prefix for predecessor (byte count = hex length / 2)
	sb.WriteString(PadLeft32(IntToHex(len(predecessor) / 2)))
	sb.WriteString(predecessor)

	// Length prefix for salt (byte count = hex length / 2)
	sb.WriteString(PadLeft32(IntToHex(len(salt) / 2)))
	sb.WriteString(salt)

	data, err := hex.DecodeString(sb.String())
	if err != nil {
		panic(fmt.Sprintf("HashTimelockOpId: invalid hex encoding: %v", err))
	}

	return hex.EncodeToString(crypto.Keccak256(data))
}

// EncodeOperationDataForHash matches the on-chain MCMS.Crypto.encodeOperationData:
// - If operationData is valid hex (even length, hex digits only), treat it as raw bytes (already hex)
// - Otherwise, treat it as text and hex-encode it (UTF-8)
func EncodeOperationDataForHash(operationData string) string {
	if isValidHex(operationData) {
		return operationData
	}

	return TextToHex(operationData)
}

func isValidHex(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return false
		}
	}

	return true
}

// ===========================================================================
// MERKLE TREE (matching Canton Crypto.daml and MCMS tree.go)
// ===========================================================================

// MerkleTree represents a Merkle tree with sorted hash pairs
type MerkleTree struct {
	Root      string
	Layers    [][]string     // layers[0] = leaves, layers[n] = root
	LeafIndex map[string]int // Map from leaf hash to index in sorted leaves
}

// NewMerkleTree builds a Merkle tree from leaf hashes
// Uses sorted hash pairs (OpenZeppelin style)
func NewMerkleTree(leaves []string) *MerkleTree {
	if len(leaves) == 0 {
		return &MerkleTree{Root: strings.Repeat("0", 64)}
	}

	// Sort leaves lexicographically
	sortedLeaves := make([]string, len(leaves))
	copy(sortedLeaves, leaves)
	sort.Strings(sortedLeaves)

	// Build leaf index map
	leafIndex := make(map[string]int)
	for i, leaf := range sortedLeaves {
		leafIndex[leaf] = i
	}

	// Build layers - store the working layers with padding
	layers := [][]string{}
	currentLayer := make([]string, len(sortedLeaves))
	copy(currentLayer, sortedLeaves)
	layers = append(layers, currentLayer)

	for len(currentLayer) > 1 {
		// Create working copy with padding if needed
		workingLayer := make([]string, len(currentLayer))
		copy(workingLayer, currentLayer)
		if len(workingLayer)%2 != 0 {
			workingLayer = append(workingLayer, workingLayer[len(workingLayer)-1]) //nolint:makezero // This is intentional padding
		}

		var nextLayer []string
		for i := 0; i < len(workingLayer); i += 2 {
			hash := HashPair(workingLayer[i], workingLayer[i+1])
			nextLayer = append(nextLayer, hash)
		}
		layers = append(layers, nextLayer)
		currentLayer = nextLayer
	}

	return &MerkleTree{
		Root:      currentLayer[0],
		Layers:    layers,
		LeafIndex: leafIndex,
	}
}

// GetProof returns the Merkle proof for a leaf
func (t *MerkleTree) GetProof(leafHash string) ([]string, error) {
	idx, ok := t.LeafIndex[leafHash]
	if !ok {
		return nil, fmt.Errorf("leaf not found in tree: %s", leafHash)
	}

	var proof []string
	currentIdx := idx

	for layerNum := range len(t.Layers) - 1 {
		layer := t.Layers[layerNum]
		layerLen := len(layer)

		// Handle odd-length layers by conceptually duplicating last element
		var sibling string
		siblingIdx := currentIdx ^ 1

		if siblingIdx < layerLen {
			sibling = layer[siblingIdx]
		} else {
			// Odd layer, sibling is duplicate of current (which is last element)
			sibling = layer[currentIdx]
		}

		proof = append(proof, sibling)

		// Move to parent index
		currentIdx = currentIdx / 2
	}

	return proof, nil
}

// HashPair hashes two nodes in sorted order (OpenZeppelin style)
// Matches Canton's hashPair function
func HashPair(a, b string) string {
	var first, second string
	if a < b {
		first, second = a, b
	} else {
		first, second = b, a
	}

	data, _ := hex.DecodeString(first + second)

	return hex.EncodeToString(crypto.Keccak256(data))
}

// VerifyMerkleProof verifies a Merkle proof on the Go side
// Matches Canton's verifyMerkleProof
func VerifyMerkleProof(root, leaf string, proof []string) bool {
	computedHash := leaf
	for _, proofElement := range proof {
		computedHash = HashPair(computedHash, proofElement)
	}

	return computedHash == root
}

// ===========================================================================
// SIGNING HELPERS
// ===========================================================================

// SignMCMSRoot signs a root with multiple signers
// Returns signatures sorted by signer address (required by MCMS)
func SignMCMSRoot(root string, validUntil time.Time, signers []*MCMSSigner) ([]RawSignature, error) {
	signedHash := ComputeSignedHash(root, validUntil)
	hashBytes, err := hex.DecodeString(signedHash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed hash: %w", err)
	}

	// Sort signers by address
	sortedSigners := SortSignersByAddress(signers)

	var signatures []RawSignature
	for _, signer := range sortedSigners {
		sig, err := signer.Sign(hashBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to sign with signer %s: %w", signer.Address, err)
		}
		signatures = append(signatures, *sig)
	}

	return signatures, nil
}

// ===========================================================================
// HEX ENCODING HELPERS (matching Canton Crypto.daml)
// ===========================================================================

// PadLeft32 pads hex string to 64 chars (32 bytes)
func PadLeft32(hexStr string) string {
	if len(hexStr) >= 64 {
		return hexStr[:64]
	}

	return strings.Repeat("0", 64-len(hexStr)) + hexStr
}

// IntToHex converts int to hex string (without padding)
func IntToHex(n int) string {
	if n == 0 {
		return "0"
	}

	return fmt.Sprintf("%x", n)
}

// TextToHex converts a string to hex via UTF-8 encoding.
// Matches Daml's DA.Crypto.Text.toHex function.
func TextToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

// EncodeMinDelay encodes a delay value (in seconds) for UpdateMinDelay operationData.
// Uses 8-byte big-endian encoding (16 hex chars)
// Example: 120 → "0000000000000078"
func EncodeMinDelay(seconds int) string {
	return fmt.Sprintf("%016x", seconds)
}

// ===========================================================================
// CONFIG HELPERS
// ===========================================================================

const NumGroups = 32

// NewEmptyConfig creates an empty MCMS config
func NewEmptyConfig() MCMSConfig {
	return MCMSConfig{
		Signers:      []SignerInfo{},
		GroupQuorums: make([]int, NumGroups),
		GroupParents: make([]int, NumGroups),
	}
}

// New2of3Config creates a simple 2-of-3 multisig config
// All signers in group 0 (root group)
func New2of3Config(signers []*MCMSSigner) MCMSConfig {
	if len(signers) < 3 {
		panic("need at least 3 signers for 2-of-3")
	}

	// Sort signers by address for consistent ordering
	sorted := SortSignersByAddress(signers)

	signerInfos := make([]SignerInfo, 3)
	for i := range 3 {
		signerInfos[i] = sorted[i].ToSignerInfo(i, 0) // All in group 0
	}

	groupQuorums := make([]int, NumGroups)
	groupQuorums[0] = 2 // Root group needs 2

	groupParents := make([]int, NumGroups)
	// All groups point to root (0)

	return MCMSConfig{
		Signers:      signerInfos,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
	}
}

// ===========================================================================
// PROPOSAL BUILDER
// ===========================================================================

// MCMSProposal helps build a complete MCMS proposal with operations and metadata
type MCMSProposal struct {
	ChainId    int
	MultisigId string
	Operations []MCMSOp
	Metadata   MCMSRootMetadata
	Tree       *MerkleTree
}

// NewMCMSProposal creates a new proposal builder
func NewMCMSProposal(chainId int, multisigId string, preOpCount int, overridePreviousRoot bool) *MCMSProposal {
	return &MCMSProposal{
		ChainId:    chainId,
		MultisigId: multisigId,
		Operations: []MCMSOp{},
		Metadata: MCMSRootMetadata{
			ChainId:              chainId,
			MultisigId:           multisigId,
			PreOpCount:           preOpCount,
			PostOpCount:          preOpCount,
			OverridePreviousRoot: overridePreviousRoot,
		},
	}
}

// AddOperation adds an operation to the proposal
func (p *MCMSProposal) AddOperation(targetInstanceAddress, functionName, operationData string) *MCMSProposal {
	op := MCMSOp{
		ChainId:               p.ChainId,
		MultisigId:            p.MultisigId,
		Nonce:                 p.Metadata.PostOpCount,
		TargetInstanceAddress: targetInstanceAddress,
		FunctionName:          functionName,
		OperationData:         operationData,
	}
	p.Operations = append(p.Operations, op)
	p.Metadata.PostOpCount++

	return p
}

// SetOverride sets whether to override previous root
func (p *MCMSProposal) SetOverride(override bool) *MCMSProposal {
	p.Metadata.OverridePreviousRoot = override
	return p
}

// Build builds the Merkle tree for the proposal
func (p *MCMSProposal) Build() *MCMSProposal {
	leaves := make([]string, 0, len(p.Operations)+1)

	// Add metadata leaf
	metadataLeaf := HashMetadataLeaf(p.Metadata)
	leaves = append(leaves, metadataLeaf)

	// Add operation leaves
	for _, op := range p.Operations {
		opLeaf := HashOpLeaf(op)
		leaves = append(leaves, opLeaf)
	}

	p.Tree = NewMerkleTree(leaves)

	return p
}

// GetRoot returns the Merkle root
func (p *MCMSProposal) GetRoot() string {
	if p.Tree == nil {
		p.Build()
	}

	return p.Tree.Root
}

// GetMetadataProof returns the Merkle proof for metadata
func (p *MCMSProposal) GetMetadataProof() ([]string, error) {
	if p.Tree == nil {
		p.Build()
	}
	metadataLeaf := HashMetadataLeaf(p.Metadata)

	return p.Tree.GetProof(metadataLeaf)
}

// GetOpProof returns the Merkle proof for an operation by index
func (p *MCMSProposal) GetOpProof(opIndex int) ([]string, error) {
	if p.Tree == nil {
		p.Build()
	}
	if opIndex >= len(p.Operations) {
		return nil, fmt.Errorf("operation index %d out of range", opIndex)
	}
	opLeaf := HashOpLeaf(p.Operations[opIndex])

	return p.Tree.GetProof(opLeaf)
}

// Sign signs the proposal with the given signers
func (p *MCMSProposal) Sign(validUntil time.Time, signers []*MCMSSigner) ([]RawSignature, error) {
	return SignMCMSRoot(p.GetRoot(), validUntil, signers)
}

// ===========================================================================
// DEBUG HELPERS
// ===========================================================================

// PublicKeyToAddress converts public key to Ethereum address
// Matches Canton's publicKeyToAddress
func PublicKeyToAddress(pubKeyHex string) string {
	// Remove 04 prefix if present
	keyData := pubKeyHex
	if len(keyData) >= 2 && keyData[:2] == "04" {
		keyData = keyData[2:]
	}

	data, _ := hex.DecodeString(keyData)
	hash := crypto.Keccak256(data)

	// Last 20 bytes
	return hex.EncodeToString(hash[12:])
}

// FormatSignerInfo formats signer info for debugging
func FormatSignerInfo(s *MCMSSigner) string {
	return fmt.Sprintf("Address: %s, PubKey: %s...", s.Address, s.PublicKey[:20])
}

// VerifySignature verifies that a signature is valid (for testing)
func VerifySignature(messageHash []byte, sig RawSignature) (string, error) {
	pubKeyBytes, err := hex.DecodeString(sig.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	rBytes, _ := hex.DecodeString(sig.R)
	sBytes, _ := hex.DecodeString(sig.S)

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	if !ecdsa.Verify(pubKey, messageHash, r, s) {
		return "", fmt.Errorf("signature verification failed")
	}

	address := crypto.PubkeyToAddress(*pubKey)

	return strings.ToLower(address.Hex()[2:]), nil
}

// Common returns the common.Hash from a hex string
func HexToHash(hexStr string) common.Hash {
	return common.HexToHash(hexStr)
}

// ===========================================================================
// OPERATION DATA ENCODING (like Aptos BCS)
// ===========================================================================

// wireToHex encodes raw MCMS wire bytes (go-daml HexCodec.Marshal / MarshalHex output) as hex ASCII.
func wireToHex(wire string) string {
	return hex.EncodeToString([]byte(wire))
}

// MustEncodeBlockedFunction encodes a BlockedFunction to hex bytes using binding's MarshalHex.
func MustEncodeBlockedFunction(t testing.TB, bf mcms.BlockedFunction) string {
	t.Helper()
	wire, err := bf.MarshalHex()
	require.NoError(t, err, "failed to encode BlockedFunction")

	return wireToHex(wire)
}

// EncodeSelfDispatchSetConfig encodes role + SetConfigParams for self-dispatch set_config
// Format matches Canton MCMS.Codec.encodeSelfDispatchSetConfig:
//
//	uint8(roleToInt(role)) <> encodeSetConfigParams(params)
func EncodeSelfDispatchSetConfig(role MCMSRole, params mcms.SetConfigParams) (string, error) {
	wire, err := params.MarshalHex()
	if err != nil {
		return "", fmt.Errorf("MarshalHex failed: %w", err)
	}

	// uint8(role) <> encodeSetConfigParams(params), then hex for TimelockCall hex:"bytes16".
	payload := append([]byte{byte(role)}, []byte(wire)...) //nolint:gosec // MCMSRole is a controlled value
	return hex.EncodeToString(payload), nil
}

// EncodePruneSeenHashesParams encodes the role byte for a PruneSeenHashes self-dispatch call.
// PruneTimelockTimestamps takes no parameters — pass "" as operationData.
func EncodePruneSeenHashesParams(role MCMSRole) string {
	return hex.EncodeToString([]byte{byte(role)}) //nolint:gosec // MCMSole is a controlled value
}

// ===========================================================================
// TIMELOCK OPERATION DATA ENCODING
// ===========================================================================

// ScheduleBatchParams matches Canton MCMS.Codec.ScheduleBatchParams
type ScheduleBatchParams struct {
	Calls       []TimelockCall
	Predecessor string // hex hash (64 chars) or empty for no predecessor
	Salt        string // arbitrary text for uniqueness
	DelaySecs   int    // additional delay in seconds
}

// BypasserExecuteBatchParams matches Canton MCMS.Codec.BypasserExecuteBatchParams
type BypasserExecuteBatchParams struct {
	Calls []TimelockCall
}

// CancelBatchParams matches Canton MCMS.Codec.CancelBatchParams
type CancelBatchParams struct {
	OpId string // operation ID to cancel (hex hash)
}

// encodeText encodes a string as length-prefixed bytes
// Format: uint8(len) + asciiBytes
// Matches Canton MCMS.Codec.encodeText
func encodeText(s string) []byte {
	textBytes := []byte(s)
	buf := make([]byte, 0, 1+len(textBytes))
	buf = append(buf, byte(len(textBytes))) //nolint:gosec
	buf = append(buf, textBytes...)

	return buf
}

// encodeTextUint16 encodes a hex string with uint16 byte-count prefix
// Format: uint16(byteCount) + hexString (where byteCount = len(hexString) / 2)
// Matches Canton MCMS.Codec usage for operationData, predecessor, salt
func encodeTextUint16(s string) []byte {
	textBytes := []byte(s)
	byteCount := len(textBytes) / 2 // Hex string: 2 chars per byte
	buf := make([]byte, 0, 2+len(textBytes))
	//nolint:gosec
	buf = append(buf, byte(byteCount>>8), byte(byteCount&0xff)) // uint16 big-endian
	buf = append(buf, textBytes...)

	return buf
}

// encodeTimelockCall encodes a TimelockCall
// Format: encodeText(targetInstanceAddress) + encodeText(functionName) + encodeTextUint16(operationData)
// Matches Canton MCMS.Codec.encodeTimelockCall
func encodeTimelockCall(call TimelockCall) []byte {
	// Calculate size: 2 uint8 length bytes + 1 uint16 length (2 bytes) + string contents
	size := 4 + len(call.TargetInstanceAddress) + len(call.FunctionName) + len(call.OperationData)
	buf := make([]byte, 0, size)
	buf = append(buf, encodeText(call.TargetInstanceAddress)...)
	buf = append(buf, encodeText(call.FunctionName)...)
	buf = append(buf, encodeTextUint16(call.OperationData)...) // uint16 prefix for operationData

	return buf
}

// EncodeScheduleBatchParams encodes ScheduleBatchParams to hex bytes
// Format: numCalls (1 byte) + calls + encodeTextUint16(predecessor) + encodeTextUint16(salt) + int64(delaySecs)
// Matches Canton MCMS.Codec.encodeScheduleBatchParams
func EncodeScheduleBatchParams(params ScheduleBatchParams) string {
	// Estimate size: 1 (numCalls) + calls + predecessor (2+len) + salt (2+len) + 8 (delay as int64)
	estimatedSize := 1 + len(params.Predecessor) + 2 + len(params.Salt) + 2 + 8
	for _, call := range params.Calls {
		estimatedSize += 4 + len(call.TargetInstanceAddress) + len(call.FunctionName) + len(call.OperationData)
	}
	buf := make([]byte, 0, estimatedSize)

	// Encode calls list: uint8(numCalls) + encoded calls
	buf = append(buf, byte(len(params.Calls))) //nolint:gosec
	for _, call := range params.Calls {
		buf = append(buf, encodeTimelockCall(call)...)
	}

	// Encode predecessor with uint16 byte-count prefix (matches Daml Codec.daml line 379)
	buf = append(buf, encodeTextUint16(params.Predecessor)...)

	// Encode salt with uint16 byte-count prefix (matches Daml Codec.daml line 380)
	buf = append(buf, encodeTextUint16(params.Salt)...)

	// Encode delaySecs (8 bytes, big-endian INT64 to match DAML INT type)
	delayBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(delayBytes, uint64(params.DelaySecs)) //nolint:gosec
	buf = append(buf, delayBytes...)

	return hex.EncodeToString(buf)
}

// EncodeBypasserExecuteBatchParams encodes BypasserExecuteBatchParams to hex bytes
// Format: numCalls (1 byte) + encoded calls
// Matches Canton MCMS.Codec.encodeBypasserExecuteBatchParams
func EncodeBypasserExecuteBatchParams(params BypasserExecuteBatchParams) string {
	// Estimate size: 1 (numCalls) + calls
	estimatedSize := 1
	for _, call := range params.Calls {
		estimatedSize += 3 + len(call.TargetInstanceAddress) + len(call.FunctionName) + len(call.OperationData)
	}
	buf := make([]byte, 0, estimatedSize)

	// Encode calls list: uint8(numCalls) + encoded calls
	buf = append(buf, byte(len(params.Calls))) //nolint:gosec
	for _, call := range params.Calls {
		buf = append(buf, encodeTimelockCall(call)...)
	}

	return hex.EncodeToString(buf)
}

// EncodeCancelBatchParams encodes CancelBatchParams to hex bytes
// Format: encodeText(opId)
// Matches Canton MCMS.Codec.encodeCancelBatchParams
func EncodeCancelBatchParams(params CancelBatchParams) string {
	return hex.EncodeToString(encodeText(params.OpId))
}
