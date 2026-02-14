package ccipreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

// ApplyDestChainConfigUpdates is a Record type
type ApplyDestChainConfigUpdates struct {
	DestChainConfigArgs []DestChainConfigArgs `json:"destChainConfigArgs"`
}

// ToMap converts ApplyDestChainConfigUpdates to a map for DAML arguments
func (t ApplyDestChainConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t ApplyDestChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

// ToMap converts ApplyFeeTokenUpdates to a map for DAML arguments
func (t ApplyFeeTokenUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokensToRemove"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToRemove))
		for _, e := range t.FeeTokensToRemove {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["feeTokensToAdd"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToAdd))
		for _, e := range t.FeeTokensToAdd {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ApplyFeeTokenUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ApplyFeeTokenUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfig2 is a Record type
type DestChainConfig2 struct {
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

// ToMap converts DestChainConfig2 to a map for DAML arguments
func (t DestChainConfig2) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["maxDataBytes"] = int64(t.MaxDataBytes)

	m["maxPerMsgGasLimit"] = int64(t.MaxPerMsgGasLimit)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destGasPerPayloadByteBase"] = int64(t.DestGasPerPayloadByteBase)

	m["chainFamilySelector"] = string(t.ChainFamilySelector)

	m["defaultTxGasLimit"] = int64(t.DefaultTxGasLimit)

	m["networkFeeUSD"] = t.NetworkFeeUSD

	m["defaultTokenFeeUSD"] = t.DefaultTokenFeeUSD

	m["defaultTokenDestGasOverhead"] = int64(t.DefaultTokenDestGasOverhead)

	return m
}

func (t DestChainConfig2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *DestChainConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfigArgs is a Record type
type DestChainConfigArgs struct {
	DestChainSelector NUMERIC          `json:"destChainSelector"`
	DestChainConfig   DestChainConfig2 `json:"destChainConfig"`
}

// ToMap converts DestChainConfigArgs to a map for DAML arguments
func (t DestChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["destChainConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DestChainConfig).(mapper); ok {
			return m.toMap()
		}
		return t.DestChainConfig
	}()

	return m
}

func (t DestChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *DestChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeQuoter is a Template type
type FeeQuoter struct {
	InstanceId                       TEXT    `json:"instanceId"`
	Owner                            PARTY   `json:"owner"`
	FeeTokens                        GENMAP  `json:"feeTokens"`
	DestChainConfigs                 GENMAP  `json:"destChainConfigs"`
	TokenTransferFeeConfigs          GENMAP  `json:"tokenTransferFeeConfigs"`
	UsdPerUnitGasByDestChainSelector GENMAP  `json:"usdPerUnitGasByDestChainSelector"`
	UsdPerToken                      GENMAP  `json:"usdPerToken"`
	PriceUpdaters                    []PARTY `json:"priceUpdaters"`
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
	args["feeTokens"] = func() any {
		if t.FeeTokens == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.FeeTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

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
	args["feeTokens"] = func() any {
		if t.FeeTokens == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.FeeTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

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
	return jsonCodec.Marshall(t)
}

func (t *FeeQuoter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for FeeQuoter

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

// FeeQuoterGetTokenTransferFee exercises the FeeQuoter_GetTokenTransferFee choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) FeeQuoterGetTokenTransferFee(contractID string, args FeeQuoterGetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterGetTokenTransferFeeWithPackageID exercises the FeeQuoter_GetTokenTransferFee choice using the provided package ID instead of package name
func (t FeeQuoter) FeeQuoterGetTokenTransferFeeWithPackageID(contractID string, packageID string, args FeeQuoterGetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterFinalizeFee exercises the FeeQuoter_FinalizeFee choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) FeeQuoterFinalizeFee(contractID string, args FeeQuoterFinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterFinalizeFeeWithPackageID exercises the FeeQuoter_FinalizeFee choice using the provided package ID instead of package name
func (t FeeQuoter) FeeQuoterFinalizeFeeWithPackageID(contractID string, packageID string, args FeeQuoterFinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "FeeQuoter_FinalizeFee",
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

// GetPremiumMultiplierWeiPerEth exercises the GetPremiumMultiplierWeiPerEth choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetPremiumMultiplierWeiPerEth(contractID string, args GetPremiumMultiplierWeiPerEth) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetPremiumMultiplierWeiPerEth",
		Arguments:  argsToMap(args),
	}
}

// GetPremiumMultiplierWeiPerEthWithPackageID exercises the GetPremiumMultiplierWeiPerEth choice using the provided package ID instead of package name
func (t FeeQuoter) GetPremiumMultiplierWeiPerEthWithPackageID(contractID string, packageID string, args GetPremiumMultiplierWeiPerEth) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetPremiumMultiplierWeiPerEth",
		Arguments:  argsToMap(args),
	}
}

// ApplyFeeTokenUpdates exercises the ApplyFeeTokenUpdates choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyFeeTokenUpdates(contractID string, args ApplyFeeTokenUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyFeeTokenUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyFeeTokenUpdatesWithPackageID exercises the ApplyFeeTokenUpdates choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyFeeTokenUpdatesWithPackageID(contractID string, packageID string, args ApplyFeeTokenUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyFeeTokenUpdates",
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

// Archive exercises the Archive choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FeeQuoter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ApplyDestChainConfigUpdates exercises the ApplyDestChainConfigUpdates choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyDestChainConfigUpdates(contractID string, args ApplyDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainConfigUpdatesWithPackageID exercises the ApplyDestChainConfigUpdates choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyDestChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// FeeQuoterFinalizeFee is a Record type
type FeeQuoterFinalizeFee struct {
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	Caller            PARTY       `json:"caller"`
}

// ToMap converts FeeQuoterFinalizeFee to a map for DAML arguments
func (t FeeQuoterFinalizeFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterFinalizeFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *FeeQuoterFinalizeFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeQuoterGetTokenTransferFee is a Record type
type FeeQuoterGetTokenTransferFee struct {
	DestChainSelector NUMERIC      `json:"destChainSelector"`
	Token             InstrumentId `json:"token"`
	Caller            PARTY        `json:"caller"`
}

// ToMap converts FeeQuoterGetTokenTransferFee to a map for DAML arguments
func (t FeeQuoterGetTokenTransferFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["token"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Token).(mapper); ok {
			return m.toMap()
		}
		return t.Token
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterGetTokenTransferFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *FeeQuoterGetTokenTransferFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeTokenArgs is a Record type
type FeeTokenArgs struct {
	InstrumentId      InstrumentId `json:"instrumentId"`
	PremiumMultiplier NUMERIC      `json:"premiumMultiplier"`
}

// ToMap converts FeeTokenArgs to a map for DAML arguments
func (t FeeTokenArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["premiumMultiplier"] = t.PremiumMultiplier

	return m
}

func (t FeeTokenArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *FeeTokenArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GasPriceUpdate is a Record type
type GasPriceUpdate struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	UsdPerUnitGas     NUMERIC `json:"usdPerUnitGas"`
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
	return jsonCodec.Marshall(t)
}

func (t *GasPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestChainConfig2 is a Record type
type GetDestChainConfig2 struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *GetDestChainConfig2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestinationChainGasPrice is a Record type
type GetDestinationChainGasPrice struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *GetDestinationChainGasPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetFeeTokens is a Record type
type GetFeeTokens struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts GetFeeTokens to a map for DAML arguments
func (t GetFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetPremiumMultiplierWeiPerEth is a Record type
type GetPremiumMultiplierWeiPerEth struct {
	Token  InstrumentId `json:"token"`
	Caller PARTY        `json:"caller"`
}

// ToMap converts GetPremiumMultiplierWeiPerEth to a map for DAML arguments
func (t GetPremiumMultiplierWeiPerEth) ToMap() map[string]any {
	m := make(map[string]any)

	m["token"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Token).(mapper); ok {
			return m.toMap()
		}
		return t.Token
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetPremiumMultiplierWeiPerEth) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetPremiumMultiplierWeiPerEth) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetTokenPrice is a Record type
type GetTokenPrice struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
}

// ToMap converts GetTokenPrice to a map for DAML arguments
func (t GetTokenPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetTokenPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["gasPriceUpdates"] = func() []any {
		res := make([]any, 0, len(t.GasPriceUpdates))
		for _, e := range t.GasPriceUpdates {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t PriceUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PriceUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TimestampedPrice is a Record type
type TimestampedPrice struct {
	Price     NUMERIC   `json:"price"`
	Timestamp TIMESTAMP `json:"timestamp"`
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
	return jsonCodec.Marshall(t)
}

func (t *TimestampedPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPriceUpdate is a Record type
type TokenPriceUpdate struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	UsdPerToken  NUMERIC      `json:"usdPerToken"`
}

// ToMap converts TokenPriceUpdate to a map for DAML arguments
func (t TokenPriceUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["usdPerToken"] = t.UsdPerToken

	return m
}

func (t TokenPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

// ToMap converts TokenTransferFeeConfig to a map for DAML arguments
func (t TokenTransferFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeUSD"] = t.FeeUSD

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t TokenTransferFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenTransferFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdatePrices is a Record type
type UpdatePrices struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       PARTY        `json:"caller"`
}

// ToMap converts UpdatePrices to a map for DAML arguments
func (t UpdatePrices) ToMap() map[string]any {
	m := make(map[string]any)

	m["priceUpdates"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PriceUpdates).(mapper); ok {
			return m.toMap()
		}
		return t.PriceUpdates
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t UpdatePrices) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UpdatePrices) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
