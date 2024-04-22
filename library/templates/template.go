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
	Title    string
	TempID   int64
	Layer    int
	Keyframe []Keyframe
	Geometry []IGeometry
}

func NewTemplate(title string, id int64, layer, num_keyframe, num_geo int) *Template {
	temp := &Template{
		Title:  title,
		TempID: id,
		Layer:  layer,
	}

	temp.Keyframe = make([]Keyframe, 0, num_keyframe)
	temp.Geometry = make([]IGeometry, 0, num_geo)
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
	for _, geo := range temp.Geometry {
		if !first {
			b.WriteString(",")
		}

		first = false
		propStr = EncodeGeometry(geo)
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
