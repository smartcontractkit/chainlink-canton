package cmd

import (
	"math/big"
	"strings"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
)

const (
	testRegistrar = "ccipOwner::1220aa"
	testPackageID = "8107899ac4723ce986bf7d27416534e576e54b92161e46150a595fb78ff3d3a1"
)

var testLink = splice_api_token_holding_v1.InstrumentId{Id: "LINK", Admin: testRegistrar}

func refFor(binding contracts.TemplateID, packageID string) templateRef {
	return templateRef{packageID: packageID, packageName: strings.TrimPrefix(binding.PackageID, "#"), module: binding.ModuleName, entity: binding.EntityName}
}

func genuine(cid, owner, amount string, locked bool) holdingEvent {
	return holdingEvent{cid: cid, template: refFor(registryHoldingTemplate, testPackageID), owner: types.PARTY(owner), instrument: testLink, amount: types.NUMERIC(amount), locked: locked}
}

func TestCheckHoldings(t *testing.T) {
	t.Parallel()
	foreign := holdingEvent{cid: "cid-foreign", template: templateRef{"deadbeef", "evil-token", "Evil.Token", "Holding"}, owner: "bob", instrument: testLink, amount: "1000.0"}
	sameName := genuine("cid-same-name", "bob", "1000.0", false)
	sameName.template.packageID = "f00d"
	viewFailed := holdingEvent{cid: "cid-view-failed", template: refFor(registryHoldingTemplate, testPackageID), viewError: "interpretation error"}
	legacy := holdingEvent{cid: "cid-legacy", template: refFor(legacyHoldingTemplate, "abc"), owner: "carol", instrument: testLink, amount: "5.0"}
	legacyLocked := holdingEvent{cid: "cid-legacy-locked", template: refFor(legacyLockedLinkTemplate, "abc"), owner: "carol", instrument: testLink, amount: "2.0", locked: true}
	legacyOtherAdmin := legacy
	legacyOtherAdmin.cid = "cid-legacy-other"
	legacyOtherAdmin.instrument = splice_api_token_holding_v1.InstrumentId{Id: "LINK", Admin: "someoneElse::1220bb"}
	pkgs := []string{testPackageID}
	otherToken := genuine("cid-other", "dave", "7.0", false)
	otherToken.instrument = splice_api_token_holding_v1.InstrumentId{Id: "USDX", Admin: testRegistrar}
	clean := []holdingEvent{
		genuine("cid-1", "alice", "10.5", false),
		genuine("cid-2", "bob", "4.5000000000", true),
		genuine("cid-3", "alice", "0.0000000001", false),
		otherToken,
	}

	tests := []struct {
		name       string
		cfg        holdingCheckConfig
		events     []holdingEvent
		wantKinds  []string
		wantCids   []string
		wantInfo   int
		wantTotal  string
		wantHolds  int
		wantOwners int
		wantLocked int
	}{
		{
			name:      "clean registry holdings with matching supply",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true, expectedSupply: big.NewRat(150000000001, 10000000000)},
			events:    clean,
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "foreign template reporting LINK view is a finding and not counted",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true},
			events:    append([]holdingEvent{foreign}, clean...),
			wantKinds: []string{"foreign template"}, wantCids: []string{"cid-foreign"},
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "registry name and module from an unexpected package id is foreign",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true},
			events:    append([]holdingEvent{sameName}, clean...),
			wantKinds: []string{"foreign template"}, wantCids: []string{"cid-same-name"},
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "an older accepted registry package id counts",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: []string{testPackageID, "f00d"}, postMigration: true},
			events:    []holdingEvent{sameName},
			wantTotal: "1000.0000000000", wantHolds: 1, wantOwners: 1,
		},
		{
			name:      "a failed holding view is a finding",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true},
			events:    append([]holdingEvent{viewFailed}, clean...),
			wantKinds: []string{"holding view failed"}, wantCids: []string{"cid-view-failed"},
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "legacy holdings after migration are findings",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true},
			events:    append([]holdingEvent{legacy, legacyLocked, legacyOtherAdmin}, clean...),
			wantKinds: []string{"legacy LINK not drained", "legacy LINK not drained"}, wantCids: []string{"cid-legacy", "cid-legacy-locked"},
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "legacy holdings before migration are info only",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs},
			events:    append([]holdingEvent{legacy, legacyLocked}, clean...),
			wantInfo:  2,
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
		{
			name:      "supply mismatch",
			cfg:       holdingCheckConfig{instrument: testLink, registryPackageIDs: pkgs, postMigration: true, expectedSupply: big.NewRat(15, 1)},
			events:    clean,
			wantKinds: []string{"supply mismatch"}, wantCids: []string{"-"},
			wantTotal: "15.0000000001", wantHolds: 3, wantOwners: 2, wantLocked: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := checkHoldings(tc.cfg, tc.events)
			require.NoError(t, err)
			var kinds, cids []string
			for _, f := range r.findings {
				kinds = append(kinds, f.kind)
				cids = append(cids, f.cid)
			}
			require.Equal(t, tc.wantKinds, kinds)
			require.Equal(t, tc.wantCids, cids)
			require.Len(t, r.info, tc.wantInfo)
			require.Equal(t, tc.wantTotal, r.total.FloatString(damlDecimalFractionDigits))
			require.Equal(t, tc.wantHolds, r.holdings)
			require.Equal(t, tc.wantOwners, r.owners)
			require.Equal(t, tc.wantLocked, r.locked)
		})
	}
}

func TestCheckHoldingsRejectsBadAmount(t *testing.T) {
	t.Parallel()
	_, err := checkHoldings(holdingCheckConfig{instrument: testLink, registryPackageIDs: []string{testPackageID}}, []holdingEvent{genuine("cid-bad", "alice", "1e", false)})
	require.ErrorContains(t, err, "cid-bad")
}

func TestHoldingEventsFromViewsReportsFailedView(t *testing.T) {
	t.Parallel()
	active := []*apiv2.ActiveContract{{CreatedEvent: &apiv2.CreatedEvent{
		ContractId:     "cid-1",
		InterfaceViews: []*apiv2.InterfaceView{{ViewStatus: &status.Status{Code: 3, Message: "interpretation error"}}},
	}}}
	events, err := holdingEventsFromViews(active)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "cid-1", events[0].cid)
	require.Equal(t, "status 3: interpretation error", events[0].viewError)
}
