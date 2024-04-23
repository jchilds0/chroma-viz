package templates

import (
	"bufio"
	"chroma-viz/library/parser"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// T -> {'id': 123, 'num_geo': 123, 'layer': 123, 'geometry': [G]}
func parseTemplate(buf *bufio.Reader) (temp Template, err error) {
	parser.NextToken(buf)
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value

		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		switch name {
		case "id":
			temp.TempID, err = strconv.ParseInt(parser.C_tok.Value, 10, 64)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)
		case "num_keyframe":
			var numKey int
			numKey, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)
			temp.Keyframe = make([]Keyframe, 0, numKey)
		case "num_geo":
			var numGeo int
			numGeo, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)
			temp.Geometry = make([]IGeometry, 0, numGeo)
		case "keyframe":
			parser.MatchToken('[', buf)

			var frame Keyframe
			for parser.C_tok.Tok == '{' {
				frame, err = parseKeyframe(buf)
				if err != nil {
					return
				}

				temp.Keyframe = append(temp.Keyframe, frame)

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

			parser.MatchToken(']', buf)
		case "geometry":
			parser.MatchToken('[', buf)

			var geo IGeometry
			for parser.C_tok.Tok == '{' {
				geo, err = parseGeometry(buf)
				if err != nil {
					return
				}

				temp.Geometry = append(temp.Geometry, geo)

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

			parser.MatchToken(']', buf)
		default:
			log.Printf("Unknown template attribute %s", name)
		}
	}

	return
}

func parseKeyframe(buf *bufio.Reader) (frame Keyframe, err error) {
	var num int
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value
		err = parser.MatchToken(parser.STRING, buf)
		if err != nil {
			return
		}

		err = parser.MatchToken(':', buf)
		if err != nil {
			return
		}

		switch name {
		case "frame_num":
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameNum = num
		case "frame_geo":
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameGeo = num

		case "frame_attr":
			frame.FrameAttr = parser.C_tok.Value

		case "mask":
			frame.Mask = (parser.C_tok.Value == "true")

		case "expand":
			frame.Expand = (parser.C_tok.Value == "true")

		case "user_frame":
			frame.FrameType = USER_FRAME

		case "value":
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameType = SET_FRAME
			frame.SetValue = num

		case "bind_frame":
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.BindFrame = num
			frame.FrameType = BIND_FRAME
		case "bind_geo":
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.BindGeo = num

		case "bind_attr":
			frame.BindAttr = parser.C_tok.Value

		default:
			log.Printf("Unknown keyframe attribute %s", name)
		}

		err = parser.NextToken(buf)
		if err != nil {
			return
		}

		if parser.C_tok.Tok == ',' {
			err = parser.MatchToken(',', buf)
			if err != nil {
				return
			}
		}
	}

	return
}

func parseGeometry(buf *bufio.Reader) (geo IGeometry, err error) {
	data := make(map[string]string, 10)
	var geom Geometry
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value
		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		switch name {
		case "id":
			geom.GeoID, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			err = parser.MatchToken(parser.INT, buf)
			if err != nil {
				return
			}

		case "name":
			geom.Name = parser.C_tok.Value
			err = parser.MatchToken(parser.STRING, buf)
			if err != nil {
				return
			}

		case "geo_type":
			geom.GeoType, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			err = parser.MatchToken(parser.INT, buf)
			if err != nil {
				return
			}

		case "prop_type":
			geom.PropType, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			err = parser.MatchToken(parser.INT, buf)
			if err != nil {
				return
			}

		case "attr":
			parser.MatchToken('[', buf)

			for parser.C_tok.Tok == '{' {
				name := parser.C_tok.Value
				parser.MatchToken(parser.STRING, buf)
				parser.MatchToken(',', buf)

				data[name] = parser.C_tok.Value
				parser.NextToken(buf)

				parser.MatchToken('}', buf)

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

		default:
			log.Printf("Unknown geometry attribute %s", name)
		}
	}

	parser.MatchToken('}', buf)

	geom.RelX, err = strconv.Atoi(data["rel_x"])
	if err != nil {
		return
	}

	geom.RelY, err = strconv.Atoi(data["rel_y"])
	if err != nil {
		return
	}

	color := strings.Split(data["color"], " ")
	if len(color) != 4 {
		err = fmt.Errorf("Incorrect number of colors (%s)", data["color"])
	}

	geom.Parent, err = strconv.Atoi(data["parent"])
	if err != nil {
		return
	}

	switch geom.GeoType {
	case GEO_RECT:
		width, _ := strconv.Atoi(data["width"])
		height, _ := strconv.Atoi(data["height"])
		rounding, _ := strconv.Atoi(data["rounding"])

		geo = NewRectangle(geom, width, height, rounding, data["color"])
	case GEO_CIRCLE:
		inner, _ := strconv.Atoi(data["inner_radius"])
		outer, _ := strconv.Atoi(data["outer_radius"])
		start, _ := strconv.Atoi(data["start_angle"])
		end, _ := strconv.Atoi(data["end_angle"])

		geo = NewCircle(geom, inner, outer, start, end, data["color"])
	case GEO_TEXT:
		geo = NewText(geom, data["string"], data["color"])
	}

	return
}

func formatColor(s string) byte {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Print(err)
		return 0
	}

	return byte(f * 255)
}
