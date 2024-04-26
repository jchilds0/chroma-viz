package templates

import (
	"strconv"
	"strings"
)

const (
	SET_FRAME = iota
	USER_FRAME
	BIND_FRAME
)

type Keyframe struct {
	FrameNum int
	GeoID    int
	GeoAttr  int
	Type     int
	Mask     bool
	Expand   bool
}

func NewKeyFrame(num, geo, attr, ftype int, mask, expand bool) *Keyframe {
	frame := &Keyframe{
		FrameNum: num,
		GeoID:    geo,
		GeoAttr:  attr,
		Type:     ftype,
		Mask:     mask,
		Expand:   expand,
	}

	return frame
}

func (key *Keyframe) Key() *Keyframe {
	return key
}

func EncodeKeyframe(frame Keyframe, attr map[string]string) (s string, err error) {
	var b strings.Builder
	key := frame.Key()

	b.WriteString("{")

	b.WriteString("'frame_num': ")
	b.WriteString(strconv.Itoa(key.FrameNum))
	b.WriteString(", ")

	b.WriteString("'frame_geo': ")
	b.WriteString(strconv.Itoa(key.GeoID))
	b.WriteString(", ")

	b.WriteString("'frame_attr': ")
	b.WriteString(strconv.Itoa(key.GeoAttr))
	b.WriteString(", ")

	b.WriteString("'frame_type': ")
	b.WriteString(strconv.Itoa(key.Type))
	b.WriteString(", ")

	b.WriteString("'mask': ")
	if key.Mask {
		b.WriteString("'true'")
	} else {
		b.WriteString("'false'")
	}
	b.WriteString(", ")

	b.WriteString("'expand': ")
	if key.Expand {
		b.WriteString("'true'")
	} else {
		b.WriteString("'false'")
	}
	b.WriteString(", ")

	first := true
	for name, value := range attr {
		if !first {
			b.WriteString(",")
		}

		first = false
		b.WriteString("'")
		b.WriteString(name)
		b.WriteString("': ")
		b.WriteString(value)
	}

	b.WriteString("}")
	s = b.String()
	return
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

func (frame *BindFrame) Attributes() map[string]string {
	bind := frame.Bind

	return map[string]string{
		"bind_frame": strconv.Itoa(bind.FrameNum),
		"bind_geo":   strconv.Itoa(bind.GeoID),
		"bind_attr":  strconv.Itoa(bind.GeoAttr),
	}
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

func (set *SetFrame) Attributes() map[string]string {
	return map[string]string{
		"value": strconv.Itoa(set.Value),
	}
}

type UserFrame struct {
	Keyframe
}

func NewUserFrame(frame Keyframe) *UserFrame {
	frame.Type = USER_FRAME

	return &UserFrame{Keyframe: frame}
}

func (user *UserFrame) Attributes() map[string]string {
	p := map[string]string{
		"user_frame": "'true'",
	}

	return p
}
