package commands

import (
	"context"
	"fmt"
)

const EchoCommand = "echo"

func Echo(ctx context.Context, cmd CommandLine) error {
	fmt.Println(cmd.ArgumentLine())
	return nil
}
