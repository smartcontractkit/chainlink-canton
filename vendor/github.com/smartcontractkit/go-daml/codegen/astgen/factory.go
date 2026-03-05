package astgen

import (
	"fmt"

	"github.com/smartcontractkit/go-daml/codegen/astgen/v3"
	model2 "github.com/smartcontractkit/go-daml/codegen/model"
)

const (
	V1 = "1."
	V2 = "2."
	V3 = "3."
)

type AstGen interface {
	GetInterfaces() (map[string]*model2.TmplStruct, error)
	GetTemplateStructs(ifcByModule map[string]model2.InterfaceMap) (map[string]*model2.TmplStruct, model2.ExternalPackages, error)
}

func GetAstGenFromVersion(payload []byte, ext model2.ExternalPackages, ver string) (AstGen, error) {
	switch ver {
	case V3:
		return v3.NewCodegenAst(payload, ext), nil
	default:
		return nil, fmt.Errorf("none supported version")
	}
}
