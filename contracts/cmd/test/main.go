package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/fatih/color"

	"github.com/schollz/progressbar/v3"
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

	rootDir := flag.String("root", "", "Path to the contracts root directory, must contain multi-package.yaml file, defaults to current working directory")
	noColor := flag.Bool("no-color", false, "Disables color output")
	plain := flag.Bool("plain", false, "Whether or not to use plain output")
	summaryOutput := flag.String("output", "", "Path to output summary file")

	flag.Parse()

	if *noColor {
		color.NoColor = true
	}

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

	// Collect all required SDK versions and package paths
	var packages []string
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
		packages = append(packages, damlYamlPath)
	}

	log.Println("Testing packages:")
	slices.Sort(packages)
	for _, p := range packages {
		log.Println(color.YellowString("- %s", p))
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

	wg := sync.WaitGroup{}
	var progress *progressbar.ProgressBar
	if !*plain {
		progress = progressbar.NewOptions(len(packages),
			progressbar.OptionSetDescription("Testing packages..."),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionEnableColorCodes(!*noColor),
			progressbar.OptionSetPredictTime(false),
		)
	} else {
		log.Println("Running tests...")
	}
	results := make([]packageResult, len(packages))
	for i, s := range packages {
		wg.Go(func() {
			directory, _ := filepath.Split(s)
			cmd := exec.CommandContext(ctx, "dpm", "test")
			if !*noColor {
				cmd.Args = append(cmd.Args, "--color")
			}
			cmd.Dir = directory
			out, err := cmd.CombinedOutput()
			var exitCode int
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					exitCode = exitErr.ExitCode()
				} else {
					panic(err)
				}
			}
			results[i] = packageResult{
				path:     directory,
				output:   string(out),
				exitCode: exitCode,
			}
			if !*plain {
				_ = progress.Add(1)
			} else {
				log.Println("✓ Done: ", directory)
			}
		})
	}
	wg.Wait()

	var summaries []string //nolint:prealloc
	log.Println("Raw Outputs:")
	for _, r := range results {
		log.Println("==============================")
		log.Printf("Package: %s\n", r.path)
		log.Printf("Exit Code: %d\n\n", r.exitCode)
		log.Println(r.output)

		// Collect summary
		for _, summary := range collectSummaries(r.output) {
			summaries = append(summaries, r.path+summary)
		}
	}
	log.Println("==============================")
	log.Println("Test Summary:")
	failedTests := 0
	for _, summary := range summaries {
		log.Println(summary)
		if strings.Contains(stripAnsi(summary), ": failed") {
			failedTests++
		}
	}

	if failedTests > 0 {
		log.Println(color.RedString("Failed Tests:", failedTests))
	} else {
		log.Println(color.GreenString("All tests passed."))
	}

	if *summaryOutput != "" {
		markdownSummary := strings.Builder{}
		markdownSummary.WriteString("### Daml Contracts Test Summary\n\n")
		markdownSummary.WriteString("| Status | Test | Active Contracts | Transactions |\n")
		markdownSummary.WriteString("|--------|------|------------------|--------------|\n")
		for _, summary := range summaries {
			summary = stripAnsi(strings.TrimSpace(summary))
			test, status, contracts, transactions := parseTestSummary(summary)
			if status == "ok" {
				markdownSummary.WriteString(fmt.Sprintf("|🟢 %s|%s|%s|%s|\n", status, test, contracts, transactions))
			} else {
				markdownSummary.WriteString(fmt.Sprintf("|🔴 %s|%s|  |  |\n", status, test))
			}
		}
		_ = os.WriteFile(*summaryOutput, []byte(markdownSummary.String()), 0600)
	}

	// If any test failed, exit with the code of the first failed test, otherwise exit with 0
	for _, result := range results {
		if result.exitCode != 0 {
			os.Exit(result.exitCode)
		}
	}
	os.Exit(0)
}

type packageResult struct {
	path     string
	output   string
	exitCode int
}

func collectSummaries(output string) []string {
	var summaries []string
	foundSummary := 0
	for line := range strings.Lines(output) {
		line = strings.TrimSpace(line)
		// The summary section starts with "Test Summary" followed by an empty line
		if strings.Contains(line, "Test Summary") || foundSummary == 1 {
			foundSummary++
			continue
		}
		// 2 = summary start + 1 skipped line
		if foundSummary == 2 {
			// Summary section ends with "Modules internal to this package:" or an empty line
			if strings.Contains(line, "Modules internal to this package") || strings.TrimSpace(line) == "" {
				return summaries
			}
			summaries = append(summaries, strings.TrimSpace(line))
		}
	}

	return summaries
}

var re = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

// stripAnsi strips all ANSI control characters from the given string
func stripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

var summaryRegexp = regexp.MustCompile(`^(.*):\s(failed|ok)(?:, (\d*) active contracts)?(?:, (\d*) transactions.)?$`)

func parseTestSummary(line string) (test string, status string, activeContracts string, transactions string) {
	matches := summaryRegexp.FindStringSubmatch(line)
	if len(matches) < 5 {
		return
	}

	return matches[1], matches[2], matches[3], matches[4]
}
