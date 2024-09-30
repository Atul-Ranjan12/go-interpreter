package parser

import (
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// ExprStatement parses an Expression statement
func (p *Parser) ExprStatement() (expressions.Expr, error) {
	// Parse the expression
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}

	// Check for semicolon at the end
	p.Consume(token.SEMICOLON, "Expect ; at the end of statement")

	return &expressions.Stmt{
		Expression: value,
	}, nil
}
