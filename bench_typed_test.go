package option_test

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type OptionalInt struct {
	value  int
	exists bool
}

func SomeOptionalInt(value int) OptionalInt {
	return OptionalInt{value: value, exists: true}
}

func NoneInt() OptionalInt {
	return OptionalInt{value: zero[int](), exists: false}
}

func (opt *OptionalInt) Get() (int, bool) {
	return opt.value, opt.exists
}

func (opt *OptionalInt) IsSome() bool {
	return opt.exists
}

func (opt *OptionalInt) EncodeMsgpack(encoder *msgpack.Encoder) error {
	if !opt.exists {
		return encoder.EncodeNil()
	}

	return encoder.EncodeInt(int64(opt.value))
}

func (opt *OptionalInt) DecodeMsgpack(decoder *msgpack.Decoder) error {
	val, err := decoder.DecodeInt()
	if err != nil {
		return err
	}

	opt.value = val
	opt.exists = true

	return nil
}

type ValueType struct {
	Value1 string
	Value2 int
}

type OptionalStruct struct {
	value  ValueType
	exists bool
}

func SomeOptionalStruct(value ValueType) OptionalStruct {
	return OptionalStruct{value: value, exists: true}
}

func NoneOptionalStruct() OptionalStruct {
	return OptionalStruct{value: zero[ValueType](), exists: false}
}

func (opt *OptionalStruct) Get() (ValueType, bool) {
	return opt.value, opt.exists
}

func (opt *OptionalStruct) HasValue() bool {
	return opt.exists
}

func (opt *OptionalStruct) EncodeMsgpack(encoder *msgpack.Encoder) error {
	var err error

	if !opt.exists {
		return encoder.EncodeNil()
	}

	err = encoder.EncodeArrayLen(2)
	if err != nil {
		return err
	}

	err = encoder.EncodeString(opt.value.Value1)
	if err != nil {
		return err
	}

	err = encoder.EncodeInt(int64(opt.value.Value2))
	if err != nil {
		return err
	}

	return nil
}

func (opt *OptionalStruct) DecodeMsgpack(decoder *msgpack.Decoder) error {
	arrLen, err := decoder.DecodeArrayLen()
	switch {
	case err != nil:
		return err
	case arrLen == -1:
		opt.exists = false
	case arrLen != 2:
		return fmt.Errorf("%w: unexpected array length: %d", errDecUnexpected, arrLen)
	}

	opt.value.Value1, err = decoder.DecodeString()
	if err != nil {
		return err
	}

	opt.value.Value2, err = decoder.DecodeInt()
	if err != nil {
		return err
	}

	opt.exists = true

	return nil
}
