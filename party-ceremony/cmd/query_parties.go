package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// queryPartiesCmd queries and displays all decentralized namespace definitions
// visible to the local participant. Its output provides the identifiers needed
// to run the kick ceremony (decentralized party ID, participant UIDs, and
// namespace fingerprints).
//
// Usage:
//
//	canton-party-ceremony query-parties \
//	  --synchronizer-id global \
//	  --config ./participant-config.json
var queryPartiesCmd = &cobra.Command{
	Use:   "query-parties",
	Short: "List all decentralized parties visible to this participant",
	Long: `Query the Canton topology store and print all active decentralized namespace
definitions with their owner fingerprints and current PartyToParticipant state.

The output provides the exact identifiers required for the kick ceremony:
  - decentralized-party-id  (--decentralized-party-id)
  - participant UIDs         (--remaining-participants, --kicked-participant-id)
  - namespace fingerprints   (--kicked-namespace-fingerprint)`,
	RunE: runQueryParties,
}

func init() {
	f := queryPartiesCmd.Flags()

	f.String("synchronizer-id", "", "Canton synchronizer ID to query (required)")
	f.String("config", "participant-config.json", "Path to participant config JSON file")

	_ = queryPartiesCmd.MarkFlagRequired("synchronizer-id")

	rootCmd.AddCommand(queryPartiesCmd)
}

func runQueryParties(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	synchronizerID, _ := f.GetString("synchronizer-id")
	configPath, _ := f.GetString("config")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	c := client.NewGRPCClient(conn)
	ctx := cmd.Context()

	dnsEntries, err := c.ListDecentralizedNamespaces(ctx, synchronizerID)
	if err != nil {
		return fmt.Errorf("listing decentralized namespaces: %w", err)
	}

	if len(dnsEntries) == 0 {
		fmt.Fprintln(os.Stdout, "No active decentralized namespaces found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	for _, dns := range dnsEntries {
		fmt.Fprintf(w, "Namespace:\t%s\n", dns.DecentralizedNamespace)
		fmt.Fprintf(w, "  Threshold:\t%d-of-%d\n", dns.Threshold, len(dns.Owners))
		fmt.Fprintf(w, "  Serial:\t%d\n", dns.Serial)
		fmt.Fprintf(w, "  Owners:\n")
		for _, fp := range dns.Owners {
			fmt.Fprintf(w, "    - fingerprint:\t%s\n", fp)
		}

		// Attempt to look up the P2P mapping. Party IDs contain the namespace,
		// but the prefix is not stored in the DNS. We use the namespace itself as
		// the filter which Canton partial-matches on.
		p2pState, p2pErr := c.GetP2P(ctx, dns.DecentralizedNamespace, synchronizerID)
		if p2pErr == nil {
			fmt.Fprintf(w, "  Party ID:\t%s\n", p2pState.Party)
			fmt.Fprintf(w, "  P2P Threshold:\t%d\n", p2pState.Threshold)
			fmt.Fprintf(w, "  P2P Serial:\t%d\n", p2pState.Serial)
			fmt.Fprintf(w, "  Hosting Participants:\n")
			for _, p := range p2pState.Participants {
				perm := strings.TrimPrefix(p.Permission, "PARTICIPANT_PERMISSION_")
				fmt.Fprintf(w, "    - uid:\t%s\tpermission: %s\n", p.ParticipantUID, perm)
			}
		}
		fmt.Fprintln(w, "---")
	}

	return nil
}
