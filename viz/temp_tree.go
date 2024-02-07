package viz

import (
	"chroma-viz/props"
	"chroma-viz/templates"
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

type TempTree struct {
    *gtk.TreeView
    treeList        *gtk.ListStore
    Temps           map[int]*templates.Template
    sendTemplate    func(*templates.Template)
}

func NewTempTree(templateToShow func(*templates.Template)) *TempTree {
    var err error
    temp := &TempTree{sendTemplate: templateToShow}

    temp.TreeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.Temps = make(map[int]*templates.Template)

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
                log.Fatalf("Error sending template to show (%s)", err)
            }

            val, err := id.GoValue()
            if err != nil {
                log.Fatalf("Error sending template to show (%s)", err)
            }

            tempID, err := strconv.Atoi(val.(string))
            if err != nil {
                log.Fatalf("Error sending template to show (%s)", err)
            }

            temp.sendTemplate(temp.Temps[tempID])
        })

    return temp
}

func (temp *TempTree) AddTemplate(title string, id, layer, num_geo int) (*templates.Template, error) {
    temp.Temps[id] = templates.NewTemplate(title, id, layer, num_geo)

    err := temp.treeList.Set(
        temp.treeList.Append(), 
        []int{0, 1}, 
        []interface{}{title, id})

    return temp.Temps[id], err
}

func (temp *TempTree) exampleHub() {
    var page *templates.Template

    page, _ = temp.AddTemplate("Red Box", 1, TOP_LEFT, 10)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Title", 1, props.TEXT_PROP)
    page.AddProp("Subtitle", 2, props.TEXT_PROP)

    page, _ = temp.AddTemplate("Orange Box", 2, TOP_LEFT, 10)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Title", 1, props.TEXT_PROP)
    page.AddProp("Subtitle", 2, props.TEXT_PROP)

    page, _ = temp.AddTemplate("Blue Box", 3, LOWER_FRAME, 10)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Logo", 1, props.CIRCLE_PROP)
    page.AddProp("Title", 2, props.TEXT_PROP)
    page.AddProp("Subtitle", 3, props.TEXT_PROP)

    page, _ = temp.AddTemplate("Clock Box", 4, TOP_LEFT, 20)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Circle", 1, props.CIRCLE_PROP)
    page.AddProp("Left Split", 2, props.RECT_PROP)
    page.AddProp("Team 1", 3, props.TEXT_PROP)
    page.AddProp("Score 1", 4, props.TEXT_PROP)
    page.AddProp("Mid Split", 5, props.RECT_PROP)
    page.AddProp("Team 2", 6, props.TEXT_PROP)
    page.AddProp("Score 2", 7, props.TEXT_PROP)
    page.AddProp("Right Split", 8, props.RECT_PROP)
    page.AddProp("Clock", 9, props.CLOCK_PROP)
    page.AddProp("Clock Title", 10, props.TEXT_PROP)

    page, _ = temp.AddTemplate("White Circle", 5, LOWER_FRAME, 10)
    page.AddProp("Circle", 0, props.CIRCLE_PROP)

    page, _ = temp.AddTemplate("Graph", 6, LOWER_FRAME, 10)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Graph", 0, props.GRAPH_PROP)
    page.AddProp("Title", 1, props.TEXT_PROP)

    page, _ = temp.AddTemplate("Ticker", 7, TICKER, 10)
    page.AddProp("Background", 0, props.RECT_PROP)
    page.AddProp("Box", 1, props.RECT_PROP)
    page.AddProp("Text", 2, props.TICKER_PROP)
}

