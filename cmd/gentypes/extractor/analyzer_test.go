package extractor_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-option/cmd/gentypes/extractor"
)

type MockPackage struct {
	NameValue    string
	PkgPathValue string
	SyntaxValue  []*ast.File
}

func (p *MockPackage) Name() string {
	return p.NameValue
}

func (p *MockPackage) PkgPath() string {
	return p.PkgPathValue
}

func (p *MockPackage) Syntax() []*ast.File {
	return p.SyntaxValue
}

func TestNewAnalyzerFromPackage_Success(t *testing.T) {
	t.Parallel()

	pkg := &MockPackage{
		SyntaxValue: []*ast.File{
			astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
		},
		NameValue:    "pkg",
		PkgPathValue: "some-pkg-path",
	}

	analyzer, err := extractor.NewAnalyzerFromPackage(pkg)
	require.NoError(t, err)
	require.NotNil(t, analyzer)
}

func TestNewAnalyzerFromPackage_PkgInfo(t *testing.T) {
	t.Parallel()

	pkg := &MockPackage{
		SyntaxValue: []*ast.File{
			astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
		},
		NameValue:    "pkg",
		PkgPathValue: "some-pkg-path",
	}

	analyzer, err := extractor.NewAnalyzerFromPackage(pkg)
	require.NoError(t, err)

	assert.Equal(t, pkg.Name(), analyzer.PackageName())
	assert.Equal(t, pkg.PkgPath(), analyzer.PackagePath())
}

func TestNewAnalyzerFromPackage_TypeInfo(t *testing.T) {
	t.Parallel()

	pkg := &MockPackage{
		SyntaxValue: []*ast.File{
			astFromString(t, s("package pkg", "type T struct{}", "func (t *T) Method() {}")),
		},
		NameValue:    "pkg",
		PkgPathValue: "some-pkg-path",
	}

	analyzer, err := extractor.NewAnalyzerFromPackage(pkg)
	require.NoError(t, err)

	entry, found := analyzer.TypeSpecEntryByName("T")
	assert.True(t, found)

	assert.Equal(t, "T", entry.Name)
	assert.Equal(t, []string{"Method"}, entry.Methods)
	assert.True(t, entry.HasMethod("Method"))

	_, found = analyzer.TypeSpecEntryByName("U")
	assert.False(t, found)
}

func TestNewAnalyzerFromPackage_NilPackage(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		_, _ = extractor.NewAnalyzerFromPackage(nil)
	})
}
