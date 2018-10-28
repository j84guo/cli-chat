package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i:=0; i<2; i++ {

		// similar to the select() system call, select in Go blocks on multiple
		// channels and waits for the first one to be ready for the specified
		// operation (R or W)
		select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
		}
	}

	// timers can be implemented using select and channels
	c3 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
		case res := <-c3:
			fmt.Println(res)
		case <- time.After(1 * time.Second):
			fmt.Println("timeout 1")
	}
}
