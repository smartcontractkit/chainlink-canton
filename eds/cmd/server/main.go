package main

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/service"
)

func main() {
	ctx := context.Background()
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel)

	cfgPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	cfgReader, err := os.Open(*cfgPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open config file")
	}
	cfg, err := config.Read(cfgReader)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read config file")
	}

	logger.Fatal().Err(service.RunEDS(ctx, logger, cfg)).Msg("Running EDS server")
}
