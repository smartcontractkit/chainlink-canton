package committeeverifier

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	extensionapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	api "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
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
	PackageName = "ccip-committee-verifier"
	PackageID   = "57d3b9724bd6b2cafb2567f098ac6e62ba201fc72dc738f1278ad99103f7425e"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	VersionTagV200           = types.TEXT("e9a05a20")
	SignatureLength          = types.INT64(64)
	SignatureComponentLength = types.INT64(32)
)

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

// AcceptStorageLocationsAdmin is a Record type
type AcceptStorageLocationsAdmin struct {
}

// ToMap converts AcceptStorageLocationsAdmin to a map for DAML arguments
func (t AcceptStorageLocationsAdmin) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t AcceptStorageLocationsAdmin) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptStorageLocationsAdmin) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptStorageLocationsAdmin to hex string (Canton MCMS format)
func (t AcceptStorageLocationsAdmin) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptStorageLocationsAdmin from hex string (Canton MCMS format)
func (t *AcceptStorageLocationsAdmin) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllowListConfigArgs is a Record type
type AllowListConfigArgs struct {
	DestChainSelector         types.NUMERIC `json:"destChainSelector"`
	AllowListEnabled          types.BOOL    `json:"allowListEnabled"`
	AddedAllowListedSenders   []types.PARTY `json:"addedAllowListedSenders"`
	RemovedAllowListedSenders []types.PARTY `json:"removedAllowListedSenders"`
}

// ToMap converts AllowListConfigArgs to a map for DAML arguments
func (t AllowListConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["allowListEnabled"] = bool(t.AllowListEnabled)

	m["addedAllowListedSenders"] = func() []any {
		res := make([]any, 0, len(t.AddedAllowListedSenders))
		for _, e := range t.AddedAllowListedSenders {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["removedAllowListedSenders"] = func() []any {
		res := make([]any, 0, len(t.RemovedAllowListedSenders))
		for _, e := range t.RemovedAllowListedSenders {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AllowListConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllowListConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllowListConfigArgs to hex string (Canton MCMS format)
func (t AllowListConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllowListConfigArgs from hex string (Canton MCMS format)
func (t *AllowListConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyAllowListUpdates is a Record type
type ApplyAllowListUpdates struct {
	AllowListConfigArgsItems []AllowListConfigArgs `json:"allowListConfigArgsItems"`
	Caller                   types.PARTY           `json:"caller"`
}

// ToMap converts ApplyAllowListUpdates to a map for DAML arguments
func (t ApplyAllowListUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["allowListConfigArgsItems"] = func() []any {
		res := make([]any, 0, len(t.AllowListConfigArgsItems))
		for _, e := range t.AllowListConfigArgsItems {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ApplyAllowListUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyAllowListUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyAllowListUpdates to hex string (Canton MCMS format)
func (t ApplyAllowListUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyAllowListUpdates from hex string (Canton MCMS format)
func (t *ApplyAllowListUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyAllowListUpdatesMCMSParams is ApplyAllowListUpdates without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ApplyAllowListUpdatesMCMSParams struct {
	AllowListConfigArgsItems []AllowListConfigArgs `json:"allowListConfigArgsItems"`
}

// MarshalHex encodes ApplyAllowListUpdatesMCMSParams to hex string for MCMS operationData.
func (t ApplyAllowListUpdatesMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyAllowListUpdatesMCMSParams from hex string.
func (t *ApplyAllowListUpdatesMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyAllowListUpdatesParams is a Record type
type ApplyAllowListUpdatesParams struct {
	AllowListConfigArgs []AllowListConfigArgs `json:"allowListConfigArgs"`
	Caller              types.PARTY           `json:"caller"`
}

// ToMap converts ApplyAllowListUpdatesParams to a map for DAML arguments
func (t ApplyAllowListUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["allowListConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.AllowListConfigArgs))
		for _, e := range t.AllowListConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ApplyAllowListUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyAllowListUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyAllowListUpdatesParams to hex string (Canton MCMS format)
func (t ApplyAllowListUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyAllowListUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyAllowListUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyRemoteChainConfigUpdates is a Record type
type ApplyRemoteChainConfigUpdates struct {
	RemoteChainConfigArgs []RemoteChainConfigArgs `json:"remoteChainConfigArgs"`
}

// ToMap converts ApplyRemoteChainConfigUpdates to a map for DAML arguments
func (t ApplyRemoteChainConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RemoteChainConfigArgs))
		for _, e := range t.RemoteChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyRemoteChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyRemoteChainConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyRemoteChainConfigUpdates to hex string (Canton MCMS format)
func (t ApplyRemoteChainConfigUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyRemoteChainConfigUpdates from hex string (Canton MCMS format)
func (t *ApplyRemoteChainConfigUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyRemoteChainConfigUpdatesParams is a Record type
type ApplyRemoteChainConfigUpdatesParams struct {
	RemoteChainConfigArgs []RemoteChainConfigArgs `json:"remoteChainConfigArgs"`
}

// ToMap converts ApplyRemoteChainConfigUpdatesParams to a map for DAML arguments
func (t ApplyRemoteChainConfigUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.RemoteChainConfigArgs))
		for _, e := range t.RemoteChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyRemoteChainConfigUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyRemoteChainConfigUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyRemoteChainConfigUpdatesParams to hex string (Canton MCMS format)
func (t ApplyRemoteChainConfigUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyRemoteChainConfigUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyRemoteChainConfigUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplySignatureConfigs is a Record type
type ApplySignatureConfigs struct {
	SourceChainSelectorsToRemove []types.NUMERIC   `json:"sourceChainSelectorsToRemove"`
	SignatureConfigs             []SignatureConfig `json:"signatureConfigs"`
}

// ToMap converts ApplySignatureConfigs to a map for DAML arguments
func (t ApplySignatureConfigs) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.SourceChainSelectorsToRemove))
		for _, e := range t.SourceChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["signatureConfigs"] = func() []any {
		res := make([]any, 0, len(t.SignatureConfigs))
		for _, e := range t.SignatureConfigs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplySignatureConfigs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplySignatureConfigs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplySignatureConfigs to hex string (Canton MCMS format)
func (t ApplySignatureConfigs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplySignatureConfigs from hex string (Canton MCMS format)
func (t *ApplySignatureConfigs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplySignatureConfigsParams is a Record type
type ApplySignatureConfigsParams struct {
	SourceChainSelectorsToRemove []types.NUMERIC   `json:"sourceChainSelectorsToRemove"`
	SignatureConfigs             []SignatureConfig `json:"signatureConfigs"`
}

// ToMap converts ApplySignatureConfigsParams to a map for DAML arguments
func (t ApplySignatureConfigsParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.SourceChainSelectorsToRemove))
		for _, e := range t.SourceChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["signatureConfigs"] = func() []any {
		res := make([]any, 0, len(t.SignatureConfigs))
		for _, e := range t.SignatureConfigs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplySignatureConfigsParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplySignatureConfigsParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplySignatureConfigsParams to hex string (Canton MCMS format)
func (t ApplySignatureConfigsParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplySignatureConfigsParams from hex string (Canton MCMS format)
func (t *ApplySignatureConfigsParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFee is a Record type
type CalculateFee struct {
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts CalculateFee to a map for DAML arguments
func (t CalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CalculateFee to hex string (Canton MCMS format)
func (t CalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFee from hex string (Canton MCMS format)
func (t *CalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFeeMCMSParams is CalculateFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CalculateFeeMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
}

// MarshalHex encodes CalculateFeeMCMSParams to hex string for MCMS operationData.
func (t CalculateFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFeeMCMSParams from hex string.
func (t *CalculateFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifier is a Template type
type CommitteeVerifier struct {
	InstanceId                   types.TEXT                          `json:"instanceId"`
	Owner                        types.PARTY                         `json:"owner"`
	CcipOwner                    types.PARTY                         `json:"ccipOwner"`
	VersionTag                   types.TEXT                          `json:"versionTag" hex:"bytes"`
	AllowListAdmin               *types.PARTY                        `json:"allowListAdmin" hex:"optional"`
	MessageSentObservers         []types.PARTY                       `json:"messageSentObservers"`
	StorageLocations             []types.TEXT                        `json:"storageLocations"`
	StorageLocationsAdmin        types.PARTY                         `json:"storageLocationsAdmin"`
	PendingStorageLocationsAdmin types.PARTY                         `json:"pendingStorageLocationsAdmin"`
	RemoteChainConfigs           map[types.NUMERIC]RemoteChainConfig `json:"remoteChainConfigs"`
	SignerConfigs                map[types.NUMERIC]SignatureConfig   `json:"signerConfigs"`
	Deps                         CommitteeVerifierDeps               `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CommitteeVerifier) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CommitteeVerifier) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CommitteeVerifier) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["versionTag"] = string(t.VersionTag)

	if t.AllowListAdmin != nil {
		args["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.AllowListAdmin).ToMap(),
		}
	} else {
		args["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocationsAdmin"] = t.StorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingStorageLocationsAdmin"] = t.PendingStorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signerConfigs"] = func() any {
		if t.SignerConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SignerConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CommitteeVerifier) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["versionTag"] = string(t.VersionTag)

	if t.AllowListAdmin != nil {
		args["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.AllowListAdmin).ToMap(),
		}
	} else {
		args["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocationsAdmin"] = t.StorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingStorageLocationsAdmin"] = t.PendingStorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signerConfigs"] = func() any {
		if t.SignerConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SignerConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CommitteeVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifier to hex string (Canton MCMS format)
func (t CommitteeVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifier from hex string (Canton MCMS format)
func (t *CommitteeVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CommitteeVerifier

// VerifyMessage exercises the VerifyMessage choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) VerifyMessage(contractID string, args VerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// VerifyMessageWithPackageID exercises the VerifyMessage choice using the provided package ID instead of package name
func (t CommitteeVerifier) VerifyMessageWithPackageID(contractID string, packageID string, args VerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// ApplyAllowListUpdates exercises the ApplyAllowListUpdates choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) ApplyAllowListUpdates(contractID string, args ApplyAllowListUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplyAllowListUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyAllowListUpdatesWithPackageID exercises the ApplyAllowListUpdates choice using the provided package ID instead of package name
func (t CommitteeVerifier) ApplyAllowListUpdatesWithPackageID(contractID string, packageID string, args ApplyAllowListUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplyAllowListUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyRemoteChainConfigUpdates exercises the ApplyRemoteChainConfigUpdates choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) ApplyRemoteChainConfigUpdates(contractID string, args ApplyRemoteChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplyRemoteChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyRemoteChainConfigUpdatesWithPackageID exercises the ApplyRemoteChainConfigUpdates choice using the provided package ID instead of package name
func (t CommitteeVerifier) ApplyRemoteChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyRemoteChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplyRemoteChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplySignatureConfigs exercises the ApplySignatureConfigs choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) ApplySignatureConfigs(contractID string, args ApplySignatureConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplySignatureConfigs",
		Arguments:  argsToMap(args),
	}
}

// ApplySignatureConfigsWithPackageID exercises the ApplySignatureConfigs choice using the provided package ID instead of package name
func (t CommitteeVerifier) ApplySignatureConfigsWithPackageID(contractID string, packageID string, args ApplySignatureConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ApplySignatureConfigs",
		Arguments:  argsToMap(args),
	}
}

// AcceptStorageLocationsAdmin exercises the AcceptStorageLocationsAdmin choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) AcceptStorageLocationsAdmin(contractID string, args AcceptStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "AcceptStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// AcceptStorageLocationsAdminWithPackageID exercises the AcceptStorageLocationsAdmin choice using the provided package ID instead of package name
func (t CommitteeVerifier) AcceptStorageLocationsAdminWithPackageID(contractID string, packageID string, args AcceptStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "AcceptStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// CalculateFee exercises the CalculateFee choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CalculateFee(contractID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CalculateFeeWithPackageID exercises the CalculateFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CalculateFeeWithPackageID(contractID string, packageID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// ForwardToVerifier exercises the ForwardToVerifier choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) ForwardToVerifier(contractID string, args ForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// ForwardToVerifierWithPackageID exercises the ForwardToVerifier choice using the provided package ID instead of package name
func (t CommitteeVerifier) ForwardToVerifierWithPackageID(contractID string, packageID string, args ForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CrossChainVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CommitteeVerifier) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CrossChainVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetFee exercises the GetFee choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfig exercises the SetDynamicConfig choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) SetDynamicConfig(contractID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfigWithPackageID exercises the SetDynamicConfig choice using the provided package ID instead of package name
func (t CommitteeVerifier) SetDynamicConfigWithPackageID(contractID string, packageID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetDeps exercises the SetDeps choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t CommitteeVerifier) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// TransferStorageLocationsAdmin exercises the TransferStorageLocationsAdmin choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) TransferStorageLocationsAdmin(contractID string, args TransferStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "TransferStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// TransferStorageLocationsAdminWithPackageID exercises the TransferStorageLocationsAdmin choice using the provided package ID instead of package name
func (t CommitteeVerifier) TransferStorageLocationsAdminWithPackageID(contractID string, packageID string, args TransferStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "TransferStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// UpdateStorageLocations exercises the UpdateStorageLocations choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) UpdateStorageLocations(contractID string, args UpdateStorageLocations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "UpdateStorageLocations",
		Arguments:  argsToMap(args),
	}
}

// UpdateStorageLocationsWithPackageID exercises the UpdateStorageLocations choice using the provided package ID instead of package name
func (t CommitteeVerifier) UpdateStorageLocationsWithPackageID(contractID string, packageID string, args UpdateStorageLocations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "UpdateStorageLocations",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this CommitteeVerifier contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CommitteeVerifier) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessage exercises the CrossChainVerifier_VerifyMessage choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierVerifyMessage(contractID string, args extensionapi.CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessageWithPackageID exercises the CrossChainVerifier_VerifyMessage choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierVerifyMessageWithPackageID(contractID string, packageID string, args extensionapi.CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFee exercises the CrossChainVerifier_CalculateFee choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierCalculateFee(contractID string, args extensionapi.CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFeeWithPackageID exercises the CrossChainVerifier_CalculateFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierCalculateFeeWithPackageID(contractID string, packageID string, args extensionapi.CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierGetFee exercises the CrossChainVerifier_GetFee choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierGetFee(contractID string, args extensionapi.CrossChainVerifierGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_GetFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierGetFeeWithPackageID exercises the CrossChainVerifier_GetFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierGetFeeWithPackageID(contractID string, packageID string, args extensionapi.CrossChainVerifierGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_GetFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifier exercises the CrossChainVerifier_ForwardToVerifier choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifier(contractID string, args extensionapi.CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifierWithPackageID exercises the CrossChainVerifier_ForwardToVerifier choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifierWithPackageID(contractID string, packageID string, args extensionapi.CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CommitteeVerifier

var _ api.IMCMSReceiver = (*CommitteeVerifier)(nil)

var _ extensionapi.IICrossChainVerifier = (*CommitteeVerifier)(nil)

// CommitteeVerifierDeps is a Record type
type CommitteeVerifierDeps struct {
	RmnRemote chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts CommitteeVerifierDeps to a map for DAML arguments
func (t CommitteeVerifierDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	return m
}

func (t CommitteeVerifierDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierDeps to hex string (Canton MCMS format)
func (t CommitteeVerifierDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierDeps from hex string (Canton MCMS format)
func (t *CommitteeVerifierDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DynamicConfig is a Record type
type DynamicConfig struct {
	AllowListAdmin       *types.PARTY  `json:"allowListAdmin" hex:"optional"`
	MessageSentObservers []types.PARTY `json:"messageSentObservers"`
}

// ToMap converts DynamicConfig to a map for DAML arguments
func (t DynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	if t.AllowListAdmin != nil {
		m["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.AllowListAdmin).ToMap(),
		}
	} else {
		m["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t DynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DynamicConfig to hex string (Canton MCMS format)
func (t DynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DynamicConfig from hex string (Canton MCMS format)
func (t *DynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForwardToVerifier is a Record type
type ForwardToVerifier struct {
	RmnRemoteCid      types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT                                 `json:"verifierArgs"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts ForwardToVerifier to a map for DAML arguments
func (t ForwardToVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["verifierArgs"] = string(t.VerifierArgs)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ForwardToVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForwardToVerifier to hex string (Canton MCMS format)
func (t ForwardToVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForwardToVerifier from hex string (Canton MCMS format)
func (t *ForwardToVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForwardToVerifierMCMSParams is ForwardToVerifier without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ForwardToVerifierMCMSParams struct {
	RmnRemoteCid      types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT                                 `json:"verifierArgs"`
}

// MarshalHex encodes ForwardToVerifierMCMSParams to hex string for MCMS operationData.
func (t ForwardToVerifierMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForwardToVerifierMCMSParams from hex string.
func (t *ForwardToVerifierMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee is a Record type
type GetFee struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee to hex string (Canton MCMS format)
func (t GetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee from hex string (Canton MCMS format)
func (t *GetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeMCMSParams is GetFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFeeMCMSParams struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// MarshalHex encodes GetFeeMCMSParams to hex string for MCMS operationData.
func (t GetFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeMCMSParams from hex string.
func (t *GetFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoteChainConfig is a Record type
type RemoteChainConfig struct {
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
	GasForVerification types.INT64   `json:"gasForVerification"`
	PayloadSizeBytes   types.INT64   `json:"payloadSizeBytes"`
	AllowListEnabled   types.BOOL    `json:"allowListEnabled"`
	AllowedSendersList []types.PARTY `json:"allowedSendersList"`
}

// ToMap converts RemoteChainConfig to a map for DAML arguments
func (t RemoteChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeUSDCents"] = t.FeeUSDCents

	m["gasForVerification"] = int64(t.GasForVerification)

	m["payloadSizeBytes"] = int64(t.PayloadSizeBytes)

	m["allowListEnabled"] = bool(t.AllowListEnabled)

	m["allowedSendersList"] = func() []any {
		res := make([]any, 0, len(t.AllowedSendersList))
		for _, e := range t.AllowedSendersList {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoteChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoteChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoteChainConfig to hex string (Canton MCMS format)
func (t RemoteChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoteChainConfig from hex string (Canton MCMS format)
func (t *RemoteChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoteChainConfigArgs is a Record type
type RemoteChainConfigArgs struct {
	RemoteChainSelector types.NUMERIC `json:"remoteChainSelector"`
	FeeUSDCents         types.NUMERIC `json:"feeUSDCents"`
	GasForVerification  types.INT64   `json:"gasForVerification"`
	PayloadSizeBytes    types.INT64   `json:"payloadSizeBytes"`
	AllowListEnabled    types.BOOL    `json:"allowListEnabled"`
}

// ToMap converts RemoteChainConfigArgs to a map for DAML arguments
func (t RemoteChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["feeUSDCents"] = t.FeeUSDCents

	m["gasForVerification"] = int64(t.GasForVerification)

	m["payloadSizeBytes"] = int64(t.PayloadSizeBytes)

	m["allowListEnabled"] = bool(t.AllowListEnabled)

	return m
}

func (t RemoteChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoteChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoteChainConfigArgs to hex string (Canton MCMS format)
func (t RemoteChainConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoteChainConfigArgs from hex string (Canton MCMS format)
func (t *RemoteChainConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDeps is a Record type
type SetDeps struct {
	NewDeps SetDepsParams `json:"newDeps"`
}

// ToMap converts SetDeps to a map for DAML arguments
func (t SetDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = model.NestedToDAMLValue(t.NewDeps)

	return m
}

func (t SetDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDeps to hex string (Canton MCMS format)
func (t SetDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDeps from hex string (Canton MCMS format)
func (t *SetDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDepsParams is a Record type
type SetDepsParams struct {
	RmnRemote *chainlinkapi.RawInstanceAddress `json:"rmnRemote" hex:"optional"`
}

// ToMap converts SetDepsParams to a map for DAML arguments
func (t SetDepsParams) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RmnRemote != nil {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RmnRemote),
		}
	} else {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetDepsParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDepsParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDepsParams to hex string (Canton MCMS format)
func (t SetDepsParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDepsParams from hex string (Canton MCMS format)
func (t *SetDepsParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDynamicConfig is a Record type
type SetDynamicConfig struct {
	DynamicConfig DynamicConfig `json:"dynamicConfig"`
}

// ToMap converts SetDynamicConfig to a map for DAML arguments
func (t SetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["dynamicConfig"] = model.NestedToDAMLValue(t.DynamicConfig)

	return m
}

func (t SetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfig to hex string (Canton MCMS format)
func (t SetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfig from hex string (Canton MCMS format)
func (t *SetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDynamicConfigParams is a Record type
type SetDynamicConfigParams struct {
	DynamicConfig DynamicConfig `json:"dynamicConfig"`
}

// ToMap converts SetDynamicConfigParams to a map for DAML arguments
func (t SetDynamicConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["dynamicConfig"] = model.NestedToDAMLValue(t.DynamicConfig)

	return m
}

func (t SetDynamicConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfigParams to hex string (Canton MCMS format)
func (t SetDynamicConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfigParams from hex string (Canton MCMS format)
func (t *SetDynamicConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SignatureConfig is a Record type
type SignatureConfig struct {
	SourceChainSelector types.NUMERIC `json:"sourceChainSelector"`
	Threshold           types.INT64   `json:"threshold"`
	SignerKeys          []types.TEXT  `json:"signerKeys"`
}

// ToMap converts SignatureConfig to a map for DAML arguments
func (t SignatureConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["threshold"] = int64(t.Threshold)

	m["signerKeys"] = func() []any {
		res := make([]any, 0, len(t.SignerKeys))
		for _, e := range t.SignerKeys {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t SignatureConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SignatureConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SignatureConfig to hex string (Canton MCMS format)
func (t SignatureConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SignatureConfig from hex string (Canton MCMS format)
func (t *SignatureConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferStorageLocationsAdmin is a Record type
type TransferStorageLocationsAdmin struct {
	NewAdmin types.PARTY `json:"newAdmin"`
}

// ToMap converts TransferStorageLocationsAdmin to a map for DAML arguments
func (t TransferStorageLocationsAdmin) ToMap() map[string]any {
	m := make(map[string]any)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t TransferStorageLocationsAdmin) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferStorageLocationsAdmin) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferStorageLocationsAdmin to hex string (Canton MCMS format)
func (t TransferStorageLocationsAdmin) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferStorageLocationsAdmin from hex string (Canton MCMS format)
func (t *TransferStorageLocationsAdmin) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferStorageLocationsAdminParams is a Record type
type TransferStorageLocationsAdminParams struct {
	NewAdmin types.PARTY `json:"newAdmin"`
}

// ToMap converts TransferStorageLocationsAdminParams to a map for DAML arguments
func (t TransferStorageLocationsAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t TransferStorageLocationsAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferStorageLocationsAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferStorageLocationsAdminParams to hex string (Canton MCMS format)
func (t TransferStorageLocationsAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferStorageLocationsAdminParams from hex string (Canton MCMS format)
func (t *TransferStorageLocationsAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateStorageLocations is a Record type
type UpdateStorageLocations struct {
	NewLocations []types.TEXT `json:"newLocations"`
}

// ToMap converts UpdateStorageLocations to a map for DAML arguments
func (t UpdateStorageLocations) ToMap() map[string]any {
	m := make(map[string]any)

	m["newLocations"] = func() []any {
		res := make([]any, 0, len(t.NewLocations))
		for _, e := range t.NewLocations {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UpdateStorageLocations) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateStorageLocations) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateStorageLocations to hex string (Canton MCMS format)
func (t UpdateStorageLocations) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateStorageLocations from hex string (Canton MCMS format)
func (t *UpdateStorageLocations) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateStorageLocationsParams is a Record type
type UpdateStorageLocationsParams struct {
	NewLocations []types.TEXT `json:"newLocations"`
}

// ToMap converts UpdateStorageLocationsParams to a map for DAML arguments
func (t UpdateStorageLocationsParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["newLocations"] = func() []any {
		res := make([]any, 0, len(t.NewLocations))
		for _, e := range t.NewLocations {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UpdateStorageLocationsParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateStorageLocationsParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateStorageLocationsParams to hex string (Canton MCMS format)
func (t UpdateStorageLocationsParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateStorageLocationsParams from hex string (Canton MCMS format)
func (t *UpdateStorageLocationsParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifyMessage is a Record type
type VerifyMessage struct {
	RmnRemoteCid        types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext        splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	VerifierResults     types.TEXT                                 `json:"verifierResults"`
	Caller              types.PARTY                                `json:"caller"`
}

// ToMap converts VerifyMessage to a map for DAML arguments
func (t VerifyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["verifierResults"] = string(t.VerifierResults)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t VerifyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *VerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes VerifyMessage to hex string (Canton MCMS format)
func (t VerifyMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyMessage from hex string (Canton MCMS format)
func (t *VerifyMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifyMessageMCMSParams is VerifyMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type VerifyMessageMCMSParams struct {
	RmnRemoteCid        types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext        splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	VerifierResults     types.TEXT                                 `json:"verifierResults"`
}

// MarshalHex encodes VerifyMessageMCMSParams to hex string for MCMS operationData.
func (t VerifyMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifyMessageMCMSParams from hex string.
func (t *VerifyMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AcceptStorageLocationsAdmin(args AcceptStorageLocationsAdmin) (*bind.EncodedChoice, error)
	ApplyAllowListUpdates(args ApplyAllowListUpdates) (*bind.EncodedChoice, error)
	ApplyAllowListUpdatesMCMSParams(args ApplyAllowListUpdatesMCMSParams) (*bind.EncodedChoice, error)
	ApplyAllowListUpdatesParams(args ApplyAllowListUpdatesParams) (*bind.EncodedChoice, error)
	ApplyRemoteChainConfigUpdates(args ApplyRemoteChainConfigUpdates) (*bind.EncodedChoice, error)
	ApplyRemoteChainConfigUpdatesParams(args ApplyRemoteChainConfigUpdatesParams) (*bind.EncodedChoice, error)
	ApplySignatureConfigs(args ApplySignatureConfigs) (*bind.EncodedChoice, error)
	ApplySignatureConfigsParams(args ApplySignatureConfigsParams) (*bind.EncodedChoice, error)
	CalculateFee(args CalculateFee) (*bind.EncodedChoice, error)
	CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	ForwardToVerifier(args ForwardToVerifier) (*bind.EncodedChoice, error)
	ForwardToVerifierMCMSParams(args ForwardToVerifierMCMSParams) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error)
	SetDeps(args SetDeps) (*bind.EncodedChoice, error)
	SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error)
	SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error)
	SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error)
	TransferStorageLocationsAdmin(args TransferStorageLocationsAdmin) (*bind.EncodedChoice, error)
	TransferStorageLocationsAdminParams(args TransferStorageLocationsAdminParams) (*bind.EncodedChoice, error)
	UpdateStorageLocations(args UpdateStorageLocations) (*bind.EncodedChoice, error)
	UpdateStorageLocationsParams(args UpdateStorageLocationsParams) (*bind.EncodedChoice, error)
	VerifyMessage(args VerifyMessage) (*bind.EncodedChoice, error)
	VerifyMessageMCMSParams(args VerifyMessageMCMSParams) (*bind.EncodedChoice, error)
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

// AcceptStorageLocationsAdmin encodes parameters for the AcceptStorageLocationsAdmin choice.
func (e *encoder) AcceptStorageLocationsAdmin(args AcceptStorageLocationsAdmin) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptStorageLocationsAdmin", args)
}

// ApplyAllowListUpdates encodes parameters for the ApplyAllowListUpdates choice.
func (e *encoder) ApplyAllowListUpdates(args ApplyAllowListUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyAllowListUpdates", args)
}

// ApplyAllowListUpdatesMCMSParams encodes MCMS parameters (without Caller) for the ApplyAllowListUpdates choice.
func (e *encoder) ApplyAllowListUpdatesMCMSParams(args ApplyAllowListUpdatesMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyAllowListUpdates", args)
}

// ApplyAllowListUpdatesParams encodes parameters for the ApplyAllowListUpdates choice.
func (e *encoder) ApplyAllowListUpdatesParams(args ApplyAllowListUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyAllowListUpdates", args)
}

// ApplyRemoteChainConfigUpdates encodes parameters for the ApplyRemoteChainConfigUpdates choice.
func (e *encoder) ApplyRemoteChainConfigUpdates(args ApplyRemoteChainConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyRemoteChainConfigUpdates", args)
}

// ApplyRemoteChainConfigUpdatesParams encodes parameters for the ApplyRemoteChainConfigUpdates choice.
func (e *encoder) ApplyRemoteChainConfigUpdatesParams(args ApplyRemoteChainConfigUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyRemoteChainConfigUpdates", args)
}

// ApplySignatureConfigs encodes parameters for the ApplySignatureConfigs choice.
func (e *encoder) ApplySignatureConfigs(args ApplySignatureConfigs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplySignatureConfigs", args)
}

// ApplySignatureConfigsParams encodes parameters for the ApplySignatureConfigs choice.
func (e *encoder) ApplySignatureConfigsParams(args ApplySignatureConfigsParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplySignatureConfigs", args)
}

// CalculateFee encodes parameters for the CalculateFee choice.
func (e *encoder) CalculateFee(args CalculateFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// CalculateFeeMCMSParams encodes MCMS parameters (without Caller) for the CalculateFee choice.
func (e *encoder) CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// ForwardToVerifier encodes parameters for the ForwardToVerifier choice.
func (e *encoder) ForwardToVerifier(args ForwardToVerifier) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ForwardToVerifier", args)
}

// ForwardToVerifierMCMSParams encodes MCMS parameters (without Caller) for the ForwardToVerifier choice.
func (e *encoder) ForwardToVerifierMCMSParams(args ForwardToVerifierMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ForwardToVerifier", args)
}

// GetFee encodes parameters for the GetFee choice.
func (e *encoder) GetFee(args GetFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetFeeMCMSParams encodes MCMS parameters (without Caller) for the GetFee choice.
func (e *encoder) GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// SetDeps encodes parameters for the SetDeps choice.
func (e *encoder) SetDeps(args SetDeps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// SetDepsParams encodes parameters for the SetDeps choice.
func (e *encoder) SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// SetDynamicConfig encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// SetDynamicConfigParams encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// TransferStorageLocationsAdmin encodes parameters for the TransferStorageLocationsAdmin choice.
func (e *encoder) TransferStorageLocationsAdmin(args TransferStorageLocationsAdmin) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferStorageLocationsAdmin", args)
}

// TransferStorageLocationsAdminParams encodes parameters for the TransferStorageLocationsAdmin choice.
func (e *encoder) TransferStorageLocationsAdminParams(args TransferStorageLocationsAdminParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferStorageLocationsAdmin", args)
}

// UpdateStorageLocations encodes parameters for the UpdateStorageLocations choice.
func (e *encoder) UpdateStorageLocations(args UpdateStorageLocations) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateStorageLocations", args)
}

// UpdateStorageLocationsParams encodes parameters for the UpdateStorageLocations choice.
func (e *encoder) UpdateStorageLocationsParams(args UpdateStorageLocationsParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateStorageLocations", args)
}

// VerifyMessage encodes parameters for the VerifyMessage choice.
func (e *encoder) VerifyMessage(args VerifyMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyMessage", args)
}

// VerifyMessageMCMSParams encodes MCMS parameters (without Caller) for the VerifyMessage choice.
func (e *encoder) VerifyMessageMCMSParams(args VerifyMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("VerifyMessage", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
