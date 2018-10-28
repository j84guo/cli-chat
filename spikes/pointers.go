package main

import (
	"fmt"
)

func zeroVal(ival int) {
	ival = 0
}

func zeroPtr(ptr *int) {
	*ptr = 0
}

func main() {
	var i int = 1
	fmt.Println("initial:", i)

	zeroVal(i)
	fmt.Println("after value:", i)

	zeroPtr(&i)
	fmt.Println("after pointer:", i)
}
