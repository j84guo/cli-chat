package main

import (
	"fmt"
)

func main() {
	nums := []int{2, 3, 4}
	sum := 0
	for _, x := range nums {
		sum += x
	}
	fmt.Println("sum:", sum)

	for i, x := range nums {
		if x == 3 {
			fmt.Println("index:", i)
		}
	}

	m := map[string]string{"a":"apple", "b":"banana"}
	for k, v := range m {
		fmt.Printf("%s -> %s\n", k, v)
	}

	for k := range m {
		fmt.Println("key:", k)
	}

	// range on a string returns unicode code points, specifically the byte
	// index and the code point's value
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}
