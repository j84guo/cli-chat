package main

import (
	"fmt"
	"time"
)

// with channels synchronization can be done by message passing (we could
// achieve the same in C using semaphore or condition variable)
func worker(done chan bool) {
	fmt.Println("worker...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true
}

// when channels are function parameters, a direction may be specified
func ping(pings chan<-string, msg string) {
	pings <- msg
}

// chan <- type means send, <- chan type means receive
func pong(pings <-chan string, pongs chan<-string) {
	msg := <-pings
	pongs <- msg
}

func main() {
	// channels are unbuffered by default, a send may only occur when other
	// goroutines are running
	var mailbox chan string = make(chan string)
	fmt.Println(mailbox)

	go func() {
		mailbox<-"new mail"
	}()

	fmt.Println(<-mailbox)

	// we can also make a buffered channel
	mb := make(chan string, 2)
	mb <- "buffered"
	mb <- "channel"
	fmt.Println(<-mb)
	fmt.Println(<-mb)

	// wait until worker indicates it's done
	done := make(chan bool, 1)
	go worker(done)
	<-done
	fmt.Println("done waiting for worker")

	// bouncing message
	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	ping(pings, "ping message")
	pong(pings, pongs)
	fmt.Println(<-pongs)

}
