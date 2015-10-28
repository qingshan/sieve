package parse

import (
	"fmt"
)

// Pos represents a byte position in the original input text
type Pos int

const NoPos Pos = 0

func (p Pos) Position() Pos {
	return p
}

// token represents a token or text string returned from the scanner.
type Token struct {
	Typ TokenType // The type of this token.
	Pos Pos       // The starting position, in bytes, of this item in the input string.
	Val string    // The value of this item.
}

func (t Token) String() string {
	switch {
	case t.Typ == EOF:
		return "EOF"
	case t.Typ == ERROR:
		return t.Val
	case t.Typ > COMMAND:
		return fmt.Sprintf("<%s>", t.Val)
	case len(t.Val) > 10:
		return fmt.Sprintf("%.10q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}

type TokenType int

const (
	ERROR TokenType = iota // error occurred; value is text of error
	EOF
	LINECOMMENT  // // ..... includes symbol
	BLOCKCOMMENT // /* block comment includes surrounding symbols*/
	LEFTPAREN    // '('
	RIGHTPAREN // ')'
	LEFTBRACKET    // '['
	RIGHTBRACKET // ']'
	LEFTCURLY    // '{'
	RIGHTCURLY // '}'
	SEMICOLON  // ';'
	COMMA      // ','
	ARGUMENT   // used only to delimit the arguments
	NUMBER     // simple number, including imaginary
	TAG        // an int
	STRING     // a string literal
	COMMAND    // used only to delimit the keywords
	IF      // if keyword
	ELSIF    // else keyword
	ELSE    // else keyword
	STOP    // stop keyword
	TEST      // include keyword
	TRUE      // true keyword
	FALSE      // false keyword
	NOT      // not keyword
	ANYOF      // anyof keyword
	ALLOF      // allof keyword
	IDENTIFIER // alphanumeric identifier
)

const eof = -1

var key = map[string]TokenType{
	"if":   IF,
	"else": ELSE,
	"elsif": ELSIF,
	"stop": STOP,
	"true": TRUE,
	"false": FALSE,
	"not": NOT,
	"allof": ALLOF,
	"anyof": ANYOF,
}

// IsArgument returns true for tokens corresponding to arguments and
// delimiters; it returns false otherwise.
//
func (tok Token) IsArgument() bool { return tok.Typ > ARGUMENT && tok.Typ < COMMAND }

// IsCommand returns true for tokens corresponding to commands;
// it returns false otherwise.
//
func (tok Token) IsCommand() bool { return tok.Typ > COMMAND && tok.Typ < TEST }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return tok.Typ > COMMAND && tok.Typ < IDENTIFIER }

// Compares Typ and Val but not position
// Used for debugging and testing
func (t Token) Equals(ot Token) bool { return t.Val == ot.Val && t.Typ == ot.Typ }