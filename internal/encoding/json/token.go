// Package json
package json

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

type Token struct {
	Type   int
	Value  []byte
	Line   int // line number for error reporting
	Column int // column number for error reporting
}
