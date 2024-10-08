package interpreter

import (
	"errors"
	"fmt"

	"github.com/Atul-Ranjan12/parser/expressions"
	"github.com/Atul-Ranjan12/token"
)

const CLASS_CONSTRUCTOR_NAME string = "construct"

// Class represents a class in runtime
type Class struct {
	Name    string
	Methods map[string]*Function
}

// Instance represents an instance of the class
// in runtime
type Instance struct {
	// Associated to the class
	ClassName *Class
	// Add fields
	Fields map[string]interface{}
}

// Class implements the callable interface
var _ Callable = (*Class)(nil)

// NewClass creates a new class
func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{
		Name:    name,
		Methods: methods,
	}
}

// NewInstance creates a new instance of a class
func NewInstance(className *Class) *Instance {
	return &Instance{
		ClassName: className,
		Fields:    make(map[string]interface{}),
	}
}

func (ins *Instance) ToString() string {
	return "Instance of " + ins.ClassName.Name
}

func (c *Class) ToString() string {
	return c.Name
}

// Call method for the class creates an instance
func (c *Class) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	classInstance := NewInstance(c)

	constructor, err := c.FindMethod(CLASS_CONSTRUCTOR_NAME)
	if err != nil {
		return nil, err
	}

	if constructor != nil {
		// There exists a constructor, bind the constructor to the
		// method and execute it
		_, err := constructor.Bind(classInstance).Call(interpreter, arguments)
		if err != nil {
			return nil, err
		}
	}

	return classInstance, nil
}

// Arity returns the arity for the class
func (c *Class) Arity() int {
	constructor, err := c.FindMethod(CLASS_CONSTRUCTOR_NAME)
	if err != nil || constructor == nil {
		return 0
	}
	return constructor.Arity()
}

// FindMethod finds the method and returns it
func (c *Class) FindMethod(name string) (*Function, error) {
	if fn, ok := c.Methods[name]; ok {
		return fn, nil
	}

	return nil, errors.New("Cannot find method in class")
}

// Get function returns a property of an instance
func (ins *Instance) Get(name *token.Token) (interface{}, error) {
	if v, ok := ins.Fields[name.Lexeme]; ok {
		return v, nil
	}

	if fn, err := ins.ClassName.FindMethod(name.Lexeme); err == nil {
		if fn != nil {
			// There is a method
			// log.Println("Reaching here")
			return fn.Bind(ins), nil
		}
	} else {
		return nil, err
	}

	return nil, errors.New(fmt.Sprintf("Property %s does not exist: ", name.Lexeme))
}

// Set function sets a property of an instance
func (ins *Instance) Set(name *token.Token, value interface{}) {
	ins.Fields[name.Lexeme] = value
}

// VisitClassStmt handles interpretation of calss
func (i *Interpreter) VisitClassStmt(stmt *expressions.Class) (interface{}, error) {
	i.Environment.Define(stmt.Name.Lexeme, nil)

	var methods map[string]*Function = make(map[string]*Function)
	for _, method := range stmt.Methods {
		fn := NewFunction(method, i.Environment)

		methods[method.Name.Lexeme] = fn
	}

	class := NewClass(stmt.Name.Lexeme, methods)

	i.Environment.Assign(stmt.Name, class)

	return nil, nil
}

// VisitClassStmt handles interpretation of calss
func (i *Interpreter) VisitGetExpr(expr *expressions.Get) (interface{}, error) {
	object, err := i.Evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	// Object has to be an instance of Instance
	if objectInstance, ok := object.(*Instance); ok {
		val, err := objectInstance.Get(&expr.Name)
		if err != nil {
			return nil, err
		}
		// log.Println("Getting val: ", val)
		return val, nil
	}

	return nil, errors.New("Only objects have properties")
}

// VisitSetExpr handles setting fields in objects
// VisitSetExpr handles setting fields in objects
func (i *Interpreter) VisitSetExpr(expr *expressions.Set) (interface{}, error) {
	object, err := i.Evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	objectInstance, ok := object.(*Instance)
	if !ok {
		return nil, errors.New("Only instances have fields")
	}

	val, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	objectInstance.Set(&expr.Name, val)

	return val, nil
}

// VisitThisExpr handles when it encounters a this function
func (i *Interpreter) VisitThisExpr(expr *expressions.This) (interface{}, error) {

	// log.Println("Reached here in visit this")

	val, err := i.LookupVariable(&expr.Keyword, expr)
	if err != nil {
		return nil, err
	}

	return val, nil
}
