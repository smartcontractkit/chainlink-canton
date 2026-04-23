package feequoter

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
	PackageName = "ccip-feequoter"
	PackageID   = "9c27fde3dc7b96748debeed08e37e8e3dcc73ae8c6df9cc53f24b4694fcb59f3"
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

// AddPriceUpdaters is a Record type
type AddPriceUpdaters struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddPriceUpdaters to a map for DAML arguments
func (t AddPriceUpdaters) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddPriceUpdaters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPriceUpdaters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPriceUpdaters to hex string (Canton MCMS format)
func (t AddPriceUpdaters) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPriceUpdaters from hex string (Canton MCMS format)
func (t *AddPriceUpdaters) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyDestChainConfigUpdates2 is a Record type
type ApplyDestChainConfigUpdates2 struct {
	DestChainConfigArgs []DestChainConfigArgs2 `json:"destChainConfigArgs"`
}

// ToMap converts ApplyDestChainConfigUpdates2 to a map for DAML arguments
func (t ApplyDestChainConfigUpdates2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyDestChainConfigUpdates2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainConfigUpdates2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainConfigUpdates2 to hex string (Canton MCMS format)
func (t ApplyDestChainConfigUpdates2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainConfigUpdates2 from hex string (Canton MCMS format)
func (t *ApplyDestChainConfigUpdates2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyDestChainConfigUpdatesParams2 is a Record type
type ApplyDestChainConfigUpdatesParams2 struct {
	DestChainConfigArgs []DestChainConfigArgs2 `json:"destChainConfigArgs"`
}

// ToMap converts ApplyDestChainConfigUpdatesParams2 to a map for DAML arguments
func (t ApplyDestChainConfigUpdatesParams2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyDestChainConfigUpdatesParams2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainConfigUpdatesParams2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainConfigUpdatesParams2 to hex string (Canton MCMS format)
func (t ApplyDestChainConfigUpdatesParams2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainConfigUpdatesParams2 from hex string (Canton MCMS format)
func (t *ApplyDestChainConfigUpdatesParams2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyPriceUpdatersUpdate is a Record type
type ApplyPriceUpdatersUpdate struct {
	AddedPriceUpdaters   []types.PARTY `json:"addedPriceUpdaters"`
	RemovedPriceUpdaters []types.PARTY `json:"removedPriceUpdaters"`
}

// ToMap converts ApplyPriceUpdatersUpdate to a map for DAML arguments
func (t ApplyPriceUpdatersUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["addedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.AddedPriceUpdaters))
		for _, e := range t.AddedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["removedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.RemovedPriceUpdaters))
		for _, e := range t.RemovedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t ApplyPriceUpdatersUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyPriceUpdatersUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyPriceUpdatersUpdate to hex string (Canton MCMS format)
func (t ApplyPriceUpdatersUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyPriceUpdatersUpdate from hex string (Canton MCMS format)
func (t *ApplyPriceUpdatersUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyPriceUpdatersUpdateParams is a Record type
type ApplyPriceUpdatersUpdateParams struct {
	AddedPriceUpdaters   []types.PARTY `json:"addedPriceUpdaters"`
	RemovedPriceUpdaters []types.PARTY `json:"removedPriceUpdaters"`
}

// ToMap converts ApplyPriceUpdatersUpdateParams to a map for DAML arguments
func (t ApplyPriceUpdatersUpdateParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["addedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.AddedPriceUpdaters))
		for _, e := range t.AddedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["removedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.RemovedPriceUpdaters))
		for _, e := range t.RemovedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t ApplyPriceUpdatersUpdateParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyPriceUpdatersUpdateParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyPriceUpdatersUpdateParams to hex string (Canton MCMS format)
func (t ApplyPriceUpdatersUpdateParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyPriceUpdatersUpdateParams from hex string (Canton MCMS format)
func (t *ApplyPriceUpdatersUpdateParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DestChainConfig2 is a Record type
type DestChainConfig2 struct {
	IsEnabled                   types.BOOL    `json:"isEnabled"`
	MaxDataBytes                types.INT64   `json:"maxDataBytes"`
	MaxPerMsgGasLimit           types.INT64   `json:"maxPerMsgGasLimit"`
	DestGasOverhead             types.INT64   `json:"destGasOverhead"`
	DestGasPerPayloadByteBase   types.INT64   `json:"destGasPerPayloadByteBase"`
	DefaultTxGasLimit           types.INT64   `json:"defaultTxGasLimit"`
	LinkFeeMultiplierPercent    types.NUMERIC `json:"linkFeeMultiplierPercent"`
	DefaultTokenFeeUSD          types.NUMERIC `json:"defaultTokenFeeUSD"`
	DefaultTokenDestGasOverhead types.INT64   `json:"defaultTokenDestGasOverhead"`
}

// ToMap converts DestChainConfig2 to a map for DAML arguments
func (t DestChainConfig2) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["maxDataBytes"] = int64(t.MaxDataBytes)

	m["maxPerMsgGasLimit"] = int64(t.MaxPerMsgGasLimit)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destGasPerPayloadByteBase"] = int64(t.DestGasPerPayloadByteBase)

	m["defaultTxGasLimit"] = int64(t.DefaultTxGasLimit)

	m["linkFeeMultiplierPercent"] = t.LinkFeeMultiplierPercent

	m["defaultTokenFeeUSD"] = t.DefaultTokenFeeUSD

	m["defaultTokenDestGasOverhead"] = int64(t.DefaultTokenDestGasOverhead)

	return m
}

func (t DestChainConfig2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DestChainConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DestChainConfig2 to hex string (Canton MCMS format)
func (t DestChainConfig2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DestChainConfig2 from hex string (Canton MCMS format)
func (t *DestChainConfig2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DestChainConfigArgs2 is a Record type
type DestChainConfigArgs2 struct {
	DestChainSelector types.NUMERIC    `json:"destChainSelector"`
	DestChainConfig   DestChainConfig2 `json:"destChainConfig"`
}

// ToMap converts DestChainConfigArgs2 to a map for DAML arguments
func (t DestChainConfigArgs2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["destChainConfig"] = model.NestedToDAMLValue(t.DestChainConfig)

	return m
}

func (t DestChainConfigArgs2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DestChainConfigArgs2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DestChainConfigArgs2 to hex string (Canton MCMS format)
func (t DestChainConfigArgs2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DestChainConfigArgs2 from hex string (Canton MCMS format)
func (t *DestChainConfigArgs2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoter is a Template type
type FeeQuoter struct {
	InstanceId                       types.TEXT                               `json:"instanceId"`
	Owner                            types.PARTY                              `json:"owner"`
	FeeTokens                        types.SET                                `json:"feeTokens"`
	DestChainConfigs                 map[types.NUMERIC]DestChainConfig2       `json:"destChainConfigs"`
	TokenTransferFeeConfigs          map[types.NUMERIC]types.GENMAP           `json:"tokenTransferFeeConfigs"`
	UsdPerUnitGasByDestChainSelector map[types.NUMERIC]TimestampedPrice       `json:"usdPerUnitGasByDestChainSelector"`
	UsdPerToken                      types.GENMAP                             `json:"usdPerToken"`
	LinkTokenInstrumentId            splice_api_token_holding_v1.InstrumentId `json:"linkTokenInstrumentId"`
	PriceUpdaters                    []types.PARTY                            `json:"priceUpdaters"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FeeQuoter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FeeQuoter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FeeQuoter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeTokens"] = model.NestedToDAMLValue(t.FeeTokens)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["linkTokenInstrumentId"] = model.NestedToDAMLValue(t.LinkTokenInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["priceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.PriceUpdaters))
		for _, e := range t.PriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FeeQuoter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeTokens"] = model.NestedToDAMLValue(t.FeeTokens)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["linkTokenInstrumentId"] = model.NestedToDAMLValue(t.LinkTokenInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["priceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.PriceUpdaters))
		for _, e := range t.PriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FeeQuoter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoter to hex string (Canton MCMS format)
func (t FeeQuoter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoter from hex string (Canton MCMS format)
func (t *FeeQuoter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FeeQuoter

// QuoteGasForExec exercises the QuoteGasForExec choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) QuoteGasForExec(contractID string, args QuoteGasForExec) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "QuoteGasForExec",
		Arguments:  argsToMap(args),
	}
}

// QuoteGasForExecWithPackageID exercises the QuoteGasForExec choice using the provided package ID instead of package name
func (t FeeQuoter) QuoteGasForExecWithPackageID(contractID string, packageID string, args QuoteGasForExec) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "QuoteGasForExec",
		Arguments:  argsToMap(args),
	}
}

// GetTokenPrice exercises the GetTokenPrice choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetTokenPrice(contractID string, args GetTokenPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenPrice",
		Arguments:  argsToMap(args),
	}
}

// GetTokenPriceWithPackageID exercises the GetTokenPrice choice using the provided package ID instead of package name
func (t FeeQuoter) GetTokenPriceWithPackageID(contractID string, packageID string, args GetTokenPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenPrice",
		Arguments:  argsToMap(args),
	}
}

// GetDestinationChainGasPrice exercises the GetDestinationChainGasPrice choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetDestinationChainGasPrice(contractID string, args GetDestinationChainGasPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestinationChainGasPrice",
		Arguments:  argsToMap(args),
	}
}

// GetDestinationChainGasPriceWithPackageID exercises the GetDestinationChainGasPrice choice using the provided package ID instead of package name
func (t FeeQuoter) GetDestinationChainGasPriceWithPackageID(contractID string, packageID string, args GetDestinationChainGasPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestinationChainGasPrice",
		Arguments:  argsToMap(args),
	}
}

// UpdatePrices exercises the UpdatePrices choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) UpdatePrices(contractID string, args UpdatePrices) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "UpdatePrices",
		Arguments:  argsToMap(args),
	}
}

// UpdatePricesWithPackageID exercises the UpdatePrices choice using the provided package ID instead of package name
func (t FeeQuoter) UpdatePricesWithPackageID(contractID string, packageID string, args UpdatePrices) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "UpdatePrices",
		Arguments:  argsToMap(args),
	}
}

// GetTokenTransferFee exercises the GetTokenTransferFee choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetTokenTransferFee(contractID string, args GetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// GetTokenTransferFeeWithPackageID exercises the GetTokenTransferFee choice using the provided package ID instead of package name
func (t FeeQuoter) GetTokenTransferFeeWithPackageID(contractID string, packageID string, args GetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfig exercises the GetDestChainConfig choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetDestChainConfig(contractID string, args GetDestChainConfig2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfigWithPackageID exercises the GetDestChainConfig choice using the provided package ID instead of package name
func (t FeeQuoter) GetDestChainConfigWithPackageID(contractID string, packageID string, args GetDestChainConfig2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainConfigUpdates exercises the ApplyDestChainConfigUpdates choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyDestChainConfigUpdates(contractID string, args ApplyDestChainConfigUpdates2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainConfigUpdatesWithPackageID exercises the ApplyDestChainConfigUpdates choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyDestChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyDestChainConfigUpdates2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FeeQuoter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeQuoter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FeeQuoter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) Get(contractID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t FeeQuoter) GetWithPackageID(contractID string, packageID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetFeeTokens exercises the GetFeeTokens choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetFeeTokens(contractID string, args GetFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// GetFeeTokensWithPackageID exercises the GetFeeTokens choice using the provided package ID instead of package name
func (t FeeQuoter) GetFeeTokensWithPackageID(contractID string, packageID string, args GetFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// ApplyPriceUpdatersUpdate exercises the ApplyPriceUpdatersUpdate choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyPriceUpdatersUpdate(contractID string, args ApplyPriceUpdatersUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyPriceUpdatersUpdate",
		Arguments:  argsToMap(args),
	}
}

// ApplyPriceUpdatersUpdateWithPackageID exercises the ApplyPriceUpdatersUpdate choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyPriceUpdatersUpdateWithPackageID(contractID string, packageID string, args ApplyPriceUpdatersUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyPriceUpdatersUpdate",
		Arguments:  argsToMap(args),
	}
}

// AddPriceUpdaters exercises the AddPriceUpdaters choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) AddPriceUpdaters(contractID string, args AddPriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "AddPriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// AddPriceUpdatersWithPackageID exercises the AddPriceUpdaters choice using the provided package ID instead of package name
func (t FeeQuoter) AddPriceUpdatersWithPackageID(contractID string, packageID string, args AddPriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "AddPriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemovePriceUpdaters exercises the RemovePriceUpdaters choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) RemovePriceUpdaters(contractID string, args RemovePriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemovePriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemovePriceUpdatersWithPackageID exercises the RemovePriceUpdaters choice using the provided package ID instead of package name
func (t FeeQuoter) RemovePriceUpdatersWithPackageID(contractID string, packageID string, args RemovePriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemovePriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemoveFeeTokens exercises the RemoveFeeTokens choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) RemoveFeeTokens(contractID string, args RemoveFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemoveFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// RemoveFeeTokensWithPackageID exercises the RemoveFeeTokens choice using the provided package ID instead of package name
func (t FeeQuoter) RemoveFeeTokensWithPackageID(contractID string, packageID string, args RemoveFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemoveFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this FeeQuoter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeQuoter) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t FeeQuoter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for FeeQuoter

var _ mcms.IMCMSReceiver = (*FeeQuoter)(nil)

// GasPriceUpdate is a Record type
type GasPriceUpdate struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	UsdPerUnitGas     types.NUMERIC `json:"usdPerUnitGas"`
}

// ToMap converts GasPriceUpdate to a map for DAML arguments
func (t GasPriceUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["usdPerUnitGas"] = t.UsdPerUnitGas

	return m
}

func (t GasPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GasPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GasPriceUpdate to hex string (Canton MCMS format)
func (t GasPriceUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GasPriceUpdate from hex string (Canton MCMS format)
func (t *GasPriceUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Get is a Record type
type Get struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts Get to a map for DAML arguments
func (t Get) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Get) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Get) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Get to hex string (Canton MCMS format)
func (t Get) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get from hex string (Canton MCMS format)
func (t *Get) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMCMSParams is Get without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetMCMSParams struct {
}

// MarshalHex encodes GetMCMSParams to hex string for MCMS operationData.
func (t GetMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMCMSParams from hex string.
func (t *GetMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestChainConfig2 is a Record type
type GetDestChainConfig2 struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts GetDestChainConfig2 to a map for DAML arguments
func (t GetDestChainConfig2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDestChainConfig2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDestChainConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDestChainConfig2 to hex string (Canton MCMS format)
func (t GetDestChainConfig2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChainConfig2 from hex string (Canton MCMS format)
func (t *GetDestChainConfig2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestChainConfig2MCMSParams is GetDestChainConfig2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDestChainConfig2MCMSParams struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// MarshalHex encodes GetDestChainConfig2MCMSParams to hex string for MCMS operationData.
func (t GetDestChainConfig2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChainConfig2MCMSParams from hex string.
func (t *GetDestChainConfig2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestinationChainGasPrice is a Record type
type GetDestinationChainGasPrice struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts GetDestinationChainGasPrice to a map for DAML arguments
func (t GetDestinationChainGasPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDestinationChainGasPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDestinationChainGasPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDestinationChainGasPrice to hex string (Canton MCMS format)
func (t GetDestinationChainGasPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestinationChainGasPrice from hex string (Canton MCMS format)
func (t *GetDestinationChainGasPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestinationChainGasPriceMCMSParams is GetDestinationChainGasPrice without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDestinationChainGasPriceMCMSParams struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// MarshalHex encodes GetDestinationChainGasPriceMCMSParams to hex string for MCMS operationData.
func (t GetDestinationChainGasPriceMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestinationChainGasPriceMCMSParams from hex string.
func (t *GetDestinationChainGasPriceMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeTokens is a Record type
type GetFeeTokens struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetFeeTokens to a map for DAML arguments
func (t GetFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeTokens to hex string (Canton MCMS format)
func (t GetFeeTokens) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeTokens from hex string (Canton MCMS format)
func (t *GetFeeTokens) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeTokensMCMSParams is GetFeeTokens without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFeeTokensMCMSParams struct {
}

// MarshalHex encodes GetFeeTokensMCMSParams to hex string for MCMS operationData.
func (t GetFeeTokensMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeTokensMCMSParams from hex string.
func (t *GetFeeTokensMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenPrice is a Record type
type GetTokenPrice struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenPrice to a map for DAML arguments
func (t GetTokenPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenPrice to hex string (Canton MCMS format)
func (t GetTokenPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenPrice from hex string (Canton MCMS format)
func (t *GetTokenPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenPriceMCMSParams is GetTokenPrice without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenPriceMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes GetTokenPriceMCMSParams to hex string for MCMS operationData.
func (t GetTokenPriceMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenPriceMCMSParams from hex string.
func (t *GetTokenPriceMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenTransferFee is a Record type
type GetTokenTransferFee struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	Token             splice_api_token_holding_v1.InstrumentId `json:"token"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenTransferFee to a map for DAML arguments
func (t GetTokenTransferFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["token"] = model.NestedToDAMLValue(t.Token)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenTransferFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenTransferFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenTransferFee to hex string (Canton MCMS format)
func (t GetTokenTransferFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenTransferFee from hex string (Canton MCMS format)
func (t *GetTokenTransferFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenTransferFeeMCMSParams is GetTokenTransferFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenTransferFeeMCMSParams struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	Token             splice_api_token_holding_v1.InstrumentId `json:"token"`
}

// MarshalHex encodes GetTokenTransferFeeMCMSParams to hex string for MCMS operationData.
func (t GetTokenTransferFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenTransferFeeMCMSParams from hex string.
func (t *GetTokenTransferFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PriceUpdates is a Record type
type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate `json:"tokenPriceUpdates"`
	GasPriceUpdates   []GasPriceUpdate   `json:"gasPriceUpdates"`
}

// ToMap converts PriceUpdates to a map for DAML arguments
func (t PriceUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenPriceUpdates"] = func() []any {
		res := make([]any, 0, len(t.TokenPriceUpdates))
		for _, e := range t.TokenPriceUpdates {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["gasPriceUpdates"] = func() []any {
		res := make([]any, 0, len(t.GasPriceUpdates))
		for _, e := range t.GasPriceUpdates {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t PriceUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PriceUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PriceUpdates to hex string (Canton MCMS format)
func (t PriceUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PriceUpdates from hex string (Canton MCMS format)
func (t *PriceUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExec is a Record type
type QuoteGasForExec struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	NonCalldataGas    types.INT64                              `json:"nonCalldataGas"`
	CalldataSize      types.INT64                              `json:"calldataSize"`
	FeeToken          splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts QuoteGasForExec to a map for DAML arguments
func (t QuoteGasForExec) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["nonCalldataGas"] = int64(t.NonCalldataGas)

	m["calldataSize"] = int64(t.CalldataSize)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t QuoteGasForExec) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *QuoteGasForExec) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes QuoteGasForExec to hex string (Canton MCMS format)
func (t QuoteGasForExec) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExec from hex string (Canton MCMS format)
func (t *QuoteGasForExec) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExecMCMSParams is QuoteGasForExec without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type QuoteGasForExecMCMSParams struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	NonCalldataGas    types.INT64                              `json:"nonCalldataGas"`
	CalldataSize      types.INT64                              `json:"calldataSize"`
	FeeToken          splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
}

// MarshalHex encodes QuoteGasForExecMCMSParams to hex string for MCMS operationData.
func (t QuoteGasForExecMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExecMCMSParams from hex string.
func (t *QuoteGasForExecMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExecResult is a Record type
type QuoteGasForExecResult struct {
	TotalGas          types.INT64   `json:"totalGas"`
	GasCostUSDCents   types.NUMERIC `json:"gasCostUSDCents"`
	FeeTokenPrice     types.NUMERIC `json:"feeTokenPrice"`
	PremiumMultiplier types.NUMERIC `json:"premiumMultiplier"`
}

// ToMap converts QuoteGasForExecResult to a map for DAML arguments
func (t QuoteGasForExecResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["totalGas"] = int64(t.TotalGas)

	m["gasCostUSDCents"] = t.GasCostUSDCents

	m["feeTokenPrice"] = t.FeeTokenPrice

	m["premiumMultiplier"] = t.PremiumMultiplier

	return m
}

func (t QuoteGasForExecResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *QuoteGasForExecResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes QuoteGasForExecResult to hex string (Canton MCMS format)
func (t QuoteGasForExecResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExecResult from hex string (Canton MCMS format)
func (t *QuoteGasForExecResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveFeeTokens is a Record type
type RemoveFeeTokens struct {
	FeeTokensToRemove []splice_api_token_holding_v1.InstrumentId `json:"feeTokensToRemove"`
}

// ToMap converts RemoveFeeTokens to a map for DAML arguments
func (t RemoveFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokensToRemove"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToRemove))
		for _, e := range t.FeeTokensToRemove {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RemoveFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveFeeTokens to hex string (Canton MCMS format)
func (t RemoveFeeTokens) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveFeeTokens from hex string (Canton MCMS format)
func (t *RemoveFeeTokens) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveFeeTokensParams is a Record type
type RemoveFeeTokensParams struct {
	FeeTokensToRemove []splice_api_token_holding_v1.InstrumentId `json:"feeTokensToRemove"`
}

// ToMap converts RemoveFeeTokensParams to a map for DAML arguments
func (t RemoveFeeTokensParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokensToRemove"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToRemove))
		for _, e := range t.FeeTokensToRemove {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RemoveFeeTokensParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveFeeTokensParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveFeeTokensParams to hex string (Canton MCMS format)
func (t RemoveFeeTokensParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveFeeTokensParams from hex string (Canton MCMS format)
func (t *RemoveFeeTokensParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemovePriceUpdaters is a Record type
type RemovePriceUpdaters struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemovePriceUpdaters to a map for DAML arguments
func (t RemovePriceUpdaters) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemovePriceUpdaters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemovePriceUpdaters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemovePriceUpdaters to hex string (Canton MCMS format)
func (t RemovePriceUpdaters) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemovePriceUpdaters from hex string (Canton MCMS format)
func (t *RemovePriceUpdaters) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TimestampedPrice is a Record type
type TimestampedPrice struct {
	Price     types.NUMERIC   `json:"price"`
	Timestamp types.TIMESTAMP `json:"timestamp"`
}

// ToMap converts TimestampedPrice to a map for DAML arguments
func (t TimestampedPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["price"] = t.Price

	m["timestamp"] = t.Timestamp

	return m
}

func (t TimestampedPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TimestampedPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TimestampedPrice to hex string (Canton MCMS format)
func (t TimestampedPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TimestampedPrice from hex string (Canton MCMS format)
func (t *TimestampedPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPriceUpdate is a Record type
type TokenPriceUpdate struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	UsdPerToken  types.NUMERIC                            `json:"usdPerToken"`
}

// ToMap converts TokenPriceUpdate to a map for DAML arguments
func (t TokenPriceUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["usdPerToken"] = t.UsdPerToken

	return m
}

func (t TokenPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPriceUpdate to hex string (Canton MCMS format)
func (t TokenPriceUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPriceUpdate from hex string (Canton MCMS format)
func (t *TokenPriceUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferFeeConfig is a Record type
type TokenTransferFeeConfig struct {
	IsEnabled         types.BOOL    `json:"isEnabled"`
	FeeUSD            types.NUMERIC `json:"feeUSD"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts TokenTransferFeeConfig to a map for DAML arguments
func (t TokenTransferFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["feeUSD"] = t.FeeUSD

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t TokenTransferFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferFeeConfig to hex string (Canton MCMS format)
func (t TokenTransferFeeConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferFeeConfig from hex string (Canton MCMS format)
func (t *TokenTransferFeeConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePrices is a Record type
type UpdatePrices struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       types.PARTY  `json:"caller"`
}

// ToMap converts UpdatePrices to a map for DAML arguments
func (t UpdatePrices) ToMap() map[string]any {
	m := make(map[string]any)

	m["priceUpdates"] = model.NestedToDAMLValue(t.PriceUpdates)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t UpdatePrices) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdatePrices) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdatePrices to hex string (Canton MCMS format)
func (t UpdatePrices) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePrices from hex string (Canton MCMS format)
func (t *UpdatePrices) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePricesMCMSParams is UpdatePrices without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type UpdatePricesMCMSParams struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
}

// MarshalHex encodes UpdatePricesMCMSParams to hex string for MCMS operationData.
func (t UpdatePricesMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePricesMCMSParams from hex string.
func (t *UpdatePricesMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePricesParams is a Record type
type UpdatePricesParams struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       types.PARTY  `json:"caller"`
}

// ToMap converts UpdatePricesParams to a map for DAML arguments
func (t UpdatePricesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["priceUpdates"] = model.NestedToDAMLValue(t.PriceUpdates)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t UpdatePricesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdatePricesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdatePricesParams to hex string (Canton MCMS format)
func (t UpdatePricesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePricesParams from hex string (Canton MCMS format)
func (t *UpdatePricesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddPriceUpdaters(args AddPriceUpdaters) (*bind.EncodedChoice, error)
	ApplyDestChainConfigUpdates2(args ApplyDestChainConfigUpdates2) (*bind.EncodedChoice, error)
	ApplyPriceUpdatersUpdate(args ApplyPriceUpdatersUpdate) (*bind.EncodedChoice, error)
	ApplyPriceUpdatersUpdateParams(args ApplyPriceUpdatersUpdateParams) (*bind.EncodedChoice, error)
	Get(args Get) (*bind.EncodedChoice, error)
	GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error)
	GetDestChainConfig2(args GetDestChainConfig2) (*bind.EncodedChoice, error)
	GetDestChainConfig2MCMSParams(args GetDestChainConfig2MCMSParams) (*bind.EncodedChoice, error)
	GetDestinationChainGasPrice(args GetDestinationChainGasPrice) (*bind.EncodedChoice, error)
	GetDestinationChainGasPriceMCMSParams(args GetDestinationChainGasPriceMCMSParams) (*bind.EncodedChoice, error)
	GetFeeTokens(args GetFeeTokens) (*bind.EncodedChoice, error)
	GetFeeTokensMCMSParams(args GetFeeTokensMCMSParams) (*bind.EncodedChoice, error)
	GetTokenPrice(args GetTokenPrice) (*bind.EncodedChoice, error)
	GetTokenPriceMCMSParams(args GetTokenPriceMCMSParams) (*bind.EncodedChoice, error)
	GetTokenTransferFee(args GetTokenTransferFee) (*bind.EncodedChoice, error)
	GetTokenTransferFeeMCMSParams(args GetTokenTransferFeeMCMSParams) (*bind.EncodedChoice, error)
	QuoteGasForExec(args QuoteGasForExec) (*bind.EncodedChoice, error)
	QuoteGasForExecMCMSParams(args QuoteGasForExecMCMSParams) (*bind.EncodedChoice, error)
	RemoveFeeTokens(args RemoveFeeTokens) (*bind.EncodedChoice, error)
	RemoveFeeTokensParams(args RemoveFeeTokensParams) (*bind.EncodedChoice, error)
	RemovePriceUpdaters(args RemovePriceUpdaters) (*bind.EncodedChoice, error)
	UpdatePrices(args UpdatePrices) (*bind.EncodedChoice, error)
	UpdatePricesMCMSParams(args UpdatePricesMCMSParams) (*bind.EncodedChoice, error)
	UpdatePricesParams(args UpdatePricesParams) (*bind.EncodedChoice, error)
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

// AddPriceUpdaters encodes parameters for the AddPriceUpdaters choice.
func (e *encoder) AddPriceUpdaters(args AddPriceUpdaters) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPriceUpdaters", args)
}

// ApplyDestChainConfigUpdates2 encodes parameters for the ApplyDestChainConfigUpdates2 choice.
func (e *encoder) ApplyDestChainConfigUpdates2(args ApplyDestChainConfigUpdates2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyDestChainConfigUpdates2", args)
}

// ApplyPriceUpdatersUpdate encodes parameters for the ApplyPriceUpdatersUpdate choice.
func (e *encoder) ApplyPriceUpdatersUpdate(args ApplyPriceUpdatersUpdate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyPriceUpdatersUpdate", args)
}

// ApplyPriceUpdatersUpdateParams encodes parameters for the ApplyPriceUpdatersUpdate choice.
func (e *encoder) ApplyPriceUpdatersUpdateParams(args ApplyPriceUpdatersUpdateParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyPriceUpdatersUpdate", args)
}

// Get encodes parameters for the Get choice.
func (e *encoder) Get(args Get) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetMCMSParams encodes MCMS parameters (without Caller) for the Get choice.
func (e *encoder) GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetDestChainConfig2 encodes parameters for the GetDestChainConfig2 choice.
func (e *encoder) GetDestChainConfig2(args GetDestChainConfig2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig2", args)
}

// GetDestChainConfig2MCMSParams encodes MCMS parameters (without Caller) for the GetDestChainConfig2 choice.
func (e *encoder) GetDestChainConfig2MCMSParams(args GetDestChainConfig2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig2", args)
}

// GetDestinationChainGasPrice encodes parameters for the GetDestinationChainGasPrice choice.
func (e *encoder) GetDestinationChainGasPrice(args GetDestinationChainGasPrice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestinationChainGasPrice", args)
}

// GetDestinationChainGasPriceMCMSParams encodes MCMS parameters (without Caller) for the GetDestinationChainGasPrice choice.
func (e *encoder) GetDestinationChainGasPriceMCMSParams(args GetDestinationChainGasPriceMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestinationChainGasPrice", args)
}

// GetFeeTokens encodes parameters for the GetFeeTokens choice.
func (e *encoder) GetFeeTokens(args GetFeeTokens) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFeeTokens", args)
}

// GetFeeTokensMCMSParams encodes MCMS parameters (without Caller) for the GetFeeTokens choice.
func (e *encoder) GetFeeTokensMCMSParams(args GetFeeTokensMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFeeTokens", args)
}

// GetTokenPrice encodes parameters for the GetTokenPrice choice.
func (e *encoder) GetTokenPrice(args GetTokenPrice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenPrice", args)
}

// GetTokenPriceMCMSParams encodes MCMS parameters (without Caller) for the GetTokenPrice choice.
func (e *encoder) GetTokenPriceMCMSParams(args GetTokenPriceMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenPrice", args)
}

// GetTokenTransferFee encodes parameters for the GetTokenTransferFee choice.
func (e *encoder) GetTokenTransferFee(args GetTokenTransferFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenTransferFee", args)
}

// GetTokenTransferFeeMCMSParams encodes MCMS parameters (without Caller) for the GetTokenTransferFee choice.
func (e *encoder) GetTokenTransferFeeMCMSParams(args GetTokenTransferFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenTransferFee", args)
}

// QuoteGasForExec encodes parameters for the QuoteGasForExec choice.
func (e *encoder) QuoteGasForExec(args QuoteGasForExec) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("QuoteGasForExec", args)
}

// QuoteGasForExecMCMSParams encodes MCMS parameters (without Caller) for the QuoteGasForExec choice.
func (e *encoder) QuoteGasForExecMCMSParams(args QuoteGasForExecMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("QuoteGasForExec", args)
}

// RemoveFeeTokens encodes parameters for the RemoveFeeTokens choice.
func (e *encoder) RemoveFeeTokens(args RemoveFeeTokens) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveFeeTokens", args)
}

// RemoveFeeTokensParams encodes parameters for the RemoveFeeTokens choice.
func (e *encoder) RemoveFeeTokensParams(args RemoveFeeTokensParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveFeeTokens", args)
}

// RemovePriceUpdaters encodes parameters for the RemovePriceUpdaters choice.
func (e *encoder) RemovePriceUpdaters(args RemovePriceUpdaters) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemovePriceUpdaters", args)
}

// UpdatePrices encodes parameters for the UpdatePrices choice.
func (e *encoder) UpdatePrices(args UpdatePrices) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// UpdatePricesMCMSParams encodes MCMS parameters (without Caller) for the UpdatePrices choice.
func (e *encoder) UpdatePricesMCMSParams(args UpdatePricesMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// UpdatePricesParams encodes parameters for the UpdatePrices choice.
func (e *encoder) UpdatePricesParams(args UpdatePricesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
