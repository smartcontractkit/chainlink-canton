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

// toMap converts BurnMintFactoryView to a map for DAML arguments
func (t BurnMintFactoryView) toMap() map[string]interface{} {
	return map[string]interface{}{

		"admin": t.Admin.ToMap(),
		"meta":  t.Meta,
	}
}

// MarshalJSON implements custom JSON marshaling for BurnMintFactoryView using JsonCodec
func (t BurnMintFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BurnMintFactoryView using JsonCodec
func (t *BurnMintFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryBurnMint is a Record type
type BurnMintFactoryBurnMint struct {
	ExpectedAdmin    PARTY            `json:"expectedAdmin"`
	InstrumentId     InstrumentId     `json:"instrumentId"`
	InputHoldingCids []string         `json:"inputHoldingCids"`
	Outputs          []BurnMintOutput `json:"outputs"`
	ExtraActors      []PARTY          `json:"extraActors"`
	ExtraArgs        ExtraArgs        `json:"extraArgs"`
}

// toMap converts BurnMintFactoryBurnMint to a map for DAML arguments
func (t BurnMintFactoryBurnMint) toMap() map[string]interface{} {
	return map[string]interface{}{

		"expectedAdmin":    t.ExpectedAdmin.ToMap(),
		"instrumentId":     t.InstrumentId,
		"inputHoldingCids": t.InputHoldingCids,
		"outputs":          t.Outputs,
		"extraActors":      t.ExtraActors,
		"extraArgs":        t.ExtraArgs,
	}
}

// MarshalJSON implements custom JSON marshaling for BurnMintFactoryBurnMint using JsonCodec
func (t BurnMintFactoryBurnMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BurnMintFactoryBurnMint using JsonCodec
func (t *BurnMintFactoryBurnMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryBurnMintResult is a Record type
type BurnMintFactoryBurnMintResult struct {
	OutputCids []string `json:"outputCids"`
}

// toMap converts BurnMintFactoryBurnMintResult to a map for DAML arguments
func (t BurnMintFactoryBurnMintResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"outputCids": t.OutputCids,
	}
}

// MarshalJSON implements custom JSON marshaling for BurnMintFactoryBurnMintResult using JsonCodec
func (t BurnMintFactoryBurnMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BurnMintFactoryBurnMintResult using JsonCodec
func (t *BurnMintFactoryBurnMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// BurnMintFactoryPublicFetch is a Record type
type BurnMintFactoryPublicFetch struct {
	ExpectedAdmin PARTY `json:"expectedAdmin"`
	Actor         PARTY `json:"actor"`
}

// toMap converts BurnMintFactoryPublicFetch to a map for DAML arguments
func (t BurnMintFactoryPublicFetch) toMap() map[string]interface{} {
	return map[string]interface{}{

		"expectedAdmin": t.ExpectedAdmin.ToMap(),
		"actor":         t.Actor.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for BurnMintFactoryPublicFetch using JsonCodec
func (t BurnMintFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BurnMintFactoryPublicFetch using JsonCodec
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

// toMap converts BurnMintOutput to a map for DAML arguments
func (t BurnMintOutput) toMap() map[string]interface{} {
	return map[string]interface{}{

		"owner":   t.Owner.ToMap(),
		"amount":  (*big.Int)(t.Amount),
		"context": t.Context,
	}
}

// MarshalJSON implements custom JSON marshaling for BurnMintOutput using JsonCodec
func (t BurnMintOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for BurnMintOutput using JsonCodec
func (t *BurnMintOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IBurnMintFactoryInterfaceID returns the interface ID for the IBurnMintFactory interface
func IBurnMintFactoryInterfaceID(packageID *string) string {
	pkgName := packageName
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "BurnMintV1", "BurnMintFactory")
}
