package commands

import (
	"context"
	"strings"
	"unicode"
)

const CommandStorageKey = "coomandstorage"

type CommandLine []string

func NewCommandLine(line string) CommandLine {
	result := make(CommandLine, 0, len(line))
	var quoteOpen bool
	var sb strings.Builder
	for _, r := range line {
		if quoteOpen {
			if isQuote(r) {
				quoteOpen = false
			} else {
				sb.WriteRune(r)
			}
		} else {
			if unicode.IsSpace(r) {
				part := sb.String()
				if strings.TrimSpace(part) != "" {
					result = append(result, sb.String())
				}
				sb.Reset()
			} else if isQuote(r) {
				quoteOpen = true
			} else {
				sb.WriteRune(r)
			}
		}
	}
	return result
}

func (c CommandLine) IsEmpty() bool {
	if len(c) == 0 {
		return true
	}
	if len(c) == 1 && strings.TrimSpace(c[0]) == "" {
		return true
	}
	return false
}

func (c CommandLine) Args() []string {
	if len(c) <= 1 {
		return []string{}
	}
	return c[1:]
}
func (c CommandLine) ArgumentLine() string {
	if len(c) < 1 {
		panic("no arguments")
	}
	return strings.Join(c[1:], " ")
}

func (c CommandLine) Name() string {
	if len(c) == 0 {
		panic("empty command line")
	}
	return c[0]
}

type CommandHandler func(context.Context, CommandLine) error

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
