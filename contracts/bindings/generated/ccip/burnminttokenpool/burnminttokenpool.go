package burnminttokenpool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipcodec "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	extensionapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	api "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "ccip-burn-mint-token-pool-v2"
	PackageID   = "fad3d7f03cafedfa5060726d0981d9edd1aa91ebd2e120a2fe3b8b7e9fbf49c3"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	BurnMintFactoryExtraArgsMetaValuesContextKey    = types.TEXT("burn-mint-factory-extra-args-meta-values")
	BurnMintFactoryExtraArgsContextValuesContextKey = types.TEXT("burn-mint-factory-extra-args-values")
	BurnMintFactoryContextKey                       = types.TEXT("burn-mint-factory")
	BpsDenominator                                  = types.NUMERIC("10000.")
)

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

// AddPoolReceiveContextContractValue is a Record type
type AddPoolReceiveContextContractValue struct {
	ContextKey       types.TEXT        `json:"contextKey"`
	ReferredContract types.CONTRACT_ID `json:"referredContract"`
}

// ToMap converts AddPoolReceiveContextContractValue to a map for DAML arguments
func (t AddPoolReceiveContextContractValue) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	m["referredContract"] = model.NestedToDAMLValue(t.ReferredContract)

	return m
}

func (t AddPoolReceiveContextContractValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPoolReceiveContextContractValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPoolReceiveContextContractValue to hex string (Canton MCMS format)
func (t AddPoolReceiveContextContractValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPoolReceiveContextContractValue from hex string (Canton MCMS format)
func (t *AddPoolReceiveContextContractValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddPoolReceiveContextContractValueParams is a Record type
type AddPoolReceiveContextContractValueParams struct {
	ContextKey              types.TEXT                      `json:"contextKey"`
	ReferentInstanceAddress chainlinkapi.RawInstanceAddress `json:"referentInstanceAddress"`
}

// ToMap converts AddPoolReceiveContextContractValueParams to a map for DAML arguments
func (t AddPoolReceiveContextContractValueParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	m["referentInstanceAddress"] = model.NestedToDAMLValue(t.ReferentInstanceAddress)

	return m
}

func (t AddPoolReceiveContextContractValueParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPoolReceiveContextContractValueParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPoolReceiveContextContractValueParams to hex string (Canton MCMS format)
func (t AddPoolReceiveContextContractValueParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPoolReceiveContextContractValueParams from hex string (Canton MCMS format)
func (t *AddPoolReceiveContextContractValueParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddPoolReceiveContextNonContractValue is a Record type
type AddPoolReceiveContextNonContractValue struct {
	ContextKey types.TEXT                            `json:"contextKey"`
	Value      splice_api_token_metadata_v1.AnyValue `json:"value"`
}

// ToMap converts AddPoolReceiveContextNonContractValue to a map for DAML arguments
func (t AddPoolReceiveContextNonContractValue) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	m["value"] = model.NestedToDAMLValue(t.Value)

	return m
}

func (t AddPoolReceiveContextNonContractValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPoolReceiveContextNonContractValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPoolReceiveContextNonContractValue to hex string (Canton MCMS format)
func (t AddPoolReceiveContextNonContractValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPoolReceiveContextNonContractValue from hex string (Canton MCMS format)
func (t *AddPoolReceiveContextNonContractValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddPoolReceiveContextNonContractValueParams is a Record type
type AddPoolReceiveContextNonContractValueParams struct {
	ContextKey   types.TEXT `json:"contextKey"`
	ValuePayload types.TEXT `json:"valuePayload"`
}

// ToMap converts AddPoolReceiveContextNonContractValueParams to a map for DAML arguments
func (t AddPoolReceiveContextNonContractValueParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	m["valuePayload"] = string(t.ValuePayload)

	return m
}

func (t AddPoolReceiveContextNonContractValueParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPoolReceiveContextNonContractValueParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPoolReceiveContextNonContractValueParams to hex string (Canton MCMS format)
func (t AddPoolReceiveContextNonContractValueParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPoolReceiveContextNonContractValueParams from hex string (Canton MCMS format)
func (t *AddPoolReceiveContextNonContractValueParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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

// BurnMintTokenPool is a Template type
type BurnMintTokenPool struct {
	InstanceId              types.TEXT                                 `json:"instanceId"`
	PoolOwner               types.PARTY                                `json:"poolOwner"`
	CcipOwner               types.PARTY                                `json:"ccipOwner"`
	InstrumentId            splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Decimals                types.INT64                                `json:"decimals"`
	RateLimitAdmin          *types.PARTY                               `json:"rateLimitAdmin" hex:"optional"`
	RemoteChainConfigs      map[types.NUMERIC]RemoteChainConfig        `json:"remoteChainConfigs"`
	TokenTransferFeeConfigs map[types.NUMERIC]TokenTransferFeeConfig   `json:"tokenTransferFeeConfigs"`
	PoolReceiveContext      splice_api_token_metadata_v1.ChoiceContext `json:"poolReceiveContext"`
	TransferTimeout         TransferTimeout                            `json:"transferTimeout"`
	Deps                    BurnMintTokenPoolDeps                      `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t BurnMintTokenPool) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t BurnMintTokenPool) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t BurnMintTokenPool) CreateCommand() *model.CreateCommand {
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
func (t BurnMintTokenPool) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
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

func (t BurnMintTokenPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintTokenPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnMintTokenPool to hex string (Canton MCMS format)
func (t BurnMintTokenPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnMintTokenPool from hex string (Canton MCMS format)
func (t *BurnMintTokenPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for BurnMintTokenPool

// LockOrBurn exercises the LockOrBurn choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) LockOrBurn(contractID string, args LockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// LockOrBurnWithPackageID exercises the LockOrBurn choice using the provided package ID instead of package name
func (t BurnMintTokenPool) LockOrBurnWithPackageID(contractID string, packageID string, args LockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// ReleaseFromTicket exercises the ReleaseFromTicket choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) ReleaseFromTicket(contractID string, args ReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// ReleaseFromTicketWithPackageID exercises the ReleaseFromTicket choice using the provided package ID instead of package name
func (t BurnMintTokenPool) ReleaseFromTicketWithPackageID(contractID string, packageID string, args ReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimitConfig exercises the SetRateLimitConfig choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) SetRateLimitConfig(contractID string, args SetRateLimitConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimitConfig",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimitConfigWithPackageID exercises the SetRateLimitConfig choice using the provided package ID instead of package name
func (t BurnMintTokenPool) SetRateLimitConfigWithPackageID(contractID string, packageID string, args SetRateLimitConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimitConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyTokenTransferFeeConfigUpdates exercises the ApplyTokenTransferFeeConfigUpdates choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) ApplyTokenTransferFeeConfigUpdates(contractID string, args ApplyTokenTransferFeeConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyTokenTransferFeeConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyTokenTransferFeeConfigUpdatesWithPackageID exercises the ApplyTokenTransferFeeConfigUpdates choice using the provided package ID instead of package name
func (t BurnMintTokenPool) ApplyTokenTransferFeeConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyTokenTransferFeeConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyTokenTransferFeeConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimiterReferences exercises the SetRateLimiterReferences choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) SetRateLimiterReferences(contractID string, args SetRateLimiterReferences) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimiterReferences",
		Arguments:  argsToMap(args),
	}
}

// SetRateLimiterReferencesWithPackageID exercises the SetRateLimiterReferences choice using the provided package ID instead of package name
func (t BurnMintTokenPool) SetRateLimiterReferencesWithPackageID(contractID string, packageID string, args SetRateLimiterReferences) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetRateLimiterReferences",
		Arguments:  argsToMap(args),
	}
}

// ApplyChainUpdates exercises the ApplyChainUpdates choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) ApplyChainUpdates(contractID string, args ApplyChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyChainUpdatesWithPackageID exercises the ApplyChainUpdates choice using the provided package ID instead of package name
func (t BurnMintTokenPool) ApplyChainUpdatesWithPackageID(contractID string, packageID string, args ApplyChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ApplyChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfig exercises the SetDynamicConfig choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) SetDynamicConfig(contractID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfigWithPackageID exercises the SetDynamicConfig choice using the provided package ID instead of package name
func (t BurnMintTokenPool) SetDynamicConfigWithPackageID(contractID string, packageID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// AddPoolReceiveContextNonContractValue exercises the AddPoolReceiveContextNonContractValue choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) AddPoolReceiveContextNonContractValue(contractID string, args AddPoolReceiveContextNonContractValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "AddPoolReceiveContextNonContractValue",
		Arguments:  argsToMap(args),
	}
}

// AddPoolReceiveContextNonContractValueWithPackageID exercises the AddPoolReceiveContextNonContractValue choice using the provided package ID instead of package name
func (t BurnMintTokenPool) AddPoolReceiveContextNonContractValueWithPackageID(contractID string, packageID string, args AddPoolReceiveContextNonContractValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "AddPoolReceiveContextNonContractValue",
		Arguments:  argsToMap(args),
	}
}

// AddPoolReceiveContextContractValue exercises the AddPoolReceiveContextContractValue choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) AddPoolReceiveContextContractValue(contractID string, args AddPoolReceiveContextContractValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "AddPoolReceiveContextContractValue",
		Arguments:  argsToMap(args),
	}
}

// AddPoolReceiveContextContractValueWithPackageID exercises the AddPoolReceiveContextContractValue choice using the provided package ID instead of package name
func (t BurnMintTokenPool) AddPoolReceiveContextContractValueWithPackageID(contractID string, packageID string, args AddPoolReceiveContextContractValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "AddPoolReceiveContextContractValue",
		Arguments:  argsToMap(args),
	}
}

// RemovePoolReceiveContextValue exercises the RemovePoolReceiveContextValue choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) RemovePoolReceiveContextValue(contractID string, args RemovePoolReceiveContextValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "RemovePoolReceiveContextValue",
		Arguments:  argsToMap(args),
	}
}

// RemovePoolReceiveContextValueWithPackageID exercises the RemovePoolReceiveContextValue choice using the provided package ID instead of package name
func (t BurnMintTokenPool) RemovePoolReceiveContextValueWithPackageID(contractID string, packageID string, args RemovePoolReceiveContextValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "RemovePoolReceiveContextValue",
		Arguments:  argsToMap(args),
	}
}

// ClearPoolReceiveContext exercises the ClearPoolReceiveContext choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) ClearPoolReceiveContext(contractID string, args ClearPoolReceiveContext) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ClearPoolReceiveContext",
		Arguments:  argsToMap(args),
	}
}

// ClearPoolReceiveContextWithPackageID exercises the ClearPoolReceiveContext choice using the provided package ID instead of package name
func (t BurnMintTokenPool) ClearPoolReceiveContextWithPackageID(contractID string, packageID string, args ClearPoolReceiveContext) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "ClearPoolReceiveContext",
		Arguments:  argsToMap(args),
	}
}

// SetTransferTimeout exercises the SetTransferTimeout choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) SetTransferTimeout(contractID string, args SetTransferTimeout) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetTransferTimeout",
		Arguments:  argsToMap(args),
	}
}

// SetTransferTimeoutWithPackageID exercises the SetTransferTimeout choice using the provided package ID instead of package name
func (t BurnMintTokenPool) SetTransferTimeoutWithPackageID(contractID string, packageID string, args SetTransferTimeout) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "SetTransferTimeout",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVs exercises the GetRequiredCCVs choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) GetRequiredCCVs(contractID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsWithPackageID exercises the GetRequiredCCVs choice using the provided package ID instead of package name
func (t BurnMintTokenPool) GetRequiredCCVsWithPackageID(contractID string, packageID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// VerifyInboundMessage exercises the VerifyInboundMessage choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) VerifyInboundMessage(contractID string, args VerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// VerifyInboundMessageWithPackageID exercises the VerifyInboundMessage choice using the provided package ID instead of package name
func (t BurnMintTokenPool) VerifyInboundMessageWithPackageID(contractID string, packageID string, args VerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// VerifyOutboundCCVs exercises the VerifyOutboundCCVs choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) VerifyOutboundCCVs(contractID string, args VerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// VerifyOutboundCCVsWithPackageID exercises the VerifyOutboundCCVs choice using the provided package ID instead of package name
func (t BurnMintTokenPool) VerifyOutboundCCVsWithPackageID(contractID string, packageID string, args VerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// CalculateFee exercises the CalculateFee choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) CalculateFee(contractID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CalculateFeeWithPackageID exercises the CalculateFee choice using the provided package ID instead of package name
func (t BurnMintTokenPool) CalculateFeeWithPackageID(contractID string, packageID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "TokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t BurnMintTokenPool) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "TokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetFee exercises the GetFee choice on this BurnMintTokenPool contract
// This method uses the package name in the template ID
func (t BurnMintTokenPool) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t BurnMintTokenPool) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "BurnMintTokenPool"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this BurnMintTokenPool contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t BurnMintTokenPool) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetRequiredCCVs exercises the TokenPool_GetRequiredCCVs choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolGetRequiredCCVs(contractID string, args extensionapi.TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetRequiredCCVsWithPackageID exercises the TokenPool_GetRequiredCCVs choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolGetRequiredCCVsWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundMessage exercises the TokenPool_VerifyInboundMessage choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolVerifyInboundMessage(contractID string, args extensionapi.TokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundMessageWithPackageID exercises the TokenPool_VerifyInboundMessage choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolVerifyInboundMessageWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyOutboundCCVs exercises the TokenPool_VerifyOutboundCCVs choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolVerifyOutboundCCVs(contractID string, args extensionapi.TokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyOutboundCCVsWithPackageID exercises the TokenPool_VerifyOutboundCCVs choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolVerifyOutboundCCVsWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicket exercises the TokenPool_ReleaseFromTicket choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolReleaseFromTicket(contractID string, args extensionapi.TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicketWithPackageID exercises the TokenPool_ReleaseFromTicket choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolReleaseFromTicketWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurn exercises the TokenPool_LockOrBurn choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolLockOrBurn(contractID string, args extensionapi.TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurnWithPackageID exercises the TokenPool_LockOrBurn choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolLockOrBurnWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFee exercises the TokenPool_CalculateFee choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolCalculateFee(contractID string, args extensionapi.TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFeeWithPackageID exercises the TokenPool_CalculateFee choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolCalculateFeeWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetFee exercises the TokenPool_GetFee choice on this BurnMintTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t BurnMintTokenPool) TokenPoolGetFee(contractID string, args extensionapi.TokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetFeeWithPackageID exercises the TokenPool_GetFee choice using the provided package ID instead of package name
func (t BurnMintTokenPool) TokenPoolGetFeeWithPackageID(contractID string, packageID string, args extensionapi.TokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.BurnMintTokenPoolV2", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for BurnMintTokenPool

var _ api.IMCMSReceiver = (*BurnMintTokenPool)(nil)

var _ extensionapi.IITokenPool = (*BurnMintTokenPool)(nil)

// BurnMintTokenPoolDeps is a Record type
type BurnMintTokenPoolDeps struct {
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
}

// ToMap converts BurnMintTokenPoolDeps to a map for DAML arguments
func (t BurnMintTokenPoolDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	return m
}

func (t BurnMintTokenPoolDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintTokenPoolDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnMintTokenPoolDeps to hex string (Canton MCMS format)
func (t BurnMintTokenPoolDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnMintTokenPoolDeps from hex string (Canton MCMS format)
func (t *BurnMintTokenPoolDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFee is a Record type
type CalculateFee struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	FeeQuoterCid          types.CONTRACT_ID                          `json:"feeQuoterCid"`
	TokenInstrumentId     splice_api_token_holding_v1.InstrumentId   `json:"tokenInstrumentId"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts CalculateFee to a map for DAML arguments
func (t CalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type CalculateFeeMCMSParams struct {
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId   `json:"tokenInstrumentId"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	RemoteChainSelector                        types.NUMERIC                     `json:"remoteChainSelector"`
	RemotePools                                []types.TEXT                      `json:"remotePools" hex:"[]bytes"`
	RemoteTokenAddress                         types.TEXT                        `json:"remoteTokenAddress" hex:"bytes"`
	InboundCCVs                                []chainlinkapi.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []chainlinkapi.RawInstanceAddress `json:"outboundCCVs"`
	FinalityConfig                             ccipcodec.FinalityConfig          `json:"finalityConfig"`
	InboundRateLimiter                         chainlinkapi.RawInstanceAddress   `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter chainlinkapi.RawInstanceAddress   `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        chainlinkapi.RawInstanceAddress   `json:"outboundRateLimiter"`
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

// ClearPoolReceiveContext is a Record type
type ClearPoolReceiveContext struct {
}

// ToMap converts ClearPoolReceiveContext to a map for DAML arguments
func (t ClearPoolReceiveContext) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ClearPoolReceiveContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ClearPoolReceiveContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ClearPoolReceiveContext to hex string (Canton MCMS format)
func (t ClearPoolReceiveContext) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ClearPoolReceiveContext from hex string (Canton MCMS format)
func (t *ClearPoolReceiveContext) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee is a Record type
type GetFee struct {
	FeeQuoterCid      types.CONTRACT_ID                          `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId   `json:"tokenInstrumentId"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["destChainSelector"] = t.DestChainSelector

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type GetFeeMCMSParams struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId   `json:"tokenInstrumentId"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	RemoteChainSelector types.NUMERIC                              `json:"remoteChainSelector"`
	SourceAmount        types.TEXT                                 `json:"sourceAmount"`
	Finality            ccipcodec.FinalityConfig                   `json:"finality"`
	ExtraData           types.TEXT                                 `json:"extraData"`
	Direction           extensionapi.TransferDirection             `json:"direction"`
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller              types.PARTY                                `json:"caller"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["sourceAmount"] = string(t.SourceAmount)

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = model.NestedToDAMLValue(t.Direction)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type GetRequiredCCVsMCMSParams struct {
	RemoteChainSelector types.NUMERIC                              `json:"remoteChainSelector"`
	SourceAmount        types.TEXT                                 `json:"sourceAmount"`
	Finality            ccipcodec.FinalityConfig                   `json:"finality"`
	ExtraData           types.TEXT                                 `json:"extraData"`
	Direction           extensionapi.TransferDirection             `json:"direction"`
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	SenderInputCids       []types.CONTRACT_ID                        `json:"senderInputCids"`
	Amount                types.NUMERIC                              `json:"amount"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts LockOrBurn to a map for DAML arguments
func (t LockOrBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type LockOrBurnMCMSParams struct {
	SenderInputCids []types.CONTRACT_ID                        `json:"senderInputCids"`
	Amount          types.NUMERIC                              `json:"amount"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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

// RateLimitConfigArgs is a Record type
type RateLimitConfigArgs struct {
	RemoteChainSelector                        types.NUMERIC                   `json:"remoteChainSelector"`
	InboundRateLimiter                         chainlinkapi.RawInstanceAddress `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter chainlinkapi.RawInstanceAddress `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        chainlinkapi.RawInstanceAddress `json:"outboundRateLimiter"`
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
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	TokenReceiveTicketCid types.CONTRACT_ID                          `json:"tokenReceiveTicketCid"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts ReleaseFromTicket to a map for DAML arguments
func (t ReleaseFromTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type ReleaseFromTicketMCMSParams struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	RemotePools                                []types.TEXT                      `json:"remotePools" hex:"[]bytes"`
	RemoteTokenAddress                         types.TEXT                        `json:"remoteTokenAddress" hex:"bytes"`
	InboundCCVs                                []chainlinkapi.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []chainlinkapi.RawInstanceAddress `json:"outboundCCVs"`
	FinalityConfig                             ccipcodec.FinalityConfig          `json:"finalityConfig"`
	InboundRateLimiter                         chainlinkapi.RawInstanceAddress   `json:"inboundRateLimiter"`
	InboundCustomBlockConfirmationsRateLimiter chainlinkapi.RawInstanceAddress   `json:"inboundCustomBlockConfirmationsRateLimiter"`
	OutboundRateLimiter                        chainlinkapi.RawInstanceAddress   `json:"outboundRateLimiter"`
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

// RemovePoolReceiveContextValue is a Record type
type RemovePoolReceiveContextValue struct {
	ContextKey types.TEXT `json:"contextKey"`
}

// ToMap converts RemovePoolReceiveContextValue to a map for DAML arguments
func (t RemovePoolReceiveContextValue) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	return m
}

func (t RemovePoolReceiveContextValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemovePoolReceiveContextValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemovePoolReceiveContextValue to hex string (Canton MCMS format)
func (t RemovePoolReceiveContextValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemovePoolReceiveContextValue from hex string (Canton MCMS format)
func (t *RemovePoolReceiveContextValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemovePoolReceiveContextValueParams is a Record type
type RemovePoolReceiveContextValueParams struct {
	ContextKey types.TEXT `json:"contextKey"`
}

// ToMap converts RemovePoolReceiveContextValueParams to a map for DAML arguments
func (t RemovePoolReceiveContextValueParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["contextKey"] = string(t.ContextKey)

	return m
}

func (t RemovePoolReceiveContextValueParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemovePoolReceiveContextValueParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemovePoolReceiveContextValueParams to hex string (Canton MCMS format)
func (t RemovePoolReceiveContextValueParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemovePoolReceiveContextValueParams from hex string (Canton MCMS format)
func (t *RemovePoolReceiveContextValueParams) UnmarshalHex(data string) error {
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

// SetRateLimitConfig is a Record type
type SetRateLimitConfig struct {
	Caller         types.PARTY       `json:"caller"`
	RateLimiterCid types.CONTRACT_ID `json:"rateLimiterCid"`
	NewIsEnabled   types.BOOL        `json:"newIsEnabled"`
	NewCapacity    types.NUMERIC     `json:"newCapacity"`
	NewRate        types.NUMERIC     `json:"newRate"`
}

// ToMap converts SetRateLimitConfig to a map for DAML arguments
func (t SetRateLimitConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	m["rateLimiterCid"] = model.NestedToDAMLValue(t.RateLimiterCid)

	m["newIsEnabled"] = bool(t.NewIsEnabled)

	m["newCapacity"] = t.NewCapacity

	m["newRate"] = t.NewRate

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type SetRateLimitConfigMCMSParams struct {
	NewIsEnabled types.BOOL    `json:"newIsEnabled"`
	NewCapacity  types.NUMERIC `json:"newCapacity"`
	NewRate      types.NUMERIC `json:"newRate"`
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
	Caller                     types.PARTY                     `json:"caller"`
	RateLimiterInstanceAddress chainlinkapi.RawInstanceAddress `json:"rateLimiterInstanceAddress"`
	NewIsEnabled               types.BOOL                      `json:"newIsEnabled"`
	NewCapacity                types.NUMERIC                   `json:"newCapacity"`
	NewRate                    types.NUMERIC                   `json:"newRate"`
}

// ToMap converts SetRateLimitConfigParams to a map for DAML arguments
func (t SetRateLimitConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	m["rateLimiterInstanceAddress"] = model.NestedToDAMLValue(t.RateLimiterInstanceAddress)

	m["newIsEnabled"] = bool(t.NewIsEnabled)

	m["newCapacity"] = t.NewCapacity

	m["newRate"] = t.NewRate

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

// SetRateLimiterReferences is a Record type
type SetRateLimiterReferences struct {
	RateLimitConfigArgs []RateLimitConfigArgs `json:"rateLimitConfigArgs"`
}

// ToMap converts SetRateLimiterReferences to a map for DAML arguments
func (t SetRateLimiterReferences) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimitConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RateLimitConfigArgs))
		for _, e := range t.RateLimitConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t SetRateLimiterReferences) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRateLimiterReferences) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRateLimiterReferences to hex string (Canton MCMS format)
func (t SetRateLimiterReferences) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRateLimiterReferences from hex string (Canton MCMS format)
func (t *SetRateLimiterReferences) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRateLimiterReferencesParams is a Record type
type SetRateLimiterReferencesParams struct {
	RateLimitConfigArgs []RateLimitConfigArgs `json:"rateLimitConfigArgs"`
}

// ToMap converts SetRateLimiterReferencesParams to a map for DAML arguments
func (t SetRateLimiterReferencesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimitConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RateLimitConfigArgs))
		for _, e := range t.RateLimitConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t SetRateLimiterReferencesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRateLimiterReferencesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRateLimiterReferencesParams to hex string (Canton MCMS format)
func (t SetRateLimiterReferencesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRateLimiterReferencesParams from hex string (Canton MCMS format)
func (t *SetRateLimiterReferencesParams) UnmarshalHex(data string) error {
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

// TokenTransferFeeConfig is a Record type
type TokenTransferFeeConfig struct {
	IsEnabled         types.BOOL    `json:"isEnabled"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	FeeBps            types.NUMERIC `json:"feeBps"`
}

// ToMap converts TokenTransferFeeConfig to a map for DAML arguments
func (t TokenTransferFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeUSDCents"] = t.FeeUSDCents

	m["feeBps"] = t.FeeBps

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
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	ExecutingMessageCid   types.CONTRACT_ID                          `json:"executingMessageCid"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts VerifyInboundMessage to a map for DAML arguments
func (t VerifyInboundMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type VerifyInboundMessageMCMSParams struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	Amount                types.NUMERIC                              `json:"amount"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts VerifyOutboundCCVs to a map for DAML arguments
func (t VerifyOutboundCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["amount"] = t.Amount

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type VerifyOutboundCCVsMCMSParams struct {
	Amount  types.NUMERIC                              `json:"amount"`
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	AddPoolReceiveContextContractValue(args AddPoolReceiveContextContractValue) (*bind.EncodedChoice, error)
	AddPoolReceiveContextContractValueParams(args AddPoolReceiveContextContractValueParams) (*bind.EncodedChoice, error)
	AddPoolReceiveContextNonContractValue(args AddPoolReceiveContextNonContractValue) (*bind.EncodedChoice, error)
	AddPoolReceiveContextNonContractValueParams(args AddPoolReceiveContextNonContractValueParams) (*bind.EncodedChoice, error)
	ApplyChainUpdates(args ApplyChainUpdates) (*bind.EncodedChoice, error)
	ApplyChainUpdatesParams(args ApplyChainUpdatesParams) (*bind.EncodedChoice, error)
	ApplyTokenTransferFeeConfigUpdates(args ApplyTokenTransferFeeConfigUpdates) (*bind.EncodedChoice, error)
	ApplyTokenTransferFeeConfigUpdatesParams(args ApplyTokenTransferFeeConfigUpdatesParams) (*bind.EncodedChoice, error)
	CalculateFee(args CalculateFee) (*bind.EncodedChoice, error)
	CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	ClearPoolReceiveContext(args ClearPoolReceiveContext) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error)
	GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error)
	GetRequiredCCVsMCMSParams(args GetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error)
	LockOrBurn(args LockOrBurn) (*bind.EncodedChoice, error)
	LockOrBurnMCMSParams(args LockOrBurnMCMSParams) (*bind.EncodedChoice, error)
	ReleaseFromTicket(args ReleaseFromTicket) (*bind.EncodedChoice, error)
	ReleaseFromTicketMCMSParams(args ReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error)
	RemovePoolReceiveContextValue(args RemovePoolReceiveContextValue) (*bind.EncodedChoice, error)
	RemovePoolReceiveContextValueParams(args RemovePoolReceiveContextValueParams) (*bind.EncodedChoice, error)
	SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error)
	SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error)
	SetRateLimitConfig(args SetRateLimitConfig) (*bind.EncodedChoice, error)
	SetRateLimitConfigMCMSParams(args SetRateLimitConfigMCMSParams) (*bind.EncodedChoice, error)
	SetRateLimitConfigParams(args SetRateLimitConfigParams) (*bind.EncodedChoice, error)
	SetRateLimiterReferences(args SetRateLimiterReferences) (*bind.EncodedChoice, error)
	SetRateLimiterReferencesParams(args SetRateLimiterReferencesParams) (*bind.EncodedChoice, error)
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

// AddPoolReceiveContextContractValue encodes parameters for the AddPoolReceiveContextContractValue choice.
func (e *encoder) AddPoolReceiveContextContractValue(args AddPoolReceiveContextContractValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPoolReceiveContextContractValue", args)
}

// AddPoolReceiveContextContractValueParams encodes parameters for the AddPoolReceiveContextContractValue choice.
func (e *encoder) AddPoolReceiveContextContractValueParams(args AddPoolReceiveContextContractValueParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPoolReceiveContextContractValue", args)
}

// AddPoolReceiveContextNonContractValue encodes parameters for the AddPoolReceiveContextNonContractValue choice.
func (e *encoder) AddPoolReceiveContextNonContractValue(args AddPoolReceiveContextNonContractValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPoolReceiveContextNonContractValue", args)
}

// AddPoolReceiveContextNonContractValueParams encodes parameters for the AddPoolReceiveContextNonContractValue choice.
func (e *encoder) AddPoolReceiveContextNonContractValueParams(args AddPoolReceiveContextNonContractValueParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPoolReceiveContextNonContractValue", args)
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

// ClearPoolReceiveContext encodes parameters for the ClearPoolReceiveContext choice.
func (e *encoder) ClearPoolReceiveContext(args ClearPoolReceiveContext) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ClearPoolReceiveContext", args)
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

// RemovePoolReceiveContextValue encodes parameters for the RemovePoolReceiveContextValue choice.
func (e *encoder) RemovePoolReceiveContextValue(args RemovePoolReceiveContextValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemovePoolReceiveContextValue", args)
}

// RemovePoolReceiveContextValueParams encodes parameters for the RemovePoolReceiveContextValue choice.
func (e *encoder) RemovePoolReceiveContextValueParams(args RemovePoolReceiveContextValueParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemovePoolReceiveContextValue", args)
}

// SetDynamicConfig encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// SetDynamicConfigParams encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
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

// SetRateLimiterReferences encodes parameters for the SetRateLimiterReferences choice.
func (e *encoder) SetRateLimiterReferences(args SetRateLimiterReferences) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRateLimiterReferences", args)
}

// SetRateLimiterReferencesParams encodes parameters for the SetRateLimiterReferences choice.
func (e *encoder) SetRateLimiterReferencesParams(args SetRateLimiterReferencesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRateLimiterReferences", args)
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
