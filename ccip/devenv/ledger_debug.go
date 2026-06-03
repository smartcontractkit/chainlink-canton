package devenv

import (
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

// debugPrintSplitHoldingSubmitResponseJSON prints a ledger update/response as JSON (temporary debugging).
// Uses gogo jsonpb: dazl ledger protos implement github.com/gogo/protobuf/proto.Message, not google.golang.org/protobuf.
func debugPrintSplitHoldingSubmitResponseJSON(res proto.Message) {
	m := jsonpb.Marshaler{Indent: "  ", EmitDefaults: true}
	s, err := m.MarshalToString(res)
	if err != nil {
		fmt.Printf("debugPrintSplitHoldingSubmitResponseJSON: %v\n", err)
		return
	}
	fmt.Printf("debugPrintSplitHoldingSubmitResponseJSON:\n%s\n", s)
}
