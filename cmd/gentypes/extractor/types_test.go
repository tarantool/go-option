package extractor

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestExtractTypeSpecsFromPackage(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		pkg := &packages.Package{}
		typeSpecs := ExtractTypeSpecsFromPackage(pkg)
		require.Len(t, typeSpecs, 0)
	})

	t.Run("single file, single type", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
			},
		}

		typeSpecs := ExtractTypeSpecsFromPackage(pkg)
		require.Len(t, typeSpecs, 1)
		require.Equal(t, "T", typeSpecs[0].Name.String())
	})

	t.Run("multiple files, multiple types", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
				astFromString(t, s("package pkg", "type U[K any, V any] struct{}")),
			},
		}

		typeSpecs := ExtractTypeSpecsFromPackage(pkg)
		require.Len(t, typeSpecs, 2)
		require.Equal(t, "T", typeSpecs[0].Name.String())
		require.Equal(t, "U", typeSpecs[1].Name.String())
	})
}
