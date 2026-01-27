package ledger

import (
	"context"
	"errors"
	"fmt"
	"io"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

type ActiveContract struct {
	CreatedEvent   *apiv2.CreatedEvent
	SynchronizerID string
}

// query all active contracts for a party using wildcard filter
func (c *Client) GetAllContractsForParty(ctx context.Context, party string) ([]*ActiveContract, error) {
	ctx, err := c.AuthContext(ctx)
	if err != nil {
		return nil, err
	}

	offset, err := c.GetCurrentOffset(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := c.stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
								WildcardFilter: &apiv2.WildcardFilter{
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts: %w", err)
	}

	var contracts []*ActiveContract
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contract: %w", err)
		}

		if ac, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			contracts = append(contracts, &ActiveContract{
				CreatedEvent:   ac.ActiveContract.GetCreatedEvent(),
				SynchronizerID: ac.ActiveContract.GetSynchronizerId(),
			})
		}
	}

	return contracts, nil
}
