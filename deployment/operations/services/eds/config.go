package eds

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	edsConfig "github.com/smartcontractkit/chainlink-canton/eds/config"
)

type GenerateEDSConfigInput struct {
	ChainSelector uint64 `json:"chain_selector"`
	Participant   int    `json:"participant"`
}

type GenerateEDSConfigOutput struct {
	NodeConfig edsConfig.NodeConfig
	Contracts  edsConfig.Contracts
}

type BuildConfigDeps struct {
	Env deployment.Environment
}

var BuildConfig = operations.NewOperation(
	"build-eds-config",
	semver.MustParse("1.0.0"),
	"Builds the EDS node + contracts config from datastore contract addresses and participants in the environment",
	func(b operations.Bundle, deps BuildConfigDeps, input GenerateEDSConfigInput) (output GenerateEDSConfigOutput, err error) {
		env := deps.Env
		chain, ok := env.BlockChains.CantonChains()[input.ChainSelector]
		if !ok {
			return GenerateEDSConfigOutput{}, fmt.Errorf("chain selector %d not found", input.ChainSelector)
		}
		participant := chain.Participants[input.Participant]
		jwt, err := participant.TokenSource.Token()
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("getting token: %w", err)
		}

		// Get all relevant addresses
		globalConfig, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("getting global config: %w", err)
		}
		onRamp, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(onramp.ContractType), onramp.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get OnRamp address: %w", err)
		}
		offRamp, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(offramp.ContractType), offramp.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get OffRamp address: %w", err)
		}
		tokenAdminRegistry, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get TokenAdminRegistry address: %w", err)
		}
		rmnRemote, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get RMNRemote address: %w", err)
		}
		perPartyRouterFactory, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(per_party_router_factory.ContractType), per_party_router_factory.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get PerPartyRouterFactory address: %w", err)
		}
		feeQuoter, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version, ""))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get FeeQuoter address: %w", err)
		}
		executor, err := env.DataStore.Addresses().Get(datastore.NewAddressRefKey(input.ChainSelector, datastore.ContractType(executor.ContractType), executor.Version, "default"))
		if err != nil {
			return GenerateEDSConfigOutput{}, fmt.Errorf("failed to get Executor address: %w", err)
		}

		// CCVs
		refs := env.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(committee_verifier.ContractType)),
		)
		if len(refs) == 0 {
			return GenerateEDSConfigOutput{}, fmt.Errorf("no CommitteeVerifier contracts found in datastore")
		}
		ccvs := make([]edsConfig.ContractIdentifier, len(refs))
		for i, ref := range refs {
			ccvs[i] = edsConfig.ContractIdentifier{
				PartyID:         participant.PartyID,
				InstanceAddress: contracts.HexToInstanceAddress(ref.Address),
			}
		}

		// Token pools on Canton are LockReleaseTokenPool (not EVM BurnMintTokenPool).
		refs = env.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(lock_release_token_pool.ContractType)),
		)
		inboundRefs := env.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(rate_limiter.ContractTypeInbound)),
			datastore.AddressRefByVersion(rate_limiter.Version),
		)
		outboundRefs := env.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(rate_limiter.ContractTypeOutbound)),
			datastore.AddressRefByVersion(rate_limiter.Version),
		)

		tokenPools := make([]edsConfig.TokenPoolContracts, 0, len(refs))
		for _, ref := range refs {
			// Rate limiter refs are keyed by qualifier as:
			//   <poolAddress>-inbound-<remoteSelector>
			//   <poolAddress>-inbound-custom-<remoteSelector>
			//   <poolAddress>-outbound-<remoteSelector>
			// so we recover the correct limiters for this pool by filtering on the
			// pool-address prefix first, then grouping the remaining refs by remote selector.
			prefix := ref.Address + "-"
			inboundByRemote := make(map[string]datastore.AddressRef)
			inboundCustomByRemote := make(map[string]datastore.AddressRef)
			outboundByRemote := make(map[string]datastore.AddressRef)

			for _, rlRef := range inboundRefs {
				if !strings.HasPrefix(rlRef.Qualifier, prefix) {
					continue
				}
				remoteSelector := strings.TrimPrefix(rlRef.Qualifier, prefix+"inbound-custom-")
				if remoteSelector != rlRef.Qualifier {
					inboundCustomByRemote[remoteSelector] = rlRef
					continue
				}
				remoteSelector = strings.TrimPrefix(rlRef.Qualifier, prefix+"inbound-")
				if remoteSelector == rlRef.Qualifier {
					continue
				}
				inboundByRemote[remoteSelector] = rlRef
			}
			for _, rlRef := range outboundRefs {
				if !strings.HasPrefix(rlRef.Qualifier, prefix+"outbound-") {
					continue
				}
				remoteSelector := strings.TrimPrefix(rlRef.Qualifier, prefix+"outbound-")
				outboundByRemote[remoteSelector] = rlRef
			}

			if len(outboundByRemote) == 0 {
				return GenerateEDSConfigOutput{}, fmt.Errorf("no rate limiter refs found in datastore for token pool %s", ref.Address)
			}

			for remoteSelector, outboundRef := range outboundByRemote {
				inboundRef, ok := inboundByRemote[remoteSelector]
				if !ok {
					return GenerateEDSConfigOutput{}, fmt.Errorf("missing inbound rate limiter ref for token pool %s remote %s", ref.Address, remoteSelector)
				}
				inboundAddress := contracts.HexToInstanceAddress(inboundRef.Address)
				inboundCustomAddress := contracts.HexToInstanceAddress("0x0")
				if inboundCustomRef, ok := inboundCustomByRemote[remoteSelector]; ok {
					inboundCustomAddress = contracts.HexToInstanceAddress(inboundCustomRef.Address)
				}
				outboundAddress := contracts.HexToInstanceAddress(outboundRef.Address)

				tokenPools = append(tokenPools, edsConfig.TokenPoolContracts{
					TokenPool: edsConfig.ContractIdentifier{
						PartyID:         participant.PartyID,
						InstanceAddress: contracts.HexToInstanceAddress(ref.Address),
					},
					InboundRateLimiter: edsConfig.ContractIdentifier{
						PartyID:         participant.PartyID,
						InstanceAddress: inboundAddress,
					},
					InboundCustomBlockConfirmationsRateLimiter: edsConfig.ContractIdentifier{
						PartyID:         participant.PartyID,
						InstanceAddress: inboundCustomAddress,
					},
					OutboundRateLimiter: edsConfig.ContractIdentifier{
						PartyID:         participant.PartyID,
						InstanceAddress: outboundAddress,
					},
				})
			}
		}

		var nodeConfig edsConfig.NodeConfig
		if participant.InternalEndpoints == nil {
			// No internal endpoints, so we assume the node is externally accessible and can be reached via the GRPC ledger API URL
			nodeConfig = edsConfig.NodeConfig{
				URL: participant.Endpoints.GRPCLedgerAPIURL,
				AuthConfig: commonconfig.AuthConfig{
					Type:   commonconfig.AuthTypeInsecureStatic,
					UserID: participant.UserID,
					JWT:    jwt.AccessToken,
				},
				MaxRetries: 0,
			}
		} else {
			nodeConfig = edsConfig.NodeConfig{
				URL: participant.InternalEndpoints.GRPCLedgerAPIURL,
				AuthConfig: commonconfig.AuthConfig{
					Type:   commonconfig.AuthTypeInsecureStatic,
					UserID: participant.UserID,
					JWT:    jwt.AccessToken,
				},
				MaxRetries: 0,
			}
		}

		return GenerateEDSConfigOutput{
			NodeConfig: nodeConfig,
			Contracts: edsConfig.Contracts{
				PerPartyRouterFactory: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(perPartyRouterFactory.Address),
				},
				OnRamp: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(onRamp.Address),
				},
				OffRamp: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(offRamp.Address),
				},
				GlobalConfig: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(globalConfig.Address),
				},
				TokenAdminRegistry: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistry.Address),
				},
				RMNRemote: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(rmnRemote.Address),
				},
				FeeQuoter: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(feeQuoter.Address),
				},
				DefaultExecutor: edsConfig.ContractIdentifier{
					PartyID:         participant.PartyID,
					InstanceAddress: contracts.HexToInstanceAddress(executor.Address),
				},
				CCVs:               ccvs,
				TokenPoolContracts: tokenPools,
				PoolOwner:          participant.PartyID,
			},
		}, nil
	},
)
