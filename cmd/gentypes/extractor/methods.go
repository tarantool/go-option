package extractor

import (
	"go/ast"
)

type methodVisitor struct {
	Methods []*ast.FuncDecl
}

func (t *methodVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Recv == nil {
		return t
	}

	t.Methods = append(t.Methods, funcDecl)

	return t
}

// ExtractMethodsFromPackage is a function to extract methods from package.
func ExtractMethodsFromPackage(pkg Package) []*ast.FuncDecl {
	visitor := &methodVisitor{
		Methods: nil,
	}
	for _, file := range pkg.Syntax() {
		ast.Walk(visitor, file)
	}

	return visitor.Methods
}

// ExtractRecvTypeName is a helper function to extract receiver type name (string) from method.
func ExtractRecvTypeName(method *ast.FuncDecl) string {
	if method.Recv == nil {
		return ""
	}

	name := method.Recv.List[0]
	tpExpr := name.Type

	// This is used to remove pointer from type.
	if star, ok := tpExpr.(*ast.StarExpr); ok {
		tpExpr = star.X
	}

	switch convertedExpr := tpExpr.(type) {
	case *ast.IndexExpr: // This is used for generic structs or typedefs.
		tpExpr = convertedExpr.X
	case *ast.IndexListExpr: // This is used for multi-type generic structs or typedefs.
		tpExpr = convertedExpr.X
	}

	switch rawTp := tpExpr.(type) {
	case *ast.Ident: // This is used for usual structs or typedefs.
		return rawTp.Name
	default:
		panic("unexpected type")
	}
}
