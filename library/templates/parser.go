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

		case "name":
			temp.Title = parser.C_tok.Value
			parser.MatchToken(parser.STRING, buf)

		case "layer":
			temp.Layer, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)

		case "num_keyframe":
			num, err := strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				num = 10
			}

			parser.MatchToken(parser.INT, buf)
			temp.BindFrame = make([]BindFrame, 0, num)
			temp.SetFrame = make([]SetFrame, 0, num)
			temp.UserFrame = make([]UserFrame, 0, num)

		case "num_geo":
			num, err := strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				num = 10
			}

			parser.MatchToken(parser.INT, buf)
			temp.Rectangle = make([]Rectangle, 0, num)
			temp.Circle = make([]Circle, 0, num)
			temp.Text = make([]Text, 0, num)

		case "keyframe":
			parser.MatchToken('[', buf)

			for parser.C_tok.Tok == '{' {
				parseKeyframe(&temp, buf)
				if err != nil {
					return
				}

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

			parser.MatchToken(']', buf)

		case "geometry":
			parser.MatchToken('[', buf)

			for parser.C_tok.Tok == '{' {
				err = parseGeometry(&temp, buf)
				if err != nil {
					return
				}

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

			parser.MatchToken(']', buf)

		default:
			log.Printf("Unknown template attribute %s", name)
			parser.NextToken(buf)
		}

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
		}
	}

	return
}

func parseKeyframe(temp *Template, buf *bufio.Reader) {
	data := make(map[string]string, 10)
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value

		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		value := parser.C_tok.Value
		parser.NextToken(buf)

		data[name] = value

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
		}
	}

	parser.MatchToken('}', buf)

	frameNum, _ := strconv.Atoi(data["frame_num"])
	geoID, _ := strconv.Atoi(data["frame_geo"])

	mask := (data["mask"] == "true")
	expand := (data["expand"] == "true")

	keyframe := NewKeyFrame(frameNum, geoID, data["frame_attr"], mask, expand)

	if data["user_frame"] == "true" {
		frame := NewUserFrame(*keyframe)
		temp.UserFrame = append(temp.UserFrame, *frame)
	} else if _, ok := data["value"]; ok {
		value, _ := strconv.Atoi(data["value"])
		frame := NewSetFrame(*keyframe, value)
		temp.SetFrame = append(temp.SetFrame, *frame)
	} else {
		bindNum, _ := strconv.Atoi(data["bind_frame"])
		bindGeo, _ := strconv.Atoi(data["bind_geo"])

		bind := NewKeyFrame(bindNum, bindGeo, data["bind_attr"], false, false)
		frame := NewBindFrame(*keyframe, *bind)

		temp.BindFrame = append(temp.BindFrame, *frame)
	}

	return
}

func parseGeometry(temp *Template, buf *bufio.Reader) (err error) {
	data := make(map[string]string, 10)
	var geom Geometry
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value
		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		switch name {
		case "id":
			geom.GeoNum, err = strconv.Atoi(parser.C_tok.Value)
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
			geom.GeoType = parser.C_tok.Value
			err = parser.MatchToken(parser.STRING, buf)
			if err != nil {
				return
			}

		case "prop_type":
			geom.PropType = parser.C_tok.Value
			err = parser.MatchToken(parser.STRING, buf)
			if err != nil {
				return
			}

		case "attr":
			parser.MatchToken('[', buf)

			for parser.C_tok.Tok == '{' {
				parser.MatchToken('{', buf)

				if parser.C_tok.Value != "name" {
					err = fmt.Errorf("Incorrect attribute %s", parser.C_tok.Value)
					return
				}

				parser.MatchToken(parser.STRING, buf)
				parser.MatchToken(':', buf)

				name := parser.C_tok.Value
				parser.MatchToken(parser.STRING, buf)

				parser.MatchToken(',', buf)

				if parser.C_tok.Value != "value" {
					err = fmt.Errorf("Incorrect attribute %s", parser.C_tok.Value)
					return
				}

				parser.MatchToken(parser.STRING, buf)
				parser.MatchToken(':', buf)

				data[name] = parser.C_tok.Value
				parser.NextToken(buf)

				parser.MatchToken('}', buf)

				if parser.C_tok.Tok == ',' {
					parser.MatchToken(',', buf)
				}
			}

			parser.MatchToken(']', buf)
		default:
			log.Printf("Unknown geometry attribute %s", name)
			parser.NextToken(buf)
		}

		if parser.C_tok.Tok == ',' {
			parser.MatchToken(',', buf)
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

		rect := NewRectangle(geom, width, height, rounding, data["color"])
		temp.Rectangle = append(temp.Rectangle, *rect)

	case GEO_CIRCLE:
		inner, _ := strconv.Atoi(data["inner_radius"])
		outer, _ := strconv.Atoi(data["outer_radius"])
		start, _ := strconv.Atoi(data["start_angle"])
		end, _ := strconv.Atoi(data["end_angle"])

		circle := NewCircle(geom, inner, outer, start, end, data["color"])
		temp.Circle = append(temp.Circle, *circle)

	case GEO_TEXT:
		text := NewText(geom, data["string"], data["color"])
		temp.Text = append(temp.Text, *text)

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
