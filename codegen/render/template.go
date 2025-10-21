package render

import (
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf/types"
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type TemplateData struct {
	// Name of the template and of the type constructor.
	// Check Module.DataTypes to retrieve the type constructor.
	Name string

	Choices []*ChoiceData
}

func TemplateFromLF(template *v2.DefTemplate, pkg *v2.Package) (*TemplateData, error) {
	templateName, err := types.GetDottedName(template.GetTyconInternedDname(), pkg)
	if err != nil {
		return nil, err
	}

	choices := make([]*ChoiceData, len(template.GetChoices()))
	for i, choice := range template.GetChoices() {
		choices[i], err = ChoiceFromLF(choice, pkg)
		if err != nil {
			return nil, err
		}
	}

	return &TemplateData{
		Name:    templateName,
		Choices: choices,
	}, nil
}
