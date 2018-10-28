package main

import (
	"fmt"
)

type rect struct {
	width int
	height int
}

// methods have the syntax
// func <receiver> <name>() <return> {
//     <body>
// }
//
// the receiver may be passed by value or as a pointer, using a pointer permits
// the object to be modified and also saves copying, either way Go handles
// passing a value or pointer as the method receiver
func (r *rect) area() int {
	return r.width * r.height
}

func (r rect) perim() int {
	return 2 * r.width + 2 * r.height
}

func main() {
	var r rect = rect{width: 10, height: 20}
	fmt.Println("area:", r.area())
	fmt.Println("perim:", r.perim())

	var p *rect = &r
	fmt.Println("area:", p.area())
	fmt.Println("perim:", p.perim())
}
