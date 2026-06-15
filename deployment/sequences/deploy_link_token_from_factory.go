package sequences

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	linkregistry "github.com/smartcontractkit/chainlink-canton/deployment/operations/linkregistry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

const defaultCantonLinkTokenInstrumentID = "link-token"

// DeployLinkTokenFromFactoryParams deploys Canton LinkRegistry via the ccip-qualified CCIPFactory.
type DeployLinkTokenFromFactoryParams struct {
	CCIPOwnerParty      string `json:"ccipOwnerParty,omitempty" yaml:"ccipOwnerParty,omitempty"`
	FactoryAddressRef   datastore.AddressRef
	ProposalDriven      bool `json:"proposalDriven,omitempty" yaml:"proposalDriven,omitempty"`
	LinkTokenInstanceID types.TEXT `json:"linkTokenInstanceId,omitempty" yaml:"linkTokenInstanceId,omitempty"`
	InstrumentID        splice_api_token_holding_v1.InstrumentId `json:"instrumentId,omitempty" yaml:"instrumentId,omitempty"`
	TokenQualifier      string `json:"tokenQualifier,omitempty" yaml:"tokenQualifier,omitempty"`
}

// DeployLinkTokenFromFactory emits an MCMS batch for DeployLinkToken when ProposalDriven is true.
var DeployLinkTokenFromFactory = operations.NewSequence(
	"canton/ccip/deploy_link_token_from_factory",
	semver.MustParse("2.0.0"),
	"Deploys LinkRegistry on Canton through CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployLinkTokenFromFactoryParams) (sequences.OnChainOutput, error) {
		if input.CCIPOwnerParty == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("CCIPOwnerParty is required")
		}
		ccipOwnerParty := types.PARTY(input.CCIPOwnerParty)

		factoryRaw, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		instrumentID := input.InstrumentID
		if instrumentID.Admin == "" {
			instrumentID = splice_api_token_holding_v1.InstrumentId{
				Admin: ccipOwnerParty,
				Id:    types.TEXT(defaultCantonLinkTokenInstrumentID),
			}
		}
		if instrumentID.Admin != ccipOwnerParty {
			return sequences.OnChainOutput{}, fmt.Errorf("instrument admin must match ccipOwnerParty")
		}

		linkInstanceID, err := ensureInstanceID(input.LinkTokenInstanceID, "link-registry")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("link token instance ID: %w", err)
		}

		deployReport, err := operations.ExecuteOperation(b, factoryops.DeployLinkToken, deps, newChoiceInput(factoryRaw, factorybindings.DeployLinkToken{
			Contract: link.LinkRegistry{
				RegistryAdmin:        ccipOwnerParty,
				RegistryInstrumentId: instrumentID,
				InstanceId:           types.TEXT(linkInstanceID),
				RegistryMeta:         splice_api_token_metadata_v1.Metadata{},
			},
		}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("deploy LinkRegistry from factory: %w", err)
		}

		linkRawAddr := linkInstanceID.RawInstanceAddress(ccipOwnerParty)
		tokenEncodedAddress := contracts.EncodeInstrumentID(instrumentID).Hex()
		qualifier := stringsTrimOrDefault(input.TokenQualifier, "LINK")

		out := sequences.OnChainOutput{
			Addresses: []datastore.AddressRef{
				newAddressRef(deps.ChainSelector(), linkRawAddr, linkregistry.ContractType, linkregistry.Version, ""),
				{
					Address: tokenEncodedAddress,
					Labels: datastore.NewLabelSet(
						fmt.Sprintf("instrument-admin:%s", instrumentID.Admin),
						fmt.Sprintf("instrument-id:%s", instrumentID.Id),
						fmt.Sprintf("ccip-owner:%s", ccipOwnerParty),
					),
					Type:          datastore.ContractType("Token"),
					Version:       linkregistry.Version,
					Qualifier:     qualifier,
					ChainSelector: deps.ChainSelector(),
				},
			},
		}

		if input.ProposalDriven && !deployReport.Output.Executed() {
			batchOp, err := contract.NewBatchOperationFromExercises([]contract.ExerciseOutput{deployReport.Output})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for link token deploy: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				out.BatchOps = []mcms_types.BatchOperation{batchOp}
			}
		}

		return out, nil
	},
)

func stringsTrimOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
