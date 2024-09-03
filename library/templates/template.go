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
	List      []*geometry.List
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
	temp.List = make([]*geometry.List, 0, numGeo)

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
	updateGeometryEntry[*geometry.List](&temp, temp.List)

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
	temp.Title = ""
	temp.TempID = 0
	temp.Layer = 0

	numGeo := 10
	numKey := 10

	temp.Geos = make(map[int]*geometry.Geometry, numGeo)

	temp.Rectangle = make([]*geometry.Rectangle, 0, numGeo)
	temp.Circle = make([]*geometry.Circle, 0, numGeo)
	temp.Clock = make([]*geometry.Clock, 0, numGeo)
	temp.Image = make([]*geometry.Image, 0, numGeo)
	temp.Polygon = make([]*geometry.Polygon, 0, numGeo)
	temp.Text = make([]*geometry.Text, 0, numGeo)
	temp.List = make([]*geometry.List, 0, numGeo)

	temp.UserFrame = make([]UserFrame, 0, numKey)
	temp.SetFrame = make([]SetFrame, 0, numKey)
	temp.BindFrame = make([]BindFrame, 0, numKey)
}

func (temp *Template) AddGeometry(geoType, geoName string) (id int, err error) {
	if temp.Geos == nil {
		temp.Geos = make(map[int]*geometry.Geometry, 128)
	}

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
		poly := geometry.NewPolygon(geo)
		temp.Geos[geo.GeometryID] = &poly.Geometry
		temp.Polygon = append(temp.Polygon)

	case geometry.GEO_LIST:
		ticker := geometry.NewList(geo)
		temp.Geos[geo.GeometryID] = &ticker.Geometry
		temp.List = append(temp.List, ticker)

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
	temp.Rectangle = removeGeometry(temp.Rectangle, geoID)
	temp.Circle = removeGeometry(temp.Circle, geoID)
	temp.Clock = removeGeometry(temp.Clock, geoID)
	temp.Image = removeGeometry(temp.Image, geoID)
	temp.Polygon = removeGeometry(temp.Polygon, geoID)
	temp.Text = removeGeometry(temp.Text, geoID)
	temp.List = removeGeometry(temp.List, geoID)

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

	return geos
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
		Title     string
		TempID    int64
		Layer     int
		UserFrame []UserFrame
		SetFrame  []SetFrame
		BindFrame []BindFrame

		Rectangle []geometry.Rectangle
		Circle    []geometry.Circle
		Clock     []geometry.Clock
		Image     []geometry.Image
		Polygon   []geometry.Polygon
		Text      []geometry.Text
		List      []geometry.List
	}

	tempJSON.Title = temp.Title
	tempJSON.TempID = temp.TempID
	tempJSON.Layer = temp.Layer
	tempJSON.UserFrame = temp.UserFrame
	tempJSON.SetFrame = temp.SetFrame
	tempJSON.BindFrame = temp.BindFrame

	tempJSON.Rectangle = removeNil(temp.Rectangle)
	tempJSON.Circle = removeNil(temp.Circle)
	tempJSON.Clock = removeNil(temp.Clock)
	tempJSON.Image = removeNil(temp.Image)
	tempJSON.Polygon = removeNil(temp.Polygon)
	tempJSON.Text = removeNil(temp.Text)
	tempJSON.List = removeNil(temp.List)

	return json.Marshal(tempJSON)
}

func removeNil[T any](geos []*T) (retval []T) {
	retval = make([]T, 0, len(geos))

	for _, geo := range geos {
		if geo == nil {
			continue
		}

		retval = append(retval, *geo)
	}

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

func GetTemplate(conn net.Conn, tempid int) (temp Template, err error) {
	s := fmt.Sprintf("ver 0 1 temp %d;", tempid)

	_, err = conn.Write([]byte(s))
	if err != nil {
		return
	}

	buf := bufio.NewReader(conn)
	data, err := buf.ReadBytes(6)
	if err != nil {
		return
	}

	err = json.Unmarshal(data[:len(data)-1], &temp)
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
	updateGeometryEntry[*geometry.List](&temp, temp.List)

	return
}

func (temp *Template) NumGeometry() (maxID int) {
	maxID = max(maxID, maxGeoNum(temp.Rectangle))
	maxID = max(maxID, maxGeoNum(temp.Circle))
	maxID = max(maxID, maxGeoNum(temp.Clock))
	maxID = max(maxID, maxGeoNum(temp.Image))
	maxID = max(maxID, maxGeoNum(temp.Polygon))
	maxID = max(maxID, maxGeoNum(temp.Text))
	maxID = max(maxID, maxGeoNum(temp.List))

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

func (temp *Template) Encode(b *strings.Builder) {
	parser.EngineAddKeyValue(b, "temp", temp.TempID)
	parser.EngineAddKeyValue(b, "layer", temp.Layer)

	encodeGeometry(b, temp.Rectangle)
	encodeGeometry(b, temp.Circle)
	encodeGeometry(b, temp.Clock)
	encodeGeometry(b, temp.Image)
	encodeGeometry(b, temp.Polygon)
	encodeGeometry(b, temp.Text)
	encodeGeometry(b, temp.List)
}

type encoder interface {
	Encode(b *strings.Builder)
}

func encodeGeometry[T encoder](b *strings.Builder, geos []T) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

		geo.Encode(b)
	}
}

type geoInterface interface {
	GetGeometryID() int
	GetGeometry() *geometry.Geometry
}

func (temp *Template) ApplyGeometryFunc(geoID int, f func(*geometry.Geometry)) {
	applyFunction(temp.Rectangle, geoID, f)
	applyFunction(temp.Circle, geoID, f)
	applyFunction(temp.Clock, geoID, f)
	applyFunction(temp.Image, geoID, f)
	applyFunction(temp.Polygon, geoID, f)
	applyFunction(temp.Text, geoID, f)
	applyFunction(temp.List, geoID, f)
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
