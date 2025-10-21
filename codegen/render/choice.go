package render

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf/types"
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type ChoiceData struct {
	// Name of the template and of the type constructor.
	// Check Module.DataTypes to retrieve the type constructor.
	Name string

	Consuming bool
	Argument  types.DamlType
	RetType   types.DamlType
}

func ChoiceFromLF(choice *v2.TemplateChoice, pkg *v2.Package) (*ChoiceData, error) {
	choiceName := pkg.GetInternedStrings()[choice.GetNameInternedStr()]

	arg, err := types.TypeFromLF(choice.GetArgBinder().GetType(), pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to get choice argument type: %w", err)
	}
	retType, err := types.TypeFromLF(choice.GetRetType(), pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to get choice return type: %w", err)
	}

	return &ChoiceData{
		Name:      choiceName,
		Consuming: choice.Consuming,
		Argument:  arg,
		RetType:   retType,
	}, nil
}
