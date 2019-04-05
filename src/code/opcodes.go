package code

type OpCode int

// List of OpCodes
const (
	OpConstant OpCode = iota
	OpReturn
)
