package tokenpool

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
)

type preapprovalFactory func(ctx context.Context) (string, *apiv2.ActiveContract, error)

func getPreapprovalFactory(acs store.ActiveContractStoreInterface, contextKey string, party types.PARTY, cfg config.TransferPreapproval) (preapprovalFactory, error) {
	templateId, err := contracts.TemplateIDFromString(cfg.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("invalid TemplateId for TransferPreapproval: %w", err)
	}
	acs.RegisterTemplates(store.RegisteredTemplate{
		TemplateID: templateId,
		PartyID:    string(party),
	})

	return func(ctx context.Context) (string, *apiv2.ActiveContract, error) {
		activePreapproval, ok := acs.GetByTemplateId(party, templateId)
		if !ok {
			return "", nil, fmt.Errorf("no preapproval found for user %s and template %s", party, templateId.String())
		}

		return contextKey, activePreapproval, nil
	}, nil
}
