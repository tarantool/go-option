package test_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/tarantool/go-option/cmd/gentypes/internal/test"
)

var errUnmarshalFailed = errors.New("unmarshal failed")

func TestGeneratedTypes_MarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			"struct_some",
			test.SomeOptionalFullMsgpackExtType(test.FullMsgpackExtType{A: 412, B: "bababa"}),
			"a: 412\nb: bababa\n",
		},
		{
			"struct_some_zero",
			test.SomeOptionalFullMsgpackExtType(test.NewEmptyFullMsgpackExtType()),
			"a: 0\nb: \"\"\n",
		},
		{"struct_none", test.NoneOptionalFullMsgpackExtType(), "null\n"},
		{"alias_some", test.SomeOptionalHiddenTypeAlias(test.HiddenTypeAlias{Hidden: "x"}), "hidden: x\n"},
		{"alias_none", test.NoneOptionalHiddenTypeAlias(), "null\n"},
		{"pointer_receiver_some", test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), "text:x\n"},
		{"pointer_receiver_none", test.NoneOptionalPointerTextType(), "null\n"},
		{"third_party_none", test.NoneOptionalUUID(), "null\n"},
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

func TestOptionalFullMsgpackExtType_YAMLRoundTrip(t *testing.T) {
	t.Parallel()

	for _, original := range []test.OptionalFullMsgpackExtType{
		test.SomeOptionalFullMsgpackExtType(test.FullMsgpackExtType{A: 412, B: "bababa"}),
		test.SomeOptionalFullMsgpackExtType(test.NewEmptyFullMsgpackExtType()),
		test.NoneOptionalFullMsgpackExtType(),
	} {
		data, err := yaml.Marshal(original)
		require.NoError(t, err)

		var decoded test.OptionalFullMsgpackExtType

		require.NoError(t, yaml.Unmarshal(data, &decoded))
		assert.Equal(t, original, decoded)
	}
}

func TestOptionalHiddenTypeAlias_YAMLRoundTrip(t *testing.T) {
	t.Parallel()

	original := test.SomeOptionalHiddenTypeAlias(test.HiddenTypeAlias{Hidden: "x"})

	data, err := yaml.Marshal(original)
	require.NoError(t, err)

	var decoded test.OptionalHiddenTypeAlias

	require.NoError(t, yaml.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}

func TestOptionalUUID_YAMLRoundTrip(t *testing.T) {
	t.Parallel()

	const text = "f47ac10b-58cc-4372-a567-0e02b2c3d479\n"

	var opt test.OptionalUUID

	require.NoError(t, yaml.Unmarshal([]byte(text), &opt))
	require.True(t, opt.IsSome())

	data, err := yaml.Marshal(opt)
	require.NoError(t, err)
	assert.Equal(t, text, string(data))
}

func TestOptionalPointerTextType_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		var opt test.OptionalPointerTextType

		require.NoError(t, yaml.Unmarshal([]byte("text:x"), &opt))
		require.True(t, opt.IsSome())
		assert.Equal(t, test.PointerTextType{Value: "x"}, opt.Unwrap())
	})

	t.Run("null_node_keeps_present_value", func(t *testing.T) {
		t.Parallel()

		// gopkg.in/yaml.v3 does not call UnmarshalYAML for a null node.
		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		require.NoError(t, yaml.Unmarshal([]byte("null"), &opt))
		assert.Equal(t, test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), opt)
	})

	t.Run("null_by_direct_call", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		err := opt.UnmarshalYAML(func(v any) error {
			return yaml.Unmarshal([]byte("null"), v)
		})
		require.NoError(t, err)
		assert.False(t, opt.IsSome())
		assert.Zero(t, opt.Unwrap())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		err := yaml.Unmarshal([]byte("bad"), &opt)
		require.ErrorIs(t, err, test.ErrMissingTextPrefix)
		assert.Equal(t, test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), opt)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		err := opt.UnmarshalYAML(func(any) error {
			return errUnmarshalFailed
		})
		require.ErrorIs(t, err, errUnmarshalFailed)
		assert.Equal(t, test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), opt)
	})
}

func TestGeneratedTypes_YAMLStructField(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		Some    test.OptionalPointerTextType `yaml:"some"`
		None    test.OptionalPointerTextType `yaml:"none"`
		Omitted test.OptionalPointerTextType `yaml:"omitted,omitempty"`
		Zero    test.OptionalPointerTextType `yaml:"zero,omitempty"`
	}

	original := wrapper{
		Some:    test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}),
		None:    test.NoneOptionalPointerTextType(),
		Omitted: test.NoneOptionalPointerTextType(),
		Zero:    test.SomeOptionalPointerTextType(test.PointerTextType{Value: ""}),
	}

	data, err := yaml.Marshal(original)
	require.NoError(t, err)
	assert.Equal(t, "some: text:x\nnone: null\nzero: 'text:'\n", string(data))

	var decoded wrapper

	require.NoError(t, yaml.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}
