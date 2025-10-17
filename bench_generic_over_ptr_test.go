package option_test

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

type GenericOverPtr[T any] struct {
	val *T
}

func SomeOverPtr[T any](value T) GenericOverPtr[T] {
	return GenericOverPtr[T]{val: &value}
}

func NoneOverPtr[T any]() GenericOverPtr[T] {
	return GenericOverPtr[T]{val: nil}
}

func (opt *GenericOverPtr[T]) Get() (T, bool) {
	if opt.val == nil {
		return zero[T](), false
	}

	return *opt.val, true
}

func (opt *GenericOverPtr[T]) IsSome() bool {
	return opt.val != nil
}

func (opt *GenericOverPtr[T]) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if !opt.IsSome() {
		return encoder.EncodeNil()
	}

	return encoder.Encode(*opt.val)
}

func (opt *GenericOverPtr[T]) DecodeMsgpack(decoder *msgpack.Decoder) error {
	code, err := decoder.PeekCode()
	switch {
	case err != nil:
		return err
	case code == msgpcode.Nil:
		opt.val = nil
		return nil
	}

	var val T

	err = decoder.Decode(&val)
	if err != nil {
		return err
	}

	opt.val = &val

	return nil
}
