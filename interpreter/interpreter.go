package interpreter

import (
	"fmt"
	"reflect"

	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Interpreter struct represents the interpreter
type Interpreter struct {
}

// Interpreter implements the expressionVisitor interface
var _ expressions.ExprVisitor = (*Interpreter)(nil)

func (i *Interpreter) RuntimeError(token token.Token, message string) error {
	return fmt.Errorf("Runtime Error at '%v': %s", token.Lexeme, message)
}

// Evaluate is the helper method for all evaluation
func (i *Interpreter) Evaluate(expr expressions.Expr) (interface{}, error) {
	return expr.Accept(i)
}

// Interpret evaluates the expression and returns the result as a string
func (i *Interpreter) Interpret(expression expressions.Expr) (string, error) {
	result, err := i.Evaluate(expression)
	if err != nil {
		return "", err
	}
	return i.stringify(result), nil
}

// stringify converts a value to its string representation
func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	if num, ok := value.(float64); ok {
		return fmt.Sprintf("%g", num)
	}
	return fmt.Sprintf("%v", value)
}

// VisitBinaryExpr handles Binary Operations
func (i *Interpreter) VisitBinaryExpr(expr *expressions.Binary) (interface{}, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	operator := expr.Operator.Type

	if i.IsString(left) && i.IsString(right) && operator == token.PLUS {
		return left.(string) + right.(string), nil
	}

	if !i.IsNumber(left) || !i.IsNumber(right) {
		return nil, i.RuntimeError(expr.Operator, "Binary operations require both operands to be numbers or strings")
	}

	leftNum, rightNum := left.(float64), right.(float64)

	switch operator {
	case token.MINUS:
		return leftNum - rightNum, nil
	case token.STAR:
		return leftNum * rightNum, nil
	case token.SLASH:
		if rightNum == 0 {
			return nil, i.RuntimeError(expr.Operator, "Division by zero")
		}
		return leftNum / rightNum, nil
	case token.PLUS:
		return leftNum + rightNum, nil
	case token.GREATER:
		return leftNum > rightNum, nil
	case token.GREATER_EQUAL:
		return leftNum >= rightNum, nil
	case token.LESS:
		return leftNum < rightNum, nil
	case token.LESS_EQUAL:
		return leftNum <= rightNum, nil
	case token.EQUAL_EQUAL:
		return i.IsEqual(left, right), nil
	case token.BANG_EQUAL:
		return !i.IsEqual(left, right), nil
	}

	return nil, i.RuntimeError(expr.Operator, "Unknown operator")
}

// VisitGroupingExpr handles Grouping Operations
func (i *Interpreter) VisitGroupingExpr(expr *expressions.Grouping) (interface{}, error) {
	return i.Evaluate(expr.Expression)
}

// VisitUnaryExpr handles unary operations
func (i *Interpreter) VisitUnaryExpr(expr *expressions.Unary) (interface{}, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		if !i.IsNumber(right) {
			return nil, i.RuntimeError(expr.Operator, "Operand must be a number")
		}
		return -right.(float64), nil
	case token.BANG:
		return !i.IsTruthy(right), nil
	}

	return nil, i.RuntimeError(expr.Operator, "Unknown operator")
}

// VisitLiteralExpr handles literal operations
func (i *Interpreter) VisitLiteralExpr(expr *expressions.Literal) (interface{}, error) {
	return expr.Value, nil
}

// IsNumber checks if an object is a number
func (i *Interpreter) IsNumber(object interface{}) bool {
	_, ok := object.(float64)
	return ok
}

// IsString checks if an object is a string
func (i *Interpreter) IsString(object interface{}) bool {
	_, ok := object.(string)
	return ok
}

// IsTruthy checks if the condition holds
func (i *Interpreter) IsTruthy(expr interface{}) bool {
	if expr == nil {
		return false
	}
	switch v := expr.(type) {
	case bool:
		return v
	default:
		return true
	}
}

// IsEqual checks if two objects are equal
func (i *Interpreter) IsEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return reflect.DeepEqual(a, b)
}
