package templates

import (
	"bufio"
	"chroma-viz/library/geometry"
	"chroma-viz/library/parser"
	"encoding/json"
	"fmt"
	"net"
	"os"
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

	Rect   []*geometry.Rectangle
	Circle []*geometry.Circle
	Clock  []*geometry.Clock
	Image  []*geometry.Image
	Poly   []*geometry.Polygon
	Text   []*geometry.Text
	Ticker []*geometry.Ticker
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

	temp.Rect = make([]*geometry.Rectangle, 0, 10)
	temp.Circle = make([]*geometry.Circle, 0, 10)
	temp.Clock = make([]*geometry.Clock, 0, 10)
	temp.Image = make([]*geometry.Image, 0, 10)
	temp.Poly = make([]*geometry.Polygon, 0, 10)
	temp.Text = make([]*geometry.Text, 0, 10)
	temp.Ticker = make([]*geometry.Ticker, 0, 10)

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
func (temp *Template) Encode(b strings.Builder) {
	b.WriteString("{")

	parser.AddAttribute(b, "id", temp.TempID)
	b.WriteString(", ")

	parser.AddAttribute(b, "num_geo", temp.NumGeometry()+1)
	b.WriteString(", ")

	parser.AddAttribute(b, "max_keyframe", temp.MaxKeyframe())
	b.WriteString(", ")

	parser.AddAttribute(b, "name", temp.Title)
	b.WriteString("', ")

	parser.AddAttribute(b, "layer", temp.Layer)
	b.WriteString(", ")

	{

		b.WriteString("'geometry': [")

		first := true

		encodeSlice(b, temp.Rect, &first)
		encodeSlice(b, temp.Circle, &first)
		encodeSlice(b, temp.Clock, &first)
		encodeSlice(b, temp.Image, &first)
		encodeSlice(b, temp.Poly, &first)
		encodeSlice(b, temp.Text, &first)
		encodeSlice(b, temp.Ticker, &first)

		b.WriteString("]")

	}

	{

		b.WriteString("'keyframe': [")

		first := true
		encodeSlice(b, temp.UserFrame, &first)
		encodeSlice(b, temp.SetFrame, &first)
		encodeSlice(b, temp.BindFrame, &first)

		b.WriteString("],")

	}

	b.WriteString("}")
}

type encoder interface {
	EncodeJSON(strings.Builder)
}

func encodeSlice[T encoder](b strings.Builder, geos []T, first *bool) {
	for _, geo := range geos {
		if !(*first) {
			b.WriteString(",")
		}
		*first = false

		geo.EncodeJSON(b)
	}
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
	maxID = max(maxID, maxGeoNum[*geometry.Rectangle](temp.Rect))
	maxID = max(maxID, maxGeoNum[*geometry.Circle](temp.Circle))
	maxID = max(maxID, maxGeoNum[*geometry.Clock](temp.Clock))
	maxID = max(maxID, maxGeoNum[*geometry.Image](temp.Image))
	maxID = max(maxID, maxGeoNum[*geometry.Polygon](temp.Poly))
	maxID = max(maxID, maxGeoNum[*geometry.Text](temp.Text))
	maxID = max(maxID, maxGeoNum[*geometry.Ticker](temp.Ticker))

	return
}

func maxGeoNum[T interface{ GetGeometryID() int }](geos []T) (maxID int) {
	for _, geo := range geos {
		maxID = max(maxID, geo.GetGeometryID())
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

func (temp *Template) EncodeEngine(b strings.Builder) {
	parser.EngineAddKeyValue(b, "temp", temp.TempID)
	parser.EngineAddKeyValue(b, "layer", temp.Layer)
}
