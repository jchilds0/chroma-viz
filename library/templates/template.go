package templates

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
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

func NewTemplateFromFile(fileName string, removeNonVisible bool) (temp Template, err error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &temp)
	if err != nil {
		return
	}

	err = temp.Init(removeNonVisible)
	return
}

func updateGeometryEntry[T geometry.Geometer[S], S any](temp *Template, geos []T) {
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

	geo := geometry.NewGeometry(id, geoName, geoType, true)

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
		temp.Polygon = append(temp.Polygon, poly)

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

func (temp *Template) CopyGeometry(fromGeoID, toGeoID int) (err error) {
	fromGeo := temp.Geos[fromGeoID]
	toGeo := temp.Geos[toGeoID]

	if fromGeo.GeoType != toGeo.GeoType {
		err = fmt.Errorf("Copy geometry from %s to %s is invalid", fromGeo.GeoType, toGeo.GeoType)
		return
	}

	copyGeometry(temp.Rectangle, fromGeoID, toGeoID, geometry.NewRectangleEditor)
	copyGeometry(temp.Circle, fromGeoID, toGeoID, geometry.NewCircleEditor)
	copyGeometry(temp.Text, fromGeoID, toGeoID, geometry.NewTextEditor)
	copyGeometry(temp.Image, fromGeoID, toGeoID, geometry.NewImageEditor)
	copyGeometry(temp.Polygon, fromGeoID, toGeoID, geometry.NewPolygonEditor)
	copyGeometry(temp.Clock, fromGeoID, toGeoID, geometry.NewClockEditor)
	copyGeometry(temp.List, fromGeoID, toGeoID, geometry.NewListEditor)

	return
}

func copyGeometry[T geometry.Geometer[S], S geometry.Editor[T]](
	geos []T, geoID1, geoID2 int, init func() (S, error)) {

	var geo1, geo2 T

	for _, geo := range geos {
		if geo.GetGeometryID() == geoID1 {
			geo1 = geo
		}

		if geo.GetGeometryID() == geoID2 {
			geo2 = geo
		}
	}

	if isNil(geo1) || isNil(geo2) {
		return
	}

	editor, err := init()
	if err != nil {
		log.Print(err)
		return
	}

	err = editor.UpdateEditor(geo1)
	if err != nil {
		log.Print(err)
		return
	}

	err = geo2.UpdateGeometry(editor)
	if err != nil {
		log.Print(err)
		return
	}
}

// T -> {'id': num, 'num_geo': num, 'layer': num, 'geometry': [G]} | T, T
func (temp *Template) MarshalJSON() (buf []byte, err error) {
	var tempJSON struct {
		Title       string
		TempID      int64
		Layer       int
		MaxGeometry int
		MaxKeyframe int

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

	tempJSON.MaxGeometry = temp.NumGeometry() + 1
	tempJSON.MaxKeyframe = temp.MaxKeyframe() + 1

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

func filterGeometry[T geometry.Geometer[S], S any](geo T) bool {
	geom := geo.GetGeometry()
	return geom.Visible
}

func (temp *Template) Init(removeNonVisible bool) (err error) {
	temp.Geos = make(map[int]*geometry.Geometry, 10)

	if removeNonVisible {
		temp.Rectangle = util.Filter(temp.Rectangle, filterGeometry)
		temp.Circle = util.Filter(temp.Circle, filterGeometry)
		temp.Clock = util.Filter(temp.Clock, filterGeometry)
		temp.Image = util.Filter(temp.Image, filterGeometry)
		temp.Polygon = util.Filter(temp.Polygon, filterGeometry)
		temp.Text = util.Filter(temp.Text, filterGeometry)
		temp.List = util.Filter(temp.List, filterGeometry)
	}

	updateGeometryEntry(temp, temp.Rectangle)
	updateGeometryEntry(temp, temp.Circle)
	updateGeometryEntry(temp, temp.Clock)
	updateGeometryEntry(temp, temp.Image)
	updateGeometryEntry(temp, temp.Polygon)
	updateGeometryEntry(temp, temp.Text)
	updateGeometryEntry(temp, temp.List)

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
	util.EngineAddKeyValue(b, "temp", temp.TempID)
	util.EngineAddKeyValue(b, "layer", temp.Layer)

	encodeGeometry(b, temp.Rectangle)
	encodeGeometry(b, temp.Circle)
	encodeGeometry(b, temp.Clock)
	encodeGeometry(b, temp.Image)
	encodeGeometry(b, temp.Polygon)
	encodeGeometry(b, temp.Text)
	encodeGeometry(b, temp.List)
}

func encodeGeometry[T geometry.Geometer[S], S any](b *strings.Builder, geos []T) {
	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

		geo.Encode(b)
	}
}

func isNil(id any) bool {
	v := reflect.ValueOf(id)
	if v.Kind() != reflect.Pointer {
		log.Fatalf("Incorrect type %s", v.Kind().String())
	}

	return v.IsNil()
}
