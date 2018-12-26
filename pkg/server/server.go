package server

import (
	"io"
	"os"
	"fmt"
	"net"
	"cli-chat/pkg/utils"
	"cli-chat/pkg/chatmp"
)

type ClientMessage struct {
	From string
	Frame *chatmp.Frame
}

var (
	accepted = make(chan *chatmp.Conn)
	messages = make(chan *ClientMessage)
	terminated = make(chan string)
	active = make(map[string]*chatmp.Conn)
)

func selectForever() {
	for {
		select {
		case conn := <-accepted:
			go handleClient(conn)
		case msg := <-messages:
			relayMessage(msg)
		case username := <-terminated:
			cleanupClient(username)
		}
	}
}

func acceptForever(listener net.Listener) {
	for {
		transport, e := listener.Accept()
		utils.CheckError(e)
		accepted <- chatmp.NewConn(transport)
	}
}

func handleClient(conn *chatmp.Conn) {
	defer conn.Close()
	username, e := registerClient(conn)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Client error: %s\n", e.Error())
		return
	}

	for {
		f, e := conn.ReadFrame()
		if e != nil {
			if e == io.EOF {
				terminated <- username
			} else {
				fmt.Fprintf(os.Stderr, "Client error: %s\n", e.Error())
			}
			break
		}

		if f.Type != chatmp.CHAT_MSG {
			fmt.Fprintf(os.Stderr, "Unexpected frame: %s\n", f.Type)
		} else {
			messages <- &ClientMessage{username, f}
		}
	}
}

func relayMessage(msg *ClientMessage) {
	for username, conn := range active {
		if username == msg.From {
			continue
		}
		if e := conn.WriteFrame(msg.Frame); e != nil {
			fmt.Fprintf(os.Stderr, "Client error:", e.Error())
		}
	}
}

/** Todo: R/W lock */
func cleanupClient(username string) {
	delete(active, username)
}

func registerClient(conn *chatmp.Conn) (string, error) {
	username, e := chatmp.ParseClaimUsername(conn)
	if e != nil {
		return "", e
	}

	_, exists := active[username]
	if exists {
		if e := chatmp.YNResult(conn, false); e != nil {
			fmt.Fprintf(os.Stderr, "Client error: %s\n", e.Error())
		}
		return "", chatmp.ERR_NAMETAKEN
	}

	active[username] = conn
	if e := chatmp.YNResult(conn, true); e != nil {
		fmt.Fprintf(os.Stderr, "Client error: %s\n", e.Error())
	}
	return username, nil
}

/**
 * Goroutine for:
 * Accepting clients
 * Handling clients
 * Selecting on channels for terminated and accepted connections
 * Reading from new messages and relaying them
 */
func Run() {
	listener, e := net.Listen("tcp", "0.0.0.0:8000")
	if e != nil {
		utils.FatalError("net.Listen", e)
	}
	fmt.Println("Server started:", listener.Addr())

	go selectForever()
	acceptForever(listener)
}
