package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

type Client struct {
	conn         *grpc.ClientConn
	stateService apiv2.StateServiceClient
	jwtSecret    string
	jwtAudience  string
}

func NewClient(target string, jwtSecret, jwtAudience string) (*Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ledger API: %w", err)
	}

	return &Client{
		conn:         conn,
		stateService: apiv2.NewStateServiceClient(conn),
		jwtSecret:    jwtSecret,
		jwtAudience:  jwtAudience,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) AuthContext(ctx context.Context) (context.Context, error) {
	token, err := c.generateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
	return metadata.NewOutgoingContext(ctx, md), nil
}

func (c *Client) generateJWT() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "ledger-api-user",
		Audience:  []string{c.jwtAudience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return t.SignedString([]byte(c.jwtSecret))
}

func (c *Client) GetCurrentOffset(ctx context.Context) (int64, error) {
	resp, err := c.stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger end: %w", err)
	}
	return resp.GetOffset(), nil
}

func (c *Client) StateService() apiv2.StateServiceClient {
	return c.stateService
}
