package executor

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type Executor struct {
	Address contracts.RawInstanceAddress
}

func ParseExecutor(createdEvent *apiv2.CreatedEvent) (*Executor, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[executor.Executor](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal executor: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.Owner)

	return &Executor{
		Address: address,
	}, nil
}
