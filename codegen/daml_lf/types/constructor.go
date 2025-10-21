package types

import (
	"fmt"
	"strings"

	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

func TypeConstructorFromLF(tyCon *v2.Type_Con, pkg *v2.Package) (DamlType, error) {
	tyConName, err := GetDottedName(tyCon.GetTycon().GetNameInternedDname(), pkg)
	if err != nil {
		return nil, err
	}
	// TODO
	ss := strings.Split(tyConName, ".")
	tyConName = ss[len(ss)-1]
	moduleName, err := GetDottedName(tyCon.GetTycon().GetModule().GetModuleNameInternedDname(), pkg)
	if err != nil {
		return nil, err
	}
	switch typeConstructor := tyCon.GetTycon().GetModule().GetPackageId().GetSum().(type) {
	case *v2.SelfOrImportedPackageId_SelfPackageId:
		return TypeConstructorSelf{
			Name: tyConName,
		}, nil
	case *v2.SelfOrImportedPackageId_ImportedPackageIdInternedStr:
		return TypeConstructorImportedInternedStr{
			Name:               tyConName,
			ImportedPackageId:  pkg.GetInternedStrings()[typeConstructor.ImportedPackageIdInternedStr],
			ImportedModuleName: moduleName,
		}, nil
	case *v2.SelfOrImportedPackageId_PackageImportId:
		return TypeConstructorImported{
			Name:      tyConName,
			PackageId: "externalPackage",
		}, nil
	default:
		return nil, fmt.Errorf("unknown type constructor type: %T", typeConstructor)
	}
}

type TypeConstructorSelf struct {
	noImportOrDefinition
	Name string
}

func (t TypeConstructorSelf) GoType(ctx GeneratorContext) string {
	return t.Name
}

type TypeConstructorImportedInternedStr struct {
	Name               string
	ImportedPackageId  string
	ImportedModuleName string
}

func (t TypeConstructorImportedInternedStr) GoType(ctx GeneratorContext) string {
	moduleInfo, err := ctx.GetModuleInformation(t.ImportedPackageId, t.ImportedModuleName)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s.%s", moduleInfo.PackageName, t.Name)
}

func (t TypeConstructorImportedInternedStr) GoImports(ctx GeneratorContext) []string {
	ctx.MarkDependencyAsUsed(t.ImportedPackageId)
	moduleInfo, err := ctx.GetModuleInformation(t.ImportedPackageId, t.ImportedModuleName)
	if err != nil {
		panic(err)
	}
	return []string{fmt.Sprintf("\"%s\"", moduleInfo.ImportPath)}
}

func (t TypeConstructorImportedInternedStr) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}

type TypeConstructorImported struct {
	Name      string
	PackageId string
}

func (t TypeConstructorImported) GoType(ctx GeneratorContext) string {
	return t.Name
}

func (t TypeConstructorImported) GoImports(ctx GeneratorContext) []string {
	// TODO implement me
	panic("implement me")
}

func (t TypeConstructorImported) GoTypeDefinitions(ctx GeneratorContext) []string {
	// TODO implement me
	panic("implement me")
}
