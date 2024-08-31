package pages

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/parser"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

/*

   Pages are the highest object in the graphics hierarchy.
   Pages consist of a number of Properties, which represent
   components of the graphic such at Title, Background or Clock.
   The Properties of a Page are defined by the template which
   the Page is built from.

   See props/props.go for information about Properties.

*/

type Page struct {
	PageNum    int
	Title      string
	TemplateID int
	Layer      int
	geo        map[int]*geometry.Geometry

	Rect   []*geometry.Rectangle
	Circle []*geometry.Circle
	Clock  []*geometry.Clock
	Image  []*geometry.Image
	Poly   []*geometry.Polygon
	Text   []*geometry.Text
	Ticker []*geometry.Ticker
}

func NewPage(pageNum, tempID, layer, numGeo int, title string) *Page {
	page := &Page{
		PageNum:    pageNum,
		Title:      title,
		TemplateID: tempID,
		Layer:      layer,
	}

	page.geo = make(map[int]*geometry.Geometry, numGeo)

	page.Rect = make([]*geometry.Rectangle, 0, numGeo)
	page.Circle = make([]*geometry.Circle, 0, numGeo)
	page.Clock = make([]*geometry.Clock, 0, numGeo)
	page.Image = make([]*geometry.Image, 0, numGeo)
	page.Poly = make([]*geometry.Polygon, 0, numGeo)
	page.Text = make([]*geometry.Text, 0, numGeo)
	page.Ticker = make([]*geometry.Ticker, 0, numGeo)

	return page
}

func NewPageFromTemplate(temp *templates.Template) *Page {
	page := NewPage(0, int(temp.TempID), temp.Layer, temp.NumGeometry(), temp.Title)

	page.Rect = temp.Rectangle
	page.Circle = temp.Circle
	page.Clock = temp.Clock
	page.Image = temp.Image
	page.Poly = temp.Polygon
	page.Text = temp.Text
	page.Ticker = temp.Ticker

	return page
}

func (page *Page) GetGeometry(geoID int) *geometry.Geometry {
	return page.geo[geoID]
}

type geoInterface interface {
	GetGeometryID() int
	GetGeometry() *geometry.Geometry
}

func AddGeometry(page *Page, geo geoInterface) (err error) {
	page.geo[geo.GetGeometryID()] = geo.GetGeometry()

	switch g := geo.(type) {
	case *geometry.Rectangle:
		page.Rect = append(page.Rect, g)

	case *geometry.Circle:
		page.Circle = append(page.Circle, g)

	case *geometry.Clock:
		page.Clock = append(page.Clock, g)

	case *geometry.Image:
		page.Image = append(page.Image, g)

	case *geometry.Polygon:
		page.Poly = append(page.Poly, g)

	case *geometry.Text:
		page.Text = append(page.Text, g)

	case *geometry.Ticker:
		page.Ticker = append(page.Ticker, g)

	default:
		err = fmt.Errorf("Unknown type to add to page")
	}

	return
}

func (page *Page) PageToListRow() (row *gtk.ListBoxRow, err error) {
	row, err = gtk.ListBoxRowNew()
	if err != nil {
		return
	}

	pageText, err := templates.TextToBuffer(strconv.Itoa(page.PageNum))
	if err != nil {
		return
	}

	row.Add(pageText)

	titleText, err := templates.TextToBuffer(page.Title)
	if err != nil {
		return
	}

	row.Add(titleText)

	return
}

func (page *Page) ImportPage(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, page)
	if err != nil {
		return err
	}

	return nil
}

func ExportPage(page *Page, filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := json.Marshal(page)
	if err != nil {
		return err
	}

	_, err = file.Write(buf)
	if err != nil {
		return err
	}
	return
}

func (page *Page) Encode(b *strings.Builder) {
	parser.EngineAddKeyValue(b, "temp", page.TemplateID)
	parser.EngineAddKeyValue(b, "layer", page.Layer)

	encodeGeometry(b, page.Rect)
	encodeGeometry(b, page.Circle)
	encodeGeometry(b, page.Clock)
	encodeGeometry(b, page.Image)
	encodeGeometry(b, page.Poly)
	encodeGeometry(b, page.Text)
	encodeGeometry(b, page.Ticker)
}

type encoder interface {
	Encode(b *strings.Builder)
}

func encodeGeometry[T encoder](b *strings.Builder, geos []T) {
	if isNil(geos) {
		return
	}

	for _, geo := range geos {
		if isNil(geo) {
			continue
		}

		geo.Encode(b)
	}
}

func isNil(id any) bool {
	v := reflect.ValueOf(id)
	if v.Kind() != reflect.Pointer && v.Kind() != reflect.Slice {
		log.Fatalf("Incorrect type %s", v.Kind().String())
	}

	return v.IsNil()
}
