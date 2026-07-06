package events

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipapi"
	ccipcodec "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	clientapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/clientapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
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
	PackageName = "ccip-events"
	PackageID   = "7f9e43c98504032e5b703d80a9045c269e721feb26a8fd10adb4bf80d8928337"
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

// CCIPMessageSent is a Template type
type CCIPMessageSent struct {
	CcipOwner types.PARTY          `json:"ccipOwner"`
	CcvOwners []types.PARTY        `json:"ccvOwners"`
	Sender    types.PARTY          `json:"sender"`
	Observers []types.PARTY        `json:"observers"`
	Event     CCIPMessageSentEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPMessageSent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "CCIPMessageSent")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPMessageSent) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.EventsV1.Events", "CCIPMessageSent")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPMessageSent) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPMessageSent) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPMessageSent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPMessageSent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPMessageSent to hex string (Canton MCMS format)
func (t CCIPMessageSent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPMessageSent from hex string (Canton MCMS format)
func (t *CCIPMessageSent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPMessageSent

// Archive exercises the Archive choice on this CCIPMessageSent contract via the IICCIPMessageSent interface
// This method uses the package name in the template ID
func (t CCIPMessageSent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPMessageSent) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.EventsV1.Events", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Verify interface implementations for CCIPMessageSent

var _ clientapi.IICCIPMessageSent = (*CCIPMessageSent)(nil)

// CCIPMessageSentEvent is a Record type
type CCIPMessageSentEvent struct {
	DestChainSelector              types.NUMERIC                            `json:"destChainSelector"`
	SequenceNumber                 types.NUMERIC                            `json:"sequenceNumber"`
	MessageId                      types.TEXT                               `json:"messageId"`
	FeeToken                       splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	TokenAmountBeforeTokenPoolFees types.NUMERIC                            `json:"tokenAmountBeforeTokenPoolFees"`
	EncodedMessage                 types.TEXT                               `json:"encodedMessage"`
	VerifierBlobs                  []types.TEXT                             `json:"verifierBlobs"`
	Receipts                       []Receipt                                `json:"receipts"`
}

// ToMap converts CCIPMessageSentEvent to a map for DAML arguments
func (t CCIPMessageSentEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["messageId"] = string(t.MessageId)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["tokenAmountBeforeTokenPoolFees"] = t.TokenAmountBeforeTokenPoolFees

	m["encodedMessage"] = string(t.EncodedMessage)

	m["verifierBlobs"] = func() []any {
		res := make([]any, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["receipts"] = func() []any {
		res := make([]any, 0, len(t.Receipts))
		for _, e := range t.Receipts {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t CCIPMessageSentEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPMessageSentEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPMessageSentEvent to hex string (Canton MCMS format)
func (t CCIPMessageSentEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPMessageSentEvent from hex string (Canton MCMS format)
func (t *CCIPMessageSentEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutionStateChanged is a Template type
type ExecutionStateChanged struct {
	CcipOwner types.PARTY                `json:"ccipOwner"`
	CcvOwners []types.PARTY              `json:"ccvOwners"`
	Receiver  types.PARTY                `json:"receiver"`
	Event     ExecutionStateChangedEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutionStateChanged) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "ExecutionStateChanged")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutionStateChanged) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.EventsV1.Events", "ExecutionStateChanged")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutionStateChanged) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutionStateChanged) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutionStateChanged) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutionStateChanged) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutionStateChanged to hex string (Canton MCMS format)
func (t ExecutionStateChanged) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutionStateChanged from hex string (Canton MCMS format)
func (t *ExecutionStateChanged) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutionStateChanged

// Archive exercises the Archive choice on this ExecutionStateChanged contract via the IIExecutionStateChanged interface
// This method uses the package name in the template ID
func (t ExecutionStateChanged) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutionStateChanged) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.EventsV1.Events", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Verify interface implementations for ExecutionStateChanged

var _ clientapi.IIExecutionStateChanged = (*ExecutionStateChanged)(nil)

// ExecutionStateChangedEvent is a Record type
type ExecutionStateChangedEvent struct {
	SourceChainSelector types.NUMERIC                 `json:"sourceChainSelector"`
	SequenceNumber      types.NUMERIC                 `json:"sequenceNumber"`
	MessageId           types.TEXT                    `json:"messageId"`
	State               ccipapi.MessageExecutionState `json:"state"`
	ReturnData          types.TEXT                    `json:"returnData"`
}

// ToMap converts ExecutionStateChangedEvent to a map for DAML arguments
func (t ExecutionStateChangedEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["messageId"] = string(t.MessageId)

	m["state"] = model.NestedToDAMLValue(t.State)

	m["returnData"] = string(t.ReturnData)

	return m
}

func (t ExecutionStateChangedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutionStateChangedEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutionStateChangedEvent to hex string (Canton MCMS format)
func (t ExecutionStateChangedEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutionStateChangedEvent from hex string (Canton MCMS format)
func (t *ExecutionStateChangedEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IssuerType is an enum type
type IssuerType string

const (
	IssuerTypeIssuerType_CCV IssuerType = "IssuerType_CCV"

	IssuerTypeIssuerType_Pool IssuerType = "IssuerType_Pool"

	IssuerTypeIssuerType_Executor IssuerType = "IssuerType_Executor"

	IssuerTypeIssuerType_Network IssuerType = "IssuerType_Network"
)

func (e IssuerType) GetEnumConstructor() string { return string(e) }

func (e IssuerType) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Receipts", "IssuerType")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e IssuerType) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.EventsV1.Receipts", "IssuerType")
}

func (e IssuerType) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *IssuerType) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes IssuerType to hex string (Canton MCMS format)
func (e IssuerType) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes IssuerType from hex string (Canton MCMS format)
func (e *IssuerType) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = IssuerType("")

// Receipt is a Record type
type Receipt struct {
	IssuerType        IssuerType    `json:"issuerType"`
	IssuerAddress     types.TEXT    `json:"issuerAddress"`
	VersionTag        *types.TEXT   `json:"versionTag" hex:"optional"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeTokenAmount    types.NUMERIC `json:"feeTokenAmount"`
	ExtraArgs         types.TEXT    `json:"extraArgs"`
}

// ToMap converts Receipt to a map for DAML arguments
func (t Receipt) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuerType"] = model.NestedToDAMLValue(t.IssuerType)

	m["issuerAddress"] = string(t.IssuerAddress)

	if t.VersionTag != nil {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": string(*t.VersionTag),
		}
	} else {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeTokenAmount"] = t.FeeTokenAmount

	m["extraArgs"] = string(t.ExtraArgs)

	return m
}

func (t Receipt) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Receipt) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Receipt to hex string (Canton MCMS format)
func (t Receipt) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Receipt from hex string (Canton MCMS format)
func (t *Receipt) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketClaimed is a Template type
type TokenReceiveTicketClaimed struct {
	CcipOwner             types.PARTY                    `json:"ccipOwner"`
	CcvOwners             []types.PARTY                  `json:"ccvOwners"`
	PoolOwner             types.PARTY                    `json:"poolOwner"`
	Receiver              types.PARTY                    `json:"receiver"`
	TokenReceiver         types.PARTY                    `json:"tokenReceiver"`
	TokenReceiveTicketCid types.CONTRACT_ID              `json:"tokenReceiveTicketCid"`
	Event                 TokenReceiveTicketClaimedEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenReceiveTicketClaimed) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "TokenReceiveTicketClaimed")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.EventsV1.Events", "TokenReceiveTicketClaimed")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenReceiveTicketClaimed) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenReceiveTicketClaimed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimed to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimed) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimed from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimed) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenReceiveTicketClaimed

// Archive exercises the Archive choice on this TokenReceiveTicketClaimed contract
// This method uses the package name in the template ID
func (t TokenReceiveTicketClaimed) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.EventsV1.Events", "TokenReceiveTicketClaimed"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.EventsV1.Events", "TokenReceiveTicketClaimed"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TokenReceiveTicketClaimedEvent is a Record type
type TokenReceiveTicketClaimedEvent struct {
	VerifiedCCVs                 []chainlinkapi.RawInstanceAddress        `json:"verifiedCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	PoolAddress                  chainlinkapi.RawInstanceAddress          `json:"poolAddress"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	SourceAmount                 types.TEXT                               `json:"sourceAmount"`
	LocalAmount                  types.NUMERIC                            `json:"localAmount"`
	SourcePoolData               types.TEXT                               `json:"sourcePoolData"`
	MessageId                    types.TEXT                               `json:"messageId"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
	Finality                     ccipcodec.FinalityConfig                 `json:"finality"`
	Output                       TokenReceiveTicketClaimedOutput          `json:"output"`
}

// ToMap converts TokenReceiveTicketClaimedEvent to a map for DAML arguments
func (t TokenReceiveTicketClaimedEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	m["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["sourceAmount"] = string(t.SourceAmount)

	m["localAmount"] = t.LocalAmount

	m["sourcePoolData"] = string(t.SourcePoolData)

	m["messageId"] = string(t.MessageId)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["output"] = model.NestedToDAMLValue(t.Output)

	return m
}

func (t TokenReceiveTicketClaimedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedEvent to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedEvent from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketClaimedCompleted is a Record type
type TokenReceiveTicketClaimedCompleted struct {
	ReceiverHoldingCids []types.CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts TokenReceiveTicketClaimedCompleted to a map for DAML arguments
func (t TokenReceiveTicketClaimedCompleted) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverHoldingCids))
		for _, e := range t.ReceiverHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t TokenReceiveTicketClaimedCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedCompleted to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedCompleted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedCompleted from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedCompleted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketClaimedOutput is a variant/union type
type TokenReceiveTicketClaimedOutput struct {
	TokenReceiveTicketClaimedPending   *TokenReceiveTicketClaimedPending   `json:"TokenReceiveTicketClaimed_Pending,omitempty"`
	TokenReceiveTicketClaimedCompleted *TokenReceiveTicketClaimedCompleted `json:"TokenReceiveTicketClaimed_Completed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TokenReceiveTicketClaimedOutput
func (v TokenReceiveTicketClaimedOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for TokenReceiveTicketClaimedOutput
func (v *TokenReceiveTicketClaimedOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes TokenReceiveTicketClaimedOutput to hex string (Canton MCMS format)
func (v TokenReceiveTicketClaimedOutput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedOutput from hex string (Canton MCMS format)
func (v *TokenReceiveTicketClaimedOutput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TokenReceiveTicketClaimedOutput) GetVariantTag() string {

	if v.TokenReceiveTicketClaimedPending != nil {
		return "TokenReceiveTicketClaimed_Pending"
	}

	if v.TokenReceiveTicketClaimedCompleted != nil {
		return "TokenReceiveTicketClaimed_Completed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TokenReceiveTicketClaimedOutput) GetVariantValue() any {

	if v.TokenReceiveTicketClaimedPending != nil {
		return v.TokenReceiveTicketClaimedPending
	}

	if v.TokenReceiveTicketClaimedCompleted != nil {
		return v.TokenReceiveTicketClaimedCompleted
	}

	return nil
}

var _ types.VARIANT = (*TokenReceiveTicketClaimedOutput)(nil)

// TokenReceiveTicketClaimedPending is a Record type
type TokenReceiveTicketClaimedPending struct {
	TransferInstructionCid types.CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts TokenReceiveTicketClaimedPending to a map for DAML arguments
func (t TokenReceiveTicketClaimedPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferInstructionCid"] = model.NestedToDAMLValue(t.TransferInstructionCid)

	return m
}

func (t TokenReceiveTicketClaimedPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedPending to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedPending from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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
