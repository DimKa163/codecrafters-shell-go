package lex

import (
	"iter"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenTypeWord TokenType = iota
	TokenTypeExpression
)

type Token struct {
	Value string
	Type  TokenType
}
type Lexer struct {
	input []rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) All() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var sb strings.Builder
		q := false
		dq := false
		for _, r := range l.input {
			if q {
				if isQuote(r) {
					q = false
				} else {
					sb.WriteRune(r)
				}
			} else if dq {
				if isDoubleQuote(r) {
					dq = false
				} else {
					sb.WriteRune(r)
				}
			} else {
				if unicode.IsSpace(r) {
					p := sb.String()
					if strings.TrimSpace(p) != "" {
						if !yield(Token{Value: p, Type: TokenTypeWord}) {
							continue
						}
						sb.Reset()
					}
				} else if isQuote(r) {
					q = true
				} else if isDoubleQuote(r) {
					dq = true
				} else {
					sb.WriteRune(r)
				}
			}
		}
	}
}

func isQuote(r rune) bool {
	return r == '\''
}

func isDoubleQuote(r rune) bool {
	return r == '"'
}
