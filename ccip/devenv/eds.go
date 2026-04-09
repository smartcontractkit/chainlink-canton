package devenv

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/exp/maps"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/util"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cldfdeployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

func GetCCIPSendDisclosures(
	ctx context.Context,
	edsClient *edsv1.ClientWithResponses,
	ccvs []contracts.InstanceAddress,
) (*SendDisclosures, error) {
	// Add the requested CCVs to the API request
	request := edsv1.CCIPSendRequest{
		Ccvs: make([]string, len(ccvs)),
	}
	for i, ccv := range ccvs {
		request.Ccvs[i] = ccv.String()
	}

	resp, err := edsClient.CcipSendWithResponse(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling CCIPSend: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// The top-level choice context has the following disclosures:
	// - OnRamp
	// - GlobalConfig
	// - TokenAdminRegistry
	// - RMNRemote
	// - FeeQuoter
	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		disclosedContract, err := disclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	// Build sendContext from the return choice context data
	// TODO: ledger.RecordToStruct isn't working for some reason on this.
	choiceCtxValues, ok := resp.JSON200.ChoiceContext.ChoiceContextData["values"]
	if !ok {
		return nil, fmt.Errorf("key 'values' not found in choice context")
	}
	choiceCtxValuesMap, ok := choiceCtxValues.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("values is not a map[string]any")
	}
	values := types.TEXTMAP{}
	for contractKey, entry := range choiceCtxValuesMap {
		entryMap, ok := entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("entry is not a map[string]any")
		}
		value := entryMap["value"].(string)
		valueContractID := types.CONTRACT_ID(value)
		values[contractKey] = common.AnyValue{AVContractId: &valueContractID}
	}
	sendContext := common.CCIPContext{Values: values}

	// Handle CCVs
	ccvContractIDs := make([]*apiv2.Value_ContractId, len(ccvs))
	for i, ccv := range ccvs {
		// Check if the API returned an explicit disclosure for the requested CCV
		disclosure, ok := resp.JSON200.Ccvs[ccv.String()]
		if !ok || disclosure.DisclosedContract == nil {
			return nil, fmt.Errorf("failed to get disclosure for ccv: %s", ccv.String())
		}
		ccvContractIDs[i] = &apiv2.Value_ContractId{
			ContractId: disclosure.DisclosedContract.ContractId,
		}
		// Add the CCV's explicit disclosure to disclosedContracts
		disclosedContract, err := disclosedContractToProto(*disclosure.DisclosedContract)
		if err != nil {
			return nil, fmt.Errorf("failed to convert disclosed contract for ccv %s: %w", ccv.String(), err)
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	// Handle executor
	var executorContractID *apiv2.Value_ContractId
	if len(resp.JSON200.DefaultExecutor.ContractId) > 0 {
		executorContractID = &apiv2.Value_ContractId{
			ContractId: resp.JSON200.DefaultExecutor.ContractId,
		}
		// don't need to add to disclosedContracts, its already part of the choice context
		// disclosures.
	}

	return &SendDisclosures{
		DisclosedContracts: disclosedContracts,
		SendContext:        sendContext,
		CCVContractIDs:     ccvContractIDs,
		ExecutorContractID: executorContractID,
	}, nil
}

type SendDisclosures struct {
	DisclosedContracts []*apiv2.DisclosedContract
	SendContext        common.CCIPContext
	CCVContractIDs     []*apiv2.Value_ContractId
	ExecutorContractID *apiv2.Value_ContractId
}

func (c *Chain) GetDisclosuresForSend(ctx context.Context, ccvs []contracts.InstanceAddress) (*SendDisclosures, error) {
	// Get the EDS output from generic services output, will contain the EDS API URL
	edsOut, err := util.OpaqueToConcreteStrict[output](c.cfg.GenericServices[c.ChainSelector()].Output)
	if err != nil {
		return nil, fmt.Errorf("failed to get generic services output for chain %d: %w", c.ChainSelector(), err)
	}
	edsClient, err := edsv1.NewClientWithResponses(edsOut.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds client: %w", err)
	}

	return GetCCIPSendDisclosures(ctx, edsClient, ccvs)
}

func GetDisclosedContractByTemplateId(ctx context.Context, participant canton.Participant, templateId *apiv2.Identifier) (*apiv2.DisclosedContract, error) {
	activeContracts, err := ListActiveContractsByTemplateId(ctx, participant, templateId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active contracts with templateId %v found on participant %vfor party %s", templateId, participant.Name, participant.PartyID)
	}

	contract := activeContracts[len(activeContracts)-1]

	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}

func GetCurrentOffset(ctx context.Context, stateService apiv2.StateServiceClient) (int64, error) {
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger end: %w", err)
	}

	return ledgerEndResp.GetOffset(), nil
}

func ListActiveContractsByTemplateId(ctx context.Context, participant canton.Participant, templateId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant.LedgerServices.State)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
								TemplateId:              templateId,
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts for party %v: %w", participant.PartyID, err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			activeContracts = append(activeContracts, c.ActiveContract)
		}
	}
	slices.SortFunc(activeContracts, func(a, b *apiv2.ActiveContract) int {
		return a.GetCreatedEvent().GetCreatedAt().AsTime().Compare(b.GetCreatedEvent().GetCreatedAt().AsTime())
	})

	return activeContracts, nil
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
	// Get the EDS output from generic services output, will contain the EDS API URL
	edsOut, err := util.OpaqueToConcreteStrict[output](c.cfg.GenericServices[c.ChainSelector()].Output)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get generic services output for chain %d: %w", c.ChainSelector(), err)
	}
	edsClient, err := edsv1.NewClientWithResponses(edsOut.EDSURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create eds client: %w", err)
	}

	// Add the verifier addresses to the request - required for EDS to return explicit disclosures for them
	request := edsv1.CCIPExecuteRequest{
		Ccvs:           make([]string, len(verifiers)),
		EncodedMessage: "", // not used (yet)
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

	// Handle Verifiers
	ccvContractIDs := make([]*apiv2.Value_ContractId, len(verifiers))
	for i, verifier := range verifiers {
		// Check if the API returned an explicit disclosure for the requested verifier
		disclosure, ok := resp.JSON200.Ccvs[verifier.String()]
		if !ok || disclosure.DisclosedContract == nil {
			return nil, nil, nil, fmt.Errorf("failed to get disclosure for verifier: %s", verifier.String())
		}
		ccvContractIDs[i] = &apiv2.Value_ContractId{
			ContractId: disclosure.DisclosedContract.ContractId,
		}
		// Add the CCV's explicit disclosure to disclosedContracts
		disclosedContract, err := disclosedContractToProto(*disclosure.DisclosedContract)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to convert disclosed contract for verifier %s: %w", verifier.String(), err)
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
		WaitingFor: wait.ForLog("Backfill complete"), // TODO: properly implement readiness endpoint
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
	EDSConfig edsConfig.Config `toml:"eds_config"`
}

type output struct {
	EDSConfig edsConfig.Config `toml:"eds_config"`
	EDSURL    string           `toml:"eds_url"`
}

func (l *launcher) Launch(
	ctx context.Context,
	env *cldfdeployment.Environment,
	chains []*blockchain.Output,
	definition *ccv.GenericServiceDefinition,
) (util.OpaqueConfig, error) {
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
	out, err := cantonChangesets.GenerateEDSConfig{}.Apply(*env, cantonChangesets.CantonCSDeps[edsConfig.Config]{
		ChainSelector: chainDetails.ChainSelector,
		Participant:   0,
		Config:        in.EDSConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate EDS config for selector %d: %w", chainDetails.ChainSelector, err)
	}

	edsCfg, err := deployment.GetEDSConfig(out.DataStore.Seal())
	if err != nil {
		return nil, fmt.Errorf("failed to get EDS config: %w", err)
	}

	// start the EDS
	edsContainer, err := startEDS(ctx, edsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to start EDS: %w", err)
	}
	host, err := edsContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get EDS container host: %w", err)
	}
	edsMappedPort, err := edsContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", edsCfg.Server.Port)))
	if err != nil {
		pm, _ := edsContainer.Ports(ctx) // Add all existing ports to error
		return nil, fmt.Errorf("failed to get EDS container mapped port (ports: %v): %w", maps.Keys(pm), err)
	}

	// return the EDS config
	// Add the config and container URL to the output
	opaqueConfig, err := util.ConcreteToOpaque(output{
		EDSConfig: *edsCfg,
		EDSURL:    fmt.Sprintf("http://%s:%s", host, edsMappedPort.Port()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert output to opaque config: %w", err)
	}

	return opaqueConfig, nil
}
