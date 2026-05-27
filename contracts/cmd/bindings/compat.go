package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type generatedPackage struct {
	dir        []string
	importAs   string
	packageAs  string
	sourceFile string
}

type exportedSymbols struct {
	consts []string
	types  []string
	funcs  []string
}

func writeCompatibilityPackages(artifactsDir, basePath string) error {
	singleTargetPackages := []struct {
		outDir     []string
		packageAs  string
		target     generatedPackage
		outputFile string
	}{
		{[]string{"ccip", "common"}, "common", generatedPackage{dir: []string{"ccip", "core"}, importAs: "core"}, "common.go"},
		{[]string{"ccip", "client"}, "client", generatedPackage{dir: []string{"ccip", "core"}, importAs: "core"}, "client.go"},
		{[]string{"ccip", "feequoter"}, "feequoter", generatedPackage{dir: []string{"ccip", "core"}, importAs: "core"}, "feequoter.go"},
		{[]string{"ccip", "rmn"}, "rmn", generatedPackage{dir: []string{"ccip", "core"}, importAs: "core"}, "rmn.go"},
		{[]string{"ccip", "tokenadminregistry"}, "tokenadminregistry", generatedPackage{dir: []string{"ccip", "core"}, importAs: "core"}, "tokenadminregistry.go"},
		{[]string{"ccip", "offramp"}, "offramp", generatedPackage{dir: []string{"ccip", "ccipruntime"}, importAs: "ccipruntime"}, "offramp.go"},
		{[]string{"ccip", "onramp"}, "onramp", generatedPackage{dir: []string{"ccip", "ccipruntime"}, importAs: "ccipruntime"}, "onramp.go"},
		{[]string{"ccip", "perpartyrouter"}, "perpartyrouter", generatedPackage{dir: []string{"ccip", "ccipruntime"}, importAs: "ccipruntime"}, "perpartyrouter.go"},
		{[]string{"ccip", "ccipmain"}, "ccipmain", generatedPackage{dir: []string{"ccip", "ccipruntime"}, importAs: "ccipruntime"}, "ccipmain.go"},
		{[]string{"ccip", "interfaces"}, "interfaces", generatedPackage{dir: []string{"ccip", "extensionapi"}, importAs: "extensionapi"}, "interfaces.go"},
		{[]string{"ccip", "ccvs"}, "ccvs", generatedPackage{dir: []string{"ccip", "committeeverifier"}, importAs: "committeeverifier"}, "ccvs.go"},
		{[]string{"ccip", "ccipsender"}, "ccipsender", generatedPackage{dir: []string{"ccip", "sender"}, importAs: "sender"}, "ccipsender.go"},
		{[]string{"ccip", "ccipreceiver"}, "ccipreceiver", generatedPackage{dir: []string{"ccip", "receiver"}, importAs: "receiver"}, "ccipreceiver.go"},
		{[]string{"chainlink", "instanceapi"}, "instanceapi", generatedPackage{dir: []string{"chainlink", "chainlinkapi"}, importAs: "chainlinkapi"}, "instanceapi.go"},
	}
	for _, pkg := range singleTargetPackages {
		if err := writeSingleTargetCompatibilityPackage(artifactsDir, basePath, pkg.outDir, pkg.packageAs, pkg.target, pkg.outputFile); err != nil {
			return err
		}
	}

	return writeMCMSCompatibilityPackage(artifactsDir, basePath)
}

func writeSingleTargetCompatibilityPackage(artifactsDir, basePath string, outDir []string, packageAs string, target generatedPackage, outputFile string) error {
	target.sourceFile = generatedPackageFile(artifactsDir, target.dir)
	symbols, err := exportedIdentifiers(target.sourceFile)
	if err != nil {
		return err
	}

	var b strings.Builder
	writeGeneratedHeader(&b, packageAs)
	fmt.Fprintf(&b, "import %s %q\n\n", target.importAs, generatedPackageImport(basePath, target.dir))
	writeAliases(&b, target.importAs, symbols, nil)

	return writeCompatibilityFile(artifactsDir, outDir, outputFile, b.String())
}

func writeMCMSCompatibilityPackage(artifactsDir, basePath string) error {
	targets := []generatedPackage{
		{dir: []string{"mcms", "core"}, importAs: "core", packageAs: "core"},
		{dir: []string{"mcms", "api"}, importAs: "api", packageAs: "api"},
	}
	for i := range targets {
		targets[i].sourceFile = generatedPackageFile(artifactsDir, targets[i].dir)
	}

	skip := map[string]bool{
		"Contract":    true,
		"MCMSEncoder": true,
		"NewContract": true,
	}
	seenConsts := map[string]bool{}
	seenTypes := map[string]bool{"RawInstanceAddress": true}
	seenFuncs := map[string]bool{}

	var b strings.Builder
	writeGeneratedHeader(&b, "mcms")
	fmt.Fprintf(&b, "import (\n")
	fmt.Fprintf(&b, "\tapi %q\n", generatedPackageImport(basePath, []string{"mcms", "api"}))
	fmt.Fprintf(&b, "\tcore %q\n", generatedPackageImport(basePath, []string{"mcms", "core"}))
	fmt.Fprintf(&b, "\tchainlinkapi %q\n", generatedPackageImport(basePath, []string{"chainlink", "chainlinkapi"}))
	fmt.Fprintf(&b, ")\n\n")
	fmt.Fprintf(&b, "type RawInstanceAddress = chainlinkapi.RawInstanceAddress\n\n")

	for _, target := range targets {
		symbols, err := exportedIdentifiers(target.sourceFile)
		if err != nil {
			return err
		}
		writeAliases(&b, target.importAs, symbols, aliasSkipSet(skip, seenConsts, seenTypes, seenFuncs))
	}

	fmt.Fprintf(&b, "type MCMSEncoder interface {\n\tcore.MCMSEncoder\n\tapi.MCMSEncoder\n\tSetConfigParams(args SetConfigParams) (*bind.EncodedChoice, error)\n}\n\n")
	fmt.Fprintf(&b, "type Contract struct {\n\tcore     *core.Contract\n\tapi      *api.Contract\n\ttemplate *bind.BoundTemplate\n}\n\n")
	fmt.Fprintf(&b, "func NewContract(packageID, moduleName, templateName string) *Contract {\n")
	fmt.Fprintf(&b, "\treturn &Contract{\n")
	fmt.Fprintf(&b, "\t\tcore:     core.NewContract(packageID, moduleName, templateName),\n")
	fmt.Fprintf(&b, "\t\tapi:      api.NewContract(packageID, moduleName, templateName),\n")
	fmt.Fprintf(&b, "\t\ttemplate: bind.NewBoundTemplate(packageID, moduleName, templateName),\n")
	fmt.Fprintf(&b, "\t}\n")
	fmt.Fprintf(&b, "}\n\n")
	fmt.Fprintf(&b, "func (c *Contract) Encoder() MCMSEncoder {\n")
	fmt.Fprintf(&b, "\treturn encoder{coreEncoder: c.core.Encoder(), apiEncoder: c.api.Encoder(), template: c.template}\n")
	fmt.Fprintf(&b, "}\n\n")
	fmt.Fprintf(&b, "type encoder struct {\n\tcoreEncoder core.MCMSEncoder\n\tapiEncoder  api.MCMSEncoder\n\ttemplate    *bind.BoundTemplate\n}\n\n")
	fmt.Fprintf(&b, "func (e encoder) BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error) { return e.apiEncoder.BypasserExecuteBatch(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error) { return e.apiEncoder.CancelBatch(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error) { return e.coreEncoder.CanExecuteOp(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) ExecuteOp(args ExecuteOp) (*bind.EncodedChoice, error) { return e.coreEncoder.ExecuteOp(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) ExecuteScheduledBatch(args ExecuteScheduledBatch) (*bind.EncodedChoice, error) { return e.coreEncoder.ExecuteScheduledBatch(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) GetBlockedFunctions(args GetBlockedFunctions) (*bind.EncodedChoice, error) { return e.coreEncoder.GetBlockedFunctions(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) GetBlockedFunctionsCount(args GetBlockedFunctionsCount) (*bind.EncodedChoice, error) { return e.coreEncoder.GetBlockedFunctionsCount(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error) { return e.coreEncoder.GetMinDelay(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) GetState(args GetState) (*bind.EncodedChoice, error) { return e.coreEncoder.GetState(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error) { return e.coreEncoder.GetTimestamp(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) IsOperation(args IsOperation) (*bind.EncodedChoice, error) { return e.coreEncoder.IsOperation(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) IsOperationDone(args IsOperationDone) (*bind.EncodedChoice, error) { return e.coreEncoder.IsOperationDone(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) IsOperationPending(args IsOperationPending) (*bind.EncodedChoice, error) { return e.coreEncoder.IsOperationPending(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) IsOperationReady(args IsOperationReady) (*bind.EncodedChoice, error) { return e.coreEncoder.IsOperationReady(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) PruneSeenHashes(args PruneSeenHashes) (*bind.EncodedChoice, error) { return e.coreEncoder.PruneSeenHashes(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) PruneTimelockTimestamps(args PruneTimelockTimestamps) (*bind.EncodedChoice, error) { return e.coreEncoder.PruneTimelockTimestamps(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error) { return e.apiEncoder.ScheduleBatch(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) { return e.coreEncoder.SetConfig(args) }\n")
	fmt.Fprintf(&b, "func (e encoder) SetConfigParams(args SetConfigParams) (*bind.EncodedChoice, error) { return e.template.EncodeChoiceArgs(\"SetConfig\", args) }\n")
	fmt.Fprintf(&b, "func (e encoder) SetRoot(args SetRoot) (*bind.EncodedChoice, error) { return e.coreEncoder.SetRoot(args) }\n")

	content := b.String()
	content = strings.Replace(content, "import (\n", "import (\n\t\"github.com/smartcontractkit/go-daml/pkg/bind\"\n", 1)

	return writeCompatibilityFile(artifactsDir, []string{"mcms"}, "mcms.go", content)
}

func writeAliases(b *strings.Builder, importAs string, symbols exportedSymbols, skip func(kind, name string) bool) {
	for _, name := range symbols.consts {
		if skip != nil && skip("const", name) {
			continue
		}
		fmt.Fprintf(b, "const %s = %s.%s\n", name, importAs, name)
	}
	if len(symbols.consts) > 0 {
		fmt.Fprintln(b)
	}
	for _, name := range symbols.types {
		if skip != nil && skip("type", name) {
			continue
		}
		fmt.Fprintf(b, "type %s = %s.%s\n", name, importAs, name)
	}
	if len(symbols.types) > 0 {
		fmt.Fprintln(b)
	}
	for _, name := range symbols.funcs {
		if skip != nil && skip("func", name) {
			continue
		}
		fmt.Fprintf(b, "var %s = %s.%s\n", name, importAs, name)
	}
	if len(symbols.funcs) > 0 {
		fmt.Fprintln(b)
	}
}

func aliasSkipSet(alwaysSkip map[string]bool, seenConsts, seenTypes, seenFuncs map[string]bool) func(kind, name string) bool {
	return func(kind, name string) bool {
		if alwaysSkip[name] {
			return true
		}
		switch kind {
		case "const":
			if seenConsts[name] {
				return true
			}
			seenConsts[name] = true
		case "type":
			if seenTypes[name] {
				return true
			}
			seenTypes[name] = true
		case "func":
			if seenFuncs[name] {
				return true
			}
			seenFuncs[name] = true
		}
		return false
	}
}

func exportedIdentifiers(filename string) (exportedSymbols, error) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		return exportedSymbols{}, fmt.Errorf("parse generated bindings %q: %w", filename, err)
	}

	symbols := exportedSymbols{}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.CONST:
				for _, spec := range d.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range valueSpec.Names {
						if ast.IsExported(name.Name) {
							symbols.consts = append(symbols.consts, name.Name)
						}
					}
				}
			case token.TYPE:
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if ok && ast.IsExported(typeSpec.Name.Name) {
						symbols.types = append(symbols.types, typeSpec.Name.Name)
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv == nil && ast.IsExported(d.Name.Name) {
				symbols.funcs = append(symbols.funcs, d.Name.Name)
			}
		}
	}

	symbols.consts = sortedUnique(symbols.consts)
	symbols.types = sortedUnique(symbols.types)
	symbols.funcs = sortedUnique(symbols.funcs)

	return symbols, nil
}

func sortedUnique(in []string) []string {
	sort.Strings(in)
	out := in[:0]
	var last string
	for i, v := range in {
		if i == 0 || v != last {
			out = append(out, v)
			last = v
		}
	}
	return out
}

func generatedPackageFile(artifactsDir string, dir []string) string {
	return filepath.Join(append([]string{artifactsDir}, append(dir, dir[len(dir)-1]+".go")...)...)
}

func generatedPackageImport(basePath string, dir []string) string {
	return basePath + "/" + strings.Join(dir, "/")
}

func writeGeneratedHeader(b *strings.Builder, packageAs string) {
	fmt.Fprintf(b, "package %s\n\n", packageAs)
	fmt.Fprintf(b, "// Code generated by contracts/cmd/bindings. DO NOT EDIT.\n\n")
}

func writeCompatibilityFile(artifactsDir string, dir []string, filename string, content string) error {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format compatibility package %q: %w", filepath.Join(append([]string{artifactsDir}, append(dir, filename)...)...), err)
	}
	outputDir := filepath.Join(append([]string{artifactsDir}, dir...)...)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create compatibility package dir %q: %w", outputDir, err)
	}
	outputFile := filepath.Join(outputDir, filename)
	if err := os.WriteFile(outputFile, formatted, 0o644); err != nil {
		return fmt.Errorf("write compatibility package %q: %w", outputFile, err)
	}
	return nil
}
