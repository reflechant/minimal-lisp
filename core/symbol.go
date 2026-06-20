package core

import "fmt"

// compile-time interface check
var _ SExpr = new(Symbol)

// Symbol represents itself. It's purely symbolical :)
// For all intents and purposes of this implementation all atoms are symbols.
// Even numbers are not supported :)
type Symbol struct {
	srcName string
	line    uint
	pos     uint
	name    string
}

func NewSymbol(srcName string, line, pos uint, val string) Symbol {
	return Symbol{
		srcName: srcName,
		line:    line,
		pos:     pos,
		name:    val,
	}
}

// Eval for an Atom returns it's value
func (s Symbol) Eval(scope Scope) (SExpr, error) {
	// lookup atom among bounded symbols in scope (that includes built-in functions)
	if v, ok := scope.SymbolValue(s.name); ok {
		return v, nil
	}

	loc := location(s.srcName, s.line, s.pos)
	return nil, fmt.Errorf("%s: unbound symbol %v", loc, s.name)
}

func (s Symbol) String() string {
	return s.name
}
