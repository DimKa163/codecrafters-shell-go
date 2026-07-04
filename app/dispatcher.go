package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/commands/lex"
)

type Shell interface {
	Exec(ctx context.Context, key string, args ...string) error
	ExecExternalProgram(ctx context.Context, name string, args ...string) error
	Check(name string) bool
	SetStdout(stdout *os.File)
	SetStderr(stderr *os.File)
}
type dipatcher struct {
	shell  Shell
	reader bufio.Reader
}

func NewDispatcher() *dipatcher {
	return &dipatcher{shell: commands.NewShell(), reader: *bufio.NewReader(os.Stdin)}
}

func (d *dipatcher) Execute(ctx context.Context) error {
	out := os.Stdout
	errOut := os.Stderr
	var closeOut bool
	var closeErr bool
	defer func() {
		if closeOut {
			out.Close()
		}
		if closeErr {
			errOut.Close()
		}
	}()
	fmt.Print("$ ")
	line, err := d.reader.ReadString('\n')
	if err != nil {
		return err
	}
	lexer := lex.NewLexer(line)
	var cmdName string
	args := make([]string, 0)

	var setStdout bool
	var setStderrout bool
	for tkn := range lexer.All() {
		switch {
		case setStdout:
			if err = os.MkdirAll(filepath.Dir(tkn.Value), 0755); err != nil {
				return err
			}
			out, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			setStdout = false
			closeOut = true
		case setStderrout:
			errOut, err = os.OpenFile(tkn.Value, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			setStderrout = false
			closeErr = true
		case tkn.Type == lex.TokenTypeName:
			cmdName = tkn.Value
		case tkn.Type == lex.TokenTypeRedirect:
			setStdout = true
		case tkn.Type == lex.TokenTypeErrorRedirect:
			setStderrout = true
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

	}
	return nil
}
