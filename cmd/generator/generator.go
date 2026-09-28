package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

	// JSONTestingValueOutput is the JSON encoding of TestingValues[0].
	JSONTestingValueOutput string
	// JSONZeroValue is an empty value of Type whose JSON and YAML encodings
	// are not null.
	JSONZeroValue string
	// JSONZeroValueOutput is the JSON encoding of JSONZeroValue.
	JSONZeroValueOutput string
	// JSONNullable is true when nil is a valid value of Type, so a present
	// value may still encode as JSON null.
	JSONNullable bool

	// YAMLTestingValueOutput is the YAML encoding of TestingValues[0].
	YAMLTestingValueOutput string
	// YAMLZeroValueOutput is the YAML encoding of JSONZeroValue.
	YAMLZeroValueOutput string
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

		"JSONTestingValueOutput": def.JSONTestingValueOutput,
		"JSONZeroValue":          def.JSONZeroValue,
		"JSONZeroValueOutput":    def.JSONZeroValueOutput,
		"JSONNullable":           def.JSONNullable,

		"YAMLTestingValueOutput": def.YAMLTestingValueOutput,
		"YAMLZeroValueOutput":    def.YAMLZeroValueOutput,

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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: "12",
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           false,

		YAMLTestingValueOutput: "12\n",
		YAMLZeroValueOutput:    "0\n",
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

		JSONTestingValueOutput: `"hello"`,
		JSONZeroValue:          `""`,
		JSONZeroValueOutput:    `""`,
		JSONNullable:           false,

		YAMLTestingValueOutput: "hello\n",
		YAMLZeroValueOutput:    "\"\"\n",
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

		JSONTestingValueOutput: `"Aw4P"`, // Base64 of {3, 14, 15}.
		JSONZeroValue:          "[]byte{}",
		JSONZeroValueOutput:    `""`,
		JSONNullable:           true,

		// YAML encodes a byte slice as a sequence of numbers, and a nil one as [].
		YAMLTestingValueOutput: "- 3\n- 14\n- 15\n",
		YAMLZeroValueOutput:    "[]\n",
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

		JSONTestingValueOutput: "true",
		JSONZeroValue:          "false",
		JSONZeroValueOutput:    "false",
		JSONNullable:           false,

		YAMLTestingValueOutput: "true\n",
		YAMLZeroValueOutput:    "false\n",
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

		// The zero value of any is nil, which encodes as null (see JSONNullable),
		// so zero of a concrete type stands in for "empty but present".
		JSONTestingValueOutput: `"hello"`,
		JSONZeroValue:          "0",
		JSONZeroValueOutput:    "0",
		JSONNullable:           true,

		YAMLTestingValueOutput: "hello\n",
		YAMLZeroValueOutput:    "0\n",
	},
}

var tplText = `
// Code generated by github.com/tarantool/go-option; DO NOT EDIT.

package {{ .packageName }}

import (
	{{ range $i, $import := .imports }}
	"{{ $import }}"
	{{ end }}

	"bytes"
	"encoding/json"

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
}

// MarshalJSON implements the json.Marshaler interface.
//   - If the value is present, it is encoded exactly as json.Marshal encodes
//     the {{.Type}} value itself.
//   - If the value is absent (None), it is encoded as JSON null.
{{- if .JSONNullable }}
//
// A present nil {{.Type}} is encoded as JSON null too, so it decodes back
// as an absent value.
{{- end }}
//
// A struct field of type {{.Name}} with the "omitzero" JSON tag option is
// omitted when the value is absent, since IsZero reports absence.
// The "omitempty" option has no effect on it.
func (o {{.Name}}) MarshalJSON() ([]byte, error) {
	if !o.exists {
		return []byte("null"), nil
	}

	return json.Marshal(&o.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//   - JSON null is decoded as no value (None{{.Name}}).
//   - Any other JSON value is decoded as {{.Type}} and stored as a present value.
//
// An error of decoding {{.Type}} is returned unchanged, and the receiver
// is left intact.
func (o *{{.Name}}) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*o = None{{.Name}}()

		return nil
	}

	var value {{.Type}}

	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	*o = Some{{.Name}}(value)

	return nil
}

// MarshalYAML implements the Marshaler interface of gopkg.in/yaml.v3
// (and yaml.v2) without depending on it.
//   - If the value is present, it returns a pointer to it, so the value is
//     encoded exactly as the YAML encoder encodes {{.Type}} itself.
//   - If the value is absent (None), it returns nil, which is encoded as null.
//
// A struct field of type {{.Name}} with the "omitempty" YAML tag option is
// omitted when the value is absent, since the encoder consults IsZero.
func (o {{.Name}}) MarshalYAML() (any, error) {
	if !o.exists {
		return nil, nil
	}

	return &o.value, nil
}

// UnmarshalYAML implements the obsolete Unmarshaler interface of
// gopkg.in/yaml.v3 (the only one of yaml.v2) without depending on it.
//   - A null value is decoded as no value (None{{.Name}}).
//   - Any other value is decoded as {{.Type}} and stored as a present value.
//
// gopkg.in/yaml.v3 never calls this method for a null node (null, ~ or an
// empty value): it leaves the receiver unchanged. Decoding null into a zero
// {{.Name}} therefore yields None, but decoding it into a present value
// keeps that value.
//
// An error of decoding {{.Type}} is returned unchanged, and the receiver
// is left intact.
func (o *{{.Name}}) UnmarshalYAML(unmarshal func(any) error) error {
	var value *{{.Type}}

	err := unmarshal(&value)
	if err != nil {
		return err
	}

	if value == nil {
		*o = None{{.Name}}()

		return nil
	}

	*o = Some{{.Name}}(*value)

	return nil
}`

var tplTestText = `
// Code generated by github.com/tarantool/go-option; DO NOT EDIT.

package {{ .packageName }}_test

import (
	{{ range $i, $import := .imports }}
	"{{ $import }}"
	{{ end }}

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"

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

func Test{{.Name}}_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.Some{{.Name}}({{.TestingValue}}))
		require.NoError(t, err)
		assert.Equal(t, {{ raw .JSONTestingValueOutput }}, string(data))
	})

	{{ range $i, $value := .TestingValues }}
	t.Run("some_as_inner_{{ $i }}", func(t *testing.T) {
		t.Parallel()

		expected, err := json.Marshal({{$.Type}}({{ $value }}))
		require.NoError(t, err)

		data, err := json.Marshal(option.Some{{$.Name}}({{ $value }}))
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
	})
	{{ end }}

	t.Run("some_zero", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.Some{{.Name}}({{.JSONZeroValue}}))
		require.NoError(t, err)
		assert.Equal(t, {{ raw .JSONZeroValueOutput }}, string(data))
	})

	{{- if .JSONNullable }}

	t.Run("some_nil", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.Some{{.Name}}(nil))
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})
	{{- end }}

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.None{{.Name}}())
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("struct_field", func(t *testing.T) {
		t.Parallel()

		type wrapper struct {
			Some option.{{.Name}} {{ raw "json:\"some\"" }}
			None option.{{.Name}} {{ raw "json:\"none\"" }}
		}

		data, err := json.Marshal(wrapper{
			Some: option.Some{{.Name}}({{.TestingValue}}),
			None: option.None{{.Name}}(),
		})
		require.NoError(t, err)
		assert.Equal(t, {{ raw (printf "{\"some\":%s,\"none\":null}" .JSONTestingValueOutput) }}, string(data))
	})

	t.Run("omitzero", func(t *testing.T) {
		t.Parallel()

		type wrapper struct {
			Value option.{{.Name}} {{ raw "json:\"value,omitzero\"" }}
		}

		data, err := json.Marshal(wrapper{Value: option.None{{.Name}}()})
		require.NoError(t, err)
		assert.Equal(t, "{}", string(data))

		data, err = json.Marshal(wrapper{Value: option.Some{{.Name}}({{.JSONZeroValue}})})
		require.NoError(t, err)
		assert.Equal(t, {{ raw (printf "{\"value\":%s}" .JSONZeroValueOutput) }}, string(data))
	})
}

func Test{{.Name}}_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		var opt option.{{.Name}}
		require.NoError(t, json.Unmarshal([]byte({{ raw .JSONTestingValueOutput }}), &opt))
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, opt.Unwrap())
	})

	{{ range $i, $value := .TestingValues }}
	t.Run("roundtrip_{{ $i }}", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.Some{{$.Name}}({{ $value }}))
		require.NoError(t, err)

		var opt option.{{$.Name}}
		require.NoError(t, json.Unmarshal(data, &opt))
		require.True(t, opt.IsSome())
		{{- $output := index $.TestingValueOutputs $i }}
		assert.EqualValues(t, {{ $output }}, opt.Unwrap())
	})
	{{ end }}

	t.Run("some_zero", func(t *testing.T) {
		t.Parallel()

		var opt option.{{.Name}}
		require.NoError(t, json.Unmarshal([]byte({{ raw .JSONZeroValueOutput }}), &opt))
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{.JSONZeroValue}}, opt.Unwrap())
	})

	{{- if .JSONNullable }}

	t.Run("some_nil", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(option.Some{{.Name}}(nil))
		require.NoError(t, err)

		opt := option.Some{{.Name}}({{.TestingValue}})
		require.NoError(t, json.Unmarshal(data, &opt))
		assert.False(t, opt.IsSome())
	})
	{{- end }}

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := option.Some{{.Name}}({{.TestingValue}})
		require.NoError(t, json.Unmarshal([]byte("null"), &opt))
		assert.False(t, opt.IsSome())
		assert.Zero(t, opt.Unwrap())
	})

	t.Run("none_with_spaces", func(t *testing.T) {
		t.Parallel()

		opt := option.Some{{.Name}}({{.TestingValue}})
		require.NoError(t, opt.UnmarshalJSON([]byte(" null\n")))
		assert.False(t, opt.IsSome())
	})

	t.Run("struct_field", func(t *testing.T) {
		t.Parallel()

		type wrapper struct {
			Some   option.{{.Name}} {{ raw "json:\"some\"" }}
			None   option.{{.Name}} {{ raw "json:\"none\"" }}
			Absent option.{{.Name}} {{ raw "json:\"absent\"" }}
		}

		var out wrapper
		err := json.Unmarshal([]byte({{ raw (printf "{\"some\":%s,\"none\":null}" .JSONTestingValueOutput) }}), &out)
		require.NoError(t, err)
		require.True(t, out.Some.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, out.Some.Unwrap())
		assert.False(t, out.None.IsSome())
		assert.False(t, out.Absent.IsSome())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		opt := option.Some{{.Name}}({{.TestingValue}})
		err := opt.UnmarshalJSON([]byte("{"))

		var syntaxErr *json.SyntaxError
		require.ErrorAs(t, err, &syntaxErr)
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, opt.Unwrap())
	})
}

func Test{{.Name}}_MarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		data, err := yaml.Marshal(option.Some{{.Name}}({{.TestingValue}}))
		require.NoError(t, err)
		assert.Equal(t, {{ printf "%q" .YAMLTestingValueOutput }}, string(data))
	})

	{{ range $i, $value := .TestingValues }}
	t.Run("some_as_inner_{{ $i }}", func(t *testing.T) {
		t.Parallel()

		expected, err := yaml.Marshal({{$.Type}}({{ $value }}))
		require.NoError(t, err)

		data, err := yaml.Marshal(option.Some{{$.Name}}({{ $value }}))
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(data))
	})
	{{ end }}

	t.Run("some_zero", func(t *testing.T) {
		t.Parallel()

		data, err := yaml.Marshal(option.Some{{.Name}}({{.JSONZeroValue}}))
		require.NoError(t, err)
		assert.Equal(t, {{ printf "%q" .YAMLZeroValueOutput }}, string(data))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		data, err := yaml.Marshal(option.None{{.Name}}())
		require.NoError(t, err)
		assert.Equal(t, "null\n", string(data))
	})

	t.Run("struct_field", func(t *testing.T) {
		t.Parallel()

		type wrapper struct {
			Some    option.{{.Name}} {{ raw "yaml:\"some\"" }}
			None    option.{{.Name}} {{ raw "yaml:\"none\"" }}
			Omitted option.{{.Name}} {{ raw "yaml:\"omitted,omitempty\"" }}
			Zero    option.{{.Name}} {{ raw "yaml:\"zero,omitempty\"" }}
		}

		data, err := yaml.Marshal(wrapper{
			Some:    option.Some{{.Name}}({{.TestingValue}}),
			None:    option.None{{.Name}}(),
			Omitted: option.None{{.Name}}(),
			Zero:    option.Some{{.Name}}({{.JSONZeroValue}}),
		})
		require.NoError(t, err)

		var fields map[string]any
		require.NoError(t, yaml.Unmarshal(data, &fields))
		assert.Contains(t, fields, "some")
		assert.NotNil(t, fields["some"])
		assert.Contains(t, fields, "none")
		assert.Nil(t, fields["none"])
		assert.NotContains(t, fields, "omitted")
		assert.Contains(t, fields, "zero")
		assert.NotNil(t, fields["zero"])
	})
}

func Test{{.Name}}_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	{{ range $i, $value := .TestingValues }}
	t.Run("roundtrip_{{ $i }}", func(t *testing.T) {
		t.Parallel()

		data, err := yaml.Marshal(option.Some{{$.Name}}({{ $value }}))
		require.NoError(t, err)

		var opt option.{{$.Name}}
		require.NoError(t, yaml.Unmarshal(data, &opt))
		require.True(t, opt.IsSome())
		{{- $output := index $.TestingValueOutputs $i }}
		assert.EqualValues(t, {{ $output }}, opt.Unwrap())
	})
	{{ end }}

	t.Run("some_zero", func(t *testing.T) {
		t.Parallel()

		var opt option.{{.Name}}
		require.NoError(t, yaml.Unmarshal([]byte({{ printf "%q" .YAMLZeroValueOutput }}), &opt))
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{.JSONZeroValue}}, opt.Unwrap())
	})

	t.Run("none_roundtrip", func(t *testing.T) {
		t.Parallel()

		data, err := yaml.Marshal(option.None{{.Name}}())
		require.NoError(t, err)

		var opt option.{{.Name}}
		require.NoError(t, yaml.Unmarshal(data, &opt))
		assert.False(t, opt.IsSome())
	})

	t.Run("null_node_keeps_present_value", func(t *testing.T) {
		t.Parallel()

		// gopkg.in/yaml.v3 does not call UnmarshalYAML for a null node.
		opt := option.Some{{.Name}}({{.TestingValue}})
		require.NoError(t, yaml.Unmarshal([]byte("null"), &opt))
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, opt.Unwrap())
	})

	t.Run("null_by_direct_call", func(t *testing.T) {
		t.Parallel()

		opt := option.Some{{.Name}}({{.TestingValue}})
		err := opt.UnmarshalYAML(func(v any) error {
			return yaml.Unmarshal([]byte("null"), v)
		})
		require.NoError(t, err)
		assert.False(t, opt.IsSome())
		assert.Zero(t, opt.Unwrap())
	})

	t.Run("struct_field", func(t *testing.T) {
		t.Parallel()

		type wrapper struct {
			Some   option.{{.Name}} {{ raw "yaml:\"some\"" }}
			Null   option.{{.Name}} {{ raw "yaml:\"nulled\"" }}
			Tilde  option.{{.Name}} {{ raw "yaml:\"tilde\"" }}
			Empty  option.{{.Name}} {{ raw "yaml:\"empty\"" }}
			Absent option.{{.Name}} {{ raw "yaml:\"absent\"" }}
		}

		some, err := yaml.Marshal(map[string]any{"some": option.Some{{.Name}}({{.TestingValue}})})
		require.NoError(t, err)

		var out wrapper
		require.NoError(t, yaml.Unmarshal(append(some, []byte("nulled: null\ntilde: ~\nempty:\n")...), &out))
		require.True(t, out.Some.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, out.Some.Unwrap())
		assert.False(t, out.Null.IsSome())
		assert.False(t, out.Tilde.IsSome())
		assert.False(t, out.Empty.IsSome())
		assert.False(t, out.Absent.IsSome())
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		errUnmarshal := errors.New("unmarshal failed")

		opt := option.Some{{.Name}}({{.TestingValue}})
		err := opt.UnmarshalYAML(func(any) error {
			return errUnmarshal
		})
		require.ErrorIs(t, err, errUnmarshal)
		require.True(t, opt.IsSome())
		assert.EqualValues(t, {{ index .TestingValueOutputs 0 }}, opt.Unwrap())
	})
}

func Example{{.Name}}_MarshalJSON() {
	for _, opt := range []option.{{.Name}}{option.Some{{.Name}}({{.TestingValue}}), option.None{{.Name}}()} {
		data, err := json.Marshal(opt)
		if err != nil {
			fmt.Println("error:", err)

			return
		}

		fmt.Println(string(data))
	}
	// Output:
	// {{.JSONTestingValueOutput}}
	// null
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

var errBacktickInRawString = errors.New("raw string literal cannot contain a backtick")

// rawString renders s as a Go raw string literal. The templates are Go raw
// strings themselves, so they cannot spell a backtick (struct tags, JSON
// literals) directly.
func rawString(s string) (string, error) {
	if strings.Contains(s, "`") {
		return "", fmt.Errorf("%w: %q", errBacktickInRawString, s)
	}

	return "`" + s + "`", nil
}

var templateFuncs = template.FuncMap{
	"raw": rawString,
}

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

	tmpl_test, err := template.New("internal_test").Funcs(templateFuncs).Parse(tplTestText)
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
