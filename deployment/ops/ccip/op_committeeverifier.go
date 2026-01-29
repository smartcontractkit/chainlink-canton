package ccip

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/ccvs"
)

// ============================================================================
// Committee Verifier Operations
// ============================================================================

// CommitteeVerifierForwardToVerifierInput contains input for forwarding a message to verifier
type CommitteeVerifierForwardToVerifierInput struct {
	CommitteeVerifierContractID string
	CommitteeVerifierTemplateID string
	CcvRegistryCid              string
	Message                     ccvs.MessageV1
	MessageId                   string
	FeeToken                    ccvs.InstrumentId
	FeeTokenAmount              string // Numeric value as string
	VerifierArgs                string
}

// CommitteeVerifierForwardToVerifierOutput contains the transaction ID and CCVTicket contract ID
type CommitteeVerifierForwardToVerifierOutput struct {
	TransactionID       string
	CCVTicketContractID string
	CCVTicketTemplateID string
}

var committeeVerifierForwardToVerifierHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input CommitteeVerifierForwardToVerifierInput) (output CantonOpResult[CommitteeVerifierForwardToVerifierOutput], err error) {
	ctx := b.GetContext()

	// Parse fee token amount
	feeTokenAmount, ok := new(big.Int).SetString(input.FeeTokenAmount, 10)
	if !ok {
		return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{}, fmt.Errorf("invalid feeTokenAmount: %s", input.FeeTokenAmount)
	}
	scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	feeTokenAmountMantissa := new(big.Int).Mul(feeTokenAmount, scale10)

	// Create forward command
	forwardArgs := ccvs.CommitteeVerifierForwardToVerifier{
		CcvRegistryCid: types.CONTRACT_ID(input.CcvRegistryCid),
		Message:        input.Message,
		MessageId:      types.TEXT(input.MessageId),
		FeeToken:       input.FeeToken,
		FeeTokenAmount: types.NUMERIC(feeTokenAmountMantissa),
		VerifierArgs:   types.TEXT(input.VerifierArgs),
		Caller:         types.PARTY(deps.Party),
	}

	committeeVerifier := ccvs.CommitteeVerifier{}
	exerciseCmd := committeeVerifier.CommitteeVerifierForwardToVerifier(input.CommitteeVerifierContractID, forwardArgs)

	// List known packages to find the package ID for ccip-committeeverifier
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipCommitteeVerifierPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-committeeverifier") {
			ccipCommitteeVerifierPkgID = p.PackageID
			break
		}
	}
	if ccipCommitteeVerifierPkgID == "" {
		return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{}, fmt.Errorf("failed to find ccip-committeeverifier package")
	}

	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: commandID,
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", ccipCommitteeVerifierPkgID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{}, fmt.Errorf("failed to forward message to verifier: %w", err)
	}

	// Extract CCVTicket contract ID from Created event
	ccvTicketContractID := ""
	ccvTicketTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == "CCIP.Tickets:CCVTicket" {
			ccvTicketContractID = ev.Created.ContractID
			ccvTicketTemplateID = ev.Created.TemplateID

			break
		}
	}

	if ccvTicketContractID == "" {
		return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{}, fmt.Errorf("forward-to-verifier tx had no Created CCVTicket event")
	}

	return CantonOpResult[CommitteeVerifierForwardToVerifierOutput]{
		UpdateID: submitResp.UpdateID,
		Output: CommitteeVerifierForwardToVerifierOutput{
			TransactionID:       commandID,
			CCVTicketContractID: ccvTicketContractID,
			CCVTicketTemplateID: ccvTicketTemplateID,
		},
	}, nil
}

// CommitteeVerifierVerifyMessageInput contains input for verifying a message
type CommitteeVerifierVerifyMessageInput struct {
	CommitteeVerifierContractID string
	CommitteeVerifierTemplateID string
	CcvRegistryCid              string
	Message                     ccvs.MessageV1
	MessageId                   string
	VerifierResults             string // BytesHex containing version tag and signatures
	Receiver                    string // Party ID
}

// CommitteeVerifierVerifyMessageOutput contains the transaction ID and CCVVerifyTicket contract ID
type CommitteeVerifierVerifyMessageOutput struct {
	TransactionID             string
	CCVVerifyTicketContractID string
	CCVVerifyTicketTemplateID string
}

var committeeVerifierVerifyMessageHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input CommitteeVerifierVerifyMessageInput) (output CantonOpResult[CommitteeVerifierVerifyMessageOutput], err error) {
	ctx := b.GetContext()

	// Create verify command
	verifyArgs := ccvs.CommitteeVerifierVerifyMessage{
		CcvRegistryCid:  types.CONTRACT_ID(input.CcvRegistryCid),
		Message:         input.Message,
		MessageId:       types.TEXT(input.MessageId),
		VerifierResults: types.TEXT(input.VerifierResults),
		Receiver:        types.PARTY(input.Receiver),
		Caller:          types.PARTY(deps.Party),
	}

	committeeVerifier := ccvs.CommitteeVerifier{}
	exerciseCmd := committeeVerifier.CommitteeVerifierVerifyMessage(input.CommitteeVerifierContractID, verifyArgs)

	// List known packages to find the package ID for ccip-committeeverifier
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[CommitteeVerifierVerifyMessageOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipCommitteeVerifierPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-committeeverifier") {
			ccipCommitteeVerifierPkgID = p.PackageID
			break
		}
	}
	if ccipCommitteeVerifierPkgID == "" {
		return CantonOpResult[CommitteeVerifierVerifyMessageOutput]{}, fmt.Errorf("failed to find ccip-committeeverifier package")
	}

	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: commandID,
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", ccipCommitteeVerifierPkgID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[CommitteeVerifierVerifyMessageOutput]{}, fmt.Errorf("failed to verify message: %w", err)
	}

	// Extract CCVVerifyTicket contract ID from Created event
	ccvVerifyTicketContractID := ""
	ccvVerifyTicketTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == "CCIP.Tickets:CCVVerifyTicket" {
			ccvVerifyTicketContractID = ev.Created.ContractID
			ccvVerifyTicketTemplateID = ev.Created.TemplateID

			break
		}
	}

	if ccvVerifyTicketContractID == "" {
		return CantonOpResult[CommitteeVerifierVerifyMessageOutput]{}, fmt.Errorf("verify-message tx had no Created CCVVerifyTicket event")
	}

	return CantonOpResult[CommitteeVerifierVerifyMessageOutput]{
		UpdateID: submitResp.UpdateID,
		Output: CommitteeVerifierVerifyMessageOutput{
			TransactionID:             commandID,
			CCVVerifyTicketContractID: ccvVerifyTicketContractID,
			CCVVerifyTicketTemplateID: ccvVerifyTicketTemplateID,
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var CommitteeVerifierForwardToVerifierOp = cld_ops.NewOperation(
	"canton/ccip/committeeverifier/forward-to-verifier",
	semver.MustParse("0.1.0"),
	"Forwards a message to the CommitteeVerifier for verification and issues a CCVTicket",
	committeeVerifierForwardToVerifierHandler,
)

var CommitteeVerifierVerifyMessageOp = cld_ops.NewOperation(
	"canton/ccip/committeeverifier/verify-message",
	semver.MustParse("0.1.0"),
	"Verifies a message using CommitteeVerifier and issues a CCVVerifyTicket",
	committeeVerifierVerifyMessageHandler,
)
