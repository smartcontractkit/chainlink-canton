package deployment

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
)

type CantonConfigs struct {
	EDS *edsConfig.Config `json:"eds"`
}

type CantonEnvMetadata struct {
	CantonConfigs *CantonConfigs `json:"cantonConfigs,omitempty"`
}

func loadCantonEnvMetadata(ds datastore.DataStore) (*CantonEnvMetadata, error) {
	envMeta, err := ds.EnvMetadata().Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get env metadata: %w", err)
	}

	return parseCantonEnvMetadata(envMeta.Metadata)
}

func loadOrCreateCantonEnvMetadata(ds datastore.MutableDataStore) (*CantonEnvMetadata, error) {
	envMeta, err := ds.EnvMetadata().Get()
	if err != nil {
		//nolint:nilerr // Returning an empty Env
		return &CantonEnvMetadata{}, nil
	}

	return parseCantonEnvMetadata(envMeta.Metadata)
}

func parseCantonEnvMetadata(metadata any) (*CantonEnvMetadata, error) {
	if metadata == nil {
		return &CantonEnvMetadata{}, nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal env metadata: %w", err)
	}

	var cantonMeta CantonEnvMetadata
	if err := json.Unmarshal(data, &cantonMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CCV env metadata: %w", err)
	}

	return &cantonMeta, nil
}

func saveCantonEnvMetadata(ds datastore.MutableDataStore, cantonMeta *CantonEnvMetadata) error {
	var base json.RawMessage = []byte(`{}`)

	if envMeta, err := ds.EnvMetadata().Get(); err == nil && envMeta.Metadata != nil {
		b, err := json.Marshal(envMeta.Metadata)
		if err != nil {
			return err
		}
		base = b
	}

	patch, err := json.Marshal(cantonMeta)
	if err != nil {
		return err
	}

	merged, err := jsonpatch.MergePatch(base, patch)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := json.Unmarshal(merged, &result); err != nil {
		return err
	}

	return ds.EnvMetadata().Set(datastore.EnvMetadata{Metadata: result})
}

func GetEDSConfig(ds datastore.DataStore) (*edsConfig.Config, error) {
	cantonMeta, err := loadCantonEnvMetadata(ds)
	if err != nil {
		return nil, err
	}

	if cantonMeta.CantonConfigs == nil || cantonMeta.CantonConfigs.EDS == nil {
		return nil, fmt.Errorf("EDS config not found in env metadata")
	}

	return cantonMeta.CantonConfigs.EDS, nil
}

func SaveEDSConfig(ds datastore.MutableDataStore, cfg *edsConfig.Config) error {
	cantonMeta, err := loadOrCreateCantonEnvMetadata(ds)
	if err != nil {
		return err
	}

	if cantonMeta.CantonConfigs == nil {
		cantonMeta.CantonConfigs = &CantonConfigs{}
	}

	cantonMeta.CantonConfigs.EDS = cfg

	return saveCantonEnvMetadata(ds, cantonMeta)
}
