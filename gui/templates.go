package gui

import (
	"chroma-viz/props"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
    PREVIEW = iota
    TOP_LEFT
    LOWER_FRAME
    TICKER
)

type Template struct {
    Box         *gtk.ListBoxRow
    title       string
    templateID  int
    numProps    int
    layer       int
    propType    []int
    propName    []string
}

func NewTemplate(title string, id int, layer int, n int) *Template {
    temp := &Template{title: title, templateID: id, layer: layer}

    temp.propType = make([]int, n)
    temp.propName = make([]string, n)
    return temp
}

func (temp *Template) templateToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting template to row (%s)", err)
    }

    row1.Add(textToBuffer(temp.title))
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

type TempTree struct {
    *gtk.TreeView
    treeList  *gtk.ListStore
    temps     map[int]*Template
    show      *ShowTree
}

func NewTempList(show *ShowTree) *TempTree {
    var err error
    temp := &TempTree{}

    temp.TreeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.temps = make(map[int]*Template)
    temp.show = show

    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.AppendColumn(column)
    column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.AppendColumn(column)

    temp.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.SetModel(temp.treeList)

    temp.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, err := temp.treeList.GetIter(path)
            if err != nil {
                log.Fatalf("Error sending template to show (%s)", err)
            }

            id, err := temp.treeList.GetValue(iter, 1)
            if err != nil {
            }

            val, err := id.GoValue()
            if err != nil {
                log.Fatalf("Error sending template to show (%s)", err)
            }

            tempID, err := strconv.Atoi(val.(string))
            if err != nil {
                log.Fatalf("Error sending template to show (%s)", err)
            }


            temp.show.NewShowPage(temp.temps[tempID]) 
        })

    //temp.exampleHub()

    return temp
}

func (temp *TempTree) AddTemplate(title string, id, layer, n int) *Template {
    temp.temps[id] = NewTemplate(title, id, layer, n)

    temp.treeList.Set(
        temp.treeList.Append(), 
        []int{0, 1}, 
        []interface{}{title, id})

    return temp.temps[id]
}

func (temp *TempTree) exampleHub() {
    var page *Template

    page = temp.AddTemplate("Red Box", 1, TOP_LEFT, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Title", props.TEXT_PROP)
    page.AddProp("Subtitle", props.TEXT_PROP)

    page = temp.AddTemplate("Orange Box", 2, TOP_LEFT, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Title", props.TEXT_PROP)
    page.AddProp("Subtitle", props.TEXT_PROP)

    page = temp.AddTemplate("Blue Box", 3, LOWER_FRAME, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Logo", props.CIRCLE_PROP)
    page.AddProp("Title", props.TEXT_PROP)
    page.AddProp("Subtitle", props.TEXT_PROP)

    page = temp.AddTemplate("Clock Box", 4, TOP_LEFT, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Clock", props.CLOCK_PROP)

    page = temp.AddTemplate("White Circle", 5, LOWER_FRAME, 10)
    page.AddProp("Circle", props.CLOCK_PROP)

    page = temp.AddTemplate("Graph", 6, LOWER_FRAME, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Graph", props.GRAPH_PROP)
    page.AddProp("Title", props.TEXT_PROP)

    page = temp.AddTemplate("Ticker", 7, TICKER, 10)
    page.AddProp("Background", props.RECT_PROP)
    page.AddProp("Text", props.TICKER_PROP)
}

