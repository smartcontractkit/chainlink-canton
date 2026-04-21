package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/chainlink/canton-party-ceremony/ceremony"
)

var rootCmd = &cobra.Command{
	Use:   "canton-party-ceremony",
	Short: "CLI for managing Canton decentralized party onboarding ceremonies",
	Long: `canton-party-ceremony manages the multi-step decentralized party
onboarding ceremony workflow for Canton distributed ledgers.`,
}

func init() {
	rootCmd.PersistentFlags().Bool("confirm", false, "Prompt for interactive confirmation before signing transactions")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// confirmerFromFlags returns an InteractiveConfirmer when --confirm is set,
// or nil otherwise (no prompts).
func confirmerFromFlags(cmd *cobra.Command) ceremony.Confirmer {
	confirm, _ := cmd.Flags().GetBool("confirm")
	if !confirm {
		return nil
	}

	return &ceremony.InteractiveConfirmer{In: os.Stdin, Out: os.Stderr}
}
