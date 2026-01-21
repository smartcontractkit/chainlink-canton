package client

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"
)

const (
	UserName = "ledger-api-user"
	Audience = "https://canton.network.global"
)

// Config holds configuration for setting up a Canton client
type Config struct {
	LedgerAPIURL      string
	AdminAPIURL       string
	JWTSecret         string // Optional, defaults to "unsafe"
	DeployerParty     string // Optional, will allocate if not provided
	DeployerPartyHint string // Optional hint for party allocation
}

// SetupResult contains the result of client setup
type SetupResult struct {
	BindingClient *client.DamlBindingClient
	Party         string
	UserID        string
}

// CantonOpDeps contains dependencies for Canton operations
// This is defined here to avoid import cycles between testenv and ops packages
type CantonOpDeps struct {
	BindingClient *client.DamlBindingClient
	Party         string
	UserID        string
}

func GetJWT() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   UserName,
		Audience:  []string{Audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	})
	return t.SignedString([]byte("unsafe"))
}

// NewBindingClient creates a new DAML binding client
func NewBindingClient(ctx context.Context, jwtToken, ledgerAPIURL, adminAPIURL string) (*client.DamlBindingClient, error) {
	bindingClient, err := client.NewDamlClient(jwtToken, ledgerAPIURL).
		WithAdminAddress(adminAPIURL).
		Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create DAML binding client: %w", err)
	}
	return bindingClient, nil
}

// GetOrAllocateParty gets the primary party for a user or allocates a new party
func GetOrAllocateParty(ctx context.Context, bindingClient *client.DamlBindingClient, party, partyHint string) (string, error) {
	if party != "" {
		return party, nil
	}

	// Try to get primary party for user
	user, err := bindingClient.UserMng.GetUser(ctx, UserName)
	if err == nil && user != nil && user.PrimaryParty != "" {
		return user.PrimaryParty, nil
	}

	// Allocate new party
	if partyHint == "" {
		partyHint = "canton-deployer"
	}
	resp, err := bindingClient.PartyMng.AllocateParty(ctx, partyHint, map[string]string{"type": "testing"}, "")
	if err != nil {
		return "", fmt.Errorf("failed to allocate party: %w", err)
	}
	return resp.Party, nil
}

// EnsureUserRights grants the user rights to act as and read as the specified party
func EnsureUserRights(ctx context.Context, bindingClient *client.DamlBindingClient, party string) error {
	_, err := bindingClient.UserMng.GrantUserRights(ctx, UserName, "", []*model.Right{
		{Type: model.CanActAs{Party: party}},
		{Type: model.CanReadAs{Party: party}},
	})
	if err != nil {
		// Ignore error if rights already exist
		// This is a common case and not necessarily an error
		return nil
	}
	return nil
}

// Setup creates a Canton client setup with all necessary dependencies
// It handles JWT generation, client creation, party management, and user rights
func Setup(ctx context.Context, config Config) (*SetupResult, error) {
	// Get JWT token
	jwtToken, err := GetJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Create DAML binding client
	bindingClient, err := NewBindingClient(ctx, jwtToken, config.LedgerAPIURL, config.AdminAPIURL)
	if err != nil {
		return nil, err
	}

	// Get or allocate deployer party
	deployerParty, err := GetOrAllocateParty(ctx, bindingClient, config.DeployerParty, config.DeployerPartyHint)
	if err != nil {
		bindingClient.Close()
		return nil, err
	}

	// Ensure user has rights to act as the party
	if err := EnsureUserRights(ctx, bindingClient, deployerParty); err != nil {
		bindingClient.Close()
		return nil, fmt.Errorf("failed to ensure user rights: %w", err)
	}

	return &SetupResult{
		BindingClient: bindingClient,
		Party:         deployerParty,
		UserID:        UserName,
	}, nil
}
