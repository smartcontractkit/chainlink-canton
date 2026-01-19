package contracts

import (
	"embed"
	"path/filepath"
)

//go:embed coin dependencies mcms multi-package.yaml
var Embed embed.FS

type Package string

const (
	Coin = Package("coin")

	MCMS     = Package("mcms")
	MCMSTest = Package("mcms_test")
)

var Contracts map[Package]string = map[Package]string{
	Coin:     filepath.Join("coin"),
	MCMS:     filepath.Join("mcms"),
	MCMSTest: filepath.Join("mcms", "test"),
}
