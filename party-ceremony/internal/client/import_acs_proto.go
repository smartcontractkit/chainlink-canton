package client

import (
	"fmt"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"google.golang.org/protobuf/encoding/protowire"
)

// importAcsSynchronizerIDField is synchronizer_id on ImportAcsRequest in Canton ≥3.4.
// dazl-client v8.9.0 bindings predate this field; wire it until generated protos catch up.
const importAcsSynchronizerIDField protowire.Number = 6

func setImportAcsFirstChunkMetadata(req *participantv30.ImportAcsRequest, synchronizerID string) error {
	if synchronizerID == "" {
		return fmt.Errorf("synchronizer_id is required on the first ImportAcs request chunk")
	}

	req.ContractImportMode = participantv30.ContractImportMode_CONTRACT_IMPORT_MODE_ACCEPT

	wire := protowire.AppendString(
		protowire.AppendTag(nil, importAcsSynchronizerIDField, protowire.BytesType),
		synchronizerID,
	)
	existing := req.ProtoReflect().GetUnknown()
	req.ProtoReflect().SetUnknown(append(append([]byte(nil), existing...), wire...))

	return nil
}

// importAcsSynchronizerID reads synchronizer_id from a decoded ImportAcsRequest, including
// when the value was set via unknown-field wire encoding.
func importAcsSynchronizerID(req *participantv30.ImportAcsRequest) (string, bool) {
	if req == nil {
		return "", false
	}

	b := req.ProtoReflect().GetUnknown()
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return "", false
		}
		b = b[n:]
		if num != importAcsSynchronizerIDField || typ != protowire.BytesType {
			n = protowire.ConsumeFieldValue(num, typ, b)
			if n < 0 {
				return "", false
			}
			b = b[n:]

			continue
		}
		v, n := protowire.ConsumeString(b)
		if n < 0 {
			return "", false
		}

		return v, true
	}

	return "", false
}
