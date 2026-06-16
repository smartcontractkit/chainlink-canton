package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	linkregistry "github.com/smartcontractkit/chainlink-canton/deployment/operations/linkregistry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

const defaultLinkTokenInstrumentID = "link-token"

// DeployLinkTokenAndBurnMintPoolFromFactoryParams deploys Canton LinkRegistry and a
// BurnMintTokenPool for the same instrument via the ccip-qualified CCIPFactory.
type DeployLinkTokenAndBurnMintPoolFromFactoryParams struct {
	CCIPOwnerParty string
	FactoryAddressRef datastore.AddressRef
	ProposalDriven  bool

	LinkTokenInstanceID types.TEXT
	InstrumentID        splice_api_token_holding_v1.InstrumentId

	PoolQualifier  string
	PoolInstanceID types.TEXT
	Decimals       int64

	TokenAdminRegistryRef datastore.AddressRef
	TokenAdminRegistryRaw contracts.RawInstanceAddress
	RmnRemoteRaw          contracts.RawInstanceAddress
	FeeQuoterRaw          contracts.RawInstanceAddress

	RegisterWithTAR bool
}

// DeployLinkTokenAndBurnMintPoolFromFactory emits one MCMS batch with DeployLinkToken
// followed by DeployBurnMintTokenPool when ProposalDriven is true.
var DeployLinkTokenAndBurnMintPoolFromFactory = operations.NewSequence(
	"canton/ccip/deploy_link_token_and_burn_mint_pool_from_factory",
	semver.MustParse("2.0.0"),
	"Deploys LinkRegistry and BurnMintTokenPool on Canton through CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployLinkTokenAndBurnMintPoolFromFactoryParams) (sequences.OnChainOutput, error) {
		ccipOwnerParty, err := requireDeployParty(input.CCIPOwnerParty)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}
		if input.PoolQualifier == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("pool qualifier is required")
		}
		if input.Decimals <= 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("decimals must be positive")
		}

		factoryRaw, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		instrumentID := input.InstrumentID
		if instrumentID.Admin == "" {
			instrumentID = splice_api_token_holding_v1.InstrumentId{
				Admin: ccipOwnerParty,
				Id:    types.TEXT(defaultLinkTokenInstrumentID),
			}
		}
		if instrumentID.Admin != ccipOwnerParty {
			return sequences.OnChainOutput{}, fmt.Errorf("instrument admin must match ccipOwnerParty")
		}

		linkInstanceID, err := ensureInstanceID(input.LinkTokenInstanceID, "link-registry")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("link token instance ID: %w", err)
		}
		poolInstanceID, err := ensureInstanceID(input.PoolInstanceID, fmt.Sprintf("burnminttokenpool-%s", input.PoolQualifier))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("burn/mint pool instance ID: %w", err)
		}

		if input.TokenAdminRegistryRaw == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("TokenAdminRegistryRaw is required")
		}
		if input.RmnRemoteRaw == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("RmnRemoteRaw is required")
		}
		if input.FeeQuoterRaw == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("FeeQuoterRaw is required")
		}

		tokenAdminRegistryRef := input.TokenAdminRegistryRef
		tokenAdminRegistryRaw := input.TokenAdminRegistryRaw
		rmnRemoteRaw := input.RmnRemoteRaw
		feeQuoterRaw := input.FeeQuoterRaw

		var proposalOutputs []contract.ExerciseOutput

		linkDeployReport, err := operations.ExecuteOperation(b, factoryops.DeployLinkToken, deps, newChoiceInput(factoryRaw, factorybindings.DeployLinkToken{
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
		proposalOutputs = appendExerciseOutput(proposalOutputs, linkDeployReport.Output, input.ProposalDriven)
		linkRawAddr := linkInstanceID.RawInstanceAddress(ccipOwnerParty)

		poolDeployReport, err := operations.ExecuteOperation(b, factoryops.DeployBurnMintTokenPool, deps, newChoiceInput(factoryRaw, factorybindings.DeployBurnMintTokenPool{
			Contract: burnminttokenpool.BurnMintTokenPool{
				InstanceId:              types.TEXT(poolInstanceID),
				CcipOwner:               ccipOwnerParty,
				PoolOwner:               ccipOwnerParty,
				InstrumentId:            instrumentID,
				Decimals:                types.INT64(input.Decimals),
				RemoteChainConfigs:      map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{},
				TokenTransferFeeConfigs: map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{},
				PoolReceiveContext: splice_api_token_metadata_v1.ChoiceContext{
					Values: map[string]splice_api_token_metadata_v1.AnyValue{},
				},
				TransferTimeout: burnminttokenpool.TransferTimeout{
					RelativeHours: new(types.INT64(24)),
				},
				Deps: burnminttokenpool.BurnMintTokenPoolDeps{
					TokenAdminRegistry: tokenAdminRegistryRaw.Binding(),
					RmnRemote:          rmnRemoteRaw.Binding(),
					FeeQuoter:          feeQuoterRaw.Binding(),
				},
			},
		}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("deploy BurnMintTokenPool from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, poolDeployReport.Output, input.ProposalDriven)
		poolRawAddr := poolInstanceID.RawInstanceAddress(ccipOwnerParty)

		var batchOps []mcms_types.BatchOperation
		if input.RegisterWithTAR {
			regOut, err := operations.ExecuteSequence(b, RegisterTokenPool, deps, RegisterTokenPoolInput{
				TokenAdminRegistryInstanceAddress:    contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
				TokenAdminRegistryRawInstanceAddress: tokenAdminRegistryRaw,
				InstrumentId:                         instrumentID,
				PoolInstanceID:                       poolRawAddr.InstanceID(),
				CcipParty:                            string(ccipOwnerParty),
				PoolOwnerParty:                       string(ccipOwnerParty),
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("register BurnMintTokenPool with TAR: %w", err)
			}
			batchOps = append(batchOps, regOut.Output.BatchOps...)
		}

		tokenEncodedAddress := contracts.EncodeInstrumentID(instrumentID).Hex()
		qualifier := input.PoolQualifier
		addresses := []datastore.AddressRef{
			newAddressRef(deps.ChainSelector(), linkRawAddr, linkregistry.ContractType, linkregistry.Version, ""),
			newAddressRef(deps.ChainSelector(), poolRawAddr, burn_mint_token_pool.ContractType, burn_mint_token_pool.Version, qualifier),
			{
				Address: tokenEncodedAddress,
				Labels: datastore.NewLabelSet(
					fmt.Sprintf("instrument-admin:%s", instrumentID.Admin),
					fmt.Sprintf("instrument-id:%s", instrumentID.Id),
					fmt.Sprintf("ccip-owner:%s", ccipOwnerParty),
				),
				Type:          datastore.ContractType("Token"),
				Version:       burn_mint_token_pool.Version,
				Qualifier:     qualifier,
				ChainSelector: deps.ChainSelector(),
			},
			{
				Address:       poolRawAddr.InstanceAddress().String(),
				Labels:        datastore.NewLabelSet(poolRawAddr.String()),
				Type:          burnMintPoolType,
				Version:       burn_mint_token_pool.Version,
				Qualifier:     qualifier,
				ChainSelector: deps.ChainSelector(),
			},
		}

		if input.ProposalDriven && len(proposalOutputs) > 0 {
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for link token and pool deploy: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				batchOps = append([]mcms_types.BatchOperation{batchOp}, batchOps...)
			}
		}

		return sequences.OnChainOutput{
			Addresses: addresses,
			BatchOps:  batchOps,
		}, nil
	},
)

func requireDeployParty(party string) (types.PARTY, error) {
	if party == "" {
		return "", fmt.Errorf("ccipOwnerParty is required")
	}

	return types.PARTY(party), nil
}
