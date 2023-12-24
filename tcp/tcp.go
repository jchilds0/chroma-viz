package tcp

import (
	"chroma-viz/props"
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
    addr        string
    port        int
    connected   bool
    conn        net.Conn
    SetPage     chan *props.Page
    SetAction   chan int
}

func NewConnection(addr string, port int) *Connection {
    conn := &Connection{addr: addr, port: port}
    conn.SetPage = make(chan *props.Page, 1)
    conn.SetAction = make(chan int, 1)

    go conn.SendPage()
    return conn
}

func (conn *Connection) Connect() {
    var err error
    conn.conn, err = net.Dial("tcp", conn.addr + ":" + strconv.Itoa(conn.port))

    if err != nil {
        log.Print(err)
        conn.connected = false
    }

    conn.connected = true 
}

// TCP Format: ver%d#len%d#action%d#page%d#attr%s#val%d ... END_OF_MESSAGE
func (conn *Connection) SendPage() {
    var page *props.Page

    for {
        action := <-conn.SetAction

        select {
        case page = <-conn.SetPage:
        default:
        }

        if page == nil {
            //log.Println("No page selected")
            continue
        }

        // switch (action) {
        // case ANIMATE_ON, ANIMATE_OFF:
        //     currentPage = page
        // case CONTINUE:
        //     if currentPage != nil && page.pageNum == currentPage.pageNum {
        //         action = CONTINUE 
        //     } else {
        //         currentPage = page
        //         action = ANIMATE_ON
        //     }
        // }

        if conn.IsConnected() == false {
            //log.Printf("%s:%d is not connected", conn.addr, conn.port)
            continue
        }

        version := [...]int{1, 4}

        header := fmt.Sprintf("ver=%d,%d#layer=%d#action=%d#temp=%d#", 
            version[0], version[1], page.Layer, action, page.TemplateID)

        str := header

        for i, prop := range page.PropMap {
            str = fmt.Sprintf("%sgeo_num=%d#%s", str, i + 1, prop.String())
        }

        str = str + string(END_OF_MESSAGE)
        conn.conn.Write([]byte(str))
    }
}

func (conn *Connection) CloseConn() {
    if conn.conn != nil {
        conn.conn.Write([]byte(string(END_OF_CONN)))
        conn.conn.Close()
    }
}

func (conn *Connection) IsConnected() bool {
    return conn.connected
}

func (conn *Connection) Read() (string, error) {
    buf := make([]byte, 100)
    conn.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
    _, err := conn.conn.Read(buf)

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
        if conn.conn == nil {
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

        log.Printf("(%s : %d): %s\n", conn.addr, conn.port, string);
        conn.connected = true
    }
}
