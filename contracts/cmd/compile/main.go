package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type packageDaml struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	SdkVersion string `yaml:"sdk-version"`
}

func main() {
	ctx := context.Background()

	artifactsDir := flag.String("artifacts", "dar", "Path to the artifacts directory")
	rootDir := flag.String("root", "", "Path to the contracts root directory, must contain multi-package.yaml file, defaults to current working directory")
	flag.Parse()

	// Detect and install all required SDK versions
	// Parse multi-package.yaml to get all packages to-be-compiled
	content, err := os.ReadFile(filepath.Join(*rootDir, "multi-package.yaml"))
	if err != nil {
		log.Fatalf("failed to read multi-package.yaml: %v", err)
	}
	var multiPackage struct {
		Packages []string `yaml:"packages"`
	}
	err = yaml.Unmarshal(content, &multiPackage)
	if err != nil {
		log.Fatalf("failed to parse multi-package.yaml: %v", err)
	}

	sdkVersions := make(map[string]struct{})
	for _, pkg := range multiPackage.Packages {
		// Read daml.yaml for package name and version
		damlYamlPath := filepath.Join(*rootDir, pkg, "daml.yaml")
		content, err := os.ReadFile(damlYamlPath)
		if err != nil {
			log.Fatalf("failed to read daml.yaml for package %q: %v", pkg, err)
		}
		var config packageDaml
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			log.Fatalf("failed to parse daml.yaml for package %q: %v", pkg, err)
		}
		sdkVersions[config.SdkVersion] = struct{}{}
	}

	// Install required SDK versions
	versions := maps.Keys(sdkVersions)
	slices.Sort(versions)
	log.Printf("Installing required DAML SDK versions: %v...", versions)
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
			log.Println("dpm install command failed:")
			log.Println(stdErr.String())
			log.Println(stdOut.String())
			log.Fatalf("failed to run dpm install command for version %q: %v", version, err)
		}
		log.Printf("DAML SDK version %q installed successfully", version)
	}

	// Compile contracts using dpm
	log.Printf("Compiling contracts using dpm...")
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
		log.Println("dpm build command failed:")
		log.Println(stdErr.String())
		log.Println(stdOut.String())
		log.Fatalf("failed to run dpm build command: %v", err)
	}

	log.Println("Contracts compiled successfully")

	// Assemble artifacts

	err = os.MkdirAll(*artifactsDir, 0755)
	if err != nil {
		log.Fatalf("failed to create artifacts directory: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(*artifactsDir, "*.dar"))
	if err != nil {
		log.Fatalf("failed to list existing DAR artifacts: %v", err)
	}
	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			log.Fatalf("failed to remove stale DAR artifact %q: %v", match, err)
		}
	}

	// For each package, read the compiled DAR file from .daml/dist/
	// Write the DAR files to the artifacts directory

	for _, pkg := range multiPackage.Packages {
		// Read daml.yaml for package name and version
		damlYamlPath := filepath.Join(*rootDir, pkg, "daml.yaml")
		content, err := os.ReadFile(damlYamlPath)
		if err != nil {
			log.Fatalf("failed to read daml.yaml for package %q: %v", pkg, err)
		}
		var config packageDaml
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			log.Fatalf("failed to parse daml.yaml for package %q: %v", pkg, err)
		}

		// Read compiled DAR file
		darPath := filepath.Join(*rootDir, pkg, ".daml", "dist", config.Name+"-"+config.Version+".dar")
		darBytes, err := os.ReadFile(darPath)
		if err != nil {
			log.Fatalf("failed to read compiled DAR file for package %q: %v", pkg, err)
		}

		// Write DAR file to artifacts directory

		// Create two files, one with the version included and one with `current`
		for _, suffix := range []string{config.Version, "current"} {
			outPath := filepath.Join(*artifactsDir, config.Name+"-"+suffix+".dar")
			err = os.WriteFile(outPath, darBytes, 0600) //nolint:gosec // Only used during compilation
			if err != nil {
				log.Fatalf("failed to write DAR file for package %q to artifacts directory: %v", pkg, err)
			}
			log.Printf("added DAR file for package %q version %q", pkg, suffix)
		}
	}

	log.Println("Artifacts assembled successfully")
}
