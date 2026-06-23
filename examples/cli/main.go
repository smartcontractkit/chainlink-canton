package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdpkg "github.com/smartcontractkit/chainlink-canton/examples/cli/cmd"
)

func main() {
	root := &cobra.Command{
		Use:           "ccip-demo",
		Short:         "CCIP Canton <-> EVM demo CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	var (
		configPath string
		network    string
	)
	root.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "path to YAML config file")
	root.PersistentFlags().StringVarP(&network, "network", "n", "mainnet", "network profile (mainnet|testnet)")

	g := &cmdpkg.Globals{
		ConfigPath: &configPath,
		Network:    &network,
	}

	root.AddCommand(cmdpkg.NewEVMCmd(g))
	root.AddCommand(cmdpkg.NewCantonCmd(g))

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
