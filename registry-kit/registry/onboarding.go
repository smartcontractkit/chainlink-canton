package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/go-daml/pkg/types"

	registryapp "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_app_v0"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

// OnboardingParties lists Registry role parties for devnet onboarding.
type OnboardingParties struct {
	Operator  string
	Provider  string
	Registrar string
}

// RequestProviderService creates a ProviderServiceRequest (DA operator must accept).
func RequestProviderService(ctx context.Context, client ledger.Client, parties OnboardingParties) (string, error) {
	req := registryapp.ProviderServiceRequest{
		Operator: types.PARTY(parties.Operator),
		Provider: types.PARTY(parties.Provider),
	}
	res, err := client.SubmitCreate(ctx, parties.Provider, req)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "ProviderServiceRequest")
	if !ok {
		return "", fmt.Errorf("ProviderServiceRequest not created")
	}

	return cid, nil
}

// WaitForProviderService polls ACS until ProviderService exists for the provider party.
func WaitForProviderService(ctx context.Context, client ledger.Client, providerParty string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for {
		cid, err := FindFirstContractByEntity(ctx, client, providerParty, registryapp.ProviderService{}, "ProviderService")
		if err != nil {
			return "", err
		}
		if cid != "" {
			return cid, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timed out waiting for ProviderService for %s", providerParty)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

// OnboardRegistrarResult holds CIDs from provider-led registrar onboarding.
type OnboardRegistrarResult struct {
	ProviderConfigurationCID   string
	RegistrarServiceRequestCID string
	RegistrarServiceCID        string
	AllocationFactoryCID       string
	TransferRuleCID            string
}

// OnboardRegistrar runs ProviderConfiguration → RegistrarServiceRequest → AcceptRegistrarServiceRequest.
func OnboardRegistrar(ctx context.Context, client ledger.Client, parties OnboardingParties, providerServiceCID string) (OnboardRegistrarResult, error) {
	providerCfgCID, err := createDevnetProviderConfiguration(ctx, client, parties.Provider, providerServiceCID)
	if err != nil {
		return OnboardRegistrarResult{}, err
	}

	registrarReqCID, err := createRegistrarServiceRequestDevnet(ctx, client, parties)
	if err != nil {
		return OnboardRegistrarResult{}, err
	}

	regSvcCID, allocCID, transferRuleCID, err := acceptRegistrarServiceRequestDevnet(ctx, client, parties.Provider, providerServiceCID, registrarReqCID, providerCfgCID)
	if err != nil {
		return OnboardRegistrarResult{}, err
	}

	return OnboardRegistrarResult{
		ProviderConfigurationCID:   providerCfgCID,
		RegistrarServiceRequestCID: registrarReqCID,
		RegistrarServiceCID:        regSvcCID,
		AllocationFactoryCID:       allocCID,
		TransferRuleCID:            transferRuleCID,
	}, nil
}

func createDevnetProviderConfiguration(ctx context.Context, client ledger.Client, providerParty, providerServiceCID string) (string, error) {
	args := registryapp.ProviderServiceCreateProviderConfiguration{
		RegistrarRequirements: nil,
		HolderRequirements:    nil,
	}
	res, err := client.SubmitExercise(ctx, providerParty, registryapp.ProviderService{}, providerServiceCID, "ProviderService_CreateProviderConfiguration", args)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "ProviderConfiguration")
	if !ok {
		return "", fmt.Errorf("ProviderConfiguration not created")
	}

	return cid, nil
}

func createRegistrarServiceRequestDevnet(ctx context.Context, client ledger.Client, parties OnboardingParties) (string, error) {
	createTransferRule := types.BOOL(true)
	createAllocationFactory := types.BOOL(true)
	req := registryapp.RegistrarServiceRequest{
		Operator:                types.PARTY(parties.Operator),
		Provider:                types.PARTY(parties.Provider),
		Registrar:               types.PARTY(parties.Registrar),
		CreateTransferRule:      &createTransferRule,
		CreateAllocationFactory: &createAllocationFactory,
	}
	res, err := client.SubmitCreate(ctx, parties.Provider, req)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "RegistrarServiceRequest")
	if !ok {
		return "", fmt.Errorf("RegistrarServiceRequest not created")
	}

	return cid, nil
}

func acceptRegistrarServiceRequestDevnet(ctx context.Context, client ledger.Client, providerParty, providerServiceCID, registrarReqCID, providerCfgCID string) (string, string, string, error) {
	args := registryapp.ProviderServiceAcceptRegistrarServiceRequest{
		Cid: types.CONTRACT_ID(registrarReqCID),
		Payload: registryapp.RegistrarServiceRequestAccept{
			ProviderConfigurationCid: types.CONTRACT_ID(providerCfgCID),
			CredentialCids:           nil,
		},
	}
	res, err := client.SubmitExercise(ctx, providerParty, registryapp.ProviderService{}, providerServiceCID, "ProviderService_AcceptRegistrarServiceRequest", args)
	if err != nil {
		return "", "", "", err
	}

	registrarSvcCID, ok := ledger.CreatedContractID(res.GetTransaction(), "RegistrarService")
	if !ok {
		return "", "", "", fmt.Errorf("RegistrarService not created")
	}
	allocCID, ok := ledger.CreatedContractID(res.GetTransaction(), "AllocationFactory")
	if !ok {
		return "", "", "", fmt.Errorf("AllocationFactory not created")
	}
	transferRuleCID, _ := ledger.CreatedContractID(res.GetTransaction(), "TransferRule")

	return registrarSvcCID, allocCID, transferRuleCID, nil
}

// CreateInstrumentConfiguration exercises RegistrarService_CreateInstrumentConfiguration.
func CreateInstrumentConfiguration(ctx context.Context, client ledger.Client, registrarParty, registrarServiceCID, instrumentID string) (string, error) {
	args := registryapp.RegistrarServiceCreateInstrumentConfiguration{
		InstrumentId:          types.TEXT(instrumentID),
		AdditionalIdentifiers: nil,
		HolderRequirements:    nil,
		IssuerRequirements:    nil,
	}
	res, err := client.SubmitExercise(ctx, registrarParty, registryapp.RegistrarService{}, registrarServiceCID, "RegistrarService_CreateInstrumentConfiguration", args)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "InstrumentConfiguration")
	if !ok {
		return "", fmt.Errorf("InstrumentConfiguration not created")
	}

	return cid, nil
}
