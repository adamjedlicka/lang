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

func (env *Env) Define(name Token, value interface{}) error {
	if _, ok := env.values[name.lexeme]; !ok {
		env.values[name.lexeme] = value
		return nil
	}

	return NewRuntimeError(
		name.line,
		fmt.Sprintf("Variable '%s' already defined.", name.lexeme))
}

func (env *Env) Assign(name Token, value interface{}) error {
	if _, ok := env.values[name.lexeme]; ok {
		env.values[name.lexeme] = value

		return nil
	}

	if env.enclosing != nil {
		return env.enclosing.Assign(name, value)
	}

	return NewRuntimeError(
		name.line,
		fmt.Sprintf("Cannot assign to undefined variable '%s'.", name.lexeme))
}

func (env *Env) AssignAt(distance int, name Token, value interface{}) error {
	if _, ok := env.ancestor(distance).values[name.lexeme]; ok {
		env.ancestor(distance).values[name.lexeme] = value

		return nil
	}

	return NewRuntimeError(
		name.line,
		fmt.Sprintf("Cannot assign to undefined variable '%s'.", name.lexeme))
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

func (env *Env) GetAt(distance int, name Token) (interface{}, error) {
	if value, ok := env.ancestor(distance).values[name.lexeme]; ok {
		return value, nil
	}

	return nil, NewRuntimeError(
		name.line,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}

func (env *Env) ancestor(distance int) *Env {
	e := env

	for i := 0; i < distance; i++ {
		e = e.enclosing
	}

	return e
}
