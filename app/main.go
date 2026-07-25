package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/chzyer/readline"
)

func main() {
	liner, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		AutoComplete: readline.NewPrefixCompleter(
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
	}
	dispatcher := NewDispatcher(liner)
	ctx := context.Background()
	for {
		if err := dispatcher.Execute(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
		}
	}
}
