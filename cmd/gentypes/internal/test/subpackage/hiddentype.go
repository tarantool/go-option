// Package subpackage contains a hidden type for testing type aliases generation.
package subpackage

// Hidden is a hidden type for testing type aliases generation.
type Hidden struct {
	Hidden string
}

// MarshalMsgpack implements the MsgpackMarshaler interface.
func (h *Hidden) MarshalMsgpack() ([]byte, error) {
	return []byte(h.Hidden), nil
}

// UnmarshalMsgpack implements the MsgpackUnmarshaler interface.
func (h *Hidden) UnmarshalMsgpack(bytes []byte) error {
	h.Hidden = string(bytes)
	return nil
}
