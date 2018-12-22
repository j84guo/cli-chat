package main

import (
	"net"
	"fmt"
	"bufio"
	"errors"
	"strings"
	"strconv"
)

const (
	/**
	 * Frame Types
	 * 
	 * By default the client will first ask the server to LIST available rooms
	 * on startup (assume public for now). The user may JOIN an existing room
	 * or CREATE one. After the user is done, he/she may LEAVE.
	 *
	 * At any point in time, the user will be checked into a room. This
	 * determines the contents of the UI and is the last joined room unless the
	 * user explicitly goes to a different one.
	 *
	 * During a session, messages from joined rooms are relayed to the client.
	 * It is up to the client to determine how to display the messages, i.e.
	 * based on which room is checked out. Messages for background rooms should
	 * be saved locally until later checkout.
	 */
	JOIN_ROOM string = "JOIN_ROOM"
	LEAVE_ROOM string = "LEAVE_ROOM"
	LIST_ROOMS string = "LIST_ROOMS"
	CREATE_ROOM string = "CREATE_ROOM"

	CHAT_MSG string = "CHAT_MSG"

	/** Headers */
	BODY_TYPE string = "Body-Type"
	BODY_LEN string = "Body-Length"

	/** Delimiters */
	NEWLN byte = '\n'
	KVDLM byte = ':'
)

/**
 * Chat Message Protocol
 *
 * Full-duplex application layer protocol, consisting of lightweight frames
 * (Hello WebSockets)
 */
type Frame struct {
	Type string
	Head map[string]string
	Body []byte
}

/**
 * Todo:
 * Connection re-connect, keep-alive (i.e. MQTT)
 */
type Conn struct {
	conn net.Conn
	in *bufio.Reader
	out *bufio.Writer
}

func NewFrame() *Frame {
	return &Frame{"", make(map[string]string), nil}
}

func Dial(addrStr string) (*Conn, error) {
	tcpConn, e := net.Dial("tcp", addrStr)
	if e != nil {
		return nil, e
	}
	return wrap(tcpConn), nil
}

func wrap(conn net.Conn) (*Conn) {
	if conn == nil {
		panic("cmp.wrap: expected non-nil")
	}
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)
	return &Conn{conn, in, out}
}

func (chat *Conn) ReadFrame() (*Frame, error) {
	f := NewFrame()

	ftype, e := chat.in.ReadString(NEWLN)
	if e != nil {
		return nil, e
	}
	f.Type = ftype[:len(ftype) - 1]

	for {
		cont, e := chat.ReadHeader(f)
		if e != nil {
			return nil, e
		} else if !cont {
			break
		}
	}

	if _, ok := f.Head[BODY_LEN]; ok {
		if e := chat.ReadBody(f); e != nil {
			return nil, e
		}
	}

	return f, nil
}

func (chat *Conn) ReadHeader(f *Frame) (bool, error) {
	hstr, e := chat.in.ReadString(NEWLN)
	if e != nil {
		return false, e
	}

	hstr = hstr[:len(hstr) - 1]
	if hstr == "" {
		return false, nil
	}

	htokens := strings.Split(hstr, string(KVDLM))
	if len(htokens) != 2 {
		return false, errors.New("Bad frame: header missing delim")
	}

	f.Head[htokens[0]] = htokens[1]
	return true, nil
}

func (chat *Conn) ReadBody(f *Frame) error {
	blenStr, ok := f.Head[BODY_LEN]
	if !ok {
		return errors.New("Bad frame: Body-Len")
	}

	blen, e := strconv.Atoi(blenStr)
	if e != nil {
		errors.New("Bad frame: Body-Len")
	}

	f.Body = make([]byte, blen)
	_, e = chat.in.Read(f.Body)
	return e
}

/**
 * Simple blocking send
 *
 * Todo:
 * Encoding
 * Message type (text vs binary)
 * Cleanup delimiters
 */
func (chat *Conn) WriteFrame(f *Frame) (error) {
	if e := chat.WriteType(f.Type); e != nil {
		return e
	}

	for hkey, hval := range f.Head {
		if e := chat.WriteHeader(hkey, hval); e != nil {
			return nil
		}
	}
	if _, ok := f.Head[BODY_LEN]; !ok && len(f.Body) > 0 {
		if e := chat.WriteHeader(BODY_LEN, strconv.Itoa(len(f.Body))); e != nil {
			return nil
		}
	}

	if e := chat.out.WriteByte(NEWLN); e != nil {
		return e
	}
	if _, e := chat.out.Write(f.Body); e != nil {
		return e
	}

	e := chat.out.Flush()
	return e
}

func (chat *Conn) WriteType(ftype string) error {
	buf := []byte(ftype)
	buf = append(buf, NEWLN)
	_, e := chat.out.Write(buf)
	return e
}

func (chat *Conn) WriteHeader(hkey, hval string) error {
	buf := []byte(hkey)
	buf = append(buf, KVDLM)
	buf = append(buf, hval...)
	buf = append(buf, NEWLN)
	 _, e := chat.out.Write(buf)
	return e
}

func main() {
	conn, e := Dial("localhost:8000")
	if e != nil {
		FatalError("Dial", e)
	}

	frame := &Frame{
		"SIGNUP",
		map[string]string{
			"SomeKey": "SomeValue",
		},
		[]byte("Hello, world!"),
	}

	e = conn.WriteFrame(frame)
	if e != nil {
		FatalError("Conn.WriteFrame", e)
	}

	resp, e := conn.ReadFrame()
	if e != nil {
		FatalError("Conn.ReadFrame", e)
	}
	fmt.Println(resp)
}