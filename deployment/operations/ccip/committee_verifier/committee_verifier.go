package committee_verifier

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CommitteeVerifier")

var Version = semver.MustParse("2.0.0")

var ccvsEncoder = committeeverifier.NewContract("", "CCIP.CommitteeVerifierV2", "CommitteeVerifier").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[committeeverifier.CommitteeVerifier]{
	Name:           "canton/ccip/committee_verifier/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP CommitteeVerifier contract on Canton",
	Validate: func(template committeeverifier.CommitteeVerifier) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}
		if template.VersionTag == "" {
			return errors.New("version tag cannot be empty")
		}
		if template.StorageLocationsAdmin == "" {
			return errors.New("storage locations admin cannot be empty")
		}
		if template.PendingStorageLocationsAdmin != template.StorageLocationsAdmin {
			return errors.New("pending storage locations admin should not be set, set to the same value as StorageLocationsAdmin use two-step transfer instead")
		}
		// These configs must not be set as part of deployment
		// If required, they should be set after deployment using ApplySignatureConfig, ApplyRemoteChainConfigUpdates, and ApplyAllowListUpdates respectively
		if len(template.RemoteChainConfigs) != 0 {
			return errors.New("remote chain configs should not be set during deployment")
		}
		if len(template.SignerConfigs) != 0 {
			return errors.New("signer configs should not be set during deployment")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommitteeVerifierV2),
	Prefix:      "committeeverifier",
})

var ApplySignatureConfigs = contract.NewExercise(contract.ExerciseParams[committeeverifier.ApplySignatureConfigs]{
	Name:         "canton/ccip/committee_verifier/apply_signature_configs",
	Version:      Version,
	Description:  "Applies new signature configs to a CommitteeVerifier instance by adding/removing configs",
	ContractType: ContractType,
	Validate: func(input committeeverifier.ApplySignatureConfigs) error {
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
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.ApplySignatureConfigs,
	EncodeMethod: ccvsEncoder.ApplySignatureConfigs,
})

var UpdateStorageLocations = contract.NewExercise(contract.ExerciseParams[committeeverifier.UpdateStorageLocations]{
	Name:         "canton/ccip/committee_verifier/update_storage_locations",
	Version:      Version,
	Description:  "Updates the storage locations of a CommitteeVerifier instance",
	ContractType: ContractType,
	Validate:     nil,
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.UpdateStorageLocations,
	EncodeMethod: ccvsEncoder.UpdateStorageLocations,
})

var TransferStorageLocationsAdmin = contract.NewExercise(contract.ExerciseParams[committeeverifier.TransferStorageLocationsAdmin]{
	Name:         "canton/ccip/committee_verifier/transfer_storage_locations_admin",
	Version:      Version,
	Description:  "Initiates the two-step transfer of the storage locations admin role",
	ContractType: ContractType,
	Validate: func(input committeeverifier.TransferStorageLocationsAdmin) error {
		if input.NewAdmin == "" {
			return errors.New("newAdmin cannot be empty")
		}

		return nil
	},
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.TransferStorageLocationsAdmin,
	EncodeMethod: ccvsEncoder.TransferStorageLocationsAdmin,
})

var AcceptStorageLocationsAdminRole = contract.NewExercise(contract.ExerciseParams[committeeverifier.AcceptStorageLocationsAdmin]{
	Name:         "canton/ccip/committee_verifier/accept_storage_locations_admin",
	Version:      Version,
	Description:  "Accepts a pending transfer of the storage locations admin role",
	ContractType: ContractType,
	Validate:     nil,
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.AcceptStorageLocationsAdmin,
	EncodeMethod: ccvsEncoder.AcceptStorageLocationsAdmin,
})

var SetDynamicConfig = contract.NewExercise(contract.ExerciseParams[committeeverifier.SetDynamicConfig]{
	Name:         "canton/ccip/committee_verifier/set_dynamic_config",
	Version:      Version,
	Description:  "Sets the dynamic config",
	ContractType: ContractType,
	Validate:     nil,
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.SetDynamicConfig,
	EncodeMethod: ccvsEncoder.SetDynamicConfig,
})

var ApplyRemoteChainConfigUpdates = contract.NewExercise(contract.ExerciseParams[committeeverifier.ApplyRemoteChainConfigUpdates]{
	Name:         "canton/ccip/committee_verifier/apply_remote_chain_config_updates",
	Version:      Version,
	Description:  "Applies remote chain configs to a CommitteeVerifier instance by adding/removing configs",
	ContractType: ContractType,
	Validate:     nil,
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.ApplyRemoteChainConfigUpdates,
	EncodeMethod: ccvsEncoder.ApplyRemoteChainConfigUpdates,
})

var ApplyAllowListUpdates = contract.NewExercise(contract.ExerciseParams[committeeverifier.ApplyAllowListUpdates]{
	Name:         "canton/ccip/committee_verifier/apply_allow_list_updates",
	Version:      Version,
	Description:  "Applies allow lists updates to a Canton CommitteeVerifier instance",
	ContractType: ContractType,
	Validate:     nil,
	Template:     committeeverifier.CommitteeVerifier{},
	Method:       committeeverifier.CommitteeVerifier{}.ApplyAllowListUpdates,
	EncodeMethod: ccvsEncoder.ApplyAllowListUpdates,
})
