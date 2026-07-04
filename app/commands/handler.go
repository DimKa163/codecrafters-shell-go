package commands

import (
	"context"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands/lex"
)

const CommandStorageKey = "coomandstorage"

type CommandLine []string
type Command struct {
	command string
	line    CommandLine
}

func NewCommandLine(lexer *lex.Lexer) CommandLine {
	result := make(CommandLine, 0)
	for p := range lexer.All() {
		result = append(result, p.Value)
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

type InternalCommand interface {
	Run() error
}
type (
	CommandHandler func(context.Context, ...string) error
)

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
