package terminal

import (
	"context"
)

type (
	CommandHandler func(context.Context, ...string) error
)
