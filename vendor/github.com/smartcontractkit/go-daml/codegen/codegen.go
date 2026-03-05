package codegen

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/go-daml/codegen/astgen"
	model2 "github.com/smartcontractkit/go-daml/codegen/model"
	"github.com/stretchr/testify/assert/yaml"
)

func GetManifest(dar fs.FS) (*model2.Manifest, error) {
	file, err := dar.Open("META-INF/MANIFEST.MF")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	manifestYaml := struct {
		ManifestVersion string `yaml:"Manifest-Version"`
		CreatedBy       string `yaml:"Created-By"`
		Name            string `yaml:"Name"`
		SdkVersion      string `yaml:"Sdk-Version"`
		DalfMain        string `yaml:"Main-Dalf"`
		Dalfs           string `yaml:"Dalfs"`
		Format          string `yaml:"Format"`
		Encryption      string `yaml:"Encryption"`
	}{}

	if err := yaml.Unmarshal(b, &manifestYaml); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	manifest := &model2.Manifest{
		Version:    strings.ReplaceAll(manifestYaml.ManifestVersion, " ", ""),
		CreatedBy:  strings.ReplaceAll(manifestYaml.CreatedBy, " ", ""),
		Name:       strings.ReplaceAll(manifestYaml.Name, " ", ""),
		SdkVersion: strings.ReplaceAll(manifestYaml.SdkVersion, " ", ""),
		MainDalf:   strings.ReplaceAll(manifestYaml.DalfMain, " ", ""),
		Dalfs:      strings.Split(strings.ReplaceAll(manifestYaml.Dalfs, " ", ""), ","),
		Format:     strings.ReplaceAll(manifestYaml.Format, " ", ""),
		Encryption: strings.ReplaceAll(manifestYaml.Encryption, " ", ""),
	}

	if manifest.MainDalf == "" {
		return nil, fmt.Errorf("main-dalf not found in manifest")
	}

	return manifest, nil
}

func CodegenDalfs(dalfToProcess []string, dar fs.FS, pkgFile string, dalfManifest *model2.Manifest, generateHexCodec bool, externalPackages model2.ExternalPackages) (map[string]string, error) {
	//  ensure stable processing order across runs
	sort.Strings(dalfToProcess)

	ifcByModule := make(map[string]model2.InterfaceMap)
	result := make(map[string]string)

	// -------- 1) INTERFACES: deterministic traversal, no renaming logic change --------
	for _, dalf := range dalfToProcess {
		dalfFile, err := dar.Open(dalf)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to open dalf '%s': %s", dalf, err)
			continue
		}
		dalfContent, err := io.ReadAll(dalfFile)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read dalf '%s': %s", dalf, err)
		}

		interfaces, err := GetInterfaces(dalfContent, dalfManifest)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to extract interfaces from dalf: %s", dalf)
			continue
		}

		//  iterate interfaces deterministically
		keys := make([]string, 0, len(interfaces))
		for k := range interfaces {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			val := interfaces[key]

			equalNames := 0
			for _, ifcName := range ifcByModule {
				for ifcKey := range ifcName {
					res, found := strings.CutPrefix(ifcKey, key)
					_, atoiErr := strconv.Atoi(res)
					if found && (res == "" || atoiErr == nil) {
						equalNames++
					}
				}
			}
			if equalNames > 0 {
				equalNames++
				val.Name = fmt.Sprintf("%s%d", key, equalNames) // keep your existing suffix scheme
				// If you also want to avoid "...22" visually, change to: fmt.Sprintf("%s_%d", key, equalNames)
			}

			m, ok := ifcByModule[val.ModuleName]
			if !ok {
				m = make(model2.InterfaceMap)
				ifcByModule[val.ModuleName] = m
			}
			m[val.Name] = val
		}
	}

	// -------- 2) STRUCTS: deterministic traversal + do not mutate map while ranging --------
	allStructNames := make(map[string]int)

	for _, dalf := range dalfToProcess {
		dalfFile, err := dar.Open(dalf)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to open dalf '%s': %s", dalf, err)
			continue
		}
		dalfContent, err := io.ReadAll(dalfFile)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read dalf '%s': %s", dalf, err)
			continue
		}

		pkg, err := GetAST(dalfContent, dalfManifest, ifcByModule, externalPackages)
		if err != nil {
			return nil, fmt.Errorf("failed to generate AST: %w", err)
		}

		currentModules := make(map[string]bool)
		for _, structDef := range pkg.Structs {
			if structDef.ModuleName != "" {
				currentModules[structDef.ModuleName] = true
			}
		}

		log.Info().Msgf("adding interfaces for dalf %s from modules: %v", dalf, currentModules)

		//  iterate modules deterministically
		moduleNames := make([]string, 0, len(currentModules))
		for m := range currentModules {
			moduleNames = append(moduleNames, m)
		}
		sort.Strings(moduleNames)

		for _, moduleName := range moduleNames {
			if ifcMap, exists := ifcByModule[moduleName]; exists {
				//  iterate interfaces deterministically
				ifcKeys := make([]string, 0, len(ifcMap))
				for k := range ifcMap {
					ifcKeys = append(ifcKeys, k)
				}
				sort.Strings(ifcKeys)

				for _, key := range ifcKeys {
					val := ifcMap[key]
					log.Debug().Msgf("adding interface %s from module %s to output", key, moduleName)
					pkg.Structs[key] = val
				}
			}
		}

		//  deterministic renaming (plan + apply)
		type rename struct {
			orig string
			new  string
			def  *model2.TmplStruct
		}
		planned := make([]rename, 0)
		renamedStructs := make(map[string]*model2.TmplStruct)

		//  iterate struct keys deterministically
		structKeys := make([]string, 0, len(pkg.Structs))
		for k := range pkg.Structs {
			structKeys = append(structKeys, k)
		}
		sort.Strings(structKeys)

		for _, structName := range structKeys {
			structDef := pkg.Structs[structName]
			if structDef.IsInterface {
				continue
			}

			equalNames := 0
			for existingName := range allStructNames {
				res, found := strings.CutPrefix(existingName, structName)
				_, atoiErr := strconv.Atoi(res)
				if found && (res == "" || atoiErr == nil) {
					equalNames++
				}
			}

			if equalNames > 0 {
				equalNames++
				newName := fmt.Sprintf("%s%d", structName, equalNames) // keep your existing suffix scheme
				// If you also want to avoid "...22" visually, change to: fmt.Sprintf("%s_%d", structName, equalNames)
				planned = append(planned, rename{orig: structName, new: newName, def: structDef})
			} else {
				allStructNames[structName] = 0
			}
		}

		for _, r := range planned {
			r.def.Name = r.new
			delete(pkg.Structs, r.orig)
			pkg.Structs[r.new] = r.def
			renamedStructs[r.orig] = r.def
			allStructNames[r.new] = 1
		}

		// Update references (unchanged)
		for _, structDef := range pkg.Structs {
			for _, field := range structDef.Fields {
				if renamed, exists := renamedStructs[field.Type.GoType()]; exists {
					field.Type = model2.Unknown{String: renamed.Name}
				}
				trimmedType := strings.TrimPrefix(field.Type.GoType(), "*")
				trimmedType = strings.TrimPrefix(trimmedType, "[]")
				if renamed, exists := renamedStructs[trimmedType]; exists {
					field.Type = model2.Unknown{String: strings.Replace(field.Type.GoType(), trimmedType, renamed.Name, 1)}
				}
			}

			for _, choice := range structDef.Choices {
				if renamed, exists := renamedStructs[choice.ArgType.GoType()]; exists {
					choice.ArgType = model2.Unknown{String: renamed.Name}
				}
				// if renamed, exists := renamedStructs[choice.ReturnType.GoType()]; exists {
				// 	choice.ReturnType = model2.Unknown{String: renamed.Name}
				// }
			}
		}

		code, err := Bind(pkgFile, pkg, dalfManifest.SdkVersion, dalf == dalfManifest.MainDalf, generateHexCodec)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Go code: %w", err)
		}

		result[dalf] = code
	}

	return result, nil
}

func GetInterfaces(payload []byte, manifest *model2.Manifest) (map[string]*model2.TmplStruct, error) {
	var version string
	if strings.HasPrefix(manifest.SdkVersion, astgen.V3) {
		version = astgen.V3
	} else if strings.HasPrefix(manifest.SdkVersion, astgen.V2) || strings.HasPrefix(manifest.SdkVersion, astgen.V1) {
		version = astgen.V2
	} else {
		return nil, fmt.Errorf("unsupported sdk version %s", manifest.SdkVersion)
	}

	gen, err := astgen.GetAstGenFromVersion(payload, model2.ExternalPackages{}, version)
	if err != nil {
		return nil, err
	}

	return gen.GetInterfaces()
}

func GetAST(payload []byte, manifest *model2.Manifest, ifcByModule map[string]model2.InterfaceMap, externalPackages model2.ExternalPackages) (*model2.Package, error) {
	var version string
	if strings.HasPrefix(manifest.SdkVersion, astgen.V3) {
		version = astgen.V3
	} else if strings.HasPrefix(manifest.SdkVersion, astgen.V2) || strings.HasPrefix(manifest.SdkVersion, astgen.V1) {
		version = astgen.V2
	} else {
		return nil, fmt.Errorf("unsupported sdk version %s", manifest.SdkVersion)
	}

	gen, err := astgen.GetAstGenFromVersion(payload, externalPackages, version)
	if err != nil {
		return nil, err
	}
	structs, importedPackages, err := gen.GetTemplateStructs(ifcByModule)
	if err != nil {
		return nil, err
	}

	packageID := GetPackageID(manifest.MainDalf)
	if packageID == "" {
		return nil, fmt.Errorf("could not extract package ID from MainDalf: %s", manifest.MainDalf)
	}

	packageName := manifest.Name
	if packageName == "" {
		packageName = GetPackageName(manifest.MainDalf)
	}
	// Always strip version from package name, regardless of source
	packageName = stripVersionFromPackageName(packageName)

	// Collect all imported packages into a slice and sort them by import path for deterministic output
	importedPackagesSlice := slices.Collect(maps.Values(importedPackages.Packages))
	slices.SortFunc(importedPackagesSlice, func(a, b model2.ExternalPackage) int {
		return strings.Compare(a.Import, b.Import)
	})

	fmt.Println(packageID, packageName)

	return &model2.Package{
		Name:             packageName,
		PackageID:        packageID,
		Structs:          structs,
		ImportedPackages: importedPackagesSlice,
	}, nil
}

func GetPackageID(dalf string) string {
	parts := strings.Split(dalf, "/")
	filename := strings.TrimSuffix(parts[len(parts)-1], ".dalf")

	lastHyphen := strings.LastIndex(filename, "-")
	if lastHyphen != -1 && lastHyphen < len(filename)-1 {
		return filename[lastHyphen+1:]
	}

	return ""
}

func stripVersionFromPackageName(name string) string {
	// Strip version pattern like "-1.0.0", "-2.9.1", etc.
	// Pattern: hyphen followed by digits and dots (version number)
	versionPattern := regexp.MustCompile(`-\d+(\.\d+)*$`)
	return versionPattern.ReplaceAllString(name, "")
}

func GetPackageName(dalf string) string {
	parts := strings.Split(dalf, "/")
	filename := strings.TrimSuffix(parts[len(parts)-1], ".dalf")

	lastHyphen := strings.LastIndex(filename, "-")
	if lastHyphen == -1 {
		return strings.ToLower(filename)
	}

	potentialHash := filename[lastHyphen+1:]
	if len(potentialHash) == 64 {
		allHex := true
		for _, ch := range potentialHash {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				allHex = false
				break
			}
		}
		if allHex {
			filename = filename[:lastHyphen]
		}
	}

	// Strip version pattern like "-1.0.0", "-2.9.1", etc.
	filename = stripVersionFromPackageName(filename)

	return strings.ToLower(filename)
}
