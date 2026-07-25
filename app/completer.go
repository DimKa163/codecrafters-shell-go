package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
)

type ShellCompleter struct {
	readline.AutoCompleter
	stderr io.Writer
}

func NewCompleter(stdout io.Writer, completers ...readline.PrefixCompleterInterface) *ShellCompleter {
	paths := filepath.SplitList(os.Getenv("PATH"))
	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if info.Mode().Perm()&0111 == 0 {
				continue
			}

			completers = append(
				completers,
				readline.PcItem(e.Name()),
			)
		}
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
