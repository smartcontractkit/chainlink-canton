package daml_lf

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	ManifestFilePath = "META-INF/MANIFEST.MF"
)

type DarFile struct {
	Manifest     DarManifest
	Main         DamlLfArchive
	Dependencies []DamlLfArchive
}

func (f *DarFile) FromBytes(b []byte) error {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}

	// Parse Manifest
	manifest, err := r.Open(ManifestFilePath)
	if err != nil {
		return err
	}
	manifestContent, err := io.ReadAll(manifest)
	if err != nil {
		return err
	}
	if err := f.Manifest.Parse(manifestContent); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	// fmt.Printf("%#v\n", f.Manifest)

	// Parse Main Dalf
	if err := ParseArchiveFromFile(&f.Main, r, f.Manifest.DalfMain); err != nil {
		return fmt.Errorf("failed to parse main archive: %w", err)
	}
	// fmt.Printf("%#v\n", f.Main)

	// Parse Dependency Dalfs
	f.Dependencies = make([]DamlLfArchive, len(f.Manifest.DalfDependencies))
	for i, dalfPath := range f.Manifest.DalfDependencies {
		if err := ParseArchiveFromFile(&f.Dependencies[i], r, dalfPath); err != nil {
			return fmt.Errorf("failed to parse dependency archive: %w", err)
		}
	}
	// for i, dependency := range f.Dependencies {
	// 	fmt.Printf("Dependency %v: %#v\n", i, dependency)
	// }

	return nil
}

func ParseArchiveFromFile(archive *DamlLfArchive, r *zip.Reader, dalfPath string) error {
	dalfName := dalfPath[strings.LastIndex(dalfPath, "/")+1:]
	d, err := r.Open(dalfPath)
	if err != nil {
		return err
	}
	dalfContent, err := io.ReadAll(d)
	if err != nil {
		return err
	}
	return archive.ParseNamed(dalfName, dalfContent)
}
