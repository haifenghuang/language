package ast

import "github.com/davidsbond/dave/token"

type (
	// The ExpressionStatement type represents a statement in the source
	// code that contains an expression. For example, assigning a value to an
	// existing variable:
	// 	a = 1
	ExpressionStatement struct {
		Token      *token.Token
		Expression Node
	}
)

// TokenLiteral returns the literal value of the token for this
// statement.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}