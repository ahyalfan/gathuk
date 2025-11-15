// Package json
package json

type ASTNode interface {
	Type() string // use feature switch node.(type)
}

type ObjectNode struct {
	Value map[string]ASTNode
}

type ArrayNode struct {
	Value []ASTNode
}

type StringNode struct {
	Value string
}

type NumberNode struct {
	Value float64
}

type BooleanNode struct {
	Value bool
}

type NullNode struct{}

func (o ObjectNode) Type() string {
	return "Object"
}

func (a ArrayNode) Type() string {
	return "Array"
}

func (s StringNode) Type() string {
	return "String"
}

func (n NumberNode) Type() string {
	return "Number"
}

func (b BooleanNode) Type() string {
	return "Boolean"
}

func (n NullNode) Type() string {
	return "Null"
}
