package render

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf/types"
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type DataType struct {
	Name            string
	DataConstructor DataConstructor
}

func DataTypeFromLF(dataType *v2.DefDataType, pkg *v2.Package) (*DataType, error) {
	typeName, err := types.GetDottedName(dataType.NameInternedDname, pkg)
	if err != nil {
		return nil, fmt.Errorf("could not determine data type name: %w", err)
	}

	switch typ := dataType.GetDataCons().(type) {
	case *v2.DefDataType_Record:
		// TODO
		ss := strings.Split(typeName, ".")
		typeName = ss[len(ss)-1]
		r, err := RecordFromLF(typ, pkg, typeName)
		if err != nil {
			return nil, fmt.Errorf("could not construct record type: %w", err)
		}
		return &DataType{
			Name:            typeName,
			DataConstructor: r,
		}, nil
	case *v2.DefDataType_Variant:
		// TODO
		v, err := VariantFromLF(typ, pkg, typeName)
		if err != nil {
			return nil, fmt.Errorf("could not construct record type: %w", err)
		}
		return &DataType{
			Name:            typeName,
			DataConstructor: v,
		}, nil
	case *v2.DefDataType_Enum:
	case *v2.DefDataType_Interface:
	default:
		return nil, fmt.Errorf("unknown data type: %T", typ)
	}
	// TODO
	return nil, nil
	return nil, fmt.Errorf("unsupported data type: %T", dataType.GetDataCons())
}

type DataConstructor interface {
	Name(ctx types.GeneratorContext) string
	// GoTypes returns the type as referenced from within a Go program.
	// It's the string the type would be referred to when instantiating it:
	//  var a <GoType>
	// Of e.g. when using it as a function argument:
	//  func Test(a <GoType>)
	// If it is non-local, it must include the package prefix:
	//  var a <packagename.Type>
	GoTypes(ctx types.GeneratorContext) []string
	// GoImports returns additional imports that are required to use the type.
	// In case of a local type, nil will be returned
	//  []string{"github.com/smartcontractkit/chainlink-canton-internal/pkg/types"}
	GoImports(ctx types.GeneratorContext) []string
	// GoTypeDefinitions returns the full definition for the given type.
	// E.g. in case of a struct:
	//  type TypeName struct {
	//  	FieldA string
	//  }
	// If it is a non-local type, or no type definition is required to use it, nil will be returned.
	GoTypeDefinitions(ctx types.GeneratorContext, pkg *PackageData, mod *ModuleData) []string
}

// Record

func RecordFromLF(record *v2.DefDataType_Record, pkg *v2.Package, name string) (*Record, error) {
	fields := make(map[string]types.DamlType)

	for _, fieldWithType := range record.Record.GetFields() {
		fieldName := pkg.GetInternedStrings()[fieldWithType.GetFieldInternedStr()]
		fieldType, err := types.TypeFromLF(fieldWithType.GetType(), pkg)
		if err != nil {
			return nil, err
		}
		fields[fieldName] = fieldType
	}

	return &Record{
		StructName: name,
		Fields:     fields,
	}, nil
}

type Record struct {
	StructName string
	Fields     map[string]types.DamlType
}

func (r Record) Name(ctx types.GeneratorContext) string {
	return r.StructName
}

func (r Record) GoTypes(ctx types.GeneratorContext) []string {
	// TODO should return fields instead?
	return []string{r.StructName}
}

func (r Record) GoImports(ctx types.GeneratorContext) []string {
	imports := make([]string, 0, len(r.Fields))
	for damlType := range maps.Values(r.Fields) {
		imports = append(imports, damlType.GoImports(ctx)...)
	}
	return imports
}

func (r Record) GoTypeDefinitions(ctx types.GeneratorContext, pkg *PackageData, mod *ModuleData) []string {
	fields := slices.Collect(maps.Keys(r.Fields))
	slices.Sort(fields)

	fStrings := make([]string, 0, len(r.Fields))
	for _, fieldName := range fields {
		damlType := r.Fields[fieldName]
		// Capitalize struct fields
		fieldName = cases.Title(language.Und).String(fieldName)
		fStrings = append(fStrings, fmt.Sprintf("%s %s", fieldName, damlType.GoType(ctx)))
	}
	return []string{fmt.Sprintf(`
	type %s struct {
		%s
	}
	`, r.StructName, strings.Join(fStrings, "\n")),
	}
}

// Variant

func VariantFromLF(record *v2.DefDataType_Variant, pkg *v2.Package, name string) (*Variant, error) {
	fields := make(map[string]types.DamlType)

	for _, fieldWithType := range record.Variant.GetFields() {
		fieldName := pkg.GetInternedStrings()[fieldWithType.GetFieldInternedStr()]
		fieldType, err := types.TypeFromLF(fieldWithType.GetType(), pkg)
		if err != nil {
			return nil, err
		}
		fields[fieldName] = fieldType
	}

	return &Variant{
		InterfaceName: name,
		Fields:        fields,
	}, nil
}

type Variant struct {
	InterfaceName string
	Fields        map[string]types.DamlType
}

func (r Variant) Name(ctx types.GeneratorContext) string {
	return r.InterfaceName
}

func (r Variant) GoTypes(ctx types.GeneratorContext) []string {
	return []string{r.InterfaceName}
}

func (r Variant) GoImports(ctx types.GeneratorContext) []string {
	imports := make([]string, 0, len(r.Fields))
	for damlType := range maps.Values(r.Fields) {
		imports = append(imports, damlType.GoImports(ctx)...)
	}
	return imports
}

func (r Variant) GoTypeDefinitions(ctx types.GeneratorContext, pkg *PackageData, mod *ModuleData) []string {
	var newTypeDefs []string
	var interfaceHints []string

	// Sort keys, to make output deterministic
	fields := slices.Collect(maps.Keys(r.Fields))
	slices.Sort(fields)

	for _, fieldName := range fields {
		damlType := r.Fields[fieldName]
		// Check if this is a DataType defined in the module
		// If it is, it has been skipped when generating the module's DataTypes, since
		// we'd like to define it next to the variant interface that it implements
		isDataType := false
		for _, dataType := range mod.DataTypes {
			// TODO
			if dataType == nil {
				continue
			}
			if damlType.GoType(ctx) == dataType.DataConstructor.Name(ctx) {
				// Found
				newTypeDefs = append(newTypeDefs, dataType.DataConstructor.GoTypeDefinitions(ctx, pkg, mod)...)
				isDataType = true
			}
		}

		// If it isn't a defined DataType, generate a type definition manually
		if !isDataType {
			newTypeDefs = append(newTypeDefs, fmt.Sprintf(`
type %s %s
		`, fieldName, damlType.GoType(ctx)))
		}

		// Interface implementation
		newTypeDefs = append(newTypeDefs, fmt.Sprintf(`
func (v %s) _is%s() {}
		`, fieldName, r.InterfaceName))

		interfaceHints = append(interfaceHints, fieldName)
	}

	typeHint := strings.Builder{}
	for _, hint := range interfaceHints {
		typeHint.WriteString(fmt.Sprintf("\n//  %s", hint))
	}

	typeDefs := []string{
		fmt.Sprintf(`
// Variant %s
// Types that are valid to be assigned to %s:%s
type %s interface {
	_is%s()
}	`, r.InterfaceName, r.InterfaceName, typeHint.String(), r.InterfaceName, r.InterfaceName),
	}

	return append(typeDefs, newTypeDefs...)
}
