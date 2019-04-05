package val

type Value interface {
	String() string

	Add(Value) Value
	Subtract(Value) Value
	Multiply(Value) Value
	Divide(Value) Value
	Negate() Value
}
