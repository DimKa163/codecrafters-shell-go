package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type CommandLine []string

func NewCommandLine(line string) CommandLine {
	return strings.Split(strings.Trim(line, "\n"), " ")
}

func (c CommandLine) IsEmpty() bool {
	if len(c) == 0 {
		return true
	}
	if len(c) == 1 && c[0] == "" {
		return true
	}
	return false
}
func (c CommandLine) ArgumentLine() string {
	if len(c) < 1 {
		panic("no arguments")
	}
	return strings.Join(c[1:], " ")
}

func (c CommandLine) Name() string {
	if len(c) == 0 {
		panic("empty command line")
	}
	return c[0]
}

type CommandHandler func(CommandLine) error

const (
	ExitCommand = "exit"
	EchoCommand = "echo"
	TypeCommand = "type"
)

func main() {
	commandMap := initCommandMap()
	reader := bufio.NewReader(os.Stdin)
	var cmd CommandLine
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				_, _ = fmt.Fprintln(os.Stderr, "Error reading input:", err)
				return
			}
			return
		}
		cmd = NewCommandLine(line)
		if cmd.IsEmpty() {
			continue
		}
		h, ok := commandMap[cmd.Name()]
		if !ok {
			fmt.Printf("%s: command not found\n", cmd.Name())
		}
		if err = h(cmd); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println(err)
		}
	}
}

func initCommandMap() map[string]CommandHandler {
	m := make(map[string]CommandHandler)
	m[EchoCommand] = func(cmd CommandLine) error {
		fmt.Println(cmd.ArgumentLine())
		return nil
	}
	m[ExitCommand] = func(cmd CommandLine) error {
		return io.EOF
	}
	m[TypeCommand] = func(cmd CommandLine) error {
		_, ok := m[cmd.ArgumentLine()]
		if !ok {
			return fmt.Errorf("%s: not found", cmd.ArgumentLine())
		}
		fmt.Printf("%s is a shell builtin\n", cmd.ArgumentLine())
		return nil
	}
	return m
}
