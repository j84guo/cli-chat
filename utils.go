package main

import (
	"os"
	"fmt"
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
