package types

import (
	v2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/digitalasset/daml/lf/archive/v2"
)

func InternedFromLF(interned *v2.Type_InternedType, pkg *v2.Package) (DamlType, error) {
	it := pkg.GetInternedTypes()[interned.InternedType]
	return TypeFromLF(it, pkg)
}
