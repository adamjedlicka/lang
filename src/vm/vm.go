package vm

import (
	"fmt"

	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/compiler"
	"github.com/adamjedlicka/lang/src/config"
	"github.com/adamjedlicka/lang/src/debug"
	"github.com/adamjedlicka/lang/src/val"
)

type VM struct {
	chunk *code.Chunk
	ip    int
	stack []val.Value
}

func NewVM() *VM {
	vm := new(VM)
	vm.ip = 0
	vm.stack = make([]val.Value, 0)

	return vm
}

func (vm *VM) Interpret(source []rune) {
	vm.chunk = compiler.NewCompiler(source).Compile()
	vm.ip = 0

	// vm.run()
}

func (vm *VM) run() {
	for {
		if config.FlagDebug {
			if config.FlagStack {
				fmt.Print("STACK :: [")
				for i, value := range vm.stack {
					if i > 0 {
						fmt.Print(", ")
					}

					fmt.Printf("%s", value)
				}
				fmt.Println("]")
			}
			debug.DisassembleInstruction(vm.chunk, vm.ip)
		}

		instruction := vm.readInstruction()

		switch instruction {
		case code.OpConstant:
			constant := vm.readConstant()
			vm.push(constant)
		case code.OpAdd:
			vm.push(vm.pop().Add(vm.pop()))
		case code.OpSubtract:
			vm.push(vm.pop().Subtract(vm.pop()))
		case code.OpMultiply:
			vm.push(vm.pop().Multiply(vm.pop()))
		case code.OpDivide:
			vm.push(vm.pop().Divide(vm.pop()))
		case code.OpNegate:
			vm.push(vm.pop().Negate())
		case code.OpReturn:
			return
		}
	}
}

func (vm *VM) push(value val.Value) {
	vm.stack = append(vm.stack, value)
}

func (vm *VM) pop() val.Value {
	value := vm.stack[len(vm.stack)-1]

	vm.stack = vm.stack[:len(vm.stack)-1]

	return value
}

func (vm *VM) readInstruction() code.OpCode {
	instruction := vm.chunk.Get(vm.ip)

	vm.ip++

	return instruction
}

func (vm *VM) readConstant() val.Value {
	return vm.chunk.GetConstant(uint8(vm.readInstruction()))
}
