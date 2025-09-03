// Package main is a binary, that generates optional types for types with support for MessagePack Extensions
// fast encoding/decoding.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/tarantool/go-option/cmd/gentypes/extractor"
	"github.com/tarantool/go-option/cmd/gentypes/generator"
)

const (
	defaultGoPermissions = 0644
)

var (
	packagePath         string
	extCode             int
	verbose             bool
	force               bool
	imports             stringListFlag
	customMarshalFunc   string
	customUnmarshalFunc string
)

func logfuncf(format string, args ...interface{}) {
	if verbose {
		fmt.Printf("> "+format+"\n", args...)
	}
}

func readGoFiles(ctx context.Context, folder string) ([]*packages.Package, error) {
	return packages.Load(&packages.Config{ //nolint:wrapcheck
		Mode:    packages.LoadAllSyntax,
		Context: ctx,
		Logf:    logfuncf,
		Dir:     folder,

		Env:        nil,
		BuildFlags: nil,
		Fset:       nil,
		ParseFile:  nil,
		Tests:      false,
		Overlay:    nil,
	})
}

func extractFirstPackageFromList(packageList []*packages.Package, name string) *packages.Package {
	if len(packageList) == 0 {
		panic("no packages found")
	}

	if name == "" {
		for _, pkg := range packageList {
			if !strings.HasSuffix(pkg.Name, "_test") {
				return pkg
			}
		}

		return packageList[0] // If no non-test packages found, return the first one.
	}

	for _, pkg := range packageList {
		if pkg.Name == name {
			return pkg
		}
	}

	fmt.Println("failed to find package with name:", name)
	fmt.Println("available packages:")

	for _, pkg := range packageList {
		fmt.Println("    ", pkg.Name)
	}

	os.Exit(1)

	return nil // Unreachable.
}

const (
	undefinedExtCode = math.MinInt8 - 1
)

func checkMsgpackExtCode(code int) bool {
	return code >= math.MinInt8 && code <= math.MaxInt8
}

func printFile(prefix string, data []byte) {
	for lineNo, line := range bytes.Split(data, []byte("\n")) {
		fmt.Printf("%03d%s%s\n", lineNo, prefix, string(line))
	}
}

func isExternalDep(name string) bool {
	return strings.Contains(name, ".")
}

const (
	maxNameParts = 2
)

func constructFileName(name string) string {
	parts := strings.SplitN(name, ".", maxNameParts)
	switch {
	case len(parts) == 1:
		name = parts[0]
	case len(parts) == maxNameParts:
		name = parts[1]
	}

	return strings.ToLower(name) + "_gen.go"
}

func main() { //nolint:funlen
	generator.InitializeTemplates()

	ctx := context.Background()

	flag.StringVar(&packagePath, "package", "./", "input and output path")
	flag.IntVar(&extCode, "ext-code", undefinedExtCode, "extension code")
	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.BoolVar(&force, "force", false, "generate files even if methods do not exist")
	flag.Var(&imports, "imports", "imports to add to generated files")
	flag.StringVar(&customMarshalFunc, "marshal-func", "", "custom marshal function")
	flag.StringVar(&customUnmarshalFunc, "unmarshal-func", "", "custom unmarshal function")
	flag.Parse()

	switch {
	case extCode == undefinedExtCode:
		fmt.Println("extension code is not set")

		flag.PrintDefaults()
		os.Exit(1)
	case !checkMsgpackExtCode(extCode):
		fmt.Println("invalid extension code:", extCode)
		fmt.Println("extension code must be in range [-128, 127]")

		flag.PrintDefaults()
		os.Exit(1)
	}

	packageList, err := readGoFiles(ctx, packagePath)
	switch {
	case err != nil:
		fmt.Println("failed to parse packages:")
		fmt.Println("    ", err)
		os.Exit(1)
	case packages.PrintErrors(packageList) > 0:
		os.Exit(1)
	case len(packageList) == 0:
		fmt.Println("no packages found")
		os.Exit(1)
	}

	pkg := extractFirstPackageFromList(packageList, "")

	analyzer, err := extractor.NewAnalyzerFromPackage(extractor.NewPackage(pkg))
	if err != nil {
		fmt.Println("failed to extract types and methods:")
		fmt.Println("    ", err)

		os.Exit(1)
	}

	args := flag.Args() // Args contains names of struct to generate optional types.
	switch {
	case len(args) == 0:
		fmt.Println("no struct name provided")

		flag.PrintDefaults()
		os.Exit(1)
	case len(args) > 1:
		fmt.Println("too many arguments")

		flag.PrintDefaults()
		os.Exit(1)
	}

	typeName := args[0]

	// Check for existence of all types that we want to generate.
	typeSpecDef, ok := analyzer.TypeSpecEntryByName(typeName)
	switch {
	case isExternalDep(typeName):
		fmt.Println("typename contains dot, probably third party type:", typeName)
	case !ok:
		fmt.Println("failed to find struct:", typeName)
		os.Exit(1)
	}

	fmt.Println("generating optional:", typeName)

	switch {
	case force || isExternalDep(typeName):
		// Skipping check for MarshalMsgpack and UnmarshalMsgpack methods.
	case !typeSpecDef.HasMethod("MarshalMsgpack") || !typeSpecDef.HasMethod("UnmarshalMsgpack"):
		fmt.Println("failed to find MarshalMsgpack or UnmarshalMsgpack method for struct:", typeName)
		os.Exit(1)
	}

	generatedGoSources, err := generator.GenerateByType(generator.GenerateOptions{
		TypeName:            typeName,
		ExtCode:             extCode,
		PackageName:         pkg.Name,
		Imports:             imports,
		CustomMarshalFunc:   customMarshalFunc,
		CustomUnmarshalFunc: customUnmarshalFunc,
	})
	if err != nil {
		fmt.Println("failed to generate optional types:")
		fmt.Println("    ", err)
		os.Exit(1)
	}

	formattedGoSource, err := format.Source(generatedGoSources)
	if err != nil {
		fmt.Println("failed to format generated code: ", err)
		printFile("> ", generatedGoSources)
		os.Exit(1)
	}

	err = os.WriteFile(
		filepath.Join(packagePath, constructFileName(typeName)),
		formattedGoSource,
		defaultGoPermissions,
	)
	if err != nil {
		fmt.Println("failed to write generated code:")
		fmt.Println("    ", err)
		os.Exit(1)
	}
}
