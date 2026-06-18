package commands

import (
	"context"
	"io"
)

const (
	ExitCommand = "exit"
)

func Exit(ctx context.Context, cmd CommandLine) error {
	return io.EOF
}
