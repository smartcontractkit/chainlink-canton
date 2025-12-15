package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type EnvironmentsConfig struct {
	Environments map[string]EnvironmentConfig `yaml:"environments"`
}

type EnvironmentConfig struct {
	Party       string              `yaml:"party"`
	Description string              `yaml:"description"`
	Contracts   ContractIdentifiers `yaml:"contracts"`
}

type ContractIdentifiers struct {
	Router             string `yaml:"router"`
	OnRamp             string `yaml:"onRamp"`
	FeeQuoter          string `yaml:"feeQuoter"`
	OffRamp            string `yaml:"offRamp"`
	CCV                string `yaml:"ccv"`
	TokenAdminRegistry string `yaml:"tokenAdminRegistry"`
	TokenPool          string `yaml:"tokenPool"`
}

func LoadEnvironments(path string) (*EnvironmentsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments config file: %w", err)
	}

	var config EnvironmentsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse environments config: %w", err)
	}

	if len(config.Environments) == 0 {
		return nil, fmt.Errorf("no environments defined in config")
	}

	return &config, nil
}

func (c *EnvironmentsConfig) EnvironmentNames() []string {
	names := make([]string, 0, len(c.Environments))
	for name := range c.Environments {
		names = append(names, name)
	}
	return names
}

func (c *EnvironmentsConfig) GetEnvironment(name string) (EnvironmentConfig, bool) {
	env, ok := c.Environments[name]
	return env, ok
}
