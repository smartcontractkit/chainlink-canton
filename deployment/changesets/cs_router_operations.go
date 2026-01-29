package changesets

import (
	"context"
	"fmt"

	cantonclient "github.com/smartcontractkit/chainlink-canton/deployment/client"
	ccipops "github.com/smartcontractkit/chainlink-canton/deployment/ops/ccip"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type RouterOperationsConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	LedgerAPIURL  string `yaml:"ledgerApiUrl"`
	AdminAPIURL   string `yaml:"adminApiUrl"`

	JWTSecret         string `yaml:"jwtSecret"`         // Optional, defaults to "unsafe"
	DeployerParty     string `yaml:"deployerParty"`     // Optional, will allocate if not provided
	DeployerPartyHint string `yaml:"deployerPartyHint"` // Optional hint for party allocation

	// Contract IDs (required)
	PerPartyRouterContractID string `yaml:"perPartyRouterContractId"`
	PerPartyRouterTemplateID string `yaml:"perPartyRouterTemplateId"`

	// CCIPSend Operation (optional)
	CCIPSend *RouterCCIPSendInput `yaml:"ccipSend"`

	// Execute Operation (optional)
	Execute *RouterExecuteInput `yaml:"execute"`
}

// RouterCCIPSendInput represents the input structure for CCIPSend
type RouterCCIPSendInput struct {
	OnRampCid             string   `yaml:"onRampCid"`
	GlobalConfigCid       string   `yaml:"globalConfigCid"`
	TokenAdminRegistryCid string   `yaml:"tokenAdminRegistryCid"`
	DestChainSelector     string   `yaml:"destChainSelector"` // Chain selector as string
	Receiver              string   `yaml:"receiver"`
	Payload               string   `yaml:"payload"`
	ExecutionGasLimit     int64    `yaml:"executionGasLimit"`
	CcipReceiveGasLimit   int64    `yaml:"ccipReceiveGasLimit"`
	TokenSendTicket       *string  `yaml:"tokenSendTicket"` // Optional
	CcvTickets            []string `yaml:"ccvTickets"`      // Optional
}

// RouterExecuteInput represents the input structure for Execute
type RouterExecuteInput struct {
	OffRampCid             string   `yaml:"offRampCid"`
	GlobalConfigCid        string   `yaml:"globalConfigCid"`
	TokenAdminRegistryCid  string   `yaml:"tokenAdminRegistryCid"`
	EncodedMessage         string   `yaml:"encodedMessage"`
	CcvVerifyTickets       []string `yaml:"ccvVerifyTickets"`
	TokenPoolCCVTicket     *string  `yaml:"tokenPoolCCVTicket"`     // Optional
	ReceiverRequiredCCVIds []string `yaml:"receiverRequiredCCVIds"` // Optional
}

var _ cldf.ChangeSetV2[RouterOperationsConfig] = RouterOperations{}

// RouterOperations performs router operations (CCIPSend and/or Execute)
type RouterOperations struct{}

// Apply implements deployment.ChangeSetV2.
func (r RouterOperations) Apply(e cldf.Environment, config RouterOperationsConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()
	ab := cldf.NewMemoryAddressBook()
	var seqReports []cld_ops.Report[any, any]

	// Setup Canton client
	setupResult, err := cantonclient.Setup(ctx, cantonclient.Config{
		LedgerAPIURL:      config.LedgerAPIURL,
		AdminAPIURL:       config.AdminAPIURL,
		JWTSecret:         config.JWTSecret,
		DeployerParty:     config.DeployerParty,
		DeployerPartyHint: config.DeployerPartyHint,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to setup Canton client: %w", err)
	}
	defer setupResult.BindingClient.Close()

	// Create Canton operation dependencies
	deps := ccipops.CantonOpDeps{
		BindingClient: setupResult.BindingClient,
		Party:         setupResult.Party,
	}

	// --------------------------
	// CCIPSend Operation
	// --------------------------
	if config.CCIPSend != nil {
		result, err := cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.RouterCCIPSendOp, deps, ccipops.RouterCCIPSendInput{
			PerPartyRouterContractID: config.PerPartyRouterContractID,
			PerPartyRouterTemplateID: config.PerPartyRouterTemplateID,
			OnRampCid:                config.CCIPSend.OnRampCid,
			GlobalConfigCid:          config.CCIPSend.GlobalConfigCid,
			TokenAdminRegistryCid:    config.CCIPSend.TokenAdminRegistryCid,
			DestChainSelector:        config.CCIPSend.DestChainSelector,
			Receiver:                 config.CCIPSend.Receiver,
			Payload:                  config.CCIPSend.Payload,
			ExecutionGasLimit:        config.CCIPSend.ExecutionGasLimit,
			CcipReceiveGasLimit:      config.CCIPSend.CcipReceiveGasLimit,
			TokenSendTicket:          config.CCIPSend.TokenSendTicket,
			CcvTickets:               config.CCIPSend.CcvTickets,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute CCIPSend: %w", err)
		}
		// Convert typed report to any,any for seqReports
		seqReports = append(seqReports, cld_ops.Report[any, any]{
			Input:  any(result.Input),
			Output: any(result.Output),
		})
	}

	// --------------------------
	// Execute Operation
	// --------------------------
	if config.Execute != nil {
		result, err := cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.RouterExecuteOp, deps, ccipops.RouterExecuteInput{
			PerPartyRouterContractID: config.PerPartyRouterContractID,
			PerPartyRouterTemplateID: config.PerPartyRouterTemplateID,
			OffRampCid:               config.Execute.OffRampCid,
			GlobalConfigCid:          config.Execute.GlobalConfigCid,
			TokenAdminRegistryCid:    config.Execute.TokenAdminRegistryCid,
			EncodedMessage:           config.Execute.EncodedMessage,
			CcvVerifyTickets:         config.Execute.CcvVerifyTickets,
			TokenPoolCCVTicket:       config.Execute.TokenPoolCCVTicket,
			ReceiverRequiredCCVIds:   config.Execute.ReceiverRequiredCCVIds,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute router Execute: %w", err)
		}
		// Convert typed report to any,any for seqReports
		seqReports = append(seqReports, cld_ops.Report[any, any]{
			Input:  any(result.Input),
			Output: any(result.Output),
		})

		// Save TokenReceiveTicket contract ID to address book if present
		// if result.Output.Output.TokenReceiveTicketID != nil {
		// 	typeAndVersionTokenReceiveTicket := cldf.NewTypeAndVersion("CantonTokenReceiveTicket", "1.0.0")
		// 	err = ab.Save(config.ChainSelector, *result.Output.Output.TokenReceiveTicketID, typeAndVersionTokenReceiveTicket)
		// 	if err != nil {
		// 		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save TokenReceiveTicket contract ID: %w", err)
		// 	}
		// }
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (r RouterOperations) VerifyPreconditions(e cldf.Environment, config RouterOperationsConfig) error {
	if config.LedgerAPIURL == "" {
		return fmt.Errorf("ledgerApiUrl is required")
	}
	if config.AdminAPIURL == "" {
		return fmt.Errorf("adminApiUrl is required")
	}
	if config.PerPartyRouterContractID == "" {
		return fmt.Errorf("perPartyRouterContractId is required")
	}
	if config.PerPartyRouterTemplateID == "" {
		return fmt.Errorf("perPartyRouterTemplateId is required")
	}
	if config.CCIPSend == nil && config.Execute == nil {
		return fmt.Errorf("at least one of ccipSend or execute must be provided")
	}

	// Validate CCIPSend if provided
	if config.CCIPSend != nil {
		if config.CCIPSend.OnRampCid == "" {
			return fmt.Errorf("ccipSend.onRampCid is required")
		}
		if config.CCIPSend.GlobalConfigCid == "" {
			return fmt.Errorf("ccipSend.globalConfigCid is required")
		}
		if config.CCIPSend.TokenAdminRegistryCid == "" {
			return fmt.Errorf("ccipSend.tokenAdminRegistryCid is required")
		}
		if config.CCIPSend.DestChainSelector == "" {
			return fmt.Errorf("ccipSend.destChainSelector is required")
		}
		if config.CCIPSend.Receiver == "" {
			return fmt.Errorf("ccipSend.receiver is required")
		}
	}

	// Validate Execute if provided
	if config.Execute != nil {
		if config.Execute.OffRampCid == "" {
			return fmt.Errorf("execute.offRampCid is required")
		}
		if config.Execute.GlobalConfigCid == "" {
			return fmt.Errorf("execute.globalConfigCid is required")
		}
		if config.Execute.TokenAdminRegistryCid == "" {
			return fmt.Errorf("execute.tokenAdminRegistryCid is required")
		}
		if config.Execute.EncodedMessage == "" {
			return fmt.Errorf("execute.encodedMessage is required")
		}
	}

	return nil
}
