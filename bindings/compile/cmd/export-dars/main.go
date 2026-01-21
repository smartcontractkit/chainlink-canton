// nolint
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/compile"
	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
)

func main() {
	outDir := flag.String("out", "artifacts/dars", "output directory for compiled DARs")
	includeTests := flag.Bool("include-tests", false, "include *_test packages (ccip_test, mcms_test)")
	only := flag.String("only", "", "comma-separated list of contracts.Package to export (e.g. coin,mcms,ccip_common)")
	flag.Parse()

	if *outDir == "" {
		log.Fatal("out directory is required")
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("failed to create out dir %q: %v", *outDir, err)
	}

	onlySet := map[contracts.Package]bool{}
	if strings.TrimSpace(*only) != "" {
		for _, s := range strings.Split(*only, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			onlySet[contracts.Package(s)] = true
		}
	}

	// stable order (maps are random)
	pkgs := make([]contracts.Package, 0, len(contracts.Contracts))
	for p := range contracts.Contracts {
		// filter test packages unless requested
		if !*includeTests && (p == contracts.CCIPTest || p == contracts.MCMSTest) {
			continue
		}
		// filter if --only is set
		if len(onlySet) > 0 && !onlySet[p] {
			continue
		}
		pkgs = append(pkgs, p)
	}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i] < pkgs[j] })

	if len(pkgs) == 0 {
		log.Fatal("no packages selected to export (check --only / --include-tests)")
	}

	wrote := 0
	for _, p := range pkgs {
		fmt.Printf("Compiling %s...\n", p)

		cp, err := compile.Package(p)
		if err != nil {
			log.Fatalf("failed compiling %s: %v", p, err)
		}
		if len(cp.Dar) == 0 {
			log.Fatalf("empty DAR output for %s", p)
		}

		// cp.PackageName comes from daml.yaml "name:" (e.g. coin, mcms, ccip-common, ...)
		outPath := filepath.Join(*outDir, fmt.Sprintf("%s.dar", cp.PackageName))
		if err := os.WriteFile(outPath, cp.Dar, 0o644); err != nil {
			log.Fatalf("failed writing %q: %v", outPath, err)
		}
		fmt.Printf("Wrote %s\n", outPath)
		wrote++
	}

	fmt.Printf("Done. Wrote %d DAR(s) to %s\n", wrote, *outDir)
}
