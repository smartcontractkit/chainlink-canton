package lockreleasetokenpool

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

// IHolding is a DAML interface
type IHolding interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// HoldingView is a Record type
type HoldingView struct {
	Owner        PARTY        `json:"owner"`
	InstrumentId InstrumentId `json:"instrumentId"`
	Amount       NUMERIC      `json:"amount"`
	Lock         *Lock        `json:"lock"`
	Meta         Metadata     `json:"meta"`
}

// ToMap converts HoldingView to a map for DAML arguments
func (t HoldingView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = t.Amount

	if t.Lock != nil {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.Lock,
		}
	} else {
		m["lock"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["meta"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.toMap()
		}
		return t.Meta
	}()

	return m
}

func (t HoldingView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *HoldingView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	Admin PARTY `json:"admin"`
	Id    TEXT  `json:"id"`
}

// ToMap converts InstrumentId to a map for DAML arguments
func (t InstrumentId) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["admin"] = t.Admin.ToMap()

	m["id"] = string(t.Id)

	return m
}

func (t InstrumentId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *InstrumentId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	Holders      []PARTY    `json:"holders"`
	ExpiresAt    *TIMESTAMP `json:"expiresAt"`
	ExpiresAfter RELTIME    `json:"expiresAfter"`
	Context      *TEXT      `json:"context"`
}

// ToMap converts Lock to a map for DAML arguments
func (t Lock) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["holders"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Holders))
		for _, e := range t.Holders {
			res = append(res, e.ToMap())
		}
		return res
	}()

	if t.ExpiresAt != nil {
		m["expiresAt"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.ExpiresAt,
		}
	} else {
		m["expiresAt"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["expiresAfter"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExpiresAfter).(mapper); ok {
			return m.toMap()
		}
		return t.ExpiresAfter
	}()

	if t.Context != nil {
		m["context"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.Context),
		}
	} else {
		m["context"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

func (t Lock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Lock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
