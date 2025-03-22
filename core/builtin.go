package core

import (
	"errors"
	"fmt"
	"strings"
)

var (
	True  = Symbol{name: "t"}
	False = List{}
)

// BuiltinScope returns the default environment for all evaluations that is always present.
// It contains the 7 basic operators from "The Roots of LISP" + `lambda` + `defun`
func BuiltinScope() Scope {
	fns := map[string]SExpr{
		"quote": Fn{name: "quote", fn: quote},
		"atom":  Fn{name: "atom", fn: atom},
		"eq":    Fn{name: "eq", fn: eq},
		"car":   Fn{name: "car", fn: car},
		"cdr":   Fn{name: "cdr", fn: cdr},
		"cons":  Fn{name: "cons", fn: cons},
		"cond":  Fn{name: "cond", fn: cond},
		// lambda and defun are placed here for convenience
		"lambda": Fn{name: "lambda", fn: lambda},
		"label":  Fn{name: "label", fn: label},
		"defun":  Fn{name: "defun", fn: defun},
		// print - for a rudimentary REPL
		"print": Fn{name: "print", fn: print},
	}

	return Scope{
		parent: nil, // this is supposed to be the root scope
		vals:   fns, // no built-in values defined so far
	}
}

// The following 7 operators are the "Maxwell equations of programming" as Paul Graham called them.
// In LISP 1.5 manual they are called "elementary functions of S-expressions"

// quote returns it's parameter unchanged. (quote x) returns x.
// Exists mostly to prevent list evaluation (which is their default behaviour)
func quote(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("quote: expects 1 argument, %d given", len(args)))
	}

	return args[0], nil
}

// atom returns the atom t if the value of x is an atom or the empty
// list. Otherwise it returns ().
func atom(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("atom: expects 1 argument, %d given", len(args)))
	}
	val, err := args[0].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("atom: evaluation error: %w", err)
	}
	switch v := val.(type) {
	case Symbol:
		return True, nil
	case List:
		if v.IsEmpty() {
			return True, nil
		}
	}

	return False, nil
}

// eq returns t if the values of x and y are the same atom or both the
// empty list, and () otherwise
func eq(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("eq: expects 2 arguments, got %d", len(args)))
	}

	// evaluate arguments
	arg1, err := args[0].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("eq: argument 1 evaluation error: %w", err)
	}
	arg2, err := args[1].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("eq: argument 2 evaluation error: %w", err)
	}

	// if equal atoms return t
	// (the only atoms present now are symbols)
	a1, ok1 := arg1.(Symbol)
	a2, ok2 := arg2.(Symbol)
	if ok1 && ok2 && a1.name == a2.name {
		return True, nil
	}

	// if both are empty lists return t
	l1, ok1 := arg1.(List)
	l2, ok2 := arg2.(List)
	if ok1 && ok2 && l1.IsEmpty() && l2.IsEmpty() {
		return True, nil
	}

	// return '()
	return False, nil
}

// car expects it's only argument to be a list, and returns its first element
func car(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("car: expects 1 argument, got %d", len(args)))
	}
	// evaluate argument
	arg, err := args[0].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("car: evaluation error: %w", err)
	}
	l, ok := arg.(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("car: argument must be a list, got %v", args[0]))
	}
	if l.First() == nil {
		return List{}, nil
	}

	return l.First(), nil
}

// cdr expects its only argument to be a list, and returns everything after the first element (may be an empty list).
func cdr(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 1 {
		return nil, errors.New(fmt.Sprintf("cdr: expects 1 argument, got %d", len(args)))
	}
	// evaluate argument
	arg, err := args[0].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("cdr: evaluation error: %w", err)
	}
	l, ok := arg.(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("cdr: argument must be a list, got %v", args[0]))
	}

	if l.Rest() == nil {
		return List{}, nil
	}

	return l.Rest(), nil
}

// cons adds an element to the front of a list.
//
// (cons x y) expects the value of y to be a list, and returns a list
// containing the value of x followed by the elements of the value of y
func cons(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("cons: expects 2 arguments, got %d", len(args)))
	}
	// evaluate arguments
	arg1, err := args[0].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("cons: argument 1 evaluation error: %w", err)
	}
	arg2, err := args[1].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("cons: argument 2 evaluation error: %w", err)
	}
	rest, ok := arg2.(List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("cons: 2nd argument must be a list, got %v", arg2))
	}

	return List{
		first:  arg1,
		second: rest,
	}, nil
}

// cond performs conditional evaluation.
//
// (cond (p1 e1) ... (pn en)) is evaluated as follows. The p
// expressions are evaluated in order until one returns t. When one is
// found, the value of the corresponding e expression is returned as the
// value of the whole cond expression.
func cond(scope Scope, args ...SExpr) (SExpr, error) {
	fmt.Println("cond", len(args), args)
	for i, arg := range args {
		fmt.Printf("arg %T %v\n", arg, arg)
		p, ok := arg.(List)
		if !ok {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is not a list, it's %v", i+1, args[i]))
		}
		pred := p.First()
		fmt.Println("predicate", pred)
		if pred == nil {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is missing a predicate", i+1))
		}
		val := p.Rest()
		fmt.Println("rest", val)
		if val == nil {
			return nil, errors.New(fmt.Sprintf("cond: argument #%d is missing a return value", i+1))
		}
		condition, err := pred.Eval(scope)
		fmt.Printf("condition: %T, %v\n", condition, condition)
		if err != nil {
			return nil, fmt.Errorf("cond: evaluation error in condition #%d: %w", i+1, err)
		}
		switch v := condition.(type) {
		case Symbol:
			if v.name == True.name {
				fmt.Println("condition is true")
				x, _ := val.Eval(scope)
				fmt.Printf("returning: %T %v", x, x)
				return val.Eval(scope)
			}
		}
	}

	return False, nil
}

// lambda creates an anonymous function and returns it
// example: (lambda (a b) (cons a b))
func lambda(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) < 1 {
		return nil, errors.New("lambda: expects at least 1 argument")
	}
	paramList, ok := args[0].(List)
	if !ok {
		return nil, errors.New("lambda: first parameter is not a list")
	}

	// get function body
	if len(args) < 2 {
		// if function body is empty, lambda will return an empty list
		return Fn{
			line: paramList.line,
			pos:  paramList.pos,
			fn: func(scope Scope, args ...SExpr) (SExpr, error) {
				return List{}, nil
			},
		}, nil
	}
	body := args[1]

	params := []Symbol{}
	for p := range paramList.Items() {
		p, ok := p.(Symbol)
		if !ok {
			return nil, errors.New(fmt.Sprintf("lambda: parameter #%d in parameter list is not a symbol", len(params)+1))
		}
		params = append(params, p)
	}

	return Fn{
		line: paramList.line,
		pos:  paramList.pos,
		fn: func(scope Scope, args ...SExpr) (SExpr, error) {
			if len(params) != len(args) {
				return nil, errors.New(fmt.Sprintf("lambda: arity error: expected %d parameters, got %d", len(params), len(args)))
			}

			// evaluate operands and bind them to parameter symbols
			scope = scope.NewLayer()
			for i, a := range args {
				v, err := a.Eval(scope)
				if err != nil {
					return nil, fmt.Errorf("lambda: error evaluating parameter #%d=%v: %w", i+1, a, err)
				}
				scope.Bind(params[i].name, v)
			}

			// evaluate function body
			return body.Eval(scope)
		},
	}, nil
}

// label creates a named function in scope and returns it
func label(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 2 {
		return nil, errors.New("label: expects 2 arguments (function name and a lambda expression)")
	}
	fnSym, ok := args[0].(Symbol)
	if !ok {
		return nil, errors.New(fmt.Sprintf("label: 1st parameter (function name) is not a symbol but %v", args[0]))
	}

	// args[1] is expected to be a lambda (but could be another function)
	fnVal, err := args[1].Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("label: %w", err)
	}
	fn, ok := fnVal.(Fn)
	if !ok {
		return nil, errors.New(fmt.Sprintf("label: second parameter is not a function but %v", fnVal))
	}
	scope.Bind(fnSym.name, fn)

	return fn, nil
}

// defun is a syntactic sugar for `label`
func defun(scope Scope, args ...SExpr) (SExpr, error) {
	if len(args) != 3 {
		return nil, errors.New("defun: expects 3 arguments (function name, parameter list, function body)")
	}
	fn, err := label(scope, args[0], List{
		first: Symbol{
			name: "lambda",
		},
		second: List{
			first:  args[1],
			second: args[2],
		},
	})
	// label binds function to the name for us
	if err != nil {
		return nil, fmt.Errorf("defun: %w", err)
	}

	return fn, nil
}

func print(scope Scope, args ...SExpr) (SExpr, error) {
	argsValStrs := []string{}
	for i, a := range args {
		aVal, err := a.Eval(scope)
		if err != nil {
			return nil, fmt.Errorf("eval: argument #%d: %w", i, err)
		}
		argsValStrs = append(argsValStrs, aVal.String())
	}

	fmt.Println(strings.Join(argsValStrs, " "))
	return nil, nil
}
