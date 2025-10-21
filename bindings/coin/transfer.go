// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coin

import (
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
	stdtime "time"
)

type CoinTransferInstruction struct {
	Executebefore stdtime.Time
	Holding       CoinHolding
	Newowner      daml.Party
	Requestedat   stdtime.Time
}
