package extractor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func astFromString(t *testing.T, s string) *ast.File {
	t.Helper()

	f, err := parser.ParseFile(token.NewFileSet(), "", s, 0)
	require.NoError(t, err)

	return f
}

func s(in ...string) string {
	return strings.Join(in, "\n")
}

func TestExtractMethodsFromPackage(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		pkg := &packages.Package{}
		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 0)
	})

	t.Run("single file, zero methods", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
			},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 0)
	})

	t.Run("single file, single method", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
			},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "Method", funcDecls[0].Name.Name)
	})

	t.Run("multiple files, couple of methods", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method1() {}")),
				astFromString(t, s("package pkg", "func (t *T) Method2() {}")),
			},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 2)
		assert.Equal(t, "Method1", funcDecls[0].Name.Name)
		assert.Equal(t, "Method2", funcDecls[1].Name.Name)
	})

	t.Run("function is ignored", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{
				astFromString(t, s("package pkg", "func Method() {}")),
			},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 0)
	})
}

func TestExtractRecvTypeName(t *testing.T) {
	t.Run("method", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T struct{}", "func (t T) Method() {}"),
			)},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})

	t.Run("ptr method", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T struct{}", "func (t *T) Method() {}"),
			)},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})

	t.Run("single-type generic method", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any] struct{}", "func (t T[K]) Method() {}"),
			)},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})

	t.Run("multi-type generic method", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any, V any] struct{}", "func (t T[K, V]) Method() {}"),
			)},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})

	t.Run("single-type generic method with ptr receiver", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any] struct{}", "func (t *T[K]) Method() {}"),
			)},
		}

		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})

	t.Run("multi-type generic method with ptr receiver", func(t *testing.T) {
		pkg := &packages.Package{
			Syntax: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any, V any] struct{}", "func (t *T[K, V]) Method() {}"),
			)},
		}
		funcDecls := ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractRecvTypeName(funcDecls[0]))
	})
}
