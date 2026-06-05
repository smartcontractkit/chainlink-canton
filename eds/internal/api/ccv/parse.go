package ccv

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type CommitteeVerifier struct {
	Address contracts.RawInstanceAddress
}

func ParseCommitteeVerifier(createdEvent *apiv2.CreatedEvent) (*CommitteeVerifier, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[committeeverifier.CommitteeVerifier](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal committee verifier: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.Owner)

	return &CommitteeVerifier{
		Address: address,
	}, nil
}
