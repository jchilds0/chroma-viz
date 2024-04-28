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

const (
	ATTR_COLOR = iota
	ATTR_POS_X
	ATTR_POS_Y
	ATTR_REL_X
	ATTR_REL_Y
	ATTR_PARENT
	ATTR_WIDTH
	ATTR_HEIGHT
	ATTR_ROUNDING
	ATTR_INNER_RADIUS
	ATTR_OUTER_RADIUS
	ATTR_START_ANGLE
	ATTR_END_ANGLE
	ATTR_TEXT
	ATTR_SCALE
	ATTR_GRAPH_NODE
	ATTR_NUM_NODE
	ATTR_GRAPH_TYPE
	ATTR_IMAGE_ID
)

var StringToAttr = map[string]int{
	"pos_x":        ATTR_POS_X,
	"pos_y":        ATTR_POS_Y,
	"rel_x":        ATTR_REL_X,
	"rel_y":        ATTR_REL_Y,
	"parent":       ATTR_PARENT,
	"width":        ATTR_WIDTH,
	"height":       ATTR_HEIGHT,
	"rounding":     ATTR_ROUNDING,
	"inner_radius": ATTR_INNER_RADIUS,
	"outer_radius": ATTR_OUTER_RADIUS,
	"start_angle":  ATTR_START_ANGLE,
	"end_angle":    ATTR_END_ANGLE,
}

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
