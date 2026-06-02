package contract

import (
	"context"
	"io"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// TestFindActiveContractByInstanceAddress_multiPartyVisibility verifies that when several
// parties each own a distinct contract, a single ACS query spanning ActAs + ReadAs parties
// can resolve the correct contract for each instance address.
func TestFindActiveContractByInstanceAddress_multiPartyVisibility(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	templateID := common.GlobalConfig{}.GetTemplateID()

	const (
		operator   = "operator"
		ccipOwner  = "ccip-owner"
		poolOwnerA = "pool-owner-a"
		poolOwnerB = "pool-owner-b"
	)

	type ownedContract struct {
		instanceID string
		signatory  string
		contractID string
	}

	ledger := []ownedContract{
		{instanceID: "globalconfig-main", signatory: ccipOwner, contractID: "cid-global-config"},
		{instanceID: "lockreleasetokenpool-usdc", signatory: poolOwnerA, contractID: "cid-pool-usdc"},
		{instanceID: "lockreleasetokenpool-link", signatory: poolOwnerB, contractID: "cid-pool-link"},
	}

	acsStream := make([]*apiv2.GetActiveContractsResponse, len(ledger))
	for i, c := range ledger {
		acsStream[i] = makeActiveContractACSResponse(c.instanceID, c.signatory, c.contractID)
	}

	queryParties := LedgerQueryParties(canton.Participant{
		PartyID:        operator,
		ReadAsPartyIDs: []string{ccipOwner, poolOwnerA, poolOwnerB},
	})

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil).
		Times(len(ledger))
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.MatchedBy(func(req *apiv2.GetActiveContractsRequest) bool {
		fbp := req.GetEventFormat().GetFiltersByParty()
		return len(fbp) == 4 &&
			fbp[operator] != nil &&
			fbp[ccipOwner] != nil &&
			fbp[poolOwnerA] != nil &&
			fbp[poolOwnerB] != nil
	}), mock.Anything).
		RunAndReturn(func(context.Context, *apiv2.GetActiveContractsRequest, ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse], error) {
			return newACSStreamMock(t, acsStream), nil
		}).
		Times(len(ledger))

	for _, want := range ledger {
		target := contracts.InstanceID(want.instanceID).
			RawInstanceAddress(types.PARTY(want.signatory)).
			InstanceAddress()

		got, err := FindActiveContractByInstanceAddress(ctx, stateClient, queryParties, templateID, target)
		require.NoError(t, err)
		require.Equal(t, want.contractID, got.GetCreatedEvent().GetContractId())
	}
}

func makeActiveContractACSResponse(instanceID, signatory, contractID string) *apiv2.GetActiveContractsResponse {
	return &apiv2.GetActiveContractsResponse{
		ContractEntry: &apiv2.GetActiveContractsResponse_ActiveContract{
			ActiveContract: &apiv2.ActiveContract{
				CreatedEvent: &apiv2.CreatedEvent{
					ContractId: contractID,
					CreateArguments: &apiv2.Record{
						Fields: []*apiv2.RecordField{
							{
								Label: "instanceId",
								Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}},
							},
						},
					},
					Signatories: []string{signatory},
				},
			},
		},
	}
}

func newACSStreamMock(t *testing.T, responses []*apiv2.GetActiveContractsResponse) *mocks.MockServerStreamingClient[apiv2.GetActiveContractsResponse] {
	t.Helper()

	stream := mocks.NewMockServerStreamingClient[apiv2.GetActiveContractsResponse](t)
	for _, resp := range responses {
		stream.EXPECT().Recv().Return(resp, nil).Once()
	}
	stream.EXPECT().Recv().Return(nil, io.EOF).Once()
	stream.EXPECT().CloseSend().Return(nil).Once()

	return stream
}
