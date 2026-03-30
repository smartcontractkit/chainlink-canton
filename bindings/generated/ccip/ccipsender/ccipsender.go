package ccipsender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	client "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
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
	PackageName = "ccip-sender"
	PackageID   = "9973959449d84e140adebe87c20d77df82265d175ce0019e4eb53f7b494b104b"
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

// CCIPSender is a Template type
type CCIPSender struct {
	InstanceId types.TEXT  `json:"instanceId"`
	Owner      types.PARTY `json:"owner"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPSender) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPSender) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPSender) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPSender) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPSender) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSender) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSender to hex string (Canton MCMS format)
func (t CCIPSender) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSender from hex string (Canton MCMS format)
func (t *CCIPSender) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPSender

// GetFee exercises the GetFee choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) GetFee(contractID string, args GetFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t CCIPSender) GetFeeWithPackageID(contractID string, packageID string, args GetFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPSender) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Send exercises the Send choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) Send(contractID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// SendWithPackageID exercises the Send choice using the provided package ID instead of package name
func (t CCIPSender) SendWithPackageID(contractID string, packageID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// GatheredQuote is a Record type
type GatheredQuote struct {
	QuotedCCVFees       []common.CCVFee      `json:"quotedCCVFees"`
	QuotedTokenSendFee  *common.TokenSendFee `json:"quotedTokenSendFee" hex:"optional"`
	QuotedExecutionMode common.ExecutionMode `json:"quotedExecutionMode"`
	QuotedExecutorFee   *common.ExecutorFee  `json:"quotedExecutorFee" hex:"optional"`
	CcipReceiveGasLimit types.INT64          `json:"ccipReceiveGasLimit"`
	ExecutorArgs        types.TEXT           `json:"executorArgs"`
	PoolFeeBps          types.NUMERIC        `json:"poolFeeBps"`
}

// ToMap converts GatheredQuote to a map for DAML arguments
func (t GatheredQuote) ToMap() map[string]any {
	m := make(map[string]any)

	m["quotedCCVFees"] = func() []any {
		res := make([]any, 0, len(t.QuotedCCVFees))
		for _, e := range t.QuotedCCVFees {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.QuotedTokenSendFee != nil {
		m["quotedTokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": *t.QuotedTokenSendFee,
		}
	} else {
		m["quotedTokenSendFee"] = map[string]any{
			"_type": "optional",
		}
	}

	m["quotedExecutionMode"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.QuotedExecutionMode).(mapper); ok {
			return m.toMap()
		}
		return t.QuotedExecutionMode
	}()

	if t.QuotedExecutorFee != nil {
		m["quotedExecutorFee"] = map[string]any{
			"_type": "optional",
			"value": *t.QuotedExecutorFee,
		}
	} else {
		m["quotedExecutorFee"] = map[string]any{
			"_type": "optional",
		}
	}

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["poolFeeBps"] = t.PoolFeeBps

	return m
}

func (t GatheredQuote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GatheredQuote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GatheredQuote to hex string (Canton MCMS format)
func (t GatheredQuote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GatheredQuote from hex string (Canton MCMS format)
func (t *GatheredQuote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee2 is a Record type
type GetFee2 struct {
	DestinationChainSelector types.NUMERIC            `json:"destinationChainSelector"`
	Message                  client.Canton2AnyMessage `json:"message"`
	Ccvs                     []client.CCVSendInput    `json:"ccvs"`
	Context                  common.CCIPContext       `json:"context"`
	RouterCid                types.CONTRACT_ID        `json:"routerCid"`
}

// ToMap converts GetFee2 to a map for DAML arguments
func (t GetFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["ccvs"] = func() []any {
		res := make([]any, 0, len(t.Ccvs))
		for _, e := range t.Ccvs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["routerCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	return m
}

func (t GetFee2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee2 to hex string (Canton MCMS format)
func (t GetFee2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee2 from hex string (Canton MCMS format)
func (t *GetFee2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeResult is a Record type
type GetFeeResult struct {
	FeeTokenAmount     types.NUMERIC `json:"feeTokenAmount"`
	PoolFeeTokenAmount types.NUMERIC `json:"poolFeeTokenAmount"`
}

// ToMap converts GetFeeResult to a map for DAML arguments
func (t GetFeeResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenAmount"] = t.FeeTokenAmount

	m["poolFeeTokenAmount"] = t.PoolFeeTokenAmount

	return m
}

func (t GetFeeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeResult to hex string (Canton MCMS format)
func (t GetFeeResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeResult from hex string (Canton MCMS format)
func (t *GetFeeResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Send is a Record type
type Send struct {
	DestinationChainSelector types.NUMERIC            `json:"destinationChainSelector"`
	Message                  client.Canton2AnyMessage `json:"message"`
	Ccvs                     []client.CCVSendInput    `json:"ccvs"`
	Context                  common.CCIPContext       `json:"context"`
	RouterCid                types.CONTRACT_ID        `json:"routerCid"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["ccvs"] = func() []any {
		res := make([]any, 0, len(t.Ccvs))
		for _, e := range t.Ccvs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["routerCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	return m
}

func (t Send) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Send) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Send to hex string (Canton MCMS format)
func (t Send) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Send from hex string (Canton MCMS format)
func (t *Send) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	GetFee2(args GetFee2) (*bind.EncodedChoice, error)
	Send(args Send) (*bind.EncodedChoice, error)
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

// GetFee2 encodes parameters for the GetFee2 choice.
func (e *encoder) GetFee2(args GetFee2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee2", args)
}

// Send encodes parameters for the Send choice.
func (e *encoder) Send(args Send) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Send", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
