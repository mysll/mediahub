package utils

import (
	"io"
	"regexp"
)

var (
	SplitChars = regexp.MustCompile(`\.|\s+|\(|\)|\[|]|-|\+|【|】|/|～|;|&|\||#|_|「|」|~`)
)

type Tokenizer struct {
	text   string
	tokens []string
	index  int
}

func NewToken(t string) *Tokenizer {
	tokens := SplitChars.Split(t, -1)
	return &Tokenizer{
		text:   t,
		tokens: tokens,
	}
}

func (t *Tokenizer) Next() (string, error) {
	if t.index < len(t.tokens)-1 {
		t.index++
		return t.tokens[t.index], nil
	}
	return "", io.EOF
}

func (t *Tokenizer) Peek() (string, error) {
	if t.index < len(t.tokens)-1 {
		return t.tokens[t.index+1], nil
	}
	return "", io.EOF
}

func (t *Tokenizer) Cur() (string, error) {
	if t.index < len(t.tokens) {
		return t.tokens[t.index], nil
	}
	return "", io.EOF
}

func (t *Tokenizer) Reset() {
	t.index = 0
}
