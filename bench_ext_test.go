//nolint:varnamelen,gocognit,funlen
package option_test

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/tarantool/go-option"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

type BenchExt struct {
	data []byte
}

func (e *BenchExt) MarshalMsgpack() ([]byte, error) {
	return e.data, nil
}

func (e *BenchExt) UnmarshalMsgpack(bytes []byte) error {
	e.data = slices.Clone(bytes)
	return nil
}

func (e *BenchExt) ExtType() int8 {
	return 8
}

type Optional1BenchExt struct {
	value  BenchExt
	exists bool
}

func SomeOptional1BenchExt(value BenchExt) Optional1BenchExt {
	return Optional1BenchExt{value: value, exists: true}
}

func NoneOptional1BenchExt() Optional1BenchExt {
	return Optional1BenchExt{value: zero[BenchExt](), exists: false}
}

var (
	//nolint: gochecknoglobals
	NilBytes = []byte{msgpcode.Nil}
)

func (opt *Optional1BenchExt) MarshalMsgpack() ([]byte, error) {
	if opt.exists {
		return opt.value.MarshalMsgpack()
	}

	return NilBytes, nil
}

func (opt *Optional1BenchExt) UnmarshalMsgpack(bytes []byte) error {
	if bytes[0] == msgpcode.Nil {
		opt.exists = false
		return nil
	}

	opt.exists = true

	return opt.value.UnmarshalMsgpack(bytes)
}

func (opt *Optional1BenchExt) EncodeMsgpack(enc *msgpack.Encoder) error {
	switch {
	case !opt.exists:
		return enc.EncodeNil()
	default:
		mpdata, err := opt.value.MarshalMsgpack()
		if err != nil {
			return err
		}

		err = enc.EncodeExtHeader(opt.value.ExtType(), len(mpdata))
		if err != nil {
			return err
		}

		mpdataLen, err := enc.Writer().Write(mpdata)
		switch {
		case err != nil:
			return err
		case mpdataLen != len(mpdata):
			return fmt.Errorf("%w: failed to write mpdata", errEncWrite)
		}

		return nil
	}
}

func (opt *Optional1BenchExt) DecodeMsgpack(dec *msgpack.Decoder) error {
	code, err := dec.PeekCode()
	if err != nil {
		return err
	}

	switch {
	case code == msgpcode.Nil:
		opt.exists = false
	case msgpcode.IsExt(code):
		extID, extLen, err := dec.DecodeExtHeader()
		switch {
		case err != nil:
			return err
		case extID != opt.value.ExtType():
			return fmt.Errorf("%w: %d", errDecUnexpectedExtType, extID)
		default:
			ext := make([]byte, extLen)

			err := dec.ReadFull(ext)
			if err != nil {
				return err
			}

			err = opt.value.UnmarshalMsgpack(ext)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%w: %x", errDecUnexpectedCode, code)
	}

	return nil
}

type Optional2BenchExt struct {
	value  BenchExt
	exists bool
}

func SomeOptional2BenchExt(value BenchExt) Optional2BenchExt {
	return Optional2BenchExt{value: value, exists: true}
}

func NoneOptional2BenchExt() Optional2BenchExt {
	return Optional2BenchExt{value: zero[BenchExt](), exists: false}
}

func (opt *Optional2BenchExt) MarshalMsgpack() ([]byte, error) {
	if opt.exists {
		return opt.value.MarshalMsgpack()
	}

	return NilBytes, nil
}

func (opt *Optional2BenchExt) UnmarshalMsgpack(b []byte) error {
	if b[0] == msgpcode.Nil {
		opt.exists = false
		return nil
	}

	opt.exists = true

	return opt.value.UnmarshalMsgpack(b)
}

func (opt *Optional2BenchExt) EncodeMsgpack(enc *msgpack.Encoder) error {
	switch {
	case !opt.exists:
		return enc.EncodeNil()
	default:
		return enc.Encode(&opt.value)
	}
}

func (opt *Optional2BenchExt) DecodeMsgpack(dec *msgpack.Decoder) error {
	code, err := dec.PeekCode()
	if err != nil {
		return err
	}

	switch {
	case code == msgpcode.Nil:
		opt.exists = false
		return nil
	case msgpcode.IsExt(code):
		return dec.Decode(&opt.value)
	default:
		return fmt.Errorf("%w: %x", errDecUnexpectedCode, code)
	}
}

type MsgpackExtInterface interface {
	ExtType() int8
	msgpack.Marshaler
	msgpack.Unmarshaler
}

type OptionalGenericStructWithInterface[T MsgpackExtInterface] struct {
	value  T
	exists bool
}

func SomeOptionalGenericStructWithInterface[T MsgpackExtInterface](value T) *OptionalGenericStructWithInterface[T] {
	return &OptionalGenericStructWithInterface[T]{
		value:  value,
		exists: true,
	}
}

func NoneOptionalGenericStructWithInterface[T MsgpackExtInterface]() *OptionalGenericStructWithInterface[T] {
	return &OptionalGenericStructWithInterface[T]{
		value:  zero[T](),
		exists: false,
	}
}

func (opt *OptionalGenericStructWithInterface[T]) DecodeMsgpack(dec *msgpack.Decoder) error {
	code, err := dec.PeekCode()
	if err != nil {
		return err
	}

	switch {
	case code == msgpcode.Nil:
		opt.exists = false
	case msgpcode.IsExt(code):
		opt.exists = true

		extID, extLen, err := dec.DecodeExtHeader()
		switch {
		case err != nil:
			return err
		case extID != opt.value.ExtType():
			return fmt.Errorf("%w: %d", errDecUnexpectedExtType, extID)
		default:
			ext := make([]byte, extLen)

			err := dec.ReadFull(ext)
			if err != nil {
				return err
			}

			err = opt.value.UnmarshalMsgpack(ext)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("%w: %x", errDecUnexpectedCode, code)
	}

	return nil
}

func (opt *OptionalGenericStructWithInterface[T]) EncodeMsgpack(enc *msgpack.Encoder) error {
	switch {
	case !opt.exists:
		return enc.EncodeNil()
	default:
		mpdata, err := opt.value.MarshalMsgpack()
		if err != nil {
			return err
		}

		err = enc.EncodeExtHeader(opt.value.ExtType(), len(mpdata))
		if err != nil {
			return err
		}

		mpdataLen, err := enc.Writer().Write(mpdata)
		switch {
		case err != nil:
			return err
		case mpdataLen != len(mpdata):
			return fmt.Errorf("%w: failed to write mpdata", errEncWrite)
		}

		return nil
	}
}

func BenchmarkExtension(b *testing.B) {
	msgpack.RegisterExt(8, &BenchExt{nil})

	var buf bytes.Buffer
	buf.Grow(4096)

	enc := msgpack.GetEncoder()
	enc.Reset(&buf)

	dec := msgpack.GetDecoder()
	dec.Reset(&buf)

	b.Run("Optional1Bench", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOptional1BenchExt(BenchExt{[]byte{1, 2, 3}})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Fatal(err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Optional2Bench", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOptional2BenchExt(BenchExt{[]byte{1, 2, 3}})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Fatal(err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("OptionalBenchGeneric", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some(BenchExt{[]byte{1, 2, 3}})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Fatal(err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("OptionalGenericStructWithInterface", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOptionalGenericStructWithInterface(&BenchExt{[]byte{1, 2, 3}})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Fatal(err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Default", func(b *testing.B) {
		for b.Loop() {
			opt := BenchExt{[]byte{1, 2, 3}}

			err := enc.Encode(&opt)
			if err != nil {
				b.Fatal(err)
			}

			err = dec.Decode(&opt)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
