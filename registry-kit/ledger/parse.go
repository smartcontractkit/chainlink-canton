package ledger

import (
	"fmt"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	extensionapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/extensionapi"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
)

// CreatedContractID returns the contract ID of the first created event matching entityName.
func CreatedContractID(tx *apiv2.Transaction, entityName string) (string, bool) {
	if tx == nil {
		return "", false
	}
	for _, event := range tx.GetEvents() {
		if c := event.GetCreated(); c != nil && c.GetTemplateId().GetEntityName() == entityName {
			return c.GetContractId(), true
		}
	}

	return "", false
}

// CreatedHoldingForOwner returns a created Holding contract ID for the given owner party.
func CreatedHoldingForOwner(tx *apiv2.Transaction, owner string) (string, bool) {
	if tx == nil {
		return "", false
	}
	for _, event := range tx.GetEvents() {
		c := event.GetCreated()
		if c == nil || c.GetTemplateId().GetEntityName() != "Holding" {
			continue
		}
		for _, field := range c.GetCreateArguments().GetFields() {
			if field.GetLabel() != "owner" {
				continue
			}
			if field.GetValue().GetParty() == owner {
				return c.GetContractId(), true
			}
		}
	}

	return "", false
}

// ParseReleaseOrMintResult extracts ReleaseFromTicket exercise output when present.
func ParseReleaseOrMintResult(tx *apiv2.Transaction) (extensionapi.ReleaseOrMintResultOutput, error) {
	if tx == nil {
		return extensionapi.ReleaseOrMintResultOutput{}, fmt.Errorf("transaction is nil")
	}
	for _, event := range tx.GetEvents() {
		ex := event.GetExercised()
		if ex == nil || !strings.Contains(ex.GetChoice(), "ReleaseFromTicket") {
			continue
		}
		var result extensionapi.ReleaseOrMintResult
		if err := ledger.RecordToStruct(ex.GetExerciseResult(), &result); err != nil {
			return extensionapi.ReleaseOrMintResultOutput{}, fmt.Errorf("parse ReleaseFromTicket result: %w", err)
		}

		return result.Output, nil
	}

	return extensionapi.ReleaseOrMintResultOutput{}, fmt.Errorf("ReleaseFromTicket not found in transaction")
}

// ParseLockOrBurnResult extracts LockOrBurn exercise output when present.
func ParseLockOrBurnResult(tx *apiv2.Transaction) (extensionapi.LockOrBurnResult, error) {
	if tx == nil {
		return extensionapi.LockOrBurnResult{}, fmt.Errorf("transaction is nil")
	}
	for _, event := range tx.GetEvents() {
		ex := event.GetExercised()
		if ex == nil || !strings.Contains(ex.GetChoice(), "LockOrBurn") {
			continue
		}
		var result extensionapi.LockOrBurnResult
		if err := ledger.RecordToStruct(ex.GetExerciseResult(), &result); err != nil {
			return extensionapi.LockOrBurnResult{}, fmt.Errorf("parse LockOrBurn result: %w", err)
		}

		return result, nil
	}

	return extensionapi.LockOrBurnResult{}, fmt.Errorf("LockOrBurn not found in transaction")
}

// CreatedHoldingsForOwner returns all Holding contract IDs created for owner in the transaction.
func CreatedHoldingsForOwner(tx *apiv2.Transaction, owner string) []string {
	if tx == nil {
		return nil
	}
	var cids []string
	for _, event := range tx.GetEvents() {
		c := event.GetCreated()
		if c == nil || c.GetTemplateId().GetEntityName() != "Holding" {
			continue
		}
		for _, field := range c.GetCreateArguments().GetFields() {
			if field.GetLabel() == "owner" && field.GetValue().GetParty() == owner {
				cids = append(cids, c.GetContractId())
				break
			}
		}
	}
	return cids
}
