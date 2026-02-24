package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func main() {
	ctx := context.Background()

	rootDir := flag.String("root", "", "Path to the contracts root directory, must contain multi-package.yaml file, defaults to current working directory")
	color := flag.Bool("color", false, "Whether or not to use color output")
	plain := flag.Bool("plain", false, "Whether or not to use plain output")
	summaryOutput := flag.String("output", "", "Path to output summary file")

	flag.Parse()

	fmt.Println("Testing packages:")
	var packages []string
	err := filepath.WalkDir(*rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip .daml directories
		if d.IsDir() && d.Name() == ".daml" {
			return filepath.SkipDir
		}
		if d.Name() == "daml.yaml" {
			fmt.Println(path)
			packages = append(packages, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	var progress *progressbar.ProgressBar
	if !*plain {
		progress = progressbar.NewOptions(len(packages),
			progressbar.OptionSetDescription("Testing packages..."),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionEnableColorCodes(*color),
			progressbar.OptionSetPredictTime(false),
		)
	} else {
		fmt.Println("Running tests...")
	}
	results := make([]packageResult, len(packages))
	for i, s := range packages {
		wg.Go(func() {
			directory, _ := filepath.Split(s)
			cmd := exec.CommandContext(ctx, "dpm", "test")
			if *color {
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
				fmt.Println("✓ Done: ", directory)
			}
		})
	}
	wg.Wait()

	var summaries []string
	fmt.Println("Raw Outputs:")
	for _, r := range results {
		fmt.Println("==============================")
		fmt.Printf("Package: %s\n", r.path)
		fmt.Printf("Exit Code: %d\n\n", r.exitCode)
		fmt.Println(r.output)

		// Collect summary
		for _, summary := range collectSummaries(r.output) {
			summaries = append(summaries, r.path+summary)
		}
	}
	fmt.Println("==============================")
	fmt.Println("Test Summary:")
	failedTests := 0
	for _, summary := range summaries {
		fmt.Println(summary)
		if strings.Contains(stripAnsi(summary), ": failed") {
			failedTests++
		}
	}
	fmt.Println()
	if failedTests > 0 {
		fmt.Println("Failed Tests:", failedTests)
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
		_ = os.WriteFile(*summaryOutput, []byte(markdownSummary.String()), 0644)
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
