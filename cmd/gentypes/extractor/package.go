package extractor

import "go/ast"

// Package is an interface that provides access to package data.
// It's used to abstract away the `packages.Package` type.
type Package interface {
	Name() string
	PkgPath() string
	Syntax() []*ast.File
}
