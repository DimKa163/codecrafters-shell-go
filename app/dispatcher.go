package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/commands/lex"
	"os"
	"os/exec"
)

type dipatcher struct {
	cmd    map[string]commands.CommandHandler
	reader bufio.Reader
}

func NewDispatcher() *dipatcher {
	cmd := make(map[string]commands.CommandHandler, 4)
	cmd[commands.EchoCommand] = commands.Echo
	cmd[commands.ExitCommand] = commands.Exit
	cmd[commands.TypeCommand] = commands.Type
	cmd[commands.PwdCommand] = commands.Pwd
	cmd[commands.CdCommand] = commands.Cd
	return &dipatcher{cmd: cmd, reader: *bufio.NewReader(os.Stdin)}
}

func (d *dipatcher) Execute(ctx context.Context) error {
	fmt.Print("$ ")
	line, err := d.reader.ReadString('\n')
	if err != nil {
		return err
	}

	cmd := commands.NewCommandLine(lex.NewLexer(line))
	if cmd.IsEmpty() || cmd.Name() == "" {
		return nil
	}

	h, ok := d.cmd[cmd.Name()]
	if ok {
		return h(context.WithValue(ctx, commands.CommandStorageKey, d.cmd), cmd)
	}
	return d.execExternalProgram(ctx, cmd.Name(), cmd.Args()...)
}

func (d *dipatcher) execExternalProgram(ctx context.Context, name string, args ...string) error {
	var err error
	if _, err = exec.LookPath(name); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%s: not found", name)
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
