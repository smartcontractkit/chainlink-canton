package registry

import (
	"context"
	"fmt"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

// UtilityPackages lists vendored Canton Network Utility DARs for upload (dependency order).
func UtilityPackages() []contracts.Package {
	return []contracts.Package{
		contracts.UtilityCommercialsV0,
		contracts.UtilityCredentialV0,
		contracts.UtilityCredentialAppV0,
		contracts.UtilityRegistryV0,
		contracts.UtilityRegistryHoldingV0,
		contracts.UtilityRegistryAppV0,
	}
}

// LoadUtilityDARs loads Registry utility DAR bytes for the current contracts version.
func LoadUtilityDARs() ([][]byte, error) {
	var dars [][]byte
	for _, pkg := range UtilityPackages() {
		dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
		if err != nil {
			return nil, fmt.Errorf("load DAR for %s: %w", pkg, err)
		}
		dars = append(dars, dar)
	}

	return dars, nil
}

// RequiredRegistryPackages are utility DARs needed for registry-kit Registry flows.
var RequiredRegistryPackages = map[string]string{
	"utility-registry-app-v0":     contracts.UtilityRegistryAppV0PackageID,
	"utility-registry-v0":         contracts.UtilityRegistryV0PackageID,
	"utility-registry-holding-v0": contracts.UtilityRegistryHoldingV0PackageID,
	"utility-credential-app-v0":   contracts.UtilityCredentialAppV0PackageID,
	"utility-credential-v0":       contracts.UtilityCredentialV0PackageID,
}

// PackageCheckResult is one required package presence check.
type PackageCheckResult struct {
	Name     string
	Expected string
	Found    bool
	FoundID  string
}

// CheckPackages verifies pinned Registry utility package IDs exist on the participant.
func CheckPackages(ctx context.Context, participant canton.Participant) ([]PackageCheckResult, error) {
	if participant.LedgerServices.Admin.PackageManagement == nil {
		return nil, fmt.Errorf("package management client unavailable (check admin_api_url and auth)")
	}

	resp, err := participant.LedgerServices.Admin.PackageManagement.ListKnownPackages(ctx, &adminv2.ListKnownPackagesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list known packages: %w", err)
	}
	installed := map[string]string{}
	for _, pkg := range resp.GetPackageDetails() {
		installed[pkg.GetPackageId()] = pkg.GetPackageId()
	}

	results := make([]PackageCheckResult, 0, len(RequiredRegistryPackages))
	missing := 0
	for name, expectedID := range RequiredRegistryPackages {
		_, found := installed[expectedID]
		if !found {
			missing++
		}
		results = append(results, PackageCheckResult{
			Name:     name,
			Expected: expectedID,
			Found:    found,
			FoundID:  expectedID,
		})
	}
	if missing > 0 {
		return results, fmt.Errorf("%d required Registry packages missing on participant", missing)
	}

	return results, nil
}
