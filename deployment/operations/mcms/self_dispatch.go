package mcms

import (
	"encoding/hex"
	"fmt"

	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
)

// SelfDispatchSetConfigFunction is the self-dispatch function name for SetConfig.
const SelfDispatchSetConfigFunction = "SetConfig"

// Matches MCMS.Codec.roleToInt. Role.MarshalHex encodes the role as text, so it can't be used.
var selfDispatchRoleByte = map[mcmsApi.Role]byte{
	mcmsApi.RoleProposer:  0,
	mcmsApi.RoleCanceller: 1,
	mcmsApi.RoleBypasser:  2,
}

// EncodeSelfDispatchSetConfig encodes uint8(role) || SetConfigParams, like MCMS.Codec.encodeSelfDispatchSetConfig.
func EncodeSelfDispatchSetConfig(role mcmsApi.Role, params mcmsApi.SetConfigParams) (string, error) {
	roleByte, ok := selfDispatchRoleByte[role]
	if !ok {
		return "", fmt.Errorf("unknown MCMS role %q", role)
	}

	wire, err := params.MarshalHex()
	if err != nil {
		return "", fmt.Errorf("encode SetConfig params: %w", err)
	}

	return hex.EncodeToString(append([]byte{roleByte}, wire...)), nil
}
