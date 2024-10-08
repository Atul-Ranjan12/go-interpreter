package environment

import (
	"fmt"

	"github.com/Atul-Ranjan12/token"
)

// Environment keeps track of all the states
// and values of variables
type Environment struct {
	Enclosing *Environment
	// Values stores all the values in the
	// interpreter
	Values map[string]interface{}
}

// NewEnvironment creates a new environment for the
// interpreter
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

// Define defines a value in the map
func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

// Assigns a new value to the map
func (e *Environment) Assign(name token.Token, value interface{}) error {
	if _, exists := e.Values[name.Lexeme]; exists {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return fmt.Errorf("Undefined variable: %s", name.Lexeme)
}

func (e *Environment) Get(name *token.Token) (interface{}, error) {
	// Firstly check if the value of the variable is in that
	// current environment
	value, exists := e.Values[name.Lexeme]
	if exists {
		return value, nil
	}

	// If not found in current environment, check in enclosing environment
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	// If we've reached here, the variable is not found in any scope
	return nil, fmt.Errorf("Undefined variable %s.", name.Lexeme)
}

// GetAt gets the value at a distance
func (e *Environment) GetAt(distance int, name string) interface{} {
	// log.Println("This is e.Ancestor: ", e.Ancestor(distance).Values)
	val := e.Ancestor(distance).Values[name]
	// log.Println("This is what we got for this: ", val)
	return val
}

// AssignAt assigns the value at a distance
func (e *Environment) AssignAt(distance int, name *token.Token, value interface{}) {
	e.Ancestor(distance).Values[name.Lexeme] = value
}

// Ancestor gets the environment at a distance
func (e *Environment) Ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Enclosing
	}

	return env
}
