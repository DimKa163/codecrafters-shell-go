package terminal

import "io"

type Cfg struct {
	Prompt    []rune
	Completer Completer
	Stdin     io.ReadCloser
	Stdout    io.Writer
	Stderr    io.Writer
}
