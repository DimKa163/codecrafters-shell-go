package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/pkg/terminal"
)

func main() {
	/*liner, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		AutoComplete: NewCompleter(
			os.Stderr,
			readline.PcItem("exit"),
			readline.PcItem("echo"),
		),
		EOFPrompt: "exit",
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
	})
	if err != nil {
		panic(err)
	}*/
	completers := make([]*terminal.CompletableItem, 0)
	completers = append(completers, terminal.NewItem("exit"))
	completers = append(completers, terminal.NewItem("echo"))
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
				terminal.NewItem(e.Name()),
			)
		}
	}
	completer := terminal.NewCompleter(completers...)
	cfg := terminal.Cfg{
		Prompt:    []rune("$ "),
		Completer: completer,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
	}
	t := terminal.NewShell(cfg, terminal.New(cfg))
	dispatcher := NewDispatcher(t)
	ctx := context.Background()
	for {
		if err := dispatcher.Execute(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
		}
	}
}
