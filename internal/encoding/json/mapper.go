// Package json provides encoding and decoding functionality for JSON format.
package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	utility "github.com/ahyalfan/gathuk/internal/utils"
	"github.com/ahyalfan/gathuk/shared"
)

// StructToAST converts a Go struct to an AST node.
//
// This method uses reflection to traverse the struct and build an
// equivalent AST representation. It handles:
//   - Nested structs → ObjectNode
//   - Slices/arrays → ArrayNode
//   - Maps → ObjectNode
//   - Primitive types → StringNode, NumberNode, BooleanNode
//   - Struct tags for custom field names
//
// Parameters:
//   - value: Pointer to the struct to convert
//
// Returns:
//   - ASTNode: The root AST node representing the struct
//   - error: An error if conversion fails
//
// Example:
//
//	type Config struct {
//	    Port int    `config:"port"`
//	    Host string `config:"host"`
//	}
//
//	config := &Config{Port: 8080, Host: "localhost"}
//	ast, err := codec.StructToAST(config)
//	// ast: ObjectNode{Value: {"port": NumberNode{8080}, "host": StringNode{"localhost"}}}
func (c *Codec[T]) StructToAST(value *T) (ASTNode, error) {
	if value == nil {
		return NullNode{}, nil
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, fmt.Errorf("value must be non-nil pointer")
	}
	return c.valueToNode(rv.Elem(), "")
}

// valueToNode converts a reflect.Value to an AST node.
//
// This is a helper method used during struct-to-AST conversion.
// It handles the conversion of Go values to their AST representations.
//
// Parameters:
//   - v: The reflect.Value to convert
//   - path: Current path in the struct (for error reporting)
//
// Returns:
//   - ASTNode: The converted AST node
//   - error: An error if conversion fails
func (c *Codec[T]) valueToNode(v reflect.Value, path string) (ASTNode, error) {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return NullNode{}, nil
		}
		return c.valueToNode(v.Elem(), path)
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NullNode{}, nil
		}
		return c.valueToNode(v.Elem(), path)
	}

	switch v.Kind() {
	case reflect.Struct:
		return c.structToNode(v, path)

	case reflect.Slice, reflect.Array:
		return c.sliceToNode(v, path)

	case reflect.Map:
		return c.mapToNode(v, path)

	case reflect.String:
		return StringNode{Value: v.String()}, nil

	case reflect.Bool:
		return BooleanNode{Value: v.Bool()}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NumberNode{Value: float64(v.Int())}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NumberNode{Value: float64(v.Uint())}, nil

	case reflect.Float32, reflect.Float64:
		return NumberNode{Value: v.Float()}, nil

	default:
		return nil, fmt.Errorf("unsupported type at %s: %s", path, v.Kind())
	}
}

func (c *Codec[T]) structToNode(v reflect.Value, path string) (ASTNode, error) {
	obj := make(map[string]ASTNode)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get(string(shared.GetTagName()))
		if jsonTag == "-" {
			continue
		}
		if jsonTag == "" {
			jsonTag = field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
		}

		name := jsonTag
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			name = jsonTag[:idx]
		}
		if name == "" {
			name = utility.PascalToLowerSnakeCase(field.Name)
		}

		fieldPath := path + "." + name
		if path == "" {
			fieldPath = name
		}

		node, err := c.valueToNode(v.Field(i), fieldPath)
		if err != nil {
			return nil, err
		}
		obj[name] = node
	}

	return ObjectNode{Value: obj}, nil
}

func (c *Codec[T]) sliceToNode(v reflect.Value, path string) (ASTNode, error) {
	nodes := make([]ASTNode, v.Len())
	for i := 0; i < v.Len(); i++ {
		elemPath := fmt.Sprintf("%s[%d]", path, i)
		if path == "" {
			elemPath = fmt.Sprintf("[%d]", i)
		}
		node, err := c.valueToNode(v.Index(i), elemPath)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return ArrayNode{Value: nodes}, nil
}

func (c *Codec[T]) mapToNode(v reflect.Value, path string) (ASTNode, error) {
	if v.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("map key must be string at %s, got %s", path, v.Type().Key().Kind())
	}

	obj := make(map[string]ASTNode)
	for _, key := range v.MapKeys() {
		keyStr := key.String()
		val := v.MapIndex(key)

		elemPath := path + "." + keyStr
		if path == "" {
			elemPath = keyStr
		}

		node, err := c.valueToNode(val, elemPath)
		if err != nil {
			return nil, err
		}
		obj[keyStr] = node
	}
	return ObjectNode{Value: obj}, nil
}

// ASTToStruct converts an AST node to a Go struct.
//
// This method uses reflection to populate a struct from an AST. It handles:
//   - ObjectNode → struct
//   - ArrayNode → slice
//   - ObjectNode → map (if target is map)
//   - Primitive nodes → basic Go types
//   - Type conversions (string → int, etc.)
//   - Struct tags for field mapping
//
// Parameters:
//   - node: The AST node to convert
//   - v: Pointer to the destination struct
//
// Returns:
//   - error: An error if conversion fails
//
// Example:
//
//	ast := ObjectNode{
//	    Value: map[string]ASTNode{
//	        "port": NumberNode{8080},
//	        "host": StringNode{"localhost"},
//	    },
//	}
//
//	var config Config
//	err := codec.ASTToStruct(ast, &config)
//	// config: Config{Port: 8080, Host: "localhost"}
func (c *Codec[T]) ASTToStruct(node ASTNode, v *T) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}
	return c.nodeToValue(node, rv.Elem(), "")
}

// nodeToValue converts an AST node to a reflect.Value.
//
// This is a helper method used during AST-to-struct conversion.
// It handles the conversion of AST nodes to Go values with proper
// type checking and conversion.
//
// Parameters:
//   - node: The AST node to convert
//   - v: The target reflect.Value
//   - path: Current path in the struct (for error reporting)
//
// Returns:
//   - error: An error if conversion fails
func (c *Codec[T]) nodeToValue(node ASTNode, v reflect.Value, path string) error {
	if !v.CanSet() {
		return c.newError(path, "value not settable")
	}

	// Handle interface{} / any
	if v.Kind() == reflect.Interface {
		native, err := c.toNative(node, path)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(native))
		return nil
	}

	switch node := node.(type) {
	case ObjectNode:
		switch v.Kind() {
		case reflect.Struct:
			return c.mapObject(node, v, path)
		case reflect.Map:
			return c.mapToMap(node, v, path)
		default:
			return c.newError(path, "expected struct or map, got %s", v.Kind())
		}

	case ArrayNode:
		if v.Kind() != reflect.Slice {
			return c.newError(path, "expected slice, got %s", v.Kind())
		}
		return c.mapArray(node, v, path)

	case StringNode:
		return c.stringValue(node.Value, v, path)

	case NumberNode:
		return c.numberValue(node.Value, v, path)

	case BooleanNode:
		if v.Kind() == reflect.Bool {
			v.SetBool(node.Value)
			return nil
		}
		return c.newError(path, "cannot unmarshal boolean into %s", v.Kind())

	case NullNode:
		v.Set(reflect.Zero(v.Type()))
		return nil

	default:
		return c.newError(path, "unsupported node type: %T", node)
	}
}

func (c *Codec[T]) mapObject(node ObjectNode, v reflect.Value, path string) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get(string(shared.GetTagName()))
		if jsonTag == "-" {
			continue
		}
		if jsonTag == "" {
			jsonTag = field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
		}

		name := jsonTag
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			name = jsonTag[:idx]
		}
		if name == "" {
			name = utility.PascalToLowerSnakeCase(field.Name)
		}

		fieldPath := path + "." + name
		if path == "" {
			fieldPath = name
		}

		if childNode, ok := node.Value[name]; ok {
			fieldVal := v.Field(i)
			if err := c.nodeToValue(childNode, fieldVal, fieldPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Codec[T]) mapToMap(node ObjectNode, v reflect.Value, path string) error {
	if v.Type().Key().Kind() != reflect.String {
		return c.newError(path, "map key must be string, got %s", v.Type().Key())
	}

	mapType := v.Type()
	newMap := reflect.MakeMap(mapType)

	for key, childNode := range node.Value {
		elemPath := path + "." + key
		if path == "" {
			elemPath = key
		}

		elemValue := reflect.New(mapType.Elem()).Elem()
		if err := c.nodeToValue(childNode, elemValue, elemPath); err != nil {
			return err
		}

		newMap.SetMapIndex(reflect.ValueOf(key), elemValue)
	}

	v.Set(newMap)
	return nil
}

func (c *Codec[T]) mapArray(node ArrayNode, v reflect.Value, path string) error {
	elemType := v.Type().Elem()
	newSlice := reflect.MakeSlice(reflect.SliceOf(elemType), len(node.Value), len(node.Value))

	for i, item := range node.Value {
		elemPath := fmt.Sprintf("%s[%d]", path, i)
		if path == "" {
			elemPath = fmt.Sprintf("[%d]", i)
		}
		elem := newSlice.Index(i)
		if err := c.nodeToValue(item, elem, elemPath); err != nil {
			return err
		}
	}

	v.Set(newSlice)
	return nil
}

func (c Codec[T]) stringValue(s string, v reflect.Value, path string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			if v.OverflowInt(i) {
				return c.newError(path, "string %q overflows %s", s, v.Type())
			}
			v.SetInt(i)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(s, 10, 64); err == nil {
			if v.OverflowUint(u) {
				return c.newError(path, "string %q overflows %s", s, v.Type())
			}
			v.SetUint(u)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			if v.OverflowFloat(f) {
				return c.newError(path, "string %q overflows %s", s, v.Type())
			}
			v.SetFloat(f)
			return nil
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err == nil {
			v.SetBool(b)
			return nil
		}
	}
	return c.newError(path, "cannot unmarshal string %q into %s", s, v.Type())
}

func (c *Codec[T]) numberValue(f float64, v reflect.Value, path string) error {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		if v.OverflowFloat(f) {
			return c.newError(path, "number %g overflows %s", f, v.Type())
		}
		v.SetFloat(f)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := int64(f)
		if float64(i) != f {
			return c.newError(path, "number %g cannot be converted to integer", f)
		}
		if v.OverflowInt(i) {
			return c.newError(path, "number %g overflows %s", f, v.Type())
		}
		v.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if f < 0 {
			return c.newError(path, "negative number %g cannot be assigned to unsigned type", f)
		}
		u := uint64(f)
		if float64(u) != f {
			return c.newError(path, "number %g cannot be converted to unsigned integer", f)
		}
		if v.OverflowUint(u) {
			return c.newError(path, "number %g overflows %s", f, v.Type())
		}
		v.SetUint(u)
		return nil
	}
	return c.newError(path, "cannot unmarshal number %g into %s", f, v.Type())
}

// toNative converts an AST node to native Go types for interface{}.
//
// This method is used when the target type is interface{} or any.
// It converts AST nodes to appropriate Go types:
//   - StringNode → string
//   - NumberNode → float64
//   - BooleanNode → bool
//   - NullNode → nil
//   - ArrayNode → []interface{}
//   - ObjectNode → map[string]interface{}
//
// Parameters:
//   - node: The AST node to convert
//   - path: Current path (for error reporting)
//
// Returns:
//   - interface{}: The converted native Go value
//   - error: An error if conversion fails
func (c *Codec[T]) toNative(node ASTNode, path string) (interface{}, error) {
	switch n := node.(type) {
	case StringNode:
		return n.Value, nil
	case NumberNode:
		return n.Value, nil
	case BooleanNode:
		return n.Value, nil
	case NullNode:
		return nil, nil
	case ArrayNode:
		slice := make([]interface{}, len(n.Value))
		for i, item := range n.Value {
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			if path == "" {
				elemPath = fmt.Sprintf("[%d]", i)
			}
			val, err := c.toNative(item, elemPath)
			if err != nil {
				return nil, err
			}
			slice[i] = val
		}
		return slice, nil
	case ObjectNode:
		m := make(map[string]interface{})
		for k, v := range n.Value {
			val, err := c.toNative(v, path+"."+k)
			if err != nil {
				return nil, err
			}
			m[k] = val
		}
		return m, nil
	default:
		return nil, c.newError(path, "unsupported node type for interface{}: %T", node)
	}
}

// newError creates a formatted error with path information.
//
// This helper method is used throughout the mapper to create
// informative error messages that include the path to the problematic
// field.
//
// Parameters:
//   - path: Path to the field where error occurred
//   - format: Error message format string
//   - args: Format arguments
//
// Returns:
//   - error: Formatted error with path context
func (c *Codec[T]) newError(path, format string, args ...interface{}) error {
	if path != "" {
		return fmt.Errorf("ast unmarshal error at %s: %w", path, fmt.Errorf(format, args...))
	}
	return fmt.Errorf("ast unmarshal error: "+format, args...)
}
