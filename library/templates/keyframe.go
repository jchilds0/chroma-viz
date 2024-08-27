package templates

import (
	"chroma-viz/library/parser"
	"strings"
)

const (
	BIND_FRAME = "bind-frame"
	SET_FRAME  = "set-frame"
	USER_FRAME = "user-frame"
)

type Keyframe struct {
	FrameNum int
	GeoID    int
	GeoAttr  string
	Type     string
	Expand   bool
}

func NewKeyFrame(num, geo int, attr string, expand bool) *Keyframe {
	frame := &Keyframe{
		FrameNum: num,
		GeoID:    geo,
		GeoAttr:  attr,
		Expand:   expand,
	}

	return frame
}

func (key *Keyframe) Key() *Keyframe {
	return key
}

func (frame Keyframe) encodeKeyframe(b strings.Builder) {
	parser.AddAttribute(b, "frame_num", frame.FrameNum)
	b.WriteString(", ")

	parser.AddAttribute(b, "frame_geo", frame.GeoID)
	b.WriteString(", ")

	parser.AddAttribute(b, "frame_attr", frame.GeoAttr)
	b.WriteString(", ")

	parser.AddAttribute(b, "frame_type", frame.Type)
	b.WriteString(", ")

	if frame.Expand {
		parser.AddAttribute(b, "expand", "true")
	} else {
		parser.AddAttribute(b, "expand", "false")
	}
}

type BindFrame struct {
	Keyframe
	Bind Keyframe
}

func NewBindFrame(frame, bind Keyframe) *BindFrame {
	frame.Type = BIND_FRAME

	return &BindFrame{
		Keyframe: frame,
		Bind:     bind,
	}
}

func (frame BindFrame) EncodeJSON(b strings.Builder) {
	b.WriteRune('{')

	frame.Keyframe.encodeKeyframe(b)

	b.WriteRune(',')
	parser.AddAttribute(b, "bind_frame", frame.Bind.FrameNum)

	b.WriteRune(',')
	parser.AddAttribute(b, "bind_geo", frame.Bind.GeoID)

	b.WriteRune(',')
	parser.AddAttribute(b, "bind_attr", frame.Bind.GeoAttr)

	b.WriteRune('}')
}

type SetFrame struct {
	Keyframe
	Value int
}

func NewSetFrame(frame Keyframe, value int) *SetFrame {
	frame.Type = SET_FRAME

	return &SetFrame{
		Keyframe: frame,
		Value:    value,
	}
}

func (set SetFrame) EncodeJSON(b strings.Builder) {
	b.WriteRune('{')

	set.Keyframe.encodeKeyframe(b)

	b.WriteRune(',')
	parser.AddAttribute(b, "value", set.Value)

	b.WriteRune('}')
}

type UserFrame struct {
	Keyframe
}

func NewUserFrame(frame Keyframe) *UserFrame {
	frame.Type = USER_FRAME

	return &UserFrame{Keyframe: frame}
}

func (user UserFrame) EncodeJSON(b strings.Builder) {
	b.WriteRune('{')

	user.Keyframe.encodeKeyframe(b)

	b.WriteRune('}')
}
