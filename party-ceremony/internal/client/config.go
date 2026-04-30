package client

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClientConfig holds participant-specific configuration loaded from a JSON file.
//
// Example file (participant-config.json):
//
//	{
//	  "participant_id": "p1",
//	  "admin_jwt":      "ey...",
//	  "admin_host":     "localhost",
//	  "admin_port":     5001
//	}
type ClientConfig struct {
	// ParticipantID identifies this node in the ceremony (required).
	ParticipantID string `json:"participant_id"`

	// AdminJWT is the Bearer token for the Canton Admin API (optional).
	// When empty the client connects without authentication, which works for
	// localhost / no-auth setups.
	AdminJWT string `json:"admin_jwt,omitempty"`

	// AdminHost and AdminPort are the Canton Admin gRPC endpoint (optional).
	// Defaults: localhost:5001.
	AdminHost string `json:"admin_host,omitempty"`
	AdminPort int    `json:"admin_port,omitempty"`

	// LedgerHost and LedgerPort are the Canton Ledger gRPC endpoint (optional).
	// Used by the contract-deploy ceremony to interact with the Ledger API.
	// Defaults: localhost:5002.
	LedgerHost string `json:"ledger_host,omitempty"`
	LedgerPort int    `json:"ledger_port,omitempty"`

	// LedgerJWT is the Bearer token for the Canton Ledger API (optional).
	LedgerJWT string `json:"ledger_jwt,omitempty"`

	// UserID is the Ledger API user ID that the contract-deploy ceremony will
	// grant actAs/readAs rights for the decentralized party.
	// Only required in JWT-auth environments. Leave empty for no-auth setups.
	UserID string `json:"user_id,omitempty"`

	// KmsNamespaceKeyID is the external KMS identifier (e.g. AWS KMS ARN) for a
	// pre-existing NAMESPACE signing key. When set the ceremony registers this
	// key via VaultService.RegisterKmsSigningKey instead of generating a new one.
	// Onboarding and add-participant require this together with KmsProtocolKeyID;
	// key rotation may use it independently when only the namespace key rotates.
	KmsNamespaceKeyID string `json:"kms_namespace_key_id,omitempty"`

	// KmsProtocolKeyID is the external KMS identifier (e.g. AWS KMS ARN) for a
	// pre-existing PROTOCOL (DAML) signing key. When set the ceremony registers
	// this key via VaultService.RegisterKmsSigningKey instead of generating a new one.
	// It is also used by contract-deploy to sign prepared transaction hashes.
	KmsProtocolKeyID string `json:"kms_protocol_key_id,omitempty"`
}

// KMSConfig contains the local operator's external KMS key IDs. These values
// are intentionally loaded from participant-config.json on each init/resume
// invocation and are not persisted in shared workflow state.
type KMSConfig struct {
	NamespaceKeyID string
	ProtocolKeyID  string
}

// KMS returns the local KMS configuration for dependency injection.
func (cfg ClientConfig) KMS() KMSConfig {
	return KMSConfig{
		NamespaceKeyID: cfg.KmsNamespaceKeyID,
		ProtocolKeyID:  cfg.KmsProtocolKeyID,
	}
}

// LoadConfig reads a ClientConfig from the given JSON file path.
// Returns an error if the file cannot be read, the JSON is malformed, or
// participant_id is blank.
func LoadConfig(path string) (ClientConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ClientConfig{}, fmt.Errorf("reading config %q: %w", path, err)
	}

	var cfg ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ClientConfig{}, fmt.Errorf("parsing config %q: %w", path, err)
	}

	if cfg.ParticipantID == "" {
		return ClientConfig{}, fmt.Errorf("config %q: participant_id is required", path)
	}

	return cfg, nil
}
