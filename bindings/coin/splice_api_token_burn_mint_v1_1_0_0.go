package coin

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"https://github.com/smartcontractkit/go-daml/pkg/codec"
	"https://github.com/smartcontractkit/go-daml/pkg/model"
	. "https://github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

// IBurnMintFactory is a DAML interface
type IBurnMintFactory interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// BurnMintFactoryPublicFetch executes the BurnMintFactory_PublicFetch choice
	BurnMintFactoryPublicFetch(contractID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand

	// BurnMintFactoryBurnMint executes the BurnMintFactory_BurnMint choice
	BurnMintFactoryBurnMint(contractID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand
}

// BurnMintFactoryView is a Record type
type BurnMintFactoryView struct {
	Admin PARTY    `json:"admin"`
	Meta  Metadata `json:"meta"`
}

// ToMap converts BurnMintFactoryView to a map for DAML arguments
func (t BurnMintFactoryView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["admin"] = t.Admin.ToMap()

	m["meta"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.toMap()
		}
		return t.Meta
	}()

	return m
}

func (t BurnMintFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *BurnMintFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryBurnMint is a Record type
type BurnMintFactoryBurnMint struct {
	ExpectedAdmin    PARTY            `json:"expectedAdmin"`
	InstrumentId     InstrumentId     `json:"instrumentId"`
	InputHoldingCids []CONTRACT_ID    `json:"inputHoldingCids"`
	Outputs          []BurnMintOutput `json:"outputs"`
	ExtraActors      []PARTY          `json:"extraActors"`
	ExtraArgs        ExtraArgs        `json:"extraArgs"`
}

// ToMap converts BurnMintFactoryBurnMint to a map for DAML arguments
func (t BurnMintFactoryBurnMint) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["inputHoldingCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["outputs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Outputs))
		for _, e := range t.Outputs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["extraActors"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ExtraActors))
		for _, e := range t.ExtraActors {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["extraArgs"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t BurnMintFactoryBurnMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *BurnMintFactoryBurnMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryBurnMintResult is a Record type
type BurnMintFactoryBurnMintResult struct {
	OutputCids []CONTRACT_ID `json:"outputCids"`
}

// ToMap converts BurnMintFactoryBurnMintResult to a map for DAML arguments
func (t BurnMintFactoryBurnMintResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["outputCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OutputCids))
		for _, e := range t.OutputCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t BurnMintFactoryBurnMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *BurnMintFactoryBurnMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryPublicFetch is a Record type
type BurnMintFactoryPublicFetch struct {
	ExpectedAdmin PARTY `json:"expectedAdmin"`
	Actor         PARTY `json:"actor"`
}

// ToMap converts BurnMintFactoryPublicFetch to a map for DAML arguments
func (t BurnMintFactoryPublicFetch) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t BurnMintFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *BurnMintFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintOutput is a Record type
type BurnMintOutput struct {
	Owner   PARTY         `json:"owner"`
	Amount  NUMERIC       `json:"amount"`
	Context ChoiceContext `json:"context"`
}

// ToMap converts BurnMintOutput to a map for DAML arguments
func (t BurnMintOutput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["amount"] = t.Amount

	m["context"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	return m
}

func (t BurnMintOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *BurnMintOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IBurnMintFactoryInterfaceID returns the interface ID for the IBurnMintFactory interface using the package name
func IBurnMintFactoryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.BurnMintV1", "BurnMintFactory")
}

// IBurnMintFactoryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IBurnMintFactoryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "Splice.Api.Token.BurnMintV1", "BurnMintFactory")
}
