package adapters

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/versioned_verifier_resolver"
	dsutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccv/deployment/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
)

type CantonAggregatorConfigAdapter struct{}

var _ adapters.AggregatorConfigAdapter = (*CantonAggregatorConfigAdapter)(nil)
var _ adapters.CommitteeVerifierOnchainAdapter = (*CantonCommitteeVerifierOnchain)(nil)

type CantonCommitteeVerifierOnchain struct{}

// ApplySignatureConfigs implements [adapters.CommitteeVerifierOnchainAdapter].
func (a *CantonCommitteeVerifierOnchain) ApplySignatureConfigs(ctx context.Context, env deployment.Environment, destChainSelector uint64, qualifier string, change adapters.SignatureConfigChange) error {
	panic("unimplemented")
}

// ScanCommitteeStates implements [adapters.CommitteeVerifierOnchainAdapter].
func (a *CantonCommitteeVerifierOnchain) ScanCommitteeStates(ctx context.Context, env deployment.Environment, chainSelector uint64) ([]*adapters.CommitteeState, error) {
	refs := env.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(committee_verifier.ContractType)),
		datastore.AddressRefByChainSelector(chainSelector),
	)

	if len(refs) == 0 {
		return nil, nil
	}

	cantonChains := env.BlockChains.CantonChains()
	if cantonChains == nil {
		return nil, fmt.Errorf("no canton chains found in environment")
	}

	chain, ok := cantonChains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("canton chain %d not found in environment", chainSelector)
	}
	if len(chain.Participants) == 0 {
		return nil, fmt.Errorf("no participants configured for canton chain %d", chainSelector)
	}

	participant := chain.Participants[0]

	states := make([]*adapters.CommitteeState, 0, len(refs))
	for _, ref := range refs {
		instanceAddr := contracts.HexToInstanceAddress(ref.Address)
		active, err := contract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			contract.LedgerQueryParties(participant),
			ccvs.CommitteeVerifier{}.GetTemplateID(),
			instanceAddr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find CommitteeVerifier %s on chain %d: %w", ref.Address, chainSelector, err)
		}

		cv, err := bindings.UnmarshalCreatedEvent[ccvs.CommitteeVerifier](active.GetCreatedEvent())
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal CommitteeVerifier %s on chain %d: %w", ref.Address, chainSelector, err)
		}

		sigConfigs, err := signatureConfigsFromCommitteeVerifier(cv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse signature configs from CommitteeVerifier %s on chain %d: %w", ref.Address, chainSelector, err)
		}

		states = append(states, &adapters.CommitteeState{
			Qualifier:        ref.Qualifier,
			ChainSelector:    chainSelector,
			Address:          ref.Address,
			SignatureConfigs: sigConfigs,
		})
	}

	return states, nil
}

// ResolveDestinationVerifierAddress implements [adapters.AggregatorConfigAdapter].
//
// For Canton we want to resolve the aggregator destination verifier address to
// the raw instance address of the committee verifier so that users
// see the raw address in the indexer's verifier results for a specific message.
// This appears when users send a message from X -> Canton.
func (a *CantonAggregatorConfigAdapter) ResolveDestinationVerifierAddress(ds datastore.DataStore, chainSelector uint64, qualifier string) (string, error) {
	return dsutils.FindAndFormatFirstRef(ds, chainSelector,
		func(r datastore.AddressRef) (string, error) {
			// The aggregator sends the raw destination instance address to the indexer so that users
			// can more easily execute messages on Canton.
			labels := r.Labels.List()
			if len(labels) == 0 {
				// Shouldn't happen, graceful fallback.
				return r.Address, nil
			}

			// labels[0] is the raw instance address, but since the aggregator requires hex-encoded addresses,
			// we need to hex-encode the string bytes before returning it.
			hexEncoded := hex.EncodeToString([]byte(labels[0]))

			return "0x" + hexEncoded, nil
		},
		datastore.AddressRef{
			Type:      datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
			Qualifier: qualifier,
		},
		datastore.AddressRef{
			Type:      datastore.ContractType(committee_verifier.ContractType),
			Qualifier: qualifier,
		},
	)
}

// ResolveSourceVerifierAddress implements [adapters.AggregatorConfigAdapter].
//
// For Canton we want to resolve the aggregator source verifier address to
// the hashed instance address of the committee verifier. This is because
// the CCIPMessageSentEvent onchain emits the hashed instance address,
// and the aggregator checks that the source verifier address it is configured
// with is present inside the message's CCV addresses.
// https://github.com/smartcontractkit/chainlink-ccv/blob/33b00afa3061efc6faca53e5cfb9658e4ad9d6e1/aggregator/pkg/quorum/evm_quorum_validator.go#L49-L51.
func (a *CantonAggregatorConfigAdapter) ResolveSourceVerifierAddress(ds datastore.DataStore, chainSelector uint64, qualifier string) (string, error) {
	return dsutils.FindAndFormatFirstRef(ds, chainSelector,
		func(r datastore.AddressRef) (string, error) { return r.Address, nil },
		// TODO: below is searching for the EVM verifier resolver, why?
		datastore.AddressRef{
			Type:      datastore.ContractType(versioned_verifier_resolver.CommitteeVerifierResolverType),
			Qualifier: qualifier,
		},
		datastore.AddressRef{
			Type:      datastore.ContractType(committee_verifier.ContractType),
			Qualifier: qualifier,
		},
	)
}

// GetDeployedChains implements [adapters.AggregatorConfigAdapter].
func (a *CantonAggregatorConfigAdapter) GetDeployedChains(ds datastore.DataStore, qualifier string) []uint64 {
	if ds == nil {
		return nil
	}
	refs := ds.Addresses().Filter(
		datastore.AddressRefByQualifier(qualifier),
		datastore.AddressRefByType(datastore.ContractType(committee_verifier.ContractType)),
		datastore.AddressRefByVersion(committee_verifier.Version),
	)
	seen := make(map[uint64]struct{}, len(refs))
	chains := make([]uint64, 0, len(refs))
	for _, ref := range refs {
		if _, exists := seen[ref.ChainSelector]; exists {
			continue
		}
		seen[ref.ChainSelector] = struct{}{}
		chains = append(chains, ref.ChainSelector)
	}

	return chains
}

func signatureConfigsFromCommitteeVerifier(cv *ccvs.CommitteeVerifier) ([]adapters.SignatureConfig, error) {
	if cv == nil || len(cv.SignerConfigs) == 0 {
		return nil, nil
	}

	out := make([]adapters.SignatureConfig, 0, len(cv.SignerConfigs))
	for _, v := range cv.SignerConfigs {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("marshal signer config value: %w", err)
		}
		var cfg ccvs.SignatureConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal signer config: %w", err)
		}

		sourceSel, err := numericToUint64(cfg.SourceChainSelector)
		if err != nil {
			return nil, fmt.Errorf("parse source chain selector %q: %w", cfg.SourceChainSelector, err)
		}

		signers := make([]string, 0, len(cfg.SignerKeys))
		for _, key := range cfg.SignerKeys {
			signer, err := signerKeyToAddress(string(key))
			if err != nil {
				return nil, fmt.Errorf("normalize signer key %q: %w", key, err)
			}
			signers = append(signers, signer)
		}

		thr := cfg.Threshold
		if thr < 0 || thr > 255 {
			return nil, fmt.Errorf("threshold %d out of range for uint8", thr)
		}
		thrU8 := uint8(thr)

		out = append(out, adapters.SignatureConfig{
			SourceChainSelector: sourceSel,
			Signers:             signers,
			Threshold:           thrU8,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].SourceChainSelector < out[j].SourceChainSelector
	})

	return out, nil
}

func numericToUint64(n types.NUMERIC) (uint64, error) {
	s := strings.TrimSpace(string(n))
	return internalparse.Uint64Checked(s)
}

func normalizeSignerHex(hexKey string) string {
	hexKey = strings.TrimPrefix(strings.TrimSpace(hexKey), "0x")
	return "0x" + hexKey
}

// signerKeyToAddress converts a Canton signer key into the address form expected
// by the aggregator quorum config. Canton stores raw 65-byte secp256k1 public
// keys, while the aggregator matches recovered signer addresses from verifier
// signatures against configured signer addresses.
func signerKeyToAddress(hexKey string) (string, error) {
	normalized := normalizeSignerHex(hexKey)
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(normalized, "0x"))
	if err != nil {
		return "", fmt.Errorf("decode hex: %w", err)
	}
	if len(keyBytes) != 65 {
		return "", fmt.Errorf("expected 65-byte uncompressed pubkey, got %d bytes", len(keyBytes))
	}

	pubKey, err := gethcrypto.UnmarshalPubkey(keyBytes)
	if err != nil {
		return "", fmt.Errorf("unmarshal pubkey: %w", err)
	}

	return gethcrypto.PubkeyToAddress(*pubKey).Hex(), nil
}
