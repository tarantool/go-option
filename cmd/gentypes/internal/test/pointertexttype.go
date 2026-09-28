package test

import (
	"errors"
	"strings"
)

// ErrMissingTextPrefix is returned when the text form of PointerTextType lacks its prefix.
var ErrMissingTextPrefix = errors.New("missing text: prefix")

// PointerTextType is a test type whose text and MessagePack methods all have
// a pointer receiver, like go-tarantool's BoxError.
type PointerTextType struct {
	Value string
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p *PointerTextType) MarshalText() ([]byte, error) {
	return []byte("text:" + p.Value), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *PointerTextType) UnmarshalText(text []byte) error {
	value, ok := strings.CutPrefix(string(text), "text:")
	if !ok {
		return ErrMissingTextPrefix
	}

	p.Value = value

	return nil
}

// MarshalMsgpack .
func (p *PointerTextType) MarshalMsgpack() ([]byte, error) {
	return []byte(p.Value), nil
}

// UnmarshalMsgpack .
func (p *PointerTextType) UnmarshalMsgpack(in []byte) error {
	p.Value = string(in)

	return nil
}
