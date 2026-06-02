package tests

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v1_5_0/operations/burn_mint_erc20_with_drip"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
)

const (
	envLoadTokenQualifier = "CANTON_LOAD_TOKEN_QUALIFIER"

	// Default transfer amounts (lane identity is discovery-driven; amounts are overridable later via env).
	CantonToEVMFeeAmount           int64 = 2_000
	CantonToEVMTokenTransferAmount int64 = 1_000
	EVMToCantonTransferAmount      int64 = 100_000_000_000
	EVMDecimalsScale               int64 = 1_000_000_000_000_000_000
)

// TokenLane describes a discovered token pool pairing for load or e2e tests.
type TokenLane struct {
	Combo             devenvcommon.TokenCombination
	TransferAmount    *big.Int
	ExecutionGasLimit uint32
	FinalityConfig    protocol.Finality
	SrcToken          protocol.UnknownAddress // EVM source only; empty for Canton source
	DestToken         protocol.UnknownAddress // EVM dest for balance checks
}

// TokenLaneOptions configures lane selection.
type TokenLaneOptions struct {
	Qualifier string // env CANTON_LOAD_TOKEN_QUALIFIER when empty
}

// DiscoverTokenLanes returns token lanes whose pools exist for srcSelector (local) → destSelector (remote).
func DiscoverTokenLanes(
	t *testing.T,
	in *ccv.Cfg,
	lib ccv.Lib,
	srcSelector, destSelector uint64,
) ([]TokenLane, error) {
	t.Helper()

	ds, err := lib.DataStore()
	if err != nil {
		return nil, fmt.Errorf("datastore: %w", err)
	}

	combos := devenvcommon.FilterTokenCombinations(
		devenvcommon.AllTokenCombinations(),
		in.EnvironmentTopology,
		ds,
		[]uint64{srcSelector, destSelector},
	)
	if len(combos) == 0 {
		return nil, fmt.Errorf("no token combinations for selectors %d -> %d", srcSelector, destSelector)
	}

	cantonSource := isCantonSelector(in, srcSelector)
	lanes := make([]TokenLane, 0, len(combos))
	for _, combo := range combos {
		lane, err := tokenLaneFromCombo(t, ds, combo, srcSelector, destSelector, cantonSource)
		if err != nil {
			return nil, err
		}
		lanes = append(lanes, lane)
	}

	return lanes, nil
}

// SelectTokenLane picks one lane from discovery using env override, then devenv default preference.
func SelectTokenLane(t *testing.T, lanes []TokenLane, opts TokenLaneOptions) TokenLane {
	t.Helper()

	if len(lanes) == 0 {
		t.Fatal("SelectTokenLane: no discovered token lanes")
	}

	qualifier := opts.Qualifier
	if qualifier == "" {
		qualifier = os.Getenv(envLoadTokenQualifier)
	}

	if qualifier != "" {
		matches := filterLanesByQualifier(lanes, qualifier)
		switch len(matches) {
		case 1:
			return matches[0]
		case 0:
			t.Fatalf("CANTON_LOAD_TOKEN_QUALIFIER=%q matched no lanes; discovered: %s",
				qualifier, formatLaneList(lanes))
		default:
			t.Fatalf("CANTON_LOAD_TOKEN_QUALIFIER=%q matched multiple lanes: %s",
				qualifier, formatLaneList(matches))
		}
	}

	defaults := filterDefaultDevenvLanes(lanes)
	switch len(defaults) {
	case 1:
		return defaults[0]
	case 0:
		if len(lanes) == 1 {
			return lanes[0]
		}
		t.Fatalf("ambiguous token lanes (set %s): %s", envLoadTokenQualifier, formatLaneList(lanes))
	default:
		t.Fatalf("multiple default BM 2.0 <-> LR 2.0 lanes (set %s): %s",
			envLoadTokenQualifier, formatLaneList(defaults))
	}

	return TokenLane{} // unreachable; satisfies compiler
}

// DefaultDevenvTokenLane discovers and selects the default lane for src -> dest.
func DefaultDevenvTokenLane(t *testing.T, in *ccv.Cfg, lib ccv.Lib, srcSelector, destSelector uint64) TokenLane {
	t.Helper()

	lanes, err := DiscoverTokenLanes(t, in, lib, srcSelector, destSelector)
	require.NoError(t, err)

	return SelectTokenLane(t, lanes, TokenLaneOptions{})
}

func tokenLaneFromCombo(
	t *testing.T,
	ds datastore.DataStore,
	combo devenvcommon.TokenCombination,
	srcSelector, destSelector uint64,
	cantonSource bool,
) (TokenLane, error) {
	t.Helper()

	var transferAmount int64
	var finality protocol.Finality
	var gasLimit uint64

	if cantonSource {
		transferAmount = CantonToEVMTokenTransferAmount
		finality = 1
		gasLimit = 500_000
	} else {
		transferAmount = EVMToCantonTransferAmount
		finality = 0
		gasLimit = 200_000
	}

	lane := TokenLane{
		Combo:             combo,
		TransferAmount:    big.NewInt(transferAmount),
		ExecutionGasLimit: uint32(gasLimit),
		FinalityConfig:    finality,
	}

	if cantonSource {
		destToken, err := resolveBurnMintToken(ds, destSelector, combo.RemotePoolAddressRef().Qualifier)
		if err != nil {
			return TokenLane{}, fmt.Errorf("resolve dest token on %d: %w", destSelector, err)
		}
		lane.DestToken = destToken
	} else {
		srcToken, err := resolveBurnMintToken(ds, srcSelector, combo.LocalPoolAddressRef().Qualifier)
		if err != nil {
			return TokenLane{}, fmt.Errorf("resolve src token on %d: %w", srcSelector, err)
		}
		lane.SrcToken = srcToken
	}

	return lane, nil
}

func resolveBurnMintToken(ds datastore.DataStore, chainSelector uint64, qualifier string) (protocol.UnknownAddress, error) {
	return tcapi.GetContractAddress(
		ds,
		chainSelector,
		datastore.ContractType(burn_mint_erc20_with_drip.ContractType),
		burn_mint_erc20_with_drip.Deploy.Version(),
		qualifier,
		"burn mint erc677",
	)
}

func isCantonSelector(cfg *ccv.Cfg, selector uint64) bool {
	for _, bc := range cfg.Blockchains {
		if bc.Type != blockchain.TypeCanton {
			continue
		}
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(bc.ChainID, chainsel.FamilyCanton)
		if err == nil && details.ChainSelector == selector {
			return true
		}
	}
	return false
}

func filterLanesByQualifier(lanes []TokenLane, qualifier string) []TokenLane {
	out := make([]TokenLane, 0, len(lanes))
	for _, lane := range lanes {
		localQ := lane.Combo.LocalPoolAddressRef().Qualifier
		remoteQ := lane.Combo.RemotePoolAddressRef().Qualifier
		if localQ == qualifier || remoteQ == qualifier ||
			strings.Contains(localQ, qualifier) || strings.Contains(remoteQ, qualifier) {
			out = append(out, lane)
		}
	}
	return out
}

func filterDefaultDevenvLanes(lanes []TokenLane) []TokenLane {
	out := make([]TokenLane, 0, len(lanes))
	for _, lane := range lanes {
		if isDefaultBM20LR20Lane(lane.Combo) {
			out = append(out, lane)
		}
	}
	return out
}

func isDefaultBM20LR20Lane(combo devenvcommon.TokenCombination) bool {
	local := combo.LocalPoolAddressRef()
	remote := combo.RemotePoolAddressRef()

	localCCVs := combo.LocalPoolCCVQualifiers()
	remoteCCVs := combo.RemotePoolCCVQualifiers()
	if len(localCCVs) != 1 || localCCVs[0] != devenvcommon.DefaultCommitteeVerifierQualifier {
		return false
	}
	if len(remoteCCVs) != 1 || remoteCCVs[0] != devenvcommon.DefaultCommitteeVerifierQualifier {
		return false
	}

	bmLR := string(local.Type) == devenvcommon.BurnMintTokenPoolType &&
		local.Version.String() == "2.0.0" &&
		string(remote.Type) == devenvcommon.LockReleaseTokenPoolType &&
		remote.Version.String() == "2.0.0"
	lrBM := string(local.Type) == devenvcommon.LockReleaseTokenPoolType &&
		local.Version.String() == "2.0.0" &&
		string(remote.Type) == devenvcommon.BurnMintTokenPoolType &&
		remote.Version.String() == "2.0.0"

	return bmLR || lrBM
}

func formatLaneList(lanes []TokenLane) string {
	parts := make([]string, 0, len(lanes))
	for _, lane := range lanes {
		local := lane.Combo.LocalPoolAddressRef()
		remote := lane.Combo.RemotePoolAddressRef()
		parts = append(parts, fmt.Sprintf("%s %s -> %s %s (%q)",
			local.Type, local.Version, remote.Type, remote.Version, local.Qualifier))
	}
	return strings.Join(parts, "; ")
}
