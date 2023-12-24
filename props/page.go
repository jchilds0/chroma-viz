package props

import (
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
    PropMap     []Property
}

func NewPage(pageNum int, title string, temp *Template) *Page {
    page := &Page{
        PageNum: pageNum, 
        Title: title, 
        TemplateID: temp.templateID,
        Layer: temp.layer,
    }
    page.PropMap = make([]Property, temp.numProps)

    for i := 0; i < temp.numProps; i++ {
        prop := temp.propType[i]
        name := temp.propName[i]

        switch (prop) {
        case RECT_PROP:
            page.PropMap[i] = NewRectProp(name)

        case TEXT_PROP:
            page.PropMap[i] = NewTextProp(name)

        case CIRCLE_PROP:
            page.PropMap[i] = NewCircleProp(name)

        case CLOCK_PROP:
            page.PropMap[i] = NewClockProp(name)

        case GRAPH_PROP:
            page.PropMap[i] = NewGraphProp(name)

        case TICKER_PROP:
            page.PropMap[i] = NewTickerProp(name)

        default:
            log.Printf("Page %d: Unknown property %d", pageNum, prop)
        }
    }

    return page
}

func (page *Page) pageToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting page to list (%s)", err)
    }

    row1.Add(textToBuffer(strconv.Itoa(page.PageNum)))
    row1.Add(textToBuffer(page.Title))

    return row1
}


