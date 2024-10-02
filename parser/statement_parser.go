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

	// Match for the return statement
	if p.Match(token.RETURN) {
		return p.ReturnStatement()
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
	if err != nil {
		return nil, err
	}
	return statements, nil
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
	if err != nil {
		return nil, err
	}

	return &expressions.ExprStatement{
		Expression: value,
	}, nil
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

// Or function handles the parsing of or logical expressions
func (p *Parser) Or() (expressions.Expr, error) {
	// Parse the left operations
	left, err := p.And()
	if err != nil {
		return nil, err
	}

	if p.Match(token.OR) {
		operator := p.Prev()
		right, err := p.And()
		if err != nil {
			return nil, err
		}

		return &expressions.Logical{
			Left:     left,
			Operator: *operator,
			Right:    right,
		}, nil
	}

	// Reaching here means there is no or
	return left, nil
}

// And function handles parsing the and operations
func (p *Parser) And() (expressions.Expr, error) {
	// Parse the left expression
	left, err := p.Equality()
	if err != nil {
		return nil, err
	}

	// Parse the right expression
	if p.Match(token.AND) {
		operator := p.Prev()
		right, err := p.Equality()
		if err != nil {
			return nil, err
		}

		return &expressions.Logical{
			Left:     left,
			Operator: *operator,
			Right:    right,
		}, nil
	}

	// Reaching here means there is no and
	return left, nil
}

// ForStatement parses a for statement by converting it
// to a while statement
func (p *Parser) ForStatement() (expressions.Stmt, error) {
	p.Consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer expressions.Stmt
	var err error
	if p.Match(token.SEMICOLON) {
		initializer = nil
	} else if p.Match(token.VAR) {
		initializer, err = p.VariableDeclaration()
	} else {
		initializer, err = p.Statement()
	}
	if err != nil {
		return nil, err
	}

	var condition expressions.Expr
	if !p.Check(token.SEMICOLON) {
		condition, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	p.Consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment expressions.Expr
	if !p.Check(token.RIGHT_PAREN) {
		increment, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	p.Consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	// Desugar for loop into while loop
	if increment != nil {
		body = &expressions.Block{
			Statements: []expressions.Stmt{
				body,
				&expressions.ExprStatement{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &expressions.Literal{Value: true}
	}
	body = &expressions.WhileStatement{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &expressions.Block{
			Statements: []expressions.Stmt{initializer, body},
		}
	}

	return body, nil
}
