package option

import (
	"github.com/vmihailenco/msgpack/v5"
)

// commonInterface is the interface that must be implemented by all optional types (generated and hand-written).
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
}
