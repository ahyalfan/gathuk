// Package json provides encoding and decoding functionality for JSON format.
package json

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

// Tokenize converts JSON bytes into a sequence of tokens.
//
// This is the first phase of JSON parsing (lexical analysis). It scans the
// input and produces a flat list of tokens that represent the structure and
// values in the JSON.
//
// Supported tokens:
//   - Structural: { } [ ] : ,
//   - Literals: "string", 123, true, false, null
//
// The tokenizer handles:
//   - Whitespace (spaces, tabs, newlines) - ignored
//   - String literals with proper quote handling
//   - Numbers (integers, floats, scientific notation)
//   - Boolean literals (true, false)
//   - Null literal
//
// Parameters:
//   - input: JSON bytes to tokenize
//
// Returns:
//   - []Token: Sequence of tokens
//   - error: An error if tokenization fails (e.g., unterminated string, invalid number)
//
// Example:
//
//	input := []byte(`{"name": "John", "age": 30}`)
//	tokens, err := Tokenize(input)
//	// tokens: [{BraceOpen}, {String "name"}, {Colon}, {String "John"}, ...]
func Tokenize(input []byte) ([]Token, error) {
	var (
		current     = 0
		tokens      []Token
		char        byte
		inputLength = len(input)
	)

	for current < len(input) {
		char = input[current]

		if unicode.IsSpace(rune(char)) {
			current++
			continue
		}

		switch char {
		case '{':
			v := Token{Type: BraceOpen, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case '}':
			v := Token{Type: BraceClose, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case '[':
			v := Token{Type: BracketOpen, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case ']':
			v := Token{Type: BracketClose, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case ':':
			v := Token{Type: Colon, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case ',':
			v := Token{Type: Comma, Value: []byte{char}}
			tokens = append(tokens, v)
			current++
		case '"':
			current++
			start := current

			for current < inputLength && input[current] != '"' {
				current++
			}

			temp := input[start:current]

			if current >= inputLength {
				return nil, fmt.Errorf("untermineted string: %s", string(temp))
			}

			v := Token{Type: String, Value: temp}
			tokens = append(tokens, v)

			current++
		default:
			rest := input[current : current+8] // get prefix

			if bytes.HasPrefix(rest, []byte("true")) {
				v := Token{Type: True, Value: []byte("true")}
				tokens = append(tokens, v)
				current += 4
			} else if bytes.HasPrefix(rest, []byte("false")) {
				v := Token{Type: False, Value: []byte("false")}
				tokens = append(tokens, v)
				current += 5
			} else if bytes.HasPrefix(rest, []byte("null")) {
				v := Token{Type: Null, Value: []byte("null")}
				tokens = append(tokens, v)
				current += 4
			} else if unicode.IsNumber(rune(char)) || char == '-' {
				start := current
				current++

				var (
					hasDot, hasExp bool
					expDigits      int
				)

				for current < inputLength {
					te := input[current]

					if unicode.IsDigit(rune(te)) {
						current++
						if hasExp {
							expDigits++
						}
					} else if te == '.' {
						if hasDot || hasExp {
							return nil, fmt.Errorf("invalid number: multiple dots or dot after exponent at position %d", current)
						}
						hasDot = true
						current++
					} else if te == 'e' || te == 'E' {
						if hasExp {
							return nil, fmt.Errorf("invalid number: multiple exponents at position %d", current)
						}
						hasExp = true
						current++

						if current < inputLength && (input[current] == '+' || input[current] == '-') {
							current++
						}

						expDigits = 0
					} else {
						break
					}

				}

				num := input[start:current]

				if hasExp && expDigits == 0 {
					return nil, fmt.Errorf("invalid number: exponent missing digits in '%s'", num)
				}
				if _, err := strconv.ParseFloat(string(num), 64); err != nil {
					return nil, fmt.Errorf("invalid number: %s", string(num))
				}
				tokens = append(tokens, Token{Type: Number, Value: num})
			} else {
				return nil, fmt.Errorf("unexpected character: %c, position: %d", char, current)
			}

		}
	}

	return tokens, nil
}
