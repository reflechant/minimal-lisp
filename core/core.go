// Package lisp implements the core LISP interpreter
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

import "strings"

// Fn is a universal function type
type Fn func(scope Scope, args ...Expr) (Expr, error)

// Scope stores known functions and values. We use lexical scope. This is a
// so-called LISP-1 - though we store functions and values separately, name
// conflicts are not allowed and one symbol can only be either a function or a
// value. Shadowing is allowed and you can reassign a function to a value and vice
// versa.
type Scope struct {
	parent *Scope // to enable lexical scope, shadowing and immutability
	fns    map[string]Fn
	vals   map[string]Expr
}

type Expr interface {
	Eval(scope Scope) (Expr, error)
	Print() string
}

type Atom string

func (a Atom) Eval(scope Scope) (Expr, error) {
	return a, nil
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
	return nil, nil
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
