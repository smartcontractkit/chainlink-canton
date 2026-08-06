package ccip

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	burnminttokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	ccipcore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// LockOrBurnInput identifies contracts for a pool LockOrBurn exercise.
type LockOrBurnInput struct {
	PoolCID               string
	TokenAdminRegistryCID string
	TokenConfigCID        string
	RMNRemoteCID          string
	SendingMessageCID     string
	SenderInputCids       []string
	Amount                string
	Sender                string
	PoolOwner             string
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext
	Disclosures           PoolSendDisclosures
}

// LockOrBurnResult is the parsed LockOrBurn exercise output.
type LockOrBurnResult struct {
	SenderChangeCids  []string
	PoolChangeCids    []string
	SendingMessageCid string
}

// LockOrBurn exercises BurnMintTokenPool.LockOrBurn and returns change holdings plus the updated SendingMessage CID.
func LockOrBurn(ctx context.Context, client ledger.Client, input LockOrBurnInput) (LockOrBurnResult, error) {
	senderInputs := make([]types.CONTRACT_ID, 0, len(input.SenderInputCids))
	for _, cid := range input.SenderInputCids {
		senderInputs = append(senderInputs, types.CONTRACT_ID(cid))
	}

	args := burnminttokenpool.LockOrBurn{
		TokenAdminRegistryCid: types.CONTRACT_ID(input.TokenAdminRegistryCID),
		TokenConfigCid:        types.CONTRACT_ID(input.TokenConfigCID),
		RmnRemoteCid:          types.CONTRACT_ID(input.RMNRemoteCID),
		Context:               input.ExtraContext,
		SendingMessageCid:     types.CONTRACT_ID(input.SendingMessageCID),
		SenderInputCids:       senderInputs,
		Amount:                types.NUMERIC(input.Amount),
		Caller:                types.PARTY(input.Sender),
	}

	actAs := []string{input.Sender}
	if input.PoolOwner != "" && input.PoolOwner != input.Sender {
		actAs = append(actAs, input.PoolOwner)
	}

	res, err := client.SubmitExerciseMulti(ctx, actAs, burnminttokenpool.BurnMintTokenPool{}, input.PoolCID, "LockOrBurn", args, input.Disclosures.All())
	if err != nil {
		return LockOrBurnResult{}, fmt.Errorf("lock or burn: %w", err)
	}

	if parsed, err := ledger.ParseLockOrBurnResult(res.GetTransaction()); err == nil {
		result := LockOrBurnResult{SendingMessageCid: string(parsed.SendingMessageCid)}
		for _, cid := range parsed.SenderChangeCids {
			result.SenderChangeCids = append(result.SenderChangeCids, string(cid))
		}
		for _, cid := range parsed.PoolChangeCids {
			result.PoolChangeCids = append(result.PoolChangeCids, string(cid))
		}

		return result, nil
	}

	return lockOrBurnResultFromCreatedHoldings(res.GetTransaction(), input)
}

func lockOrBurnResultFromCreatedHoldings(tx *apiv2.Transaction, input LockOrBurnInput) (LockOrBurnResult, error) {
	senderChanges := ledger.CreatedHoldingsForOwner(tx, input.Sender)
	poolChanges := ledger.CreatedHoldingsForOwner(tx, input.PoolOwner)
	if len(senderChanges) == 0 {
		return LockOrBurnResult{}, fmt.Errorf("sender change Holding not created")
	}

	sendingMessageCID := input.SendingMessageCID
	if cid, ok := ledger.CreatedContractID(tx, "SendingMessage"); ok {
		sendingMessageCID = cid
	}

	return LockOrBurnResult{
		SenderChangeCids:  senderChanges,
		PoolChangeCids:    poolChanges,
		SendingMessageCid: sendingMessageCID,
	}, nil
}

// CreateSendingMessageInput identifies fields for a FeeFinalized SendingMessage used in LockOrBurn tests.
type CreateSendingMessageInput struct {
	CcipParty           string
	Sender              string
	InstrumentId        splice_api_token_holding_v1.InstrumentId
	DestChainSelector   string
	SourceChainSelector string
	PoolInstanceID      string
	PoolOwner           string
	RmnRemote           contracts.RawInstanceAddress
	TokenAdminRegistry  contracts.RawInstanceAddress
	FeeQuoter           contracts.RawInstanceAddress
}

// CreateSendingMessage creates a FeeFinalized SendingMessage on the CCIP party for outbound pool tests.
func CreateSendingMessage(ctx context.Context, client ledger.Client, input CreateSendingMessageInput) (string, error) {
	sourceChainSelector := input.SourceChainSelector
	if sourceChainSelector == "" {
		sourceChainSelector = "123"
	}

	executionMode := ccipcore.ExecutionModeExecutionMode_NoExecutor
	expectedInstrument := input.InstrumentId
	emptyOutboundCCVs := []chainlinkapi.RawInstanceAddress{}
	tokenSendFee := ccipcore.TokenSendFee{
		PoolInstanceId:    types.TEXT(input.PoolInstanceID),
		PoolOwner:         types.PARTY(input.PoolOwner),
		FeeUSDCents:       types.NUMERIC("0"),
		DestGasOverhead:   types.INT64(0),
		DestBytesOverhead: types.INT64(64),
	}

	msg := ccipcore.SendingMessage{
		CcipOwner: types.PARTY(input.CcipParty),
		Sender:    types.PARTY(input.Sender),
		Deps: ccipcore.SendingMessageDeps{
			Router:             rawInstanceAddressBinding("test-router", input.Sender),
			OnRamp:             rawInstanceAddressBinding("test-onramp", input.CcipParty),
			GlobalConfig:       rawInstanceAddressBinding("test-globalconfig", input.CcipParty),
			RmnRemote:          bindRawInstanceAddress(input.RmnRemote),
			TokenAdminRegistry: bindRawInstanceAddress(input.TokenAdminRegistry),
			FeeQuoter:          bindRawInstanceAddress(input.FeeQuoter),
		},
		DestChainSelector:              types.NUMERIC(input.DestChainSelector),
		DestAddressBytesLength:         types.INT64(20),
		SequenceNumber:                 types.NUMERIC("1"),
		DestDefaultCCVs:                nil,
		RequiredCCVs:                   nil,
		ExecutorAddress:                types.TEXT(instanceAddressHex("default-executor", input.CcipParty)),
		ExecutionMode:                  &executionMode,
		SourceChainSelector:            types.NUMERIC(sourceChainSelector),
		SenderAddress:                  types.TEXT("0000000000000000000000000000000000000001"),
		Receiver:                       types.TEXT("0000000000000000000000000000000000000001"),
		Payload:                        types.TEXT(""),
		ExecutionGasLimit:              types.INT64(0),
		CcipReceiveGasLimit:            types.INT64(100000),
		CcvAndExecutorHash:             types.TEXT(""),
		OnRampAddress:                  types.TEXT("0000000000000000000000000000000000000001"),
		OffRampAddress:                 types.TEXT("0000000000000000000000000000000000000002"),
		TokenReceiver:                  types.TEXT("0000000000000000000000000000000000000003"),
		TokenArgs:                      types.TEXT(""),
		FeeToken:                       input.InstrumentId,
		NetworkFeeUSDCents:             types.NUMERIC("0"),
		ExpectedTokenInstrumentId:      &expectedInstrument,
		TokenAmountBeforeTokenPoolFees: types.NUMERIC("0"),
		OutboundPoolCCVs:               &emptyOutboundCCVs,
		ExecutorArgs:                   types.TEXT(""),
		ExecutorDestGasLimit:           types.INT64(0),
		ExecutorDestBytesOverhead:      types.INT64(0),
		ExecutorFeeTokenAmount:         types.NUMERIC("0"),
		ObservingParties:               []types.PARTY{types.PARTY(input.Sender)},
		CcvFees:                        nil,
		TokenSendFee:                   &tokenSendFee,
		CcvFeeTokenAmounts:             nil,
		TokenSendFeeTokenAmount:        types.NUMERIC("0"),
		NetworkFeeTokenAmount:          types.NUMERIC("0"),
		VerifierData:                   nil,
		CcvOwners:                      nil,
		EncodedMessage:                 types.TEXT(""),
		MessageId:                      types.TEXT(""),
		State:                          ccipcore.SendingMessageStateSendingMessageState_FeeFinalized,
	}

	res, err := client.SubmitCreate(ctx, input.CcipParty, msg)
	if err != nil {
		return "", fmt.Errorf("create SendingMessage: %w", err)
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "SendingMessage")
	if !ok {
		return "", fmt.Errorf("SendingMessage not created")
	}

	return cid, nil
}

// FetchSendingMessage loads an active SendingMessage by contract ID.
func FetchSendingMessage(ctx context.Context, client ledger.Client, ccipParty, contractID string) (ccipcore.SendingMessage, error) {
	tpl := contracts.IdentifierFromBinding(ccipcore.SendingMessage{})
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.ForParty(ccipParty), tpl)
	if err != nil {
		return ccipcore.SendingMessage{}, fmt.Errorf("list SendingMessage: %w", err)
	}
	for _, ac := range active {
		if ac.GetCreatedEvent().GetContractId() != contractID {
			continue
		}
		msg, err := bindings.UnmarshalCreatedEvent[ccipcore.SendingMessage](ac.GetCreatedEvent())
		if err != nil {
			return ccipcore.SendingMessage{}, fmt.Errorf("unmarshal SendingMessage: %w", err)
		}

		return *msg, nil
	}

	return ccipcore.SendingMessage{}, fmt.Errorf("SendingMessage %s not in ACS", contractID)
}

// SetBurnMintFactoryInput identifies contracts for pointing TokenConfig at a burn-mint factory.
type SetBurnMintFactoryInput struct {
	TokenAdminRegistryCID string
	TokenConfigCID        string
	InstrumentId          splice_api_token_holding_v1.InstrumentId
	BurnMintFactoryCID    string
	CcipParty             string
	PoolOwnerParty        string
	CcipClient            ledger.Client
	PoolOwnerClient       ledger.Client
}

// SetBurnMintFactory exercises TokenAdminRegistry.SetBurnMintFactory and returns the updated TokenConfig CID.
func SetBurnMintFactory(ctx context.Context, client ledger.Client, input SetBurnMintFactoryInput) (string, error) {
	ccipClient := input.CcipClient
	if ccipClient == nil {
		ccipClient = client
	}
	poolOwnerClient := input.PoolOwnerClient
	if poolOwnerClient == nil {
		poolOwnerClient = client
	}

	tarDisclosed, err := registry.DiscloseByID(ctx, ccipClient, input.CcipParty, input.TokenAdminRegistryCID)
	if err != nil {
		return "", err
	}

	factoryCID := types.CONTRACT_ID(input.BurnMintFactoryCID)
	res, err := poolOwnerClient.SubmitExerciseMulti(ctx, []string{input.PoolOwnerParty}, ccipcore.TokenAdminRegistry{}, input.TokenAdminRegistryCID, "SetBurnMintFactory",
		ccipcore.SetBurnMintFactory{
			TokenConfigCid:  types.CONTRACT_ID(input.TokenConfigCID),
			InstrumentId:    input.InstrumentId,
			BurnMintFactory: &factoryCID,
			Caller:          types.PARTY(input.PoolOwnerParty),
		}, []*apiv2.DisclosedContract{tarDisclosed})
	if err != nil {
		return "", fmt.Errorf("set burn mint factory: %w", err)
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "TokenConfig")
	if !ok {
		return "", fmt.Errorf("TokenConfig not created after SetBurnMintFactory")
	}

	return cid, nil
}

// EncodeUint256Hex encodes a base-10 integer string as a 32-byte hex string (CCIP MessageCodec style).
func EncodeUint256Hex(decimalValue string) string {
	n := new(big.Int)
	_, ok := n.SetString(decimalValue, 10)
	if !ok {
		panic("invalid decimal: " + decimalValue)
	}
	hexStr := n.Text(16)

	return strings.Repeat("0", 64-len(hexStr)) + hexStr
}

func rawInstanceAddressBinding(instanceID, owner string) chainlinkapi.RawInstanceAddress {
	return bindRawInstanceAddress(contracts.NewRawInstanceAddress(contracts.InstanceID(instanceID), types.PARTY(owner)))
}

func bindRawInstanceAddress(raw contracts.RawInstanceAddress) chainlinkapi.RawInstanceAddress {
	return chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(raw.String())}
}

func instanceAddressHex(instanceID, owner string) string {
	raw := contracts.NewRawInstanceAddress(contracts.InstanceID(instanceID), types.PARTY(owner))

	return strings.TrimPrefix(raw.InstanceAddress().Hex(), "0x")
}
