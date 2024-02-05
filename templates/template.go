package templates 

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Template struct {
    Box         *gtk.ListBoxRow
    Title       string
    TempID  int
    NumProps    int
    Layer       int
    PropType    map[int]int
    PropName    map[int]string
}

func NewTemplate(title string, id int, layer int, num_geo int) *Template {
    temp := &Template{Title: title, TempID: id, Layer: layer}

    temp.PropType = make(map[int]int, num_geo)
    temp.PropName = make(map[int]string, num_geo)
    return temp
}

func (temp *Template) TemplateToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting template to row (%s)", err)
    }

    row1.Add(TextToBuffer(temp.Title))
    return row1
}

func (temp *Template) AddProp(name string, geo_id, typed int) {
    temp.PropName[geo_id] = name
    temp.PropType[geo_id] = typed
    temp.NumProps++
}

func TextToBuffer(text string) *gtk.TextView {
    text1, err := gtk.TextViewNew()
    if err != nil {
        log.Fatalf("Error creating text buffer (%s)", err)
    }

    buffer, err := text1.GetBuffer()
    if err != nil {
        log.Fatalf("Error creating text buffer (%s)", err)
    }

    buffer.SetText(text)
    return text1
}

