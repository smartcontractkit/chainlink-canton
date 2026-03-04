package devenv

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/util"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfdeployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	offramp2 "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

func convertToDisclosedContract(contract *apiv2.ActiveContract) *apiv2.DisclosedContract {
	if contract == nil {
		return nil
	}

	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}
}

// GetDisclosuresForExecutionFromACS uses the active contract set directly to get all necessary disclosures to execute a message.
func (c *Chain) GetDisclosuresForExecutionFromACS(ctx context.Context, verifiers []contracts.InstanceAddress) (
	// List of disclosed contracts to be used during the execute call
	[]*apiv2.DisclosedContract,
	// The choiceContext value
	*apiv2.Value,
	// The CCV ContractIDs, in the same order as the input ccvs
	[]*apiv2.Value_ContractId,
	error,
) {
	// Use only a single participant for now
	participant := c.chain.Participants[0]

	var disclosedContracts []*apiv2.DisclosedContract

	// OffRamp
	offRampRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(offramp.ContractType),
			offramp.Version,
			"",
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get offramp address ref: %w", err)
	}
	offRampAddress := contracts.HexToInstanceAddress(offRampRef.Address)
	activeOffRamp, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, offramp2.OffRamp{}.GetTemplateID(), offRampAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get offramp contract ID: %w", err)
	}
	disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeOffRamp))
	c.logger.Debug().Str("InstanceAddress", offRampAddress.String()).Str("ContractId", activeOffRamp.GetCreatedEvent().GetContractId()).Msg("Resolved OffRamp contract")

	// GlobalConfig
	globalConfigRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(global_config.ContractType),
			global_config.Version,
			"",
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get global config address ref: %w", err)
	}
	globalConfigAddress := contracts.HexToInstanceAddress(globalConfigRef.Address)
	activeGlobalConfig, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, common.GlobalConfig{}.GetTemplateID(), globalConfigAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get global config contract ID: %w", err)
	}
	disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeGlobalConfig))
	c.logger.Debug().Str("InstanceAddress", globalConfigAddress.String()).Str("ContractId", activeGlobalConfig.GetCreatedEvent().GetContractId()).Msg("Resolved GlobalConfig contract")

	// Token Admin Registry
	tokenAdminRegistryRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(token_admin_registry.ContractType),
			token_admin_registry.Version,
			"",
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get token admin registry address ref: %w", err)
	}
	tokenAdminRegistryAddress := contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address)
	activeTokenAdminRegistry, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, tokenadminregistry.TokenAdminRegistry{}.GetTemplateID(), tokenAdminRegistryAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get token admin registry contract ID: %w", err)
	}
	disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeTokenAdminRegistry))
	c.logger.Debug().Str("InstanceAddress", tokenAdminRegistryAddress.String()).Str("ContractId", activeTokenAdminRegistry.GetCreatedEvent().GetContractId()).Msg("Resolved TokenAdminRegistry contract")

	// RMN Remote
	rmnRemoteRef, err := c.e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			c.chainDetails.ChainSelector,
			datastore.ContractType(rmn_remote.ContractType),
			rmn_remote.Version,
			"",
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get rmn remote address ref: %w", err)
	}
	rmnRemoteAddress := contracts.HexToInstanceAddress(rmnRemoteRef.Address)
	activeRMNRemote, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, rmn.RMNRemote{}.GetTemplateID(), rmnRemoteAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get rmn remote contract ID: %w", err)
	}
	disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeRMNRemote))
	c.logger.Debug().Str("InstanceAddress", rmnRemoteAddress.String()).Str("ContractId", activeRMNRemote.GetCreatedEvent().GetContractId()).Msg("Resolved RMNRemote contract")

	// Verifiers
	ccvContractIDs := make([]*apiv2.Value_ContractId, len(verifiers))
	for i, verifierAddr := range verifiers {
		activeVerifier, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, ccvs.CommitteeVerifier{}.GetTemplateID(), verifierAddr)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get committee verifier contract ID for address %s: %w", verifierAddr.String(), err)
		}
		disclosedContracts = append(disclosedContracts, convertToDisclosedContract(activeVerifier))
		ccvContractIDs[i] = &apiv2.Value_ContractId{
			ContractId: activeVerifier.GetCreatedEvent().GetContractId(),
		}
		c.logger.Debug().Str("InstanceAddress", verifierAddr.String()).Str("ContractId", activeVerifier.GetCreatedEvent().GetContractId()).Msg("Resolved CCV contract")
	}

	choiceContext := map[string]any{
		"values": map[string]any{
			"off-ramp": map[string]any{
				"tag":   "AV_ContractId",
				"value": activeOffRamp.GetCreatedEvent().GetContractId(),
			},
			"global-config": map[string]any{
				"tag":   "AV_ContractId",
				"value": activeGlobalConfig.GetCreatedEvent().GetContractId(),
			},
			"token-admin-registry": map[string]any{
				"tag":   "AV_ContractId",
				"value": activeTokenAdminRegistry.GetCreatedEvent().GetContractId(),
			},
			"rmn-remote": map[string]any{
				"tag":   "AV_ContractId",
				"value": activeRMNRemote.GetCreatedEvent().GetContractId(),
			},
		},
	}
	choiceContextValue, err := ChoiceContextFromData(choiceContext)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create choice context: %w", err)
	}

	return disclosedContracts, choiceContextValue, ccvContractIDs, nil
}

// GetDisclosuresForExecution returns all the necessary disclosed contracts to execute a message on Canton using the EDS API.
func (c *Chain) GetDisclosuresForExecution(ctx context.Context, verifiers []contracts.InstanceAddress) (
	// List of disclosed contracts to be used during the execute call
	[]*apiv2.DisclosedContract,
	// The choiceContext value
	*apiv2.Value,
	// The CCV ContractIDs, in the same order as the input ccvs
	[]*apiv2.Value_ContractId,
	error,
) {
	edsCfg, err := deployment.GetEDSConfig(c.e.DataStore)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get eds config: %w", err)
	}

	// The container exposes the API on the same port as internally
	url := fmt.Sprintf("http://localhost:%d", edsCfg.Server.Port)
	edsClient, err := edsv1.NewClientWithResponses(url)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create eds client: %w", err)
	}

	request := edsv1.CCIPExecuteRequest{
		Ccvs:      make([]string, len(verifiers)),
		MessageID: "", // not used (yet)
	}
	for i, verifier := range verifiers {
		request.Ccvs[i] = verifier.String()
	}

	resp, err := edsClient.CcipExecuteWithResponse(ctx, request)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get disclosures from EDS: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("failed to get disclosures from EDS: %s", resp.Status())
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		disclosedContract, err := disclosedContractToProto(contract)
		if err != nil {
			return nil, nil, nil, err
		}

		disclosedContracts = append(disclosedContracts, disclosedContract)
	}
	choiceContext, err := ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	// Handle CCVs
	ccvContractIDs := make([]*apiv2.Value_ContractId, len(verifiers))
	for i, ccv := range verifiers {
		// Check if the API returned an explicit disclosure for the requested CCV
		disclosure, ok := resp.JSON200.Ccvs[ccv.String()]
		if !ok || disclosure.DisclosedContract == nil {
			return nil, nil, nil, fmt.Errorf("failed to get disclosure for ccv: %s", ccv.String())
		}
		ccvContractIDs[i] = &apiv2.Value_ContractId{
			ContractId: disclosure.DisclosedContract.ContractId,
		}
		// Add the CCV's explicit disclosure to disclosedContracts
		disclosedContract, err := disclosedContractToProto(*disclosure.DisclosedContract)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to convert disclosed contract for ccv %s: %w", ccv.String(), err)
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	return disclosedContracts, choiceContext, ccvContractIDs, nil
}

func TemplateIdFromString(s string) (*apiv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}

	return &apiv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func disclosedContractToProto(contract edsv1.DisclosedContract) (*apiv2.DisclosedContract, error) {
	id, err := TemplateIdFromString(contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       id,
		ContractId:       contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   contract.SynchronizerId,
	}, nil
}

const (
	DefaultEDSName  = "canton-eds"
	DefaultEDSImage = "canton-eds:latest"
)

func startEDS(ctx context.Context, cfg *edsConfig.Config) (testcontainers.Container, error) {
	configToml, err := toml.Marshal(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EDS config: %w", err)
	}
	req := testcontainers.ContainerRequest{
		Image:    DefaultEDSImage,
		Name:     DefaultEDSName,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {DefaultEDSName},
		},
		ExposedPorts: []string{
			fmt.Sprintf("%d/tcp", cfg.Server.Port),
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port(fmt.Sprintf("%d/tcp", cfg.Server.Port)): []nat.PortBinding{
					{HostPort: ""}, // Docker assigns a random free host port.
				},
			}
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            bytes.NewReader(configToml),
				ContainerFilePath: "/app/config.toml",
				FileMode:          0755,
			},
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return c, nil
}

type launcher struct{}

var _ ccv.Launcher = &launcher{}

func NewLauncher() *launcher {
	return &launcher{}
}

type input struct {
	EDSServerConfig edsConfig.ServerConfig `toml:"eds_server_config"`
}

func (l *launcher) Launch(
	ctx context.Context,
	env *cldfdeployment.Environment,
	chains []*blockchain.Output,
	definition *ccv.GenericServiceDefinition,
) (output util.OpaqueConfig, err error) {
	var cantonOutput *blockchain.Output
	for _, chain := range chains {
		if chain.Type == chain_selectors.FamilyCanton {
			cantonOutput = chain
			break
		}
	}
	if cantonOutput == nil {
		return nil, fmt.Errorf("canton chain not found")
	}
	chainDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(cantonOutput.ChainID, chain_selectors.FamilyCanton)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for canton chain: %w", err)
	}

	// parse input from definition
	in, err := util.OpaqueToConcreteStrict[input](definition.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}
	out, err := cantonChangesets.GenerateEDSConfig{}.Apply(*env, cantonChangesets.CantonCSDeps[edsConfig.ServerConfig]{
		ChainSelector: chainDetails.ChainSelector,
		Participant:   0,
		Config: edsConfig.ServerConfig{
			Host: in.EDSServerConfig.Host,
			Port: in.EDSServerConfig.Port,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate EDS config for selector %d: %w", chainDetails.ChainSelector, err)
	}

	edsCfg, err := deployment.GetEDSConfig(out.DataStore.Seal())
	if err != nil {
		return nil, fmt.Errorf("failed to get EDS config: %w", err)
	}

	// start the EDS
	_, err = startEDS(ctx, edsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to start EDS: %w", err)
	}

	// return the EDS config
	return util.OpaqueConfig{
		"eds_config": edsCfg,
	}, nil
}
