package adapters

import (
	"context"
	"fmt"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
)

// buildCantonChainMetadata queries live MCMS state and builds proposal chain metadata
// with Canton additionalFields (pre/post op counts, multisigId, instanceId, etc.).
func buildCantonChainMetadata(
	ctx context.Context,
	stateClient apiv2.StateServiceClient,
	partyID string,
	mcmsAddrHex string,
	role cantonsdk.TimelockRole,
	txCount uint64,
	overridePreviousRoot bool,
) (mcms_types.ChainMetadata, error) {
	inspector := cantonsdk.NewInspector(stateClient, []string{partyID}, role)

	opCount, err := inspector.GetOpCount(ctx, mcmsAddrHex)
	if err != nil {
		return mcms_types.ChainMetadata{}, fmt.Errorf("get op count: %w", err)
	}

	mcmsContract, err := cantonsdk.GetMCMSContract(ctx, stateClient, []string{partyID}, mcmsAddrHex)
	if err != nil {
		return mcms_types.ChainMetadata{}, fmt.Errorf("get MCMS contract: %w", err)
	}

	multisigID := cantonMultisigID(string(mcmsContract.InstanceId), string(mcmsContract.Owner), role)

	return cantonsdk.NewChainMetadata(
		opCount,
		opCount+txCount,
		int64(mcmsContract.ChainId),
		multisigID,
		mcmsAddrHex,
		overridePreviousRoot,
		string(mcmsContract.InstanceId),
	)
}

func cantonMultisigID(instanceID, party string, role cantonsdk.TimelockRole) string {
	return fmt.Sprintf("%s@%s-%s", instanceID, party, strings.ToLower(role.String()))
}
