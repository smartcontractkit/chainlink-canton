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
	"google.golang.org/grpc/credentials/insecure"
)

const defaultLedgerPort = 5002

// DialLedger opens a gRPC connection to a Canton Ledger API endpoint using
// the ledger host, port, and optional JWT bearer token from cfg.
// TLS is not yet supported — connections are plaintext.
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

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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
	interactive interactive.InteractiveSubmissionServiceClient
	party       apiv2admin.PartyManagementServiceClient
	state       apiv2.StateServiceClient
}

// NewGRPCLedgerClient creates a new [GRPCLedgerClient] from an established gRPC connection.
func NewGRPCLedgerClient(conn *grpc.ClientConn) *GRPCLedgerClient {
	return &GRPCLedgerClient{
		interactive: interactive.NewInteractiveSubmissionServiceClient(conn),
		party:       apiv2admin.NewPartyManagementServiceClient(conn),
		state:       apiv2.NewStateServiceClient(conn),
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
		SynchronizerId: synchronizerID,
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
