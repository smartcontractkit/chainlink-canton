package codegen

import (
	"archive/zip"
	"io"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-canton/bindings/codegen/model"
	"github.com/stretchr/testify/require"
)

func TestGetMainDalfV3(t *testing.T) {
	srcPath := "../test-data/all-kinds-of-1.0.0_lf.dar"

	reader, err := zip.OpenReader(srcPath)
	require.NoError(t, err)

	manifest, err := GetManifest(reader)
	require.NoError(t, err)
	require.Equal(t, "all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948/all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948.dalf", manifest.MainDalf)
	require.NotNil(t, manifest)
	require.Equal(t, "1.0", manifest.Version)
	require.Equal(t, "damlc", manifest.CreatedBy)
	require.Equal(t, "all-kinds-of-1.0.0", manifest.Name)
	require.Equal(t, "3.3.0-snapshot.20250417.0", manifest.SdkVersion)
	require.Equal(t, "daml-lf", manifest.Format)
	require.Equal(t, "non-encrypted", manifest.Encryption)
	require.Len(t, manifest.Dalfs, 30)

	dalfFile, err := reader.Open(manifest.MainDalf)
	require.NoError(t, err)
	dalfContent, err := io.ReadAll(dalfFile)
	require.NoError(t, err)
	require.NotNil(t, dalfContent)

	ast, err := GetAST(dalfContent, manifest, nil, model.ExternalPackages{})
	require.Nil(t, err)
	require.NotEmpty(t, ast.Structs)

	// Test MappyContract template
	pkgMappy, exists := ast.Structs["MappyContract"]
	require.True(t, exists)
	require.Equal(t, pkgMappy.Name, "MappyContract")
	require.Equal(t, "Template", pkgMappy.RawType)
	require.Len(t, pkgMappy.Fields, 2)
	require.Equal(t, pkgMappy.Fields[0].Name, "operator")
	require.Equal(t, pkgMappy.Fields[1].Name, "value")

	// Test OneOfEverything template
	pkgEverything, exists := ast.Structs["OneOfEverything"]
	require.True(t, exists)
	require.Equal(t, pkgEverything.Name, "OneOfEverything")
	require.Equal(t, "Template", pkgEverything.RawType)
	require.Len(t, pkgEverything.Fields, 16) // Based on the generated output
	require.Equal(t, pkgEverything.Fields[0].Name, "operator")
	require.Equal(t, pkgEverything.Fields[1].Name, "someBoolean")
	require.Equal(t, pkgEverything.Fields[2].Name, "someInteger")

	// Test Accept struct
	pkgAccept, exists := ast.Structs["Accept"]
	require.True(t, exists)
	require.Equal(t, pkgAccept.Name, "Accept")
	require.Equal(t, "Record", pkgAccept.RawType)

	// Test Color enum
	colorStruct, exists := ast.Structs["Color"]
	require.True(t, exists)
	require.Equal(t, "Enum", colorStruct.RawType)
	require.Len(t, colorStruct.Fields, 3)
	require.Equal(t, colorStruct.Fields[0].Name, "Red")
	require.Equal(t, colorStruct.Fields[1].Name, "Green")
	require.Equal(t, colorStruct.Fields[2].Name, "Blue")

	res, err := Bind("codegen_test", ast, manifest.SdkVersion, true, false)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	testRes := "../test-data/all_kinds_of_1_0_0.go_gen"
	expectedCode, err := os.ReadFile(testRes)
	require.NoError(t, err)

	require.Equal(t, string(expectedCode), res, "generated code should match expected output")
}

func TestGetPackageName(t *testing.T) {
	require.Equal(t, "all-kinds-of",
		GetPackageName("all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948/all-kinds-of-1.0.0-6d7e83e81a0a7960eec37340f5b11e7a61606bd9161f413684bc345c3f387948.dalf"))
	require.Equal(t, "my-package",
		GetPackageName("my-package-1.0.0-1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef.dalf"))
	require.Equal(t, "my-package",
		GetPackageName("My-Package-1.0.0-1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef.dalf"))
}
