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
			participant.PartyID,
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

// ResolveVerifierAddress implements [adapters.AggregatorConfigAdapter].
func (a *CantonAggregatorConfigAdapter) ResolveVerifierAddress(ds datastore.DataStore, chainSelector uint64, qualifier string) (string, error) {
	return dsutils.FindAndFormatFirstRef(ds, chainSelector,
		func(r datastore.AddressRef) (string, error) { return r.Address, nil },
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
