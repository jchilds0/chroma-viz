package viz

import (
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
    treeView        *gtk.TreeView
    treeList        *gtk.ListStore
    //Temps           map[int]*templates.Template
    Temps           *templates.Temps
    sendTemplate    func(*templates.Template)
}

func NewTempTree(templateToShow func(*templates.Template)) *TempTree {
    var err error
    temp := &TempTree{sendTemplate: templateToShow}

    temp.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.Temps = templates.NewTemps()

    // create tree columns
    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.treeView.AppendColumn(column)
    column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.treeView.AppendColumn(column)

    temp.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
    if err != nil {
        log.Fatalf("Error creating temp list (%s)", err)
    }

    temp.treeView.SetModel(temp.treeList)

    // send template to show on double click
    temp.treeView.Connect("row-activated", 
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

            temp.sendTemplate(temp.Temps.Temps[tempID])
        })

    return temp
}

func (temp *TempTree) AddTemplate(title string, id, layer, num_geo int) (*templates.Template, error) {
    temp.Temps.SetTemplate(id, layer, num_geo, title)

    err := temp.treeList.Set(
        temp.treeList.Append(), 
        []int{0, 1}, 
        []interface{}{title, id},
    )

    return temp.Temps.Temps[id], err
}

