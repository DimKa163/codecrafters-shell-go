package commands

import (
	"context"
)

type (
	CommandHandler func(context.Context, ...string) error
)
