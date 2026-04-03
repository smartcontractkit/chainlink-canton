package rmn

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
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
	PackageName = "ccip-rmn"
	PackageID   = "60d7876833854dc278dd529aa09ce8bd0cc9cbcf360a8dc59e1e9d957c712356"
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

// AddCustomObservers is a Record type
type AddCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObservers to a map for DAML arguments
func (t AddCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObservers to hex string (Canton MCMS format)
func (t AddCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObservers from hex string (Canton MCMS format)
func (t *AddCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCustomObserversParams is a Record type
type AddCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObserversParams to a map for DAML arguments
func (t AddCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObserversParams to hex string (Canton MCMS format)
func (t AddCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObserversParams from hex string (Canton MCMS format)
func (t *AddCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Curse is a Record type
type Curse struct {
	Subject types.TEXT `json:"subject"`
}

// ToMap converts Curse to a map for DAML arguments
func (t Curse) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t Curse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Curse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Curse to hex string (Canton MCMS format)
func (t Curse) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Curse from hex string (Canton MCMS format)
func (t *Curse) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseChain is a Record type
type CurseChain struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts CurseChain to a map for DAML arguments
func (t CurseChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t CurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseChain to hex string (Canton MCMS format)
func (t CurseChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseChain from hex string (Canton MCMS format)
func (t *CurseChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseChainParams is a Record type
type CurseChainParams struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts CurseChainParams to a map for DAML arguments
func (t CurseChainParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t CurseChainParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseChainParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseChainParams to hex string (Canton MCMS format)
func (t CurseChainParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseChainParams from hex string (Canton MCMS format)
func (t *CurseChainParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseGlobal is a Record type
type CurseGlobal struct {
}

// ToMap converts CurseGlobal to a map for DAML arguments
func (t CurseGlobal) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseGlobal to hex string (Canton MCMS format)
func (t CurseGlobal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseGlobal from hex string (Canton MCMS format)
func (t *CurseGlobal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseMultiple is a Record type
type CurseMultiple struct {
	Subjects []types.TEXT `json:"subjects"`
}

// ToMap converts CurseMultiple to a map for DAML arguments
func (t CurseMultiple) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseMultiple to hex string (Canton MCMS format)
func (t CurseMultiple) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseMultiple from hex string (Canton MCMS format)
func (t *CurseMultiple) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseMultipleParams is a Record type
type CurseMultipleParams struct {
	Subjects []types.TEXT `json:"subjects"`
}

// ToMap converts CurseMultipleParams to a map for DAML arguments
func (t CurseMultipleParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CurseMultipleParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseMultipleParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseMultipleParams to hex string (Canton MCMS format)
func (t CurseMultipleParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseMultipleParams from hex string (Canton MCMS format)
func (t *CurseMultipleParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseParams is a Record type
type CurseParams struct {
	Subject types.TEXT `json:"subject"`
}

// ToMap converts CurseParams to a map for DAML arguments
func (t CurseParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t CurseParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseParams to hex string (Canton MCMS format)
func (t CurseParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseParams from hex string (Canton MCMS format)
func (t *CurseParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Get is a Record type
type Get struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts Get to a map for DAML arguments
func (t Get) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Get) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Get) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Get to hex string (Canton MCMS format)
func (t Get) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get from hex string (Canton MCMS format)
func (t *Get) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMCMSParams is Get without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetMCMSParams struct {
}

// MarshalHex encodes GetMCMSParams to hex string for MCMS operationData.
func (t GetMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMCMSParams from hex string.
func (t *GetMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetCursedSubjects is a Record type
type GetCursedSubjects struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetCursedSubjects to a map for DAML arguments
func (t GetCursedSubjects) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetCursedSubjects) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetCursedSubjects) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetCursedSubjects to hex string (Canton MCMS format)
func (t GetCursedSubjects) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetCursedSubjects from hex string (Canton MCMS format)
func (t *GetCursedSubjects) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetCursedSubjectsMCMSParams is GetCursedSubjects without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetCursedSubjectsMCMSParams struct {
}

// MarshalHex encodes GetCursedSubjectsMCMSParams to hex string for MCMS operationData.
func (t GetCursedSubjectsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetCursedSubjectsMCMSParams from hex string.
func (t *GetCursedSubjectsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursed is a Record type
type IsCursed struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts IsCursed to a map for DAML arguments
func (t IsCursed) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsCursed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsCursed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsCursed to hex string (Canton MCMS format)
func (t IsCursed) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursed from hex string (Canton MCMS format)
func (t *IsCursed) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedMCMSParams is IsCursed without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsCursedMCMSParams struct {
}

// MarshalHex encodes IsCursedMCMSParams to hex string for MCMS operationData.
func (t IsCursedMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedMCMSParams from hex string.
func (t *IsCursedMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedForChain is a Record type
type IsCursedForChain struct {
	Caller        types.PARTY   `json:"caller"`
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts IsCursedForChain to a map for DAML arguments
func (t IsCursedForChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t IsCursedForChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsCursedForChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsCursedForChain to hex string (Canton MCMS format)
func (t IsCursedForChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedForChain from hex string (Canton MCMS format)
func (t *IsCursedForChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedForChainMCMSParams is IsCursedForChain without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsCursedForChainMCMSParams struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// MarshalHex encodes IsCursedForChainMCMSParams to hex string for MCMS operationData.
func (t IsCursedForChainMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedForChainMCMSParams from hex string.
func (t *IsCursedForChainMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RMNRemote is a Template type
type RMNRemote struct {
	InstanceId      types.TEXT    `json:"instanceId"`
	RmnOwner        types.PARTY   `json:"rmnOwner"`
	CcipOwner       types.PARTY   `json:"ccipOwner"`
	CustomObservers []types.PARTY `json:"customObservers"`
	CursedSubjects  []types.TEXT  `json:"cursedSubjects"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RMNRemote) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RMNRemote) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RMNRemote) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []any {
		res := make([]any, 0, len(t.CursedSubjects))
		for _, e := range t.CursedSubjects {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RMNRemote) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []any {
		res := make([]any, 0, len(t.CursedSubjects))
		for _, e := range t.CursedSubjects {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RMNRemote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemote to hex string (Canton MCMS format)
func (t RMNRemote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemote from hex string (Canton MCMS format)
func (t *RMNRemote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RMNRemote

// IsCursedForChain exercises the IsCursedForChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) IsCursedForChain(contractID string, args IsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// IsCursedForChainWithPackageID exercises the IsCursedForChain choice using the provided package ID instead of package name
func (t RMNRemote) IsCursedForChainWithPackageID(contractID string, packageID string, args IsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// UncurseChain exercises the UncurseChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseChain(contractID string, args UncurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseChain",
		Arguments:  argsToMap(args),
	}
}

// UncurseChainWithPackageID exercises the UncurseChain choice using the provided package ID instead of package name
func (t RMNRemote) UncurseChainWithPackageID(contractID string, packageID string, args UncurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseChain",
		Arguments:  argsToMap(args),
	}
}

// CurseChain exercises the CurseChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseChain(contractID string, args CurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseChain",
		Arguments:  argsToMap(args),
	}
}

// CurseChainWithPackageID exercises the CurseChain choice using the provided package ID instead of package name
func (t RMNRemote) CurseChainWithPackageID(contractID string, packageID string, args CurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseChain",
		Arguments:  argsToMap(args),
	}
}

// UncurseGlobal exercises the UncurseGlobal choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseGlobal(contractID string, args UncurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// UncurseGlobalWithPackageID exercises the UncurseGlobal choice using the provided package ID instead of package name
func (t RMNRemote) UncurseGlobalWithPackageID(contractID string, packageID string, args UncurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// CurseGlobal exercises the CurseGlobal choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseGlobal(contractID string, args CurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// CurseGlobalWithPackageID exercises the CurseGlobal choice using the provided package ID instead of package name
func (t RMNRemote) CurseGlobalWithPackageID(contractID string, packageID string, args CurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// Curse exercises the Curse choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Curse(contractID string, args Curse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Curse",
		Arguments:  argsToMap(args),
	}
}

// CurseWithPackageID exercises the Curse choice using the provided package ID instead of package name
func (t RMNRemote) CurseWithPackageID(contractID string, packageID string, args Curse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Curse",
		Arguments:  argsToMap(args),
	}
}

// Uncurse exercises the Uncurse choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Uncurse(contractID string, args Uncurse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Uncurse",
		Arguments:  argsToMap(args),
	}
}

// UncurseWithPackageID exercises the Uncurse choice using the provided package ID instead of package name
func (t RMNRemote) UncurseWithPackageID(contractID string, packageID string, args Uncurse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Uncurse",
		Arguments:  argsToMap(args),
	}
}

// CurseMultiple exercises the CurseMultiple choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseMultiple(contractID string, args CurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// CurseMultipleWithPackageID exercises the CurseMultiple choice using the provided package ID instead of package name
func (t RMNRemote) CurseMultipleWithPackageID(contractID string, packageID string, args CurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// UncurseMultiple exercises the UncurseMultiple choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseMultiple(contractID string, args UncurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// UncurseMultipleWithPackageID exercises the UncurseMultiple choice using the provided package ID instead of package name
func (t RMNRemote) UncurseMultipleWithPackageID(contractID string, packageID string, args UncurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// IsCursed exercises the IsCursed choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) IsCursed(contractID string, args IsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursed",
		Arguments:  argsToMap(args),
	}
}

// IsCursedWithPackageID exercises the IsCursed choice using the provided package ID instead of package name
func (t RMNRemote) IsCursedWithPackageID(contractID string, packageID string, args IsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursed",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RMNRemote contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RMNRemote) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RMNRemote) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Get(contractID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t RMNRemote) GetWithPackageID(contractID string, packageID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetCursedSubjects exercises the GetCursedSubjects choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) GetCursedSubjects(contractID string, args GetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// GetCursedSubjectsWithPackageID exercises the GetCursedSubjects choice using the provided package ID instead of package name
func (t RMNRemote) GetCursedSubjectsWithPackageID(contractID string, packageID string, args GetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObservers exercises the AddCustomObservers choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) AddCustomObservers(contractID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObserversWithPackageID exercises the AddCustomObservers choice using the provided package ID instead of package name
func (t RMNRemote) AddCustomObserversWithPackageID(contractID string, packageID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObservers exercises the RemoveCustomObservers choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) RemoveCustomObservers(contractID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObserversWithPackageID exercises the RemoveCustomObservers choice using the provided package ID instead of package name
func (t RMNRemote) RemoveCustomObserversWithPackageID(contractID string, packageID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// UpdateCCIPOwner exercises the UpdateCCIPOwner choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UpdateCCIPOwner(contractID string, args UpdateCCIPOwner) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UpdateCCIPOwner",
		Arguments:  argsToMap(args),
	}
}

// UpdateCCIPOwnerWithPackageID exercises the UpdateCCIPOwner choice using the provided package ID instead of package name
func (t RMNRemote) UpdateCCIPOwnerWithPackageID(contractID string, packageID string, args UpdateCCIPOwner) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UpdateCCIPOwner",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this RMNRemote contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RMNRemote) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t RMNRemote) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for RMNRemote

var _ mcms.IMCMSReceiver = (*RMNRemote)(nil)

// RemoveCustomObservers is a Record type
type RemoveCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObservers to a map for DAML arguments
func (t RemoveCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObservers to hex string (Canton MCMS format)
func (t RemoveCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObservers from hex string (Canton MCMS format)
func (t *RemoveCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveCustomObserversParams is a Record type
type RemoveCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObserversParams to a map for DAML arguments
func (t RemoveCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObserversParams to hex string (Canton MCMS format)
func (t RemoveCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObserversParams from hex string (Canton MCMS format)
func (t *RemoveCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Uncurse is a Record type
type Uncurse struct {
	Subject types.TEXT `json:"subject"`
}

// ToMap converts Uncurse to a map for DAML arguments
func (t Uncurse) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t Uncurse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Uncurse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Uncurse to hex string (Canton MCMS format)
func (t Uncurse) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Uncurse from hex string (Canton MCMS format)
func (t *Uncurse) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseChain is a Record type
type UncurseChain struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts UncurseChain to a map for DAML arguments
func (t UncurseChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t UncurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseChain to hex string (Canton MCMS format)
func (t UncurseChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseChain from hex string (Canton MCMS format)
func (t *UncurseChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseChainParams is a Record type
type UncurseChainParams struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts UncurseChainParams to a map for DAML arguments
func (t UncurseChainParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t UncurseChainParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseChainParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseChainParams to hex string (Canton MCMS format)
func (t UncurseChainParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseChainParams from hex string (Canton MCMS format)
func (t *UncurseChainParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseGlobal is a Record type
type UncurseGlobal struct {
}

// ToMap converts UncurseGlobal to a map for DAML arguments
func (t UncurseGlobal) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UncurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseGlobal to hex string (Canton MCMS format)
func (t UncurseGlobal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseGlobal from hex string (Canton MCMS format)
func (t *UncurseGlobal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseMultiple is a Record type
type UncurseMultiple struct {
	Subjects []types.TEXT `json:"subjects"`
}

// ToMap converts UncurseMultiple to a map for DAML arguments
func (t UncurseMultiple) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UncurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseMultiple to hex string (Canton MCMS format)
func (t UncurseMultiple) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseMultiple from hex string (Canton MCMS format)
func (t *UncurseMultiple) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseMultipleParams is a Record type
type UncurseMultipleParams struct {
	Subjects []types.TEXT `json:"subjects"`
}

// ToMap converts UncurseMultipleParams to a map for DAML arguments
func (t UncurseMultipleParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UncurseMultipleParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseMultipleParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseMultipleParams to hex string (Canton MCMS format)
func (t UncurseMultipleParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseMultipleParams from hex string (Canton MCMS format)
func (t *UncurseMultipleParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseParams is a Record type
type UncurseParams struct {
	Subject types.TEXT `json:"subject"`
}

// ToMap converts UncurseParams to a map for DAML arguments
func (t UncurseParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t UncurseParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseParams to hex string (Canton MCMS format)
func (t UncurseParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseParams from hex string (Canton MCMS format)
func (t *UncurseParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateCCIPOwner is a Record type
type UpdateCCIPOwner struct {
	NewCCIPOwner types.PARTY `json:"newCCIPOwner"`
}

// ToMap converts UpdateCCIPOwner to a map for DAML arguments
func (t UpdateCCIPOwner) ToMap() map[string]any {
	m := make(map[string]any)

	m["newCCIPOwner"] = t.NewCCIPOwner.ToMap()

	return m
}

func (t UpdateCCIPOwner) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateCCIPOwner) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateCCIPOwner to hex string (Canton MCMS format)
func (t UpdateCCIPOwner) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateCCIPOwner from hex string (Canton MCMS format)
func (t *UpdateCCIPOwner) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error)
	AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error)
	Curse(args Curse) (*bind.EncodedChoice, error)
	CurseChain(args CurseChain) (*bind.EncodedChoice, error)
	CurseChainParams(args CurseChainParams) (*bind.EncodedChoice, error)
	CurseGlobal(args CurseGlobal) (*bind.EncodedChoice, error)
	CurseMultiple(args CurseMultiple) (*bind.EncodedChoice, error)
	CurseMultipleParams(args CurseMultipleParams) (*bind.EncodedChoice, error)
	CurseParams(args CurseParams) (*bind.EncodedChoice, error)
	Get(args Get) (*bind.EncodedChoice, error)
	GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error)
	GetCursedSubjects(args GetCursedSubjects) (*bind.EncodedChoice, error)
	GetCursedSubjectsMCMSParams(args GetCursedSubjectsMCMSParams) (*bind.EncodedChoice, error)
	IsCursed(args IsCursed) (*bind.EncodedChoice, error)
	IsCursedMCMSParams(args IsCursedMCMSParams) (*bind.EncodedChoice, error)
	IsCursedForChain(args IsCursedForChain) (*bind.EncodedChoice, error)
	IsCursedForChainMCMSParams(args IsCursedForChainMCMSParams) (*bind.EncodedChoice, error)
	RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error)
	RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error)
	Uncurse(args Uncurse) (*bind.EncodedChoice, error)
	UncurseChain(args UncurseChain) (*bind.EncodedChoice, error)
	UncurseChainParams(args UncurseChainParams) (*bind.EncodedChoice, error)
	UncurseGlobal(args UncurseGlobal) (*bind.EncodedChoice, error)
	UncurseMultiple(args UncurseMultiple) (*bind.EncodedChoice, error)
	UncurseMultipleParams(args UncurseMultipleParams) (*bind.EncodedChoice, error)
	UncurseParams(args UncurseParams) (*bind.EncodedChoice, error)
	UpdateCCIPOwner(args UpdateCCIPOwner) (*bind.EncodedChoice, error)
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

// AddCustomObservers encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// AddCustomObserversParams encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// Curse encodes parameters for the Curse choice.
func (e *encoder) Curse(args Curse) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Curse", args)
}

// CurseChain encodes parameters for the CurseChain choice.
func (e *encoder) CurseChain(args CurseChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseChain", args)
}

// CurseChainParams encodes parameters for the CurseChain choice.
func (e *encoder) CurseChainParams(args CurseChainParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseChain", args)
}

// CurseGlobal encodes parameters for the CurseGlobal choice.
func (e *encoder) CurseGlobal(args CurseGlobal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseGlobal", args)
}

// CurseMultiple encodes parameters for the CurseMultiple choice.
func (e *encoder) CurseMultiple(args CurseMultiple) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseMultiple", args)
}

// CurseMultipleParams encodes parameters for the CurseMultiple choice.
func (e *encoder) CurseMultipleParams(args CurseMultipleParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseMultiple", args)
}

// CurseParams encodes parameters for the Curse choice.
func (e *encoder) CurseParams(args CurseParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Curse", args)
}

// Get encodes parameters for the Get choice.
func (e *encoder) Get(args Get) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetMCMSParams encodes MCMS parameters (without Caller) for the Get choice.
func (e *encoder) GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetCursedSubjects encodes parameters for the GetCursedSubjects choice.
func (e *encoder) GetCursedSubjects(args GetCursedSubjects) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetCursedSubjects", args)
}

// GetCursedSubjectsMCMSParams encodes MCMS parameters (without Caller) for the GetCursedSubjects choice.
func (e *encoder) GetCursedSubjectsMCMSParams(args GetCursedSubjectsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetCursedSubjects", args)
}

// IsCursed encodes parameters for the IsCursed choice.
func (e *encoder) IsCursed(args IsCursed) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursed", args)
}

// IsCursedMCMSParams encodes MCMS parameters (without Caller) for the IsCursed choice.
func (e *encoder) IsCursedMCMSParams(args IsCursedMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursed", args)
}

// IsCursedForChain encodes parameters for the IsCursedForChain choice.
func (e *encoder) IsCursedForChain(args IsCursedForChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursedForChain", args)
}

// IsCursedForChainMCMSParams encodes MCMS parameters (without Caller) for the IsCursedForChain choice.
func (e *encoder) IsCursedForChainMCMSParams(args IsCursedForChainMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursedForChain", args)
}

// RemoveCustomObservers encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// RemoveCustomObserversParams encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// Uncurse encodes parameters for the Uncurse choice.
func (e *encoder) Uncurse(args Uncurse) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Uncurse", args)
}

// UncurseChain encodes parameters for the UncurseChain choice.
func (e *encoder) UncurseChain(args UncurseChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseChain", args)
}

// UncurseChainParams encodes parameters for the UncurseChain choice.
func (e *encoder) UncurseChainParams(args UncurseChainParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseChain", args)
}

// UncurseGlobal encodes parameters for the UncurseGlobal choice.
func (e *encoder) UncurseGlobal(args UncurseGlobal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseGlobal", args)
}

// UncurseMultiple encodes parameters for the UncurseMultiple choice.
func (e *encoder) UncurseMultiple(args UncurseMultiple) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseMultiple", args)
}

// UncurseMultipleParams encodes parameters for the UncurseMultiple choice.
func (e *encoder) UncurseMultipleParams(args UncurseMultipleParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseMultiple", args)
}

// UncurseParams encodes parameters for the Uncurse choice.
func (e *encoder) UncurseParams(args UncurseParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Uncurse", args)
}

// UpdateCCIPOwner encodes parameters for the UpdateCCIPOwner choice.
func (e *encoder) UpdateCCIPOwner(args UpdateCCIPOwner) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateCCIPOwner", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
