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

	// ErrNegativeSize is returned when an offset or length is negative
	ErrNegativeSize = errors.New("negative size")

	// ErrDirtyPadding is returned when padding bytes are not expected
	ErrDirtyPadding = errors.New("dirty padding")
)
