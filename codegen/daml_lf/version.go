package daml_lf

import (
	"fmt"
)

type LanguageVersion struct {
	Major int32
	Minor string
}

func (v LanguageVersion) String() string {
	return fmt.Sprintf("v%d_%s", v.Major, v.Minor)
}
