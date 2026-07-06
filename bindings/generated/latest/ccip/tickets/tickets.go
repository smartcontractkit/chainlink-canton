package tickets

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipapi"
	ccipcodec "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
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
	PackageName = "ccip-tickets"
	PackageID   = "c2a6263f6b20be4ec65698c9d28180c629bcbc4fb63bda63665b5ab1f31f90a0"
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

// Consume2 is a Record type
type Consume2 struct {
}

// ToMap converts Consume2 to a map for DAML arguments
func (t Consume2) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t Consume2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Consume2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Consume2 to hex string (Canton MCMS format)
func (t Consume2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Consume2 from hex string (Canton MCMS format)
func (t *Consume2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner                    types.PARTY                                `json:"ccipOwner"`
	CcvOwners                    []types.PARTY                              `json:"ccvOwners"`
	VerifiedCCVs                 []chainlinkapi.RawInstanceAddress          `json:"verifiedCCVs"`
	RequiredInboundPoolCCVs      []chainlinkapi.RawInstanceAddress          `json:"requiredInboundPoolCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                                 `json:"tokenAdminRegistryInstanceId"`
	PoolAddress                  chainlinkapi.RawInstanceAddress            `json:"poolAddress"`
	PoolOwner                    types.PARTY                                `json:"poolOwner"`
	Receiver                     types.PARTY                                `json:"receiver"`
	TokenReceiver                types.PARTY                                `json:"tokenReceiver"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Amount                       types.TEXT                                 `json:"amount"`
	SourcePoolData               types.TEXT                                 `json:"sourcePoolData"`
	MessageId                    types.TEXT                                 `json:"messageId"`
	SourceChainSelector          types.NUMERIC                              `json:"sourceChainSelector"`
	Finality                     ccipcodec.FinalityConfig                   `json:"finality"`
	Context                      splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenReceiveTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TicketsV1", "TokenReceiveTicket")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenReceiveTicket) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.TicketsV1", "TokenReceiveTicket")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenReceiveTicket) CreateCommand() *model.CreateCommand {
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
	args["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredInboundPoolCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredInboundPoolCCVs))
		for _, e := range t.RequiredInboundPoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["amount"] = string(t.Amount)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = model.NestedToDAMLValue(t.Finality)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["context"] = model.NestedToDAMLValue(t.Context)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenReceiveTicket) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
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
	args["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredInboundPoolCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredInboundPoolCCVs))
		for _, e := range t.RequiredInboundPoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["amount"] = string(t.Amount)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = model.NestedToDAMLValue(t.Finality)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["context"] = model.NestedToDAMLValue(t.Context)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicket to hex string (Canton MCMS format)
func (t TokenReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicket from hex string (Canton MCMS format)
func (t *TokenReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenReceiveTicket

// Consume exercises the Consume choice on this TokenReceiveTicket contract via the IITokenReceiveTicket interface
// This method uses the package name in the template ID
func (t TokenReceiveTicket) Consume(contractID string, args ccipapi.Consume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TicketsV1", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Consume",
		Arguments:  argsToMap(args),
	}
}

// ConsumeWithPackageID exercises the Consume choice using the provided package ID instead of package name
func (t TokenReceiveTicket) ConsumeWithPackageID(contractID string, packageID string, args ccipapi.Consume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TicketsV1", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Consume",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenReceiveTicket contract via the IITokenReceiveTicket interface
// This method uses the package name in the template ID
func (t TokenReceiveTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TicketsV1", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenReceiveTicket) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TicketsV1", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Verify interface implementations for TokenReceiveTicket

var _ ccipapi.IITokenReceiveTicket = (*TokenReceiveTicket)(nil)

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
