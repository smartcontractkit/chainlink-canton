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
	PackageID   = "76c4d4618dd9adfa08de309740a2302325988b4f55eb34f2ac450c5568ef97ba"
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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
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

// ChainUpdate is a Record type
type ChainUpdate struct {
	RemoteChainSelector                        types.NUMERIC             `json:"remoteChainSelector"`
	RemotePools                                []types.TEXT              `json:"remotePools"`
	RemoteTokenAddress                         types.TEXT                `json:"remoteTokenAddress"`
	InboundCCVs                                []mcms.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []mcms.RawInstanceAddress `json:"outboundCCVs"`
	MinBlockConfirmations                      types.INT64               `json:"minBlockConfirmations"`
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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["outboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.OutboundCCVs))
		for _, e := range t.OutboundCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["minBlockConfirmations"] = int64(t.MinBlockConfirmations)

	m["inboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundRateLimiter
	}()

	m["inboundCustomBlockConfirmationsRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundCustomBlockConfirmationsRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundCustomBlockConfirmationsRateLimiter
	}()

	m["outboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OutboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.OutboundRateLimiter
	}()

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

// LockReleaseTokenPool is a Template type
type LockReleaseTokenPool struct {
	InstanceId              types.TEXT                               `json:"instanceId"`
	PoolOwner               types.PARTY                              `json:"poolOwner"`
	CcipOwner               types.PARTY                              `json:"ccipOwner"`
	InstrumentId            splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Decimals                types.INT64                              `json:"decimals"`
	RateLimitAdmin          *types.PARTY                             `json:"rateLimitAdmin" hex:"optional"`
	RemoteChainConfigs      types.GENMAP                             `json:"remoteChainConfigs"`
	TokenTransferFeeConfigs types.GENMAP                             `json:"tokenTransferFeeConfigs"`
	PoolReceiveContext      common.CCIPContext                       `json:"poolReceiveContext"`
	TransferTimeout         TransferTimeout                          `json:"transferTimeout"`
	Deps                    LockReleaseTokenPoolDeps                 `json:"deps"`
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
	args["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

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
	args["poolReceiveContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolReceiveContext
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TransferTimeout).(mapper); ok {
			return m.toMap()
		}
		return t.TransferTimeout
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

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
	args["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

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
	args["poolReceiveContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolReceiveContext
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TransferTimeout).(mapper); ok {
			return m.toMap()
		}
		return t.TransferTimeout
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

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

// LockReleaseTokenPoolReleaseFromTicket exercises the LockReleaseTokenPool_ReleaseFromTicket choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolReleaseFromTicket(contractID string, args LockReleaseTokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolReleaseFromTicketWithPackageID exercises the LockReleaseTokenPool_ReleaseFromTicket choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolReleaseFromTicketWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolLockOrBurn exercises the LockReleaseTokenPool_LockOrBurn choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolLockOrBurn(contractID string, args LockReleaseTokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolLockOrBurnWithPackageID exercises the LockReleaseTokenPool_LockOrBurn choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolLockOrBurnWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_LockOrBurn",
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

// LockReleaseTokenPoolVerifyInboundMessage exercises the LockReleaseTokenPool_VerifyInboundMessage choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyInboundMessage(contractID string, args LockReleaseTokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolVerifyInboundMessageWithPackageID exercises the LockReleaseTokenPool_VerifyInboundMessage choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyInboundMessageWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolVerifyInboundMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyInboundMessage",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolVerifyOutboundCCVs exercises the LockReleaseTokenPool_VerifyOutboundCCVs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyOutboundCCVs(contractID string, args LockReleaseTokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolVerifyOutboundCCVsWithPackageID exercises the LockReleaseTokenPool_VerifyOutboundCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyOutboundCCVsWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolVerifyOutboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyOutboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolGetFee exercises the LockReleaseTokenPool_GetFee choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolGetFee(contractID string, args LockReleaseTokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolGetFeeWithPackageID exercises the LockReleaseTokenPool_GetFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolGetFeeWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_GetFee",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolGetRequiredCCVs exercises the LockReleaseTokenPool_GetRequiredCCVs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolGetRequiredCCVs(contractID string, args LockReleaseTokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolGetRequiredCCVsWithPackageID exercises the LockReleaseTokenPool_GetRequiredCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolGetRequiredCCVsWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_GetRequiredCCVs",
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

	m["tokenAdminRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistry
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["feeQuoter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoter).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoter
	}()

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

// LockReleaseTokenPoolGetFee is a Record type
type LockReleaseTokenPoolGetFee struct {
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolGetFee to a map for DAML arguments
func (t LockReleaseTokenPoolGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeQuoterCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

	m["destChainSelector"] = t.DestChainSelector

	m["tokenInstrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInstrumentId
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolGetFee to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolGetFee from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolGetFeeMCMSParams is LockReleaseTokenPoolGetFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolGetFeeMCMSParams struct {
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
}

// MarshalHex encodes LockReleaseTokenPoolGetFeeMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolGetFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolGetFeeMCMSParams from hex string.
func (t *LockReleaseTokenPoolGetFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolGetRequiredCCVs is a Record type
type LockReleaseTokenPoolGetRequiredCCVs struct {
	RemoteChainSelector types.NUMERIC                `json:"remoteChainSelector"`
	Amount              types.NUMERIC                `json:"amount"`
	Finality            types.INT64                  `json:"finality"`
	ExtraData           types.TEXT                   `json:"extraData"`
	Direction           interfaces.TransferDirection `json:"direction"`
	Caller              types.PARTY                  `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolGetRequiredCCVs to a map for DAML arguments
func (t LockReleaseTokenPoolGetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["amount"] = t.Amount

	m["finality"] = int64(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Direction).(mapper); ok {
			return m.toMap()
		}
		return t.Direction
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolGetRequiredCCVs to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolGetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolGetRequiredCCVs from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolGetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolGetRequiredCCVsMCMSParams is LockReleaseTokenPoolGetRequiredCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolGetRequiredCCVsMCMSParams struct {
	RemoteChainSelector types.NUMERIC                `json:"remoteChainSelector"`
	Amount              types.NUMERIC                `json:"amount"`
	Finality            types.INT64                  `json:"finality"`
	ExtraData           types.TEXT                   `json:"extraData"`
	Direction           interfaces.TransferDirection `json:"direction"`
}

// MarshalHex encodes LockReleaseTokenPoolGetRequiredCCVsMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolGetRequiredCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolGetRequiredCCVsMCMSParams from hex string.
func (t *LockReleaseTokenPoolGetRequiredCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolLockOrBurn is a Record type
type LockReleaseTokenPoolLockOrBurn struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids       []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount                types.NUMERIC         `json:"amount"`
	Caller                types.PARTY           `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolLockOrBurn to a map for DAML arguments
func (t LockReleaseTokenPoolLockOrBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

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

func (t LockReleaseTokenPoolLockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolLockOrBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolLockOrBurn to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolLockOrBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolLockOrBurn from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolLockOrBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolLockOrBurnMCMSParams is LockReleaseTokenPoolLockOrBurn without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolLockOrBurnMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids       []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount                types.NUMERIC         `json:"amount"`
}

// MarshalHex encodes LockReleaseTokenPoolLockOrBurnMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolLockOrBurnMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolLockOrBurnMCMSParams from hex string.
func (t *LockReleaseTokenPoolLockOrBurnMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolReleaseFromTicket is a Record type
type LockReleaseTokenPoolReleaseFromTicket struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	TokenReceiveTicketCid types.CONTRACT_ID     `json:"tokenReceiveTicketCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
	Caller                types.PARTY           `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolReleaseFromTicket to a map for DAML arguments
func (t LockReleaseTokenPoolReleaseFromTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["tokenReceiveTicketCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenReceiveTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenReceiveTicketCid
	}()

	m["tokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolReleaseFromTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolReleaseFromTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolReleaseFromTicket to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolReleaseFromTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolReleaseFromTicket from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolReleaseFromTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolReleaseFromTicketMCMSParams is LockReleaseTokenPoolReleaseFromTicket without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolReleaseFromTicketMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID     `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext          common.CCIPContext    `json:"extraContext"`
	TokenReceiveTicketCid types.CONTRACT_ID     `json:"tokenReceiveTicketCid"`
	TokenInput            interfaces.TokenInput `json:"tokenInput"`
}

// MarshalHex encodes LockReleaseTokenPoolReleaseFromTicketMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolReleaseFromTicketMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolReleaseFromTicketMCMSParams from hex string.
func (t *LockReleaseTokenPoolReleaseFromTicketMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolVerifyInboundMessage is a Record type
type LockReleaseTokenPoolVerifyInboundMessage struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid   types.CONTRACT_ID  `json:"executingMessageCid"`
	Caller                types.PARTY        `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolVerifyInboundMessage to a map for DAML arguments
func (t LockReleaseTokenPoolVerifyInboundMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolVerifyInboundMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolVerifyInboundMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolVerifyInboundMessage to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolVerifyInboundMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolVerifyInboundMessage from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolVerifyInboundMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolVerifyInboundMessageMCMSParams is LockReleaseTokenPoolVerifyInboundMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolVerifyInboundMessageMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid   types.CONTRACT_ID  `json:"executingMessageCid"`
}

// MarshalHex encodes LockReleaseTokenPoolVerifyInboundMessageMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolVerifyInboundMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolVerifyInboundMessageMCMSParams from hex string.
func (t *LockReleaseTokenPoolVerifyInboundMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolVerifyOutboundCCVs is a Record type
type LockReleaseTokenPoolVerifyOutboundCCVs struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID  `json:"sendingMessageCid"`
	Amount                types.NUMERIC      `json:"amount"`
	Caller                types.PARTY        `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolVerifyOutboundCCVs to a map for DAML arguments
func (t LockReleaseTokenPoolVerifyOutboundCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["amount"] = t.Amount

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolVerifyOutboundCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolVerifyOutboundCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolVerifyOutboundCCVs to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolVerifyOutboundCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolVerifyOutboundCCVs from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolVerifyOutboundCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams is LockReleaseTokenPoolVerifyOutboundCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams struct {
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	ExtraContext          common.CCIPContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID  `json:"sendingMessageCid"`
	Amount                types.NUMERIC      `json:"amount"`
}

// MarshalHex encodes LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams from hex string.
func (t *LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams) UnmarshalHex(data string) error {
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

	m["inboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundRateLimiter
	}()

	m["inboundCustomBlockConfirmationsRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundCustomBlockConfirmationsRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundCustomBlockConfirmationsRateLimiter
	}()

	m["outboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OutboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.OutboundRateLimiter
	}()

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

// RemoteChainConfig is a Record type
type RemoteChainConfig struct {
	RemotePools                                []types.TEXT              `json:"remotePools"`
	RemoteTokenAddress                         types.TEXT                `json:"remoteTokenAddress"`
	InboundCCVs                                []mcms.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs                               []mcms.RawInstanceAddress `json:"outboundCCVs"`
	MinBlockConfirmations                      types.INT64               `json:"minBlockConfirmations"`
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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["outboundCCVs"] = func() []any {
		res := make([]any, 0, len(t.OutboundCCVs))
		for _, e := range t.OutboundCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["minBlockConfirmations"] = int64(t.MinBlockConfirmations)

	m["inboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundRateLimiter
	}()

	m["inboundCustomBlockConfirmationsRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InboundCustomBlockConfirmationsRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.InboundCustomBlockConfirmationsRateLimiter
	}()

	m["outboundRateLimiter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OutboundRateLimiter).(mapper); ok {
			return m.toMap()
		}
		return t.OutboundRateLimiter
	}()

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

// SetPoolReceiveContext is a Record type
type SetPoolReceiveContext struct {
	NewPoolReceiveContext common.CCIPContext `json:"newPoolReceiveContext"`
}

// ToMap converts SetPoolReceiveContext to a map for DAML arguments
func (t SetPoolReceiveContext) ToMap() map[string]any {
	m := make(map[string]any)

	m["newPoolReceiveContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewPoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.NewPoolReceiveContext
	}()

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

// SetTransferTimeout is a Record type
type SetTransferTimeout struct {
	NewTransferTimeout TransferTimeout `json:"newTransferTimeout"`
}

// ToMap converts SetTransferTimeout to a map for DAML arguments
func (t SetTransferTimeout) ToMap() map[string]any {
	m := make(map[string]any)

	m["newTransferTimeout"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewTransferTimeout).(mapper); ok {
			return m.toMap()
		}
		return t.NewTransferTimeout
	}()

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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ApplyChainUpdates(args ApplyChainUpdates) (*bind.EncodedChoice, error)
	ApplyTokenTransferFeeConfigUpdates(args ApplyTokenTransferFeeConfigUpdates) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetFee(args LockReleaseTokenPoolGetFee) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetFeeMCMSParams(args LockReleaseTokenPoolGetFeeMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetRequiredCCVs(args LockReleaseTokenPoolGetRequiredCCVs) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetRequiredCCVsMCMSParams(args LockReleaseTokenPoolGetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolLockOrBurn(args LockReleaseTokenPoolLockOrBurn) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolLockOrBurnMCMSParams(args LockReleaseTokenPoolLockOrBurnMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolReleaseFromTicket(args LockReleaseTokenPoolReleaseFromTicket) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolReleaseFromTicketMCMSParams(args LockReleaseTokenPoolReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyInboundMessage(args LockReleaseTokenPoolVerifyInboundMessage) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyInboundMessageMCMSParams(args LockReleaseTokenPoolVerifyInboundMessageMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyOutboundCCVs(args LockReleaseTokenPoolVerifyOutboundCCVs) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams(args LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams) (*bind.EncodedChoice, error)
	SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error)
	SetPoolReceiveContext(args SetPoolReceiveContext) (*bind.EncodedChoice, error)
	SetRateLimitConfig(args SetRateLimitConfig) (*bind.EncodedChoice, error)
	SetRateLimitConfigMCMSParams(args SetRateLimitConfigMCMSParams) (*bind.EncodedChoice, error)
	SetTransferTimeout(args SetTransferTimeout) (*bind.EncodedChoice, error)
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

// ApplyTokenTransferFeeConfigUpdates encodes parameters for the ApplyTokenTransferFeeConfigUpdates choice.
func (e *encoder) ApplyTokenTransferFeeConfigUpdates(args ApplyTokenTransferFeeConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyTokenTransferFeeConfigUpdates", args)
}

// LockReleaseTokenPoolGetFee encodes parameters for the LockReleaseTokenPoolGetFee choice.
func (e *encoder) LockReleaseTokenPoolGetFee(args LockReleaseTokenPoolGetFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolGetFee", args)
}

// LockReleaseTokenPoolGetFeeMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolGetFee choice.
func (e *encoder) LockReleaseTokenPoolGetFeeMCMSParams(args LockReleaseTokenPoolGetFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolGetFee", args)
}

// LockReleaseTokenPoolGetRequiredCCVs encodes parameters for the LockReleaseTokenPoolGetRequiredCCVs choice.
func (e *encoder) LockReleaseTokenPoolGetRequiredCCVs(args LockReleaseTokenPoolGetRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolGetRequiredCCVs", args)
}

// LockReleaseTokenPoolGetRequiredCCVsMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolGetRequiredCCVs choice.
func (e *encoder) LockReleaseTokenPoolGetRequiredCCVsMCMSParams(args LockReleaseTokenPoolGetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolGetRequiredCCVs", args)
}

// LockReleaseTokenPoolLockOrBurn encodes parameters for the LockReleaseTokenPoolLockOrBurn choice.
func (e *encoder) LockReleaseTokenPoolLockOrBurn(args LockReleaseTokenPoolLockOrBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolLockOrBurn", args)
}

// LockReleaseTokenPoolLockOrBurnMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolLockOrBurn choice.
func (e *encoder) LockReleaseTokenPoolLockOrBurnMCMSParams(args LockReleaseTokenPoolLockOrBurnMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolLockOrBurn", args)
}

// LockReleaseTokenPoolReleaseFromTicket encodes parameters for the LockReleaseTokenPoolReleaseFromTicket choice.
func (e *encoder) LockReleaseTokenPoolReleaseFromTicket(args LockReleaseTokenPoolReleaseFromTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolReleaseFromTicket", args)
}

// LockReleaseTokenPoolReleaseFromTicketMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolReleaseFromTicket choice.
func (e *encoder) LockReleaseTokenPoolReleaseFromTicketMCMSParams(args LockReleaseTokenPoolReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolReleaseFromTicket", args)
}

// LockReleaseTokenPoolVerifyInboundMessage encodes parameters for the LockReleaseTokenPoolVerifyInboundMessage choice.
func (e *encoder) LockReleaseTokenPoolVerifyInboundMessage(args LockReleaseTokenPoolVerifyInboundMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolVerifyInboundMessage", args)
}

// LockReleaseTokenPoolVerifyInboundMessageMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolVerifyInboundMessage choice.
func (e *encoder) LockReleaseTokenPoolVerifyInboundMessageMCMSParams(args LockReleaseTokenPoolVerifyInboundMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolVerifyInboundMessage", args)
}

// LockReleaseTokenPoolVerifyOutboundCCVs encodes parameters for the LockReleaseTokenPoolVerifyOutboundCCVs choice.
func (e *encoder) LockReleaseTokenPoolVerifyOutboundCCVs(args LockReleaseTokenPoolVerifyOutboundCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolVerifyOutboundCCVs", args)
}

// LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolVerifyOutboundCCVs choice.
func (e *encoder) LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams(args LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolVerifyOutboundCCVs", args)
}

// SetDynamicConfig encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error) {
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

// SetTransferTimeout encodes parameters for the SetTransferTimeout choice.
func (e *encoder) SetTransferTimeout(args SetTransferTimeout) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferTimeout", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
