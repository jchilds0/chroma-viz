package templates

import (
	"bufio"
	"chroma-viz/library/props"
	"fmt"
	"log"
	"strconv"
	"strings"
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
func (temp *Template) parseTemplate(buf *bufio.Reader) (err error) {
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
				return fmt.Errorf("Error reading num geo from template (%s)", err)
			}

			temp_id, err = strconv.Atoi(data["id"])
			if err != nil {
				return fmt.Errorf("Error reading temp id from template (%s)", err)
			}

			layer, err = strconv.Atoi(data["layer"])
			if err != nil {
				return fmt.Errorf("Error reading layer from template (%s)", err)
			}

			if data["name"] == "" {
				data["name"] = "Template"
			}

            temp.TempID = temp_id
            temp.Layer = layer
            temp.NumGeo = num_geo
            temp.Title = data["name"]

			parseProperty(temp, buf)

			matchToken(']', buf)
		}

		if c_tok.tok == ',' {
			matchToken(',', buf)
		}
	}

	matchToken('}', buf)
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
