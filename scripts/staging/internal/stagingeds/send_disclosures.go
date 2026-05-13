// Package stagingeds mirrors the CCIP → token pool (optional) → CCV → Executor EDS sequence used in ccip/devenv.
package stagingeds

import (
	"context"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// TokenPoolSendEDS holds token-pool-scoped send disclosures from EDS (when sending with TokenTransfer).
type TokenPoolSendEDS struct {
	ContractID         types.CONTRACT_ID
	PoolExtraContext   splice_api_token_metadata_v1.ChoiceContext
	RequiredCCVs       []string
	DisclosedContracts []*apiv2.DisclosedContract
}

// SendEDSOutcome aggregates EDS disclosures for a Canton → remote CCIP send.
type SendEDSOutcome struct {
	SendContext        splice_api_token_metadata_v1.ChoiceContext
	CcvSendInputs      []ccipsender.CCVSendInput
	CcvExtraArgs       []ccipclient.CCVExtraArg
	ExecutorInput      *ccipsender.ExecutorInput
	DisclosedContracts []*apiv2.DisclosedContract
	TokenPoolSend      *TokenPoolSendEDS
}

func newEDSClients(edsBase string, doer *http.Client) (
	ccipCl *oapiCCIP.ClientWithResponses,
	ccvCl *oapiCCV.ClientWithResponses,
	execCl *oapiExecutor.ClientWithResponses,
	tpCl *oapiTokenPool.ClientWithResponses,
	err error,
) {
	ccipCl, err = oapiCCIP.NewClientWithResponses(edsBase, oapiCCIP.WithHTTPClient(doer))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("ccip eds client: %w", err)
	}
	ccvCl, err = oapiCCV.NewClientWithResponses(edsBase, oapiCCV.WithHTTPClient(doer))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("ccv eds client: %w", err)
	}
	execCl, err = oapiExecutor.NewClientWithResponses(edsBase, oapiExecutor.WithHTTPClient(doer))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("executor eds client: %w", err)
	}
	tpCl, err = oapiTokenPool.NewClientWithResponses(edsBase, oapiTokenPool.WithHTTPClient(doer))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("token pool eds client: %w", err)
	}
	return ccipCl, ccvCl, execCl, tpCl, nil
}

// CollectSendDisclosures runs the EDS send disclosure chain (see ccip/devenv Send path).
// senderRequiredRawCCVs are RawInstanceAddress strings (e.g. instanceId@party).
// When tokenPoolEDSOverride is non-nil, GetTokenPoolSendDisclosure is skipped and the override supplies
// token pool contract id, choice context, required CCVs, and disclosed contracts (e.g. built from ledger).
// Otherwise when tokenPoolAddress is non-nil, GetTokenPoolSendDisclosure runs first and supplies tokenPoolRequiredCCVs for CCIP send.
func CollectSendDisclosures(
	ctx context.Context,
	edsBase string,
	httpClient *http.Client,
	outgoing oapiCommon.Message,
	senderRequiredRawCCVs []string,
	tokenPoolAddress *contracts.InstanceAddress,
	tokenPoolEDSOverride *TokenPoolSendEDS,
) (*SendEDSOutcome, error) {
	ccipCl, ccvCl, execCl, tpCl, err := newEDSClients(edsBase, httpClient)
	if err != nil {
		return nil, err
	}

	var tokenPoolRequired []string
	out := &SendEDSOutcome{}

	tokenPoolFetch := tokenPoolAddress != nil && tokenPoolEDSOverride == nil
	if tokenPoolEDSOverride != nil {
		tokenPoolRequired = tokenPoolEDSOverride.RequiredCCVs
		out.TokenPoolSend = tokenPoolEDSOverride
		out.DisclosedContracts = append(out.DisclosedContracts, tokenPoolEDSOverride.DisclosedContracts...)
	} else if tokenPoolFetch {
		tpd, err := eds.GetTokenPoolSendDisclosure(ctx, tpCl, outgoing, *tokenPoolAddress)
		if err != nil {
			return nil, fmt.Errorf("token pool send disclosure: %w", err)
		}
		tokenPoolRequired = tpd.RequiredCCVs
		out.TokenPoolSend = &TokenPoolSendEDS{
			ContractID:         types.CONTRACT_ID(tpd.ContractId),
			PoolExtraContext:   tpd.ChoiceContext,
			RequiredCCVs:       tpd.RequiredCCVs,
			DisclosedContracts: tpd.DisclosedContracts,
		}
		out.DisclosedContracts = append(out.DisclosedContracts, tpd.DisclosedContracts...)
	}
	if !tokenPoolFetch {
		_ = tpCl
	}

	ccipDisc, err := eds.GetCCIPSendDisclosure(ctx, ccipCl, outgoing, senderRequiredRawCCVs, tokenPoolRequired)
	if err != nil {
		return nil, fmt.Errorf("ccip send disclosure: %w", err)
	}
	out.SendContext = ccipDisc.ChoiceContext
	out.DisclosedContracts = append(out.DisclosedContracts, ccipDisc.DisclosedContracts...)

	var allCCVHex []string
	for _, v := range ccipDisc.CCVs {
		ccvAddr, err := ccvInstanceAddressFromEDS(v)
		if err != nil {
			return nil, fmt.Errorf("parse ccv address %q: %w", v, err)
		}
		ccvSend, err := eds.GetCCVSendDisclosure(ctx, ccvCl, outgoing, ccvAddr)
		if err != nil {
			return nil, fmt.Errorf("ccv send disclosure for %q: %w", v, err)
		}
		out.CcvSendInputs = append(out.CcvSendInputs, ccipsender.CCVSendInput{
			CcvAddress:      ccvSend.Address.Binding(),
			CcvCid:          types.CONTRACT_ID(ccvSend.ContractId),
			CcvExtraContext: ccvSend.ChoiceContext,
		})
		out.CcvExtraArgs = append(out.CcvExtraArgs, ccipclient.CCVExtraArg{
			CcvAddress: ccvSend.Address.Binding(),
			CcvArgs:    types.TEXT(""),
		})
		allCCVHex = append(allCCVHex, ccvSend.Address.InstanceAddress().Hex())
		out.DisclosedContracts = append(out.DisclosedContracts, ccvSend.DisclosedContracts...)
	}

	if ccipDisc.Executor != nil {
		var execAddr contracts.InstanceAddress
		if rawInstanceAddress, err := contracts.RawInstanceAddressFromString(*ccipDisc.Executor); err == nil {
			execAddr = rawInstanceAddress.InstanceAddress()
		} else {
			execAddr = contracts.HexToInstanceAddress(*ccipDisc.Executor)
		}
		execDisc, err := eds.GetExecutorSendDisclosure(ctx, execCl, outgoing, execAddr, allCCVHex)
		if err != nil {
			return nil, fmt.Errorf("executor send disclosure: %w", err)
		}
		out.ExecutorInput = &ccipsender.ExecutorInput{
			ExecutorCid:          types.CONTRACT_ID(execDisc.ContractId),
			ExecutorExtraContext: execDisc.ChoiceContext,
		}
		out.DisclosedContracts = append(out.DisclosedContracts, execDisc.DisclosedContracts...)
	}

	return out, nil
}

// FindDisclosedContractByContractID returns the EDS disclosed contract for the given ledger contract id, if present.
func FindDisclosedContractByContractID(disclosed []*apiv2.DisclosedContract, contractID string) *apiv2.DisclosedContract {
	for _, dc := range disclosed {
		if dc != nil && dc.GetContractId() == contractID {
			return dc
		}
	}
	return nil
}

func ccvInstanceAddressFromEDS(v string) (contracts.InstanceAddress, error) {
	if rawInstanceAddress, err := contracts.RawInstanceAddressFromString(v); err == nil {
		return rawInstanceAddress.InstanceAddress(), nil
	}
	return contracts.HexToInstanceAddress(v), nil
}
