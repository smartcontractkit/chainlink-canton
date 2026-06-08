package fee_quoter

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("FeeQuoter")

var Version = semver.MustParse("0.1.0")

var feeQuoterEncoder = core.NewContract("", "CCIP.FeeQuoter", "FeeQuoter").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[core.FeeQuoter]{
	Name:           "canton/ccip/fee_quoter/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP FeeQuoter contract on Canton",
	Validate: func(template core.FeeQuoter) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPFeeQuoter),
	Prefix:      "feequoter",
})

var ApplyPriceUpdatersUpdate = contract.NewExercise(contract.ExerciseParams[core.ApplyPriceUpdatersUpdate]{
	Name:         "canton/ccip/fee_quoter/apply_price_updates",
	Version:      Version,
	Description:  "Adds and removed prices updaters on a Canton FeeQuoter",
	ContractType: ContractType,
	Validate:     nil,
	Template:     core.FeeQuoter{},
	Method:       core.FeeQuoter{}.ApplyPriceUpdatersUpdate,
	EncodeMethod: feeQuoterEncoder.ApplyPriceUpdatersUpdate,
})

var UpdatePrices = contract.NewExercise(contract.ExerciseParams[core.UpdatePrices]{
	Name:         "canton/ccip/fee_quoter/update_prices",
	Version:      Version,
	Description:  "Updates the FeeQuoter's prices",
	ContractType: ContractType,
	Validate: func(input core.UpdatePrices) error {
		// TODO add validation
		return nil
	},
	Modifier: func(chain canton.Chain, input core.UpdatePrices) (core.UpdatePrices, error) {
		// Automatically set the caller
		input.Caller = types.PARTY(chain.Participants[0].PartyID)

		return input, nil
	},
	Template:     core.FeeQuoter{},
	Method:       core.FeeQuoter{}.UpdatePrices,
	EncodeMethod: feeQuoterEncoder.UpdatePrices,
})

var RemoveFeeTokens = contract.NewExercise(contract.ExerciseParams[core.RemoveFeeTokens]{
	Name:         "canton/ccip/fee_quoter/remove_fee_tokens",
	Version:      Version,
	Description:  "Removes fee tokens from the FeeQuoter and clears their prices",
	ContractType: ContractType,
	Validate: func(input core.RemoveFeeTokens) error {
		// TODO add validation
		return nil
	},
	Template:     core.FeeQuoter{},
	Method:       core.FeeQuoter{}.RemoveFeeTokens,
	EncodeMethod: feeQuoterEncoder.RemoveFeeTokens,
})

var ApplyDestChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[core.ApplyFeeQuoterDestChainConfigUpdates]{
	Name:         "canton/ccip/fee_quoter/apply_dest_chain_config_updates",
	Version:      Version,
	Description:  "Applies destination chain configuration updates to the FeeQuoter",
	ContractType: ContractType,
	Validate: func(input core.ApplyFeeQuoterDestChainConfigUpdates) error {
		for _, cfg := range input.DestChainConfigArgs {
			if cfg.DestChainConfig.LinkFeeMultiplierPercent == "" {
				return fmt.Errorf("linkFeeMultiplierPercent cannot be empty for dest chain %s", cfg.DestChainSelector)
			}
		}

		return nil
	},
	Template:     core.FeeQuoter{},
	Method:       core.FeeQuoter{}.ApplyFeeQuoterDestChainConfigUpdates,
	EncodeMethod: feeQuoterEncoder.ApplyFeeQuoterDestChainConfigUpdates,
})
