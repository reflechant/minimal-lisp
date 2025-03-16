package core

import "fmt"

// compile-time interface check
var _ SExpr = new(Atom)

type Atom struct {
	line uint
	pos  uint
	text string
}

func NewAtom(line, pos uint, text string) *Atom {
	return &Atom{
		line: line,
		pos:  pos,
		text: text,
	}
}

// func (a *Atom) Eq(e SExpr) bool {
// 	a2, ok := e.(*Atom)
// 	if !ok {
// 		return false
// 	}

// 	return a.text == a2.text
// }

func (a *Atom) Eval(scope Scope) (SExpr, error) {
	e, err := scope.FindVal(a.text)
	if err != nil {
		return nil, fmt.Errorf("%d:%d: %w", a.line, a.pos, err)
	}

	return e, nil
}

func (a *Atom) Print() string {
	return string(a.text)
}
