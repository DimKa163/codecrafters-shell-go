package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/pkg/terminal"
)

const (
	EchoCommand = "echo"
	ExitCommand = "exit"
	TypeCommand = "type"
	PwdCommand  = "pwd"
	CdCommand   = "cd"
)

type writer struct {
	io.Writer
}

func (w *writer) WriteString(p string) (n int, err error) {
	return w.Write([]byte(p))
}

type shell struct {
	cmd    map[string]terminal.CommandHandler
	Stdout *writer
	Stderr *writer
}

func NewShell() *shell {
	s := &shell{}
	cmd := make(map[string]terminal.CommandHandler, 4)
	cmd[EchoCommand] = s.echo
	cmd[ExitCommand] = s.exit
	cmd[TypeCommand] = s.check
	cmd[PwdCommand] = s.pwd
	cmd[CdCommand] = s.cd
	s.cmd = cmd
	return s
}

func (s *shell) Check(name string) bool {
	_, ok := s.cmd[name]
	return ok
}

func (s *shell) SetStdout(stdout io.Writer) {
	s.Stdout = &writer{stdout}
}

func (s *shell) SetStderr(stderr io.Writer) {
	s.Stderr = &writer{stderr}
}

func (s *shell) Exec(ctx context.Context, key string, args ...string) error {
	return s.cmd[key](ctx, args...)
}

func (s *shell) ExecExternalProgram(ctx context.Context, name string, args ...string) error {
	var err error
	if _, err = exec.LookPath(name); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		if _, err = s.Stderr.WriteString(fmt.Sprintf("%s: not found\n", name)); err != nil {
			return err
		}
		return nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr
	return cmd.Run()
}

func (s *shell) cd(ctx context.Context, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("cd: no args")
	}
	arg := args[0]
	if strings.HasPrefix(arg, "~") {
		hmDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		arg = strings.Replace(arg, "~", hmDir, 1)
	}
	if err := os.Chdir(arg); err != nil {
		if pathErr, ok := errors.AsType[*os.PathError](err); ok && pathErr.Op == "chdir" {
			if _, err = s.Stderr.WriteString(fmt.Sprintf("cd: %s: No such file or directory\n", arg)); err != nil {
				return err
			}
			return fmt.Errorf("cd: %s: No such file or directory", arg)
		}
		return err
	}
	return nil
}

func (s *shell) pwd(ctx context.Context, args ...string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if _, err = os.Stdout.WriteString(fmt.Sprintf("%s\n", dir)); err != nil {
		return err
	}
	return nil
}

func (s *shell) check(ctx context.Context, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("cd: no args")
	}
	_, ok := s.cmd[args[0]]
	if ok {
		if _, err := s.Stdout.WriteString(fmt.Sprintf("%s is a shell builtin\n", args[0])); err != nil {
			return err
		}
		return nil
	}
	var err error
	var path string
	if path, err = exec.LookPath(args[0]); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		if _, err = s.Stderr.WriteString(fmt.Sprintf("%s: not found\n", args[0])); err != nil {
			return err
		}
		return nil
	}
	if _, err = s.Stdout.WriteString(fmt.Sprintf("%s is %s\n", args[0], path)); err != nil {
		return err
	}
	return nil
}

func (s *shell) exit(ctx context.Context, args ...string) error {
	return io.EOF
}

func (s *shell) echo(ctx context.Context, args ...string) error {
	if _, err := s.Stdout.WriteString(fmt.Sprintf("%s\n", strings.Join(args, " "))); err != nil {
		return err
	}
	return nil
}
