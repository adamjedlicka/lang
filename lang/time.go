package lang

import (
	"time"
)

type Time struct {
}

func (t Time) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixNano() / 1000), nil
}

func (t Time) Arity() int {
	return 0
}

func (t Time) String() string {
	return "<native fn>"
}
