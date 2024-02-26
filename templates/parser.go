package templates 

import (
	"bufio"
	"chroma-viz/props"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// tokens
const (
    INT = iota + 256
    STRING
)

type Token struct {
    tok     int 
    value   string
    buf     []rune
}

var c_tok Token 

// S -> {'num_temp': 123, 'templates': [T]}
func (temp *Temps) ImportTemplates(conn net.Conn) error {
    // check conn is open
    if conn == nil {
        return fmt.Errorf("Graphics hub not connected")
    }

    buf := bufio.NewReader(conn)
    nextToken(buf)
    matchToken('{', buf)

    for c_tok.tok != '}' {
        // 'num_temp': ...
        if (c_tok.tok == STRING && c_tok.value == "num_temp") {
            matchToken(STRING, buf)
            matchToken(':', buf)
            n, err := strconv.Atoi(c_tok.value)
            if err != nil {
                log.Printf("Error reading %s (%s)", c_tok.value, err)
            }

            fmt.Printf("Number of templates %d\n", n)
            matchToken(INT, buf)
        }

        // 'templates': [...]
        if (c_tok.tok == STRING && c_tok.value == "templates") {
            matchToken(STRING, buf)
            matchToken(':', buf)
            matchToken('[', buf)

            temp.parseTemplate(buf)

            matchToken(']', buf)
        }

        if (c_tok.tok == ',') {
            matchToken(',', buf)
        }
    }

    return nil
}

// T -> {'id': 123, 'num_geo': 123, 'geometry': [G]} | T, T
func (temp *Temps) parseTemplate(buf *bufio.Reader) (err error) {
    data := make(map[string]string)
    matchToken('{', buf)

    for c_tok.tok == STRING {
        name := c_tok.value

        matchToken(STRING, buf)
        matchToken(':', buf)

        if name != "geometry" {
            data[name] = c_tok.value 

            nextToken(buf)
        } else {
            matchToken('[', buf)

            var num_geo, temp_id, layer int
            num_geo, err = strconv.Atoi(data["num_geo"])
            if err != nil {
                return
            }

            temp_id, err = strconv.Atoi(data["id"])
            if err != nil {
                return
            }

            layer, err = strconv.Atoi(data["layer"])

            if data["name"] == "" {
                data["name"] = "Template"
            }

            temp.SetTemplate(temp_id, layer, num_geo, data["name"])
            template := temp.Temps[temp_id]

            if err != nil {
                return
            }

            parseProperty(template, buf)

            matchToken(']', buf)
        }

        if c_tok.tok == ',' {
            matchToken(',', buf)
        }
    }

    matchToken('}', buf)

    if (c_tok.tok == ',') {
        matchToken(',', buf)
        temp.parseTemplate(buf)
    }

    return
}

// G -> {'id': 123, 'prop_type': 'abc', 'parent': 123, 'attr': [A]} | G, G
func parseProperty(temp *Template, buf *bufio.Reader) (err error) {
    data := make(map[string]string)
    matchToken('{', buf)

    for c_tok.tok == STRING {
        name := c_tok.value

        matchToken(STRING, buf)
        matchToken(':', buf)

        if name != "attr" {
            data[name] = c_tok.value 

            nextToken(buf)
        } else {
            matchToken('[', buf)

            var prop_id int
            prop_type := props.StringToProp[data["prop_type"]]
            prop_id, err = strconv.Atoi(data["id"])
            if err != nil {
                return
            }

            if data["name"] == "" {
                data["name"] = "Property"
            }

            temp.AddProp(data["name"], prop_id, prop_type)
            parseAttributes(temp, buf)

            matchToken(']', buf)
        }

        if c_tok.tok == ',' {
            matchToken(',', buf)
        }
    }

    matchToken('}', buf)

    if (c_tok.tok == ',') {
        matchToken(',', buf)
        parseProperty(temp, buf)
    }

    return nil
}

func parseAttributes(temp *Template, buf *bufio.Reader) (err error) {
    data := make(map[string]string)
    matchToken('{', buf)

    for c_tok.tok == STRING {
        name := c_tok.value

        matchToken(STRING, buf)
        matchToken(':', buf)

        data[name] = c_tok.value 
        nextToken(buf)

        if c_tok.tok == ',' {
            matchToken(',', buf)
        }
    }

    matchToken('}', buf)

    if (c_tok.tok == ',') {
        matchToken(',', buf)
        parseAttributes(temp, buf)
    }

    return nil
}

func matchToken(tok int, buf *bufio.Reader) {
    if tok != c_tok.tok {
        log.Printf("Incorrect token in graphics hub parsing (%d)", tok)
        return
    }

    nextToken(buf)
}

func nextToken(buf *bufio.Reader) {
    var err error
    c_tok, err = getToken(buf)
    if err != nil {
        log.Fatalf("Error getting next token (%s)", err)
    }
}

var peek = ' '

func getToken(buf *bufio.Reader) (tok Token, err error) {
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

    tok.buf = make([]rune, 256)

    switch peek {
    case '\'':
        var c rune 
        tok.tok = STRING
        peek = ' '

        for i := 0; ; i++ {
            c, _, err = buf.ReadRune()
            if err != nil {
                return 
            }

            if c == '\'' {
                break
            }

            tok.buf[i] = c
        }
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        tok.tok = INT
        for i := 0; '0' <= peek && peek <= '9'; i++ {
            tok.buf[i] = peek 

            peek, _, err = buf.ReadRune()
            if err != nil {
                return
            }
        }
    default:
        tok.tok = int(peek)
        tok.buf[0] = peek
        peek = ' '
    }

    tok.value = strings.TrimRight(string(tok.buf[:]), "\x00")
    return
}
