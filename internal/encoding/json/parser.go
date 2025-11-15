// Package json provides encoding and decoding functionality for JSON format.
package json

import (
	"errors"
	"fmt"
	"strconv"
)

// Parser converts a sequence of tokens into an Abstract Syntax Tree (AST).
//
// This is the second phase of JSON parsing (syntax analysis). It takes the
// flat list of tokens produced by the tokenizer and builds a hierarchical
// tree structure that represents the JSON document.
//
// The parser uses recursive descent to build the AST:
//   - parseValue: Entry point that dispatches to specific parsers
//   - parseObject: Handles JSON objects
//   - parseArray: Handles JSON arrays
//   - Primitive values: Handled directly in parseValue
//
// Parameters:
//   - tokens: Sequence of tokens from the tokenizer
//
// Returns:
//   - ASTNode: Root node of the AST
//   - error: An error if parsing fails (e.g., syntax error, unexpected token)
//
// Example:
//
//	tokens := []Token{
//	    {Type: BraceOpen},
//	    {Type: String, Value: []byte("name")},
//	    {Type: Colon},
//	    {Type: String, Value: []byte("John")},
//	    {Type: BraceClose},
//	}
//	ast, err := Parser(tokens)
//	// ast: ObjectNode{Value: {"name": StringNode{"John"}}}
func Parser(tokens []Token) (ASTNode, error) {
	if len(tokens) == 0 {
		return nil, errors.New("nothing to parse")
	}
	current := 0
	return parseValue(&current, tokens)
}

// parseValue parses any JSON value from the token stream.
//
// This is the main dispatcher that determines the type of value based on
// the current token and calls the appropriate parsing function.
//
// Parameters:
//   - current: Pointer to current position in token stream
//   - tokens: Complete token stream
//
// Returns:
//   - ASTNode: The parsed value as an AST node
//   - error: An error if parsing fails
func parseValue(current *int, tokens []Token) (ASTNode, error) {
	if *current >= len(tokens) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	token := tokens[*current]

	switch token.Type {
	case String:
		*current++
		return StringNode{Value: string(token.Value)}, nil
	case Number:
		num, _ := strconv.ParseFloat(string(token.Value), 64)
		*current++
		return NumberNode{Value: num}, nil
	case True:
		*current++
		return BooleanNode{Value: true}, nil
	case False:
		*current++
		return BooleanNode{Value: false}, nil
	case Null:
		*current++
		return NullNode{}, nil
	case BraceOpen:
		return parseObject(current, tokens)
	case BracketOpen:
		return parseArray(current, tokens)
	default:
		return nil, fmt.Errorf("invalid token type: %v", token.Type)
	}
}

// parseObject parses a JSON object from the token stream.
//
// Expected token sequence: { "key" : value , "key" : value ... }
//
// Parameters:
//   - current: Pointer to current position in token stream
//   - tokens: Complete token stream
//
// Returns:
//   - ASTNode: ObjectNode containing the parsed key-value pairs
//   - error: An error if parsing fails
func parseObject(current *int, tokens []Token) (ASTNode, error) {
	node := ObjectNode{
		Value: make(map[string]ASTNode),
	}

	*current++

	for *current < len(tokens) && tokens[*current].Type != BraceClose {
		currToken := tokens[*current]

		if currToken.Type != String {
			return nil, fmt.Errorf("expected string key in object, got: %v", currToken.Type)
		}

		key := currToken.Value
		*current++

		if *current >= len(tokens) || tokens[*current].Type != Colon {
			return nil, fmt.Errorf("expected : in key value pair, got: %v", currToken.Type)
		}

		*current++

		value, err := parseValue(current, tokens)
		if err != nil {
			return nil, err
		}
		node.Value[string(key)] = value

		if *current < len(tokens) && tokens[*current].Type == Comma {
			*current++
		}
	}

	if *current < len(tokens) && tokens[*current].Type != BraceClose {
		return nil, fmt.Errorf("expected closing brace, got: %v", tokens[*current].Type)
	}

	*current++

	return node, nil
}

// parseArray parses a JSON array from the token stream.
//
// Expected token sequence: [ value , value , value ... ]
//
// Parameters:
//   - current: Pointer to current position in token stream
//   - tokens: Complete token stream
//
// Returns:
//   - ASTNode: ArrayNode containing the parsed elements
//   - error: An error if parsing fails
func parseArray(current *int, tokens []Token) (ASTNode, error) {
	node := ArrayNode{
		Value: make([]ASTNode, 0),
	}

	*current++
	for *current < len(tokens) && tokens[*current].Type != BracketClose {
		val, err := parseValue(current, tokens)
		if err != nil {
			return nil, err
		}
		node.Value = append(node.Value, val)
		if *current < len(tokens) && tokens[*current].Type == Comma {
			*current++
		}
	}

	if *current < len(tokens) && tokens[*current].Type != BracketClose {
		return nil, fmt.Errorf("expected closing bracket, got: %v", tokens[*current].Type)
	}
	*current++
	return node, nil
}
