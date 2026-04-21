package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/aws/smithy-go/ptr"
	gethcommon "github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	tokenadaptersfinality "github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ConfigureTokenForTransfers = operations.NewSequence(
	"canton/token-adapter/configure-token-for-transfers",
	semver.MustParse("2.0.0"),
	"Configures a Canton lock/release pool for cross-chain transfers",
	func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.ConfigureTokenForTransfersInput) (ccipsequences.OnChainOutput, error) {
		if input.ExistingDataStore == nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("existing datastore is required")
		}

		cantonChain, ok := chains.CantonChains()[input.ChainSelector]
		if !ok || len(cantonChain.Participants) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.ChainSelector)
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
			return ccipsequences.OnChainOutput{}, fmt.Errorf("find active lock/release pool %s: %w", input.TokenPoolAddress, err)
		}
		parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("parse active lock/release pool %s: %w", input.TokenPoolAddress, err)
		}

		registryRef, err := input.ExistingDataStore.Addresses().Get(datastore.NewAddressRefKey(
			input.ChainSelector,
			datastore.ContractType(token_admin_registry.ContractType),
			token_admin_registry.Version,
			"",
		))
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
		}
		_, err = operations.ExecuteSequence(b, RegisterTokenPool, cantonChain, RegisterTokenPoolInput{
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(registryRef.Address),
			InstrumentId:                      parsedPool.InstrumentId,
			PoolInstanceID:                    string(parsedPool.InstanceId),
			CcipParty:                         string(parsedPool.CcipOwner),
			PoolOwnerParty:                    string(parsedPool.PoolOwner),
		})
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("register lock/release pool with token admin registry: %w", err)
		}

		out := ccipsequences.OnChainOutput{}
		committeeVerifierRefs := input.ExistingDataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByType(datastore.ContractType(committee_verifier.ContractType)),
			datastore.AddressRefByVersion(committee_verifier.Version),
		)
		updates := make([]lockreleasetokenpool.ChainUpdate, 0, len(input.RemoteChains))
		for remoteSelector, remoteCfg := range input.RemoteChains {
			remoteSelectorKey := strconv.FormatUint(remoteSelector, 10)
			if _, found := parsedPool.RemoteChainConfigs[remoteSelectorKey]; found {
				return out, fmt.Errorf("remote chain %d is already configured on token pool", remoteSelector)
			}

			inboundCCVs := make([]mcms.RawInstanceAddress, 0, len(remoteCfg.InboundCCVs))
			for _, inboundCCVAddress := range remoteCfg.InboundCCVs {
				var matchedRef *datastore.AddressRef
				for _, ccvRef := range committeeVerifierRefs {
					if strings.EqualFold(ccvRef.Address, inboundCCVAddress) {
						refCopy := ccvRef
						matchedRef = &refCopy

						break
					}
				}
				if matchedRef == nil {
					return out, fmt.Errorf("resolve inbound CCV ref for address %s", inboundCCVAddress)
				}
				inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(*matchedRef)
				if err != nil {
					return out, fmt.Errorf("resolve inbound CCV raw address for %s: %w", inboundCCVAddress, err)
				}
				inboundCCVs = append(inboundCCVs, inboundCCV.Binding())
			}

			outboundCCVs := make([]mcms.RawInstanceAddress, 0, len(remoteCfg.OutboundCCVs))
			for _, outboundCCVAddress := range remoteCfg.OutboundCCVs {
				var matchedRef *datastore.AddressRef
				for _, ccvRef := range committeeVerifierRefs {
					if strings.EqualFold(ccvRef.Address, outboundCCVAddress) {
						refCopy := ccvRef
						matchedRef = &refCopy

						break
					}
				}
				if matchedRef == nil {
					return out, fmt.Errorf("resolve outbound CCV ref for address %s", outboundCCVAddress)
				}
				outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(*matchedRef)
				if err != nil {
					return out, fmt.Errorf("resolve outbound CCV raw address for %s: %w", outboundCCVAddress, err)
				}
				outboundCCVs = append(outboundCCVs, outboundCCV.Binding())
			}

			outboundDefaultCfg, inboundDefaultCfg := tokenadapters.GenerateTPRLConfigs(
				remoteCfg.DefaultFinalityOutboundRateLimiterConfig,
				remoteCfg.DefaultFinalityInboundRateLimiterConfig,
				uint8(parsedPool.Decimals),
				remoteCfg.RemoteDecimals,
				"canton",
				semver.MustParse("2.0.0"),
			)
			_, inboundCustomCfg := tokenadapters.GenerateTPRLConfigs(
				remoteCfg.CustomFinalityOutboundRateLimiterConfig,
				remoteCfg.CustomFinalityInboundRateLimiterConfig,
				uint8(parsedPool.Decimals),
				remoteCfg.RemoteDecimals,
				"canton",
				semver.MustParse("2.0.0"),
			)

			outboundRef, outboundRaw, err := deployTokenPoolRateLimiter(
				b,
				cantonChain,
				input.ExistingDataStore,
				parsedPool,
				input.TokenPoolAddress,
				remoteSelectorKey,
				"outbound",
				outboundDefaultCfg,
				common.RateLimitModeRateLimitMode_DefaultFinality,
			)
			if err != nil {
				return out, fmt.Errorf("deploy outbound rate limiter for remote chain %d: %w", remoteSelector, err)
			}
			out.Addresses = append(out.Addresses, outboundRef)

			inboundRef, inboundRaw, err := deployTokenPoolRateLimiter(
				b,
				cantonChain,
				input.ExistingDataStore,
				parsedPool,
				input.TokenPoolAddress,
				remoteSelectorKey,
				"inbound",
				inboundDefaultCfg,
				common.RateLimitModeRateLimitMode_DefaultFinality,
			)
			if err != nil {
				return out, fmt.Errorf("deploy inbound rate limiter for remote chain %d: %w", remoteSelector, err)
			}
			out.Addresses = append(out.Addresses, inboundRef)

			customRef, customInboundRaw, err := deployTokenPoolRateLimiter(
				b,
				cantonChain,
				input.ExistingDataStore,
				parsedPool,
				input.TokenPoolAddress,
				remoteSelectorKey,
				"inbound-custom",
				inboundCustomCfg,
				common.RateLimitModeRateLimitMode_CustomFinality,
			)
			if err != nil {
				return out, fmt.Errorf("deploy custom inbound rate limiter for remote chain %d: %w", remoteSelector, err)
			}
			out.Addresses = append(out.Addresses, customRef)

			remoteFamily, err := chain_selectors.GetSelectorFamily(remoteSelector)
			if err != nil {
				return out, fmt.Errorf("get remote chain family for %d: %w", remoteSelector, err)
			}
			remotePoolAddress := strings.ToLower(hex.EncodeToString(remoteCfg.RemotePool))
			remoteTokenAddress := strings.ToLower(hex.EncodeToString(remoteCfg.RemoteToken))
			if remoteFamily == chain_selectors.FamilyEVM {
				remotePoolAddress = strings.TrimPrefix(strings.ToLower(gethcommon.BytesToHash(remoteCfg.RemotePool).Hex()), "0x")
				remoteTokenAddress = strings.TrimPrefix(strings.ToLower(gethcommon.BytesToAddress(remoteCfg.RemoteToken).Hex()), "0x")
			}

			updates = append(updates, lockreleasetokenpool.ChainUpdate{
				RemoteChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
				RemotePools:         []types.TEXT{types.TEXT(remotePoolAddress)},
				RemoteTokenAddress:  types.TEXT(remoteTokenAddress),
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
				if strings.Contains(err.Error(), "ApplyChainUpdates: chain already exists:") {
					return out, nil
				}
				return out, fmt.Errorf("apply remote chain updates to lock/release pool: %w", err)
			}
		}

		return out, nil
	},
)

// Canton wires rate limiters during initial remote-chain setup in ConfigureTokenForTransfers
// via ApplyChainUpdates, so this follow-up sequence is intentionally unused.
var SetTokenPoolRateLimits = operations.NewSequence(
	"canton/token-adapter/set-token-pool-rate-limits",
	semver.MustParse("2.0.0"),
	"Updates Canton lock/release pool rate limiter config for a remote chain",
	func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.TPRLRemotes) (ccipsequences.OnChainOutput, error) {
		_ = b
		_ = chains
		_ = input

		return ccipsequences.OnChainOutput{}, fmt.Errorf(
			"SetTokenPoolRateLimits is disabled: initial Canton rate limiter setup happens during ConfigureTokenForTransfers",
		)
	},
)

var DeployTokenPoolForToken = operations.NewSequence(
	"canton/token-adapter/deploy-token-pool-for-token",
	semver.MustParse("2.0.0"),
	"Deploys a Canton lock/release pool and returns both the canonical and logical datastore refs",
	func(b operations.Bundle, chains chain.BlockChains, input tokenadapters.DeployTokenPoolInput) (ccipsequences.OnChainOutput, error) {
		lockReleasePoolType := datastore.ContractType("LockReleaseTokenPool")
		if input.TokenPoolVersion == nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("TokenPoolVersion is required")
		}
		if datastore.ContractType(input.PoolType) != lockReleasePoolType {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("unsupported Canton token pool type %q", input.PoolType)
		}
		if input.ExistingDataStore == nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("existing datastore is required")
		}

		qualifier := strings.TrimSpace(input.TokenPoolQualifier)
		if qualifier == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("token pool qualifier is required")
		}
		if input.TokenRef == nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("tokenRef is required and must include instrument labels")
		}
		var instrumentAdmin, instrumentIDText string
		for _, label := range input.TokenRef.Labels.List() {
			switch {
			case strings.HasPrefix(label, "instrument-admin:"):
				instrumentAdmin = strings.TrimSpace(strings.TrimPrefix(label, "instrument-admin:"))
			case strings.HasPrefix(label, "instrument-id:"):
				instrumentIDText = strings.TrimSpace(strings.TrimPrefix(label, "instrument-id:"))
			}
		}
		if instrumentAdmin == "" || instrumentIDText == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf(
				"tokenRef labels must include instrument-admin:<party> and instrument-id:<id>",
			)
		}
		instrumentID := splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(instrumentAdmin),
			Id:    types.TEXT(instrumentIDText),
		}
		matches := input.ExistingDataStore.Addresses().Filter(
			datastore.AddressRefByType(lockReleasePoolType),
			datastore.AddressRefByChainSelector(input.ChainSelector),
			datastore.AddressRefByQualifier(qualifier),
			datastore.AddressRefByVersion(input.TokenPoolVersion),
		)
		if len(matches) > 1 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("multiple Canton token pools found with qualifier %q", qualifier)
		}
		if len(matches) == 1 {
			b.Logger.Infof("Canton token pool already deployed at %s", matches[0].Address)
			return ccipsequences.OnChainOutput{}, nil
		}

		cantonChain, ok := chains.CantonChains()[input.ChainSelector]
		if !ok || len(cantonChain.Participants) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.ChainSelector)
		}
		participant := cantonChain.Participants[0]

		resolveRefAndRaw := func(name string, contractType datastore.ContractType, version *semver.Version) (datastore.AddressRef, contracts.RawInstanceAddress, error) {
			ref, err := input.ExistingDataStore.Addresses().Get(datastore.NewAddressRefKey(
				input.ChainSelector,
				contractType,
				version,
				"",
			))
			if err != nil {
				return datastore.AddressRef{}, "", fmt.Errorf("resolve %s: %w", name, err)
			}
			raw, err := dsutils.GetRawInstanceAddressFromAddressRef(ref)
			if err != nil {
				return datastore.AddressRef{}, "", fmt.Errorf("resolve %s raw address: %w", name, err)
			}

			return ref, raw, nil
		}

		tokenAdminRegistryRef, tokenAdminRegistryRaw, err := resolveRefAndRaw("token admin registry", datastore.ContractType(token_admin_registry.ContractType), token_admin_registry.Version)
		if err != nil {
			return ccipsequences.OnChainOutput{}, err
		}
		_, rmnRemoteRaw, err := resolveRefAndRaw("rmn remote", datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version)
		if err != nil {
			return ccipsequences.OnChainOutput{}, err
		}
		_, feeQuoterRaw, err := resolveRefAndRaw("fee quoter", datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version)
		if err != nil {
			return ccipsequences.OnChainOutput{}, err
		}

		relativeHours := types.INT64(24)
		deployReport, err := operations.ExecuteOperation(b, lock_release_token_pool.Deploy, cantonChain, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
			Qualifier: ptr.String(qualifier),
			Template: lockreleasetokenpool.LockReleaseTokenPool{
				CcipOwner:               types.PARTY(participant.PartyID),
				PoolOwner:               types.PARTY(participant.PartyID),
				InstrumentId:            instrumentID,
				Decimals:                types.INT64(10),
				RemoteChainConfigs:      types.GENMAP{},
				TokenTransferFeeConfigs: types.GENMAP{},
				PoolReceiveContext: common.CCIPContext{
					Values: types.TEXTMAP{},
				},
				TransferTimeout: lockreleasetokenpool.TransferTimeout{
					RelativeHours: &relativeHours,
				},
				Deps: lockreleasetokenpool.LockReleaseTokenPoolDeps{
					TokenAdminRegistry: tokenAdminRegistryRaw.Binding(),
					RmnRemote:          rmnRemoteRaw.Binding(),
					FeeQuoter:          feeQuoterRaw.Binding(),
				},
			},
			OwnerParty: types.PARTY(participant.PartyID),
		})
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("deploy Canton lock/release pool: %w", err)
		}
		if len(deployReport.Output.Labels.List()) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("missing raw lock/release pool label in deploy output")
		}
		rawPoolAddr, err := contracts.RawInstanceAddressFromString(deployReport.Output.Labels.List()[0])
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("parse raw lock/release pool label: %w", err)
		}
		_, err = operations.ExecuteSequence(b, RegisterTokenPool, cantonChain, RegisterTokenPoolInput{
			TokenAdminRegistryInstanceAddress: contracts.HexToInstanceAddress(tokenAdminRegistryRef.Address),
			InstrumentId:                      instrumentID,
			PoolInstanceID:                    rawPoolAddr.InstanceID(),
			CcipParty:                         participant.PartyID,
			PoolOwnerParty:                    participant.PartyID,
		})
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("register Canton lock/release pool: %w", err)
		}

		logicalRef := datastore.AddressRef{
			Address:       deployReport.Output.Address,
			Labels:        deployReport.Output.Labels,
			Type:          datastore.ContractType(input.PoolType),
			Version:       input.TokenPoolVersion,
			Qualifier:     qualifier,
			ChainSelector: input.ChainSelector,
		}
		tokenAddress := strings.TrimSpace(input.TokenRef.Address)
		if tokenAddress == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("tokenRef.address is required")
		}
		tokenRef := datastore.AddressRef{
			Address:       tokenAddress,
			Type:          datastore.ContractType("Token"),
			Qualifier:     qualifier,
			ChainSelector: input.ChainSelector,
		}

		return ccipsequences.OnChainOutput{
			Addresses: []datastore.AddressRef{
				logicalRef,
				deployReport.Output,
				tokenRef,
			},
		}, nil
	},
)

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

func deployTokenPoolRateLimiter(
	b operations.Bundle,
	cantonChain cldfcanton.Chain,
	existingDataStore datastore.DataStore,
	parsedPool *lockreleasetokenpool.LockReleaseTokenPool,
	tokenPoolAddress string,
	remoteSelectorKey string,
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

	qualifier := fmt.Sprintf("%s-%s-%s", tokenPoolAddress, direction, remoteSelectorKey)
	report, err := operations.ExecuteOperation(b, deployOp, cantonChain, contract.DeployInput[common.RateLimiter]{
		Qualifier: &qualifier,
		Template: common.RateLimiter{
			PoolInstanceId:      parsedPool.InstanceId,
			PoolOwner:           parsedPool.PoolOwner,
			RemoteChainSelector: types.NUMERIC(remoteSelectorKey),
			Direction:           dirEnum,
			Mode:                mode,
			IsEnabled:           types.BOOL(cfg.IsEnabled),
			Capacity:            capacity,
			Rate:                rate,
			Tokens:              capacity,
			LastUpdated:         types.TIMESTAMP(time.Now()),
		},
		OwnerParty: parsedPool.PoolOwner,
	})
	if err != nil {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, err
	}
	addresses, ok := existingDataStore.Addresses().(datastore.MutableAddressRefStore)
	if !ok {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, fmt.Errorf("existing datastore addresses are not mutable")
	}
	if err := addresses.Add(report.Output); err != nil {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, err
	}

	rawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(report.Output)
	if err != nil {
		return datastore.AddressRef{}, mcms.RawInstanceAddress{}, err
	}

	return report.Output, rawAddr.Binding(), nil
}
