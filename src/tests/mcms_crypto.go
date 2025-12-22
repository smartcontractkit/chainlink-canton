package tests

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ===========================================================================
// TYPES (matching Canton MCMSPoc.Types)
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
	ChainId       int
	MultisigId    string
	Nonce         int
	TargetAddress string
	FunctionName  string
	OperationData string // hex encoded
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

// MakeMcmsId creates mcmsId from base ID and role (e.g., "mcms-001-proposer")
func MakeMcmsId(baseId string, role MCMSRole) string {
	return baseId + "-" + role.String()
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

// TimeToHex converts time to 32-byte hex (matching Canton's placeholder)
// Canton currently uses padLeft32 "0" - we match that for compatibility
func TimeToHex(t time.Time) string {
	// For now, match Canton's placeholder implementation
	// TODO: Implement proper timestamp encoding when Canton updates
	return strings.Repeat("0", 64)
}

// TimeToHexUnix converts time to 32-byte hex with Unix timestamp
// Use this if Canton updates to use actual timestamps
func TimeToHexUnix(t time.Time) string {
	timestamp := uint64(t.Unix())
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, timestamp)
	// Pad to 32 bytes (64 hex chars)
	return strings.Repeat("0", 48) + hex.EncodeToString(buf)
}

// HashOpLeaf hashes an operation to get its Merkle leaf
// Matches Canton's hashOpLeafNative
func HashOpLeaf(op MCMSOp) string {
	encoded := PadLeft32(IntToHex(op.ChainId)) +
		AsciiToHex(op.MultisigId) +
		PadLeft32(IntToHex(op.Nonce)) +
		AsciiToHex(op.TargetAddress) +
		AsciiToHex(op.FunctionName) +
		op.OperationData

	data, _ := hex.DecodeString(encoded)
	return hex.EncodeToString(crypto.Keccak256(data))
}

// HashMetadataLeaf hashes metadata to get its Merkle leaf
// Matches Canton's hashMetadataLeafNative
func HashMetadataLeaf(meta MCMSRootMetadata) string {
	overrideFlag := "00"
	if meta.OverridePreviousRoot {
		overrideFlag = "01"
	}

	encoded := PadLeft32(IntToHex(meta.ChainId)) +
		AsciiToHex(meta.MultisigId) +
		PadLeft32(IntToHex(meta.PreOpCount)) +
		PadLeft32(IntToHex(meta.PostOpCount)) +
		overrideFlag

	data, _ := hex.DecodeString(encoded)
	return hex.EncodeToString(crypto.Keccak256(data))
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
			workingLayer = append(workingLayer, workingLayer[len(workingLayer)-1])
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

	for layerNum := 0; layerNum < len(t.Layers)-1; layerNum++ {
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

// AsciiToHex converts ASCII string to hex
func AsciiToHex(s string) string {
	return hex.EncodeToString([]byte(s))
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
	for i := 0; i < 3; i++ {
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
// CANTON VALUE BUILDERS
// ===========================================================================

// These helpers create Canton API values for MCMS contracts

// BuildSignerInfoValue creates Canton Value for SignerInfo
func BuildSignerInfoValue(si SignerInfo) map[string]interface{} {
	return map[string]interface{}{
		"signerAddress": si.SignerAddress,
		"signerIndex":   si.SignerIndex,
		"signerGroup":   si.SignerGroup,
	}
}

// BuildRawSignatureValue creates Canton Value for RawSignature
func BuildRawSignatureValue(sig RawSignature) map[string]interface{} {
	return map[string]interface{}{
		"publicKey": sig.PublicKey,
		"r":         sig.R,
		"s":         sig.S,
	}
}

// BuildOpValue creates Canton Value for Op
func BuildOpValue(op MCMSOp) map[string]interface{} {
	return map[string]interface{}{
		"chainId":       op.ChainId,
		"multisigId":    op.MultisigId,
		"nonce":         op.Nonce,
		"targetAddress": op.TargetAddress,
		"functionName":  op.FunctionName,
		"operationData": op.OperationData,
	}
}

// BuildMetadataValue creates Canton Value for RootMetadata
func BuildMetadataValue(meta MCMSRootMetadata) map[string]interface{} {
	return map[string]interface{}{
		"chainId":              meta.ChainId,
		"multisigId":           meta.MultisigId,
		"preOpCount":           meta.PreOpCount,
		"postOpCount":          meta.PostOpCount,
		"overridePreviousRoot": meta.OverridePreviousRoot,
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
func (p *MCMSProposal) AddOperation(targetAddress, functionName, operationData string) *MCMSProposal {
	op := MCMSOp{
		ChainId:       p.ChainId,
		MultisigId:    p.MultisigId,
		Nonce:         p.Metadata.PostOpCount,
		TargetAddress: targetAddress,
		FunctionName:  functionName,
		OperationData: operationData,
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
	var leaves []string

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

// SetConfigParams matches Canton SetConfigParams for encoding
type SetConfigParams struct {
	Signers      []SignerInfo
	GroupQuorums []int // Length 32
	GroupParents []int // Length 32
	ClearRoot    bool
}

// EncodeSetConfigParams encodes SetConfigParams to hex bytes
// Format matches Canton MCMSPoc.Codec.encodeSetConfigParams
func EncodeSetConfigParams(params SetConfigParams) string {
	var buf []byte

	// Encode signers list
	buf = append(buf, byte(len(params.Signers))) // numSigners (1 byte)
	for _, signer := range params.Signers {
		// Address: length + hex bytes
		addrBytes, _ := hex.DecodeString(signer.SignerAddress)
		buf = append(buf, byte(len(addrBytes))) // addressLen (1 byte)
		buf = append(buf, addrBytes...)         // address bytes

		// SignerIndex (4 bytes, big-endian)
		indexBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(indexBytes, uint32(signer.SignerIndex))
		buf = append(buf, indexBytes...)

		// SignerGroup (4 bytes, big-endian)
		groupBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(groupBytes, uint32(signer.SignerGroup))
		buf = append(buf, groupBytes...)
	}

	// Encode group quorums
	buf = append(buf, byte(len(params.GroupQuorums))) // numQuorums (1 byte)
	for _, quorum := range params.GroupQuorums {
		quorumBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(quorumBytes, uint32(quorum))
		buf = append(buf, quorumBytes...)
	}

	// Encode group parents
	buf = append(buf, byte(len(params.GroupParents))) // numParents (1 byte)
	for _, parent := range params.GroupParents {
		parentBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(parentBytes, uint32(parent))
		buf = append(buf, parentBytes...)
	}

	// Encode clearRoot (1 byte)
	if params.ClearRoot {
		buf = append(buf, 0x01)
	} else {
		buf = append(buf, 0x00)
	}

	return hex.EncodeToString(buf)
}

// DecodeSetConfigParams decodes SetConfigParams from hex bytes
// Format matches Canton MCMSPoc.Codec.decodeSetConfigParams
func DecodeSetConfigParams(hexData string) (*SetConfigParams, error) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	if len(data) < 1 {
		return nil, fmt.Errorf("data too short: need at least 1 byte for signer count")
	}

	offset := 0
	result := &SetConfigParams{}

	// Decode signers list
	numSigners := int(data[offset])
	offset++

	result.Signers = make([]SignerInfo, numSigners)
	for i := 0; i < numSigners; i++ {
		if offset >= len(data) {
			return nil, fmt.Errorf("data truncated at signer %d address length", i)
		}

		// Address length (in bytes)
		addrLen := int(data[offset])
		offset++

		if offset+addrLen > len(data) {
			return nil, fmt.Errorf("data truncated at signer %d address", i)
		}

		// Address bytes -> hex string
		addrBytes := data[offset : offset+addrLen]
		result.Signers[i].SignerAddress = hex.EncodeToString(addrBytes)
		offset += addrLen

		// SignerIndex (4 bytes, big-endian)
		if offset+4 > len(data) {
			return nil, fmt.Errorf("data truncated at signer %d index", i)
		}
		result.Signers[i].SignerIndex = int(binary.BigEndian.Uint32(data[offset : offset+4]))
		offset += 4

		// SignerGroup (4 bytes, big-endian)
		if offset+4 > len(data) {
			return nil, fmt.Errorf("data truncated at signer %d group", i)
		}
		result.Signers[i].SignerGroup = int(binary.BigEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}

	// Decode group quorums
	if offset >= len(data) {
		return nil, fmt.Errorf("data truncated at quorums count")
	}
	numQuorums := int(data[offset])
	offset++

	result.GroupQuorums = make([]int, numQuorums)
	for i := 0; i < numQuorums; i++ {
		if offset+4 > len(data) {
			return nil, fmt.Errorf("data truncated at quorum %d", i)
		}
		result.GroupQuorums[i] = int(binary.BigEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}

	// Decode group parents
	if offset >= len(data) {
		return nil, fmt.Errorf("data truncated at parents count")
	}
	numParents := int(data[offset])
	offset++

	result.GroupParents = make([]int, numParents)
	for i := 0; i < numParents; i++ {
		if offset+4 > len(data) {
			return nil, fmt.Errorf("data truncated at parent %d", i)
		}
		result.GroupParents[i] = int(binary.BigEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}

	// Decode clearRoot
	if offset >= len(data) {
		return nil, fmt.Errorf("data truncated at clearRoot")
	}
	result.ClearRoot = data[offset] == 0x01
	offset++

	// Verify all data was consumed
	if offset != len(data) {
		return nil, fmt.Errorf("trailing bytes after decoding: %d extra bytes", len(data)-offset)
	}

	return result, nil
}
