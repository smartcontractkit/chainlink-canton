package main

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

func getJWT() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   "ledger-api-user",
		Audience:  []string{"https://canton.network.global"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	})
	return t.SignedString([]byte("unsafe"))
}

func main() {
	jwtToken, err := getJWT()
	if err != nil {
		panic(err)
	}
	fmt.Println("JWT:", jwtToken)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwtToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	grpcClient, err := grpc.NewClient("participant2.grpc-ledger-api.localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// grpcClient, err := grpc.NewClient("localhost:1301", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	vClient := apiv2.NewVersionServiceClient(grpcClient)
	versionResp, err := vClient.GetLedgerApiVersion(ctx, &apiv2.GetLedgerApiVersionRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Ledger API version: ", versionResp.Version)

	packageClient := apiv2.NewPackageServiceClient(grpcClient)
	packageResp, err := packageClient.ListPackages(ctx, &apiv2.ListPackagesRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Available Package IDs: ", packageResp.GetPackageIds())
}
