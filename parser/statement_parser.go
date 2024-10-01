package parser

import (
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Statement handles print statement expression statement or block
func (p *Parser) Statement() (expressions.Stmt, error) {
	if p.Match(token.PRINT) {
		return p.PrintStatement()
	}

	if p.Match(token.LEFT_BRACE) {
		// Return the block
		if statements, err := p.Block(); err == nil {
			return &expressions.Block{
				Statements: statements,
			}, nil
		} else {
			return nil, err
		}
	}

	return p.ExprStatement()
}

// Block gets all the block of statements
func (p *Parser) Block() ([]expressions.Stmt, error) {
	var statements []expressions.Stmt

	for !p.Check(token.RIGHT_BRACE) && !p.IsAtEnd() {
		declarationStatement, err := p.Declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declarationStatement)
	}

	// Check for the right brace
	_, err := p.Consume(token.RIGHT_BRACE, "Expect '}' after a block")
	return statements, err
}

// ExprStatement parses an Expression statement
func (p *Parser) ExprStatement() (expressions.Stmt, error) {
	// Parse the expression
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}

	// Check for semicolon at the end
	_, err = p.Consume(token.SEMICOLON, "Expect ; at the end of statement")

	return &expressions.ExprStatement{
		Expression: value,
	}, err
}

// PrintStatement parses a print statement
func (p *Parser) PrintStatement() (expressions.Stmt, error) {
	// Parse the expression
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}

	// Check for semicolon at the end
	_, err = p.Consume(token.SEMICOLON, "Expect ; at the end of statement")

	return &expressions.PrintStatement{
		Expression: value,
	}, err
}
