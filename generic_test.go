package option_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/tarantool/go-option"
)

// TestSomeAndIsSome verifies that Some() creates a valid optional and IsSome returns true.
func TestSomeAndIsSome(t *testing.T) {
	t.Parallel()

	opt := option.Some(42)
	assert.True(t, opt.IsSome())
	assert.False(t, opt.IsZero())
	assert.False(t, opt.IsNil())
}

// TestNoneAndIsZero verifies that None() creates an empty optional and IsZero/IsNil returns true.
func TestNoneAndIsZero(t *testing.T) {
	t.Parallel()

	opt := option.None[int]()
	assert.True(t, opt.IsZero())
	assert.True(t, opt.IsNil())
	assert.False(t, opt.IsSome())
}

// TestZeroValueIsZero verifies that the zero value of Generic[T] behaves as None.
func TestZeroValueIsZero(t *testing.T) {
	t.Parallel()

	var opt option.Generic[string]
	assert.True(t, opt.IsZero())
	assert.True(t, opt.IsNil())
	assert.False(t, opt.IsSome())
}

// TestGetWithValue verifies Get returns the value and true when present.
func TestGetWithValue(t *testing.T) {
	t.Parallel()

	opt := option.Some("hello")
	value, ok := opt.Get()
	assert.True(t, ok)
	assert.Equal(t, "hello", value)
}

// TestGetWithoutValue verifies Get returns zero value and false when absent.
func TestGetWithoutValue(t *testing.T) {
	t.Parallel()

	opt := option.None[float64]()
	value, ok := opt.Get()
	assert.False(t, ok)
	assert.InDelta(t, 0.0, value, 1e-6) // Zero value of float64.
}

// TestMustGetWithValue verifies MustGet returns the value when present.
func TestMustGetWithValue(t *testing.T) {
	t.Parallel()

	opt := option.Some(99)
	value := opt.MustGet()
	assert.Equal(t, 99, value)
}

// TestMustGetPanic verifies MustGet panics when no value is present.
func TestMustGetPanic(t *testing.T) {
	t.Parallel()

	opt := option.None[bool]()
	assert.Panics(t, func() { opt.MustGet() }) //nolint:wsl_v5
}

// TestUnwrapAlias verifies Unwrap is an alias for MustGet.
func TestUnwrapAlias(t *testing.T) {
	t.Parallel()

	opt := option.Some("test")
	assert.Equal(t, "test", opt.Unwrap())

	emptyOpt := option.None[string]()
	assert.Empty(t, emptyOpt.Unwrap())
}

// TestUnwrapOrWithValue verifies UnwrapOr returns inner value when present.
func TestUnwrapOrWithValue(t *testing.T) {
	t.Parallel()

	opt := option.Some("actual")
	assert.Equal(t, "actual", opt.UnwrapOr("default"))
}

// TestUnwrapOrWithoutValue verifies UnwrapOr returns default when absent.
func TestUnwrapOrWithoutValue(t *testing.T) {
	t.Parallel()

	opt := option.None[string]()
	assert.Equal(t, "default", opt.UnwrapOr("default"))
}

// TestUnwrapOrElseWithValue verifies UnwrapOrElse returns inner value when present.
func TestUnwrapOrElseWithValue(t *testing.T) {
	t.Parallel()

	opt := option.Some(100)
	result := opt.UnwrapOrElse(func() int {
		return 200 // Should not be called.
	})
	assert.Equal(t, 100, result)
}

// TestUnwrapOrElseWithoutValue verifies UnwrapOrElse calls func when absent.
func TestUnwrapOrElseWithoutValue(t *testing.T) {
	t.Parallel()

	var (
		called bool
		opt    = option.None[int]()
	)

	result := opt.UnwrapOrElse(func() int {
		called = true

		return 42
	})

	assert.True(t, called)
	assert.Equal(t, 42, result)
}

// TestMsgpackEncodeSome verifies that a Some value is correctly encoded to MessagePack.
func TestMsgpackEncodeSome(t *testing.T) {
	t.Parallel()

	opt := option.Some("hello")

	data, err := msgpack.Marshal(opt)
	require.NoError(t, err)

	var decoded string

	err = msgpack.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "hello", decoded)
}

// TestMsgpackEncodeNone verifies that a None value is encoded as nil in MessagePack.
func TestMsgpackEncodeNone(t *testing.T) {
	t.Parallel()

	opt := option.None[string]()

	data, err := msgpack.Marshal(opt)
	require.NoError(t, err)

	var result *string

	err = msgpack.Unmarshal(data, &result)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestMsgpackDecodeSome verifies that MessagePack data can be decoded into a Some optional.
func TestMsgpackDecodeSome(t *testing.T) {
	t.Parallel()

	data, _ := msgpack.Marshal("hello")

	var opt option.Generic[string]

	err := msgpack.Unmarshal(data, &opt)
	require.NoError(t, err)

	assert.True(t, opt.IsSome())
	assert.Equal(t, "hello", opt.Unwrap())
}

// TestMsgpackDecodeNone verifies that nil MessagePack data becomes a None optional.
func TestMsgpackDecodeNone(t *testing.T) {
	t.Parallel()

	data, _ := msgpack.Marshal(nil)

	var opt option.Generic[int]

	err := msgpack.Unmarshal(data, &opt)
	require.NoError(t, err)

	assert.True(t, opt.IsZero())
	assert.True(t, opt.IsNil())
	assert.False(t, opt.IsSome())
}

// TestRoundTrip verifies full encode-decode roundtrip for both Some and None.
func TestRoundTrip(t *testing.T) {
	t.Parallel()

	// Test roundtrip of Some.
	originalSome := option.Some("roundtrip-test")
	data, err := msgpack.Marshal(originalSome)
	require.NoError(t, err)

	var decodedSome option.Generic[string]

	err = msgpack.Unmarshal(data, &decodedSome)
	require.NoError(t, err)
	assert.True(t, decodedSome.IsSome())
	assert.Equal(t, "roundtrip-test", decodedSome.Unwrap())

	// Test roundtrip of None.
	originalNone := option.None[string]()

	data, err = msgpack.Marshal(originalNone)
	require.NoError(t, err)

	var decodedNone option.Generic[string]

	err = msgpack.Unmarshal(data, &decodedNone)
	require.NoError(t, err)
	assert.True(t, decodedNone.IsZero())
}

// TestCustomEncoderDecoder demonstrates support for types with custom (de)serialization.
type CustomType struct {
	Value string
}

func (c *CustomType) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeString("custom:" + c.Value)
}

func (c *CustomType) DecodeMsgpack(dec *msgpack.Decoder) error {
	s, err := dec.DecodeString()
	if err != nil {
		return err
	}

	c.Value = s[7:] // Strip "custom:".

	return nil
}

func TestCustomEncoderDecoder(t *testing.T) {
	t.Parallel()

	opt := option.Some(CustomType{Value: "test"})

	data, err := msgpack.Marshal(opt)
	require.NoError(t, err)

	var decodedOpt option.Generic[CustomType]

	err = msgpack.Unmarshal(data, &decodedOpt)
	require.NoError(t, err)

	assert.True(t, decodedOpt.IsSome())
	assert.Equal(t, "test", decodedOpt.Unwrap().Value)
}

// TestUnwrapOrElseLazyEvaluation verifies that the default function is not called unnecessarily.
func TestUnwrapOrElseLazyEvaluation(t *testing.T) {
	t.Parallel()

	opt := option.Some(10)

	var called bool

	result := opt.UnwrapOrElse(func() int {
		called = true
		return 999
	})

	assert.Equal(t, 10, result)
	assert.False(t, called, "default function should not be called when value exists")
}
