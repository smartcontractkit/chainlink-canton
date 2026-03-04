package changesets

import (
	"fmt"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	deploymentadapters "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

type ConfigureChainForLanesConfig struct {
	Input sequences.ConfigureChainForLanesInput
}

var _ cldf.ChangeSetV2[CantonCSDeps[ConfigureChainForLanesConfig]] = ConfigureChainForLanes{}

type ConfigureChainForLanes struct{}

func (c ConfigureChainForLanes) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (c ConfigureChainForLanes) Apply(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chainFamily := deploymentadapters.NewCantonChainFamilyAdapter()
	ccipInput := toCCIPConfigureChainForLanesInput(config.Config.Input)
	out, err := operations.ExecuteSequence(e.OperationsBundle, chainFamily.ConfigureChainForLanes(), e.BlockChains, ccipInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployChainContracts sequence: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{},
	}, nil
}

func toCCIPConfigureChainForLanesInput(input sequences.ConfigureChainForLanesInput) ccipadapters.ConfigureChainForLanesInput {
	out := ccipadapters.ConfigureChainForLanesInput{
		ChainSelector:      input.ChainSelector,
		Router:             input.GlobalConfig.Hex(), // Canton-specific: adapter maps this back to GlobalConfig.
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
