package props

import (
	"bufio"
	"chroma-viz/library/parser"
	"fmt"
	"strconv"
)

// G -> {'id': 123, 'name': 'abc', 'prop_type': 'abc', 'geo_type': 'abc', 'visible': [...], 'attr': [A]} | G, G
func ParseProperty(buf *bufio.Reader, numGeo int) (propMap map[int]*Property, err error) {
	propMap = make(map[int]*Property, numGeo)
	data := make(map[string]string, 20)
	visible := make(map[string]bool, numGeo)

	for {
		parser.MatchToken('{', buf)

		for parser.C_tok.Tok == parser.STRING {
			name := parser.C_tok.Value

			parser.MatchToken(parser.STRING, buf)
			parser.MatchToken(':', buf)

			if name == "attr" {
				parser.MatchToken('[', buf)

				var prop_id int
				prop_type := StringToProp[data["prop_type"]]
				prop_id, err = strconv.Atoi(data["id"])
				if err != nil {
					err = fmt.Errorf("Error reading prop id from property (%s)", err)
					return
				}

				if data["name"] == "" {
					data["name"] = "Property"
				}

				propMap[prop_id] = NewProperty(prop_type, data["name"], false, visible)
				parseAttributes(propMap[prop_id], buf)

				parser.MatchToken(']', buf)
			} else if name == "visible" {
				parser.MatchToken('[', buf)

				for parser.C_tok.Tok == parser.STRING {
					attr := parser.C_tok.Value
					parser.MatchToken(parser.STRING, buf)
					parser.MatchToken(':', buf)

					visible[attr] = (parser.C_tok.Value == "true")
					parser.NextToken(buf)

					if parser.C_tok.Tok == ',' {
						parser.MatchToken(',', buf)
					}
				}

				parser.MatchToken(']', buf)
			} else {
				data[name] = parser.C_tok.Value

				parser.NextToken(buf)
			}

			if parser.C_tok.Tok == ',' {
				parser.MatchToken(',', buf)
			}
		}
		parser.MatchToken('}', buf)

		if parser.C_tok.Tok != ',' {
			break
		}

		parser.MatchToken(',', buf)
	}

	return
}

// A -> {'name': string, 'value': string} | A, A
func parseAttributes(prop *Property, buf *bufio.Reader) (err error) {
	data := make(map[string]string)
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value

		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		data[name] = parser.C_tok.Value
		parser.NextToken(buf)

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
		}
	}

	attr := prop.Attr[data["name"]]
	if attr != nil {
		attr.Decode(data["value"])
	}

	parser.MatchToken('}', buf)

	if parser.C_tok.Tok == ',' {
		parser.MatchToken(',', buf)
		parseAttributes(prop, buf)
	}

	return nil
}
