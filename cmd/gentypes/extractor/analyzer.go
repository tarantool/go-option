// Package extractor is a package, that extracts type specs and methods from given ast tree.
package extractor

import (
	"go/ast"
)

// TypeSpecEntry is an entry, that defines ast's TypeSpec and contains type name and methods.
type TypeSpecEntry struct {
	Name    string
	Methods []string

	methodMap map[string]struct{}

	rawType    *ast.TypeSpec
	rawMethods []*ast.FuncDecl
}

// HasMethod returns true if type spec has method with given name.
func (e TypeSpecEntry) HasMethod(name string) bool {
	_, ok := e.methodMap[name]
	return ok
}

// Analyzer is an analyzer, that extracts type specs and methods from package and groups
// them for quick access.
type Analyzer struct {
	pkgPath string
	pkgName string
	entries map[string]*TypeSpecEntry
}

// NewAnalyzerFromPackage parses ast tree for TypeSpecs and associated methods.
func NewAnalyzerFromPackage(pkg Package) (*Analyzer, error) {
	typeSpecs := ExtractTypeSpecsFromPackage(pkg)
	methodsDefs := ExtractMethodsFromPackage(pkg)

	analyzer := &Analyzer{
		entries: make(map[string]*TypeSpecEntry, len(typeSpecs)),
		pkgPath: pkg.PkgPath(),
		pkgName: pkg.Name(),
	}

	for _, typeSpec := range typeSpecs {
		tsName := typeSpec.Name.String()
		if _, ok := analyzer.entries[tsName]; ok {
			// Duplicate type spec, skipping.
			continue
		}

		entry := &TypeSpecEntry{
			Name:       tsName,
			Methods:    nil,
			methodMap:  make(map[string]struct{}),
			rawType:    typeSpec,
			rawMethods: nil,
		}

		for _, methodDef := range methodsDefs {
			typeName := ExtractRecvTypeName(methodDef)
			if typeName != tsName {
				continue
			}

			entry.Methods = append(entry.Methods, methodDef.Name.String())
			entry.rawMethods = append(entry.rawMethods, methodDef)
			entry.methodMap[methodDef.Name.String()] = struct{}{}
		}

		analyzer.entries[tsName] = entry
	}

	return analyzer, nil
}

// PackagePath returns package path of analyzed package.
func (a Analyzer) PackagePath() string {
	return a.pkgPath
}

// PackageName returns package name of analyzed package.
func (a Analyzer) PackageName() string {
	return a.pkgName
}

// TypeSpecEntryByName returns TypeSpecEntry entry by name.
func (a Analyzer) TypeSpecEntryByName(name string) (*TypeSpecEntry, bool) {
	structEntry, ok := a.entries[name]
	return structEntry, ok
}
