package store

import (
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

var (
	_ ActiveContractStoreInterface    = (*ActiveContractStore)(nil)
	_ InstrumentHoldingStoreInterface = (*InstrumentHoldingStore)(nil)
)

// ActiveContractStoreInterface provides read access to active contracts and the ability to register templates.
type ActiveContractStoreInterface interface {
	// Get returns the active contract for the given instance address.
	Get(address contracts.InstanceAddress) (*apiv2.ActiveContract, bool)
	// GetByTemplateId returns the active contract for the given party and template ID.
	GetByTemplateId(party types.PARTY, templateId contracts.TemplateID) (*apiv2.ActiveContract, bool)
	// GetByContractId returns the active contract for the given contract ID.
	GetByContractId(contractId types.CONTRACT_ID) (*apiv2.ActiveContract, bool)
	// RegisterTemplates registers templates to be tracked by the store.
	RegisterTemplates(templates ...RegisteredTemplate)
}

// InstrumentHoldingStoreInterface provides read access to instrument holdings and the ability to register parties.
type InstrumentHoldingStoreInterface interface {
	// GetHolding returns the active holdings for the given party and instrument ID.
	GetHolding(party types.PARTY, instrumentId splice_api_token_holding_v1.InstrumentId) ([]*apiv2.ActiveContract, bool)
	// RegisterParty registers parties whose holdings should be tracked by the store.
	RegisterParty(parties ...string)
}
