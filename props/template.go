package props

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Template struct {
    Box         *gtk.ListBoxRow
    Title       string
    templateID  int
    numProps    int
    layer       int
    propType    []int
    propName    []string
}

func NewTemplate(title string, id int, layer int, n int) *Template {
    temp := &Template{Title: title, templateID: id, layer: layer}

    temp.propType = make([]int, n)
    temp.propName = make([]string, n)
    return temp
}

func (temp *Template) templateToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting template to row (%s)", err)
    }

    row1.Add(textToBuffer(temp.Title))
    return row1
}

func (temp *Template) AddProp(name string, typed int) {
    if temp.numProps == len(temp.propName) {
        log.Println("Ran out of memory in template")
        return
    }

    temp.propName[temp.numProps] = name
    temp.propType[temp.numProps] = typed
    temp.numProps++
}

func textToBuffer(text string) *gtk.TextView {
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


