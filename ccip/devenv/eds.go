package devenv

import (
	"bytes"
	"context"
	"fmt"
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
	cldfdeployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func (c *Chain) GetDisclosuresForSend(ctx context.Context, ccvs []contracts.InstanceAddress) (*testhelpers.SendDisclosures, error) {
	// Get the EDS output from generic services output, will contain the EDS API URL
	edsOut, err := util.OpaqueToConcreteStrict[output](c.cfg.GenericServices[c.ChainSelector()].Output)
	if err != nil {
		return nil, fmt.Errorf("failed to get generic services output for chain %d: %w", c.ChainSelector(), err)
	}
	edsClient, err := edsv1.NewClientWithResponses(edsOut.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds client: %w", err)
	}

	return testhelpers.GetCCIPSendDisclosures(ctx, edsClient, ccvs)
}

// GetDisclosuresForExecution returns all the necessary disclosed contracts to execute a message on Canton using the EDS API.
func (c *Chain) GetDisclosuresForExecution(ctx context.Context, encodedMessageHex string, verifiers []contracts.InstanceAddress) (*testhelpers.ExecuteDisclosures, error) {
	// Get the EDS output from generic services output, will contain the EDS API URL
	edsOut, err := util.OpaqueToConcreteStrict[output](c.cfg.GenericServices[c.ChainSelector()].Output)
	if err != nil {
		return nil, fmt.Errorf("failed to get generic services output for chain %d: %w", c.ChainSelector(), err)
	}
	edsClient, err := edsv1.NewClientWithResponses(edsOut.EDSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eds client: %w", err)
	}

	return testhelpers.GetCCIPExecuteDisclosures(ctx, encodedMessageHex, edsClient, verifiers)
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
