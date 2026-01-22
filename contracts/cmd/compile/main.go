package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	artifactsDir := flag.String("artifacts", "dar", "Path to the artifacts directory")
	rootDir := flag.String("root", "", "Path to the contracts root directory, must contain multi-package.yaml file, defaults to current working directory")
	flag.Parse()

	// Compile contracts using dpm
	cmd := exec.Command("dpm", "build", "--all")
	if *rootDir != "" {
		cmd.Dir = *rootDir
	}
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err := cmd.Run()
	if err != nil {
		log.Println("dpm build command failed:")
		log.Println(stdErr.String())
		log.Println(stdOut.String())
		log.Fatalf("failed to run dpm build command: %v", err)
	}

	log.Println("Contracts compiled successfully")

	// Assemble artifacts

	// Read multi-package.yaml to get package names
	// For each package, read the compiled DAR file from .daml/dist/
	// Write the DAR files to the artifacts directory

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

	for _, pkg := range multiPackage.Packages {
		// Read daml.yaml for package name and version
		damlYamlPath := filepath.Join(*rootDir, pkg, "daml.yaml")
		content, err := os.ReadFile(damlYamlPath)
		if err != nil {
			log.Fatalf("failed to read daml.yaml for package %q: %v", pkg, err)
		}
		var damlConfig struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
		}
		err = yaml.Unmarshal(content, &damlConfig)
		if err != nil {
			log.Fatalf("failed to parse daml.yaml for package %q: %v", pkg, err)
		}

		// Read compiled DAR file
		darPath := filepath.Join(*rootDir, pkg, ".daml", "dist", damlConfig.Name+"-"+damlConfig.Version+".dar")
		darBytes, err := os.ReadFile(darPath)
		if err != nil {
			log.Fatalf("failed to read compiled DAR file for package %q: %v", pkg, err)
		}

		// Write DAR file to artifacts directory

		// Create two files, one with the version included and one with `current`
		for _, suffix := range []string{damlConfig.Version, "current"} {
			outPath := filepath.Join(*artifactsDir, damlConfig.Name+"-"+suffix+".dar")
			err = os.WriteFile(outPath, darBytes, 0o644)
			if err != nil {
				log.Fatalf("failed to write DAR file for package %q to artifacts directory: %v", pkg, err)
			}
			log.Printf("added DAR file for package %q version %q", pkg, suffix)
		}
	}

	log.Println("Artifacts assembled successfully")
}
