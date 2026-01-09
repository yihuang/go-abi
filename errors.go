package abi

import "errors"

// Global error instances to avoid dynamic error creation in generated code.
//
// These errors are used throughout the generated ABI code to provide
// consistent and reusable error handling.

var (
	// ErrInvalidOffsetForDynamicField is returned when the offset for a dynamic field is invalid
	ErrInvalidOffsetForDynamicField = errors.New("invalid offset for dynamic field")

	// ErrInvalidNumberOfTopics is returned when the number of event topics is invalid
	ErrInvalidNumberOfTopics = errors.New("invalid number of topics")

	// ErrInvalidEventTopic is returned when an event topic is invalid
	ErrInvalidEventTopic = errors.New("invalid event topic")

	// ErrInvalidOffsetForSliceElement is returned when the offset for a slice element is invalid
	ErrInvalidOffsetForSliceElement = errors.New("invalid offset for slice element")

	// ErrInvalidOffsetForArrayElement is returned when the offset for an array element is invalid
	ErrInvalidOffsetForArrayElement = errors.New("invalid offset for array element")

	// ErrSizeOverflow is returned when an offset or length is negative
	ErrSizeOverflow = errors.New("size overflow")

	// ErrDirtyPadding is returned when padding bytes are not expected
	ErrDirtyPadding = errors.New("dirty padding")

	// ErrNegativeValue is returned when a negative value is provided for an unsigned type
	ErrNegativeValue = errors.New("negative value for unsigned type")

	// ErrIntegerTooLarge is returned when an integer value exceeds 256 bits
	ErrIntegerTooLarge = errors.New("integer too large")

	// ErrViewIndexOutOfBounds is returned when accessing a slice view with an invalid index
	ErrViewIndexOutOfBounds = errors.New("view index out of bounds")
)
