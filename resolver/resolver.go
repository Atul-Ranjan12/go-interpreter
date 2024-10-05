package resolver

import (
	"fmt"
	"log"
	"strings"

	"github.com/Atul-Ranjan12/interpreter"
	"github.com/Atul-Ranjan12/parser/astprinter"
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

type FunctionType int

func PrintAST(stmt expressions.Stmt, depth int) {
	printer := astprinter.NewAstPrinter()
	result, err := stmt.Accept(printer)
	if err != nil {
		log.Printf("Error printing AST: %v", err)
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s\n", indent, result)

	switch s := stmt.(type) {
	case *expressions.Block:
		for _, subStmt := range s.Statements {
			PrintAST(subStmt, depth+1)
		}
	case *expressions.If:
		PrintAST(s.ThenBranch, depth+1)
		if s.ElseBranch != nil {
			PrintAST(s.ElseBranch, depth+1)
		}
	case *expressions.WhileStatement:
		PrintAST(s.Body, depth+1)
	}
}

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

type Resolver struct {
	Interpreter     *interpreter.Interpreter
	Scopes          []map[string]bool
	CurrentFunction FunctionType
	FunctionDepth   int
}

var _ expressions.ExprVisitor = (*Resolver)(nil)
var _ expressions.StmtVisitor = (*Resolver)(nil)

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		Scopes:          []map[string]bool{},
		CurrentFunction: FunctionTypeNone,
	}
}

func (r *Resolver) Error(token token.Token, message string) error {
	return fmt.Errorf("Resolution Error at '%v': %s", token.Lexeme, message)
}

func (r *Resolver) ResolveStatement(stmt expressions.Stmt) error {
	_, err := stmt.Accept(r)
	return err
}

func (r *Resolver) ResolveExpression(expr expressions.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) ResolveStatements(statements []expressions.Stmt) error {
	for _, statement := range statements {
		// log.Println("Resolving for statement: ")
		PrintAST(statement, 0)
		if err := r.ResolveStatement(statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) BeginScope() {
	r.Scopes = append(r.Scopes, make(map[string]bool))
}

func (r *Resolver) EndScope() {
	if len(r.Scopes) > 0 {
		r.Scopes = r.Scopes[:len(r.Scopes)-1]
	}
}

func (r *Resolver) Declare(name token.Token) error {
	// log.Println("This is r.Scopes right now: ", r.Scopes)
	if len(r.Scopes) == 0 {
		return nil
	}
	scope := r.Scopes[len(r.Scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return r.Error(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false

	return nil
}

func (r *Resolver) Define(name token.Token) {
	if len(r.Scopes) == 0 {
		return
	}
	r.Scopes[len(r.Scopes)-1][name.Lexeme] = true
}

func (r *Resolver) ResolveLocal(expr expressions.Expr, name token.Token) {
	for i := len(r.Scopes) - 1; i >= 0; i-- {
		if _, ok := r.Scopes[i][name.Lexeme]; ok {
			depth := len(r.Scopes) - 1 - i
			if r.FunctionDepth > 0 && depth >= r.FunctionDepth {
				r.Interpreter.Resolve(expr, 0) // It's a global variable for this function
			} else {
				r.Interpreter.Resolve(expr, depth)
			}
			return
		}
	}
	// Not found in any scope, it's a global
}

func (r *Resolver) ResolveFunction(function *expressions.Function) error {
	r.BeginScope()
	for _, param := range function.Params {
		if err := r.Declare(param); err != nil {
			return err
		}
		r.Define(param)
	}
	if err := r.ResolveStatements(function.Body); err != nil {
		return err
	}
	r.EndScope()
	return nil
}

// Visitor methods

func (r *Resolver) VisitBlockStmt(stmt *expressions.Block) (interface{}, error) {
	// log.Println("This is called: VisitBlockStmt")
	// log.Println("Scopes are: ", r.Scopes)
	r.BeginScope()
	err := r.ResolveStatements(stmt.Statements)
	r.EndScope()
	return nil, err
}

func (r *Resolver) VisitVarStmt(stmt *expressions.Var) (interface{}, error) {
	if err := r.Declare(stmt.Name); err != nil {
		return nil, err
	}
	if stmt.Initializer != nil {
		if err := r.ResolveExpression(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.Define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *expressions.Variable) (interface{}, error) {
	// log.Println("Reaching here for: ", expr.Name.Lexeme)
	if len(r.Scopes) > 0 {
		if initialized, ok := r.Scopes[len(r.Scopes)-1][expr.Name.Lexeme]; ok && !initialized {
			return nil, r.Error(expr.Name, "Can't read local variable in its own initializer.")
		}
	}
	r.ResolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *expressions.Assign) (interface{}, error) {
	if err := r.ResolveExpression(expr.Value); err != nil {
		return nil, err
	}
	r.ResolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *expressions.Function) (interface{}, error) {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	r.FunctionDepth++
	defer func() { r.FunctionDepth-- }()

	r.BeginScope()
	for _, param := range stmt.Params {
		r.Declare(param)
		r.Define(param)
	}
	err := r.ResolveStatements(stmt.Body)
	r.EndScope()

	return nil, err
}

func (r *Resolver) VisitExprStatementStmt(stmt *expressions.ExprStatement) (interface{}, error) {
	return nil, r.ResolveExpression(stmt.Expression)
}

func (r *Resolver) VisitIfStmt(stmt *expressions.If) (interface{}, error) {
	if err := r.ResolveExpression(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.ResolveStatement(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if err := r.ResolveStatement(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStatementStmt(stmt *expressions.PrintStatement) (interface{}, error) {
	return nil, r.ResolveExpression(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt *expressions.Return) (interface{}, error) {
	if stmt.Value != nil {
		return nil, r.ResolveExpression(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStatementStmt(stmt *expressions.WhileStatement) (interface{}, error) {
	if err := r.ResolveExpression(stmt.Condition); err != nil {
		return nil, err
	}
	return nil, r.ResolveStatement(stmt.Body)
}

func (r *Resolver) VisitBinaryExpr(expr *expressions.Binary) (interface{}, error) {
	if err := r.ResolveExpression(expr.Left); err != nil {
		return nil, err
	}
	return nil, r.ResolveExpression(expr.Right)
}

func (r *Resolver) VisitCallExpr(expr *expressions.Call) (interface{}, error) {
	if err := r.ResolveExpression(expr.Callee); err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		if err := r.ResolveExpression(argument); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *expressions.Grouping) (interface{}, error) {
	return nil, r.ResolveExpression(expr.Expression)
}

func (r *Resolver) VisitLiteralExpr(expr *expressions.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *expressions.Logical) (interface{}, error) {
	if err := r.ResolveExpression(expr.Left); err != nil {
		return nil, err
	}
	return nil, r.ResolveExpression(expr.Right)
}

func (r *Resolver) VisitUnaryExpr(expr *expressions.Unary) (interface{}, error) {
	return nil, r.ResolveExpression(expr.Right)
}

func (r *Resolver) VisitBreakExprExpr(expr *expressions.BreakExpr) (interface{}, error) {
	return nil, nil
}
