package option

import (
	"fmt"
)

// DecodeError is returned when decoding failed due to invalid code in msgpack stream.
type DecodeError struct {
	Type   string
	Code   Byte
	Parent error
}

// Error returns the text representation of error.
func (d DecodeError) Error() string {
	if d.Code.IsSome() {
		return fmt.Sprintf("failed to decode %s, invalid code: %d", d.Type, d.Code.Unwrap())
	}

	return fmt.Sprintf("failed to decode %s: %s", d.Type, d.Parent)
}

func newDecodeWithCodeError(operationType string, code byte) error {
	return DecodeError{
		Type:   operationType,
		Code:   SomeByte(code),
		Parent: nil,
	}
}

func newDecodeError(operationType string, err error) error {
	if err == nil {
		return nil
	}

	return DecodeError{
		Type:   operationType,
		Code:   NoneByte(),
		Parent: err,
	}
}

// EncodeError is returned when encoding failed due to stream errors.
type EncodeError struct {
	Type   string
	Parent error
}

// Error returns the text representation of error.
func (e EncodeError) Error() string {
	return fmt.Sprintf("failed to encode %s: %s", e.Type, e.Parent)
}

func newEncodeError(operationType string, err error) error {
	if err == nil {
		return nil
	}

	return EncodeError{Type: operationType, Parent: err}
}
