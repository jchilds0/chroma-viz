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
    setPage chan *Page
    sendPage chan int
}

func NewConnection(addr string, port int) *Connection {
    conn := &Connection{addr: addr, port: port}
    conn.setPage = make(chan *Page, 1)
    conn.sendPage = make(chan int, 1)

    go conn.SendPage()
    return conn
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
func (conn *Connection) SendPage() {
    var page *Page

    for {
        action := <-conn.sendPage

        select {
        case page = <-conn.setPage:
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

        version := [...]int{1, 2}
        length := 2

        header := fmt.Sprintf("ver=%d,%d#len=%d#action=%d#temp=%d#", 
            version[0], version[1], length, action, page.templateID)

        str := header

        for _, prop := range page.propMap {
            str = str + prop.String()
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
