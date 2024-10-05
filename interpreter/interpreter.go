package interpreter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Atul-Ranjan12/environment"
	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

// Interpreter struct represents the interpreter
type Interpreter struct {
	// Takes the environment class instance
	// Set up a global Environment
	Globals     *environment.Environment
	Environment *environment.Environment // Current environment
	Locals      map[expressions.Expr]int
}

func (i *Interpreter) Define(env *environment.Environment, callable Callable, callableName string) {
	env.Define(callableName, callable)
}

// NewInterpreter is the initializer for ther interpreter
func NewInterpreter() *Interpreter {
	// Create new environment
	globalEnvironment := environment.NewEnvironment(nil)
	// If it does not exist on the interpreter environment
	// it should be in the global environment, hence
	// global environment is an enclosing of the interpreter
	// environment

	// Define the interpreter
	i := &Interpreter{
		Globals:     globalEnvironment,
		Environment: globalEnvironment,
		Locals:      make(map[expressions.Expr]int),
	}

	// Define the native functions
	i.Define(i.Globals, &Clock{}, "clock")

	return i
}

// Interpreter implements the expressionVisitor interface
var _ expressions.ExprVisitor = (*Interpreter)(nil)
var _ expressions.StmtVisitor = (*Interpreter)(nil)

// Adds in values in the locals map
func (i *Interpreter) Resolve(expr expressions.Expr, depth int) {
	i.Locals[expr] = depth
}

// LookUpVariable looks up a variable in the locals map
func (i *Interpreter) LookupVariable(name *token.Token, expr expressions.Expr) (interface{}, error) {
	distance, ok := i.Locals[expr]
	if !ok {
		return i.Globals.Get(name)
	}
	return i.Environment.GetAt(distance, name.Lexeme), nil
}

func (i *Interpreter) RuntimeError(token token.Token, message string) error {
	return fmt.Errorf("Runtime Error at '%v': %s", token.Lexeme, message)
}

// Evaluate is the helper method for all evaluation
func (i *Interpreter) Evaluate(expr expressions.Expr) (interface{}, error) {
	return expr.Accept(i)
}

// Execute is the method for all statements
func (i *Interpreter) Execute(expr expressions.Stmt) (interface{}, error) {
	return expr.Accept(i)
}

// Interpret evaluates the expression and returns the result as a string
func (i *Interpreter) Interpret(statements []expressions.Stmt) error {
	for _, statement := range statements {
		_, err := i.Execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
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

// VisitExprStatement evaluates each expression statement
func (i *Interpreter) VisitExprStatementStmt(stmt *expressions.ExprStatement) (interface{}, error) {
	_, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// VisitPrintStatement evaluates each print statement
func (i *Interpreter) VisitPrintStatementStmt(stmt *expressions.PrintStatement) (interface{}, error) {
	// log.Println("Reaching here for print statement: ")
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(value)
	return nil, nil
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

// VisitVariableExpr implements the missing method for ExprVisitor
func (i *Interpreter) VisitVariableExpr(expr *expressions.Variable) (interface{}, error) {
	value, err := i.LookupVariable(&expr.Name, expr)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// VisitVarStmt implements the missing method for StmtVisitor
func (i *Interpreter) VisitVarStmt(stmt *expressions.Var) (interface{}, error) {
	// Here we evaluate the value of the initializer
	// if it has the initializer, else we put nil
	var value interface{}
	var err error
	if stmt.Initializer != nil {
		value, err = i.Evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.Environment.Define(stmt.Name.Lexeme, value)
	return value, nil
}

// VisitAssignExpr handles interpretation of an assignment
func (i *Interpreter) VisitAssignExpr(expr *expressions.Assign) (interface{}, error) {
	value, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.Locals[expr]
	if !ok {
		i.Environment.Assign(expr.Name, value)
	} else {
		// Assign at the particular scope
		i.Environment.AssignAt(distance, &expr.Name, value)
	}

	// Assign in the environment

	return value, nil
}

// VisitBlockStmt handles the interpretation of a block
func (i *Interpreter) VisitBlockStmt(block *expressions.Block) (interface{}, error) {
	// log.Println("This is called second")
	err := i.ExecuteBlock(block.Statements, environment.NewEnvironment(i.Environment))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) ExecuteBlock(statements []expressions.Stmt, environment *environment.Environment) error {
	previous := i.Environment
	i.Environment = environment
	defer func() { i.Environment = previous }()

	for _, statement := range statements {
		_, err := i.Execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitIfStmt handles the execution of an if else block
func (i *Interpreter) VisitIfStmt(statement *expressions.If) (interface{}, error) {
	condition, err := i.Evaluate(statement.Condition)
	if err != nil {
		return nil, err
	}

	// If the condition is true we evaluate the if then block
	if i.IsTruthy(condition) {
		_, err = i.Execute(statement.ThenBranch)
		if err != nil {
			return nil, err
		}
	} else if statement.ElseBranch != nil {
		_, err = i.Execute(statement.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// VisitLogicalExpr handles the execution of logical statements
func (i *Interpreter) VisitLogicalExpr(expr *expressions.Logical) (interface{}, error) {
	// Evaluate the left section of the exrpression
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	// Check for if the logical expression is OR
	// or AND
	if expr.Operator.Type == token.OR {
		// Or expression :: Short circuit the or expression
		// i.e if left is true return immediately
		if i.IsTruthy(left) {
			return left, nil
		}
	} else {
		// And expression :: Short circuit the and expression
		// i.e if left is false, return left (false) immediately
		if !i.IsTruthy(left) {
			return left, nil
		}
	}

	// If we reach here, it means that it is an
	// or expression where left = false, or an
	// and expression where left = true. Therefore,
	// the value depends on what is the value of the right
	// expression
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	return right, nil
}

// VisitWhileStatementStmt handles execution of while statements
func (i *Interpreter) VisitWhileStatementStmt(stmt *expressions.WhileStatement) (interface{}, error) {
	for {
		val, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}

		if !i.IsTruthy(val) {
			break
		}

		_, err = i.Execute(stmt.Body)
		if err != nil {

			// Check if a break statement was issued
			if err.Error() == "CODE_999_LOOP_BREAK" {
				break
			}

			// We had some other error
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitBreakExprExpr(expr *expressions.BreakExpr) (interface{}, error) {
	return nil, errors.New("CODE_999_LOOP_BREAK")
}
