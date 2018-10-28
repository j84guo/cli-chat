package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

const (
	TYPE string = "tcp"
	ADDR string = "0.0.0.0:8888"
)

type ChatMsg struct {
	text string
	sender int
}

type ChatState struct {
	next int
	messages chan *ChatMsg
	accepted chan net.Conn
	terminated chan int
	clients map[int]net.Conn
}

func newChatState() *ChatState {
	var state ChatState

	state.next = 0
	state.messages = make(chan *ChatMsg)
	state.accepted = make(chan net.Conn)
	state.terminated = make(chan int)
	state.clients = make(map[int]net.Conn)

	return &state
}

func newChatMsg(text string, id int) *ChatMsg {
	var msg ChatMsg

	msg.text = text
	msg.sender = id

	return &msg
}

func acceptForever(server net.Listener, accepted *chan net.Conn) {
	for {
		con, e := server.Accept()
		if e != nil {
			fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			os.Exit(1)
		}

		*accepted <-con
	}
}

func loopForever(state *ChatState) {
	for {
		loopOne(state)
	}
}

func loopOne(state *ChatState) {
	select {
		case con := <-state.accepted:
			go handleClient(con, state.next, state)

		case msg := <-state.messages:
			handleMsg(msg, state)

		// todo: does this need to be synchronized
		case id := <-state.terminated:
			delete(state.clients, id)
	}
}

func handleClient(con net.Conn, id int, state *ChatState) {
	state.clients[id] = con
	state.next++
	reader := bufio.NewReader(con)

	for {
		data, e := reader.ReadString('\n')
		if e != nil {
			break
		}

		state.messages <-newChatMsg(data, id)
	}

	state.terminated <-id
}

func handleMsg(msg *ChatMsg, state *ChatState) {
	for id, con := range state.clients {
		if id == msg.sender {
			continue
		}

		go relayMsg(con, id, msg.text, state)
	}
}

func relayMsg(con net.Conn, id int, text string, state *ChatState) {
	_, e := con.Write([]byte(text))

	if e != nil {
		state.terminated <-id
	}
}

func main() {
	var state *ChatState = newChatState()

	server, e := net.Listen(TYPE, ADDR)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Error())
		os.Exit(1)
	}

	go acceptForever(server, &state.accepted)
	loopForever(state)
}
