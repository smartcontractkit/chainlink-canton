package tmpl

type ModuleData struct {
	ModuleName string
	Imports    []string

	TypeDefinitions []string

	Templates []TemplateData
}
