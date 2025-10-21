package types

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

func BuiltinFromLF(builtin *v2.Type_Builtin, pkg *v2.Package) (DamlType, error) {
	switch builtin.Builtin {
	// Builtin type 'Unit'
	case v2.BuiltinType_UNIT:
		return Unit{}, nil
	// Builtin type 'Bool'
	case v2.BuiltinType_BOOL:
		return Boolean{}, nil
	// Builtin type 'Int64'
	case v2.BuiltinType_INT64:
		return Int64{}, nil
	// Builtin type 'Date'
	case v2.BuiltinType_DATE:
		return Date{}, nil
	// Builtin type 'Timestamp'
	case v2.BuiltinType_TIMESTAMP:
		return Timestamp{}, nil
	// Builtin type 'Numeric'
	case v2.BuiltinType_NUMERIC:
		return Numeric{}, nil
	// Builtin tpe 'Party'
	case v2.BuiltinType_PARTY:
		return Party{}, nil
	// Builtin type 'Text'
	case v2.BuiltinType_TEXT:
		return Text{}, nil
	// Builtin type 'ContractId'
	case v2.BuiltinType_CONTRACT_ID:
		return ContractId{}, nil
	// Builtin type 'Optional'
	case v2.BuiltinType_OPTIONAL:
		if len(builtin.GetArgs()) <= 0 {
			return nil, fmt.Errorf("got Optional without any args")
		}
		if len(builtin.GetArgs()) > 1 {
			return nil, fmt.Errorf("got Optional with more than one arg")
		}
		vType, err := TypeFromLF(builtin.GetArgs()[0], pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to determine type of Optional: %w", err)
		}
		return Optional{
			Value: vType,
		}, nil
	// Builtin type 'List'
	case v2.BuiltinType_LIST:
		if len(builtin.GetArgs()) <= 0 {
			return nil, fmt.Errorf("got List without any args")
		}
		if len(builtin.GetArgs()) > 1 {
			return nil, fmt.Errorf("got List with more than one arg")
		}
		vType, err := TypeFromLF(builtin.GetArgs()[0], pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to determine type of List: %w", err)
		}
		return List{
			Value: vType,
		}, nil
	// Builtin type 'TGenMap`
	case v2.BuiltinType_GENMAP:
		if len(builtin.GetArgs()) <= 0 {
			return nil, fmt.Errorf("got GenMap without any args")
		}
		if len(builtin.GetArgs()) > 2 {
			return nil, fmt.Errorf("got GenMap with more than two args")
		}
		kType, err := TypeFromLF(builtin.GetArgs()[0], pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to determine key type of GenMap: %w", err)
		}
		vType, err := TypeFromLF(builtin.GetArgs()[1], pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to determine value type of GenMap: %w", err)
		}
		return GenMap{
			Key:   kType,
			Value: vType,
		}, nil
	// Builtin type 'Any'
	case v2.BuiltinType_ANY:
		return Any{}, nil
	// Builtin type 'TAnyException'
	case v2.BuiltinType_ANY_EXCEPTION:
		return AnyException{}, nil
	// Builtin type 'TypeRep'
	case v2.BuiltinType_TYPE_REP:
		return TypeRep{}, nil
	// Builtin type `TArrow`
	case v2.BuiltinType_ARROW:
		return Arrow{}, nil
	// Builtin type 'Update'
	case v2.BuiltinType_UPDATE:
		return Update{}, nil
	// Builtin type for FailureCategory
	case v2.BuiltinType_FAILURE_CATEGORY:
		return FailureCategory{}, nil
	// Builtin type 'TTextMap`
	case v2.BuiltinType_TEXTMAP:
		if len(builtin.GetArgs()) <= 0 {
			return nil, fmt.Errorf("got List without any args")
		}
		if len(builtin.GetArgs()) > 1 {
			return nil, fmt.Errorf("got List with more than one arg")
		}
		vType, err := TypeFromLF(builtin.GetArgs()[0], pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to determine type of TextMap: %w", err)
		}
		return TextMap{
			Value: vType,
		}, nil
	// We use fields above 1000 for dev features
	// Builtin type 'TBigNumeric'
	// *Available in versions >= 2.dev*
	case v2.BuiltinType_BIGNUMERIC:
		return BigNumeric{}, nil
	// Builtin type 'TRoundingMode'
	// *Available in versions >= 2.dev*
	case v2.BuiltinType_ROUNDING_MODE:
		return RoundingMode{}, nil
	default:
		return nil, fmt.Errorf("unknown builtin type: %s", builtin.Builtin.String())
	}
}

type Unit struct{ noImportOrDefinition }

func (Unit) GoType(ctx GeneratorContext) string {
	return "struct{}"
}

type Boolean struct{ noImportOrDefinition }

func (Boolean) GoType(ctx GeneratorContext) string {
	return "bool"
}

type Int64 struct{ noImportOrDefinition }

func (Int64) GoType(ctx GeneratorContext) string {
	return "int64"
}

type Date struct{}

func (Date) GoType(ctx GeneratorContext) string {
	return "stdtime.Time"
}

func (d Date) GoImports(ctx GeneratorContext) []string {
	return []string{"stdtime \"time\""}
}

func (d Date) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}

type Timestamp struct{}

func (Timestamp) GoType(ctx GeneratorContext) string {
	return "stdtime.Time"
}

func (t Timestamp) GoImports(ctx GeneratorContext) []string {
	return []string{"stdtime \"time\""}
}

func (t Timestamp) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}

type Numeric struct{ noImportOrDefinition }

func (Numeric) GoType(ctx GeneratorContext) string {
	return "string"
}

type Party struct{}

func (Party) GoType(ctx GeneratorContext) string {
	return "daml.Party"
}

func (p Party) GoImports(ctx GeneratorContext) []string {
	return []string{daml.PackageName}
}

func (p Party) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}

type Text struct{ noImportOrDefinition }

func (Text) GoType(ctx GeneratorContext) string {
	return "string"
}

type ContractId struct{}

func (ContractId) GoType(ctx GeneratorContext) string {
	return "daml.ContractId"
}

func (i ContractId) GoImports(ctx GeneratorContext) []string {
	return []string{daml.PackageName}
}

func (i ContractId) GoTypeDefinitions(ctx GeneratorContext) []string {
	return nil
}

type Optional struct {
	Value DamlType
}

func (o Optional) GoType(ctx GeneratorContext) string {
	vType := o.Value.GoType(ctx)
	return fmt.Sprintf("*%v", vType)
}

func (o Optional) GoImports(ctx GeneratorContext) []string {
	return o.Value.GoImports(ctx)
}

func (o Optional) GoTypeDefinitions(ctx GeneratorContext) []string {
	return o.Value.GoTypeDefinitions(ctx)
}

type List struct {
	Value DamlType
}

func (l List) GoType(ctx GeneratorContext) string {
	vType := l.Value.GoType(ctx)
	return fmt.Sprintf("[]%v", vType)
}

func (l List) GoImports(ctx GeneratorContext) []string {
	return l.Value.GoImports(ctx)
}

func (l List) GoTypeDefinitions(ctx GeneratorContext) []string {
	return l.Value.GoTypeDefinitions(ctx)
}

type GenMap struct {
	Key   DamlType
	Value DamlType
}

func (m GenMap) GoType(ctx GeneratorContext) string {
	kType := m.Key.GoType(ctx)
	vType := m.Value.GoType(ctx)
	return fmt.Sprintf("map[%v]%v", kType, vType)
}

func (m GenMap) GoImports(ctx GeneratorContext) []string {
	var imports []string
	imports = append(imports, m.Key.GoImports(ctx)...)
	imports = append(imports, m.Value.GoImports(ctx)...)
	return imports
}

func (m GenMap) GoTypeDefinitions(ctx GeneratorContext) []string {
	var typeDefs []string
	typeDefs = append(typeDefs, m.Key.GoTypeDefinitions(ctx)...)
	typeDefs = append(typeDefs, m.Value.GoTypeDefinitions(ctx)...)
	return typeDefs
}

type TextMap struct {
	Value DamlType
}

func (m TextMap) GoType(ctx GeneratorContext) string {
	vType := m.Value.GoType(ctx)
	return fmt.Sprintf("map[string]%v", vType)
}

func (m TextMap) GoImports(ctx GeneratorContext) []string {
	return m.Value.GoImports(ctx)
}

func (m TextMap) GoTypeDefinitions(ctx GeneratorContext) []string {
	return m.Value.GoTypeDefinitions(ctx)
}

type Any struct{ noImportOrDefinition }

func (a Any) GoType(ctx GeneratorContext) string {
	return "any"
}

type AnyException struct{ nonRenderable }

type TypeRep struct{ nonRenderable }

type Arrow struct{ nonRenderable }

type Update struct{ nonRenderable }

type FailureCategory struct{ nonRenderable }

type BigNumeric struct{ nonRenderable }

type RoundingMode struct{ nonRenderable }
