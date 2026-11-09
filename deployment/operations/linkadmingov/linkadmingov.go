// Package linkadmingov holds the MCMS-governed calls on LinkTokenAdmin and its authorizations.
package linkadmingov

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/model"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("LinkTokenAdmin")

var Version = semver.MustParse("1.0.0")

// Kind names one of the four authorization templates LinkTokenAdmin creates.
type Kind string

const (
	KindMint     Kind = "mint"
	KindBurn     Kind = "burn"
	KindTransfer Kind = "transfer"
	KindAllocate Kind = "allocate"
)

// Every authorization template exposes this one MCMS function, whatever its Daml choice is called.
const withdrawFunction = "WithdrawAuthorization"

var adminEncoder = tokenadmin.NewContract("", "Link.TokenAdmin", "LinkTokenAdmin").Encoder()

var ApproveMint = newAdminCall("approve_mint", "Approves a LINK mint on LinkTokenAdmin",
	tokenadmin.LinkTokenAdmin{}.ApproveMint, adminEncoder.ApproveMint)

var ApproveBurn = newAdminCall("approve_burn", "Approves a LINK burn on LinkTokenAdmin",
	tokenadmin.LinkTokenAdmin{}.ApproveBurn, adminEncoder.ApproveBurn)

var ApproveTransfer = newAdminCall("approve_transfer", "Approves a transfer of ccipOwner-held LINK on LinkTokenAdmin",
	tokenadmin.LinkTokenAdmin{}.ApproveTransfer, adminEncoder.ApproveTransfer)

var ApproveAllocate = newAdminCall("approve_allocate", "Approves an allocation of ccipOwner-held LINK on LinkTokenAdmin",
	tokenadmin.LinkTokenAdmin{}.ApproveAllocate, adminEncoder.ApproveAllocate)

var SetPaused = newAdminCall("set_paused", "Pauses or unpauses LinkTokenAdmin",
	tokenadmin.LinkTokenAdmin{}.SetPaused, adminEncoder.SetPaused)

var SetObservers = newAdminCall("set_observers", "Replaces the LinkTokenAdmin observers",
	tokenadmin.LinkTokenAdmin{}.SetObservers, adminEncoder.SetObservers)

var WithdrawMintAuthorization = newWithdraw(KindMint, "LinkMintAuthorization",
	tokenadmin.LinkMintAuthorization{}, tokenadmin.LinkMintAuthorization{}.WithdrawMintAuthorization)

var WithdrawBurnAuthorization = newWithdraw(KindBurn, "LinkBurnAuthorization",
	tokenadmin.LinkBurnAuthorization{}, tokenadmin.LinkBurnAuthorization{}.WithdrawBurnAuthorization)

var WithdrawTransferAuthorization = newWithdraw(KindTransfer, "LinkTransferAuthorization",
	tokenadmin.LinkTransferAuthorization{}, tokenadmin.LinkTransferAuthorization{}.WithdrawTransferAuthorization)

var WithdrawAllocationAuthorization = newWithdraw(KindAllocate, "LinkAllocationAuthorization",
	tokenadmin.LinkAllocationAuthorization{}, tokenadmin.LinkAllocationAuthorization{}.WithdrawAllocationAuthorization)

type exerciseOp[ARGS any] = operations.Operation[contract.ChoiceInput[ARGS], contract.ExerciseOutput, canton.Chain]

func newAdminCall[ARGS any](
	name, description string,
	method func(contractID string, args ARGS) *model.ExerciseCommand,
	encode func(args ARGS) (*bind.EncodedChoice, error),
) *exerciseOp[ARGS] {
	return contract.NewExercise(contract.ExerciseParams[ARGS]{
		Name:         "canton/link_token_admin/" + name,
		Version:      Version,
		Description:  description,
		ContractType: ContractType,
		Template:     tokenadmin.LinkTokenAdmin{},
		Method:       method,
		EncodeMethod: encode,
	})
}

func newWithdraw[ARGS any](
	kind Kind,
	contractType deployment.ContractType,
	template core.Template,
	method func(contractID string, args ARGS) *model.ExerciseCommand,
) *exerciseOp[ARGS] {
	return contract.NewExercise(contract.ExerciseParams[ARGS]{
		Name:         "canton/link_token_admin/withdraw_" + string(kind) + "_authorization",
		Version:      Version,
		Description:  "Withdraws an unexecuted LinkTokenAdmin " + string(kind) + " authorization",
		ContractType: contractType,
		Template:     template,
		Method:       method,
		EncodeMethod: func(ARGS) (*bind.EncodedChoice, error) {
			return &bind.EncodedChoice{Choice: withdrawFunction}, nil
		},
	})
}
