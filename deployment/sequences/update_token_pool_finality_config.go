package sequences

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	tokenadaptersfinality "github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/registry_burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// ValidateCantonFinalityConfig checks that exactly one finality mode is configured. The Canton
// finality config is a single-variant value, so combinations cannot be expressed on-ledger and
// would be silently narrowed by toCantonFinalityConfig.
func ValidateCantonFinalityConfig(cfg tokenadaptersfinality.Config) error {
	modes := 0
	if cfg.WaitForFinality {
		modes++
	}
	if cfg.WaitForSafe {
		modes++
	}
	if cfg.BlockDepth > 0 {
		modes++
	}
	if modes != 1 {
		return fmt.Errorf(
			"exactly one finality mode must be set (waitForFinality, waitForSafe, or blockDepth); got waitForFinality=%t waitForSafe=%t blockDepth=%d",
			cfg.WaitForFinality, cfg.WaitForSafe, cfg.BlockDepth,
		)
	}

	return nil
}

// ResolveCantonTokenPoolTypeFromRef maps a token pool datastore ref type to the logical pool type
// used to dispatch pool choices. Both logical names (e.g. "BurnMintTokenPool") and canonical
// datastore names (e.g. "CantonBurnMintTokenPool") are accepted. An empty type defaults to the
// lock/release pool, matching resolveCantonTokenPoolType.
func ResolveCantonTokenPoolTypeFromRef(poolType datastore.ContractType) (datastore.ContractType, error) {
	trimmed := strings.TrimSpace(poolType.String())
	switch datastore.ContractType(trimmed) {
	case "":
		return lockReleasePoolType, nil
	case lockReleasePoolType, datastore.ContractType(lock_release_token_pool.ContractType):
		return lockReleasePoolType, nil
	case burnMintPoolType, datastore.ContractType(burn_mint_token_pool.ContractType):
		return burnMintPoolType, nil
	case registryBurnMintPoolType, datastore.ContractType(registry_burn_mint_token_pool.ContractType):
		return registryBurnMintPoolType, nil
	default:
		return "", fmt.Errorf("unsupported Canton token pool type %q", poolType)
	}
}

// laneFinalityUpdate is one lane's preserved config with the finality config replaced.
type laneFinalityUpdate struct {
	// remoteChainSelectorKey is the ledger Numeric key of the existing lane config, reused for
	// both removal and re-addition in the ApplyChainUpdates exercise.
	remoteChainSelectorKey                     types.NUMERIC
	remotePools                                []types.TEXT
	remoteTokenAddress                         types.TEXT
	inboundCCVs                                []chainlinkapi.RawInstanceAddress
	outboundCCVs                               []chainlinkapi.RawInstanceAddress
	finalityConfig                             ccipcodec.FinalityConfig
	inboundRateLimiter                         chainlinkapi.RawInstanceAddress
	inboundCustomBlockConfirmationsRateLimiter chainlinkapi.RawInstanceAddress
	outboundRateLimiter                        chainlinkapi.RawInstanceAddress
}

type UpdateTokenPoolFinalityInput struct {
	CantonChainSelector uint64
	// TokenPoolRef points at the deployed Canton token pool whose lane finality configs are updated.
	TokenPoolRef datastore.AddressRef
	// RemoteChainSelectors lists the remote chains whose finality config is updated.
	RemoteChainSelectors []uint64
	// AllowedFinalityConfig is the finality config set on every listed remote chain.
	AllowedFinalityConfig tokenadaptersfinality.Config
}

var UpdateTokenPoolFinalityConfig = operations.NewSequence(
	"canton/token-adapter/update-token-pool-finality-config",
	semver.MustParse("2.0.0"),
	"Updates the finality config of configured remote chains on a Canton token pool",
	func(b operations.Bundle, chains chain.BlockChains, input UpdateTokenPoolFinalityInput) (ccipsequences.OnChainOutput, error) {
		if len(input.RemoteChainSelectors) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("at least one remote chain selector is required")
		}
		if input.TokenPoolRef.Address == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("token pool ref address is required")
		}
		if err := ValidateCantonFinalityConfig(input.AllowedFinalityConfig); err != nil {
			return ccipsequences.OnChainOutput{}, err
		}

		cantonChain, ok := chains.CantonChains()[input.CantonChainSelector]
		if !ok || len(cantonChain.Participants) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found", input.CantonChainSelector)
		}
		participant := cantonChain.Participants[0]
		// ReadAsPartyIDs are CanReadAs rights for parties the operator cannot ActAs (e.g. ccip owner).
		// When present, exercises must be encoded as MCMS proposals instead of submitted directly.
		mcmsEnabled := len(participant.ReadAsPartyIDs) > 0

		logicalPoolType, err := ResolveCantonTokenPoolTypeFromRef(input.TokenPoolRef.Type)
		if err != nil {
			return ccipsequences.OnChainOutput{}, err
		}

		poolAddress := contracts.HexToInstanceAddress(input.TokenPoolRef.Address)
		newFinality := toCantonFinalityConfig(input.AllowedFinalityConfig)
		out := ccipsequences.OnChainOutput{}

		switch logicalPoolType {
		case lockReleasePoolType:
			pool, err := findActiveCantonContract[lockreleasetokenpool.LockReleaseTokenPool](
				b.GetContext(), participant, lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(), poolAddress,
			)
			if err != nil {
				return out, err
			}
			updates, skipped, err := buildLaneFinalityUpdates(lockReleaseLaneConfigs(pool.RemoteChainConfigs), input.RemoteChainSelectors, newFinality)
			if err != nil {
				return out, err
			}
			if len(updates) == 0 {
				b.Logger.Infof("all %d requested lanes on pool %s already allow the requested finality; nothing to update", skipped, poolAddress.Hex())

				return out, nil
			}
			poolRaw := contracts.NewRawInstanceAddress(contracts.InstanceID(pool.InstanceId), pool.PoolOwner)
			chainUpdates := make([]lockreleasetokenpool.ChainUpdate, 0, len(updates))
			for _, update := range updates {
				chainUpdates = append(chainUpdates, lockreleasetokenpool.ChainUpdate{
					RemoteChainSelector: update.remoteChainSelectorKey,
					RemotePools:         update.remotePools,
					RemoteTokenAddress:  update.remoteTokenAddress,
					InboundCCVs:         update.inboundCCVs,
					OutboundCCVs:        update.outboundCCVs,
					FinalityConfig:      update.finalityConfig,
					InboundRateLimiter:  update.inboundRateLimiter,
					InboundCustomBlockConfirmationsRateLimiter: update.inboundCustomBlockConfirmationsRateLimiter,
					OutboundRateLimiter:                        update.outboundRateLimiter,
				})
			}
			out.BatchOps, err = executePoolExercise(b, lock_release_token_pool.ApplyChainUpdates, cantonChain, contract.ChoiceInput[lockreleasetokenpool.ApplyChainUpdates]{
				InstanceAddress:    poolRaw.InstanceAddress(),
				RawInstanceAddress: poolRaw.String(),
				MCMSEnabled:        mcmsEnabled,
				Args: lockreleasetokenpool.ApplyChainUpdates{
					RemoteChainSelectorsToRemove: laneSelectorKeys(updates),
					ChainsToAdd:                  chainUpdates,
				},
			}, mcmsEnabled, out.BatchOps)
			if err != nil {
				return out, fmt.Errorf("apply finality config update to lock/release pool: %w", err)
			}
		case burnMintPoolType:
			pool, err := findActiveCantonContract[burnminttokenpool.BurnMintTokenPool](
				b.GetContext(), participant, burnminttokenpool.BurnMintTokenPool{}.GetTemplateID(), poolAddress,
			)
			if err != nil {
				return out, err
			}
			updates, skipped, err := buildLaneFinalityUpdates(burnMintLaneConfigs(pool.RemoteChainConfigs), input.RemoteChainSelectors, newFinality)
			if err != nil {
				return out, err
			}
			if len(updates) == 0 {
				b.Logger.Infof("all %d requested lanes on pool %s already allow the requested finality; nothing to update", skipped, poolAddress.Hex())

				return out, nil
			}
			poolRaw := contracts.NewRawInstanceAddress(contracts.InstanceID(pool.InstanceId), pool.PoolOwner)
			chainUpdates := make([]burnminttokenpool.ChainUpdate, 0, len(updates))
			for _, update := range updates {
				chainUpdates = append(chainUpdates, burnminttokenpool.ChainUpdate{
					RemoteChainSelector: update.remoteChainSelectorKey,
					RemotePools:         update.remotePools,
					RemoteTokenAddress:  update.remoteTokenAddress,
					InboundCCVs:         update.inboundCCVs,
					OutboundCCVs:        update.outboundCCVs,
					FinalityConfig:      update.finalityConfig,
					InboundRateLimiter:  update.inboundRateLimiter,
					InboundCustomBlockConfirmationsRateLimiter: update.inboundCustomBlockConfirmationsRateLimiter,
					OutboundRateLimiter:                        update.outboundRateLimiter,
				})
			}
			out.BatchOps, err = executePoolExercise(b, burn_mint_token_pool.ApplyChainUpdates, cantonChain, contract.ChoiceInput[burnminttokenpool.ApplyChainUpdates]{
				InstanceAddress:    poolRaw.InstanceAddress(),
				RawInstanceAddress: poolRaw.String(),
				MCMSEnabled:        mcmsEnabled,
				Args: burnminttokenpool.ApplyChainUpdates{
					RemoteChainSelectorsToRemove: laneSelectorKeys(updates),
					ChainsToAdd:                  chainUpdates,
				},
			}, mcmsEnabled, out.BatchOps)
			if err != nil {
				return out, fmt.Errorf("apply finality config update to burn/mint pool: %w", err)
			}
		case registryBurnMintPoolType:
			pool, err := findActiveCantonContract[registrybnm.BurnMintTokenPool](
				b.GetContext(), participant, registrybnm.BurnMintTokenPool{}.GetTemplateID(), poolAddress,
			)
			if err != nil {
				return out, err
			}
			updates, skipped, err := buildLaneFinalityUpdates(registryLaneConfigs(pool.RemoteChainConfigs), input.RemoteChainSelectors, newFinality)
			if err != nil {
				return out, err
			}
			if len(updates) == 0 {
				b.Logger.Infof("all %d requested lanes on pool %s already allow the requested finality; nothing to update", skipped, poolAddress.Hex())

				return out, nil
			}
			poolRaw := contracts.NewRawInstanceAddress(contracts.InstanceID(pool.InstanceId), pool.PoolOwner)
			chainUpdates := make([]registrybnm.ChainUpdate, 0, len(updates))
			for _, update := range updates {
				chainUpdates = append(chainUpdates, registrybnm.ChainUpdate{
					RemoteChainSelector: update.remoteChainSelectorKey,
					RemotePools:         update.remotePools,
					RemoteTokenAddress:  update.remoteTokenAddress,
					InboundCCVs:         update.inboundCCVs,
					OutboundCCVs:        update.outboundCCVs,
					FinalityConfig:      update.finalityConfig,
					InboundRateLimiter:  update.inboundRateLimiter,
					InboundCustomBlockConfirmationsRateLimiter: update.inboundCustomBlockConfirmationsRateLimiter,
					OutboundRateLimiter:                        update.outboundRateLimiter,
				})
			}
			out.BatchOps, err = executePoolExercise(b, registry_burn_mint_token_pool.ApplyChainUpdates, cantonChain, contract.ChoiceInput[registrybnm.ApplyChainUpdates]{
				InstanceAddress:    poolRaw.InstanceAddress(),
				RawInstanceAddress: poolRaw.String(),
				MCMSEnabled:        mcmsEnabled,
				Args: registrybnm.ApplyChainUpdates{
					RemoteChainSelectorsToRemove: laneSelectorKeys(updates),
					ChainsToAdd:                  chainUpdates,
				},
			}, mcmsEnabled, out.BatchOps)
			if err != nil {
				return out, fmt.Errorf("apply finality config update to registry burn/mint pool: %w", err)
			}
		default:
			return out, fmt.Errorf("unsupported Canton token pool type %q", logicalPoolType)
		}

		return out, nil
	},
)

// executePoolExercise submits a pool choice exercise and, in proposal mode, converts the encoded
// exercise into an MCMS batch operation.
func executePoolExercise[ARGS any](
	b operations.Bundle,
	op *operations.Operation[contract.ChoiceInput[ARGS], contract.ExerciseOutput, cldfcanton.Chain],
	cantonChain cldfcanton.Chain,
	input contract.ChoiceInput[ARGS],
	mcmsEnabled bool,
	batchOps []mcms_types.BatchOperation,
) ([]mcms_types.BatchOperation, error) {
	report, err := operations.ExecuteOperation(b, op, cantonChain, input)
	if err != nil {
		return batchOps, err
	}
	if !mcmsEnabled || report.Output.Executed() {
		return batchOps, nil
	}

	batchOp, err := contract.NewBatchOperationFromExercises([]contract.ExerciseOutput{report.Output})
	if err != nil {
		return batchOps, fmt.Errorf("build MCMS batch: %w", err)
	}

	return append(batchOps, batchOp), nil
}

// buildLaneFinalityUpdates reads the pool's current lane configs and, for each requested selector,
// preserves every field except the finality config. Lanes already set to the requested finality
// are skipped.
func buildLaneFinalityUpdates(
	remoteChainConfigs map[types.NUMERIC]tokenPoolChainUpdate,
	remoteChainSelectors []uint64,
	newFinality ccipcodec.FinalityConfig,
) ([]laneFinalityUpdate, int, error) {
	updates := make([]laneFinalityUpdate, 0, len(remoteChainSelectors))
	skipped := 0
	for _, selector := range remoteChainSelectors {
		key, cfg, ok := matchRemoteChainConfig(remoteChainConfigs, selector)
		if !ok {
			return nil, 0, fmt.Errorf("remote chain %d is not configured on token pool", selector)
		}
		if finalityConfigEqual(cfg.FinalityConfig, newFinality) {
			skipped++

			continue
		}
		if err := validateLaneCustomFinalityRateLimiter(
			selector, cfg.InboundRateLimiter, cfg.InboundCustomBlockConfirmationsRateLimiter, cfg.OutboundRateLimiter, newFinality,
		); err != nil {
			return nil, 0, err
		}
		cfg.FinalityConfig = newFinality
		updates = append(updates, laneFinalityUpdate{
			remoteChainSelectorKey: key,
			remotePools:            cfg.RemotePools,
			remoteTokenAddress:     cfg.RemoteTokenAddress,
			inboundCCVs:            cfg.InboundCCVs,
			outboundCCVs:           cfg.OutboundCCVs,
			finalityConfig:         cfg.FinalityConfig,
			inboundRateLimiter:     cfg.InboundRateLimiter,
			inboundCustomBlockConfirmationsRateLimiter: cfg.InboundCustomBlockConfirmationsRateLimiter,
			outboundRateLimiter:                        cfg.OutboundRateLimiter,
		})
	}

	return updates, skipped, nil
}

// matchRemoteChainConfig returns the ledger Numeric key and config for a bare chain selector,
// e.g. key "16015285229680338985." for selector 16015285229680338985.
func matchRemoteChainConfig(configs map[types.NUMERIC]tokenPoolChainUpdate, selector uint64) (types.NUMERIC, tokenPoolChainUpdate, bool) {
	selectorStr := strconv.FormatUint(selector, 10)
	for key, cfg := range configs {
		if whole, _, _ := strings.Cut(string(key), "."); whole == selectorStr {
			return key, cfg, true
		}
	}

	return "", tokenPoolChainUpdate{}, false
}

// lockReleaseLaneConfigs, burnMintLaneConfigs and registryLaneConfigs convert the pool-type
// specific lane config maps into the pool-agnostic shape; the on-ledger structs are field-for-field
// identical apart from the selector key.
func lockReleaseLaneConfigs(configs map[types.NUMERIC]lockreleasetokenpool.RemoteChainConfig) map[types.NUMERIC]tokenPoolChainUpdate {
	out := make(map[types.NUMERIC]tokenPoolChainUpdate, len(configs))
	for key, cfg := range configs {
		out[key] = tokenPoolChainUpdate{
			RemotePools:        cfg.RemotePools,
			RemoteTokenAddress: cfg.RemoteTokenAddress,
			InboundCCVs:        cfg.InboundCCVs,
			OutboundCCVs:       cfg.OutboundCCVs,
			FinalityConfig:     cfg.FinalityConfig,
			InboundRateLimiter: cfg.InboundRateLimiter,
			InboundCustomBlockConfirmationsRateLimiter: cfg.InboundCustomBlockConfirmationsRateLimiter,
			OutboundRateLimiter:                        cfg.OutboundRateLimiter,
		}
	}

	return out
}

func burnMintLaneConfigs(configs map[types.NUMERIC]burnminttokenpool.RemoteChainConfig) map[types.NUMERIC]tokenPoolChainUpdate {
	out := make(map[types.NUMERIC]tokenPoolChainUpdate, len(configs))
	for key, cfg := range configs {
		out[key] = tokenPoolChainUpdate{
			RemotePools:        cfg.RemotePools,
			RemoteTokenAddress: cfg.RemoteTokenAddress,
			InboundCCVs:        cfg.InboundCCVs,
			OutboundCCVs:       cfg.OutboundCCVs,
			FinalityConfig:     cfg.FinalityConfig,
			InboundRateLimiter: cfg.InboundRateLimiter,
			InboundCustomBlockConfirmationsRateLimiter: cfg.InboundCustomBlockConfirmationsRateLimiter,
			OutboundRateLimiter:                        cfg.OutboundRateLimiter,
		}
	}

	return out
}

func registryLaneConfigs(configs map[types.NUMERIC]registrybnm.RemoteChainConfig) map[types.NUMERIC]tokenPoolChainUpdate {
	out := make(map[types.NUMERIC]tokenPoolChainUpdate, len(configs))
	for key, cfg := range configs {
		out[key] = tokenPoolChainUpdate{
			RemotePools:        cfg.RemotePools,
			RemoteTokenAddress: cfg.RemoteTokenAddress,
			InboundCCVs:        cfg.InboundCCVs,
			OutboundCCVs:       cfg.OutboundCCVs,
			FinalityConfig:     cfg.FinalityConfig,
			InboundRateLimiter: cfg.InboundRateLimiter,
			InboundCustomBlockConfirmationsRateLimiter: cfg.InboundCustomBlockConfirmationsRateLimiter,
			OutboundRateLimiter:                        cfg.OutboundRateLimiter,
		}
	}

	return out
}

func laneSelectorKeys(updates []laneFinalityUpdate) []types.NUMERIC {
	keys := make([]types.NUMERIC, 0, len(updates))
	for _, update := range updates {
		keys = append(keys, update.remoteChainSelectorKey)
	}

	return keys
}

// findActiveCantonContract finds an active contract by instance address and unmarshals its create
// arguments into the generated bindings type T.
func findActiveCantonContract[T model.CreateCommander](
	ctx context.Context,
	participant cldfcanton.Participant,
	templateID string,
	instanceAddress contracts.InstanceAddress,
) (T, error) {
	active, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		templateID,
		instanceAddress,
	)
	if err != nil {
		var zero T

		return zero, fmt.Errorf("find active contract %s at %s: %w", templateID, instanceAddress.Hex(), err)
	}

	parsed, err := bindings.UnmarshalCreatedEvent[T](active.GetCreatedEvent())
	if err != nil {
		var zero T

		return zero, fmt.Errorf("parse active contract %s at %s: %w", templateID, instanceAddress.Hex(), err)
	}

	return *parsed, nil
}

// finalityConfigEqual reports whether two Canton finality configs are identical.
func finalityConfigEqual(a, b ccipcodec.FinalityConfig) bool {
	return reflect.DeepEqual(a, b)
}

// customFinalityRateLimiterRequired mirrors CCIP.CodecV2.FinalityConfig.customFinalityRateLimiterRequired:
// any config other than finalized-only uses the custom-finality inbound limiter path.
func customFinalityRateLimiterRequired(finality ccipcodec.FinalityConfig) bool {
	return finality.GetVariantTag() != "WaitForFinality"
}

// validateLaneCustomFinalityRateLimiter checks that a lane's current custom-finality inbound rate
// limiter address is usable when the new finality config is faster than finality. The DAML
// ApplyChainUpdates choice only validates this address when the finality config requires it, so a
// lane set up with WaitForFinality may carry an empty or duplicate address.
func validateLaneCustomFinalityRateLimiter(
	selector uint64,
	inboundRateLimiter, inboundCustomRateLimiter, outboundRateLimiter chainlinkapi.RawInstanceAddress,
	newFinality ccipcodec.FinalityConfig,
) error {
	if !customFinalityRateLimiterRequired(newFinality) {
		return nil
	}

	custom := string(inboundCustomRateLimiter.Unpack)
	if custom == "" {
		return fmt.Errorf(
			"remote chain %d: custom finality inbound rate limiter is not configured on the pool; deploy one before allowing faster-than-finality transfers",
			selector,
		)
	}
	if custom == string(inboundRateLimiter.Unpack) || custom == string(outboundRateLimiter.Unpack) {
		return fmt.Errorf(
			"remote chain %d: custom finality inbound rate limiter must be distinct from the default inbound and outbound rate limiters",
			selector,
		)
	}

	return nil
}
