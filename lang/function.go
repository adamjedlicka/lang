package lang

type Function struct {
	declaration FnStmnt
	closure     *Env
	isInit      bool
}

func MakeFunction(declaration FnStmnt, closure *Env, isInit bool) Function {
	return Function{
		declaration: declaration,
		closure:     closure,
		isInit:      isInit,
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

	// Automatic calls of super initializers. Not good code, dunno if I should leave it here...
	// if f.isInit {
	// 	super, err := env.Get(Token{lexeme: "super"})
	// 	if err == nil {
	// 		if init, ok := (super.(*BluClass)).methods["init"]; ok {
	// 			if init.Arity() == 0 {
	// 				this, err := env.Get(Token{lexeme: "this"})
	// 				if err != nil {
	// 					return nil, err
	// 				}

	// 				_, err = init.bind(this.(*BluInstance)).Call(i, []interface{}{})
	// 				if err != nil {
	// 					return nil, err
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	err := i.executeBlock(f.declaration.body, env)
	if returner, ok := err.(Returner); ok {
		if f.isInit {
			return f.closure.GetAt(0, Token{lexeme: "this"})
		}

		return returner.value, nil
	} else if err != nil {
		return nil, err
	}

	if f.isInit {
		return f.closure.GetAt(0, Token{lexeme: "this"})
	}

	return nil, nil
}

func (f Function) Arity() int {
	return len(f.declaration.params)
}

func (f Function) bind(instance *BluInstance) Function {
	env := MakeEnv(f.closure)
	token := MakeToken(This, "this", nil, -1, -1, -1)
	_ = env.Define(token, instance)

	return MakeFunction(f.declaration, env, f.isInit)
}

func (f Function) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
