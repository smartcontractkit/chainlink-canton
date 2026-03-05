package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/smartcontractkit/go-daml/codegen/model"
)

type tmplData struct {
	Package          string
	PackageID        string
	PackageName      string
	SdkVersion       string
	Structs          map[string]*model.TmplStruct
	IsMainDalf       bool
	GenerateHexCodec bool
	ChoiceArgTypes   map[string]bool // Types used as choice arguments (for Encode functions)
	ImportedPackages []model.ExternalPackage
}

//go:embed source.go.tmpl
var tmplSource string

func Bind(genPkg string, pkg *model.Package, sdkVersion string, isMainDalf bool, generateHexCodec bool) (string, error) {
	// Collect all types used as choice arguments
	choiceArgTypes := make(map[string]bool)
	for _, tmpl := range pkg.Structs {
		if tmpl.IsTemplate {
			for _, choice := range tmpl.Choices {
				if choice.ArgType != nil && choice.ArgType.GoType() != "" && choice.ArgType.GoType() != "UNIT" {
					argType := choice.ArgType.GoType()
					// Special case: SET type uses the choice name as the struct name
					// This matches the transformation done in the template for choice methods
					if argType == "SET" {
						argType = capitalize(choice.Name)
					}
					choiceArgTypes[argType] = true
				}
			}
		}
	}

	data := &tmplData{
		Package:          genPkg,
		PackageID:        pkg.PackageID,
		PackageName:      pkg.Name,
		SdkVersion:       sdkVersion,
		Structs:          pkg.Structs,
		IsMainDalf:       isMainDalf,
		GenerateHexCodec: generateHexCodec,
		ChoiceArgTypes:   choiceArgTypes,
		ImportedPackages: pkg.ImportedPackages,
	}
	buffer := new(bytes.Buffer)

	funcs := map[string]interface{}{
		"capitalise":        capitalize,
		"decapitalize":      decapitalize,
		"stringsHasPrefix":  strings.HasPrefix,
		"stringsTrimPrefix": strings.TrimPrefix,
		"stringsHasSuffix":  strings.HasSuffix,
		"stringsTrimSuffix": strings.TrimSuffix,
	}
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(tmplSource))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", err
	}
	// Pass the code through gofmt to clean it up
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, buffer)
	}
	return string(code), nil
}

func capitalize(input string) string {
	if len(input) == 0 {
		return input
	}

	hasSeparators := strings.ContainsAny(input, "_- ")

	if !hasSeparators && len(input) > 0 && input[0] >= 'A' && input[0] <= 'Z' {
		return input
	}

	result := toCamelCase(input)
	return strings.ToUpper(result[:1]) + result[1:]
}

func decapitalize(input string) string {
	if len(input) == 0 {
		return input
	}

	if isAllCaps(input) {
		return strings.ToLower(input)
	}

	if len(input) > 0 && input[0] >= 'a' && input[0] <= 'z' && !strings.ContainsAny(input, "_- ") {
		return input
	}

	result := toCamelCase(input)
	return strings.ToLower(result[:1]) + result[1:]
}

func toCamelCase(input string) string {
	if len(input) == 0 {
		return input
	}

	if !strings.ContainsAny(input, "_- ") {
		return input
	}

	words := strings.FieldsFunc(input, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	if len(words) == 0 {
		return input
	}

	var result strings.Builder
	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		if isAllCaps(word) {
			if len(word) <= 3 {
				result.WriteString(word)
			} else {
				if i == 0 {
					result.WriteString(strings.ToLower(word))
				} else {
					result.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
				}
			}
		} else {
			if i == 0 {
				result.WriteString(strings.ToLower(word[:1]) + word[1:])
			} else {
				result.WriteString(strings.ToUpper(word[:1]) + word[1:])
			}
		}
	}

	return result.String()
}

func isAllCaps(input string) bool {
	if len(input) == 0 {
		return false
	}
	for _, r := range input {
		if r >= 'a' && r <= 'z' {
			return false
		}
	}
	for _, r := range input {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}
