package main

import (
	"fmt"
)

func main() {
	m := make(map[string]int)
	m["one"] = 1
	m["two"] = 2
	fmt.Println("map:", m)

	// the value type's zero value is returned on invalid key
	v1 := m["one"]
	fmt.Println("v1:", v1)
	fmt.Println("len:", len(m))

	// deleting an invalid key seems to raise no error
	delete(m, "one")
	fmt.Println("map:", m)
	fmt.Println("len:", len(m))

	// this form of access returns a boolean indicating key existence
	a, b := m["twoo"]
	fmt.Println(a, b)

	// brace initialization
	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
