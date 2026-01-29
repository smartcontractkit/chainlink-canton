package changesets

import (
	"context"
	"fmt"
	"math/big"

	cantonclient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
	ccipops "github.com/smartcontractkit/chainlink-canton-internal/deployment/ops/ccip"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/ccvs"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type CommitteeVerifierForwardConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	LedgerAPIURL  string `yaml:"ledgerApiUrl"`
	AdminAPIURL   string `yaml:"adminApiUrl"`

	JWTSecret         string `yaml:"jwtSecret"`         // Optional, defaults to "unsafe"
	DeployerParty     string `yaml:"deployerParty"`     // Optional, will allocate if not provided
	DeployerPartyHint string `yaml:"deployerPartyHint"` // Optional hint for party allocation

	// Contract IDs (required)
	CommitteeVerifierContractID string `yaml:"committeeVerifierContractId"`
	CommitteeVerifierTemplateID string `yaml:"committeeVerifierTemplateId"`
	CcvRegistryCid              string `yaml:"ccvRegistryCid"`

	// Message to forward
	Message *MessageV1Input `yaml:"message"`

	// Message metadata
	MessageId    string `yaml:"messageId"`
	VerifierArgs string `yaml:"verifierArgs"`

	// Fee token information
	FeeToken       *InstrumentIdInput `yaml:"feeToken"`
	FeeTokenAmount string             `yaml:"feeTokenAmount"` // Numeric value as string
}

// MessageV1Input represents the input structure for MessageV1
type MessageV1Input struct {
	SourceChainSelector string                `yaml:"sourceChainSelector"` // Chain selector as string
	DestChainSelector   string                `yaml:"destChainSelector"`   // Chain selector as string
	SequenceNumber      string                `yaml:"sequenceNumber"`      // Numeric value as string
	ExecutionGasLimit   int64                 `yaml:"executionGasLimit"`
	CcipReceiveGasLimit int64                 `yaml:"ccipReceiveGasLimit"`
	Finality            int64                 `yaml:"finality"`
	CcvAndExecutorHash  string                `yaml:"ccvAndExecutorHash"`
	OnRampAddress       string                `yaml:"onRampAddress"`
	OffRampAddress      string                `yaml:"offRampAddress"`
	Sender              string                `yaml:"sender"`
	Receiver            string                `yaml:"receiver"`
	DestBlob            string                `yaml:"destBlob"`
	MessageData         string                `yaml:"messageData"`
	TokenTransfer       *TokenTransferV1Input `yaml:"tokenTransfer"` // Optional
}

// TokenTransferV1Input represents the input structure for TokenTransferV1
type TokenTransferV1Input struct {
	MessageId        string             `yaml:"messageId"`
	SourceTokenData  []TokenAmountInput `yaml:"sourceTokenData"`
	DestTokenAmounts []TokenAmountInput `yaml:"destTokenAmounts"`
}

// TokenAmountInput represents the input structure for TokenAmount
type TokenAmountInput struct {
	InstrumentId InstrumentIdInput `yaml:"instrumentId"`
	Amount       string            `yaml:"amount"` // Numeric value as string
}

var _ cldf.ChangeSetV2[CommitteeVerifierForwardConfig] = CommitteeVerifierForward{}

// CommitteeVerifierForward forwards a message to the CommitteeVerifier for verification
type CommitteeVerifierForward struct{}

// Apply implements deployment.ChangeSetV2.
func (c CommitteeVerifierForward) Apply(e cldf.Environment, config CommitteeVerifierForwardConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0, 1)

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

	// Convert MessageV1Input to ccvs.MessageV1
	scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)

	// Parse source chain selector
	sourceChainSelector, ok := new(big.Int).SetString(config.Message.SourceChainSelector, 10)
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid sourceChainSelector: %s", config.Message.SourceChainSelector)
	}
	sourceChainSelectorMantissa := new(big.Int).Mul(sourceChainSelector, scale10)

	// Parse dest chain selector
	destChainSelector, ok := new(big.Int).SetString(config.Message.DestChainSelector, 10)
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid destChainSelector: %s", config.Message.DestChainSelector)
	}
	destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

	// Parse sequence number
	sequenceNumber, ok := new(big.Int).SetString(config.Message.SequenceNumber, 10)
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid sequenceNumber: %s", config.Message.SequenceNumber)
	}
	sequenceNumberMantissa := new(big.Int).Mul(sequenceNumber, scale10)

	// Convert token transfer if provided
	var tokenTransfer *ccvs.TokenTransferV1
	if config.Message.TokenTransfer != nil {
		// Parse message ID for token transfer
		// tokenTransferMessageId, ok := new(big.Int).SetString(config.Message.TokenTransfer.MessageId, 10)
		// if !ok {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("invalid tokenTransfer.messageId: %s", config.Message.TokenTransfer.MessageId)
		// }
		// tokenTransferMessageIdMantissa := new(big.Int).Mul(tokenTransferMessageId, scale10)

		// Convert source token data
		sourceTokenData := make([]ccvs.TokenAmount, len(config.Message.TokenTransfer.SourceTokenData))
		for i, st := range config.Message.TokenTransfer.SourceTokenData {
			amount, ok := new(big.Int).SetString(st.Amount, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid sourceTokenData[%d].amount: %s", i, st.Amount)
			}
			amountMantissa := new(big.Int).Mul(amount, scale10)

			sourceTokenData[i] = ccvs.TokenAmount{
				InstrumentId: ccvs.InstrumentId{
					Admin: types.PARTY(st.InstrumentId.Admin),
					Id:    types.TEXT(st.InstrumentId.Id),
				},
				Amount: types.NUMERIC(amountMantissa),
			}
		}

		// Convert dest token amounts
		destTokenAmounts := make([]ccvs.TokenAmount, len(config.Message.TokenTransfer.DestTokenAmounts))
		for i, dt := range config.Message.TokenTransfer.DestTokenAmounts {
			amount, ok := new(big.Int).SetString(dt.Amount, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid destTokenAmounts[%d].amount: %s", i, dt.Amount)
			}
			amountMantissa := new(big.Int).Mul(amount, scale10)

			destTokenAmounts[i] = ccvs.TokenAmount{
				InstrumentId: ccvs.InstrumentId{
					Admin: types.PARTY(dt.InstrumentId.Admin),
					Id:    types.TEXT(dt.InstrumentId.Id),
				},
				Amount: types.NUMERIC(amountMantissa),
			}
		}

		tokenTransfer = &ccvs.TokenTransferV1{
			//	MessageId:        types.NUMERIC(tokenTransferMessageIdMantissa),
			//	SourceTokenData:  sourceTokenData,
			//	DestTokenAmounts: destTokenAmounts,
		}
	}

	message := ccvs.MessageV1{
		SourceChainSelector: types.NUMERIC(sourceChainSelectorMantissa),
		DestChainSelector:   types.NUMERIC(destChainSelectorMantissa),
		SequenceNumber:      types.NUMERIC(sequenceNumberMantissa),
		ExecutionGasLimit:   types.INT64(config.Message.ExecutionGasLimit),
		CcipReceiveGasLimit: types.INT64(config.Message.CcipReceiveGasLimit),
		Finality:            types.INT64(config.Message.Finality),
		CcvAndExecutorHash:  types.TEXT(config.Message.CcvAndExecutorHash),
		OnRampAddress:       types.TEXT(config.Message.OnRampAddress),
		OffRampAddress:      types.TEXT(config.Message.OffRampAddress),
		Sender:              types.TEXT(config.Message.Sender),
		Receiver:            types.TEXT(config.Message.Receiver),
		DestBlob:            types.TEXT(config.Message.DestBlob),
		MessageData:         types.TEXT(config.Message.MessageData),
		TokenTransfer:       tokenTransfer,
	}

	// Convert fee token
	feeToken := ccvs.InstrumentId{
		Admin: types.PARTY(config.FeeToken.Admin),
		Id:    types.TEXT(config.FeeToken.Id),
	}

	// Execute the forward operation
	result, err := cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.CommitteeVerifierForwardToVerifierOp, deps, ccipops.CommitteeVerifierForwardToVerifierInput{
		CommitteeVerifierContractID: config.CommitteeVerifierContractID,
		CommitteeVerifierTemplateID: config.CommitteeVerifierTemplateID,
		CcvRegistryCid:              config.CcvRegistryCid,
		Message:                     message,
		MessageId:                   config.MessageId,
		FeeToken:                    feeToken,
		FeeTokenAmount:              config.FeeTokenAmount,
		VerifierArgs:                config.VerifierArgs,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to forward message to CommitteeVerifier: %w", err)
	}

	// Save CCVTicket contract ID to address book if needed
	// typeAndVersionCCVTicket := cldf.NewTypeAndVersion("CantonCCVTicket", "1.0.0")
	// err = ab.Save(config.ChainSelector, result.Output.Output.CCVTicketContractID, typeAndVersionCCVTicket)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCVTicket contract ID: %w", err)
	// }

	seqReports = append(seqReports, []cld_ops.Report[any, any]{result.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (c CommitteeVerifierForward) VerifyPreconditions(e cldf.Environment, config CommitteeVerifierForwardConfig) error {
	if config.LedgerAPIURL == "" {
		return fmt.Errorf("ledgerApiUrl is required")
	}
	if config.AdminAPIURL == "" {
		return fmt.Errorf("adminApiUrl is required")
	}
	if config.CommitteeVerifierContractID == "" {
		return fmt.Errorf("committeeVerifierContractId is required")
	}
	if config.CommitteeVerifierTemplateID == "" {
		return fmt.Errorf("committeeVerifierTemplateId is required")
	}
	if config.CcvRegistryCid == "" {
		return fmt.Errorf("ccvRegistryCid is required")
	}
	if config.Message == nil {
		return fmt.Errorf("message is required")
	}
	if config.MessageId == "" {
		return fmt.Errorf("messageId is required")
	}
	if config.FeeToken == nil {
		return fmt.Errorf("feeToken is required")
	}
	if config.FeeTokenAmount == "" {
		return fmt.Errorf("feeTokenAmount is required")
	}

	return nil
}
