package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

func convertPartySlice(in []string) []types.PARTY {
	out := make([]types.PARTY, len(in))
	for i := range in {
		out[i] = types.PARTY(in[i])
	}

	return out
}

type ConfigureCommitteeVerifierForLanesInput struct {
	ChainSelector uint64
	adapters.CommitteeVerifierConfig[contracts.InstanceAddress]
}

var ConfigureCommitteeVerifierForLanes = operations.NewSequence(
	"canton/ccip/configure_committee_verifier_for_lanes",
	semver.MustParse("1.7.0"),
	"Configures a Canton CommitteeVerifier contract for multiple remote chains",
	func(b operations.Bundle, deps dependencies.CantonDeps, input ConfigureCommitteeVerifierForLanesInput) (output sequences.OnChainOutput, err error) {
		remoteChainConfigArgs := make([]ccvs.RemoteChainConfigArgs, 0, len(input.RemoteChains))
		allowListArgs := make([]ccvs.AllowListConfigArgs, 0, len(input.RemoteChains))
		signatureConfigs := make([]ccvs.SignatureConfig, 0, len(input.RemoteChains))

		for remoteSelector, remoteConfig := range input.RemoteChains {
			// Remote Chain Config
			remoteChainConfigArgs = append(remoteChainConfigArgs, ccvs.RemoteChainConfigArgs{
				RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
				FeeUSDCents:         types.NUMERIC(strconv.FormatUint(uint64(remoteConfig.FeeUSDCents), 10)),
				GasForVerification:  types.INT64(remoteConfig.GasForVerification),
				PayloadSizeBytes:    types.INT64(remoteConfig.PayloadSizeBytes),
				AllowListEnabled:    types.BOOL(remoteConfig.AllowlistEnabled),
			})

			// Allow List
			allowListArgs = append(allowListArgs, ccvs.AllowListConfigArgs{
				DestChainSelector:         types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
				AllowListEnabled:          types.BOOL(remoteConfig.AllowlistEnabled),
				AddedAllowListedSenders:   convertPartySlice(remoteConfig.AddedAllowlistedSenders),
				RemovedAllowListedSenders: convertPartySlice(remoteConfig.RemovedAllowlistedSenders),
			})

			// Signers
			signerKeys := make([]types.TEXT, len(remoteConfig.SignatureConfig.Signers))
			for i, signer := range remoteConfig.SignatureConfig.Signers {
				// Decode and encode signer pubkeys to ensure they're in the correct format
				signerBytes, err := hex.DecodeString(strings.TrimPrefix(signer, "0x"))
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to decode signer key %d for remote chain %d: %w", i, remoteSelector, err)
				}
				signerKeys[i] = types.TEXT(hex.EncodeToString(signerBytes))
			}
			signatureConfigs = append(signatureConfigs, ccvs.SignatureConfig{
				SourceChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
				Threshold:           types.INT64(remoteConfig.SignatureConfig.Threshold),
				SignerKeys:          signerKeys,
			})
		}

		for _, address := range input.CommitteeVerifier {
			_, err := operations.ExecuteOperation(b, committee_verifier.ApplySignatureConfigs, deps, contract.ChoiceInput[ccvs.CommitteeVerifierApplySignatureConfigs]{
				ChainSelector:   deps.Chain.Selector,
				InstanceAddress: address,
				ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
				Args: ccvs.CommitteeVerifierApplySignatureConfigs{
					SourceChainSelectorsToRemove: nil, // This doesn't support removing chains
					SignatureConfigs:             signatureConfigs,
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply signature configs to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}

			_, err = operations.ExecuteOperation(b, committee_verifier.ApplyRemoteChainConfigUpdates, deps, contract.ChoiceInput[ccvs.ApplyRemoteChainConfigUpdates]{
				ChainSelector:   deps.Chain.Selector,
				InstanceAddress: address,
				ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
				Args: ccvs.ApplyRemoteChainConfigUpdates{
					RemoteChainConfigArgs: remoteChainConfigArgs,
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply remote chain configs to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}

			_, err = operations.ExecuteOperation(b, committee_verifier.ApplyAllowListUpdates, deps, contract.ChoiceInput[ccvs.ApplyAllowListUpdates]{
				ChainSelector:   deps.Chain.Selector,
				InstanceAddress: address,
				ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
				Args: ccvs.ApplyAllowListUpdates{
					AllowListConfigArgsItems: allowListArgs,
					Caller:                   types.PARTY(deps.Chain.Participants[deps.Participant].PartyID),
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply allow list updates to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)
