package main

import (
	"os"
	"path/filepath"
	"strings"
)

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
