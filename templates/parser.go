package templates

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// tokens
const (
    INT = iota + 256
    STRING
)

type Token struct {
    tok     int 
    value   []rune
}

var c_tok Token 

func ImportTemplates(conn net.Conn, temps map[int]*Template) error {
    // check conn is open

    buf := bufio.NewReader(conn)
    err := conn.SetReadDeadline(time.Now().Add(time.Second))
    if err != nil {
        log.Printf("Error setting read deadline (%s)", err)
    }

    for {
        char, err := buf.ReadByte()
        if err != nil {
            log.Printf("Done reading from tcp (%s)", err)
            break
        }

        fmt.Printf("%c", char)
    }
    
    return nil
}

func matchToken(tok int, buf bufio.Reader) {
    if tok != c_tok.tok {
        log.Printf("Incorrect token in graphics hub parsing (%d)", tok)
        return
    }
    
    var err error
    c_tok, err = getToken(buf)
    if err != nil {
        log.Printf("Error getting next token (%s)", err)
    }
}

var peek = ' '

func getToken(buf bufio.Reader) (tok Token, err error) {
WS: 
    for {
        switch peek {
        case ' ', '\t', '\r', '\n':
        default:
            break WS
        }

        peek, _, err = buf.ReadRune()
        if err != nil {
            return
        }
    }

    tok.value = make([]rune, 256)

    switch peek {
    case '\'':
        var c rune 
        i := 0
        tok.tok = STRING
        peek = ' '

        for {
            c, _, err = buf.ReadRune()
            if err != nil || c == '\'' {
                return 
            }

            tok.value[i] = c
        }
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        tok.tok = INT
        for i := 0; '0' <= peek && peek <= '9'; i++ {
            tok.value[i] = peek 

            peek, _, err = buf.ReadRune()
            if err != nil {
                return 
            }
        }
    default:
        tok.tok = int(peek)
        tok.value[0] = peek
        peek = ' '
        return
    }

    return
}
