package compiler

type Precedence uint8

// List of all precedences.
const (
	PrecedenceNone       Precedence = iota
	PrecedenceAssignment            // =
	PrecedenceOr                    // or
	PrecedenceAnd                   // and
	PrecedenceEquality              // == !=
	PrecedenceComparison            // < > <= >=
	PrecedenceTerm                  // + -
	PrecedenceFactor                // * /
	PrecedenceUnary                 // ! -
	PrecedenceCall                  // . () []
	PrecedencePrimary
)
