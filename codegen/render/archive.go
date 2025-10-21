package render

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf/types"
)

type ArchiveData struct {
	Name         string
	MainPackage  *PackageData
	Dependencies []*PackageData
}

func ArchiveDataFromLF(darFile *daml_lf.DarFile) (*ArchiveData, error) {
	mainPackage, err := PackageDataFromLFV2(darFile.Main.Hash, darFile.Main.Payload.Package.V2)
	if err != nil {
		return nil, fmt.Errorf("could not load main package: %w", err)
	}
	var dependencies []*PackageData
	for _, dependency := range darFile.Dependencies {
		pkg, err := PackageDataFromLFV2(dependency.Hash, dependency.Payload.Package.V2)
		if err != nil {
			return nil, fmt.Errorf("could not load dependency package: %w", err)
		}
		dependencies = append(dependencies, pkg)
	}
	return &ArchiveData{
		Name:         darFile.Main.Name,
		MainPackage:  mainPackage,
		Dependencies: dependencies,
	}, nil
}

func (a *ArchiveData) GetGeneratorContext(context context.Context, goModulePrefix string) types.GeneratorContext {
	ctx := types.GeneratorContext{
		Context:              context,
		Packages:             make(map[string]types.PackageContext),
		GoModulePrefix:       goModulePrefix,
		IncludePackagePrefix: false,
		UsedDependencies:     make(map[string]bool),
	}

	// Add all packages
	for _, p := range append(a.Dependencies, a.MainPackage) {
		pkg := types.PackageContext{
			Name:      p.Name,
			PackageId: p.PackageId,
			Modules:   make(map[string]types.ModuleContext),
		}
		for _, module := range p.Modules {
			pkg.Modules[module.Name] = types.ModuleContext{
				Name: module.Name,
			}
		}
		ctx.Packages[p.PackageId] = pkg
	}

	return ctx
}
