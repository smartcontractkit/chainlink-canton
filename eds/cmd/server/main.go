package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
)

const (
	// EnvConfigFile - The path to the config to use, defaults to 'config.toml'
	EnvConfigFile = "CONFIG_FILE"
)

func main() {
	ctx := context.Background()
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel)

	cfgPath := os.Getenv(EnvConfigFile)
	if cfgPath == "" {
		cfgPath = "config.toml"
	}

	logger.Info().Str("file", cfgPath).Msg("Reading config...")
	cfgReader, err := os.Open(cfgPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open config file")
	}
	cfg, err := config.Read(cfgReader)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read config file")
	}

	logger.Fatal().Err(service.RunEDS(ctx, logger, cfg)).Msg("Running EDS server")
}
