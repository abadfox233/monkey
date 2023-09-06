package evaluator

import (
	"monkey/internal/object"
	"os"
	"time"
)

var builtins = map[string]*object.Builtin{
	"len":       {Fn: lenObject},
	"timestamp": {Fn: timestamp},
	"first":     {Fn: first},
	"last":      {Fn: last},
	"rest":      {Fn: rest},
	"push":      {Fn: push},
	"puts":      {Fn: puts},
	"int":       {Fn: toInt},
	"float":     {Fn: toFloat},
	"str":       {Fn: toStr},
	"exit":      {Fn: exit},
}

func lenObject(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
		return NULL
	default:
		return newError("argument to `first` not supported, got %s", args[0].Type())
	}
}

func last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[len(arg.Elements)-1]
		}
		return NULL
	default:
		return newError("argument to `last` not supported, got %s", args[0].Type())
	}
}

func rest(args ...object.Object) object.Object {

	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)
		if length > 0 {
			newElements := make([]object.Object, length-1)
			copy(newElements, arg.Elements[1:length])
			return &object.Array{Elements: newElements}
		}
		return NULL
	default:
		return newError("argument to `rest` not supported, got %s", args[0].Type())
	}

}

func push(args ...object.Object) object.Object{
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]
	return &object.Array{Elements: newElements}
}

func timestamp(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("wrong number of arguments. got=%d, want=0", len(args))
	}
	return &object.Integer{Value: time.Now().Unix()}
}

func puts(args ...object.Object) object.Object {
	for _, arg := range args {
		println(arg.Inspect())
	}
	return NULL
}

func toInt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if number, ok := args[0].(object.Number) ; ok {
		return &object.Integer{Value: number.Integer()}
	}
	return newError("argument to `int` must be NUMBER, got %s", args[0].Type())
}

func toFloat(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if number, ok := args[0].(object.Number) ; ok {
		return &object.Float{Value: number.Float()}
	}
	return newError("argument to `float` must be NUMBER, got %s", args[0].Type())
}

func toStr(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	return &object.String{Value: args[0].Inspect()}
}

func exit(args ...object.Object) object.Object {
	os.Exit(0)
	return NULL
}