package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	shell := NewShell()
	shell.Run()
}

type CommandHandler func(shell *Shell, args []string)

type Shell struct {
	commandMap map[string]CommandHandler
}

func NewShell() *Shell {
	shell := &Shell{
		commandMap: make(map[string]CommandHandler),
	}
	shell.SetUpCommands()
	return shell
}

func (shell *Shell) SetUpCommands() {
	shell.commandMap["exit"] = Exit
	shell.commandMap["echo"] = Echo
	shell.commandMap["type"] = Type
}

func (shell *Shell) Run() {
	for {
		PrintPrompt()
		input := ReadInput()
		command, args := ParseInput(input)

		handleCommand, commandFound := shell.commandMap[command]
		if commandFound {
			handleCommand(shell, args)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func PrintPrompt() {
	fmt.Fprint(os.Stdout, "$ ")
}

func ReadInput() (input string) {
	input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}

func ParseInput(inputs string) (command string, args []string) {
	inputs = strings.Replace(inputs, "\n", "", 1)
	parts := strings.Split(inputs, " ")
	return parts[0], parts[1:]
}

func Echo(shell *Shell, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func Exit(shell *Shell, args []string) {
	os.Exit(0)
}

func Type(shell *Shell, args []string) {
	commandToCheck := args[0]

	_, commandFound := shell.commandMap[commandToCheck]
	if commandFound {
		fmt.Printf("%s is a shell builtin\n", commandToCheck)
		return
	}

	executablePath, isExecutable := GetExecutablePath(commandToCheck)
	if !isExecutable {
		fmt.Printf("%s: not found\n", commandToCheck)
		return
	}

	fmt.Printf("%s is %s\n", commandToCheck, executablePath)
}

func GetExecutablePath(command string) (executablePath string, isExecutable bool) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		executablePath := filepath.Join(dir, command)

		isExecutable, err := IsExecutable(executablePath)

		if err != nil {
			continue
		}

		if isExecutable {
			return executablePath, true
		}
	}

	return "", false
}

func IsExecutable(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.Mode()&0111 != 0, nil
}
