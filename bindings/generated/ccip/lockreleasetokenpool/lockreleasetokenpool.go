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
	PackageID   = "dce10828d896953f04ce6cecf89a8a47a04ac9b2307c95c78ab664a80f16fc5d"
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

// ChainPoolConfig is a Record type
type ChainPoolConfig struct {
	InboundCCVs  []common.RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs []common.RawInstanceAddress `json:"outboundCCVs"`
	RemotePools  []types.TEXT                `json:"remotePools"`
}

// ToMap converts ChainPoolConfig to a map for DAML arguments
func (t ChainPoolConfig) ToMap() map[string]any {
	m := make(map[string]any)

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

	m["remotePools"] = func() []any {
		res := make([]any, 0, len(t.RemotePools))
		for _, e := range t.RemotePools {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t ChainPoolConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ChainPoolConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ChainPoolConfig to hex string (Canton MCMS format)
func (t ChainPoolConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ChainPoolConfig from hex string (Canton MCMS format)
func (t *ChainPoolConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPool is a Template type
type LockReleaseTokenPool struct {
	InstanceId           types.TEXT                               `json:"instanceId"`
	CcipOwner            types.PARTY                              `json:"ccipOwner"`
	PoolOwner            types.PARTY                              `json:"poolOwner"`
	InstrumentId         splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Decimals             types.INT64                              `json:"decimals"`
	ChainPoolConfigs     types.GENMAP                             `json:"chainPoolConfigs"`
	ChainFeeConfigs      types.GENMAP                             `json:"chainFeeConfigs"`
	RemoteTokens         types.GENMAP                             `json:"remoteTokens"`
	OutboundRateLimiters types.GENMAP                             `json:"outboundRateLimiters"`
	InboundRateLimiters  types.GENMAP                             `json:"inboundRateLimiters"`
	PoolReceiveContext   common.CCIPContext                       `json:"poolReceiveContext"`
	TransferTimeout      TransferTimeout                          `json:"transferTimeout"`
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
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

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

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainPoolConfigs"] = func() any {
		if t.ChainPoolConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ChainPoolConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainFeeConfigs"] = func() any {
		if t.ChainFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ChainFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteTokens"] = func() any {
		if t.RemoteTokens == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundRateLimiters"] = func() any {
		if t.OutboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundRateLimiters}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["inboundRateLimiters"] = func() any {
		if t.InboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.InboundRateLimiters}
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
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

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

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainPoolConfigs"] = func() any {
		if t.ChainPoolConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ChainPoolConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainFeeConfigs"] = func() any {
		if t.ChainFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ChainFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteTokens"] = func() any {
		if t.RemoteTokens == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundRateLimiters"] = func() any {
		if t.OutboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundRateLimiters}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["inboundRateLimiters"] = func() any {
		if t.InboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.InboundRateLimiters}
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

// LockReleaseTokenPoolCalculateFee exercises the LockReleaseTokenPool_CalculateFee choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolCalculateFee(contractID string, args LockReleaseTokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolCalculateFeeWithPackageID exercises the LockReleaseTokenPool_CalculateFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolCalculateFeeWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
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

// LockReleaseTokenPoolUpdateChainPoolConfigs exercises the LockReleaseTokenPool_UpdateChainPoolConfigs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateChainPoolConfigs(contractID string, args LockReleaseTokenPoolUpdateChainPoolConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateChainPoolConfigs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolUpdateChainPoolConfigsWithPackageID exercises the LockReleaseTokenPool_UpdateChainPoolConfigs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateChainPoolConfigsWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolUpdateChainPoolConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateChainPoolConfigs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolUpdateRateLimiters exercises the LockReleaseTokenPool_UpdateRateLimiters choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateRateLimiters(contractID string, args LockReleaseTokenPoolUpdateRateLimiters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateRateLimiters",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolUpdateRateLimitersWithPackageID exercises the LockReleaseTokenPool_UpdateRateLimiters choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateRateLimitersWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolUpdateRateLimiters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateRateLimiters",
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

// Verify interface implementations for LockReleaseTokenPool

var _ mcms.IMCMSReceiver = (*LockReleaseTokenPool)(nil)

var _ interfaces.IITokenPool = (*LockReleaseTokenPool)(nil)

// LockReleaseTokenPoolCalculateFee is a Record type
type LockReleaseTokenPoolCalculateFee struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolCalculateFee to a map for DAML arguments
func (t LockReleaseTokenPoolCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["feeQuoterCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

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

func (t LockReleaseTokenPoolCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolCalculateFee to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolCalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolCalculateFee from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolCalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolCalculateFeeMCMSParams is LockReleaseTokenPoolCalculateFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type LockReleaseTokenPoolCalculateFeeMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
}

// MarshalHex encodes LockReleaseTokenPoolCalculateFeeMCMSParams to hex string for MCMS operationData.
func (t LockReleaseTokenPoolCalculateFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolCalculateFeeMCMSParams from hex string.
func (t *LockReleaseTokenPoolCalculateFeeMCMSParams) UnmarshalHex(data string) error {
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
	RmnRemoteCid      types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext      common.CCIPContext    `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput        interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids   []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount            types.NUMERIC         `json:"amount"`
	Caller            types.PARTY           `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolLockOrBurn to a map for DAML arguments
func (t LockReleaseTokenPoolLockOrBurn) ToMap() map[string]any {
	m := make(map[string]any)

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
	RmnRemoteCid      types.CONTRACT_ID     `json:"rmnRemoteCid"`
	ExtraContext      common.CCIPContext    `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID     `json:"sendingMessageCid"`
	TokenInput        interfaces.TokenInput `json:"tokenInput"`
	SenderInputCids   []types.CONTRACT_ID   `json:"senderInputCids"`
	Amount            types.NUMERIC         `json:"amount"`
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

// LockReleaseTokenPoolUpdateChainPoolConfigs is a Record type
type LockReleaseTokenPoolUpdateChainPoolConfigs struct {
	NewChainPoolConfigs types.GENMAP `json:"newChainPoolConfigs"`
}

// ToMap converts LockReleaseTokenPoolUpdateChainPoolConfigs to a map for DAML arguments
func (t LockReleaseTokenPoolUpdateChainPoolConfigs) ToMap() map[string]any {
	m := make(map[string]any)

	m["newChainPoolConfigs"] = func() any {
		if t.NewChainPoolConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.NewChainPoolConfigs}
	}()

	return m
}

func (t LockReleaseTokenPoolUpdateChainPoolConfigs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolUpdateChainPoolConfigs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolUpdateChainPoolConfigs to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolUpdateChainPoolConfigs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolUpdateChainPoolConfigs from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolUpdateChainPoolConfigs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockReleaseTokenPoolUpdateRateLimiters is a Record type
type LockReleaseTokenPoolUpdateRateLimiters struct {
	NewOutboundRateLimiters types.GENMAP `json:"newOutboundRateLimiters"`
	NewInboundRateLimiters  types.GENMAP `json:"newInboundRateLimiters"`
}

// ToMap converts LockReleaseTokenPoolUpdateRateLimiters to a map for DAML arguments
func (t LockReleaseTokenPoolUpdateRateLimiters) ToMap() map[string]any {
	m := make(map[string]any)

	m["newOutboundRateLimiters"] = func() any {
		if t.NewOutboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.NewOutboundRateLimiters}
	}()

	m["newInboundRateLimiters"] = func() any {
		if t.NewInboundRateLimiters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.NewInboundRateLimiters}
	}()

	return m
}

func (t LockReleaseTokenPoolUpdateRateLimiters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockReleaseTokenPoolUpdateRateLimiters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockReleaseTokenPoolUpdateRateLimiters to hex string (Canton MCMS format)
func (t LockReleaseTokenPoolUpdateRateLimiters) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockReleaseTokenPoolUpdateRateLimiters from hex string (Canton MCMS format)
func (t *LockReleaseTokenPoolUpdateRateLimiters) UnmarshalHex(data string) error {
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

// PoolFeeConfig is a Record type
type PoolFeeConfig struct {
	IsEnabled         types.BOOL    `json:"isEnabled"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts PoolFeeConfig to a map for DAML arguments
func (t PoolFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t PoolFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PoolFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PoolFeeConfig to hex string (Canton MCMS format)
func (t PoolFeeConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PoolFeeConfig from hex string (Canton MCMS format)
func (t *PoolFeeConfig) UnmarshalHex(data string) error {
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
	LockReleaseTokenPoolCalculateFee(args LockReleaseTokenPoolCalculateFee) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolCalculateFeeMCMSParams(args LockReleaseTokenPoolCalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetRequiredCCVs(args LockReleaseTokenPoolGetRequiredCCVs) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolGetRequiredCCVsMCMSParams(args LockReleaseTokenPoolGetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolLockOrBurn(args LockReleaseTokenPoolLockOrBurn) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolLockOrBurnMCMSParams(args LockReleaseTokenPoolLockOrBurnMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolReleaseFromTicket(args LockReleaseTokenPoolReleaseFromTicket) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolReleaseFromTicketMCMSParams(args LockReleaseTokenPoolReleaseFromTicketMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolUpdateChainPoolConfigs(args LockReleaseTokenPoolUpdateChainPoolConfigs) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolUpdateRateLimiters(args LockReleaseTokenPoolUpdateRateLimiters) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyInboundMessage(args LockReleaseTokenPoolVerifyInboundMessage) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyInboundMessageMCMSParams(args LockReleaseTokenPoolVerifyInboundMessageMCMSParams) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyOutboundCCVs(args LockReleaseTokenPoolVerifyOutboundCCVs) (*bind.EncodedChoice, error)
	LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams(args LockReleaseTokenPoolVerifyOutboundCCVsMCMSParams) (*bind.EncodedChoice, error)
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

// LockReleaseTokenPoolCalculateFee encodes parameters for the LockReleaseTokenPoolCalculateFee choice.
func (e *encoder) LockReleaseTokenPoolCalculateFee(args LockReleaseTokenPoolCalculateFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolCalculateFee", args)
}

// LockReleaseTokenPoolCalculateFeeMCMSParams encodes MCMS parameters (without Caller) for the LockReleaseTokenPoolCalculateFee choice.
func (e *encoder) LockReleaseTokenPoolCalculateFeeMCMSParams(args LockReleaseTokenPoolCalculateFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolCalculateFee", args)
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

// LockReleaseTokenPoolUpdateChainPoolConfigs encodes parameters for the LockReleaseTokenPoolUpdateChainPoolConfigs choice.
func (e *encoder) LockReleaseTokenPoolUpdateChainPoolConfigs(args LockReleaseTokenPoolUpdateChainPoolConfigs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolUpdateChainPoolConfigs", args)
}

// LockReleaseTokenPoolUpdateRateLimiters encodes parameters for the LockReleaseTokenPoolUpdateRateLimiters choice.
func (e *encoder) LockReleaseTokenPoolUpdateRateLimiters(args LockReleaseTokenPoolUpdateRateLimiters) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockReleaseTokenPoolUpdateRateLimiters", args)
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

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
