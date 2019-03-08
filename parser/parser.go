// Package parser contains the Parser type, it is responsible for converting tokens
// produced by the lexer into the abstract syntax tree.
package parser

import (
	"fmt"
	"io"

	"github.com/davidsbond/dave/ast"
	"github.com/davidsbond/dave/lexer"
	"github.com/davidsbond/dave/token"
)

const (
	_ int = iota
	// LOWEST is the lowest, default level of precedence
	LOWEST

	// EQUALS is precedence level for '=' operators
	EQUALS

	// LESSGREATER is the precedence level for '>'/'<' operators
	LESSGREATER

	// SUM is the precedence level for '+'/'-' operators.
	SUM

	// PRODUCT is the precedence level for '*' operators.
	PRODUCT
)

var (
	precedence = map[token.Type]int{
		token.EQUALS:   EQUALS,
		token.ASSIGN:   EQUALS,
		token.PLUS:     SUM,
		token.MINUS:    SUM,
		token.LT:       LESSGREATER,
		token.GT:       LESSGREATER,
		token.ASTERISK: PRODUCT,
		token.SLASH:    PRODUCT,
		token.MOD:      PRODUCT,
	}
)

type (
	// The Parser type is responsible for iterating over tokens provided by the lexer
	// and converting them into the abstract syntax tree
	Parser struct {
		lexer        *lexer.Lexer
		currentToken *token.Token
		peekToken    *token.Token

		prefixParsers map[token.Type]prefixParseFn
		infixParsers  map[token.Type]infixParseFn

		Errors []error
	}

	prefixParseFn func() ast.Node
	infixParseFn  func(ast.Node) ast.Node
)

// New creates a new instance of the Parser type that will parse the tokens
// returned by the provided lexer.
func New(lexer *lexer.Lexer) (parser *Parser) {
	parser = &Parser{
		lexer: lexer,
	}

	parser.prefixParsers = map[token.Type]prefixParseFn{
		token.NUMBER: parser.parseNumberLiteral,
		token.STRING: parser.parseStringLiteral,
		token.CHAR:   parser.parseCharacterLiteral,
		token.TRUE:   parser.parseBooleanLiteral,
		token.FALSE:  parser.parseBooleanLiteral,
		token.IDENT:  parser.parseIdentifier,
	}

	parser.infixParsers = map[token.Type]infixParseFn{
		token.PLUS:     parser.parseInfixExpression,
		token.MINUS:    parser.parseInfixExpression,
		token.ASTERISK: parser.parseInfixExpression,
		token.LT:       parser.parseInfixExpression,
		token.GT:       parser.parseInfixExpression,
		token.SLASH:    parser.parseInfixExpression,
		token.MOD:      parser.parseInfixExpression,
		token.EQUALS:   parser.parseInfixExpression,
	}

	parser.nextToken()
	parser.nextToken()

	return
}

// Parse attempts to convert tokens provided by the underlying lexer into an instance
// of ast.AST that can be evaluated.
func (p *Parser) Parse() *ast.AST {
	ast := &ast.AST{}

	// While we're not at the end of the file, parse
	// statements and append them to the AST.
	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()

		if stmt != nil {
			ast.Nodes = append(ast.Nodes, stmt)
		}

		p.nextToken()
	}

	return ast
}

func (p *Parser) parseStatement() ast.Node {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.ATOMIC:
		return p.parseAtomicStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	// Check if we have a prefix parsing function for the current
	// token type.
	prefix, ok := p.prefixParsers[p.currentToken.Type]

	// If not, return an error
	if !ok {
		panic("no prefix parser")
	}

	// Otherwise, parse the LHS of the expression
	leftExp := prefix()

	// While the operator's precedence is less than that of the next
	// token's precedence, parse infix expressions
	for precedence < p.peekPrecedence() {
		// Check if we have an infix parsing function for the upcoming
		// token, return the expression as is if not
		infix := p.infixParsers[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		// Otherwise, parse the infix expression
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	return &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}
}

func (p *Parser) nextToken() {
	var err error

	p.currentToken = p.peekToken
	p.peekToken, err = p.lexer.NextToken()

	if err != nil && err != io.EOF {
		p.Errors = append(p.Errors, err)
	}
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if pr, ok := precedence[p.peekToken.Type]; ok {
		return pr
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedence[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	err := fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.Errors = append(p.Errors, err)
}
