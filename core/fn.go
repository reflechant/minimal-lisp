package core

import "fmt"

// Fn is a universal function type
type Fn struct {
	line uint
	pos  uint
	name string
	fn   func(scope Scope, args ...SExpr) (SExpr, error)
}

// compile-time interface checks
func a(_ SExpr) {}
func _(fn Fn)   { a(SExpr(fn)) }

// Eval returns the function itself, the real "evaluation" happens
// on lists using Fn.Invoke
func (fn Fn) Eval(_ Scope) (SExpr, error) {
	return fn, nil
}

func (fn Fn) Invoke(scope Scope, args ...SExpr) (SExpr, error) {
	return fn.fn(scope, args...)
}

func (fn Fn) String() string {
	// memory addresses are irrelevant but it's a way to distinguish
	// functions from each other. They can be anonymous
	return fmt.Sprintf("function %s @ %d:%d", fn.name, fn.line, fn.pos)
}
