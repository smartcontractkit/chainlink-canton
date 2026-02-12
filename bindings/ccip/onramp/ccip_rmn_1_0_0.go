package onramp

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

// Curse is a Record type
type Curse struct {
	Subject TEXT `json:"subject"`
}

// ToMap converts Curse to a map for DAML arguments
func (t Curse) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["subject"] = string(t.Subject)

	return m
}

func (t Curse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Curse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	ChainSelector NUMERIC `json:"chainSelector"`
}

// ToMap converts CurseChain to a map for DAML arguments
func (t CurseChain) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t CurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// CurseGlobal is a Record type
type CurseGlobal struct {
}

// ToMap converts CurseGlobal to a map for DAML arguments
func (t CurseGlobal) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	return m
}

func (t CurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	Subjects []TEXT `json:"subjects"`
}

// ToMap converts CurseMultiple to a map for DAML arguments
func (t CurseMultiple) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["subjects"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// GetCursedSubjects is a Record type
type GetCursedSubjects struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts GetCursedSubjects to a map for DAML arguments
func (t GetCursedSubjects) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetCursedSubjects) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetCursedSubjects) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// IsCursed is a Record type
type IsCursed struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts IsCursed to a map for DAML arguments
func (t IsCursed) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsCursed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *IsCursed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// IsCursedForChain is a Record type
type IsCursedForChain struct {
	Caller        PARTY   `json:"caller"`
	ChainSelector NUMERIC `json:"chainSelector"`
}

// ToMap converts IsCursedForChain to a map for DAML arguments
func (t IsCursedForChain) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t IsCursedForChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *IsCursedForChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// RMNRemote is a Template type
type RMNRemote struct {
	InstanceId     TEXT   `json:"instanceId"`
	RmnOwner       PARTY  `json:"rmnOwner"`
	CcipOwner      PARTY  `json:"ccipOwner"`
	CursedSubjects []TEXT `json:"cursedSubjects"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CursedSubjects))
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CursedSubjects))
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
	return jsonCodec.Marshall(t)
}

func (t *RMNRemote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// Archive exercises the Archive choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RMNRemote) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
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

// Uncurse is a Record type
type Uncurse struct {
	Subject TEXT `json:"subject"`
}

// ToMap converts Uncurse to a map for DAML arguments
func (t Uncurse) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["subject"] = string(t.Subject)

	return m
}

func (t Uncurse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Uncurse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	ChainSelector NUMERIC `json:"chainSelector"`
}

// ToMap converts UncurseChain to a map for DAML arguments
func (t UncurseChain) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t UncurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UncurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// UncurseGlobal is a Record type
type UncurseGlobal struct {
}

// ToMap converts UncurseGlobal to a map for DAML arguments
func (t UncurseGlobal) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	return m
}

func (t UncurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UncurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	Subjects []TEXT `json:"subjects"`
}

// ToMap converts UncurseMultiple to a map for DAML arguments
func (t UncurseMultiple) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["subjects"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UncurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UncurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// UpdateCCIPOwner is a Record type
type UpdateCCIPOwner struct {
	NewCCIPOwner PARTY `json:"newCCIPOwner"`
}

// ToMap converts UpdateCCIPOwner to a map for DAML arguments
func (t UpdateCCIPOwner) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newCCIPOwner"] = t.NewCCIPOwner.ToMap()

	return m
}

func (t UpdateCCIPOwner) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UpdateCCIPOwner) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
