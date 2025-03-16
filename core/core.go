// Package core implements the core LISP interpreter
// described in Paul Graham's paper "The Roots of LISP" in Go.
// You may find the article here: https://paulgraham.com/rootsoflisp.html.
// It's not intended to be practically usable
// only to showcase the points made in the paper
// and demonstrate the "Maxwell equations" nature of LISP
// by implementing the Eval function in itself
// using 7 axiomatic operators
// and a little bit of other functionality (like function invocation).
// The eval.lisp file is a copy from https://paulgraham.com/rootsoflisp.html.
package core

import (
	"errors"
	"fmt"
)

// SExpr defines S-expressions as an interface. It's implemented by: atom, list
type SExpr interface {
	// Eval is the only mandatory operation needed on expressions. Returns the expression value
	Eval(scope Scope) (SExpr, error)
	// Print is used to give expressions visual representation (P in REPL)
	Print() string
}

// Fn is a universal function type
type Fn func(scope Scope, args ...SExpr) (SExpr, error)

// Scope stores known functions and values. Scope is lexical. This is
// a so-called LISP-1 - name conflicts are not allowed and one symbol can
// only be either a function or a value. Shadowing is allowed and you can
// rebind a function symbol to a value and vice versa.
type Scope struct {
	parent *Scope // to enable lexical scope, shadowing and immutability
	fns    map[string]Fn
	vals   map[string]SExpr
}

func (scope Scope) FindVal(sym string) (SExpr, error) {
	for scope := &scope; scope != nil; scope = scope.parent {
		if val, ok := scope.vals[sym]; ok {
			return val, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("symbol %s not found in scope", sym))
}

func (scope Scope) FindFn(sym string) (Fn, error) {
	for scope := &scope; scope != nil; scope = scope.parent {
		if fn, ok := scope.fns[sym]; ok {
			return fn, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("function %s not found in scope", sym))
}
