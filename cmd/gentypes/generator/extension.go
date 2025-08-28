// Package generator is a package that defines how code should be generated.
package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
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

// GenerateByType generates the code for the optional type.
func GenerateByType(typeName string, code int, packageName string) ([]byte, error) {
	var buf bytes.Buffer

	err := cTypeGenTemplate.Execute(&buf, struct {
		Name        string
		Type        string
		ExtCode     string
		PackageName string
		Imports     []string
	}{
		Name:        "Optional" + typeName,
		Type:        typeName,
		ExtCode:     strconv.Itoa(code),
		PackageName: packageName,
		Imports:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generateByType: %w", err)
	}

	return buf.Bytes(), nil
}
