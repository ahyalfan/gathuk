// Package json
package json

import (
	"errors"
	"fmt"
	"strconv"
)

func Parser(tokens []Token) (ASTNode, error) {
	if len(tokens) == 0 {
		return nil, errors.New("nothing to parse")
	}
	current := 0
	return parseValue(&current, tokens)
}

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
