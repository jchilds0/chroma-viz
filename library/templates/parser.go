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

		if name == "geometry" {
			parser.MatchToken('[', buf)

			var num_geo, num_keyframe, temp_id, layer int
			num_geo, err = strconv.Atoi(data["num_geo"])
			if err != nil {
				return
			}

			num_keyframe, err = strconv.Atoi(data["num_keyframe"])
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

			temp = NewTemplate(data["name"], temp_id, layer, num_geo, num_keyframe)
			temp.Geometry, err = parser.ParseProperty(buf, num_geo)
			if err != nil {
				return
			}

			parser.MatchToken(']', buf)
		} else if name == "keyframe" {
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
		} else {
			data[name] = parser.C_tok.Value

			parser.NextToken(buf)
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

		if name == "frame_num" {
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameNum = num
		} else if name == "frame_geo" {
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameGeo = num
		} else if name == "frame_attr" {
			frame.FrameAttr = parser.C_tok.Value
		} else if name == "mask" {
			frame.Mask = (parser.C_tok.Value == "true")
		} else if name == "expand" {
			frame.Expand = (parser.C_tok.Value == "true")
		} else if name == "user_frame" {
			frame.FrameType = USER_FRAME
		} else if name == "value" {
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.FrameType = SET_FRAME
			frame.SetValue = num
		} else if name == "bind_frame" {
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.BindFrame = num
			frame.FrameType = BIND_FRAME
		} else if name == "bind_geo" {
			num, err = strconv.Atoi(parser.C_tok.Value)
			if err != nil {
				return
			}

			frame.BindGeo = num
		} else if name == "bind_attr" {
			frame.BindAttr = parser.C_tok.Value
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
