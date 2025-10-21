package types

import (
	"strings"

	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

func GetDottedName(idx int32, pkg *v2.Package) (string, error) {
	internedSegments := pkg.GetInternedDottedNames()[idx].GetSegmentsInternedStr()
	internedStrings := make([]string, len(internedSegments))
	for i2, segment := range internedSegments {
		internedStrings[i2] = pkg.GetInternedStrings()[segment]
	}
	return strings.Join(internedStrings, "."), nil
}
