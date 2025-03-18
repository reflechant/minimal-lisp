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

// SExpr represents a S-expression (atom or a list).
// Go doesn't have union types (well, it does with generics
// but they are still useless since you can't have a slice of them)
type SExpr interface {
	// Eval is the only mandatory operation for an S-expression
	// line and pos are needed for helpful errors
	Eval(scope Scope) (SExpr, error)
	// String returns a textual representation of a value (P in REPL)
	String() string
}

// Scope stores values (functions are values) bound to names(aka symbols).
// Scope is lexical. This is a so-called LISP-1 - name conflicts are not allowed
// and one symbol can only be either a function or a value.
// Shadowing is allowed and you can rebind a function symbol to a value and vice versa.
type Scope struct {
	parent *Scope // to enable lexical scope, shadowing and immutability
	vals   map[string]SExpr
}

func (scope Scope) NewLayer() Scope {
	return Scope{
		parent: &scope,
		vals:   map[string]SExpr{},
	}
}

func (scope Scope) Bind(s string, v SExpr) {
	scope.vals[s] = v
}

func (scope Scope) SymbolValue(sym string) (SExpr, bool) {
	for scope := &scope; scope != nil; scope = scope.parent {
		if val, ok := scope.vals[sym]; ok {
			return val, true
		}
	}

	return nil, false
}
