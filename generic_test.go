package option_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/tarantool/go-option"
)

func ExampleSome() {
	opt := option.Some("hello")

	fmt.Println(opt.IsSome())
	fmt.Println(opt.Unwrap())
	// Output:
	// true
	// hello
}

func ExampleNone() {
	opt := option.None[int]()

	fmt.Println(opt.IsSome())
	fmt.Println(opt.Unwrap())
	// Output:
	// false
	// 0
}

func ExampleGeneric_IsSome() {
	some := option.Some("hello")
	none := option.None[string]()

	fmt.Println(some.IsSome())
	fmt.Println(none.IsSome())
	// Output:
	// true
	// false
}

func ExampleGeneric_IsZero() {
	some := option.Some("hello")
	none := option.None[string]()

	fmt.Println(some.IsZero())
	fmt.Println(none.IsZero())
	// Output:
	// false
	// true
}

func ExampleGeneric_IsNil() {
	some := option.Some("hello")
	none := option.None[string]()

	fmt.Println(some.IsNil() == some.IsZero())
	fmt.Println(none.IsNil() == none.IsZero())
	// Output:
	// true
	// true
}

func ExampleGeneric_Get() {
	some := option.Some(12)
	none := option.None[int]()

	val, ok := some.Get()
	fmt.Println(val, ok)

	val, ok = none.Get()
	fmt.Println(val, ok)
	// Output:
	// 12 true
	// 0 false
}

func ExampleGeneric_MustGet() {
	some := option.Some(12)
	fmt.Println(some.MustGet())
	// Output: 12
}

func ExampleGeneric_MustGet_panic() {
	none := option.None[int]()
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

func ExampleGeneric_Unwrap() {
	some := option.Some(12)
	none := option.None[int]()

	fmt.Println(some.Unwrap())
	fmt.Println(none.Unwrap())
	// Output:
	// 12
	// 0
}

func ExampleGeneric_UnwrapOr() {
	some := option.Some(12)
	none := option.None[int]()

	fmt.Println(some.UnwrapOr(13))
	fmt.Println(none.UnwrapOr(13))
	// Output:
	// 12
	// 13
}

func ExampleGeneric_UnwrapOrElse() {
	some := option.Some(12)
	none := option.None[int]()

	fmt.Println(some.UnwrapOrElse(func() int {
		return 13
	}))
	fmt.Println(none.UnwrapOrElse(func() int {
		return 13
	}))
	// Output:
	// 12
	// 13
}

// TestZeroValueIsZero verifies that the zero value of Generic[T] behaves as None.
func TestZeroValueIsZero(t *testing.T) {
	t.Parallel()

	var opt option.Generic[string]
	assert.True(t, opt.IsZero())
	assert.True(t, opt.IsNil())
	assert.False(t, opt.IsSome())
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
