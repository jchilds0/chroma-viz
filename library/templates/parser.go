package templates

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
	tok   int
	value string
	buf   []rune
}

var c_tok Token

// T -> {'id': 123, 'num_geo': 123, 'layer': 123, 'geometry': [G]}
func parseTemplate(buf *bufio.Reader) (temp *Template, err error) {
	data := make(map[string]string)
    nextToken(buf)
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
			if err != nil {
				return
			}

			if data["name"] == "" {
				data["name"] = "Template"
			}

            temp = NewTemplate(data["name"], temp_id, layer, num_geo)
			parseProperty(temp, buf)
			matchToken(']', buf)
		}

		if c_tok.tok == ',' {
			matchToken(',', buf)
		}
	}

    if temp == nil {
        err = fmt.Errorf("Template not created")
        return
    }

    if c_tok.tok != '}' {
		return
    }
	return
}

// G -> {'id': 123, 'name': 'abc', 'prop_type': 'abc', 'geo_type': 'abc', 'visible': [...], 'attr': [A]} | G, G
func parseProperty(temp *Template, buf *bufio.Reader) (err error) {
	data := make(map[string]string)
	visible := make(map[string]bool)
	matchToken('{', buf)

	for c_tok.tok == STRING {
		name := c_tok.value

		matchToken(STRING, buf)
		matchToken(':', buf)

		if name == "attr" {
			matchToken('[', buf)

			var prop_id int
			prop_type := props.StringToProp[data["prop_type"]]
			prop_id, err = strconv.Atoi(data["id"])
			if err != nil {
				return fmt.Errorf("Error reading prop id from property (%s)", err)
			}

			if data["name"] == "" {
				data["name"] = "Property"
			}

			prop := temp.AddGeometry(data["name"], prop_id, prop_type, visible)
			parseAttributes(prop, buf)

			matchToken(']', buf)
		} else if name == "visible" {
			matchToken('[', buf)

			for c_tok.tok == STRING {
				attr := c_tok.value
				matchToken(STRING, buf)
				matchToken(':', buf)

				visible[attr] = (c_tok.value == "true")
				nextToken(buf)

				if c_tok.tok == ',' {
					matchToken(',', buf)
				}
			}

			matchToken(']', buf)
		} else {
			data[name] = c_tok.value

			nextToken(buf)
		}

		if c_tok.tok == ',' {
			matchToken(',', buf)
		}
	}

	matchToken('}', buf)

	if c_tok.tok == ',' {
		matchToken(',', buf)
		parseProperty(temp, buf)
	}

	return nil
}

// A -> {'name': string, 'value': string} | A, A
func parseAttributes(prop *props.Property, buf *bufio.Reader) (err error) {
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

	attr := prop.Attr[data["name"]]
	if attr != nil {
		attr.Decode(data["value"])
	}

	matchToken('}', buf)

	if c_tok.tok == ',' {
		matchToken(',', buf)
		parseAttributes(prop, buf)
	}

	return nil
}

func matchToken(tok int, buf *bufio.Reader) {
	if tok != c_tok.tok {
		log.Printf("Incorrect token %s, expected %c", c_tok.value, tok)
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

	tok.buf = make([]rune, 64)
    bufLength := 0

	switch peek {
	case '\'':
		var c rune
		tok.tok = STRING
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
		tok.tok = INT
		for '0' <= peek && peek <= '9' {
			tok.buf[bufLength] = peek
            bufLength++

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

    tok.value = string(tok.buf[:bufLength])
	return
}
