package main

import (
	"fmt"
)

const (
	ExitCommand = "exit"
)

func main() {
	var line string
	for {
		fmt.Print("$ ")
		fmt.Scanln(&line)
		switch line {
		case ExitCommand:
			return
		default:
			fmt.Printf("%s: command not found\n", line)
		}
	}
}
