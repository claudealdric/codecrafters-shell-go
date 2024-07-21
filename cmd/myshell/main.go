package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
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
		path := os.Getenv("PATH")
		dirs := strings.Split(path, ":")
		commandToCheck := args[0]

		if _, commandFound := commandMap[commandToCheck]; commandFound {
			fmt.Printf("%s is a shell builtin\n", commandToCheck)
			return
		}

		executablePath, isExecutable := GetExecutablePath(commandToCheck, dirs)
		if !isExecutable {
			fmt.Printf("%s: not found\n", commandToCheck)
			return
		}

		fmt.Printf("%s is %s\n", commandToCheck, executablePath)
	}

	return commandMap
}

func GetExecutablePath(command string, dirs []string) (path string, isExecutable bool) {
	for _, dir := range dirs {
		executablePath := filepath.Join(dir, command)
		if _, err := os.Stat(executablePath); err == nil {
			if fileInfo, _ := os.Stat(executablePath); fileInfo.Mode()&0111 != 0 {
				return executablePath, true
			}
		}
	}

	return "", false
}

type CommandHandler func(command string, args []string)
