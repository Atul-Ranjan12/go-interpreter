package parser

import (
	"errors"
	"fmt"
	"log"

	"github.com/Atul-Ranjan12/errorHandler"
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Grammar for lang
//
// Grammar for lang
// program -> statement* EOF ;
// statement -> exprStatement | printStatement ;
// exprStatement -> expression ;
// printStatement -> print ( expression ) ;
//
//
// Grammar for expressions
//
// expression -> equality
// equality -> comparison ( ( != | == ) comparison)*
// comparison -> term ( ( > | >= | < | <= ) term )*
// term -> factor ( ( / | * ) factor)*
// factor -> unary ( ( + | - ) unary)*
// unary -> ( ! | - ) unary | primary
// primary -> NUMBER | STRING | "true" | "false" | "nil"
// 			  | "(" expression ")"

// Parser represents the parser for lang
type Parser struct {
	Tokens  []*token.Token
	Current int

	// Handle errors
	ErrorHandler errorHandler.ErrorHandler
}

// NewParser creates a new parser
func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

// Expression is the root of the tree
func (p *Parser) Expression() (expressions.Expr, error) {
	return p.Equality()
}

// Equality checks if an expression is equality
func (p *Parser) Equality() (expressions.Expr, error) {
	expr, err := p.Comparison()
	if err != nil {
		return nil, err
	}

	for p.Match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.Prev()
		right, err := p.Comparison()
		if err != nil {
			return nil, err
		}
		expr = &expressions.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

// Comparison checks if an expression is a comparison
func (p *Parser) Comparison() (expressions.Expr, error) {
	expr, err := p.Term()
	if err != nil {
		return nil, err
	}

	for p.Match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.Prev()
		right, err := p.Term()
		if err != nil {
			return nil, err
		}
		expr = &expressions.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

// Term checks if an expression is a term
func (p *Parser) Term() (expressions.Expr, error) {
	expr, err := p.Factor()
	if err != nil {
		return nil, err
	}

	for p.Match(token.MINUS, token.PLUS) {
		operator := p.Prev()
		right, err := p.Factor()
		if err != nil {
			return nil, err
		}
		expr = &expressions.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) Factor() (expressions.Expr, error) {
	expr, err := p.Unary()
	if err != nil {
		return nil, err
	}

	for p.Match(token.SLASH, token.STAR) {
		operator := p.Prev()
		right, err := p.Unary()
		if err != nil {
			return nil, err
		}
		expr = &expressions.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

// Unary checks if an expression is unary
func (p *Parser) Unary() (expressions.Expr, error) {
	if p.Match(token.BANG, token.MINUS) {
		operator := p.Prev()
		right, err := p.Unary()
		if err != nil {
			return nil, err
		}
		return &expressions.Unary{
			Operator: *operator,
			Right:    right,
		}, nil
	}

	return p.Primary()
}

// Primary checks if an expression is primary
func (p *Parser) Primary() (expressions.Expr, error) {
	if p.Match(token.FALSE) {
		return &expressions.Literal{Value: false}, nil
	}
	if p.Match(token.TRUE) {
		return &expressions.Literal{Value: true}, nil
	}
	if p.Match(token.NIL) {
		return &expressions.Literal{Value: nil}, nil
	}

	if p.Match(token.NUMBER, token.STRING) {
		return &expressions.Literal{Value: p.Prev().Literal}, nil
	}

	if p.Match(token.LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		p.Consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &expressions.Grouping{Expression: expr}, nil
	}

	// If we get here, we have an error
	// Here we have reached the EOF with incomplete parse
	err := p.Error(p.Peek(), "Expect expression.")
	return nil, err
}

// Consume checks if the current token is of the particular
// tokenType, if it is it moves ahead, else it throws an
// error
func (p *Parser) Consume(tokenType token.TokenType, message string) (*token.Token, error) {
	if p.Check(tokenType) {
		return p.Advance(), nil
	}

	err := p.Error(p.Peek(), message)
	return nil, err
}

// TODO: Handle errors here
// Error reports an error at the given token
func (p *Parser) Error(token *token.Token, message string) error {
	// TODO: implement error reporting logic here
	// For now, we'll just panic
	return errors.New(fmt.Sprintf("Error at '%v': %s", token.Lexeme, message))
}

// Match function matches a a token type in
// current position
// if a grammar has multiple tokens, it is useful to
// check for all the tokens
func (p *Parser) Match(tokenType ...token.TokenType) bool {
	for _, tt := range tokenType {
		// For each token in tokenType
		if p.Check(tt) {
			p.Advance()
			return true
		}
	}
	return false
}

// Check function returns true if the current token
// is of the given type
func (p *Parser) Check(tokenType token.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().Type == tokenType
}

// Advance moves the pointer of the parser forward
func (p *Parser) Advance() *token.Token {
	// Returns the current value and moves the pointer
	// one step forward
	if !p.IsAtEnd() {
		p.Current++
	}

	log.Println("Reaching here: this is p.Prev: ", p.Prev())
	return p.Prev()
}

// IsAtEnd Function checks if we are at the end
func (p *Parser) IsAtEnd() bool {
	return p.Peek().Type == token.EOF
}

// Peek function takes a look at the current token
func (p *Parser) Peek() *token.Token {
	return p.Tokens[p.Current]
}

// Prev function returns the previous value
func (p *Parser) Prev() *token.Token {
	return p.Tokens[p.Current-1]
}

// Parse function parses the tokens
func (p *Parser) Parse() (expressions.Expr, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}
	if !p.IsAtEnd() {
		return nil, p.Error(p.Peek(), "Unexpected tokens after expression")
	}
	return expr, nil
}

// synchronize discards tokens until it finds a likely statement boundary
func (p *Parser) synchronize() {
	p.Advance()

	for !p.IsAtEnd() {
		if p.Prev().Type == token.SEMICOLON {
			return
		}

		switch p.Peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.Advance()
	}
}
