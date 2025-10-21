package tmpl

import (
	_ "embed"
)

//go:embed module.tmpl
var EmbeddedTemplate string
