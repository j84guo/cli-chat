package chatmp

import (
	"net"
	"bufio"
	"errors"
	"strings"
	"strconv"
)

/**
 * Chat Message Protocol is a full-duplex application layer protocol. (Hello
 * WebSockets). Client and server communicate by sending frames, each
 * consisting of a type, headers and an optional payload.
 *
 * The frame type indicates the intent of the message, for example joining a
 * particular chat room, claiming a username or sending a message to a room.
 *
 * Delimiters between header key values is a colon ':' while between header
 * lines is a new line '\n'. The end of headers is indicated by another new
 * line '\n'.
 */
const (
	/** Frame Types */
	CHAT_MSG string = "CHAT_MSG"
	JOIN_ROOM string = "JOIN_ROOM"
	LEAVE_ROOM string = "LEAVE_ROOM"
	LIST_ROOMS string = "LIST_ROOMS"
	CREATE_ROOM string = "CREATE_ROOM"
	SERVER_INFO string = "SERVER_INFO"
	CLAIM_USERNAME string = "CLAIM_USERNAME"

	/** Headers */
	BODY_TYPE string = "Body-Type"
	BODY_LEN string = "Body-Length"
	YN_RESULT string = "Yes-No-Result"
	FROM_USERNAME string = "From-Username"

	/** Delimiters */
	KVDLM byte = ':'
	NEWLN byte = '\n'
)

var (
	/** Error codes */
	ERR_BADFRAME error = errors.New("Malformed frame")
	ERR_NAMETAKEN error = errors.New("Username taken")
)

/******************************************************************************
 * Frame                                                                      *
 *****************************************************************************/
type Frame struct {
	Type string
	Head map[string]string
	Body []byte
}

func NewFrame() *Frame {
	return &Frame{"", make(map[string]string), nil}
}

/******************************************************************************
 * Conn                                                                       *
 *****************************************************************************/
/** Todo: Conn re-connect, keep-alive (i.e. MQTT) */
type Conn struct {
	conn net.Conn
	in *bufio.Reader
	out *bufio.Writer
}

func NewConn(transport net.Conn) *Conn {
	if transport == nil {
		panic("NewConn: Expected non nil transport")
	}
	in := bufio.NewReader(transport)
	out := bufio.NewWriter(transport)
	return &Conn{transport, in, out}
}

func (conn *Conn) ReadFrame() (*Frame, error) {
	f := NewFrame()

	ftype, e := conn.in.ReadString(NEWLN)
	if e != nil {
		return nil, e
	}
	f.Type = ftype[:len(ftype) - 1]

	for {
		cont, e := conn.readHeader(f)
		if e != nil {
			return nil, e
		} else if !cont {
			break
		}
	}
	if _, ok := f.Head[BODY_LEN]; ok {
		if e := conn.readBody(f); e != nil {
			return nil, e
		}
	}

	return f, nil
}

func (conn *Conn) readHeader(f *Frame) (bool, error) {
	hstr, e := conn.in.ReadString(NEWLN)
	if e != nil {
		return false, e
	}

	hstr = hstr[:len(hstr) - 1]
	if hstr == "" {
		return false, nil
	}

	htokens := strings.Split(hstr, string(KVDLM))
	if len(htokens) != 2 {
		return false, ERR_BADFRAME
	}

	f.Head[htokens[0]] = htokens[1]
	return true, nil
}

func (conn *Conn) readBody(f *Frame) error {
	blenStr, ok := f.Head[BODY_LEN]
	if !ok {
		return ERR_BADFRAME
	}

	blen, e := strconv.Atoi(blenStr)
	if e != nil {
		return ERR_BADFRAME
	}

	f.Body = make([]byte, blen)
	_, e = conn.in.Read(f.Body)
	return e
}

func (conn *Conn) WriteFrame(f *Frame) error {
	if e := conn.writeType(f.Type); e != nil {
		return e
	}

	for hkey, hval := range f.Head {
		if e := conn.writeHeader(hkey, hval); e != nil {
			return nil
		}
	}

	if e := conn.out.WriteByte(NEWLN); e != nil {
		return e
	}
	if _, e := conn.out.Write(f.Body); e != nil {
		return e
	}

	e := conn.out.Flush()
	return e
}

func (conn *Conn) writeType(ftype string) error {
	buf := []byte(ftype)
	buf = append(buf, NEWLN)
	_, e := conn.out.Write(buf)
	return e
}

func (conn *Conn) writeHeader(hkey, hval string) error {
	buf := []byte(hkey)
	buf = append(buf, KVDLM)
	buf = append(buf, hval...)
	buf = append(buf, NEWLN)
	 _, e := conn.out.Write(buf)
	return e
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

/******************************************************************************
 * Client                                                                     *
 *****************************************************************************/
type Client struct {
	*Conn
	Username string
}

func NewClient(username, addrStr string) (*Client, error) {
	transport, e := net.Dial("tcp", addrStr)
	if e != nil {
		return nil, e
	}
	return &Client{NewConn(transport), username}, nil
}

func (client *Client) SendFrame(f *Frame) error {
	client.addDefaultHeaders(f)
	return client.WriteFrame(f)
}

/** Todo: Encoding, Different message types */
func (client *Client) SendText(text string) error {
	f := NewFrame()
	f.Type = CHAT_MSG
	f.Body = []byte(text)
	return client.SendFrame(f)
}

func (client *Client) RecvFrame() (*Frame, error) {
	return client.ReadFrame()
}

func (client *Client) addDefaultHeaders(f *Frame) {
	if _, ok := f.Head[BODY_LEN]; !ok && len(f.Body) > 0 {
		f.Head[BODY_LEN] = strconv.Itoa(len(f.Body))
	}
	if _, ok := f.Head[FROM_USERNAME]; !ok {
		f.Head[FROM_USERNAME] = client.Username
	}
}

func (client *Client) ClaimUsername() error {
	request := NewFrame()
	request.Type = CLAIM_USERNAME
	if e := client.SendFrame(request); e != nil {
		return e
	}

	response, e := client.RecvFrame()
	if e != nil {
		return e
	}

	if response.Type != SERVER_INFO {
		return ERR_BADFRAME
	}

	res, ok := response.Head[YN_RESULT]
	if !ok {
		return ERR_BADFRAME
	}

	if res != "Y" {
		return ERR_NAMETAKEN
	}
	return nil
}

/******************************************************************************
 * Utils                                                                      *
 *****************************************************************************/
func ParseClaimUsername(conn *Conn) (string, error) {
	f, e := conn.ReadFrame()
	if e != nil || f.Type != CLAIM_USERNAME {
		return "", e
	}
	username := f.Head[FROM_USERNAME]
	return username, nil
}

func YNResult(conn *Conn, yes bool) error {
	f := NewFrame()
	f.Type = SERVER_INFO

	if yes {
		f.Head[YN_RESULT] = "Y"
	} else {
		f.Head[YN_RESULT] = "N"
	}

	return conn.WriteFrame(f)
}
