package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	s := Shell{
		make(map[string]CommandHandler),
	}
	s.SetUpCommands()
	s.Run()
}

type CommandHandler func(s *Shell, args []string)

type Shell struct {
	commandMap map[string]CommandHandler
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
		PrintPrompt()
		input := ReadInput()
		command, args := ParseInput(input)

		handleCommand, commandFound := s.commandMap[command]
		if commandFound {
			handleCommand(s, args)
			continue
		}

		executablePath, executable := GetExecutablePath(command)

		if executable {
			ExecuteExternalCommand(executablePath, args)
			continue
		}

		fmt.Printf("%s: command not found\n", command)
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

func Cd(s *Shell, args []string) {
	dir := args[0]

	if dir == "~" {
		dir = os.Getenv("HOME")
	}

	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
	}
}

func Echo(s *Shell, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func Exit(s *Shell, args []string) {
	os.Exit(0)
}

func Pwd(s *Shell, args []string) {
	dir, _ := os.Getwd()
	fmt.Println(dir)
}

func Type(s *Shell, args []string) {
	commandToCheck := args[0]

	_, commandFound := s.commandMap[commandToCheck]
	if commandFound {
		fmt.Printf("%s is a shell builtin\n", commandToCheck)
		return
	}

	executablePath, executable := GetExecutablePath(commandToCheck)
	if !executable {
		fmt.Printf("%s: not found\n", commandToCheck)
		return
	}

	fmt.Printf("%s is %s\n", commandToCheck, executablePath)
}

func GetExecutablePath(command string) (executablePath string, executable bool) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		executablePath := filepath.Join(dir, command)

		executable, err := IsExecutable(executablePath)

		if err != nil {
			continue
		}

		if executable {
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

func ExecuteExternalCommand(executablePath string, args []string) {
	command := exec.Command(executablePath, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		fmt.Fprintf(command.Stderr, "%s: %v\n", command, err)
	}
}
