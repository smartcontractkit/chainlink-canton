package changesets

import (
	"context"
	"errors"
	"io"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

func TestRequireNoLiveAuthorization(t *testing.T) {
	t.Parallel()

	const (
		selector  = uint64(1)
		ccipOwner = "owner::1220"
		authID    = "mint-1"
	)
	adminRaw := contracts.NewRawInstanceAddress("lta-1", types.PARTY(ccipOwner))
	approveMint := sequences.LinkTokenAdminAction{ApproveMint: &tokenadmin.ApproveMint{
		AuthorizationId: authID, Recipient: ccipOwner, Amount: "1", ValidForSeconds: 1,
	}}
	liveMint := func(contractID string) *apiv2.GetActiveContractsResponse {
		return &apiv2.GetActiveContractsResponse{ContractEntry: &apiv2.GetActiveContractsResponse_ActiveContract{
			ActiveContract: &apiv2.ActiveContract{CreatedEvent: &apiv2.CreatedEvent{
				ContractId: contractID,
				CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: authID}}},
				}},
				Signatories: []string{ccipOwner},
			}},
		}}
	}

	cases := []struct {
		name    string
		actions []sequences.LinkTokenAdminAction
		acs     []*apiv2.GetActiveContractsResponse
		acsErr  error
		lookups int
		wantErr string
	}{
		{name: "live authorization", actions: []sequences.LinkTokenAdminAction{approveMint}, acs: []*apiv2.GetActiveContractsResponse{liveMint("cid-1")}, lookups: 1, wantErr: "already live"},
		{name: "two live authorizations", actions: []sequences.LinkTokenAdminAction{approveMint}, acs: []*apiv2.GetActiveContractsResponse{liveMint("cid-1"), liveMint("cid-2")}, lookups: 1, wantErr: "already live"},
		{name: "no live authorization", actions: []sequences.LinkTokenAdminAction{approveMint}, lookups: 1},
		{name: "lookup fails", actions: []sequences.LinkTokenAdminAction{approveMint}, acsErr: errors.New("unavailable"), lookups: 1, wantErr: "unavailable"},
		{name: "withdraw only", actions: []sequences.LinkTokenAdminAction{{Withdraw: &sequences.LinkAuthorizationRef{Kind: linkadmingov.KindMint, AuthorizationID: authID}}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			state := mocks.NewMockStateServiceClient(t)
			if tc.lookups > 0 {
				state.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil).Times(tc.lookups)
				state.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(context.Context, *apiv2.GetActiveContractsRequest, ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse], error) {
						if tc.acsErr != nil {
							return nil, tc.acsErr
						}

						return newLinkACSStream(t, tc.acs), nil
					}).Times(tc.lookups)
			}
			e := cldf.Environment{
				GetContext: t.Context,
				BlockChains: chain.NewBlockChainsFromSlice([]chain.BlockChain{canton.Chain{
					Selector: selector,
					Participants: []canton.Participant{{
						PartyID:        "operator::1220",
						ReadAsPartyIDs: []string{ccipOwner},
						LedgerServices: canton.LedgerServiceClients{State: state},
					}},
				}}),
			}

			err := requireNoLiveAuthorization(e, CantonCSDeps[GovernLinkTokenAdminConfig]{ChainSelector: selector}, sequences.GovernLinkTokenAdminInput{
				AdminAddress: adminRaw,
				Actions:      tc.actions,
			})
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func newLinkACSStream(t *testing.T, responses []*apiv2.GetActiveContractsResponse) *mocks.MockServerStreamingClient[apiv2.GetActiveContractsResponse] {
	t.Helper()

	stream := mocks.NewMockServerStreamingClient[apiv2.GetActiveContractsResponse](t)
	for _, resp := range responses {
		stream.EXPECT().Recv().Return(resp, nil).Once()
	}
	stream.EXPECT().Recv().Return(nil, io.EOF).Maybe()
	stream.EXPECT().CloseSend().Return(nil).Once()

	return stream
}
