package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

type ShellCompleter struct {
	readline.AutoCompleter
	stderr io.Writer
}

func NewCompleter(stdout io.Writer, completers ...readline.PrefixCompleterInterface) *ShellCompleter {
	paths := filepath.SplitList(os.Getenv("PATH"))
	for _, path := range paths {
		path = strings.ReplaceAll(path, `\`, `/`)
		s := strings.Split(path, `/`)
		completers = append(completers, readline.PcItem(s[len(s)-1]))
	}
	completer := readline.NewPrefixCompleter(completers...)
	return &ShellCompleter{AutoCompleter: completer, stderr: stdout}
}

func (c *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	sug, l := c.AutoCompleter.Do(line, pos)
	if len(sug) == 0 {
		_, _ = c.stderr.Write([]byte{'\a'})
	}
	return sug, l
}
