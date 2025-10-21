package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml_lf"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/render"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/render/tmpl"
)

func main() {
	inputDar := flag.String("dar", "", "input dar to parse")
	outputPath := flag.String("o", "", "output file path")
	goPackagePrefix := flag.String("prefix", "", "package prefix")
	flag.Parse()

	if inputDar == nil || *inputDar == "" {
		panic("inputDar is required")
	}
	if outputPath == nil || *outputPath == "" {
		panic("output file path is required")
	}
	if goPackagePrefix == nil || *goPackagePrefix == "" {
		panic("prefix is required")
	}

	f, err := os.ReadFile(*inputDar)
	if err != nil {
		panic(err)
	}

	darFile := &daml_lf.DarFile{}
	err = darFile.FromBytes(f)
	if err != nil {
		panic(err)
	}

	// Render
	ad, err := render.ArchiveDataFromLF(darFile)
	if err != nil {
		panic(err)
	}
	ctx := ad.GetGeneratorContext(context.Background(), *goPackagePrefix)

	tpl := template.Must(template.New("").Parse(tmpl.EmbeddedTemplate))

	for _, module := range ad.MainPackage.Modules {
		moduleInfo, err := ctx.GetModuleInformation(ad.MainPackage.PackageId, module.Name)
		if err != nil {
			panic(err)
		}
		outputFile := filepath.Join(*outputPath, moduleInfo.RelPath, fmt.Sprintf("%s.go", moduleInfo.Filename))

		buffer := new(bytes.Buffer)
		moduleData, err := module.GetModuleData(ctx, ad.MainPackage, ad.MainPackage.PackageId)
		if err != nil {
			panic(err)
		}
		err = tpl.Execute(buffer, moduleData)
		if err != nil {
			panic(err)
		}

		bb := buffer.Bytes()
		formatted, err := format.Source(bb)
		if err != nil {
			panic(err)
		}

		log.Printf("Writing output to %s", outputFile)
		_ = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
		if err := os.WriteFile(outputFile, []byte(formatted), 0600); err != nil {
			panic(err)
		}
	}

	for i := len(ad.Dependencies) - 1; i >= 0; i-- {
		dependency := ad.Dependencies[i]
		if _, ok := ctx.UsedDependencies[dependency.PackageId]; !ok {
			continue
		}
		for _, module := range dependency.Modules {
			moduleInfo, err := ctx.GetModuleInformation(dependency.PackageId, module.Name)
			if err != nil {
				panic(err)
			}
			outputFile := filepath.Join(*outputPath, moduleInfo.RelPath, fmt.Sprintf("%s.go", moduleInfo.Filename))

			buffer := new(bytes.Buffer)
			moduleData, err := module.GetModuleData(ctx, dependency, dependency.PackageId)
			if err != nil {
				panic(err)
			}
			err = tpl.Execute(buffer, moduleData)
			if err != nil {
				panic(err)
			}

			bb := buffer.Bytes()
			formatted, err := format.Source(bb)
			if err != nil {
				panic(err)
			}

			log.Printf("Writing output to %s", outputFile)
			_ = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
			if err := os.WriteFile(outputFile, []byte(formatted), 0600); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Done")
}
