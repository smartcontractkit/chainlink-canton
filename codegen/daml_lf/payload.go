package daml_lf

import (
	"errors"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive"
	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v1"
	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

type DamlLfArchivePayload struct {
	LanguageVersion LanguageVersion
	Package         DamlLfPackage
}

type DamlLfPackage struct {
	// TODO use an interface here?
	V1 *v1.Package
	V2 *v2.Package
}

func (p *DamlLfArchivePayload) FromBytes(data []byte) error {
	payload := &archive.ArchivePayload{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return err
	}
	p.LanguageVersion.Minor = payload.Minor

	switch payload.Sum.(type) {
	case *archive.ArchivePayload_DamlLf_1:
		p.Package.V1 = &v1.Package{}
		if err := proto.Unmarshal(payload.GetDamlLf_1(), p.Package.V1); err != nil {
			return err
		}
		p.LanguageVersion.Major = 1
	case *archive.ArchivePayload_DamlLf_2:
		p.Package.V2 = &v2.Package{}
		if err := proto.Unmarshal(payload.GetDamlLf_2(), p.Package.V2); err != nil {
			return err
		}
		p.LanguageVersion.Major = 2
	default:
		return errors.New("invalid payload type")
	}
	return nil
}
