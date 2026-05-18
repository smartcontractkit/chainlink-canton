// Command generate_all_contract_reports generates HTML reports for all CCIP contracts.
//
// This script wraps the main fetch_active_contract_by_instance_address tool
// to generate reports for all contracts in a single run.
//
// Usage:
//   go run ./scripts/staging/generate_all_contract_reports.go
//
// Or build and run:
//   go build -o generate_reports ./scripts/staging/generate_all_contract_reports.go
//   ./generate_reports

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Party constants
const (
	partyCCIPOwner          = "ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51"
	partyCCIPBootstrapOwner = "ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d"
)

// ContractInfo holds the metadata for each contract to fetch
type ContractInfo struct {
	TemplateID      string `json:"template_id"`
	InstanceAddress string `json:"instance_address"`
	InstanceID      string `json:"instance_id"`
	Type            string `json:"type"`
	Party           string `json:"party"` // Optional: defaults to ccipOwner if empty
}

// List of all contracts to fetch
var contractsToFetch = []ContractInfo{
	{
		TemplateID:      "#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter",
		InstanceAddress: `{"address":"0x3891327bf89b1621f67a720f73f8478777f2c106d95e570c5fa388f138bc0728"}`,
		InstanceID:      "feequoter-shywn",
		Type:            "FeeQuoter",
	},
	{
		TemplateID:      "#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier",
		InstanceAddress: `{"address":"0xf11b7b25ed8ac60beecb78e58fba954dd9b75f13b1b67ff0983b55aab52dfcd1"}`,
		InstanceID:      "committeeverifier-suoid",
		Type:            "CommitteeVerifier",
	},
	{
		TemplateID:      "#ccip-offramp:CCIP.OffRamp:OffRamp",
		InstanceAddress: `{"address":"0xe9c3534382c638dbd457aa92becdc61cb6c294795e176365baaa06be3dd885fa"}`,
		InstanceID:      "offramp-uaxss",
		Type:            "OffRamp",
	},
	{
		TemplateID:      "#ccip-onramp:CCIP.OnRamp:OnRamp",
		InstanceAddress: `{"address":"0x92b53bcb058aabfc52cb617230375b5dacf8bc19932de5a9f56df659e4944c7b"}`,
		InstanceID:      "onramp-tlspm",
		Type:            "OnRamp",
	},
	{
		TemplateID:      "#ccip-executor:CCIP.Executor:Executor",
		InstanceAddress: `{"address":"0xa3fecf9edeb0686bf58e17b4765a5806ff934ff8efb145a42c965a79a32f875c"}`,
		InstanceID:      "executor-zzpfy",
		Type:            "Executor",
	},
	{
		TemplateID:      "#ccip-lockreleasetokenpool:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
		InstanceAddress: `{"address":"0x9771c1e34476f3f3468c8bec25b6ac9c67bc1e43a86dc37b97cc3198382a0005"}`,
		InstanceID:      "lockreleasetokenpool-aswyq",
		Type:            "LockReleaseTokenPool",
		Party:           partyCCIPBootstrapOwner, // LockReleaseTP uses bootstrap owner
	},
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Get working directory (repo root)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Create reports directory
	reportsDir := filepath.Join(wd, "scripts", "staging", "fetch_active_contract_by_instance_address", "reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("create reports directory: %w", err)
	}

	fmt.Printf("Generating reports for %d contracts...\n", len(contractsToFetch))
	fmt.Printf("Reports will be saved to: %s\n\n", reportsDir)

	// Generate report for each contract
	successCount := 0
	failCount := 0

	for i, contract := range contractsToFetch {
		fmt.Printf("[%d/%d] Fetching %s (%s)...\n", i+1, len(contractsToFetch), contract.Type, contract.InstanceID)

		reportFile := filepath.Join(reportsDir, fmt.Sprintf("%s.html", contract.InstanceID))

		if err := generateReport(contract, reportFile); err != nil {
			fmt.Printf("    ❌ ERROR: %v\n", err)
			failCount++
			continue
		}

		fmt.Printf("    ✓ Saved to %s\n", reportFile)
		successCount++
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("✅ Done! %d succeeded, %d failed\n", successCount, failCount)
	fmt.Printf("Reports directory: %s\n", reportsDir)
	fmt.Printf("========================================\n")

	return nil
}

func generateReport(contract ContractInfo, outputFile string) error {
	// Determine party to use (default to ccipOwner if not specified)
	party := contract.Party
	if party == "" {
		party = partyCCIPOwner
	}

	// Build the command to run the main fetch tool
	cmd := exec.Command(
		"go", "run",
		"./scripts/staging/fetch_active_contract_by_instance_address",
		"--format", "html",
		"--html-out", outputFile,
		"--template", contract.TemplateID,
		"--instance-address", contract.InstanceAddress,
		"--instance-id", contract.InstanceID,
		"--party", party,
	)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	// Print the ACS scan messages
	outputStr := string(output)
	if outputStr != "" {
		// Print non-empty lines
		lines := splitLines(outputStr)
		for _, line := range lines {
			if line != "" {
				fmt.Printf("    %s\n", line)
			}
		}
	}

	return nil
}

func splitLines(s string) []string {
	var lines []string
	var current string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// PrettyPrint prints the contracts list as formatted JSON
func PrettyPrint(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}
