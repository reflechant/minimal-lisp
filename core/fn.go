package core

// Fn is a universal function type
type Fn struct {
	srcName string
	line    uint
	pos     uint
	name    string
	fn      func(scope Scope, args ...SExpr) (SExpr, error)
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
	loc := location(fn.srcName, fn.line, fn.pos)
	if fn.name != "" {
		return "function " + fn.name + " @ " + loc
	}
	return "function @ " + loc
}
