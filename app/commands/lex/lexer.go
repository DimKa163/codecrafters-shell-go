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
	runes []rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{runes: []rune(input)}
}

func (l *Lexer) All() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var sb strings.Builder
		var quote rune
		var treatAsLiteral bool
		for i := 0; i < len(l.runes); i++ {
			ch := l.runes[i]
			switch {
			case treatAsLiteral:
				sb.WriteRune(ch)
				treatAsLiteral = false
			case (isSlash(ch) && quote == 0) || (isSlash(ch) && isDoubleQuote(quote)):
				treatAsLiteral = true
			case quote != 0:
				if ch == quote {
					quote = 0
					continue
				}
				sb.WriteRune(ch)
			case unicode.IsSpace(ch):
				if sb.Len() > 0 {
					if !yield(Token{Value: sb.String(), Type: TokenTypeWord}) {
						return
					}
					sb.Reset()
				}
			case isQuote(ch):
				quote = '\''
			case isDoubleQuote(ch):
				quote = '"'
			default:
				sb.WriteRune(ch)
			}
		}
		if sb.Len() > 0 {
			yield(Token{Value: sb.String(), Type: TokenTypeWord})
		}
	}
}

func isQuote(r rune) bool {
	return r == '\''
}

func isDoubleQuote(r rune) bool {
	return r == '"'
}

func isSlash(r rune) bool {
	return r == '\\'
}
