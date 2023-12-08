package gui

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
    END_OF_CONN = iota + 1
    END_OF_MESSAGE
    ANIMATE_ON
    CONTINUE
    ANIMATE_OFF
)

type Connection struct {
    addr    string
    port    int
    conn    net.Conn
}

func NewConnection(addr string, port int) *Connection {
    return &Connection{addr: addr, port: port}
}

func (conn *Connection) Connect() bool {
    var err error
    conn.conn, err = net.Dial("tcp", conn.addr + ":" + strconv.Itoa(conn.port))

    if err != nil {
        log.Print(err)
        return false
    }

    return true
}

// TCP Format: ver%d#len%d#action%d#page%d#attr%s#val%d ... END_OF_MESSAGE
func (conn *Connection) SendPage(page *Page, action int) {
    if conn.IsConnected() == false {
        //log.Printf("%s:%d is not connected", conn.addr, conn.port)
        return
    }

    version := [...]int{1, 0}
    length := 2

    header := fmt.Sprintf("ver%d,%d#len%d#action%d#temp%d#", 
        version[0], version[1], length, action, page.templateID)

    str := header

    for _, prop := range page.props {
        str = str + prop.String()
    }

    str = str + string(END_OF_MESSAGE)
    conn.conn.Write([]byte(str))
}

func (conn *Connection) CloseConn() {
    if conn.conn != nil {
        conn.conn.Write([]byte(string(END_OF_CONN)))
        conn.conn.Close()
    }
}

func (conn *Connection) IsConnected() bool {
    if conn.conn == nil {
        return false
    }

    buf := make([]byte, 100)
    conn.conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
    _, err := conn.conn.Read(buf)

    if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
        return true
    }

    if err != nil {
        log.Print(err)
        return false
    }

    //fmt.Printf("Server Message: %s\n", buf);
    return true
}

func (conn *Connection) Read() {
    if conn.IsConnected() == false {
        return
    }

    buf := make([]byte, 100)
    conn.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
    _, err := conn.conn.Read(buf)

    if err != nil {
        log.Print(err)
        return 
    }

    fmt.Printf("Server Message: %s\n", buf);
}
