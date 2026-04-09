package adapters

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	ccvsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
)

var _ ccipadapters.DeployChainContractsAdapter = (*CantonDeployChainContractsAdapter)(nil)

const receiverWaitForFinalityConfig types.TEXT = "00000000"

type CantonDeployChainContractsAdapter struct{}

func (c *CantonDeployChainContractsAdapter) SetContractParamsFromImportedConfig() *cldf_ops.Sequence[ccipadapters.DeployChainConfigCreatorInput, ccipadapters.DeployContractParams, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton/set-contract-params-from-imported-config",
		semver.MustParse("2.0.0"),
		"Returns the user-provided contract params for Canton",
		func(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, input ccipadapters.DeployChainConfigCreatorInput) (ccipadapters.DeployContractParams, error) {
			return input.UserProvidedConfig, nil
		},
	)
}

func (c *CantonDeployChainContractsAdapter) DeployChainContracts() *cldf_ops.Sequence[ccipadapters.DeployChainContractsInput, seqcore.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton/deploy-chain-contracts",
		semver.MustParse("2.0.0"),
		"Deploys CCIP contracts on a Canton chain",
		func(bundle cldf_ops.Bundle, chains cldf_chain.BlockChains, input ccipadapters.DeployChainContractsInput) (seqcore.OnChainOutput, error) {
			return DeployCantonChainContracts(bundle.GetContext(), bundle, chains, input)
		},
	)
}

func DeployCantonChainContracts(ctx context.Context, bundle cldf_ops.Bundle, chains cldf_chain.BlockChains, input ccipadapters.DeployChainContractsInput) (seqcore.OnChainOutput, error) {
	chain, ok := chains.CantonChains()[input.ChainSelector]
	if !ok || len(chain.Participants) == 0 {
		return seqcore.OnChainOutput{}, fmt.Errorf("canton chain %d not found or has no participants", input.ChainSelector)
	}

	participant := chain.Participants[0]
	ownerParty := deployerPartyID(input.DeployerContract, participant)

	nativeInstrumentID, err := lookupNativeInstrumentID(ctx, participant)
	if err != nil {
		return seqcore.OnChainOutput{}, err
	}

	env := deployment.Environment{
		BlockChains:       chains,
		DataStore:         datastore.NewMemoryDataStore().Seal(),
		Logger:            bundle.Logger,
		OperationsBundle:  bundle,
		ExistingAddresses: nil,
	}

	cfg := cantonChangesets.CantonCSDeps[cantonChangesets.DeployChainContractsConfig]{
		ChainSelector: input.ChainSelector,
		Participant:   0,
		Config: cantonChangesets.DeployChainContractsConfig{
			Params: sequences.DeployChainContractsParams{
				CCIPOwnerParty:     ownerParty,
				CommitteeVerifiers: committeeVerifierParams(ownerParty, input.ContractParams.CommitteeVerifiers),
				GlobalConfig: sequences.GlobalConfigParams{
					Template: common.GlobalConfig{
						CcipOwner:     "",
						ChainSelector: types.NUMERIC(strconv.FormatUint(input.ChainSelector, 10)),
					},
				},
				NativeInstrumentId: nativeInstrumentID,
				FeeQuoterConfig: sequences.FeeQuoterParams{
					Template: feequoter.FeeQuoter{
						PriceUpdaters: []types.PARTY{types.PARTY(ownerParty)},
					},
					USDPerNative: big.NewInt(100_000_000),
				},
				RMNRemote: sequences.RMNRemoteParams{
					Template: rmn.RMNRemote{
						RmnOwner:       types.PARTY(ownerParty),
						CursedSubjects: nil,
					},
				},
				Executors: executorParams(ownerParty, input.ContractParams.Executors),
			},
		},
	}

	out, err := cantonChangesets.DeployChainContracts{}.Apply(env, cfg)
	if err != nil {
		return seqcore.OnChainOutput{}, fmt.Errorf("failed to deploy canton chain contracts for selector %d: %w", input.ChainSelector, err)
	}

	return datastoreToOnChainOutput(out.DataStore.Seal())
}

func deployerPartyID(deployerContract string, participant canton.Participant) string {
	if trimmed := strings.TrimPrefix(deployerContract, "canton:"); trimmed != "" && trimmed != deployerContract {
		return trimmed
	}
	if deployerContract != "" {
		return deployerContract
	}
	return participant.PartyID
}

func lookupNativeInstrumentID(ctx context.Context, participant canton.Participant) (splice_api_token_holding_v1.InstrumentId, error) {
	tokenSource := participant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}

	client, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(interceptor),
	)
	if err != nil {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("failed to create token metadata client: %w", err)
	}

	info, err := client.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("error getting registry info: %w", err)
	}
	if info.StatusCode() != http.StatusOK {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("unexpected status code from token metadata client: %d: %v", info.StatusCode(), info.Body)
	}

	return splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(info.JSON200.AdminId),
		Id:    types.TEXT("Amulet"),
	}, nil
}

func committeeVerifierParams(ownerParty string, verifiers []ccipadapters.CommitteeVerifierDeployParams) []sequences.CommitteeVerifierParams {
	params := make([]sequences.CommitteeVerifierParams, 0, len(verifiers))
	for _, verifier := range verifiers {
		storageLocations := make([]types.TEXT, len(verifier.StorageLocations))
		for i, location := range verifier.StorageLocations {
			storageLocations[i] = types.TEXT(location)
		}
		qualifier := verifier.Qualifier
		if qualifier == "" {
			qualifier = devenvcommon.DefaultCommitteeVerifierQualifier
		}
		params = append(params, sequences.CommitteeVerifierParams{
			Qualifier: qualifier,
			Template: ccvsbindings.CommitteeVerifier{
				Owner:                        types.PARTY(ownerParty),
				CcipOwner:                    types.PARTY(ownerParty),
				VersionTag:                   types.TEXT("e9a05a20"),
				MessageSentObservers:         nil,
				StorageLocations:             storageLocations,
				StorageLocationsAdmin:        types.PARTY(ownerParty),
				PendingStorageLocationsAdmin: types.PARTY(ownerParty),
				Deps:                         ccvsbindings.CommitteeVerifierDeps{},
			},
		})
	}
	return params
}

func executorParams(ownerParty string, executors []ccipadapters.ExecutorDeployParams) []sequences.ExecutorParams {
	if len(executors) == 0 {
		return []sequences.ExecutorParams{defaultExecutorParams(ownerParty)}
	}

	params := make([]sequences.ExecutorParams, 0, len(executors))
	for _, exec := range executors {
		qualifier := exec.Qualifier
		if qualifier == "" {
			qualifier = devenvcommon.DefaultExecutorQualifier
		}
		maxCCVs := exec.MaxCCVsPerMsg
		if maxCCVs == 0 {
			maxCCVs = 10
		}
		params = append(params, sequences.ExecutorParams{
			Qualifier: qualifier,
			Template: executorbindings.Executor{
				Owner:         types.PARTY(ownerParty),
				MaxCCVsPerMsg: types.INT64(maxCCVs),
				DynamicConfig: executorbindings.DynamicConfig{
					FeeAggregator:         nil,
					AllowedFinalityConfig: finalityConfig(exec.DynamicConfig.AllowedFinalityConfig),
					CcvAllowlistEnabled:   types.BOOL(exec.DynamicConfig.CcvAllowlistEnabled),
				},
				AllowedCCVs: nil,
			},
		})
	}
	return params
}

func defaultExecutorParams(ownerParty string) sequences.ExecutorParams {
	return sequences.ExecutorParams{
		Qualifier: devenvcommon.DefaultExecutorQualifier,
		Template: executorbindings.Executor{
			Owner:         types.PARTY(ownerParty),
			MaxCCVsPerMsg: types.INT64(10),
			DynamicConfig: executorbindings.DynamicConfig{
				FeeAggregator:         nil,
				AllowedFinalityConfig: receiverWaitForFinalityConfig,
				CcvAllowlistEnabled:   types.BOOL(false),
			},
			AllowedCCVs: nil,
		},
	}
}

func finalityConfig(cfg finality.Config) types.TEXT {
	if cfg.IsZero() || cfg.WaitForFinality {
		return receiverWaitForFinalityConfig
	}
	raw := cfg.Raw()
	return types.TEXT(hex.EncodeToString(raw[:]))
}

func datastoreToOnChainOutput(ds datastore.DataStore) (seqcore.OnChainOutput, error) {
	addresses, err := ds.Addresses().Fetch()
	if err != nil {
		return seqcore.OnChainOutput{}, fmt.Errorf("fetch deployed addresses: %w", err)
	}
	contractsMeta, err := ds.ContractMetadata().Fetch()
	if err != nil {
		return seqcore.OnChainOutput{}, fmt.Errorf("fetch deployed contract metadata: %w", err)
	}
	chainMeta, err := ds.ChainMetadata().Fetch()
	if err != nil {
		return seqcore.OnChainOutput{}, fmt.Errorf("fetch deployed chain metadata: %w", err)
	}

	output := seqcore.OnChainOutput{
		Addresses: addresses,
		Metadata: seqcore.Metadata{
			Contracts: contractsMeta,
		},
	}
	if len(chainMeta) > 1 {
		return seqcore.OnChainOutput{}, fmt.Errorf("expected at most one chain metadata record, got %d", len(chainMeta))
	}
	if len(chainMeta) == 1 {
		output.Metadata.Chain = &chainMeta[0]
	}

	envMeta, err := ds.EnvMetadata().Get()
	if err == nil {
		output.Metadata.Env = &envMeta
	} else if !errors.Is(err, datastore.ErrEnvMetadataNotSet) {
		return seqcore.OnChainOutput{}, fmt.Errorf("fetch deployed env metadata: %w", err)
	}

	return output, nil
}
