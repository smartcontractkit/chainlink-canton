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

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/perpartyrouter"
)

// ============================================================================
// Router Operations
// ============================================================================

// RouterCCIPSendInput contains input for sending a message via router
type RouterCCIPSendInput struct {
	PerPartyRouterContractID string
	PerPartyRouterTemplateID string
	OnRampCid                string
	GlobalConfigCid          string
	TokenAdminRegistryCid    string
	DestChainSelector        string // Chain selector as string
	Receiver                 string
	Payload                  string
	ExecutionGasLimit        int64
	CcipReceiveGasLimit      int64
	TokenSendTicket          *string  // Optional contract ID
	CcvTickets               []string // Optional CCV ticket contract IDs
}

// RouterCCIPSendOutput contains the transaction ID and result
type RouterCCIPSendOutput struct {
	RouterContractID   string
	CcipMessageSentCID string
	MessageId          string
}

var routerCCIPSendHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input RouterCCIPSendInput) (output CantonOpResult[RouterCCIPSendOutput], err error) {
	ctx := b.GetContext()

	// Parse dest chain selector
	destChainSelector, ok := new(big.Int).SetString(input.DestChainSelector, 10)
	if !ok {
		return CantonOpResult[RouterCCIPSendOutput]{}, fmt.Errorf("invalid destChainSelector: %s", input.DestChainSelector)
	}
	scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

	// Convert CCV tickets
	ccvTickets := make([]types.CONTRACT_ID, len(input.CcvTickets))
	for i, ticket := range input.CcvTickets {
		ccvTickets[i] = types.CONTRACT_ID(ticket)
	}

	// Convert optional token send ticket
	var tokenSendTicket *types.CONTRACT_ID
	if input.TokenSendTicket != nil && *input.TokenSendTicket != "" {
		cid := types.CONTRACT_ID(*input.TokenSendTicket)
		tokenSendTicket = &cid
	}

	// Create CCIPSend command
	sendArgs := perpartyrouter.CCIPSend{
		OnRampCid:             types.CONTRACT_ID(input.OnRampCid),
		GlobalConfigCid:       types.CONTRACT_ID(input.GlobalConfigCid),
		TokenAdminRegistryCid: types.CONTRACT_ID(input.TokenAdminRegistryCid),
		DestChainSelector:     types.NUMERIC(destChainSelectorMantissa),
		Receiver:              types.TEXT(input.Receiver),
		Payload:               types.TEXT(input.Payload),
		ExecutionGasLimit:     types.INT64(input.ExecutionGasLimit),
		CcipReceiveGasLimit:   types.INT64(input.CcipReceiveGasLimit),
		TokenSendTicket:       tokenSendTicket,
		CcvTickets:            ccvTickets,
	}

	perPartyRouter := perpartyrouter.PerPartyRouter{}
	exerciseCmd := perPartyRouter.CCIPSend(input.PerPartyRouterContractID, sendArgs)

	// List known packages to find the package ID for ccip-perpartyrouter
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[RouterCCIPSendOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipPerPartyRouterPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-perpartyrouter") {
			ccipPerPartyRouterPkgID = p.PackageID
			break
		}
	}
	if ccipPerPartyRouterPkgID == "" {
		return CantonOpResult[RouterCCIPSendOutput]{}, fmt.Errorf("failed to find ccip-perpartyrouter package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", ccipPerPartyRouterPkgID, "CCIP.PerPartyRouter", "PerPartyRouter"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[RouterCCIPSendOutput]{}, fmt.Errorf("failed to send message via router: %w", err)
	}

	// Extract CCIPSendResult from transaction events
	// The result is returned as a Created event for CCIPMessageSent
	ccipMessageSentCID := ""
	messageId := ""
	routerContractID := input.PerPartyRouterContractID

	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == "CCIP.PerPartyRouter:CCIPMessageSent" {
			ccipMessageSentCID = ev.Created.ContractID
			// Extract messageId from the created contract's arguments if available
			// Note: messageId might need to be extracted differently depending on the event structure
			break
		}
	}

	if ccipMessageSentCID == "" {
		return CantonOpResult[RouterCCIPSendOutput]{}, fmt.Errorf("ccip-send tx had no Created CCIPMessageSent event")
	}

	return CantonOpResult[RouterCCIPSendOutput]{
		UpdateID: submitResp.UpdateID,
		Output: RouterCCIPSendOutput{
			RouterContractID:   routerContractID,
			CcipMessageSentCID: ccipMessageSentCID,
			MessageId:          messageId, // May be empty if not extractable from events
		},
	}, nil
}

// RouterExecuteInput contains input for executing a message via router
type RouterExecuteInput struct {
	PerPartyRouterContractID string
	PerPartyRouterTemplateID string
	OffRampCid               string
	GlobalConfigCid          string
	TokenAdminRegistryCid    string
	EncodedMessage           string
	CcvVerifyTickets         []string // CCV verify ticket contract IDs
	TokenPoolCCVTicket       *string  // Optional token pool CCV ticket contract ID
	ReceiverRequiredCCVIds   []string // Required CCV IDs for receiver
}

// RouterExecuteOutput contains the transaction ID and result
type RouterExecuteOutput struct {
	RouterContractID     string
	TokenReceiveTicketID *string // Optional
	MessageId            string
	Message              perpartyrouter.MessageV1
	State                perpartyrouter.MessageExecutionState
}

var routerExecuteHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input RouterExecuteInput) (output CantonOpResult[RouterExecuteOutput], err error) {
	ctx := b.GetContext()

	// Convert CCV verify tickets
	ccvVerifyTickets := make([]types.CONTRACT_ID, len(input.CcvVerifyTickets))
	for i, ticket := range input.CcvVerifyTickets {
		ccvVerifyTickets[i] = types.CONTRACT_ID(ticket)
	}

	// Convert optional token pool CCV ticket
	var tokenPoolCCVTicket *types.CONTRACT_ID
	if input.TokenPoolCCVTicket != nil && *input.TokenPoolCCVTicket != "" {
		cid := types.CONTRACT_ID(*input.TokenPoolCCVTicket)
		tokenPoolCCVTicket = &cid
	}

	// Convert receiver required CCV IDs
	receiverRequiredCCVIds := make([]types.TEXT, len(input.ReceiverRequiredCCVIds))
	for i, ccvId := range input.ReceiverRequiredCCVIds {
		receiverRequiredCCVIds[i] = types.TEXT(ccvId)
	}

	// Create Execute command
	executeArgs := perpartyrouter.Execute{
		OffRampCid:             types.CONTRACT_ID(input.OffRampCid),
		GlobalConfigCid:        types.CONTRACT_ID(input.GlobalConfigCid),
		TokenAdminRegistryCid:  types.CONTRACT_ID(input.TokenAdminRegistryCid),
		EncodedMessage:         types.TEXT(input.EncodedMessage),
		CcvVerifyTickets:       ccvVerifyTickets,
		TokenPoolCCVTicket:     tokenPoolCCVTicket,
		ReceiverRequiredCCVIds: receiverRequiredCCVIds,
	}

	perPartyRouter := perpartyrouter.PerPartyRouter{}
	exerciseCmd := perPartyRouter.Execute(input.PerPartyRouterContractID, executeArgs)

	// List known packages to find the package ID for ccip-perpartyrouter
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[RouterExecuteOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipPerPartyRouterPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-perpartyrouter") {
			ccipPerPartyRouterPkgID = p.PackageID
			break
		}
	}
	if ccipPerPartyRouterPkgID == "" {
		return CantonOpResult[RouterExecuteOutput]{}, fmt.Errorf("failed to find ccip-perpartyrouter package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", ccipPerPartyRouterPkgID, "CCIP.PerPartyRouter", "PerPartyRouter"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[RouterExecuteOutput]{}, fmt.Errorf("failed to execute message via router: %w", err)
	}

	// Extract ExecuteResult from transaction events
	// The result may be in Created events or we may need to query the router state
	// For now, we'll extract what we can from events
	tokenReceiveTicketID := (*string)(nil)
	messageId := ""
	routerContractID := input.PerPartyRouterContractID

	// Note: ExecuteResult details might need to be extracted differently
	// depending on how the DAML contract returns the result
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		// Look for TokenReceiveTicket if present
		if strings.Contains(normalized, "TokenReceiveTicket") {
			tid := ev.Created.ContractID
			tokenReceiveTicketID = &tid
		}
	}

	return CantonOpResult[RouterExecuteOutput]{
		UpdateID: submitResp.UpdateID,
		Output: RouterExecuteOutput{
			RouterContractID:     routerContractID,
			TokenReceiveTicketID: tokenReceiveTicketID,
			MessageId:            messageId,                                     // May need to be extracted differently
			Message:              perpartyrouter.MessageV1{},                    // May need to be extracted from events
			State:                perpartyrouter.MessageExecutionStateUNTOUCHED, // Default state, may need to be extracted from events
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var RouterCCIPSendOp = cld_ops.NewOperation(
	"canton/ccip/router/ccip-send",
	semver.MustParse("0.1.0"),
	"Sends a CCIP message via PerPartyRouter",
	routerCCIPSendHandler,
)

var RouterExecuteOp = cld_ops.NewOperation(
	"canton/ccip/router/execute",
	semver.MustParse("0.1.0"),
	"Executes a CCIP message via PerPartyRouter",
	routerExecuteHandler,
)
