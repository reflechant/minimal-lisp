package core

import (
	"errors"
	"fmt"
	"strings"
)

var _ SExpr = new(List)

type ListElement struct {
	expr SExpr
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

func NewList(line, pos uint, expr SExpr) List {
	return List{
		head: &ListElement{
			expr: expr,
			next: nil,
		},
		line: line,
		pos:  pos,
	}
}

func NewListFromElements(line, pos uint, els []SExpr) List {
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

// func (l List) Eq(e SExpr) bool {
// 	l2, ok := e.(List)
// 	if !ok {
// 		return false
// 	}

// 	if l.IsEmpty() && l2.IsEmpty() {
// 		return true
// 	}

// 	if l.IsEmpty() != l2.IsEmpty() {
// 		return false
// 	}

// 	if !l.First().Eq(l2.First()) {
// 		return false
// 	}

// 	return l.Rest().Eq(l2.Rest())
// }

func (l List) Eval(scope Scope) (SExpr, error) {
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
	arguments := []SExpr{}
	for arg := l.Rest(); !arg.IsEmpty(); arg = arg.Rest() {
		arguments = append(arguments, arg.First())
	}
	// pass them to the Fn

	result, err := fn(scope, arguments...)
	if err != nil {
		return nil, fmt.Errorf("%d:%d list evaluation error: %w", l.line, l.pos, err)
	}

	return result, nil
}

func (l List) IsEmpty() bool {
	return l.head == nil
}

func (l List) First() SExpr {
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
func (l List) Cons(e SExpr) List {
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
