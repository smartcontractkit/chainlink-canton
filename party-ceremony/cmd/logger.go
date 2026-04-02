package cmd

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newCLILogger() (logger.Logger, error) {
	return logger.NewWith(func(cfg *zap.Config) {
		cfg.Encoding = "console"
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.EncoderConfig.TimeKey = ""
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	})
}
