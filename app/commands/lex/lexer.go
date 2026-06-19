package lex

import (
	"iter"
	"strings"
	"unicode"
)

type Token struct {
	Value string
}

type State interface {
	Handle(l *Lexer, r rune)
	Reset(l *Lexer)
}
type Lexer struct {
	input []rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		var sb strings.Builder
		q := false
		for _, r := range l.input {
			if q {
				if isQuote(r) {
					q = false
				} else {
					sb.WriteRune(r)
				}
			} else {
				if unicode.IsSpace(r) {
					p := sb.String()
					if strings.TrimSpace(p) != "" {
						if !yield(p) {
							continue
						}
						sb.Reset()
					}
				} else if isQuote(r) {
					q = true
				} else {
					sb.WriteRune(r)
				}
			}
		}
	}
}

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
