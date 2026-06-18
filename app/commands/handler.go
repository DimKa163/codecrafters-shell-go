package commands

import (
	"context"
	"strings"
)

const CommandStorageKey = "coomandstorage"

type CommandLine []string

func NewCommandLine(line string) CommandLine {
	return strings.Split(strings.Trim(line, "\n"), " ")
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
