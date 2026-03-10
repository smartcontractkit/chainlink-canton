package contracts

import (
	"fmt"
	"strings"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
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

// TemplateIDFromBinding is a convenience function to get a TemplateID from a generated binding.
func TemplateIDFromBinding(template common.Template) TemplateID {
	split := strings.Split(template.GetTemplateID(), ":")

	return TemplateID{
		PackageID:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}
}

// TemplateID uniquely identifies a template.
type TemplateID struct {
	// The PackageID follow the same syntax as the protos, it can contain either the package's ID or name in '#package-name' syntax.
	PackageID  string `json:"packageId" toml:"package_id"`
	ModuleName string `json:"moduleName" toml:"module_name"`
	EntityName string `json:"entityName" toml:"entity_name"`
}

// ToLedgerIdentifier returns a ledgerv2.Identifier from the TemplateID.
func (t *TemplateID) ToLedgerIdentifier() *ledgerv2.Identifier {
	return &ledgerv2.Identifier{
		PackageId:  t.PackageID,
		ModuleName: t.ModuleName,
		EntityName: t.EntityName,
	}
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
