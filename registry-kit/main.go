package main

import (
	"os"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
