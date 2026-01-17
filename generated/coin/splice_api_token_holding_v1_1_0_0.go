package coin

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

// toMap converts HoldingView to a map for DAML arguments
func (t HoldingView) toMap() map[string]interface{} {
	return map[string]interface{}{

		"owner": t.Owner.ToMap(),
		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"amount": (*big.Int)(t.Amount),
		"lock":   *t.Lock,
		"meta": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Meta).(mapper); ok {
				return m.toMap()
			}
			return t.Meta
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for HoldingView using JsonCodec
func (t HoldingView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for HoldingView using JsonCodec
func (t *HoldingView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// InstrumentId is a Record type
type InstrumentId struct {
	Admin PARTY `json:"admin"`
	Id    TEXT  `json:"id"`
}

// toMap converts InstrumentId to a map for DAML arguments
func (t InstrumentId) toMap() map[string]interface{} {
	return map[string]interface{}{

		"admin": t.Admin.ToMap(),
		"id":    string(t.Id),
	}
}

// MarshalJSON implements custom JSON marshaling for InstrumentId using JsonCodec
func (t InstrumentId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for InstrumentId using JsonCodec
func (t *InstrumentId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Lock is a Record type
type Lock struct {
	Holders      []PARTY    `json:"holders"`
	ExpiresAt    *TIMESTAMP `json:"expiresAt"`
	ExpiresAfter RELTIME    `json:"expiresAfter"`
	Context      *TEXT      `json:"context"`
}

// toMap converts Lock to a map for DAML arguments
func (t Lock) toMap() map[string]interface{} {
	return map[string]interface{}{

		"holders": func() []interface{} {
			res := make([]interface{}, 0, len(t.Holders))
			for _, e := range t.Holders {
				res = append(res, e.ToMap())
			}
			return res
		}(),
		"expiresAt": *t.ExpiresAt,
		"expiresAfter": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.ExpiresAfter).(mapper); ok {
				return m.toMap()
			}
			return t.ExpiresAfter
		}(),
		"context": string(*t.Context),
	}
}

// MarshalJSON implements custom JSON marshaling for Lock using JsonCodec
func (t Lock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Lock using JsonCodec
func (t *Lock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IHoldingInterfaceID returns the interface ID for the IHolding interface
func IHoldingInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "Splice.Api.Token.HoldingV1", "Holding")
}
