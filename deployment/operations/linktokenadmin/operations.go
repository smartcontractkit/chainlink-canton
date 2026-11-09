// Package linktokenadmin runs the submitter side of a LinkTokenAdmin authorization: after MCMS
// approves a mint, burn, transfer or allocation, any party executes it with the contract ids MCMS
// cannot carry, or cleans it up once it has expired.
package linktokenadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	spliceMeta "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/daRegistry"
)

var Version = semver.MustParse("1.0.0")

// ExecuteInput runs on a participant whose PartyID or ReadAsPartyIDs can see the authorization,
// the LinkTokenAdmin and, for a burn, the holder's holdings.
type ExecuteInput struct {
	Kind Kind `json:"kind"`
	// Authorization is the authorization's MCMS target address, "authorizationId@ccipOwner".
	Authorization string `json:"authorization"`
	// Submitter defaults to the participant's party. A mint or burn also acts as the recipient or
	// holder when that is not ccipOwner, so the participant user needs CanActAs for it.
	Submitter        string        `json:"submitter,omitempty"`
	ParticipantIndex int           `json:"participantIndex,omitempty"`
	Registry         RegistryInput `json:"registry"`
	// InputHoldingCids pins the holdings a burn, transfer or allocation consumes. When empty they
	// are selected from the unlocked holdings visible to the participant.
	InputHoldingCids     []string         `json:"inputHoldingCids,omitempty"`
	TransferTimeoutHours *int64           `json:"transferTimeoutHours,omitempty"`
	Allocation           *AllocationInput `json:"allocation,omitempty"`
}

type AllocationInput struct {
	SettlementExecutor  string `json:"settlementExecutor"`
	SettlementRefID     string `json:"settlementRefId"`
	AllocateBeforeHours int64  `json:"allocateBeforeHours"`
	SettleBeforeHours   int64  `json:"settleBeforeHours"`
}

type ExecuteOutput struct {
	UpdateID string `json:"updateId"`
	// ResultContractIDs lists the contract ids in the choice result in order, for example the
	// minted holdings, or a transfer's pending instruction followed by the sender's change.
	ResultContractIDs []string `json:"resultContractIds"`
}

type CleanupInput struct {
	Kind             Kind   `json:"kind"`
	Authorization    string `json:"authorization"`
	Submitter        string `json:"submitter,omitempty"`
	ParticipantIndex int    `json:"participantIndex,omitempty"`
}

type CleanupOutput struct {
	UpdateID string `json:"updateId"`
}

var Execute = operations.NewOperation(
	"canton/link_token_admin/execute_authorization",
	Version,
	"Executes an MCMS-approved LinkTokenAdmin mint, burn, transfer or allocation authorization",
	func(b operations.Bundle, chain canton.Chain, input ExecuteInput) (ExecuteOutput, error) {
		return execute(b.GetContext(), chain, input)
	},
)

var CleanupExpired = operations.NewOperation(
	"canton/link_token_admin/cleanup_expired_authorization",
	Version,
	"Archives an expired LinkTokenAdmin authorization",
	func(b operations.Bundle, chain canton.Chain, input CleanupInput) (CleanupOutput, error) {
		return cleanupExpired(b.GetContext(), chain, input)
	},
)

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
	if err := in.Registry.validate(in.Kind); err != nil {
		return kindSpec{}, "", err
	}
	if in.Kind == KindMint && len(in.InputHoldingCids) > 0 {
		return kindSpec{}, "", errors.New("a mint takes no input holdings")
	}
	if in.Kind == KindAllocate && in.Allocation == nil {
		return kindSpec{}, "", errors.New("allocation settlement parameters are required for an allocate")
	}
	if in.Kind != KindAllocate && in.Allocation != nil {
		return kindSpec{}, "", errors.New("allocation settlement parameters apply only to an allocate")
	}
	if in.Kind != KindTransfer && in.TransferTimeoutHours != nil {
		return kindSpec{}, "", errors.New("a transfer timeout applies only to a transfer")
	}

	return spec, addr, nil
}

func execute(ctx context.Context, chain canton.Chain, in ExecuteInput) (ExecuteOutput, error) {
	spec, addr, err := in.validate()
	if err != nil {
		return ExecuteOutput{}, err
	}
	p, err := contract.ParticipantAt(chain, in.ParticipantIndex)
	if err != nil {
		return ExecuteOutput{}, err
	}
	submitter := in.Submitter
	if submitter == "" {
		submitter = p.PartyID
	}

	auth, authDisclosed, err := findAuthorization(ctx, p, spec, addr)
	if err != nil {
		return ExecuteOutput{}, err
	}
	adminCid, adminDisclosed, err := findAdmin(ctx, p, auth)
	if err != nil {
		return ExecuteOutput{}, err
	}
	inputs, inputsDisclosed, err := resolveInputs(ctx, p, in, spec, auth)
	if err != nil {
		return ExecuteOutput{}, err
	}
	factory, err := resolveFactory(ctx, p, in, auth, inputs)
	if err != nil {
		return ExecuteOutput{}, err
	}

	cmd, err := executeCommand(in, auth.contractID, executeArgs{
		submitter:  submitter,
		adminCid:   adminCid,
		factoryCid: factory.factoryCid,
		inputs:     inputs,
		context:    factory.context,
	})
	if err != nil {
		return ExecuteOutput{}, err
	}

	disclosed := append([]*apiv2.DisclosedContract{authDisclosed, adminDisclosed}, factory.disclosed...)
	disclosed = append(disclosed, inputsDisclosed...)
	updateID, result, err := submit(ctx, p, actAs(spec, auth, submitter), cmd, dedupDisclosures(disclosed...))
	if err != nil {
		return ExecuteOutput{}, err
	}

	return ExecuteOutput{UpdateID: updateID, ResultContractIDs: resultContractIDs(result)}, nil
}

func resolveInputs(ctx context.Context, p canton.Participant, in ExecuteInput, spec kindSpec, auth authorization) ([]string, []*apiv2.DisclosedContract, error) {
	if in.Kind == KindMint {
		return nil, nil, nil
	}

	if len(in.InputHoldingCids) > 0 {
		disclosed := make([]*apiv2.DisclosedContract, 0, len(in.InputHoldingCids))
		for _, cid := range in.InputHoldingCids {
			d, err := discloseByID(ctx, p, cid)
			if err != nil {
				return nil, nil, fmt.Errorf("input holding: %w", err)
			}
			disclosed = append(disclosed, d)
		}

		return in.InputHoldingCids, disclosed, nil
	}

	owner := holdingOwner(spec, auth)
	candidates, byID, err := listHoldings(ctx, p, owner, auth.instrumentID)
	if err != nil {
		return nil, nil, err
	}
	picked, err := selectHoldings(in.Kind, candidates, auth.amount)
	if err != nil {
		return nil, nil, fmt.Errorf("select holdings of %s: %w", owner, err)
	}
	disclosed := make([]*apiv2.DisclosedContract, 0, len(picked))
	for _, cid := range picked {
		disclosed = append(disclosed, byID[cid])
	}

	return picked, disclosed, nil
}

func resolveFactory(ctx context.Context, p canton.Participant, in ExecuteInput, auth authorization, inputs []string) (registryFactory, error) {
	if in.Registry.URL != "" {
		client, err := daRegistry.NewClientWithResponses(in.Registry.URL)
		if err != nil {
			return registryFactory{}, fmt.Errorf("registry client for %q: %w", in.Registry.URL, err)
		}

		return fetchBurnMintFactory(ctx, client, burnMintRequest(in.Kind, auth, inputs))
	}

	ids := append([]string{in.Registry.FactoryCid}, in.Registry.DiscloseContractIDs...)
	disclosed := make([]*apiv2.DisclosedContract, 0, len(ids))
	for _, cid := range ids {
		d, err := discloseByID(ctx, p, cid)
		if err != nil {
			return registryFactory{}, fmt.Errorf("registry contract: %w", err)
		}
		disclosed = append(disclosed, d)
	}

	return registryFactory{factoryCid: in.Registry.FactoryCid, context: in.Registry.ChoiceContext, disclosed: disclosed}, nil
}

type executeArgs struct {
	submitter  string
	adminCid   string
	factoryCid string
	inputs     []string
	context    spliceMeta.ChoiceContext
}

func contractIDs(cids []string) []types.CONTRACT_ID {
	out := make([]types.CONTRACT_ID, len(cids))
	for i, cid := range cids {
		out[i] = types.CONTRACT_ID(cid)
	}

	return out
}

func executeCommand(in ExecuteInput, authCid string, a executeArgs) (*model.ExerciseCommand, error) {
	submitter := types.PARTY(a.submitter)
	adminCid := types.CONTRACT_ID(a.adminCid)
	factoryCid := types.CONTRACT_ID(a.factoryCid)
	meta := spliceMeta.Metadata{Values: map[string]types.TEXT{}}

	switch in.Kind {
	case KindMint:
		return tokenadmin.LinkMintAuthorization{}.ExecuteMint(authCid, tokenadmin.ExecuteMint{
			Submitter:          submitter,
			TokenAdminCid:      adminCid,
			BurnMintFactoryCid: factoryCid,
			RegistryContext:    a.context,
			Meta:               meta,
		}), nil
	case KindBurn:
		return tokenadmin.LinkBurnAuthorization{}.ExecuteBurn(authCid, tokenadmin.ExecuteBurn{
			Submitter:          submitter,
			TokenAdminCid:      adminCid,
			BurnMintFactoryCid: factoryCid,
			InputHoldingCids:   contractIDs(a.inputs),
			RegistryContext:    a.context,
			Meta:               meta,
		}), nil
	case KindTransfer:
		var timeout *types.INT64
		if in.TransferTimeoutHours != nil {
			timeout = new(types.INT64(*in.TransferTimeoutHours))
		}

		return tokenadmin.LinkTransferAuthorization{}.ExecuteTransfer(authCid, tokenadmin.ExecuteTransfer{
			Submitter:            submitter,
			TokenAdminCid:        adminCid,
			TransferFactoryCid:   factoryCid,
			InputHoldingCids:     contractIDs(a.inputs),
			TransferTimeoutHours: timeout,
			RegistryContext:      a.context,
			Meta:                 meta,
		}), nil
	case KindAllocate:
		return tokenadmin.LinkAllocationAuthorization{}.ExecuteAllocate(authCid, tokenadmin.ExecuteAllocate{
			Submitter:            submitter,
			TokenAdminCid:        adminCid,
			AllocationFactoryCid: factoryCid,
			SettlementExecutor:   types.PARTY(in.Allocation.SettlementExecutor),
			SettlementRefId:      types.TEXT(in.Allocation.SettlementRefID),
			InputHoldingCids:     contractIDs(a.inputs),
			AllocateBeforeHours:  types.INT64(in.Allocation.AllocateBeforeHours),
			SettleBeforeHours:    types.INT64(in.Allocation.SettleBeforeHours),
			RegistryContext:      a.context,
			Meta:                 meta,
		}), nil
	}

	return nil, fmt.Errorf("unknown authorization kind %q", in.Kind)
}

func cleanupExpired(ctx context.Context, chain canton.Chain, in CleanupInput) (CleanupOutput, error) {
	spec, addr, err := parseTarget(in.Kind, in.Authorization)
	if err != nil {
		return CleanupOutput{}, err
	}
	p, err := contract.ParticipantAt(chain, in.ParticipantIndex)
	if err != nil {
		return CleanupOutput{}, err
	}
	submitter := in.Submitter
	if submitter == "" {
		submitter = p.PartyID
	}

	auth, authDisclosed, err := findAuthorization(ctx, p, spec, addr)
	if err != nil {
		return CleanupOutput{}, err
	}
	updateID, _, err := submit(ctx, p, []string{submitter}, spec.cleanup(auth.contractID, types.PARTY(submitter)), []*apiv2.DisclosedContract{authDisclosed})
	if err != nil {
		return CleanupOutput{}, err
	}

	return CleanupOutput{UpdateID: updateID}, nil
}
