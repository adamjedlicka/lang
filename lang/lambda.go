package lang

type Lambda struct {
	declaration LambdaExpr
	closure     *Env
}

func MakeLambda(declaration LambdaExpr, closure *Env) Lambda {
	return Lambda{
		declaration: declaration,
		closure:     closure,
	}
}

func (f Lambda) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
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
	} else if err != nil {
		return nil, err
	}

	return nil, err
}

func (f Lambda) Arity() int {
	return len(f.declaration.params)
}

func (f Lambda) String() string {
	return "<lambda fn>"
}
