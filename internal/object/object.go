package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/internal/ast"
	"monkey/internal/code"
	"strings"
)

type BuiltinFunction func(args ...Object) Object

type ObjectType string

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	STRING_OBJ            = "STRING"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	HASH_OBJ              = "HASH"
	QUOTE_OBJ             = "QUOTE"
	FLOAT_OBJ             = "FLOAT"
	BREAK_OBJ             = "BREAK"
	CONTINUE_OBJ          = "CONTINUE"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
)

type Hashable interface {
	HashKey() HashKey
}

type Number interface {
	Float() float64
	Integer() int64
}

var _ Object = (*Integer)(nil)
var _ Object = (*Boolean)(nil)
var _ Object = (*Null)(nil)
var _ Object = (*ReturnValue)(nil)
var _ Object = (*Error)(nil)
var _ Object = (*Function)(nil)
var _ Object = (*String)(nil)
var _ Object = (*Builtin)(nil)
var _ Object = (*Array)(nil)
var _ Object = (*Float)(nil)
var _ Object = (*CompiledFunction)(nil)

var _ Hashable = (*String)(nil)
var _ Hashable = (*Boolean)(nil)
var _ Hashable = (*Integer)(nil)
var _ Hashable = (*Float)(nil)

var _ Number = (*Integer)(nil)
var _ Number = (*Float)(nil)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.15f", f.Value)))
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}
func (f *Float) Float() float64 { return f.Value }
func (f *Float) Integer() int64 { return int64(f.Value) }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (i *Integer) Float() float64 { return float64(i.Value) }
func (i *Integer) Integer() int64 { return i.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Array struct {
	Elements []Object
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Break struct{}

func (br *Break) Type() ObjectType { return BREAK_OBJ }
func (br *Break) Inspect() string  { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
	return "QUOTE(" + q.Node.String() + ")"
}

type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
