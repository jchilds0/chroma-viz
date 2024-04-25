package pages

import (
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"encoding/json"
	"os"
	"strconv"

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
	PropMap    map[int]*props.Property
}

func NewPage(pageNum, tempID, layer, numGeo int, title string) *Page {
	page := &Page{
		PageNum:    pageNum,
		Title:      title,
		TemplateID: tempID,
		Layer:      layer,
	}

	page.PropMap = make(map[int]*props.Property, numGeo)
	return page
}

func NewPageFromTemplate(temp *templates.Template) (page *Page) {
	numGeo := len(temp.Rectangle) + len(temp.Circle) + len(temp.Text)
	page = NewPage(0, int(temp.TempID), temp.Layer, numGeo, temp.Title)

	for _, geo := range temp.Rectangle {
		page.PropMap[geo.GeoID] = props.NewPropertyFromGeometry(geo.Geometry, geo.Attributes())
	}

	for _, geo := range temp.Text {
		page.PropMap[geo.GeoID] = props.NewPropertyFromGeometry(geo.Geometry, geo.Attributes())
	}

	for _, geo := range temp.Circle {
		page.PropMap[geo.GeoID] = props.NewPropertyFromGeometry(geo.Geometry, geo.Attributes())
	}

	return
}

func (page *Page) CreateTemplate() (temp *templates.Template) {
	temp = templates.NewTemplate(page.Title, 0, page.Layer, 0, len(page.PropMap))

	for geoID, prop := range page.PropMap {
		prop.CreateGeometry(temp, geoID)
	}

	return temp
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

func (page *Page) UnmarshalJSON(b []byte) error {
	var tempPage struct {
		Page
		UnmarshalJSON struct{}
	}

	err := json.Unmarshal(b, &tempPage)
	if err != nil {
		return err
	}

	*page = tempPage.Page
	return nil
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

func (page *Page) GetTemplateID() int {
	return page.TemplateID
}

func (page *Page) GetLayer() int {
	return page.Layer
}

func (page *Page) GetPropMap() map[int]*props.Property {
	return page.PropMap
}
