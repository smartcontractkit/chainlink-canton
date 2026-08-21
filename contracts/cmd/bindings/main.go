package main

import (
	"archive/zip"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/go-daml/codegen"
	"github.com/smartcontractkit/go-daml/codegen/model"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/utilitydars"
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if buildInfo == nil || !ok {
		log.Fatal().Msg("Failed to read build info")
		return
	}

	artifactsDir := flag.String("artifacts", filepath.Join("contracts", "bindings", "generated"), "Path to the bindings artifacts output directory")
	basePath := flag.String("basePath", buildInfo.Main.Path+"/contracts/bindings/generated", "Base Go import path for generated bindings")
	flag.Parse()

	ctx := context.Background()
	tmpDir, cleanup, err := utilitydars.FetchToTemp(ctx, contracts.UtilityPackageIDs())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch utility DARs for binding generation")
	}
	defer cleanup()
	utilitydars.SetResolveDir(tmpDir)
	defer utilitydars.SetResolveDir("")

	log.Info().Str("artifacts", *artifactsDir).Msg("Generating bindings...")

	// Ensure the output directory exists
	err = os.MkdirAll(*artifactsDir, 0o755)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create artifacts directory")
	}

	// Collect all external packages info
	log.Debug().Msg("Collecting external package information...")
	externalPackages := model.ExternalPackages{
		Packages: make(map[string]model.ExternalPackage, len(contracts.BindingsOutputDirs)),
	}
	for p, s := range contracts.BindingsOutputDirs {
		dar, err := contracts.GetDar(p, contracts.Versions[p][len(contracts.Versions[p])-1])
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Msg("Failed to get DAR for package")
		}
		packageId, err := getMainPackageId(dar)
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Msg("Failed to get main package ID from DAR")
		}
		externalPackage := model.ExternalPackage{
			Import: fmt.Sprintf("%s/%s", *basePath, strings.Join(s, "/")),
			Alias:  s[len(s)-1],
		}
		log.Debug().Str("packageId", packageId).Str("alias", externalPackage.Alias).Str("import", externalPackage.Import).Msg("Collected external package information")
		externalPackages.Packages[packageId] = externalPackage
	}
	log.Debug().Int("count", len(externalPackages.Packages)).Msg("Collected external package information")

	// Field hints declare which field names in the generated MCMS/CCIP contracts
	// require non-default hex encoding tags. The Daml compiler erases type synonyms
	// (e.g. BytesHex → Text), so these cannot be inferred from the compiled .dalf.
	fieldHints := model.FieldHints{
		// hex:"bytes" — fixed-size hex fields ≤255 bytes
		BytesFields: map[string]bool{
			"signerAddress":       true, // EVM signer address (20 bytes)
			"chainFamilySelector": true, // Chain family selector (4 bytes)
			"root":                true, // Merkle root hash (32 bytes)
			"newRoot":             true, // New merkle root hash (32 bytes)
			"versionTag":          true, // CCV version tag is a 4-byte BytesHex field
			"onRampAddresses":     true, // Daml [BytesHex] → hex:"[]bytes" on []types.TEXT
			"signerKeys":          true, // Daml [BytesHex] → hex:"[]bytes" on []types.TEXT
			"offRampAddress":      true, // Daml BytesHex scalar → hex:"bytes" on types.TEXT
			"subject":             true, // Daml BytesHex scalar → hex:"bytes"; RMNRemote Curse/Uncurse subject (chain selector hash)
			"remoteTokenAddress":  true, // Daml BytesHex scalar → hex:"bytes"; token pool remote token address on dest chain
			"subjects":            true, // Daml [BytesHex] → hex:"[]bytes"; RMNRemote CurseMultiple/UncurseMultiple
			"cursedSubjects":      true, // Daml [BytesHex] → hex:"[]bytes"; Factory DeployRMNRemote initial cursed subjects
			"remotePools":         true, // Daml [BytesHex] → hex:"[]bytes"; token pool allowed remote pool addresses per chain
		},
		// hex:"bytes16" — fields that may exceed 255 bytes (uint16 length prefix)
		BytesHexFields: map[string]bool{
			"operationData": true, // Serialized choice parameters in TimelockCall
			"predecessor":   true, // Hex hash in ScheduleBatchParams
			"salt":          true, // Hex value in ScheduleBatchParams
		},
		// hex:"uint32" — INT64 fields encoded as 4-byte uint32
		Uint32Fields: map[string]bool{
			"signerIndex": true, // encodeUint32 in encodeSignerInfo
			"signerGroup": true, // encodeUint32 in encodeSignerInfo
		},
		// hex:"[]uint32" — []INT64 fields where each element is a 4-byte uint32
		Uint32ListFields: map[string]bool{
			"groupQuorums":   true, // encodeUint32 list in encodeMultisigConfig
			"groupParents":   true, // encodeUint32 list in encodeMultisigConfig
			"apGroupQuorums": true, // used in AdminParams.AP_SetConfig
			"apGroupParents": true, // used in AdminParams.AP_SetConfig
		},
		// hex:"decimal" — Daml Decimal fields encoded via MCMS.Codec.encodeDecimal (sign byte + 10^10 shift).
		DecimalFields: map[string]bool{
			"usdPerUnitGas": true, // FeeQuoter gas price updates
			"usdPerToken":   true, // FeeQuoter token price updates
		},
		VariantTagByteMap: map[string]map[string]byte{
			"CCIP.LockReleaseTokenPoolV2Types.TransferTimeout": {
				"Indefinite":    0x00,
				"RelativeHours": 0x01,
			},
			"CCIP.BurnMintTokenPoolV2Types.TransferTimeout": {
				"Indefinite":    0x00,
				"RelativeHours": 0x01,
			},
			// Tags must match encode/decode in contracts/ccip/codec/daml/CCIP/CodecV2/FinalityConfig.daml.
			"CCIP.CodecV2.FinalityConfig.FinalityConfig": {
				"WaitForFinality": 0x00,
				"WaitForSafe":     0x01,
				"BlockDepth":      0x02,
			},
		},
		// EnumTagByteMap: enums decoded as a single uint8 ordinal byte by the Daml MCMS codec.
		// Tags must match the decodeXxxAt case statement in the corresponding Daml contract.
		EnumTagByteMap: map[string]map[string]byte{
			// Tags match decodeRateLimitDirectionAt in contracts/ccip/factory/daml/CCIP/FactoryV2.daml.
			"CCIP.RateLimiterV2.RateLimitDirection": {
				"RateLimitDirection_Outbound": 0x00,
				"RateLimitDirection_Inbound":  0x01,
			},
			// Tags match decodeRateLimitModeAt in contracts/ccip/factory/daml/CCIP/FactoryV2.daml.
			"CCIP.RateLimiterV2.RateLimitMode": {
				"RateLimitMode_DefaultFinality": 0x00,
				"RateLimitMode_CustomFinality":  0x01,
			},
		},
		// Dispatcher operationData payloads decoded by MCMS.Main.ExecuteOp.
		ChoiceParamEncoderNames: map[string]bool{
			"SetConfig":            true,
			"ScheduleBatch":        true,
			"CancelBatch":          true,
			"BypasserExecuteBatch": true,
		},
		// TokenAdminRegistry mcmsEntrypoint decodes Params records whose names do not match choice names.
		ChoiceOperationDataParams: map[string]string{
			"ProposeAdministrator": "ProposeAdminParams",
			"AcceptAdminRole":      "AcceptAdminParams",
			"TransferAdminRole":    "TransferAdminParams",
		},
	}

	// Generate bindings for each package
	for p, s := range contracts.BindingsOutputDirs {
		dar, err := contracts.GetDar(p, contracts.Versions[p][len(contracts.Versions[p])-1])
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Msg("Failed to get DAR for package")
		}
		log.Info().Str("package", string(p)).Msg("Generating bindings for package...")
		output, err := generatePackage(dar, s[len(s)-1], externalPackages, fieldHints)
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Str("package", string(p)).Msg("Failed to generate bindings for package")
		}
		outputFile := filepath.Join(*artifactsDir, filepath.Join(s...), fmt.Sprintf("%s.go", s[len(s)-1]))
		log.Debug().Str("package", string(p)).Str("outputFile", outputFile).Msg("Writing generated bindings to file")
		// Ensure the output subdirectory exists
		err = os.MkdirAll(filepath.Dir(outputFile), 0o755)
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Str("outputFile", outputFile).Msg("Failed to create output subdirectory for package")
		}

		err = os.WriteFile(outputFile, output, 0o644)
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Str("outputFile", outputFile).Msg("Failed to write generated bindings to file")
		}
	}

	// Utility DARs need post-processing: Tuple2 shims (prim dalfs skipped above) and
	// splice Transfer2 collision fix when registry packages define a local Transfer2 type.
	if err := writeTuple2Shims(*artifactsDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to write utility Tuple2 shims")
	}
	if err := patchUtilityBindingCollisions(*artifactsDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to patch utility binding collisions")
	}

	log.Info().Msg("Successfully generated all bindings")
}

// Gets the main package ID from the Dar file.
func getMainPackageId(dar []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(dar), int64(len(dar)))
	if err != nil {
		return "", fmt.Errorf("failed to created zip reader: %w", err)
	}
	manifest, err := codegen.GetManifest(reader)
	if err != nil {
		return "", fmt.Errorf("failed to parse manifest: %w", err)
	}

	return codegen.GetPackageID(manifest.MainDalf), nil
}

// Generated a single package's code and returns the generated code.
func generatePackage(dar []byte, pkgFile string, externalPackages model.ExternalPackages, fieldHints model.FieldHints) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(dar), int64(len(dar)))
	if err != nil {
		return nil, fmt.Errorf("failed to created zip reader: %w", err)
	}
	manifest, err := codegen.GetManifest(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	dalfs := []string{
		manifest.MainDalf,
	}
	for _, dalf := range manifest.Dalfs {
		if dalf == manifest.MainDalf {
			continue
		}

		dalfLower := strings.ToLower(dalf)
		if strings.Contains(dalfLower, "stdlib") {
			continue
		}
		// Skip prim dalfs; utility Tuple2 references are shimmed in per-package tuple2.go.
		if strings.Contains(dalfLower, "prim") {
			continue
		}

		dalfs = append(dalfs, dalf)
	}

	result, err := codegen.CodegenDalfs(dalfs, reader, pkgFile, manifest, true, externalPackages, fieldHints)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}
	res, ok := result[manifest.MainDalf]
	if !ok {
		return nil, fmt.Errorf("generated code not found for main dalf")
	}

	return []byte(res), nil
}
