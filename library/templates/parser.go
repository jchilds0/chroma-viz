package templates

import (
	"bufio"
	"chroma-viz/library/geometry"
	"chroma-viz/library/parser"
	"fmt"
	"log"
	"strconv"
)

// T -> {'id': 123, 'num_geo': 123, 'layer': 123, 'geometry': [G]}
func parseTemplate(buf *bufio.Reader) (temp Template, err error) {
	parser.NextToken(buf)
	parser.MatchToken('{', buf)

	var tempID int64
	var title string
	layer := -1
	numKey := -1
	numGeo := -1

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value

		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		switch name {
		case "id":
			tempID, err = strconv.ParseInt(parser.C_tok.Value, 10, 64)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)

		case "name":
			title = parser.C_tok.Value
			parser.MatchToken(parser.STRING, buf)

		case "layer":
			layer, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			parser.MatchToken(parser.INT, buf)

		case "max_keyframe":
			numKey, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				numKey = 10
			}

			parser.MatchToken(parser.INT, buf)

		case "num_geo":
			numGeo, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				numGeo = 10
			}

			parser.MatchToken(parser.INT, buf)

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

		if tempID != 0 && title != "" && numGeo != -1 && numKey != -1 && layer != -1 {
			temp = *NewTemplate(title, tempID, layer, numKey, numGeo)
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

	expand := (data["expand"] == "true")

	keyframe := NewKeyFrame(frameNum, geoID, data["frame_attr"], expand)

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

		bind := NewKeyFrame(bindNum, bindGeo, data["bind_attr"], false)
		frame := NewBindFrame(*keyframe, *bind)

		temp.BindFrame = append(temp.BindFrame, *frame)
	}

	return
}

func parseGeometry(temp *Template, buf *bufio.Reader) (err error) {
	data := make(map[string]string, 10)
	var geom geometry.Geometry
	parser.MatchToken('{', buf)

	for parser.C_tok.Tok == parser.STRING {
		name := parser.C_tok.Value
		parser.MatchToken(parser.STRING, buf)
		parser.MatchToken(':', buf)

		switch name {
		case "id":
			geom.GeometryID, err = strconv.Atoi(parser.C_tok.Value)
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

	geom.RelX.Value, err = strconv.Atoi(data["rel_x"])
	if err != nil {
		return
	}

	geom.RelY.Value, err = strconv.Atoi(data["rel_y"])
	if err != nil {
		return
	}

	geom.Parent.Value, err = strconv.Atoi(data["parent"])
	if err != nil {
		return
	}

	geom.Mask.Value, err = strconv.Atoi(data["mask"])
	if err != nil {
		return
	}

	switch geom.GeoType {
	case geometry.GEO_RECT:
		rect := geometry.NewRectangle(geom)
		rect.Width.Value, _ = strconv.Atoi(data["width"])
		rect.Height.Value, _ = strconv.Atoi(data["height"])
		rect.Rounding.Value, _ = strconv.Atoi(data["rounding"])
		rect.Color.FromString(data["color"])

		temp.Rectangle = append(temp.Rectangle, rect)

	case geometry.GEO_CIRCLE:
		circle := geometry.NewCircle(geom)
		circle.InnerRadius.Value, _ = strconv.Atoi(data["inner_radius"])
		circle.OuterRadius.Value, _ = strconv.Atoi(data["outer_radius"])
		circle.StartAngle.Value, _ = strconv.Atoi(data["start_angle"])
		circle.EndAngle.Value, _ = strconv.Atoi(data["end_angle"])
		circle.Color.FromString(data["color"])

		temp.Circle = append(temp.Circle, circle)

	case geometry.GEO_TEXT:
		text := geometry.NewText(geom)
		text.Scale.Value, _ = strconv.ParseFloat(data["scale"], 64)
		text.String.Value, _ = data["string"]
		text.Color.FromString(data["color"])

		temp.Text = append(temp.Text, text)

	case geometry.GEO_IMAGE:
		img := geometry.NewImage(geom)
		img.Scale.Value, _ = strconv.ParseFloat(data["scale"], 64)
		img.Image.Value, _ = strconv.Atoi(data["image_id"])

		temp.Image = append(temp.Image, img)
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
