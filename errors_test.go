package option //nolint:testpackage
// this is unit test, that checks internal logic of error constructing.

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("some error")
)

func TestDecodeError_Error(t *testing.T) {
	t.Parallel()

	t.Run("newEncodeError with error", func(t *testing.T) {
		t.Parallel()

		a := newEncodeError("Byte", errTest)

		require.Error(t, a)
		assert.Equal(t, "failed to encode Byte: some error", a.Error())
	})

	t.Run("newEncodeError without error", func(t *testing.T) {
		t.Parallel()

		a := newEncodeError("Byte", nil)

		require.NoError(t, a)
	})
}

func TestEncodeError_Error(t *testing.T) {
	t.Parallel()

	t.Run("newDecodeError with error", func(t *testing.T) {
		t.Parallel()

		a := newDecodeError("Byte", errTest)

		require.Error(t, a)
		assert.Equal(t, "failed to decode Byte: some error", a.Error())
	})

	t.Run("newDecodeError without error", func(t *testing.T) {
		t.Parallel()

		a := newDecodeError("Byte", nil)

		require.NoError(t, a)
	})

	t.Run("newDecodeWithCodeError", func(t *testing.T) {
		t.Parallel()

		a := newDecodeWithCodeError("Byte", 1)

		require.Error(t, a)
		assert.Equal(t, "failed to decode Byte, invalid code: 1", a.Error())
	})
}
