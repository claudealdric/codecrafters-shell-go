package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Shell struct {
	commandMap map[string]CommandHandler
}

func NewShell() *Shell {
	s := &Shell{
		make(map[string]CommandHandler),
	}
	return s
}

func (s *Shell) SetUpCommands() {
	s.commandMap["cd"] = Cd
	s.commandMap["exit"] = Exit
	s.commandMap["echo"] = Echo
	s.commandMap["pwd"] = Pwd
	s.commandMap["type"] = Type
}

func (s *Shell) Run() {
	for {
		printPrompt()
		input := readInput()
		command, args := parseInput(input)

		handleCommand, commandFound := s.commandMap[command]
		if commandFound {
			handleCommand(s, args)
			continue
		}

		executablePath, executable := GetExecutablePath(command)

		if executable {
			executeExternalCommand(executablePath, args)
			continue
		}

		fmt.Printf("%s: command not found\n", command)
	}
}

func executeExternalCommand(executablePath string, args []string) {
	command := exec.Command(executablePath, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		fmt.Fprintf(command.Stderr, "%s: %v\n", command, err)
	}
}

func parseInput(inputs string) (command string, args []string) {
	inputs = strings.Replace(inputs, "\n", "", 1)
	parts := strings.Split(inputs, " ")
	return parts[0], parts[1:]
}

func printPrompt() {
	fmt.Fprint(os.Stdout, "$ ")
}

func readInput() (input string) {
	input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}
