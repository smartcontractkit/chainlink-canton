package lockreleasetokenpool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

// ApplyDestChainConfigUpdates is a Record type
type ApplyDestChainConfigUpdates struct {
	DestChainConfigArgs []DestChainConfigArgs `json:"destChainConfigArgs"`
}

// toMap converts ApplyDestChainConfigUpdates to a map for DAML arguments
func (t ApplyDestChainConfigUpdates) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainConfigArgs": func() []interface{} {
			res := make([]interface{}, 0, len(t.DestChainConfigArgs))
			for _, e := range t.DestChainConfigArgs {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ApplyDestChainConfigUpdates using JsonCodec
func (t ApplyDestChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ApplyDestChainConfigUpdates using JsonCodec
func (t *ApplyDestChainConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ApplyFeeTokenUpdates is a Record type
type ApplyFeeTokenUpdates struct {
	FeeTokensToRemove []InstrumentId `json:"feeTokensToRemove"`
	FeeTokensToAdd    []FeeTokenArgs `json:"feeTokensToAdd"`
	Caller            PARTY          `json:"caller"`
}

// toMap converts ApplyFeeTokenUpdates to a map for DAML arguments
func (t ApplyFeeTokenUpdates) toMap() map[string]interface{} {
	return map[string]interface{}{

		"feeTokensToRemove": func() []interface{} {
			res := make([]interface{}, 0, len(t.FeeTokensToRemove))
			for _, e := range t.FeeTokensToRemove {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
		"feeTokensToAdd": func() []interface{} {
			res := make([]interface{}, 0, len(t.FeeTokensToAdd))
			for _, e := range t.FeeTokensToAdd {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for ApplyFeeTokenUpdates using JsonCodec
func (t ApplyFeeTokenUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ApplyFeeTokenUpdates using JsonCodec
func (t *ApplyFeeTokenUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfig22 is a Record type
type DestChainConfig22 struct {
	IsEnabled                   BOOL    `json:"isEnabled"`
	MaxDataBytes                INT64   `json:"maxDataBytes"`
	MaxPerMsgGasLimit           INT64   `json:"maxPerMsgGasLimit"`
	DestGasOverhead             INT64   `json:"destGasOverhead"`
	DestGasPerPayloadByteBase   INT64   `json:"destGasPerPayloadByteBase"`
	ChainFamilySelector         TEXT    `json:"chainFamilySelector"`
	DefaultTxGasLimit           INT64   `json:"defaultTxGasLimit"`
	NetworkFeeUSD               NUMERIC `json:"networkFeeUSD"`
	DefaultTokenFeeUSD          NUMERIC `json:"defaultTokenFeeUSD"`
	DefaultTokenDestGasOverhead INT64   `json:"defaultTokenDestGasOverhead"`
}

// toMap converts DestChainConfig22 to a map for DAML arguments
func (t DestChainConfig22) toMap() map[string]interface{} {
	return map[string]interface{}{

		"isEnabled":                   bool(t.IsEnabled),
		"maxDataBytes":                int64(t.MaxDataBytes),
		"maxPerMsgGasLimit":           int64(t.MaxPerMsgGasLimit),
		"destGasOverhead":             int64(t.DestGasOverhead),
		"destGasPerPayloadByteBase":   int64(t.DestGasPerPayloadByteBase),
		"chainFamilySelector":         string(t.ChainFamilySelector),
		"defaultTxGasLimit":           int64(t.DefaultTxGasLimit),
		"networkFeeUSD":               (*big.Int)(t.NetworkFeeUSD),
		"defaultTokenFeeUSD":          (*big.Int)(t.DefaultTokenFeeUSD),
		"defaultTokenDestGasOverhead": int64(t.DefaultTokenDestGasOverhead),
	}
}

// MarshalJSON implements custom JSON marshaling for DestChainConfig22 using JsonCodec
func (t DestChainConfig22) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DestChainConfig22 using JsonCodec
func (t *DestChainConfig22) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfigArgs is a Record type
type DestChainConfigArgs struct {
	DestChainSelector NUMERIC           `json:"destChainSelector"`
	DestChainConfig   DestChainConfig22 `json:"destChainConfig"`
}

// toMap converts DestChainConfigArgs to a map for DAML arguments
func (t DestChainConfigArgs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"destChainConfig": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.DestChainConfig).(mapper); ok {
				return m.toMap()
			}
			return t.DestChainConfig
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for DestChainConfigArgs using JsonCodec
func (t DestChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DestChainConfigArgs using JsonCodec
func (t *DestChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeQuoter is a Template type
type FeeQuoter struct {
	Owner                            PARTY   `json:"owner"`
	InstanceId                       TEXT    `json:"instanceId"`
	FeeTokens                        GENMAP  `json:"feeTokens"`
	DestChainConfigs                 GENMAP  `json:"destChainConfigs"`
	TokenTransferFeeConfigs          GENMAP  `json:"tokenTransferFeeConfigs"`
	UsdPerUnitGasByDestChainSelector GENMAP  `json:"usdPerUnitGasByDestChainSelector"`
	UsdPerToken                      GENMAP  `json:"usdPerToken"`
	PriceUpdaters                    []PARTY `json:"priceUpdaters"`
}

// GetTemplateID returns the template ID for this template
func (t FeeQuoter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter")
}

// CreateCommand returns a CreateCommand for this template
func (t FeeQuoter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	if t.FeeTokens != nil && len(t.FeeTokens) > 0 {
		args["feeTokens"] = map[string]interface{}{"_type": "genmap", "value": t.FeeTokens}
	}

	if t.DestChainConfigs != nil && len(t.DestChainConfigs) > 0 {
		args["destChainConfigs"] = map[string]interface{}{"_type": "genmap", "value": t.DestChainConfigs}
	}

	if t.TokenTransferFeeConfigs != nil && len(t.TokenTransferFeeConfigs) > 0 {
		args["tokenTransferFeeConfigs"] = map[string]interface{}{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}

	if t.UsdPerUnitGasByDestChainSelector != nil && len(t.UsdPerUnitGasByDestChainSelector) > 0 {
		args["usdPerUnitGasByDestChainSelector"] = map[string]interface{}{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}

	if t.UsdPerToken != nil && len(t.UsdPerToken) > 0 {
		args["usdPerToken"] = map[string]interface{}{"_type": "genmap", "value": t.UsdPerToken}
	}

	if len(t.PriceUpdaters) > 0 {
		args["priceUpdaters"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.PriceUpdaters))
			for _, e := range t.PriceUpdaters {
				res = append(res, e.ToMap())
			}
			return res
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for FeeQuoter using JsonCodec
func (t FeeQuoter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeeQuoter using JsonCodec
func (t *FeeQuoter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for FeeQuoter

// GetTokenPrice exercises the GetTokenPrice choice on this FeeQuoter contract
func (t FeeQuoter) GetTokenPrice(contractID string, args GetTokenPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenPrice",
		Arguments:  argsToMap(args),
	}
}

// GetDestinationChainGasPrice exercises the GetDestinationChainGasPrice choice on this FeeQuoter contract
func (t FeeQuoter) GetDestinationChainGasPrice(contractID string, args GetDestinationChainGasPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestinationChainGasPrice",
		Arguments:  argsToMap(args),
	}
}

// UpdatePrices exercises the UpdatePrices choice on this FeeQuoter contract
func (t FeeQuoter) UpdatePrices(contractID string, args UpdatePrices) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "UpdatePrices",
		Arguments:  argsToMap(args),
	}
}

// GetValidatedFee exercises the GetValidatedFee choice on this FeeQuoter contract
func (t FeeQuoter) GetValidatedFee(contractID string, args GetValidatedFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetValidatedFee",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterGetTokenTransferFee exercises the FeeQuoter_GetTokenTransferFee choice on this FeeQuoter contract
func (t FeeQuoter) FeeQuoterGetTokenTransferFee(contractID string, args FeeQuoterGetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterQuoteGasForExec exercises the FeeQuoter_QuoteGasForExec choice on this FeeQuoter contract
func (t FeeQuoter) FeeQuoterQuoteGasForExec(contractID string, args FeeQuoterQuoteGasForExec) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_QuoteGasForExec",
		Arguments:  argsToMap(args),
	}
}

// GetFeeTokens exercises the GetFeeTokens choice on this FeeQuoter contract
func (t FeeQuoter) GetFeeTokens(contractID string, args GetFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// GetPremiumMultiplierWeiPerEth exercises the GetPremiumMultiplierWeiPerEth choice on this FeeQuoter contract
func (t FeeQuoter) GetPremiumMultiplierWeiPerEth(contractID string, args GetPremiumMultiplierWeiPerEth) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetPremiumMultiplierWeiPerEth",
		Arguments:  argsToMap(args),
	}
}

// ApplyFeeTokenUpdates exercises the ApplyFeeTokenUpdates choice on this FeeQuoter contract
func (t FeeQuoter) ApplyFeeTokenUpdates(contractID string, args ApplyFeeTokenUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyFeeTokenUpdates",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfig exercises the GetDestChainConfig choice on this FeeQuoter contract
func (t FeeQuoter) GetDestChainConfig(contractID string, args GetDestChainConfig2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FeeQuoter contract
func (t FeeQuoter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ApplyDestChainConfigUpdates exercises the ApplyDestChainConfigUpdates choice on this FeeQuoter contract
func (t FeeQuoter) ApplyDestChainConfigUpdates(contractID string, args ApplyDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterGetTokenTransferFee is a Record type
type FeeQuoterGetTokenTransferFee struct {
	DestChainSelector NUMERIC      `json:"destChainSelector"`
	Token             InstrumentId `json:"token"`
	Caller            PARTY        `json:"caller"`
}

// toMap converts FeeQuoterGetTokenTransferFee to a map for DAML arguments
func (t FeeQuoterGetTokenTransferFee) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"token": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Token).(mapper); ok {
				return m.toMap()
			}
			return t.Token
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for FeeQuoterGetTokenTransferFee using JsonCodec
func (t FeeQuoterGetTokenTransferFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeeQuoterGetTokenTransferFee using JsonCodec
func (t *FeeQuoterGetTokenTransferFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeQuoterQuoteGasForExec is a Record type
type FeeQuoterQuoteGasForExec struct {
	DestChainSelector NUMERIC      `json:"destChainSelector"`
	NonCalldataGas    INT64        `json:"nonCalldataGas"`
	CalldataSize      INT64        `json:"calldataSize"`
	FeeToken          InstrumentId `json:"feeToken"`
	Caller            PARTY        `json:"caller"`
}

// toMap converts FeeQuoterQuoteGasForExec to a map for DAML arguments
func (t FeeQuoterQuoteGasForExec) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"nonCalldataGas":    int64(t.NonCalldataGas),
		"calldataSize":      int64(t.CalldataSize),
		"feeToken": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeToken).(mapper); ok {
				return m.toMap()
			}
			return t.FeeToken
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for FeeQuoterQuoteGasForExec using JsonCodec
func (t FeeQuoterQuoteGasForExec) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeeQuoterQuoteGasForExec using JsonCodec
func (t *FeeQuoterQuoteGasForExec) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeTokenArgs is a Record type
type FeeTokenArgs struct {
	InstrumentId      InstrumentId `json:"instrumentId"`
	PremiumMultiplier NUMERIC      `json:"premiumMultiplier"`
}

// toMap converts FeeTokenArgs to a map for DAML arguments
func (t FeeTokenArgs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"premiumMultiplier": (*big.Int)(t.PremiumMultiplier),
	}
}

// MarshalJSON implements custom JSON marshaling for FeeTokenArgs using JsonCodec
func (t FeeTokenArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeeTokenArgs using JsonCodec
func (t *FeeTokenArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GasPriceUpdate is a Record type
type GasPriceUpdate struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	UsdPerUnitGas     NUMERIC `json:"usdPerUnitGas"`
}

// toMap converts GasPriceUpdate to a map for DAML arguments
func (t GasPriceUpdate) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"usdPerUnitGas":     (*big.Int)(t.UsdPerUnitGas),
	}
}

// MarshalJSON implements custom JSON marshaling for GasPriceUpdate using JsonCodec
func (t GasPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GasPriceUpdate using JsonCodec
func (t *GasPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestChainConfig2 is a Record type
type GetDestChainConfig2 struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
}

// toMap converts GetDestChainConfig2 to a map for DAML arguments
func (t GetDestChainConfig2) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"caller":            t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetDestChainConfig2 using JsonCodec
func (t GetDestChainConfig2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetDestChainConfig2 using JsonCodec
func (t *GetDestChainConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestinationChainGasPrice is a Record type
type GetDestinationChainGasPrice struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
}

// toMap converts GetDestinationChainGasPrice to a map for DAML arguments
func (t GetDestinationChainGasPrice) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"caller":            t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetDestinationChainGasPrice using JsonCodec
func (t GetDestinationChainGasPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetDestinationChainGasPrice using JsonCodec
func (t *GetDestinationChainGasPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetFeeTokens is a Record type
type GetFeeTokens struct {
	Caller PARTY `json:"caller"`
}

// toMap converts GetFeeTokens to a map for DAML arguments
func (t GetFeeTokens) toMap() map[string]interface{} {
	return map[string]interface{}{

		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetFeeTokens using JsonCodec
func (t GetFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetFeeTokens using JsonCodec
func (t *GetFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetPremiumMultiplierWeiPerEth is a Record type
type GetPremiumMultiplierWeiPerEth struct {
	Token  InstrumentId `json:"token"`
	Caller PARTY        `json:"caller"`
}

// toMap converts GetPremiumMultiplierWeiPerEth to a map for DAML arguments
func (t GetPremiumMultiplierWeiPerEth) toMap() map[string]interface{} {
	return map[string]interface{}{

		"token": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Token).(mapper); ok {
				return m.toMap()
			}
			return t.Token
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetPremiumMultiplierWeiPerEth using JsonCodec
func (t GetPremiumMultiplierWeiPerEth) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetPremiumMultiplierWeiPerEth using JsonCodec
func (t *GetPremiumMultiplierWeiPerEth) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetTokenPrice is a Record type
type GetTokenPrice struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts GetTokenPrice to a map for DAML arguments
func (t GetTokenPrice) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetTokenPrice using JsonCodec
func (t GetTokenPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetTokenPrice using JsonCodec
func (t *GetTokenPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetValidatedFee is a Record type
type GetValidatedFee struct {
	DestChainSelector NUMERIC           `json:"destChainSelector"`
	Message           Canton2AnyMessage `json:"message"`
	Caller            PARTY             `json:"caller"`
}

// toMap converts GetValidatedFee to a map for DAML arguments
func (t GetValidatedFee) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetValidatedFee using JsonCodec
func (t GetValidatedFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetValidatedFee using JsonCodec
func (t *GetValidatedFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PriceUpdates is a Record type
type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate `json:"tokenPriceUpdates"`
	GasPriceUpdates   []GasPriceUpdate   `json:"gasPriceUpdates"`
}

// toMap converts PriceUpdates to a map for DAML arguments
func (t PriceUpdates) toMap() map[string]interface{} {
	return map[string]interface{}{

		"tokenPriceUpdates": func() []interface{} {
			res := make([]interface{}, 0, len(t.TokenPriceUpdates))
			for _, e := range t.TokenPriceUpdates {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
		"gasPriceUpdates": func() []interface{} {
			res := make([]interface{}, 0, len(t.GasPriceUpdates))
			for _, e := range t.GasPriceUpdates {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for PriceUpdates using JsonCodec
func (t PriceUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PriceUpdates using JsonCodec
func (t *PriceUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TimestampedPrice is a Record type
type TimestampedPrice struct {
	Price     NUMERIC   `json:"price"`
	Timestamp TIMESTAMP `json:"timestamp"`
}

// toMap converts TimestampedPrice to a map for DAML arguments
func (t TimestampedPrice) toMap() map[string]interface{} {
	return map[string]interface{}{

		"price":     (*big.Int)(t.Price),
		"timestamp": t.Timestamp,
	}
}

// MarshalJSON implements custom JSON marshaling for TimestampedPrice using JsonCodec
func (t TimestampedPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TimestampedPrice using JsonCodec
func (t *TimestampedPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPriceUpdate is a Record type
type TokenPriceUpdate struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	UsdPerToken  NUMERIC      `json:"usdPerToken"`
}

// toMap converts TokenPriceUpdate to a map for DAML arguments
func (t TokenPriceUpdate) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"usdPerToken": (*big.Int)(t.UsdPerToken),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenPriceUpdate using JsonCodec
func (t TokenPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPriceUpdate using JsonCodec
func (t *TokenPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenTransferFeeConfig is a Record type
type TokenTransferFeeConfig struct {
	FeeUSD            NUMERIC `json:"feeUSD"`
	DestGasOverhead   INT64   `json:"destGasOverhead"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
}

// toMap converts TokenTransferFeeConfig to a map for DAML arguments
func (t TokenTransferFeeConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"feeUSD":            (*big.Int)(t.FeeUSD),
		"destGasOverhead":   int64(t.DestGasOverhead),
		"destBytesOverhead": int64(t.DestBytesOverhead),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenTransferFeeConfig using JsonCodec
func (t TokenTransferFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenTransferFeeConfig using JsonCodec
func (t *TokenTransferFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdatePrices is a Record type
type UpdatePrices struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts UpdatePrices to a map for DAML arguments
func (t UpdatePrices) toMap() map[string]interface{} {
	return map[string]interface{}{

		"priceUpdates": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.PriceUpdates).(mapper); ok {
				return m.toMap()
			}
			return t.PriceUpdates
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for UpdatePrices using JsonCodec
func (t UpdatePrices) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UpdatePrices using JsonCodec
func (t *UpdatePrices) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
