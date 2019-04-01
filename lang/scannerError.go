package lang

import (
	"fmt"
)

type ScannerError struct {
	line    int
	message string
}

func NewScannerError(line int, message string) ScannerError {
	e := ScannerError{}
	e.line = line
	e.message = message

	return e
}

func (e ScannerError) Error() string {
	return fmt.Sprintf("[line %v] Error: %v", e.line, e.message)
}
