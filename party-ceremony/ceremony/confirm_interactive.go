package ceremony

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

// InteractiveConfirmer prompts the user on In/Out for confirmation before
// signing a transaction. Output is written to Out (typically os.Stderr so
// that os.Stdout stays clean for JSON/machine output).
type InteractiveConfirmer struct {
	In  io.Reader
	Out io.Writer
}

func (c *InteractiveConfirmer) ConfirmTopologySign(_ context.Context, d TopologySignDetail) error {
	fmt.Fprintln(c.Out)
	fmt.Fprintln(c.Out, "╔══════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(c.Out, "║  TOPOLOGY TRANSACTION — Review Before Signing               ║")
	fmt.Fprintln(c.Out, "╠══════════════════════════════════════════════════════════════╣")

	printField(c.Out, "Type", d.MappingType)
	printField(c.Out, "Operation", d.Operation)
	printField(c.Out, "Serial", fmt.Sprintf("%d", d.Serial))

	switch d.MappingType {
	case "DecentralizedNamespaceDefinition":
		printField(c.Out, "Namespace", d.DNSNamespace)
		printField(c.Out, "Threshold", fmt.Sprintf("%d of %d owners", d.DNSThreshold, len(d.DNSOwners)))
		fmt.Fprintln(c.Out, "║  Owners:")
		for _, o := range d.DNSOwners {
			fmt.Fprintf(c.Out, "║    - %s\n", o)
		}
	case "PartyToParticipant":
		printField(c.Out, "Party", d.P2PParty)
		printField(c.Out, "Threshold", fmt.Sprintf("%d", d.P2PThreshold))
		fmt.Fprintln(c.Out, "║  Participants:")
		for _, p := range d.P2PParticipants {
			fmt.Fprintf(c.Out, "║    - %s\n", p)
		}
	case "NamespaceDelegation":
		printField(c.Out, "Namespace", d.NSDNamespace)
	}

	if d.ProposalHash != "" {
		hash := d.ProposalHash
		if len(hash) > 32 {
			hash = hash[:32] + "…"
		}
		printField(c.Out, "Proposal Hash", hash)
	}
	printField(c.Out, "Existing Signatures", fmt.Sprintf("%d", d.ExistingSignatures))
	printField(c.Out, "Your Identity", d.SignerIdentity)

	fmt.Fprintln(c.Out, "╚══════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(c.Out)

	return c.askYesNo()
}

func (c *InteractiveConfirmer) ConfirmDAMLSign(_ context.Context, d DAMLSignDetail) error {
	fmt.Fprintln(c.Out)
	fmt.Fprintln(c.Out, "╔══════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(c.Out, "║  DAML TRANSACTION — Review Before Signing                   ║")
	fmt.Fprintln(c.Out, "╠══════════════════════════════════════════════════════════════╣")

	hash := d.TransactionHash
	if len(hash) > 32 {
		hash = hash[:32] + "…"
	}
	printField(c.Out, "Transaction Hash", hash)
	printField(c.Out, "Your Identity", d.SignerIdentity)

	fmt.Fprintln(c.Out, "╚══════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(c.Out)

	return c.askYesNo()
}

func (c *InteractiveConfirmer) askYesNo() error {
	fmt.Fprint(c.Out, "Confirm signing? [y/N]: ")

	scanner := bufio.NewScanner(c.In)
	if !scanner.Scan() {
		return ErrUserRejected
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer == "y" || answer == "yes" {
		return nil
	}

	return ErrUserRejected
}

func printField(w io.Writer, label, value string) {
	fmt.Fprintf(w, "║  %-22s %s\n", label+":", value)
}
