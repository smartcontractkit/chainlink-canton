package sequences

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// LinkTokenAdminAction is one governed call. Exactly one field must be set.
type LinkTokenAdminAction struct {
	ApproveMint     *tokenadmin.ApproveMint     `json:"approveMint,omitempty" yaml:"approveMint,omitempty"`
	ApproveBurn     *tokenadmin.ApproveBurn     `json:"approveBurn,omitempty" yaml:"approveBurn,omitempty"`
	ApproveTransfer *tokenadmin.ApproveTransfer `json:"approveTransfer,omitempty" yaml:"approveTransfer,omitempty"`
	ApproveAllocate *tokenadmin.ApproveAllocate `json:"approveAllocate,omitempty" yaml:"approveAllocate,omitempty"`
	SetPaused       *tokenadmin.SetPaused       `json:"setPaused,omitempty" yaml:"setPaused,omitempty"`
	SetObservers    *tokenadmin.SetObservers    `json:"setObservers,omitempty" yaml:"setObservers,omitempty"`
	Withdraw        *LinkAuthorizationRef       `json:"withdraw,omitempty" yaml:"withdraw,omitempty"`
}

// LinkAuthorizationRef names an authorization by kind and the authorizationId it was approved with.
type LinkAuthorizationRef struct {
	Kind            linkadmingov.Kind `json:"kind" yaml:"kind"`
	AuthorizationID string            `json:"authorizationId" yaml:"authorizationId"`
}

type GovernLinkTokenAdminInput struct {
	// AdminAddress is the LinkTokenAdmin "instanceId@ccipOwner"; authorizations share its owner.
	AdminAddress   contracts.RawInstanceAddress `json:"adminAddress" yaml:"adminAddress"`
	ProposalDriven bool                         `json:"proposalDriven,omitempty" yaml:"proposalDriven,omitempty"`
	Actions        []LinkTokenAdminAction       `json:"actions" yaml:"actions"`
}

// GovernLinkTokenAdmin runs the actions in order. When ProposalDriven is true it returns them as one
// MCMS batch, so they execute atomically and later actions see earlier ones (e.g. unpause, then approve).
var GovernLinkTokenAdmin = operations.NewSequence(
	"canton/link_token_admin/govern",
	semver.MustParse("1.0.0"),
	"Approves, pauses, sets observers on, or withdraws authorizations of LinkTokenAdmin",
	func(b operations.Bundle, deps canton.Chain, input GovernLinkTokenAdminInput) (sequences.OnChainOutput, error) {
		if err := input.Validate(); err != nil {
			return sequences.OnChainOutput{}, err
		}

		outs := make([]contract.ExerciseOutput, 0, len(input.Actions))
		for i, action := range input.Actions {
			out, err := action.exercise(b, deps, input.AdminAddress, input.ProposalDriven)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("action %d: %w", i, err)
			}
			outs = append(outs, out)
		}

		if !input.ProposalDriven {
			return sequences.OnChainOutput{}, nil
		}
		batchOp, err := contract.NewBatchOperationFromExercises(outs)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for link token admin: %w", err)
		}

		return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
	},
)

func (in GovernLinkTokenAdminInput) Validate() error {
	if _, err := contracts.RawInstanceAddressFromString(in.AdminAddress.String()); err != nil {
		return fmt.Errorf("adminAddress: %w", err)
	}
	if len(in.Actions) == 0 {
		return errors.New("at least one action is required")
	}
	approved := make(map[types.TEXT]bool)
	approvedKind := make(map[string]bool)
	withdrawn := make(map[string]bool)
	for i, action := range in.Actions {
		if err := action.validate(); err != nil {
			return fmt.Errorf("action %d: %w", i, err)
		}
		if id, _, _, ok := action.approval(); ok {
			// Two live authorizations with one id make every later withdraw of that id ambiguous.
			// This covers one proposal; GovernLinkTokenAdmin.VerifyPreconditions checks the ledger.
			key := string(action.approvalKind()) + "/" + string(id)
			if approvedKind[key] {
				return fmt.Errorf("action %d: authorization %q is approved twice in this proposal", i, id)
			}
			approvedKind[key] = true
			approved[id] = true
		}
		if action.Withdraw == nil {
			continue
		}
		// The executor resolves every target before the batch runs, so the new authorization is not found yet.
		if approved[types.TEXT(action.Withdraw.AuthorizationID)] {
			return fmt.Errorf("action %d: authorization %q is approved earlier in this proposal; withdraw it in a later one", i, action.Withdraw.AuthorizationID)
		}
		// Every kind shares the address id@ccipOwner, and the executor maps one contract per address.
		key := action.Withdraw.AuthorizationID
		if withdrawn[key] {
			return fmt.Errorf("action %d: authorization %q is withdrawn twice in this proposal; withdraw it in a later one", i, action.Withdraw.AuthorizationID)
		}
		withdrawn[key] = true
	}

	return nil
}

func (a LinkTokenAdminAction) validate() error {
	set := 0
	for _, isSet := range []bool{
		a.ApproveMint != nil, a.ApproveBurn != nil, a.ApproveTransfer != nil, a.ApproveAllocate != nil,
		a.SetPaused != nil, a.SetObservers != nil, a.Withdraw != nil,
	} {
		if isSet {
			set++
		}
	}
	if set != 1 {
		return fmt.Errorf("exactly one call must be set, got %d", set)
	}

	if id, amount, validForSeconds, ok := a.approval(); ok {
		if err := validateParty(a.approvalCounterparty()); err != nil {
			return err
		}

		return validateApproval(id, amount, validForSeconds)
	}
	if a.SetObservers != nil {
		for _, p := range a.SetObservers.Observers {
			if err := validateParty(p); err != nil {
				return fmt.Errorf("setObservers: %w", err)
			}
		}

		return nil
	}
	if a.Withdraw == nil {
		return nil
	}
	if err := validateAuthorizationID(a.Withdraw.AuthorizationID); err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}
	switch a.Withdraw.Kind {
	case linkadmingov.KindMint, linkadmingov.KindBurn, linkadmingov.KindTransfer, linkadmingov.KindAllocate:
		return nil
	default:
		return fmt.Errorf("withdraw: unknown authorization kind %q", a.Withdraw.Kind)
	}
}

func (a LinkTokenAdminAction) approval() (types.TEXT, types.NUMERIC, types.INT64, bool) {
	switch {
	case a.ApproveMint != nil:
		return a.ApproveMint.AuthorizationId, a.ApproveMint.Amount, a.ApproveMint.ValidForSeconds, true
	case a.ApproveBurn != nil:
		return a.ApproveBurn.AuthorizationId, a.ApproveBurn.MaxAmount, a.ApproveBurn.ValidForSeconds, true
	case a.ApproveTransfer != nil:
		return a.ApproveTransfer.AuthorizationId, a.ApproveTransfer.Amount, a.ApproveTransfer.ValidForSeconds, true
	case a.ApproveAllocate != nil:
		return a.ApproveAllocate.AuthorizationId, a.ApproveAllocate.Amount, a.ApproveAllocate.ValidForSeconds, true
	}

	return "", "", 0, false
}

// ApprovalRef names the authorization an Approve* action creates.
func (a LinkTokenAdminAction) ApprovalRef() (LinkAuthorizationRef, bool) {
	id, _, _, ok := a.approval()
	if !ok {
		return LinkAuthorizationRef{}, false
	}

	return LinkAuthorizationRef{Kind: a.approvalKind(), AuthorizationID: string(id)}, true
}

func (a LinkTokenAdminAction) approvalKind() linkadmingov.Kind {
	switch {
	case a.ApproveMint != nil:
		return linkadmingov.KindMint
	case a.ApproveBurn != nil:
		return linkadmingov.KindBurn
	case a.ApproveTransfer != nil:
		return linkadmingov.KindTransfer
	}

	return linkadmingov.KindAllocate
}

func (a LinkTokenAdminAction) approvalCounterparty() types.PARTY {
	switch {
	case a.ApproveMint != nil:
		return a.ApproveMint.Recipient
	case a.ApproveBurn != nil:
		return a.ApproveBurn.Holder
	case a.ApproveTransfer != nil:
		return a.ApproveTransfer.Receiver
	case a.ApproveAllocate != nil:
		return a.ApproveAllocate.Receiver
	}

	return ""
}

// validateParty rejects an empty or malformed party, which the Daml decoder would reject only at execution.
func validateParty(party types.PARTY) error {
	if !strings.Contains(string(party), "::") {
		return fmt.Errorf("party %q must be a full party id (hint::fingerprint)", party)
	}

	return nil
}

// Daml Decimal is Numeric 10. The MCMS encoder truncates extra fractional digits instead of rejecting them.
var damlDecimal = regexp.MustCompile(`^[0-9]{1,28}(\.[0-9]{1,10})?$`)

// validateApproval rejects what LinkTokenAdmin would reject, before a proposal is signed.
func validateApproval(authorizationID types.TEXT, amount types.NUMERIC, validForSeconds types.INT64) error {
	if err := validateAuthorizationID(string(authorizationID)); err != nil {
		return err
	}
	value, ok := new(big.Rat).SetString(string(amount))
	if !damlDecimal.MatchString(string(amount)) || !ok || value.Sign() <= 0 {
		return fmt.Errorf("amount %q must be a positive decimal with at most 10 fractional digits", amount)
	}
	if validForSeconds <= 0 {
		return errors.New("validForSeconds must be positive")
	}

	return nil
}

func validateAuthorizationID(id string) error {
	if id == "" {
		return errors.New("authorizationId is required")
	}
	if strings.Contains(id, "@") {
		return fmt.Errorf("authorizationId %q must not contain '@'", id)
	}

	return nil
}

func (a LinkTokenAdminAction) exercise(
	b operations.Bundle,
	chain canton.Chain,
	admin contracts.RawInstanceAddress,
	proposalDriven bool,
) (contract.ExerciseOutput, error) {
	switch {
	case a.ApproveMint != nil:
		return exerciseOn(b, chain, linkadmingov.ApproveMint, admin, *a.ApproveMint, proposalDriven)
	case a.ApproveBurn != nil:
		return exerciseOn(b, chain, linkadmingov.ApproveBurn, admin, *a.ApproveBurn, proposalDriven)
	case a.ApproveTransfer != nil:
		return exerciseOn(b, chain, linkadmingov.ApproveTransfer, admin, *a.ApproveTransfer, proposalDriven)
	case a.ApproveAllocate != nil:
		return exerciseOn(b, chain, linkadmingov.ApproveAllocate, admin, *a.ApproveAllocate, proposalDriven)
	case a.SetPaused != nil:
		return exerciseOn(b, chain, linkadmingov.SetPaused, admin, *a.SetPaused, proposalDriven)
	case a.SetObservers != nil:
		return exerciseOn(b, chain, linkadmingov.SetObservers, admin, *a.SetObservers, proposalDriven)
	}

	auth := contracts.InstanceID(a.Withdraw.AuthorizationID).RawInstanceAddress(types.PARTY(admin.Owner()))
	switch a.Withdraw.Kind {
	case linkadmingov.KindMint:
		return exerciseOn(b, chain, linkadmingov.WithdrawMintAuthorization, auth, tokenadmin.WithdrawMintAuthorization{}, proposalDriven)
	case linkadmingov.KindBurn:
		return exerciseOn(b, chain, linkadmingov.WithdrawBurnAuthorization, auth, tokenadmin.WithdrawBurnAuthorization{}, proposalDriven)
	case linkadmingov.KindTransfer:
		return exerciseOn(b, chain, linkadmingov.WithdrawTransferAuthorization, auth, tokenadmin.WithdrawTransferAuthorization{}, proposalDriven)
	case linkadmingov.KindAllocate:
		return exerciseOn(b, chain, linkadmingov.WithdrawAllocationAuthorization, auth, tokenadmin.WithdrawAllocationAuthorization{}, proposalDriven)
	}

	return contract.ExerciseOutput{}, fmt.Errorf("unknown authorization kind %q", a.Withdraw.Kind)
}

func exerciseOn[ARGS any](
	b operations.Bundle,
	chain canton.Chain,
	op *operations.Operation[contract.ChoiceInput[ARGS], contract.ExerciseOutput, canton.Chain],
	target contracts.RawInstanceAddress,
	args ARGS,
	proposalDriven bool,
) (contract.ExerciseOutput, error) {
	report, err := operations.ExecuteOperation(b, op, chain, newChoiceInput(target, args, proposalDriven))

	return report.Output, err
}
