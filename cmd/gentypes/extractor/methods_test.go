package extractor_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-option/cmd/gentypes/extractor"
)

func TestExtractMethodsFromPackageSimple(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{NameValue: "", PkgPathValue: "", SyntaxValue: nil}
		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Empty(t, funcDecls)
	})

	t.Run("single file, zero methods", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}")),
			},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Empty(t, funcDecls)
	})

	t.Run("single file, single method", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
			},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "Method", funcDecls[0].Name.Name)
	})
}

func TestExtractMethodsFromPackageMultiple(t *testing.T) {
	t.Parallel()

	t.Run("multiple files, couple of methods", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method1() {}")),
				astFromString(t, s("package pkg", "func (t *T) Method2() {}")),
			},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 2)
		assert.Equal(t, "Method1", funcDecls[0].Name.Name)
		assert.Equal(t, "Method2", funcDecls[1].Name.Name)
	})

	t.Run("function is ignored", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{
				astFromString(t, s("package pkg", "func Method() {}")),
			},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Empty(t, funcDecls)
	})
}

func TestExtractRecvTypeNameSimple(t *testing.T) {
	t.Parallel()

	t.Run("method", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T struct{}", "func (t T) Method() {}"),
			)},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})

	t.Run("ptr method", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T struct{}", "func (t *T) Method() {}"),
			)},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})
}

func TestExtractRecvTypeNameGenericSingle(t *testing.T) {
	t.Parallel()

	t.Run("single-type generic method", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any] struct{}", "func (t T[K]) Method() {}"),
			)},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})

	t.Run("single-type generic method with ptr receiver", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any] struct{}", "func (t *T[K]) Method() {}"),
			)},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})
}

func TestExtractRecvTypeNameGenericMulti(t *testing.T) {
	t.Parallel()

	t.Run("multi-type generic method", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any, V any] struct{}", "func (t T[K, V]) Method() {}"),
			)},
		}

		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})

	t.Run("multi-type generic method with ptr receiver", func(t *testing.T) {
		t.Parallel()

		pkg := &MockPackage{
			NameValue:    "",
			PkgPathValue: "",
			SyntaxValue: []*ast.File{astFromString(t,
				s("package pkg", "type T[K any, V any] struct{}", "func (t *T[K, V]) Method() {}"),
			)},
		}
		funcDecls := extractor.ExtractMethodsFromPackage(pkg)
		require.Len(t, funcDecls, 1)
		assert.Equal(t, "T", extractor.ExtractRecvTypeName(funcDecls[0]))
	})
}
