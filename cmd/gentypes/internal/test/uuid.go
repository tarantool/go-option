package test

import (
	"errors"

	"github.com/google/uuid"
)

func encodeUUID(uuid uuid.UUID) ([]byte, error) {
	return uuid[:], nil
}

var (
	// ErrInvalidLength is returned when the length of the input data is invalid.
	ErrInvalidLength = errors.New("invalid length")
)

func decodeUUID(uuid *uuid.UUID, data []byte) error {
	if len(data) != len(uuid) {
		return ErrInvalidLength
	}

	copy(uuid[:], data)

	return nil
}
