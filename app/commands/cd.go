package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

const CdCommand = "cd"

func Cd(ctx context.Context, cmd CommandLine) error {
	arg := cmd.Args()[0]
	if strings.HasPrefix(arg, "~") {
		hmDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		arg = strings.Replace(arg, "~", hmDir, 1)
	}
	if err := os.Chdir(cmd.Args()[0]); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Op == "chdir" {
			return fmt.Errorf("cd: %s: No such file or directory", cmd.Args()[0])
		}
		return err
	}
	return nil
}
