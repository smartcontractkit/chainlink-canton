package daml_lf

import (
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive"
)

type DamlLfArchive struct {
	Name    string
	Payload DamlLfArchivePayload
	Hash    string
}

func (a *DamlLfArchive) ParseNamed(name string, b []byte) error {
	arc := archive.Archive{}
	if err := proto.Unmarshal(b, &arc); err != nil {
		return err
	}

	a.Name = name
	a.Hash = arc.GetHash()

	if err := a.Payload.FromBytes(arc.GetPayload()); err != nil {
		return err
	}

	return nil
}
