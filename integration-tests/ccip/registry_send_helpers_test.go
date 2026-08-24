package tests

import (
	"context"
	"fmt"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ratelimiter"
	contractops "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ccip"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

func uploadRegistryDARs(t *testing.T, participants ...canton.Participant) {
	t.Helper()

	dars, err := registry.LoadUtilityDARs()
	require.NoError(t, err, "load registry utility DARs")
	_, err = testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, participants...)
	require.NoError(t, err, "upload registry utility DARs")
}

type registryPoolSendDeps struct {
	Client                    ledger.Client
	CcipClient                ledger.Client
	RegistrarParty            string
	CcipParty                 string
	Bootstrap                 registry.BootstrapResult
	PoolInstanceID            string
	RateLimiterInstanceID     string
	PoolAddress               contracts.RawInstanceAddress
	TokenAdminRegistryAddress contracts.RawInstanceAddress
	TokenAdminRegistryCID     string
	RMNRemoteAddress          contracts.RawInstanceAddress
}

func findContractCIDByInstanceID(
	ctx context.Context,
	client ledger.Client,
	queryParty string,
	template interface{ GetTemplateID() string },
	instanceID string,
) (string, error) {
	tpl := contracts.IdentifierFromBinding(template)
	active, err := testhelpers.ListActiveContractsByTemplateId(ctx, client.ForParty(queryParty), tpl)
	if err != nil {
		return "", err
	}
	for _, ac := range active {
		for _, field := range ac.GetCreatedEvent().GetCreateArguments().GetFields() {
			if field.GetLabel() == "instanceId" && field.GetValue().GetText() == instanceID {
				return ac.GetCreatedEvent().GetContractId(), nil
			}
		}
	}

	return "", fmt.Errorf("contract with instanceId %s not found", instanceID)
}

func contractCIDByInstance(
	t *testing.T,
	ctx context.Context,
	participant canton.Participant,
	party string,
	template interface{ GetTemplateID() string },
	instanceAddr contracts.InstanceAddress,
) string {
	t.Helper()

	cid, err := contractops.FindActiveContractIDByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		[]string{party},
		template.GetTemplateID(),
		instanceAddr,
	)
	require.NoError(t, err)

	return cid
}

func buildRegistryTokenPoolSendDisclosure(
	t *testing.T,
	ctx context.Context,
	registrarParticipant canton.Participant,
	ccipParticipant canton.Participant,
	ccipAPIClient oapiCCIP.ClientWithResponsesInterface,
	deps registryPoolSendDeps,
	hashedInstrumentID contracts.EncodedInstrumentID,
	enableResultContracts bool,
) *edsTesthelpers.TokenPoolSendDisclosure {
	t.Helper()

	poolAddressEDS, err := edsTesthelpers.GetTokenPoolForToken(ctx, ccipAPIClient, hashedInstrumentID)
	require.NoError(t, err)

	outboundRLCID, err := findContractCIDByInstanceID(ctx, deps.Client, deps.RegistrarParty, ratelimiter.RateLimiter{}, deps.RateLimiterInstanceID)
	require.NoError(t, err)

	poolCID := contractCIDByInstance(t, ctx, registrarParticipant, deps.RegistrarParty, burnminttokenpool.BurnMintTokenPool{}, deps.PoolAddress.InstanceAddress())
	rmnCID := contractCIDByInstance(t, ctx, ccipParticipant, deps.CcipParty, core.RMNRemote{}, deps.RMNRemoteAddress.InstanceAddress())

	choiceContext := ccip.RegistryPoolSendExtraContextV1(
		outboundRLCID,
		deps.Bootstrap.AllocationFactory,
		deps.Bootstrap.InstrumentConfiguration,
		enableResultContracts,
	)

	poolDisclosures, err := ccip.DisclosePoolSendContracts(ctx, deps.Client, ccip.PoolSendInput{
		RegistrarParty:          deps.RegistrarParty,
		CcipParty:               deps.CcipParty,
		CcipClient:              deps.CcipClient,
		InstrumentConfiguration: deps.Bootstrap.InstrumentConfiguration,
		AllocationFactory:       deps.Bootstrap.AllocationFactory,
		PoolCID:                 poolCID,
		TokenAdminRegistryCID:   deps.TokenAdminRegistryCID,
		OutboundRateLimiterCID:  outboundRLCID,
		RMNRemoteCID:            rmnCID,
	})
	require.NoError(t, err)

	var disclosedContracts []*apiv2.DisclosedContract
	for _, d := range poolDisclosures.All() {
		if d != nil {
			disclosedContracts = append(disclosedContracts, d)
		}
	}

	return &edsTesthelpers.TokenPoolSendDisclosure{
		ContractId:         poolCID,
		Address:            poolAddressEDS,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       nil,
	}
}

type registryPoolExecuteDeps struct {
	Client                       ledger.Client
	CcipClient                   ledger.Client
	RegistrarParty               string
	CcipParty                    string
	Bootstrap                    registry.BootstrapResult
	PoolInstanceID               string
	DefaultRateLimiterInstanceID string
	CustomRateLimiterInstanceID  string
	PoolAddress                  contracts.RawInstanceAddress
	TokenAdminRegistryCID        string
	RMNRemoteAddress             contracts.RawInstanceAddress
}

func buildRegistryTokenPoolExecuteDisclosure(
	t *testing.T,
	ctx context.Context,
	registrarParticipant canton.Participant,
	ccipParticipant canton.Participant,
	ccipAPIClient oapiCCIP.ClientWithResponsesInterface,
	deps registryPoolExecuteDeps,
	hashedInstrumentID contracts.EncodedInstrumentID,
	customFinality bool,
) *edsTesthelpers.TokenPoolExecuteDisclosure {
	t.Helper()

	poolAddressEDS, err := edsTesthelpers.GetTokenPoolForToken(ctx, ccipAPIClient, hashedInstrumentID)
	require.NoError(t, err)

	rlInstanceID := deps.DefaultRateLimiterInstanceID
	if customFinality {
		rlInstanceID = deps.CustomRateLimiterInstanceID
	}
	inboundRLCID, err := findContractCIDByInstanceID(ctx, deps.Client, deps.RegistrarParty, ratelimiter.RateLimiter{}, rlInstanceID)
	require.NoError(t, err)

	poolCID := contractCIDByInstance(t, ctx, registrarParticipant, deps.RegistrarParty, burnminttokenpool.BurnMintTokenPool{}, deps.PoolAddress.InstanceAddress())
	rmnCID := contractCIDByInstance(t, ctx, ccipParticipant, deps.CcipParty, core.RMNRemote{}, deps.RMNRemoteAddress.InstanceAddress())

	choiceContext := ccip.RegistryPoolSendExtraContextV1(
		inboundRLCID,
		deps.Bootstrap.AllocationFactory,
		deps.Bootstrap.InstrumentConfiguration,
		true,
	)

	poolDisclosures, err := ccip.DisclosePoolExecuteContracts(ctx, deps.Client, ccip.PoolExecuteInput{
		RegistrarParty:          deps.RegistrarParty,
		CcipParty:               deps.CcipParty,
		CcipClient:              deps.CcipClient,
		InstrumentConfiguration: deps.Bootstrap.InstrumentConfiguration,
		AllocationFactory:       deps.Bootstrap.AllocationFactory,
		PoolCID:                 poolCID,
		TokenAdminRegistryCID:   deps.TokenAdminRegistryCID,
		InboundRateLimiterCID:   inboundRLCID,
		RMNRemoteCID:            rmnCID,
	})
	require.NoError(t, err)

	var disclosedContracts []*apiv2.DisclosedContract
	for _, d := range poolDisclosures.All() {
		if d != nil {
			disclosedContracts = append(disclosedContracts, d)
		}
	}

	return &edsTesthelpers.TokenPoolExecuteDisclosure{
		ContractId:         poolCID,
		Address:            poolAddressEDS,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       nil,
	}
}
