package library

import (
	"chroma-viz/library/util"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Animator interface {
	Encode(b *strings.Builder)
}

const (
	BLANK = iota
	END_OF_CONN
	END_OF_MESSAGE
	ANIMATE_ON
	CONTINUE
	ANIMATE_OFF
	UPDATE
)

type Connection struct {
	Name      string
	addr      string
	port      int
	connected bool
	Conn      net.Conn
	SetPage   chan Animator
	SetAction chan int
}

func NewConnection(name, addr string, port int) *Connection {
	conn := &Connection{Name: name, addr: addr, port: port}
	conn.SetPage = make(chan Animator, 1)
	conn.SetAction = make(chan int, 1)

	go conn.SendPage()
	return conn
}

func (conn *Connection) Connect() (err error) {
	conn.Conn, err = net.Dial("tcp", conn.addr+":"+strconv.Itoa(conn.port))
	if err != nil {
		conn.connected = false
		return
	}

	conn.connected = true
	return
}

/*
Chroma message grammar (spaces are only for grammar readability)

S -> ver=%d#len=%d#action=%d#page=%d# G END_OF_MESSAGE
G -> geo_num=%d# P G
P -> attr=%s#val=%s# P
*/
func (conn *Connection) SendPage() {
	var page Animator

	for {
		action := <-conn.SetAction

		select {
		case page = <-conn.SetPage:
		default:
		}

		if page == nil {
			continue
		}

		if conn.IsConnected() == false {
			log.Printf("%s:%d is not connected", conn.addr, conn.port)
			continue
		}

		var b strings.Builder
		version := "1,4"

		util.EngineAddKeyValue(&b, "version", version)
		util.EngineAddKeyValue(&b, "action", action)

		page.Encode(&b)
		b.WriteByte(END_OF_MESSAGE)

		conn.Conn.Write([]byte(b.String()))
	}
}

func (conn *Connection) CloseConn() {
	if conn.Conn != nil {
		conn.Conn.Write([]byte(string(END_OF_CONN)))
		conn.Conn.Close()
	}
}

func (conn *Connection) IsConnected() bool {
	return conn.connected
}

func (conn *Connection) Read() (string, error) {
	buf := make([]byte, 100)
	conn.Conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err := conn.Conn.Read(buf)

	return string(buf), err
}

/*
- watch a connection for a close,
- call emit() on close,
- print any message recieved to stdout
*/
func (conn *Connection) Watcher(emit func()) {
	for {
		time.Sleep(500 * time.Millisecond)
		if conn.Conn == nil {
			emit()
			conn.connected = false
			return
		}

		string, err := conn.Read()
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			conn.connected = true
			continue
		}

		if err != nil {
			log.Print(err)
			conn.connected = false
			emit()
			return
		}

		log.Printf("(%s : %d): %s\n", conn.addr, conn.port, string)
		conn.connected = true
	}
}
