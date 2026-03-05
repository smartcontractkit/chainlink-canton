package bind

import (
	"github.com/smartcontractkit/go-daml/pkg/codec"
)

// EncodedChoice represents an encoded Daml choice invocation.
// It contains all the information needed to execute a choice on a contract.
type EncodedChoice struct {
	// TemplateID identifies the template (packageID:moduleName:templateName)
	TemplateID TemplateInformation
	// Choice is the name of the choice to exercise
	Choice string
	// OperationData is the hex-encoded choice parameters
	OperationData string
}

// TemplateInformation holds the components of a Daml template identifier.
type TemplateInformation struct {
	PackageID    string
	ModuleName   string
	TemplateName string
}

// BoundTemplate holds template metadata for encoding choice parameters.
type BoundTemplate struct {
	packageID    string
	moduleName   string
	templateName string
	hexCodec     *codec.HexCodec
}

// NewBoundTemplate creates a new BoundTemplate with the given identifiers.
func NewBoundTemplate(packageID, moduleName, templateName string) *BoundTemplate {
	return &BoundTemplate{
		packageID:    packageID,
		moduleName:   moduleName,
		templateName: templateName,
		hexCodec:     codec.NewHexCodec(),
	}
}

// EncodeChoiceArgs encodes choice parameters to hex and returns an EncodedChoice.
func (t *BoundTemplate) EncodeChoiceArgs(choice string, params any) (*EncodedChoice, error) {
	encoded, err := t.hexCodec.Marshal(params)
	if err != nil {
		return nil, err
	}
	return &EncodedChoice{
		TemplateID: TemplateInformation{
			PackageID:    t.packageID,
			ModuleName:   t.moduleName,
			TemplateName: t.templateName,
		},
		Choice:        choice,
		OperationData: encoded,
	}, nil
}

// PackageID returns the package ID of this template.
func (t *BoundTemplate) PackageID() string {
	return t.packageID
}

// ModuleName returns the module name of this template.
func (t *BoundTemplate) ModuleName() string {
	return t.moduleName
}

// TemplateName returns the template name of this template.
func (t *BoundTemplate) TemplateName() string {
	return t.templateName
}
