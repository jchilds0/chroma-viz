package templates

import (
	"bufio"
	"chroma-viz/library/props"
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
	Title       string
	TempID      int
	NumGeo      int
	NumKeyframe int
	Layer       int
	Keyframe    []Keyframe
	Geometry    map[int]*props.Property
}

func NewTemplate(title string, id, layer, num_geo, num_keyframe int) *Template {
	temp := &Template{
		Title:       title,
		TempID:      id,
		Layer:       layer,
		NumGeo:      num_geo,
		NumKeyframe: num_keyframe,
	}

	temp.Keyframe = make([]Keyframe, 0, num_keyframe)
	temp.Geometry = make(map[int]*props.Property, num_geo)
	return temp
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

func (temp *Template) AddGeometry(name string, geo_id, typed int, visible map[string]bool) *props.Property {
	temp.Geometry[geo_id] = props.NewProperty(typed, name, true, visible)
	temp.NumGeo++

	return temp.Geometry[geo_id]
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
	b.WriteString(strconv.Itoa(temp.TempID))
	b.WriteString(", ")

	b.WriteString("'num_geo': ")
	b.WriteString(strconv.Itoa(len(temp.Geometry)))
	b.WriteString(", ")

	b.WriteString("'num_keyframe': ")
	b.WriteString(strconv.Itoa(len(temp.Keyframe)))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(temp.Title)
	b.WriteString("', ")

	b.WriteString("'layer': ")
	b.WriteString(strconv.Itoa(temp.Layer))
	b.WriteString(", ")

	b.WriteString("'keyframe': [")
	first := true
	var frameStr string
	for _, frame := range temp.Keyframe {
		if !first {
			b.WriteString(",")
		}
		first = false

		frameStr, err = frame.Encode()
		if err != nil {
			return
		}
		b.WriteString(frameStr)
	}
	b.WriteString("],")

	b.WriteString("'geometry': [")
	first = true
	var propStr string
	for geo_id, prop := range temp.Geometry {
		if !first {
			b.WriteString(",")
		}

		first = false

		propStr, err = prop.Encode(geo_id)
		if err != nil {
			return
		}

		b.WriteString(propStr)
	}

	b.WriteString("]}")
	s = b.String()
	return
}

func ExportTemplate(temp *Template, filename string) error {
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

func (temp *Template) GetTemplateID() int {
	return temp.TempID
}

func (temp *Template) GetLayer() int {
	return temp.Layer
}

func (temp *Template) GetPropMap() map[int]*props.Property {
	return temp.Geometry
}

func GetTemplate(hub net.Conn, tempid int) (*Template, error) {
	s := fmt.Sprintf("ver 0 1 temp %d;", tempid)

	_, err := hub.Write([]byte(s))
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(hub)
	temp, err := parseTemplate(buf)

	return temp, err
}
