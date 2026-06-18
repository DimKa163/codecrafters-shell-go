package main

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func main() {
	dispatcher := NewDispatcher()
	ctx := context.Background()
	for {
		if err := dispatcher.Execute(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println(err)
		}
	}
}
