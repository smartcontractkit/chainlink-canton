package main

import (
	"context"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/bootstrap"
	cmd "github.com/smartcontractkit/chainlink-ccv/cmd/verifier"
	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/verifier/commit"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	"github.com/smartcontractkit/chainlink-canton/ccip/accessors"
)

const CantonConfigPathEnv = "CANTON_CONFIG_PATH"

func loadConfig(path string) (*ccip.Config, error) {
	var cfg ccip.Config
	if md, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %w", path, err)
	} else if len(md.Undecoded()) > 0 {
		return nil, fmt.Errorf("unknown fields in config: %v", md.Undecoded())
	}

	return &cfg, nil
}

func main() {
	if err := bootstrap.Run(
		"CantonCommitteeVerifier",
		cmd.NewServiceFactory(
			chainsel.FamilyCanton,
			func(
				ctx context.Context,
				lggr logger.Logger,
				infos map[string]*ccip.BlockchainInfo,
				cfg commit.Config,
			) (chainaccess.AccessorFactory, error) {
				configPath, ok := os.LookupEnv(CantonConfigPathEnv)
				if !ok {
					configPath = ccip.DefaultCantonConfigPath
				}

				cantonConfig, err := loadConfig(configPath)
				if err != nil {
					return nil, fmt.Errorf("failed to load config: %w", err)
				}

				lggr.Infow("loaded canton config", "config", cantonConfig)

				return accessors.NewFactory(lggr, cantonConfig.BlockchainInfos, cantonConfig.ReaderConfigs, cfg.RMNRemoteAddresses), nil
			}),
		bootstrap.WithLogLevel[commit.JobSpec](zapcore.InfoLevel),
	); err != nil {
		panic(fmt.Sprintf("failed to run Canton committee verifier: %s", err.Error()))
	}
}
