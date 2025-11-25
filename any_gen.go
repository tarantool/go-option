package option

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// Any represents an optional value of type any.
// It can either hold a valid Any (any type like string or int, float).
// (IsSome == true) or be empty (IsZero == true).
type Any struct { //nolint:recvcheck
	value  any
	exists bool
}

var _ commonInterface[any] = (*Any)(nil)

// SomeAny creates an optional Any with the given value.
// The returned Any (any type) will have IsSome() == true and IsZero() == false.
//
// Example:
//
//	o := SomeAny(7777.7777777)
//	if o.IsSome() {
//	    v := o.Unwrap() // v == true
//	}
func SomeAny(value any) Any {
	return Any{
		value:  value,
		exists: true,
	}
}

// NoneAny creates an empty optional any value.
// The returned Any will have IsSome() == false and IsZero() == true.
//
// Example:
//
//	o := NoneAny()
//	if o.IsZero() {
//	    fmt.Println("value is absent")
//	}
func NoneAny() Any {
	return Any{
		exists: false,
		value:  zero[any](),
	}
}

// IsSome returns true if the Any contains a value.
// This indicates the value is explicitly set (not None).
func (o Any) IsSome() bool {
	return o.exists
}

// IsZero returns true if the Any does not contain a value.
// Equivalent to !IsSome(). Useful for consistency with types where
// zero value (e.g. 0, false, zero struct) is valid and needs to be distinguished.
func (o Any) IsZero() bool {
	return !o.exists
}

// IsNil is an alias for IsZero.
//
// This method is provided for compatibility with the msgpack Encoder interface.
func (o Any) IsNil() bool {
	return o.IsZero()
}

// Get returns the stored value and a boolean flag indicating its presence.
// If the value is present, returns (value, true).
// If the value is absent, returns (zero value of byte, false).
//
// Recommended usage:
//
//	if value, ok := o.Get(); ok {
//	    // use value
//	}
func (o Any) Get() (any, bool) {
	return o.value, o.exists
}

// MustGet returns the stored value if it is present.
// Panics if the value is absent (i.e., IsZero() == true).
//
// Use with caution — only when you are certain the value exists.
//
// Panics with: "optional value is not set" if no value is set.
func (o Any) MustGet() any {
	if !o.exists {
		panic("optional value is not set!")
	}

	return o.value
}

// Unwrap returns the stored value regardless of presence.
// If no value is set, returns the zero value for byte.
//
// Warning: Does not check presence. Use IsSome() before calling if you need
// to distinguish between absent value and explicit zero value.
func (o Any) Unwrap() any {
	return o.value
}

// UnwrapOr returns the stored value if present.
// Otherwise, returns the provided default value.
//
// Example:
//
//	o := NoneAny()
//	v := o.UnwrapOr(someDefaultByte)
func (o Any) UnwrapOr(defaultValue any) any {
	if o.exists {
		return o.value
	}

	return defaultValue
}

// UnwrapOrElse returns the stored value if present.
// Otherwise, calls the provided function and returns its result.
// Useful when the default value requires computation or side effects.
//
// Example:
//
//	o := NoneAny()
//	v := o.UnwrapOrElse(func() any { return computeDefault() })
func (o Any) UnwrapOrElse(defaultValue func() any) any {
	if o.exists {
		return o.value
	}

	return defaultValue()
}

// EncodeMsgpack encodes the Any value using MessagePack format.
// - If the value is present, it is encoded as byte.
// - If the value is absent (None), it is encoded as nil.
//
// Returns an error if encoding fails.
func (o Any) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if o.exists {
		return newEncodeError("Any", encodeAny(encoder, o.value))
	}

	return newEncodeError("Any", encoder.EncodeNil())
}

// DecodeMsgpack decodes a Any value from MessagePack format.
// Supports two input types:
//   - nil: interpreted as no value (NoneByte)
//   - byte: interpreted as a present value (SomeByte)
//
// Returns an error if the input type is unsupported or decoding fails.
//
// After successful decoding:
//   - on nil: exists = false, value = default zero value
//   - on any: exists = true, value = decoded value
func (o *Any) DecodeMsgpack(decoder *msgpack.Decoder) error {
	code, err := decoder.PeekCode()
	if err != nil {
		return newDecodeError("Any", err)
	}

	switch {
	case code == msgpcode.Nil:
		o.exists = false

		return newDecodeError("Any", decoder.Skip())
	case checkAny(code):
		o.value, err = decodeAny(decoder)
		if err != nil {
			return newDecodeError("Any", err)
		}

		o.exists = true

		return err
	default:
		return newDecodeWithCodeError("Any", code)
	}
}
