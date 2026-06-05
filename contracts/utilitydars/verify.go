package utilitydars

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/smartcontractkit/go-daml/codegen"
)

// VerifyPackageID checks that darBytes contains the expected main package ID.
func VerifyPackageID(darBytes []byte, expectedID string) error {
	reader, err := zip.NewReader(bytes.NewReader(darBytes), int64(len(darBytes)))
	if err != nil {
		return fmt.Errorf("open DAR zip: %w", err)
	}

	manifest, err := codegen.GetManifest(reader)
	if err != nil {
		return fmt.Errorf("parse DAR manifest: %w", err)
	}

	actual := codegen.GetPackageID(manifest.MainDalf)
	if !strings.EqualFold(actual, expectedID) {
		return fmt.Errorf("package ID mismatch: got %s, want %s", actual, expectedID)
	}

	return nil
}

// MainPackageID returns the main package ID from a DAR byte slice.
func MainPackageID(darBytes []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(darBytes), int64(len(darBytes)))
	if err != nil {
		return "", fmt.Errorf("open DAR zip: %w", err)
	}

	manifest, err := codegen.GetManifest(reader)
	if err != nil {
		return "", fmt.Errorf("parse DAR manifest: %w", err)
	}

	return codegen.GetPackageID(manifest.MainDalf), nil
}
