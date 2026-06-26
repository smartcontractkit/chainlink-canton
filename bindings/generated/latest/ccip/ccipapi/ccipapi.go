package ccipapi

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
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
	PackageName = "ccip-api"
	PackageID   = "c352ceb205c1317b0932a5640e6b4a40d234abf40d8dff68cae45fc576440e05"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IEventEmitter is a DAML interface
type IEventEmitter interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// EventEmitterEmitCCIPMessageSentEvent executes the EventEmitter_EmitCCIPMessageSentEvent choice
	EventEmitterEmitCCIPMessageSentEvent(contractID string, args EventEmitterEmitCCIPMessageSentEvent) *model.ExerciseCommand
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

// EventEmitterView is a Record type
type EventEmitterView struct {
	CcipOwner  types.PARTY `json:"ccipOwner"`
	InstanceId types.TEXT  `json:"instanceId"`
}

// ToMap converts EventEmitterView to a map for DAML arguments
func (t EventEmitterView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t EventEmitterView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EventEmitterView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EventEmitterView to hex string (Canton MCMS format)
func (t EventEmitterView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EventEmitterView from hex string (Canton MCMS format)
func (t *EventEmitterView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EventEmitterEmitCCIPMessageSentEvent is a Record type
type EventEmitterEmitCCIPMessageSentEvent struct {
	DestChainSelector              types.NUMERIC                              `json:"destChainSelector"`
	SequenceNumber                 types.NUMERIC                              `json:"sequenceNumber"`
	Sender                         types.PARTY                                `json:"sender"`
	EncodedMessage                 types.TEXT                                 `json:"encodedMessage"`
	FeeToken                       splice_api_token_holding_v1.InstrumentId   `json:"feeToken"`
	TokenAmountBeforeTokenPoolFees types.NUMERIC                              `json:"tokenAmountBeforeTokenPoolFees"`
	Receipts                       []Receipt                                  `json:"receipts"`
	CcvOwners                      []types.PARTY                              `json:"ccvOwners"`
	VerifierBlobs                  []types.TEXT                               `json:"verifierBlobs"`
	Observers                      []types.PARTY                              `json:"observers"`
	ExtraActors                    []types.PARTY                              `json:"extraActors"`
	Context                        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts EventEmitterEmitCCIPMessageSentEvent to a map for DAML arguments
func (t EventEmitterEmitCCIPMessageSentEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["sender"] = t.Sender.ToMap()

	m["encodedMessage"] = string(t.EncodedMessage)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["tokenAmountBeforeTokenPoolFees"] = t.TokenAmountBeforeTokenPoolFees

	m["receipts"] = func() []any {
		res := make([]any, 0, len(t.Receipts))
		for _, e := range t.Receipts {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["verifierBlobs"] = func() []any {
		res := make([]any, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["extraActors"] = func() []any {
		res := make([]any, 0, len(t.ExtraActors))
		for _, e := range t.ExtraActors {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t EventEmitterEmitCCIPMessageSentEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EventEmitterEmitCCIPMessageSentEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EventEmitterEmitCCIPMessageSentEvent to hex string (Canton MCMS format)
func (t EventEmitterEmitCCIPMessageSentEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EventEmitterEmitCCIPMessageSentEvent from hex string (Canton MCMS format)
func (t *EventEmitterEmitCCIPMessageSentEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IssuerType is an enum type
type IssuerType string

const (
	IssuerTypeIssuerType_CCV IssuerType = "IssuerType_CCV"

	IssuerTypeIssuerType_Pool IssuerType = "IssuerType_Pool"

	IssuerTypeIssuerType_Executor IssuerType = "IssuerType_Executor"

	IssuerTypeIssuerType_Network IssuerType = "IssuerType_Network"
)

func (e IssuerType) GetEnumConstructor() string { return string(e) }

func (e IssuerType) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.EventEmitterV1", "IssuerType")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e IssuerType) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.API.EventEmitterV1", "IssuerType")
}

func (e IssuerType) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *IssuerType) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes IssuerType to hex string (Canton MCMS format)
func (e IssuerType) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes IssuerType from hex string (Canton MCMS format)
func (e *IssuerType) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = IssuerType("")

// Receipt is a Record type
type Receipt struct {
	IssuerType        IssuerType    `json:"issuerType"`
	IssuerAddress     types.TEXT    `json:"issuerAddress"`
	VersionTag        *types.TEXT   `json:"versionTag" hex:"optional"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeTokenAmount    types.NUMERIC `json:"feeTokenAmount"`
	ExtraArgs         types.TEXT    `json:"extraArgs"`
}

// ToMap converts Receipt to a map for DAML arguments
func (t Receipt) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuerType"] = model.NestedToDAMLValue(t.IssuerType)

	m["issuerAddress"] = string(t.IssuerAddress)

	if t.VersionTag != nil {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": string(*t.VersionTag),
		}
	} else {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeTokenAmount"] = t.FeeTokenAmount

	m["extraArgs"] = string(t.ExtraArgs)

	return m
}

func (t Receipt) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Receipt) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Receipt to hex string (Canton MCMS format)
func (t Receipt) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Receipt from hex string (Canton MCMS format)
func (t *Receipt) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IEventEmitterInterfaceID returns the interface ID for the IEventEmitter interface using the package name
func IEventEmitterInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.EventEmitterV1", "EventEmitter")
}

// IEventEmitterInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IEventEmitterInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.EventEmitterV1", "EventEmitter")
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
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

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
