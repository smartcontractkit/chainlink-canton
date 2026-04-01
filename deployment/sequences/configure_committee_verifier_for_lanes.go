package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/contracts"
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

type ConfigureCommitteeVerifierAsSourceInput struct {
	ChainSelector uint64
	lanes.CommitteeVerifierConfig[datastore.AddressRef]
}

type ConfigureCommitteeVerifierAsDestInput struct {
	ChainSelector uint64
	lanes.CommitteeVerifierConfig[datastore.AddressRef]
}

var ConfigureCommitteeVerifierAsSource = operations.NewSequence(
	"canton/ccip/committee_verifier/configure_as_source",
	semver.MustParse("2.0.0"),
	"Configures the outbound CommitteeVerifier settings for remote chains",
	func(b operations.Bundle, deps chain.BlockChains, input ConfigureCommitteeVerifierAsSourceInput) (output sequences.OnChainOutput, err error) {
		chain, ok := deps.CantonChains()[input.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		remoteChainConfigArgs := make([]ccvs.RemoteChainConfigArgs, 0, len(input.RemoteChains))
		allowListArgs := make([]ccvs.AllowListConfigArgs, 0, len(input.RemoteChains))
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
		}

		for _, addressRef := range input.CommitteeVerifier {
			address := contracts.HexToInstanceAddress(addressRef.Address)

			_, err = operations.ExecuteOperation(b, committee_verifier.ApplyRemoteChainConfigUpdates, chain, contract.ChoiceInput[ccvs.ApplyRemoteChainConfigUpdates]{
				InstanceAddress: address,
				Args: ccvs.ApplyRemoteChainConfigUpdates{
					RemoteChainConfigArgs: remoteChainConfigArgs,
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply remote chain configs to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}

			_, err = operations.ExecuteOperation(b, committee_verifier.ApplyAllowListUpdates, chain, contract.ChoiceInput[ccvs.ApplyAllowListUpdates]{
				InstanceAddress: address,
				Args: ccvs.ApplyAllowListUpdates{
					AllowListConfigArgsItems: allowListArgs,
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply allow list updates to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)

var ConfigureCommitteeVerifierAsDest = operations.NewSequence(
	"canton/ccip/committee_verifier/configure_as_dest",
	semver.MustParse("2.0.0"),
	"Configures inbound CommitteeVerifier settings for remote chains",
	func(b operations.Bundle, deps chain.BlockChains, input ConfigureCommitteeVerifierAsDestInput) (output sequences.OnChainOutput, err error) {
		chain, ok := deps.CantonChains()[input.ChainSelector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		signatureConfigs := make([]ccvs.SignatureConfig, 0, len(input.RemoteChains))
		for remoteSelector, remoteConfig := range input.RemoteChains {
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

		for _, addressRef := range input.CommitteeVerifier {
			address := contracts.HexToInstanceAddress(addressRef.Address)

			_, err := operations.ExecuteOperation(b, committee_verifier.ApplySignatureConfigs, chain, contract.ChoiceInput[ccvs.ApplySignatureConfigs]{
				InstanceAddress: address,
				Args: ccvs.ApplySignatureConfigs{
					SourceChainSelectorsToRemove: nil, // This doesn't support removing chains
					SignatureConfigs:             signatureConfigs,
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply signature configs to CommitteeVerifier at address %s: %w", address.Hex(), err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)
