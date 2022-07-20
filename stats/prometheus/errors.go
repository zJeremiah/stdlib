package prometheus

import "errors"

var (
	ErrNoKey       = errors.New("no such key")
	ErrInvalidType = errors.New("collector with key was not the expected type")
)
