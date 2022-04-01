package kv

import "fmt"

const ErrUnknownEngine = "unknown engine: %s"
const ErrConfigNotSet = "%s config not set"

type KVError struct {
	Type    string
	Message string
}

func NewError(typ string, message string) *KVError {
	return &KVError{typ, message}
}

func (e *KVError) Error() string {
	return fmt.Sprintf(e.Type, e.Message)
}
