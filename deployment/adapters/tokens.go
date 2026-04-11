package adapters

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	tokenadaptersfinality "github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	lock_release_token_pool "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	cantonsequences "github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	tokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
)

var (
	_ tokenadapters.TokenAdapter = &CantonTokenAdapter{}
)

type CantonTokenAdapter struct{}

func (c CantonTokenAdapter) ConfigureTokenForTransfersSequence() *operations.Sequence[tokenadapters.ConfigureTokenForTransfersInput, sequences.OnChainOutput, chain.BlockChains] {
	return operations.NewSequence(
		"canton/token-adapter/configure-token-for-transfers",
		semver.MustParse("2.0.0"),
		"Configures a Canton lock/release pool for cross-chain transfers",
		func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.ConfigureTokenForTransfersInput) (sequences.OnChainOutput, error) {
			ds := effectiveDataStore(input.ExistingDataStore)
			if ds == nil {
				return sequences.OnChainOutput{}, fmt.Errorf("existing datastore is required")
			}

			cantonChain, ok := chains.CantonChains()[input.ChainSelector]
			if !ok || len(cantonChain.Participants) == 0 {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.ChainSelector)
			}
			participant := cantonChain.Participants[0]
			poolAddress := contracts.HexToInstanceAddress(input.TokenPoolAddress)

			activePool, err := contract.FindActiveContractByInstanceAddress(
				b.GetContext(),
				participant.LedgerServices.State,
				participant.PartyID,
				lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
				poolAddress,
			)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("find active lock/release pool %s: %w", input.TokenPoolAddress, err)
			}
			parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("parse active lock/release pool %s: %w", input.TokenPoolAddress, err)
			}

			registryAddress := input.RegistryAddress
			if registryAddress == "" {
				registryRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
					input.ChainSelector,
					datastore.ContractType(token_admin_registry.ContractType),
					token_admin_registry.Version,
					"",
				))
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
				}
				registryAddress = registryRef.Address
			}

			if registryAddress != "" {
				_, err = operations.ExecuteSequence(b, cantonsequences.RegisterTokenPool, cantonChain, cantonsequences.RegisterTokenPoolInput{
					TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(registryAddress),
					InstrumentId:                      parsedPool.InstrumentId,
					PoolInstanceID:                    string(parsedPool.InstanceId),
					CcipParty:                         string(parsedPool.CcipOwner),
					PoolOwnerParty:                    string(parsedPool.PoolOwner),
				})
				if err != nil && !isAlreadyRegisteredTokenPoolError(err) {
					return sequences.OnChainOutput{}, fmt.Errorf("register lock/release pool with token admin registry: %w", err)
				}
			}

			out := sequences.OnChainOutput{}
			updates := make([]lockreleasetokenpool.ChainUpdate, 0, len(input.RemoteChains))
			for remoteSelector, remoteCfg := range input.RemoteChains {
				if _, found := lookupRemoteChainConfigValue(parsedPool.RemoteChainConfigs, strconv.FormatUint(remoteSelector, 10)); found {
					b.Logger.Infof("Canton token pool remote chain %d already configured; skipping duplicate chain update", remoteSelector)
					continue
				}

				inboundCCVs, err := resolveRawInstanceAddresses(ds, input.ChainSelector, datastore.ContractType(committee_verifier.ContractType), committee_verifier.Version, remoteCfg.InboundCCVs)
				if err != nil {
					return out, fmt.Errorf("resolve inbound CCVs for remote chain %d: %w", remoteSelector, err)
				}
				outboundCCVs, err := resolveRawInstanceAddresses(ds, input.ChainSelector, datastore.ContractType(committee_verifier.ContractType), committee_verifier.Version, remoteCfg.OutboundCCVs)
				if err != nil {
					return out, fmt.Errorf("resolve outbound CCVs for remote chain %d: %w", remoteSelector, err)
				}

				outboundRef, outboundRaw, err := deployRateLimiter(b, cantonChain, *parsedPool, input.TokenPoolAddress, remoteSelector, "outbound", remoteCfg.DefaultFinalityOutboundRateLimiterConfig)
				if err != nil {
					return out, fmt.Errorf("deploy outbound rate limiter for remote chain %d: %w", remoteSelector, err)
				}
				inboundRef, inboundRaw, err := deployRateLimiter(b, cantonChain, *parsedPool, input.TokenPoolAddress, remoteSelector, "inbound", remoteCfg.DefaultFinalityInboundRateLimiterConfig)
				if err != nil {
					return out, fmt.Errorf("deploy inbound rate limiter for remote chain %d: %w", remoteSelector, err)
				}
				out.Addresses = append(out.Addresses, outboundRef, inboundRef)
				customInboundRaw := mcms.RawInstanceAddress{}
				if remoteCfg.CustomFinalityInboundRateLimiterConfig.IsEnabled {
					customRef, customRaw, customErr := deployRateLimiter(b, cantonChain, *parsedPool, input.TokenPoolAddress, remoteSelector, "inbound-custom", remoteCfg.CustomFinalityInboundRateLimiterConfig)
					if customErr != nil {
						return out, fmt.Errorf("deploy custom inbound rate limiter for remote chain %d: %w", remoteSelector, customErr)
					}
					out.Addresses = append(out.Addresses, customRef)
					customInboundRaw = customRaw
				}

				updates = append(updates, lockreleasetokenpool.ChainUpdate{
					RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
					RemotePools:         []types.TEXT{types.TEXT(strings.TrimPrefix(strings.ToLower(gethcommon.BytesToHash(remoteCfg.RemotePool).Hex()), "0x"))},
					RemoteTokenAddress:  types.TEXT(strings.TrimPrefix(strings.ToLower(gethcommon.BytesToAddress(remoteCfg.RemoteToken).Hex()), "0x")),
					InboundCCVs:         inboundCCVs,
					OutboundCCVs:        outboundCCVs,
					FinalityConfig:      toCantonFinalityConfig(input.AllowedFinalityConfig),
					InboundRateLimiter:  inboundRaw,
					InboundCustomBlockConfirmationsRateLimiter: customInboundRaw,
					OutboundRateLimiter:                        outboundRaw,
				})
			}

			if len(updates) > 0 {
				_, err = operations.ExecuteOperation(b, lock_release_token_pool.ApplyChainUpdates, cantonChain, contract.ChoiceInput[lockreleasetokenpool.ApplyChainUpdates]{
					InstanceAddress: poolAddress,
					Args: lockreleasetokenpool.ApplyChainUpdates{
						RemoteChainSelectorsToRemove: []types.NUMERIC{},
						ChainsToAdd:                  updates,
					},
				})
				if err != nil {
					return out, fmt.Errorf("apply remote chain updates to lock/release pool: %w", err)
				}
			}

			return out, nil
		},
	)
}

func (c CantonTokenAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

func (c CantonTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef) ([]byte, error) {
	if tokenAddress, err := deriveInstrumentTokenAddress(e, chainSelector, poolRef); err == nil {
		return tokenAddress, nil
	}

	addr, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, poolRef.Type, poolRef.Version, poolRef.Qualifier))
	if err != nil {
		return nil, err
	}
	if rawAddr, rawErr := contracts.RawInstanceAddressFromString(addr.Address); rawErr == nil {
		return rawAddr.InstanceAddress().Bytes(), nil
	}

	addrBytes := gethcommon.FromHex(addr.Address)
	if len(addrBytes) == 0 {
		return nil, datastore.ErrAddressRefNotFound
	}

	return contracts.BytesToInstanceAddress(crypto.Keccak256(addrBytes)).Bytes(), nil
}

func (c CantonTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef, token []byte) (uint8, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok || len(cantonChain.Participants) == 0 {
		return 0, datastore.ErrAddressRefNotFound
	}
	instanceAddress, err := dsutils.ToInstanceAddress(poolRef)
	if err != nil {
		return 0, err
	}
	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}
	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		cantonChain.Participants[0].LedgerServices.State,
		cantonChain.Participants[0].PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		instanceAddress,
	)
	if err != nil {
		return 0, err
	}
	pool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return 0, err
	}
	return uint8(pool.Decimals), nil
}

func (c CantonTokenAdapter) DeriveTokenPoolCounterpart(e deployment.Environment, chainSelector uint64, tokenPool []byte, token []byte) ([]byte, error) {
	return tokenPool, nil
}

func (c CantonTokenAdapter) ManualRegistration() *operations.Sequence[tokenadapters.ManualRegistrationSequenceInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) SetTokenPoolRateLimits() *operations.Sequence[tokenadapters.TPRLRemotes, sequences.OnChainOutput, chain.BlockChains] {
	return operations.NewSequence(
		"canton/token-adapter/set-token-pool-rate-limits",
		semver.MustParse("2.0.0"),
		"Replaces Canton lock/release pool rate limiter contracts for a remote chain",
		func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.TPRLRemotes) (sequences.OnChainOutput, error) {
			ds := effectiveDataStore(input.ExistingDataStore)
			if ds == nil {
				return sequences.OnChainOutput{}, fmt.Errorf("existing datastore is required")
			}
			cantonChain, ok := chains.CantonChains()[input.ChainSelector]
			if !ok || len(cantonChain.Participants) == 0 {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.ChainSelector)
			}

			poolAddress, err := dsutils.ToInstanceAddress(input.TokenPoolRef)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve token pool address: %w", err)
			}
			activePool, err := contract.FindActiveContractByInstanceAddress(
				b.GetContext(),
				cantonChain.Participants[0].LedgerServices.State,
				cantonChain.Participants[0].PartyID,
				lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
				poolAddress,
			)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("find active lock/release pool: %w", err)
			}
			parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("parse active lock/release pool: %w", err)
			}

			remoteSelectorKey := strconv.FormatUint(input.RemoteChainSelector, 10)
			currentCfg, err := findRemoteChainConfig(parsedPool.RemoteChainConfigs, remoteSelectorKey)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}

			out := sequences.OnChainOutput{}
			outboundRef, outboundRaw, err := deployRateLimiterBigInt(b, cantonChain, *parsedPool, input.TokenPoolRef.Address, input.RemoteChainSelector, "outbound", input.DefaultFinalityOutboundRateLimiterConfig, common.RateLimitModeRateLimitMode_DefaultFinality)
			if err != nil {
				return out, fmt.Errorf("deploy outbound rate limiter: %w", err)
			}
			inboundRef, inboundRaw, err := deployRateLimiterBigInt(b, cantonChain, *parsedPool, input.TokenPoolRef.Address, input.RemoteChainSelector, "inbound", input.DefaultFinalityInboundRateLimiterConfig, common.RateLimitModeRateLimitMode_DefaultFinality)
			if err != nil {
				return out, fmt.Errorf("deploy inbound rate limiter: %w", err)
			}
			out.Addresses = append(out.Addresses, outboundRef, inboundRef)

			customInboundRaw := currentCfg.InboundCustomBlockConfirmationsRateLimiter
			if input.CustomFinalityInboundRateLimiterConfig.IsEnabled {
				customRef, customRaw, customErr := deployRateLimiterBigInt(b, cantonChain, *parsedPool, input.TokenPoolRef.Address, input.RemoteChainSelector, "inbound-custom", input.CustomFinalityInboundRateLimiterConfig, common.RateLimitModeRateLimitMode_CustomFinality)
				if customErr != nil {
					return out, fmt.Errorf("deploy custom inbound rate limiter: %w", customErr)
				}
				out.Addresses = append(out.Addresses, customRef)
				customInboundRaw = customRaw
			}

			_, err = operations.ExecuteOperation(b, lock_release_token_pool.SetRateLimitConfig, cantonChain, contract.ChoiceInput[lockreleasetokenpool.SetRateLimitConfig]{
				InstanceAddress: poolAddress,
				Args: lockreleasetokenpool.SetRateLimitConfig{
					Caller: parsedPool.PoolOwner,
					RateLimitConfigArgs: []lockreleasetokenpool.RateLimitConfigArgs{
						{
							RemoteChainSelector:                        types.NUMERIC(remoteSelectorKey),
							InboundRateLimiter:                         inboundRaw,
							InboundCustomBlockConfirmationsRateLimiter: customInboundRaw,
							OutboundRateLimiter:                        outboundRaw,
						},
					},
				},
			})
			if err != nil {
				return out, fmt.Errorf("set rate limit config on lock/release pool: %w", err)
			}

			return out, nil
		},
	)
}

func (c CantonTokenAdapter) DeployToken() *operations.Sequence[tokenadapters.DeployTokenInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeployTokenVerify(e deployment.Environment, in tokenadapters.DeployTokenInput) error {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeployTokenPoolForToken() *operations.Sequence[tokenadapters.DeployTokenPoolInput, sequences.OnChainOutput, chain.BlockChains] {
	return operations.NewSequence(
		"canton/token-adapter/deploy-token-pool-for-token",
		semver.MustParse("2.0.0"),
		"Deploys a Canton lock/release pool and returns both the canonical and logical datastore refs",
		func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.DeployTokenPoolInput) (sequences.OnChainOutput, error) {
			if input.TokenPoolVersion == nil {
				return sequences.OnChainOutput{}, fmt.Errorf("TokenPoolVersion is required")
			}
			if datastore.ContractType(input.PoolType) != datastore.ContractType("LockReleaseTokenPool") {
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported Canton token pool type %q", input.PoolType)
			}

			ds := effectiveDataStore(input.ExistingDataStore)
			if ds == nil {
				return sequences.OnChainOutput{}, fmt.Errorf("existing datastore is required")
			}
			qualifier := input.TokenPoolQualifier
			if qualifier == "" {
				qualifier = "Amulet"
			}
			matches := ds.Addresses().Filter(
				datastore.AddressRefByType(datastore.ContractType(input.PoolType)),
				datastore.AddressRefByChainSelector(input.ChainSelector),
				datastore.AddressRefByQualifier(qualifier),
				datastore.AddressRefByVersion(input.TokenPoolVersion),
			)
			if len(matches) > 1 {
				return sequences.OnChainOutput{}, fmt.Errorf("multiple Canton token pools found with qualifier %q", qualifier)
			}
			if len(matches) == 1 {
				b.Logger.Infof("Canton token pool already deployed at %s", matches[0].Address)
				return sequences.OnChainOutput{}, nil
			}

			cantonChain, ok := chains.CantonChains()[input.ChainSelector]
			if !ok || len(cantonChain.Participants) == 0 {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.ChainSelector)
			}
			participant := cantonChain.Participants[0]

			instrumentAdmin := participant.PartyID
			if registryAdmin, err := resolveRegistryAdmin(participant); err == nil && registryAdmin != "" {
				instrumentAdmin = registryAdmin
			}
			instrumentID := splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(instrumentAdmin),
				Id:    types.TEXT("Amulet"),
			}

			tokenAdminRegistryRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
				input.ChainSelector,
				datastore.ContractType(token_admin_registry.ContractType),
				token_admin_registry.Version,
				"",
			))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
			}
			rmnRemoteRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
				input.ChainSelector,
				datastore.ContractType(rmn_remote.ContractType),
				rmn_remote.Version,
				"",
			))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve rmn remote: %w", err)
			}
			feeQuoterRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
				input.ChainSelector,
				datastore.ContractType(fee_quoter.ContractType),
				fee_quoter.Version,
				"",
			))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve fee quoter: %w", err)
			}

			tokenAdminRegistryRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(tokenAdminRegistryRef)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve token admin registry raw address: %w", err)
			}
			rmnRemoteRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(rmnRemoteRef)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve rmn remote raw address: %w", err)
			}
			feeQuoterRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(feeQuoterRef)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve fee quoter raw address: %w", err)
			}

			relativeHours := types.INT64(24)
			deployReport, err := operations.ExecuteOperation(b, lock_release_token_pool.Deploy, cantonChain, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
				Qualifier: &qualifier,
				Template: lockreleasetokenpool.LockReleaseTokenPool{
					CcipOwner:    types.PARTY(participant.PartyID),
					PoolOwner:    types.PARTY(participant.PartyID),
					InstrumentId: instrumentID,
					Decimals:     types.INT64(10),
					PoolReceiveContext: common.CCIPContext{
						Values: types.TEXTMAP{},
					},
					TransferTimeout: lockreleasetokenpool.TransferTimeout{
						RelativeHours: &relativeHours,
					},
					RemoteChainConfigs:      types.GENMAP{},
					TokenTransferFeeConfigs: types.GENMAP{},
					Deps: lockreleasetokenpool.LockReleaseTokenPoolDeps{
						TokenAdminRegistry: tokenAdminRegistryRaw.Binding(),
						RmnRemote:          rmnRemoteRaw.Binding(),
						FeeQuoter:          feeQuoterRaw.Binding(),
					},
				},
				OwnerParty: types.PARTY(participant.PartyID),
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("deploy Canton lock/release pool: %w", err)
			}
			rawPoolAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(deployReport.Output)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("resolve deployed pool raw address: %w", err)
			}
			_, err = operations.ExecuteSequence(b, cantonsequences.RegisterTokenPool, cantonChain, cantonsequences.RegisterTokenPoolInput{
				TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
				InstrumentId:                      instrumentID,
				PoolInstanceID:                    rawPoolAddr.InstanceID(),
				CcipParty:                         participant.PartyID,
				PoolOwnerParty:                    participant.PartyID,
			})
			if err != nil && !isAlreadyRegisteredTokenPoolError(err) {
				return sequences.OnChainOutput{}, fmt.Errorf("register Canton lock/release pool: %w", err)
			}

			logicalRef := datastore.AddressRef{
				Address:       deployReport.Output.Address,
				Labels:        deployReport.Output.Labels,
				Type:          datastore.ContractType(input.PoolType),
				Version:       input.TokenPoolVersion,
				Qualifier:     qualifier,
				ChainSelector: input.ChainSelector,
			}
			tokenAddress := gethcommon.Bytes2Hex(cantonInstrumentTokenAddress(instrumentID))
			tokenRef := datastore.AddressRef{
				Address:       tokenAddress,
				Type:          datastore.ContractType("Token"),
				Qualifier:     qualifier,
				ChainSelector: input.ChainSelector,
			}
			return sequences.OnChainOutput{
				Addresses: []datastore.AddressRef{
					deployReport.Output,
					logicalRef,
					tokenRef,
				},
			}, nil
		},
	)
}

func (c CantonTokenAdapter) UpdateAuthorities() *operations.Sequence[tokenadapters.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *operations.Sequence[tokenadapters.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func deriveInstrumentTokenAddress(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef) ([]byte, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok || len(cantonChain.Participants) == 0 {
		return nil, datastore.ErrAddressRefNotFound
	}
	addr, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return nil, err
	}
	instanceAddress := contracts.HexToInstanceAddress(addr.Address)
	participant := cantonChain.Participants[0]
	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}
	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		instanceAddress,
	)
	if err != nil {
		return nil, err
	}
	pool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return nil, err
	}
	instrumentCombined := string(pool.InstrumentId.Id) + "@" + string(pool.InstrumentId.Admin)

	return crypto.Keccak256([]byte(instrumentCombined)), nil
}

func cantonInstrumentTokenAddress(instrumentID splice_api_token_holding_v1.InstrumentId) []byte {
	instrumentCombined := string(instrumentID.Id) + "@" + string(instrumentID.Admin)
	return crypto.Keccak256([]byte(instrumentCombined))
}

func effectiveDataStore(ds datastore.DataStore) datastore.DataStore {
	if ds != nil {
		return ds
	}
	return getRuntimeDataStore()
}

func resolveRegistryAdmin(participant cldfcanton.Participant) (string, error) {
	tokenSource := participant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}
	client, err := tokenMetadataV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL), tokenMetadataV1.WithRequestEditorFn(interceptor))
	if err != nil {
		return "", fmt.Errorf("create token metadata client: %w", err)
	}
	resp, err := client.GetRegistryInfoWithResponse(context.Background())
	if err != nil {
		return "", fmt.Errorf("get registry info: %w", err)
	}
	if resp.StatusCode() != 200 || resp.JSON200 == nil {
		return "", fmt.Errorf("unexpected registry info status %d", resp.StatusCode())
	}
	return resp.JSON200.AdminId, nil
}

func resolveRawInstanceAddresses(ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType, version *semver.Version, addresses []string) ([]mcms.RawInstanceAddress, error) {
	refs := ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(contractType),
		datastore.AddressRefByVersion(version),
	)
	out := make([]mcms.RawInstanceAddress, 0, len(addresses))
	for _, address := range addresses {
		var matched *datastore.AddressRef
		for _, ref := range refs {
			if strings.EqualFold(ref.Address, address) {
				refCopy := ref
				matched = &refCopy
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("could not resolve %s ref for address %s", contractType, address)
		}
		rawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(*matched)
		if err != nil {
			return nil, err
		}
		out = append(out, rawAddr.Binding())
	}
	return out, nil
}

func deployRateLimiter(
	b operations.Bundle,
	cantonChain cldfcanton.Chain,
	pool lockreleasetokenpool.LockReleaseTokenPool,
	poolAddress string,
	remoteSelector uint64,
	direction string,
	cfg tokenadapters.RateLimiterConfigFloatInput,
) (datastore.AddressRef, mcms.RawInstanceAddress, error) {
	outbound, inbound := tokenadapters.GenerateTPRLConfigs(
		cfg,
		cfg,
		uint8(pool.Decimals),
		uint8(pool.Decimals),
		chain_selectors.FamilyCanton,
		semver.MustParse("2.0.0"),
	)
	mode := common.RateLimitModeRateLimitMode_DefaultFinality
	if direction == "outbound" {
		return deployRateLimiterBigInt(b, cantonChain, pool, poolAddress, remoteSelector, direction, outbound, common.RateLimitModeRateLimitMode_DefaultFinality)
	}
	if direction == "inbound-custom" {
		mode = common.RateLimitModeRateLimitMode_CustomFinality
	}
	return deployRateLimiterBigInt(b, cantonChain, pool, poolAddress, remoteSelector, direction, inbound, mode)
}

func deployRateLimiterBigInt(
	b operations.Bundle,
	cantonChain cldfcanton.Chain,
	pool lockreleasetokenpool.LockReleaseTokenPool,
	poolAddress string,
	remoteSelector uint64,
	direction string,
	cfg tokenadapters.RateLimiterConfig,
	mode common.RateLimitMode,
) (datastore.AddressRef, mcms.RawInstanceAddress, error) {
	deployOp := rate_limiter.DeployInbound
	dirEnum := common.RateLimitDirectionRateLimitDirection_Inbound
	if direction == "outbound" {
		deployOp = rate_limiter.DeployOutbound
		dirEnum = common.RateLimitDirectionRateLimitDirection_Outbound
	}
	capacity := types.NUMERIC("0")
	if cfg.Capacity != nil {
		capacity = types.NUMERIC(cfg.Capacity.String())
	}
	rate := types.NUMERIC("0")
	if cfg.Rate != nil {
		rate = types.NUMERIC(cfg.Rate.String())
	}
	qualifier := fmt.Sprintf("%s-%s-%d", poolAddress, direction, remoteSelector)
	report, err := operations.ExecuteOperation(b, deployOp, cantonChain, contract.DeployInput[common.RateLimiter]{
		Qualifier: &qualifier,
		Template: common.RateLimiter{
			PoolInstanceId:      pool.InstanceId,
			PoolOwner:           pool.PoolOwner,
			RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
			Direction:           dirEnum,
			Mode:                mode,
			IsEnabled:           types.BOOL(cfg.IsEnabled),
			Capacity:            capacity,
			Rate:                rate,
			Tokens:              capacity,
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
		OwnerParty: pool.PoolOwner,
	})
	if err != nil {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, err
	}
	rawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(report.Output)
	if err != nil {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, err
	}
	return report.Output, rawAddr.Binding(), nil
}

func findRemoteChainConfig(remoteChainConfigs map[string]any, selectorKey string) (lockreleasetokenpool.RemoteChainConfig, error) {
	cfgAny, ok := lookupRemoteChainConfigValue(remoteChainConfigs, selectorKey)
	if !ok {
		return lockreleasetokenpool.RemoteChainConfig{}, fmt.Errorf("missing remote chain config for selector %s", selectorKey)
	}

	return remoteChainConfigFromAny(cfgAny)
}

func lookupRemoteChainConfigValue(remoteChainConfigs map[string]any, selectorKey string) (any, bool) {
	if cfgAny, ok := remoteChainConfigs[selectorKey]; ok {
		return cfgAny, true
	}
	normalizedSelector := normalizeNumericForCompare(selectorKey)
	for rawKey, cfgAny := range remoteChainConfigs {
		if normalizeNumericForCompare(rawKey) == normalizedSelector {
			return cfgAny, true
		}
	}
	return nil, false
}

func remoteChainConfigFromAny(v any) (lockreleasetokenpool.RemoteChainConfig, error) {
	switch cfg := v.(type) {
	case lockreleasetokenpool.RemoteChainConfig:
		return cfg, nil
	case map[string]any:
		m := cfg
		if data, ok := cfg["data"].(map[string]any); ok {
			m = data
		}
		out := lockreleasetokenpool.RemoteChainConfig{}
		if raw, ok := m["inboundCustomBlockConfirmationsRateLimiter"]; ok {
			unpack, err := extractRawInstanceAddress(raw)
			if err != nil {
				return lockreleasetokenpool.RemoteChainConfig{}, fmt.Errorf("decode inbound custom rate limiter: %w", err)
			}
			out.InboundCustomBlockConfirmationsRateLimiter = mcms.RawInstanceAddress{Unpack: types.TEXT(unpack)}
		}

		return out, nil
	default:
		return lockreleasetokenpool.RemoteChainConfig{}, fmt.Errorf("unexpected remote chain config type %T", v)
	}
}

func toCantonFinalityConfig(cfg tokenadaptersfinality.Config) common.FinalityConfig {
	switch {
	case cfg.WaitForFinality:
		return common.FinalityConfig{WaitForFinality: &types.UNIT{}}
	case cfg.WaitForSafe:
		return common.FinalityConfig{WaitForSafe: &types.UNIT{}}
	case cfg.BlockDepth > 0:
		depth := types.INT64(cfg.BlockDepth)
		return common.FinalityConfig{BlockDepth: &depth}
	default:
		return common.FinalityConfig{}
	}
}

func extractRawInstanceAddress(v any) (string, error) {
	switch value := v.(type) {
	case mcms.RawInstanceAddress:
		return string(value.Unpack), nil
	case map[string]any:
		m := value
		if data, ok := value["data"].(map[string]any); ok {
			m = data
		}
		if unpack, ok := m["unpack"].(string); ok && unpack != "" {
			return unpack, nil
		}
	}
	return "", fmt.Errorf("unexpected raw instance address type %T", v)
}

func normalizeNumericForCompare(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	if strings.Contains(v, ".") {
		parts := strings.SplitN(v, ".", 2)
		frac := strings.TrimRight(parts[1], "0")
		if frac == "" {
			return parts[0]
		}
	}
	return strings.TrimSuffix(v, ".")
}

func isAlreadyRegisteredTokenPoolError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "TokenAdminRegistry_ProposeAdministrator: Admin already set") ||
		strings.Contains(msg, "TokenAdminRegistry_SetPool")
}
