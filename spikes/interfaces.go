package main

import (
	"fmt"
	"math"
)

// interfaces are named collections of methods and fields
type geometry interface {
	area() float64
	perim() float64
}

// structs don't need any implements keyword, they conform to an interface by
// virtue of having the right methods
type rect struct {
	width float64
	height float64
}

func (r rect) area() float64 {
	return r.width * r.height
}

func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

// geometry implementatino for circle
type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

// polymorphism can be achieved using interface types
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	var c circle = circle{10}
	measure(c)

	var r rect = rect{10, 20}
	measure(r)
}
