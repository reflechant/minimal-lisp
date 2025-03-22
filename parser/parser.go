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

func Parse(srcName string, input io.Reader) ([]core.SExpr, error) {
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", srcName, err)
	}

	exprs := []core.SExpr{}

	i := 0
	for i < len(tokens) {
		expr, next, err := parse(tokens, i)
		if err != nil {
			return exprs, fmt.Errorf("%s: %w", srcName, err)
		}
		// fmt.Printf("expr: %v\n", expr)
		exprs = append(exprs, expr)
		i = next
	}

	return exprs, nil
}

func parse(tokens []lexer.Token, start int) (core.SExpr, int, error) {
	tok := tokens[start]

	switch tok.Typ {
	case lexer.Atom:
		return core.NewSymbol(tok.Line, tok.Pos, tok.Text), start + 1, nil
	case lexer.LParen:
		return parseList(tokens, start) // we start at i so that it can set line and pos for the list
	case lexer.RParen:
		return nil, start + 1, Error{
			line: tok.Line,
			pos:  tok.Pos,
			msg:  fmt.Sprintf("unexpected token %v", tok.Text),
		}
	case lexer.Quote:
		if start == len(tokens)-1 {
			return nil, start + 1, Error{
				line: tok.Line,
				pos:  tok.Pos,
				msg:  fmt.Sprintf("unexpected end of input: quote needs an argument"),
			}
		}
		quotedExpr, next, err := parse(tokens, start+1)
		if err != nil {
			return nil, next, err
		}
		return core.NewList(
			tok.Line,
			tok.Pos,
			core.NewSymbol(tok.Line, tok.Pos, "quote"),
			quotedExpr,
		), next, nil
	default:
		return nil, start + 1, Error{
			line: tok.Line,
			pos:  tok.Pos,
			msg:  fmt.Sprintf("unknown token %v", tok.Text),
		}
	}
}

func parseList(tokens []lexer.Token, start int) (core.SExpr, int, error) {
	// start points to '(' which is guaranteed to exist by the caller
	line, pos := tokens[start].Line, tokens[start].Pos

	if start == len(tokens)-1 {
		// if we're already at the end on input
		return nil, start + 1, Error{
			line: line,
			pos:  pos,
			msg:  fmt.Sprintf("unexpected end of input: list opened at %d:%d was not closed", line, pos),
		}
	}

	items := []core.SExpr{}
	i := start + 1
	for i < len(tokens) {
		tok := tokens[i]
		if tok.Typ == lexer.RParen {
			// end of the list, returning accumulated items
			return core.NewList(line, pos, items...), i + 1, nil
		}

		expr, next, err := parse(tokens, i)
		if err != nil {
			return nil, next, err
		}

		items = append(items, expr)
		i = next
	}

	return nil, i + 1, Error{
		line: tokens[i-1].Line,
		pos:  tokens[i-1].Pos,
		msg:  fmt.Sprintf("unexpected end of input: list opened at %d:%d was not closed", line, pos),
	}
}
