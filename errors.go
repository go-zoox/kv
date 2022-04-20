package kv

import "fmt"

// ErrUnknownEngine means unknown engine error.
const ErrUnknownEngine = "unknown engine: %s"

// ErrConfigNotSet means config not set error.
const ErrConfigNotSet = "%s config not set"

// Error is the error type for KV.
type Error struct {
	Type    string
	Message string
}

// NewError returns a new KVError.
func NewError(typ string, message string) *Error {
	return &Error{typ, message}
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf(e.Type, e.Message)
}
