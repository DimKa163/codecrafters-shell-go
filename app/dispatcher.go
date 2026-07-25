package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/commands/lex"
)

//go:generate mockgen -source=dispatcher.go -destination=mocks/mock_dispatcher.go -package=mocks
type Liner interface {
	Readline() (string, error)
	Stdout() io.Writer
	Stderr() io.Writer
}

type Shell interface {
	Exec(ctx context.Context, key string, args ...string) error
	ExecExternalProgram(ctx context.Context, name string, args ...string) error
	Check(name string) bool
	SetStdout(stdout io.Writer)
	SetStderr(stderr io.Writer)
}
type dipatcher struct {
	shell Shell
	liner Liner
}

func NewDispatcher(liner Liner) *dipatcher {
	return &dipatcher{shell: commands.NewShell(), liner: liner}
}

func (d *dipatcher) Execute(ctx context.Context) error {
	out := d.liner.Stdout()
	errOut := d.liner.Stderr()
	fmt.Print("$ ")
	line, err := d.liner.Readline()
	if err != nil {
		return err
	}
	lexer := lex.NewLexer(line)
	var cmdName string
	args := make([]string, 0)

	var setStdoutOverwrite bool
	var setStdoutAppend bool
	var setStderroutOverwrite bool
	var setStderroutAppend bool
	for tkn := range lexer.All() {
		switch {
		case setStdoutOverwrite:
			if err = os.MkdirAll(filepath.Dir(tkn.Value), 0755); err != nil {
				return err
			}
			out, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			setStdoutOverwrite = false
		case setStdoutAppend:
			if err = os.MkdirAll(filepath.Dir(tkn.Value), 0755); err != nil {
				return err
			}
			out, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			setStdoutAppend = false
		case setStderroutOverwrite:
			errOut, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			setStderroutOverwrite = false
		case setStderroutAppend:
			errOut, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			setStderroutAppend = false
		case tkn.Type == lex.TokenTypeName:
			cmdName = tkn.Value
		case tkn.Type == lex.TokenTypeRedirect:
			setStdoutOverwrite = true
		case tkn.Type == lex.TokenTypeErrorRedirect:
			setStderroutOverwrite = true
		case tkn.Type == lex.TokenTypeRedirectAppend:
			setStdoutAppend = true
		case tkn.Type == lex.TokenTypeErrorRedirectAppend:
			setStderroutAppend = true
		case tkn.Type == lex.TokenTypeWord:
			args = append(args, tkn.Value)
		}
	}
	d.shell.SetStdout(out)
	d.shell.SetStderr(errOut)
	if d.shell.Check(cmdName) {
		if err = d.shell.Exec(ctx, cmdName, args...); err != nil {
			return err
		}
		return nil
	}

	if err = d.shell.ExecExternalProgram(ctx, cmdName, args...); err != nil {
		return err
	}
	return nil
}
