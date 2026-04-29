package token_standard

import (
	"encoding/base64"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

func ActiveContractToDisclosedContract(activeContract *apiv2.ActiveContract, includeDebugFields bool) oapiTransferInstruction.DisclosedContract {
	contract := oapiTransferInstruction.DisclosedContract{
		TemplateId:       fmt.Sprintf("%s:%s:%s", activeContract.GetCreatedEvent().GetTemplateId().GetPackageId(), activeContract.GetCreatedEvent().GetTemplateId().GetModuleName(), activeContract.GetCreatedEvent().GetTemplateId().GetEntityName()),
		ContractId:       activeContract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: base64.StdEncoding.EncodeToString(activeContract.GetCreatedEvent().GetCreatedEventBlob()),
		SynchronizerId:   activeContract.GetSynchronizerId(),
		DebugPayload:     nil,
	}
	if includeDebugFields {
		contract.DebugCreatedAt = new(activeContract.GetCreatedEvent().GetCreatedAt().AsTime())
		contract.DebugPackageName = new(activeContract.GetCreatedEvent().GetPackageName())
		_ = ledger.RecordToStruct(activeContract.GetCreatedEvent().GetCreateArguments(), &contract.DebugPayload)
	}

	return contract
}
