// Package json provides encoding and decoding functionality for JSON format.
package json

// Token type constants representing different JSON token types.
//
// These constants are used during the tokenization phase to identify
// the type of each token in the JSON input.
const (
	BraceOpen int = iota
	BraceClose
	BracketOpen
	BracketClose
	Comma
	Colon
	String
	Number
	True
	False
	Null
)

// Token represents a single lexical token from JSON input.
//
// Tokens are the output of the tokenization phase and input to the parser.
// Each token carries information about its type, value, and position for
// error reporting.
//
// Example:
//
//	token := Token{
//	    Type:   String,
//	    Value:  []byte("hello"),
//	    Line:   1,
//	    Column: 5,
//	}
type Token struct {
	Type   int
	Value  []byte
	Line   int // line number for error reporting
	Column int // column number for error reporting
}
