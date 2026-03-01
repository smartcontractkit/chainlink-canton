package committee_verifier

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CommitteeVerifier")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[ccvs.CommitteeVerifier]{
	Name:           "canton/ccip/committee_verifier/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP CommitteeVerifier contract on Canton",
	Validate: func(template ccvs.CommitteeVerifier) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}
		if template.VersionTag == "" {
			return errors.New("version tag cannot be empty")
		}
		if template.MessageSentObserver == "" {
			return errors.New("message sent observer cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommitteeVerifier),
	Prefix:      "committeeverifier",
})

var ApplySignatureConfigs = contract.NewExercise(contract.ExerciseParams[ccvs.CommitteeVerifierApplySignatureConfigs]{
	Name:         "canton/ccip/committee_verifier/apply_signature_configs",
	Version:      Version,
	Description:  "Applies new signature configs to a CommitteeVerifier instance by adding/removing configs",
	ContractType: ContractType,
	Validate: func(input ccvs.CommitteeVerifierApplySignatureConfigs) error {
		for i, config := range input.SignatureConfigs {
			// Verify all thresholds
			if int(config.Threshold) > len(config.SignerKeys) {
				return fmt.Errorf("threshold of config at index %d cannot be greater than the number of signer keys", i)
			}
			// Verify that the signer keys are correctly encoded uncompressed pubkeys of 65 byte length (4-prefixed)
			for i2, key := range config.SignerKeys {
				pubkey, err := hex.DecodeString(strings.TrimPrefix(string(key), "0x"))
				if err != nil {
					return fmt.Errorf("invalid signer key at index %d of config at index %d: %w", i2, i, err)
				}
				if len(pubkey) != 65 {
					return fmt.Errorf("invalid signer key length at index %d of config at index %d: expected uncompressed pubkey with 65 bytes, got %d bytes", i2, i, len(pubkey))
				}
			}
		}

		return nil
	},
	Template: ccvs.CommitteeVerifier{},
	Method:   ccvs.CommitteeVerifier{}.CommitteeVerifierApplySignatureConfigs,
})
