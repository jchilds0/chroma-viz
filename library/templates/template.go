package templates

import (
	"bufio"
	"chroma-viz/library/geometry"
	"chroma-viz/library/parser"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"slices"
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
	Geos      map[int]*geometry.Geometry

	Rectangle []*geometry.Rectangle
	Circle    []*geometry.Circle
	Clock     []*geometry.Clock
	Image     []*geometry.Image
	Polygon   []*geometry.Polygon
	Text      []*geometry.Text
	Ticker    []*geometry.Ticker
}

func NewTemplate(title string, id int64, layer, numKey, numGeo int) *Template {
	temp := &Template{
		Title:  title,
		TempID: id,
		Layer:  layer,
	}

	temp.UserFrame = make([]UserFrame, 0, numKey)
	temp.SetFrame = make([]SetFrame, 0, numKey)
	temp.BindFrame = make([]BindFrame, 0, numKey)

	temp.Geos = make(map[int]*geometry.Geometry, numGeo)

	temp.Rectangle = make([]*geometry.Rectangle, 0, numGeo)
	temp.Circle = make([]*geometry.Circle, 0, numGeo)
	temp.Clock = make([]*geometry.Clock, 0, numGeo)
	temp.Image = make([]*geometry.Image, 0, numGeo)
	temp.Polygon = make([]*geometry.Polygon, 0, numGeo)
	temp.Text = make([]*geometry.Text, 0, numGeo)
	temp.Ticker = make([]*geometry.Ticker, 0, numGeo)

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

	temp.Geos = make(map[int]*geometry.Geometry, 10)

	updateGeometryEntry[*geometry.Rectangle](&temp, temp.Rectangle)
	updateGeometryEntry[*geometry.Circle](&temp, temp.Circle)
	updateGeometryEntry[*geometry.Clock](&temp, temp.Clock)
	updateGeometryEntry[*geometry.Image](&temp, temp.Image)
	updateGeometryEntry[*geometry.Polygon](&temp, temp.Polygon)
	updateGeometryEntry[*geometry.Text](&temp, temp.Text)
	updateGeometryEntry[*geometry.Ticker](&temp, temp.Ticker)

	return
}

func updateGeometryEntry[T geoInterface](temp *Template, geos []T) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

		temp.Geos[geo.GetGeometryID()] = geo.GetGeometry()
	}
}

func (temp *Template) Clean() {
	numGeo := 10
	numKey := 10

	temp.Geos = make(map[int]*geometry.Geometry, numGeo)

	temp.Rectangle = make([]*geometry.Rectangle, 0, numGeo)
	temp.Circle = make([]*geometry.Circle, 0, numGeo)
	temp.Clock = make([]*geometry.Clock, 0, numGeo)
	temp.Image = make([]*geometry.Image, 0, numGeo)
	temp.Polygon = make([]*geometry.Polygon, 0, numGeo)
	temp.Text = make([]*geometry.Text, 0, numGeo)
	temp.Ticker = make([]*geometry.Ticker, 0, numGeo)

	temp.UserFrame = make([]UserFrame, 0, numKey)
	temp.SetFrame = make([]SetFrame, 0, numKey)
	temp.BindFrame = make([]BindFrame, 0, numKey)
}

func (temp *Template) AddGeometry(geoType, geoName string) (id int, err error) {
	ok := true
	for id = 1; ok; id++ {
		geo, ok := temp.Geos[id]
		if !ok {
			break
		}

		if geo == nil {
			break
		}
	}

	geo := geometry.NewGeometry(id, geoName, geoType)

	switch geoType {
	case geometry.GEO_RECT:
		rect := geometry.NewRectangle(geo)
		temp.Geos[geo.GeometryID] = &rect.Geometry
		temp.Rectangle = append(temp.Rectangle, rect)

	case geometry.GEO_CIRCLE:
		circle := geometry.NewCircle(geo)
		temp.Geos[geo.GeometryID] = &circle.Geometry
		temp.Circle = append(temp.Circle, circle)

	case geometry.GEO_TEXT:
		text := geometry.NewText(geo)
		temp.Geos[geo.GeometryID] = &text.Geometry
		temp.Text = append(temp.Text, text)

	case geometry.GEO_IMAGE:
		img := geometry.NewImage(geo)
		temp.Geos[geo.GeometryID] = &img.Geometry
		temp.Image = append(temp.Image, img)

	case geometry.GEO_POLY:
		poly := geometry.NewPolygon(geo, 10)
		temp.Geos[geo.GeometryID] = &poly.Geometry
		temp.Polygon = append(temp.Polygon)

	case geometry.GEO_TICKER:
		ticker := geometry.NewTicker(geo)
		temp.Geos[geo.GeometryID] = &ticker.Geometry
		temp.Ticker = append(temp.Ticker, ticker)

	case geometry.GEO_CLOCK:
		clock := geometry.NewClock(geo)
		temp.Geos[geo.GeometryID] = &clock.Geometry
		temp.Clock = append(temp.Clock, clock)

	default:
		err = fmt.Errorf("Error: Unknown geometry type %s", geoType)
	}

	return
}

func (temp *Template) RemoveGeometry(geoID int) {
	temp.Rectangle = removeGeometry[*geometry.Rectangle](temp.Rectangle, geoID)
	temp.Circle = removeGeometry[*geometry.Circle](temp.Circle, geoID)
	temp.Clock = removeGeometry[*geometry.Clock](temp.Clock, geoID)
	temp.Image = removeGeometry[*geometry.Image](temp.Image, geoID)
	temp.Polygon = removeGeometry[*geometry.Polygon](temp.Polygon, geoID)
	temp.Text = removeGeometry[*geometry.Text](temp.Text, geoID)
	temp.Ticker = removeGeometry[*geometry.Ticker](temp.Ticker, geoID)

	delete(temp.Geos, geoID)
}

func removeGeometry[T interface{ GetGeometryID() int }](geos []T, geoID int) (retval []T) {
	for i, geo := range geos {
		if isNil(geo) {
			continue
		}

		if geo.GetGeometryID() != geoID {
			continue
		}

		retval = slices.Delete(geos, i, i+1)
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
func (temp *Template) MarshalJSON() (buf []byte, err error) {
	var tempJSON struct {
		Template
		NumGeometry int
		NumKeyframe int
	}

	tempJSON.Template = *temp
	tempJSON.NumGeometry = temp.NumGeometry()
	tempJSON.NumKeyframe = temp.MaxKeyframe()

	return json.Marshal(tempJSON)
}

type encoder interface {
	EncodeJSON(strings.Builder)
}

func encodeSlice[T encoder](b strings.Builder, geos []T, first *bool) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

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

func GetTemplate(conn net.Conn, tempid int) (temp Template, err error) {
	s := fmt.Sprintf("ver 0 1 temp %d;", tempid)

	_, err = conn.Write([]byte(s))
	if err != nil {
		return
	}

	buf := bufio.NewReader(conn)
	data, err := buf.ReadBytes('\n')
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &temp)
	return
}

func (temp *Template) NumGeometry() (maxID int) {
	maxID = max(maxID, maxGeoNum[*geometry.Rectangle](temp.Rectangle))
	maxID = max(maxID, maxGeoNum[*geometry.Circle](temp.Circle))
	maxID = max(maxID, maxGeoNum[*geometry.Clock](temp.Clock))
	maxID = max(maxID, maxGeoNum[*geometry.Image](temp.Image))
	maxID = max(maxID, maxGeoNum[*geometry.Polygon](temp.Polygon))
	maxID = max(maxID, maxGeoNum[*geometry.Text](temp.Text))
	maxID = max(maxID, maxGeoNum[*geometry.Ticker](temp.Ticker))

	return
}

func maxGeoNum[T interface{ GetGeometryID() int }](geos []T) (maxID int) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

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

type geoInterface interface {
	GetGeometryID() int
	GetGeometry() *geometry.Geometry
}

func (temp *Template) ApplyGeometryFunc(geoID int, f func(*geometry.Geometry)) {
	applyFunction[*geometry.Rectangle](temp.Rectangle, geoID, f)
	applyFunction[*geometry.Circle](temp.Circle, geoID, f)
	applyFunction[*geometry.Clock](temp.Clock, geoID, f)
	applyFunction[*geometry.Image](temp.Image, geoID, f)
	applyFunction[*geometry.Polygon](temp.Polygon, geoID, f)
	applyFunction[*geometry.Text](temp.Text, geoID, f)
	applyFunction[*geometry.Ticker](temp.Ticker, geoID, f)
}

func applyFunction[T geoInterface](geos []T, geoID int, f func(*geometry.Geometry)) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

		if geo.GetGeometryID() != geoID {
			continue
		}

		f(geo.GetGeometry())
	}

}

func isNil(id any) bool {
	v := reflect.ValueOf(id)
	if v.Kind() != reflect.Pointer {
		log.Fatalf("Incorrect type %s", v.Kind().String())
	}

	return v.IsNil()
}
