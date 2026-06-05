package registry

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// ContractRef is a template instance visible to a party.
type ContractRef struct {
	ContractID string
	TemplateID string
}

// FindContractsByEntity lists active contract IDs for a template entity name.
func FindContractsByEntity(ctx context.Context, client ledger.Client, party string, template any, entityName string) ([]ContractRef, error) {
	tplID := contracts.IdentifierFromBinding(template.(interface{ GetTemplateID() string }))
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.ForParty(party), tplID)
	if err != nil {
		return nil, fmt.Errorf("list %s contracts: %w", entityName, err)
	}

	refs := make([]ContractRef, 0, len(active))
	for _, ac := range active {
		ce := ac.GetCreatedEvent()
		if ce == nil {
			continue
		}
		if entityName != "" && ce.GetTemplateId().GetEntityName() != entityName {
			continue
		}
		refs = append(refs, ContractRef{
			ContractID: ce.GetContractId(),
			TemplateID: fmt.Sprintf("%s:%s:%s",
				ce.GetTemplateId().GetPackageId(),
				ce.GetTemplateId().GetModuleName(),
				ce.GetTemplateId().GetEntityName(),
			),
		})
	}

	return refs, nil
}

// FindFirstContractByEntity returns the first ACS match or empty string.
func FindFirstContractByEntity(ctx context.Context, client ledger.Client, party string, template any, entityName string) (string, error) {
	refs, err := FindContractsByEntity(ctx, client, party, template, entityName)
	if err != nil {
		return "", err
	}
	if len(refs) == 0 {
		return "", nil
	}
	return refs[0].ContractID, nil
}
