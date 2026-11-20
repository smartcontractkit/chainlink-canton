package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// queryContractsCmd lists a party's active contracts of one template, so an
// operator can check what a ceremony committed when its result was lost.
//
// Usage:
//
//	canton-party-ceremony query-contracts \
//	  --party "prefix::namespace" \
//	  --template-id "#utility-registry-app-v0:Utility.Registry.App.V0.Service.Provider:ProviderService" \
//	  --config ./participant-config.json
var queryContractsCmd = &cobra.Command{
	Use:   "query-contracts",
	Short: "List a party's active contracts of one template",
	Long: `Query the Ledger API Active Contract Set as of the current ledger end and print
each matching contract's ID, creation offset and time, and create arguments.

--template-id is #<package-name>:<module>:<entity>; package IDs are rejected.
The configured Ledger API user needs read rights for the party; the
contract-deploy and exercise-choice ceremonies grant them.

After a lost execute response, run this on the participant that submitted, and
wait 10 minutes first: a transaction still in flight is not yet listed.`,
	RunE: runQueryContracts,
}

func init() {
	f := queryContractsCmd.Flags()

	f.String("party", "", "Party whose contracts to list (required)")
	f.String("template-id", "", "Template as #<package-name>:<module>:<entity> (required)")
	f.String("config", "participant-config.json", "Path to participant config JSON file")

	_ = queryContractsCmd.MarkFlagRequired("party")
	_ = queryContractsCmd.MarkFlagRequired("template-id")

	rootCmd.AddCommand(queryContractsCmd)
}

func runQueryContracts(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	party, _ := f.GetString("party")
	templateID, _ := f.GetString("template-id")
	configPath, _ := f.GetString("config")

	pkg, module, entity, err := parseTemplateID(templateID)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(pkg, "#") {
		return fmt.Errorf("template ID %q: the package must be a #<package-name> reference", templateID)
	}

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	conn, err := client.DialLedger(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton ledger API: %w", err)
	}
	defer conn.Close()

	contracts, err := client.NewGRPCLedgerClient(conn).GetActiveContractsByTemplateForParty(cmd.Context(), party, pkg, module, entity)
	if err != nil {
		return fmt.Errorf("querying active contracts: %w", err)
	}

	return printContracts(os.Stdout, contracts)
}

func printContracts(w io.Writer, contracts []*apiv2.CreatedEvent) error {
	fmt.Fprintf(w, "Active contracts: %d\n", len(contracts))
	for _, c := range contracts {
		args, err := protojson.Marshal(c.GetCreateArguments())
		if err != nil {
			return fmt.Errorf("encoding create arguments of %s: %w", c.GetContractId(), err)
		}
		fmt.Fprintf(w, "\nContract ID: %s\nOffset:      %d\nCreated at:  %s\nArguments:   %s\n",
			c.GetContractId(), c.GetOffset(), c.GetCreatedAt().AsTime().UTC().Format("2006-01-02T15:04:05Z"), args)
	}

	return nil
}
