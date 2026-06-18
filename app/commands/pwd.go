package commands

import (
	"context"
	"fmt"
	"os"
)

const PwdCommand = "pwd"

func Pwd(ctx context.Context, cmd CommandLine) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}
