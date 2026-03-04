package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

		localInput, err := toCantonConfigureChainForLanesInput(input)
		if err != nil {
			return seqcore.OnChainOutput{}, err
		}

		out, err := operations.ExecuteSequence(
			b,
			ConfigureChainForLanes,
			dependencies.CantonDeps{Chain: chain},
			localInput,
		)
		if err != nil {
			return seqcore.OnChainOutput{}, err
		}

		return out.Output, nil
	},
)

func toCantonConfigureChainForLanesInput(input ccipadapters.ConfigureChainForLanesInput) (ConfigureChainForLanesInput, error) {
	// Tooling API's "Router" field maps to Canton's GlobalConfig address.
	if input.Router == "" {
		return ConfigureChainForLanesInput{}, fmt.Errorf("missing required Router field (mapped to Canton GlobalConfig)")
	}
	if input.OnRamp == "" {
		return ConfigureChainForLanesInput{}, fmt.Errorf("missing required OnRamp field")
	}
	if input.FeeQuoter == "" {
		return ConfigureChainForLanesInput{}, fmt.Errorf("missing required FeeQuoter field")
	}
	if input.OffRamp == "" {
		return ConfigureChainForLanesInput{}, fmt.Errorf("missing required OffRamp field")
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

	return localInput, nil
}

// ToCCIPConfigureChainForLanesInput maps Canton's native sequence input
// to the shared Tooling API input shape.
func ToCCIPConfigureChainForLanesInput(input ConfigureChainForLanesInput) ccipadapters.ConfigureChainForLanesInput {
	out := ccipadapters.ConfigureChainForLanesInput{
		ChainSelector:      input.ChainSelector,
		Router:             input.GlobalConfig.Hex(), // Tooling Router maps to Canton GlobalConfig.
		OnRamp:             input.OnRamp.Hex(),
		FeeQuoter:          input.FeeQuoter.Hex(),
		OffRamp:            input.OffRamp.Hex(),
		CommitteeVerifiers: make([]ccipadapters.CommitteeVerifierConfig[datastore.AddressRef], 0, len(input.CommitteeVerifiers)),
		RemoteChains:       make(map[uint64]ccipadapters.RemoteChainConfig[[]byte, string], len(input.RemoteChains)),
	}

	for _, committee := range input.CommitteeVerifiers {
		cv := ccipadapters.CommitteeVerifierConfig[datastore.AddressRef]{
			CommitteeVerifier: make([]datastore.AddressRef, 0, len(committee.CommitteeVerifier)),
			RemoteChains:      committee.RemoteChains,
		}
		for _, address := range committee.CommitteeVerifier {
			cv.CommitteeVerifier = append(cv.CommitteeVerifier, datastore.AddressRef{
				Address: address.Hex(),
			})
		}
		out.CommitteeVerifiers = append(out.CommitteeVerifiers, cv)
	}

	for selector, rc := range input.RemoteChains {
		remote := ccipadapters.RemoteChainConfig[[]byte, string]{
			AllowTrafficFrom:         rc.AllowTrafficFrom,
			OnRamps:                  rc.OnRamps,
			OffRamp:                  rc.OffRamp,
			DefaultInboundCCVs:       make([]string, 0, len(rc.DefaultInboundCCVs)),
			LaneMandatedInboundCCVs:  make([]string, 0, len(rc.LaneMandatedInboundCCVs)),
			DefaultOutboundCCVs:      make([]string, 0, len(rc.DefaultOutboundCCVs)),
			LaneMandatedOutboundCCVs: make([]string, 0, len(rc.LaneMandatedOutboundCCVs)),
			DefaultExecutor:          rc.DefaultExecutor.String(),
			FeeQuoterDestChainConfig: rc.FeeQuoterDestChainConfig,
			ExecutorDestChainConfig:  rc.ExecutorDestChainConfig,
			AddressBytesLength:       rc.AddressBytesLength,
			BaseExecutionGasCost:     rc.BaseExecutionGasCost,
		}
		for _, v := range rc.DefaultInboundCCVs {
			remote.DefaultInboundCCVs = append(remote.DefaultInboundCCVs, v.String())
		}
		for _, v := range rc.LaneMandatedInboundCCVs {
			remote.LaneMandatedInboundCCVs = append(remote.LaneMandatedInboundCCVs, v.String())
		}
		for _, v := range rc.DefaultOutboundCCVs {
			remote.DefaultOutboundCCVs = append(remote.DefaultOutboundCCVs, v.String())
		}
		for _, v := range rc.LaneMandatedOutboundCCVs {
			remote.LaneMandatedOutboundCCVs = append(remote.LaneMandatedOutboundCCVs, v.String())
		}
		out.RemoteChains[selector] = remote
	}

	return out
}
