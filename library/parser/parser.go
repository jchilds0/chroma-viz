package parser

import (
	"bufio"
	"chroma-viz/library/props"
	"fmt"
	"log"
	"strconv"
)

// tokens
const (
	INT = iota + 256
	STRING
)

type Token struct {
	Tok   int
	Value string
	buf   []rune
}

var C_tok Token

// G -> {'id': 123, 'name': 'abc', 'prop_type': 'abc', 'geo_type': 'abc', 'visible': [...], 'attr': [A]} | G, G
func ParseProperty(buf *bufio.Reader, isTemp bool) (propMap map[int]*props.Property, err error) {
    propMap = make(map[int]*props.Property, 10)
	data := make(map[string]string)
	visible := make(map[string]bool)

    for { 
        MatchToken('{', buf)

        for C_tok.Tok == STRING {
            name := C_tok.Value

            MatchToken(STRING, buf)
            MatchToken(':', buf)

            if name == "attr" {
                MatchToken('[', buf)

                var prop_id int
                prop_type := props.StringToProp[data["prop_type"]]
                prop_id, err = strconv.Atoi(data["id"])
                if err != nil {
                    err = fmt.Errorf("Error reading prop id from property (%s)", err)
                    return
                }

                if data["name"] == "" {
                    data["name"] = "Property"
                }

                propMap[prop_id] = props.NewProperty(prop_type, data["name"], isTemp, visible)
                parseAttributes(propMap[prop_id], buf)

                MatchToken(']', buf)
            } else if name == "visible" {
                MatchToken('[', buf)

                for C_tok.Tok == STRING {
                    attr := C_tok.Value
                    MatchToken(STRING, buf)
                    MatchToken(':', buf)

                    visible[attr] = (C_tok.Value == "true")
                    NextToken(buf)

                    if C_tok.Tok == ',' {
                        MatchToken(',', buf)
                    }
                }

                MatchToken(']', buf)
            } else {
                data[name] = C_tok.Value

                NextToken(buf)
            }

            if C_tok.Tok == ',' {
                MatchToken(',', buf)
            }
        }
        MatchToken('}', buf)

        if C_tok.Tok != ',' {
            break
        }

        MatchToken(',', buf)
    }

	return 
}

// A -> {'name': string, 'value': string} | A, A
func parseAttributes(prop *props.Property, buf *bufio.Reader) (err error) {
	data := make(map[string]string)
	MatchToken('{', buf)

	for C_tok.Tok == STRING {
		name := C_tok.Value

		MatchToken(STRING, buf)
		MatchToken(':', buf)

		data[name] = C_tok.Value
		NextToken(buf)

		if C_tok.Tok == ',' {
			MatchToken(',', buf)
		}
	}

	attr := prop.Attr[data["name"]]
	if attr != nil {
		attr.Decode(data["value"])
	}

	MatchToken('}', buf)

	if C_tok.Tok == ',' {
		MatchToken(',', buf)
		parseAttributes(prop, buf)
	}

	return nil
}

func MatchToken(tok int, buf *bufio.Reader) {
	if tok != C_tok.Tok {
		log.Printf("Incorrect token %s, expected %c", C_tok.Value, tok)
		return
	}

	NextToken(buf)
}

func NextToken(buf *bufio.Reader) {
	var err error
	C_tok, err = getToken(buf)
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

	tok.buf = make([]rune, 64)
    bufLength := 0

	switch peek {
	case '\'':
		var c rune
		tok.Tok = STRING
		peek = ' '

		for {
			c, _, err = buf.ReadRune()
			if err != nil {
				return
			}

			if c == '\'' {
				break
			}

			tok.buf[bufLength] = c
            bufLength++
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Tok = INT
		for '0' <= peek && peek <= '9' {
			tok.buf[bufLength] = peek
            bufLength++

			peek, _, err = buf.ReadRune()
			if err != nil {
				return
			}
		}
	default:
		tok.Tok = int(peek)
		tok.buf[0] = peek
		peek = ' '
	}

    tok.Value = string(tok.buf[:bufLength])
	return
}
