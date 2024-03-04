package shows

import (
	"chroma-viz/props"
	"chroma-viz/templates"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

type Page struct {
    Box         *gtk.ListBoxRow
    PageNum     int
    Title       string
    TemplateID  int
    Layer       int
    PropMap     map[int]props.Property
}

func newPage(pageNum int, title string, temp *templates.Template) *Page {
    page := &Page{
        PageNum: pageNum, 
        Title: title, 
        TemplateID: temp.TempID,
        Layer: temp.Layer,
    }
    page.PropMap = make(map[int]props.Property, temp.NumProps)

    for i, prop := range temp.Prop {
        page.PropMap[i] = props.NewProperty(prop.Type, prop.Name, prop.Visible)
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


