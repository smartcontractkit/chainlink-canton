// Cannot currently import due to github.com/nao1215/markdown colliding with chainlink-ccv
//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	md "github.com/nao1215/markdown"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var skippedPrefixes = []string{
	"daml-prim",
	"daml-stdlib",
	"ghc-stdlib",
}

func main() {
	ctx := context.Background()
	_ = ctx
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	output := flag.String("output", "", "Optional output Markdown file to write the output to")
	flag.Parse()

	arguments := flag.Args()
	if len(arguments) == 0 {
		log.Fatal().Msg("No files to inspect provided")
	}

	var files []string
	for _, argument := range arguments {
		matches, err := filepath.Glob(argument)
		if err != nil {
			log.Fatal().Err(err).Str("argument", argument).Msg("failed to glob argument")
		}
		files = append(files, matches...)
	}

	log.Info().Strs("files", files).Msg("Inspecting files")

	wg := sync.WaitGroup{}
	results := make([]inspectResult, len(files))
	for i, file := range files {
		wg.Go(func() {
			filename := filepath.Base(file)
			cmd := exec.CommandContext(ctx, "dpm", "inspect-dar", "--json", file)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal().Err(err).Str("file", file).Msgf("failed to inspect file: %s", string(out))
			}
			var output inspectOutput
			if err := json.Unmarshal(out, &output); err != nil {
				log.Fatal().Err(err).Str("file", file).Msgf("failed to parse inspect output: %s", string(out))
			}

			result := inspectResult{
				Filename:           filename,
				MainPackageId:      output.MainPackageId,
				MainPackageName:    output.Packages[output.MainPackageId].Name,
				MainPackageVersion: output.Packages[output.MainPackageId].Version,
				Packages:           nil,
			}
			for packageId, p := range output.Packages {
				// Remove stdlib packages
				skip := false
				for _, prefix := range skippedPrefixes {
					if strings.HasPrefix(p.Name, prefix) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}

				result.Packages = append(result.Packages, inspectResultPackage{
					Name:      p.Name,
					Version:   p.Version,
					PackageID: packageId,
				})
			}

			// Sort packages alphabetically
			slices.SortFunc(result.Packages, func(a, b inspectResultPackage) int {
				return strings.Compare(a.Name, b.Name)
			})

			results[i] = result
		})
	}
	wg.Wait()

	// Sort results by filename
	slices.SortFunc(results, func(a, b inspectResult) int {
		return strings.Compare(a.Filename, b.Filename)
	})

	var writer io.Writer
	if *output != "" {
		// Write to file
		file, err := os.Create(*output)
		if err != nil {
			log.Fatal().Err(err).Str("file", *output).Msg("failed to create output file")
		}
		writer = file
		defer file.Close()
	} else {
		// Print
		writer = os.Stdout
	}

	m := md.NewMarkdown(writer)
	m.H1("Packages")

	for _, result := range results {
		m.H2(fmt.Sprintf(`<a name="%s"></a>%s`, result.MainPackageId, result.Filename))

		m.BulletList(
			fmt.Sprintf("Name: `%s`", result.MainPackageName),
			fmt.Sprintf("Version: `%s`", result.MainPackageVersion),
			fmt.Sprintf("Package ID: `%s`", result.MainPackageId),
		)

		m.H3("Dependencies:")
		var rows [][]string
		for _, pkg := range result.Packages {
			rows = append(rows, []string{fmt.Sprintf("`%s`", pkg.Name), fmt.Sprintf("`%s`", pkg.Version), fmt.Sprintf("[`%s`](#%s)", pkg.PackageID, pkg.PackageID)})
		}
		m.Table(md.TableSet{
			Header: []string{"Package Name", "Version", "Package ID"},
			Rows:   rows,
		})
	}

	if err := m.Build(); err != nil {
		log.Fatal().Err(err).Msg("failed to build markdown")
	}
}

type inspectOutput struct {
	Files         []string `json:"files"`
	MainPackageId string   `json:"main_package_id"`
	Packages      map[string]struct {
		Name    string `json:"name"`
		Path    string `json:"path"`
		Version string `json:"version"`
	} `json:"packages"`
}

type inspectResultPackage struct {
	Name      string
	Version   string
	PackageID string
}

type inspectResult struct {
	Filename           string
	MainPackageId      string
	MainPackageName    string
	MainPackageVersion string
	Packages           []inspectResultPackage
}
