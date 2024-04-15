package templates

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SET_FRAME = iota
	BIND_FRAME
	USER_FRAME
)

type Keyframe struct {
	FrameNum  int
	FrameGeo  int
	FrameAttr string
	FrameType int
	SetValue  int 
	UserValue bool
	BindFrame int
	BindGeo   int
	BindAttr  string
}

func NewKeyFrame(num, geo int, attr string, ftype int) *Keyframe {
	frame := &Keyframe{
		FrameNum:  num,
		FrameGeo:  geo,
		FrameAttr: attr,
		FrameType: ftype,
	}

	return frame
}

func (frame *Keyframe) Encode() (s string, err error) {
	var b strings.Builder
	b.WriteString("{")

	b.WriteString("'frame_num': ")
	b.WriteString(strconv.Itoa(frame.FrameNum))
	b.WriteString(", ")

	b.WriteString("'frame_geo': ")
	b.WriteString(strconv.Itoa(frame.FrameGeo))
	b.WriteString(", ")

	b.WriteString("'frame_attr': '")
	b.WriteString(frame.FrameAttr)
	b.WriteString("', ")

	switch frame.FrameType {
	case USER_FRAME:
		b.WriteString("'user_frame': ")

		if frame.UserValue {
			b.WriteString("'true'")
		} else {
			b.WriteString("'false'")
		}
	case BIND_FRAME:
		b.WriteString("'bind_frame': ")
		b.WriteString(strconv.Itoa(frame.BindFrame))
		b.WriteString(", ")

		b.WriteString("'bind_geo': ")
		b.WriteString(strconv.Itoa(frame.BindGeo))
		b.WriteString(", ")

		b.WriteString("'bind_attr': '")
		b.WriteString(frame.BindAttr)
		b.WriteString("'")
	case SET_FRAME:
		b.WriteString("'value': ")
		b.WriteString(strconv.Itoa(frame.SetValue))

	default:
		err = fmt.Errorf("Unknown frame type %d", frame.FrameType)
	}

	b.WriteString("}")
	s = b.String()
	return
}
