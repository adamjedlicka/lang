package code

type OpCode uint8

// List of OpCodes
const (
	OpConstant OpCode = iota
	OpReturn
)
