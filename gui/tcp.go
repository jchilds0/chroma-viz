package gui

import (
	"log"
	"net"
	"strconv"
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

func NewConnection(addr string) *Connection {
    return &Connection{addr: addr, port: 6100}
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

func (conn *Connection) SendPage(pageNum int, action int) {
    conn.conn.Write([]byte(strconv.Itoa(pageNum) + string(action) + string(END_OF_MESSAGE)))
}

func (conn *Connection) CloseConn() {
    conn.conn.Write([]byte(string(END_OF_CONN)))
    conn.conn.Close()
}

