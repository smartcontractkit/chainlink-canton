package lockreleasetokenpool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageName = "ccip-lockreleasetokenpool"
const SDKVersion = "3.4.10"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	type mapper interface {
		ToMap() map[string]interface{}
	}
	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{"args": args}
}

// ChainCCVConfig is a Record type
type ChainCCVConfig struct {
	InboundCCVs  []RawInstanceAddress `json:"inboundCCVs"`
	OutboundCCVs []RawInstanceAddress `json:"outboundCCVs"`
}

// ToMap converts ChainCCVConfig to a map for DAML arguments
func (t ChainCCVConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["inboundCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.InboundCCVs))
		for _, e := range t.InboundCCVs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["outboundCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OutboundCCVs))
		for _, e := range t.OutboundCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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

func (t ChainCCVConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ChainCCVConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPool is a Template type
type LockReleaseTokenPool struct {
	InstanceId           TEXT            `json:"instanceId"`
	CcipOwner            PARTY           `json:"ccipOwner"`
	PoolOwner            PARTY           `json:"poolOwner"`
	InstrumentId         InstrumentId    `json:"instrumentId"`
	Decimals             INT64           `json:"decimals"`
	ChainCCVRequirements GENMAP          `json:"chainCCVRequirements"`
	ChainFeeConfigs      GENMAP          `json:"chainFeeConfigs"`
	RemoteTokens         GENMAP          `json:"remoteTokens"`
	PoolReceiveContext   ChoiceContext   `json:"poolReceiveContext"`
	TransferTimeout      TransferTimeout `json:"transferTimeout"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["decimals"] = int64(t.Decimals)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainCCVRequirements"] = func() interface{} {
		if t.ChainCCVRequirements == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ChainCCVRequirements}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainFeeConfigs"] = func() interface{} {
		if t.ChainFeeConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ChainFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteTokens"] = func() interface{} {
		if t.RemoteTokens == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RemoteTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolReceiveContext"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.PoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolReceiveContext
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["decimals"] = int64(t.Decimals)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainCCVRequirements"] = func() interface{} {
		if t.ChainCCVRequirements == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ChainCCVRequirements}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainFeeConfigs"] = func() interface{} {
		if t.ChainFeeConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ChainFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteTokens"] = func() interface{} {
		if t.RemoteTokens == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RemoteTokens}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolReceiveContext"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.PoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolReceiveContext
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferTimeout"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for LockReleaseTokenPool

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

// LockReleaseTokenPoolVerifyInboundCCVs exercises the LockReleaseTokenPool_VerifyInboundCCVs choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyInboundCCVs(contractID string, args LockReleaseTokenPoolVerifyInboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyInboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolVerifyInboundCCVsWithPackageID exercises the LockReleaseTokenPool_VerifyInboundCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolVerifyInboundCCVsWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolVerifyInboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_VerifyInboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolUpdateChainCCVRequirements exercises the LockReleaseTokenPool_UpdateChainCCVRequirements choice on this LockReleaseTokenPool contract
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateChainCCVRequirements(contractID string, args LockReleaseTokenPoolUpdateChainCCVRequirements) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateChainCCVRequirements",
		Arguments:  argsToMap(args),
	}
}

// LockReleaseTokenPoolUpdateChainCCVRequirementsWithPackageID exercises the LockReleaseTokenPool_UpdateChainCCVRequirements choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) LockReleaseTokenPoolUpdateChainCCVRequirementsWithPackageID(contractID string, packageID string, args LockReleaseTokenPoolUpdateChainCCVRequirements) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "LockReleaseTokenPool_UpdateChainCCVRequirements",
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
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
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

// TokenPoolGetRequiredCCVs exercises the TokenPool_GetRequiredCCVs choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolGetRequiredCCVs(contractID string, args TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolGetRequiredCCVsWithPackageID exercises the TokenPool_GetRequiredCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolGetRequiredCCVsWithPackageID(contractID string, packageID string, args TokenPoolGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundCCVs exercises the TokenPool_VerifyInboundCCVs choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolVerifyInboundCCVs(contractID string, args TokenPoolVerifyInboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolVerifyInboundCCVsWithPackageID exercises the TokenPool_VerifyInboundCCVs choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolVerifyInboundCCVsWithPackageID(contractID string, packageID string, args TokenPoolVerifyInboundCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_VerifyInboundCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicket exercises the TokenPool_ReleaseFromTicket choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolReleaseFromTicket(contractID string, args TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolReleaseFromTicketWithPackageID exercises the TokenPool_ReleaseFromTicket choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolReleaseFromTicketWithPackageID(contractID string, packageID string, args TokenPoolReleaseFromTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_ReleaseFromTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurn exercises the TokenPool_LockOrBurn choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolLockOrBurn(contractID string, args TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolLockOrBurnWithPackageID exercises the TokenPool_LockOrBurn choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolLockOrBurnWithPackageID(contractID string, packageID string, args TokenPoolLockOrBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_LockOrBurn",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFee exercises the TokenPool_CalculateFee choice on this LockReleaseTokenPool contract via the IITokenPool interface
// This method uses the package name in the template ID
func (t LockReleaseTokenPool) TokenPoolCalculateFee(contractID string, args TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// TokenPoolCalculateFeeWithPackageID exercises the TokenPool_CalculateFee choice using the provided package ID instead of package name
func (t LockReleaseTokenPool) TokenPoolCalculateFeeWithPackageID(contractID string, packageID string, args TokenPoolCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.LockReleaseTokenPool", "ITokenPool"),
		ContractID: contractID,
		Choice:     "TokenPool_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for LockReleaseTokenPool

var _ IITokenPool = (*LockReleaseTokenPool)(nil)

// LockReleaseTokenPoolCalculateFee is a Record type
type LockReleaseTokenPoolCalculateFee struct {
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	Caller            PARTY       `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolCalculateFee to a map for DAML arguments
func (t LockReleaseTokenPoolCalculateFee) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPoolGetRequiredCCVs is a Record type
type LockReleaseTokenPoolGetRequiredCCVs struct {
	RemoteChainSelector NUMERIC           `json:"remoteChainSelector"`
	Amount              NUMERIC           `json:"amount"`
	Finality            INT64             `json:"finality"`
	ExtraData           TEXT              `json:"extraData"`
	Direction           TransferDirection `json:"direction"`
	Caller              PARTY             `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolGetRequiredCCVs to a map for DAML arguments
func (t LockReleaseTokenPoolGetRequiredCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["amount"] = t.Amount

	m["finality"] = int64(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPoolLockOrBurn is a Record type
type LockReleaseTokenPoolLockOrBurn struct {
	SendingMessageCid CONTRACT_ID   `json:"sendingMessageCid"`
	TokenInput        TokenInput    `json:"tokenInput"`
	SenderInputCids   []CONTRACT_ID `json:"senderInputCids"`
	Amount            NUMERIC       `json:"amount"`
	RmnRemoteCid      CONTRACT_ID   `json:"rmnRemoteCid"`
	Caller            PARTY         `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolLockOrBurn to a map for DAML arguments
func (t LockReleaseTokenPoolLockOrBurn) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["senderInputCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolLockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolLockOrBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPoolReleaseFromTicket is a Record type
type LockReleaseTokenPoolReleaseFromTicket struct {
	TokenReceiveTicketCid CONTRACT_ID `json:"tokenReceiveTicketCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          CONTRACT_ID `json:"rmnRemoteCid"`
	TokenInput            TokenInput  `json:"tokenInput"`
	Caller                PARTY       `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolReleaseFromTicket to a map for DAML arguments
func (t LockReleaseTokenPoolReleaseFromTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["tokenReceiveTicketCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenReceiveTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenReceiveTicketCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolReleaseFromTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPoolUpdateChainCCVRequirements is a Record type
type LockReleaseTokenPoolUpdateChainCCVRequirements struct {
	NewChainCCVRequirements GENMAP `json:"newChainCCVRequirements"`
}

// ToMap converts LockReleaseTokenPoolUpdateChainCCVRequirements to a map for DAML arguments
func (t LockReleaseTokenPoolUpdateChainCCVRequirements) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newChainCCVRequirements"] = func() interface{} {
		if t.NewChainCCVRequirements == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.NewChainCCVRequirements}
	}()

	return m
}

func (t LockReleaseTokenPoolUpdateChainCCVRequirements) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolUpdateChainCCVRequirements) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockReleaseTokenPoolVerifyInboundCCVs is a Record type
type LockReleaseTokenPoolVerifyInboundCCVs struct {
	ExecutingMessageCid   CONTRACT_ID `json:"executingMessageCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
	Caller                PARTY       `json:"caller"`
}

// ToMap converts LockReleaseTokenPoolVerifyInboundCCVs to a map for DAML arguments
func (t LockReleaseTokenPoolVerifyInboundCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t LockReleaseTokenPoolVerifyInboundCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *LockReleaseTokenPoolVerifyInboundCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PoolFeeConfig is a Record type
type PoolFeeConfig struct {
	FeeUSDCents       NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   INT64   `json:"destGasOverhead"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
}

// ToMap converts PoolFeeConfig to a map for DAML arguments
func (t PoolFeeConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t PoolFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PoolFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferTimeout is a variant/union type
type TransferTimeout struct {
	Indefinite    *UNIT  `json:"Indefinite,omitempty"`
	RelativeHours *INT64 `json:"RelativeHours,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferTimeout
func (v TransferTimeout) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferTimeout
func (v *TransferTimeout) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
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
func (v TransferTimeout) GetVariantValue() interface{} {

	if v.Indefinite != nil {
		return v.Indefinite
	}

	if v.RelativeHours != nil {
		return v.RelativeHours
	}

	return nil
}

var _ VARIANT = (*TransferTimeout)(nil)
