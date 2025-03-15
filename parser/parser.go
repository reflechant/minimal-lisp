package parser

import (
	"fmt"
	"io"

	"github.com/reflechant/minimal-lisp/core"
	"github.com/reflechant/minimal-lisp/lexer"
)

type Error struct {
	line uint
	pos  uint
	msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("parser error at %d:%d, %s", e.line, e.pos, e.msg)
}

func Parse(input io.Reader) ([]core.Expr, error) {
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}

	exprs := []core.Expr{}

	i := 0
	for i < len(tokens) {
		expr, next, err := parse(tokens, i)
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, expr)
		i = next
	}

	return exprs, nil
}

func parse(tokens []lexer.Token, start int) (core.Expr, int, error) {
	tok := tokens[start]

	switch tok.Typ {
	case lexer.Atom:
		return core.NewAtom(tok.Line, tok.Pos, tok.Text), start + 1, nil
	case lexer.LParen:
		return parseList(tokens, start) // we start at i so that it can set line and pos for the list
	case lexer.RParen:
		return nil, start + 1, Error{
			line: tok.Line,
			pos:  tok.Pos,
			msg:  fmt.Sprintf("unexpected token %v", tok.Text),
		}
	// TODO: support 'x as syntactic sugar over (quote x)
	// case lexer.Quote:
	// exprs = append(exprs, core.Expr)
	default:
		return nil, start + 1, Error{
			line: tok.Line,
			pos:  tok.Pos,
			msg:  fmt.Sprintf("unknown token %v", tok.Text),
		}
	}
}

func parseList(tokens []lexer.Token, start int) (core.Expr, int, error) {
	// start points to '('
	line, pos := tokens[start].Line, tokens[start].Pos

	items := []core.Expr{}
	i := start + 1
	for i < len(tokens) {
		tok := tokens[i]
		switch tok.Typ {
		case lexer.Atom:
			items = append(items, core.NewAtom(tok.Line, tok.Pos, tok.Text))
			i++
		case lexer.LParen:
			lst, next, err := parseList(tokens, i)
			if err != nil {
				return nil, next, err
			}
			items = append(items, lst)
			i = next
		case lexer.RParen:
			return core.NewListFromElements(line, pos, items), i + 1, nil
		// case lexer.Quote:
		// exprs = append(exprs, core.Expr)
		default:
			return nil, i + 1, Error{
				line: tok.Line,
				pos:  tok.Pos,
				msg:  fmt.Sprintf("unknown token %v", tok.Text),
			}
		}
	}

	// return core.NewListFromElements(line, pos, items), i + 1, nil
	return nil, i + 1, Error{
		line: tokens[i].Line,
		pos:  tokens[i].Pos,
		msg:  fmt.Sprintf("list opened at %d:%d was not closed", line, pos),
	}
}
