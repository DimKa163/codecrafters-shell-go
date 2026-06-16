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

const (
	ExitCommand = "exit"
	EchoCommand = "echo"
)

func main() {
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
		switch cmd.Name() {
		case ExitCommand:
			return
		case EchoCommand:
			fmt.Println(cmd.ArgumentLine())
		default:
			fmt.Printf("%s: command not found\n", cmd.Name())
		}
	}
}
