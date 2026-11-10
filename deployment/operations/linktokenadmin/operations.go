// Package linktokenadmin runs the submitter side of a LinkTokenAdmin authorization: after MCMS
// approves a mint or burn, the participant's party executes it with the contract ids MCMS cannot
// carry, or cleans up any expired authorization.
package linktokenadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	spliceMeta "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/daRegistry"
)

var Version = semver.MustParse("1.0.0")

// ExecuteInput runs as the participant's PartyID, whose PartyID or ReadAsPartyIDs must see the
// authorization, the LinkTokenAdmin and, for a burn, the input holdings.
type ExecuteInput struct {
	// Kind is KindMint or KindBurn.
	Kind Kind `json:"kind"`
	// Authorization is the authorization's MCMS target address, "authorizationId@ccipOwner".
	Authorization    string        `json:"authorization"`
	ParticipantIndex int           `json:"participantIndex,omitempty"`
	Registry         RegistryInput `json:"registry"`
	// InputHoldingCids are the holdings a burn consumes.
	InputHoldingCids []string `json:"inputHoldingCids,omitempty"`
}

type ExecuteOutput struct {
	UpdateID string `json:"updateId"`
	// MintedHoldingCids can be empty after a committed mint when reading the update back fails.
	MintedHoldingCids []string `json:"mintedHoldingCids,omitempty"`
}

type CleanupInput struct {
	Kind             Kind   `json:"kind"`
	Authorization    string `json:"authorization"`
	ParticipantIndex int    `json:"participantIndex,omitempty"`
}

type CleanupOutput struct {
	UpdateID string `json:"updateId"`
}

var Execute = operations.NewOperation(
	"canton/link_token_admin/execute_authorization",
	Version,
	"Executes an MCMS-approved LinkTokenAdmin mint or burn authorization",
	execute,
)

var CleanupExpired = operations.NewOperation(
	"canton/link_token_admin/cleanup_expired_authorization",
	Version,
	"Archives an expired LinkTokenAdmin authorization",
	cleanupExpired,
)

var exerciseMint = contract.NewExercise(contract.ExerciseParams[tokenadmin.ExecuteMint]{
	Name:        "canton/link_token_admin/execute_mint",
	Version:     Version,
	Description: "Exercises ExecuteMint on a LinkMintAuthorization",
	Template:    tokenadmin.LinkMintAuthorization{},
	Method:      tokenadmin.LinkMintAuthorization{}.ExecuteMint,
})

var exerciseBurn = contract.NewExercise(contract.ExerciseParams[tokenadmin.ExecuteBurn]{
	Name:        "canton/link_token_admin/execute_burn",
	Version:     Version,
	Description: "Exercises ExecuteBurn on a LinkBurnAuthorization",
	Template:    tokenadmin.LinkBurnAuthorization{},
	Method:      tokenadmin.LinkBurnAuthorization{}.ExecuteBurn,
})

// Validate checks the input without touching the ledger or the registry.
func (in ExecuteInput) Validate() error {
	_, _, err := in.validate()

	return err
}

// Validate checks the input without touching the ledger.
func (in CleanupInput) Validate() error {
	_, _, err := parseTarget(in.Kind, in.Authorization)

	return err
}

func parseTarget(kind Kind, authorization string) (kindSpec, contracts.RawInstanceAddress, error) {
	spec, err := specFor(kind)
	if err != nil {
		return kindSpec{}, "", err
	}
	addr, err := contracts.RawInstanceAddressFromString(authorization)
	if err != nil {
		return kindSpec{}, "", fmt.Errorf("authorization: %w", err)
	}

	return spec, addr, nil
}

func (in ExecuteInput) validate() (kindSpec, contracts.RawInstanceAddress, error) {
	spec, addr, err := parseTarget(in.Kind, in.Authorization)
	if err != nil {
		return kindSpec{}, "", err
	}
	switch {
	case in.Kind != KindMint && in.Kind != KindBurn:
		return kindSpec{}, "", fmt.Errorf("execute supports mint and burn, not %s", in.Kind)
	case in.Kind == KindMint && len(in.InputHoldingCids) > 0:
		return kindSpec{}, "", errors.New("a mint takes no input holdings")
	case in.Kind == KindBurn && len(in.InputHoldingCids) == 0:
		return kindSpec{}, "", errors.New("a burn requires InputHoldingCids")
	}
	if err := in.Registry.validate(); err != nil {
		return kindSpec{}, "", err
	}

	return spec, addr, nil
}

// checkCounterparty rejects an authorization whose recipient or holder would have to co-sign:
// Execute acts only as the submitter, so the counterparty must be ccipOwner or the submitter.
func checkCounterparty(kind Kind, auth authorization, submitter string) error {
	if auth.counterparty == auth.ccipOwner || auth.counterparty == submitter {
		return nil
	}

	return fmt.Errorf("%s authorization names %s, which is neither ccipOwner %s nor the submitter %s; it must execute from that party's participant",
		kind, auth.counterparty, auth.ccipOwner, submitter)
}

func execute(b operations.Bundle, chain canton.Chain, in ExecuteInput) (ExecuteOutput, error) {
	spec, addr, err := in.validate()
	if err != nil {
		return ExecuteOutput{}, err
	}
	p, err := contract.ParticipantAt(chain, in.ParticipantIndex)
	if err != nil {
		return ExecuteOutput{}, err
	}
	ctx := b.GetContext()

	auth, authDisclosed, err := findAuthorization(ctx, p, spec, addr)
	if err != nil {
		return ExecuteOutput{}, err
	}
	if err := checkCounterparty(in.Kind, auth, p.PartyID); err != nil {
		return ExecuteOutput{}, err
	}
	adminCid, adminDisclosed, err := findAdmin(ctx, p, auth)
	if err != nil {
		return ExecuteOutput{}, err
	}
	inputsDisclosed, err := discloseAll(ctx, p, in.InputHoldingCids)
	if err != nil {
		return ExecuteOutput{}, fmt.Errorf("input holding: %w", err)
	}
	factory, err := resolveFactory(ctx, p, in, auth)
	if err != nil {
		return ExecuteOutput{}, err
	}

	disclosed := append([]*apiv2.DisclosedContract{authDisclosed, adminDisclosed}, factory.disclosed...)
	disclosed = append(disclosed, inputsDisclosed...)
	disclosedContracts := contract.DisclosedContractsFromProto(dedupDisclosures(disclosed...))
	submitter := types.PARTY(p.PartyID)
	meta := spliceMeta.Metadata{Values: map[string]types.TEXT{}}

	if in.Kind == KindBurn {
		report, err := operations.ExecuteOperation(b, exerciseBurn, chain, contract.ChoiceInput[tokenadmin.ExecuteBurn]{
			ContractID:         auth.contractID,
			DisclosedContracts: disclosedContracts,
			ParticipantIndex:   in.ParticipantIndex,
			Args: tokenadmin.ExecuteBurn{
				Submitter:          submitter,
				TokenAdminCid:      types.CONTRACT_ID(adminCid),
				BurnMintFactoryCid: types.CONTRACT_ID(factory.factoryCid),
				InputHoldingCids:   contractIDs(in.InputHoldingCids),
				RegistryContext:    factory.context,
				Meta:               meta,
			},
		})
		if err != nil {
			return ExecuteOutput{}, err
		}

		return ExecuteOutput{UpdateID: report.Output.ExecInfo.UpdateID}, nil
	}

	report, err := operations.ExecuteOperation(b, exerciseMint, chain, contract.ChoiceInput[tokenadmin.ExecuteMint]{
		ContractID:         auth.contractID,
		DisclosedContracts: disclosedContracts,
		ParticipantIndex:   in.ParticipantIndex,
		Args: tokenadmin.ExecuteMint{
			Submitter:          submitter,
			TokenAdminCid:      types.CONTRACT_ID(adminCid),
			BurnMintFactoryCid: types.CONTRACT_ID(factory.factoryCid),
			RegistryContext:    factory.context,
			Meta:               meta,
		},
	})
	if err != nil {
		return ExecuteOutput{}, err
	}
	updateID := report.Output.ExecInfo.UpdateID
	minted, err := mintedHoldingCids(ctx, p, updateID, auth.contractID)
	if err != nil {
		// The mint has committed; reporting a failure would invite a retry against an archived authorization.
		b.Logger.Warnf("mint committed in update %s but its holding ids could not be read: %v", updateID, err)
	}

	return ExecuteOutput{UpdateID: updateID, MintedHoldingCids: minted}, nil
}

func resolveFactory(ctx context.Context, p canton.Participant, in ExecuteInput, auth authorization) (registryFactory, error) {
	if in.Registry.URL != "" {
		client, err := daRegistry.NewClientWithResponses(in.Registry.URL)
		if err != nil {
			return registryFactory{}, fmt.Errorf("registry client for %q: %w", in.Registry.URL, err)
		}

		return fetchBurnMintFactory(ctx, client, burnMintRequest(in.Kind, auth, in.InputHoldingCids))
	}

	if len(in.Registry.Disclosed) > 0 {
		disclosed := make([]*apiv2.DisclosedContract, 0, len(in.Registry.Disclosed))
		for _, dc := range in.Registry.Disclosed {
			d, err := registryDisclosure(dc)
			if err != nil {
				return registryFactory{}, err
			}
			disclosed = append(disclosed, d)
		}

		return registryFactory{factoryCid: in.Registry.FactoryCid, context: in.Registry.ChoiceContext, disclosed: disclosed}, nil
	}

	disclosed, err := discloseAll(ctx, p, append([]string{in.Registry.FactoryCid}, in.Registry.DiscloseContractIDs...))
	if err != nil {
		return registryFactory{}, fmt.Errorf("registry contract: %w", err)
	}

	return registryFactory{factoryCid: in.Registry.FactoryCid, context: in.Registry.ChoiceContext, disclosed: disclosed}, nil
}

func contractIDs(cids []string) []types.CONTRACT_ID {
	out := make([]types.CONTRACT_ID, len(cids))
	for i, cid := range cids {
		out[i] = types.CONTRACT_ID(cid)
	}

	return out
}

func cleanupExpired(b operations.Bundle, chain canton.Chain, in CleanupInput) (CleanupOutput, error) {
	spec, addr, err := parseTarget(in.Kind, in.Authorization)
	if err != nil {
		return CleanupOutput{}, err
	}
	p, err := contract.ParticipantAt(chain, in.ParticipantIndex)
	if err != nil {
		return CleanupOutput{}, err
	}

	ac, err := findByInstance(b.GetContext(), p, spec.template.GetTemplateID(), addr.InstanceID(), addr.Owner())
	if err != nil {
		return CleanupOutput{}, fmt.Errorf("find authorization %s: %w", addr, err)
	}
	report, err := operations.ExecuteOperation(b, cleanupOps[in.Kind], chain, contract.ChoiceInput[types.PARTY]{
		ContractID:         ac.GetCreatedEvent().GetContractId(),
		Args:               types.PARTY(p.PartyID),
		DisclosedContracts: contract.DisclosedContractsFromProto([]*apiv2.DisclosedContract{disclosureOf(ac.GetCreatedEvent(), ac.GetSynchronizerId())}),
		ParticipantIndex:   in.ParticipantIndex,
	})
	if err != nil {
		return CleanupOutput{}, err
	}

	return CleanupOutput{UpdateID: report.Output.ExecInfo.UpdateID}, nil
}
