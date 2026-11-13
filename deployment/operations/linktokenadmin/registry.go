package linktokenadmin

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	spliceMeta "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/daRegistry"
)

// RegistryInput says where the registry factory and its choice context come from.
// Set URL to ask a DA Registry utility backend, or set FactoryCid and ChoiceContext for a
// registry without a backend.
type RegistryInput struct {
	URL string `json:"url,omitempty"`

	FactoryCid    string                   `json:"factoryCid,omitempty"`
	ChoiceContext spliceMeta.ChoiceContext `json:"choiceContext"`
	// DiscloseContractIDs are contracts the choice context references, such as the instrument
	// configuration. They and the factory are disclosed from the participant's view.
	DiscloseContractIDs []string `json:"discloseContractIds,omitempty"`
	// Disclosed carries the created event blobs of the factory and the contracts the choice
	// context references, in the DA Registry's disclosedContracts shape. Use it when the
	// participant's parties cannot see those contracts: they are passed as given, and nothing is
	// read from the participant.
	Disclosed []daRegistry.DisclosedContract `json:"disclosed,omitempty"`
}

func (r RegistryInput) validate() error {
	switch {
	case r.URL != "" && r.FactoryCid != "":
		return errors.New("registry: set either URL or FactoryCid, not both")
	case len(r.Disclosed) > 0 && (r.URL != "" || len(r.DiscloseContractIDs) > 0):
		return errors.New("registry: Disclosed goes with FactoryCid only, without URL or DiscloseContractIDs")
	case r.URL == "" && r.FactoryCid == "":
		return errors.New("registry: URL or FactoryCid is required")
	}

	return nil
}

// registryFactory is the factory an Execute choice exercises, with what the submitter must pass
// and disclose for it.
type registryFactory struct {
	factoryCid string
	context    spliceMeta.ChoiceContext
	// disclosed holds what the registry backend returned; static inputs are resolved by the caller.
	disclosed []*apiv2.DisclosedContract
}

// burnMintRequest is the DA Registry body for an ExecuteMint (one output, no inputs) or an
// ExecuteBurn (inputs, no outputs). The API rejects null arrays, so both are always set.
func burnMintRequest(kind Kind, auth authorization, inputHoldingCids []string) daRegistry.GetBurnMintFactoryRequest {
	req := daRegistry.GetBurnMintFactoryRequest{
		InstrumentId:     daRegistry.InstrumentId{Admin: string(auth.instrumentID.Admin), Id: string(auth.instrumentID.Id)},
		InputHoldingCids: []string{},
		Outputs:          []daRegistry.MintOutput{},
	}
	if kind == KindMint {
		req.Outputs = append(req.Outputs, daRegistry.MintOutput{Owner: auth.counterparty, Amount: string(auth.amount)})
	} else {
		req.InputHoldingCids = append(req.InputHoldingCids, inputHoldingCids...)
	}

	return req
}

func fetchBurnMintFactory(ctx context.Context, client daRegistry.ClientWithResponsesInterface, req daRegistry.GetBurnMintFactoryRequest) (registryFactory, error) {
	resp, err := client.GetBurnMintFactoryWithResponse(ctx, req)
	if err != nil {
		return registryFactory{}, fmt.Errorf("get burn-mint factory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return registryFactory{}, fmt.Errorf("get burn-mint factory: status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return registryFactory{}, fmt.Errorf("decode choice context: %w", err)
	}
	disclosed := make([]*apiv2.DisclosedContract, 0, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for _, dc := range resp.JSON200.ChoiceContext.DisclosedContracts {
		d, err := registryDisclosure(dc)
		if err != nil {
			return registryFactory{}, err
		}
		disclosed = append(disclosed, d)
	}

	return registryFactory{factoryCid: resp.JSON200.FactoryId, context: choiceContext, disclosed: disclosed}, nil
}

func registryDisclosure(dc daRegistry.DisclosedContract) (*apiv2.DisclosedContract, error) {
	templateID, err := contracts.TemplateIDFromString(dc.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("registry disclosure %s: %w", dc.ContractId, err)
	}
	blob, err := base64.StdEncoding.DecodeString(dc.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("registry disclosure %s: decode blob: %w", dc.ContractId, err)
	}
	synchronizerID := ""
	if dc.SynchronizerId != nil {
		synchronizerID = *dc.SynchronizerId
	}

	return &apiv2.DisclosedContract{
		TemplateId:       templateID.ToLedgerIdentifier(),
		ContractId:       dc.ContractId,
		CreatedEventBlob: blob,
		SynchronizerId:   synchronizerID,
	}, nil
}
