package evaluator

import (
	"monkey/internal/object"
	"os"
	"time"
)

var builtins = map[string]*object.Builtin{
	"len":       object.GetBuiltinByName("len"),
	"timestamp": {Fn: timestamp},
	"first":     object.GetBuiltinByName("first"),
	"last":      object.GetBuiltinByName("last"),
	"rest":      object.GetBuiltinByName("rest"),
	"push":      object.GetBuiltinByName("push"),
	"puts":      object.GetBuiltinByName("puts"),
	"int":       {Fn: toInt},
	"float":     {Fn: toFloat},
	"str":       {Fn: toStr},
	"exit":      {Fn: exit},
}

func timestamp(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("wrong number of arguments. got=%d, want=0", len(args))
	}
	return &object.Integer{Value: time.Now().Unix()}
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