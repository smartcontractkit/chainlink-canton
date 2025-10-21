// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

import (
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
)

type BurnMintFactoryView struct {
	Admin daml.Party
	Meta  Metadata
}

type BurnMintFactory_BurnMintResult struct {
	Outputcids []daml.ContractId
}

type BurnMintOutput struct {
	Amount  string
	Context ChoiceContext
	Owner   daml.Party
}

type BurnMintFactory_PublicFetch struct {
	Actor         daml.Party
	Expectedadmin daml.Party
}

type BurnMintFactory_BurnMint struct {
	Expectedadmin    daml.Party
	Extraactors      []daml.Party
	Extraargs        ExtraArgs
	Inputholdingcids []daml.ContractId
	Instrumentid     InstrumentId
	Outputs          []BurnMintOutput
	Sender           daml.Party
}
