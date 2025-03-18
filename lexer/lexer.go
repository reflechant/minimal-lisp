package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Error struct {
	line uint
	pos  uint
	msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("tokenizer error at %d:%d, %s", e.line, e.pos, e.msg)
}

type Token struct {
	Typ  TokenType
	Line uint
	Pos  uint
	Text string
}

type TokenType uint

const (
	Atom TokenType = iota
	LParen
	RParen
	Quote
)

// Tokenize splits the input into recognized tokens and returns them in order.
// It's hand-written and rather ad-hoc but it's good enough for now.
// For a more complex grammar I would definitely make it more generic or use ANTLR
func Tokenize(input io.Reader) ([]Token, error) {
	lineScanner := bufio.NewScanner(input)
	lineScanner.Split(bufio.ScanLines)

	tokens := []Token{}

	var lineIdx uint = 0

	finishAtom := func(b *strings.Builder, pos *int) {
		if *pos >= 0 {
			tokens = append(tokens, Token{
				Typ:  Atom,
				Line: lineIdx,
				Pos:  uint(*pos),
				Text: b.String(),
			})
			b.Reset()
			*pos = -1
		}
	}

	for lineScanner.Scan() {
		line := lineScanner.Text()
		lineIdx++

		// ignore comments
		if strings.HasPrefix(line, ";") {
			continue
		}

		var atomBuf strings.Builder
		pos := -1
		for i, r := range line {
			// if current rune is letter, start/append an the atom
			if unicode.IsLetter(r) {
				atomBuf.WriteRune(r)
				if pos < 0 {
					pos = i + 1
				}
				continue
			}

			// ignore spaces
			if unicode.IsSpace(r) {
				finishAtom(&atomBuf, &pos)
				continue
			}

			if r == '(' {
				finishAtom(&atomBuf, &pos)
				tokens = append(tokens, Token{
					Typ:  LParen,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: "(",
				})
				continue
			}
			if r == ')' {
				finishAtom(&atomBuf, &pos)
				tokens = append(tokens, Token{
					Typ:  RParen,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: ")",
				})
				continue
			}
			if r == '\'' {
				finishAtom(&atomBuf, &pos)
				tokens = append(tokens, Token{
					Typ:  Quote,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: "'",
				})
				continue
			}

			// allow numbers and other punctuation symbols inside atoms
			if unicode.IsDigit(r) || unicode.IsPunct(r) {
				atomBuf.WriteRune(r)
				if pos < 0 {
					pos = i + 1
				}
				continue
			}

			return tokens, Error{
				line: lineIdx,
				pos:  uint(i + 1),
				msg:  fmt.Sprintf("unexpected character %s", string(r)),
			}
		}
		// if reached end of the line and we're still reading an atom, finish it
		if atomBuf.Len() != 0 {
			tokens = append(tokens, Token{
				Typ:  Atom,
				Line: lineIdx,
				Pos:  uint(pos),
				Text: atomBuf.String(),
			})
		}
	}

	return tokens, nil
}
