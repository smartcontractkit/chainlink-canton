package feetreasury

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	api "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
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
	PackageName = "ccip-fee-treasury"
	PackageID   = "1fa7a19b1d741940c6f6513d28a87111dbe99c5b3533dc8ce5ab9d50dc36b983"
	SDKVersion  = "3.4.11"
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

// AuthorizeFeeWithdrawal is a Record type
type AuthorizeFeeWithdrawal struct {
	Params AuthorizeFeeWithdrawalParams `json:"params"`
}

// ToMap converts AuthorizeFeeWithdrawal to a map for DAML arguments
func (t AuthorizeFeeWithdrawal) ToMap() map[string]any {
	m := make(map[string]any)

	m["params"] = model.NestedToDAMLValue(t.Params)

	return m
}

func (t AuthorizeFeeWithdrawal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AuthorizeFeeWithdrawal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AuthorizeFeeWithdrawal to hex string (Canton MCMS format)
func (t AuthorizeFeeWithdrawal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AuthorizeFeeWithdrawal from hex string (Canton MCMS format)
func (t *AuthorizeFeeWithdrawal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AuthorizeFeeWithdrawalParams is a Record type
type AuthorizeFeeWithdrawalParams struct {
	AuthorizationId types.TEXT                               `json:"authorizationId"`
	Recipient       types.PARTY                              `json:"recipient"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	MaxAmount       types.NUMERIC                            `json:"maxAmount" hex:"decimal"`
	ValiditySecs    types.INT64                              `json:"validitySecs"`
}

// ToMap converts AuthorizeFeeWithdrawalParams to a map for DAML arguments
func (t AuthorizeFeeWithdrawalParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["authorizationId"] = string(t.AuthorizationId)

	m["recipient"] = t.Recipient.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["maxAmount"] = t.MaxAmount

	m["validitySecs"] = int64(t.ValiditySecs)

	return m
}

func (t AuthorizeFeeWithdrawalParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AuthorizeFeeWithdrawalParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AuthorizeFeeWithdrawalParams to hex string (Canton MCMS format)
func (t AuthorizeFeeWithdrawalParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AuthorizeFeeWithdrawalParams from hex string (Canton MCMS format)
func (t *AuthorizeFeeWithdrawalParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CleanupExpiredAuthorization is a Record type
type CleanupExpiredAuthorization struct {
	Submitter types.PARTY `json:"submitter"`
}

// ToMap converts CleanupExpiredAuthorization to a map for DAML arguments
func (t CleanupExpiredAuthorization) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	return m
}

func (t CleanupExpiredAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CleanupExpiredAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CleanupExpiredAuthorization to hex string (Canton MCMS format)
func (t CleanupExpiredAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CleanupExpiredAuthorization from hex string (Canton MCMS format)
func (t *CleanupExpiredAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteFeeWithdrawal is a Record type
type ExecuteFeeWithdrawal struct {
	Submitter          types.PARTY                            `json:"submitter"`
	TransferFactoryCid types.CONTRACT_ID                      `json:"transferFactoryCid"`
	InputHoldingCids   []types.CONTRACT_ID                    `json:"inputHoldingCids"`
	Amount             types.NUMERIC                          `json:"amount"`
	ExtraArgs          splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
	RequestedAt        types.TIMESTAMP                        `json:"requestedAt"`
}

// ToMap converts ExecuteFeeWithdrawal to a map for DAML arguments
func (t ExecuteFeeWithdrawal) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["transferFactoryCid"] = model.NestedToDAMLValue(t.TransferFactoryCid)

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	m["requestedAt"] = t.RequestedAt

	return m
}

func (t ExecuteFeeWithdrawal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteFeeWithdrawal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteFeeWithdrawal to hex string (Canton MCMS format)
func (t ExecuteFeeWithdrawal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteFeeWithdrawal from hex string (Canton MCMS format)
func (t *ExecuteFeeWithdrawal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeWithdrawalAuthorization is a Template type
type FeeWithdrawalAuthorization struct {
	InstanceId     types.TEXT                               `json:"instanceId"`
	FeeOwner       types.PARTY                              `json:"feeOwner"`
	McmsController types.PARTY                              `json:"mcmsController"`
	Recipient      types.PARTY                              `json:"recipient"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	MaxAmount      types.NUMERIC                            `json:"maxAmount" hex:"decimal"`
	ExpiresAt      types.TIMESTAMP                          `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FeeWithdrawalAuthorization) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FeeWithdrawalAuthorization) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeOwner"] = t.FeeOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsController"] = t.McmsController.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["recipient"] = t.Recipient.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.MaxAmount != "" {
		args["maxAmount"] = t.MaxAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeOwner"] = t.FeeOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsController"] = t.McmsController.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["recipient"] = t.Recipient.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.MaxAmount != "" {
		args["maxAmount"] = t.MaxAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FeeWithdrawalAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeWithdrawalAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeWithdrawalAuthorization to hex string (Canton MCMS format)
func (t FeeWithdrawalAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeWithdrawalAuthorization from hex string (Canton MCMS format)
func (t *FeeWithdrawalAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FeeWithdrawalAuthorization

// ExecuteFeeWithdrawal exercises the ExecuteFeeWithdrawal choice on this FeeWithdrawalAuthorization contract
// This method uses the package name in the template ID
func (t FeeWithdrawalAuthorization) ExecuteFeeWithdrawal(contractID string, args ExecuteFeeWithdrawal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteFeeWithdrawal",
		Arguments:  argsToMap(args),
	}
}

// ExecuteFeeWithdrawalWithPackageID exercises the ExecuteFeeWithdrawal choice using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) ExecuteFeeWithdrawalWithPackageID(contractID string, packageID string, args ExecuteFeeWithdrawal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteFeeWithdrawal",
		Arguments:  argsToMap(args),
	}
}

// CleanupExpiredAuthorization exercises the CleanupExpiredAuthorization choice on this FeeWithdrawalAuthorization contract
// This method uses the package name in the template ID
func (t FeeWithdrawalAuthorization) CleanupExpiredAuthorization(contractID string, args CleanupExpiredAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "CleanupExpiredAuthorization",
		Arguments:  argsToMap(args),
	}
}

// CleanupExpiredAuthorizationWithPackageID exercises the CleanupExpiredAuthorization choice using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) CleanupExpiredAuthorizationWithPackageID(contractID string, packageID string, args CleanupExpiredAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "CleanupExpiredAuthorization",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FeeWithdrawalAuthorization contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeWithdrawalAuthorization) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// WithdrawAuthorization exercises the WithdrawAuthorization choice on this FeeWithdrawalAuthorization contract
// This method uses the package name in the template ID
func (t FeeWithdrawalAuthorization) WithdrawAuthorization(contractID string, args WithdrawAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "WithdrawAuthorization",
		Arguments:  argsToMap(args),
	}
}

// WithdrawAuthorizationWithPackageID exercises the WithdrawAuthorization choice using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) WithdrawAuthorizationWithPackageID(contractID string, packageID string, args WithdrawAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "FeeWithdrawalAuthorization"),
		ContractID: contractID,
		Choice:     "WithdrawAuthorization",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this FeeWithdrawalAuthorization contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeWithdrawalAuthorization) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t FeeWithdrawalAuthorization) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for FeeWithdrawalAuthorization

var _ api.IMCMSReceiver = (*FeeWithdrawalAuthorization)(nil)

// MCMSFeeTreasury is a Template type
type MCMSFeeTreasury struct {
	InstanceId     types.TEXT  `json:"instanceId"`
	FeeOwner       types.PARTY `json:"feeOwner"`
	McmsController types.PARTY `json:"mcmsController"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MCMSFeeTreasury) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSFeeTreasury")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MCMSFeeTreasury) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSFeeTreasury")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MCMSFeeTreasury) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeOwner"] = t.FeeOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsController"] = t.McmsController.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MCMSFeeTreasury) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeOwner"] = t.FeeOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsController"] = t.McmsController.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MCMSFeeTreasury) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSFeeTreasury) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSFeeTreasury to hex string (Canton MCMS format)
func (t MCMSFeeTreasury) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSFeeTreasury from hex string (Canton MCMS format)
func (t *MCMSFeeTreasury) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MCMSFeeTreasury

// AuthorizeFeeWithdrawal exercises the AuthorizeFeeWithdrawal choice on this MCMSFeeTreasury contract
// This method uses the package name in the template ID
func (t MCMSFeeTreasury) AuthorizeFeeWithdrawal(contractID string, args AuthorizeFeeWithdrawal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSFeeTreasury"),
		ContractID: contractID,
		Choice:     "AuthorizeFeeWithdrawal",
		Arguments:  argsToMap(args),
	}
}

// AuthorizeFeeWithdrawalWithPackageID exercises the AuthorizeFeeWithdrawal choice using the provided package ID instead of package name
func (t MCMSFeeTreasury) AuthorizeFeeWithdrawalWithPackageID(contractID string, packageID string, args AuthorizeFeeWithdrawal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSFeeTreasury"),
		ContractID: contractID,
		Choice:     "AuthorizeFeeWithdrawal",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMSFeeTreasury contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t MCMSFeeTreasury) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MCMSFeeTreasury) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this MCMSFeeTreasury contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t MCMSFeeTreasury) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t MCMSFeeTreasury) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeTreasury", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for MCMSFeeTreasury

var _ api.IMCMSReceiver = (*MCMSFeeTreasury)(nil)

// WithdrawAuthorization is a Record type
type WithdrawAuthorization struct {
}

// ToMap converts WithdrawAuthorization to a map for DAML arguments
func (t WithdrawAuthorization) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t WithdrawAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithdrawAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithdrawAuthorization to hex string (Canton MCMS format)
func (t WithdrawAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithdrawAuthorization from hex string (Canton MCMS format)
func (t *WithdrawAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AuthorizeFeeWithdrawal(args AuthorizeFeeWithdrawal) (*bind.EncodedChoice, error)
	AuthorizeFeeWithdrawalParams(args AuthorizeFeeWithdrawalParams) (*bind.EncodedChoice, error)
	CleanupExpiredAuthorization(args CleanupExpiredAuthorization) (*bind.EncodedChoice, error)
	ExecuteFeeWithdrawal(args ExecuteFeeWithdrawal) (*bind.EncodedChoice, error)
	WithdrawAuthorization(args WithdrawAuthorization) (*bind.EncodedChoice, error)
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

// AuthorizeFeeWithdrawal encodes parameters for the AuthorizeFeeWithdrawal choice.
func (e *encoder) AuthorizeFeeWithdrawal(args AuthorizeFeeWithdrawal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AuthorizeFeeWithdrawal", args)
}

// AuthorizeFeeWithdrawalParams encodes parameters for the AuthorizeFeeWithdrawal choice.
func (e *encoder) AuthorizeFeeWithdrawalParams(args AuthorizeFeeWithdrawalParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AuthorizeFeeWithdrawal", args)
}

// CleanupExpiredAuthorization encodes parameters for the CleanupExpiredAuthorization choice.
func (e *encoder) CleanupExpiredAuthorization(args CleanupExpiredAuthorization) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CleanupExpiredAuthorization", args)
}

// ExecuteFeeWithdrawal encodes parameters for the ExecuteFeeWithdrawal choice.
func (e *encoder) ExecuteFeeWithdrawal(args ExecuteFeeWithdrawal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteFeeWithdrawal", args)
}

// WithdrawAuthorization encodes parameters for the WithdrawAuthorization choice.
func (e *encoder) WithdrawAuthorization(args WithdrawAuthorization) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("WithdrawAuthorization", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
