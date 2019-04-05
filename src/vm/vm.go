package vm

import (
	"fmt"
	"time"

	"github.com/adamjedlicka/lang/src/code"
	"github.com/adamjedlicka/lang/src/debug"
	"github.com/adamjedlicka/lang/src/val"
)

type VM struct {
	chunk    *code.Chunk
	ip       int
	stack    []val.Value
	stackLen int
}

type InterpretResult uint8

const (
	InterpretOK InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func NewVM() *VM {
	vm := new(VM)
	vm.ip = 0
	vm.stack = make([]val.Value, 0)
	vm.stackLen = 0

	return vm
}

func (vm *VM) Interpret(chunk *code.Chunk) InterpretResult {
	vm.chunk = chunk
	vm.ip = 0

	start := time.Now().UnixNano()

	ret := vm.run()

	end := time.Now().UnixNano()
	fmt.Printf("runtime: %dus\n", (end-start)/1000)

	return ret
}

func (vm *VM) run() InterpretResult {
	for {
		fmt.Print("          ")
		for _, value := range vm.stack {
			fmt.Printf("[%s]", value.String())
		}
		fmt.Println()
		debug.DisassembleInstruction(vm.chunk, vm.ip)

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
			fmt.Println(vm.pop())
			return InterpretOK
		}
	}
}

func (vm *VM) push(value val.Value) {
	if vm.stackLen == len(vm.stack) {
		vm.stack = append(vm.stack, value)
	} else {
		vm.stack[vm.stackLen] = value
	}

	vm.stackLen++
}

func (vm *VM) pop() val.Value {
	vm.stackLen--

	return vm.stack[vm.stackLen]
}

func (vm *VM) readInstruction() code.OpCode {
	instruction := vm.chunk.Get(vm.ip)

	vm.ip++

	return instruction
}

func (vm *VM) readConstant() val.Value {
	return vm.chunk.GetConstant(uint8(vm.readInstruction()))
}
