package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "canton-party-ceremony",
	Short: "CLI for managing Canton decentralized party onboarding ceremonies",
	Long: `canton-party-ceremony manages the multi-step decentralized party
onboarding ceremony workflow for Canton distributed ledgers.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
