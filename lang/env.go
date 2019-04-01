package lang

import (
	"fmt"
)

type Env struct {
	values map[string]interface{}
}

func MakeEnv() Env {
	env := Env{}
	env.values = make(map[string]interface{})

	return env
}

func (env *Env) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Env) Assign(name Token, value interface{}) error {
	_, ok := env.values[name.lexeme]
	if !ok {
		return NewRuntimeError(
			name.line,
			fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}

	env.values[name.lexeme] = value

	return nil
}

func (env *Env) Get(name Token) (interface{}, error) {
	value, ok := env.values[name.lexeme]
	if !ok {
		return nil, NewRuntimeError(
			name.line,
			fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}

	return value, nil
}
