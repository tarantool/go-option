package test_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tarantool/go-option/cmd/gentypes/internal/test"
)

const testUUID = `"f47ac10b-58cc-4372-a567-0e02b2c3d479"`

func TestGeneratedTypes_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			"struct_some",
			test.SomeOptionalFullMsgpackExtType(test.FullMsgpackExtType{A: 412, B: "bababa"}),
			`{"A":412,"B":"bababa"}`,
		},
		{
			"struct_some_zero",
			test.SomeOptionalFullMsgpackExtType(test.NewEmptyFullMsgpackExtType()),
			`{"A":0,"B":""}`,
		},
		{"struct_none", test.NoneOptionalFullMsgpackExtType(), `null`},
		{"alias_some", test.SomeOptionalHiddenTypeAlias(test.HiddenTypeAlias{Hidden: "x"}), `{"Hidden":"x"}`},
		{"alias_none", test.NoneOptionalHiddenTypeAlias(), `null`},
		{"pointer_receiver_some", test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), `"text:x"`},
		{"pointer_receiver_none", test.NoneOptionalPointerTextType(), `null`},
		{"third_party_none", test.NoneOptionalUUID(), `null`},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(testCase.value)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, string(data))
		})
	}
}

func TestOptionalFullMsgpackExtType_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	for _, original := range []test.OptionalFullMsgpackExtType{
		test.SomeOptionalFullMsgpackExtType(test.FullMsgpackExtType{A: 412, B: "bababa"}),
		test.SomeOptionalFullMsgpackExtType(test.NewEmptyFullMsgpackExtType()),
		test.NoneOptionalFullMsgpackExtType(),
	} {
		data, err := json.Marshal(original)
		require.NoError(t, err)

		decoded := test.SomeOptionalFullMsgpackExtType(test.FullMsgpackExtType{A: 1, B: "b"})
		require.NoError(t, json.Unmarshal(data, &decoded))
		assert.Equal(t, original, decoded)
	}
}

func TestOptionalHiddenTypeAlias_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := test.SomeOptionalHiddenTypeAlias(test.HiddenTypeAlias{Hidden: "x"})

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded test.OptionalHiddenTypeAlias

	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}

func TestOptionalUUID_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	var opt test.OptionalUUID

	require.NoError(t, json.Unmarshal([]byte(testUUID), &opt))
	require.True(t, opt.IsSome())

	data, err := json.Marshal(opt)
	require.NoError(t, err)
	assert.Equal(t, testUUID, string(data))
}

func TestOptionalPointerTextType_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		var opt test.OptionalPointerTextType

		require.NoError(t, json.Unmarshal([]byte(`"text:x"`), &opt))
		require.True(t, opt.IsSome())
		assert.Equal(t, test.PointerTextType{Value: "x"}, opt.Unwrap())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		require.NoError(t, json.Unmarshal([]byte(`null`), &opt))
		assert.False(t, opt.IsSome())
		assert.Zero(t, opt.Unwrap())
	})

	t.Run("none_with_spaces", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		require.NoError(t, opt.UnmarshalJSON([]byte(" null\n")))
		assert.False(t, opt.IsSome())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		opt := test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"})
		err := opt.UnmarshalJSON([]byte(`"bad"`))
		require.ErrorIs(t, err, test.ErrMissingTextPrefix)
		assert.Equal(t, test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}), opt)
	})
}

func TestGeneratedTypes_JSONStructField(t *testing.T) {
	t.Parallel()

	type wrapper struct {
		Some    test.OptionalPointerTextType `json:"some"`
		None    test.OptionalPointerTextType `json:"none"`
		Omitted test.OptionalPointerTextType `json:"omitted,omitzero"`
		Zero    test.OptionalPointerTextType `json:"zero,omitzero"`
	}

	original := wrapper{
		Some:    test.SomeOptionalPointerTextType(test.PointerTextType{Value: "x"}),
		None:    test.NoneOptionalPointerTextType(),
		Omitted: test.NoneOptionalPointerTextType(),
		Zero:    test.SomeOptionalPointerTextType(test.PointerTextType{Value: ""}),
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `{"some":"text:x","none":null,"zero":"text:"}`, string(data))

	var decoded wrapper

	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}
