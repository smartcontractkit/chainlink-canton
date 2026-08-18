package utilitydars

import (
	"fmt"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Manifest describes the pinned Canton Network Utility DAR bundle.
type Manifest struct {
	BundleVersion string            `yaml:"bundleVersion"`
	URL           string            `yaml:"url"`
	SHA256        string            `yaml:"sha256"`
	Packages      map[string]string `yaml:"packages"`
}

var (
	manifestMu sync.RWMutex
	manifest   *Manifest
)

// SetManifest installs the pinned bundle metadata (typically from contracts embed).
func SetManifest(m *Manifest) {
	manifestMu.Lock()
	defer manifestMu.Unlock()
	manifest = m
}

// ParseManifest parses manifest YAML bytes.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse utility manifest: %w", err)
	}
	if m.BundleVersion == "" {
		return nil, fmt.Errorf("utility manifest: bundleVersion is required")
	}
	if m.URL == "" {
		return nil, fmt.Errorf("utility manifest: url is required")
	}
	if m.SHA256 == "" {
		return nil, fmt.Errorf("utility manifest: sha256 is required")
	}
	if len(m.Packages) == 0 {
		return nil, fmt.Errorf("utility manifest: packages is required")
	}

	return &m, nil
}

func currentManifest() (*Manifest, error) {
	manifestMu.RLock()
	m := manifest
	manifestMu.RUnlock()
	if m == nil {
		return nil, fmt.Errorf("utility manifest not configured")
	}

	return m, nil
}

// SemverForPackage returns the pinned semver for a utility package name.
func SemverForPackage(packageName string) (string, error) {
	m, err := currentManifest()
	if err != nil {
		return "", err
	}

	semver, ok := m.Packages[packageName]
	if !ok {
		return "", fmt.Errorf("utility package %q not in manifest", packageName)
	}

	return semver, nil
}

// PackageNames returns all utility package names from the manifest.
func PackageNames() ([]string, error) {
	m, err := currentManifest()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(m.Packages))
	for name := range m.Packages {
		names = append(names, name)
	}

	return names, nil
}

// ResolveDarPath returns the path to a utility DAR file in dir.
func ResolveDarPath(packageName, semver, dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s-%s.dar", packageName, semver))
}
