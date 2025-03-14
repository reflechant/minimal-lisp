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
	"strings"
)

// Fn is a universal function type
type Fn func(scope Scope, args ...Expr) (Expr, error)

// Scope stores known functions and values. We use lexical scope. This is a
// so-called LISP-1 - though we store functions and values separately, name
// conflicts are not allowed and one symbol can only be either a function or a
// value. Shadowing is allowed and you can rebind a function symbol to a value and vice
// versa.
type Scope struct {
	parent *Scope // to enable lexical scope, shadowing and immutability
	fns    map[string]Fn
	vals   map[string]Expr
}

func (scope Scope) FindVal(sym string) (Expr, error) {
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

type Expr interface {
	Eval(scope Scope) (Expr, error)
	Print() string
}

type Atom string

func (a Atom) Eval(scope Scope) (Expr, error) {
	e, err := scope.FindVal(string(a))
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (a Atom) Print() string {
	return string(a)
}

type ListElement struct {
	expr Expr
	next *ListElement
}

type List struct {
	head *ListElement
}

func (l List) Eval(scope Scope) (Expr, error) {
	if l.IsEmpty() {
		// empty list evaluates to itself
		return l, nil
	}
	arg1, err := l.First().Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("list evaluation error: %w", err)
	}
	fnSym, ok := arg1.(Atom)
	if !ok {
		return nil, errors.New(fmt.Sprintf("list evaluation error: first argument is not an atom but %v", arg1))
	}
	fn, err := scope.FindFn(string(fnSym))
	if err != nil {
		return nil, fmt.Errorf("list evaluation error: %w", err)
	}
	// eval all arguments
	arguments := []Expr{}
	for arg := l.Rest(); !arg.IsEmpty(); arg = arg.Rest() {
		arguments = append(arguments, arg.First())
	}
	// pass them to the Fn

	return fn(scope, arguments...)
}

func (l List) IsEmpty() bool {
	return l.head == nil
}

func (l List) First() Expr {
	if l.head != nil {
		return l.head.expr
	}

	return nil
}

func (l List) Rest() List {
	if l.head != nil {
		return List{head: l.head.next}
	}

	return List{}
}

func (l List) Cons(e Expr) List {
	return List{
		head: &ListElement{
			expr: e,
			next: l.head,
		},
	}
}

func (l List) Print() string {
	var b strings.Builder
	b.WriteRune('(')

	for el := l.head; el != nil; el = el.next {
		b.WriteString(el.expr.Print())
		b.WriteRune(' ')
	}
	b.WriteRune(')')

	return b.String()
}
