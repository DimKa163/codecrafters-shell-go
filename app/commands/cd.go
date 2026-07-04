package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const CdCommand = "cd"

type CD struct {
	Stdin  io.Reader
	Stdout io.Writer
	Args   []string
}

func NewCD(stdin io.Reader, stdout io.Writer, args ...string) InternalCommand {
	return &CD{
		Stdin:  stdin,
		Stdout: stdout,
		Args:   args,
	}
}

func (c *CD) Run() error {
	arg := c.Args[0]
	if strings.HasPrefix(arg, "~") {
		hmDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		arg = strings.Replace(arg, "~", hmDir, 1)
	}
	if err := os.Chdir(arg); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Op == "chdir" {
			return fmt.Errorf("cd: %s: No such file or directory", arg)
		}
		return err
	}
	return nil
}
