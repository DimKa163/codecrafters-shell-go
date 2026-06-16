package main

import (
	"fmt"
)

func main() {
	var line string
	for {
		fmt.Print("$ ")
		fmt.Scanln(&line)
		switch line {
		default:
			fmt.Printf("%s: command not found\n", line)
		}
	}
}
