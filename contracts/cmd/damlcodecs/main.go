package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/go-daml/codegen/damlparser"
	"github.com/smartcontractkit/go-daml/codegen/damltemplate"
)

const contractsRoot = "contracts"

// Custom type codecs map domain-specific types to their hand-written
// encode/decode functions. Update this when a new custom type is introduced.
var customCodecs = map[string]damltemplate.CustomCodec{
	"RawInstanceAddress": {
		EncodeFunc:     "encodeRawInstanceAddress",
		DecodeFunc:     "decodeRawInstanceAddressAt",
		ImportModule:   "CCIP.CodecV2",
		EncodeListFunc: "encodeRawInstanceAddressList",
		DecodeListFunc: "decodeRawInstanceAddressList",
	},
	"InstrumentId": {
		EncodeFunc:     "encodeInstrumentId",
		DecodeFunc:     "decodeInstrumentId",
		ImportModule:   "CCIP.CodecV2",
		EncodeListFunc: "encodeInstrumentIdList",
		DecodeListFunc: "decodeInstrumentIdList",
	},
	"CCIPContext": {
		EncodeFunc:   "encodeCCIPContext",
		DecodeFunc:   "decodeCCIPContextAt",
		ImportModule: "CCIP.ContextCodec",
	},
	"BytesHex": {
		EncodeFunc:     "encodeBytesHex",
		DecodeFunc:     "decodeBytesHexAt",
		ImportModule:   "MCMS.Codec",
		EncodeListFunc: "encodeBytesHexList",
		DecodeListFunc: "decodeBytesHexList",
	},
	"FinalityConfig": {
		EncodeFunc:   "encodeRequestedFinality",
		DecodeFunc:   "decodeRequestedFinalityAt",
		ImportModule: "CCIP.CodecV2",
	},
}

// Variant tag byte overrides. By default variants use constructor index as
// the tag byte. Add entries here when the on-wire tag must differ.
var variantTagBytes = map[string]map[string]int{
	"TransferTimeout": {
		"Indefinite":    0,
		"RelativeHours": 1,
	},
}

func main() {
	log.Info().Msg("Discovering *Types.daml files...")

	targets, err := discoverTargets(contractsRoot)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to discover Types.daml files")
	}

	if len(targets) == 0 {
		log.Warn().Msg("No *Types.daml files found")
		return
	}

	log.Info().Int("count", len(targets)).Msg("Generating Daml codec files...")

	for _, target := range targets {
		if err := generateCodec(target); err != nil {
			log.Fatal().Err(err).Str("types", target.typesFile).Msg("Failed to generate codec")
		}
	}

	log.Info().Int("count", len(targets)).Msg("Successfully generated all Daml codec files")
}

type codecTarget struct {
	typesFile   string
	outputFile  string
	typesModule string
	codecModule string
}

// discoverTargets walks contractsRoot looking for *Types.daml files.
// It skips test directories and files whose Daml module name doesn't end in "Types".
func discoverTargets(root string) ([]codecTarget, error) {
	var targets []codecTarget

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip test directories and .daml build cache
		if d.IsDir() && (d.Name() == "test" || d.Name() == ".daml") {
			return filepath.SkipDir
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), "Types.daml") {
			return nil
		}

		// Bare "Types.daml" (e.g., MCMS/Types.daml) doesn't follow the convention
		if d.Name() == "Types.daml" {
			return nil
		}

		// Parse to extract the module name
		f, err := os.Open(path) //nolint:gosec // Only used during compilation
		if err != nil {
			return fmt.Errorf("opening %s: %w", path, err)
		}
		defer f.Close()

		module, err := damlparser.Parse(f)
		if err != nil {
			log.Warn().Err(err).Str("file", path).Msg("Skipping unparseable file")
			return nil
		}

		if !strings.HasSuffix(module.ModuleName, "Types") {
			return nil
		}

		baseName := strings.TrimSuffix(module.ModuleName, "Types")
		codecModule := baseName + "CodecGen"
		outputFile := strings.TrimSuffix(path, "Types.daml") + "CodecGen.daml"

		targets = append(targets, codecTarget{
			typesFile:   path,
			outputFile:  outputFile,
			typesModule: module.ModuleName,
			codecModule: codecModule,
		})

		return nil
	})

	return targets, err
}

func generateCodec(target codecTarget) error {
	f, err := os.Open(target.typesFile)
	if err != nil {
		return fmt.Errorf("opening types file %s: %w", target.typesFile, err)
	}
	defer f.Close()

	module, err := damlparser.Parse(f)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", target.typesFile, err)
	}

	config := damltemplate.DamlCodecConfig{
		ModuleName:        target.codecModule,
		TypesModule:       target.typesModule,
		CustomTypeCodecs:  customCodecs,
		VariantTagByteMap: variantTagBytes,
		// No TargetTypes: generate codecs for ALL records/variants in the file
	}

	output, err := damltemplate.Generate(module, config)
	if err != nil {
		return fmt.Errorf("generating codec for %s: %w", target.typesModule, err)
	}

	if err := os.MkdirAll(filepath.Dir(target.outputFile), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	if err := os.WriteFile(target.outputFile, []byte(output), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", target.outputFile, err)
	}

	log.Info().
		Str("types", target.typesFile).
		Str("output", target.outputFile).
		Int("records", len(module.Records)).
		Int("variants", len(module.Variants)).
		Msg("Generated codec")

	return nil
}
