package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleLetterAtom(t *testing.T) {
	input := "x"
	expected := []Token{{typ: Atom, text: "x"}}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestMultiLetterAtom(t *testing.T) {
	input := "foo"
	expected := []Token{{typ: Atom, text: "foo"}}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestEmptyList(t *testing.T) {
	input := "()"
	expected := []Token{
		{typ: LParen},
		{typ: RParen},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestSingleAtomList(t *testing.T) {
	input := "(foo)"
	expected := []Token{
		{typ: LParen},
		{typ: Atom, text: "foo"},
		{typ: RParen},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestMultiAtomList(t *testing.T) {
	input := "(foo bar)"
	expected := []Token{
		{typ: LParen},
		{typ: Atom, text: "foo"},
		{typ: Atom, text: "bar"},
		{typ: RParen},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestNestedEmptyLists(t *testing.T) {
	input := "( ())"
	expected := []Token{
		{typ: LParen},
		{typ: LParen},
		{typ: RParen},
		{typ: RParen},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}

func TestNestedNonEmptyLists(t *testing.T) {
	input := "(foo (bar) baz)"
	expected := []Token{
		{typ: LParen},
		{typ: Atom, text: "foo"},
		{typ: LParen},
		{typ: Atom, text: "bar"},
		{typ: RParen},
		{typ: Atom, text: "baz"},
		{typ: RParen},
	}
	tokens, err := Tokenize(strings.NewReader(input))
	require.NoError(t, err)
	assert.Equal(t, expected, tokens)
}
