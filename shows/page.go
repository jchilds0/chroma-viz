package shows

import (
	"chroma-viz/props"
	"chroma-viz/templates"
	"encoding/json"
	"log"
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
    Box         *gtk.ListBoxRow
    PageNum     int
    Title       string
    TemplateID  int
    Layer       int
    PropMap     map[int]*props.Property
}

func newPage(pageNum int, title string, temp *templates.Template, cont func(*Page)) *Page {
    page := &Page{
        PageNum: pageNum, 
        Title: title, 
        TemplateID: temp.TempID,
        Layer: temp.Layer,
    }
    page.PropMap = make(map[int]*props.Property, temp.NumProps)

    contPage := func() { cont(page) }

    for i, prop := range temp.Prop {
        page.PropMap[i] = props.NewProperty(prop.Type, prop.Name, prop.Visible, contPage)
    }

    return page
}

func (page *Page) PageToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting page to list (%s)", err)
    }

    row1.Add(templates.TextToBuffer(strconv.Itoa(page.PageNum)))
    row1.Add(templates.TextToBuffer(page.Title))

    return row1
}

func (page *Page) UnmarshalJSON(b []byte) error {
    var tempPage struct {
        Page
        UnmarshalJSON struct {}
    }

    err := json.Unmarshal(b, &tempPage)
    if err != nil {
        return err 
    }

    page = &tempPage.Page

    return nil
}

func ImportPage(filename string) (page *Page, err error) {
    buf, err := os.ReadFile(filename)
    if err != nil {
        return 
    }

    page = &Page{}
    err = json.Unmarshal(buf, page)
    if err != nil {
        return
    }

    return
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
