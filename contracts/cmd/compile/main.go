package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type packageDaml struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	SdkVersion string `yaml:"sdk-version"`
}

func main() {
	ctx := context.Background()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	artifactsDir := flag.String("artifacts", "dar", "Path to the artifacts directory")
	rootDir := flag.String("root", "", "Path to the contracts root directory, must contain multi-package.yaml file, defaults to current working directory")
	flag.Parse()

	// Detect and install all required SDK versions
	// Parse multi-package.yaml to get all packages to-be-compiled
	content, err := os.ReadFile(filepath.Join(*rootDir, "multi-package.yaml"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read multi-package.yaml")
	}
	var multiPackage struct {
		Packages []string `yaml:"packages"`
	}
	err = yaml.Unmarshal(content, &multiPackage)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse multi-package.yaml")
	}

	sdkVersions := make(map[string]struct{})
	for _, pkg := range multiPackage.Packages {
		// Read daml.yaml for package name and version
		damlYamlPath := filepath.Join(*rootDir, pkg, "daml.yaml")
		content, err := os.ReadFile(damlYamlPath)
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to read daml.yaml")
		}
		var config packageDaml
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to parse daml.yaml")
		}
		sdkVersions[config.SdkVersion] = struct{}{}
	}

	// Install required SDK versions
	versions := maps.Keys(sdkVersions)
	slices.Sort(versions)
	log.Info().Any("versions", versions).Msg("Installing required DAML SDK versions")
	for _, version := range versions {
		cmd := exec.CommandContext(ctx, "dpm", "install", version)
		if *rootDir != "" {
			cmd.Dir = *rootDir
		}
		stdOut := &bytes.Buffer{}
		stdErr := &bytes.Buffer{}
		cmd.Stdout = stdOut
		cmd.Stderr = stdErr
		err := cmd.Run()
		if err != nil {
			log.Error().Str("output", stdErr.String()).Msg("dpm install command failed")
			log.Error().Str("stdout", stdOut.String()).Msg("dpm install stdout")
			log.Fatal().Err(err).Str("version", version).Msg("failed to run dpm install command")
		}
		log.Info().Str("version", version).Msg("DAML SDK version installed successfully")
	}

	// Compile contracts using dpm
	log.Info().Msg("Compiling contracts using dpm")
	cmd := exec.CommandContext(ctx, "dpm", "build", "--all")
	if *rootDir != "" {
		cmd.Dir = *rootDir
	}
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err = cmd.Run()
	if err != nil {
		log.Error().Str("output", stdErr.String()).Msg("dpm build command failed")
		log.Error().Str("stdout", stdOut.String()).Msg("dpm build stdout")
		log.Fatal().Err(err).Msg("failed to run dpm build command")
	}

	log.Info().Msg("Contracts compiled successfully")

	// Assemble artifacts

	// For each package, read the compiled DAR file from .daml/dist/
	// Write the DAR files to the artifacts directory

	for _, pkg := range multiPackage.Packages {
		// Read daml.yaml for package name and version
		damlYamlPath := filepath.Join(*rootDir, pkg, "daml.yaml")
		content, err := os.ReadFile(damlYamlPath)
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to read daml.yaml")
		}
		var config packageDaml
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to parse daml.yaml")
		}

		// Read compiled DAR file
		darPath := filepath.Join(*rootDir, pkg, ".daml", "dist", config.Name+"-"+config.Version+".dar")
		darBytes, err := os.ReadFile(darPath)
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to read compiled DAR file")
		}

		// Write DAR file(s) to artifacts directory

		// Check if the package + version is marked as released
		releasedVersions := contracts.ReleasedVersions[contracts.Package(config.Name)]
		if slices.Contains(releasedVersions, config.Version) {
			// If released, write to the release directory
			outPath := filepath.Join(*artifactsDir, contracts.ReleaseDir, fmt.Sprintf("%s-%s.dar", config.Name, config.Version))
			err = os.WriteFile(outPath, darBytes, 0600) //nolint:gosec // Only used during compilation
			if err != nil {
				log.Fatal().Err(err).Str("package", pkg).Msg("failed to write DAR file to artifacts directory")
			}
			log.Info().Str("package", pkg).Str("version", config.Version).Msg("added DAR file")
		}

		// Write dev version to dev directory
		outPath := filepath.Join(*artifactsDir, contracts.DevDir, fmt.Sprintf("%s-dev.dar", config.Name))
		err = os.WriteFile(outPath, darBytes, 0600) //nolint:gosec // Only used during compilation
		if err != nil {
			log.Fatal().Err(err).Str("package", pkg).Msg("failed to write DAR file to artifacts directory")
		}
		log.Info().Str("package", pkg).Str("version", "dev").Msg("added DAR file")
	}

	log.Info().Msg("Artifacts assembled successfully")
}
