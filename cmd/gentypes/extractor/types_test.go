package extractor_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-option/cmd/gentypes/extractor"
)

func TestExtractTypeSpecsFromPackage(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{NameValue: "", PkgPathValue: "", SyntaxValue: nil}
		typeSpecs := extractor.ExtractTypeSpecsFromPackage(pkg)
		require.Empty(t, typeSpecs)
	})

	t.Run("single file, single type", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
			},
		}

		typeSpecs := extractor.ExtractTypeSpecsFromPackage(pkg)
		require.Len(t, typeSpecs, 1)
		require.Equal(t, "T", typeSpecs[0].Name.String())
	})

	t.Run("multiple files, multiple types", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
				astFromString(t, s("package pkg", "type U[K any, V any] struct{}")),
			},
		}

		typeSpecs := extractor.ExtractTypeSpecsFromPackage(pkg)
		require.Len(t, typeSpecs, 2)
		require.Equal(t, "T", typeSpecs[0].Name.String())
		require.Equal(t, "U", typeSpecs[1].Name.String())
	})
}
