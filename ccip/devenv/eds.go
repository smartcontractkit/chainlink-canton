package devenv

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/exp/maps"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/chainreg"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/util"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfdeployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment"
	cantonChangesets "github.com/smartcontractkit/chainlink-canton/deployment/changesets"
	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	edsTesthelpers "github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

func (c *Chain) GetPerPartyRouterFactoryDisclosure(ctx context.Context, partyId string) (*edsTesthelpers.PerPartyRouterFactoryDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIP EDS client: %w", err)
	}

	return edsTesthelpers.GetPerPartyRouterFactoryDisclosure(ctx, ccipAPIClient, partyId)
}

func (c *Chain) GetTokenPoolForToken(ctx context.Context, token contracts.EncodedInstrumentID) (contracts.RawInstanceAddress, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return contracts.RawInstanceAddress(""), err
	}
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(edsURL)
	if err != nil {
		return contracts.RawInstanceAddress(""), fmt.Errorf("failed to create CCIP EDS client: %w", err)
	}

	return edsTesthelpers.GetTokenPoolForToken(ctx, ccipAPIClient, token)
}

func (c *Chain) GetTokenPoolSendDisclosure(ctx context.Context, message oapiCommon.Message, tokenPoolAddress contracts.InstanceAddress) (*edsTesthelpers.TokenPoolSendDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Token Pool EDS client: %w", err)
	}

	return edsTesthelpers.GetTokenPoolSendDisclosure(ctx, tokenPoolAPIClient, message, tokenPoolAddress)
}

func (c *Chain) GetCCIPSendDisclosure(ctx context.Context, message oapiCommon.Message, senderRequiredCCVs, tokenPoolRequiredCCVs []string) (*edsTesthelpers.CCIPSendDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIP EDS client: %w", err)
	}

	return edsTesthelpers.GetCCIPSendDisclosure(ctx, ccipAPIClient, message, senderRequiredCCVs, tokenPoolRequiredCCVs)
}

func (c *Chain) GetCCVSendDisclosure(ctx context.Context, message oapiCommon.Message, ccvAddress contracts.InstanceAddress) (*edsTesthelpers.CCVSendDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCV EDS client: %w", err)
	}

	return edsTesthelpers.GetCCVSendDisclosure(ctx, ccvAPIClient, message, ccvAddress)
}

func (c *Chain) GetExecutorSendDisclosure(ctx context.Context, message oapiCommon.Message, executorAddress contracts.InstanceAddress, ccvAddresses []string) (*edsTesthelpers.ExecutorDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	executorAPIClient, err := oapiExecutor.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Executor EDS client: %w", err)
	}

	return edsTesthelpers.GetExecutorSendDisclosure(ctx, executorAPIClient, message, executorAddress, ccvAddresses)
}

func (c *Chain) GetTokenPoolExecuteDisclosure(ctx context.Context, encodedMessageHex string, tokenPoolAddress contracts.InstanceAddress) (*edsTesthelpers.TokenPoolExecuteDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	tokenPoolAPIClient, err := oapiTokenPool.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Token Pool EDS client: %w", err)
	}

	return edsTesthelpers.GetTokenPoolExecuteDisclosure(ctx, tokenPoolAPIClient, encodedMessageHex, tokenPoolAddress)
}

func (c *Chain) GetCCIPExecuteDisclosure(ctx context.Context, encodedMessageHex string) (*edsTesthelpers.CCIPExecuteDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	ccipAPIClient, err := oapiCCIP.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIP EDS client: %w", err)
	}

	return edsTesthelpers.GetCCIPExecuteDisclosure(ctx, ccipAPIClient, encodedMessageHex)
}

func (c *Chain) GetCCVExecuteDisclosure(ctx context.Context, encodedMessageHex string, ccvAddress contracts.InstanceAddress) (*edsTesthelpers.CCVExecuteDisclosure, error) {
	edsURL, err := deployment.GetEDSURL(c.e.DataStore)
	if err != nil {
		return nil, err
	}
	ccvAPIClient, err := oapiCCV.NewClientWithResponses(edsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCV EDS client: %w", err)
	}

	return edsTesthelpers.GetCCVExecuteDisclosure(ctx, ccvAPIClient, encodedMessageHex, ccvAddress)
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
			hc.PortBindings = network.PortMap{
				network.MustParsePort(fmt.Sprintf("%d/tcp", cfg.Server.Port)): []network.PortBinding{
					{HostPort: ""}, // Docker assigns a random free host port.
				},
			}
		},
		// TODO: properly implement readiness endpoint
		WaitingFor: wait.ForLog("Backfill complete").WithStartupTimeout(2 * time.Minute),
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

var _ chainreg.Launcher = &launcher{}

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
	edsMappedPort, err := edsContainer.MappedPort(ctx, fmt.Sprintf("%d/tcp", edsCfg.Server.Port))
	if err != nil {
		pm, _ := edsContainer.Ports(ctx) // Add all existing ports to error
		return nil, fmt.Errorf("failed to get EDS container mapped port (ports: %v): %w", maps.Keys(pm), err)
	}

	// return the EDS config
	// Add the config and container URL to the output
	edsURL := fmt.Sprintf("http://%s:%s", host, edsMappedPort.Port())
	opaqueConfig, err := util.ConcreteToOpaque(output{
		EDSConfig: *edsCfg,
		EDSURL:    edsURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert output to opaque config: %w", err)
	}

	// Update the CLDF env datastore with the EDS URL,
	// so that we can fetch it in the Canton chain impl.
	// env.DataStore is immutable, so we need to create a new memory data store
	// and merge the existing data store into it.
	// Then we save the EDS URL to the new data store.
	ds := datastore.NewMemoryDataStore()
	if err := ds.Merge(env.DataStore); err != nil {
		return nil, fmt.Errorf("failed to merge existing datastore: %w", err)
	}
	if err := deployment.SaveEDSURL(ds, edsURL); err != nil {
		return nil, fmt.Errorf("failed to save EDS URL: %w", err)
	}
	env.DataStore = ds.Seal()

	return opaqueConfig, nil
}
