package interpreter

import (
	"fmt"

	"github.com/Atul-Ranjan12/environment"
	"github.com/Atul-Ranjan12/parser/expressions"
)

// This file handles interpratation of functions

type Callable interface {
	// This will have an arity
	// number of arguments of a function
	Arity() int
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
}

// Type for the Return statements
// It is treated as a runtime exception like break
type ReturnValue struct {
	Value interface{}
}

// Return implements the error interface
var _ error = (*ReturnValue)(nil)

// Implement the Error function
func (r *ReturnValue) Error() string {
	return "Return statement"
}

// The Call function handles function calls
func (i *Interpreter) VisitCallExpr(expr *expressions.Call) (interface{}, error) {
	// Evaluate the callee
	callee, err := i.Evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arg, err := i.Evaluate(argument)
		if err != nil {
			return nil, err
		}

		// Append to arguments
		arguments = append(arguments, arg)
	}

	// See if the callee can be a function
	// "Not a function"() is not a function
	function, ok := callee.(Callable)
	if !ok {
		return nil, i.RuntimeError(expr.Paren, "Can only call functions and classes")
	}

	if len(arguments) != function.Arity() {
		return nil, i.RuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

// Function is the structure for a function
type Function struct {
	Declaration *expressions.Function
	// Create a closure for the environment
	Closure *environment.Environment
}

var _ Callable = (*Function)(nil)

// NewFunction is the initializer to the function
func NewFunction(declaration *expressions.Function, closure *environment.Environment) *Function {
	return &Function{
		Declaration: declaration,
		Closure:     closure,
	}
}

// Arity returns the number of arguments
func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

// ToString returns the function as what it is
func (f *Function) ToString() string {
	return "<fn " + f.Declaration.Name.Lexeme + " >"
}

// Call handles the calling of funciton
func (f *Function) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	env := environment.NewEnvironment(f.Closure)

	for i := 0; i < len(args); i++ {
		env.Define(f.Declaration.Params[i].Lexeme, args[i])
	}

	err := i.ExecuteBlock(f.Declaration.Body, env)
	if err != nil {
		// See if the error thrown is actually the return value
		returnValue, ok := err.(*ReturnValue)
		if ok {
			// It is a return value
			return returnValue.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

// The VisitFunctionStmt declares a function
func (i *Interpreter) VisitFunctionStmt(stmt *expressions.Function) (interface{}, error) {
	function := NewFunction(stmt, i.Environment)
	i.Environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

// The VisitReturnStmt handles when a function returns a value
func (i *Interpreter) VisitReturnStmt(stmt *expressions.Return) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.Value != nil {
		value, err = i.Evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, &ReturnValue{Value: value}
}
