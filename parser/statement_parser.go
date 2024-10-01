package parser

import (
	"log"

	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Statement handles print statement expression statement or block
func (p *Parser) Statement() (expressions.Stmt, error) {
	// Match the if statement
	if p.Match(token.IF) {
		return p.IfStatement()
	}

	// Match print statement
	if p.Match(token.PRINT) {
		return p.PrintStatement()
	}

	// Match the while statement
	if p.Match(token.WHILE) {
		// Return while statement here
		return p.WhileStatement()
	}

	// Match the for statement
	if p.Match(token.FOR) {
		return p.ForStatement()
	}

	// Match left brace (block)
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

// IfStatement handles parsing for if branches
func (p *Parser) IfStatement() (expressions.Stmt, error) {
	_, err := p.Consume(token.LEFT_PAREN, "Expect '(' after if")
	if err != nil {
		return nil, err
	}

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	_, err = p.Consume(token.RIGHT_PAREN, "Expect ')' after if condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.Statement()
	if err != nil {
		return nil, err
	}

	var elseBranch expressions.Stmt
	if p.Match(token.ELSE) {
		elseBranch, err = p.Statement()
		if err != nil {
			return nil, err
		}
	}

	return &expressions.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
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

// WhileStatement parses a while statement
func (p *Parser) WhileStatement() (expressions.Stmt, error) {
	// Consume for left paren
	_, err := p.Consume(token.LEFT_PAREN, "Expect '(' after while statement")
	if err != nil {
		return nil, err
	}

	// Parse the expression
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	// Consume for the right paren
	_, err = p.Consume(token.RIGHT_PAREN, "Expect ')' after while")
	if err != nil {
		return nil, err
	}

	// Reaching here meaning successful parse of while
	body, err := p.Statement()
	if err != nil {
		log.Println("Error parsing the body")
		return nil, err
	}

	return &expressions.WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}
