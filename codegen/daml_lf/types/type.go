package types

import (
	"fmt"

	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type DamlType interface {
	// GoType returns the type as referenced from within a Go program.
	// It's the string the type would be referred to when instantiating it:
	//  var a <GoType>
	// Of e.g. when using it as a function argument:
	//  func Test(a <GoType>)
	// If it is non-local, it must include the package prefix:
	//  var a <packagename.Type>
	GoType(ctx GeneratorContext) string
	// GoImports returns additional imports that are required to use the type.
	// In case of a local type, nil will be returned
	//  []string{"github.com/smartcontractkit/chainlink-canton-internal/pkg/types"}
	GoImports(ctx GeneratorContext) []string
	// GoTypeDefinitions returns the full definition for the given type.
	// E.g. in case of a struct:
	//  type TypeName struct {
	//  	FieldA string
	//  }
	// If it is a non-local type, or no type definition is required to use it, nil will be returned.
	GoTypeDefinitions(ctx GeneratorContext) []string
}

func TypeFromLF(typ *v2.Type, pkg *v2.Package) (DamlType, error) {
	switch typ := typ.GetSum().(type) {
	case *v2.Type_Con_:
		return TypeConstructorFromLF(typ.Con, pkg)
	case *v2.Type_Builtin_:
		return BuiltinFromLF(typ.Builtin, pkg)
	case *v2.Type_Forall_:
	case *v2.Type_Struct_:
	case *v2.Type_Nat:
	case *v2.Type_InternedType:
		return InternedFromLF(typ, pkg)
	case *v2.Type_Tapp:
	default:
		return nil, fmt.Errorf("unknown type: %T", typ)
	case *v2.Type_Var_:
		return VarFromLF(typ.Var, pkg)
	}
	// TODO
	return nil, fmt.Errorf("unsupported type: %T", typ.GetSum())
}

// Helper types

type nonRenderable struct{}

func (n nonRenderable) GoType(ctx GeneratorContext) string {
	panic("type is non-renderable")
}

func (n nonRenderable) GoImports(ctx GeneratorContext) []string {
	panic("type is non-renderable")
}

func (n nonRenderable) GoTypeDefinitions(ctx GeneratorContext) []string {
	panic("type is non-renderable")
}

type noImportOrDefinition struct{}

func (n noImportOrDefinition) GoImports(ctx GeneratorContext) []string {
	return nil
}

func (n noImportOrDefinition) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}
