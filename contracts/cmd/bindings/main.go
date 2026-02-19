package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/go-daml/codegen"
	"github.com/smartcontractkit/go-daml/codegen/model"
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if buildInfo == nil || !ok {
		log.Fatal().Msg("Failed to read build info")
		return
	}

	artifactsDir := flag.String("artifacts", filepath.Join("bindings", "generated"), "Path to the bindings artifacts output directory")
	basePath := flag.String("basePath", buildInfo.Main.Path+"/bindings/generated", "Base Go import path for generated bindings")
	flag.Parse()

	log.Info().Str("artifacts", *artifactsDir).Msg("Generating bindings...")

	// Ensure the output directory exists
	err := os.MkdirAll(*artifactsDir, 0o755)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create artifacts directory")
	}

	// Collect all external packages info
	log.Debug().Msg("Collecting external package information...")
	externalPackages := model.ExternalPackages{
		Packages: make(map[string]model.ExternalPackage, len(contracts.OutputDirs)),
	}
	for p, s := range contracts.OutputDirs {
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

	// Generate bindings for each package
	for p, s := range contracts.OutputDirs {
		dar, err := contracts.GetDar(p, contracts.Versions[p][len(contracts.Versions[p])-1])
		if err != nil {
			log.Fatal().Err(err).Str("package", string(p)).Msg("Failed to get DAR for package")
		}
		log.Info().Str("package", string(p)).Msg("Generating bindings for package...")
		output, err := generatePackage(dar, s[len(s)-1], externalPackages)
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
func generatePackage(dar []byte, pkgFile string, externalPackages model.ExternalPackages) ([]byte, error) {
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
		if strings.Contains(dalfLower, "prim") || strings.Contains(dalfLower, "stdlib") {
			continue
		}

		dalfs = append(dalfs, dalf)
	}

	result, err := codegen.CodegenDalfs(dalfs, reader, pkgFile, manifest, true, externalPackages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}
	res, ok := result[manifest.MainDalf]
	if !ok {
		return nil, fmt.Errorf("generated code not found for main dalf")
	}

	return []byte(res), nil
}
