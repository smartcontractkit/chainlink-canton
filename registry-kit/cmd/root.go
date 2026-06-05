package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/config"
)

var (
	configPath string
	statePath  string
)

var rootCmd = &cobra.Command{
	Use:   "canton-registry-kit",
	Short: "DevNet CLI for Canton Registry token onboarding and CCIP pool linking",
	Long: `canton-registry-kit automates Registry utility onboarding, token lifecycle,
and CCIP TokenAdminRegistry linking on devnet.cv1.

Configuration: registry-kit.toml (stable inputs) + registry-kit.state.json (progress).`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", config.DefaultConfigPath, "Path to registry-kit.toml")
	rootCmd.PersistentFlags().StringVar(&statePath, "state", "", "Path to registry-kit.state.json (default: beside --config)")

	rootCmd.AddCommand(onboardingCmd)
	rootCmd.AddCommand(issuerCmd)
	rootCmd.AddCommand(operatorCmd)
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}
