package val

import "fmt"

type Number struct {
	value float64
}

func NewNumber(value float64) Number {
	return Number{
		value: value,
	}
}

func (n Number) Add(other Value) Value {
	return NewNumber(n.value + (other.(Number)).value)
}

func (n Number) Subtract(other Value) Value {
	return NewNumber(n.value - (other.(Number)).value)
}

func (n Number) Multiply(other Value) Value {
	return NewNumber(n.value * (other.(Number)).value)
}

func (n Number) Divide(other Value) Value {
	return NewNumber(n.value / (other.(Number)).value)
}

func (n Number) Negate() Value {
	return NewNumber(-n.value)
}

func (n Number) String() string {
	return fmt.Sprintf("%v", n.value)
}
