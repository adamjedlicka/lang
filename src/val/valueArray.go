package val

type ValueArray struct {
	values []Value
}

func NewValueArray() *ValueArray {
	va := new(ValueArray)
	va.values = make([]Value, 0)

	return va
}

func (va *ValueArray) Write(value Value) {
	va.values = append(va.values, value)
}

func (va *ValueArray) GetValue(offset int) Value {
	return va.values[offset]
}

func (va *ValueArray) Len() int {
	return len(va.values)
}
