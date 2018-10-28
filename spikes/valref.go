package main

import "fmt"

func deleteKey(m map[string]int) {
	_, x := m["one"]

	if x {
		delete(m, "one")
	}
}

func main() {
	var m map[string]int = map[string]int {
		"one": 1,
		"two": 2,
	}

	deleteKey(m)
	fmt.Println(m)
}
