package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command, _ := ParseInput(input)
		if command == "exit" {
			os.Exit(0)
		}

		fmt.Printf("%s: command not found\n", command)
	}
}

func ParseInput(input string) (command string, args []string) {
	input = strings.Replace(input, "\n", "", 1)
	parts := strings.Split(input, " ")
	return parts[0], parts[1:]
}
