package gui

import (
	"log"
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
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting template to row (%s)", err)
    }

    row1.Add(textToBuffer(temp.title))
    return row1
}

func (temp *Template) AddProp(name string, typed string) {
    temp.props[name] = typed
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

    temp.AddTemplate("Red Box", 1)
    temp.temps[1].AddProp("Background", "RectProp")
    temp.temps[1].AddProp("Title", "TextProp")
    temp.temps[1].AddProp("Subtitle", "TextProp")

    temp.AddTemplate("Orange Box", 2)
    temp.temps[2].AddProp("Background", "RectProp")
    temp.temps[2].AddProp("Title", "TextProp")
    temp.temps[2].AddProp("Subtitle", "TextProp")

    temp.AddTemplate("Blue Box", 3)
    temp.temps[3].AddProp("Background", "RectProp")
    temp.temps[3].AddProp("Title", "TextProp")
    temp.temps[3].AddProp("Subtitle", "TextProp")

    temp.AddTemplate("Clock Box", 4)
    temp.temps[4].AddProp("Background", "RectProp")
    temp.temps[4].AddProp("Clock", "ClockProp")

    temp.AddTemplate("White Circle", 5)
    temp.temps[5].AddProp("Circle", "CircleProp")

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

    return temp
}

func (temp *TempTree) AddTemplate(title string, id int) {
    temp.temps[id] = NewTemplate(title, id)

    temp.treeList.Set(
        temp.treeList.Append(), 
        []int{0, 1}, 
        []interface{}{title, id})
}

