package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/pkg/terminal"
)

//go:generate mockgen -source=dispatcher.go -destination=mocks/mock_dispatcher.go -package=mocks
type Liner interface {
	ReadLine() string
	Stdout() io.Writer
	Stderr() io.Writer
	Exec(ctx context.Context, key string, args ...string) error
	ExecExternalProgram(ctx context.Context, name string, args ...string) error
	Check(name string) bool
	SetStdout(stdout io.Writer)
	SetStderr(stderr io.Writer)
	ResetOutput()
}
type dipatcher struct {
	liner Liner
}

func NewDispatcher(liner Liner) *dipatcher {
	return &dipatcher{liner: liner}
}

func (d *dipatcher) Execute(ctx context.Context) error {
	var err error
	defer d.liner.ResetOutput()
	out := d.liner.Stdout()
	errOut := d.liner.Stderr()
	line := d.liner.ReadLine()
	if line == "" {
		return nil
	}
	lexer := terminal.NewLexer(line)
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
		case tkn.Type == terminal.TokenTypeName:
			cmdName = tkn.Value
		case tkn.Type == terminal.TokenTypeRedirect:
			setStdoutOverwrite = true
		case tkn.Type == terminal.TokenTypeErrorRedirect:
			setStderroutOverwrite = true
		case tkn.Type == terminal.TokenTypeRedirectAppend:
			setStdoutAppend = true
		case tkn.Type == terminal.TokenTypeErrorRedirectAppend:
			setStderroutAppend = true
		case tkn.Type == terminal.TokenTypeWord:
			args = append(args, tkn.Value)
		}
	}
	d.liner.SetStdout(out)
	d.liner.SetStderr(errOut)
	if d.liner.Check(cmdName) {
		if err = d.liner.Exec(ctx, cmdName, args...); err != nil {
			return err
		}
		return nil
	}

	if err = d.liner.ExecExternalProgram(ctx, cmdName, args...); err != nil {
		return err
	}
	return nil
}
