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

func Tokenize(input io.Reader) ([]Token, error) {
	lineScanner := bufio.NewScanner(input)
	lineScanner.Split(bufio.ScanLines)

	tokens := []Token{}

	var lineIdx uint = 0
	for lineScanner.Scan() {
		line := lineScanner.Text()
		lineIdx++

		var atomBuf strings.Builder
		for i, r := range line {
			// if current rune is letter, append to the atom
			if unicode.IsLetter(r) {
				atomBuf.WriteRune(r)
				continue
			}
			// if current rune is not a letter, finish the atom if any
			if atomBuf.Len() != 0 {
				tokens = append(tokens, Token{typ: Atom, text: atomBuf.String()})
				atomBuf.Reset()
			}

			// ignore spaces
			if unicode.IsSpace(r) {
				continue
			}

			if r == '(' {
				tokens = append(tokens, Token{typ: LParen})
				continue
			}
			if r == ')' {
				tokens = append(tokens, Token{typ: RParen})
				continue
			}
			if r == '\'' {
				tokens = append(tokens, Token{typ: Quote})
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
			tokens = append(tokens, Token{typ: Atom, text: atomBuf.String()})
		}
	}

	return tokens, nil
}

type Token struct {
	typ  TokenType
	text string
}

type TokenType uint

const (
	Atom TokenType = iota
	LParen
	RParen
	Quote
)
