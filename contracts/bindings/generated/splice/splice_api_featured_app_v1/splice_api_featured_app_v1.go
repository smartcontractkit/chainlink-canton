package splice_api_featured_app_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	PackageName = "splice-api-featured-app-v1"
	PackageID   = "7804375fe5e4c6d5afe067bd314c42fe0b7d005a1300019c73154dd939da4dda"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IFeaturedAppActivityMarker is a DAML interface
type IFeaturedAppActivityMarker interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// IFeaturedAppRight is a DAML interface
type IFeaturedAppRight interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// FeaturedAppRightCreateActivityMarker executes the FeaturedAppRight_CreateActivityMarker choice
	FeaturedAppRightCreateActivityMarker(contractID string, args FeaturedAppRightCreateActivityMarker) *model.ExerciseCommand
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

// AppRewardBeneficiary is a Record type
type AppRewardBeneficiary struct {
	Beneficiary types.PARTY   `json:"beneficiary"`
	Weight      types.NUMERIC `json:"weight"`
}

// ToMap converts AppRewardBeneficiary to a map for DAML arguments
func (t AppRewardBeneficiary) ToMap() map[string]any {
	m := make(map[string]any)

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = t.Weight

	return m
}

func (t AppRewardBeneficiary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AppRewardBeneficiary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AppRewardBeneficiary to hex string (Canton MCMS format)
func (t AppRewardBeneficiary) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AppRewardBeneficiary from hex string (Canton MCMS format)
func (t *AppRewardBeneficiary) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeaturedAppActivityMarkerView is a Record type
type FeaturedAppActivityMarkerView struct {
	Dso         types.PARTY   `json:"dso"`
	Provider    types.PARTY   `json:"provider"`
	Beneficiary types.PARTY   `json:"beneficiary"`
	Weight      types.NUMERIC `json:"weight"`
}

// ToMap converts FeaturedAppActivityMarkerView to a map for DAML arguments
func (t FeaturedAppActivityMarkerView) ToMap() map[string]any {
	m := make(map[string]any)

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = t.Weight

	return m
}

func (t FeaturedAppActivityMarkerView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeaturedAppActivityMarkerView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeaturedAppActivityMarkerView to hex string (Canton MCMS format)
func (t FeaturedAppActivityMarkerView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeaturedAppActivityMarkerView from hex string (Canton MCMS format)
func (t *FeaturedAppActivityMarkerView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeaturedAppRightView is a Record type
type FeaturedAppRightView struct {
	Dso      types.PARTY `json:"dso"`
	Provider types.PARTY `json:"provider"`
}

// ToMap converts FeaturedAppRightView to a map for DAML arguments
func (t FeaturedAppRightView) ToMap() map[string]any {
	m := make(map[string]any)

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	return m
}

func (t FeaturedAppRightView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeaturedAppRightView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeaturedAppRightView to hex string (Canton MCMS format)
func (t FeaturedAppRightView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeaturedAppRightView from hex string (Canton MCMS format)
func (t *FeaturedAppRightView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeaturedAppRightCreateActivityMarker is a Record type
type FeaturedAppRightCreateActivityMarker struct {
	Beneficiaries []AppRewardBeneficiary `json:"beneficiaries"`
}

// ToMap converts FeaturedAppRightCreateActivityMarker to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarker) ToMap() map[string]any {
	m := make(map[string]any)

	m["beneficiaries"] = func() []any {
		res := make([]any, 0, len(t.Beneficiaries))
		for _, e := range t.Beneficiaries {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t FeaturedAppRightCreateActivityMarker) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeaturedAppRightCreateActivityMarker) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeaturedAppRightCreateActivityMarker to hex string (Canton MCMS format)
func (t FeaturedAppRightCreateActivityMarker) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeaturedAppRightCreateActivityMarker from hex string (Canton MCMS format)
func (t *FeaturedAppRightCreateActivityMarker) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeaturedAppRightCreateActivityMarkerResult is a Record type
type FeaturedAppRightCreateActivityMarkerResult struct {
	ActivityMarkerCids []types.CONTRACT_ID `json:"activityMarkerCids"`
}

// ToMap converts FeaturedAppRightCreateActivityMarkerResult to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarkerResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["activityMarkerCids"] = func() []any {
		res := make([]any, 0, len(t.ActivityMarkerCids))
		for _, e := range t.ActivityMarkerCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t FeaturedAppRightCreateActivityMarkerResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeaturedAppRightCreateActivityMarkerResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeaturedAppRightCreateActivityMarkerResult to hex string (Canton MCMS format)
func (t FeaturedAppRightCreateActivityMarkerResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeaturedAppRightCreateActivityMarkerResult from hex string (Canton MCMS format)
func (t *FeaturedAppRightCreateActivityMarkerResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IFeaturedAppActivityMarkerInterfaceID returns the interface ID for the IFeaturedAppActivityMarker interface using the package name
func IFeaturedAppActivityMarkerInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.FeaturedAppRightV1", "FeaturedAppActivityMarker")
}

// IFeaturedAppActivityMarkerInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IFeaturedAppActivityMarkerInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.FeaturedAppRightV1", "FeaturedAppActivityMarker")
}

// IFeaturedAppRightInterfaceID returns the interface ID for the IFeaturedAppRight interface using the package name
func IFeaturedAppRightInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.FeaturedAppRightV1", "FeaturedAppRight")
}

// IFeaturedAppRightInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IFeaturedAppRightInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.FeaturedAppRightV1", "FeaturedAppRight")
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
