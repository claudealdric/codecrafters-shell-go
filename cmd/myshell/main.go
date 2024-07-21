package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	path := os.Getenv("PATH")
	fmt.Println(path)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command, args := ParseInput(input)
		commandMap := SetUpCommands()

		if handleCommand, commandFound := commandMap[command]; commandFound {
			handleCommand(command, args)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func ParseInput(input string) (command string, args []string) {
	input = strings.Replace(input, "\n", "", 1)
	parts := strings.Split(input, " ")
	return parts[0], parts[1:]
}

func SetUpCommands() map[string]CommandHandler {
	commandMap := map[string]CommandHandler{
		"exit": func(command string, args []string) {
			os.Exit(0)
		},
		"echo": func(command string, args []string) {
			returnValue := strings.Join(args, " ")
			fmt.Println(returnValue)
		},
	}
	commandMap["type"] = func(command string, args []string) {
		commandToCheck := args[0]
		if _, commandFound := commandMap[commandToCheck]; commandFound {
			fmt.Printf("%s is a shell builtin\n", commandToCheck)
		} else {
			fmt.Printf("%s: not found\n", commandToCheck)
		}

	}
	return commandMap
}

type CommandHandler func(command string, args []string)
