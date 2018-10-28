package main

import (
	"fmt"
)

type person struct {
	name string
	age int
}

// structs can also be nested to indicate an "is a" relationship
func main() {
	// there aren't constructors in Go, rather brace initialization is used
	fmt.Println(person{"Bob", 20})

	// fields can also be named
	fmt.Println(person{name: "Caila", age: 21})

	// omitted fields get the zero value for their type, all types in Go have a
	// zero value
	fmt.Println(person{})

	// prints the address of a newly created struct, in GC languages, stack vs
	// heap allocation is considered an implementation detail
	fmt.Println(&person{"Alice", 28})

	// dot notation can be used to access struct fields
	a := person{"Andrew", 90}
	fmt.Println(a.name, a.age)

	// and even pointers, which are automatically de-referenced
	var p *person = &a
	fmt.Println(p.name, p.age)

	// structs are mutable
	a.name = "Andrew Jr."
	fmt.Println(a)
}
