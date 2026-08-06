package operator

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
)

// Client calls DA's hosted Utilities operator backend for Registry choice context.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChoiceContext bundles extraArgs.context and disclosed contracts from the operator backend.
type ChoiceContext struct {
	Context            splice_api_token_metadata_v1.ChoiceContext
	DisclosedContracts []*apiv2.DisclosedContract
}

type instrumentIDJSON struct {
	Admin string `json:"admin"`
	ID    string `json:"id"`
}

type mintRequestBody struct {
	Holder       string           `json:"holder"`
	InstrumentID instrumentIDJSON `json:"instrumentId"`
}

type mintContextResponse struct {
	ChoiceContext *struct {
		ChoiceContextData  splice_api_token_metadata_v1.ChoiceContext `json:"choiceContextData"`
		DisclosedContracts json.RawMessage                            `json:"disclosedContracts"`
	} `json:"choiceContext"`
	ChoiceContextData  splice_api_token_metadata_v1.ChoiceContext `json:"choiceContextData"`
	DisclosedContracts json.RawMessage                            `json:"disclosedContracts"`
}

type acceptContextBody struct {
	Meta               map[string]any `json:"meta"`
	ExcludeDebugFields bool           `json:"excludeDebugFields"`
}

// MintRequestContext fetches choice context for AllocationFactory_RequestMint.
func (c *Client) MintRequestContext(ctx context.Context, holder string, instrumentID splice_api_token_holding_v1.InstrumentId) (ChoiceContext, error) {
	body, err := json.Marshal(mintRequestBody{
		Holder: holder,
		InstrumentID: instrumentIDJSON{
			Admin: string(instrumentID.Admin),
			ID:    string(instrumentID.Id),
		},
	})
	if err != nil {
		return ChoiceContext{}, err
	}

	var resp mintContextResponse
	if err := c.postJSON(ctx, "/v0/registry/mint/v0/request", body, &resp); err != nil {
		return ChoiceContext{}, err
	}

	return resp.toChoiceContext()
}

// MintAcceptContext fetches choice context for MintRequest_Accept.
func (c *Client) MintAcceptContext(ctx context.Context, mintRequestCID string) (ChoiceContext, error) {
	body, err := json.Marshal(acceptContextBody{Meta: map[string]any{}, ExcludeDebugFields: true})
	if err != nil {
		return ChoiceContext{}, err
	}

	var resp mintContextResponse
	path := fmt.Sprintf("/v0/registry/mint/v0/request/%s/choice-contexts/accept", mintRequestCID)
	if err := c.postJSON(ctx, path, body, &resp); err != nil {
		return ChoiceContext{}, err
	}

	return resp.toChoiceContext()
}

func (r mintContextResponse) toChoiceContext() (ChoiceContext, error) {
	ctx := r.ChoiceContextData
	if r.ChoiceContext != nil {
		ctx = r.ChoiceContext.ChoiceContextData
	}

	raw := r.DisclosedContracts
	if r.ChoiceContext != nil && len(r.ChoiceContext.DisclosedContracts) > 0 {
		raw = r.ChoiceContext.DisclosedContracts
	}

	disclosed, err := parseDisclosedContracts(raw)
	if err != nil {
		return ChoiceContext{}, err
	}

	return ChoiceContext{Context: ctx, DisclosedContracts: disclosed}, nil
}

func (c *Client) postJSON(ctx context.Context, path string, body []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("operator backend %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("operator backend %s: status %d: %s", path, resp.StatusCode, string(respBody))
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode operator backend response: %w", err)
	}

	return nil
}

// disclosedContractJSON matches Canton JSON API disclosed contract shape.
type disclosedContractJSON struct {
	TemplateID       string `json:"templateId"`
	ContractID       string `json:"contractId"`
	CreatedEventBlob string `json:"createdEventBlob"`
	SynchronizerID   string `json:"synchronizerId"`
}

func parseDisclosedContracts(raw json.RawMessage) ([]*apiv2.DisclosedContract, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var items []disclosedContractJSON
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("parse disclosed contracts: %w", err)
	}

	out := make([]*apiv2.DisclosedContract, 0, len(items))
	for _, item := range items {
		blob, err := base64.StdEncoding.DecodeString(item.CreatedEventBlob)
		if err != nil {
			return nil, fmt.Errorf("decode createdEventBlob for %s: %w", item.ContractID, err)
		}
		dc := &apiv2.DisclosedContract{
			ContractId:       item.ContractID,
			CreatedEventBlob: blob,
			SynchronizerId:   item.SynchronizerID,
		}
		if item.TemplateID != "" {
			dc.TemplateId = parseTemplateID(item.TemplateID)
		}
		out = append(out, dc)
	}

	return out, nil
}

func parseTemplateID(raw string) *apiv2.Identifier {
	// Accept packageId:module:entity or #packageName:module:entity — entity name only needed for submission.
	parts := splitTemplateID(raw)
	if len(parts) < 3 {
		return &apiv2.Identifier{EntityName: raw}
	}
	id := &apiv2.Identifier{
		ModuleName: parts[1],
		EntityName: parts[2],
	}
	if parts[0] != "" && parts[0][0] != '#' {
		id.PackageId = parts[0]
	}

	return id
}

func splitTemplateID(raw string) []string {
	out := make([]string, 0, 3)
	start := 0
	for i := range len(raw) {
		if raw[i] == ':' {
			out = append(out, raw[start:i])
			start = i + 1
		}
	}
	out = append(out, raw[start:])

	return out
}
