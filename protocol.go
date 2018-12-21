package main

import (
	"net"
	"fmt"
)

/**
 * Chat Transfer Protocol
 *
 * Full-duplex application layer protocol, consisting of lightweight frames.
 * (Hello WebSockets)
 */
type CTPFrame struct {
	Type string
	Head map[string]string
	Body []byte
}

func recvCTPFrame(conn net.Conn) (*CTPFrame, error) {
	return nil, nil
}

func sendCTPFrame(conn net.Conn, ftype string, fhead map[string]string,
		fbody []byte) (error) {
	return nil
}

func main() {
	var frame CTPFrame
	fmt.Println(frame)
}
