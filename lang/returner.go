package lang

import (
	"fmt"
)

type Returner struct {
	value interface{}
}

func MakeReturner(value interface{}) Returner {
	return Returner{
		value: value,
	}
}

func (r Returner) Error() string {
	return fmt.Sprintf("RETURN: %v", r.value)
}
