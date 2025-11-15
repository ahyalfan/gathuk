// Package json provides encoding and decoding functionality for JSON format.
package json

// ASTNode is the base interface for all JSON AST nodes.
//
// Each node type implements this interface and provides a Type() method
// that returns a string identifying the node type. This allows for
// type switching when traversing the AST.
//
// Example:
//
//	switch node := astNode.(type) {
//	case ObjectNode:
//	    // Handle object
//	case ArrayNode:
//	    // Handle array
//	}
type ASTNode interface {
	Type() string // returns the type identifier for type switching
}

// ObjectNode represents a JSON object in the AST.
//
// A JSON object is a collection of key-value pairs where keys are strings
// and values can be any JSON type.
//
// Example JSON: {"name": "John", "age": 30}
type ObjectNode struct {
	Value map[string]ASTNode
}

// ArrayNode represents a JSON array in the AST.
//
// A JSON array is an ordered list of values where each value can be
// any JSON type.
//
// Example JSON: [1, "text", true, null]
type ArrayNode struct {
	Value []ASTNode
}

// StringNode represents a JSON string value in the AST.
//
// Example JSON: "hello world"
type StringNode struct {
	Value string
}

// NumberNode represents a JSON number value in the AST.
//
// JSON numbers are stored as float64 to handle both integers and
// floating-point numbers uniformly.
//
// Example JSON: 123, 45.67, 1.23e10
type NumberNode struct {
	Value float64
}

// BooleanNode represents a JSON boolean value in the AST.
//
// Example JSON: true, false
type BooleanNode struct {
	Value bool
}

// NullNode represents a JSON null value in the AST.
//
// Example JSON: null
type NullNode struct{}

// Type returns "Object" for ObjectNode.
func (o ObjectNode) Type() string {
	return "Object"
}

// Type returns "Array" for ArrayNode.
func (a ArrayNode) Type() string {
	return "Array"
}

// Type returns "String" for StringNode.
func (s StringNode) Type() string {
	return "String"
}

// Type returns "Number" for NumberNode.
func (n NumberNode) Type() string {
	return "Number"
}

// Type returns "Boolean" for BooleanNode.
func (b BooleanNode) Type() string {
	return "Boolean"
}

// Type returns "Null" for NullNode.
func (n NullNode) Type() string {
	return "Null"
}
