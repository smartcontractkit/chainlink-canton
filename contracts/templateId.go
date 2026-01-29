package contracts

import (
	"strings"
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
