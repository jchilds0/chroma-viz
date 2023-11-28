package gui

import "github.com/gotk3/gotk3/gtk"

type Page struct {
    Box         *gtk.ListBoxRow
    pageNum     int
    title       string
    templateID  int
}

func NewPage(pageNum int, title string, id int) *Page {
    return &Page{pageNum: pageNum, title: title, templateID: id}
}

func pageToListRow(page Page) *gtk.ListBoxRow {
    row1, _ := gtk.ListBoxRowNew()

    row1.Add(textToBuffer(string(page.pageNum)))
    row1.Add(textToBuffer(page.title))
    return row1
}

func textToBuffer(text string) *gtk.TextView {
    text1, _ := gtk.TextViewNew()
    buffer, _ := text1.GetBuffer()
    buffer.SetText(text)

    return text1
}


