package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/spf13/cobra"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
)

func init() {
	// Register the canton modifier for the canton family.
	// TODO: in the future, all chain-specific registrations should be done here.
	committeeverifier.RegisterModifier(chain_selectors.FamilyCanton, cantondevenv.CommitteeVerifierModifier)
}

var rootCmd = &cobra.Command{
	Use:   "ccv",
	Short: "A CCV local environment tool",
}

// TODO: the commands below are currently copy/pasted from chainlink-ccv, but instead, chainlink-ccv/devenv
// should probably export a library of CLI functions that can be selectively imported
// by other CLIs.
var upCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"u"},
	Short:   "Spin up the development environment",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		} else {
			configFile = "env.toml"
		}
		framework.L.Info().Str("Config", configFile).Msg("Creating development environment")
		_ = os.Setenv("CTF_CONFIGS", configFile)
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		_, err := ccv.NewEnvironment()
		if err != nil {
			return err
		}

		return nil
	},
}

var downCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Tear down the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		framework.L.Info().Msg("Tearing down the development environment")
		err := framework.RemoveTestContainers()
		if err != nil {
			return fmt.Errorf("failed to clean Docker resources: %w", err)
		}

		return nil
	},
}

var dumpLogsCmd = &cobra.Command{
	Use:     "dump-logs",
	Aliases: []string{"dl"},
	Short:   "Dump the logs of the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirSuffix, _ := cmd.Flags().GetString("dir-suffix")
		if dirSuffix == "" {
			return fmt.Errorf("dir-suffix is required")
		}

		framework.L.Info().Msg("Dumping the logs of all docker containers in the development environment")
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, dirSuffix))
		if err != nil {
			return fmt.Errorf("failed to dump logs: %w", err)
		}
		framework.L.Info().Msg("Logs dumped successfully")
		return nil
	},
}

func checkDockerIsRunning() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Can't create Docker client, please check if Docker daemon is running!")
		os.Exit(1)
	}
	defer cli.Close()
	_, err = cli.Ping(context.Background())
	if err != nil {
		fmt.Println("Docker is not running, please start Docker daemon first!")
		os.Exit(1)
	}
}

func main() {
	checkDockerIsRunning()
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(dumpLogsCmd)
	if err := rootCmd.Execute(); err != nil {
		ccv.Plog.Err(err).Send()
		os.Exit(1)
	}
}
