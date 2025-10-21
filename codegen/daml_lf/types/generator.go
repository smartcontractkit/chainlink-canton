package types

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type ModuleContext struct {
	Name string
}

type PackageContext struct {
	Name      string
	PackageId string
	Modules   map[string]ModuleContext
}

type GeneratorContext struct {
	context.Context
	// Whether to include the package name in the generated files
	IncludePackagePrefix bool
	GoModulePrefix       string
	Packages             map[string]PackageContext

	UsedDependencies map[string]bool
}

type ModuleInformation struct {
	// The full import path for this module
	ImportPath string
	// The path of this module (relative to the GoModulePrefix)
	RelPath string
	// The filename (without .go) to use for this module
	Filename string
	// The package name of this module
	PackageName string
}

func (ctx *GeneratorContext) GetModuleInformation(packageId string, moduleName string) (ModuleInformation, error) {
	pkg, ok := ctx.Packages[packageId]
	if !ok {
		return ModuleInformation{}, fmt.Errorf("package %s not found", packageId)
	}
	mod, ok := pkg.Modules[moduleName]
	if !ok {
		return ModuleInformation{}, fmt.Errorf("module %s not found", moduleName)
	}

	// github.com, smartcontractkit, chainlink-canton-internal
	prefix := strings.Split(ctx.GoModulePrefix, "/")
	var pkgName []string
	if ctx.IncludePackagePrefix {
		pkgName = append(pkgName, strings.ToLower(pkg.Name))
	}
	modName := strings.Split(strings.ToLower(mod.Name), ".")

	importPath := append(prefix, append(pkgName, modName...)...)
	relPath := append([]string{"."}, append(pkgName, modName[:len(modName)-1]...)...)
	filename := importPath[len(importPath)-1]
	packageName := importPath[len(importPath)-2]
	importPath = importPath[:len(importPath)-1]

	return ModuleInformation{
		ImportPath:  strings.Join(importPath, "/"),
		RelPath:     filepath.Join(relPath...),
		Filename:    filename,
		PackageName: packageName,
	}, nil
}

func (ctx *GeneratorContext) MarkDependencyAsUsed(packageId string) {
	ctx.UsedDependencies[packageId] = true
}
