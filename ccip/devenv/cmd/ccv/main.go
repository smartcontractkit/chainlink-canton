package main

import (
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cli"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // register all canton-related devenv components
)

func main() {
	cli.RunCLI()
}
