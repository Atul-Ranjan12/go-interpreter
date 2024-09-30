package astprinter

import (
	"fmt"
	"strings"

	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

type ASTPrinter struct{}

func NewAstPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

// Specify explicitly to implement ExprVisitor
var _ expressions.ExprVisitor = (*ASTPrinter)(nil)

func (p *ASTPrinter) parenthesize(name string, exprs ...expressions.Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		component, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		builder.WriteString(component.(string))
	}
	builder.WriteString(")")

	return builder.String(), nil
}

func (p *ASTPrinter) Print(expr expressions.Expr) (string, error) {
	result, err := expr.Accept(p)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (p *ASTPrinter) VisitBinaryExpr(expr *expressions.Binary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *expressions.Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitLiteralExpr(expr *expressions.Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (p *ASTPrinter) VisitUnaryExpr(expr *expressions.Unary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

// Example usage
func ExampleASTPrinter() {
	expression := &expressions.Binary{
		Left: &expressions.Unary{
			Operator: token.Token{Type: token.MINUS, Lexeme: "-"},
			Right:    &expressions.Literal{Value: 123},
		},
		Operator: token.Token{Type: token.STAR, Lexeme: "*"},
		Right: &expressions.Grouping{
			Expression: &expressions.Literal{Value: 45.67},
		},
	}

	printer := &ASTPrinter{}
	result, err := printer.Print(expression)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
	// Output: (* (- 123) (group 45.67))
}

func main() {
	ExampleASTPrinter()
}
