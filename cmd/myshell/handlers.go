package main

import (
	"fmt"
	"os"
	"strings"
)

type CommandHandler func(s *Shell, args []string)

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
