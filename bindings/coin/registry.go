// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coin

import (
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/splice/api/token"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
)

type CoinRegistry struct {
	Instrumentid token.InstrumentId
	Issuer       daml.Party
	Meta         token.Metadata
}

type MintPreapproval struct {
	Receiver daml.Party
	Sender   daml.Party
}

type MintRole struct {
	Issuer   daml.Party
	Minter   daml.Party
	Registry daml.ContractId
}
