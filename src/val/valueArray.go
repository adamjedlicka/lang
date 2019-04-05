package val

type ValueArray struct {
	values []Value
}

func NewValueArray() *ValueArray {
	va := new(ValueArray)
	va.values = make([]Value, 0)

	return va
}

func (va *ValueArray) Write(value Value) uint8 {
	va.values = append(va.values, value)

	return va.Len() - 1
}

func (va *ValueArray) GetValue(offset uint8) Value {
	return va.values[offset]
}

func (va *ValueArray) Len() uint8 {
	return uint8(len(va.values))
}
