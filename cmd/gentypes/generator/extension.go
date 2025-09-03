// Package generator is a package that defines how code should be generated.
package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

//go:embed type_gen.go.tpl
var typeGenTemplate string

//go:embed type_gen_test.go.tpl
var typeGenTestTemplate string

var (
	cTypeGenTemplate     *template.Template
	cTypeGenTestTemplate *template.Template //nolint:unused
)

// InitializeTemplates initializes the templates, should be called at the start of the main program loop.
func InitializeTemplates() {
	cTypeGenTemplate = template.Must(template.New("type_gen.go.tpl").Parse(typeGenTemplate))
	cTypeGenTestTemplate = template.Must(template.New("type_gen_test.go.tpl").Parse(typeGenTestTemplate))
}

const (
	maxNameParts = 2
)

func constructTypeName(typeName string) string {
	splittedName := strings.SplitN(typeName, ".", maxNameParts)
	switch len(splittedName) {
	case 1:
		typeName = splittedName[0]
	case maxNameParts:
		typeName = splittedName[1]
	default:
		panic("invalid type name: " + typeName)
	}

	return "Optional" + typeName
}

// GenerateOptions is the options for the code generation.
type GenerateOptions struct {
	// TypeName is the name of the type to generate optional to.
	TypeName string
	// ExtCode is the extension code.
	ExtCode int
	// PackageName is the name of the package to generate to.
	PackageName string
	// Imports is the list of imports to add to the generated code.
	Imports []string
	// CustomMarshalFunc is the name of the custom marshal function.
	CustomMarshalFunc string
	// CustomUnmarshalFunc is the name of the custom unmarshal function.
	CustomUnmarshalFunc string
}

// GenerateByType generates the code for the optional type.
func GenerateByType(opts GenerateOptions) ([]byte, error) {
	var buf bytes.Buffer

	if opts.CustomMarshalFunc == "" {
		opts.CustomMarshalFunc = "o.value.MarshalMsgpack()"
	} else {
		opts.CustomMarshalFunc += "(o.value)"
	}

	if opts.CustomUnmarshalFunc == "" {
		opts.CustomUnmarshalFunc = "o.value.UnmarshalMsgpack(a)"
	} else {
		opts.CustomUnmarshalFunc += "(&o.value, a)"
	}

	err := cTypeGenTemplate.Execute(&buf, struct {
		Name                string
		Type                string
		ExtCode             string
		PackageName         string
		Imports             []string
		CustomMarshalFunc   string
		CustomUnmarshalFunc string
	}{
		Name:                constructTypeName(opts.TypeName),
		Type:                opts.TypeName,
		ExtCode:             strconv.Itoa(opts.ExtCode),
		PackageName:         opts.PackageName,
		Imports:             opts.Imports,
		CustomMarshalFunc:   opts.CustomMarshalFunc,
		CustomUnmarshalFunc: opts.CustomUnmarshalFunc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generateByType: %w", err)
	}

	return buf.Bytes(), nil
}
