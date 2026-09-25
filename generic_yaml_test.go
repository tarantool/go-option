package option_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/tarantool/go-option"
)

var errUnmarshalFailed = errors.New("unmarshal failed")

func TestGeneric_MarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"some", option.Some(42), "42\n"},
		{"some_zero", option.Some(0), "0\n"},
		{"some_empty_string", option.Some(""), "\"\"\n"},
		{"some_struct", option.Some(CustomType{Value: "x"}), "value: x\n"},
		{"some_pointer_receiver", option.Some(pointerTextType{Value: "x"}), "text:x\n"},
		{"some_nil_pointer", option.Some[*int](nil), "null\n"},
		{"none", option.None[int](), "null\n"},
		{"none_struct", option.None[CustomType](), "null\n"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			data, err := yaml.Marshal(testCase.value)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, string(data))
		})
	}
}

func TestGeneric_MarshalYAML_StructField(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		Some    option.Generic[int] `yaml:"some"`
		None    option.Generic[int] `yaml:"none"`
		Omitted option.Generic[int] `yaml:"omitted,omitempty"`
		Zero    option.Generic[int] `yaml:"zero,omitempty"`
	}

	data, err := yaml.Marshal(wrapper{
		Some:    option.Some(12),
		None:    option.None[int](),
		Omitted: option.None[int](),
		Zero:    option.Some(0),
	})
	require.NoError(t, err)
	assert.Equal(t, "some: 12\nnone: null\nzero: 0\n", string(data))
}

func TestGeneric_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		var opt option.Generic[int]

		require.NoError(t, yaml.Unmarshal([]byte("42"), &opt))
		require.True(t, opt.IsSome())
		assert.Equal(t, 42, opt.Unwrap())
	})

	t.Run("some_zero", func(t *testing.T) {
		t.Parallel()

		var opt option.Generic[int]

		require.NoError(t, yaml.Unmarshal([]byte("0"), &opt))
		require.True(t, opt.IsSome())
		assert.Equal(t, 0, opt.Unwrap())
	})

	t.Run("some_pointer_receiver", func(t *testing.T) {
		t.Parallel()

		var opt option.Generic[pointerTextType]

		require.NoError(t, yaml.Unmarshal([]byte("text:x"), &opt))
		require.True(t, opt.IsSome())
		assert.Equal(t, pointerTextType{Value: "x"}, opt.Unwrap())
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		opt := option.Some(42)
		err := opt.UnmarshalYAML(func(any) error {
			return errUnmarshalFailed
		})
		require.ErrorIs(t, err, errUnmarshalFailed)
		assert.Equal(t, option.Some(42), opt)
	})

	t.Run("type_error", func(t *testing.T) {
		t.Parallel()

		opt := option.Some(42)
		err := yaml.Unmarshal([]byte("x"), &opt)

		var typeErr *yaml.TypeError

		require.ErrorAs(t, err, &typeErr)
		assert.Equal(t, option.Some(42), opt)
	})
}

func TestGeneric_UnmarshalYAML_Null(t *testing.T) {
	t.Parallel()

	t.Run("null_node_keeps_present_value", func(t *testing.T) {
		t.Parallel()

		// gopkg.in/yaml.v3 does not call UnmarshalYAML for a null node.
		opt := option.Some(42)
		require.NoError(t, yaml.Unmarshal([]byte("null"), &opt))
		assert.Equal(t, option.Some(42), opt)
	})

	t.Run("null_by_direct_call", func(t *testing.T) {
		t.Parallel()

		opt := option.Some(42)
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
			Some   option.Generic[int] `yaml:"some"`
			Null   option.Generic[int] `yaml:"nulled"`
			Tilde  option.Generic[int] `yaml:"tilde"`
			Empty  option.Generic[int] `yaml:"empty"`
			Absent option.Generic[int] `yaml:"absent"`
		}

		var out wrapper

		require.NoError(t, yaml.Unmarshal([]byte("some: 12\nnulled: null\ntilde: ~\nempty:\n"), &out))
		assert.Equal(t, option.Some(12), out.Some)
		assert.False(t, out.Null.IsSome())
		assert.False(t, out.Tilde.IsSome())
		assert.False(t, out.Empty.IsSome())
		assert.False(t, out.Absent.IsSome())
	})
}

func TestGeneric_YAMLRoundTrip(t *testing.T) {
	t.Parallel()

	for _, original := range []option.Generic[string]{
		option.Some("hello"),
		option.Some(""),
		option.None[string](),
	} {
		data, err := yaml.Marshal(original)
		require.NoError(t, err)

		var decoded option.Generic[string]

		require.NoError(t, yaml.Unmarshal(data, &decoded))
		assert.Equal(t, original, decoded)
	}
}

func ExampleGeneric_MarshalYAML() {
	type user struct {
		Name  string                 `yaml:"name"`
		Phone option.Generic[string] `yaml:"phone"`
		Email option.Generic[string] `yaml:"email,omitempty"`
		Age   option.Generic[int]    `yaml:"age"`
	}

	data, err := yaml.Marshal(user{
		Name:  "Maryamu Efe",
		Phone: option.None[string](),
		Email: option.None[string](),
		Age:   option.Some(0),
	})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Print(string(data))

	// Output:
	// name: Maryamu Efe
	// phone: null
	// age: 0
}
