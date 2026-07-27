package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// uploadDarsCmd uploads DAR files to a single participant via Admin PackageService.
// Unlike contract-deploy, this command never touches the Ledger API and cannot
// prepare, sign, or execute contract creation.
var uploadDarsCmd = &cobra.Command{
	Use:   "upload-dars",
	Short: "Upload DAR packages to a participant (Admin API only)",
	Long: `Upload DAML package archives (DARs) to one Canton participant using the
Admin PackageService. This is the DAR-only path — no party verification, no
ledger prepare/sign/execute, and no contract deployment.

Each CV operator runs this against their own participant config. Re-uploading
the same DAR is idempotent.

DARs are read from ./dars/<version>/<name>-<version>.dar relative to the
working directory (symlink ../contracts/dars → dars before running).`,
	RunE: runUploadDars,
}

func init() {
	f := uploadDarsCmd.Flags()

	f.String("packages", "", "Comma-separated list of packages in name:version format (required)")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("dars-dir", "dars", "Directory containing versioned DAR subdirectories (dev/, released/)")

	_ = uploadDarsCmd.MarkFlagRequired("packages")

	rootCmd.AddCommand(uploadDarsCmd)
}

func runUploadDars(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	packagesRaw, _ := f.GetString("packages")
	configPath, _ := f.GetString("config")
	darsDir, _ := f.GetString("dars-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	packages, err := parsePackageRefs(packagesRaw)
	if err != nil {
		return err
	}
	if len(packages) == 0 {
		return fmt.Errorf("--packages must list at least one package")
	}

	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	adminClient := client.NewGRPCClient(conn)
	ctx := cmd.Context()
	loader := ledger.FileDARLoader(darsDir)

	fmt.Fprintf(os.Stdout, "Uploading %d package(s) via %s:%d...\n", len(packages), cfg.AdminHost, cfg.AdminPort)

	for _, pkg := range packages {
		darBytes, err := loader(pkg.Name, pkg.Version)
		if err != nil {
			return err
		}

		pkgID, err := adminClient.UploadDar(ctx, darBytes)
		if err != nil {
			return fmt.Errorf("uploading %q@%s: %w", pkg.Name, pkg.Version, err)
		}

		fmt.Fprintf(os.Stdout, "  uploaded %s@%s → %s\n", pkg.Name, pkg.Version, pkgID)
	}

	fmt.Fprintln(os.Stdout, "Done.")

	return nil
}
