package main

import (
	"io"
	"os"
	"fmt"
	"bufio"
)

var Stdin *bufio.Scanner = bufio.NewScanner(os.Stdin)

func FatalError(prefix string, e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", prefix, e.Error())
	os.Exit(1)
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func PromptUsername() (string, error) {
	return PromptLine("New username:\n")
}

func PromptLine(msg string) (string, error) {
	fmt.Print(msg)
	if !Stdin.Scan() {
		return "", io.EOF
	}
	if e := Stdin.Err(); e != nil {
		return "", e
	}
	return Stdin.Text(), nil
}
