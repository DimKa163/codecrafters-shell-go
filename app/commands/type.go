package commands

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

const (
	TypeCommand = "type"
)

func Type(ctx context.Context, cmd CommandLine) error {
	m, ok := ctx.Value(CommandStorageKey).(map[string]CommandHandler)
	if !ok {
		return fmt.Errorf("mismatch context")
	}
	_, ok = m[cmd.ArgumentLine()]
	if ok {
		fmt.Printf("%s is a shell builtin\n", cmd.ArgumentLine())
		return nil
	}
	var err error
	var path string
	if path, err = exec.LookPath(cmd.ArgumentLine()); err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%s: not found", cmd.ArgumentLine())
	}
	fmt.Printf("%s is %s\n", cmd.ArgumentLine(), path)
	return nil
}
