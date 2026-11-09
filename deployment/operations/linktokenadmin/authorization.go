package linktokenadmin

import (
	"fmt"
	"math/big"
	"sort"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	spliceHolding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
)

// Kind names one of the four LinkTokenAdmin authorization templates.
type Kind string

const (
	KindMint     Kind = "mint"
	KindBurn     Kind = "burn"
	KindTransfer Kind = "transfer"
	KindAllocate Kind = "allocate"
)

// authorization is the part of a Link*Authorization contract that Execute needs.
type authorization struct {
	contractID      string
	ccipOwner       string
	adminInstanceID string
	instrumentID    spliceHolding.InstrumentId
	// counterparty is the recipient, holder or receiver named by the authorization.
	counterparty string
	// amount is the approved amount; for a burn it is maxAmount.
	amount types.NUMERIC
}

type kindSpec struct {
	templateID string
	decode     func(*apiv2.CreatedEvent) (authorization, error)
	cleanup    func(contractID string, submitter types.PARTY) *model.ExerciseCommand
	// coSigned reports whether the counterparty is a controller of Execute when it is not ccipOwner.
	coSigned bool
	// holdingsOwnedByCounterparty is true when input holdings belong to the counterparty, false
	// when they belong to ccipOwner; ignored for a mint, which takes no inputs.
	holdingsOwnedByCounterparty bool
}

var kinds = map[Kind]kindSpec{
	KindMint: {
		templateID: tokenadmin.LinkMintAuthorization{}.GetTemplateID(),
		decode: func(ev *apiv2.CreatedEvent) (authorization, error) {
			a, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkMintAuthorization](ev)
			if err != nil {
				return authorization{}, err
			}

			return authorization{ev.GetContractId(), string(a.CcipOwner), string(a.AdminInstanceId), a.InstrumentId, string(a.Recipient), a.Amount}, nil
		},
		cleanup: func(cid string, submitter types.PARTY) *model.ExerciseCommand {
			return tokenadmin.LinkMintAuthorization{}.CleanupExpiredMintAuthorization(cid, tokenadmin.CleanupExpiredMintAuthorization{Submitter: submitter})
		},
		coSigned: true,
	},
	KindBurn: {
		templateID: tokenadmin.LinkBurnAuthorization{}.GetTemplateID(),
		decode: func(ev *apiv2.CreatedEvent) (authorization, error) {
			a, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkBurnAuthorization](ev)
			if err != nil {
				return authorization{}, err
			}

			return authorization{ev.GetContractId(), string(a.CcipOwner), string(a.AdminInstanceId), a.InstrumentId, string(a.Holder), a.MaxAmount}, nil
		},
		cleanup: func(cid string, submitter types.PARTY) *model.ExerciseCommand {
			return tokenadmin.LinkBurnAuthorization{}.CleanupExpiredBurnAuthorization(cid, tokenadmin.CleanupExpiredBurnAuthorization{Submitter: submitter})
		},
		coSigned:                    true,
		holdingsOwnedByCounterparty: true,
	},
	KindTransfer: {
		templateID: tokenadmin.LinkTransferAuthorization{}.GetTemplateID(),
		decode: func(ev *apiv2.CreatedEvent) (authorization, error) {
			a, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkTransferAuthorization](ev)
			if err != nil {
				return authorization{}, err
			}

			return authorization{ev.GetContractId(), string(a.CcipOwner), string(a.AdminInstanceId), a.InstrumentId, string(a.Receiver), a.Amount}, nil
		},
		cleanup: func(cid string, submitter types.PARTY) *model.ExerciseCommand {
			return tokenadmin.LinkTransferAuthorization{}.CleanupExpiredTransferAuthorization(cid, tokenadmin.CleanupExpiredTransferAuthorization{Submitter: submitter})
		},
	},
	KindAllocate: {
		templateID: tokenadmin.LinkAllocationAuthorization{}.GetTemplateID(),
		decode: func(ev *apiv2.CreatedEvent) (authorization, error) {
			a, err := bindings.UnmarshalCreatedEvent[tokenadmin.LinkAllocationAuthorization](ev)
			if err != nil {
				return authorization{}, err
			}

			return authorization{ev.GetContractId(), string(a.CcipOwner), string(a.AdminInstanceId), a.InstrumentId, string(a.Receiver), a.Amount}, nil
		},
		cleanup: func(cid string, submitter types.PARTY) *model.ExerciseCommand {
			return tokenadmin.LinkAllocationAuthorization{}.CleanupExpiredAllocationAuthorization(cid, tokenadmin.CleanupExpiredAllocationAuthorization{Submitter: submitter})
		},
	},
}

func specFor(kind Kind) (kindSpec, error) {
	spec, ok := kinds[kind]
	if !ok {
		return kindSpec{}, fmt.Errorf("unknown authorization kind %q", kind)
	}

	return spec, nil
}

// actAs is the submitter plus, for a mint or burn, the counterparty when it is not ccipOwner.
func actAs(spec kindSpec, auth authorization, submitter string) []string {
	parties := []string{submitter}
	if spec.coSigned && auth.counterparty != auth.ccipOwner && auth.counterparty != submitter {
		parties = append(parties, auth.counterparty)
	}

	return parties
}

// holdingOwner is the party whose holdings Execute consumes.
func holdingOwner(spec kindSpec, auth authorization) string {
	if spec.holdingsOwnedByCounterparty {
		return auth.counterparty
	}

	return auth.ccipOwner
}

type holding struct {
	contractID string
	amount     *big.Rat
}

// selectHoldings picks unlocked holdings largest first. A burn consumes every input and its
// authorization is spent in one go, so it must hit maxAmount exactly; a smaller burn needs
// explicit InputHoldingCids. A transfer or allocation returns change, so it stops once covered.
func selectHoldings(kind Kind, candidates []holding, amount types.NUMERIC) ([]string, error) {
	target, ok := new(big.Rat).SetString(string(amount))
	if !ok {
		return nil, fmt.Errorf("invalid authorization amount %q", amount)
	}

	sorted := append([]holding(nil), candidates...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].amount.Cmp(sorted[j].amount) > 0 })

	total := new(big.Rat)
	var picked []string
	for _, h := range sorted {
		next := new(big.Rat).Add(total, h.amount)
		if kind == KindBurn {
			if next.Cmp(target) > 0 {
				continue
			}
		} else if total.Cmp(target) >= 0 {
			break
		}
		total = next
		picked = append(picked, h.contractID)
	}

	if len(picked) == 0 {
		return nil, fmt.Errorf("no unlocked holdings fit the authorization amount %s", amount)
	}
	if kind == KindBurn && total.Cmp(target) != 0 {
		return nil, fmt.Errorf("unlocked holdings that fit sum to %s, not the authorized burn amount %s; pass InputHoldingCids to burn a chosen set", total.FloatString(10), amount)
	}
	if kind != KindBurn && total.Cmp(target) < 0 {
		return nil, fmt.Errorf("unlocked holdings total %s, below the authorization amount %s", total.FloatString(10), amount)
	}

	return picked, nil
}

// resultContractIDs walks an exercise result and returns every contract id in it, in order.
func resultContractIDs(v *apiv2.Value) []string {
	var out []string
	var walk func(*apiv2.Value)
	walk = func(v *apiv2.Value) {
		switch s := v.GetSum().(type) {
		case *apiv2.Value_ContractId:
			out = append(out, s.ContractId)
		case *apiv2.Value_List:
			for _, e := range s.List.GetElements() {
				walk(e)
			}
		case *apiv2.Value_Optional:
			if s.Optional.GetValue() != nil {
				walk(s.Optional.GetValue())
			}
		case *apiv2.Value_Record:
			for _, f := range s.Record.GetFields() {
				walk(f.GetValue())
			}
		case *apiv2.Value_Variant:
			walk(s.Variant.GetValue())
		}
	}
	if v != nil {
		walk(v)
	}

	return out
}
