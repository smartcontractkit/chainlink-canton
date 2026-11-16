package linktokenadmin

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	spliceHolding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
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
	// counterparty is the recipient or holder named by the authorization.
	counterparty string
	amount       types.NUMERIC
}

type cleanupOperation = *operations.Operation[contract.ChoiceInput[types.PARTY], contract.ExerciseOutput, canton.Chain]

type cleanupCommand = func(contractID string, submitter types.PARTY) *model.ExerciseCommand

type kindSpec struct {
	template core.Template
	// decode is set for the kinds Execute supports.
	decode  func(*apiv2.CreatedEvent) (authorization, error)
	cleanup cleanupCommand
}

var cleanupOps = func() map[Kind]cleanupOperation {
	ops := make(map[Kind]cleanupOperation, len(kinds))
	for kind, spec := range kinds {
		ops[kind] = contract.NewExercise(contract.ExerciseParams[types.PARTY]{
			Name:        fmt.Sprintf("canton/link_token_admin/cleanup_expired_%s_authorization", kind),
			Version:     Version,
			Description: fmt.Sprintf("Archives an expired LinkTokenAdmin %s authorization", kind),
			Template:    spec.template,
			Method:      spec.cleanup,
		})
	}

	return ops
}()

var kinds = map[Kind]kindSpec{
	KindMint: {
		template: tokenadmin.LinkMintAuthorization{},
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
	},
	KindBurn: {
		template: tokenadmin.LinkBurnAuthorization{},
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
	},
	KindTransfer: {
		template: tokenadmin.LinkTransferAuthorization{},
		cleanup: func(cid string, submitter types.PARTY) *model.ExerciseCommand {
			return tokenadmin.LinkTransferAuthorization{}.CleanupExpiredTransferAuthorization(cid, tokenadmin.CleanupExpiredTransferAuthorization{Submitter: submitter})
		},
	},
	KindAllocate: {
		template: tokenadmin.LinkAllocationAuthorization{},
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
