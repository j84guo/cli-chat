package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
)

func manageConnIn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Received:", scanner.Text())
	}
	if e := scanner.Err(); e != nil {
		panic(e)
	}
}

func manageConnOut(conn net.Conn, connOutChan chan string) {
	for {
		userMsg := <-connOutChan
		fmt.Println("Sending:", userMsg)

		/* Here we've ignored partial writes, as error will be reported in that
           case anyways. This is in contrast with the Unix standard and socket
           API, where write() and send() returning a partial write is not
           considered an error. */
		_, e := conn.Write([]byte(userMsg + "\n"))
		if e != nil {
			panic(e)
		}
	}
}

func manageStdin(connOutChan chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		connOutChan <-scanner.Text()
	}
	if e := scanner.Err(); e != nil {
		panic(e)
	}
}

/* Todo: Use a done channel to signal goroutines to end?
   Note: net.Conn is safe to use by concurrent goroutines */
func main() {
	connOutChan := make(chan string)

	conn, e := net.Dial("tcp", "127.0.0.1:8888")
	if e != nil {
		panic(e)
	}
	defer conn.Close()

	go manageConnIn(conn)
	go manageConnOut(conn, connOutChan)
	manageStdin(connOutChan)
}
