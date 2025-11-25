package option_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/tarantool/go-option"
)

func TestAny_IsSome(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny("aaaaaa+++++")
		assert.True(t, someAny.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		assert.False(t, emptyAny.IsSome())
	})
}

func TestAny_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny(12232.777777777)
		assert.False(t, someAny.IsZero())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		assert.True(t, emptyAny.IsZero())
	})
}

func TestAny_IsNil(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny(777.111111)
		assert.False(t, someAny.IsNil())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		assert.True(t, emptyAny.IsNil())
	})
}

func TestAny_Get(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny("lllllllll")
		val, ok := someAny.Get()
		require.True(t, ok)
		assert.EqualValues(t, "lllllllll", val)
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		_, ok := emptyAny.Get()
		require.False(t, ok)
	})
}

func TestAny_MustGet(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny(1111.1000000)
		assert.InEpsilon(t, 1111.1000000, someAny.MustGet(), 0.01)
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()

		assert.Panics(t, func() {
			emptyAny.MustGet()
		})
	})
}

func TestAny_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny(
			"HH77771111111111111111111111111111111111111111111111111111111111111111111111")
		assert.EqualValues(t,
			"HH77771111111111111111111111111111111111111111111111111111111111111111111111",
			someAny.Unwrap())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()

		assert.NotPanics(t, func() {
			emptyAny.Unwrap()
		})
	})
}

func TestAny_UnwrapOr(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny("(((09,111")
		assert.EqualValues(t, "(((09,111", someAny.UnwrapOr(111111))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		assert.InEpsilon(t, 11111.8880, emptyAny.UnwrapOr(11111.8880), 0.01)
	})
}

func TestAny_UnwrapOrElse(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		someAny := option.SomeAny(34534534)
		assert.EqualValues(t, 34534534, someAny.UnwrapOrElse(func() any {
			return "EXAMPLE!!!"
		}))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		emptyAny := option.NoneAny()
		assert.EqualValues(t, 145, emptyAny.UnwrapOrElse(func() any {
			return 145
		}))
	})
}

func TestAny_EncodeDecodeMsgpack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    any
		expected any
	}{
		{"string", "test string", "test string"},
		{"int", 42, int64(42)},
		{"float", 3.14, 3.14},
		{"bool", true, true},
		{"slice", []int{1, 2, 3}, []any{int8(1), int8(2), int8(3)}},
		{"map", map[string]int{"a": 1}, map[string]any{"a": int8(1)}},
	}

	for _, testCase := range testCases {
		t.Run("some_"+testCase.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			enc := msgpack.NewEncoder(&buf)
			dec := msgpack.NewDecoder(&buf)

			// Encode.
			someAny := option.SomeAny(testCase.value)
			err := someAny.EncodeMsgpack(enc)
			require.NoError(t, err)

			// Decode.
			var unmarshaled option.Any

			err = unmarshaled.DecodeMsgpack(dec)
			require.NoError(t, err)
			assert.True(t, unmarshaled.IsSome())
			assert.Equal(t, testCase.expected, unmarshaled.Unwrap())
		})
	}

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		enc := msgpack.NewEncoder(&buf)
		dec := msgpack.NewDecoder(&buf)

		// Encode nil.
		emptyAny := option.NoneAny()

		err := emptyAny.EncodeMsgpack(enc)
		require.NoError(t, err)

		// Decode.
		var unmarshaled option.Any

		err = unmarshaled.DecodeMsgpack(dec)
		require.NoError(t, err)

		// Verify it's none.
		assert.False(t, unmarshaled.IsSome())
		assert.Nil(t, unmarshaled.Unwrap())
	})
}
