package pages

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

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
	templates.Template
	PageNum int
	Title   string
	lock    sync.Mutex
}

func NewPage(temp *templates.Template) *Page {
	page := &Page{
		Template: *temp,
		Title:    temp.Title,
		PageNum:  1,
	}

	return page
}

func (page *Page) GetGeometry(geoID int) *geometry.Geometry {
	return page.Geos[geoID]
}

func (page *Page) PageToListRow() (row *gtk.ListBoxRow, err error) {
	row, err = gtk.ListBoxRowNew()
	if err != nil {
		return
	}

	pageText, err := TextToBuffer(strconv.Itoa(page.PageNum))
	if err != nil {
		return
	}

	row.Add(pageText)

	titleText, err := TextToBuffer(page.Title)
	if err != nil {
		return
	}

	row.Add(titleText)

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

func (page *Page) ImportPage(filename string) error {
	page.lock.Lock()
	defer page.lock.Unlock()

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
	page.lock.Lock()
	defer page.lock.Unlock()

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
	page.lock.Lock()
	defer page.lock.Unlock()

	util.EngineAddKeyValue(b, "temp", page.TempID)
	util.EngineAddKeyValue(b, "layer", page.Layer)

	encodeGeometry(b, page.Rectangle)
	encodeGeometry(b, page.Circle)
	encodeGeometry(b, page.Clock)
	encodeGeometry(b, page.Image)
	encodeGeometry(b, page.Polygon)
	encodeGeometry(b, page.Text)
	encodeGeometry(b, page.List)
}

func encodeGeometry[T geometry.Geometer[S], S any](b *strings.Builder, geos []T) {
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
