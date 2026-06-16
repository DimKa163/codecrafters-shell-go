package main

import (
	"fmt"
)

func main() {
	fmt.Print("$ ")
	var line string
	fmt.Scanln(&line)
	switch line {
	default:
		fmt.Printf("%s: command not found\n", line)
	}
}
