package core

import (
	"errors"
	"fmt"
)

var (
	true_  = Atom("t")
	false_ = List{}
)

// BuiltinScope defines the default scope for all evaluations that is always present
func BuiltinScope() Scope {
	fns := map[string]Fn{
		"quote": quote,
		"atom":  atom,
		"eq":    eq,
		"car":   car,
		"cdr":   cdr,
		"cons":  cons,
		"cond":  cond,
	}

	return Scope{
		parent: nil, // this is supposed to be the root scope
		fns:    fns,
		vals:   nil, // no built-in values defined so far
	}
}

// The following 7 operators are the "Maxwell equations of programming" as Paul Graham called them

// quote returns it's parameter unchanged. (quote x) returns x.
// Exists mostly to prevent list evaluation (which is their default behaviour)
func quote(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("quote: expects 1 argument, %d given", len(args)))
	}

	return args[0], nil
}

// atom returns the atom t if the value of x is an atom or the empty
// list. Otherwise it returns ().
func atom(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("atom: expects 1 argument, %d given", len(args)))
	}
	switch v := args[0].(type) {
	case Atom:
		return true_, nil
	case List:
		if v.IsEmpty() {
			return true_, nil
		}
	}

	return false_, nil
}

// eq returns t if the values of x and y are the same atom or both the
// empty list, and () otherwise
func eq(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("eq: expects 2 arguments, %d given", len(args)))
	}

	// if equal atoms return t
	a1, ok1 := args[0].(Atom)
	a2, ok2 := args[1].(Atom)
	if ok1 && ok2 && a1 == a2 {
		return true_, nil
	}

	// if both are empty lists return t
	l1, ok1 := args[0].(List)
	l2, ok2 := args[1].(List)
	if ok1 && ok2 && l1.IsEmpty() && l2.IsEmpty() {
		return true_, nil
	}

	// return ()
	return false_, nil
}

// car expects it's only argument to be a list, and returns its first element
func car(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("car: expects 1 argument, %d given", len(args)))
	}
	l, ok := args[0].(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("car: argument must be a list, instead was given %v", args[0]))
	}
	if l.IsEmpty() {
		return nil, errors.New("car: can't return the 1st element of an empty list")
	}

	return l.head.expr, nil
}

// cdr expects its only argument to be a list, and returns everything after the first element (may be an empty list).
func cdr(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("cdr: expects 1 argument, %d given", len(args)))
	}
	l, ok := args[0].(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("cdr: argument must be a list, instead was given %v", args[0]))
	}

	// we return an empty list if there is no 1st element or there is nothing after it
	return l.Rest(), nil
}

// cons adds an element to the front of a list.
//
// (cons x y) expects the value of y to be a list, and returns a list
// containing the value of x followed by the elements of the value of y
func cons(scope Scope, args ...Expr) (Expr, error) {
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("cons: expects 2 arguments, %d given", len(args)))
	}
	rest, ok := args[1].(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("cons: 2nd argument must be a list, instead was given %v", args[1]))
	}

	return rest.Cons(args[0]), nil
}

// cond performs conditional evaluation.
//
// (cond (p1 e1 ) ... (pn en)) is evaluated as follows. The p
// expressions are evaluated in order until one returns t. When one is
// found, the value of the corresponding e expression is returned as the
// value of the whole cond expression.
func cond(scope Scope, args ...Expr) (Expr, error) {
	for i, arg := range args {
		l, ok := arg.(List)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is not a list, it's %v", i+1, args[i]))
		}
		pred := l.First()
		if pred == nil {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is missing a predicate", i+1))
		}
		val := l.Rest().First()
		if val == nil {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is missing a value", i+1))
		}
		condition, err := pred.Eval(scope)
		if err != nil {
			return nil, fmt.Errorf("cond: evaluation error in condition #%d: %w", i+1, err)
		}
		if condition == true_ {
			return val.Eval(scope)
		}
	}

	return false_, nil
}
