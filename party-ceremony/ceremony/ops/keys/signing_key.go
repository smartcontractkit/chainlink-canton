package keys

import (
	"context"
	"fmt"

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
		key, err := c.RegisterKmsSigningKey(ctx, kmsKeyID, name, usage)
		if err == nil {
			return key, nil
		}
		if client.IsKmsPublicKeyConflict(err) {
			existing, lookupErr := c.LookupKmsSigningKey(ctx, kmsKeyID, usage)
			if lookupErr != nil {
				return nil, fmt.Errorf("register KMS key %q as %q: %w; lookup existing key: %v", kmsKeyID, name, err, lookupErr)
			}

			return existing, nil
		}

		return nil, err
	}

	return c.GenerateSigningKey(ctx, name, usage)
}

func vaultRegistrationName(namespaceName, kmsVaultName string) string {
	if kmsVaultName != "" {
		return kmsVaultName
	}

	return namespaceName
}
