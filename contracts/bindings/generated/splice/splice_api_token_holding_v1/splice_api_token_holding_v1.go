package splice_api_token_holding_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "splice-api-token-holding-v1"
	PackageID   = "718a0f77e505a8de22f188bd4c87fe74101274e9d4cb1bfac7d09aec7158d35b"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IHolding is a DAML interface
type IHolding interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
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

// HoldingView is a Record type
type HoldingView struct {
	Owner        types.PARTY                           `json:"owner"`
	InstrumentId InstrumentId                          `json:"instrumentId"`
	Amount       types.NUMERIC                         `json:"amount"`
	Lock         *Lock                                 `json:"lock" hex:"optional"`
	Meta         splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts HoldingView to a map for DAML arguments
func (t HoldingView) ToMap() map[string]any {
	m := make(map[string]any)

	m["owner"] = t.Owner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = t.Amount

	if t.Lock != nil {
		m["lock"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Lock),
		}
	} else {
		m["lock"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t HoldingView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingView to hex string (Canton MCMS format)
func (t HoldingView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingView from hex string (Canton MCMS format)
func (t *HoldingView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentId is a Record type
type InstrumentId struct {
	Admin types.PARTY `json:"admin"`
	Id    types.TEXT  `json:"id"`
}

// ToMap converts InstrumentId to a map for DAML arguments
func (t InstrumentId) ToMap() map[string]any {
	m := make(map[string]any)

	m["admin"] = t.Admin.ToMap()

	m["id"] = string(t.Id)

	return m
}

func (t InstrumentId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentId to hex string (Canton MCMS format)
func (t InstrumentId) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentId from hex string (Canton MCMS format)
func (t *InstrumentId) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Lock is a Record type
type Lock struct {
	Holders      []types.PARTY    `json:"holders"`
	ExpiresAt    *types.TIMESTAMP `json:"expiresAt" hex:"optional"`
	ExpiresAfter *types.RELTIME   `json:"expiresAfter" hex:"optional"`
	Context      *types.TEXT      `json:"context" hex:"optional"`
}

// ToMap converts Lock to a map for DAML arguments
func (t Lock) ToMap() map[string]any {
	m := make(map[string]any)

	m["holders"] = func() []any {
		res := make([]any, 0, len(t.Holders))
		for _, e := range t.Holders {
			res = append(res, e.ToMap())
		}
		return res
	}()

	if t.ExpiresAt != nil {
		m["expiresAt"] = map[string]any{
			"_type": "optional",
			"value": *t.ExpiresAt,
		}
	} else {
		m["expiresAt"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExpiresAfter != nil {
		m["expiresAfter"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExpiresAfter),
		}
	} else {
		m["expiresAfter"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.Context != nil {
		m["context"] = map[string]any{
			"_type": "optional",
			"value": string(*t.Context),
		}
	} else {
		m["context"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t Lock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Lock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Lock to hex string (Canton MCMS format)
func (t Lock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Lock from hex string (Canton MCMS format)
func (t *Lock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IHoldingInterfaceID returns the interface ID for the IHolding interface using the package name
func IHoldingInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.HoldingV1", "Holding")
}

// IHoldingInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IHoldingInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.HoldingV1", "Holding")
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
