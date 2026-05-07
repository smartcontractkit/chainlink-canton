package main

import (
	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // register all canton-related devenv components
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cli"
)

func main() {
	cli.RunCLI()
}
