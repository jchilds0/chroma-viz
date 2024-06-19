package templates

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

/*

   Templates form the basis of pages. Each page corresponds to
   one page which specifies the shape of the page, usually with
   a number of properties that can be edited by the user.

*/

type Template struct {
	Title     string
	TempID    int64
	Layer     int
	UserFrame []UserFrame
	SetFrame  []SetFrame
	BindFrame []BindFrame
	Rectangle []Rectangle
	Circle    []Circle
	Text      []Text
	Asset     []Asset
}

func NewTemplate(title string, id int64, layer, num_key, num_geo int) *Template {
	temp := &Template{
		Title:  title,
		TempID: id,
		Layer:  layer,
	}

	temp.UserFrame = make([]UserFrame, 0, num_key)
	temp.SetFrame = make([]SetFrame, 0, num_key)
	temp.BindFrame = make([]BindFrame, 0, num_key)

	temp.Rectangle = make([]Rectangle, 0, num_geo)
	temp.Circle = make([]Circle, 0, num_geo)
	temp.Text = make([]Text, 0, num_geo)
	temp.Asset = make([]Asset, 0, num_geo)

	return temp
}

func NewTemplateFromFile(fileName string) (temp Template, err error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &temp)
	if err != nil {
		return
	}

	return
}

func (temp *Template) TemplateToListRow() (row *gtk.ListBoxRow, err error) {
	row, err = gtk.ListBoxRowNew()
	if err != nil {
		return
	}

	textView, err := TextToBuffer(temp.Title)
	if err != nil {
		return
	}

	row.Add(textView)
	return
}

func TextToBuffer(text string) (textView *gtk.TextView, err error) {
	textView, err = gtk.TextViewNew()
	if err != nil {
		return
	}

	buffer, err := textView.GetBuffer()
	if err != nil {
		return
	}

	buffer.SetText(text)
	return
}

// T -> {'id': num, 'num_geo': num, 'layer': num, 'geometry': [G]} | T, T
func (temp *Template) Encode() (s string, err error) {
	var b strings.Builder

	b.WriteString("{")

	b.WriteString("'id': ")
	b.WriteString(strconv.FormatInt(temp.TempID, 10))
	b.WriteString(", ")

	b.WriteString("'num_geo': ")
	b.WriteString(strconv.Itoa(temp.NumGeometry() + 1))
	b.WriteString(", ")

	b.WriteString("'max_keyframe': ")
	b.WriteString(strconv.Itoa(temp.MaxKeyframe()))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(temp.Title)
	b.WriteString("', ")

	b.WriteString("'layer': ")
	b.WriteString(strconv.Itoa(temp.Layer))
	b.WriteString(", ")

	{

		b.WriteString("'geometry': [")

		first := true
		var propStr string

		for _, geo := range temp.Rectangle {
			if !first {
				b.WriteString(",")
			}

			first = false
			propStr = EncodeGeometry(geo.Geometry, geo.Attributes())
			b.WriteString(propStr)
		}

		for _, geo := range temp.Circle {
			if !first {
				b.WriteString(",")
			}

			first = false
			propStr = EncodeGeometry(geo.Geometry, geo.Attributes())
			b.WriteString(propStr)
		}

		for _, geo := range temp.Text {
			if !first {
				b.WriteString(",")
			}

			first = false
			propStr = EncodeGeometry(geo.Geometry, geo.Attributes())
			b.WriteString(propStr)
		}

		for _, geo := range temp.Asset {
			if !first {
				b.WriteString(",")
			}

			first = false
			propStr = EncodeGeometry(geo.Geometry, geo.Attributes())
			b.WriteString(propStr)
		}

		b.WriteString("]")

	}

	{

		b.WriteString("'keyframe': [")

		first := true
		var frameStr string

		for _, frame := range temp.UserFrame {
			if !first {
				b.WriteString(",")
			}
			first = false

			frameStr, err = EncodeKeyframe(frame.Keyframe, frame.Attributes())
			if err != nil {
				return
			}
			b.WriteString(frameStr)
		}

		for _, frame := range temp.BindFrame {
			if !first {
				b.WriteString(",")
			}
			first = false

			frameStr, err = EncodeKeyframe(frame.Keyframe, frame.Attributes())
			if err != nil {
				return
			}
			b.WriteString(frameStr)
		}

		for _, frame := range temp.SetFrame {
			if !first {
				b.WriteString(",")
			}
			first = false

			frameStr, err = EncodeKeyframe(frame.Keyframe, frame.Attributes())
			if err != nil {
				return
			}
			b.WriteString(frameStr)
		}

		b.WriteString("],")

	}

	b.WriteString("}")
	s = b.String()
	return
}

func (temp *Template) ExportTemplate(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := json.Marshal(temp)
	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func GetTemplate(hub net.Conn, tempid int) (temp Template, err error) {
	s := fmt.Sprintf("ver 0 1 temp %d;", tempid)

	_, err = hub.Write([]byte(s))
	if err != nil {
		return
	}

	buf := bufio.NewReader(hub)
	temp, err = parseTemplate(buf)

	return temp, err
}

func (temp *Template) NumGeometry() (maxID int) {
	for _, rect := range temp.Rectangle {
		maxID = max(maxID, rect.GeoNum)
	}

	for _, circle := range temp.Circle {
		maxID = max(maxID, circle.GeoNum)
	}

	for _, text := range temp.Text {
		maxID = max(maxID, text.GeoNum)
	}

	for _, img := range temp.Asset {
		maxID = max(maxID, img.GeoNum)
	}

	return
}

func (temp *Template) MaxKeyframe() (maxFrameNum int) {
	for _, user := range temp.UserFrame {
		maxFrameNum = max(maxFrameNum, user.FrameNum)
	}

	for _, set := range temp.SetFrame {
		maxFrameNum = max(maxFrameNum, set.FrameNum)
	}

	for _, bind := range temp.BindFrame {
		maxFrameNum = max(maxFrameNum, bind.FrameNum)
	}

	return
}
