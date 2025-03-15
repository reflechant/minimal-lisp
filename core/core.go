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
	Line() uint
	Pos() uint
	Eq(e Expr) bool
}

// interface checks
var _ Expr = List{}
var _ Expr = Atom{}

type Atom struct {
	line uint
	pos  uint
	text string
}

func NewAtom(line, pos uint, text string) Atom {
	return Atom{
		line: line,
		pos:  pos,
		text: text,
	}
}

func (a Atom) Line() uint { return a.line }

func (a Atom) Pos() uint { return a.pos }

func (a Atom) Eq(e Expr) bool {
	a2, ok := e.(Atom)
	if !ok {
		return false
	}

	return a.text == a2.text
}

func (a Atom) Eval(scope Scope) (Expr, error) {
	e, err := scope.FindVal(a.text)
	if err != nil {
		return nil, fmt.Errorf("%d:%d: %w", a.line, a.pos, err)
	}

	return e, nil
}

func (a Atom) Print() string {
	return string(a.text)
}

type ListElement struct {
	expr Expr
	next *ListElement
}

type List struct {
	head *ListElement
	line uint
	pos  uint
}

func NewEmptyList(line, pos uint) List {
	return List{
		head: nil,
		line: line,
		pos:  pos,
	}
}

func NewList(line, pos uint, expr Expr) List {
	return List{
		head: &ListElement{
			expr: expr,
			next: nil,
		},
		line: line,
		pos:  pos,
	}
}

func NewListFromElements(line, pos uint, els []Expr) List {
	lst := NewEmptyList(line, pos)

	var el *ListElement = nil
	// going in reverse direction because it's more suitable for a singly linked list
	for i := len(els) - 1; i >= 0; i-- {
		newEl := &ListElement{
			expr: els[i],
			next: el,
		}
		el = newEl
	}

	lst.head = el

	return lst
}

func (l List) Line() uint { return l.line }

func (l List) Pos() uint { return l.pos }

func (l List) Eq(e Expr) bool {
	l2, ok := e.(List)
	if !ok {
		return false
	}

	if l.IsEmpty() && l2.IsEmpty() {
		return true
	}

	if l.IsEmpty() != l2.IsEmpty() {
		return false
	}

	if !l.First().Eq(l2.First()) {
		return false
	}

	return l.Rest().Eq(l2.Rest())
}

func (l List) Eval(scope Scope) (Expr, error) {
	if l.IsEmpty() {
		// empty list evaluates to itself
		return l, nil
	}
	arg1, err := l.First().Eval(scope)
	if err != nil {
		return nil, fmt.Errorf("%d:%d list evaluation error: %w", l.line, l.pos, err)
	}
	fnSym, ok := arg1.(Atom)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%d:%d list evaluation error: first argument is not an atom but %v", l.line, l.pos, arg1))
	}
	fn, err := scope.FindFn(fnSym.text)
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

// Cons adds to the front
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

	els := []string{}
	for el := l.head; el != nil; el = el.next {
		els = append(els, el.expr.Print())
	}
	b.WriteString(strings.Join(els, " "))

	b.WriteRune(')')

	return b.String()
}
