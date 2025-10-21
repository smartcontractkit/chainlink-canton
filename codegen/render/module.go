package render

import (
	"fmt"
	"maps"
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf/types"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/render/tmpl"
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type ModuleData struct {
	Name      string
	DataTypes []*DataType
	Templates []*TemplateData
}

func ModuleDataFromLF(module *v2.Module, pkg *v2.Package) (*ModuleData, error) {
	moduleName, err := types.GetDottedName(module.GetNameInternedDname(), pkg)
	if err != nil {
		return nil, fmt.Errorf("could not determine module name: %w", err)
	}

	var dataTypes []*DataType
	for _, dataType := range module.GetDataTypes() {
		dt, err := DataTypeFromLF(dataType, pkg)
		if err != nil {
			return nil, err
		}
		dataTypes = append(dataTypes, dt)
	}

	var templates []*TemplateData
	for _, template := range module.Templates {
		t, err := TemplateFromLF(template, pkg)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return &ModuleData{
		Name:      moduleName,
		DataTypes: dataTypes,
		Templates: templates,
	}, nil
}

func (m *ModuleData) GetDependencies() []string {
	var deps []string

	for _, dataType := range m.DataTypes {
		_ = dataType
	}

	return deps
}

func (m *ModuleData) GetModuleData(ctx types.GeneratorContext, pkg *PackageData, packageId string) (tmpl.ModuleData, error) {
	moduleInfo, err := ctx.GetModuleInformation(packageId, m.Name)
	if err != nil {
		return tmpl.ModuleData{}, err
	}

	var imports []string

	// Generate Type Definitions
	var typeDefs []string
	for _, dataType := range m.DataTypes {
		// TODO
		if dataType == nil {
			continue
		}

		// Skip DataTypes that are arguments to templates/choices
		skip := false
		for _, template := range m.Templates {
			if dataType.Name == template.Name {
				skip = true
				break
			}
			for _, choice := range template.Choices {
				if choice.Name == dataType.Name {
					skip = true
					break
				}
			}
		}

		// Skip DataTypes that are part of the fields of a variant
		for _, d := range m.DataTypes {
			// TODO
			if d == nil {
				continue
			}
			if variant, ok := d.DataConstructor.(*Variant); ok {
				for _, damlType := range variant.Fields {
					if damlType.GoType(ctx) == dataType.DataConstructor.Name(ctx) {
						skip = true
						break
					}
				}
			}
		}

		if skip {
			continue
		}

		typeDefs = append(typeDefs, dataType.DataConstructor.GoTypeDefinitions(ctx, pkg, m)...)
		imports = append(imports, dataType.DataConstructor.GoImports(ctx)...)
	}

	var templates []tmpl.TemplateData
	for _, template := range m.Templates {
		var constructorType *DataType
		for _, dataType := range m.DataTypes {
			if dataType.Name == template.Name {
				constructorType = dataType
				break
			}
		}
		if constructorType == nil {
			return tmpl.ModuleData{}, fmt.Errorf("could not find constructor type for template: %s", template.Name)
		}

		var constructorArgs []tmpl.Argument
		switch conType := constructorType.DataConstructor.(type) {
		case *Record:
			fields := slices.Collect(maps.Keys(conType.Fields))
			slices.Sort(fields)
			for _, fieldName := range fields {
				damlType := conType.Fields[fieldName]
				fieldName = cases.Title(language.Und).String(fieldName)
				argType := damlType.GoType(ctx)
				constructorArgs = append(constructorArgs, tmpl.Argument{
					Name: fieldName,
					Type: argType,
				})
				imports = append(imports, damlType.GoImports(ctx)...)
			}
		}

		templates = append(templates, tmpl.TemplateData{
			Name:            template.Name,
			ConstructorArgs: constructorArgs,
		})
	}

	// Remove empty + duplicate imports and sort them
	imports = slices.DeleteFunc(imports, func(s string) bool {
		return s == ""
	})
	slices.Sort(imports)
	imports = slices.Compact(imports)

	return tmpl.ModuleData{
		ModuleName:      moduleInfo.PackageName,
		Imports:         imports,
		TypeDefinitions: typeDefs,
		Templates:       templates,
	}, nil
}
