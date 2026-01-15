package types

// contract disclosure for off-ledger distribution
type DisclosedContract struct {
	ContractID       string     `json:"contractId"`
	TemplateID       TemplateID `json:"templateId"`
	CreatedEventBlob string     `json:"createdEventBlob"` // base64-encoded
	SynchronizerID   string     `json:"synchronizerId,omitempty"`
}

// daml template identifier
type TemplateID struct {
	PackageID  string `json:"packageId"`
	ModuleName string `json:"moduleName"`
	EntityName string `json:"entityName"`
}

// contracts needed for CCIPSend
type CCIPSendDisclosures struct {
	InstanceID string            `json:"instanceId"`
	Contracts  CCIPSendContracts `json:"contracts"`
}

type CCIPSendContracts struct {
	Router    *DisclosedContract `json:"router"`
	OnRamp    *DisclosedContract `json:"onRamp"`
	FeeQuoter *DisclosedContract `json:"feeQuoter"`
}

// contracts needed for CCIPExecute
type CCIPExecuteDisclosures struct {
	InstanceID string               `json:"instanceId"`
	Contracts  CCIPExecuteContracts `json:"contracts"`
}

type CCIPExecuteContracts struct {
	OffRamp            *DisclosedContract `json:"offRamp"`
	CCV                *DisclosedContract `json:"ccv"`
	TokenAdminRegistry *DisclosedContract `json:"tokenAdminRegistry"`
}

type EnvironmentInfo struct {
	ID          string `json:"id"`
	Party       string `json:"party"`
	Description string `json:"description"`
}

type EnvironmentsResponse struct {
	Environments []EnvironmentInfo `json:"environments"`
}

type HealthResponse struct {
	Status             string   `json:"status"`
	LedgerAPIConnected bool     `json:"ledgerApiConnected"`
	Environments       []string `json:"environments"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
