package lang

type BluClass struct {
	name         string
	superclass   *BluClass
	declarations map[string]Expr
	methods      map[string]Function
}

func MakeBluClass(name string, superclass *BluClass, declarations map[string]Expr, methods map[string]Function) *BluClass {
	return &BluClass{
		name:         name,
		superclass:   superclass,
		declarations: declarations,
		methods:      methods,
	}
}

func (c *BluClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := MakeBluInstance(c)

	err := c.init(instance, interpreter, arguments)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (c *BluClass) findMethod(name Token) (interface{}, error) {
	if method, ok := c.methods[name.lexeme]; ok {
		return method, nil
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil, NewRuntimeError(name.line, "Undefined property '"+name.lexeme+"'.")
}

func (c *BluClass) init(instance *BluInstance, interpreter *Interpreter, arguments []interface{}) error {
	if c.superclass != nil {
		err := c.superclass.init(instance, interpreter, arguments)
		if err != nil {
			return err
		}
	}

	for key, declaration := range c.declarations {
		if declaration == nil {
			instance.fields[key] = nil
			continue
		}

		value, err := interpreter.evaluate(declaration)
		if err != nil {
			return err
		}

		instance.fields[key] = value
	}

	if init, ok := c.methods["init"]; ok {
		_, err := init.bind(instance).Call(interpreter, arguments)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BluClass) Arity() int {
	if init, ok := c.methods["init"]; ok {
		return init.Arity()
	}

	return 0
}

func (c *BluClass) String() string {
	return c.name
}

var _ Callable = &BluClass{}
