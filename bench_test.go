//nolint:wsl_v5,varnamelen,gocognit,forbidigo,funlen
package option_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/tarantool/go-option"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	errDecUnexpected        = errors.New("unexpected")
	errDecUnexpectedCode    = errors.New("unexpected code")
	errDecUnexpectedExtType = errors.New("unexpected extension type")
	errEncWrite             = errors.New("write failure")
)

func BenchmarkNoneString(b *testing.B) {
	var val string
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.NoneString()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.None[string]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverPtr[string]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverSlice[string]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkNoneInt(b *testing.B) {
	var val int
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.NoneInt()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.None[int]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverPtr[int]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverSlice[int]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkNoneStruct(b *testing.B) {
	var val ValueType
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOptionalStruct()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.None[ValueType]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverPtr[ValueType]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := NoneOverSlice[ValueType]()
			val, ok = opt.Get()
			if ok {
				b.Fatal("Get() returned true")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkSomeInt(b *testing.B) {
	var val int
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.SomeInt(42)
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some(42)
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr(42)
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice(42)
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkSomeString(b *testing.B) {
	var val string
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.SomeString("Hello!")
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some("Hello!")
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr("Hello!")
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice("Hello!")
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkSomeStruct(b *testing.B) {
	var val ValueType
	var ok bool

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOptionalStruct(ValueType{"foo", 42})
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some(ValueType{"foo", 42})
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr(ValueType{"foo", 42})
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice(ValueType{"foo", 42})
			val, ok = opt.Get()
			if !ok {
				b.Fatal("Get() returned false")
			}
		}
	})

	fmt.Println(val, ok)
}

func BenchmarkEncodeDecodeInt(b *testing.B) {
	var buf bytes.Buffer
	buf.Grow(4096)

	enc := msgpack.GetEncoder()
	enc.Reset(&buf)

	dec := msgpack.GetDecoder()
	dec.Reset(&buf)

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.SomeInt(42)

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some(42)

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr(42)

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice(42)

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})
}

func BenchmarkEncodeDecodeString(b *testing.B) {
	var buf bytes.Buffer
	buf.Grow(4096)

	enc := msgpack.GetEncoder()
	enc.Reset(&buf)

	dec := msgpack.GetDecoder()
	dec.Reset(&buf)

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := option.SomeString("Hello!")

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some("Hello!")

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr("Hello!")

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice("Hello!")

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})
}

func BenchmarkEncodeDecodeStruct(b *testing.B) {
	var buf bytes.Buffer
	buf.Grow(4096)

	enc := msgpack.GetEncoder()
	enc.Reset(&buf)

	dec := msgpack.GetDecoder()
	dec.Reset(&buf)

	b.Run("Typed", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOptionalStruct(ValueType{"foo", 42})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("Generic", func(b *testing.B) {
		for b.Loop() {
			opt := option.Some(ValueType{"foo", 42})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericPtr", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverPtr(ValueType{"foo", 42})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})

	b.Run("GenericSlice", func(b *testing.B) {
		for b.Loop() {
			opt := SomeOverSlice(ValueType{"foo", 42})

			err := opt.EncodeMsgpack(enc)
			if err != nil {
				b.Errorf("EncodeMsgpack() failed: %v", err)
			}

			err = opt.DecodeMsgpack(dec)
			if err != nil {
				b.Errorf("DecodeMsgpack() failed: %v", err)
			}

			buf.Reset()
		}
	})
}
