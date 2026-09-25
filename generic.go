package option

import (
	"bytes"
	"encoding/json"

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
	var zero T

	return Generic[T]{value: zero, exists: false}
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

// MarshalJSON implements the json.Marshaler interface.
//   - If the optional contains a value (Some), it is encoded exactly as json.Marshal
//     encodes a *T, so MarshalJSON and MarshalText methods of T with a pointer
//     receiver are honoured as well.
//   - If the optional is empty (None), it is encoded as JSON null.
//
// A present value that itself encodes as null (a nil pointer, slice, map or
// interface) is encoded as null too, so it decodes back as None.
//
// A struct field of type Generic[T] with the "omitzero" JSON tag option is
// omitted when the optional is empty, since IsZero reports absence.
// The "omitempty" option has no effect on it.
func (o Generic[T]) MarshalJSON() ([]byte, error) {
	if !o.exists {
		return []byte("null"), nil
	}

	return json.Marshal(&o.value) //nolint:wrapcheck // Errors of T are returned unchanged.
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//   - JSON null is decoded as None.
//   - Any other JSON value is decoded as T and stored as Some.
//
// An error of decoding T is returned unchanged, and the receiver is left intact.
//
// Note: This method modifies the receiver and must be called on a pointer.
func (o *Generic[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*o = None[T]()

		return nil
	}

	var value T

	err := json.Unmarshal(data, &value)
	if err != nil {
		return err //nolint:wrapcheck // Errors of T are returned unchanged.
	}

	*o = Some(value)

	return nil
}

// MarshalYAML implements the Marshaler interface of gopkg.in/yaml.v3
// (and yaml.v2) without depending on it.
//   - If the optional contains a value (Some), it returns a *T, so the value is
//     encoded exactly as the YAML encoder encodes a *T: MarshalYAML and
//     MarshalText methods of T with a pointer receiver are honoured as well.
//   - If the optional is empty (None), it returns nil, which is encoded as null.
//
// A struct field of type Generic[T] with the "omitempty" YAML tag option is
// omitted when the optional is empty, since the encoder consults IsZero.
func (o Generic[T]) MarshalYAML() (any, error) {
	if !o.exists {
		return nil, nil //nolint:nilnil // A nil value is how YAML Marshaler spells null.
	}

	return &o.value, nil
}

// UnmarshalYAML implements the obsolete Unmarshaler interface of
// gopkg.in/yaml.v3 (the only one of yaml.v2) without depending on it.
//   - A null value is decoded as None.
//   - Any other value is decoded as T and stored as Some.
//
// gopkg.in/yaml.v3 never calls this method for a null node (null, ~ or an
// empty value): it leaves the receiver unchanged. Decoding null into a zero
// Generic[T] therefore yields None, but decoding it into Some keeps the value.
//
// An error of decoding T is returned unchanged, and the receiver is left intact.
//
// Note: This method modifies the receiver and must be called on a pointer.
func (o *Generic[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var value *T

	err := unmarshal(&value)
	if err != nil {
		return err
	}

	if value == nil {
		*o = None[T]()

		return nil
	}

	*o = Some(*value)

	return nil
}
