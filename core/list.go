package core

import (
	"fmt"
	"iter"
	"strings"
)

// compile-time interface checks
var _ SExpr = new(List)

// List represents one "cons cell" or "pair" from which lists are
// constructed.  Another way of implementing lists would be to use binary
// trees. It's interesting to check which is more convenient and
// performant.
type List struct {
	// first item in pair, traditionally called "car"
	first SExpr
	// second item in pair, traditionally called "cdr"
	second SExpr
	line   uint
	pos    uint
}

func NewList(line, pos uint, els ...SExpr) List {
	var prev List
	for i := len(els) - 1; i >= 0; i-- {
		p := List{
			first:  els[i],
			second: prev,
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

	// get the function to evaluate
	fnSExpr, err := l.First().Eval(scope)
	if err != nil {
		return nil, l.error("", err)
	}
	fn, ok := fnSExpr.(Fn)
	if !ok {
		return nil, l.error(fmt.Sprintf("can not call `%v` as a function", fnSExpr), nil)
	}

	// get the parameters
	args := []SExpr{}

	rest := l.Rest()
	switch v := rest.(type) {
	case List:
		args = append(args, v.Flatten()...)
	default:
		args = append(args, v)
	}
	// fmt.Println(fnSExpr, ":", len(args), args)

	// pass arguments to the Fn (unevaluated)
	result, err := fn.Invoke(scope, args...)
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

func (l List) Rest() SExpr {
	return l.second
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

			if l.Rest() == nil {
				return
			}
			next, ok := l.Rest().(List)
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
		line:       l.line,
		pos:        l.pos,
		wrappedErr: err,
		msg:        msg,
	}
}

type ListEvalError struct {
	line       uint
	pos        uint
	wrappedErr error
	msg        string
}

func (e ListEvalError) Error() string {
	var errBuilder strings.Builder

	if e.line > 0 && e.pos > 0 {
		errBuilder.WriteString(fmt.Sprintf("%d:%d ", e.line, e.pos))
	}

	errBuilder.WriteString("evaluation error: ")

	if e.msg != "" {
		errBuilder.WriteString(e.msg)
		errBuilder.WriteString(": ")
	}

	if e.wrappedErr != nil {
		errBuilder.WriteString(e.wrappedErr.Error())
	}

	return errBuilder.String()
}

func (e ListEvalError) Unwrap() error {
	return e.wrappedErr
}
