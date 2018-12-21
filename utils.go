package main

import (
	"os"
	"fmt"
	"strings"
)

func FatalError(msg string) {
	fmt.Fprintf(os.Stderr, msg + "\n")
	os.Exit(1)
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func PromptUsername() (string, error) {
	return PromptLine("new username:\n")
}

func PromptLine(msg string) (string, error) {
	var line string
	fmt.Print(msg)
	if _, e := fmt.Scanln(&line); e != nil {
		return "", e
	}
	return strings.TrimRight(line, "\r\n"), nil
}
