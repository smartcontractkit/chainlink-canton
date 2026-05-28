package devenv

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	ledgerv2admin "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/services/committeeverifier"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/util"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/ccip"
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

const (
	DefaultCantonCommitteVerifierImage = "committeeverifier-canton:latest"
)

func CommitteeVerifierConfigLoader(outputs []*blockchain.Output) (map[string]any, error) {
	ret := make(map[string]any)
	for _, output := range outputs {
		if output.Family != chainsel.FamilyCanton {
			continue
		}

		chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(output.ChainID, output.Family)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain details for chain %s, family %s: %w", output.ChainID, output.Family, err)
		}

		strSelector := strconv.FormatUint(chainDetails.ChainSelector, 10)

		// Return placeholder values here, the real values will be pulled from the mounted config in the container.
		ret[strSelector] = ccip.BlockchainInfo{
			GRPCLedgerAPIURL: "dontuse",
			Auth: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeStatic,
				JWT:  "dontuse",
			},
		}
	}

	return ret, nil
}

// CommitteeVerifierModifier modifies a testcontainers.ContainerRequest for canton.
func CommitteeVerifierModifier(req testcontainers.ContainerRequest, verifierInput *committeeverifier.Input, outputs []*blockchain.Output) (testcontainers.ContainerRequest, error) {
	// Use the canton committee verifier image to properly read from Canton.
	req.Image = DefaultCantonCommitteVerifierImage

	// Update name to reflect chain family.
	req.Name = fmt.Sprintf("canton-%s", verifierInput.ContainerName)

	// Marshal the canton config into TOML bytes.
	cantonConfigBytes, err := hydrateAndMarshalCantonConfig(verifierInput, outputs)
	if err != nil {
		return req, fmt.Errorf("failed to hydrate and marshal canton config: %w", err)
	}

	// Copy the hydrated config into the container using the supported Files API.
	req.Files = append(req.Files, testcontainers.ContainerFile{
		Reader:            bytes.NewReader(cantonConfigBytes),
		ContainerFilePath: ccip.DefaultCantonConfigPath,
		FileMode:          0o644,
	})

	return req, nil
}

// hydrateAndMarshalCantonConfig hydrates the canton config with the full party ID for the CCIPOwnerParty.
func hydrateAndMarshalCantonConfig(in *committeeverifier.Input, outputs []*blockchain.Output) ([]byte, error) {
	cantonConfigs, err := util.OpaqueToConcreteStrict[ccip.Config](in.OpaqueConfigs[chainsel.FamilyCanton])
	if err != nil {
		return nil, fmt.Errorf("failed to get canton config from opaque: %w", err)
	}
	for _, output := range outputs {
		if output.Family != chainsel.FamilyCanton {
			continue
		}

		chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(output.ChainID, output.Family)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain details for chain %s, family %s: %w", output.ChainID, output.Family, err)
		}

		strSelector := strconv.FormatUint(chainDetails.ChainSelector, 10)
		readerConfig, ok := cantonConfigs.ReaderConfigs[strSelector]
		if !ok {
			return nil, fmt.Errorf("no canton config found for chain %s, please update the config appropriately if you're using canton", strSelector)
		}
		if readerConfig.CCIPOwnerParty == "" {
			return nil, fmt.Errorf("CCIPOwnerParty is not set for chain %s, please update the config appropriately if you're using canton", strSelector)
		}
		if readerConfig.CCIPMessageSentTemplateID == (contracts.TemplateID{}) {
			return nil, fmt.Errorf("CCIPMessageSentTemplateID is not set for chain %s, please update the config appropriately if you're using canton", strSelector)
		}
		if readerConfig.RMNRemoteTemplateID == (contracts.TemplateID{}) {
			return nil, fmt.Errorf("RMNRemoteTemplateID is not set for chain %s, please update the config appropriately if you're using canton", strSelector)
		}

		// Get the full party ID (name + hex id) from the canton participant.
		// TODO: how to support multiple participants?
		participant := output.NetworkSpecificData.CantonData.ExternalEndpoints.Participants[0]
		grpcURL := participant.GRPCLedgerAPIURL
		if grpcURL == "" {
			return nil, fmt.Errorf("GRPC ledger API URL is not set for chain %s, please update the config appropriately if you're using canton", strSelector)
		}

		authCfg, err := buildParticipantAuthConfig(participant)
		if err != nil {
			return nil, fmt.Errorf("participant auth config for chain %s: %w", strSelector, err)
		}
		authProvider, err := participantAuthProvider(context.Background(), authCfg)
		if err != nil {
			return nil, fmt.Errorf("participant auth provider for chain %s: %w", strSelector, err)
		}

		cantonConfigs.BlockchainInfos[strSelector] = ccip.BlockchainInfo{
			GRPCLedgerAPIURL: output.NetworkSpecificData.CantonData.InternalEndpoints.Participants[0].GRPCLedgerAPIURL,
			Auth:             authCfg,
		}

		// find the party that starts with the prefix that is listed in the canton config.
		conn, err := grpc.NewClient(
			grpcURL,
			grpc.WithTransportCredentials(authProvider.TransportCredentials()),
			grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
		}
		defer conn.Close() //nolint:revive // defer is used to close the connection after the function returns, and we don't spin up that many chains

		resp, err := ledgerv2admin.NewPartyManagementServiceClient(conn).ListKnownParties(context.Background(), &ledgerv2admin.ListKnownPartiesRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		authority := grpcURL
		if idx := strings.LastIndex(authority, ":"); idx != -1 {
			authority = authority[:idx]
		}

		var found bool
		for _, partyDetail := range resp.PartyDetails {
			if strings.HasPrefix(partyDetail.GetParty(), readerConfig.CCIPOwnerParty) {
				cantonConfigs.ReaderConfigs[strSelector] = sourcereader.ReaderConfig{
					// TODO: ideally node operator party and ccip owner party should be separate parties.
					NodeOperatorParty:         partyDetail.GetParty(),
					CCIPOwnerParty:            partyDetail.GetParty(),
					CCIPMessageSentTemplateID: readerConfig.CCIPMessageSentTemplateID,
					RMNRemoteTemplateID:       readerConfig.RMNRemoteTemplateID,
					Authority:                 authority,
				}
				found = true

				break
			}
		}
		if !found {
			return nil, fmt.Errorf("expected CCIPOwnerParty %s not found for canton chain %s, please update the config appropriately if you're using canton", readerConfig.CCIPOwnerParty, strSelector)
		}
	}

	// Marshal the canton config into TOML.
	cantonConfigBytes, err := toml.Marshal(cantonConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal canton config: %w", err)
	}

	// to reduce confusion, re-set the config in the input to what we generated here
	// so that env-out.toml has the real values, not the placeholder values.
	newOpaque := make(util.OpaqueConfig)
	newOpaque["reader_configs"] = cantonConfigs.ReaderConfigs
	newOpaque["blockchain_infos"] = cantonConfigs.BlockchainInfos
	in.OpaqueConfigs[chainsel.FamilyCanton] = newOpaque

	return cantonConfigBytes, nil
}
