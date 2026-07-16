package adapters

import (
	"fmt"

	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

// cantonSigningIdentityReader reads OnchainSigningPubKey from the JD OCRKeyBundle.
// Canton's signer identity is the full uncompressed secp256k1 public key, not the
// EVM-derived address.
type cantonSigningIdentityReader struct{}

func (cantonSigningIdentityReader) FromBundle(b *nodev1.OCR2Config_OCRKeyBundle) (string, error) {
	if b == nil {
		return "", fmt.Errorf("nil OCR key bundle")
	}
	if b.OnchainSigningPubKey == "" {
		return "", fmt.Errorf("OnchainSigningPubKey is empty")
	}

	return b.OnchainSigningPubKey, nil
}
