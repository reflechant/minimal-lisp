package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleLetterAtom(t *testing.T) {
	input := "x"
	expected := []Token{{
		Typ:  Atom,
		Line: 1,
		Pos:  1,
		Text: "x",
	}}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestMultiLetterAtom(t *testing.T) {
	input := "foo"
	expected := []Token{{
		Typ:  Atom,
		Line: 1,
		Pos:  1,
		Text: "foo",
	}}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestEmptyList(t *testing.T) {
	input := "()"
	expected := []Token{
		{
			Typ:  LParen,
			Line: 1,
			Pos:  1,
			Text: "",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  2,
			Text: "",
		},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestSingleAtomList(t *testing.T) {
	input := "(foo)"
	expected := []Token{
		{
			Typ:  LParen,
			Line: 1,
			Pos:  1,
			Text: "",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  2,
			Text: "foo",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  5,
			Text: "",
		},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestMultiAtomList(t *testing.T) {
	input := "(foo bar)"
	expected := []Token{
		{
			Typ:  LParen,
			Line: 1,
			Pos:  1,
			Text: "",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  2,
			Text: "foo",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  6,
			Text: "bar",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  9,
			Text: "",
		},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestNestedEmptyLists(t *testing.T) {
	input := "( ())"
	expected := []Token{
		{
			Typ:  LParen,
			Line: 1,
			Pos:  1,
			Text: "",
		},
		{
			Typ:  LParen,
			Line: 1,
			Pos:  3,
			Text: "",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  4,
			Text: "",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  5,
			Text: "",
		},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestNestedNonEmptyLists(t *testing.T) {
	input := "(foo (bar) baz)"
	expected := []Token{
		{
			Typ:  LParen,
			Line: 1,
			Pos:  1,
			Text: "",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  2,
			Text: "foo",
		},
		{
			Typ:  LParen,
			Line: 1,
			Pos:  6,
			Text: "",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  7,
			Text: "bar",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  10,
			Text: "",
		},
		{
			Typ:  Atom,
			Line: 1,
			Pos:  12,
			Text: "baz",
		},
		{
			Typ:  RParen,
			Line: 1,
			Pos:  15,
			Text: "",
		},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}
