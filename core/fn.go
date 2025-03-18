package core

import "fmt"

// Fn is a universal function type
type Fn struct {
	line uint
	pos  uint
	fn   func(scope Scope, args ...SExp) (SExp, error)
}

// compile-time interface checks
func a(_ SExp) {}
func _(fn Fn)  { a(SExp(fn)) }

// Eval returns the function itself, the real "evaluation" happens
// on lists
func (fn Fn) Eval(scope Scope) (SExp, error) {
	return fn, nil
}

func (fn Fn) Invoke(scope Scope, args ...SExp) (SExp, error) {
	return fn.fn(scope, args...)
}

func (fn Fn) String() string {
	// memory addresses are irrelevant but it's a way to distinguish
	// functions from each other. They can be anonymous
	return fmt.Sprintf("function @ %d:%d", fn.line, fn.pos)
}
