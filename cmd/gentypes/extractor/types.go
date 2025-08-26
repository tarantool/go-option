package extractor

import (
	"go/ast"
	"go/token"
)

type typeSpecVisitor struct {
	Types []*ast.TypeSpec
}

func (t *typeSpecVisitor) Visit(node ast.Node) ast.Visitor {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.TYPE {
		return t
	}

	for _, spec := range genDecl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		t.Types = append(t.Types, ts)
	}

	return nil
}

// ExtractTypeSpecsFromPackage extracts type specs from a ast tree.
func ExtractTypeSpecsFromPackage(pkg Package) []*ast.TypeSpec {
	visitor := &typeSpecVisitor{
		Types: nil,
	}
	for _, file := range pkg.Syntax() {
		ast.Walk(visitor, file)
	}

	return visitor.Types
}
