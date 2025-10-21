// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

import (
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
	stdtime "time"
)

type TransferFactoryView struct {
	Admin daml.Party
	Meta  Metadata
}

type TransferFactory_PublicFetch struct {
	Actor         daml.Party
	Expectedadmin daml.Party
}

type TransferFactory_Transfer struct {
	Expectedadmin daml.Party
	Extraargs     ExtraArgs
	Transfer      Transfer
}

type TransferInstruction_Update struct {
	Extraactors []daml.Party
	Extraargs   ExtraArgs
}

type TransferInstruction_Withdraw struct {
	Extraargs ExtraArgs
}

type TransferInstruction_Reject struct {
	Extraargs ExtraArgs
}

type TransferInstruction_Accept struct {
	Extraargs ExtraArgs
}

type TransferInstructionView struct {
	Meta                   Metadata
	Originalinstructioncid *daml.ContractId
	Status                 TransferInstructionStatus
	Transfer               Transfer
}

// Variant TransferInstructionStatus
// Types that are valid to be assigned to TransferInstructionStatus:
//
//	TransferPendingInternalWorkflow
//	TransferPendingReceiverAcceptance
type TransferInstructionStatus interface {
	_isTransferInstructionStatus()
}

type TransferPendingInternalWorkflow struct {
	Pendingactions map[daml.Party]string
}

func (v TransferPendingInternalWorkflow) _isTransferInstructionStatus() {}

type TransferPendingReceiverAcceptance struct{}

func (v TransferPendingReceiverAcceptance) _isTransferInstructionStatus() {}

// Variant TransferInstructionResult_Output
// Types that are valid to be assigned to TransferInstructionResult_Output:
//
//	TransferInstructionResult_Completed
//	TransferInstructionResult_Failed
//	TransferInstructionResult_Pending
type TransferInstructionResult_Output interface {
	_isTransferInstructionResult_Output()
}

type TransferInstructionResult_Completed struct {
	Receiverholdingcids []daml.ContractId
}

func (v TransferInstructionResult_Completed) _isTransferInstructionResult_Output() {}

type TransferInstructionResult_Failed struct{}

func (v TransferInstructionResult_Failed) _isTransferInstructionResult_Output() {}

type TransferInstructionResult_Pending struct {
	Transferinstructioncid daml.ContractId
}

func (v TransferInstructionResult_Pending) _isTransferInstructionResult_Output() {}

type TransferInstructionResult struct {
	Meta             Metadata
	Output           TransferInstructionResult_Output
	Senderchangecids []daml.ContractId
}

type Transfer struct {
	Amount           string
	Executebefore    stdtime.Time
	Inputholdingcids []daml.ContractId
	Instrumentid     InstrumentId
	Meta             Metadata
	Receiver         daml.Party
	Requestedat      stdtime.Time
	Sender           daml.Party
}
