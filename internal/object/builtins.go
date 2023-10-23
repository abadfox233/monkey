package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}

				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					if len(arg.Elements) > 0 {
						return arg.Elements[0]
					}
					return nil
				case *String:
					if len(arg.Value) > 0 {
						return &String{Value: string([]rune(arg.Value)[0])}
					}
					return nil
				default:
					return newError("argument to `first` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					length := len(arg.Elements)
					if length > 0 {
						return arg.Elements[length-1]
					}
					return nil
				case *String:
					length := len(arg.Value)
					if length > 0 {
						return &String{Value: string([]rune(arg.Value)[length-1])}
					}
					return nil
				default:
					return newError("argument to `last` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		Name: "rest",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					length := len(arg.Elements)
					if length > 0 {
						newElements := make([]Object, length-1)
						copy(newElements, arg.Elements[1:length])
						return &Array{Elements: newElements}
					}
					return nil
				case *String:
					length := len(arg.Value)
					if length > 0 {
						return &String{Value: string([]rune(arg.Value)[1:length])}
					}
					return nil
				default:
					return newError("argument to `rest` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		Name: "push",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					length := len(arg.Elements)
					newElements := make([]Object, length+1)
					copy(newElements, arg.Elements)
					newElements[length] = args[1]
					return &Array{Elements: newElements}
				case *String:
					return &String{Value: arg.Value + args[1].Inspect()}
				default:
					return newError("argument to `push` not supported, got %s", args[0].Type())
				}
			},
		},
	},
}


func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, builtin := range Builtins {
		if builtin.Name == name {
			return builtin.Builtin
		}
	}
	return nil
}