package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
)

// ConfigureChainForLanesAdapter exposes the CCIP tooling API input shape while
// delegating to Canton's native ConfigureChainForLanes sequence.
var ConfigureChainForLanesAdapter = operations.NewSequence(
	"canton/ccip/configure_chain_for_lanes_adapter",
	semver.MustParse("1.7.0"),
	"Adapter sequence for Canton ConfigureChainForLanes",
	func(
		b operations.Bundle,
		bcs cldfchain.BlockChains,
		input ccipadapters.ConfigureChainForLanesInput,
	) (seqcore.OnChainOutput, error) {
		chain, ok := bcs.CantonChains()[input.ChainSelector]
		if !ok {
			return seqcore.OnChainOutput{}, fmt.Errorf("canton chain %d not found in blockchains", input.ChainSelector)
		}

		localInput := ConfigureChainForLanesInput{
			ChainSelector:      input.ChainSelector,
			GlobalConfig:       contracts.HexToInstanceAddress(input.Router),
			FeeQuoter:          contracts.HexToInstanceAddress(input.FeeQuoter),
			OnRamp:             contracts.HexToInstanceAddress(input.OnRamp),
			OffRamp:            contracts.HexToInstanceAddress(input.OffRamp),
			CommitteeVerifiers: make([]ccipadapters.CommitteeVerifierConfig[contracts.InstanceAddress], 0, len(input.CommitteeVerifiers)),
			RemoteChains:       make(map[uint64]ccipadapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress], len(input.RemoteChains)),
		}

		for _, committee := range input.CommitteeVerifiers {
			cv := ccipadapters.CommitteeVerifierConfig[contracts.InstanceAddress]{
				CommitteeVerifier: make([]contracts.InstanceAddress, 0, len(committee.CommitteeVerifier)),
				RemoteChains:      committee.RemoteChains,
			}
			for _, ref := range committee.CommitteeVerifier {
				cv.CommitteeVerifier = append(cv.CommitteeVerifier, contracts.HexToInstanceAddress(ref.Address))
			}
			localInput.CommitteeVerifiers = append(localInput.CommitteeVerifiers, cv)
		}

		for selector, rc := range input.RemoteChains {
			localRC := ccipadapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
				AllowTrafficFrom:         rc.AllowTrafficFrom,
				OnRamps:                  rc.OnRamps,
				OffRamp:                  rc.OffRamp,
				DefaultInboundCCVs:       make([]contracts.RawInstanceAddress, 0, len(rc.DefaultInboundCCVs)),
				LaneMandatedInboundCCVs:  make([]contracts.RawInstanceAddress, 0, len(rc.LaneMandatedInboundCCVs)),
				DefaultOutboundCCVs:      make([]contracts.RawInstanceAddress, 0, len(rc.DefaultOutboundCCVs)),
				LaneMandatedOutboundCCVs: make([]contracts.RawInstanceAddress, 0, len(rc.LaneMandatedOutboundCCVs)),
				DefaultExecutor:          contracts.RawInstanceAddress(rc.DefaultExecutor),
				FeeQuoterDestChainConfig: rc.FeeQuoterDestChainConfig,
				ExecutorDestChainConfig:  rc.ExecutorDestChainConfig,
				AddressBytesLength:       rc.AddressBytesLength,
				BaseExecutionGasCost:     rc.BaseExecutionGasCost,
			}
			for _, v := range rc.DefaultInboundCCVs {
				localRC.DefaultInboundCCVs = append(localRC.DefaultInboundCCVs, contracts.RawInstanceAddress(v))
			}
			for _, v := range rc.LaneMandatedInboundCCVs {
				localRC.LaneMandatedInboundCCVs = append(localRC.LaneMandatedInboundCCVs, contracts.RawInstanceAddress(v))
			}
			for _, v := range rc.DefaultOutboundCCVs {
				localRC.DefaultOutboundCCVs = append(localRC.DefaultOutboundCCVs, contracts.RawInstanceAddress(v))
			}
			for _, v := range rc.LaneMandatedOutboundCCVs {
				localRC.LaneMandatedOutboundCCVs = append(localRC.LaneMandatedOutboundCCVs, contracts.RawInstanceAddress(v))
			}
			localInput.RemoteChains[selector] = localRC
		}

		out, err := operations.ExecuteSequence(
			b,
			ConfigureChainForLanes,
			dependencies.CantonDeps{Chain: chain, Participant: 0},
			localInput,
		)
		if err != nil {
			return seqcore.OnChainOutput{}, err
		}

		return out.Output, nil
	},
)
