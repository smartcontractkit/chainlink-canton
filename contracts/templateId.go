package contracts

import (
	"fmt"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

// TODO these shouldn't be needed

func StripPackageNameFromTemplateID(templateID string) string {
	templateID = strings.TrimPrefix(templateID, "#") // TODO: this shouldn't have been added in the bindings, remove
	parts := strings.Split(templateID, ":")
	if len(parts) != 3 {
		return templateID
	}

	return parts[1] + ":" + parts[2]
}

func ReplacePackageIdWithNameInTemplateID(templateID, packageName string) string {
	templateID = strings.TrimPrefix(templateID, "#") // TODO: this shouldn't have been added in the bindings, remove
	parts := strings.Split(templateID, ":")
	if len(parts) != 3 {
		return templateID
	}

	return packageName + ":" + parts[1] + ":" + parts[2]
}

type templateBinding interface {
	GetTemplateID() string
}

// TemplateIDFromBinding is a convenience function to get a TemplateID from a generated binding.
func TemplateIDFromBinding(template templateBinding) TemplateID {
	templateId, err := TemplateIDFromString(template.GetTemplateID())
	if err != nil {
		panic(err)
	}

	return templateId
}

func TemplateIDFromString(s string) (TemplateID, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return TemplateID{}, fmt.Errorf("invalid template id format: %s", s)
	}

	return TemplateID{
		PackageID:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

// TemplateID uniquely identifies a template.
type TemplateID struct {
	// The PackageID follow the same syntax as the protos, it can contain either the package's ID or name in '#package-name' syntax.
	PackageID  string `json:"packageId" toml:"package_id"`
	ModuleName string `json:"moduleName" toml:"module_name"`
	EntityName string `json:"entityName" toml:"entity_name"`
}

// ToLedgerIdentifier returns a ledgerv2.Identifier from the TemplateID.
func (t TemplateID) ToLedgerIdentifier() *apiv2.Identifier {
	return &apiv2.Identifier{
		PackageId:  t.PackageID,
		ModuleName: t.ModuleName,
		EntityName: t.EntityName,
	}
}

func (t TemplateID) String() string {
	return fmt.Sprintf("%s:%s:%s", t.PackageID, t.ModuleName, t.EntityName)
}

// ParseTemplateIDFromString parses a template ID string like "#package:Module:Entity" into its components
func ParseTemplateIDFromString(templateID string) (packageID, moduleName, entityName string, err error) {
	if !strings.HasPrefix(templateID, "#") {
		return "", "", "", fmt.Errorf("template ID must start with #")
	}
	parts := strings.Split(templateID, ":")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("template ID must have format #package:module:entity, got: %s", templateID)
	}

	return parts[0], parts[1], parts[2], nil
}
