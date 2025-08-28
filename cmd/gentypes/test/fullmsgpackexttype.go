// Package test contains testing types and scenarios.
package test

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

// FullMsgpackExtType is a test type with both MarshalMsgpack and UnmarshalMsgpack methods.
type FullMsgpackExtType struct {
	A int
	B string
}

// NewEmptyFullMsgpackExtType is an empty constructor for FullMsgpackExtType.
func NewEmptyFullMsgpackExtType() (out FullMsgpackExtType) { //nolint:nonamedreturns
	return
}

// MarshalMsgpack .
func (t *FullMsgpackExtType) MarshalMsgpack() ([]byte, error) {
	var buf bytes.Buffer

	enc := msgpack.NewEncoder(&buf)

	err := enc.EncodeInt(int64(t.A))
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	err = enc.EncodeString(t.B)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return buf.Bytes(), nil
}

// UnmarshalMsgpack .
func (t *FullMsgpackExtType) UnmarshalMsgpack(in []byte) error {
	dec := msgpack.NewDecoder(bytes.NewReader(in))

	a, err := dec.DecodeInt()
	if err != nil {
		return err //nolint:wrapcheck
	}

	t.A = a

	b, err := dec.DecodeString()
	if err != nil {
		return err //nolint:wrapcheck
	}

	t.B = b

	return nil
}
