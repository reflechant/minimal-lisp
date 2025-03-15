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

func Tokenize(input io.Reader) ([]Token, error) {
	lineScanner := bufio.NewScanner(input)
	lineScanner.Split(bufio.ScanLines)

	tokens := []Token{}

	var lineIdx uint = 0
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
			// if current rune is letter, append to the atom
			if unicode.IsLetter(r) {
				atomBuf.WriteRune(r)
				if pos < 0 {
					pos = i + 1
				}
				continue
			}
			// if current rune is not a letter, finish the atom if any
			if atomBuf.Len() != 0 {
				tokens = append(tokens, Token{
					Typ:  Atom,
					Line: lineIdx,
					Pos:  uint(pos),
					Text: atomBuf.String(),
				})
				atomBuf.Reset()
				pos = -1
			}

			// ignore spaces
			if unicode.IsSpace(r) {
				continue
			}

			if r == '(' {
				tokens = append(tokens, Token{
					Typ:  LParen,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: "(",
				})
				continue
			}
			if r == ')' {
				tokens = append(tokens, Token{
					Typ:  RParen,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: ")",
				})
				continue
			}
			if r == '\'' {
				tokens = append(tokens, Token{
					Typ:  Quote,
					Line: lineIdx,
					Pos:  uint(i + 1),
					Text: "'",
				})
				continue
			}

			return tokens, Error{
				line: lineIdx,
				pos:  uint(i + 1),
				msg:  fmt.Sprintf("unexpected symbol %s", string(r)),
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
