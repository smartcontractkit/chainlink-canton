package splice_api_token_burn_mint_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-canton/bindings/codec"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	model "github.com/smartcontractkit/chainlink-canton/bindings/ledger"
	"github.com/smartcontractkit/chainlink-canton/bindings/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
)

const (
	PackageName = "splice-api-token-burn-mint-v1"
	PackageID   = "9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IBurnMintFactory is a DAML interface
type IBurnMintFactory interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// BurnMintFactoryPublicFetch executes the BurnMintFactory_PublicFetch choice
	BurnMintFactoryPublicFetch(contractID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand

	// BurnMintFactoryBurnMint executes the BurnMintFactory_BurnMint choice
	BurnMintFactoryBurnMint(contractID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand
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

// BurnMintFactoryView is a Record type
type BurnMintFactoryView struct {
	Admin types.PARTY                           `json:"admin"`
	Meta  splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts BurnMintFactoryView to a map for DAML arguments
func (t BurnMintFactoryView) ToMap() map[string]any {
	m := make(map[string]any)

	m["admin"] = t.Admin.ToMap()

	m["meta"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.toMap()
		}
		return t.Meta
	}()

	return m
}

func (t BurnMintFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// BurnMintFactoryBurnMint is a Record type
type BurnMintFactoryBurnMint struct {
	ExpectedAdmin    types.PARTY                              `json:"expectedAdmin"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	InputHoldingCids []types.CONTRACT_ID                      `json:"inputHoldingCids"`
	Outputs          []BurnMintOutput                         `json:"outputs"`
	ExtraActors      []types.PARTY                            `json:"extraActors"`
	ExtraArgs        splice_api_token_metadata_v1.ExtraArgs   `json:"extraArgs"`
}

// ToMap converts BurnMintFactoryBurnMint to a map for DAML arguments
func (t BurnMintFactoryBurnMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["outputs"] = func() []any {
		res := make([]any, 0, len(t.Outputs))
		for _, e := range t.Outputs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
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

	m["extraArgs"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t BurnMintFactoryBurnMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintFactoryBurnMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// BurnMintFactoryBurnMintResult is a Record type
type BurnMintFactoryBurnMintResult struct {
	OutputCids []types.CONTRACT_ID `json:"outputCids"`
}

// ToMap converts BurnMintFactoryBurnMintResult to a map for DAML arguments
func (t BurnMintFactoryBurnMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["outputCids"] = func() []any {
		res := make([]any, 0, len(t.OutputCids))
		for _, e := range t.OutputCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t BurnMintFactoryBurnMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintFactoryBurnMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// BurnMintFactoryPublicFetch is a Record type
type BurnMintFactoryPublicFetch struct {
	ExpectedAdmin types.PARTY `json:"expectedAdmin"`
	Actor         types.PARTY `json:"actor"`
}

// ToMap converts BurnMintFactoryPublicFetch to a map for DAML arguments
func (t BurnMintFactoryPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t BurnMintFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// BurnMintOutput is a Record type
type BurnMintOutput struct {
	Owner   types.PARTY                                `json:"owner"`
	Amount  types.NUMERIC                              `json:"amount"`
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts BurnMintOutput to a map for DAML arguments
func (t BurnMintOutput) ToMap() map[string]any {
	m := make(map[string]any)

	m["owner"] = t.Owner.ToMap()

	m["amount"] = t.Amount

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	return m
}

func (t BurnMintOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnMintOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// IBurnMintFactoryInterfaceID returns the interface ID for the IBurnMintFactory interface using the package name
func IBurnMintFactoryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.BurnMintV1", "BurnMintFactory")
}

// IBurnMintFactoryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IBurnMintFactoryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.BurnMintV1", "BurnMintFactory")
}
