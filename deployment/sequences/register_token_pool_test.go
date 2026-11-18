package sequences

import (
	"context"
	"encoding/hex"
	"io"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// The legacy LINK pool is already registered; under MCMS the new pool must not get a second TokenConfig.
func TestRegisterTokenPool_proposalModeSkipsExistingTokenConfig(t *testing.T) {
	t.Parallel()

	ccipOwner := types.PARTY("ccip-owner::1220aa")
	instrumentID := splice_api_token_holding_v1.InstrumentId{Admin: ccipOwner, Id: "link-token"}
	existing := core.TokenConfig{
		InstanceId:         types.TEXT(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentID).Bytes())),
		RegistryInstanceId: "tar",
		RegistryOwner:      ccipOwner,
		InstrumentId:       instrumentID,
		Admin:              &ccipOwner,
		TokenPool:          &core.PoolRegistration2{PoolOwner: ccipOwner, PoolInstanceId: "legacy-pool"},
	}
	chain := cldfcanton.Chain{
		Participants: []cldfcanton.Participant{{
			PartyID:        "operator::1220dd",
			ReadAsPartyIDs: []string{string(ccipOwner)},
			LedgerServices: cldfcanton.LedgerServiceClients{State: ledgerStateWith(t, ccipOwner, existing)},
		}},
	}
	bundle := operations.NewBundle(t.Context, logger.Test(t), operations.NewMemoryReporter())
	tarRaw := contracts.InstanceID("tar").RawInstanceAddress(ccipOwner)

	out, err := operations.ExecuteSequence(bundle, RegisterTokenPool, chain, RegisterTokenPoolInput{
		TokenAdminRegistryInstanceAddress:    tarRaw.InstanceAddress(),
		TokenAdminRegistryRawInstanceAddress: tarRaw,
		InstrumentId:                         instrumentID,
		PoolInstanceID:                       "registry-pool",
		CcipParty:                            string(ccipOwner),
		PoolOwnerParty:                       string(ccipOwner),
		PoolAdminParty:                       string(ccipOwner),
	})

	require.NoError(t, err)
	require.Empty(t, out.Output.BatchOps)
}

// ledgerStateWith answers every active-contract query with the given contract, signed by signatory.
func ledgerStateWith(t *testing.T, signatory types.PARTY, template core.TokenConfig) *mocks.MockStateServiceClient {
	t.Helper()

	state := mocks.NewMockStateServiceClient(t)
	state.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil).Maybe()
	state.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(context.Context, *apiv2.GetActiveContractsRequest, ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse], error) {
			stream := mocks.NewMockServerStreamingClient[apiv2.GetActiveContractsResponse](t)
			stream.EXPECT().Recv().Return(&apiv2.GetActiveContractsResponse{
				ContractEntry: &apiv2.GetActiveContractsResponse_ActiveContract{ActiveContract: &apiv2.ActiveContract{
					CreatedEvent: &apiv2.CreatedEvent{
						ContractId:      "legacy-token-config",
						CreateArguments: bindings.MarshalTemplateToRecord(template),
						Signatories:     []string{string(signatory)},
					},
				}},
			}, nil).Once()
			stream.EXPECT().Recv().Return(nil, io.EOF).Once()
			stream.EXPECT().CloseSend().Return(nil).Maybe()

			return stream, nil
		}).Maybe()

	return state
}
