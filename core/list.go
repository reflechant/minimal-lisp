package core

import (
	"fmt"
	"iter"
	"strconv"
	"strings"
)

// compile-time interface checks
var _ SExpr = new(List)

// List represents a list S-expression as one "cons cell" or "pair".
// Another way of implementing lists would be to use binary trees.
// It's interesting to check which is more convenient and performant.
type List struct {
	// first item in pair, traditionally called "car"
	first SExpr
	// second item in pair, traditionally called "cdr"
	second  SExpr
	srcName string
	line    uint
	pos     uint
}

func NewList(srcName string, line, pos uint, els ...SExpr) List {
	var prev List
	for i := len(els) - 1; i >= 0; i-- {
		p := List{
			first:   els[i],
			second:  prev,
			srcName: srcName,
			line:    line,
			pos:     pos,
		}
		prev = p
	}

	return prev
}

// Eval returns the value of a list S-expression.
// Usually a form like `(f a b c)` is called a "function call" but depending
// on what f is different rules may apply (specifically for lambda, label and defun).
// That's why we let the function to evaluate the arguments because in some cases
// they shouldn't be evaluated (e.g. parameter list for lambda)
// See "The Roots of LISP" for details.
func (l List) Eval(scope Scope) (SExpr, error) {
	// fmt.Println(l)
	if l.IsEmpty() {
		// empty list evaluates to itself
		return l, nil
	}

	items := l.Flatten()

	// get the function to evaluate
	fnSExpr, err := items[0].Eval(scope)
	if err != nil {
		return nil, l.error("", err)
	}
	fn, ok := fnSExpr.(Fn)
	if !ok {
		return nil, l.error(fmt.Sprintf("can not call `%v` as a function", fnSExpr), nil)
	}

	// pass arguments to the Fn (unevaluated)
	result, err := fn.Invoke(scope, items[1:]...)
	if err != nil {
		return nil, l.error("", err)
	}

	return result, nil
}

func (l List) IsEmpty() bool {
	return l.first == nil
}

func (l List) First() SExpr {
	return l.first
}

func (l List) Second() SExpr {
	return l.second
}

// Rest returns an iterator of all list elements after the first
func (l List) Rest() iter.Seq[SExpr] {
	if l.second == nil {
		return func(yield func(SExpr) bool) {}
	}
	if v, ok := l.second.(List); ok {
		return v.Items()
	}

	return func(yield func(SExpr) bool) {
		if !yield(l.second) {
			return
		}
	}
}

func (l List) Cons(v SExpr) List {
	return List{
		first:  v,
		second: l,
	}
}

func (l List) Items() iter.Seq[SExpr] {
	return func(yield func(SExpr) bool) {
		for !l.IsEmpty() {
			if !yield(l.First()) {
				return
			}

			if l.second == nil {
				return
			}
			next, ok := l.second.(List)
			if !ok {
				if !yield(next) {
					return
				}
			}
			l = next
		}
	}
}

func (l List) Flatten() []SExpr {
	items := []SExpr{}
	for item := range l.Items() {
		items = append(items, item)
	}

	return items
}

func (l List) String() string {
	var b strings.Builder
	b.WriteRune('(')

	itemStrs := []string{}
	for item := range l.Items() {
		itemStrs = append(itemStrs, item.String())
	}
	b.WriteString(strings.Join(itemStrs, " "))

	b.WriteRune(')')

	return b.String()
}

func (l List) error(msg string, err error) ListEvalError {
	return ListEvalError{
		srcName:    l.srcName,
		line:       l.line,
		pos:        l.pos,
		wrappedErr: err,
		msg:        msg,
	}
}

type ListEvalError struct {
	srcName    string
	line       uint
	pos        uint
	wrappedErr error
	msg        string
}

func (e ListEvalError) Error() string {
	var b strings.Builder

	b.WriteString("evaluation error at ")
	b.WriteString(e.srcName)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(int(e.line)))
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(int(e.pos)))

	if e.msg != "" {
		b.WriteString(": ")
		b.WriteString(e.msg)
	}

	if e.wrappedErr != nil {
		b.WriteString(": ")
		b.WriteString(e.wrappedErr.Error())
	}

	return b.String()
}

func (e ListEvalError) Unwrap() error {
	return e.wrappedErr
}
