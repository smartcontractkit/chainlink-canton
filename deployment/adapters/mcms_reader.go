package adapters

import (
	"context"
	"fmt"

	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	ccipchangesets "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var _ ccipchangesets.MCMSReader = (*CantonMCMSReader)(nil)

type CantonMCMSReader struct{}

func (r *CantonMCMSReader) GetChainMetadata(
	e cldf.Environment,
	chainSelector uint64,
	input ccipmcms.Input,
) (mcms_types.ChainMetadata, error) {
	mcmsRef, err := r.GetMCMSRef(e, chainSelector, input)
	if err != nil {
		return mcms_types.ChainMetadata{}, err
	}

	chain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok || len(chain.Participants) == 0 {
		return mcms_types.ChainMetadata{}, fmt.Errorf("canton chain %d not found or has no participants", chainSelector)
	}

	participant := chain.Participants[0]
	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}

	inspector := cantonsdk.NewInspector(
		participant.LedgerServices.State,
		opcontract.LedgerQueryParties(participant),
		cantonmcms.RoleForAction(input.TimelockAction).SDKRole,
	)

	metadata, err := inspector.GetRootMetadata(ctx, mcmsRef.Address)
	if err != nil {
		return mcms_types.ChainMetadata{}, fmt.Errorf("failed to read Canton MCMS metadata for chain %d: %w", chainSelector, err)
	}

	return metadata, nil
}

func (r *CantonMCMSReader) GetTimelockRef(
	e cldf.Environment,
	chainSelector uint64,
	input ccipmcms.Input,
) (datastore.AddressRef, error) {
	timelockRef, err := findMCMSAddressRef(e.DataStore, chainSelector, datastore.ContractType(ccipdeploymentutils.RBACTimelock), input.Qualifier)
	if err == nil {
		return timelockRef, nil
	}

	// Canton uses the MCMS instance address as the timelock target for proposal execution.
	return r.GetMCMSRef(e, chainSelector, input)
}

func (r *CantonMCMSReader) GetMCMSRef(
	e cldf.Environment,
	chainSelector uint64,
	input ccipmcms.Input,
) (datastore.AddressRef, error) {
	return findMCMSAddressRef(e.DataStore, chainSelector, cantonmcms.RoleForAction(input.TimelockAction).ContractType, input.Qualifier)
}

func findMCMSAddressRef(
	ds datastore.DataStore,
	chainSelector uint64,
	contractType datastore.ContractType,
	qualifier string,
) (datastore.AddressRef, error) {
	if ds == nil {
		return datastore.AddressRef{}, fmt.Errorf("datastore is required to resolve %s on chain %d", contractType, chainSelector)
	}

	filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(contractType),
	}
	if qualifier != "" {
		filters = append(filters, datastore.AddressRefByQualifier(qualifier))
	}

	refs := ds.Addresses().Filter(filters...)
	switch len(refs) {
	case 1:
		return refs[0], nil
	case 0:
		if qualifier != "" {
			return datastore.AddressRef{}, fmt.Errorf("no %s ref found on chain %d with qualifier %q", contractType, chainSelector, qualifier)
		}

		return datastore.AddressRef{}, fmt.Errorf("no %s ref found on chain %d", contractType, chainSelector)
	default:
		if qualifier != "" {
			return datastore.AddressRef{}, fmt.Errorf("multiple %s refs found on chain %d with qualifier %q", contractType, chainSelector, qualifier)
		}

		return datastore.AddressRef{}, fmt.Errorf("multiple %s refs found on chain %d; specify MCMS qualifier", contractType, chainSelector)
	}
}
