package val

import "fmt"

type Number struct {
	value float64
}

func NewNumber(value float64) *Number {
	n := new(Number)
	n.value = value

	return n
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}
