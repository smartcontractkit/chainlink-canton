package render

import (
	"fmt"

	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type PackageData struct {
	Name      string
	PackageId string
	Modules   []*ModuleData
}

func PackageDataFromLFV2(packageId string, pkg *v2.Package) (*PackageData, error) {
	packageName := pkg.GetInternedStrings()[pkg.GetMetadata().GetNameInternedStr()]

	var modules []*ModuleData
	for _, module := range pkg.GetModules() {
		m, err := ModuleDataFromLF(module, pkg)
		if err != nil {
			return nil, fmt.Errorf("could not determine module data: %w", err)
		}
		modules = append(modules, m)
	}

	return &PackageData{
		Name:      packageName,
		PackageId: packageId,
		Modules:   modules,
	}, nil
}
