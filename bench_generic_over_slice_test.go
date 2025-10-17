package option_test

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

type GenericOverSlice[T any] []T

func SomeOverSlice[T any](value T) GenericOverSlice[T] {
	return []T{value}
}

func NoneOverSlice[T any]() GenericOverSlice[T] {
	return []T{}
}

func (opt *GenericOverSlice[T]) Get() (T, bool) {
	if len(*opt) == 0 {
		return zero[T](), false
	}

	return (*opt)[0], true
}

func (opt *GenericOverSlice[T]) IsSome() bool {
	return len(*opt) > 0
}

func (opt *GenericOverSlice[T]) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if !opt.IsSome() {
		return encoder.EncodeNil()
	}

	return encoder.Encode((*opt)[0])
}

func (opt *GenericOverSlice[T]) DecodeMsgpack(decoder *msgpack.Decoder) error {
	code, err := decoder.PeekCode()
	if err != nil {
		return err
	}

	if code == msgpcode.Nil {
		*opt = nil
		return nil
	}

	var val T

	err = decoder.Decode(&val)
	if err != nil {
		return err
	}

	*opt = []T{val}

	return nil
}
