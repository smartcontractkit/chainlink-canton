package ccip

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	burnminttokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	ccipcore "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// PoolReleaseDisclosures holds disclosed contracts required for ReleaseFromTicket.
type PoolReleaseDisclosures struct {
	InstrumentConfiguration *apiv2.DisclosedContract
	AllocationFactory       *apiv2.DisclosedContract
	BurnMintTokenPool       *apiv2.DisclosedContract
	TokenAdminRegistry      *apiv2.DisclosedContract
	RateLimiter             *apiv2.DisclosedContract
	RMNRemote               *apiv2.DisclosedContract
	TokenReceiveTicket      *apiv2.DisclosedContract
}

func (d PoolReleaseDisclosures) All() []*apiv2.DisclosedContract {
	return []*apiv2.DisclosedContract{
		d.InstrumentConfiguration,
		d.AllocationFactory,
		d.BurnMintTokenPool,
		d.TokenAdminRegistry,
		d.RateLimiter,
		d.RMNRemote,
		d.TokenReceiveTicket,
	}
}

// PoolReleaseInput identifies contracts needed for pool release disclosures.
type PoolReleaseInput struct {
	RegistrarParty          string
	// CcipParty discloses CCIP-owned contracts (TAR, RMN) when they differ from the registrar.
	CcipParty               string
	InstrumentConfiguration string
	AllocationFactory       string
	PoolCID                 string
	TokenAdminRegistryCID   string
	InboundRateLimiterCID   string
	RMNRemoteCID            string
	TokenReceiveTicketCID   string
}

func registryPoolBurnMintExtraContext(
	rateLimiterCID, allocationFactoryCID, instrumentConfigCID string,
	enableResultContracts bool,
) splice_api_token_metadata_v1.ChoiceContext {
	nestedValues := registry.MintChoiceContext(instrumentConfigCID, enableResultContracts).Values

	return splice_api_token_metadata_v1.ChoiceContext{
		Values: map[string]splice_api_token_metadata_v1.AnyValue{
			string(ccipcore.RateLimiterKey): {
				AVContractId: new(types.CONTRACT_ID(rateLimiterCID)),
			},
			string(burnminttokenpool.BurnMintFactoryContextKey): {
				AVContractId: new(types.CONTRACT_ID(allocationFactoryCID)),
			},
			string(burnminttokenpool.BurnMintFactoryExtraArgsContextValuesContextKey): {
				AVMap: &nestedValues,
			},
		},
	}
}

// RegistryPoolExtraContext builds pool extraContext for ReleaseFromTicket with AllocationFactory
// as burn-mint-factory and nested Registry mint choice-context values.
func RegistryPoolExtraContext(inboundRateLimiterCID, allocationFactoryCID, instrumentConfigCID string) splice_api_token_metadata_v1.ChoiceContext {
	return registryPoolBurnMintExtraContext(inboundRateLimiterCID, allocationFactoryCID, instrumentConfigCID, true)
}

// RegistryPoolSendExtraContext builds pool extraContext for LockOrBurn with AllocationFactory
// as burn-mint-factory and nested Registry burn choice-context values.
func RegistryPoolSendExtraContext(outboundRateLimiterCID, allocationFactoryCID, instrumentConfigCID string, enableResultContracts bool) splice_api_token_metadata_v1.ChoiceContext {
	return registryPoolBurnMintExtraContext(outboundRateLimiterCID, allocationFactoryCID, instrumentConfigCID, enableResultContracts)
}

// PoolSendDisclosures holds disclosed contracts required for LockOrBurn.
type PoolSendDisclosures struct {
	InstrumentConfiguration *apiv2.DisclosedContract
	AllocationFactory       *apiv2.DisclosedContract
	BurnMintTokenPool       *apiv2.DisclosedContract
	TokenAdminRegistry      *apiv2.DisclosedContract
	RateLimiter             *apiv2.DisclosedContract
	RMNRemote               *apiv2.DisclosedContract
	SendingMessage          *apiv2.DisclosedContract
}

func (d PoolSendDisclosures) All() []*apiv2.DisclosedContract {
	out := []*apiv2.DisclosedContract{
		d.InstrumentConfiguration,
		d.AllocationFactory,
		d.BurnMintTokenPool,
		d.TokenAdminRegistry,
		d.RateLimiter,
		d.RMNRemote,
	}
	if d.SendingMessage != nil {
		out = append(out, d.SendingMessage)
	}

	return out
}

// PoolSendInput identifies contracts needed for pool send disclosures.
type PoolSendInput struct {
	RegistrarParty            string
	CcipParty                 string
	InstrumentConfiguration   string
	AllocationFactory         string
	PoolCID                   string
	TokenAdminRegistryCID     string
	OutboundRateLimiterCID    string
	RMNRemoteCID              string
	SendingMessageCID         string
	// CcipClient optionally supplies the ledger client for CCIP-party disclosures when it
	// differs from the client used for registrar-party disclosures (multi-participant CTF).
	CcipClient ledger.Client
}

// DisclosePoolSendContracts fetches stakeholder-party disclosures for the LockOrBurn path.
func DisclosePoolSendContracts(ctx context.Context, client ledger.Client, input PoolSendInput) (PoolSendDisclosures, error) {
	registrar := input.RegistrarParty
	ccipClient := input.CcipClient
	if ccipClient == nil {
		ccipClient = client
	}

	instDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.InstrumentConfiguration)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose InstrumentConfiguration: %w", err)
	}
	allocDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.AllocationFactory)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose AllocationFactory: %w", err)
	}
	poolDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.PoolCID)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose BurnMintTokenPool: %w", err)
	}
	ccipParty := input.CcipParty
	if ccipParty == "" {
		ccipParty = registrar
	}
	tarDisclosed, err := registry.DiscloseByID(ctx, ccipClient, ccipParty, input.TokenAdminRegistryCID)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose TokenAdminRegistry: %w", err)
	}
	rlDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.OutboundRateLimiterCID)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose RateLimiter: %w", err)
	}
	rmnDisclosed, err := registry.DiscloseByID(ctx, ccipClient, ccipParty, input.RMNRemoteCID)
	if err != nil {
		return PoolSendDisclosures{}, fmt.Errorf("disclose RMNRemote: %w", err)
	}
	var msgDisclosed *apiv2.DisclosedContract
	if input.SendingMessageCID != "" {
		msgDisclosed, err = registry.DiscloseByID(ctx, ccipClient, ccipParty, input.SendingMessageCID)
		if err != nil {
			return PoolSendDisclosures{}, fmt.Errorf("disclose SendingMessageV1: %w", err)
		}
	}

	return PoolSendDisclosures{
		InstrumentConfiguration: instDisclosed,
		AllocationFactory:       allocDisclosed,
		BurnMintTokenPool:       poolDisclosed,
		TokenAdminRegistry:      tarDisclosed,
		RateLimiter:             rlDisclosed,
		RMNRemote:               rmnDisclosed,
		SendingMessage:          msgDisclosed,
	}, nil
}

// DisclosePoolReleaseContracts fetches stakeholder-party disclosures for the release path.
func DisclosePoolReleaseContracts(ctx context.Context, client ledger.Client, input PoolReleaseInput) (PoolReleaseDisclosures, error) {
	registrar := input.RegistrarParty

	instDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.InstrumentConfiguration)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose InstrumentConfiguration: %w", err)
	}
	allocDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.AllocationFactory)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose AllocationFactory: %w", err)
	}
	poolDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.PoolCID)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose BurnMintTokenPool: %w", err)
	}
	ccipParty := input.CcipParty
	if ccipParty == "" {
		ccipParty = registrar
	}
	tarDisclosed, err := registry.DiscloseByID(ctx, client, ccipParty, input.TokenAdminRegistryCID)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose TokenAdminRegistry: %w", err)
	}
	rlDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.InboundRateLimiterCID)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose RateLimiter: %w", err)
	}
	rmnDisclosed, err := registry.DiscloseByID(ctx, client, ccipParty, input.RMNRemoteCID)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose RMNRemote: %w", err)
	}
	ticketDisclosed, err := registry.DiscloseByID(ctx, client, ccipParty, input.TokenReceiveTicketCID)
	if err != nil {
		return PoolReleaseDisclosures{}, fmt.Errorf("disclose TokenReceiveTicket: %w", err)
	}

	return PoolReleaseDisclosures{
		InstrumentConfiguration: instDisclosed,
		AllocationFactory:       allocDisclosed,
		BurnMintTokenPool:       poolDisclosed,
		TokenAdminRegistry:      tarDisclosed,
		RateLimiter:             rlDisclosed,
		RMNRemote:               rmnDisclosed,
		TokenReceiveTicket:      ticketDisclosed,
	}, nil
}

// PoolExecuteDisclosures holds disclosed contracts required for CCIPReceiver.Execute mint path.
type PoolExecuteDisclosures struct {
	InstrumentConfiguration *apiv2.DisclosedContract
	AllocationFactory       *apiv2.DisclosedContract
	BurnMintTokenPool       *apiv2.DisclosedContract
	TokenAdminRegistry      *apiv2.DisclosedContract
	RateLimiter             *apiv2.DisclosedContract
	RMNRemote               *apiv2.DisclosedContract
}

func (d PoolExecuteDisclosures) All() []*apiv2.DisclosedContract {
	return []*apiv2.DisclosedContract{
		d.InstrumentConfiguration,
		d.AllocationFactory,
		d.BurnMintTokenPool,
		d.TokenAdminRegistry,
		d.RateLimiter,
		d.RMNRemote,
	}
}

// PoolExecuteInput identifies contracts needed for pool execute disclosures.
type PoolExecuteInput struct {
	RegistrarParty          string
	CcipParty               string
	CcipClient              ledger.Client
	InstrumentConfiguration string
	AllocationFactory       string
	PoolCID                 string
	TokenAdminRegistryCID   string
	InboundRateLimiterCID   string
	RMNRemoteCID            string
}

// DisclosePoolExecuteContracts fetches stakeholder-party disclosures for the execute mint path.
func DisclosePoolExecuteContracts(ctx context.Context, client ledger.Client, input PoolExecuteInput) (PoolExecuteDisclosures, error) {
	registrar := input.RegistrarParty
	ccipClient := input.CcipClient
	if ccipClient == nil {
		ccipClient = client
	}

	instDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.InstrumentConfiguration)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose InstrumentConfiguration: %w", err)
	}
	allocDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.AllocationFactory)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose AllocationFactory: %w", err)
	}
	poolDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.PoolCID)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose BurnMintTokenPool: %w", err)
	}
	ccipParty := input.CcipParty
	if ccipParty == "" {
		ccipParty = registrar
	}
	tarDisclosed, err := registry.DiscloseByID(ctx, ccipClient, ccipParty, input.TokenAdminRegistryCID)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose TokenAdminRegistry: %w", err)
	}
	rlDisclosed, err := registry.DiscloseByID(ctx, client, registrar, input.InboundRateLimiterCID)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose RateLimiter: %w", err)
	}
	rmnDisclosed, err := registry.DiscloseByID(ctx, ccipClient, ccipParty, input.RMNRemoteCID)
	if err != nil {
		return PoolExecuteDisclosures{}, fmt.Errorf("disclose RMNRemote: %w", err)
	}

	return PoolExecuteDisclosures{
		InstrumentConfiguration: instDisclosed,
		AllocationFactory:       allocDisclosed,
		BurnMintTokenPool:       poolDisclosed,
		TokenAdminRegistry:      tarDisclosed,
		RateLimiter:             rlDisclosed,
		RMNRemote:               rmnDisclosed,
	}, nil
}
