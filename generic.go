package option

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// Generic represents an optional value: it may contain a value of type T (Some),
// or it may be empty (None).
//
// This type is useful for safely handling potentially absent values without relying on
// nil pointers or sentinel values, and avoids panics when proper checks are used.
//
// Example:
//
//	opt := option.Some(42)
//	if opt.IsSome() {
//	    fmt.Println(opt.Unwrap()) // prints 42
//	}
//
//	var empty option.Generic[string]
//	fmt.Println(empty.IsZero()) // true
type Generic[T any] struct {
	value  T
	exists bool
}

var _ commonInterface[any] = (*Generic[any])(nil)

// Some creates a Generic[T] containing the given value.
//
// The returned Generic is in the "some" state, meaning IsSome() will return true.
func Some[T any](value T) Generic[T] {
	return Generic[T]{
		value:  value,
		exists: true,
	}
}

// None creates an Generic[T] that does not contain a value.
//
// The returned Generic is in the "none" state, meaning IsZero() will return true.
func None[T any]() Generic[T] {
	return Generic[T]{exists: false} //nolint:exhaustruct
}

// IsSome returns true if the optional contains a value.
func (o Generic[T]) IsSome() bool {
	return o.exists
}

// IsZero returns true if the optional does not contain a value.
func (o Generic[T]) IsZero() bool {
	return !o.exists
}

// IsNil is an alias for IsZero.
//
// This method is provided for compatibility with the msgpack Encoder interface.
func (o Generic[T]) IsNil() bool {
	return o.IsZero()
}

// Get returns the contained value and a boolean indicating whether the value exists.
//
// This is the safest way to extract the value. The second return value is true if a value
// is present, false otherwise. The first return value is the zero value of T when no value exists.
func (o Generic[T]) Get() (T, bool) {
	return o.value, o.exists
}

// MustGet returns the contained value if present.
//
// Panics if the optional is in the "none" state (i.e., no value is present).
//
// Only use this method when you are certain the value exists.
// For safer access, use Get() instead.
func (o Generic[T]) MustGet() T {
	if !o.exists {
		panic("optional value is not set")
	}

	return o.value
}

// Unwrap returns the stored value regardless of presence.
// If no value is set, returns the zero value for T.
//
// Warning: Does not check presence. Use IsSome() before calling if you need
// to distinguish between absent value and explicit zero value.
func (o Generic[T]) Unwrap() T {
	return o.value
}

// UnwrapOr returns the contained value if present, otherwise returns the provided default value.
//
// This is useful when you want to provide a simple fallback value.
func (o Generic[T]) UnwrapOr(defaultValue T) T {
	if o.exists {
		return o.value
	}

	return defaultValue
}

// UnwrapOrElse returns the contained value if present, otherwise calls the provided function
// to compute a default value.
//
// This is useful when the default value is expensive to compute, or requires dynamic logic.
func (o Generic[T]) UnwrapOrElse(defaultValueFunc func() T) T {
	if o.exists {
		return o.value
	}

	return defaultValueFunc()
}

// convertToEncoder checks whether the given value implements msgpack.CustomEncoder.
//
// Used internally during encoding to support custom MessagePack encoding logic.
func convertToEncoder(v any) (msgpack.CustomEncoder, bool) {
	enc, ok := v.(msgpack.CustomEncoder)
	return enc, ok
}

// EncodeMsgpack implements the msgpack.CustomEncoder interface.
//
// If the optional is empty (None), it encodes as a MessagePack nil.
// If the optional contains a value (Some), it attempts to use a custom encoder if the value
// implements msgpack.CustomEncoder; otherwise, it uses the standard encoder.
func (o Generic[T]) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if !o.exists {
		err := encoder.EncodeNil()
		if err != nil {
			return newEncodeGenericError[T](err)
		}

		return nil
	}

	encoderValue, ok := convertToEncoder(&o.value)

	var err error
	if ok {
		err = encoderValue.EncodeMsgpack(encoder)
	} else {
		err = encoder.Encode(&o.value)
	}

	return newEncodeGenericError[T](err)
}

// convertToDecoder checks whether the given value implements msgpack.CustomDecoder.
//
// Used internally during decoding to support custom MessagePack decoding logic.
func convertToDecoder(v any) (msgpack.CustomDecoder, bool) {
	dec, ok := v.(msgpack.CustomDecoder)
	return dec, ok
}

// DecodeMsgpack implements the msgpack.CustomDecoder interface.
//
// It reads a MessagePack value and decodes it into the Generic.
//   - If the encoded value is nil, the optional is set to None.
//   - Otherwise, it decodes into the internal value, using a custom decoder if available,
//     and marks the optional as Some.
//
// Note: This method modifies the receiver and must be called on a pointer.
func (o *Generic[T]) DecodeMsgpack(decoder *msgpack.Decoder) error {
	code, err := decoder.PeekCode()
	switch {
	case err != nil:
		return newDecodeGenericError[T](err)
	case code == msgpcode.Nil:
		o.exists = false

		err := decoder.Skip()
		if err != nil {
			return newDecodeGenericError[T](err)
		}

		return nil
	}

	decoderValue, ok := convertToDecoder(&o.value)
	if ok {
		err = decoderValue.DecodeMsgpack(decoder)
	} else {
		err = decoder.Decode(&o.value)
	}

	if err != nil {
		return newDecodeGenericError[T](err)
	}

	o.exists = true

	return nil
}
