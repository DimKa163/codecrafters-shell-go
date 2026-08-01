package terminal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	EchoCommand = "echo"
	ExitCommand = "exit"
	TypeCommand = "type"
	PwdCommand  = "pwd"
	CdCommand   = "cd"
)

type Shell struct {
	cfg          *Cfg
	cmd          map[string]CommandHandler
	t            *Terminal
	wg           sync.WaitGroup
	outchan      chan string
	buf          *RuneBuffer
	inSelectMode bool
	candidates   []*CompletableItem
}

func NewShell(cfg *Cfg, t *Terminal) *Shell {
	s := &Shell{cfg: cfg, t: t, outchan: make(chan string), buf: NewBuffer(cfg)}
	cmd := make(map[string]CommandHandler, 4)
	cmd[EchoCommand] = s.echo
	cmd[ExitCommand] = s.exit
	cmd[TypeCommand] = s.check
	cmd[PwdCommand] = s.pwd
	cmd[CdCommand] = s.cd
	s.cmd = cmd
	go s.exec()
	return s
}

func (s *Shell) Stdout() io.Writer {
	return s.cfg.Stdout
}

func (s *Shell) Stderr() io.Writer {
	return s.cfg.Stderr
}

func (s *Shell) exec() {
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		close(s.outchan)
	}()

	for {
		r := s.t.ReadRune()
		if r == 0 {
			continue
		}

		switch r {
		case CharTab:
			if s.cfg.Completer == nil || len(s.buf.buf) == 0 {
				s.t.Bell()
				continue
			}
			if len(s.candidates) == 0 {
				s.Complete()
			} else {
				s.PrintSuggestions()
				s.candidates = nil
			}
		case CharEnter, CharCtrlJ:
			s.t.WriteString("\r\n")
			data := s.buf.Reset()
			s.outchan <- string(data)
		case MetaUp:
		case MetaDown:
		case MetaBackward:
			s.buf.MoveLeft()
		case MetaForward:
			s.buf.MoveRight()
		case CharBackspace:
			s.buf.Backspace()
		default:
			s.buf.WriteRune(r)
		}
	}
}

func (s *Shell) Complete() {
	runes := s.buf.Runes()
	candidates, _ := s.cfg.Completer.Do(runes, s.buf.idx)
	if len(candidates) == 0 {
		s.t.Bell()
		return
	}
	if len(candidates) == 1 {
		c := candidates[0]
		a := append(c.value[len(runes):], ' ')
		s.buf.WriteRunes(a)
	}
	if len(candidates) > 1 {
		prefix := commonPrefix(candidates)
		if len(prefix) > len(runes) {
			s.buf.WriteRunes(prefix[len(runes):])
			return
		}
		s.t.Bell()
		s.candidates = candidates
	}
}

func (s *Shell) SetStdout(stdout io.Writer) {
	s.t.SetStdout(stdout)
}

func (s *Shell) SetStderr(stderr io.Writer) {
	s.t.SetStderr(stderr)
}

func (s *Shell) ResetOutput() {
	s.t.Reset()
}

func commonPrefix(candidates []*CompletableItem) []rune {
	if len(candidates) == 0 {
		return nil
	}
	prefix := append([]rune{}, candidates[0].value...)
	for _, c := range candidates[1:] {
		for len(prefix) > 0 && !HasPrefix(prefix, c.value) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func (s *Shell) PrintSuggestions() {
	data := s.buf.Reset()
	suggestions := make([]string, 0, len(s.candidates))
	for _, i := range s.candidates {
		suggestions = append(suggestions, strings.TrimSpace(string(i.value)))
	}
	s.t.WriteString("\r\n")
	s.t.WriteStringLine(strings.Join(suggestions, "  "))
	s.buf.WriteRunes(data)
	//s.buf.Refresh(nil)
}

func (s *Shell) ReadLine() string {
	s.buf.Refresh(nil)
	l := <-s.outchan
	return l
}
func (s *Shell) Check(name string) bool {
	_, ok := s.cmd[name]
	return ok
}
func (s *Shell) Exec(ctx context.Context, key string, args ...string) error {
	return s.cmd[key](ctx, args...)
}

func (s *Shell) ExecExternalProgram(ctx context.Context, name string, args ...string) error {
	var err error
	if _, err = exec.LookPath(name); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		s.t.WriteErrorStringLine(fmt.Sprintf("%s: not found", name))
		return nil
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = s.cfg.Stdout
	cmd.Stderr = s.cfg.Stderr
	return cmd.Run()
}

func (s *Shell) WriteLine(data ...string) {
	s.t.WriteStringLine(strings.Join(data, ""))
	s.buf.Refresh(nil)
}

func (s *Shell) cd(ctx context.Context, args ...string) error {
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
			s.t.WriteErrorStringLine(fmt.Sprintf("cd: %s: No such file or directory", arg))
			return fmt.Errorf("cd: %s: No such file or directory", arg)
		}
		return err
	}
	return nil
}

func (s *Shell) pwd(ctx context.Context, args ...string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	s.t.WriteStringLine(fmt.Sprintf("%s", dir))
	return nil
}

func (s *Shell) check(ctx context.Context, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("cd: no args")
	}
	_, ok := s.cmd[args[0]]
	if ok {
		s.t.WriteStringLine(fmt.Sprintf("%s is a shell builtin", args[0]))
		return nil
	}
	var err error
	var path string
	if path, err = exec.LookPath(args[0]); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		s.t.WriteErrorStringLine(fmt.Sprintf("%s: not found", args[0]))
		return nil
	}
	s.t.WriteStringLine(fmt.Sprintf("%s is %s", args[0], path))
	return nil
}

func (s *Shell) exit(ctx context.Context, args ...string) error {
	return io.EOF
}

func (s *Shell) echo(ctx context.Context, args ...string) error {
	s.t.WriteStringLine(fmt.Sprintf("%s", strings.Join(args, " ")))
	return nil
}
