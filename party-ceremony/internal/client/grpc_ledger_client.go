package client

import (
	"context"
	"errors"
	"fmt"
	"io"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	apiv2admin "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const defaultLedgerPort = 5002

// DialLedger opens a gRPC connection to a Canton Ledger API endpoint using
// the ledger host, port, and optional JWT bearer token from cfg.
func DialLedger(cfg ClientConfig) (*grpc.ClientConn, error) {
	host := cfg.LedgerHost
	if host == "" {
		host = defaultHost
	}
	port := cfg.LedgerPort
	if port == 0 {
		port = defaultLedgerPort
	}
	target := fmt.Sprintf("%s:%d", host, port)

	creds := transportCredentials(host, port)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	jwt := cfg.LedgerJWT
	if jwt == "" {
		jwt = cfg.AdminJWT // fall back to admin JWT if no separate ledger JWT
	}
	if jwt != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(jwtUnaryInterceptor(jwt)))
		opts = append(opts, grpc.WithStreamInterceptor(jwtStreamInterceptor(jwt)))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial ledger %s: %w", target, err)
	}

	return conn, nil
}

// GRPCLedgerClient implements [LedgerClient] using real Canton Ledger gRPC APIs.
type GRPCLedgerClient struct {
	interactive    interactive.InteractiveSubmissionServiceClient
	party          apiv2admin.PartyManagementServiceClient
	userManagement apiv2admin.UserManagementServiceClient
	state          apiv2.StateServiceClient
}

// NewGRPCLedgerClient creates a new [GRPCLedgerClient] from an established gRPC connection.
func NewGRPCLedgerClient(conn *grpc.ClientConn) *GRPCLedgerClient {
	return &GRPCLedgerClient{
		interactive:    interactive.NewInteractiveSubmissionServiceClient(conn),
		party:          apiv2admin.NewPartyManagementServiceClient(conn),
		userManagement: apiv2admin.NewUserManagementServiceClient(conn),
		state:          apiv2.NewStateServiceClient(conn),
	}
}

func (c *GRPCLedgerClient) PartyExists(ctx context.Context, partyID string) (bool, error) {
	resp, err := c.party.ListKnownParties(ctx, &apiv2admin.ListKnownPartiesRequest{
		FilterParty: partyID,
	})
	if err != nil {
		return false, fmt.Errorf("ListKnownParties: %w", err)
	}

	for _, pd := range resp.GetPartyDetails() {
		if pd.GetParty() == partyID {
			return true, nil
		}
	}

	return false, nil
}

func (c *GRPCLedgerClient) PrepareSubmission(
	ctx context.Context,
	commands []*apiv2.Command,
	actAs []string,
	readAs []string,
	synchronizerID string,
) (*interactive.PrepareSubmissionResponse, error) {
	resp, err := c.interactive.PrepareSubmission(ctx, &interactive.PrepareSubmissionRequest{
		CommandId:      uuid.New().String(),
		Commands:       commands,
		ActAs:          actAs,
		ReadAs:         readAs,
		SynchronizerId: LogicalSynchronizerID(synchronizerID),
	})
	if err != nil {
		return nil, fmt.Errorf("PrepareSubmission: %w", err)
	}

	return resp, nil
}

func (c *GRPCLedgerClient) GetActiveContractsByTemplate(
	ctx context.Context,
	partyID string,
	packageID string,
	moduleName string,
	entityName string,
) ([]*apiv2.CreatedEvent, error) {
	stream, err := c.state.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		EventFormat: &apiv2.EventFormat{
			// FiltersForAnyParty returns contracts visible to any party the
			// participant hosts that match the template. Using this instead of
			// FiltersByParty because decentralized-namespace parties may not
			// be listed in the participant's local party registry even though
			// they are hosted via the P2P topology mapping.
			FiltersForAnyParty: &apiv2.Filters{
				Cumulative: []*apiv2.CumulativeFilter{{
					IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
						TemplateFilter: &apiv2.TemplateFilter{
							TemplateId: &apiv2.Identifier{
								PackageId:  packageID,
								ModuleName: moduleName,
								EntityName: entityName,
							},
							IncludeCreatedEventBlob: true,
						},
					},
				}},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("GetActiveContracts: %w", err)
	}

	var contracts []*apiv2.CreatedEvent
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("GetActiveContracts stream: %w", err)
		}
		if ac := resp.GetActiveContract(); ac != nil {
			if ce := ac.GetCreatedEvent(); ce != nil {
				contracts = append(contracts, ce)
			}
		}
	}

	return contracts, nil
}

// GetActiveContractsByTemplateForParty queries the ACS using FiltersByParty,
// which authorizes against CanReadAs rights for the given party — sufficient
// for users granted per-party rights via [GRPCLedgerClient.GrantPartyRights].
// FiltersForAnyParty (used by [GRPCLedgerClient.GetActiveContractsByTemplate])
// requires participant/IDP admin claims, which per-party-granted users lack.
func (c *GRPCLedgerClient) GetActiveContractsByTemplateForParty(
	ctx context.Context,
	partyID string,
	packageID string,
	moduleName string,
	entityName string,
) ([]*apiv2.CreatedEvent, error) {
	// active_at_offset is required; per proto, zero returns the empty set.
	// Use the current ledger end so we evaluate "active now".
	endResp, err := c.state.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetLedgerEnd: %w", err)
	}
	activeAtOffset := endResp.GetOffset()

	stream, err := c.state.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: activeAtOffset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				partyID: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
							TemplateFilter: &apiv2.TemplateFilter{
								TemplateId: &apiv2.Identifier{
									PackageId:  packageID,
									ModuleName: moduleName,
									EntityName: entityName,
								},
								IncludeCreatedEventBlob: true,
							},
						},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("GetActiveContracts: %w", err)
	}

	var contracts []*apiv2.CreatedEvent
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("GetActiveContracts stream: %w", err)
		}
		if ac := resp.GetActiveContract(); ac != nil {
			if ce := ac.GetCreatedEvent(); ce != nil {
				contracts = append(contracts, ce)
			}
		}
	}

	return contracts, nil
}

// ExecuteSubmission calls InteractiveSubmissionService.ExecuteSubmissionAndWaitForTransaction
// to submit the prepared transaction and block until it is committed to the ledger.
// Returns the contract ID of the first created contract in the committed transaction.
func (c *GRPCLedgerClient) ExecuteSubmission(
	ctx context.Context,
	preparedTx *interactive.PreparedTransaction,
	partySignatures *interactive.PartySignatures,
	hashingSchemeVersion interactive.HashingSchemeVersion,
) (string, error) {
	resp, err := c.interactive.ExecuteSubmissionAndWaitForTransaction(ctx, &interactive.ExecuteSubmissionAndWaitForTransactionRequest{
		PreparedTransaction:  preparedTx,
		PartySignatures:      partySignatures,
		HashingSchemeVersion: hashingSchemeVersion,
		SubmissionId:         uuid.New().String(),
	})
	if err != nil {
		return "", fmt.Errorf("ExecuteSubmission: %w", err)
	}
	// Extract the first created contract ID from the committed transaction.
	for _, event := range resp.GetTransaction().GetEvents() {
		if ce := event.GetCreated(); ce != nil && ce.GetContractId() != "" {
			return ce.GetContractId(), nil
		}
	}

	return "", nil
}

func (c *GRPCLedgerClient) GrantPartyRights(ctx context.Context, userID, partyID string) error {
	_, err := c.userManagement.GrantUserRights(ctx, &apiv2admin.GrantUserRightsRequest{
		UserId: userID,
		Rights: []*apiv2admin.Right{
			{Kind: &apiv2admin.Right_CanActAs_{CanActAs: &apiv2admin.Right_CanActAs{Party: partyID}}},
			{Kind: &apiv2admin.Right_CanReadAs_{CanReadAs: &apiv2admin.Right_CanReadAs{Party: partyID}}},
		},
	})
	if err != nil {
		return fmt.Errorf("GrantUserRights for user %q party %q: %w", userID, partyID, err)
	}

	return nil
}
