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

	return instance, nil
}

func (c *BluClass) Arity() int {
	return 0
}

func (c *BluClass) String() string {
	return c.name
}

var _ Callable = &BluClass{}
