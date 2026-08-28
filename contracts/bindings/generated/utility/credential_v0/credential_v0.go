package credential_v0

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
	PackageName = "utility-credential-v0"
	PackageID   = "5a29ead611a0abd5f5b3fc3caf7d0f67c0ff802032ab6d392824aa9060e56d70"
	SDKVersion  = "3.4.9"
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

// Claim is a Record type
type Claim struct {
	Subject  types.TEXT `json:"subject" hex:"bytes"`
	Property types.TEXT `json:"property"`
	Value    types.TEXT `json:"value"`
}

// ToMap converts Claim to a map for DAML arguments
func (t Claim) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	m["property"] = string(t.Property)

	m["value"] = string(t.Value)

	return m
}

func (t Claim) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Claim) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Claim to hex string (Canton MCMS format)
func (t Claim) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Claim from hex string (Canton MCMS format)
func (t *Claim) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Credential is a Template type
type Credential struct {
	Issuer      types.PARTY      `json:"issuer"`
	Holder      types.PARTY      `json:"holder"`
	Id          types.TEXT       `json:"id"`
	Description types.TEXT       `json:"description"`
	ValidFrom   *types.TIMESTAMP `json:"validFrom" hex:"optional"`
	ValidUntil  *types.TIMESTAMP `json:"validUntil" hex:"optional"`
	Claims      []Claim          `json:"claims"`
	Observers   types.SET        `json:"observers"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t Credential) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Credential.V0.Credential", "Credential")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t Credential) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Credential.V0.Credential", "Credential")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t Credential) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["id"] = string(t.Id)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["description"] = string(t.Description)

	if t.ValidFrom != nil {
		args["validFrom"] = map[string]any{
			"_type": "optional",
			"value": *t.ValidFrom,
		}
	} else {
		args["validFrom"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ValidUntil != nil {
		args["validUntil"] = map[string]any{
			"_type": "optional",
			"value": *t.ValidUntil,
		}
	} else {
		args["validUntil"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["claims"] = func() []any {
		res := make([]any, 0, len(t.Claims))
		for _, e := range t.Claims {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = model.NestedToDAMLValue(t.Observers)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t Credential) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["id"] = string(t.Id)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["description"] = string(t.Description)

	if t.ValidFrom != nil {
		args["validFrom"] = map[string]any{
			"_type": "optional",
			"value": *t.ValidFrom,
		}
	} else {
		args["validFrom"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ValidUntil != nil {
		args["validUntil"] = map[string]any{
			"_type": "optional",
			"value": *t.ValidUntil,
		}
	} else {
		args["validUntil"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["claims"] = func() []any {
		res := make([]any, 0, len(t.Claims))
		for _, e := range t.Claims {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = model.NestedToDAMLValue(t.Observers)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t Credential) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Credential) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Credential to hex string (Canton MCMS format)
func (t Credential) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Credential from hex string (Canton MCMS format)
func (t *Credential) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for Credential

// CredentialRevoke exercises the Credential_Revoke choice on this Credential contract
// This method uses the package name in the template ID
func (t Credential) CredentialRevoke(contractID string, args CredentialRevoke) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Credential_Revoke",
		Arguments:  argsToMap(args),
	}
}

// CredentialRevokeWithPackageID exercises the Credential_Revoke choice using the provided package ID instead of package name
func (t Credential) CredentialRevokeWithPackageID(contractID string, packageID string, args CredentialRevoke) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Credential_Revoke",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this Credential contract
// This method uses the package name in the template ID
func (t Credential) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t Credential) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CredentialGet exercises the Credential_Get choice on this Credential contract
// This method uses the package name in the template ID
func (t Credential) CredentialGet(contractID string, args CredentialGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Credential_Get",
		Arguments:  argsToMap(args),
	}
}

// CredentialGetWithPackageID exercises the Credential_Get choice using the provided package ID instead of package name
func (t Credential) CredentialGetWithPackageID(contractID string, packageID string, args CredentialGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Credential.V0.Credential", "Credential"),
		ContractID: contractID,
		Choice:     "Credential_Get",
		Arguments:  argsToMap(args),
	}
}

// CredentialGet is a Record type
type CredentialGet struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts CredentialGet to a map for DAML arguments
func (t CredentialGet) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t CredentialGet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CredentialGet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CredentialGet to hex string (Canton MCMS format)
func (t CredentialGet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CredentialGet from hex string (Canton MCMS format)
func (t *CredentialGet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CredentialGetResult is a Record type
type CredentialGetResult struct {
	Credential Credential `json:"credential"`
}

// ToMap converts CredentialGetResult to a map for DAML arguments
func (t CredentialGetResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["credential"] = model.NestedToDAMLValue(t.Credential)

	return m
}

func (t CredentialGetResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CredentialGetResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CredentialGetResult to hex string (Canton MCMS format)
func (t CredentialGetResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CredentialGetResult from hex string (Canton MCMS format)
func (t *CredentialGetResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CredentialRevoke is a Record type
type CredentialRevoke struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts CredentialRevoke to a map for DAML arguments
func (t CredentialRevoke) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t CredentialRevoke) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CredentialRevoke) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CredentialRevoke to hex string (Canton MCMS format)
func (t CredentialRevoke) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CredentialRevoke from hex string (Canton MCMS format)
func (t *CredentialRevoke) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CredentialRevokeResult is a Record type
type CredentialRevokeResult struct {
}

// ToMap converts CredentialRevokeResult to a map for DAML arguments
func (t CredentialRevokeResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CredentialRevokeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CredentialRevokeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CredentialRevokeResult to hex string (Canton MCMS format)
func (t CredentialRevokeResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CredentialRevokeResult from hex string (Canton MCMS format)
func (t *CredentialRevokeResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PartyCredentialRequirement is a Record type
type PartyCredentialRequirement struct {
	Issuer         types.PARTY `json:"issuer"`
	RequiredClaims []Tuple2    `json:"requiredClaims"`
}

// ToMap converts PartyCredentialRequirement to a map for DAML arguments
func (t PartyCredentialRequirement) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuer"] = t.Issuer.ToMap()

	m["requiredClaims"] = func() []any {
		res := make([]any, 0, len(t.RequiredClaims))
		for _, e := range t.RequiredClaims {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t PartyCredentialRequirement) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PartyCredentialRequirement) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PartyCredentialRequirement to hex string (Canton MCMS format)
func (t PartyCredentialRequirement) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PartyCredentialRequirement from hex string (Canton MCMS format)
func (t *PartyCredentialRequirement) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithIssuer is a Record type
type WithIssuer struct {
	Issuer types.PARTY `json:"issuer"`
}

// ToMap converts WithIssuer to a map for DAML arguments
func (t WithIssuer) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuer"] = t.Issuer.ToMap()

	return m
}

func (t WithIssuer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithIssuer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithIssuer to hex string (Canton MCMS format)
func (t WithIssuer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithIssuer from hex string (Canton MCMS format)
func (t *WithIssuer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithIssuerHolder is a Record type
type WithIssuerHolder struct {
	Issuer types.PARTY `json:"issuer"`
	Holder types.PARTY `json:"holder"`
}

// ToMap converts WithIssuerHolder to a map for DAML arguments
func (t WithIssuerHolder) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuer"] = t.Issuer.ToMap()

	m["holder"] = t.Holder.ToMap()

	return m
}

func (t WithIssuerHolder) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithIssuerHolder) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithIssuerHolder to hex string (Canton MCMS format)
func (t WithIssuerHolder) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithIssuerHolder from hex string (Canton MCMS format)
func (t *WithIssuerHolder) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	CredentialGet(args CredentialGet) (*bind.EncodedChoice, error)
	CredentialRevoke(args CredentialRevoke) (*bind.EncodedChoice, error)
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

// CredentialGet encodes parameters for the Credential_Get choice.
func (e *encoder) CredentialGet(args CredentialGet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Credential_Get", args)
}

// CredentialRevoke encodes parameters for the Credential_Revoke choice.
func (e *encoder) CredentialRevoke(args CredentialRevoke) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Credential_Revoke", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
