package main

import (
	"io"

	"github.com/chzyer/readline"
)

type ShellCompleter struct {
	readline.AutoCompleter
	stderr io.Writer
}

func NewCompleter(completer readline.AutoCompleter, stdout io.Writer) *ShellCompleter {
	return &ShellCompleter{AutoCompleter: completer, stderr: stdout}
}

func (c *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	sug, l := c.AutoCompleter.Do(line, pos)
	if len(sug) == 0 {
		_, _ = c.stderr.Write([]byte{'\a'})
	}
	return sug, l
}
