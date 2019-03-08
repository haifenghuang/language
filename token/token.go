// Package token contains constants used to identify collections of
// lexemes within the source code.
package token

const (
	// IDENT is the token type used for an identifier
	IDENT = "identifier"

	// NUMBER is the token type used for a number
	NUMBER = "number"

	// EOF is the token type used when the end of the source
	// code is reached.
	EOF = "EOF"

	// VAR is the token type used when declaring a mutable
	// variable
	VAR = "var"

	// CONST is the token type used when declaring an immutable
	// variable.
	CONST = "const"

	// ATOMIC is the token type used when declaring an atomic variable
	ATOMIC = "atomic"

	// ASSIGN is the token type used for assignment statements
	ASSIGN = "="

	// EQUALS is the token type used for equality statements.
	EQUALS = "=="

	// STRING is the token type used for opening/closing literal strings.
	STRING = `"`

	// CHAR is the token type used for opening/closing literal characters.
	CHAR = `'`

	// PLUS is the token type used for addition
	PLUS = "+"

	// MINUS is the token type used for subtraction/negative numbers
	MINUS = "-"

	// LT is the token type used for the less than symbol
	LT = "<"

	// GT is the token type used for the greater than symbol
	GT = ">"

	// ASTERISK is the token type used for the '*' symbol.
	ASTERISK = "*"

	// SLASH is the token type used for the '/' symbol.
	SLASH = "/"

	// MOD is the token type used for the percent symbol.
	MOD = "%"

	// TRUE is the token type used for a boolean 'true'
	TRUE = "true"

	// FALSE is the token type used for a boolean 'false'
	FALSE = "false"

	// NEWLINE is the token type used for newline characters.
	NEWLINE = "\n"
)

var (
	keywords = map[string]Type{
		"const":  CONST,
		"var":    VAR,
		"atomic": ATOMIC,
		"true":   TRUE,
		"false":  FALSE,
	}
)

type (
	// The Type type is used to differentiate between tokens
	Type string

	// The Token type represents a collection of lexemes
	// within the source code.
	Token struct {
		Type    Type
		Literal string
		Line    int
		Column  int
	}
)

// New creates a new instance of the Token type using the provided literal, token
// type, line & column numbers.
func New(literal string, t Type, l, c int) *Token {
	return &Token{
		Literal: literal,
		Type:    t,
		Line:    l,
		Column:  c,
	}
}

func (t *Token) String() string {
	return t.Literal
}

// LookupIdentifier checks to see if a given identifier has a matching
// token type. These are usually used for keywords.
func LookupIdentifier(ident string) (t Type) {
	t, ok := keywords[ident]

	if !ok {
		t = IDENT
	}

	return
}
