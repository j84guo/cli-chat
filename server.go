package main

import (
	"fmt"
	"net"
	"bufio"
)

const (
	CONN_TYPE string = "tcp"
	CONN_ADDR string = "0.0.0.0:8888"
)

/** sender ID and message */
type ChatMsg struct {
	text string
	sender int
}

/**
 * next - next ID
 * messages - channel for incoming messages
 * accepted - channel for accepted clients
 * terminated - channel for disconnected clients
 * clients - map of ID to net.Conn
 */
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
	return &ChatMsg{text: text, sender: id}
}

func acceptForever(server net.Listener, accepted *chan net.Conn) {
	for {
		con, e := server.Accept()
		if e != nil {
			FatalError(e.Error())
		}

		fmt.Println("Accepted:", con.RemoteAddr())
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

		/* TODO: does this need to be synchronized */
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
	state := newChatState()

	server, e := net.Listen(CONN_TYPE, CONN_ADDR)
	if e != nil {
		FatalError(e.Error())
	}

	go acceptForever(server, &state.accepted)
	loopForever(state)
}
