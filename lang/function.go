package lang

type Function struct {
	declaration FnStmnt
	closure     *Env
}

func MakeFunction(declaration FnStmnt, closure *Env) Function {
	return Function{
		declaration: declaration,
		closure:     closure,
	}
}

func (f Function) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	env := MakeEnv(f.closure)

	for i, argument := range arguments {
		err := env.Define(f.declaration.params[i], argument)
		if err != nil {
			return nil, err
		}
	}

	err := i.executeBlock(f.declaration.body, env)
	if returner, ok := err.(Returner); ok {
		return returner.value, nil
	}

	return nil, err
}

func (f Function) Arity() int {
	return len(f.declaration.params)
}

func (f Function) bind(instance *BluInstance) Function {
	env := MakeEnv(f.closure)
	token := MakeToken(This, "this", nil, -1, -1, -1)
	_ = env.Define(token, instance)

	return MakeFunction(f.declaration, env)
}

func (f Function) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
