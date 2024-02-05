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

func NewPage(pageNum int, title string, temp *templates.Template) *Page {
    page := &Page{
        PageNum: pageNum, 
        Title: title, 
        TemplateID: temp.TempID,
        Layer: temp.Layer,
    }
    page.PropMap = make(map[int]props.Property, temp.NumProps)

    for i, prop := range temp.PropType {
        name := temp.PropName[i]

        switch (prop) {
        case props.RECT_PROP:
            page.PropMap[i] = props.NewRectProp(name)

        case props.TEXT_PROP:
            page.PropMap[i] = props.NewTextProp(name)

        case props.CIRCLE_PROP:
            page.PropMap[i] = props.NewCircleProp(name)

        case props.CLOCK_PROP:
            page.PropMap[i] = props.NewClockProp(name)

        case props.GRAPH_PROP:
            page.PropMap[i] = props.NewGraphProp(name)

        case props.TICKER_PROP:
            page.PropMap[i] = props.NewTickerProp(name)

        default:
            log.Printf("Page %d: Unknown property %d", pageNum, prop)
        }
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


