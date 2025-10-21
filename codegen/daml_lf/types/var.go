package types

import (
	"fmt"
	"strings"

	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

func VarFromLF(v *v2.Type_Var, pkg *v2.Package) (DamlType, error) {
	varName := pkg.GetInternedStrings()[v.GetVarInternedStr()]
	tArgs := make([]DamlType, len(v.GetArgs()))
	for i, t := range v.GetArgs() {
		tArg, err := TypeFromLF(t, pkg)
		if err != nil {
			return nil, err
		}
		tArgs[i] = tArg
	}
	return DamlVar{
		Name:          varName,
		TypeArguments: tArgs,
	}, nil
}

type DamlVar struct {
	Name          string
	TypeArguments []DamlType
}

func (v DamlVar) GoType(ctx GeneratorContext) string {
	tArgs := make([]string, len(v.TypeArguments))
	for i, arg := range v.TypeArguments {
		tArgs[i] = arg.GoType(ctx)
	}
	return fmt.Sprintf("%s[%s]", v.Name, strings.Join(tArgs, ","))
}

func (v DamlVar) GoImports(ctx GeneratorContext) []string {
	typeArgs := make([]string, 0, len(v.TypeArguments))
	for _, arg := range v.TypeArguments {
		typeArgs = append(typeArgs, arg.GoImports(ctx)...)
	}
	return typeArgs
}

func (v DamlVar) GoTypeDefinitions(ctx GeneratorContext) []string {
	typeDefs := make([]string, len(v.TypeArguments))
	for i, arg := range v.TypeArguments {
		typeDefs[i] = arg.GoType(ctx)
	}
	return typeDefs
}
