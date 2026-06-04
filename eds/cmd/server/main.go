package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
)

const (
	// EnvConfigFile - The path to the config to use, defaults to 'config.toml'
	// Multiple, comma-separated, files can be specified, in which case they take precedence in increasing order:
	//  If e.g. CONFIG_FILE=config1.toml,config2.toml are specified, values from config2.toml will override values from config1.toml
	EnvConfigFile = "CONFIG_FILE"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel)

	cfgPaths := os.Getenv(EnvConfigFile)
	if cfgPaths == "" {
		cfgPaths = "config.toml"
	}

	paths := strings.Split(cfgPaths, ",")
	readers := make([]io.Reader, len(paths))
	for i, s := range paths {
		s := strings.TrimSpace(s)
		logger.Info().Str("file", s).Int("index", i).Msg("Reading config...")
		cfgReader, err := os.Open(s) //nolint:gosec
		if err != nil {
			stop()
			logger.Fatal().Err(err).Str("file", s).Int("index", i).Msg("failed to open config file")
		}

		readers[i] = cfgReader
	}

	cfg, err := config.ReadAndMerge(readers...)
	if err != nil {
		stop()
		logger.Fatal().Err(err).Msg("failed to parse config files")
	}

	runErr := service.RunEDS(ctx, logger, cfg)
	stop()
	if runErr != nil {
		logger.Fatal().Err(runErr).Msg("EDS server exited with error")
	}
	logger.Info().Msg("EDS server exited without an error")
}
