package vm

import (
	"fmt"
	"monkey/internal/code"
	"monkey/internal/compiler"
	"monkey/internal/object"
)


const StackSize = 2048

type VM struct {
	constants []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向下一个空闲的栈位置。 栈顶的值是 stack[sp-1]。
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants: bytecode.Constants,
		instructions: bytecode.Instructions,
		stack: make([]object.Object, StackSize),
		sp: 0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {

	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftVal, ok := left.(*object.Integer)
			if !ok {
				return fmt.Errorf("unsupported type for add: %T", left)
			}
			rightVal, ok := right.(*object.Integer)
			if !ok {
				return fmt.Errorf("unsupported type for add: %T", right)
			}
			result := leftVal.Value + rightVal.Value
			vm.push(&object.Integer{Value: result})
		}
	}
	return nil

}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++
	return nil

}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj
}