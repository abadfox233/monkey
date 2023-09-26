package vm

import (
	"fmt"
	"monkey/internal/code"
	"monkey/internal/compiler"
	"monkey/internal/object"
)

const StackSize = 2048

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向下一个空闲的栈位置。 栈顶的值是 stack[sp-1]。
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
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
		case code.OpPop:
			vm.pop()
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
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

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftVal, ok := left.(*object.Integer)
	if !ok {
		return fmt.Errorf("unsupported type for add: %T", left)
	}
	rightVal, ok := right.(*object.Integer)
	if !ok {
		return fmt.Errorf("unsupported type for add: %T", right)
	}
	var result int64
	switch op {
	case code.OpAdd:
		result = leftVal.Value + rightVal.Value
	case code.OpSub:
		result = leftVal.Value - rightVal.Value
	case code.OpMul:
		result = leftVal.Value * rightVal.Value
	case code.OpDiv:
		if rightVal.Value == 0 {
			return fmt.Errorf("division by zero")
		}
		result = leftVal.Value / rightVal.Value
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	vm.push(&object.Integer{Value: result})
	return nil
}
