package lang

import (
	"fmt"
)

type RuntimeError struct {
	line    int
	message string
}

func NewRuntimeError(line int, message string) RuntimeError {
	e := RuntimeError{}
	e.line = line
	e.message = message

	return e
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %v] RuntimeError: %v", e.line, e.message)
}
