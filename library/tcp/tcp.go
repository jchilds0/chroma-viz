package tcp

import (
	"chroma-viz/library/props"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Animator interface {
	GetTemplateID() int
	GetLayer() int
	GetPropMap() map[int]*props.Property
}

const (
	END_OF_CONN = iota + 1
	END_OF_MESSAGE
	ANIMATE_ON
	CONTINUE
	ANIMATE_OFF
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

func (conn *Connection) Connect() {
	var err error
	conn.Conn, err = net.Dial("tcp", conn.addr+":"+strconv.Itoa(conn.port))

	if err != nil {
		log.Print(err)
		conn.connected = false
	}

	conn.connected = true
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

		version := [...]int{1, 4}

		header := fmt.Sprintf("ver=%d,%d#layer=%d#action=%d#temp=%d#",
			version[0], version[1], page.GetLayer(), action, page.GetTemplateID())

		geo := ""
		for i, prop := range page.GetPropMap() {
			if prop == nil {
				continue
			}

			geo = geo + fmt.Sprintf("geo_num=%d#%s", i, prop.String())
		}

		str := header + geo + string(END_OF_MESSAGE)
		conn.Conn.Write([]byte(str))
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
