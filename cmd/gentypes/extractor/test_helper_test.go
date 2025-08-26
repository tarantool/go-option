package extractor_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func s(lines ...string) string {
	return strings.Join(lines, "\n")
}

func astFromString(t *testing.T, s string) *ast.File {
	t.Helper()

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "test.go", s, parser.AllErrors)
	require.NoError(t, err)

	return f
}
