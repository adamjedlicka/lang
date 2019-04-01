package lang

import (
	"fmt"
)

type Env struct {
	enclosing *Env
	values    map[string]interface{}
}

func MakeEnv(enclosing *Env) *Env {
	env := new(Env)
	env.enclosing = enclosing
	env.values = make(map[string]interface{})

	return env
}

func (env *Env) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Env) Assign(name Token, value interface{}) error {
	if _, ok := env.values[name.lexeme]; ok {
		env.values[name.lexeme] = value

		return nil
	}

	return NewRuntimeError(
		name.line,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}

func (env *Env) Get(name Token) (interface{}, error) {
	if value, ok := env.values[name.lexeme]; ok {
		return value, nil
	}

	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}

	return nil, NewRuntimeError(
		name.line,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}
