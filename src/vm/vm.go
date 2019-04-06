package vm

import (
	"fmt"
	"time"

	"github.com/adamjedlicka/lang/src/code"
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

func (vm *VM) Interpret(chunk *code.Chunk) {
	vm.chunk = chunk

	start := time.Now().UnixNano()
	out := vm.run()
	end := time.Now().UnixNano()

	fmt.Println(out.String())
	fmt.Printf("time: %dns\n", end-start)
}

func (vm *VM) run() val.Value {
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
			right := vm.pop()
			left := vm.pop()
			vm.push(left.Add(right))
		case code.OpSubtract:
			right := vm.pop()
			left := vm.pop()
			vm.push(left.Subtract(right))
		case code.OpMultiply:
			right := vm.pop()
			left := vm.pop()
			vm.push(left.Multiply(right))
		case code.OpDivide:
			right := vm.pop()
			left := vm.pop()
			vm.push(left.Divide(right))
		case code.OpNegate:
			vm.push(vm.pop().Negate())
		case code.OpReturn:
			return vm.pop()
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
