package testhelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	splice_api_token_transfer_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

// ListedHolding is one active HoldingV1 contract that matched list filters, with parsed amount.
type ListedHolding struct {
	ContractID string
	View       splice_api_token_holding_v1.HoldingView
	Amount     *big.Rat
}

// Filter selects HoldingV1 rows using the ledger contract id (created event) and the interface view.
type Filter func(contractID string, hv splice_api_token_holding_v1.HoldingView) bool

// WithHoldingOwner restricts the sum to holdings whose interface view owner equals party.
func WithHoldingOwner(party string) Filter {
	p := strings.TrimSpace(party)
	if p == "" {
		return func(string, splice_api_token_holding_v1.HoldingView) bool { return true }
	}

	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return string(hv.Owner) == p
	}
}

// WithUnlockedHoldingsOnly skips holdings with a non-nil Lock in the interface view.
func WithUnlockedHoldingsOnly() Filter {
	return func(_ string, hv splice_api_token_holding_v1.HoldingView) bool {
		return hv.Lock == nil
	}
}

// WithUnlockedHoldingsOnly skips holdings with a non-nil Lock in the interface view.
func ExcludeCIDs(cids []string) Filter {
	return func(contractID string, _ splice_api_token_holding_v1.HoldingView) bool {
		return !slices.Contains(cids, contractID)
	}
}

// GetHoldingsBalance returns the sum of holding amounts in ledger Numeric units (*big.Rat).
// If instrument is nil, amounts for all instruments are summed (filters still apply).
func GetHoldingsBalance(
	ctx context.Context,
	participant canton.Participant,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) (*big.Rat, error) {
	rows, err := listHoldingsForInstrument(ctx, participant, instrument, filters...)
	if err != nil {
		return nil, err
	}

	total := new(big.Rat)
	for _, row := range rows {
		total.Add(total, row.Amount)
	}

	return total, nil
}

// PickHoldingsForInstrument returns len(minAmounts) distinct HoldingV1 rows for instrument.
// The i-th returned row corresponds to minAmounts[i]: each chosen holding has Amount >= minAmounts[i]
// (nil entries in minAmounts are treated as zero minimum).
//
// Callers choose how many holdings they need by the length of minAmounts (e.g. one entry for a single fee leg,
// two entries for two same-instrument legs). For another instrument, call again with that instrument’s ID and filters.
//
// filters should narrow the ACS row set (e.g. WithHoldingOwner(party), WithUnlockedHoldingsOnly(), WithHoldingContractID(cid)).
// Candidates are sorted by amount descending; slots with larger minimums are satisfied first, then each slot
// gets the largest still-unused holding that meets its minimum.
func PickHoldingsForInstrument(
	ctx context.Context,
	participant canton.Participant,
	instrument splice_api_token_holding_v1.InstrumentId,
	minAmounts []*big.Rat,
	filters ...Filter,
) ([]ListedHolding, error) {
	if len(minAmounts) == 0 {
		return nil, fmt.Errorf("minAmounts must be non-empty (each entry is one required holding)")
	}

	rows, err := listHoldingsForInstrument(ctx, participant, &instrument, filters...)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(rows, func(a, b ListedHolding) int {
		return b.Amount.Cmp(a.Amount)
	})

	type slot struct {
		idx int
		min *big.Rat
	}
	slots := make([]slot, len(minAmounts))
	for i := range minAmounts {
		m := minAmounts[i]
		if m == nil {
			m = big.NewRat(0, 1)
		}
		slots[i] = slot{idx: i, min: m}
	}
	slices.SortFunc(slots, func(a, b slot) int {
		return b.min.Cmp(a.min)
	})

	out := make([]ListedHolding, len(minAmounts))
	used := make(map[string]struct{})
	for _, s := range slots {
		var chosen ListedHolding
		var found bool
		for _, row := range rows {
			if _, ok := used[row.ContractID]; ok {
				continue
			}
			if row.Amount.Cmp(s.min) < 0 {
				continue
			}
			chosen, found = row, true
			break
		}
		if !found {
			return nil, fmt.Errorf(
				"could not pick holding for slot %d (min=%s): insufficient distinct holdings for instrument",
				s.idx, s.min.FloatString(8),
			)
		}
		used[chosen.ContractID] = struct{}{}
		out[s.idx] = chosen
	}

	return out, nil
}

// SplitHoldingSelfTransfer splits holding into two equal halves via a self-transfer (same sender and
// receiver). The parent holding is consumed; the submit transaction must contain exactly two new
// Splice.Amulet contracts for the same owner and instrument (parsed from createArguments).
//
// registryAdmin is the instrument registry admin (see GetRegistryAdmin).
func SplitHoldingSelfTransfer(
	ctx context.Context,
	participant canton.Participant,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	registryAdmin string,
	holding ListedHolding,
) (ListedHolding, ListedHolding, error) {
	if holding.Amount == nil || holding.Amount.Sign() <= 0 {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("holding amount must be positive")
	}

	owner := strings.TrimSpace(string(holding.View.Owner))
	if owner == "" {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("holding has empty owner")
	}

	half := new(big.Rat).Quo(new(big.Rat).Set(holding.Amount), big.NewRat(2, 1))
	halfStr := half.FloatString(10)
	if _, ok := new(big.Rat).SetString(halfStr); !ok {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("split: invalid half numeric %q", halfStr)
	}

	now := time.Now()
	transferParams := splice_api_token_transfer_instruction_v1.Transfer{
		Sender:        types.PARTY(owner),
		Receiver:      types.PARTY(owner),
		Amount:        types.NUMERIC(halfStr),
		InstrumentId:  holding.View.InstrumentId,
		RequestedAt:   types.TIMESTAMP(now.Add(-time.Hour)),
		ExecuteBefore: types.TIMESTAMP(now.Add(24 * time.Hour)),
		InputHoldingCids: []types.CONTRACT_ID{
			types.CONTRACT_ID(holding.ContractID),
		},
		Meta: splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	}

	tf, err := GetTransferFactoryV2(ctx, transferInstructionClient, registryAdmin, transferParams)
	if err != nil {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("get transfer factory: %w", err)
	}

	extraArgs := splice_api_token_metadata_v1.ExtraArgs{
		Context: splice_api_token_metadata_v1.ChoiceContext{Values: map[string]splice_api_token_metadata_v1.AnyValue{}},
		Meta:    splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	}
	if len(tf.ChoiceContextData) > 0 {
		raw, mErr := json.Marshal(tf.ChoiceContextData)
		if mErr != nil {
			return ListedHolding{}, ListedHolding{}, fmt.Errorf("marshal choice context: %w", mErr)
		}
		if uErr := json.Unmarshal(raw, &extraArgs.Context); uErr != nil {
			return ListedHolding{}, ListedHolding{}, fmt.Errorf("unmarshal choice context: %w", uErr)
		}
	}

	exercise := splice_api_token_transfer_instruction_v1.TransferFactoryTransfer{
		ExpectedAdmin: types.PARTY(registryAdmin),
		Transfer:      transferParams,
		ExtraArgs:     extraArgs,
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{
						PackageId:  "#" + splice_api_token_transfer_instruction_v1.PackageName,
						ModuleName: "Splice.Api.Token.TransferInstructionV1",
						EntityName: "TransferFactory",
					},
					ContractId:     tf.FactoryID,
					Choice:         "TransferFactory_Transfer",
					ChoiceArgument: ledger.MapToValue(exercise),
				}},
			}},
			ActAs:              []string{owner},
			DisclosedContracts: tf.DisclosedContracts,
		},
	})
	if err != nil {
		return ListedHolding{}, ListedHolding{}, fmt.Errorf("submit transfer: %w", err)
	}

	debugPrintSplitHoldingSubmitResponseJSON(res)

	children, err := splitChildrenFromSelfTransferTransaction(res.GetTransaction(), owner, holding.View.InstrumentId)
	if err != nil {
		return ListedHolding{}, ListedHolding{}, err
	}

	return children[0], children[1], nil
}

// debugPrintSplitHoldingSubmitResponseJSON prints SubmitAndWaitForTransactionResponse as JSON (temporary debugging).
// Uses gogo jsonpb: dazl ledger protos implement github.com/gogo/protobuf/proto.Message, not google.golang.org/protobuf.
func debugPrintSplitHoldingSubmitResponseJSON(res proto.Message) {
	m := jsonpb.Marshaler{Indent: "  ", EmitDefaults: true}
	s, err := m.MarshalToString(res)
	if err != nil {
		fmt.Printf("debugPrintSplitHoldingSubmitResponseJSON: %v\n", err)
		return
	}
	fmt.Printf("debugPrintSplitHoldingSubmitResponseJSON:\n%s\n", s)
}

func splitChildrenFromSelfTransferTransaction(
	tx *apiv2.Transaction,
	owner string,
	instrument splice_api_token_holding_v1.InstrumentId,
) ([]ListedHolding, error) {
	if tx == nil {
		return nil, fmt.Errorf("split transfer: nil transaction")
	}

	var out []ListedHolding
	for _, event := range tx.GetEvents() {
		created := event.GetCreated()
		if created == nil {
			continue
		}
		if strings.TrimSpace(created.GetContractId()) == "" {
			continue
		}
		row, matched, err := listedHoldingFromSpliceAmuletCreatedEvent(created, owner, instrument)
		if err != nil {
			return nil, fmt.Errorf("split transfer: %w", err)
		}
		if matched {
			out = append(out, row)
		}
	}
	if len(out) != 2 {
		return nil, fmt.Errorf("split transfer: expected exactly 2 Amulet child creates (got %d)", len(out))
	}
	return out, nil
}

// listedHoldingFromSpliceAmuletCreatedEvent builds ListedHolding from Splice.Amulet:Amulet createArguments.
func listedHoldingFromSpliceAmuletCreatedEvent(
	created *apiv2.CreatedEvent,
	owner string,
	instrument splice_api_token_holding_v1.InstrumentId,
) (ListedHolding, bool, error) {
	tid := created.GetTemplateId()
	if tid == nil || tid.GetModuleName() != "Splice.Amulet" || tid.GetEntityName() != "Amulet" {
		return ListedHolding{}, false, nil
	}
	args := created.GetCreateArguments()
	if args == nil {
		return ListedHolding{}, false, nil
	}

	var dsoParty, ownerParty, initialAmount string
	for _, f := range args.GetFields() {
		switch f.GetLabel() {
		case "dso":
			dsoParty = strings.TrimSpace(f.GetValue().GetParty())
		case "owner":
			ownerParty = strings.TrimSpace(f.GetValue().GetParty())
		case "amount":
			expiring := f.GetValue().GetRecord()
			if expiring == nil {
				continue
			}
			for _, ef := range expiring.GetFields() {
				if ef.GetLabel() == "initialAmount" {
					initialAmount = strings.TrimSpace(ef.GetValue().GetNumeric())
				}
			}
		}
	}
	if ownerParty != owner {
		return ListedHolding{}, false, nil
	}
	if strings.TrimSpace(string(instrument.Admin)) != dsoParty {
		return ListedHolding{}, false, nil
	}

	amt, ok := new(big.Rat).SetString(initialAmount)
	if !ok {
		return ListedHolding{}, false, fmt.Errorf("amulet create: invalid initialAmount %q", initialAmount)
	}

	cid := strings.TrimSpace(created.GetContractId())
	hv := splice_api_token_holding_v1.HoldingView{
		Owner:        types.PARTY(ownerParty),
		InstrumentId: instrument,
		Amount:       types.NUMERIC(initialAmount),
		Lock:         nil,
		Meta:         splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	}

	return ListedHolding{
		ContractID: cid,
		View:       hv,
		Amount:     new(big.Rat).Set(amt),
	}, true, nil
}

// listHoldingsForInstrument walks active HoldingV1 contracts, applies optional instrument and filters,
// and returns rows with a non-empty contract id and a parseable amount (including zero).
// If instrument is nil, all instruments match (same semantics as GetHoldingsBalance).
func listHoldingsForInstrument(
	ctx context.Context,
	participant canton.Participant,
	instrument *splice_api_token_holding_v1.InstrumentId,
	filters ...Filter,
) ([]ListedHolding, error) {
	contracts, err := ListActiveContractsByInterfaceId(ctx, participant, holdingV1InterfaceID())
	if err != nil {
		return nil, err
	}

	var out []ListedHolding
	for _, ac := range contracts {
		created := ac.GetCreatedEvent()
		if created == nil || strings.TrimSpace(created.GetContractId()) == "" {
			continue
		}
		contractID := strings.TrimSpace(created.GetContractId())

		hv, ok, err := parseHoldingV1View(ac.GetCreatedEvent())
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if instrument != nil && hv.InstrumentId != *instrument {
			continue
		}
		if !holdingPassesFilters(contractID, hv, filters) {
			continue
		}

		amountRaw := strings.TrimSpace(string(hv.Amount))
		if amountRaw == "" {
			continue
		}
		amt, ok := new(big.Rat).SetString(amountRaw)
		if !ok {
			return nil, fmt.Errorf("invalid holding amount %q", amountRaw)
		}

		out = append(out, ListedHolding{
			ContractID: contractID,
			View:       hv,
			Amount:     amt,
		})
	}

	return out, nil
}

func holdingPassesFilters(contractID string, hv splice_api_token_holding_v1.HoldingView, filters []Filter) bool {
	for _, f := range filters {
		if !f(contractID, hv) {
			return false
		}
	}

	return true
}

func parseHoldingV1View(created *apiv2.CreatedEvent) (splice_api_token_holding_v1.HoldingView, bool, error) {
	if created == nil {
		return splice_api_token_holding_v1.HoldingView{}, false, nil
	}
	want := holdingV1InterfaceID()
	for _, iv := range created.GetInterfaceViews() {
		iface := iv.GetInterfaceId()
		if iface.GetModuleName() != want.GetModuleName() || iface.GetEntityName() != want.GetEntityName() {
			continue
		}
		vv := iv.GetViewValue()
		if vv == nil {
			continue
		}
		var hv splice_api_token_holding_v1.HoldingView
		if err := ledger.RecordToStruct(vv, &hv); err != nil {
			return splice_api_token_holding_v1.HoldingView{}, false, fmt.Errorf("parse HoldingV1 view: %w", err)
		}

		return hv, true, nil
	}

	return splice_api_token_holding_v1.HoldingView{}, false, nil
}

// holdingV1InterfaceID is the canonical HoldingV1 ledger interface id (single place for module/entity strings).
func holdingV1InterfaceID() *apiv2.Identifier {
	return &apiv2.Identifier{
		PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}
}
