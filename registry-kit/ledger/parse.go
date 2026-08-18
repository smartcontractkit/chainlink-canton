package ledger

import (
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
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
