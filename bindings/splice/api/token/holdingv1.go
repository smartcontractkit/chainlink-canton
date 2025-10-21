// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

import (
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/da/time"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
	stdtime "time"
)

type HoldingView struct {
	Amount       string
	Instrumentid InstrumentId
	Lock         *Lock
	Meta         Metadata
	Owner        daml.Party
}

type Lock struct {
	Context      *string
	Expiresafter *time.RelTime
	Expiresat    *stdtime.Time
	Holders      []daml.Party
}

type InstrumentId struct {
	Admin daml.Party
	Id    string
}
