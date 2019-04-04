package lang

type BluClass struct {
	name    string
	methods map[string]Function
}

func MakeBluClass(name string, methods map[string]Function) *BluClass {
	return &BluClass{
		name:    name,
		methods: methods,
	}
}

func (c *BluClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := MakeBluInstance(c)

	if init, ok := c.methods["init"]; ok {
		_, err := init.bind(instance).Call(interpreter, arguments)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
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
