package core

import "fmt"

// compile-time interface check
var _ SExpr = new(Symbol)

// Symbol represents itself. It's purely symbolical :)
// For all intents and purposes of this implementation all atoms are symbols.
// Even numbers are not supported :)
type Symbol struct {
	line uint
	pos  uint
	name string
}

func NewSymbol(line, pos uint, val string) Symbol {
	return Symbol{
		line: line,
		pos:  pos,
		name: val,
	}
}

// Eval for an Atom returns it's value
func (s Symbol) Eval(scope Scope) (SExpr, error) {
	// lookup atom among bounded symbols in scope (that includes built-in functions)
	if v, ok := scope.SymbolValue(s.name); ok {
		return v, nil
	}

	return nil, fmt.Errorf("unbound symbol %v", s.name)
}

func (s Symbol) String() string {
	return s.name
}
