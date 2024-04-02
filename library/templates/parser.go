package templates

import (
	"bufio"
	"chroma-viz/library/parser"
	"fmt"
	"strconv"
)

// T -> {'id': 123, 'num_geo': 123, 'layer': 123, 'geometry': [G]}
func parseTemplate(buf *bufio.Reader) (temp *Template, err error) {
	data := make(map[string]string)
	parser.NextToken(buf)
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value

		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		if name != "geometry" {
			data[name] = parser.C_tok.Value

			parser.NextToken(buf)
		} else {
			parser.MatchToken('[', buf)

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
			temp.Geometry, err = parser.ParseProperty(buf, false)
			if err != nil {
				return
			}

			parser.MatchToken(']', buf)
		}

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
		}
	}

	if temp == nil {
		err = fmt.Errorf("Template not created")
		return
	}

	return
}
