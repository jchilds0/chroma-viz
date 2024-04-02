package shows

import (
	"bufio"
	"chroma-viz/library/parser"
	"fmt"
	"strconv"
)

// T -> {'id': 123, 'num_geo': 123, 'layer': 123, 'geometry': [G]}
func parsePage(buf *bufio.Reader) (page *Page, err error) {
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

            page = NewPage(0, temp_id, layer, num_geo, data["name"])
            page.PropMap, err = parser.ParseProperty(buf, false)
            if err != nil {
                return
            }

			parser.MatchToken(']', buf)
		}

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
		}
	}

    if page == nil {
        err = fmt.Errorf("Page not created")
        return
    }

    return 
}

