package commands

import (
	"fmt"
	"io"
	"os"
)

const PwdCommand = "pwd"

type pwdCmd struct {
	Stdin  io.Reader
	Stdout io.Writer
	Args   []string
}

func NewPwd(stdin io.Reader, stdout io.Writer, args ...string) InternalCommand {
	return &pwdCmd{
		Stdin:  stdin,
		Stdout: stdout,
		Args:   args,
	}
}

func (c *pwdCmd) Run() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}
