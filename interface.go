package option

import (
	"github.com/vmihailenco/msgpack/v5"
)

// commonInterface is the interface that must be implemented by all optional types (generated and hand-written).
//
//nolint:interfacebloat // It is a compile-time check of the full method set, never used as a parameter.
type commonInterface[T any] interface {
	IsSome() bool
	IsZero() bool
	IsNil() bool
	Get() (T, bool)
	MustGet() T
	Unwrap() T
	UnwrapOr(def T) T
	UnwrapOrElse(defCb func() T) T

	EncodeMsgpack(enc *msgpack.Encoder) error
	DecodeMsgpack(dec *msgpack.Decoder) error

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error

	MarshalYAML() (any, error)
	UnmarshalYAML(unmarshal func(any) error) error
}
