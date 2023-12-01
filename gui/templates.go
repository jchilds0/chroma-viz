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
}

func NewTemplate(title string, id int) *Template {
    return &Template{title: title, templateID: id}
}

func (temp *Template) templateToListRow() *gtk.ListBoxRow {
    row1, _ := gtk.ListBoxRowNew()
    row1.Add(textToBuffer(temp.title))

    return row1
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

    temp.AppendColumn(createColumn("Name", 0))
    temp.AppendColumn(createColumn("Template ID", 1))

    temp.treeList, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
    temp.SetModel(temp.treeList)

    temp.AddTemplate("template 1", 1)
    temp.AddTemplate("template 2", 2)
    temp.AddTemplate("template 3", 3)

    temp.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, _ := temp.treeList.GetIter(path)
            id, _ := temp.treeList.GetValue(iter, 1)
            val, _ := id.GoValue()
            tempID, _ := strconv.Atoi(val.(string))

            temp.show.NewShowPage(*temp.temps[tempID]) 
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

