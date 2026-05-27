package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/smartcontractkit/chainlink-ccv/bootstrap"
	cmd "github.com/smartcontractkit/chainlink-ccv/cmd/verifier"
	"go.uber.org/zap/zapcore"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/accessors" // register canton accessor factory
)

func main() {
	if err := bootstrap.Run(
		"CantonCommitteeVerifier",
		cmd.NewCommitteeVerifierServiceFactory(),
		bootstrap.WithLogLevel(zapcore.InfoLevel),
	); err != nil {
		panic(fmt.Sprintf("failed to run Canton committee verifier: %s", err.Error()))
	}
}
