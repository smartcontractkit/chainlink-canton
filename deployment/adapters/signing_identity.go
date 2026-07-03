package adapters

import (
	"fmt"
	"strings"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	ccvshared "github.com/smartcontractkit/chainlink-ccv/deployment/shared"
)

type CantonSigningIdentityReader struct{}

var _ ccvshared.SigningIdentityReader = (*CantonSigningIdentityReader)(nil)

func (CantonSigningIdentityReader) FromBundle(bundle *nodev1.OCR2Config_OCRKeyBundle) (string, error) {
	pubKey := strings.TrimSpace(bundle.OnchainSigningPubKey)
	if pubKey == "" {
		return "", fmt.Errorf("canton: missing onchain_signing_pub_key")
	}
	pubKey = strings.ToLower(strings.TrimPrefix(pubKey, "0x"))
	return "0x" + pubKey, nil
}
