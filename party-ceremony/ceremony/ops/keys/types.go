package keys

// CreateMemberKeyInput is the input to [CreateMemberKeyOp].
type CreateMemberKeyInput struct {
	NamespaceName string `json:"namespace_name"`
	ParticipantID string `json:"participant_id"`
}

// CreateMemberKeyOutput is the output of [CreateMemberKeyOp].
type CreateMemberKeyOutput struct {
	ParticipantID        string `json:"participant_id"`
	ParticipantUID       string `json:"participant_uid"`
	NamespaceFingerprint string `json:"namespace_fingerprint"`
	// SigningKeyB64 is the base64-encoded proto-marshalled SigningPublicKey
	// returned by GenerateSigningKey (NamespaceOnly usage).
	SigningKeyB64 string `json:"signing_key_b64"`
	// DamlKeyB64 is the base64-encoded proto-marshalled PROTOCOL SigningPublicKey
	// (SigningKeyUsage_PROTOCOL). This key is registered in
	// PartyToParticipant.PartySigningKeys and is used to authorise DAML
	// transactions submitted via InteractiveSubmissionService.
	DamlKeyB64         string `json:"daml_key_b64"`
	DamlKeyFingerprint string `json:"daml_key_fingerprint"`
}

// ResolveProtocolSigningKeyInput is the input to [ResolveProtocolSigningKeyOp].
type ResolveProtocolSigningKeyInput struct {
	// ParticipantID identifies the participant whose local protocol key should
	// be resolved. Only this participant's node can execute this operation.
	ParticipantID string `json:"participant_id"`

	// KnownSigningKeysB64 is the current PartyToParticipant signing-key set
	// from topology. The local vault is cross-referenced with this list.
	KnownSigningKeysB64 []string `json:"known_signing_keys_b64"`
}

// ResolveProtocolSigningKeyOutput is the output of [ResolveProtocolSigningKeyOp].
type ResolveProtocolSigningKeyOutput struct {
	ParticipantID      string `json:"participant_id"`
	KeyB64             string `json:"key_b64"`
	KeyFingerprint     string `json:"key_fingerprint"`
	KnownSigningKeyIdx int    `json:"known_signing_key_idx"`
}

// GenerateRotatedKeyInput is the input to [GenerateRotatedKeyOp].
type GenerateRotatedKeyInput struct {
	// ParticipantID identifies the target participant. Only this participant's
	// node can execute this operation (UID check enforced).
	ParticipantID string `json:"participant_id"`

	// SynchronizerID is the Canton synchronizer to query when discovering the
	// existing key name from the vault.
	SynchronizerID string `json:"synchronizer_id"`

	// DNSOwners is the list of current namespace fingerprints from the
	// DecentralizedNamespaceDefinition. Used to look up the existing
	// namespace key name from the vault via GetNamespaceKeyName.
	DNSOwners []string `json:"dns_owners"`

	// RotateNamespaceKey controls whether a new namespace key is generated.
	RotateNamespaceKey bool `json:"rotate_namespace_key"`

	// RotateDamlKey controls whether a new DAML key is generated.
	RotateDamlKey bool `json:"rotate_daml_key"`

	// KnownSigningKeysB64 is the list of current party signing keys
	// (base64-encoded proto). Used to discover the target's old DAML key
	// via vault cross-reference when RotateDamlKey is true.
	KnownSigningKeysB64 []string `json:"known_signing_keys_b64,omitempty"`
}

// GenerateRotatedKeyOutput is the output of [GenerateRotatedKeyOp].
type GenerateRotatedKeyOutput struct {
	ParticipantID  string `json:"participant_id"`
	ParticipantUID string `json:"participant_uid"`

	// Namespace key rotation fields (populated when RotateNamespaceKey is true).
	NewNamespaceKeyB64      string `json:"new_namespace_key_b64,omitempty"`
	NewNamespaceFingerprint string `json:"new_namespace_fingerprint,omitempty"`

	// DAML key rotation fields (populated when RotateDamlKey is true).
	NewDamlKeyB64         string `json:"new_daml_key_b64,omitempty"`
	NewDamlKeyFingerprint string `json:"new_daml_key_fingerprint,omitempty"`

	// Old DAML key info (auto-discovered from vault, populated when RotateDamlKey is true).
	OldDamlKeyFingerprint string `json:"old_daml_key_fingerprint,omitempty"`
	OldDamlKeyB64         string `json:"old_daml_key_b64,omitempty"`
}
