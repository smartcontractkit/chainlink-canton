package daml_lf

import (
	"strings"

	"github.com/goccy/go-yaml"
)

type DarManifest struct {
	ManifestVersion  string
	CreatedBy        string
	Name             string
	SdkVersion       string
	DalfMain         string
	DalfDependencies []string
	Format           string
	Encryption       string
}

type manifestYaml struct {
	ManifestVersion string `yaml:"Manifest-Version"`
	CreatedBy       string `yaml:"Created-By"`
	Name            string `yaml:"Name"`
	SdkVersion      string `yaml:"Sdk-Version"`
	DalfMain        string `yaml:"Main-Dalf"`
	Dalfs           string `yaml:"Dalfs"`
	Format          string `yaml:"Format"`
	Encryption      string `yaml:"Encryption"`
}

func (m *DarManifest) Parse(b []byte) error {
	parsed := manifestYaml{}
	err := yaml.Unmarshal(b, &parsed)
	if err != nil {
		return err
	}

	// Strip all whitespaces from string
	parsed.ManifestVersion = strings.ReplaceAll(parsed.ManifestVersion, " ", "")
	parsed.CreatedBy = strings.ReplaceAll(parsed.CreatedBy, " ", "")
	parsed.Name = strings.ReplaceAll(parsed.Name, " ", "")
	parsed.SdkVersion = strings.ReplaceAll(parsed.SdkVersion, " ", "")
	parsed.DalfMain = strings.ReplaceAll(parsed.DalfMain, " ", "")
	parsed.Dalfs = strings.ReplaceAll(parsed.Dalfs, " ", "")
	parsed.Format = strings.ReplaceAll(parsed.Format, " ", "")
	parsed.Encryption = strings.ReplaceAll(parsed.Encryption, " ", "")

	dalfDependencies := strings.Split(parsed.Dalfs, ",")

	m.ManifestVersion = parsed.ManifestVersion
	m.CreatedBy = parsed.CreatedBy
	m.Name = parsed.Name
	m.SdkVersion = parsed.SdkVersion
	m.DalfMain = parsed.DalfMain
	m.DalfDependencies = dalfDependencies
	m.Format = parsed.Format
	m.Encryption = parsed.Encryption

	return nil
}
