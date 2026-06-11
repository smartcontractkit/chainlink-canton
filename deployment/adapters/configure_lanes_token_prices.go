package adapters

import (
	"fmt"
	"maps"
	"math"
	"math/big"
	"strconv"
	"strings"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
)

// CantonRemoteTokenPricesFamilyExtraKey is an optional ConfigureChainForLanesInput.FamilyExtras
// entry for per-remote fee-token USD prices pushed via FeeQuoter::UpdatePrices during lane configure.
//
// Shape (YAML-friendly):
//
//	cantonRemoteTokenPrices:
//	  "<remoteChainSelector>":
//	    "ccipOwner::1220...:link-token": "10"
//
// Keys use "<adminParty>:<instrumentId>"; values are USD per whole token (e.g. "10" for $10 LINK).
// Canton FeeQuoter stores usdPerToken as DAML Decimal (see FeeQuoter.daml tests: 20.0 = $20/LINK).
const CantonRemoteTokenPricesFamilyExtraKey = "cantonRemoteTokenPrices"

const defaultLinkTokenInstrumentID = "link-token"

// defaultLinkUsdPerTokenDollars is the nominal LINK/USD spot used when FamilyExtras omit a price.
const defaultLinkUsdPerTokenDollars int64 = 10

// defaultNativeUsdPerTokenDollars is the nominal Amulet/USD spot used when FamilyExtras omit a price.
const defaultNativeUsdPerTokenDollars int64 = 1

// cantonUsdPerTokenScale is the internal fixed-point scale (USD * 1e8) before formatting to Decimal.
const cantonUsdPerTokenScale int64 = 100_000_000

func resolveTokenPricesForRemoteDest(
	ds datastore.DataStore,
	input ccipadapters.ConfigureChainForLanesInput,
	remoteSelector uint64,
	nativeInstrument *splice_api_token_holding_v1.InstrumentId,
) (map[string]*big.Int, error) {
	ccipOwner, err := resolveCcipOwnerParty(ds, input.ChainSelector)
	if err != nil {
		return nil, err
	}

	prices := map[string]*big.Int{
		fmt.Sprintf("%s:%s", ccipOwner, defaultLinkTokenInstrumentID): usdPerTokenToScaled(defaultLinkUsdPerTokenDollars),
	}
	if nativeInstrument != nil && nativeInstrument.Admin != "" && nativeInstrument.Id != "" {
		key := instrumentPriceKey(nativeInstrument.Admin, nativeInstrument.Id)
		prices[key] = usdPerTokenToScaled(defaultNativeUsdPerTokenDollars)
	}

	extras, err := tokenPricesFromFamilyExtras(input.FamilyExtras, remoteSelector)
	if err != nil {
		return nil, err
	}
	maps.Copy(prices, extras)

	return prices, nil
}

func usdPerTokenToScaled(usdDollars int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(usdDollars), big.NewInt(cantonUsdPerTokenScale))
}

func tokenPricesFromFamilyExtras(extras map[string]any, remoteSelector uint64) (map[string]*big.Int, error) {
	if extras == nil {
		return map[string]*big.Int{}, nil
	}
	raw, ok := extras[CantonRemoteTokenPricesFamilyExtraKey]
	if !ok || raw == nil {
		return map[string]*big.Int{}, nil
	}

	byRemote, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%q must be a map keyed by remote chain selector", CantonRemoteTokenPricesFamilyExtraKey)
	}

	remoteKey := strconv.FormatUint(remoteSelector, 10)
	instruments, ok := byRemote[remoteKey]
	if !ok {
		return map[string]*big.Int{}, nil
	}

	instrumentMap, ok := instruments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%q entry for remote %s must be a map of instrument to price", CantonRemoteTokenPricesFamilyExtraKey, remoteKey)
	}

	out := make(map[string]*big.Int, len(instrumentMap))
	for instrument, priceRaw := range instrumentMap {
		price, err := parseUsdPerTokenPrice(priceRaw)
		if err != nil {
			return nil, fmt.Errorf("%q remote %s instrument %q: %w", CantonRemoteTokenPricesFamilyExtraKey, remoteKey, instrument, err)
		}
		out[instrument] = price
	}

	return out, nil
}

func parseUsdPerTokenPrice(raw any) (*big.Int, error) {
	switch v := raw.(type) {
	case string:
		return parseUsdPerTokenPriceString(strings.TrimSpace(v))
	case int:
		if v <= 0 {
			return nil, fmt.Errorf("price must be positive")
		}

		return usdPerTokenToScaled(int64(v)), nil
	case int64:
		if v <= 0 {
			return nil, fmt.Errorf("price must be positive")
		}

		return usdPerTokenToScaled(v), nil
	case uint64:
		if v == 0 {
			return nil, fmt.Errorf("price must be positive")
		}
		if v > math.MaxInt64 {
			return nil, fmt.Errorf("price exceeds supported range")
		}

		return usdPerTokenToScaled(int64(v)), nil
	case float64:
		if v <= 0 {
			return nil, fmt.Errorf("price must be positive")
		}

		return parseUsdPerTokenPriceString(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return nil, fmt.Errorf("unsupported price type %T", raw)
	}
}

func parseUsdPerTokenPriceString(raw string) (*big.Int, error) {
	if raw == "" {
		return nil, fmt.Errorf("price must be non-empty")
	}

	r, ok := new(big.Rat).SetString(raw)
	if !ok {
		return nil, fmt.Errorf("invalid USD price %q", raw)
	}
	if r.Sign() <= 0 {
		return nil, fmt.Errorf("price must be positive")
	}

	scaled := new(big.Rat).Mul(r, new(big.Rat).SetInt64(cantonUsdPerTokenScale))
	if !scaled.IsInt() {
		return nil, fmt.Errorf("price %q exceeds supported precision", raw)
	}

	return scaled.Num(), nil
}

func resolveCcipOwnerParty(ds datastore.DataStore, chainSelector uint64) (string, error) {
	feeQuoterRef, err := findContractRef(
		ds,
		chainSelector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	)
	if err != nil {
		return "", fmt.Errorf("resolve fee quoter for ccipOwner party: %w", err)
	}

	if party, ok := partyFromDeployLabels(feeQuoterRef.Labels); ok {
		return party, nil
	}

	return "", fmt.Errorf("ccipOwner party not found in FeeQuoter labels on chain %d", chainSelector)
}

func partyFromDeployLabels(labels datastore.LabelSet) (string, bool) {
	for _, label := range labels.List() {
		at := strings.LastIndex(label, "@")
		if at < 0 || at+1 >= len(label) {
			continue
		}
		party := label[at+1:]
		if strings.Contains(party, "::") {
			return party, true
		}
	}
	return "", false
}
