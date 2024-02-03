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
    PropType    []int
    PropName    []string
}

func NewTemplate(title string, id int, layer int, n int) *Template {
    temp := &Template{Title: title, TempID: id, Layer: layer}

    temp.PropType = make([]int, n)
    temp.PropName = make([]string, n)
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

func (temp *Template) AddProp(name string, typed int) {
    if temp.NumProps == len(temp.PropName) {
        log.Println("Ran out of memory in template")
        return
    }

    temp.PropName[temp.NumProps] = name
    temp.PropType[temp.NumProps] = typed
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

