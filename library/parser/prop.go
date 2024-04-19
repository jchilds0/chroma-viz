package parser

import (
	"bufio"
	"chroma-viz/library/props"
	"fmt"
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
func ParseProperty(buf *bufio.Reader, numGeo int) (propMap map[int]*props.Property, err error) {
	propMap = make(map[int]*props.Property, numGeo)
	data := make(map[string]string, 20)
	visible := make(map[string]bool, numGeo)

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

				propMap[prop_id] = props.NewProperty(prop_type, data["name"], false, visible)
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
