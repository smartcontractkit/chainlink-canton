package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultStatePath = "registry-kit.state.json"

// State tracks contract IDs produced by CLI steps (machine-written progress).
type State struct {
	ProviderServiceRequestCID  string `json:"provider_service_request_cid,omitempty"`
	ProviderServiceCID         string `json:"provider_service_cid,omitempty"`
	ProviderConfigurationCID   string `json:"provider_configuration_cid,omitempty"`
	RegistrarServiceRequestCID string `json:"registrar_service_request_cid,omitempty"`
	RegistrarServiceCID        string `json:"registrar_service_cid,omitempty"`
	AllocationFactoryCID       string `json:"allocation_factory_cid,omitempty"`
	TransferRuleCID            string `json:"transfer_rule_cid,omitempty"`
	InstrumentID               string `json:"instrument_id,omitempty"`
	InstrumentConfigurationCID string `json:"instrument_configuration_cid,omitempty"`
	LastMintRequestCID         string `json:"last_mint_request_cid,omitempty"`
	TokenConfigCID             string `json:"token_config_cid,omitempty"`
	TokenAdminRegistryCID      string `json:"token_admin_registry_cid,omitempty"`
}

// LoadState reads registry-kit.state.json. Missing file yields an empty state.
func LoadState(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, fmt.Errorf("read state %q: %w", path, err)
	}

	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, fmt.Errorf("parse state %q: %w", path, err)
	}

	return st, nil
}

// Save writes state atomically next to the config directory when path is relative.
func (s State) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return fmt.Errorf("create state dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	data = append(data, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write state temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename state: %w", err)
	}

	return nil
}

// StatePathNextTo returns registry-kit.state.json beside the config file.
func StatePathNextTo(configPath string) string {
	dir := filepath.Dir(configPath)
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, DefaultStatePath)
}
