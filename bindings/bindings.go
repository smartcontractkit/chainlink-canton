package bindings

import (
	"errors"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/service/ledger"
)

// UnmarshalActiveContract unmarshals a Canton ActiveContract response into a typed DAML binding struct.
// This provides type-safe access to contract fields instead of manual string-based field parsing.
//
// Example usage:
//
//	mcmsContract, err := UnmarshalActiveContract[mcms.MCMS](activeContract)
//	if err != nil {
//	    return err
//	}
//	// Now use type-safe fields: mcmsContract.McmsId, mcmsContract.Config.Signers, etc.
func UnmarshalActiveContract[T any](ac *apiv2.GetActiveContractsResponse_ActiveContract) (*T, error) {
	if ac == nil {
		return nil, errors.New("active contract is nil")
	}

	createdEvent := ac.ActiveContract.GetCreatedEvent()
	if createdEvent == nil {
		return nil, errors.New("no created event in active contract")
	}

	return UnmarshalCreatedEvent[T](createdEvent)
}

// UnmarshalCreatedEvent unmarshals a CreatedEvent into a typed DAML binding struct.
// This is useful when you have a CreatedEvent directly (e.g., from transaction events).
func UnmarshalCreatedEvent[T any](event *apiv2.CreatedEvent) (*T, error) {
	if event == nil {
		return nil, errors.New("created event is nil")
	}

	createArgs := event.GetCreateArguments()
	if createArgs == nil {
		return nil, errors.New("no create arguments in created event")
	}

	// Use go-daml's RecordToStruct function
	var result T
	err := ledger.RecordToStruct(createArgs, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert record to type %T: %w", result, err)
	}

	return &result, nil
}
