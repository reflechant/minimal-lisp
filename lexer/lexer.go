package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Error struct {
	srcName string
	line    uint
	pos     uint
	msg     string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%d:%d: lex error: %s", e.srcName, e.line, e.pos, e.msg)
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
// srcName is used in error messages (file path, "repl", etc.).
// lineOffset is added to every line number, enabling the REPL to report
// cumulative line numbers across multiple inputs.
func Tokenize(srcName string, lineOffset uint, input io.Reader) ([]Token, error) {
	lineScanner := bufio.NewScanner(input)
	lineScanner.Split(bufio.ScanLines)

	tokens := []Token{}

	var lineIdx uint = 0

	finishAtom := func(b *strings.Builder, pos *int) {
		if *pos >= 0 {
			tokens = append(tokens, Token{
				Typ:  Atom,
				Line: lineIdx + lineOffset,
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
					Line: lineIdx + lineOffset,
					Pos:  uint(i + 1),
					Text: "(",
				})
				continue
			}
			if r == ')' {
				finishAtom(&atomBuf, &pos)
				tokens = append(tokens, Token{
					Typ:  RParen,
					Line: lineIdx + lineOffset,
					Pos:  uint(i + 1),
					Text: ")",
				})
				continue
			}
			if r == '\'' {
				finishAtom(&atomBuf, &pos)
				tokens = append(tokens, Token{
					Typ:  Quote,
					Line: lineIdx + lineOffset,
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
				srcName: srcName,
				line:    lineIdx + lineOffset,
				pos:     uint(i + 1),
				msg:     fmt.Sprintf("unexpected character %s", string(r)),
			}
		}
		// if reached end of the line and we're still reading an atom, finish it
		if atomBuf.Len() != 0 {
			tokens = append(tokens, Token{
				Typ:  Atom,
				Line: lineIdx + lineOffset,
				Pos:  uint(pos),
				Text: atomBuf.String(),
			})
		}
	}

	return tokens, nil
}
