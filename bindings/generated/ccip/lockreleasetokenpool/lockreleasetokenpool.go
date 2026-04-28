package lockreleasetokenpool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
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
	PackageName = "ccip-lockreleasetokenpool"
	PackageID   = "81d9209c78c7867f4fc8f49cfbe4b0bc1f06087903381eb257a1d88cff0166b6"
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

// ApplyChainUpdates is a Record type
type ApplyChainUpdates struct {
	RemoteChainSelectorsToRemove []types.NUMERIC `json:"remoteChainSelectorsToRemove"`
	ChainsToAdd                  []ChainUpdate   `json:"chainsToAdd"`
}

// ToMap converts ApplyChainUpdates to a map for DAML arguments
func (t ApplyChainUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.RemoteChainSelectorsToRemove))
		for _, e := range t.RemoteChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["chainsToAdd"] = func() []any {
		res := make([]any, 0, len(t.ChainsToAdd))
		for _, e := range t.ChainsToAdd {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyChainUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyChainUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyChainUpdates to hex string (Canton MCMS format)
func (t ApplyChainUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyChainUpdates from hex string (Canton MCMS format)
func (t *ApplyChainUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyChainUpdatesParams is a Record type
type ApplyChainUpdatesParams struct {
	RemoteChainSelectorsToRemove []types.NUMERIC `json:"remoteChainSelectorsToRemove"`
	ChainsToAdd                  []ChainUpdate   `json:"chainsToAdd"`
}

// ToMap converts ApplyChainUpdatesParams to a map for DAML arguments
func (t ApplyChainUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.RemoteChainSelectorsToRemove))
		for _, e := range t.RemoteChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["chainsToAdd"] = func() []any {
		res := make([]any, 0, len(t.ChainsToAdd))
		for _, e := range t.ChainsToAdd {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyChainUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyChainUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyChainUpdatesParams to hex string (Canton MCMS format)
func (t ApplyChainUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyChainUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyChainUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyTokenTransferFeeConfigUpdates is a Record type
type ApplyTokenTransferFeeConfigUpdates struct {
	TokenTransferFeeConfigArgs        []TokenTransferFeeConfigArgs `json:"tokenTransferFeeConfigArgs"`
	DisableTokenTransferFeeConfigArgs []types.NUMERIC              `json:"disableTokenTransferFeeConfigArgs"`
}

// ToMap converts ApplyTokenTransferFeeConfigUpdates to a map for DAML arguments
func (t ApplyTokenTransferFeeConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenTransferFeeConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.TokenTransferFeeConfigArgs))
		for _, e := range t.TokenTransferFeeConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["disableTokenTransferFeeConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DisableTokenTransferFeeConfigArgs))
		for _, e := range t.DisableTokenTransferFeeConfigArgs {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t ApplyTokenTransferFeeConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyTokenTransferFeeConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyTokenTransferFeeConfigUpdates to hex string (Canton MCMS format)
func (t ApplyTokenTransferFeeConfigUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyTokenTransferFeeConfigUpdates from hex string (Canton MCMS format)
func (t *ApplyTokenTransferFeeConfigUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyTokenTransferFeeConfigUpdatesParams is a Record type
type ApplyTokenTransferFeeConfigUpdatesParams struct {
	TokenTransferFeeConfigArgs        []TokenTransferFeeConfigArgs `json:"tokenTransferFeeConfigArgs"`
	DisableTokenTransferFeeConfigArgs []types.NUMERIC              `json:"disableTokenTransferFeeConfigArgs"`
}

// ToMap converts ApplyTokenTransferFeeConfigUpdatesParams to a map for DAML arguments
func (t ApplyTokenTransferFeeConfigUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenTransferFeeConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.TokenTransferFeeConfigArgs))
		for _, e := range t.TokenTransferFeeConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["disableTokenTransferFeeConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DisableTokenTransferFeeConfigArgs))
		for _, e := range t.DisableTokenTransferFeeConfigArgs {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t ApplyTokenTransferFeeConfigUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyTokenTransferFeeConfigUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyTokenTransferFeeConfigUpdatesParams to hex string (Canton MCMS format)
func (t ApplyTokenTransferFeeConfigUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyTokenTransferFeeConfigUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyTokenTransferFeeConfigUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFee is a Record type
type CalculateFee struct {
	TokenAdminRegistryCid types.CONTRACT_ID                        `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                        `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext                       `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID                        `json:"sendingMessageCid"`
	FeeQuoterCid          types.CONTRACT_ID                        `json:"feeQuoterCid"`
	TokenInstrumentId     splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	Caller                types.PARTY                              `json:"caller"`
}

// ToMap converts CalculateFee to a map for DAML arguments
func (t CalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CalculateFee to hex string (Canton MCMS format)
func (t CalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFee from hex string (Canton MCMS format)
func (t *CalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFeeMCMSParams is CalculateFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CalculateFeeMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID                        `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                        `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext                       `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID                        `json:"sendingMessageCid"`
	FeeQuoterCid          types.CONTRACT_ID                        `json:"feeQuoterCid"`
	TokenInstrumentId     splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
}

// MarshalHex encodes CalculateFeeMCMSParams to hex string for MCMS operationData.
func (t CalculateFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFeeMCMSParams from hex string.
func (t *CalculateFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ChainUpdate is a Record type
type ChainUpdate struct {
	RemoteChainSelector                        types.NUMERIC             `json:"remoteChainSelector"`
	RemotePools                                []types.TEXT              `json:"remotePools"`
	RemoteTokenAddress                         types.TEXT                `json:"remoteTokenAddress"`
	InboundCCVs                                []mcms.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []mcms.RawInstanceAddress `json:"outboundCCVs"`
	FinalityConfig                             common.FinalityConfig     `json:"finalityConfig"`
	InboundRateLimiter                         mcms.RawInstanceAddress   `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter mcms.RawInstanceAddress   `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        mcms.RawInstanceAddress   `json:"outboundRateLimiter"`
}

// ToMap converts ChainUpdate to a map for DAML arguments
func (t ChainUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["remotePools"] = func() []any {
		res := make([]any, 0, len(t.RemotePools))
		for _, e := range t.RemotePools {
			res = append(res, string(e))
		}
		return res
	}()

	m["remoteTokenAddress"] = string(t.RemoteTokenAddress)

	m["inboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.InboundCCVs))
		for _, e := range t.InboundCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["outboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.OutboundCCVs))
		for _, e := range t.OutboundCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["finalityConfig"] = model.NestedToDAMLValue(t.FinalityConfig)

	m["inboundRateLimiter"] = model.NestedToDAMLValue(t.InboundRateLimiter)

	m["inboundCustomBlockConfirmationsRateLimiter"] = model.NestedToDAMLValue(t.InboundCustomBlockConfirmationsRateLimiter)

	m["outboundRateLimiter"] = model.NestedToDAMLValue(t.OutboundRateLimiter)

	return m
}

func (t ChainUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ChainUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ChainUpdate to hex string (Canton MCMS format)
func (t ChainUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ChainUpdate from hex string (Canton MCMS format)
func (t *ChainUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee is a Record type
type GetFee struct {
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["destChainSelector"] = t.DestChainSelector

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee to hex string (Canton MCMS format)
func (t GetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee from hex string (Canton MCMS format)
func (t *GetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeMCMSParams is GetFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFeeMCMSParams struct {
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
}

// MarshalHex encodes GetFeeMCMSParams to hex string for MCMS operationData.
func (t GetFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeMCMSParams from hex string.
func (t *GetFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVs is a Record type
type GetRequiredCCVs struct {
	RemoteChainSelector types.NUMERIC                `json:"remoteChainSelector"`
	Amount              types.NUMERIC                `json:"amount"`
	Finality            common.FinalityConfig        `json:"finality"`
	ExtraData           types.TEXT                   `json:"extraData"`
	Direction           interfaces.TransferDirection `json:"direction"`
	Caller              types.PARTY                  `json:"caller"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["amount"] = t.Amount

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = model.NestedToDAMLValue(t.Direction)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVs to hex string (Canton MCMS format)
func (t GetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVs from hex string (Canton MCMS format)
func (t *GetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsMCMSParams is GetRequiredCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetRequiredCCVsMCMSParams struct {
	RemoteChainSelector types.NUMERIC                `json:"remoteChainSelector"`
	Amount              types.NUMERIC                `json:"amount"`
	Finality            common.FinalityConfig        `json:"finality"`
	ExtraData           types.TEXT                   `json:"extraData"`
	Direction           interfaces.TransferDirection `json:"direction"`
}

// MarshalHex encodes GetRequiredCCVsMCMSParams to hex string for MCMS operationData.
func (t GetRequiredCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsMCMSParams from hex string.
func (t *GetRequiredCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOrBurn is a Record type
type LockOrBurn struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID     `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids       []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount                types.NUMERIC         `json:"amount"`
	Caller                types.PARTY           `json:"caller"`
}

// ToMap converts LockOrBurn to a map for DAML arguments
func (t LockOrBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["tokenInput"] = model.NestedToDAMLValue(t.TokenInput)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOrBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOrBurn to hex string (Canton MCMS format)
func (t LockOrBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOrBurn from hex string (Canton MCMS format)
func (t *LockOrBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOrBurnMCMSParams is LockOrBurn without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockOrBurnMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID     `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids       []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount                types.NUMERIC         `json:"amount"`
}

// MarshalHex encodes LockOrBurnMCMSParams to hex string for MCMS operationData.
func (t LockOrBurnMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOrBurnMCMSParams from hex string.
func (t *LockOrBurnMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPool is a Template type
type LockReleaseTokenPool struct {
	InstanceId              types.TEXT                                `json:"instanceId"`
	PoolOwner               types.PARTY                               `json:"poolOwner"`
	CcipOwner               types.PARTY                               `json:"ccipOwner"`
	InstrumentId            splice_api_token_holding_v1.InstrumentId  `json:"instrumentId"`
	Decimals                types.INT64                               `json:"decimals"`
	RateLimitAdmin          *types.PARTY                              `json:"rateLimitAdmin" hex:"optional"`
	RemoteChainConfigs      map[types.NUMERIC]RemoteChainConfig       `json:"remoteChainConfigs"`
	TokenTransferFeeConfigs map[types.NUMERIC]TokenTransferFeeConfig2 `json:"tokenTransferFeeConfigs"`
	PoolReceiveContext      common.CCIPContext                        `json:"poolReceiveContext"`
	TransferTimeout         TransferTimeout                           `json:"transferTimeout"`
	Deps                    LockReleaseTokenPoolDeps                  `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LockReleaseTokenPool) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LockReleaseTokenPool) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LockReleaseTokenPool) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["decimals"] = int64(t.Decimals)

	if t.RateLimitAdmin != nil {
		args["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		args["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolReceiveContext"] = model.NestedToDAMLValue(t.PoolReceiveContext)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = model.NestedToDAMLValue(t.TransferTimeout)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LockReleaseTokenPool) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["decimals"] = int64(t.Decimals)

	if t.RateLimitAdmin != nil {
		args["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		args["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolReceiveContext"] = model.NestedToDAMLValue(t.PoolReceiveContext)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = model.NestedToDAMLValue(t.TransferTimeout)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LockReleaseTokenPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPool to hex string (Canton MCMS format)
func (t LockReleaseTokenPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPool from hex string (Canton MCMS format)
func (t *LockReleaseTokenPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LockReleaseTokenPool

// ReleaseFromTicket exercises the ReleaseFromTicket choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) ReleaseFromTicket(contractID string, args ReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// ReleaseFromTicketWithPackageID exercises the ReleaseFromTicket choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ReleaseFromTicketWithPackageID(contractID string, packageID string, args ReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// LockOrBurn exercises the LockOrBurn choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockOrBurn(contractID string, args LockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// LockOrBurnWithPackageID exercises the LockOrBurn choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockOrBurnWithPackageID(contractID string, packageID string, args LockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// ApplyTokenTransferFeeConfigUpdates exercises the ApplyTokenTransferFeeConfigUpdates choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) ApplyTokenTransferFeeConfigUpdates(contractID string, args ApplyTokenTransferFeeConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyTokenTransferFeeConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyTokenTransferFeeConfigUpdatesWithPackageID exercises the ApplyTokenTransferFeeConfigUpdates choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ApplyTokenTransferFeeConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyTokenTransferFeeConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyTokenTransferFeeConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimitConfig exercises the SetRateLimitConfig choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) SetRateLimitConfig(contractID string, args SetRateLimitConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimitConfig",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimitConfigWithPackageID exercises the SetRateLimitConfig choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) SetRateLimitConfigWithPackageID(contractID string, packageID string, args SetRateLimitConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimitConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyChainUpdates exercises the ApplyChainUpdates choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) ApplyChainUpdates(contractID string, args ApplyChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyChainUpdatesWithPackageID exercises the ApplyChainUpdates choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ApplyChainUpdatesWithPackageID(contractID string, packageID string, args ApplyChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// VerifyInboundMessage exercises the VerifyInboundMessage choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) VerifyInboundMessage(contractID string, args VerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// VerifyInboundMessageWithPackageID exercises the VerifyInboundMessage choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) VerifyInboundMessageWithPackageID(contractID string, packageID string, args VerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// VerifyOutboundCCVs exercises the VerifyOutboundCCVs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) VerifyOutboundCCVs(contractID string, args VerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// VerifyOutboundCCVsWithPackageID exercises the VerifyOutboundCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) VerifyOutboundCCVsWithPackageID(contractID string, packageID string, args VerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// CalculateFee exercises the CalculateFee choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) CalculateFee(contractID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CalculateFeeWithPackageID exercises the CalculateFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) CalculateFeeWithPackageID(contractID string, packageID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// GetFee exercises the GetFee choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVs exercises the GetRequiredCCVs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) GetRequiredCCVs(contractID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsWithPackageID exercises the GetRequiredCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) GetRequiredCCVsWithPackageID(contractID string, packageID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "TokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "TokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetDynamicConfig exercises the SetDynamicConfig choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) SetDynamicConfig(contractID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfigWithPackageID exercises the SetDynamicConfig choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) SetDynamicConfigWithPackageID(contractID string, packageID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetPoolReceiveContext exercises the SetPoolReceiveContext choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) SetPoolReceiveContext(contractID string, args SetPoolReceiveContext) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetPoolReceiveContext",
		Arguments:  argsToMap(args),
	}
}

// SetPoolReceiveContextWithPackageID exercises the SetPoolReceiveContext choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) SetPoolReceiveContextWithPackageID(contractID string, packageID string, args SetPoolReceiveContext) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetPoolReceiveContext",
		Arguments:  argsToMap(args),
	}
}

// SetTransferTimeout exercises the SetTransferTimeout choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) SetTransferTimeout(contractID string, args SetTransferTimeout) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetTransferTimeout",
		Arguments:  argsToMap(args),
	}
}

// SetTransferTimeoutWithPackageID exercises the SetTransferTimeout choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) SetTransferTimeoutWithPackageID(contractID string, packageID string, args SetTransferTimeout) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "SetTransferTimeout",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this LockReleaseTokenPool contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetRequiredCCVs exercises the TokenPool_GetRequiredCCVs choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolGetRequiredCCVs(contractID string, args interfaces.TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetRequiredCCVsWithPackageID exercises the TokenPool_GetRequiredCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolGetRequiredCCVsWithPackageID(contractID string, packageID string, args interfaces.TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundMessage exercises the TokenPool_VerifyInboundMessage choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolVerifyInboundMessage(contractID string, args interfaces.TokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundMessageWithPackageID exercises the TokenPool_VerifyInboundMessage choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolVerifyInboundMessageWithPackageID(contractID string, packageID string, args interfaces.TokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyOutboundCCVs exercises the TokenPool_VerifyOutboundCCVs choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolVerifyOutboundCCVs(contractID string, args interfaces.TokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyOutboundCCVsWithPackageID exercises the TokenPool_VerifyOutboundCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolVerifyOutboundCCVsWithPackageID(contractID string, packageID string, args interfaces.TokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicket exercises the TokenPool_ReleaseFromTicket choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolReleaseFromTicket(contractID string, args interfaces.TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicketWithPackageID exercises the TokenPool_ReleaseFromTicket choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolReleaseFromTicketWithPackageID(contractID string, packageID string, args interfaces.TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurn exercises the TokenPool_LockOrBurn choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolLockOrBurn(contractID string, args interfaces.TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurnWithPackageID exercises the TokenPool_LockOrBurn choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolLockOrBurnWithPackageID(contractID string, packageID string, args interfaces.TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFee exercises the TokenPool_CalculateFee choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolCalculateFee(contractID string, args interfaces.TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFeeWithPackageID exercises the TokenPool_CalculateFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolCalculateFeeWithPackageID(contractID string, packageID string, args interfaces.TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetFee exercises the TokenPool_GetFee choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolGetFee(contractID string, args interfaces.TokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetFeeWithPackageID exercises the TokenPool_GetFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolGetFeeWithPackageID(contractID string, packageID string, args interfaces.TokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for LockReleaseTokenPool

var _ mcms.IMCMSReceiver = (*LockReleaseTokenPool)(nil)

var _ interfaces.IITokenPool = (*LockReleaseTokenPool)(nil)

// LockReleaseTokenPoolDeps is a Record type
type LockReleaseTokenPoolDeps struct {
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
	FeeQuoter          mcms.RawInstanceAddress `json:"feeQuoter"`
}

// ToMap converts LockReleaseTokenPoolDeps to a map for DAML arguments
func (t LockReleaseTokenPoolDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	return m
}

func (t LockReleaseTokenPoolDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolDeps to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolDeps from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RateLimitConfigArgs is a Record type
type RateLimitConfigArgs struct {
	RemoteChainSelector                        types.NUMERIC           `json:"remoteChainSelector"`
	InboundRateLimiter                         mcms.RawInstanceAddress `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter mcms.RawInstanceAddress `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        mcms.RawInstanceAddress `json:"outboundRateLimiter"`
}

// ToMap converts RateLimitConfigArgs to a map for DAML arguments
func (t RateLimitConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["inboundRateLimiter"] = model.NestedToDAMLValue(t.InboundRateLimiter)

	m["inboundCustomBlockConfirmationsRateLimiter"] = model.NestedToDAMLValue(t.InboundCustomBlockConfirmationsRateLimiter)

	m["outboundRateLimiter"] = model.NestedToDAMLValue(t.OutboundRateLimiter)

	return m
}

func (t RateLimitConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RateLimitConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RateLimitConfigArgs to hex string (Canton MCMS format)
func (t RateLimitConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RateLimitConfigArgs from hex string (Canton MCMS format)
func (t *RateLimitConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReleaseFromTicket is a Record type
type ReleaseFromTicket struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID     `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	TokenReceiveTicketCid types.CONTRACT_ID     `json:"tokenReceiveTicketCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	Caller                types.PARTY           `json:"caller"`
}

// ToMap converts ReleaseFromTicket to a map for DAML arguments
func (t ReleaseFromTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	m["tokenInput"] = model.NestedToDAMLValue(t.TokenInput)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ReleaseFromTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReleaseFromTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReleaseFromTicket to hex string (Canton MCMS format)
func (t ReleaseFromTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReleaseFromTicket from hex string (Canton MCMS format)
func (t *ReleaseFromTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReleaseFromTicketMCMSParams is ReleaseFromTicket without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ReleaseFromTicketMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID     `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	TokenReceiveTicketCid types.CONTRACT_ID     `json:"tokenReceiveTicketCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
}

// MarshalHex encodes ReleaseFromTicketMCMSParams to hex string for MCMS operationData.
func (t ReleaseFromTicketMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReleaseFromTicketMCMSParams from hex string.
func (t *ReleaseFromTicketMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoteChainConfig is a Record type
type RemoteChainConfig struct {
	RemotePools                                []types.TEXT              `json:"remotePools"`
	RemoteTokenAddress                         types.TEXT                `json:"remoteTokenAddress"`
	InboundCCVs                                []mcms.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []mcms.RawInstanceAddress `json:"outboundCCVs"`
	FinalityConfig                             common.FinalityConfig     `json:"finalityConfig"`
	InboundRateLimiter                         mcms.RawInstanceAddress   `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter mcms.RawInstanceAddress   `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        mcms.RawInstanceAddress   `json:"outboundRateLimiter"`
}

// ToMap converts RemoteChainConfig to a map for DAML arguments
func (t RemoteChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["remotePools"] = func() []any {
		res := make([]any, 0, len(t.RemotePools))
		for _, e := range t.RemotePools {
			res = append(res, string(e))
		}
		return res
	}()

	m["remoteTokenAddress"] = string(t.RemoteTokenAddress)

	m["inboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.InboundCCVs))
		for _, e := range t.InboundCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["outboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.OutboundCCVs))
		for _, e := range t.OutboundCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["finalityConfig"] = model.NestedToDAMLValue(t.FinalityConfig)

	m["inboundRateLimiter"] = model.NestedToDAMLValue(t.InboundRateLimiter)

	m["inboundCustomBlockConfirmationsRateLimiter"] = model.NestedToDAMLValue(t.InboundCustomBlockConfirmationsRateLimiter)

	m["outboundRateLimiter"] = model.NestedToDAMLValue(t.OutboundRateLimiter)

	return m
}

func (t RemoteChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoteChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoteChainConfig to hex string (Canton MCMS format)
func (t RemoteChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoteChainConfig from hex string (Canton MCMS format)
func (t *RemoteChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDynamicConfig is a Record type
type SetDynamicConfig struct {
	RateLimitAdmin *types.PARTY `json:"rateLimitAdmin" hex:"optional"`
}

// ToMap converts SetDynamicConfig to a map for DAML arguments
func (t SetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RateLimitAdmin != nil {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfig to hex string (Canton MCMS format)
func (t SetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfig from hex string (Canton MCMS format)
func (t *SetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDynamicConfigParams is a Record type
type SetDynamicConfigParams struct {
	RateLimitAdmin *types.PARTY `json:"rateLimitAdmin" hex:"optional"`
}

// ToMap converts SetDynamicConfigParams to a map for DAML arguments
func (t SetDynamicConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RateLimitAdmin != nil {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetDynamicConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfigParams to hex string (Canton MCMS format)
func (t SetDynamicConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfigParams from hex string (Canton MCMS format)
func (t *SetDynamicConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPoolReceiveContext is a Record type
type SetPoolReceiveContext struct {
	NewPoolReceiveContext common.CCIPContext `json:"newPoolReceiveContext"`
}

// ToMap converts SetPoolReceiveContext to a map for DAML arguments
func (t SetPoolReceiveContext) ToMap() map[string]any {
	m := make(map[string]any)

	m["newPoolReceiveContext"] = model.NestedToDAMLValue(t.NewPoolReceiveContext)

	return m
}

func (t SetPoolReceiveContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPoolReceiveContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPoolReceiveContext to hex string (Canton MCMS format)
func (t SetPoolReceiveContext) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPoolReceiveContext from hex string (Canton MCMS format)
func (t *SetPoolReceiveContext) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRateLimitConfig is a Record type
type SetRateLimitConfig struct {
	RateLimitConfigArgs []RateLimitConfigArgs `json:"rateLimitConfigArgs"`
	Caller              types.PARTY           `json:"caller"`
}

// ToMap converts SetRateLimitConfig to a map for DAML arguments
func (t SetRateLimitConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimitConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RateLimitConfigArgs))
		for _, e := range t.RateLimitConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetRateLimitConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRateLimitConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRateLimitConfig to hex string (Canton MCMS format)
func (t SetRateLimitConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRateLimitConfig from hex string (Canton MCMS format)
func (t *SetRateLimitConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRateLimitConfigMCMSParams is SetRateLimitConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetRateLimitConfigMCMSParams struct {
	RateLimitConfigArgs []RateLimitConfigArgs `json:"rateLimitConfigArgs"`
}

// MarshalHex encodes SetRateLimitConfigMCMSParams to hex string for MCMS operationData.
func (t SetRateLimitConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRateLimitConfigMCMSParams from hex string.
func (t *SetRateLimitConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRateLimitConfigParams is a Record type
type SetRateLimitConfigParams struct {
	RateLimitConfigArgs []RateLimitConfigArgs `json:"rateLimitConfigArgs"`
	Caller              types.PARTY           `json:"caller"`
}

// ToMap converts SetRateLimitConfigParams to a map for DAML arguments
func (t SetRateLimitConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimitConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RateLimitConfigArgs))
		for _, e := range t.RateLimitConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetRateLimitConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRateLimitConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRateLimitConfigParams to hex string (Canton MCMS format)
func (t SetRateLimitConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRateLimitConfigParams from hex string (Canton MCMS format)
func (t *SetRateLimitConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferTimeout is a Record type
type SetTransferTimeout struct {
	NewTransferTimeout TransferTimeout `json:"newTransferTimeout"`
}

// ToMap converts SetTransferTimeout to a map for DAML arguments
func (t SetTransferTimeout) ToMap() map[string]any {
	m := make(map[string]any)

	m["newTransferTimeout"] = model.NestedToDAMLValue(t.NewTransferTimeout)

	return m
}

func (t SetTransferTimeout) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetTransferTimeout) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetTransferTimeout to hex string (Canton MCMS format)
func (t SetTransferTimeout) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferTimeout from hex string (Canton MCMS format)
func (t *SetTransferTimeout) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferTimeoutParams is a Record type
type SetTransferTimeoutParams struct {
	TransferTimeout TransferTimeout `json:"transferTimeout"`
}

// ToMap converts SetTransferTimeoutParams to a map for DAML arguments
func (t SetTransferTimeoutParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferTimeout"] = model.NestedToDAMLValue(t.TransferTimeout)

	return m
}

func (t SetTransferTimeoutParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetTransferTimeoutParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetTransferTimeoutParams to hex string (Canton MCMS format)
func (t SetTransferTimeoutParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferTimeoutParams from hex string (Canton MCMS format)
func (t *SetTransferTimeoutParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferFeeConfig2 is a Record type
type TokenTransferFeeConfig2 struct {
	IsEnabled         types.BOOL    `json:"isEnabled"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	FeeBps            types.NUMERIC `json:"feeBps"`
}

// ToMap converts TokenTransferFeeConfig2 to a map for DAML arguments
func (t TokenTransferFeeConfig2) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeUSDCents"] = t.FeeUSDCents

	m["feeBps"] = t.FeeBps

	return m
}

func (t TokenTransferFeeConfig2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferFeeConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferFeeConfig2 to hex string (Canton MCMS format)
func (t TokenTransferFeeConfig2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferFeeConfig2 from hex string (Canton MCMS format)
func (t *TokenTransferFeeConfig2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferFeeConfigArgs is a Record type
type TokenTransferFeeConfigArgs struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	IsEnabled         types.BOOL    `json:"isEnabled"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	FeeBps            types.NUMERIC `json:"feeBps"`
}

// ToMap converts TokenTransferFeeConfigArgs to a map for DAML arguments
func (t TokenTransferFeeConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["isEnabled"] = bool(t.IsEnabled)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeUSDCents"] = t.FeeUSDCents

	m["feeBps"] = t.FeeBps

	return m
}

func (t TokenTransferFeeConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferFeeConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferFeeConfigArgs to hex string (Canton MCMS format)
func (t TokenTransferFeeConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferFeeConfigArgs from hex string (Canton MCMS format)
func (t *TokenTransferFeeConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferTimeout is a variant/union type
type TransferTimeout struct {
	Indefinite    *types.UNIT  `json:"Indefinite,omitempty"`
	RelativeHours *types.INT64 `json:"RelativeHours,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferTimeout
func (v TransferTimeout) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for TransferTimeout
func (v *TransferTimeout) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes TransferTimeout to hex string (Canton MCMS format)
func (v TransferTimeout) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes TransferTimeout from hex string (Canton MCMS format)
func (v *TransferTimeout) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferTimeout) GetVariantTag() string {

	if v.Indefinite != nil {
		return "Indefinite"
	}

	if v.RelativeHours != nil {
		return "RelativeHours"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferTimeout) GetVariantValue() any {

	if v.Indefinite != nil {
		return v.Indefinite
	}

	if v.RelativeHours != nil {
		return v.RelativeHours
	}

	return nil
}

var _ types.VARIANT = (*TransferTimeout)(nil)

// GetVariantTagByte implements types.VariantWithTagByte interface for MCMS numeric tag encoding
func (v TransferTimeout) GetVariantTagByte() byte {

	if v.Indefinite != nil {
		return 0
	}

	if v.RelativeHours != nil {
		return 1
	}

	return 0xFF // Invalid/unknown variant
}

var _ types.VariantWithTagByte = (*TransferTimeout)(nil)

// VerifyInboundMessage is a Record type
type VerifyInboundMessage struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID  `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid   types.CONTRACT_ID  `json:"executingMessageCid"`
	Caller                types.PARTY        `json:"caller"`
}

// ToMap converts VerifyInboundMessage to a map for DAML arguments
func (t VerifyInboundMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t VerifyInboundMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *VerifyInboundMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes VerifyInboundMessage to hex string (Canton MCMS format)
func (t VerifyInboundMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyInboundMessage from hex string (Canton MCMS format)
func (t *VerifyInboundMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifyInboundMessageMCMSParams is VerifyInboundMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type VerifyInboundMessageMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID  `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid   types.CONTRACT_ID  `json:"executingMessageCid"`
}

// MarshalHex encodes VerifyInboundMessageMCMSParams to hex string for MCMS operationData.
func (t VerifyInboundMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyInboundMessageMCMSParams from hex string.
func (t *VerifyInboundMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifyOutboundCCVs is a Record type
type VerifyOutboundCCVs struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID  `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID  `json:"sendingMessageCid"`
	Amount                types.NUMERIC      `json:"amount"`
	Caller                types.PARTY        `json:"caller"`
}

// ToMap converts VerifyOutboundCCVs to a map for DAML arguments
func (t VerifyOutboundCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["amount"] = t.Amount

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t VerifyOutboundCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *VerifyOutboundCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes VerifyOutboundCCVs to hex string (Canton MCMS format)
func (t VerifyOutboundCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyOutboundCCVs from hex string (Canton MCMS format)
func (t *VerifyOutboundCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifyOutboundCCVsMCMSParams is VerifyOutboundCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type VerifyOutboundCCVsMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID  `json:"tokenConfigCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID  `json:"sendingMessageCid"`
	Amount                types.NUMERIC      `json:"amount"`
}

// MarshalHex encodes VerifyOutboundCCVsMCMSParams to hex string for MCMS operationData.
func (t VerifyOutboundCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyOutboundCCVsMCMSParams from hex string.
func (t *VerifyOutboundCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ApplyChainUpdates(args ApplyChainUpdates) (*bind.EncodedChoice, error)
	ApplyChainUpdatesParams(args ApplyChainUpdatesParams) (*bind.EncodedChoice, error)
	ApplyTokenTransferFeeConfigUpdates(args ApplyTokenTransferFeeConfigUpdates) (*bind.EncodedChoice, error)
	ApplyTokenTransferFeeConfigUpdatesParams(args ApplyTokenTransferFeeConfigUpdatesParams) (*bind.EncodedChoice, error)
	CalculateFee(args CalculateFee) (*bind.EncodedChoice, error)
	CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error)
	GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error)
	GetRequiredCCVsMCMSParams(args GetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error)
	LockOrBurn(args LockOrBurn) (*bind.EncodedChoice, error)
	LockOrBurnMCMSParams(args LockOrBurnMCMSParams) (*bind.EncodedChoice, error)
	ReleaseFromTicket(args ReleaseFromTicket) (*bind.EncodedChoice, error)
	ReleaseFromTicketMCMSParams(args ReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error)
	SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error)
	SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error)
	SetPoolReceiveContext(args SetPoolReceiveContext) (*bind.EncodedChoice, error)
	SetRateLimitConfig(args SetRateLimitConfig) (*bind.EncodedChoice, error)
	SetRateLimitConfigMCMSParams(args SetRateLimitConfigMCMSParams) (*bind.EncodedChoice, error)
	SetRateLimitConfigParams(args SetRateLimitConfigParams) (*bind.EncodedChoice, error)
	SetTransferTimeout(args SetTransferTimeout) (*bind.EncodedChoice, error)
	SetTransferTimeoutParams(args SetTransferTimeoutParams) (*bind.EncodedChoice, error)
	VerifyInboundMessage(args VerifyInboundMessage) (*bind.EncodedChoice, error)
	VerifyInboundMessageMCMSParams(args VerifyInboundMessageMCMSParams) (*bind.EncodedChoice, error)
	VerifyOutboundCCVs(args VerifyOutboundCCVs) (*bind.EncodedChoice, error)
	VerifyOutboundCCVsMCMSParams(args VerifyOutboundCCVsMCMSParams) (*bind.EncodedChoice, error)
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

// ApplyChainUpdates encodes parameters for the ApplyChainUpdates choice.
func (e *encoder) ApplyChainUpdates(args ApplyChainUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyChainUpdates", args)
}

// ApplyChainUpdatesParams encodes parameters for the ApplyChainUpdates choice.
func (e *encoder) ApplyChainUpdatesParams(args ApplyChainUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyChainUpdates", args)
}

// ApplyTokenTransferFeeConfigUpdates encodes parameters for the ApplyTokenTransferFeeConfigUpdates choice.
func (e *encoder) ApplyTokenTransferFeeConfigUpdates(args ApplyTokenTransferFeeConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyTokenTransferFeeConfigUpdates", args)
}

// ApplyTokenTransferFeeConfigUpdatesParams encodes parameters for the ApplyTokenTransferFeeConfigUpdates choice.
func (e *encoder) ApplyTokenTransferFeeConfigUpdatesParams(args ApplyTokenTransferFeeConfigUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyTokenTransferFeeConfigUpdates", args)
}

// CalculateFee encodes parameters for the CalculateFee choice.
func (e *encoder) CalculateFee(args CalculateFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// CalculateFeeMCMSParams encodes MCMS parameters (without Caller) for the CalculateFee choice.
func (e *encoder) CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// GetFee encodes parameters for the GetFee choice.
func (e *encoder) GetFee(args GetFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetFeeMCMSParams encodes MCMS parameters (without Caller) for the GetFee choice.
func (e *encoder) GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetRequiredCCVs encodes parameters for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// GetRequiredCCVsMCMSParams encodes MCMS parameters (without Caller) for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVsMCMSParams(args GetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// LockOrBurn encodes parameters for the LockOrBurn choice.
func (e *encoder) LockOrBurn(args LockOrBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOrBurn", args)
}

// LockOrBurnMCMSParams encodes MCMS parameters (without Caller) for the LockOrBurn choice.
func (e *encoder) LockOrBurnMCMSParams(args LockOrBurnMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOrBurn", args)
}

// ReleaseFromTicket encodes parameters for the ReleaseFromTicket choice.
func (e *encoder) ReleaseFromTicket(args ReleaseFromTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ReleaseFromTicket", args)
}

// ReleaseFromTicketMCMSParams encodes MCMS parameters (without Caller) for the ReleaseFromTicket choice.
func (e *encoder) ReleaseFromTicketMCMSParams(args ReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ReleaseFromTicket", args)
}

// SetDynamicConfig encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// SetDynamicConfigParams encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// SetPoolReceiveContext encodes parameters for the SetPoolReceiveContext choice.
func (e *encoder) SetPoolReceiveContext(args SetPoolReceiveContext) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPoolReceiveContext", args)
}

// SetRateLimitConfig encodes parameters for the SetRateLimitConfig choice.
func (e *encoder) SetRateLimitConfig(args SetRateLimitConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRateLimitConfig", args)
}

// SetRateLimitConfigMCMSParams encodes MCMS parameters (without Caller) for the SetRateLimitConfig choice.
func (e *encoder) SetRateLimitConfigMCMSParams(args SetRateLimitConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRateLimitConfig", args)
}

// SetRateLimitConfigParams encodes parameters for the SetRateLimitConfig choice.
func (e *encoder) SetRateLimitConfigParams(args SetRateLimitConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRateLimitConfig", args)
}

// SetTransferTimeout encodes parameters for the SetTransferTimeout choice.
func (e *encoder) SetTransferTimeout(args SetTransferTimeout) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferTimeout", args)
}

// SetTransferTimeoutParams encodes parameters for the SetTransferTimeout choice.
func (e *encoder) SetTransferTimeoutParams(args SetTransferTimeoutParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferTimeout", args)
}

// VerifyInboundMessage encodes parameters for the VerifyInboundMessage choice.
func (e *encoder) VerifyInboundMessage(args VerifyInboundMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyInboundMessage", args)
}

// VerifyInboundMessageMCMSParams encodes MCMS parameters (without Caller) for the VerifyInboundMessage choice.
func (e *encoder) VerifyInboundMessageMCMSParams(args VerifyInboundMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyInboundMessage", args)
}

// VerifyOutboundCCVs encodes parameters for the VerifyOutboundCCVs choice.
func (e *encoder) VerifyOutboundCCVs(args VerifyOutboundCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyOutboundCCVs", args)
}

// VerifyOutboundCCVsMCMSParams encodes MCMS parameters (without Caller) for the VerifyOutboundCCVs choice.
func (e *encoder) VerifyOutboundCCVsMCMSParams(args VerifyOutboundCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyOutboundCCVs", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
