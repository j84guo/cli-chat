package main

import (
	"io"
	"fmt"
	"net"
	"bufio"
)

func manageConnIn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Received:", scanner.Text())
	}
	CheckError(scanner.Err())
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
		CheckError(e)
	}
}

func manageStdin(connOutChan chan string) {
	for Stdin.Scan() {
		connOutChan <-Stdin.Text()
	}
	e := Stdin.Err()
	if e == io.EOF {
		fmt.Println("EOF detected")
	} else {
		CheckError(e)
	}
}

/* Todo: Use a done channel to signal goroutines to end?
   Note: net.Conn is safe to use by concurrent goroutines */
func main() {
	config, e := LoadOrPromptConfig()
	if e == io.EOF {
		return
	} else if e != nil {
		FatalError("LoadOrPromptConfig", e)
	}
	fmt.Println(config)

	connOutChan := make(chan string)
	conn, e := net.Dial("tcp", "127.0.0.1:8888")
	CheckError(e)
	defer conn.Close()

	go manageConnIn(conn)
	go manageConnOut(conn, connOutChan)
	manageStdin(connOutChan)
}
