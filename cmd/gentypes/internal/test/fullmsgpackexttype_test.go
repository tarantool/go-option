package test_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/tarantool/go-option/cmd/gentypes/internal/test"
)

func TestOptionalMsgpackExtType_RoundtripLL(t *testing.T) {
	t.Parallel()

	input := test.FullMsgpackExtType{
		A: 412,
		B: "bababa",
	}

	opt := test.SomeOptionalFullMsgpackExtType(input)

	b := bytes.Buffer{}
	enc := msgpack.NewEncoder(&b)
	dec := msgpack.NewDecoder(&b)

	require.NoError(t, opt.EncodeMsgpack(enc))

	opt2 := test.NoneOptionalFullMsgpackExtType()
	require.NoError(t, opt2.DecodeMsgpack(dec))

	assert.Equal(t, opt, opt2)
	assert.Equal(t, input, opt2.Unwrap())
}

func TestOptionalMsgpackExtType_RoundtripHL(t *testing.T) {
	t.Parallel()

	input := test.FullMsgpackExtType{
		A: 412,
		B: "bababa",
	}

	opt := test.SomeOptionalFullMsgpackExtType(input)

	b := bytes.Buffer{}
	enc := msgpack.NewEncoder(&b)
	dec := msgpack.NewDecoder(&b)

	require.NoError(t, enc.Encode(opt))

	opt2 := test.NoneOptionalFullMsgpackExtType()
	require.NoError(t, dec.Decode(&opt2))

	assert.Equal(t, opt, opt2)
	assert.Equal(t, input, opt2.Unwrap())
}

func TestOptionalFullMsgpackExtType_IsSome(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		assert.True(t, opt.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.NoneOptionalFullMsgpackExtType()

		assert.False(t, opt.IsSome())
	})
}

func TestOptionalFullMsgpackExtType_IsZero(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		assert.False(t, opt.IsZero())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.NoneOptionalFullMsgpackExtType()

		assert.True(t, opt.IsZero())
	})
}

func TestOptionalFullMsgpackExtType_Get(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		val, ok := opt.Get()
		require.True(t, ok)
		assert.Equal(t, input, val)
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.NoneOptionalFullMsgpackExtType()
		val, ok := opt.Get()
		require.False(t, ok)
		assert.Equal(t, test.NewEmptyFullMsgpackExtType(), val)
	})
}

func TestOptionalFullMsgpackExtType_MustGet(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		var val test.FullMsgpackExtType

		require.NotPanics(t, func() {
			val = opt.MustGet()
		})
		assert.Equal(t, input, val)
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.NoneOptionalFullMsgpackExtType()

		require.Panics(t, func() { opt.MustGet() })
	})
}

func TestOptionalFullMsgpackExtType_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		assert.Equal(t, input, opt.Unwrap())
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		opt := test.NoneOptionalFullMsgpackExtType()
		assert.Equal(t, test.NewEmptyFullMsgpackExtType(), opt.Unwrap())
	})
}

func TestOptionalFullMsgpackExtType_UnwrapOr(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		assert.Equal(t, input, opt.UnwrapOr(test.NewEmptyFullMsgpackExtType()))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		alt := test.FullMsgpackExtType{
			A: 1,
			B: "b",
		}

		opt := test.NoneOptionalFullMsgpackExtType()
		assert.Equal(t, alt, opt.UnwrapOr(alt))
	})
}

func TestOptionalFullMsgpackExtType_UnwrapOrElse(t *testing.T) {
	t.Parallel()

	t.Run("some", func(t *testing.T) {
		t.Parallel()

		input := test.FullMsgpackExtType{
			A: 412,
			B: "bababa",
		}

		opt := test.SomeOptionalFullMsgpackExtType(input)

		assert.Equal(t, input, opt.UnwrapOrElse(test.NewEmptyFullMsgpackExtType))
	})

	t.Run("none", func(t *testing.T) {
		t.Parallel()

		alt := test.FullMsgpackExtType{
			A: 1,
			B: "b",
		}

		opt := test.NoneOptionalFullMsgpackExtType()

		assert.Equal(t, alt, opt.UnwrapOrElse(func() test.FullMsgpackExtType {
			return alt
		}))
	})
}
