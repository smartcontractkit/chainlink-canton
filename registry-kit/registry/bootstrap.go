package registry

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/go-daml/pkg/types"

	registryapp "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/utility/registry_app_v0"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

// BootstrapResult holds CIDs produced by the LocalNet Registry bootstrap (operator=provider=registrar).
type BootstrapResult struct {
	Party                   string
	InstrumentID            string
	OperatorConfiguration   string
	ProviderService         string
	ProviderConfiguration   string
	RegistrarService        string
	AllocationFactory       string
	InstrumentConfiguration string
}

// BootstrapServices bootstraps Registry utility services without the DA operator backend.
func BootstrapServices(ctx context.Context, client ledger.Client, party, instrumentID string) (BootstrapResult, error) {
	opCfgCID, err := createOperatorConfiguration(ctx, client, party)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("create operator configuration: %w", err)
	}
	providerReqCID, err := createProviderServiceRequest(ctx, client, party)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("create provider service request: %w", err)
	}
	providerSvcCID, err := acceptProviderServiceRequest(ctx, client, party, providerReqCID, opCfgCID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("accept provider service request: %w", err)
	}
	providerCfgCID, err := createProviderConfiguration(ctx, client, party, providerSvcCID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("create provider configuration: %w", err)
	}
	registrarReqCID, err := createRegistrarServiceRequest(ctx, client, party)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("create registrar service request: %w", err)
	}
	regSvcCID, allocFactoryCID, err := acceptRegistrarServiceRequest(ctx, client, party, providerSvcCID, registrarReqCID, providerCfgCID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("accept registrar service request: %w", err)
	}
	instCfgCID, err := createInstrumentConfiguration(ctx, client, party, regSvcCID, instrumentID)
	if err != nil {
		return BootstrapResult{}, fmt.Errorf("create instrument configuration: %w", err)
	}

	return BootstrapResult{
		Party:                   party,
		InstrumentID:            instrumentID,
		OperatorConfiguration:   opCfgCID,
		ProviderService:         providerSvcCID,
		ProviderConfiguration:   providerCfgCID,
		RegistrarService:        regSvcCID,
		AllocationFactory:       allocFactoryCID,
		InstrumentConfiguration: instCfgCID,
	}, nil
}

func createOperatorConfiguration(ctx context.Context, client ledger.Client, party string) (string, error) {
	opCfg := registryapp.OperatorConfiguration{
		Operator:             types.PARTY(party),
		ProviderRequirements: nil,
	}
	res, err := client.SubmitCreate(ctx, party, opCfg)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "OperatorConfiguration")
	if !ok {
		return "", fmt.Errorf("OperatorConfiguration not created")
	}

	return cid, nil
}

func createProviderServiceRequest(ctx context.Context, client ledger.Client, party string) (string, error) {
	req := registryapp.ProviderServiceRequest{
		Operator: types.PARTY(party),
		Provider: types.PARTY(party),
	}
	res, err := client.SubmitCreate(ctx, party, req)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "ProviderServiceRequest")
	if !ok {
		return "", fmt.Errorf("ProviderServiceRequest not created")
	}

	return cid, nil
}

func acceptProviderServiceRequest(ctx context.Context, client ledger.Client, party, reqCID, opCfgCID string) (string, error) {
	args := registryapp.ProviderServiceRequestAccept{
		OperatorConfigurationCid:      types.CONTRACT_ID(opCfgCID),
		CredentialCids:                nil,
		AppRewardConfigurationDetails: nil,
	}
	res, err := client.SubmitExercise(ctx, party, registryapp.ProviderServiceRequest{}, reqCID, "ProviderServiceRequest_Accept", args)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "ProviderService")
	if !ok {
		return "", fmt.Errorf("ProviderService not created")
	}

	return cid, nil
}

func createProviderConfiguration(ctx context.Context, client ledger.Client, party, providerSvcCID string) (string, error) {
	args := registryapp.ProviderServiceCreateProviderConfiguration{
		RegistrarRequirements: nil,
		HolderRequirements:    nil,
	}
	res, err := client.SubmitExercise(ctx, party, registryapp.ProviderService{}, providerSvcCID, "ProviderService_CreateProviderConfiguration", args)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "ProviderConfiguration")
	if !ok {
		return "", fmt.Errorf("ProviderConfiguration not created")
	}

	return cid, nil
}

func createRegistrarServiceRequest(ctx context.Context, client ledger.Client, party string) (string, error) {
	createTransferRule := types.BOOL(true)
	createAllocationFactory := types.BOOL(true)
	req := registryapp.RegistrarServiceRequest{
		Operator:                types.PARTY(party),
		Provider:                types.PARTY(party),
		Registrar:               types.PARTY(party),
		CreateTransferRule:      &createTransferRule,
		CreateAllocationFactory: &createAllocationFactory,
	}
	res, err := client.SubmitCreate(ctx, party, req)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "RegistrarServiceRequest")
	if !ok {
		return "", fmt.Errorf("RegistrarServiceRequest not created")
	}

	return cid, nil
}

func acceptRegistrarServiceRequest(ctx context.Context, client ledger.Client, party, providerSvcCID, registrarReqCID, providerCfgCID string) (string, string, error) {
	args := registryapp.ProviderServiceAcceptRegistrarServiceRequest{
		Cid: types.CONTRACT_ID(registrarReqCID),
		Payload: registryapp.RegistrarServiceRequestAccept{
			ProviderConfigurationCid: types.CONTRACT_ID(providerCfgCID),
			CredentialCids:           nil,
		},
	}
	res, err := client.SubmitExercise(ctx, party, registryapp.ProviderService{}, providerSvcCID, "ProviderService_AcceptRegistrarServiceRequest", args)
	if err != nil {
		return "", "", err
	}

	registrarSvcCID, ok := ledger.CreatedContractID(res.GetTransaction(), "RegistrarService")
	if !ok {
		return "", "", fmt.Errorf("RegistrarService not created")
	}
	allocFactoryCID, ok := ledger.CreatedContractID(res.GetTransaction(), "AllocationFactory")
	if !ok {
		return "", "", fmt.Errorf("AllocationFactory not created")
	}

	return registrarSvcCID, allocFactoryCID, nil
}

func createInstrumentConfiguration(ctx context.Context, client ledger.Client, party, registrarSvcCID, instrumentID string) (string, error) {
	args := registryapp.RegistrarServiceCreateInstrumentConfiguration{
		InstrumentId:          types.TEXT(instrumentID),
		AdditionalIdentifiers: nil,
		HolderRequirements:    nil,
		IssuerRequirements:    nil,
	}
	res, err := client.SubmitExercise(ctx, party, registryapp.RegistrarService{}, registrarSvcCID, "RegistrarService_CreateInstrumentConfiguration", args)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "InstrumentConfiguration")
	if !ok {
		return "", fmt.Errorf("InstrumentConfiguration not created")
	}

	return cid, nil
}
