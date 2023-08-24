package evaluator

import (
	"monkey/internal/ast"
	"monkey/internal/object"
)

func quote(node ast.Node) object.Object {

	return &object.Quote{Node: node}
}
