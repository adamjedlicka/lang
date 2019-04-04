package lang

type BluInstance struct {
	class  *BluClass
	fields map[string]interface{}
}

func MakeBluInstance(class *BluClass) *BluInstance {
	return &BluInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (i *BluInstance) get(name Token) (interface{}, error) {
	if value, ok := i.fields[name.lexeme]; ok {
		return value, nil
	}

	if method, ok := i.class.methods[name.lexeme]; ok {
		return method.bind(i), nil
	}

	return nil, NewRuntimeError(name.line, "Undefined property '"+name.lexeme+"'.")
}

func (i *BluInstance) set(name Token, value interface{}) {
	i.fields[name.lexeme] = value
}

func (i *BluInstance) String() string {
	return i.class.String() + " instance"
}
