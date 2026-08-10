package ccip

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
)

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
	RegistrarParty          string
	CcipParty               string
	InstrumentConfiguration string
	AllocationFactory       string
	PoolCID                 string
	TokenAdminRegistryCID   string
	OutboundRateLimiterCID  string
	RMNRemoteCID            string
	SendingMessageCID       string
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
			return PoolSendDisclosures{}, fmt.Errorf("disclose SendingMessage: %w", err)
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
