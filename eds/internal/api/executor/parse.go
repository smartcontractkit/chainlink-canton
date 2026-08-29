package executor

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/parse"
)

type Executor struct {
	Address             contracts.RawInstanceAddress
	MaxCCVsPerMessage   int64
	CCVAllowlistEnabled bool
	AllowedCCVs         []contracts.RawInstanceAddress
}

func ParseExecutor(createdEvent *apiv2.CreatedEvent) (*Executor, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[executor.Executor](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal executor: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.Owner)

	allowedCCVs, err := parse.RawInstanceAddressList(boundContract.AllowedCCVs)
	if err != nil {
		return nil, fmt.Errorf("invalid allowed CCVs: %w", err)
	}

	return &Executor{
		Address:             address,
		MaxCCVsPerMessage:   int64(boundContract.MaxCCVsPerMsg),
		CCVAllowlistEnabled: bool(boundContract.DynamicConfig.CcvAllowlistEnabled),
		AllowedCCVs:         allowedCCVs,
	}, nil
}
