package lang

type Function struct {
	declaration FnStmnt
}

func MakeFunction(declaration FnStmnt) Function {
	return Function{
		declaration: declaration,
	}
}

func (f Function) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	env := MakeEnv(i.globals)

	for i, argument := range arguments {
		env.Define(f.declaration.params[i], argument)
	}

	block := f.declaration.body.(BlockStmnt)

	return nil, i.executeBlock(block.stmnts, env)
}

func (f Function) Arity() int {
	return len(f.declaration.params)
}

func (f Function) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
