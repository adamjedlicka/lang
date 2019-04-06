package val

import (
	"math/big"
)

type Number big.Float

func NewNumber(lexeme string) *Number {
	f, _, err := new(big.Float).Parse(lexeme, 10)
	if err != nil {
		panic(err)
	}

	return (*Number)(f)
}

func (n *Number) Add(other Value) Value {
	return (*Number)((*big.Float)(n).Add((*big.Float)(n), (*big.Float)(other.(*Number))))
}

func (n *Number) Subtract(other Value) Value {
	return (*Number)((*big.Float)(n).Sub((*big.Float)(n), (*big.Float)(other.(*Number))))
}

func (n *Number) Multiply(other Value) Value {
	return (*Number)((*big.Float)(n).Mul((*big.Float)(n), (*big.Float)(other.(*Number))))
}

func (n *Number) Divide(other Value) Value {
	return (*Number)((*big.Float)(n).Quo((*big.Float)(n), (*big.Float)(other.(*Number))))
}

func (n *Number) Negate() Value {
	return (*Number)((*big.Float)(n).Neg((*big.Float)(n)))
}

func (n *Number) String() string {
	return (*big.Float)(n).String()
}
