package gui

import (
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Template struct {
    Box         *gtk.ListBoxRow
    title       string
    templateID  int
    props       map[string]string
}

func NewTemplate(title string, id int) *Template {
    temp := &Template{title: title, templateID: id}
    temp.props = make(map[string]string)

    return temp
}

func (temp *Template) templateToListRow() *gtk.ListBoxRow {
    row1, _ := gtk.ListBoxRowNew()
    row1.Add(textToBuffer(temp.title))

    return row1
}

func (temp *Template) AddProp(name string, typed string) {
    temp.props[name] = typed
}

func textToBuffer(text string) *gtk.TextView {
    text1, _ := gtk.TextViewNew()
    buffer, _ := text1.GetBuffer()
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
    temp := &TempTree{}
    temp.TreeView, _ = gtk.TreeViewNew()
    temp.temps = make(map[int]*Template)
    temp.show = show

    cell, _ := gtk.CellRendererTextNew()
    column, _ := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
    temp.AppendColumn(column)
    column, _ = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
    temp.AppendColumn(column)

    temp.treeList, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
    temp.SetModel(temp.treeList)

    temp.AddTemplate("Red Box", 1)
    temp.temps[1].AddProp("Background", "RectProp")
    temp.temps[1].AddProp("2", "TextProp")

    temp.AddTemplate("Orange Box", 2)
    temp.temps[2].AddProp("Background", "RectProp")
    temp.temps[2].AddProp("2", "TextProp")

    temp.AddTemplate("Blue Box", 3)
    temp.temps[3].AddProp("Background", "RectProp")
    temp.temps[3].AddProp("2", "TextProp")

    temp.AddTemplate("Clock Box", 4)
    temp.temps[4].AddProp("Background", "RectProp")
    temp.temps[4].AddProp("Clock", "ClockProp")

    temp.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, _ := temp.treeList.GetIter(path)
            id, _ := temp.treeList.GetValue(iter, 1)
            val, _ := id.GoValue()
            tempID, _ := strconv.Atoi(val.(string))

            temp.show.NewShowPage(temp.temps[tempID]) 
        })

    return temp
}

func (temp *TempTree) AddTemplate(title string, id int) {
    temp.temps[id] = NewTemplate(title, id)

    temp.treeList.Set(
        temp.treeList.Append(), 
        []int{0, 1}, 
        []interface{}{title, id})
}
