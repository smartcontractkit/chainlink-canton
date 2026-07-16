package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/smartcontractkit/chainlink-ccv/bootstrap"
	cmd "github.com/smartcontractkit/chainlink-ccv/cmd/verifier"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/accessors" // register canton accessor factory
)

func main() {
	if err := bootstrap.Run(
		"CantonCommitteeVerifier",
		cmd.NewCommitteeVerifierServiceFactory(),
	); err != nil {
		panic(fmt.Sprintf("failed to run Canton committee verifier: %s", err.Error()))
	}
}
