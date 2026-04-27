package tokenadminregistry

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

const (
	PackageName = "ccip-tokenadminregistry"
	PackageID   = "51c3e6adc0534fd50889cad4b07d63c7523263b162c045157ed2ab2fd9df8117"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

func argsToMap(args any) map[string]any {
	if args == nil {
		return map[string]any{}
	}

	if m, ok := args.(map[string]any); ok {
		return m
	}

	type mapper interface {
		ToMap() map[string]any
	}
	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]any{"args": args}
}

// AcceptAdminParams is a Record type
type AcceptAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// ToMap converts AcceptAdminParams to a map for DAML arguments
func (t AcceptAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	return m
}

func (t AcceptAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptAdminParams to hex string (Canton MCMS format)
func (t AcceptAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminParams from hex string (Canton MCMS format)
func (t *AcceptAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptAdminRole is a Record type
type AcceptAdminRole struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts AcceptAdminRole to a map for DAML arguments
func (t AcceptAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AcceptAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptAdminRole to hex string (Canton MCMS format)
func (t AcceptAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminRole from hex string (Canton MCMS format)
func (t *AcceptAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptAdminRoleMCMSParams is AcceptAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AcceptAdminRoleMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes AcceptAdminRoleMCMSParams to hex string for MCMS operationData.
func (t AcceptAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminRoleMCMSParams from hex string.
func (t *AcceptAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSend2 is a Record type
type AddTokenSend2 struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                               `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount            types.NUMERIC                            `json:"amount"`
	DestTokenAddress  types.TEXT                               `json:"destTokenAddress"`
	ExtraData         types.TEXT                               `json:"extraData"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts AddTokenSend2 to a map for DAML arguments
func (t AddTokenSend2) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = t.Amount

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["extraData"] = string(t.ExtraData)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddTokenSend2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddTokenSend2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddTokenSend2 to hex string (Canton MCMS format)
func (t AddTokenSend2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSend2 from hex string (Canton MCMS format)
func (t *AddTokenSend2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSend2MCMSParams is AddTokenSend2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddTokenSend2MCMSParams struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                               `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount            types.NUMERIC                            `json:"amount"`
	DestTokenAddress  types.TEXT                               `json:"destTokenAddress"`
	ExtraData         types.TEXT                               `json:"extraData"`
}

// MarshalHex encodes AddTokenSend2MCMSParams to hex string for MCMS operationData.
func (t AddTokenSend2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSend2MCMSParams from hex string.
func (t *AddTokenSend2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSendFee2 is a Record type
type AddTokenSendFee2 struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT        `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC     `json:"feeUSDCents"`
	DestGasOverhead   types.INT64       `json:"destGasOverhead"`
	DestBytesOverhead types.INT64       `json:"destBytesOverhead"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts AddTokenSendFee2 to a map for DAML arguments
func (t AddTokenSendFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddTokenSendFee2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddTokenSendFee2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddTokenSendFee2 to hex string (Canton MCMS format)
func (t AddTokenSendFee2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSendFee2 from hex string (Canton MCMS format)
func (t *AddTokenSendFee2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSendFee2MCMSParams is AddTokenSendFee2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddTokenSendFee2MCMSParams struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT        `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC     `json:"feeUSDCents"`
	DestGasOverhead   types.INT64       `json:"destGasOverhead"`
	DestBytesOverhead types.INT64       `json:"destBytesOverhead"`
}

// MarshalHex encodes AddTokenSendFee2MCMSParams to hex string for MCMS operationData.
func (t AddTokenSendFee2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSendFee2MCMSParams from hex string.
func (t *AddTokenSendFee2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeReceiveTicket is a Record type
type ConsumeReceiveTicket struct {
	TicketCid    types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts ConsumeReceiveTicket to a map for DAML arguments
func (t ConsumeReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["ticketCid"] = model.NestedToDAMLValue(t.TicketCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ConsumeReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeReceiveTicket to hex string (Canton MCMS format)
func (t ConsumeReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeReceiveTicket from hex string (Canton MCMS format)
func (t *ConsumeReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeReceiveTicketMCMSParams is ConsumeReceiveTicket without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ConsumeReceiveTicketMCMSParams struct {
	TicketCid    types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes ConsumeReceiveTicketMCMSParams to hex string for MCMS operationData.
func (t ConsumeReceiveTicketMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeReceiveTicketMCMSParams from hex string.
func (t *ConsumeReceiveTicketMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeExecute2 is a Record type
type FinalizeExecute2 struct {
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	TicketReceiver      types.PARTY       `json:"ticketReceiver"`
	ReturnData          types.TEXT        `json:"returnData"`
	Caller              types.PARTY       `json:"caller"`
}

// ToMap converts FinalizeExecute2 to a map for DAML arguments
func (t FinalizeExecute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["ticketReceiver"] = t.TicketReceiver.ToMap()

	m["returnData"] = string(t.ReturnData)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FinalizeExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeExecute2 to hex string (Canton MCMS format)
func (t FinalizeExecute2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeExecute2 from hex string (Canton MCMS format)
func (t *FinalizeExecute2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeExecute2MCMSParams is FinalizeExecute2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type FinalizeExecute2MCMSParams struct {
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	TicketReceiver      types.PARTY       `json:"ticketReceiver"`
	ReturnData          types.TEXT        `json:"returnData"`
}

// MarshalHex encodes FinalizeExecute2MCMSParams to hex string for MCMS operationData.
func (t FinalizeExecute2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeExecute2MCMSParams from hex string.
func (t *FinalizeExecute2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Get2 is a Record type
type Get2 struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts Get2 to a map for DAML arguments
func (t Get2) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Get2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Get2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Get2 to hex string (Canton MCMS format)
func (t Get2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get2 from hex string (Canton MCMS format)
func (t *Get2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Get2MCMSParams is Get2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type Get2MCMSParams struct {
}

// MarshalHex encodes Get2MCMSParams to hex string for MCMS operationData.
func (t Get2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get2MCMSParams from hex string.
func (t *Get2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenConfig is a Record type
type GetTokenConfig struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenConfig to a map for DAML arguments
func (t GetTokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenConfig to hex string (Canton MCMS format)
func (t GetTokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenConfig from hex string (Canton MCMS format)
func (t *GetTokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenConfigMCMSParams is GetTokenConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenConfigMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes GetTokenConfigMCMSParams to hex string for MCMS operationData.
func (t GetTokenConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenConfigMCMSParams from hex string.
func (t *GetTokenConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsAdministrator is a Record type
type IsAdministrator struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Administrator types.PARTY                              `json:"administrator"`
	Caller        types.PARTY                              `json:"caller"`
}

// ToMap converts IsAdministrator to a map for DAML arguments
func (t IsAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["administrator"] = t.Administrator.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsAdministrator to hex string (Canton MCMS format)
func (t IsAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsAdministrator from hex string (Canton MCMS format)
func (t *IsAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsAdministratorMCMSParams is IsAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsAdministratorMCMSParams struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Administrator types.PARTY                              `json:"administrator"`
}

// MarshalHex encodes IsAdministratorMCMSParams to hex string for MCMS operationData.
func (t IsAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsAdministratorMCMSParams from hex string.
func (t *IsAdministratorMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PoolRegistration is a Record type
type PoolRegistration struct {
	PoolOwner      types.PARTY `json:"poolOwner"`
	PoolInstanceId types.TEXT  `json:"poolInstanceId"`
}

// ToMap converts PoolRegistration to a map for DAML arguments
func (t PoolRegistration) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	return m
}

func (t PoolRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PoolRegistration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PoolRegistration to hex string (Canton MCMS format)
func (t PoolRegistration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PoolRegistration from hex string (Canton MCMS format)
func (t *PoolRegistration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdminParams is a Record type
type ProposeAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// ToMap converts ProposeAdminParams to a map for DAML arguments
func (t ProposeAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t ProposeAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProposeAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProposeAdminParams to hex string (Canton MCMS format)
func (t ProposeAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdminParams from hex string (Canton MCMS format)
func (t *ProposeAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministrator is a Record type
type ProposeAdministrator struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts ProposeAdministrator to a map for DAML arguments
func (t ProposeAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ProposeAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProposeAdministrator to hex string (Canton MCMS format)
func (t ProposeAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdministrator from hex string (Canton MCMS format)
func (t *ProposeAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministratorMCMSParams is ProposeAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ProposeAdministratorMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes ProposeAdministratorMCMSParams to hex string for MCMS operationData.
func (t ProposeAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdministratorMCMSParams from hex string.
func (t *ProposeAdministratorMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetInboundPoolCCVs2 is a Record type
type SetInboundPoolCCVs2 struct {
	ExecutingMessageCid types.CONTRACT_ID         `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                `json:"poolInstanceId"`
	PoolCCVs            []mcms.RawInstanceAddress `json:"poolCCVs"`
	Caller              types.PARTY               `json:"caller"`
}

// ToMap converts SetInboundPoolCCVs2 to a map for DAML arguments
func (t SetInboundPoolCCVs2) ToMap() map[string]any {
	m := make(map[string]any)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetInboundPoolCCVs2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetInboundPoolCCVs2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetInboundPoolCCVs2 to hex string (Canton MCMS format)
func (t SetInboundPoolCCVs2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetInboundPoolCCVs2 from hex string (Canton MCMS format)
func (t *SetInboundPoolCCVs2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetInboundPoolCCVs2MCMSParams is SetInboundPoolCCVs2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetInboundPoolCCVs2MCMSParams struct {
	ExecutingMessageCid types.CONTRACT_ID         `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                `json:"poolInstanceId"`
	PoolCCVs            []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// MarshalHex encodes SetInboundPoolCCVs2MCMSParams to hex string for MCMS operationData.
func (t SetInboundPoolCCVs2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetInboundPoolCCVs2MCMSParams from hex string.
func (t *SetInboundPoolCCVs2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetOutboundPoolCCVs2 is a Record type
type SetOutboundPoolCCVs2 struct {
	SendingMessageCid types.CONTRACT_ID         `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                `json:"poolInstanceId"`
	PoolCCVs          []mcms.RawInstanceAddress `json:"poolCCVs"`
	Caller            types.PARTY               `json:"caller"`
}

// ToMap converts SetOutboundPoolCCVs2 to a map for DAML arguments
func (t SetOutboundPoolCCVs2) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetOutboundPoolCCVs2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetOutboundPoolCCVs2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetOutboundPoolCCVs2 to hex string (Canton MCMS format)
func (t SetOutboundPoolCCVs2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetOutboundPoolCCVs2 from hex string (Canton MCMS format)
func (t *SetOutboundPoolCCVs2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetOutboundPoolCCVs2MCMSParams is SetOutboundPoolCCVs2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetOutboundPoolCCVs2MCMSParams struct {
	SendingMessageCid types.CONTRACT_ID         `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                `json:"poolInstanceId"`
	PoolCCVs          []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// MarshalHex encodes SetOutboundPoolCCVs2MCMSParams to hex string for MCMS operationData.
func (t SetOutboundPoolCCVs2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetOutboundPoolCCVs2MCMSParams from hex string.
func (t *SetOutboundPoolCCVs2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPool is a Record type
type SetPool struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts SetPool to a map for DAML arguments
func (t SetPool) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPool to hex string (Canton MCMS format)
func (t SetPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPool from hex string (Canton MCMS format)
func (t *SetPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPoolMCMSParams is SetPool without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetPoolMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// MarshalHex encodes SetPoolMCMSParams to hex string for MCMS operationData.
func (t SetPoolMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPoolMCMSParams from hex string.
func (t *SetPoolMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPoolParams is a Record type
type SetPoolParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// ToMap converts SetPoolParams to a map for DAML arguments
func (t SetPoolParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetPoolParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPoolParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPoolParams to hex string (Canton MCMS format)
func (t SetPoolParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPoolParams from hex string (Canton MCMS format)
func (t *SetPoolParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	InstanceId   types.TEXT                 `json:"instanceId"`
	Owner        types.PARTY                `json:"owner"`
	TokenConfigs map[types.TEXT]TokenConfig `json:"tokenConfigs"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenAdminRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenAdminRegistry) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenAdminRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() any {
		if t.TokenConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenAdminRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() any {
		if t.TokenConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistry to hex string (Canton MCMS format)
func (t TokenAdminRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistry from hex string (Canton MCMS format)
func (t *TokenAdminRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenAdminRegistry

// ConsumeReceiveTicket exercises the ConsumeReceiveTicket choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) ConsumeReceiveTicket(contractID string, args ConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// ConsumeReceiveTicketWithPackageID exercises the ConsumeReceiveTicket choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ConsumeReceiveTicketWithPackageID(contractID string, packageID string, args ConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVs exercises the SetOutboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetOutboundPoolCCVs(contractID string, args SetOutboundPoolCCVs2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVsWithPackageID exercises the SetOutboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetOutboundPoolCCVsWithPackageID(contractID string, packageID string, args SetOutboundPoolCCVs2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFee exercises the AddTokenSendFee choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AddTokenSendFee(contractID string, args AddTokenSendFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFeeWithPackageID exercises the AddTokenSendFee choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AddTokenSendFeeWithPackageID(contractID string, packageID string, args AddTokenSendFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSend exercises the AddTokenSend choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AddTokenSend(contractID string, args AddTokenSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendWithPackageID exercises the AddTokenSend choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AddTokenSendWithPackageID(contractID string, packageID string, args AddTokenSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVs exercises the SetInboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetInboundPoolCCVs(contractID string, args SetInboundPoolCCVs2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVsWithPackageID exercises the SetInboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetInboundPoolCCVsWithPackageID(contractID string, packageID string, args SetInboundPoolCCVs2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecute exercises the FinalizeExecute choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) FinalizeExecute(contractID string, args FinalizeExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecuteWithPackageID exercises the FinalizeExecute choice using the provided package ID instead of package name
func (t TokenAdminRegistry) FinalizeExecuteWithPackageID(contractID string, packageID string, args FinalizeExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// GetTokenConfig exercises the GetTokenConfig choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) GetTokenConfig(contractID string, args GetTokenConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "GetTokenConfig",
		Arguments:  argsToMap(args),
	}
}

// GetTokenConfigWithPackageID exercises the GetTokenConfig choice using the provided package ID instead of package name
func (t TokenAdminRegistry) GetTokenConfigWithPackageID(contractID string, packageID string, args GetTokenConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "GetTokenConfig",
		Arguments:  argsToMap(args),
	}
}

// SetPool exercises the SetPool choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetPool(contractID string, args SetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetPool",
		Arguments:  argsToMap(args),
	}
}

// SetPoolWithPackageID exercises the SetPool choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetPoolWithPackageID(contractID string, packageID string, args SetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetPool",
		Arguments:  argsToMap(args),
	}
}

// AcceptAdminRole exercises the AcceptAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AcceptAdminRole(contractID string, args AcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// AcceptAdminRoleWithPackageID exercises the AcceptAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AcceptAdminRoleWithPackageID(contractID string, packageID string, args AcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TransferAdminRole exercises the TransferAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TransferAdminRole(contractID string, args TransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TransferAdminRoleWithPackageID exercises the TransferAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TransferAdminRoleWithPackageID(contractID string, packageID string, args TransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// IsAdministrator exercises the IsAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) IsAdministrator(contractID string, args IsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// IsAdministratorWithPackageID exercises the IsAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) IsAdministratorWithPackageID(contractID string, packageID string, args IsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// ProposeAdministrator exercises the ProposeAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) ProposeAdministrator(contractID string, args ProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// ProposeAdministratorWithPackageID exercises the ProposeAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ProposeAdministratorWithPackageID(contractID string, packageID string, args ProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Get(contractID string, args Get2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t TokenAdminRegistry) GetWithPackageID(contractID string, packageID string, args Get2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t TokenAdminRegistry) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for TokenAdminRegistry

var _ mcms.IMCMSReceiver = (*TokenAdminRegistry)(nil)

// TokenConfig is a Record type
type TokenConfig struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Admin        *types.PARTY                             `json:"admin" hex:"optional"`
	PendingAdmin *types.PARTY                             `json:"pendingAdmin" hex:"optional"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// ToMap converts TokenConfig to a map for DAML arguments
func (t TokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.Admin != nil {
		m["admin"] = map[string]any{
			"_type": "optional",
			"value": (*t.Admin).ToMap(),
		}
	} else {
		m["admin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.PendingAdmin != nil {
		m["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.PendingAdmin).ToMap(),
		}
	} else {
		m["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfig to hex string (Canton MCMS format)
func (t TokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfig from hex string (Canton MCMS format)
func (t *TokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferAdminParams is a Record type
type TransferAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// ToMap converts TransferAdminParams to a map for DAML arguments
func (t TransferAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t TransferAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferAdminParams to hex string (Canton MCMS format)
func (t TransferAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminParams from hex string (Canton MCMS format)
func (t *TransferAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferAdminRole is a Record type
type TransferAdminRole struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TransferAdminRole to a map for DAML arguments
func (t TransferAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TransferAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferAdminRole to hex string (Canton MCMS format)
func (t TransferAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminRole from hex string (Canton MCMS format)
func (t *TransferAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferAdminRoleMCMSParams is TransferAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TransferAdminRoleMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes TransferAdminRoleMCMSParams to hex string for MCMS operationData.
func (t TransferAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminRoleMCMSParams from hex string.
func (t *TransferAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AcceptAdminRole(args AcceptAdminRole) (*bind.EncodedChoice, error)
	AcceptAdminRoleMCMSParams(args AcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error)
	AddTokenSend(args AddTokenSend2) (*bind.EncodedChoice, error)
	AddTokenSendMCMSParams(args AddTokenSend2MCMSParams) (*bind.EncodedChoice, error)
	AddTokenSendFee(args AddTokenSendFee2) (*bind.EncodedChoice, error)
	AddTokenSendFeeMCMSParams(args AddTokenSendFee2MCMSParams) (*bind.EncodedChoice, error)
	ConsumeReceiveTicket(args ConsumeReceiveTicket) (*bind.EncodedChoice, error)
	ConsumeReceiveTicketMCMSParams(args ConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error)
	FinalizeExecute(args FinalizeExecute2) (*bind.EncodedChoice, error)
	FinalizeExecuteMCMSParams(args FinalizeExecute2MCMSParams) (*bind.EncodedChoice, error)
	Get(args Get2) (*bind.EncodedChoice, error)
	GetMCMSParams(args Get2MCMSParams) (*bind.EncodedChoice, error)
	GetTokenConfig(args GetTokenConfig) (*bind.EncodedChoice, error)
	GetTokenConfigMCMSParams(args GetTokenConfigMCMSParams) (*bind.EncodedChoice, error)
	IsAdministrator(args IsAdministrator) (*bind.EncodedChoice, error)
	IsAdministratorMCMSParams(args IsAdministratorMCMSParams) (*bind.EncodedChoice, error)
	ProposeAdministrator(args ProposeAdministrator) (*bind.EncodedChoice, error)
	ProposeAdministratorMCMSParams(args ProposeAdministratorMCMSParams) (*bind.EncodedChoice, error)
	SetInboundPoolCCVs(args SetInboundPoolCCVs2) (*bind.EncodedChoice, error)
	SetInboundPoolCCVsMCMSParams(args SetInboundPoolCCVs2MCMSParams) (*bind.EncodedChoice, error)
	SetOutboundPoolCCVs(args SetOutboundPoolCCVs2) (*bind.EncodedChoice, error)
	SetOutboundPoolCCVsMCMSParams(args SetOutboundPoolCCVs2MCMSParams) (*bind.EncodedChoice, error)
	SetPool(args SetPool) (*bind.EncodedChoice, error)
	SetPoolMCMSParams(args SetPoolMCMSParams) (*bind.EncodedChoice, error)
	SetPoolParams(args SetPoolParams) (*bind.EncodedChoice, error)
	TransferAdminRole(args TransferAdminRole) (*bind.EncodedChoice, error)
	TransferAdminRoleMCMSParams(args TransferAdminRoleMCMSParams) (*bind.EncodedChoice, error)
}

// encoder provides typed encoding methods for choice parameters (unexported).
// It wraps bind.BoundTemplate to encode parameters to hex-encoded operation data.
type encoder struct {
	*bind.BoundTemplate
}

// Contract wraps template operations with Sui-style API access.
// Use NewContract to create instances, then call Encoder() for encoding methods.
type Contract struct {
	enc *encoder
}

// NewContract creates a Contract with encoder for the given template.
// This provides Sui-style API: contract.Encoder().Method(args)
func NewContract(packageID, moduleName, templateName string) *Contract {
	return &Contract{
		enc: &encoder{
			BoundTemplate: bind.NewBoundTemplate(packageID, moduleName, templateName),
		},
	}
}

// Encoder returns the encoder for Sui-style contract.Encoder().Method() usage.
func (c *Contract) Encoder() MCMSEncoder {
	return c.enc
}

// AcceptAdminRole encodes parameters for the AcceptAdminRole choice.
func (e *encoder) AcceptAdminRole(args AcceptAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptAdminRole", args)
}

// AcceptAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the AcceptAdminRole choice.
func (e *encoder) AcceptAdminRoleMCMSParams(args AcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptAdminRole", args)
}

// AddTokenSend encodes parameters for the AddTokenSend choice.
func (e *encoder) AddTokenSend(args AddTokenSend2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSend", args)
}

// AddTokenSendMCMSParams encodes MCMS parameters (without Caller) for the AddTokenSend choice.
func (e *encoder) AddTokenSendMCMSParams(args AddTokenSend2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSend", args)
}

// AddTokenSendFee encodes parameters for the AddTokenSendFee choice.
func (e *encoder) AddTokenSendFee(args AddTokenSendFee2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSendFee", args)
}

// AddTokenSendFeeMCMSParams encodes MCMS parameters (without Caller) for the AddTokenSendFee choice.
func (e *encoder) AddTokenSendFeeMCMSParams(args AddTokenSendFee2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSendFee", args)
}

// ConsumeReceiveTicket encodes parameters for the ConsumeReceiveTicket choice.
func (e *encoder) ConsumeReceiveTicket(args ConsumeReceiveTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeReceiveTicket", args)
}

// ConsumeReceiveTicketMCMSParams encodes MCMS parameters (without Caller) for the ConsumeReceiveTicket choice.
func (e *encoder) ConsumeReceiveTicketMCMSParams(args ConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeReceiveTicket", args)
}

// FinalizeExecute encodes parameters for the FinalizeExecute choice.
func (e *encoder) FinalizeExecute(args FinalizeExecute2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeExecute", args)
}

// FinalizeExecuteMCMSParams encodes MCMS parameters (without Caller) for the FinalizeExecute choice.
func (e *encoder) FinalizeExecuteMCMSParams(args FinalizeExecute2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeExecute", args)
}

// Get encodes parameters for the Get choice.
func (e *encoder) Get(args Get2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetMCMSParams encodes MCMS parameters (without Caller) for the Get choice.
func (e *encoder) GetMCMSParams(args Get2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetTokenConfig encodes parameters for the GetTokenConfig choice.
func (e *encoder) GetTokenConfig(args GetTokenConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenConfig", args)
}

// GetTokenConfigMCMSParams encodes MCMS parameters (without Caller) for the GetTokenConfig choice.
func (e *encoder) GetTokenConfigMCMSParams(args GetTokenConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenConfig", args)
}

// IsAdministrator encodes parameters for the IsAdministrator choice.
func (e *encoder) IsAdministrator(args IsAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsAdministrator", args)
}

// IsAdministratorMCMSParams encodes MCMS parameters (without Caller) for the IsAdministrator choice.
func (e *encoder) IsAdministratorMCMSParams(args IsAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsAdministrator", args)
}

// ProposeAdministrator encodes parameters for the ProposeAdministrator choice.
func (e *encoder) ProposeAdministrator(args ProposeAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProposeAdministrator", args)
}

// ProposeAdministratorMCMSParams encodes MCMS parameters (without Caller) for the ProposeAdministrator choice.
func (e *encoder) ProposeAdministratorMCMSParams(args ProposeAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProposeAdministrator", args)
}

// SetInboundPoolCCVs encodes parameters for the SetInboundPoolCCVs choice.
func (e *encoder) SetInboundPoolCCVs(args SetInboundPoolCCVs2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetInboundPoolCCVs", args)
}

// SetInboundPoolCCVsMCMSParams encodes MCMS parameters (without Caller) for the SetInboundPoolCCVs choice.
func (e *encoder) SetInboundPoolCCVsMCMSParams(args SetInboundPoolCCVs2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetInboundPoolCCVs", args)
}

// SetOutboundPoolCCVs encodes parameters for the SetOutboundPoolCCVs choice.
func (e *encoder) SetOutboundPoolCCVs(args SetOutboundPoolCCVs2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetOutboundPoolCCVs", args)
}

// SetOutboundPoolCCVsMCMSParams encodes MCMS parameters (without Caller) for the SetOutboundPoolCCVs choice.
func (e *encoder) SetOutboundPoolCCVsMCMSParams(args SetOutboundPoolCCVs2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetOutboundPoolCCVs", args)
}

// SetPool encodes parameters for the SetPool choice.
func (e *encoder) SetPool(args SetPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// SetPoolMCMSParams encodes MCMS parameters (without Caller) for the SetPool choice.
func (e *encoder) SetPoolMCMSParams(args SetPoolMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// SetPoolParams encodes parameters for the SetPool choice.
func (e *encoder) SetPoolParams(args SetPoolParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// TransferAdminRole encodes parameters for the TransferAdminRole choice.
func (e *encoder) TransferAdminRole(args TransferAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferAdminRole", args)
}

// TransferAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the TransferAdminRole choice.
func (e *encoder) TransferAdminRoleMCMSParams(args TransferAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferAdminRole", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
