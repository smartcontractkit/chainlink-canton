package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/smartcontractkit/chainlink-ccv/bootstrap"
	cmd "github.com/smartcontractkit/chainlink-ccv/cmd/verifier"
	"github.com/smartcontractkit/chainlink-ccv/verifier/pkg/commit"
	"github.com/smartcontractkit/chainlink-common/keystore"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/accessors" // register canton accessor factory
	_ "github.com/smartcontractkit/chainlink-canton/deployment/adapters" // register canton adapters
)

func main() {
	if err := bootstrap.Run(
		"CantonCommitteeVerifier",
		cmd.NewCommitteeVerifierServiceFactory(),
		bootstrap.WithKey(commit.DefaultECDSASigningKeyName, "signing", keystore.ECDSA_S256), // ECDSA key for signing verification results
	); err != nil {
		panic(fmt.Sprintf("failed to run Canton committee verifier: %s", err.Error()))
	}
}
