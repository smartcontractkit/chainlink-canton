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
