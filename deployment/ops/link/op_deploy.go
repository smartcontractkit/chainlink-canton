package linkops

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/coin"
	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// CantonOpDeps is an alias for the shared type in the client package
// This maintains backward compatibility while avoiding import cycles
type CantonOpDeps = compileClient.CantonOpDeps

// DeployLinkTokenOutput contains the deployed LINK token registry contract ID
type DeployLinkTokenOutput struct {
	RegistryInstanceID contracts.InstanceID
	RegistryContractID string
}

// CantonOpResult wraps the output for Canton operations
type CantonOpResult[T any] struct {
	UpdateID string
	Output   T
}

// normalizeTemplateKey normalizes template ID to match the pattern used in tests
func normalizeTemplateKey(tid string) string {
	tid = strings.TrimPrefix(tid, "#")
	parts := strings.Split(tid, ":")
	if len(parts) < 3 {
		return tid
	}

	return parts[len(parts)-2] + ":" + parts[len(parts)-1]
}

var handler = func(b cld_ops.Bundle, deps CantonOpDeps, input cld_ops.EmptyInput) (output CantonOpResult[DeployLinkTokenOutput], err error) {
	ctx := b.GetContext()

	instanceID, err := contracts.NewInstanceID("LINK", deps.Party)
	if err != nil {
		return CantonOpResult[DeployLinkTokenOutput]{}, fmt.Errorf("failed to create instance ID: %w", err)
	}

	// Create CoinRegistry contract for LINK token (following TestCoin pattern)
	reg := coin.CoinRegistry{
		Issuer: types.PARTY(deps.Party),
		InstrumentId: coin.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    "LINK",
		},
		InstanceId: types.TEXT(instanceID.String()),
		Meta: coin.Metadata{
			Values: types.TEXTMAP{},
		},
	}

	// Submit via binding client's CommandService
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands:  []*model.Command{{Command: reg.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployLinkTokenOutput]{}, fmt.Errorf("failed to submit CoinRegistry creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	registryContractID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		// TODO: add template name to bindings
		if contracts.ReplacePackageIdWithNameInTemplateID(reg.GetTemplateID(), string(contracts.Coin)) == contracts.ReplacePackageIdWithNameInTemplateID(event.Created.TemplateID, event.Created.PackageName) {
			registryContractID = event.Created.ContractID
			break
		}
	}

	if registryContractID == "" {
		return CantonOpResult[DeployLinkTokenOutput]{}, fmt.Errorf("failed to find CoinRegistry contract in transaction events")
	}

	fmt.Printf("Deployed LINK token registry contract   id=%s\n", registryContractID)

	return CantonOpResult[DeployLinkTokenOutput]{
		UpdateID: submitResp.UpdateID,
		Output: DeployLinkTokenOutput{
			RegistryContractID: registryContractID,
			RegistryInstanceID: instanceID,
		},
	}, nil
}

type MintLinkTokenInput struct {
	RegistryContractID string
	ReceiverParty      string
	Amount             string
}

type MintLinkTokenOutput struct {
	TokenHoldingContractID string
}

var handlerMint = func(b cld_ops.Bundle, deps CantonOpDeps, input MintLinkTokenInput) (output CantonOpResult[MintLinkTokenOutput], err error) {
	ctx := b.GetContext()

	// Parse amount
	amount, ok := new(big.Int).SetString(input.Amount, 10)
	if !ok {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("invalid amount: %s", input.Amount)
	}

	// Create MintPreapproval first (following TestCoin pattern)
	mintPreapproval := coin.MintPreapproval{
		Receiver: types.PARTY(input.ReceiverParty),
		Sender:   types.PARTY(deps.Party),
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{input.ReceiverParty},
			Commands:  []*model.Command{{Command: mintPreapproval.CreateCommand()}},
		},
	}

	preapprovalCommandID := uuid.Must(uuid.NewUUID()).String()
	cmds.Commands.CommandID = preapprovalCommandID
	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to create MintPreapproval: %w", err)
	}

	// Find the MintPreapproval contract ID
	mintPreapprovalCID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		// TODO: add template name to bindings
		if contracts.ReplacePackageIdWithNameInTemplateID(mintPreapproval.GetTemplateID(), string(contracts.Coin)) == contracts.ReplacePackageIdWithNameInTemplateID(event.Created.TemplateID, event.Created.PackageName) {
			mintPreapprovalCID = event.Created.ContractID
			break
		}
	}
	if mintPreapprovalCID == "" {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to find MintPreapproval contract")
	}

	fmt.Printf("MintPreapproval contract ID   id=%s\n", mintPreapprovalCID)
	fmt.Printf("Minting tokens to receiver party   party=%s\n", input.ReceiverParty)

	// Now mint tokens using BurnMintFactory_BurnMint choice
	mintArgs := coin.BurnMintFactoryBurnMint{
		ExpectedAdmin: types.PARTY(deps.Party),
		InstrumentId: coin.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    "LINK",
		},
		InputHoldingCids: []types.CONTRACT_ID{},
		Outputs: []coin.BurnMintOutput{
			{
				Owner:  types.PARTY(input.ReceiverParty),
				Amount: types.NUMERIC(amount),
				Context: coin.ChoiceContext{
					Values: types.TEXTMAP{
						"mint-preapproval": coin.AnyValue{
							AVContractId: func() *types.CONTRACT_ID {
								cid := types.CONTRACT_ID(mintPreapprovalCID)
								return &cid
							}(),
						},
					},
				},
			},
		},
		ExtraActors: []types.PARTY{},
		ExtraArgs: coin.ExtraArgs{
			Context: coin.ChoiceContext{Values: types.TEXTMAP{}},
			Meta:    coin.Metadata{Values: types.TEXTMAP{}},
		},
	}

	// Exercise the choice on the registry contract
	exerciseCmd := coin.CoinRegistry{}.BurnMintFactoryBurnMint(input.RegistryContractID, mintArgs)

	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var burnMintPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "splice-api-token-burn-mint") {
			burnMintPkgID = p.PackageID
			break
		}
	}

	cmds = &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				// TODO find a better way rather than this templateID override hack which exposes PackageID to the client
				TemplateID: fmt.Sprintf("%s:%s:%s", burnMintPkgID, "Splice.Api.Token.BurnMintV1", "BurnMintFactory"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	mintCommandID := uuid.Must(uuid.NewUUID()).String()
	cmds.Commands.CommandID = mintCommandID
	submitResp, err = deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to mint LINK tokens: %w", err)
	}

	// Find the token holding contract ID from the created events
	tokenHoldingCID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}
		// Look for Splice.Api.Token.HoldingV1:Holding contract
		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if strings.Contains(normalizedTemplateID, "Holding") {
			tokenHoldingCID = event.Created.ContractID
			break
		}
	}

	if tokenHoldingCID == "" {
		return CantonOpResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to find token holding contract in mint transaction")
	}

	fmt.Printf("Minted token to tokenHoldingCID   id=%s\n", tokenHoldingCID)

	return CantonOpResult[MintLinkTokenOutput]{
		UpdateID: submitResp.UpdateID,
		Output: MintLinkTokenOutput{
			TokenHoldingContractID: tokenHoldingCID,
		},
	}, nil
}

var DeployLINKOp = cld_ops.NewOperation(
	"canton/link/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the LINK Token CoinRegistry contract on Canton",
	handler,
)

var MintLINKPreApprovalOp = cld_ops.NewOperation(
	"canton/link/mint",
	semver.MustParse("0.1.0"),
	"Mint LINK tokens on Canton",
	handlerMint,
)
