package parser

import (
	"errors"

	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// ClassDeclaration parses the ClassDeclaration
func (p *Parser) ClassDeclaration() (expressions.Stmt, error) {
	// Consume name to be an identifier
	name, err := p.Consume(token.IDENTIFIER, "Expect Identifier after class declaration")
	if err != nil {
		return nil, err
	}

	// Consume left brace
	_, err = p.Consume(token.LEFT_BRACE, "Expect { after class identifier")
	if err != nil {
		return nil, err
	}

	var functions []*expressions.Function
	for !p.Check(token.RIGHT_BRACE) && !p.IsAtEnd() {
		fn, err := p.Function("method")
		if err != nil {
			return nil, err
		}

		if fn, ok := fn.(*expressions.Function); ok {
			functions = append(functions, fn)
		} else {
			return nil, errors.New("Expect to be a method in fucntion body")
		}
	}

	// Consume left brace
	_, err = p.Consume(token.RIGHT_BRACE, "Expect { after class identifier")
	if err != nil {
		return nil, err
	}

	return &expressions.Class{Name: *name, Methods: functions}, nil
}
