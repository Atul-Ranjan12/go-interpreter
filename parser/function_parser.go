package parser

import (
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Call function handles parsing a function call
func (p *Parser) Call() (expressions.Expr, error) {
	// Parse the primary expression
	// should give an identifier as an output
	expr, err := p.Primary()
	if err != nil {
		return nil, err
	}

	for true {
		if p.Match(token.LEFT_PAREN) {
			// Finish parsing the arguments
			expr, err = p.FinishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.Match(token.DOT) {
			// Return a get expression
			name, err := p.Consume(token.IDENTIFIER, "Expect property name after .")
			if err != nil {
				return nil, err
			}
			expr = &expressions.Get{Name: *name, Object: expr}
		} else {
			break
		}
	}

	return expr, nil
}

// FinishCall function handles parsing the arguments
func (p *Parser) FinishCall(callee expressions.Expr) (expressions.Expr, error) {
	var arguments []expressions.Expr

	if !p.Check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, p.Error(p.Peek(), "Can't have more than 255 arguments.")
			}

			arg, err := p.Expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.Match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.Consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &expressions.Call{Callee: callee, Paren: *paren, Arguments: arguments}, nil
}

// Function handles the parsing of function declarations
func (p *Parser) Function(kind string) (expressions.Stmt, error) {
	// At this point we have already recognized the fun
	// keyword, recognize an identifier after the keyword
	name, err := p.Consume(token.IDENTIFIER, "Expect "+kind+" name")
	if err != nil {
		return nil, err
	}

	// Parse the left paren
	_, err = p.Consume(token.LEFT_PAREN, "Expect ( after function name")
	if err != nil {
		return nil, err
	}

	// Get the parameters
	var parameters []token.Token
	if !p.Check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, p.Error(p.Peek(), "Can't have more than 255 arguments in a function")
			}

			param, err := p.Consume(token.IDENTIFIER, "Expect parameter name")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, *param)

			if !p.Match(token.COMMA) {
				break
			}
		}
	}

	// Consume the right bracket
	_, err = p.Consume(token.RIGHT_PAREN, "Expect ) after parameters")
	if err != nil {
		return nil, err
	}

	// Consume the right brace before parsing a block
	_, err = p.Consume(token.LEFT_BRACE, "Expect '{' before function body")
	if err != nil {
		return nil, err
	}

	// Get the body
	body, err := p.Block()
	if err != nil {
		return nil, err
	}

	return &expressions.Function{
		Name:   *name,
		Params: parameters,
		Body:   body,
	}, nil
}

// ReturnStatement handles parsing a return statement
func (p *Parser) ReturnStatement() (expressions.Stmt, error) {
	// Get the token
	returnToken := p.Prev()

	var value expressions.Expr = nil
	var err error

	// If no semicolon immediately, there is a value to be
	// returned
	if !p.Check(token.SEMICOLON) {
		value, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.Consume(token.SEMICOLON, "Expect ; after retrurn statement")
	if err != nil {
		return nil, err
	}

	return &expressions.Return{
		Keyword: *returnToken,
		Value:   value,
	}, nil
}
