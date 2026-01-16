package compile

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
)

type CompiledPackage struct {
	PackageName string
	Dar         []byte
}

func Package(packageName contracts.Package) (CompiledPackage, error) {
	packageDir, ok := contracts.Contracts[packageName]
	if !ok {
		return CompiledPackage{}, fmt.Errorf("package %s not found", packageName)
	}

	// Create a random temporary directory path
	dstDir, err := os.MkdirTemp("", "canton-build-*")
	if err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		// Clean up the temporary directory
		err := os.RemoveAll(dstDir)
		if err != nil {
			fmt.Printf("warning: failed to remove temporary directory %s: %v\n", dstDir, err)
		}
	}()

	srcDir := filepath.Join(".")
	dstRoot := filepath.Join(dstDir, "contracts")

	err = writeEFS(contracts.Embed, srcDir, dstRoot)
	if err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to copy embedded files to %q: %w", dstRoot, err)
	}

	damlYamlPath := filepath.Join(dstRoot, packageDir, "daml.yaml")
	_, err = os.Stat(damlYamlPath)
	if os.IsNotExist(err) {
		return CompiledPackage{}, fmt.Errorf("daml.yaml not found in %q", damlYamlPath)
	} else if err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to stat daml.yaml in %q: %w", damlYamlPath, err)
	}

	// Read SDK version, package name and version
	content, err := os.ReadFile(damlYamlPath)
	if err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to read daml.yaml in %q: %w", damlYamlPath, err)
	}
	// Parse yaml
	var damlConfig struct {
		Name    string `yaml:"name"`
		SDK     string `yaml:"sdk-version"`
		Version string `yaml:"version"`
	}
	if err := yaml.Unmarshal(content, &damlConfig); err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to parse daml.yaml in %q: %w", damlYamlPath, err)
	}

	// Install Daml SDK version
	cmd := exec.Command("dpm", "install", "package")
	cmd.Dir = filepath.Join(dstRoot, packageDir)
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err = cmd.Run()
	if err != nil {
		fmt.Println("daml install command failed:")
		fmt.Println(stdErr.String())
		fmt.Println(stdOut.String())
		return CompiledPackage{}, fmt.Errorf("failed to run daml install command: %w", err)
	}

	// Compile the package using daml CLI
	cmd = exec.Command("dpm", "build")
	cmd.Dir = filepath.Join(dstRoot, packageDir)
	// Buffer stdErr and stdOut
	stdOut = &bytes.Buffer{}
	stdErr = &bytes.Buffer{}
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	err = cmd.Run()
	if err != nil {
		fmt.Println("dpm command failed:")
		fmt.Println(stdErr.String())
		fmt.Println(stdOut.String())
		return CompiledPackage{}, fmt.Errorf("failed to run dpm command: %w", err)
	}

	// Read the compiled DAR file
	darPath := filepath.Join(dstRoot, packageDir, ".daml", "dist", fmt.Sprintf("%s-%s.dar", damlConfig.Name, damlConfig.Version))
	darBytes, err := os.ReadFile(darPath)
	if err != nil {
		return CompiledPackage{}, fmt.Errorf("failed to read compiled DAR file %q: %w", darPath, err)
	}

	return CompiledPackage{
		PackageName: damlConfig.Name,
		Dar:         darBytes,
	}, nil
}

func writeEFS(efs embed.FS, srcDir, dstDir string) error {
	return fs.WalkDir(efs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, path)

		if d.IsDir() {
			err := os.MkdirAll(dstPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %q: %w", dstPath, err)
			}
			return nil
		}

		srcFile, err := efs.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open src file %q: %w", path, err)
		}
		defer func(srcFile fs.File) {
			_ = srcFile.Close()
		}(srcFile)

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("failed to create dst file %q: %w", dstPath, err)
		}
		defer func(dstFile *os.File) {
			_ = dstFile.Close()
		}(dstFile)

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy %q to %q: %w", path, dstPath, err)
		}

		return nil
	})
}
