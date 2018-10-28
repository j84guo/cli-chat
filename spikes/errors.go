package main

import (
	"fmt"
	"errors"
)

func f1(arg int) (int, error) {
	if arg == 42 {
		return -1, errors.New("42 is invalid")
	}

	// nil is the zero value for pointers, interfaces, maps, slices
	// error is an interface which can be nil, custom errors are structs which
	// implement an interface, we should return custom structs as pointers
	// since error methods have pointer receiver
	return arg + 3, nil
}

type ArgError struct {
	arg int
	msg string
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("%d, %s", e.arg, e.msg)
}

func f2(arg int) (int, error) {
	if arg == 42 {
		return -1, &ArgError{arg, "invalid input"}
	}

	return arg + 5, nil
}

func main() {
	for _, i := range []int{7, 42} {
		if r, e := f1(i); e != nil {
			fmt.Println("f1 failed:", e)
		} else {
			fmt.Println("f1 worked:", r)
		}
	}

	for _, i := range []int{7, 42} {
		if r, e := f2(i); e != nil {
			fmt.Println("f2 failed:", e)
		} else {
			fmt.Println("f2 worked:", r)
		}
	}
}
