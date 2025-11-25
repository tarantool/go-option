package option

// This file provides utility functions for decoding and encoding basic types used in MessagePack serialization.
// It includes type checks (e.g., checkNumber, checkString) to ensure compatibility with expected data types and
// wrappers for consistent error handling.
// Note: encodeInt and encodeUint differ from others by converting to int64/uint64 explicitly,
// as the underlying library requires it for correct encoding.

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

func checkNumber(code byte) bool {
	switch {
	case msgpcode.IsFixedNum(code):
		return true
	case code == msgpcode.Int8 || code == msgpcode.Int16 || code == msgpcode.Int32 || code == msgpcode.Int64:
		return true
	case code == msgpcode.Uint8 || code == msgpcode.Uint16 || code == msgpcode.Uint32 || code == msgpcode.Uint64:
		return true
	default:
		return false
	}
}

func decodeInt(decoder *msgpack.Decoder) (int, error) {
	return decoder.DecodeInt() //nolint:wrapcheck
}

func decodeInt8(decoder *msgpack.Decoder) (int8, error) {
	return decoder.DecodeInt8() //nolint:wrapcheck
}

func decodeInt16(decoder *msgpack.Decoder) (int16, error) {
	return decoder.DecodeInt16() //nolint:wrapcheck
}

func decodeInt32(decoder *msgpack.Decoder) (int32, error) {
	return decoder.DecodeInt32() //nolint:wrapcheck
}

func decodeInt64(decoder *msgpack.Decoder) (int64, error) {
	return decoder.DecodeInt64() //nolint:wrapcheck
}

func decodeUint(decoder *msgpack.Decoder) (uint, error) {
	return decoder.DecodeUint() //nolint:wrapcheck
}

func decodeUint8(decoder *msgpack.Decoder) (uint8, error) {
	return decoder.DecodeUint8() //nolint:wrapcheck
}

func decodeUint16(decoder *msgpack.Decoder) (uint16, error) {
	return decoder.DecodeUint16() //nolint:wrapcheck
}

func decodeUint32(decoder *msgpack.Decoder) (uint32, error) {
	return decoder.DecodeUint32() //nolint:wrapcheck
}

func decodeUint64(decoder *msgpack.Decoder) (uint64, error) {
	return decoder.DecodeUint64() //nolint:wrapcheck
}

func encodeInt(encoder *msgpack.Encoder, val int) error {
	return encoder.EncodeInt(int64(val)) //nolint:wrapcheck
}

func encodeInt8(encoder *msgpack.Encoder, val int8) error {
	return encoder.EncodeInt8(val) //nolint:wrapcheck
}

func encodeInt16(encoder *msgpack.Encoder, val int16) error {
	return encoder.EncodeInt16(val) //nolint:wrapcheck
}

func encodeInt32(encoder *msgpack.Encoder, val int32) error {
	return encoder.EncodeInt32(val) //nolint:wrapcheck
}

func encodeInt64(encoder *msgpack.Encoder, val int64) error {
	return encoder.EncodeInt64(val) //nolint:wrapcheck
}

func encodeUint(encoder *msgpack.Encoder, val uint) error {
	return encoder.EncodeUint(uint64(val)) //nolint:wrapcheck
}

func encodeUint8(encoder *msgpack.Encoder, val uint8) error {
	return encoder.EncodeUint8(val) //nolint:wrapcheck
}

func encodeUint16(encoder *msgpack.Encoder, val uint16) error {
	return encoder.EncodeUint16(val) //nolint:wrapcheck
}

func encodeUint32(encoder *msgpack.Encoder, val uint32) error {
	return encoder.EncodeUint32(val) //nolint:wrapcheck
}

func encodeUint64(encoder *msgpack.Encoder, val uint64) error {
	return encoder.EncodeUint64(val) //nolint:wrapcheck
}

func checkFloat(code byte) bool {
	return checkNumber(code) || code == msgpcode.Float || code == msgpcode.Double
}

func decodeFloat32(decoder *msgpack.Decoder) (float32, error) {
	return decoder.DecodeFloat32() //nolint:wrapcheck
}

func encodeFloat32(encoder *msgpack.Encoder, val float32) error {
	return encoder.EncodeFloat32(val) //nolint:wrapcheck
}

func decodeFloat64(decoder *msgpack.Decoder) (float64, error) {
	return decoder.DecodeFloat64() //nolint:wrapcheck
}

func encodeFloat64(encoder *msgpack.Encoder, val float64) error {
	return encoder.EncodeFloat64(val) //nolint:wrapcheck
}

func checkString(code byte) bool {
	return msgpcode.IsBin(code) || msgpcode.IsString(code)
}

func decodeString(decoder *msgpack.Decoder) (string, error) {
	return decoder.DecodeString() //nolint:wrapcheck
}

func encodeString(encoder *msgpack.Encoder, val string) error {
	return encoder.EncodeString(val) //nolint:wrapcheck
}

func checkBytes(code byte) bool {
	return msgpcode.IsBin(code) || msgpcode.IsString(code)
}

func decodeBytes(decoder *msgpack.Decoder) ([]byte, error) {
	return decoder.DecodeBytes() //nolint:wrapcheck
}

func encodeBytes(encoder *msgpack.Encoder, b []byte) error {
	return encoder.EncodeBytes(b) //nolint:wrapcheck
}

func checkBool(code byte) bool {
	return code == msgpcode.True || code == msgpcode.False
}

func decodeBool(decoder *msgpack.Decoder) (bool, error) {
	return decoder.DecodeBool() //nolint:wrapcheck
}

func encodeBool(encoder *msgpack.Encoder, b bool) error {
	return encoder.EncodeBool(b) //nolint:wrapcheck
}

func decodeByte(decoder *msgpack.Decoder) (byte, error) {
	return decoder.DecodeUint8() //nolint:wrapcheck
}

func encodeByte(encoder *msgpack.Encoder, b byte) error {
	return encoder.EncodeUint8(b) //nolint:wrapcheck
}

func checkAny(code byte) bool {
	return code != msgpcode.Nil
}
func encodeAny(encoder *msgpack.Encoder, val any) error {
	return encoder.Encode(val) //nolint:wrapcheck
}

func decodeAny(decoder *msgpack.Decoder) (any, error) {
	return decoder.DecodeInterfaceLoose() //nolint:wrapcheck
}
