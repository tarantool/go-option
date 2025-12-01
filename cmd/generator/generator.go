package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	defaultGoPermissions = 0644
)

var (
	outputDirectory string
	verbose         bool
)

type generatorDef struct {
	Name        string
	Type        string
	DecodeFunc  string
	EncoderFunc string
	CheckerFunc string

	TestingValues                []string
	TestingValueOutputs          []string
	ExampleValueOutputs          []string
	UnexpectedTestingValue       string
	UnexpectedTestingValueOutput string
	ZeroTestingValueOutput       string
}

func structToMap(def generatorDef) map[string]any {
	caser := cases.Title(language.English)

	// Using first value for test.
	testingValue := def.TestingValues[0]
	testingValueOutput := def.TestingValueOutputs[0]

	out := map[string]any{
		"Name":        caser.String(def.Name),
		"Type":        def.Name,
		"DecodeFunc":  def.DecodeFunc,
		"EncoderFunc": def.EncoderFunc,
		"CheckerFunc": def.CheckerFunc,

		"TestingValue":                 testingValue,
		"TestingValueOutput":           testingValueOutput,
		"UnexpectedTestingValue":       def.UnexpectedTestingValue,
		"UnexpectedTestingValueOutput": def.UnexpectedTestingValueOutput,
		"ZeroTestingValueOutput":       def.ZeroTestingValueOutput,
		"ExampleValueOutput":           def.ExampleValueOutputs[0],

		// Adding arrays for EncodeDecodeMsgpack tests.
		"TestingValues":       def.TestingValues,
		"TestingValueOutputs": def.TestingValueOutputs,
	}

	if def.Type != "" {
		out["Type"] = def.Type
	}

	if def.UnexpectedTestingValueOutput != "" {
		out["UnexpectedTestingValueOutput"] = def.UnexpectedTestingValueOutput
	}

	return out
}

func zeroOutput[T any]() string {
	var zero T

	return fmt.Sprint(zero)
}

var defaultTypes = []generatorDef{
	{
		Name:        "byte",
		Type:        "byte",
		DecodeFunc:  "decodeByte",
		EncoderFunc: "encodeByte",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[byte](),
	},
	{
		Name:        "int",
		Type:        "int",
		DecodeFunc:  "decodeInt",
		EncoderFunc: "encodeInt",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[int](),
	},
	{
		Name:        "int8",
		Type:        "int8",
		DecodeFunc:  "decodeInt8",
		EncoderFunc: "encodeInt8",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[int8](),
	},
	{
		Name:        "int16",
		Type:        "int16",
		DecodeFunc:  "decodeInt16",
		EncoderFunc: "encodeInt16",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[int16](),
	},
	{
		Name:        "int32",
		Type:        "int32",
		DecodeFunc:  "decodeInt32",
		EncoderFunc: "encodeInt32",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[int32](),
	},
	{
		Name:        "int64",
		Type:        "int64",
		DecodeFunc:  "decodeInt64",
		EncoderFunc: "encodeInt64",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[int64](),
	},
	{
		Name:        "uint",
		Type:        "uint",
		DecodeFunc:  "decodeUint",
		EncoderFunc: "encodeUint",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[uint](),
	},
	{
		Name:        "uint8",
		Type:        "uint8",
		DecodeFunc:  "decodeUint8",
		EncoderFunc: "encodeUint8",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[uint8](),
	},
	{
		Name:        "uint16",
		Type:        "uint16",
		DecodeFunc:  "decodeUint16",
		EncoderFunc: "encodeUint16",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[uint16](),
	},
	{
		Name:        "uint32",
		Type:        "uint32",
		DecodeFunc:  "decodeUint32",
		EncoderFunc: "encodeUint32",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[uint32](),
	},
	{
		Name:        "uint64",
		Type:        "uint64",
		DecodeFunc:  "decodeUint64",
		EncoderFunc: "encodeUint64",
		CheckerFunc: "checkNumber",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[uint64](),
	},
	{
		Name:        "float32",
		Type:        "float32",
		DecodeFunc:  "decodeFloat32",
		EncoderFunc: "encodeFloat32",
		CheckerFunc: "checkFloat",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[float32](),
	},
	{
		Name:        "float64",
		Type:        "float64",
		DecodeFunc:  "decodeFloat64",
		EncoderFunc: "encodeFloat64",
		CheckerFunc: "checkFloat",

		TestingValues:                []string{"12"},
		TestingValueOutputs:          []string{"12"},
		ExampleValueOutputs:          []string{"12"},
		UnexpectedTestingValue:       "13",
		UnexpectedTestingValueOutput: "13",
		ZeroTestingValueOutput:       zeroOutput[float64](),
	},
	{
		Name:        "string",
		Type:        "string",
		DecodeFunc:  "decodeString",
		EncoderFunc: "encodeString",
		CheckerFunc: "checkString",

		TestingValues:                []string{"\"hello\""},
		TestingValueOutputs:          []string{"\"hello\""},
		ExampleValueOutputs:          []string{"hello"},
		UnexpectedTestingValue:       "\"bye\"",
		UnexpectedTestingValueOutput: "bye",
		ZeroTestingValueOutput:       zeroOutput[string](),
	},
	{
		Name:        "bytes",
		Type:        "[]byte",
		DecodeFunc:  "decodeBytes",
		EncoderFunc: "encodeBytes",
		CheckerFunc: "checkBytes",

		TestingValues:                []string{"[]byte{3, 14, 15}"},
		TestingValueOutputs:          []string{"[]byte{3, 14, 15}"},
		ExampleValueOutputs:          []string{"[3 14 15]"},
		UnexpectedTestingValue:       "[]byte{3, 14, 15, 9, 26}",
		UnexpectedTestingValueOutput: "[3 14 15 9 26]",
		ZeroTestingValueOutput:       zeroOutput[[]byte](),
	},
	{
		Name:        "bool",
		Type:        "bool",
		DecodeFunc:  "decodeBool",
		EncoderFunc: "encodeBool",
		CheckerFunc: "checkBool",

		TestingValues:                []string{"true"},
		TestingValueOutputs:          []string{"true"},
		ExampleValueOutputs:          []string{"true"},
		UnexpectedTestingValue:       "false",
		UnexpectedTestingValueOutput: "false",
		ZeroTestingValueOutput:       zeroOutput[bool](),
	},
	{
		Name:        "any",
		Type:        "any",
		DecodeFunc:  "decodeAny",
		EncoderFunc: "encodeAny",
		CheckerFunc: "checkAny",

		TestingValues:                []string{"\"hello\"", "123", "true", "123.456"},
		TestingValueOutputs:          []string{"\"hello\"", "123", "true", "123.456"},
		ExampleValueOutputs:          []string{"hello", "123", "true", "123.456"},
		UnexpectedTestingValue:       "\"bye\"",
		UnexpectedTestingValueOutput: "bye",
		ZeroTestingValueOutput:       zeroOutput[any](),
	},
}

var tplText = `
// Code generated by github.com/tarantool/go-option; DO NOT EDIT.

package {{ .packageName }}

import (
	{{ range $i, $import := .imports }}
	"{{ $import }}"
	{{ end }}

	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// {{.Name}} represents an optional value of type {{.Type}}.
// It can either hold a valid {{.Type}} (IsSome == true) or be empty (IsZero == true).
type {{.Name}} struct {
	value  {{.Type}}
	exists bool
}

var _ commonInterface[{{.Type}}] = (*{{.Name}})(nil)

// Some{{.Name}} creates an optional {{.Name}} with the given {{.Type}} value.
// The returned {{.Name}} will have IsSome() == true and IsZero() == false.
func Some{{.Name}}(value {{.Type}}) {{.Name}} {
	return {{.Name}}{
		value: value,
		exists: true,
	}
}

// None{{.Name}} creates an empty optional {{.Name}} value.
// The returned {{.Name}} will have IsSome() == false and IsZero() == true.
func None{{.Name}}() {{.Name}} {
	return {{.Name}}{
		exists: false,
		value:  zero[{{.Type}}](),
	}
}

// IsSome returns true if the {{.Name}} contains a value.
// This indicates the value is explicitly set (not None).
func (o {{.Name}}) IsSome() bool {
	return o.exists
}

// IsZero returns true if the {{.Name}} does not contain a value.
// Equivalent to !IsSome(). Useful for consistency with types where
// zero value (e.g. 0, false, zero struct) is valid and needs to be distinguished.
func (o {{.Name}}) IsZero() bool {
	return !o.exists
}

// IsNil is an alias for IsZero.
//
// This method is provided for compatibility with the msgpack Encoder interface.
func (o {{.Name}}) IsNil() bool {
	return o.IsZero()
}

// Get returns the stored value and a boolean flag indicating its presence.
// If the value is present, returns (value, true).
// If the value is absent, returns (zero value of {{.Type}}, false).
//
// Recommended usage:
//
//	if value, ok := o.Get(); ok {
//	    // use value
//	}
func (o {{.Name}}) Get() ({{.Type}}, bool) {
	return o.value, o.exists
}

// MustGet returns the stored value if it is present.
// Panics if the value is absent (i.e., IsZero() == true).
//
// Use with caution — only when you are certain the value exists.
//
// Panics with: "optional value is not set" if no value is set.
func (o {{.Name}}) MustGet() {{.Type}} {
	if !o.exists {
		panic("optional value is not set")
	}

	return o.value
}

// Unwrap returns the stored value regardless of presence.
// If no value is set, returns the zero value for {{.Type}}.
//
// Warning: Does not check presence. Use IsSome() before calling if you need
// to distinguish between absent value and explicit zero value.
func (o {{.Name}}) Unwrap() {{.Type}} {
	return o.value
}

// UnwrapOr returns the stored value if present.
// Otherwise, returns the provided default value.
func (o {{.Name}}) UnwrapOr(defaultValue {{.Type}}) {{.Type}} {
	if o.exists {
		return o.value
	}

	return defaultValue
}

// UnwrapOrElse returns the stored value if present.
// Otherwise, calls the provided function and returns its result.
// Useful when the default value requires computation or side effects.
func (o {{.Name}}) UnwrapOrElse(defaultValue func() {{.Type}}) {{.Type}} {
	if o.exists {
		return o.value
	}

	return defaultValue()
}

// EncodeMsgpack encodes the {{.Name}} value using MessagePack format.
// - If the value is present, it is encoded as {{.Type}}.
// - If the value is absent (None), it is encoded as nil.
//
// Returns an error if encoding fails.
func (o {{.Name}}) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if o.exists {
		return newEncodeError("{{.Name}}", {{ .EncoderFunc }}(encoder, o.value))
	}

	return newEncodeError("{{.Name}}", encoder.EncodeNil())
}

// DecodeMsgpack decodes a {{.Name}} value from MessagePack format.
// Supports two input types:
//   - nil: interpreted as no value (None{{.Name}})
//   - {{.Type}}: interpreted as a present value (Some{{.Name}})
//
// Returns an error if the input type is unsupported or decoding fails.
//
// After successful decoding:
//   - on nil: exists = false, value = default zero value
//   - on {{.Type}}: exists = true, value = decoded value
func (o *{{.Name}}) DecodeMsgpack(decoder *msgpack.Decoder) error {
	code, err := decoder.PeekCode()
	if err != nil {
		return newDecodeError("{{.Name}}", err)
	}

	switch {
	case code == msgpcode.Nil:
		o.exists = false

		return newDecodeError("{{.Name}}", decoder.Skip())
	case {{ .CheckerFunc }}(code):
		o.value, err = {{ .DecodeFunc }}(decoder)
		if err != nil {
			return newDecodeError("{{.Name}}", err)
		}
		o.exists = true

		return err
	default:
		return newDecodeWithCodeError("{{.Name}}", code)
	}
}`

var tplTestText = `
// Code generated by github.com/tarantool/go-option; DO NOT EDIT.

package {{ .packageName }}_test

import (
	{{ range $i, $import := .imports }}
	"{{ $import }}"
	{{ end }}

	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/tarantool/go-option"
)

func Test{{.Name}}_IsSome(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.True(t, some{{.Name}}.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.False(t, empty{{.Name}}.IsSome())
	})
}

func Test{{.Name}}_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.False(t, some{{.Name}}.IsZero())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.True(t, empty{{.Name}}.IsZero())
	})
}

func Test{{.Name}}_IsNil(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.False(t, some{{.Name}}.IsNil())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.True(t, empty{{.Name}}.IsNil())
	})
}

func Test{{.Name}}_Get(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		val, ok := some{{.Name}}.Get()
		require.True(t, ok)
		assert.EqualValues(t, {{.TestingValue}}, val)
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		_, ok := empty{{.Name}}.Get()
		require.False(t, ok)
	})
}

func Test{{.Name}}_MustGet(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.EqualValues(t, {{.TestingValue}}, some{{.Name}}.MustGet())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.Panics(t, func() {
			empty{{.Name}}.MustGet()
		})
	})
}

func Test{{.Name}}_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.EqualValues(t, {{.TestingValue}}, some{{.Name}}.Unwrap())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.NotPanics(t, func() {
			empty{{.Name}}.Unwrap()
		})
	})
}

func Test{{.Name}}_UnwrapOr(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.EqualValues(t, {{.TestingValue}}, some{{.Name}}.UnwrapOr({{.UnexpectedTestingValue}}))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.EqualValues(t, {{.UnexpectedTestingValue}}, empty{{.Name}}.UnwrapOr({{.UnexpectedTestingValue}}))
	})
}

func Test{{.Name}}_UnwrapOrElse(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		some{{.Name}} := option.Some{{.Name}}({{.TestingValue}})
		assert.EqualValues(t, {{.TestingValue}}, some{{.Name}}.UnwrapOrElse(func() {{.Type}} {
			return {{.UnexpectedTestingValue}}
		}))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		empty{{.Name}} := option.None{{.Name}}()
		assert.EqualValues(t, {{.UnexpectedTestingValue}}, empty{{.Name}}.UnwrapOrElse(func() {{.Type}} {
			return {{.UnexpectedTestingValue}}
		}))
	})
}

func Test{{.Name}}_EncodeDecodeMsgpack(t *testing.T) {
	t.Parallel()

	{{ range $i, $value := .TestingValues }}

	{{ if (eq $i 0) }}
	t.Run("some", func(t *testing.T) {
	{{ else }}
	t.Run("some_{{ $i }}", func(t *testing.T) {
	{{ end -}}

		t.Parallel()

		var buf bytes.Buffer

		enc := msgpack.NewEncoder(&buf)
		dec := msgpack.NewDecoder(&buf)

		some{{$.Name}} := option.Some{{$.Name}}({{ $value }})
		err := some{{$.Name}}.EncodeMsgpack(enc)
		require.NoError(t, err)

		var unmarshaled option.{{$.Name}}
		err = unmarshaled.DecodeMsgpack(dec)
		require.NoError(t, err)
		assert.True(t, unmarshaled.IsSome())
		{{- $output := index $.TestingValueOutputs $i }}
		assert.EqualValues(t, {{ $output }}, unmarshaled.Unwrap())
	})
	{{ end }}

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		enc := msgpack.NewEncoder(&buf)
		dec := msgpack.NewDecoder(&buf)

		empty{{.Name}} := option.None{{.Name}}()
		err := empty{{.Name}}.EncodeMsgpack(enc)
		require.NoError(t, err)

		var unmarshaled option.{{.Name}}
		err = unmarshaled.DecodeMsgpack(dec)

		require.NoError(t, err)
		assert.False(t, unmarshaled.IsSome())
	})
}

func ExampleSome{{.Name}}() {
	opt := option.Some{{.Name}}({{.TestingValue}})
	if opt.IsSome() {
		fmt.Println(opt.Unwrap())
	}
	// Output: {{.ExampleValueOutput}}
}

func ExampleNone{{.Name}}() {
	opt := option.None{{.Name}}()
	if opt.IsZero() {
		fmt.Println("value is absent")
	}
	// Output: value is absent
}

func Example{{.Name}}_IsSome() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.IsSome())
	fmt.Println(none.IsSome())
	// Output:
	// true
	// false
}

func Example{{.Name}}_IsZero() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.IsZero())
	fmt.Println(none.IsZero())
	// Output:
	// false
	// true
}

func Example{{.Name}}_IsNil() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.IsNil() == some.IsZero())
	fmt.Println(none.IsNil() == none.IsZero())
	// Output:
	// true
	// true
}

func Example{{.Name}}_Get() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	val, ok := some.Get()
	fmt.Println(val, ok)
	val, ok = none.Get()
	fmt.Println(val, ok)
	// Output:
	// {{.ExampleValueOutput}} true
	// {{.ZeroTestingValueOutput}} false
}

func Example{{.Name}}_MustGet() {
	some := option.Some{{.Name}}({{.TestingValue}})
	fmt.Println(some.MustGet())
	// Output: {{.ExampleValueOutput}}
}

func Example{{.Name}}_MustGet_panic() {
	none := option.None{{.Name}}()
	eof := false
	defer func() {
		if !eof {
			fmt.Println("panic!", recover())
		}
	}()
	fmt.Println(none.MustGet())
	eof = true
	// Output: panic! optional value is not set
}

func Example{{.Name}}_Unwrap() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.Unwrap())
	fmt.Println(none.Unwrap())
	// Output:
	// {{.ExampleValueOutput}}
	// {{.ZeroTestingValueOutput}}
}

func Example{{.Name}}_UnwrapOr() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.UnwrapOr({{.UnexpectedTestingValue}}))
	fmt.Println(none.UnwrapOr({{.UnexpectedTestingValue}}))
	// Output:
	// {{.ExampleValueOutput}}
	// {{.UnexpectedTestingValueOutput}}
}

func Example{{.Name}}_UnwrapOrElse() {
	some := option.Some{{.Name}}({{.TestingValue}})
	none := option.None{{.Name}}()
	fmt.Println(some.UnwrapOrElse(func() {{.Type}} {
		return {{.UnexpectedTestingValue}}
	}))
	fmt.Println(none.UnwrapOrElse(func() {{.Type}} {
		return {{.UnexpectedTestingValue}}
	}))
	// Output:
	// {{.ExampleValueOutput}}
	// {{.UnexpectedTestingValueOutput}}
}
`

func printFile(prefix string, data []byte) {
	for lineNo, line := range bytes.Split(data, []byte("\n")) {
		fmt.Printf("%03d%s%s\n", lineNo, prefix, string(line))
	}
}

func generateAndWrite() error {
	tmpl, err := template.New("internal").Parse(tplText)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	tmpl_test, err := template.New("internal_test").Parse(tplTestText)
	if err != nil {
		return fmt.Errorf("failed to parse testing template: %w", err)
	}

	outputData := make(map[string][]byte, 2*len(defaultTypes)) //nolint:mnd

	var data bytes.Buffer

	// 1. Generate code for each type.
	for _, generatedType := range defaultTypes {
		tmplData := structToMap(generatedType)

		tmplData["packageName"] = "option" // Package name is option, since generator is used for option only right now.
		tmplData["imports"] = []string{}   // No additional imports are needed right now.

		// Generate code of an Optional type.
		{
			err := tmpl.Execute(&data, tmplData)
			if err != nil {
				return fmt.Errorf("failed to execute template: %w", err)
			}

			outputData[generatedType.Name+"_gen.go"] = slices.Clone(data.Bytes())
			data.Reset()
		}

		// Generate code for tests of an Optional type.
		{
			err := tmpl_test.Execute(&data, tmplData)
			if err != nil {
				return fmt.Errorf("failed to execute test template: %w", err)
			}

			outputData[generatedType.Name+"_gen_test.go"] = slices.Clone(data.Bytes())
			data.Reset()
		}
	}

	// 2. Just in case format code using gofmt.
	for name, origData := range outputData {
		data, err := format.Source(origData)
		if err != nil {
			if verbose {
				printFile("> ", origData)
			}

			return fmt.Errorf("failed to format code: %w", err)
		}

		outputData[name] = data
	}

	// 3. Write resulting code to files.
	for name, data := range outputData {
		err = os.WriteFile(filepath.Join(outputDirectory, name), data, defaultGoPermissions)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

func main() {
	flag.StringVar(&outputDirectory, "output", ".", "output directory")
	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.Parse()

	// Get absolute path for output directory.
	absOutputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		fmt.Printf("failed to get absolute path for output directory (%s): %s\n", outputDirectory, err)
		os.Exit(1)
	}

	// Check if output directory exists and is directory.
	switch fInfo, err := os.Stat(absOutputDirectory); {
	case err != nil:
		fmt.Println("failed to stat output directory: ", err.Error())
		os.Exit(1)
	case !fInfo.IsDir():
		fmt.Printf("output directory '%s' is not a directory\n", absOutputDirectory)
		os.Exit(1)
	}

	err = generateAndWrite()
	if err != nil {
		fmt.Println("failed to generate or write code: ", err)
		os.Exit(1)
	}
}
