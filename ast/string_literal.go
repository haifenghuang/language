package ast

import "github.com/davidsbond/dave/token"

type (
	// The StringLiteral type represents a literal string within the source code.
	// For example, in a variable assignment:
	// 	var x = "test"
	// The literal value of 'test' is stored in the StringLiteral type.
	StringLiteral struct {
		Token *token.Token
		Value string
	}
)

// TokenLiteral returns the literal value of the token for this
// statement.
func (nl *StringLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *StringLiteral) String() string {
	return nl.Value
}