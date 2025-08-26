package extractor

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestNewAnalyzerFromPackage(t *testing.T) {
	pkg := &packages.Package{
		Syntax: []*ast.File{
			astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
		},
		Name:    "pkg",
		PkgPath: "some-pkg-path",
	}

	analyzer, err := NewAnalyzerFromPackage(pkg)
	require.NoError(t, err)
	require.NotNil(t, analyzer)

	assert.Equal(t, pkg.Name, analyzer.PackageName())
	assert.Equal(t, pkg.PkgPath, analyzer.PackagePath())

	entry, ok := analyzer.TypeSpecEntryByName("T")
	assert.True(t, ok)

	assert.Equal(t, "T", entry.Name)
	assert.Equal(t, []string{"Method"}, entry.Methods)
	assert.True(t, entry.HasMethod("Method"))

	_, ok = analyzer.TypeSpecEntryByName("U")
	assert.False(t, ok)
}
