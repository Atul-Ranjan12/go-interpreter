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

// Specify explicitly to implement ExprVisitor and StmtVisitor
var _ expressions.ExprVisitor = (*ASTPrinter)(nil)
var _ expressions.StmtVisitor = (*ASTPrinter)(nil)

func (p *ASTPrinter) parenthesize(name string, exprs ...expressions.Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		if expr != nil {
			component, err := expr.Accept(p)
			if err != nil {
				return "", err
			}
			builder.WriteString(component.(string))
		} else {
			builder.WriteString("nil")
		}
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

// Existing ExprVisitor methods...

// New StmtVisitor methods

func (p *ASTPrinter) VisitBlockStmt(stmt *expressions.Block) (interface{}, error) {
	var builder strings.Builder
	builder.WriteString("(block")
	for _, statement := range stmt.Statements {
		builder.WriteString(" ")
		result, err := statement.Accept(p)
		if err != nil {
			return nil, err
		}
		builder.WriteString(result.(string))
	}
	builder.WriteString(")")
	return builder.String(), nil
}

func (p *ASTPrinter) VisitExprStatementStmt(stmt *expressions.ExprStatement) (interface{}, error) {
	return p.parenthesize("expr", stmt.Expression)
}

func (p *ASTPrinter) VisitPrintStatementStmt(stmt *expressions.PrintStatement) (interface{}, error) {
	return p.parenthesize("print", stmt.Expression)
}

func (p *ASTPrinter) VisitWhileStatementStmt(stmt *expressions.WhileStatement) (interface{}, error) {
	condition, err := stmt.Condition.Accept(p)
	if err != nil {
		return nil, err
	}
	body, err := stmt.Body.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(while %s %s)", condition, body), nil
}

func (p *ASTPrinter) VisitVarStmt(stmt *expressions.Var) (interface{}, error) {
	if stmt.Initializer == nil {
		return fmt.Sprintf("(var %s)", stmt.Name.Lexeme), nil
	}
	initializer, err := stmt.Initializer.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(var %s %s)", stmt.Name.Lexeme, initializer), nil
}

func (p *ASTPrinter) VisitIfStmt(stmt *expressions.If) (interface{}, error) {
	condition, err := stmt.Condition.Accept(p)
	if err != nil {
		return nil, err
	}
	thenBranch, err := stmt.ThenBranch.Accept(p)
	if err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		elseBranch, err := stmt.ElseBranch.Accept(p)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("(if %s %s %s)", condition, thenBranch, elseBranch), nil
	}
	return fmt.Sprintf("(if %s %s)", condition, thenBranch), nil
}

func (p *ASTPrinter) VisitBinaryExpr(expr *expressions.Binary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *expressions.Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitExprStatementExpr(expr *expressions.ExprStatement) (interface{}, error) {
	return p.parenthesize("expression-statement", expr.Expression)
}

func (p *ASTPrinter) VisitPrintStatementExpr(expr *expressions.PrintStatement) (interface{}, error) {
	return p.parenthesize("print-statement", expr.Expression)
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

func (p *ASTPrinter) VisitAssignExpr(expr *expressions.Assign) (interface{}, error) {
	return p.parenthesize("=", expr.Value)
}

func (p *ASTPrinter) VisitLogicalExpr(expr *expressions.Logical) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitVariableExpr(expr *expressions.Variable) (interface{}, error) {
	return expr.Name.Lexeme, nil
}

func (p *ASTPrinter) VisitBreakExprExpr(expr *expressions.BreakExpr) (interface{}, error) {
	return "break", nil
}

func (p *ASTPrinter) VisitCallExpr(expr *expressions.Call) (interface{}, error) {
	args := make([]expressions.Expr, len(expr.Arguments)+1)
	args[0] = expr.Callee
	copy(args[1:], expr.Arguments)
	return p.parenthesize("call", args...)
}

func (p *ASTPrinter) VisitFunctionStmt(stmt *expressions.Function) (interface{}, error) {
	var builder strings.Builder
	builder.WriteString("(fun ")
	builder.WriteString(stmt.Name.Lexeme)
	builder.WriteString(" (")

	for i, param := range stmt.Params {
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(param.Lexeme)
	}

	builder.WriteString(") ")

	for _, bodyStmt := range stmt.Body {
		result, err := bodyStmt.Accept(p)
		if err != nil {
			return nil, err
		}
		builder.WriteString(result.(string))
		builder.WriteString(" ")
	}

	builder.WriteString(")")

	return builder.String(), nil
}

func (p *ASTPrinter) VisitReturnStmt(stmt *expressions.Return) (interface{}, error) {
	return p.parenthesize("return", stmt.Value)
}

func (p *ASTPrinter) VisitClassStmt(stmt *expressions.Class) (interface{}, error) {
	return p.parenthesize("class", nil)
}

func (p *ASTPrinter) VisitGetExpr(expr *expressions.Get) (interface{}, error) {
	object, err := expr.Object.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(get %v %s)", object, expr.Name.Lexeme), nil
}

func (p *ASTPrinter) VisitSetExpr(expr *expressions.Set) (interface{}, error) {
	object, err := expr.Object.Accept(p)
	if err != nil {
		return nil, err
	}
	value, err := expr.Value.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(set %v %s %v)", object, expr.Name.Lexeme, value), nil
}

func main() {
	ExampleASTPrinter()
}
