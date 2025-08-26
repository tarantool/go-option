package extractor

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type packageImpl struct {
	pkg *packages.Package
}

// NewPackage creates a new Package from a packages.Package.
func NewPackage(pkg *packages.Package) Package {
	return &packageImpl{pkg: pkg}
}

func (p *packageImpl) Name() string {
	return p.pkg.Name
}

func (p *packageImpl) PkgPath() string {
	return p.pkg.PkgPath
}

func (p *packageImpl) Syntax() []*ast.File {
	return p.pkg.Syntax
}
