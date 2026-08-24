package devenv

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/ccip/devenv/ledgertarget"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// GetMaxDataBytes implements cciptestinterfaces.CCIP17.
func (c *Chain) GetMaxDataBytes(ctx context.Context, remoteChainSelector uint64) (uint32, error) {
	cfg, err := c.feeQuoterDestConfig(ctx, remoteChainSelector)
	if err != nil {
		return 0, err
	}
	if cfg.MaxDataBytes < 0 {
		return 0, fmt.Errorf("negative maxDataBytes: %d", cfg.MaxDataBytes)
	}
	if cfg.MaxDataBytes > math.MaxUint32 {
		return 0, fmt.Errorf("maxDataBytes overflows uint32: %d", cfg.MaxDataBytes)
	}

	return uint32(cfg.MaxDataBytes), nil
}

// GetMaxPerMsgGasLimit returns the FeeQuoter maxPerMsgGasLimit for a destination chain.
func (c *Chain) GetMaxPerMsgGasLimit(ctx context.Context, remoteChainSelector uint64) (uint32, error) {
	cfg, err := c.feeQuoterDestConfig(ctx, remoteChainSelector)
	if err != nil {
		return 0, err
	}
	if cfg.MaxPerMsgGasLimit < 0 {
		return 0, fmt.Errorf("negative maxPerMsgGasLimit: %d", cfg.MaxPerMsgGasLimit)
	}
	if cfg.MaxPerMsgGasLimit > math.MaxUint32 {
		return 0, fmt.Errorf("maxPerMsgGasLimit overflows uint32: %d", cfg.MaxPerMsgGasLimit)
	}

	return uint32(cfg.MaxPerMsgGasLimit), nil
}

func (c *Chain) feeQuoterDestConfig(ctx context.Context, remoteChainSelector uint64) (ledgertarget.FeeQuoterDestChainConfig, error) {
	if c.e == nil {
		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("canton chain environment is nil")
	}
	if len(c.chain.Participants) == 0 {
		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("no canton participants configured")
	}

	feeQuoterRef, err := c.e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		c.ChainSelector(),
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	))
	if err != nil {
		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("resolve FeeQuoter address: %w", err)
	}

	participant := c.chain.Participants[0]
	party := participant.PartyID
	feeQuoterAddress := contracts.HexToInstanceAddress(feeQuoterRef.Address)

	activeContract, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		[]string{party},
		ledgertarget.FeeQuoter{}.GetTemplateID(),
		feeQuoterAddress,
	)
	if err != nil {
		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("find active FeeQuoter contract: %w", err)
	}

	createArgs := activeContract.GetCreatedEvent().GetCreateArguments()
	if createArgs == nil {
		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("FeeQuoter create arguments missing")
	}

	return destChainConfigFromFeeQuoterCreateArgs(createArgs, remoteChainSelector)
}

func destChainConfigFromFeeQuoterCreateArgs(createArgs *apiv2.Record, remoteChainSelector uint64) (ledgertarget.FeeQuoterDestChainConfig, error) {
	selectorKey := strconv.FormatUint(remoteChainSelector, 10)
	for _, field := range createArgs.GetFields() {
		if field.GetLabel() != "destChainConfigs" {
			continue
		}
		genMap := field.GetValue().GetGenMap()
		if genMap == nil {
			return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("destChainConfigs is not a GenMap")
		}
		for _, entry := range genMap.GetEntries() {
			key := strings.TrimSuffix(entry.GetKey().GetNumeric(), ".")
			if key != selectorKey {
				continue
			}
			record := entry.GetValue().GetRecord()
			if record == nil {
				return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("dest chain config value is not a record")
			}
			var cfg ledgertarget.FeeQuoterDestChainConfig
			if err := ledger.RecordToStruct(record, &cfg); err != nil {
				return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("parse dest chain config record: %w", err)
			}

			return cfg, nil
		}

		return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("no FeeQuoter dest config for chain selector %d", remoteChainSelector)
	}

	return ledgertarget.FeeQuoterDestChainConfig{}, fmt.Errorf("destChainConfigs field not found on FeeQuoter")
}
