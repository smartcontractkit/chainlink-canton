package keys

import (
	"context"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

func obtainSigningKey(
	ctx context.Context,
	c client.CantonClient,
	kmsKeyID string,
	name string,
	usage []cryptov30.SigningKeyUsage,
) (*cryptov30.SigningPublicKey, error) {
	if kmsKeyID != "" {
		return c.RegisterKmsSigningKey(ctx, kmsKeyID, name, usage)
	}

	return c.GenerateSigningKey(ctx, name, usage)
}
