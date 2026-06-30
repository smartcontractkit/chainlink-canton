package adapters

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	committeeverifierop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	executorop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

// OutboundDestFeeOverrides sets Canton→remote dest fee fields for outbound hardening.
type OutboundDestFeeOverrides struct {
	MessageNetworkFeeUSDCents *uint16
	TokenNetworkFeeUSDCents   *uint16
	DefaultTokenFeeUSDCents   *uint16
}

// ResolveCantonOutboundDestLane builds source/dest lane definitions for Canton outbound
// dest fee updates (GlobalConfig + FeeQuoter only). Non-fee fields come from adapter defaults
// and deployed contract refs so existing lane wiring is preserved.
func ResolveCantonOutboundDestLane(
	ds datastore.DataStore,
	cantonSelector, remoteDestSelector uint64,
	overrides OutboundDestFeeOverrides,
	defaultOutboundCCVQualifiers []string,
	executorQualifier string,
) (source *lanes.ChainDefinition, dest *lanes.ChainDefinition, err error) {
	if remoteDestSelector == 0 {
		return nil, nil, fmt.Errorf("remoteDestSelector is required")
	}
	if executorQualifier == "" {
		executorQualifier = "default"
	}
	if len(defaultOutboundCCVQualifiers) == 0 {
		defaultOutboundCCVQualifiers = []string{"default"}
	}

	remoteFamily, err := chainsel.GetSelectorFamily(remoteDestSelector)
	if err != nil {
		return nil, nil, fmt.Errorf("remote chain family: %w", err)
	}
	remoteAdapter, ok := ccipadapters.GetChainFamilyRegistry().GetChainFamily(remoteFamily)
	if !ok {
		return nil, nil, fmt.Errorf("no chain family adapter for %q", remoteFamily)
	}

	remoteOnRamp, err := remoteAdapter.GetOnRampAddress(ds, remoteDestSelector)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve remote OnRamp: %w", err)
	}
	remoteOffRamp, err := remoteAdapter.GetOffRampAddress(ds, remoteDestSelector)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve remote OffRamp: %w", err)
	}

	remoteDefaults := remoteAdapter.GetDefaultRemoteChainConfig(cantonSelector, remoteDestSelector)
	messageNetworkFee := remoteDefaults.MessageNetworkFeeUSDCents
	if overrides.MessageNetworkFeeUSDCents != nil {
		messageNetworkFee = *overrides.MessageNetworkFeeUSDCents
	}
	tokenNetworkFee := remoteDefaults.TokenNetworkFeeUSDCents
	if overrides.TokenNetworkFeeUSDCents != nil {
		tokenNetworkFee = *overrides.TokenNetworkFeeUSDCents
	}

	fqOverrides := remoteAdapter.GetDefaultFeeQuoterDestChainConfig(
		cantonSelector,
		remoteDestSelector,
		remoteAdapter.GetChainFamilySelector(),
	)
	if overrides.DefaultTokenFeeUSDCents != nil {
		fqOverrides.DefaultTokenFeeUSDCents = overrides.DefaultTokenFeeUSDCents
	}
	fqOverrides.ChainFamilySelector = remoteAdapter.GetChainFamilySelector()

	defaultExecutor, err := findContractRef(
		ds,
		cantonSelector,
		datastore.ContractType(executorop.ContractType),
		executorop.Version,
		executorQualifier,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve default executor: %w", err)
	}

	defaultOutboundCCVs := make([]datastore.AddressRef, 0, len(defaultOutboundCCVQualifiers))
	for _, qualifier := range defaultOutboundCCVQualifiers {
		ref, err := findContractRef(
			ds,
			cantonSelector,
			datastore.ContractType(committeeverifierop.ContractType),
			committeeverifierop.Version,
			qualifier,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve default outbound CCV qualifier %q: %w", qualifier, err)
		}
		defaultOutboundCCVs = append(defaultOutboundCCVs, ref)
	}

	dest, err = BuildRemoteChainDefinition(remoteDestSelector, ccipadapters.RemoteChainConfig[[]byte, string]{
		OnRamps:                   [][]byte{remoteOnRamp},
		OffRamp:                   remoteOffRamp,
		BaseExecutionGasCost:      remoteDefaults.BaseExecutionGasCost,
		TokenReceiverAllowed:      &remoteDefaults.TokenReceiverAllowed,
		MessageNetworkFeeUSDCents: messageNetworkFee,
		TokenNetworkFeeUSDCents:   tokenNetworkFee,
		FeeQuoterDestChainConfig:  fqOverrides,
		AddressBytesLength:        remoteAdapter.GetAddressBytesLength(),
	})
	if err != nil {
		return nil, nil, err
	}

	source = &lanes.ChainDefinition{
		Selector:                 cantonSelector,
		DefaultExecutor:          defaultExecutor,
		DefaultOutboundCCVs:      defaultOutboundCCVs,
		LaneMandatedOutboundCCVs: nil,
	}

	return source, dest, nil
}
