package interpreter

import "time"

// This function defines all native functions
// of the interpreter -> implements -> Callable

// Clock is the callable for time
type Clock struct {
}

var _ Callable = (*Clock)(nil)

// Returns the number of arguments of the fucntion
func (c *Clock) Arity() int {
	return 0
}

// Implements the call function of Clock
func (c *Clock) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

// Implements the string function of the clock
func (c *Clock) String() string {
	return "<native fn: clock>"
}
