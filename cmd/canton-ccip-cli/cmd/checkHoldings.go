package cmd

import (
	"fmt"
	"math/big"
	"slices"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/utility/registry_holding_v0"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

var (
	registryHoldingTemplate  = contracts.TemplateIDFromBinding(registry_holding_v0.Holding{})
	legacyHoldingTemplate    = contracts.TemplateIDFromBinding(link.LinkHolding{})
	legacyLockedLinkTemplate = contracts.TemplateIDFromBinding(link.LockedLinkHolding{})
)

const damlDecimalFractionDigits = 10

type templateRef struct {
	packageID, packageName, module, entity string
}

func templateRefOf(ev *apiv2.CreatedEvent) templateRef {
	id := ev.GetTemplateId()

	return templateRef{packageID: id.GetPackageId(), packageName: ev.GetPackageName(), module: id.GetModuleName(), entity: id.GetEntityName()}
}

// matches compares by package name because binding template ids carry "#<package-name>", not a package id.
func (t templateRef) matches(binding contracts.TemplateID) bool {
	return t.packageName == strings.TrimPrefix(binding.PackageID, "#") && t.module == binding.ModuleName && t.entity == binding.EntityName
}

type holdingEvent struct {
	cid        string
	template   templateRef
	owner      types.PARTY
	instrument splice_api_token_holding_v1.InstrumentId
	amount     types.NUMERIC
	locked     bool
	viewError  string
}

type holdingCheckConfig struct {
	instrument         splice_api_token_holding_v1.InstrumentId
	registryPackageIDs []string
	postMigration      bool
	expectedSupply     *big.Rat
}

type holdingFinding struct {
	kind, cid, detail string
}

type holdingReport struct {
	holdings, owners, locked int
	total                    *big.Rat
	findings, info           []holdingFinding
}

func checkHoldings(cfg holdingCheckConfig, events []holdingEvent) (holdingReport, error) {
	r := holdingReport{total: new(big.Rat)}
	owners := map[types.PARTY]struct{}{}
	for _, e := range events {
		isLegacy := e.template.matches(legacyHoldingTemplate) || e.template.matches(legacyLockedLinkTemplate)
		switch {
		case e.viewError != "":
			r.findings = append(r.findings, holdingFinding{"holding view failed", e.cid, e.viewError})
		case isLegacy:
			if e.instrument.Admin != cfg.instrument.Admin {
				continue
			}
			f := holdingFinding{"legacy LINK not drained", e.cid, e.template.module + ":" + e.template.entity}
			if cfg.postMigration {
				r.findings = append(r.findings, f)
			} else {
				r.info = append(r.info, f)
			}
		case e.instrument != cfg.instrument:
			continue
		case !e.template.matches(registryHoldingTemplate) || !slices.Contains(cfg.registryPackageIDs, e.template.packageID):
			detail := fmt.Sprintf("%s:%s:%s (package %s)", e.template.packageName, e.template.module, e.template.entity, e.template.packageID)
			r.findings = append(r.findings, holdingFinding{"foreign template", e.cid, detail})
		default:
			amount, ok := new(big.Rat).SetString(string(e.amount))
			if !ok {
				return holdingReport{}, fmt.Errorf("holding %s: invalid amount %q", e.cid, e.amount)
			}
			r.total.Add(r.total, amount)
			r.holdings++
			owners[e.owner] = struct{}{}
			if e.locked {
				r.locked++
			}
		}
	}
	r.owners = len(owners)
	if cfg.expectedSupply != nil && r.total.Cmp(cfg.expectedSupply) != 0 {
		detail := fmt.Sprintf("holdings total %s, expected %s", r.total.FloatString(damlDecimalFractionDigits), cfg.expectedSupply.FloatString(damlDecimalFractionDigits))
		r.findings = append(r.findings, holdingFinding{"supply mismatch", "-", detail})
	}

	return r, nil
}

func holdingEventsFromViews(active []*apiv2.ActiveContract) ([]holdingEvent, error) {
	var out []holdingEvent
	for _, c := range active {
		ev := c.GetCreatedEvent()
		for _, view := range ev.GetInterfaceViews() {
			// A failed view has no value; skipping it would hide the holding from the check.
			if st := view.GetViewStatus(); st.GetCode() != 0 {
				out = append(out, holdingEvent{cid: ev.GetContractId(), template: templateRefOf(ev), viewError: fmt.Sprintf("status %d: %s", st.GetCode(), st.GetMessage())})
				continue
			}
			var hv splice_api_token_holding_v1.HoldingView
			if err := ledger.RecordToStruct(view.GetViewValue(), &hv); err != nil {
				return nil, fmt.Errorf("decode holding view %s: %w", ev.GetContractId(), err)
			}
			out = append(out, holdingEvent{cid: ev.GetContractId(), template: templateRefOf(ev), owner: hv.Owner, instrument: hv.InstrumentId, amount: hv.Amount, locked: hv.Lock != nil})
		}
	}

	return out, nil
}

// LockedLinkHolding does not implement HoldingV1, so the interface query never returns it.
func holdingEventsFromLockedLink(active []*apiv2.ActiveContract) ([]holdingEvent, error) {
	out := make([]holdingEvent, 0, len(active))
	for _, c := range active {
		ev := c.GetCreatedEvent()
		l, err := bindings.UnmarshalCreatedEvent[link.LockedLinkHolding](ev)
		if err != nil {
			return nil, fmt.Errorf("unmarshal LockedLinkHolding %s: %w", ev.GetContractId(), err)
		}
		out = append(out, holdingEvent{cid: ev.GetContractId(), template: templateRefOf(ev), owner: l.LockedSender, instrument: l.LockedInstrumentId, amount: l.LockedAmount, locked: true})
	}

	return out, nil
}

func newCantonCheckHoldingsCmd(g *Globals) *cobra.Command {
	var (
		instrument, admin, expectedSupply string
		extraRegistryPackageIDs           []string
		postMigration                     bool
	)
	c := &cobra.Command{
		Use:   "check-holdings",
		Short: "Check registry holdings of an instrument for foreign templates, legacy LINK leftovers and supply drift",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cfg := holdingCheckConfig{
				instrument:         splice_api_token_holding_v1.InstrumentId{Id: types.TEXT(instrument), Admin: types.PARTY(admin)},
				registryPackageIDs: append([]string{registry_holding_v0.PackageID}, extraRegistryPackageIDs...),
				postMigration:      postMigration,
			}
			if expectedSupply != "" {
				v, ok := new(big.Rat).SetString(expectedSupply)
				if !ok {
					return fmt.Errorf("invalid --expected-supply %q", expectedSupply)
				}
				cfg.expectedSupply = v
			}
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			viewed, err := testhelpers.ListActiveContractsByInterfaceId(ctx, b.Participant, contracts.MustTemplateIDFromString(splice_api_token_holding_v1.IHoldingInterfaceID()).ToLedgerIdentifier())
			if err != nil {
				return fmt.Errorf("list holdings: %w", err)
			}
			lockedLegacy, err := testhelpers.ListActiveContractsByTemplateId(ctx, b.Participant, legacyLockedLinkTemplate.ToLedgerIdentifier())
			if err != nil {
				return fmt.Errorf("list LockedLinkHolding: %w", err)
			}
			events, err := holdingEventsFromViews(viewed)
			if err != nil {
				return err
			}
			lockedEvents, err := holdingEventsFromLockedLink(lockedLegacy)
			if err != nil {
				return err
			}
			r, err := checkHoldings(cfg, append(events, lockedEvents...))
			if err != nil {
				return err
			}

			expected := "-"
			if cfg.expectedSupply != nil {
				expected = cfg.expectedSupply.FloatString(damlDecimalFractionDigits)
			}
			summary := table.NewWriter()
			summary.SetStyle(table.StyleLight)
			summary.SetTitle(fmt.Sprintf("%s@%s holdings", instrument, admin))
			summary.AppendRows([]table.Row{
				{"Party", b.Participant.PartyID},
				{"Holdings", r.holdings},
				{"Distinct owners", r.owners},
				{"Locked", r.locked},
				{"Total", r.total.FloatString(damlDecimalFractionDigits)},
				{"Expected total", expected},
			})
			fmt.Println(summary.Render())
			printHoldingFindings("Findings", r.findings)
			printHoldingFindings("Info: legacy holdings (findings when --post-migration is set)", r.info)
			if len(r.findings) > 0 {
				return fmt.Errorf("%d finding(s)", len(r.findings))
			}

			return nil
		},
	}
	c.Flags().StringVar(&instrument, "instrument", "", "instrument id, e.g. LINK")
	c.Flags().StringVar(&admin, "admin", "", "instrument admin (registrar) party")
	// The package name is author-chosen, so only the package id proves a holding is the genuine registry template.
	c.Flags().StringSliceVar(&extraRegistryPackageIDs, "registry-package-id", nil, "extra accepted utility-registry-holding-v0 package ids, e.g. older registry versions; the pinned id is always accepted")
	c.Flags().StringVar(&expectedSupply, "expected-supply", "", "expected sum of registry holdings, locked included")
	c.Flags().BoolVar(&postMigration, "post-migration", true, "treat active legacy LinkHolding or LockedLinkHolding as findings; set false before migration to report them as info")
	_ = c.MarkFlagRequired("instrument")
	_ = c.MarkFlagRequired("admin")

	return c
}

func printHoldingFindings(title string, findings []holdingFinding) {
	if len(findings) == 0 {
		return
	}
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.SetTitle(title)
	tw.AppendHeader(table.Row{"Kind", "Contract ID", "Detail"})
	for _, f := range findings {
		tw.AppendRow(table.Row{f.kind, f.cid, f.detail})
	}
	fmt.Println(tw.Render())
}
